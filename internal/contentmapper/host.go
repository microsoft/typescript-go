package contentmapper

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/spanmap"
)

// Result is the outcome of transforming a foreign file's content into TypeScript.
type Result struct {
	// Text is the transformed TypeScript source text that is parsed into the program.
	Text string
	// ScriptKind is how Text should be parsed.
	ScriptKind core.ScriptKind
	// Diagnostics are syntax errors in the original content.
	Diagnostics []*ast.Diagnostic
	// Mappings maps positions in Text back to the original content, so that diagnostics the compiler
	// produces against the transformed text can be reported at their original locations. A successful
	// transform must return a non-nil map; an empty map describes fully synthesized output.
	Mappings *spanmap.SpanMap
}

// Request carries the inputs for transforming one foreign file.
type Request struct {
	// FileName is the foreign file being transformed.
	FileName string
	// Content is the foreign file's text.
	Content string
	// CompilerOptions is the project's compiler options. The host marshals and forwards only the subset
	// each mapper declared it depends on; a mapper that declares none receives an empty object.
	CompilerOptions *core.CompilerOptions
}

// Host transforms foreign file content into TypeScript during program construction, by driving the
// configured content mappers. Create one with NewHost; Close tears down every mapper it spawned.
type Host interface {
	// Acquire retains the processes for the given mapper identities until the returned lease is released.
	// Acquiring a mapper does not start its process; processes remain lazy until Transform is called.
	Acquire(mappers []*Mapper) (release func())
	// Transform maps a foreign file's content to TypeScript using the given content mapper.
	//
	// A non-nil error indicates the mapper itself failed to produce a result — for example the
	// host hit a broken pipe, a process crash, or could not deserialize the mapper's response.
	Transform(mapper *Mapper, request Request) (result Result, err error)
	// Close shuts down every mapper process the host spawned.
	Close() error
}
