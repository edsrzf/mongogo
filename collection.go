// Copyright 2010, Evan Shaw.  All rights reserved.
// Use of this source code is governed by a BSD-style License
// that can be found in the LICENSE file.

package mongo

import (
	"bytes"
	"encoding/binary"
	"os"
	"github.com/edsrzf/go-bson"
)

// common message header size
// 16-byte header
const headerSize = 16

// A Collection represents a MongoDB collection.
type Collection struct {
	db       *Database
	name     string
	fullName []byte
}

// Drop deletes c from the database.
func (c *Collection) Drop() os.Error {
	cmd := Query{"drop": string(c.fullName)}
	_, err := c.db.Command(cmd)
	return err
}

// Update updates a single document selected by query, according to doc.
func (c *Collection) Update(query, doc bson.Doc) os.Error {
	return c.update(query, doc, false, false)
}

// Upsert updates or inserts a single document selected by query,
// according to doc.
func (c *Collection) Upsert(query, doc bson.Doc) os.Error {
	return c.update(query, doc, true, false)
}

// Update updates multiple documents selected by query, according to doc.
func (c *Collection) UpdateAll(query, doc bson.Doc) os.Error {
	return c.update(query, doc, false, true)
}

// UpsertAll updates or inserts multiple documents selected by query,
// according to doc.
func (c *Collection) UpsertAll(query, doc bson.Doc) os.Error {
	return c.update(query, doc, true, true)
}

func (c *Collection) update(query, doc bson.Doc, upsert, multi bool) os.Error {
	selData, err := bson.Marshal(query)
	if err != nil {
		return err
	}
	docData, err := bson.Marshal(doc)
	if err != nil {
		return err
	}
	cap := headerSize + 4 + len(c.fullName) + 4 + len(selData) + len(docData)
	payload := make([]byte, headerSize+4, cap)
	buf := bytes.NewBuffer(payload)
	buf.Write(c.fullName)
	var flags int32
	if upsert {
		flags |= 1
	}
	if multi {
		flags |= 2
	}
	binary.Write(buf, order, flags)
	buf.Write(selData)
	buf.Write(docData)
	payload = payload[:cap]
	return c.db.conn.sendMessage(2001, 0, payload)
}

// Insert creates a new document in c.
func (c *Collection) Insert(doc bson.Doc) os.Error {
	data, err := bson.Marshal(doc)
	if err != nil {
		return err
	}
	cap := headerSize + 4 + len(c.fullName) + len(data)
	payload := make([]byte, headerSize+4, cap)
	buf := bytes.NewBuffer(payload)
	buf.Write(c.fullName)
	buf.Write(data)
	payload = payload[:cap]
	return c.db.conn.sendMessage(2002, 0, payload)
}

// Find searches c for any documents matching a query. It skips the first skip
// documents and limits the search to limit.
func (c *Collection) Find(query Query, skip, limit int32) (*Cursor, os.Error) {
	return c.FindFields(query, nil, skip, limit)
}

// FindFields performs a query that returns only specified fields. It skips the
// first skip documents and limits the search to limit.
// The fields specified can be inclusive or exclusive, but not both. That is,
// the values in the fields parameter must be all true or all false with no
// mixing. Fields with true values will be returned, while fields with false
// values will be excluded.
func (c *Collection) FindFields(query Query, fields map[string]interface{}, skip, limit int32) (*Cursor, os.Error) {
	conn := c.db.conn
	data, err := bson.Marshal(bson.Doc(query))
	if err != nil {
		return nil, err
	}
	var fieldData []byte
	if fields != nil {
		fieldData, err = bson.Marshal(bson.Doc(fields))
		if err != nil {
			return nil, err
		}
	}
	cap := headerSize + 4 + len(c.fullName) + 8 + len(data) + len(fieldData)
	payload := make([]byte, headerSize, cap)
	buf := bytes.NewBuffer(payload[headerSize:])
	// TODO(eds): Consider supporting flags
	binary.Write(buf, order, int32(0))
	buf.Write(c.fullName)
	binary.Write(buf, order, skip)
	binary.Write(buf, order, limit)
	buf.Write(data)
	buf.Write(fieldData)
	payload = payload[:cap]
	if err := conn.sendMessage(2004, 0, payload); err != nil {
		return nil, err
	}

	reply, err := conn.readReply()
	if err != nil {
		return nil, err
	}

	return &Cursor{c, reply.cursorID, 0, reply.docs}, nil
}

func (c *Collection) FindOneFields(query Query, fields map[string]interface{}) (bson.Doc, os.Error) {
	cursor, err := c.FindFields(query, fields, 0, 1)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	return cursor.Next(), nil
}

// FindAll returns all documents in c matching a query.
func (c *Collection) FindAll(query Query) (*Cursor, os.Error) {
	return c.Find(query, 0, 0)
}

// FindOne returns the first document in c that matches a query.
func (c *Collection) FindOne(query Query) (bson.Doc, os.Error) {
	cursor, err := c.Find(query, 0, 1)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	return cursor.Next(), nil
}

// Count returns the number of documents in c that match a query.
func (c *Collection) Count(query bson.Doc) (int64, os.Error) {
	cmd := Query{"count": c.name, "query": query}
	reply, err := c.db.Command(cmd)
	if reply == nil || err != nil {
		return -1, err
	}

	// NOTE(eds): Mongo returns count as a double? Really? That seems silly.
	return int64(reply["n"].(float64)), nil
}

func (c *Collection) remove(query bson.Doc, singleRemove bool) os.Error {
	data, err := bson.Marshal(query)
	if err != nil {
		return err
	}
	l := len(c.fullName)
	payload := make([]byte, headerSize+4+l+4+len(data))
	copy(payload[headerSize+4:], c.fullName)
	if singleRemove {
		payload[headerSize+4+l] |= 1
	}
	copy(payload[headerSize+4+l+4:], data)
	return c.db.conn.sendMessage(2006, 0, payload)
}

// Remove removes all documents in c that match a query.
func (c *Collection) Remove(query bson.Doc) os.Error {
	return c.remove(query, false)
}

// RemoveFirst removes the first document in c that matches a query.
func (c *Collection) RemoveFirst(query bson.Doc) os.Error {
	return c.remove(query, true)
}

// EnsureIndex ensures that an index exists on this collection.
func (c *Collection) EnsureIndex(name string, keys map[string]int32, unique bool) os.Error {
	col := c.db.Collection("system.indexes")
	id := bson.Doc{"name": name, "ns": string(c.fullName), "key": keys, "unique": unique}
	return col.Insert(id)
}

// DropIndexes deletes all indexes on c.
func (c *Collection) DropIndexes() os.Error {
	return c.DropIndex("*")
}

// DropIndex deletes a single index.
func (c *Collection) DropIndex(name string) os.Error {
	cmd := Query{"deleteIndexes": string(c.fullName), "index": name}
	_, err := c.db.Command(cmd)
	return err
}
