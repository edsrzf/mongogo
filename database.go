// Copyright 2010, Evan Shaw.  All rights reserved.
// Use of this source code is governed by a BSD-style License
// that can be found in the LICENSE file.

package mongo

import (
	"os"
	"github.com/edsrzf/go-bson"
)

// Database represents a MongoDB database.
type Database struct {
	conn *Conn
	name string
}

// Collection returns a Collection specified by name.
func (db *Database) Collection(name string) *Collection {
	return &Collection{db, name, []byte(db.name + "." + name + "\x00")}
}

// Drop deletes db.
func (db *Database) Drop() os.Error {
	cmd := Query{"dropDatabase": int32(1)}
	_, err := db.Command(cmd)
	return err
}

// Eval evaluates a JavaScript expression or function on the MongoDB server.
func (db *Database) Eval(code *bson.JavaScript, args string) (bson.Doc, os.Error) {
	cmd := Query{"$eval": code}
	if args != "" {
		cmd["args"] = args
	}
	return db.Command(cmd)
}

// Stats returns database statistics for db.
func (db *Database) Stats() (bson.Doc, os.Error) {
	cmd := Query{"dbstats": int32(1)}
	return db.Command(cmd)
}

// Repair checks for and repairs corruption in the database.
func (db *Database) Repair() os.Error {
	cmd := Query{"repairDatabase": int32(1)}
	_, err := db.Command(cmd)
	return err
}

// Command sends an arbitrary command to the database.
// It is equivalent to
//	col := db.Collection("$cmd")
//	doc := col.FindOne(cmd)
// If the $err key is not null in the reply, Command returns an error.
func (db *Database) Command(cmd Query) (bson.Doc, os.Error) {
	col := db.Collection("$cmd")
	reply, err := col.FindOne(cmd)
	if reply["$err"] != nil {
		msg, ok := reply["$err"].(string)
		if !ok {
			// this probably shouldn't ever happen
			msg = "non-string error message"
		}
		return nil, os.NewError(msg)
	}
	return reply, err
}
