package lsconv_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"gotest.tools/v3/assert"
)

func TestDocumentURIToFileName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		uri      lsproto.DocumentUri
		fileName string
	}{
		{"file:///path/to/file.ts", "/path/to/file.ts"},
		{"file://server/share/file.ts", "//server/share/file.ts"},
		{"file:///d%3A/work/tsgo932/lib/utils.ts", "d:/work/tsgo932/lib/utils.ts"},
		{"file:///D%3A/work/tsgo932/lib/utils.ts", "d:/work/tsgo932/lib/utils.ts"},
		{"file:///d%3A/work/tsgo932/app/%28test%29/comp/comp-test.tsx", "d:/work/tsgo932/app/(test)/comp/comp-test.tsx"},
		{"file:///path/to/file.ts#section", "/path/to/file.ts"},
		{"file:///c:/test/me", "c:/test/me"},
		{"file://shares/files/c%23/p.cs", "//shares/files/c#/p.cs"},
		{"file:///c:/Source/Z%C3%BCrich%20or%20Zurich%20(%CB%88zj%CA%8A%C9%99r%C9%AAk,/Code/resources/app/plugins/c%23/plugin.json", "c:/Source/Zürich or Zurich (ˈzjʊərɪk,/Code/resources/app/plugins/c#/plugin.json"},
		{"file:///c:/test %25/path", "c:/test %/path"},
		// {"file:?q", "/"},
		{"file:///_:/path", "/_:/path"},
		{"file:///users/me/c%23-projects/", "/users/me/c#-projects/"},
		{"file://localhost/c%24/GitDevelopment/express", "//localhost/c$/GitDevelopment/express"},
		{"file:///c%3A/test%20with%20%2525/c%23code", "c:/test with %25/c#code"},

		{"untitled:Untitled-1", "^/untitled/ts-nul-authority/Untitled-1"},
		{"untitled:Untitled-1#fragment", "^/untitled/ts-nul-authority/Untitled-1#fragment"},
		{"untitled:c:/Users/jrieken/Code/abc.txt", "^/untitled/ts-nul-authority/c:/Users/jrieken/Code/abc.txt"},
		{"untitled:C:/Users/jrieken/Code/abc.txt", "^/untitled/ts-nul-authority/C:/Users/jrieken/Code/abc.txt"},
		{"untitled://wsl%2Bubuntu/home/jabaile/work/TypeScript-go/newfile.ts", "^/untitled/wsl%2Bubuntu/home/jabaile/work/TypeScript-go/newfile.ts"},
	}

	for _, test := range tests {
		t.Run(string(test.uri), func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.uri.FileName(), test.fileName)
		})
	}
}

func TestFileNameToDocumentURI(t *testing.T) {
	t.Parallel()

	tests := []struct {
		fileName string
		uri      lsproto.DocumentUri
	}{
		{"/path/to/file.ts", "file:///path/to/file.ts"},
		{"//server/share/file.ts", "file://server/share/file.ts"},
		{"d:/work/tsgo932/lib/utils.ts", "file:///d%3A/work/tsgo932/lib/utils.ts"},
		{"d:/work/tsgo932/lib/utils.ts", "file:///d%3A/work/tsgo932/lib/utils.ts"},
		{"d:/work/tsgo932/app/(test)/comp/comp-test.tsx", "file:///d%3A/work/tsgo932/app/%28test%29/comp/comp-test.tsx"},
		{"/path/to/file.ts", "file:///path/to/file.ts"},
		{"c:/test/me", "file:///c%3A/test/me"},
		{"//shares/files/c#/p.cs", "file://shares/files/c%23/p.cs"},
		{"c:/Source/Zürich or Zurich (ˈzjʊərɪk,/Code/resources/app/plugins/c#/plugin.json", "file:///c%3A/Source/Z%C3%BCrich%20or%20Zurich%20%28%CB%88zj%CA%8A%C9%99r%C9%AAk%2C/Code/resources/app/plugins/c%23/plugin.json"},
		{"c:/test %/path", "file:///c%3A/test%20%25/path"},
		{"/", "file:///"},
		{"/_:/path", "file:///_%3A/path"},
		{"/users/me/c#-projects/", "file:///users/me/c%23-projects/"},
		{"//localhost/c$/GitDevelopment/express", "file://localhost/c%24/GitDevelopment/express"},
		{"c:/test with %25/c#code", "file:///c%3A/test%20with%20%2525/c%23code"},

		{"^/untitled/ts-nul-authority/Untitled-1", "untitled:Untitled-1"},
		{"^/untitled/ts-nul-authority/c:/Users/jrieken/Code/abc.txt", "untitled:c:/Users/jrieken/Code/abc.txt"},
		{"^/untitled/ts-nul-authority///wsl%2Bubuntu/home/jabaile/work/TypeScript-go/newfile.ts", "untitled://wsl%2Bubuntu/home/jabaile/work/TypeScript-go/newfile.ts"},
	}

	for _, test := range tests {
		t.Run(test.fileName, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, lsconv.FileNameToDocumentURI(test.fileName), test.uri)
		})
	}
}

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
