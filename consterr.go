package gotvm

// ConstError is a constant error type.
type ConstError string

// Error returns the description of the error.
func (e ConstError) Error() string {
	return string(e)
}
