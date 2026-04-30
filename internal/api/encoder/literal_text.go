package encoder

import (
	"strings"
	"unicode/utf8"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/scanner"
)

const (
	surr1    = 0xd800
	surr2    = 0xdc00
	surr3    = 0xe000
	surrSelf = 0x10000
)

func encodeLiteralTextForJS(text string, node *ast.Node, strs *stringTable) string {
	raw, ok := rawQuotedLiteralText(node, strs)
	if !ok {
		return text
	}
	decoded, hasSurrogate, ok := decodeQuotedLiteralText(raw)
	if !ok || !hasSurrogate {
		return text
	}
	return decoded
}

func rawQuotedLiteralText(node *ast.Node, strs *stringTable) (string, bool) {
	if node.End() <= 0 || node.End() > len(strs.fileText) {
		return "", false
	}
	start := scanner.SkipTrivia(strs.fileText, node.Pos())
	if start >= node.End() {
		return "", false
	}
	switch strs.fileText[start] {
	case '\'', '"', '`':
		if node.End()-start < 2 || strs.fileText[node.End()-1] != strs.fileText[start] {
			return "", false
		}
		return strs.fileText[start:node.End()], true
	default:
		return "", false
	}
}

func decodeQuotedLiteralText(raw string) (text string, hasSurrogate bool, ok bool) {
	if len(raw) < 2 {
		return "", false, false
	}
	var out strings.Builder
	for i := 1; i < len(raw)-1; {
		if raw[i] != '\\' {
			out.WriteByte(raw[i])
			i++
			continue
		}
		ch, next, ok := decodeEscape(raw, i, len(raw)-1)
		if !ok {
			return "", false, false
		}
		if codePointIsHighSurrogate(ch) {
			hasSurrogate = true
			if nextCh, nextNext, ok := decodeUnicodeEscape(raw, next, len(raw)-1); ok && codePointIsLowSurrogate(nextCh) {
				out.WriteRune(surrogatePairToCodepoint(ch, nextCh))
				i = nextNext
				continue
			}
		} else if codePointIsLowSurrogate(ch) {
			hasSurrogate = true
		}
		out.WriteString(encodeCodePointForJS(ch))
		i = next
	}
	return out.String(), hasSurrogate, true
}

func decodeEscape(raw string, start int, end int) (rune, int, bool) {
	if start+1 >= end {
		return 0, 0, false
	}
	switch raw[start+1] {
	case '0':
		if start+2 >= end || !isDigit(raw[start+2]) {
			return 0, start + 2, true
		}
		return decodeOctalEscape(raw, start, end, 3)
	case '1', '2', '3':
		return decodeOctalEscape(raw, start, end, 3)
	case '4', '5', '6', '7':
		return decodeOctalEscape(raw, start, end, 2)
	case 'u':
		return decodeUnicodeEscape(raw, start, end)
	case 'x':
		if start+4 > end {
			return 0, 0, false
		}
		hi, ok := hexValue(raw[start+2])
		if !ok {
			return 0, 0, false
		}
		lo, ok := hexValue(raw[start+3])
		if !ok {
			return 0, 0, false
		}
		return rune(hi<<4 | lo), start + 4, true
	case 'b':
		return '\b', start + 2, true
	case 't':
		return '\t', start + 2, true
	case 'n':
		return '\n', start + 2, true
	case 'v':
		return '\v', start + 2, true
	case 'f':
		return '\f', start + 2, true
	case 'r':
		return '\r', start + 2, true
	case '\r':
		next := start + 2
		if next < end && raw[next] == '\n' {
			next++
		}
		return -1, next, true
	case '\n':
		return -1, start + 2, true
	default:
		ch, size := utf8.DecodeRuneInString(raw[start+1 : end])
		if ch == utf8.RuneError && size == 0 {
			return 0, 0, false
		}
		return ch, start + 1 + size, true
	}
}

func decodeOctalEscape(raw string, start int, end int, maxDigits int) (rune, int, bool) {
	next := start + 2
	for digits := 1; digits < maxDigits && next < end && isOctalDigit(raw[next]); digits++ {
		next++
	}
	return parseOctalEscape(raw[start+1 : next]), next, true
}

func parseOctalEscape(text string) rune {
	value := rune(0)
	for i := range len(text) {
		value = value*8 + rune(text[i]-'0')
	}
	return value
}

func isDigit(b byte) bool {
	return '0' <= b && b <= '9'
}

func isOctalDigit(b byte) bool {
	return '0' <= b && b <= '7'
}

func decodeUnicodeEscape(raw string, start int, end int) (rune, int, bool) {
	if start+1 >= end || raw[start] != '\\' || raw[start+1] != 'u' {
		return 0, 0, false
	}
	if start+2 < end && raw[start+2] == '{' {
		value := 0
		i := start + 3
		for ; i < end && raw[i] != '}'; i++ {
			digit, ok := hexValue(raw[i])
			if !ok {
				return 0, 0, false
			}
			value = value*16 + digit
		}
		if i >= end || raw[i] != '}' || value > 0x10FFFF {
			return 0, 0, false
		}
		return rune(value), i + 1, true
	}
	if start+6 > end {
		return 0, 0, false
	}
	value := 0
	for i := start + 2; i < start+6; i++ {
		digit, ok := hexValue(raw[i])
		if !ok {
			return 0, 0, false
		}
		value = value*16 + digit
	}
	return rune(value), start + 6, true
}

func hexValue(b byte) (int, bool) {
	switch {
	case '0' <= b && b <= '9':
		return int(b - '0'), true
	case 'a' <= b && b <= 'f':
		return int(b-'a') + 10, true
	case 'A' <= b && b <= 'F':
		return int(b-'A') + 10, true
	default:
		return 0, false
	}
}

func codePointIsHighSurrogate(r rune) bool {
	return surr1 <= r && r < surr2
}

func codePointIsLowSurrogate(r rune) bool {
	return surr2 <= r && r < surr3
}

func surrogatePairToCodepoint(r1, r2 rune) rune {
	return (r1-surr1)<<10 | (r2 - surr2) + surrSelf
}

func encodeCodePointForJS(r rune) string {
	if r < 0 {
		return ""
	}
	if codePointIsHighSurrogate(r) || codePointIsLowSurrogate(r) {
		return string([]byte{
			0xed,
			byte(0x80 | ((r >> 6) & 0x3f)),
			byte(0x80 | (r & 0x3f)),
		})
	}
	return string(r)
}
