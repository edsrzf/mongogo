Mongogo
=======

Mongogo is a MongoDB driver for the [Go programming language](http://golang.org/).

This project is still in development. It's not well tested, but the basics seem to
work well enough.

Dependencies
------------

Mongogo's sole external dependency is [Go-BSON](edsrzf/go-bson).

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
    cursor := col.Query(q, 0, 0)
    defer cursor.Close()

See the documentation in the source for more information.

Contributing
------------

Simply use GitHub as usual to create a fork, make your changes, and create a pull
request. Code is expected to be formatted with gofmt and to adhere to the usual Go
conventions -- that is, the conventions used by Go's core libraries.
