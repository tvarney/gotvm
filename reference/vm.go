package reference

import (
	"github.com/tvarney/gotvm/op"
	"github.com/tvarney/gotvm/vmerr"
)

// VirtualMachine is a reference implementation of the virtual machine in the
// main package.
//
// The main differences between this VM and the one in the main package are
// that this VM supports stepping through bytecode, and implements the opcodes
// without any attempt at making it fast.
//
// This allows the machine to be much simpler to test, and can be used to
// inspect the execution of a program much closer, at the expense of speed.
type VirtualMachine struct {
	Stack     []interface{}
	FrameBase int

	code op.ByteCode
	idx  int
}

// New returns a new VirtualMachine instance with a pre-allocated stack of 1024
// items.
func New() *VirtualMachine {
	return NewWithSize(1024)
}

// NewWithSize returns a new VirtualMachine instance with a pre-allocated stack
// of the given size.
func NewWithSize(size uint32) *VirtualMachine {
	return &VirtualMachine{
		Stack:     make([]interface{}, 0, int(size)),
		FrameBase: 0,
	}
}

// Push adds the given value to the stack.
//
// The value given _must_ be one of a int64, uint64, or float64. This is not
// checked by the function; pushing a different value will result in errors
// when the value is popped from the stack.
//
// TODO: Support string, list, map, and objects.
func (vm *VirtualMachine) Push(v interface{}) {
	vm.Stack = append(vm.Stack, v)
}

// Pop removes the topmost value from the stack and returns it.
func (vm *VirtualMachine) Pop(opcode string) (interface{}, error) {
	if len(vm.Stack) <= 0 {
		return nil, vmerr.TooFewValuesError{OpCode: opcode}
	}
	v := vm.Stack[len(vm.Stack)-1]
	vm.Stack = vm.Stack[:len(vm.Stack)-1]
	return v, nil
}

// PopInt pops the topmost value from the stack and coerces it to an int.
func (vm *VirtualMachine) PopInt(opcode string) (int64, error) {
	ival, err := vm.Pop(opcode)
	if err != nil {
		return 0, err
	}
	switch v := ival.(type) {
	case int64:
		return v, nil
	case uint64:
		return int64(v), nil
	case float64:
		return int64(v), nil
	}
	return 0, vmerr.InvalidTypeError{OpCode: opcode}
}

// PopUint pops the topmost value from the stack and coerces it to an unsigned
// int.
func (vm *VirtualMachine) PopUint(opcode string) (uint64, error) {
	ival, err := vm.Pop(opcode)
	if err != nil {
		return 0, err
	}
	switch v := ival.(type) {
	case int64:
		return uint64(v), nil
	case uint64:
		return v, nil
	case float64:
		return uint64(v), nil
	}
	return 0, vmerr.InvalidTypeError{OpCode: opcode}
}

// PopFloat pops the topmost value from the stack and coerces it to a float.
func (vm *VirtualMachine) PopFloat(opcode string) (float64, error) {
	ival, err := vm.Pop(opcode)
	if err != nil {
		return 0, err
	}
	switch v := ival.(type) {
	case int64:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case float64:
		return v, nil
	}
	return 0, vmerr.InvalidTypeError{OpCode: opcode}
}
