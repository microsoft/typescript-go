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
	"github.com/microsoft/typescript-go/internal/contentmapper"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ipc"
	"github.com/microsoft/typescript-go/internal/json"
	"github.com/microsoft/typescript-go/internal/spanmap"
)

const (
	// TransformingMapper selects the full transforming mapper (synthesized preamble + verbatim body + option-token
	// atoms). This is the "realistic" mapper used by most content-mapper tests.
	TransformingMapper = "compiler-test-mapper"
	// VerbatimMapper selects a mapper that copies its input through unchanged with an identity span map. The
	// input is expected to already be valid TypeScript.
	VerbatimMapper = "verbatim-mapper"
	// FailingMapper selects a mapper that initializes but fails every transform request, exercising the
	// per-file failure and mapper-disabled paths.
	FailingMapper = "failing-mapper"
	// SynthesizingMapper selects a mapper that emits synthesized TypeScript with no original counterpart, so
	// compiler diagnostics in it are reported against the generated text.
	SynthesizingMapper = "synthesizing-mapper"
	// ComponentMapper selects a Vue-like mapper that extracts <script> contents verbatim, lowers identifiers
	// in {{ template expressions }} as atom mappings, omits markup, and synthesizes component glue.
	ComponentMapper = "component-mapper"
	// DuplicateMapper maps one original identifier to separate semantic and navigation projections.
	DuplicateMapper = "duplicate-mapper"
	// PackageName is the conventional npm package name for the mapper in test fixtures.
	PackageName = "mapper"
)

// preamble is synthesized glue the mapper injects ahead of the body. It has no counterpart in the
// original source, so it maps to a single synthesized segment anchored at the start of the original.
const preamble = "const __VERSION = \"1.0.0\";\n"

// DeclaredOptions are the compiler options the mapper depends on, declared in its package.json manifest.
// Their values are substituted into the body wherever a #{name} token appears.
var DeclaredOptions = []string{"target", "jsx"}

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

type duplicateHandler struct {
	noNotifications
}

func (h duplicateHandler) HandleRequest(ctx context.Context, method string, params json.Value) (any, error) {
	switch method {
	case contentmapper.MethodInitialize:
		return contentmapper.InitializeResult{
			ProtocolVersion:  contentmapper.ProtocolVersion,
			PositionEncoding: contentmapper.PositionEncodingUTF8,
		}, nil
	case contentmapper.MethodTransform:
		var p contentmapper.TransformParams
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		generated := "export const " + p.Content + " = 1;\n" + p.Content + ";\n"
		first := len("export const ")
		second := first + len(p.Content) + len(" = 1;\n")
		disabled := strings.Contains(p.FileName, "disabled")
		semanticPurpose := spanmap.PurposeSemantic
		navigationPurpose := spanmap.PurposeNavigation
		if disabled {
			semanticPurpose = spanmap.PurposeNone
			navigationPurpose = spanmap.PurposeNone
		}
		mappings, err := spanmap.New([]spanmap.Segment{
			{GenStart: core.TextPos(first), GenEnd: core.TextPos(first + len(p.Content)), OrigStart: 0, OrigEnd: core.TextPos(len(p.Content)), Kind: spanmap.KindVerbatim, Purpose: semanticPurpose},
			{GenStart: core.TextPos(second), GenEnd: core.TextPos(second + len(p.Content)), OrigStart: 0, OrigEnd: core.TextPos(len(p.Content)), Kind: spanmap.KindVerbatim, Purpose: navigationPurpose},
		}).Marshal()
		if err != nil {
			return nil, err
		}
		return contentmapper.TransformResult{Text: generated, ScriptKind: core.ScriptKindTS, Mappings: json.Value(mappings)}, nil
	default:
		return nil, fmt.Errorf("contentmappertest: unexpected method %q", method)
	}
}

// Handler implements the content mapper protocol (see internal/contentmapper). It answers the
// initialize handshake and transforms foreign file content into TypeScript.
type Handler struct{ noNotifications }

var _ ipc.Handler = Handler{}

