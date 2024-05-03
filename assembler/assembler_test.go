package assembler_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tvarney/gotvm/assembler"
	"github.com/tvarney/gotvm/op"
)

func TestDefinition(t *testing.T) {
	t.Parallel()

	r := func(s string) []rune {
		return []rune(s)
	}
	b := func(values ...op.Op) op.ByteCode {
		return values
	}

	defNoArgs := assembler.Definition{
		Name:      "Noop",
		Value:     0,
		Arguments: nil,
	}
	defOneArg := assembler.Definition{
		Name:      "PushI32",
		Value:     1,
		Arguments: []assembler.ArgType{assembler.ArgInt32},
	}
	defTwoArgs := assembler.Definition{
		Name:  "Swap",
		Value: 2,
		Arguments: []assembler.ArgType{
			assembler.ArgUint32,
			assembler.ArgUint32,
		},
	}

	tests := []struct {
		Name             string
		Definition       assembler.Definition
		Line             []rune
		ExpectedByteCode op.ByteCode
		ExpectedError    error
	}{
		{"no-args-nil", defNoArgs, nil, b(0), nil},
		{"no-args-empty", defNoArgs, r(""), b(0), nil},
		{"no-args-with-1", defNoArgs, r("1"), b(0), assembler.ErrInvalidArgCount},
		{"one-arg-with-0", defOneArg, nil, b(1, 0), assembler.ErrInvalidArgCount},
		{"one-arg-with-1", defOneArg, r("42"), b(1, 42), nil},
		{"one-arg-with-2", defOneArg, r("41 43"), b(1, 41), assembler.ErrInvalidArgCount},
		{"two-args-with-0", defTwoArgs, nil, b(2, 0, 0), assembler.ErrInvalidArgCount},
		{"two-args-with-1", defTwoArgs, r("40"), b(2, 40, 0), assembler.ErrInvalidArgCount},
		{"two-args-with-2", defTwoArgs, r("41 42"), b(2, 41, 42), nil},
		{"two-args-with-3", defTwoArgs, r("43 44 45"), b(2, 43, 44), assembler.ErrInvalidArgCount},
	}

	for _, test := range tests {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			code, err := test.Definition.Parse(nil, test.Line)
			assert.Equal(t, test.ExpectedByteCode, code)
			if test.ExpectedError != nil {
				assert.ErrorIs(t, err, test.ExpectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAssemble(t *testing.T) {
	t.Parallel()

	b := func(values ...op.Op) op.ByteCode {
		return values
	}
	ae := func(values ...assembler.AssembleError) []assembler.AssembleError {
		return values
	}
	a := func(lineno int, line, msg string) assembler.AssembleError {
		return assembler.AssembleError{lineno, line, msg}
	}

	tests := []struct {
		Name             string
		Lines            string
		ExpectedByteCode op.ByteCode
		Errors           []assembler.AssembleError
	}{
		{"empty", "", nil, nil},
		{"comments-only", "; something\n; another comment\n  ; comment after spaces\n", nil, nil},
		{"no-op", "noop\nnoop ; second no-op\nnoop", b(op.Noop, op.Noop, op.Noop), nil},
		{"invalid-opcode", "noop\noops\nhalt", b(op.Noop, op.Halt), ae(a(2, "oops", "invalid opcode \"oops\""))},
		{
			"opcode-error",
			"Noop 12\nPushI32 abc\nhalt\n",
			b(op.Noop, op.PushInt32, 0, op.Halt),
			ae(
				a(1, "Noop 12", "incorrect number of arguments: Noop takes no arguments"),
				a(2, "PushI32 abc", "invalid argument value: invalid integer value: strconv.ParseInt: parsing \"abc\": invalid syntax"),
			),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			var errors []assembler.AssembleError
			r := func(err assembler.AssembleError) {
				errors = append(errors, err)
			}

			lines := strings.Split(test.Lines, "\n")
			code := assembler.Assemble(lines, r)
			assert.Equal(t, test.ExpectedByteCode, code)
			assert.Equal(t, test.Errors, errors)
		})
	}
}
