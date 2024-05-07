package reference

import (
	"math"

	"github.com/tvarney/gotvm/op"
	"github.com/tvarney/gotvm/vmerr"
)

// OpPushInt32 implements the PushI32 opcode for the reference VM.
//
// This function will take the next value in the bytecode and convert it to an
// int32 value before pushing it to the stack as an int64.
func (vm *VirtualMachine) OpPushInt32() error {
	value, err := op.ConstArgU32(vm.code, vm.idx+1)
	if err != nil {
		return err
	}
	vm.Push(int64(int32(value)))
	vm.idx += 2
	return nil
}

// OpPushInt64 implements the PushI64 opcode for the reference VM.
//
// This function will take the next 2 values in the bytecode and convert them
// to an int64 value before pushing it to the stack.
func (vm *VirtualMachine) OpPushInt64() error {
	value, err := op.ConstArgU64(vm.code, vm.idx+1)
	if err != nil {
		return err
	}
	vm.Push(int64(value))
	vm.idx += 3
	return nil
}

// OpPushUint32 implements the PushU32 opcode for the reference VM.
//
// This function will take the next value in the bytecode and convert it to a
// uint32 value before pushing it to the stack as a uint64.
func (vm *VirtualMachine) OpPushUint32() error {
	value, err := op.ConstArgU32(vm.code, vm.idx+1)
	if err != nil {
		return err
	}
	vm.Push(uint64(value))
	vm.idx += 2
	return nil
}

// OpPushUint64 implements the PushU64 opcode for the reference VM.
//
// This function will take the next 2 values in the bytecode and convert them
// to a uint64 value before pushing it to the stack.
func (vm *VirtualMachine) OpPushUint64() error {
	value, err := op.ConstArgU64(vm.code, vm.idx+1)
	if err != nil {
		return err
	}
	vm.Push(value)
	vm.idx += 3
	return nil
}

// OpPushFloat32 implements the PushF32 opcode for the reference VM.
//
// This function will take the next value in the bytecode and convert it to a
// float32 value before pushing it to the stack as a float64.
func (vm *VirtualMachine) OpPushFloat32() error {
	value, err := op.ConstArgU32(vm.code, vm.idx+1)
	if err != nil {
		return err
	}
	vm.Push(float64(math.Float32frombits(value)))
	vm.idx += 2
	return nil
}

// OpPushFloat64 implements the PushF64 opcode for the reference VM.
//
// This function will take the next 2 values in the bytecode and convert them
// to a float64 value before pushing it to the stack.
func (vm *VirtualMachine) OpPushFloat64() error {
	value, err := op.ConstArgU64(vm.code, vm.idx+1)
	if err != nil {
		return err
	}
	vm.Push(math.Float64frombits(value))
	vm.idx += 3
	return nil
}

// OpPop implements the Pop opcode for the reference VM.
//
// This function will remove the topmost value of the stack. If popping the
// value would reveal the FrameBase value on the stack, this function generates
// an index out of bounds error.
func (vm *VirtualMachine) OpPop() error {
	if len(vm.Stack)-1 == vm.FrameBase {
		return vmerr.IndexOutOfBoundsError{OpCode: "Pop"}
	}
	_, err := vm.Pop("Pop")
	return err
}

// OpPopN implements the PopN opcode for the reference VM.
//
// This function will take the next value in the bytecode and convert it to a
// u32 `N`, then pop that many values from the stack.
//
// If popping that many values from the stack would result in popping past the
// FrameBase value, an index out of bounds error is generated.
func (vm *VirtualMachine) OpPopN() error {
	n, err := op.ConstArgU32(vm.code, vm.idx+1)
	if err != nil {
		return err
	}

	last := len(vm.Stack) - int(n)
	if last < 0 || last <= vm.FrameBase || last >= len(vm.Stack) {
		return vmerr.IndexOutOfBoundsError{OpCode: "PopN"}
	}

	vm.Stack = vm.Stack[:last]
	vm.idx += 2
	return nil
}

