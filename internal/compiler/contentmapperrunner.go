package compiler

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/spanmap"
)

// ContentMapperResult is the outcome of transforming a foreign file's content into TypeScript.
type ContentMapperResult struct {
	// Text is the transformed TypeScript source text that is parsed into the program.
	Text string
	// ScriptKind is how Text should be parsed.
	ScriptKind core.ScriptKind
	// Diagnostics are syntax errors in the original content.
	Diagnostics []*ast.Diagnostic
	// Mappings maps positions in Text back to the original content, so that diagnostics the compiler
	// produces against the transformed text can be reported at their original locations. A nil map
	// means positions are used as-is.
	Mappings *spanmap.SpanMap
}

// ContentMapperRunner transforms foreign file content into TypeScript during program construction.
type ContentMapperRunner interface {
	// Transform maps a foreign file's content to TypeScript.
	//
	// A non-nil error indicates the mapper itself failed to produce a result — for example the
	// runner/coordinator hit a broken pipe, a process crash, or could not deserialize the mapper's
	// response.
	Transform(fileName string, content string) (result ContentMapperResult, err error)
}
