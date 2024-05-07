package reference

import (
	"github.com/tvarney/gotvm/op"
	"github.com/tvarney/gotvm/vmerr"
)

const (
	ErrHalt vmerr.ConstError = "halt"
)

// Execute takes a chunk of ByteCode and runs it in the VirtualMachine
// instance.
//
// Under the hood this calls the (*VirtualMachine).Start() function with the
// code, then calls (*VirtualMachine).Step() until it returns an error.
func (vm *VirtualMachine) Execute(code op.ByteCode) error {
	vm.Start(code)
	for {
		if err := vm.Step(); err != nil {
			if err == ErrHalt {
				return nil
			}
			return err
		}
	}
}

// Start initializes the VirtualMachine instance to run the given bytecode.
func (vm *VirtualMachine) Start(code op.ByteCode) {
	vm.code = code
	vm.idx = 0
}

// Step executes the next opcode in the code that the VirtualMachine instance
// is running.
func (vm *VirtualMachine) Step() error {
	if vm.idx < 0 || vm.idx >= len(vm.code) {
		return ErrHalt
	}

	switch vm.code[vm.idx] {
	case op.Noop:
		vm.idx++
	case op.Halt:
		vm.idx++
		return ErrHalt
	case op.PushInt32:
		return vm.OpPushInt32()
	case op.PushInt64:
		return vm.OpPushInt64()
	case op.PushUint32:
		return vm.OpPushUint32()
	case op.PushUint64:
		return vm.OpPushUint64()
	case op.PushFloat32:
		return vm.OpPushFloat32()
	case op.PushFloat64:
		return vm.OpPushFloat64()
	case op.Pop:
		return vm.OpPop()
	case op.PopN:
		return vm.OpPopN()
	case op.Copy:
		return vm.OpCopy()
	case op.Swap:
		return vm.OpSwap()
	case op.Negative:
		return vm.OpNegative()
	case op.AddInt:
		return vm.OpAddInt()
	case op.SubInt:
		return vm.OpSubInt()
	case op.MulInt:
		return vm.OpMulInt()
	case op.DivInt:
		return vm.OpDivInt()
	case op.AddConstInt32:
		return vm.OpAddConstInt32()
	case op.AddConstInt64:
		return vm.OpAddConstInt64()
	case op.AddConstUint32:
		return vm.OpAddConstUint32()
	case op.AddConstUint64:
		return vm.OpAddConstUint64()
	case op.AddConstFloat32:
		return vm.OpAddConstFloat32()
	case op.AddConstFloat64:
		return vm.OpAddConstFloat64()
	case op.SubConstInt32:
		return vm.OpSubConstInt32()
	case op.SubConstInt64:
		return vm.OpSubConstInt64()
	case op.SubConstUint32:
		return vm.OpSubConstUint32()
	case op.SubConstUint64:
		return vm.OpSubConstUint64()
	case op.SubConstFloat32:
		return vm.OpSubConstFloat32()
	case op.SubConstFloat64:
		return vm.OpSubConstFloat64()
	case op.MulConstInt32:
		return vm.OpMulConstInt32()
	case op.MulConstInt64:
		return vm.OpMulConstInt64()
	case op.MulConstUint32:
		return vm.OpMulConstUint32()
	case op.MulConstUint64:
		return vm.OpMulConstUint64()
	case op.MulConstFloat32:
		return vm.OpMulConstFloat32()
	case op.MulConstFloat64:
		return vm.OpMulConstFloat64()
	case op.DivConstInt32:
		return vm.OpDivConstInt32()
	case op.DivConstInt64:
		return vm.OpDivConstInt64()
	case op.DivConstUint32:
		return vm.OpDivConstUint32()
	case op.DivConstUint64:
		return vm.OpDivConstUint64()
	case op.DivConstFloat32:
		return vm.OpDivConstFloat32()
	case op.DivConstFloat64:
		return vm.OpDivConstFloat64()
	case op.Increment:
		return vm.OpIncrement()
	case op.Decrement:
		return vm.OpDecrement()
	default:
		return vmerr.InvalidOpcodeError{OpCode: uint32(vm.code[vm.idx])}
	}

	return nil
}
