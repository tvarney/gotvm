package assembler

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func CutSpace(line []rune) (string, []rune, bool) {
	if len(line) == 0 {
		return "", nil, false
	}

	firstIdx := -1
	for idx, r := range line {
		if unicode.IsSpace(r) {
			firstIdx = idx
			break
		}
	}

	if firstIdx == -1 {
		return string(line), nil, false
	}

	for idx := firstIdx + 1; idx < len(line); idx++ {
		if !unicode.IsSpace(line[idx]) {
			return string(line[:firstIdx]), line[idx:], true
		}
	}

	return string(line[:firstIdx]), nil, true
}

// ParseInt parses an integer with several possible base prefixes.
func ParseInt(value string, bitsize int) (int64, error) {
	if value == "" {
		return 0, fmt.Errorf("%w: integer value may not be empty", ErrInvalidArgValue)
	}
	if value == "-" || value == "+" {
		return 0, fmt.Errorf("%w: integer value may not be just sign character", ErrInvalidArgValue)
	}

	negative := false
	if value[0] == '+' {
		value = value[1:]
	} else if value[0] == '-' {
		value = value[1:]
		negative = true
	}

	// Special case to handle ambiguity between `0633` and `0` when considering
	// base prefixes
	if value == "0" {
		return 0, nil
	}

	// Parse the base we expect; the order should be longest to shortest in
	// this switch statement
	base := 10
	switch {
	case strings.HasPrefix(value, "0x"):
		base = 16
		value = value[2:]
	case strings.HasPrefix(value, "0o"):
		base = 8
		value = value[2:]
	case strings.HasPrefix(value, "0b"):
		base = 2
		value = value[2:]
	case strings.HasPrefix(value, "o"):
		base = 8
		value = value[1:]
	case strings.HasPrefix(value, "0"):
		base = 8
		value = value[1:]
	case strings.HasPrefix(value, "b"):
		base = 2
		value = value[1:]
	}

	// Validate we didn't _just_ get a base identifier
	if value == "" {
		return 0, fmt.Errorf("%w: integer may not be only base prefix", ErrInvalidArgValue)
	}

	// Parse the value now
	ival, err := strconv.ParseInt(value, base, bitsize)
	if err != nil {
		return 0, fmt.Errorf("%w: invalid integer value: %s", ErrInvalidArgValue, err.Error())
	}
	if negative {
		ival = -ival
	}
	return ival, nil
}

// ParseUint parses an unsigned integer with several possible base prefixes.
func ParseUint(value string, bitsize int) (uint64, error) {
	if value == "" {
		return 0, fmt.Errorf("%w: unsigned integer value may not be empty", ErrInvalidArgValue)
	}

	// Special case to handle ambiguity between `0633` and `0` when considering
	// base prefixes
	if value == "0" {
		return 0, nil
	}

	// Parse the base we expect; the order should be longest to shortest in
	// this switch statement
	base := 10
	switch {
	case strings.HasPrefix(value, "0x"):
		base = 16
		value = value[2:]
	case strings.HasPrefix(value, "0o"):
		base = 8
		value = value[2:]
	case strings.HasPrefix(value, "0b"):
		base = 2
		value = value[2:]
	case strings.HasPrefix(value, "o"):
		base = 8
		value = value[1:]
	case strings.HasPrefix(value, "0"):
		base = 8
		value = value[1:]
	case strings.HasPrefix(value, "b"):
		base = 2
		value = value[1:]
	}

	// Validate we didn't _just_ get a base identifier
	if value == "" {
		return 0, fmt.Errorf("%w: unsigned integer may not be only base prefix", ErrInvalidArgValue)
	}

	// Parse the value now
	uval, err := strconv.ParseUint(value, base, bitsize)
	if err != nil {
		return 0, fmt.Errorf("%w: invalid unsigned integer value", ErrInvalidArgValue)
	}
	return uval, nil
}
