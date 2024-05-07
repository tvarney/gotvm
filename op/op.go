package op

import (
	"github.com/tvarney/gotvm/vmerr"
)

type Op uint32

type ByteCode []Op

const (
	Noop = iota
	Halt

	// Stack Operations
	PushInt32   // Push a constant i32 to the stack
	PushInt64   // Push a constant i64 to the stack
	PushUint32  // Push a constant u32 to the stack
	PushUint64  // Push a constant u64 to the stack
	PushFloat32 // Push a constant f32 to the stack
	PushFloat64 // Push a constant f64 to the stack
	Pop         // Pop a single value from the stack
	PopN        // Pop constant N values from the stack
	Copy        // Copy a value within the stack
	Swap        // Swap the topmost value on the stack with another

	// Unary Operations
	Negative // Negate the topmost value on the stack

	// Binary Operations
	AddInt // Add the topmost two values on the stack
	SubInt // Subtract the topmost two values on the stack
	MulInt // Multiply the topmost two values on the stack
	DivInt // Divide the topmost two values on the stack

	// Const Binary Operations
	AddConstInt32   // Add a constant value to the topmost value on the stack
	AddConstInt64   // Add a constant value to the topmost value on the stack
	AddConstUint32  // Add a constant value to the topmost value on the stack
	AddConstUint64  // Add a constant value to the topmost value on the stack
	AddConstFloat32 // Add a constant value to the topmost value on the stack
	AddConstFloat64 // Add a constant value to the topmost value on the stack

	SubConstInt32   // Subtract a constant value from the topmost value on the stack
	SubConstInt64   // Subtract a constant value from the topmost value on the stack
	SubConstUint32  // Subtract a constant value from the topmost value on the stack
	SubConstUint64  // Subtract a constant value from the topmost value on the stack
	SubConstFloat32 // Subtract a constant value from the topmost value on the stack
	SubConstFloat64 // Subtract a constant value from the topmost value on the stack

	MulConstInt32   // Multiple the topmost value on the stack by a constant value
	MulConstInt64   // Multiple the topmost value on the stack by a constant value
	MulConstUint32  // Multiple the topmost value on the stack by a constant value
	MulConstUint64  // Multiple the topmost value on the stack by a constant value
	MulConstFloat32 // Multiple the topmost value on the stack by a constant value
	MulConstFloat64 // Multiple the topmost value on the stack by a constant value

	DivConstInt32   // Divide the topmost value on the stack by a constant value
	DivConstInt64   // Divide the topmost value on the stack by a constant value
	DivConstUint32  // Divide the topmost value on the stack by a constant value
	DivConstUint64  // Divide the topmost value on the stack by a constant value
	DivConstFloat32 // Divide the topmost value on the stack by a constant value
	DivConstFloat64 // Divide the topmost value on the stack by a constant value

	// Generic increment/decrement
	Increment // Increment the topmost value on the stack
	Decrement // Decrement the topmost value on the stack

	// Functions
	Call
	NativeCall
)

// ConstArgU32 converts the Op value at the given index to a uint32 value.
//
// The intended use of this is when an opcode takes a constant argument which
// is encoded in the bytecode itself; this function may be used to fetch and
// convert the value.
func ConstArgU32(code ByteCode, idx int) (uint32, error) {
	if idx < 0 || idx >= len(code) {
		return 0, vmerr.ErrMissingConstArg
	}
	return uint32(code[idx]), nil
}

// ConstArgU64 converts the OpValue at the given index to a uint64 value.
//
// The intended use of this is when an opcode takes a constant argument which
// is encoded in the bytecode itself; this function may be used to fetch and
// convert the value.
func ConstArgU64(code ByteCode, idx int) (uint64, error) {
	idx2 := idx + 1
	// If we just used unsigned values for indices, we could simplify this
	// *so* much - we have to be overly verbose to account for overflow to
	// avoid runtime.panicIndex() being inserted
	if idx < 0 || idx >= len(code) || idx2 < 0 || idx2 >= len(code) {
		return 0, vmerr.ErrMissingConstArg
	}
	return uint64(code[idx])<<32 + uint64(code[idx2]), nil
}

// ConstArgI32 converts the OpValue at the given index to an int32 value.
func ConstArgI32(code ByteCode, idx int) (int32, error) {
	if idx < 0 || idx >= len(code) {
		return 0, vmerr.ErrMissingConstArg
	}
	return int32(uint32(code[idx])), nil
}

// ConstArgI64 converts the OpValue at the given index to an int64 value.
func ConstArgI64(code ByteCode, idx int) (int64, error) {
	idx2 := idx + 1
	if idx < 0 || idx >= len(code) || idx2 < 0 || idx2 >= len(code) {
		return 0, vmerr.ErrMissingConstArg
	}
	return int64(uint64(code[idx])<<32 + uint64(code[idx2])), nil

}
