package lsconv_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"gotest.tools/v3/assert"
)

type testScript struct {
	name string
	text string
}

func (s *testScript) FileName() string { return s.name }
func (s *testScript) Text() string     { return s.text }

func TestLineAndCharacterToPosition_ValidUTF8(t *testing.T) {
	t.Parallel()

	// Text with em-dash (U+2014, 3 bytes in UTF-8, 1 UTF-16 code unit)
	// "ab—cd\nef"
	//  a(0) b(1) —(2,3,4) c(5) d(6) \n(7) e(8) f(9)
	text := "ab\u2014cd\nef"
	script := &testScript{name: "test.ts", text: text}
	lineMap := lsconv.ComputeLSPLineStarts(text)

	converters := lsconv.NewConverters(lsproto.PositionEncodingKindUTF16, func(_ string) *lsconv.LSPLineMap {
		return lineMap
	})

	// Line 0, char 0 → byte 0 (before 'a')
	pos := converters.LineAndCharacterToPosition(script, lsproto.Position{Line: 0, Character: 0})
	assert.Equal(t, pos, core.TextPos(0))

	// Line 0, char 2 → byte 2 (before em-dash)
	pos = converters.LineAndCharacterToPosition(script, lsproto.Position{Line: 0, Character: 2})
	assert.Equal(t, pos, core.TextPos(2))

	// Line 0, char 3 → byte 5 (after em-dash, before 'c')
	pos = converters.LineAndCharacterToPosition(script, lsproto.Position{Line: 0, Character: 3})
	assert.Equal(t, pos, core.TextPos(5))

	// Line 0, char 5 → byte 7 (after 'd', at newline)
	pos = converters.LineAndCharacterToPosition(script, lsproto.Position{Line: 0, Character: 5})
	assert.Equal(t, pos, core.TextPos(7))

	// Line 1, char 0 → byte 8 (start of 'e')
	pos = converters.LineAndCharacterToPosition(script, lsproto.Position{Line: 1, Character: 0})
	assert.Equal(t, pos, core.TextPos(8))
}

func TestLineAndCharacterToPosition_InvalidUTF8(t *testing.T) {
	t.Parallel()

	// Text with invalid UTF-8 byte 0x80 (continuation byte without start byte).
	// range produces RuneError for it, advancing by 1 byte.
	// The old code used utf8.RuneLen(RuneError)==3, overshooting the byte offset.
	text := "a\x80b\ncd"
	script := &testScript{name: "test.ts", text: text}
	lineMap := lsconv.ComputeLSPLineStarts(text)

	converters := lsconv.NewConverters(lsproto.PositionEncodingKindUTF16, func(_ string) *lsconv.LSPLineMap {
		return lineMap
	})

	// Line 0, char 0 → byte 0 ('a')
	pos := converters.LineAndCharacterToPosition(script, lsproto.Position{Line: 0, Character: 0})
	assert.Equal(t, pos, core.TextPos(0))

	// Line 0, char 1 → byte 1 (the invalid byte 0x80; RuneError = 1 UTF-16 code unit)
	pos = converters.LineAndCharacterToPosition(script, lsproto.Position{Line: 0, Character: 1})
	assert.Equal(t, pos, core.TextPos(1))

	// Line 0, char 2 → byte 2 ('b')
	pos = converters.LineAndCharacterToPosition(script, lsproto.Position{Line: 0, Character: 2})
	assert.Equal(t, pos, core.TextPos(2))

	// Line 0, char 3 → byte 3 (at newline, end of line content)
	pos = converters.LineAndCharacterToPosition(script, lsproto.Position{Line: 0, Character: 3})
	assert.Equal(t, pos, core.TextPos(3))

	// Line 1, char 0 → byte 4 ('c')
	pos = converters.LineAndCharacterToPosition(script, lsproto.Position{Line: 1, Character: 0})
	assert.Equal(t, pos, core.TextPos(4))
}

func TestPositionToLineAndCharacter_ValidUTF8(t *testing.T) {
	t.Parallel()

	// Same text as above: "ab—cd\nef"
	text := "ab\u2014cd\nef"
	script := &testScript{name: "test.ts", text: text}
	lineMap := lsconv.ComputeLSPLineStarts(text)

	converters := lsconv.NewConverters(lsproto.PositionEncodingKindUTF16, func(_ string) *lsconv.LSPLineMap {
		return lineMap
	})

	// Byte 0 → (0, 0)
	lc := converters.PositionToLineAndCharacter(script, 0)
	assert.Equal(t, lc, lsproto.Position{Line: 0, Character: 0})

	// Byte 5 → (0, 3) — after em-dash
	lc = converters.PositionToLineAndCharacter(script, 5)
	assert.Equal(t, lc, lsproto.Position{Line: 0, Character: 3})

	// Byte 8 → (1, 0) — start of second line
	lc = converters.PositionToLineAndCharacter(script, 8)
	assert.Equal(t, lc, lsproto.Position{Line: 1, Character: 0})
}

func TestRoundTrip_ValidUTF8(t *testing.T) {
	t.Parallel()

	text := "ab\u2014cd\nef"
	script := &testScript{name: "test.ts", text: text}
	lineMap := lsconv.ComputeLSPLineStarts(text)

	converters := lsconv.NewConverters(lsproto.PositionEncodingKindUTF16, func(_ string) *lsconv.LSPLineMap {
		return lineMap
	})

	// Round-trip only at valid UTF-8 character boundaries.
	// Positions in the middle of multi-byte chars (e.g. byte 3,4 inside em-dash)
	// snap to the next character, so round-tripping them isn't meaningful.
	validPositions := []core.TextPos{0, 1, 2, 5, 6, 7, 8, 9, core.TextPos(len(text))}
	for _, bytePos := range validPositions {
		lc := converters.PositionToLineAndCharacter(script, bytePos)
		rtPos := converters.LineAndCharacterToPosition(script, lc)
		assert.Equal(t, rtPos, bytePos, "round-trip failed for byte position %d", bytePos)
	}
}

func TestRoundTrip_InvalidUTF8(t *testing.T) {
	t.Parallel()

	// Text with invalid UTF-8 byte
	text := "a\x80b\ncd"
	script := &testScript{name: "test.ts", text: text}
	lineMap := lsconv.ComputeLSPLineStarts(text)

	converters := lsconv.NewConverters(lsproto.PositionEncodingKindUTF16, func(_ string) *lsconv.LSPLineMap {
		return lineMap
	})

	// Round-trip: position → line/char → position
	for bytePos := core.TextPos(0); bytePos <= core.TextPos(len(text)); bytePos++ {
		lc := converters.PositionToLineAndCharacter(script, bytePos)
		rtPos := converters.LineAndCharacterToPosition(script, lc)
		assert.Equal(t, rtPos, bytePos, "round-trip failed for byte position %d", bytePos)
	}
}
