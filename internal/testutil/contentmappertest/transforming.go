package contentmappertest

import (
	"context"
	"fmt"
	"strings"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/contentmapper"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ipc"
	"github.com/microsoft/typescript-go/internal/json"
	"github.com/microsoft/typescript-go/internal/spanmap"
)

const preamble = "const __VERSION = \"1.0.0\";\n"

var DeclaredOptions = []string{"target", "jsx"}

const (
	diagnosticSource          = "box"
	unclosedInterpolationCode = 1000
)

// Handler implements the transforming content mapper protocol.
type Handler struct{ noNotifications }

var _ ipc.Handler = Handler{}

func (Handler) HandleRequest(ctx context.Context, method string, params json.Value) (any, error) {
	switch method {
	case contentmapper.MethodInitialize:
		return initializeResult(diagnosticSource), nil
	case contentmapper.MethodTransform:
		var p contentmapper.TransformParams
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		text, mappings, diagnostics, diagnosticDirectives, err := transform(p.Content, p.CompilerOptions)
		if err != nil {
			return nil, err
		}
		return contentmapper.TransformResult{
			MappedOutput: contentmapper.MappedOutput{Text: text, Mappings: mappings, DiagnosticDirectives: diagnosticDirectives},
			Diagnostics:  diagnostics,
		}, nil
	default:
		return nil, fmt.Errorf("contentmappertest: unexpected method %q", method)
	}
}

func transform(content string, options *collections.OrderedMap[string, json.Value]) (string, json.Value, []contentmapper.Diagnostic, []contentmapper.MappedDiagnosticDirective, error) {
	var virtual strings.Builder
	var segments []spanmap.Segment
	var diagnostics []contentmapper.Diagnostic

	virtual.WriteString(preamble)

	writeVerbatim := func(from, to int) {
		if to <= from {
			return
		}
		virtualStart := core.TextPos(virtual.Len())
		virtual.WriteString(content[from:to])
		segments = append(segments, spanmap.Segment{
			VirtualStart:  virtualStart,
			VirtualEnd:    core.TextPos(virtual.Len()),
			OriginalStart: core.TextPos(from),
			OriginalEnd:   core.TextPos(to),
			Kind:          spanmap.KindVerbatim,
			Features:      spanmap.FeatureAll,
		})
	}

	writeAtom := func(value string, from, to int) {
		virtualStart := core.TextPos(virtual.Len())
		virtual.WriteString(value)
		segments = append(segments, spanmap.Segment{
			VirtualStart:  virtualStart,
			VirtualEnd:    core.TextPos(virtual.Len()),
			OriginalStart: core.TextPos(from),
			OriginalEnd:   core.TextPos(to),
			Kind:          spanmap.KindAtom,
			Features:      spanmap.FeatureAll,
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

	spanMap := spanmap.New(segments)
	mappings, err := spanMap.Marshal()
	if err != nil {
		return "", nil, nil, nil, err
	}
	return virtual.String(), json.Value(mappings), diagnostics, diagnosticDirectives(content, spanMap), nil
}

func diagnosticDirectives(content string, mappings *spanmap.SpanMap) []contentmapper.MappedDiagnosticDirective {
	const invalidPrefix = "// @box-invalid-directive:"
	if strings.HasPrefix(content, invalidPrefix) {
		switch strings.TrimSpace(strings.TrimPrefix(strings.SplitN(content, "\n", 2)[0], invalidPrefix)) {
		case "invalid-range":
			return []contentmapper.MappedDiagnosticDirective{{VirtualStart: -1, Policy: contentmapper.DiagnosticDirectivePolicyIgnore}}
		case "original-range-out-of-bounds":
			return []contentmapper.MappedDiagnosticDirective{{OriginalStart: len(content) + 1, Policy: contentmapper.DiagnosticDirectivePolicyIgnore}}
		case "virtual-range-out-of-bounds":
			return []contentmapper.MappedDiagnosticDirective{{VirtualStart: 1 << 20, Policy: contentmapper.DiagnosticDirectivePolicyIgnore}}
		case "invalid-policy":
			return []contentmapper.MappedDiagnosticDirective{{Policy: "invalid"}}
		case "ignore-with-unused-diagnostic":
			return []contentmapper.MappedDiagnosticDirective{{Policy: contentmapper.DiagnosticDirectivePolicyIgnore, UnusedDiagnostic: &contentmapper.UnusedDirectiveDiagnostic{}}}
		case "expect-without-unused-diagnostic":
			return []contentmapper.MappedDiagnosticDirective{{Policy: contentmapper.DiagnosticDirectivePolicyExpect}}
		case "overlap":
			return []contentmapper.MappedDiagnosticDirective{
				{VirtualLength: 2, Policy: contentmapper.DiagnosticDirectivePolicyIgnore},
				{VirtualStart: 1, VirtualLength: 2, Policy: contentmapper.DiagnosticDirectivePolicyIgnore},
			}
		}
	}
	const ignorePrefix = "// @box-ignore"
	const expectPrefix = "// @box-expect-error"
	var result []contentmapper.MappedDiagnosticDirective
	for lineStart := 0; lineStart < len(content); {
		lineEnd := strings.IndexByte(content[lineStart:], '\n')
		if lineEnd < 0 {
			lineEnd = len(content)
		} else {
			lineEnd += lineStart
		}
		line := content[lineStart:lineEnd]
		trimmed := strings.TrimSpace(line)
		policy := contentmapper.DiagnosticDirectivePolicy("")
		var unused *contentmapper.UnusedDirectiveDiagnostic
		switch {
		case trimmed == ignorePrefix:
			policy = contentmapper.DiagnosticDirectivePolicyIgnore
		case strings.HasPrefix(trimmed, expectPrefix+":"):
			policy = contentmapper.DiagnosticDirectivePolicyExpect
			unused = &contentmapper.UnusedDirectiveDiagnostic{
				Code:        2578,
				MessageText: strings.TrimSpace(strings.TrimPrefix(trimmed, expectPrefix+":")),
			}
		}
		if policy != "" && lineEnd < len(content) {
			affectedStart := lineEnd + 1
			affectedLength := strings.IndexByte(content[affectedStart:], '\n')
			if affectedLength < 0 {
				affectedLength = len(content) - affectedStart
			}
			virtualSpans := mappings.OriginalToVirtualSpans(core.NewTextRange(affectedStart, affectedStart+affectedLength), spanmap.FeatureAll)
			if len(virtualSpans) == 1 {
				result = append(result, contentmapper.MappedDiagnosticDirective{
					OriginalStart:    lineStart,
					OriginalLength:   lineEnd - lineStart,
					VirtualStart:     virtualSpans[0].Span.Pos(),
					VirtualLength:    virtualSpans[0].Span.Len(),
					Policy:           policy,
					UnusedDiagnostic: unused,
				})
			}
		}
		if lineEnd == len(content) {
			break
		}
		lineStart = lineEnd + 1
	}
	return result
}

func renderOption(options *collections.OrderedMap[string, json.Value], name string) string {
	if options != nil {
		if value, ok := options.Get(name); ok && len(value) > 0 {
			return string(value)
		}
	}
	return "undefined"
}
