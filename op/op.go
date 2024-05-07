package op

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
	PopN        // Pop N values from the stack
	PopC        // Pop a constant C values from the stack
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
	AddConstInt // Add a constant value to the topmost value on the stack
	SubConstInt // Subtract a constant value from the topmost value on the stack
	MulConstInt // Multiple the topmost value on the stack by a constant value
	DivConstInt // Divide the topmost value on the stack by a constant value

	// Generic increment/decrement
	Increment // Increment the topmost value on the stack
	Decrement // Decrement the topmost value on the stack

	// Functions
	Call
	NativeCall
)
