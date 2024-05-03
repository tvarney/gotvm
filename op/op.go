package op

type Op uint32

type ByteCode []Op

const (
	Noop = iota
	Halt

	// Stack Operations
	PushInt32   // STACK.push(i32)
	PushInt64   // STACK.push(i64)
	PushUint32  // STACK.push(u32)
	PushUint64  // STACK.push(u64)
	PushFloat32 // STACK.push(f32)
	PushFloat64 // STACK.push(f64)
	Pop         // STACK.pop(1)
	PopN        // STACK.pop(u32)
	Copy        // STACK.push(stack[-int(u32)])
	Swap        // v = STACK[-int(u32)]; STACK[-int(u32)] = STACK[-1]; STACK[-1] = v

	// Unary Operations
	Negative // STACK[-1] = -STACK[-1]

	// Binary Operations
	AddInt
	SubInt
	MulInt
	DivInt

	// Functions
	Call
	NativeCall
)
