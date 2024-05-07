package gotvm

import (
	"github.com/tvarney/gotvm/vmerr"
)

// VirtualMachine is a stack based VM which executes bytecode.
type VirtualMachine struct {
	Stack     []interface{}
	FrameBase int
}

// New returns a new VirtualMachine instance with a pre-allocated stack.
func New() *VirtualMachine {
	return &VirtualMachine{
		Stack:     make([]interface{}, 0, 1024),
		FrameBase: 0,
	}
}

// NewWithSize returns a new VirtualMachine instance with a Stack pre-allocated
// to the given size.
func NewWithSize(size int) *VirtualMachine {
	if size < 0 {
		size = 1024
	}
	return &VirtualMachine{
		Stack:     make([]interface{}, 0, size),
		FrameBase: 0,
	}
}

// Helpers
// =======

func (vm *VirtualMachine) push(v interface{}) {
	vm.Stack = append(vm.Stack, v)
}

func coerceInt(v interface{}) (int64, error) {
	switch value := v.(type) {
	case int64:
		return value, nil
	case uint64:
		return int64(value), nil
	case float64:
		return int64(value), nil
	}
	return 0, vmerr.ConstError("can't coerce value to int")
}
