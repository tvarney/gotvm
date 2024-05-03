package assembler

import (
	"fmt"
	"math"
	"strconv"

	"github.com/tvarney/gotvm/cerr"
	"github.com/tvarney/gotvm/op"
)

const (
	ErrInvalidArgValue = cerr.Error("invalid argument value")
	ErrInvalidArgType  = cerr.Error("internal error: invalid argument type")
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
		parseArgInt32,
		parseArgInt64,
		parseArgUint32,
		parseArgUint64,
		parseArgFloat32,
		parseArgFloat64,
	}
)

// Parse takes a 'line' of runes and parses a value according to the arg type.
func (a ArgType) Parse(rest []rune, code op.ByteCode) ([]rune, op.ByteCode, error) {
	if int(a) < 0 || int(a) > len(argParseLookup) {
		return rest, code, fmt.Errorf("%w: Unknown arg type %d", ErrInvalidArgType, int(a))
	}
	return argParseLookup[int(a)](rest, code)
}

// parseArgInt32 implements the argument parsing logic for a 32-bit int.
func parseArgInt32(rest []rune, code op.ByteCode) ([]rune, op.ByteCode, error) {
	if len(rest) <= 0 {
		code = append(code, 0)
		return nil, code, ErrInvalidArgCount
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

// parseArgInt64 implements the argument parsing logic for a 64-bit int.
func parseArgInt64(rest []rune, code op.ByteCode) ([]rune, op.ByteCode, error) {
	if len(rest) == 0 {
		code = append(code, 0, 0)
		return nil, code, ErrInvalidArgCount
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

// parseArgUint32 implements the argument parsing logic for a 32-bit unsigned
// int.
func parseArgUint32(rest []rune, code op.ByteCode) ([]rune, op.ByteCode, error) {
	if len(rest) <= 0 {
		code = append(code, 0)
		return nil, code, ErrInvalidArgCount
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

// parseArgUint64 implements the argument parsing logic for a 64-bit unsigned
// int.
func parseArgUint64(rest []rune, code op.ByteCode) ([]rune, op.ByteCode, error) {
	if len(rest) <= 0 {
		code = append(code, 0, 0)
		return nil, code, ErrInvalidArgCount
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

// parseArgFloat32 implements the argument parsing logic for a 32-bit float.
func parseArgFloat32(rest []rune, code op.ByteCode) ([]rune, op.ByteCode, error) {
	if len(rest) <= 0 {
		code = append(code, 0)
		return nil, code, ErrInvalidArgCount
	}
	strval, rest, _ := CutSpace(rest)
	fval, err := strconv.ParseFloat(strval, 32)
	if err != nil {
		code = append(code, 0)
		return rest, code, fmt.Errorf("%w: invalid float value %q", ErrInvalidArgValue, strval)
	}
	code = append(code, op.Op(math.Float32bits(float32(fval))))
	return rest, code, nil
}

// parseArgFloat64 implements the argument parsing logic for a 64-bit float.
func parseArgFloat64(rest []rune, code op.ByteCode) ([]rune, op.ByteCode, error) {
	if len(rest) <= 0 {
		code = append(code, 0, 0)
		return nil, code, ErrInvalidArgCount
	}
	strval, rest, _ := CutSpace(rest)
	fval, err := strconv.ParseFloat(strval, 64)
	if err != nil {
		code = append(code, 0, 0)
		return rest, code, fmt.Errorf("%w: invalid float value %q", ErrInvalidArgValue, strval)
	}
	uval := math.Float64bits(fval)
	code = append(
		code,
		op.Op(uint32((uval&0xFFFFFFFF00000000)>>32)),
		op.Op(uint32(uval&0x00000000FFFFFFFF)),
	)
	return rest, code, nil

}
