package assembler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tvarney/gotvm/assembler"
	"github.com/tvarney/gotvm/op"
)

func TestArgType(t *testing.T) {
	t.Parallel()
	t.Run("Parse", testArgTypeParse)
}

func testArgTypeParse(t *testing.T) {
	r := func(s string) []rune {
		return []rune(s)
	}
	b := func(values ...op.Op) op.ByteCode {
		return op.ByteCode(values)
	}

	t.Parallel()
	tests := []struct {
		Name             string
		Arg              assembler.ArgType
		ArgString        []rune
		ExpectedByteCode op.ByteCode
		ExpectedRest     []rune
		ExpectedError    error
	}{
		{"invalid-op-negative", assembler.ArgType(-1), nil, nil, nil, assembler.ErrInvalidArgType},
		{"invalid-op-very-large", assembler.ArgType(1000), nil, nil, nil, assembler.ErrInvalidArgType},

		{"i32-nil", assembler.ArgInt32, nil, b(0), nil, assembler.ErrInvalidArgCount},
		{"i32-empty", assembler.ArgInt32, r(""), b(0), nil, assembler.ErrInvalidArgCount},
		{"i32-invalid-only-arg", assembler.ArgInt32, r("abc"), b(0), nil, assembler.ErrInvalidArgValue},
		{"i32-invalid-multiple-args", assembler.ArgInt32, r("abc 10 12"), b(0), r("10 12"), assembler.ErrInvalidArgValue},
		{"i32-invalid-out-of-range", assembler.ArgInt32, r("0xDEADBEEF"), b(0), nil, assembler.ErrInvalidArgValue},
		{"i32-valid-only-arg", assembler.ArgInt32, r("10"), b(10), nil, nil},
		{"i32-valid-multiple-args", assembler.ArgInt32, r("10 11 12"), b(10), r("11 12"), nil},
		{"i32-valid-negative", assembler.ArgInt32, r("-10"), b(0xFFFFFFF6), nil, nil},

		{"i64-nil", assembler.ArgInt64, nil, b(0, 0), nil, assembler.ErrInvalidArgCount},
		{"i64-empty", assembler.ArgInt64, r(""), b(0, 0), nil, assembler.ErrInvalidArgCount},
		{"i64-invalid-only-arg", assembler.ArgInt64, r("abc"), b(0, 0), nil, assembler.ErrInvalidArgValue},
		{"i64-invalid-multiple-args", assembler.ArgInt64, r("abc 10 12"), b(0, 0), r("10 12"), assembler.ErrInvalidArgValue},
		{"i64-invalid-out-of-range", assembler.ArgInt64, r("0xDEADBEEFCAFED00D"), b(0, 0), nil, assembler.ErrInvalidArgValue},
		{"i64-valid-only-arg", assembler.ArgInt64, r("10"), b(0, 10), nil, nil},
		{"i64-valid-multiple-args", assembler.ArgInt64, r("10 11 12"), b(0, 10), r("11 12"), nil},
		{"i64-valid-large", assembler.ArgInt64, r("0x000FF1CECAFED00D"), b(0x000FF1CE, 0xCAFED00D), nil, nil},
		{"i64-valid-negative", assembler.ArgInt64, r("-10"), b(0xFFFFFFFF, 0xFFFFFFF6), nil, nil},

		{"u32-nil", assembler.ArgUint32, nil, b(0), nil, assembler.ErrInvalidArgCount},
		{"u32-empty", assembler.ArgUint32, r(""), b(0), nil, assembler.ErrInvalidArgCount},
		{"u32-invalid-only-arg", assembler.ArgUint32, r("abc"), b(0), nil, assembler.ErrInvalidArgValue},
		{"u32-invalid-multiple-args", assembler.ArgUint32, r("abc 10 12"), b(0), r("10 12"), assembler.ErrInvalidArgValue},
		{"u32-invalid-out-of-range", assembler.ArgUint32, r("0x1DEADBEEF"), b(0), nil, assembler.ErrInvalidArgValue},
		{"u32-valid-only-arg", assembler.ArgUint32, r("10"), b(10), nil, nil},
		{"u32-valid-multiple-args", assembler.ArgUint32, r("10 11 12"), b(10), r("11 12"), nil},

		{"u64-nil", assembler.ArgUint64, nil, b(0, 0), nil, assembler.ErrInvalidArgCount},
		{"u64-empty", assembler.ArgUint64, r(""), b(0, 0), nil, assembler.ErrInvalidArgCount},
		{"u64-invalid-only-arg", assembler.ArgUint64, r("abc"), b(0, 0), nil, assembler.ErrInvalidArgValue},
		{"u64-invalid-multiple-args", assembler.ArgUint64, r("abc 10 12"), b(0, 0), r("10 12"), assembler.ErrInvalidArgValue},
		{"u64-invalid-out-of-range", assembler.ArgUint64, r("0x1DEADBEEFCAFED00D"), b(0, 0), nil, assembler.ErrInvalidArgValue},
		{"u64-valid-only-arg", assembler.ArgUint64, r("10"), b(0, 10), nil, nil},
		{"u64-valid-multiple-args", assembler.ArgUint64, r("10 11 12"), b(0, 10), r("11 12"), nil},
		{"u64-valid-large", assembler.ArgUint64, r("0xDEADBEEFCAFED00D"), b(0xDEADBEEF, 0xCAFED00D), nil, nil},

		{"f32-nil", assembler.ArgFloat32, nil, b(0), nil, assembler.ErrInvalidArgCount},
		{"f32-empty", assembler.ArgFloat32, r(""), b(0), nil, assembler.ErrInvalidArgCount},
		{"f32-invalid-only-arg", assembler.ArgFloat32, r("abc"), b(0), nil, assembler.ErrInvalidArgValue},
		{"f32-invalid-multiple-args", assembler.ArgFloat32, r("abc 10 12"), b(0), r("10 12"), assembler.ErrInvalidArgValue},
		{"f32-valid-only-arg", assembler.ArgFloat32, r("10.5"), b(0x41280000), nil, nil},
		{"f32-valid-multiple-args", assembler.ArgFloat32, r("10.5 12.0"), b(0x41280000), r("12.0"), nil},

		{"f32-nil", assembler.ArgFloat32, nil, b(0), nil, assembler.ErrInvalidArgCount},
		{"f32-empty", assembler.ArgFloat32, r(""), b(0), nil, assembler.ErrInvalidArgCount},
		{"f32-invalid-only-arg", assembler.ArgFloat32, r("abc"), b(0), nil, assembler.ErrInvalidArgValue},
		{"f32-invalid-multiple-args", assembler.ArgFloat32, r("abc 10 12"), b(0), r("10 12"), assembler.ErrInvalidArgValue},
		{"f32-valid-only-arg", assembler.ArgFloat32, r("10.5"), b(0x41280000), nil, nil},
		{"f32-valid-multiple-args", assembler.ArgFloat32, r("10.5 12.0"), b(0x41280000), r("12.0"), nil},

		{"f64-nil", assembler.ArgFloat64, nil, b(0, 0), nil, assembler.ErrInvalidArgCount},
		{"f64-empty", assembler.ArgFloat64, r(""), b(0, 0), nil, assembler.ErrInvalidArgCount},
		{"f64-invalid-only-arg", assembler.ArgFloat64, r("abc"), b(0, 0), nil, assembler.ErrInvalidArgValue},
		{"f64-invalid-multiple-args", assembler.ArgFloat64, r("abc 10"), b(0, 0), r("10"), assembler.ErrInvalidArgValue},
		{"f64-valid-only-arg", assembler.ArgFloat64, r("9123456789.0"), b(0x4200FE67, 0x38A80000), nil, nil},
		{"f64-valid-multiple-args", assembler.ArgFloat64, r("70123456789.0 1.5"), b(0x423053AF, 0x09150000), r("1.5"), nil},
	}

	for _, test := range tests {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()
			var code op.ByteCode
			rest, result, err := test.Arg.Parse([]rune(test.ArgString), code)
			assert.Equal(t, test.ExpectedByteCode, result, "ByteCode does not match expected value")
			assert.Equal(t, test.ExpectedRest, rest, "Rest does not match expected value")
			if test.ExpectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, test.ExpectedError)
			}
		})
	}
}
