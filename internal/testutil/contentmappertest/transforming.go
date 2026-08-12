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
		text, mappings, diagnostics, err := transform(p.Content, p.CompilerOptions)
		if err != nil {
			return nil, err
		}
		return contentmapper.TransformResult{
			MappedOutput: contentmapper.MappedOutput{Text: text, Mappings: mappings},
			Diagnostics:  diagnostics,
		}, nil
	default:
		return nil, fmt.Errorf("contentmappertest: unexpected method %q", method)
	}
}

func transform(content string, options *collections.OrderedMap[string, json.Value]) (string, json.Value, []contentmapper.Diagnostic, error) {
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

	mappings, err := spanmap.New(segments).Marshal()
	if err != nil {
		return "", nil, nil, err
	}
	return virtual.String(), json.Value(mappings), diagnostics, nil
}

func renderOption(options *collections.OrderedMap[string, json.Value], name string) string {
	if options != nil {
		if value, ok := options.Get(name); ok && len(value) > 0 {
			return string(value)
		}
	}
	return "undefined"
}
