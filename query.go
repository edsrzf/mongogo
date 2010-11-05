// Copyright 2010, Evan Shaw.  All rights reserved.
// Use of this source code is governed by a BSD-style License
// that can be found in the LICENSE file.

package mongo

// Query represents a MongoDB database query.
// Simple queries can be constructed just like regular maps. For example the query
//
//	q := Query{"_id": 14}
//
// would search a collection for a document with an _id of 14.
// Query's methods allow construction of more advanced queries.
type Query map[string]interface{}

// makeComplex makes query into a "complex" query, where the inner query is the
// value stored under the key "$query" and other options can be specified.
func (query *Query) makeComplex() {
	q := *query
	if _, ok := q["$query"]; !ok {
		outer := make(Query)
		outer["$query"] = q
		*query = outer
	}
}

// Explain causes MongoDB to return information on how this query is performed
// instead of returning the actual results of the query.
// After calling this method on a Query, the Query should not be updated through
// index expressions.
func (q *Query) Explain() {
	q.makeComplex()
	(*q)["$explain"] = true
}

// Hint forces MongoDB to use a particular index for this query.
// After calling this method on a Query, the Query should not be updated through
// index expressions.
func (q *Query) Hint(index string) {
	q.makeComplex()
	(*q)["$hint"] = index
}

func (q *Query) MinKey(min map[string]int32) {
	q.makeComplex()
	(*q)["$min"] = min
}

func (q *Query) MaxKey(max map[string]int32) {
	q.makeComplex()
	(*q)["$max"] = max
}

func (q *Query) MaxScan(max int32) {
	q.makeComplex()
	(*q)["$maxScan"] = max
}

// ShowDiskLocation will cause the returned documents to contain a key $diskLoc
// which shows the location of that document on disk.
// After calling this method on a Query, the Query should not be updated through
// index expressions.
func (q *Query) ShowDiskLocation() {
	q.makeComplex()
	(*q)["$showDiskLoc"] = true
}

// After calling this method on a Query, the Query should not be updated through
// index expressions.
func (q *Query) Snapshot() {
	q.makeComplex()
	(*q)["$snapshot"] = true
}

func (q *Query) Sort(keys map[string]int32) {
	q.makeComplex()
	(*q)["$orderby"] = keys
}
