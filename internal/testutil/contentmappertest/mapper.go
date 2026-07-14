// Package contentmappertest provides a small but realistic content mapper used by tests. Unlike the
// verbatim mapper baked into the tsc test harness, this mapper performs a genuine transform: it injects
// synthesized glue, copies the body verbatim, and substitutes compiler-option-driven tokens, producing a
// span map that exercises every span kind (verbatim, atom, synthesized). The same handler can be driven
// in-process over a net.Pipe (via NewSpawner) or out-of-process over stdio (via Serve), so the identical
// mapper code backs both the hermetic compiler-test path and the occasional real-subprocess e2e test.
package contentmappertest

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/contentmapperhost"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ipc"
	"github.com/microsoft/typescript-go/internal/json"
	"github.com/microsoft/typescript-go/internal/spanmap"
)

const (
	// ExecName selects the full transforming mapper (synthesized preamble + verbatim body + option-token
	// atoms). This is the "realistic" mapper used by most content-mapper tests.
	ExecName = "compiler-test-mapper"
	// VerbatimExec selects a mapper that copies its input through unchanged with an identity span map. The
	// input is expected to already be valid TypeScript.
	VerbatimExec = "verbatim-mapper"
	// FailingExec selects a mapper that initializes but fails every transform request, exercising the
	// per-file failure and mapper-disabled paths.
	FailingExec = "failing-mapper"
	// SynthesizingExec selects a mapper that emits synthesized TypeScript with no original counterpart, so
	// compiler diagnostics in it are reported against the generated text.
	SynthesizingExec = "synthesizing-mapper"
	// PackageName is the conventional npm package name for the mapper in test fixtures.
	PackageName = "mapper"
)

// preamble is synthesized glue the mapper injects ahead of the body. It has no counterpart in the
// original source, so it maps to a single synthesized segment anchored at the start of the original.
const preamble = "const __VERSION = \"1.0.0\";\n"

// declaredOptions are the compiler options the mapper asks to receive on each transform. Their values are
// substituted into the body wherever a #{name} token appears.
var declaredOptions = []string{"target", "jsx"}

const (
	// diagnosticSource is the prefix the mapper's own diagnostics render with (e.g. "box1000").
	diagnosticSource = "box"
	// unclosedInterpolationCode is the code reported for an interpolation token with no closing brace.
	unclosedInterpolationCode = 1000
)

// noNotifications provides the no-op HandleNotification shared by every mapper: the content mapper
// protocol is request/response only.
type noNotifications struct{}

func (noNotifications) HandleNotification(ctx context.Context, method string, params json.Value) error {
	return nil
}

// Handler implements the content mapper protocol (see internal/contentmapperhost). It answers the
// initialize handshake and transforms foreign file content into TypeScript.
type Handler struct{ noNotifications }

var _ ipc.Handler = Handler{}

func (Handler) HandleRequest(ctx context.Context, method string, params json.Value) (any, error) {
	switch method {
	case contentmapperhost.MethodInitialize:
		return contentmapperhost.InitializeResult{
			ProtocolVersion: contentmapperhost.ProtocolVersion,
			CompilerOptions: declaredOptions,
		}, nil
	case contentmapperhost.MethodTransform:
		var p contentmapperhost.TransformParams
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		text, mappings, diagnostics, err := transform(p.Content, p.CompilerOptions)
		if err != nil {
			return nil, err
		}
		return contentmapperhost.TransformResult{
			Text:        text,
			ScriptKind:  core.ScriptKindTS,
			Mappings:    mappings,
			Diagnostics: diagnostics,
		}, nil
	default:
		return nil, fmt.Errorf("contentmappertest: unexpected method %q", method)
	}
}