// OpCopy implements the Copy opcode for the reference VM.
//
// This function will take the next value in the bytecode and convert it to a
// uint32 `OFFSET` which indicates the element offset from FrameBase to copy
// to the top of the stack.
func (vm *VirtualMachine) OpCopy() error {
	offset, err := op.ConstArgU32(vm.code, vm.idx+1)
	if err != nil {
		return err
	}
	idx := vm.FrameBase + int(offset)
	if idx < 0 || idx >= len(vm.Stack) {
		return vmerr.IndexOutOfBoundsError{OpCode: "Copy"}
	}
	vm.Push(vm.Stack[idx])
	vm.idx += 2
	return nil
}

// OpSwap implements the Swap opcode for the reference VM.
//
// This function will take the next value in the bytecode and convert it to a
// uint32 `OFFSET` which indicates the element offset from FrameBase to swap
// with the topmost value on the stack.
//
// As an example, given the stack [0, 1, 2, 3, 4] with FrameBase 1, a OpSwap
// with `OFFSET` of 1 will swap the values at indices 2 and 4, for a final
// stack of [0, 1, 4, 3, 2].
func (vm *VirtualMachine) OpSwap() error {
	offset, err := op.ConstArgU32(vm.code, vm.idx+1)
	if err != nil {
		return err
	}
	idx := vm.FrameBase + int(offset)
	last := len(vm.Stack) - 1
	if idx < 0 || idx >= len(vm.Stack) || last < 0 || last >= len(vm.Stack) {
		return vmerr.IndexOutOfBoundsError{OpCode: "Swap"}
	}
	vm.Stack[idx], vm.Stack[last] = vm.Stack[last], vm.Stack[idx]
	vm.idx += 2
	return nil
}

// OpNegative implements the Negative opcode for the reference VM.
//
// This function will take the topmost value of the stack and negate it, then
// push it back to the stack. If the value is a uint64, this will the negative
// bit pattern.
func (vm *VirtualMachine) OpNegative() error {
	top, err := vm.Pop("Negative")
	if err != nil {
		return err
	}

	switch v := top.(type) {
	case int64:
		vm.Push(-v)
	case uint64:
		vm.Push(-v)
	case float64:
		vm.Push(-v)
	default:
		return vmerr.InvalidTypeError{OpCode: "Negative"}
	}
	vm.idx++
	return nil
}

// OpAddInt implements the AddInt opcode for the reference VM.
func (vm *VirtualMachine) OpAddInt() error {
	v1, err := vm.PopInt("AddInt")
	if err != nil {
		return err
	}
	v2, err := vm.PopInt("AddInt")
	if err != nil {
		return err
	}
	vm.Push(v1 + v2)
	vm.idx++
	return nil
}

// OpSubInt implements the SubInt opcode for the reference VM.
func (vm *VirtualMachine) OpSubInt() error {
	v1, err := vm.PopInt("AddInt")
	if err != nil {
		return err
	}
	v2, err := vm.PopInt("AddInt")
	if err != nil {
		return err
	}
	vm.Push(v1 - v2)
	vm.idx++
	return nil
}

// OpMulInt implements the MulInt opcode for the reference VM.
func (vm *VirtualMachine) OpMulInt() error {
	v1, err := vm.PopInt("AddInt")
	if err != nil {
		return err
	}
	v2, err := vm.PopInt("AddInt")
	if err != nil {
		return err
	}
	vm.Push(v1 * v2)
	vm.idx++
	return nil
}

// OpDivInt implements the DivInt opcode for the reference VM.
func (vm *VirtualMachine) OpDivInt() error {
	v1, err := vm.PopInt("AddInt")
	if err != nil {
		return err
	}
	v2, err := vm.PopInt("AddInt")
	if err != nil {
		return err
	}
	vm.Push(v1 / v2)
	vm.idx++
	return nil
}

// OpAddConstInt32 implements the AddConstI32 opcode for the reference VM.
func (vm *VirtualMachine) OpAddConstInt32() error {
	v1, err := op.ConstArgU32(vm.code, vm.idx+1)
	if err != nil {
		return vmerr.MissingConstArgError{OpCode: "AddConstI32"}
	}
	v2, err := vm.PopInt("AddConstI32")
	if err != nil {
		return err
	}
	vm.Push(int64(v1) + v2)
	vm.idx += 2
	return nil
}

