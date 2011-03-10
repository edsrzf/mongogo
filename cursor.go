// Copyright 2010, Evan Shaw.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mongo

import (
	"bytes"
	"encoding/binary"
	"os"
)

// A Cursor is the result of a query.
type Cursor struct {
	collection *Collection
	id         int64
	pos        int
	docs       []interface{}
}

// Peek returns the next document or nil if the Cursor is at the end.
func (c *Cursor) Peek() interface{} {
	if !c.More() {
		return nil
	}
	return c.docs[c.pos]
}

// Next is like Peek, but also iterates to the next document.
func (c *Cursor) Next() interface{} {
	doc := c.Peek()
	if doc != nil {
		c.pos++
	}
	return doc
}

// More indicates whether the Cursor still has more documents to iterate through.
func (c *Cursor) More() bool {
	if c.pos < len(c.docs) {
		return true
	}

	if err := c.getMore(0); err != nil {
		return false
	}

	return c.pos < len(c.docs)
}

func (c *Cursor) getMore(limit int32) os.Error {
	if c.id == 0 {
		return os.NewError("no cursorID")
	}

	cap := headerSize + 4 + len(c.collection.fullName) + 4 + 8
	payload := make([]byte, headerSize+4, cap)
	buf := bytes.NewBuffer(payload)
	buf.Write(c.collection.fullName)
	binary.Write(buf, order, limit)
	binary.Write(buf, order, c.id)
	payload = payload[:cap]

	conn := c.collection.db.conn
	if err := conn.sendMessage(2005, 0, payload); err != nil {
		return err
	}

	reply, err := conn.readReply(c.collection.form)
	if err != nil {
		return err
	}

	c.pos = 0
	c.docs = reply.docs

	return nil
}

// Close tells the server that c is no longer in use and makes c invalid.
func (c *Cursor) Close() os.Error {
	if c.id == 0 {
		// not open on server
		return nil
	}

	cap := headerSize + 16
	payload := make([]byte, headerSize+4, cap)
	buf := bytes.NewBuffer(payload)
	binary.Write(buf, order, int32(1))
	binary.Write(buf, order, c.id)
	payload = payload[:cap]
	c.id = 0
	return c.collection.db.conn.sendMessage(2007, 0, payload)
}
