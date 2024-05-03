package assembler

import (
	"fmt"
	"strings"

	"github.com/tvarney/gotvm"
	"github.com/tvarney/gotvm/op"
)

const (
	ErrInvalidArgCount = gotvm.ConstError("incorrect number of arguments")
)

type Definition struct {
	Name      string
	Value     op.Op
	Arguments []ArgType
}

func newdef(name string, opcode op.Op, args ...ArgType) Definition {
	return Definition{
		Name:      name,
		Value:     opcode,
		Arguments: args,
	}
}

func (d *Definition) Parse(code op.ByteCode, argvalues []rune) (op.ByteCode, error) {
	// Append our bytecode
	code = append(code, d.Value)

	if len(d.Arguments) == 0 && len(argvalues) > 0 {
		return code, fmt.Errorf("%w: %s takes no arguments", ErrInvalidArgCount, d.Name)
	}

	var err error
	for _, arg := range d.Arguments {
		rest, result, parseErr := arg.Parse(argvalues, code)
		argvalues, code = rest, result
		if parseErr != nil && err == nil {
			err = parseErr
		}
	}

	return code, err
}

var (
	ops = []Definition{
		newdef("Noop", op.Noop),
		newdef("Halt", op.Halt),
		newdef("PushI32", op.PushInt32, ArgInt32),
		newdef("PushI64", op.PushInt64, ArgInt64),
		newdef("PushU32", op.PushUint32, ArgUint32),
		newdef("PushU64", op.PushUint64, ArgUint64),
		newdef("PushF32", op.PushFloat32, ArgFloat32),
		newdef("PushF64", op.PushFloat64, ArgFloat64),
		newdef("Pop", op.Pop),
		newdef("PopN", op.PopN, ArgUint32),
		newdef("Copy", op.Copy, ArgUint32),
		newdef("Swap", op.Swap, ArgUint32),
		newdef("Negative", op.Negative),
		newdef("AddInt", op.AddInt),
	}
	definitions = map[string]Definition{}
)

// RemoveComment removes a comment from the line
func RemoveComment(line string) string {
	idx := strings.IndexRune(line, ';')
	if idx < 0 {
		return line
	}
	return line[:idx]
}

// Assemble takes a text file and converts it to bytecode.
//
// The syntax of a line in the assembly is:
//
//	OPCODE [ARG ARG ...] [; Comment]
func Assemble(lines []string, report func(string, int, string)) op.ByteCode {
	// Build our definition table if it hasn't already been done
	if len(definitions) == 0 {
		for _, def := range ops {
			definitions[strings.ToLower(def.Name)] = def
		}
	}

	if report == nil {
		report = func(string, int, string) {}
	}

	size := len(lines) * 3 / 4
	if size < 10 {
		size = 10
	}
	code := make(op.ByteCode, 0, size)

	for idx, line := range lines {
		line = strings.TrimSpace(RemoveComment(line))
		if line == "" {
			// Skip empty lines, comment-only lines, etc.
			continue
		}
		runes := []rune(line)

		opcode, rest, _ := CutSpace(runes)
		def, ok := definitions[strings.ToLower(opcode)]
		if !ok {
			report(fmt.Sprintf("Invalid opcode %q", opcode), idx, line)
			continue
		}

		result, err := def.Parse(code, rest)
		code = result
		if err != nil {
			report(err.Error(), idx, line)
			continue
		}
	}

	return code
}