// OpAddConstInt64 implements the AddConstI64 opcode for the reference VM.
func (vm *VirtualMachine) OpAddConstInt64() error {
	v1, err := op.ConstArgI64(vm.code, vm.idx+1)
	if err != nil {
		return vmerr.MissingConstArgError{OpCode: "AddConstI64"}
	}
	v2, err := vm.PopInt("AddConstI64")
	if err != nil {
		return err
	}
	vm.Push(v1 + v2)
	vm.idx += 3
	return nil
}

// OpAddConstUint32 implements the AddConstU32 opcode for the reference VM.
func (vm *VirtualMachine) OpAddConstUint32() error {
	v1, err := op.ConstArgU32(vm.code, vm.idx+1)
	if err != nil {
		return vmerr.MissingConstArgError{OpCode: "AddConstU32"}
	}
	v2, err := vm.PopUint("AddConstU32")
	if err != nil {
		return err
	}
	vm.Push(uint64(v1) + v2)
	vm.idx += 2
	return nil
}

// OpAddConstUint64 implements the AddConstU64 opcode for the reference VM.
func (vm *VirtualMachine) OpAddConstUint64() error {
	v1, err := op.ConstArgU64(vm.code, vm.idx+1)
	if err != nil {
		return vmerr.MissingConstArgError{OpCode: "AddConstU64"}
	}
	v2, err := vm.PopUint("AddConstU64")
	if err != nil {
		return err
	}
	vm.Push(v1 + v2)
	vm.idx += 3
	return nil
}

// OpAddConstFloat32 implements the AddConstF32 opcode for the reference VM.
func (vm *VirtualMachine) OpAddConstFloat32() error {
	v1, err := op.ConstArgU32(vm.code, vm.idx+1)
	if err != nil {
		return vmerr.MissingConstArgError{OpCode: "AddConstF32"}
	}
	v2, err := vm.PopFloat("AddConstF32")
	if err != nil {
		return err
	}
	vm.Push(float64(math.Float32frombits(v1)) + v2)
	vm.idx += 2
	return nil
}

// OpAddConstFloat64 implements the AddConstF64 opcode for the reference VM.
func (vm *VirtualMachine) OpAddConstFloat64() error {
	v1, err := op.ConstArgU64(vm.code, vm.idx+1)
	if err != nil {
		return vmerr.MissingConstArgError{OpCode: "AddConstF64"}
	}
	v2, err := vm.PopFloat("AddConstF64")
	if err != nil {
		return err
	}
	vm.Push(math.Float64frombits(v1) + v2)
	vm.idx += 3
	return nil
}

// OpSubConstInt32 implements the SubConstI32 opcode for the reference VM.
func (vm *VirtualMachine) OpSubConstInt32() error {
	v1, err := op.ConstArgI32(vm.code, vm.idx+1)
	if err != nil {
		return vmerr.MissingConstArgError{OpCode: "SubConstI32"}
	}
	v2, err := vm.PopInt("SubConstI32")
	if err != nil {
		return err
	}
	vm.Push(v2 - int64(v1))
	vm.idx += 2
	return nil
}

// OpSubConstInt64 implements the SubConstI64 opcode for the reference VM.
func (vm *VirtualMachine) OpSubConstInt64() error {
	v1, err := op.ConstArgI64(vm.code, vm.idx+1)
	if err != nil {
		return vmerr.MissingConstArgError{OpCode: "SubConstI64"}
	}
	v2, err := vm.PopInt("SubConstI64")
	if err != nil {
		return err
	}
	vm.Push(v2 - v1)
	vm.idx += 3
	return nil
}

// OpSubConstUint32 implements the SubConstU32 opcode for the reference VM.
func (vm *VirtualMachine) OpSubConstUint32() error {
	v1, err := op.ConstArgU32(vm.code, vm.idx+1)
	if err != nil {
		return vmerr.MissingConstArgError{OpCode: "SubConstU32"}
	}
	v2, err := vm.PopUint("SubConstU32")
	if err != nil {
		return err
	}
	vm.Push(v2 - uint64(v1))
	vm.idx += 2
	return nil
}

// OpSubConstUint64 implements the SubConstU64 opcode for the reference VM.
func (vm *VirtualMachine) OpSubConstUint64() error {
	v1, err := op.ConstArgU64(vm.code, vm.idx+1)
	if err != nil {
		return vmerr.MissingConstArgError{OpCode: "SubConstU64"}
	}
	v2, err := vm.PopUint("SubConstU64")
	if err != nil {
		return err
	}
	vm.Push(v2 - v1)
	vm.idx += 3
	return nil
}

