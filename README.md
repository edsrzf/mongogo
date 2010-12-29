Mongogo
=======

Mongogo is a MongoDB driver for the [Go programming language](http://golang.org/).

This project is still in development. It's been tested on Arch and Ubuntu Linux for
the amd64 architecture, but there's no reason it shouldn't work on other architectures
as well.

Dependencies
------------

Mongogo compiles with Go release 2010-10-27 or newer, barring any recent language or
library changes.

Mongogo works with MongoDB version 2.6 or newer. It may partially work with older versions.

Mongogo's only non-core Go dependency is [Go-BSON](go-bson).
You can install it with goinstall by running
    goinstall github.com/edsrzf/go-bson

Usage
-----

Create a connection:

    conn := mongo.Dial("127.0.0.1:27017")

Get a database:

    db := conn.Database("blog")

Get a collection:

    col := db.Collection("posts")

Insert a document into the collection:

    doc := map[string]interface{}{"title": "Hello", "body": "World!"}
    col.Insert(doc)

Query the database:

    q := mongo.Query{"title": "Hello"}
    cursor := col.Find(q, 0, 0)
    defer cursor.Close()

See the documentation in the source for more information.

Contributing
------------

Simply use GitHub as usual to create a fork, make your changes, and create a pull
request. Code is expected to be formatted with gofmt and to adhere to the usual Go
conventions -- that is, the conventions used by Go's core libraries.
