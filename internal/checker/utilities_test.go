package checker

import (
	"testing"
)

func TestIsValidBigIntString(t *testing.T) {
	tests := []struct {
		input       string
		roundTrip   bool
		expectValid bool
	}{
		// Regular decimal integers
		{"5", false, true},
		{"42", false, true},
		{"0", false, true},
		{"-123", false, true},
		{"9007199254740991", false, true}, // MAX_SAFE_INTEGER

		// Hexadecimal format
		{"0x1A", false, true},
		{"0xF", false, true},
		{"0x0", false, true},
		{"-0xFF", false, true},

		// Binary format
		{"0b101", false, true},
		{"0b0", false, true},
		{"-0b1", false, true},

		// Octal format
		{"0o7", false, true},
		{"0o0", false, true},
		{"-0o77", false, true},

		// Invalid cases
		{"", false, false},          // Empty string
		{"-", false, false},         // Just a minus sign
		{"5.5", false, false},       // Decimal point not allowed
		{"5e10", false, false},      // Scientific notation not allowed
		{"abc", false, false},       // Letters not allowed
		{"0xG", false, false},       // Invalid hex digit
		{"0b2", false, false},       // Invalid binary digit
		{"0o8", false, false},       // Invalid octal digit
		{"5n", false, false},        // Trailing n not allowed in string representation
	}

	for _, test := range tests {
		result := isValidBigIntString(test.input, test.roundTrip)
		if result != test.expectValid {
			t.Errorf("isValidBigIntString(%q, %v) = %v, expected %v", 
				test.input, test.roundTrip, result, test.expectValid)
		}
	}
} 