// OpSubConstFloat32 implements the SubConstF32 opcode for the reference VM.
func (vm *VirtualMachine) OpSubConstFloat32() error {
	v1, err := op.ConstArgU32(vm.code, vm.idx+1)
	if err != nil {
		return vmerr.MissingConstArgError{OpCode: "SubConstF32"}
	}
	v2, err := vm.PopFloat("SubConstF32")
	if err != nil {
		return err
	}
	vm.Push(v2 - float64(math.Float32frombits(v1)))
	vm.idx += 2
	return nil
}

// OpSubConstFloat64 implements the SubConstF64 opcode for the reference VM.
func (vm *VirtualMachine) OpSubConstFloat64() error {
	v1, err := op.ConstArgU64(vm.code, vm.idx+1)
	if err != nil {
		return vmerr.MissingConstArgError{OpCode: "SubConstF64"}
	}
	v2, err := vm.PopFloat("SubConstF64")
	if err != nil {
		return err
	}
	vm.Push(v2 - math.Float64frombits(v1))
	vm.idx += 3
	return nil
}

// OpMulConstInt32 implements the MulConstI32 opcode for the reference VM.
func (vm *VirtualMachine) OpMulConstInt32() error {
	v1, err := op.ConstArgI32(vm.code, vm.idx+1)
	if err != nil {
		return vmerr.MissingConstArgError{OpCode: "MulConstI32"}
	}
	v2, err := vm.PopInt("MulConstI32")
	if err != nil {
		return err
	}
	vm.Push(int64(v1) * v2)
	vm.idx += 2
	return nil
}

// OpMulConstInt64 implements the MulConstI64 opcode for the reference VM.
func (vm *VirtualMachine) OpMulConstInt64() error {
	v1, err := op.ConstArgI64(vm.code, vm.idx+1)
	if err != nil {
		return vmerr.MissingConstArgError{OpCode: "MulConstI64"}
	}
	v2, err := vm.PopInt("MulConstI64")
	if err != nil {
		return err
	}
	vm.Push(v1 * v2)
	vm.idx += 3
	return nil
}

// OpMulConstUint32 implements the MulConstU32 opcode for the reference VM.
func (vm *VirtualMachine) OpMulConstUint32() error {
	v1, err := op.ConstArgU32(vm.code, vm.idx+1)
	if err != nil {
		return vmerr.MissingConstArgError{OpCode: "MulConstU32"}
	}
	v2, err := vm.PopUint("MulConstU32")
	if err != nil {
		return err
	}
	vm.Push(uint64(v1) * v2)
	vm.idx += 2
	return nil
}

// OpMulConstUint64 implements the MulConstU64 opcode for the reference VM.
func (vm *VirtualMachine) OpMulConstUint64() error {
	v1, err := op.ConstArgU64(vm.code, vm.idx+1)
	if err != nil {
		return vmerr.MissingConstArgError{OpCode: "MulConstU64"}
	}
	v2, err := vm.PopUint("MulConstU64")
	if err != nil {
		return err
	}
	vm.Push(v1 * v2)
	vm.idx += 3
	return nil
}

// OpMulConstFloat32 implements the MulConstF32 opcode for the reference VM.
func (vm *VirtualMachine) OpMulConstFloat32() error {
	v1, err := op.ConstArgU32(vm.code, vm.idx+1)
	if err != nil {
		return vmerr.MissingConstArgError{OpCode: "MulConstF32"}
	}
	v2, err := vm.PopFloat("MulConstF32")
	if err != nil {
		return err
	}
	vm.Push(float64(math.Float32frombits(v1)) * v2)
	vm.idx += 2
	return nil
}

// OpMulConstFloat64 implements the MulConstF64 opcode for the reference VM.
func (vm *VirtualMachine) OpMulConstFloat64() error {
	v1, err := op.ConstArgU64(vm.code, vm.idx+1)
	if err != nil {
		return vmerr.MissingConstArgError{OpCode: "MulConstF64"}
	}
	v2, err := vm.PopFloat("MulConstF64")
	if err != nil {
		return err
	}
	vm.Push(math.Float64frombits(v1) * v2)
	vm.idx += 3
	return nil
}

