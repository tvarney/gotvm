package cerr

// ConstError is a constant error type.
type Error string

// Error returns the description of the error.
func (e Error) Error() string {
	return string(e)
}
