package mongo

// Conforms to os.Error interface.
// This error is returned when the error is connection related.
//
// Example usage:
// err := mongo.Insert(... )
// if err, ok := err.(*mongo.ConnError); ok {
// 	...
// }
type ConnError struct {
	str string
}

func NewConnError(s string) *ConnError {
	return &ConnError{str: s}
}

func (this *ConnError) String() string {
	return this.str
}