// OpDivConstInt32 implements the DivConstI32 opcode for the reference VM.
func (vm *VirtualMachine) OpDivConstInt32() error {
	v1, err := op.ConstArgI32(vm.code, vm.idx+1)
	if err != nil {
		return vmerr.MissingConstArgError{OpCode: "DivConstI32"}
	}
	v2, err := vm.PopInt("DivConstI32")
	if err != nil {
		return err
	}
	vm.Push(v2 / int64(v1))
	vm.idx += 2
	return nil
}

// OpDivConstInt64 implements the DivConstI64 opcode for the reference VM.
func (vm *VirtualMachine) OpDivConstInt64() error {
	v1, err := op.ConstArgI64(vm.code, vm.idx+1)
	if err != nil {
		return vmerr.MissingConstArgError{OpCode: "DivConstI64"}
	}
	v2, err := vm.PopInt("DivConstI64")
	if err != nil {
		return err
	}
	vm.Push(v2 / v1)
	vm.idx += 3
	return nil
}

// OpDivConstUint32 implements the DivConstU32 opcode for the reference VM.
func (vm *VirtualMachine) OpDivConstUint32() error {
	v1, err := op.ConstArgU32(vm.code, vm.idx+1)
	if err != nil {
		return vmerr.MissingConstArgError{OpCode: "DivConstU32"}
	}
	v2, err := vm.PopUint("DivConstU32")
	if err != nil {
		return err
	}
	vm.Push(v2 / uint64(v1))
	vm.idx += 2
	return nil
}

// OpDivConstUint64 implements the DivConstU64 opcode for the reference VM.
func (vm *VirtualMachine) OpDivConstUint64() error {
	v1, err := op.ConstArgU64(vm.code, vm.idx+1)
	if err != nil {
		return vmerr.MissingConstArgError{OpCode: "DivConstU64"}
	}
	v2, err := vm.PopUint("DivConstU64")
	if err != nil {
		return err
	}
	vm.Push(v2 / v1)
	vm.idx += 3
	return nil
}

// OpDivConstFloat32 implements the DivConstF32 opcode for the reference VM.
func (vm *VirtualMachine) OpDivConstFloat32() error {
	v1, err := op.ConstArgU32(vm.code, vm.idx+1)
	if err != nil {
		return vmerr.MissingConstArgError{OpCode: "DivConstF32"}
	}
	v2, err := vm.PopFloat("DivConstF32")
	if err != nil {
		return err
	}
	vm.Push(v2 / float64(math.Float32frombits(v1)))
	vm.idx += 2
	return nil
}

// OpDivConstFloat64 implements the DivConstF64 opcode for the reference VM.
func (vm *VirtualMachine) OpDivConstFloat64() error {
	v1, err := op.ConstArgU64(vm.code, vm.idx+1)
	if err != nil {
		return vmerr.MissingConstArgError{OpCode: "DivConstF64"}
	}
	v2, err := vm.PopFloat("DivConstF64")
	if err != nil {
		return err
	}
	vm.Push(v2 / math.Float64frombits(v1))
	vm.idx += 3
	return nil
}

// OpIncrement implements the Increment opcode for the reference VM.
func (vm *VirtualMachine) OpIncrement() error {
	top, err := vm.Pop("Increment")
	if err != nil {
		return err
	}
	switch v := top.(type) {
	case int64:
		vm.Push(v + 1)
	case uint64:
		vm.Push(v + 1)
	case float64:
		vm.Push(v + 1)
	default:
		return vmerr.InvalidTypeError{OpCode: "Increment"}
	}
	vm.idx++
	return nil
}

// OpDecrement implements the Decrement opcode for the reference VM.
func (vm *VirtualMachine) OpDecrement() error {
	top, err := vm.Pop("Decrement")
	if err != nil {
		return err
	}
	switch v := top.(type) {
	case int64:
		vm.Push(v - 1)
	case uint64:
		vm.Push(v - 1)
	case float64:
		vm.Push(v - 1)
	default:
		return vmerr.InvalidTypeError{OpCode: "Decrement"}
	}
	vm.idx++
	return nil
}
