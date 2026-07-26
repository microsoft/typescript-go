package parser

import (
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
)

func parseSourceFileScanWork(sourceText string, scriptKind core.ScriptKind) int {
	p := getParser()
	defer putParser(p)
	p.initializeState(ast.SourceFileParseOptions{
		FileName: "/index.ts",
		Path:     "/index.ts",
	}, sourceText, scriptKind)
	work := 0
	p.scanWork = &work
	p.nextToken()
	p.parseSourceFileWorker()
	return work
}

func TestParserScanWorkIsLinear(t *testing.T) {
	t.Parallel()

	tests := map[string]func(int) string{
		"unclosed parentheses": func(n int) string {
			return strings.Repeat("(", n) + "x;"
		},
		"balanced parentheses": func(n int) string {
			return strings.Repeat("(", n) + "x" + strings.Repeat(")", n) + ";"
		},
		"ambiguous parameters": func(n int) string {
			return strings.Repeat("(a,", n) + "x" + strings.Repeat(")", n) + ";"
		},
		"failed type arguments": func(n int) string {
			return "f" + strings.Repeat("<T>", n) + ";"
		},
		"nested type arguments": func(n int) string {
			return "f" + strings.Repeat("<T", n) + strings.Repeat(">", n) + ";"
		},
		"async arrow ambiguity": func(n int) string {
			return "async " + strings.Repeat("x + ", n) + "x;"
		},
		"declaration modifiers": func(n int) string {
			return strings.Repeat("abstract ", n) + "x;"
		},
		"type member modifiers": func(n int) string {
			return "type T = { " + strings.Repeat("readonly ", n) + "x: string }"
		},
		"function type binding pattern": func(n int) string {
			return "type T = (" + strings.Repeat("{x:", n) + "x" + strings.Repeat("}", n) + ") => void;"
		},
		"infer constraint ambiguity": func(n int) string {
			return "type T<X> = X extends infer U extends " + strings.Repeat("(", n) + "string" + strings.Repeat(")", n) + " ? U : never;"
		},
		"nested infer constraints": func(n int) string {
			return "type T<X> = X extends " + strings.Repeat("infer U extends ", n) + "string ? X : never;"
		},
		"trailing jsdoc whitespace": func(n int) string {
			return "/** x" + strings.Repeat("\n * ", n) + "*/\nconst x = 1;"
		},
		"jsdoc child tags": func(n int) string {
			return "/**\n * @typedef {Object} T\n" + strings.Repeat(" * @property {Object} x\n", n) + " */\nconst x = 1;"
		},
	}

	for name, generate := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			previousWork := 0
			for _, size := range []int{64, 128, 256, 512, 1024} {
				sourceText := generate(size)
				work := parseSourceFileScanWork(sourceText, core.ScriptKindTS)
				if work > len(sourceText)*8 {
					t.Fatalf("size %d used %d scan work for %d input bytes", size, work, len(sourceText))
				}
				if previousWork != 0 && work > previousWork*3 {
					t.Fatalf("doubling from size %d grew scan work from %d to %d", size/2, previousWork, work)
				}
				previousWork = work
			}
		})
	}
}

func TestParserRepeatedFragmentScanWork(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		prefix string
		repeat string
		suffix string
	}{
		{name: "abstract", repeat: "abstract ", suffix: "x;"},
		{name: "accessor", repeat: "accessor ", suffix: "x;"},
		{name: "async", repeat: "async ", suffix: "x;"},
		{name: "declare", repeat: "declare ", suffix: "x;"},
		{name: "private", repeat: "private ", suffix: "x;"},
		{name: "protected", repeat: "protected ", suffix: "x;"},
		{name: "public", repeat: "public ", suffix: "x;"},
		{name: "readonly", repeat: "readonly ", suffix: "x;"},
		{name: "static", repeat: "static ", suffix: "x;"},
		{name: "export", repeat: "export ", suffix: "x;"},
		{name: "mixed modifiers", repeat: "abstract async declare readonly ", suffix: "x;"},
		{name: "open paren", repeat: "(", suffix: "x;"},
		{name: "ambiguous paren", repeat: "(a,", suffix: "x;"},
		{name: "failed type arguments", prefix: "f", repeat: "<T>", suffix: ";"},
		{name: "async binary", prefix: "async ", repeat: "x + ", suffix: "x;"},
		{name: "conditional", repeat: "x ? x : ", suffix: "x;"},
		{name: "unary", repeat: "!~+-", suffix: "x;"},
		{name: "new", repeat: "new ", suffix: "x;"},
		{name: "type operator", prefix: "type T = ", repeat: "keyof ", suffix: "string;"},
		{name: "infer constraint", prefix: "type T<X> = X extends infer U extends ", repeat: "keyof ", suffix: "string ? U : never;"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			small := parseSourceFileScanWork(test.prefix+strings.Repeat(test.repeat, 64)+test.suffix, core.ScriptKindTS)
			large := parseSourceFileScanWork(test.prefix+strings.Repeat(test.repeat, 256)+test.suffix, core.ScriptKindTS)
			if large > small*6 {
				t.Errorf("scan work grew superlinearly: size 64=%d, size 256=%d", small, large)
			}
		})
	}
}

func TestJSXParserScanWorkIsLinear(t *testing.T) {
	t.Parallel()

	tests := map[string]func(int) string{
		"attributes": func(n int) string {
			return "const x = <div " + strings.Repeat(`a="" `, n) + "/>;"
		},
		"text": func(n int) string {
			return "const x = <div>" + strings.Repeat("text ", n) + "</div>;"
		},
		"children": func(n int) string {
			return "const x = <div>" + strings.Repeat("<span />", n) + "</div>;"
		},
	}

	for name, generate := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			previousWork := 0
			for _, size := range []int{64, 128, 256, 512, 1024} {
				sourceText := generate(size)
				work := parseSourceFileScanWork(sourceText, core.ScriptKindTSX)
				if work > len(sourceText)*8 {
					t.Fatalf("size %d used %d scan work for %d input bytes", size, work, len(sourceText))
				}
				if previousWork != 0 && work > previousWork*3 {
					t.Fatalf("doubling from size %d grew scan work from %d to %d", size/2, previousWork, work)
				}
				previousWork = work
			}
		})
	}
}