func (Handler) HandleRequest(ctx context.Context, method string, params json.Value) (any, error) {
	switch method {
	case contentmapper.MethodInitialize:
		return contentmapper.InitializeResult{
			ProtocolVersion:  contentmapper.ProtocolVersion,
			PositionEncoding: contentmapper.PositionEncodingUTF8,
		}, nil
	case contentmapper.MethodTransform:
		var p contentmapper.TransformParams
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		text, mappings, diagnostics, err := transform(p.Content, p.CompilerOptions)
		if err != nil {
			return nil, err
		}
		return contentmapper.TransformResult{
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
func transform(content string, options *collections.OrderedMap[string, json.Value]) (string, json.Value, []contentmapper.Diagnostic, error) {
	var gen strings.Builder
	var segments []spanmap.Segment
	var diagnostics []contentmapper.Diagnostic

	// The preamble is synthesized glue with no original counterpart; it is left uncovered (a gap), which
	// the span map treats as synthesized.
	gen.WriteString(preamble)

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
			Purpose:   spanmap.PurposeAll,
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
			Purpose:   spanmap.PurposeAll,
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
			diagnostics = append(diagnostics, contentmapper.Diagnostic{
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
	case contentmapper.MethodInitialize:
		return contentmapper.InitializeResult{ProtocolVersion: contentmapper.ProtocolVersion, PositionEncoding: contentmapper.PositionEncodingUTF8}, nil
	case contentmapper.MethodTransform:
		var p contentmapper.TransformParams
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		mappings, err := spanmap.New([]spanmap.Segment{{
			GenEnd:  core.TextPos(len(p.Content)),
			OrigEnd: core.TextPos(len(p.Content)),
			Kind:    spanmap.KindVerbatim,
			Purpose: spanmap.PurposeAll,
		}}).Marshal()
		if err != nil {
			return nil, err
		}
		return contentmapper.TransformResult{Text: p.Content, Mappings: json.Value(mappings)}, nil
	default:
		return nil, fmt.Errorf("contentmappertest: unexpected method %q", method)
	}
}

// failingHandler initializes successfully but fails every transform request, exercising the compiler's
// per-file failure handling and the mapper-disabled-after-repeated-failures path.
type failingHandler struct{ noNotifications }

func (failingHandler) HandleRequest(ctx context.Context, method string, params json.Value) (any, error) {
	switch method {
	case contentmapper.MethodInitialize:
		return contentmapper.InitializeResult{ProtocolVersion: contentmapper.ProtocolVersion, PositionEncoding: contentmapper.PositionEncodingUTF8}, nil
	case contentmapper.MethodTransform:
		return nil, errors.New("content mapper failed to transform the file")
	default:
		return nil, fmt.Errorf("contentmappertest: unexpected method %q", method)
	}
}

// synthesizedOutput is TypeScript with no counterpart in any original file. It references undeclared names
// so the compiler reports errors that, being fully synthesized, can only be shown against this text.
const synthesizedOutput = "export const el = jsxRuntime(Widget);\n"

// synthesizingHandler emits synthesizedOutput with an empty span map (fully synthesized), so compiler
// diagnostics in it map to no original location.
type synthesizingHandler struct{ noNotifications }

func (synthesizingHandler) HandleRequest(ctx context.Context, method string, params json.Value) (any, error) {
	switch method {
	case contentmapper.MethodInitialize:
		return contentmapper.InitializeResult{ProtocolVersion: contentmapper.ProtocolVersion, PositionEncoding: contentmapper.PositionEncodingUTF8}, nil
	case contentmapper.MethodTransform:
		mappings, err := spanmap.New(nil).Marshal()
		if err != nil {
			return nil, err
		}
		return contentmapper.TransformResult{
			Text:       synthesizedOutput,
			ScriptKind: core.ScriptKindTS,
			Mappings:   json.Value(mappings),
		}, nil
	default:
		return nil, fmt.Errorf("contentmappertest: unexpected method %q", method)
	}
}

type componentHandler struct{ noNotifications }

func (componentHandler) HandleRequest(ctx context.Context, method string, params json.Value) (any, error) {
	switch method {
	case contentmapper.MethodInitialize:
		return contentmapper.InitializeResult{ProtocolVersion: contentmapper.ProtocolVersion, PositionEncoding: contentmapper.PositionEncodingUTF8}, nil
	case contentmapper.MethodTransform:
		var p contentmapper.TransformParams
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		text, mappings, err := transformComponent(p.Content)
		if err != nil {
			return nil, err
		}
		return contentmapper.TransformResult{
			Text:       text,
			ScriptKind: core.ScriptKindTS,
			Mappings:   mappings,
		}, nil
	default:
		return nil, fmt.Errorf("contentmappertest: unexpected method %q", method)
	}
}

func transformComponent(content string) (string, json.Value, error) {
	var gen strings.Builder
	var segments []spanmap.Segment

	writeSynthesized := func(text string) {
		gen.WriteString(text)
	}
	writeMapped := func(text string, origStart, origEnd int, kind spanmap.Kind) {
		genStart := core.TextPos(gen.Len())
		gen.WriteString(text)
		segments = append(segments, spanmap.Segment{
			GenStart:  genStart,
			GenEnd:    core.TextPos(gen.Len()),
			OrigStart: core.TextPos(origStart),
			OrigEnd:   core.TextPos(origEnd),
			Kind:      kind,
			Purpose:   spanmap.PurposeAll,
		})
	}

	scriptOpen := strings.Index(content, "<script")
	if scriptOpen >= 0 {
		openEndRel := strings.IndexByte(content[scriptOpen:], '>')
		if openEndRel < 0 {
			return "", nil, errors.New("contentmappertest: unclosed <script> tag")
		}
		scriptStart := scriptOpen + openEndRel + 1
		closeRel := strings.Index(content[scriptStart:], "</script>")
		if closeRel < 0 {
			return "", nil, errors.New("contentmappertest: missing </script> tag")
		}
		scriptEnd := scriptStart + closeRel
		writeMapped(content[scriptStart:scriptEnd], scriptStart, scriptEnd, spanmap.KindVerbatim)
	}

	writeSynthesized("\nfunction __render() {\n")
	for searchStart := 0; searchStart < len(content); {
		openRel := strings.Index(content[searchStart:], "{{")
		if openRel < 0 {
			break
		}
		exprStart := searchStart + openRel + len("{{")
		closeRel := strings.Index(content[exprStart:], "}}")
		if closeRel < 0 {
			return "", nil, errors.New("contentmappertest: unclosed template expression")
		}
		exprEnd := exprStart + closeRel
		writeSynthesized("  void (")
		for pos := exprStart; pos < exprEnd; {
			if !isIdentifierStart(content[pos]) {
				writeSynthesized(content[pos : pos+1])
				pos++
				continue
			}
			end := pos + 1
			for end < exprEnd && isIdentifierPart(content[end]) {
				end++
			}
			writeMapped(content[pos:end], pos, end, spanmap.KindAtom)
			pos = end
		}
		writeSynthesized(");\n")
		searchStart = exprEnd + len("}}")
	}
	writeSynthesized("}\n")
	if nameStart, nameEnd, ok := componentNameRange(content); ok {
		writeSynthesized("export class ")
		writeMapped(content[nameStart:nameEnd], nameStart, nameEnd, spanmap.KindAtom)
		writeSynthesized(" {}\n")
	}
	writeSynthesized("export default {};\n")

	mappings, err := spanmap.New(segments).Marshal()
	if err != nil {
		return "", nil, err
	}
	return gen.String(), json.Value(mappings), nil
}

func componentNameRange(content string) (start, end int, ok bool) {
	componentStart := strings.Index(content, "<component")
	if componentStart < 0 {
		return 0, 0, false
	}
	tagEndRel := strings.IndexByte(content[componentStart:], '>')
	if tagEndRel < 0 {
		return 0, 0, false
	}
	tag := content[componentStart : componentStart+tagEndRel]
	nameRel := strings.Index(tag, `name="`)
	if nameRel < 0 {
		return 0, 0, false
	}
	start = componentStart + nameRel + len(`name="`)
	endRel := strings.IndexByte(content[start:], '"')
	if endRel < 0 {
		return 0, 0, false
	}
	return start, start + endRel, true
}

func isIdentifierStart(ch byte) bool {
	return ch == '_' || ch == '$' || ch >= 'A' && ch <= 'Z' || ch >= 'a' && ch <= 'z'
}

func isIdentifierPart(ch byte) bool {
	return isIdentifierStart(ch) || ch >= '0' && ch <= '9'
}

// handlerForMapper selects the mapper implementation named by a package.json's tsContentMapper.exec command.
func handlerForMapper(command []string) (ipc.Handler, error) {
	if len(command) == 0 {
		return nil, errors.New("contentmappertest: empty mapper command")
	}
	switch command[0] {
	case TransformingMapper:
		return Handler{}, nil
	case VerbatimMapper:
		return verbatimHandler{}, nil
	case FailingMapper:
		return failingHandler{}, nil
	case SynthesizingMapper:
		return synthesizingHandler{}, nil
	case ComponentMapper:
		return componentHandler{}, nil
	case DuplicateMapper:
		return duplicateHandler{}, nil
	default:
		return nil, fmt.Errorf("contentmappertest: unknown mapper command %v", command)
	}
}

// NewSpawner returns a contentmapper.Spawner that serves the fake mappers in-process. It selects the
// implementation named by the mapper package, standing it up
// over a net.Pipe so tests exercise the full IPC stack without spawning a real subprocess.
func NewSpawner() contentmapper.Spawner {
	return spawner{}
}

type spawner struct{}

func (spawner) Spawn(command []string, dir string) (io.ReadWriteCloser, error) {
	handler, err := handlerForMapper(command)
	if err != nil {
		return nil, err
	}
	client, server := net.Pipe()
	go func() { _ = ipc.NewAsyncConn(server, handler).Run(context.Background()) }()
	return client, nil
}

// PackageJSON returns the contents of a package.json that selects the given mapper via its
// tsContentMapper.exec command, for use as a node_modules fixture in tests. TransformingMapper declares
// the compiler options it depends on; the others declare none.
func PackageJSON(mapper string) string {
	compilerOptions := ""
	if mapper == TransformingMapper {
		compilerOptions = `, "compilerOptions": ["target", "jsx"]`
	}
	return fmt.Sprintf(`{
	"name": %q,
	"version": "1.0.0",
	"tsContentMapper": { "exec": [%q]%s }
}`, PackageName, mapper, compilerOptions)
}
