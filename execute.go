package gotvm

import (
	"github.com/tvarney/gotvm/cerr"
	"github.com/tvarney/gotvm/op"
)

// Execute takes a chunk of ByteCode and runs it in the VirtualMachine
// instance.
func (vm *VirtualMachine) Execute(code op.ByteCode) error {
	vm.FrameBase = 0

	// This is the hot loop and should be optimized heavily.
	//
	// Check for the compiler inserting bounds checks via:
	//      go build -gcflags="-d=ssa/check_bce" .\pkg/vm
	//
	// At some point measuring how using functions effect this loop should be
	// done. Generally a function call introduces a lot of extra overhead for
	// shorter functions which makes it slower, but this function is very
	// likely too big for the instruction cache - making it shorter could help
	// keep the loop resident in the instruction cache.
	idx := 0
	for idx >= 0 && idx < len(code) {
		opcode := code[idx]
		switch opcode {
		case op.Noop:
			idx++
		case op.PushInt32:
			v, err := constArgU32(code, idx+1)
			if err != nil {
				return cerr.Error("missing const arg to op.PushInt32")
			}
			vm.push(int64(int32(v)))
			idx += 2
		case op.PushInt64:
			v, err := constArgU64(code, idx+1)
			if err != nil {
				return cerr.Error("missing const arg to op.PushInt64")
			}
			vm.push(int64(v))
			idx += 3
		case op.PushUint32:
			v, err := constArgU32(code, idx+1)
			if err != nil {
				return cerr.Error("missing const arg to op.PushUint32")
			}
			vm.push(uint64(v))
			idx += 2
		case op.PushUint64:
			v, err := constArgU64(code, idx+1)
			if err != nil {
				return cerr.Error("missing const arg to op.PushUint64")
			}
			vm.push(uint64(v))
			idx += 3
		case op.Pop:
			if len(vm.Stack) <= 0 || idx < 0 {
				return cerr.Error("too few arguments on stack for op.Pop")
			}
			vm.Stack = vm.Stack[:len(vm.Stack)-1]
			idx++
		case op.PopN:
			v, err := constArgU32(code, idx+1)
			if err != nil {
				return cerr.Error("missing const arg to op.PopN")
			}
			last := len(vm.Stack) - int(v)
			if last < 0 || last >= len(vm.Stack) {
				return cerr.Error("too few arguments on stack for op.PopN")
			}
			vm.Stack = vm.Stack[:last]
			idx += 2
		case op.Copy:
			v, err := constArgU32(code, idx+1)
			if err != nil {
				return cerr.Error("missing const arg to op.Copy")
			}
			ref := vm.FrameBase + int(v)
			if ref >= len(vm.Stack) || ref < 0 {
				return cerr.Error("copy index out of bounds")
			}
			vm.push(vm.Stack[ref])
			idx += 2
		case op.Swap:
			v, err := constArgU32(code, idx+1)
			if err != nil {
				return cerr.Error("missing const arg to op.Swap")
			}
			ref := vm.FrameBase + int(v)
			if ref >= len(vm.Stack) || ref < 0 {
				return cerr.Error("swap index out of bounds")
			}
			last := len(vm.Stack) - 1
			if last <= 0 {
				return cerr.Error("nothing in stack to swap for op.Swap")
			}
			vm.Stack[ref], vm.Stack[last] = vm.Stack[last], vm.Stack[ref]
			idx += 2
		case op.Negative:
			if len(vm.Stack) < 1 {
				return cerr.Error("too few arguments on stack for op.Negative")
			}
			iv := vm.Stack[len(vm.Stack)-1]
			vm.Stack = vm.Stack[:len(vm.Stack)-1]
			switch v := iv.(type) {
			case float64:
				vm.Stack = append(vm.Stack, -v)
			case uint64:
				vm.Stack = append(vm.Stack, uint64(-int64(v)))
			case int64:
				vm.Stack = append(vm.Stack, -v)
			default:
				return cerr.Error("invalid type for op.Negative")
			}
			idx++
		case op.AddInt:
			if len(vm.Stack) < 2 {
				return cerr.Error("too few arguments on stack for op.AddInt")
			}
			v1, err := coerceInt(vm.Stack[len(vm.Stack)-1])
			if err != nil {
				return err
			}
			v2, err := coerceInt(vm.Stack[len(vm.Stack)-2])
			if err != nil {
				return err
			}
			vm.Stack = vm.Stack[:len(vm.Stack)-2]
			vm.Stack = append(vm.Stack, v1+v2)
			idx += 1
		default:
			idx++
		}
	}
	return nil
}
