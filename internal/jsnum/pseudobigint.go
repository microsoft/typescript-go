package jsnum

import (
	"cmp"
	"fmt"
	"math/big"
	"strings"
)

// PseudoBigInt represents a JS-like bigint. The zero state of the struct represents the value 0.
type PseudoBigInt struct {
	Negative    bool   // true if the value is a non-zero negative number.
	Base10Value string // The absolute value in base 10 with no leading zeros. The value zero is represented as an empty string.
}

func NewPseudoBigInt(value string, negative bool) PseudoBigInt {
	value = strings.TrimLeft(value, "0")
	return PseudoBigInt{Negative: negative && len(value) != 0, Base10Value: value}
}

func (value PseudoBigInt) String() string {
	if len(value.Base10Value) == 0 {
		return "0"
	}
	if value.Negative {
		return "-" + value.Base10Value
	}
	return value.Base10Value
}

func (value PseudoBigInt) Sign() int {
	if len(value.Base10Value) == 0 {
		return 0
	}
	if value.Negative {
		return -1
	}
	return 1
}

// Compare returns -1, 0, or 1 depending on whether value is less than, equal to, or greater than other.
func (value PseudoBigInt) Compare(other PseudoBigInt) int {
	if c := cmp.Compare(value.Sign(), other.Sign()); c != 0 {
		return c
	}
	// Same sign. Compare absolute magnitudes: longer Base10Value means larger magnitude;
	// equal length sorts lexicographically (ASCII digits compare numerically).
	mag := cmp.Compare(len(value.Base10Value), len(other.Base10Value))
	if mag == 0 {
		mag = strings.Compare(value.Base10Value, other.Base10Value)
	}
	if value.Negative {
		return -mag
	}
	return mag
}

func ParseValidBigInt(text string) PseudoBigInt {
	text, negative := strings.CutPrefix(text, "-")
	return NewPseudoBigInt(ParsePseudoBigInt(text), negative)
}

func ParsePseudoBigInt(stringValue string) string {
	stringValue = strings.TrimSuffix(stringValue, "n")
	var b1 byte
	if len(stringValue) > 1 {
		b1 = stringValue[1]
	}
	switch b1 {
	case 'b', 'B', 'o', 'O', 'x', 'X':
		// Not decimal.
	default:
		stringValue = strings.TrimLeft(stringValue, "0")
		if stringValue == "" {
			return "0"
		}
		return stringValue
	}
	bi, ok := new(big.Int).SetString(stringValue, 0)
	if !ok {
		panic(fmt.Sprintf("Failed to parse big int: %q", stringValue))
	}
	return bi.String() // !!!
}
