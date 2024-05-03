package assembler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tvarney/gotvm/assembler"
)

func TestCutSpace(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Name              string
		Value             string
		ExpectedValue     string
		ExpectedRemainder []rune
		ExpectedFound     bool
	}{
		{"empty-string", "", "", nil, false},
		{"no-cut", "field1", "field1", nil, false},
		{"simple", "field1 field2", "field1", []rune("field2"), true},
		{"consume-spaces", "field1   field2", "field1", []rune("field2"), true},
		{"tabs", "field1\tfield2", "field1", []rune("field2"), true},
		{"consume-tabs", "field1\t\tfield2", "field1", []rune("field2"), true},
		{"mixed-whitespace", "field1\t \v field2", "field1", []rune("field2"), true},
		{"trailing-whitespace", "field1\t  ", "field1", nil, true},
	}

	for _, test := range tests {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			line := []rune(test.Value)
			result, rest, found := assembler.CutSpace(line)
			assert.Equal(t, test.ExpectedValue, result)
			assert.Equal(t, test.ExpectedRemainder, rest)
			assert.Equal(t, test.ExpectedFound, found)
		})
	}
}

func TestParseInt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Name          string
		Value         string
		BitSize       int
		ExpectedInt   int64
		ExpectedError error
	}{
		{"empty-string", "", 32, 0, assembler.ErrInvalidArgValue},
		{"bad-sign-negative", "-", 32, 0, assembler.ErrInvalidArgValue},
		{"bad-sign-positive", "+", 32, 0, assembler.ErrInvalidArgValue},
		{"zero", "0", 32, 0, nil},
		{"negative-zero", "-0", 32, 0, nil},
		{"positive-zero", "+0", 32, 0, nil},
		{"base-only-0x", "0x", 32, 0, assembler.ErrInvalidArgValue},
		{"base-only-negative-0x", "-0x", 32, 0, assembler.ErrInvalidArgValue},
		{"base-only-positive-0x", "+0x", 32, 0, assembler.ErrInvalidArgValue},
		{"base-only-0o", "0o", 32, 0, assembler.ErrInvalidArgValue},
		{"base-only-negative-0o", "-0o", 32, 0, assembler.ErrInvalidArgValue},
		{"base-only-positive-0o", "+0o", 32, 0, assembler.ErrInvalidArgValue},
		{"base-only-0b", "0b", 32, 0, assembler.ErrInvalidArgValue},
		{"base-only-negative-0b", "-0b", 32, 0, assembler.ErrInvalidArgValue},
		{"base-only-positive-0b", "+0b", 32, 0, assembler.ErrInvalidArgValue},
		{"base-only-o", "o", 32, 0, assembler.ErrInvalidArgValue},
		{"base-only-negative-o", "-o", 32, 0, assembler.ErrInvalidArgValue},
		{"base-only-positive-o", "+o", 32, 0, assembler.ErrInvalidArgValue},
		{"base-only-b", "b", 32, 0, assembler.ErrInvalidArgValue},
		{"base-only-negative-b", "-b", 32, 0, assembler.ErrInvalidArgValue},
		{"base-only-positive-b", "+b", 32, 0, assembler.ErrInvalidArgValue},
		{"base16-valid", "0x10", 32, 0x10, nil},
		{"base16-valid-negative", "-0x10", 32, -0x10, nil},
		{"base16-valid-positive", "+0x10", 32, 0x10, nil},
		{"base16-invalid", "0xGOODBYE", 32, 0, assembler.ErrInvalidArgValue},
		{"base10-valid", "10", 32, 10, nil},
		{"base10-valid-negative", "-10", 32, -10, nil},
		{"base10-valid-positive", "+10", 32, 10, nil},
		{"base10-invalid", "10FF", 32, 0, assembler.ErrInvalidArgValue},
		{"base8-valid", "0o10", 32, 8, nil},
		{"base8-valid-negative", "-0o10", 32, -8, nil},
		{"base8-valid-positive", "+0o10", 32, 8, nil},
		{"base8-invalid", "0o79", 32, 0, assembler.ErrInvalidArgValue},
		{"base8-short-valid", "o10", 32, 8, nil},
		{"base8-short-valid-negative", "-o10", 32, -8, nil},
		{"base8-short-valid-positive", "+o10", 32, 8, nil},
		{"base8-short-invalid", "o79", 32, 0, assembler.ErrInvalidArgValue},
		{"base8-prefix-zero-valid", "010", 32, 8, nil},
		{"base8-prefix-zero-valid-negative", "-010", 32, -8, nil},
		{"base8-prefix-zero-valid-positive", "+010", 32, 8, nil},
		{"base8-prefix-zero-invalid", "079", 32, 0, assembler.ErrInvalidArgValue},
		{"base2-valid", "0b111", 32, 7, nil},
		{"base2-valid-negative", "-0b101", 32, -5, nil},
		{"base2-valid-positive", "+0b101", 32, 5, nil},
		{"base2-invalid", "0b123", 32, 0, assembler.ErrInvalidArgValue},
		{"base2-short-valid", "b111", 32, 7, nil},
		{"base2-short-valid-negative", "-b101", 32, -5, nil},
		{"base2-short-valid-positive", "+b101", 32, 5, nil},
		{"base2-short-invalid", "b123", 32, 0, assembler.ErrInvalidArgValue},
	}

	for _, test := range tests {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			value, err := assembler.ParseInt(test.Value, test.BitSize)
			assert.Equal(t, test.ExpectedInt, value)
			if test.ExpectedError != nil {
				assert.ErrorIs(t, err, test.ExpectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseUint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Name          string
		Value         string
		BitSize       int
		ExpectedInt   uint64
		ExpectedError error
	}{
		{"empty-string", "", 32, 0, assembler.ErrInvalidArgValue},
		{"zero", "0", 32, 0, nil},
		{"base-only-0x", "0x", 32, 0, assembler.ErrInvalidArgValue},
		{"base-only-0o", "0o", 32, 0, assembler.ErrInvalidArgValue},
		{"base-only-0b", "0b", 32, 0, assembler.ErrInvalidArgValue},
		{"base-only-o", "o", 32, 0, assembler.ErrInvalidArgValue},
		{"base-only-b", "b", 32, 0, assembler.ErrInvalidArgValue},
		{"base16-valid", "0x10", 32, 0x10, nil},
		{"base16-invalid", "0xGOODBYE", 32, 0, assembler.ErrInvalidArgValue},
		{"base10-valid", "10", 32, 10, nil},
		{"base10-invalid", "10FF", 32, 0, assembler.ErrInvalidArgValue},
		{"base8-valid", "0o10", 32, 8, nil},
		{"base8-invalid", "0o79", 32, 0, assembler.ErrInvalidArgValue},
		{"base8-short-valid", "o10", 32, 8, nil},
		{"base8-short-invalid", "o79", 32, 0, assembler.ErrInvalidArgValue},
		{"base8-prefix-zero-valid", "010", 32, 8, nil},
		{"base8-prefix-zero-invalid", "079", 32, 0, assembler.ErrInvalidArgValue},
		{"base2-valid", "0b111", 32, 7, nil},
		{"base2-invalid", "0b123", 32, 0, assembler.ErrInvalidArgValue},
		{"base2-short-valid", "b111", 32, 7, nil},
		{"base2-short-invalid", "b123", 32, 0, assembler.ErrInvalidArgValue},
	}

	for _, test := range tests {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			value, err := assembler.ParseUint(test.Value, test.BitSize)
			assert.Equal(t, test.ExpectedInt, value)
			if test.ExpectedError != nil {
				assert.ErrorIs(t, err, test.ExpectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
