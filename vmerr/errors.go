package vmerr

import "strconv"

const (
	ErrTooFewValues     ConstError = "too few arguments on stack"
	ErrInvalidType      ConstError = "invalid type"
	ErrIndexOutOfBounds ConstError = "index out of bounds"
	ErrMissingConstArg  ConstError = "missing const arg"
	ErrInvalidOpcode    ConstError = "invalid opcode"
)

// TooFewValuesError is an error type wrapping the ErrTooFewValues constant
// error with the opcode which expected arguments.
type TooFewValuesError struct {
	OpCode string
}

func (e TooFewValuesError) Unwrap() error {
	return ErrTooFewValues
}

func (e TooFewValuesError) Error() string {
	return string(ErrTooFewValues) + " for " + e.OpCode
}

// InvalidTypeError is an error type wrapping the ErrInvalidType constant
// error with the opcode which expected arguments.
type InvalidTypeError struct {
	OpCode string
}

func (e InvalidTypeError) Unwrap() error {
	return ErrInvalidType
}

func (e InvalidTypeError) Error() string {
	return string(ErrInvalidType) + " for " + e.OpCode
}

// IndexOutOfBoundsError is an error type wrapping the ErrIndexOutOfBounds
// constant error with the opcode which used the index.
type IndexOutOfBoundsError struct {
	OpCode string
}

func (e IndexOutOfBoundsError) Unwrap() error {
	return ErrIndexOutOfBounds
}

func (e IndexOutOfBoundsError) Error() string {
	return string(ErrIndexOutOfBounds) + " for " + e.OpCode
}

// MissingConstArgError is an error type wrapping the ErrMissingConstArg
// constant error with the opcode which required the const arg.
type MissingConstArgError struct {
	OpCode string
}

func (e MissingConstArgError) Unwrap() error {
	return ErrMissingConstArg
}

func (e MissingConstArgError) Error() string {
	return string(ErrMissingConstArg) + " for " + e.OpCode
}

// InvalidOpcodeError is an error type wrapping the ErrInvalidOpcode constant
// error with the opcode which is unknown.
type InvalidOpcodeError struct {
	OpCode uint32
}

func (e InvalidOpcodeError) Unwrap() error {
	return ErrInvalidOpcode
}

func (e InvalidOpcodeError) Error() string {
	return string(ErrInvalidOpcode) + "0x" + strconv.FormatUint(uint64(e.OpCode), 16)
}