// transform turns foreign content into TypeScript and records how the generated text maps back to the
// original. It emits three kinds of segments: a synthesized preamble, verbatim copies of the body, and
// atom substitutions for #{option} tokens whose value comes from a declared compiler option. An
// interpolation token that is never closed on its line is reported as a mapper syntax error (in original
// coordinates) and substituted with `undefined` so the generated TypeScript stays well-formed.
func transform(content string, options *collections.OrderedMap[string, json.Value]) (string, json.Value, []contentmapperhost.Diagnostic, error) {
	var gen strings.Builder
	var segments []spanmap.Segment
	var diagnostics []contentmapperhost.Diagnostic

	gen.WriteString(preamble)
	segments = append(segments, spanmap.Segment{
		GenStart:  0,
		GenEnd:    core.TextPos(gen.Len()),
		OrigStart: 0,
		OrigEnd:   0,
		Kind:      spanmap.KindSynthesized,
	})

	writeVerbatim := func(from, to int) {
		if to <= from {
			return
		}
		genStart := core.TextPos(gen.Len())
		gen.WriteString(content[from:to])
		segments = append(segments, spanmap.Segment{
			GenStart:  genStart,
			GenEnd:    core.TextPos(gen.Len()),
			OrigStart: core.TextPos(from),
			OrigEnd:   core.TextPos(to),
			Kind:      spanmap.KindVerbatim,
		})
	}

	// writeAtom substitutes generated text for the original span [from, to), recording an atom segment so
	// positions within map back to that original span as a whole.
	writeAtom := func(value string, from, to int) {
		genStart := core.TextPos(gen.Len())
		gen.WriteString(value)
		segments = append(segments, spanmap.Segment{
			GenStart:  genStart,
			GenEnd:    core.TextPos(gen.Len()),
			OrigStart: core.TextPos(from),
			OrigEnd:   core.TextPos(to),
			Kind:      spanmap.KindAtom,
		})
	}

	pos := 0
	for pos < len(content) {
		rel := strings.Index(content[pos:], "#{")
		if rel < 0 {
			writeVerbatim(pos, len(content))
			break
		}
		tokenStart := pos + rel

		// The closing brace must appear on the same line; otherwise the interpolation is unclosed.
		lineEnd := tokenStart + strings.IndexByte(content[tokenStart:], '\n')
		if lineEnd < tokenStart {
			lineEnd = len(content)
		}
		closeRel := strings.IndexByte(content[tokenStart:lineEnd], '}')
		if closeRel < 0 {
			writeVerbatim(pos, tokenStart)
			writeAtom("undefined", tokenStart, lineEnd)
			diagnostics = append(diagnostics, contentmapperhost.Diagnostic{
				MessageText: "Unclosed interpolation.",
				Start:       tokenStart,
				Length:      lineEnd - tokenStart,
				Code:        unclosedInterpolationCode,
				Source:      diagnosticSource,
			})
			pos = lineEnd
			continue
		}
		tokenEnd := tokenStart + closeRel + 1
		name := content[tokenStart+len("#{") : tokenEnd-len("}")]

		writeVerbatim(pos, tokenStart)
		writeAtom(renderOption(options, name), tokenStart, tokenEnd)
		pos = tokenEnd
	}

	mappings, err := spanmap.New(segments).Marshal()
	if err != nil {
		return "", nil, nil, err
	}
	return gen.String(), json.Value(mappings), diagnostics, nil
}

// renderOption renders the value of a declared compiler option as a TypeScript expression, or "undefined"
// when the option was not supplied. The option values arrive as raw JSON (e.g. a numeric enum), which is
// already a valid TypeScript expression.
func renderOption(options *collections.OrderedMap[string, json.Value], name string) string {
	if options != nil {
		if value, ok := options.Get(name); ok && len(value) > 0 {
			return string(value)
		}
	}
	return "undefined"
}

// Serve drives the mapper over the given connection until the connection closes or ctx is cancelled. It
// is used both by the in-process spawner and by an out-of-process mapper binary wired to stdio.
func Serve(ctx context.Context, rwc io.ReadWriteCloser) error {
	return ipc.NewAsyncConn(rwc, Handler{}).Run(ctx)
}

// verbatimHandler copies its input through unchanged with an identity span map. The input is expected to
// already be valid TypeScript.
type verbatimHandler struct{ noNotifications }

