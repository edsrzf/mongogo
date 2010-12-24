// This test relies upon there being a mongo server running on localhost on the
// default port (27017).
package mongo

import (
	"fmt"
	"github.com/edsrzf/go-bson"
	"runtime"
	"testing"
)

const (
	mongoAddr = "localhost:27017"
	testDb = "mongogoTestDB"
	testCol = "mongogoTestCol"
)

func info() string {
	_, file, line, ok := runtime.Caller(1)
	if ok != true {
		return "[err]"
	}
	return fmt.Sprintf("[%v][%v]", file, line)
}

func TestConnect(t *testing.T) {
	conn, err := Dial(mongoAddr)
	if err != nil {
		t.Fatal(info())
	}
	err = conn.Close()
	if err != nil {
		t.Fatal(info())
	}
}

func TestFindOneFields(t *testing.T) {
	conn, err := Dial(mongoAddr)
	if err != nil {
		t.Fatal(info())
	}
	db := conn.Database(testDb)
	col := db.Collection(testCol)

	//insert test docs
	err = col.Insert(bson.Doc{"name": "Joe", "age": 28})
	if err != nil {
		t.Fatal(info())
	}
	err = col.Insert(bson.Doc{"name": "Jim", "age": 25})
	if err != nil {
		t.Fatal(info())
	}

	//find test doc
	doc, err := col.FindOneFields(Query{"name": "Joe"},
		map[string]interface{}{"age": "1"})
	if err != nil {
		t.Fatal(info())
	}
	if doc == nil {
		t.Fatal(info())
	}
	ageIf, ok := doc["age"]
	if !ok {
		t.Fatal(info())
	}
	switch age := ageIf.(type) {
	case int64:
		if age != 28 {
			t.Fatal(info())
		}
	default:
		t.Fatal(info())
	}

	err = conn.Close()
	if err != nil {
		t.Fatal(info())
	}
}
