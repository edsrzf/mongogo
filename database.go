// Copyright 2010, Evan Shaw.  All rights reserved.
// Use of this source code is governed by a BSD-style License
// that can be found in the LICENSE file.

package mongo

import (
	"os"
	"reflect"
	"github.com/edsrzf/go-bson"
)

// Database represents a MongoDB database.
type Database struct {
	conn *Conn
	name string
}

// Collection returns a Collection specified by name.
// The form parameter specifies the type to use for query results. It should be
// a poitner to a map with a string key type or a struct. If form is nil, the type
// map[string]interface{} will be used.
func (db *Database) Collection(name string, form interface{}) (*Collection, os.Error) {
	if form == nil {
		form = new(map[string]interface{})
	}
	ptrType, ok := reflect.Typeof(form).(*reflect.PtrType)
	if !ok {
		return nil, os.NewError("form must be a pointer type")
	}
	return &Collection{db, name, ptrType, []byte(db.name + "." + name + "\x00")}, nil
}

// Drop deletes db.
func (db *Database) Drop() os.Error {
	cmd := Query{"dropDatabase": 1}
	_, err := db.Command(cmd, struct{}{})
	return err
}

// Eval evaluates a JavaScript expression or function on the MongoDB server.
func (db *Database) Eval(code *bson.JavaScript, args string) (interface{}, os.Error) {
	cmd := Query{"$eval": code}
	if args != "" {
		cmd["args"] = args
	}
	reply, err := db.Command(cmd, nil)
	return reply.(map[string]interface{})["retval"], err
}

// A Stats structure provides statistical information about a database.
type Stats struct {
	Collections int "collections"
	Objects int "objects"
	AvgObjSize float64 "avgObjSize"
	DataSize int "dataSize"
	StorageSize int "storageSize"
	NumExtents int "numExtents"
	Indexes int "indexes"
	IndexSize int "indexSize"
	FileSize int64 "fileSize"
}

// Stats returns database statistics for db.
func (db *Database) Stats() (*Stats, os.Error) {
	cmd := Query{"dbstats": 1}
	s := new(Stats)
	iface, err := db.Command(cmd, s)
	return iface.(*Stats), err
}

// Repair checks for and repairs corruption in the database.
func (db *Database) Repair() os.Error {
	cmd := Query{"repairDatabase": 1}
	_, err := db.Command(cmd, struct{}{})
	return err
}

// Command sends an arbitrary command to the database.
// It is equivalent to
//	col, err := db.Collection("$cmd", form)
//	result, err = col.FindOne(cmd)
// If the $err key is not null in the reply, Command returns an error.
func (db *Database) Command(cmd Query, form interface{}) (result interface{}, err os.Error) {
	col, err := db.Collection("$cmd", form)
	if err != nil {
		return nil, err
	}
	result, err = col.FindOne(cmd)
	// TODO: do a better job of handling errors
	reply, ok := result.(map[string]interface{})
	if ok && reply["$err"] != nil {
		msg, ok := reply["$err"].(string)
		if !ok {
			// this probably shouldn't ever happen
			msg = "non-string error message"
		}
		err = os.NewError(msg)
	}
	return
}
