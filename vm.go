package gotvm

import (
	"github.com/tvarney/gotvm/cerr"
	"github.com/tvarney/gotvm/op"
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
	return 0, cerr.Error("can't coerce value to int")
}

func constArgU32(code op.ByteCode, idx int) (uint32, error) {
	if idx < 0 || idx >= len(code) {
		return 0, cerr.Error("malformed bytecode; missing const arg")
	}
	return uint32(code[idx]), nil
}

func constArgU64(code op.ByteCode, idx int) (uint64, error) {
	idx2 := idx + 1
	// If we just used unsigned values for indices, we could simplify this
	// *so* much - we have to be overly verbose to account for overflow to
	// avoid runtime.panicIndex() being inserted
	if idx < 0 || idx >= len(code) || idx2 < 0 || idx2 >= len(code) {
		return 0, cerr.Error("malformed bytecode; missing const arg")
	}
	return uint64(code[idx])<<32 + uint64(code[idx2]), nil
}
