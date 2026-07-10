package compiler

import (
	"slices"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/spanmap"
	"github.com/microsoft/typescript-go/internal/tsoptions"
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

	// Identity returns a stable identifier for the given mapper (e.g. "name@version"), used both to
	// attribute diagnostics and to detect when a mapper implementation has changed between builds. How
	// the identity is resolved — declared in the config, read from a manifest, reported during an
	// initialize handshake, or some combination — is the runner's concern. An empty string means the
	// mapper has no identity to report.
	Identity(mapper *core.ContentMapper) string
}

// ContentMapperIdentities returns the sorted identities that runner assigns to the content mappers in
// config, used to detect when a mapper implementation has changed between builds. Mappers with no
// identity are omitted, and the result is sorted so that merely reordering content mappers in tsconfig
// does not force a rebuild. Returns nil when there is no runner or no mapper reports an identity.
func ContentMapperIdentities(runner ContentMapperRunner, config *tsoptions.ParsedCommandLine) []string {
	if runner == nil {
		return nil
	}
	mappers := config.ContentMappers()
	identities := make([]string, 0, len(mappers))
	for _, mapper := range mappers {
		if identity := runner.Identity(mapper); identity != "" {
			identities = append(identities, identity)
		}
	}
	if len(identities) == 0 {
		return nil
	}
	slices.Sort(identities)
	return identities
}
