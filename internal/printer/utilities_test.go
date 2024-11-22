package printer

// import (
// 	"fmt"
// 	"testing"

// 	"gotest.tools/v3/assert"
// )

// func TestEscapeString(t *testing.T) {
// 	data := []struct {
// 		s         string
// 		quoteChar QuoteChar
// 		expected  string
// 	}{
// 		{s: "", quoteChar: QuoteCharDoubleQuote, expected: ``},
// 		{s: "abc", quoteChar: QuoteCharDoubleQuote, expected: `abc`},
// 		{s: "ab\"c", quoteChar: QuoteCharDoubleQuote, expected: `ab\"c`},
// 		{s: "ab\tc", quoteChar: QuoteCharDoubleQuote, expected: `ab\tc`},
// 		{s: "ab\nc", quoteChar: QuoteCharDoubleQuote, expected: `ab\nc`},
// 		{s: "ab'c", quoteChar: QuoteCharDoubleQuote, expected: `ab'c`},
// 		{s: "ab'c", quoteChar: QuoteCharSingleQuote, expected: `ab\'c`},
// 		{s: "ab\"c", quoteChar: QuoteCharSingleQuote, expected: `ab"c`},
// 		{s: "ab`c", quoteChar: QuoteCharBacktick, expected: "ab\\`c"},
// 	}
// 	for i, rec := range data {
// 		t.Run(fmt.Sprintf("[%d] escapeString(%q, %v)", i, rec.s, rec.quoteChar), func(t *testing.T) {
// 			t.Parallel()
// 			actual := escapeString(rec.s, rec.quoteChar)
// 			assert.Equal(t, actual, rec.expected)
// 		})
// 	}
// }

// func TestEscapeNonAsciiString(t *testing.T) {
// 	data := []struct {
// 		s         string
// 		quoteChar QuoteChar
// 		expected  string
// 	}{
// 		{s: "", quoteChar: QuoteCharDoubleQuote, expected: ``},
// 		{s: "abc", quoteChar: QuoteCharDoubleQuote, expected: `abc`},
// 		{s: "ab\"c", quoteChar: QuoteCharDoubleQuote, expected: `ab\"c`},
// 		{s: "ab\tc", quoteChar: QuoteCharDoubleQuote, expected: `ab\tc`},
// 		{s: "ab\nc", quoteChar: QuoteCharDoubleQuote, expected: `ab\nc`},
// 		{s: "ab'c", quoteChar: QuoteCharDoubleQuote, expected: `ab'c`},
// 		{s: "ab'c", quoteChar: QuoteCharSingleQuote, expected: `ab\'c`},
// 		{s: "ab\"c", quoteChar: QuoteCharSingleQuote, expected: `ab"c`},
// 		{s: "ab`c", quoteChar: QuoteCharBacktick, expected: "ab\\`c"},
// 		{s: "ab\u008fc", quoteChar: QuoteCharDoubleQuote, expected: `ab\u008Fc`},
// 		{s: "𝟘𝟙", quoteChar: QuoteCharDoubleQuote, expected: `\uD835\uDFD8\uD835\uDFD9`},
// 	}
// 	for i, rec := range data {
// 		t.Run(fmt.Sprintf("[%d] escapeNonAsciiString(%q, %v)", i, rec.s, rec.quoteChar), func(t *testing.T) {
// 			t.Parallel()
// 			actual := escapeNonAsciiString(rec.s, rec.quoteChar)
// 			assert.Equal(t, actual, rec.expected)
// 		})
// 	}
// }

// func TestEscapeJsxAttributeString(t *testing.T) {
// 	data := []struct {
// 		s         string
// 		quoteChar QuoteChar
// 		expected  string
// 	}{
// 		{s: "", quoteChar: QuoteCharDoubleQuote, expected: ""},
// 		{s: "abc", quoteChar: QuoteCharDoubleQuote, expected: "abc"},
// 		{s: "ab\"c", quoteChar: QuoteCharDoubleQuote, expected: "ab&quot;c"},
// 		{s: "ab\tc", quoteChar: QuoteCharDoubleQuote, expected: "ab&#x9;c"},
// 		{s: "ab\nc", quoteChar: QuoteCharDoubleQuote, expected: "ab&#xA;c"},
// 		{s: "ab'c", quoteChar: QuoteCharDoubleQuote, expected: "ab'c"},
// 		{s: "ab'c", quoteChar: QuoteCharSingleQuote, expected: "ab&apos;c"},
// 		{s: "ab\"c", quoteChar: QuoteCharSingleQuote, expected: "ab\"c"},
// 		{s: "ab\u008fc", quoteChar: QuoteCharDoubleQuote, expected: "ab\u008Fc"},
// 		{s: "𝟘𝟙", quoteChar: QuoteCharDoubleQuote, expected: "𝟘𝟙"},
// 	}
// 	for i, rec := range data {
// 		t.Run(fmt.Sprintf("[%d] escapeJsxAttributeString(%q, %v)", i, rec.s, rec.quoteChar), func(t *testing.T) {
// 			t.Parallel()
// 			actual := escapeJsxAttributeString(rec.s, rec.quoteChar)
// 			assert.Equal(t, actual, rec.expected)
// 		})
// 	}
// }
