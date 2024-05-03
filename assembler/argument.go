package assembler

import (
	"fmt"
	"strconv"

	"github.com/tvarney/gotvm"
	"github.com/tvarney/gotvm/op"
)

const (
	ErrInvalidArgValue = gotvm.ConstError("invalid argument value")
	ErrInvalidArgType  = gotvm.ConstError("internal error: invalid argument type")
)

// ArgType is a value indicating the expected argument type to an opcode.
type ArgType int

const (
	ArgInt32 ArgType = iota
	ArgInt64
	ArgUint32
	ArgUint64
	ArgFloat32
	ArgFloat64
)

var (
	argParseLookup = []func([]rune, op.ByteCode) ([]rune, op.ByteCode, error){
		ParseArgInt32,
		ParseArgInt64,
		ParseArgUint32,
		ParseArgUint64,
		ParseArgFloat32,
		ParseArgFloat64,
	}
)

// Parse takes a 'line' of runes and parses a value according to the arg type.
func (a ArgType) Parse(rest []rune, code op.ByteCode) ([]rune, op.ByteCode, error) {
	if int(a) < 0 || int(a) > len(argParseLookup) {
		return rest, code, fmt.Errorf("%w: Unknown arg type %d", ErrInvalidArgType, int(a))
	}
	return argParseLookup[int(a)](rest, code)
}

// ParseArgInt32 implements the argument parsing logic for a 32-bit int.
func ParseArgInt32(rest []rune, code op.ByteCode) ([]rune, op.ByteCode, error) {
	if rest == nil {
		code = append(code, 0)
		return rest, code, ErrInvalidArgCount
	}
	strval, rest, _ := CutSpace(rest)
	ival, err := ParseInt(strval, 32)
	if err != nil {
		code = append(code, 0)
		return rest, code, err
	}
	code = append(code, op.Op(int32(ival)))
	return rest, code, nil
}

// ParseArgInt64 implements the argument parsing logic for a 64-bit int.
func ParseArgInt64(rest []rune, code op.ByteCode) ([]rune, op.ByteCode, error) {
	if rest == nil {
		code = append(code, 0, 0)
		return rest, code, ErrInvalidArgCount
	}
	strval, rest, _ := CutSpace(rest)
	ival, err := ParseInt(strval, 64)
	if err != nil {
		code = append(code, 0, 0)
		return rest, code, err
	}
	uval := uint64(ival)
	code = append(
		code,
		op.Op(uint32((uval&0xFFFFFFFF00000000)>>32)),
		op.Op(uint32(uval&0x00000000FFFFFFFF)),
	)
	return rest, code, nil
}

// ParseArgUint32 implements the argument parsing logic for a 32-bit unsigned
// int.
func ParseArgUint32(rest []rune, code op.ByteCode) ([]rune, op.ByteCode, error) {
	if rest == nil {
		code = append(code, 0)
		return rest, code, ErrInvalidArgCount
	}
	strval, rest, _ := CutSpace(rest)
	uval, err := ParseUint(strval, 32)
	if err != nil {
		code = append(code, 0)
		return rest, code, err
	}
	code = append(code, op.Op(uval))
	return rest, code, nil
}

// ParseArgUint64 implements the argument parsing logic for a 64-bit unsigned
// int.
func ParseArgUint64(rest []rune, code op.ByteCode) ([]rune, op.ByteCode, error) {
	if rest == nil {
		code = append(code, 0, 0)
		return rest, code, ErrInvalidArgCount
	}
	strval, rest, _ := CutSpace(rest)
	uval, err := ParseUint(strval, 64)
	if err != nil {
		code = append(code, 0, 0)
		return rest, code, err
	}
	code = append(
		code,
		op.Op(uint32((uval&0xFFFFFFFF00000000)>>32)),
		op.Op(uint32(uval&0x00000000FFFFFFFF)),
	)
	return rest, code, nil
}

// ParseArgFloat32 implements the argument parsing logic for a 32-bit float.
func ParseArgFloat32(rest []rune, code op.ByteCode) ([]rune, op.ByteCode, error) {
	if rest == nil {
		code = append(code, 0)
		return rest, code, ErrInvalidArgCount
	}
	strval, rest, _ := CutSpace(rest)
	fval, err := strconv.ParseFloat(strval, 32)
	if err != nil {
		code = append(code, 0)
		return rest, code, fmt.Errorf("%w: invalid float value %q", ErrInvalidArgValue, strval)
	}
	code = append(code, op.Op(uint32(fval)))
	return rest, code, nil
}

// ParseArgFloat64 implements the argument parsing logic for a 64-bit float.
func ParseArgFloat64(rest []rune, code op.ByteCode) ([]rune, op.ByteCode, error) {
	if rest == nil {
		code = append(code, 0)
		return rest, code, ErrInvalidArgCount
	}
	strval, rest, _ := CutSpace(rest)
	fval, err := strconv.ParseFloat(strval, 32)
	if err != nil {
		code = append(code, 0)
		return rest, code, fmt.Errorf("%w: invalid float value %q", ErrInvalidArgValue, strval)
	}
	uval := uint64(fval)
	code = append(
		code,
		op.Op(uint32((uval&0xFFFFFFFF00000000)>>32)),
		op.Op(uint32(uval&0x00000000FFFFFFFF)),
	)
	return rest, code, nil

}
