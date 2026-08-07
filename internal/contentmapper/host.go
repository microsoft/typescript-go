package contentmapper

import (
	"fmt"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/locale"
	"github.com/microsoft/typescript-go/internal/spanmap"
)

// TransformErrorKind identifies the stage at which a content mapper transform failed.
type TransformErrorKind uint8

const (
	TransformErrorKindUnknown TransformErrorKind = iota
	TransformErrorKindInitialize
	TransformErrorKindProject
	TransformErrorKindCompilerOptions
	TransformErrorKindRequest
	TransformErrorKindResponse
	TransformErrorKindMappings
)

// TransformError reports a failure while preparing, requesting, or decoding a transform.
type TransformError struct {
	Kind TransformErrorKind
	err  error
}

// NewTransformError creates a transform error for the given stage and underlying error.
func NewTransformError(kind TransformErrorKind, err error) *TransformError {
	return &TransformError{Kind: kind, err: err}
}

func (e *TransformError) Error() string {
	return fmt.Sprintf("content mapper transform failed: %v", e.err)
}

func (e *TransformError) Unwrap() error { return e.err }

// ProjectErrorKind identifies why a mapper's openProject response was rejected.
type ProjectErrorKind uint8

const (
	ProjectErrorKindMalformedResponse ProjectErrorKind = iota
	ProjectErrorKindMissingConfigIdentity
	ProjectErrorKindNonAbsoluteWatchedFile
)

// ProjectError reports an invalid mapper openProject response.
type ProjectError struct {
	Kind ProjectErrorKind
}

func (e *ProjectError) Error() string {
	switch e.Kind {
	case ProjectErrorKindMalformedResponse:
		return "content mapper returned a malformed project response"
	case ProjectErrorKindMissingConfigIdentity:
		return "content mapper did not return configIdentity for dynamic configuration"
	case ProjectErrorKindNonAbsoluteWatchedFile:
		return "content mapper returned a non-absolute path in watchedFiles"
	default:
		return "content mapper returned an invalid project response"
	}
}

// InitializeErrorKind identifies why a mapper's initialize response was rejected.
type InitializeErrorKind uint8

const (
	InitializeErrorKindProtocolVersion InitializeErrorKind = iota
	InitializeErrorKindPositionEncoding
	InitializeErrorKindEmptyDiagnosticSource
	InitializeErrorKindReservedDiagnosticSource
)

// InitializeError reports an invalid or unsupported mapper initialize response.
type InitializeError struct {
	Kind             InitializeErrorKind
	ProtocolVersion  int
	PositionEncoding PositionEncoding
	DiagnosticSource string
}

// SupplementalFileCollisionError reports a compiler-assigned supplemental filename that already exists.
type SupplementalFileCollisionError struct {
	FileName string
}

func (e *SupplementalFileCollisionError) Error() string {
	return fmt.Sprintf("content mapper supplemental output file %q already exists", e.FileName)
}

func (e *InitializeError) Error() string {
	switch e.Kind {
	case InitializeErrorKindProtocolVersion:
		return fmt.Sprintf("unsupported protocol version %d (expected %d)", e.ProtocolVersion, ProtocolVersion)
	case InitializeErrorKindPositionEncoding:
		return fmt.Sprintf("unsupported position encoding %q", e.PositionEncoding)
	case InitializeErrorKindEmptyDiagnosticSource:
		return "diagnostic source must not be empty"
	case InitializeErrorKindReservedDiagnosticSource:
		return fmt.Sprintf("diagnostic source %q is reserved by TypeScript", e.DiagnosticSource)
	default:
		return "content mapper initialization failed"
	}
}

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
	// Supplemental contains additional unnamed outputs associated with the canonical result.
	Supplemental []MappedResult
}

// MappedResult is one generated output and its mapping to the original input.
type MappedResult struct {
	Text       string
	ScriptKind core.ScriptKind
	Mappings   *spanmap.SpanMap
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

// ProjectSpec describes the project configuration visible to its content mappers.
type ProjectSpec struct {
	// ConfigFileName is the absolute project configuration file name, or empty for a project without one.
	ConfigFileName string
	// Mappers are the resolved content mapper entries configured for the project.
	Mappers []*Mapper
	// CompilerOptions are the project's effective compiler options.
	CompilerOptions *core.CompilerOptions
}

// Project is the project-scoped view of a Host. It owns mapper configuration handles and provides the
// identities and watch dependencies needed for caching and incremental builds. Dynamic-config mapper projects
// are opened lazily when an identity, watch dependency, or transform is first requested.
type Project interface {
	// Refresh invalidates dynamic mapper configuration so it is reopened on the next operation that needs it.
	Refresh() error
	// Identities returns sorted transform identities for all configured mappers. It returns an error if
	// dynamic project configuration cannot be opened or validated.
	Identities() ([]string, error)
	// Identity returns the transform identity for mapper, or an empty string if mapper is not in this
	// project. It returns an error if dynamic project configuration cannot be opened or validated.
	Identity(mapper *Mapper) (string, error)
	// WatchedFiles returns the absolute files on which dynamic mapper configuration depends. It returns an
	// error if dynamic project configuration cannot be opened or validated.
	WatchedFiles() ([]string, error)
	// Transform transforms one foreign file using mapper in this project's configuration.
	Transform(mapper *Mapper, request Request) (result Result, err error)
	// Close releases this project reference and closes mapper project handles when no references remain.
	Close() error
}

// Host transforms foreign file content into TypeScript during program construction, by driving the
// configured content mappers. Create one with NewHost; Close tears down every mapper it spawned.
type Host interface {
	// Project returns a retained project-scoped view for spec. Equivalent specs share underlying mapper
	// configuration state; the caller must close the returned Project.
	Project(spec ProjectSpec) Project
	// Acquire retains the processes for the given mapper identities until the returned lease is released.
	// Acquiring a mapper does not start its process; processes remain lazy until Transform is called.
	Acquire(mappers []*Mapper) (release func())
	// SetLocale updates the locale used to initialize mapper processes. Existing processes are stopped
	// and respawned lazily so subsequent transforms use the new locale.
	SetLocale(locale locale.Locale)
	// Transform maps a foreign file's content to TypeScript using the given content mapper.
	//
	// A non-nil error indicates the mapper itself failed to produce a result — for example the
	// host hit a broken pipe, a process crash, or could not deserialize the mapper's response.
	Transform(mapper *Mapper, request Request) (result Result, err error)
	// Close shuts down every mapper process the host spawned.
	Close() error
}