func (verbatimHandler) HandleRequest(ctx context.Context, method string, params json.Value) (any, error) {
	switch method {
	case contentmapperhost.MethodInitialize:
		return contentmapperhost.InitializeResult{ProtocolVersion: contentmapperhost.ProtocolVersion}, nil
	case contentmapperhost.MethodTransform:
		var p contentmapperhost.TransformParams
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		mappings, err := spanmap.New([]spanmap.Segment{{
			GenEnd:  core.TextPos(len(p.Content)),
			OrigEnd: core.TextPos(len(p.Content)),
			Kind:    spanmap.KindVerbatim,
		}}).Marshal()
		if err != nil {
			return nil, err
		}
		return contentmapperhost.TransformResult{Text: p.Content, Mappings: json.Value(mappings)}, nil
	default:
		return nil, fmt.Errorf("contentmappertest: unexpected method %q", method)
	}
}

// failingHandler initializes successfully but fails every transform request, exercising the compiler's
// per-file failure handling and the mapper-disabled-after-repeated-failures path.
type failingHandler struct{ noNotifications }

func (failingHandler) HandleRequest(ctx context.Context, method string, params json.Value) (any, error) {
	switch method {
	case contentmapperhost.MethodInitialize:
		return contentmapperhost.InitializeResult{ProtocolVersion: contentmapperhost.ProtocolVersion}, nil
	case contentmapperhost.MethodTransform:
		return nil, errors.New("content mapper failed to transform the file")
	default:
		return nil, fmt.Errorf("contentmappertest: unexpected method %q", method)
	}
}

// synthesizedOutput is TypeScript with no counterpart in any original file. It references undeclared names
// so the compiler reports errors that, being fully synthesized, can only be shown against this text.
const synthesizedOutput = "export const el = jsxRuntime(Widget);\n"

// synthesizingHandler emits synthesizedOutput with a fully synthesized span map, so compiler diagnostics
// in it map to no original location.
type synthesizingHandler struct{ noNotifications }

func (synthesizingHandler) HandleRequest(ctx context.Context, method string, params json.Value) (any, error) {
	switch method {
	case contentmapperhost.MethodInitialize:
		return contentmapperhost.InitializeResult{ProtocolVersion: contentmapperhost.ProtocolVersion}, nil
	case contentmapperhost.MethodTransform:
		mappings, err := spanmap.New([]spanmap.Segment{{
			GenEnd:    core.TextPos(len(synthesizedOutput)),
			OrigStart: 0,
			OrigEnd:   0,
			Kind:      spanmap.KindSynthesized,
		}}).Marshal()
		if err != nil {
			return nil, err
		}
		return contentmapperhost.TransformResult{
			Text:       synthesizedOutput,
			ScriptKind: core.ScriptKindTS,
			Mappings:   json.Value(mappings),
		}, nil
	default:
		return nil, fmt.Errorf("contentmappertest: unexpected method %q", method)
	}
}

// handlerForExec selects the mapper implementation named by a package.json's tsContentMapper.exec command.
func handlerForExec(command []string) (ipc.Handler, error) {
	if len(command) == 0 {
		return nil, errors.New("contentmappertest: empty mapper command")
	}
	switch command[0] {
	case ExecName:
		return Handler{}, nil
	case VerbatimExec:
		return verbatimHandler{}, nil
	case FailingExec:
		return failingHandler{}, nil
	case SynthesizingExec:
		return synthesizingHandler{}, nil
	default:
		return nil, fmt.Errorf("contentmappertest: unknown mapper command %v", command)
	}
}

// NewSpawner returns a contentmapperhost.Spawner that serves the fake mappers in-process. It selects the
// implementation by the exec command a mapper package declares (see the *Exec constants), standing it up
// over a net.Pipe so tests exercise the full IPC stack without spawning a real subprocess.
func NewSpawner() contentmapperhost.Spawner {
	return spawner{}
}

type spawner struct{}

func (spawner) Spawn(command []string, dir string) (io.ReadWriteCloser, error) {
	handler, err := handlerForExec(command)
	if err != nil {
		return nil, err
	}
	client, server := net.Pipe()
	go func() { _ = ipc.NewAsyncConn(server, handler).Run(context.Background()) }()
	return client, nil
}

// PackageJSON returns the contents of a package.json that selects the given mapper via its
// tsContentMapper.exec command (one of the *Exec constants), for use as a node_modules fixture in tests.
func PackageJSON(exec string) string {
	return fmt.Sprintf(`{
	"name": %q,
	"version": "1.0.0",
	"tsContentMapper": { "exec": [%q] }
}`, PackageName, exec)
}
