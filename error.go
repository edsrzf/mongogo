package mongo

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

func NewConnError(s string) *connError {
	return &connError{str: s}
}

func (this *connError) String() string {
	return this.str
}
