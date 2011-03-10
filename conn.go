// Copyright 2010, Evan Shaw.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The mongo package provides a MongoDB driver implementation.
package mongo

import (
	"bytes"
	"encoding/binary"
	"os"
	"net"
	"rand"
	"reflect"
	"github.com/edsrzf/go-bson"
)

var order = binary.LittleEndian

// ConnError is an error related to the connection to the MongoDB server.
type ConnError string

func (e ConnError) String() string { return string(e) }

// A Conn represents a connection to a MongoDB server.
type Conn struct {
	conn net.Conn
}

type reply struct {
	requestID      int32
	responseTo     int32
	responseFlags  int32
	cursorID       int64
	startingFrom   int32
	numberReturned int32
	docs           []interface{}
}

// Dial connects to a MongoDB server at the remote address addr.
func Dial(addr string) (*Conn, os.Error) {
	c, err := net.Dial("tcp", "", addr)
	if err != nil {
		return nil, ConnError(err.String())
	}
	return &Conn{c}, nil
}

// Close closes the connection.
func (c *Conn) Close() os.Error {
	err := c.conn.Close()
	if err != nil {
		return ConnError(err.String())
	}
	return nil
}

// Database returns the Database object for a name.
func (c *Conn) Database(name string) *Database {
	return &Database{c, name}
}

func (c *Conn) sendMessage(opCode, responseId int32, message []byte) os.Error {
	messageLength := int32(len(message))
	message = message[:0]
	buf := bytes.NewBuffer(message)
	binary.Write(buf, order, messageLength)
	// request ID
	binary.Write(buf, order, rand.Int31())
	// response to
	binary.Write(buf, order, responseId)
	binary.Write(buf, order, opCode)
	message = message[:messageLength]
	_, err := c.conn.Write(message)
	if err != nil {
		return ConnError(err.String())
	}
	return nil
}

func (c *Conn) readReply(form *reflect.PtrType) (*reply, os.Error) {
	var size uint32
	err := binary.Read(c.conn, order, &size)
	if err != nil {
		return nil, ConnError(err.String())
	}
	raw := make([]byte, size)
	_, err = c.conn.Read(raw)
	if err != nil {
		return nil, ConnError(err.String())
	}
	buf := bytes.NewBuffer(raw)
	r := new(reply)
	binary.Read(buf, order, &r.requestID)
	binary.Read(buf, order, &r.responseTo)
	var opCode int32
	binary.Read(buf, order, &opCode)
	if opCode != 1 {
		return nil, os.NewError("expected OP_REPLY opCode")
	}
	binary.Read(buf, order, &r.responseFlags)
	binary.Read(buf, order, &r.cursorID)
	binary.Read(buf, order, &r.startingFrom)
	binary.Read(buf, order, &r.numberReturned)
	elType := form.Elem()
	r.docs = make([]interface{}, r.numberReturned)
	for i := range r.docs {
		raw := buf.Bytes()
		size := order.Uint32(raw)
		doc := reflect.MakeZero(elType)
		ptr := reflect.MakeZero(form)
		ptr.(*reflect.PtrValue).PointTo(doc)
		r.docs[i] = ptr.Interface()
		err = bson.Unmarshal(raw, r.docs[i])
		if err != nil {
			break
		}
		buf.Next(int(size))
	}
	return r, err
}
