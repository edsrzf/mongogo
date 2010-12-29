package mongo

import (
	"os"
)

// Conforms to os.Error interface.
// This error is returned when the error is connection related.
//
// Example usage:
// err := mongo.Insert(... )
// if err, ok := err.(mongo.connError); ok {
// 	...
// }
type connError struct {
	str string
}

func NewConnError(err os.Error) *connError {
	return &connError{str: err.String()}
}

func (this *connError) String() string {
	return this.str
}
