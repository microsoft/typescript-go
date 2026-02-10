package scanner

import (
	"testing"
)

func TestUtf16Length(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"empty", "", 0},
		{"ascii", "hello", 5},
		{"ascii with spaces", "hello world", 11},
		{"2-byte utf8 (latin)", "café", 4},                          // é is U+00E9, 1 UTF-16 unit
		{"3-byte utf8 (cjk)", "你好", 2},                              // Each is 1 UTF-16 unit
		{"4-byte utf8 (emoji)", "😀", 2},                              // U+1F600, needs surrogate pair (2 UTF-16 units)
		{"mixed ascii and emoji", "a😀b", 4},                          // a=1, 😀=2, b=1
		{"mixed ascii and cjk", "hello世界", 7},                       // 5 ASCII + 2 CJK
		{"supplementary char", "\U0001F4A9", 2},                      // U+1F4A9 pile of poo, surrogate pair
		{"multiple supplementary", "\U0001F600\U0001F601", 4},        // Two emoji, 4 UTF-16 units
		{"ascii + supplementary", "abc\U0001F600def", 8},             // 3 + 2 + 3
		{"tab and newline", "a\tb\n", 4},                             // All ASCII
		{"bmp non-ascii", "\u00e9\u00e8\u00ea", 3},                   // 3 BMP chars, 3 UTF-16 units
		{"mixed bmp and supplementary", "\u00e9\U0001F600\u00e8", 4}, // 1 + 2 + 1
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utf16Length(tt.input)
			if got != tt.expected {
				t.Errorf("utf16Length(%q) = %d, want %d", tt.input, got, tt.expected)
			}
		})
	}
}
