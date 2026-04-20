package lsconv_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"gotest.tools/v3/assert"
)

type mockScript struct {
	fileName string
	text     string
}

func (m *mockScript) FileName() string { return m.fileName }
func (m *mockScript) Text() string     { return m.text }

func TestLineAndCharacterToPosition_ClampsStaleLineMap(t *testing.T) {
	t.Parallel()

	// Simulate a stale line map that was computed from a longer text.
	// Old text was "hello\nworld\nextra text here\n" (28 chars)
	// New text is "hello\nwor" (9 chars)
	// Line map from old text: [0, 6, 12, 28]
	staleLineMap := &lsconv.LSPLineMap{
		LineStarts: lsconv.LSPLineStarts{0, 6, 12, 28},
		AsciiOnly:  false,
	}

	script := &mockScript{
		fileName: "test.ts",
		text:     "hello\nwor",
	}

	converters := lsconv.NewConverters(lsproto.PositionEncodingKindUTF16, func(_ string) *lsconv.LSPLineMap {
		return staleLineMap
	})

	// Request position on line 1 (stale lineEnd=12 > textLen=9) — should not panic.
	pos := converters.LineAndCharacterToPosition(script, lsproto.Position{Line: 1, Character: 2})
	assert.Assert(t, pos <= core.TextPos(len(script.Text())), "position %d should not exceed text length %d", pos, len(script.Text()))
}

func TestLineAndCharacterToPosition_ClampsStaleLineMapUTF8(t *testing.T) {
	t.Parallel()

	staleLineMap := &lsconv.LSPLineMap{
		LineStarts: lsconv.LSPLineStarts{0, 6, 12, 28},
		AsciiOnly:  true,
	}

	script := &mockScript{
		fileName: "test.ts",
		text:     "hello\nwor",
	}

	converters := lsconv.NewConverters(lsproto.PositionEncodingKindUTF8, func(_ string) *lsconv.LSPLineMap {
		return staleLineMap
	})

	// Request position on line 1 (stale lineEnd=12 > textLen=9) — should not panic.
	pos := converters.LineAndCharacterToPosition(script, lsproto.Position{Line: 1, Character: 2})
	assert.Assert(t, pos <= core.TextPos(len(script.Text())), "position %d should not exceed text length %d", pos, len(script.Text()))
}

func TestLineAndCharacterToPosition_NonASCIIClamped(t *testing.T) {
	t.Parallel()

	// Text with em-dash (U+2014, 3 bytes in UTF-8, 1 UTF-16 code unit)
	// "ab\u2014\n" = 6 bytes: 'a'(1) + 'b'(1) + '\u2014'(3) + '\n'(1)
	// But stale line map says newline is at byte 10 (beyond actual text length of 6).
	staleLineMap := &lsconv.LSPLineMap{
		LineStarts: lsconv.LSPLineStarts{0, 10},
		AsciiOnly:  false,
	}

	script := &mockScript{
		fileName: "test.ts",
		text:     "ab\u2014\n",
	}

	converters := lsconv.NewConverters(lsproto.PositionEncodingKindUTF16, func(_ string) *lsconv.LSPLineMap {
		return staleLineMap
	})

	// Line 0, char 1 — lineEnd from stale map (10) exceeds text length (6). Should not panic.
	pos := converters.LineAndCharacterToPosition(script, lsproto.Position{Line: 0, Character: 1})
	assert.Assert(t, pos <= core.TextPos(len(script.Text())), "position %d should not exceed text length %d", pos, len(script.Text()))
}

func TestPositionToLineAndCharacter_ClampsStaleLineMap(t *testing.T) {
	t.Parallel()

	// Stale line map from a shorter text.
	// Old text: "ab\ncd" (5 chars), line starts: [0, 3]
	// New text: "ab\ncd\nefghij" (12 chars)
	// Position 10 is on what would be line 2 in the new text, but the stale
	// line map only knows about 2 lines. The binary search will place it on
	// line 1 with start=3, which is fine (start <= position).
	staleLineMap := &lsconv.LSPLineMap{
		LineStarts: lsconv.LSPLineStarts{0, 3},
		AsciiOnly:  false,
	}

	script := &mockScript{
		fileName: "test.ts",
		text:     "ab\ncd\nefghij",
	}

	converters := lsconv.NewConverters(lsproto.PositionEncodingKindUTF16, func(_ string) *lsconv.LSPLineMap {
		return staleLineMap
	})

	// Should not panic even with stale line map.
	result := converters.PositionToLineAndCharacter(script, 10)
	_ = result
}

func TestPositionToLineAndCharacter_ClampsStartToPosition(t *testing.T) {
	t.Parallel()

	// Stale line map from a longer text where line starts exceed current text.
	// Stale line map: [0, 20, 40] (from a 40+ char text)
	// Current text: "short" (5 chars)
	// position=3, clamped to 3 (within text)
	// Binary search for 3 in [0, 20, 40] → between 0 and 20, line=0
	// start=0, which is fine. But let's test line map where start > position.
	staleLineMap := &lsconv.LSPLineMap{
		LineStarts: lsconv.LSPLineStarts{0, 2, 40},
		AsciiOnly:  false,
	}

	script := &mockScript{
		fileName: "test.ts",
		text:     "short",
	}

	converters := lsconv.NewConverters(lsproto.PositionEncodingKindUTF16, func(_ string) *lsconv.LSPLineMap {
		return staleLineMap
	})

	// position=3: clamped to 3, binary search in [0,2,40] finds between 2 and 40,
	// line=1, start=2, start(2) < position(3) — fine.
	result := converters.PositionToLineAndCharacter(script, 3)
	_ = result

	// position=5 (textLen): clamped to 5, binary search in [0,2,40] finds between 2 and 40,
	// line=1, start=2, start(2) < position(5) — fine.
	result = converters.PositionToLineAndCharacter(script, 5)
	_ = result

	// position=100: clamped to 5 (textLen), binary search in [0,2,40] finds between 2 and 40,
	// line=1, start=2, start(2) < position(5) — fine.
	result = converters.PositionToLineAndCharacter(script, 100)
	_ = result
}

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
