package contentmapper_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/microsoft/typescript-go/internal/contentmapper"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ipc"
	"github.com/microsoft/typescript-go/internal/json"
	"github.com/microsoft/typescript-go/internal/locale"
	"github.com/microsoft/typescript-go/internal/spanmap"
	"gotest.tools/v3/assert"
)

// fakeMapper is an in-process mapper that transforms content verbatim and reports one diagnostic.
type fakeMapper struct{}

func (fakeMapper) HandleRequest(ctx context.Context, method string, params json.Value) (any, error) {
	switch method {
	case contentmapper.MethodInitialize:
		return contentmapper.InitializeResult{ProtocolVersion: contentmapper.ProtocolVersion, PositionEncoding: contentmapper.PositionEncodingUTF8, DiagnosticSource: "vue"}, nil
	case contentmapper.MethodTransform:
		var p contentmapper.TransformParams
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
		return contentmapper.TransformResult{
			Text:     p.Content,
			Mappings: json.Value(mappings),
			Diagnostics: []contentmapper.Diagnostic{{
				MessageText: "boom",
				Start:       0,
				Length:      min(3, len(p.Content)),
				Code:        9999,
			}},
		}, nil
	default:
		return nil, fmt.Errorf("unexpected method %s", method)
	}
}

type unicodeMapper struct {
	encoding contentmapper.PositionEncoding
	source   *string
}

func (m unicodeMapper) HandleRequest(ctx context.Context, method string, params json.Value) (any, error) {
	switch method {
	case contentmapper.MethodInitialize:
		var p contentmapper.InitializeParams
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		offered := slices.Contains(p.PositionEncodings, m.encoding)
		if !offered && (m.encoding == contentmapper.PositionEncodingUTF8 || m.encoding == contentmapper.PositionEncodingUTF16) {
			return nil, fmt.Errorf("position encoding %q was not offered", m.encoding)
		}
		source := "mapper"
		if m.source != nil {
			source = *m.source
		}
		return contentmapper.InitializeResult{ProtocolVersion: contentmapper.ProtocolVersion, PositionEncoding: m.encoding, DiagnosticSource: source}, nil
	case contentmapper.MethodTransform:
		var p contentmapper.TransformParams
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		var emojiLength, textLength int
		switch m.encoding {
		case contentmapper.PositionEncodingUTF8:
			emojiLength, textLength = 2, 3
		case contentmapper.PositionEncodingUTF16:
			emojiLength, textLength = 1, 2
		default:
			return contentmapper.TransformResult{Text: p.Content}, nil
		}
		mappings, err := json.Marshal([][5]int{
			{0, emojiLength, 0, emojiLength, int(spanmap.KindVerbatim)},
			{emojiLength, textLength - emojiLength, emojiLength, textLength - emojiLength, int(spanmap.KindVerbatim)},
		})
		if err != nil {
			return nil, err
		}
		return contentmapper.TransformResult{
			Text:     p.Content,
			Mappings: mappings,
			Diagnostics: []contentmapper.Diagnostic{{
				MessageText: "after non-ASCII character",
				Start:       emojiLength,
				Length:      textLength - emojiLength,
			}},
		}, nil
	default:
		return nil, fmt.Errorf("unexpected method %s", method)
	}
}

func (unicodeMapper) HandleNotification(ctx context.Context, method string, params json.Value) error {
	return nil
}

type invalidDiagnosticMapper struct {
	encoding contentmapper.PositionEncoding
}

func (m invalidDiagnosticMapper) HandleRequest(ctx context.Context, method string, params json.Value) (any, error) {
	switch method {
	case contentmapper.MethodInitialize:
		return contentmapper.InitializeResult{ProtocolVersion: contentmapper.ProtocolVersion, PositionEncoding: m.encoding, DiagnosticSource: "mapper"}, nil
	case contentmapper.MethodTransform:
		return contentmapper.TransformResult{
			Text: "",
			Diagnostics: []contentmapper.Diagnostic{{
				MessageText: "invalid boundary",
				Start:       1,
			}},
		}, nil
	default:
		return nil, fmt.Errorf("unexpected method %s", method)
	}
}

func (invalidDiagnosticMapper) HandleNotification(ctx context.Context, method string, params json.Value) error {
	return nil
}

func (fakeMapper) HandleNotification(ctx context.Context, method string, params json.Value) error {
	return nil
}

// fakeSpawner serves each spawn request with an in-process mapper over a net.Pipe, counting spawns so
// tests can assert process consolidation. When handler is nil it serves a fakeMapper.
type fakeSpawner struct {
	spawns  atomic.Int32
	closes  atomic.Int32
	handler ipc.Handler
}

func (s *fakeSpawner) Spawn(command []string, dir string) (io.ReadWriteCloser, error) {
	s.spawns.Add(1)
	handler := s.handler
	if handler == nil {
		handler = fakeMapper{}
	}
	client, server := net.Pipe()
	go func() { _ = ipc.NewAsyncConn(server, handler).Run(context.Background()) }()
	return &countingReadWriteCloser{ReadWriteCloser: client, closes: &s.closes}, nil
}

type countingReadWriteCloser struct {
	io.ReadWriteCloser
	closes *atomic.Int32
	once   sync.Once
}

func (c *countingReadWriteCloser) Close() error {
	c.once.Do(func() { c.closes.Add(1) })
	return c.ReadWriteCloser.Close()
}

func TestRunnerTransform(t *testing.T) {
	t.Parallel()
	r := contentmapper.NewHost(t.Context(), &fakeSpawner{}, locale.Default)
	defer r.Close()

	mapper := &contentmapper.Mapper{Manifest: contentmapper.Manifest{Name: "vue", Version: "1.0.0", Exec: []string{"vue-mapper"}}}
	result, err := r.Transform(mapper, contentmapper.Request{FileName: "/a.vue", Content: "export const x = 1;"})
	assert.NilError(t, err)
	assert.Equal(t, result.Text, "export const x = 1;")
	assert.Equal(t, result.ScriptKind, core.ScriptKindTS)
	assert.Assert(t, result.Mappings != nil)
	assert.Equal(t, len(result.Diagnostics), 1)
	assert.Equal(t, result.Diagnostics[0].Code(), int32(9999))
	assert.Equal(t, result.Diagnostics[0].Source(), "vue")
}

func TestRunnerPositionEncodings(t *testing.T) {
	t.Parallel()
	for _, encoding := range []contentmapper.PositionEncoding{
		contentmapper.PositionEncodingUTF8,
		contentmapper.PositionEncodingUTF16,
	} {
		t.Run(string(encoding), func(t *testing.T) {
			t.Parallel()
			r := contentmapper.NewHost(t.Context(), &fakeSpawner{handler: unicodeMapper{encoding: encoding}}, locale.Default)
			defer r.Close()
			mapper := &contentmapper.Mapper{Manifest: contentmapper.Manifest{Name: string(encoding), Exec: []string{"mapper"}}}
			result, err := r.Transform(mapper, contentmapper.Request{FileName: "/a.vue", Content: "éx"})
			assert.NilError(t, err)
			segments := result.Mappings.Segments()
			assert.Equal(t, len(segments), 2)
			assert.Equal(t, int(segments[0].GenEnd), 2)
			assert.Equal(t, int(segments[0].OrigEnd), 2)
			assert.Equal(t, int(segments[1].GenStart), 2)
			assert.Equal(t, int(segments[1].OrigStart), 2)
			assert.Equal(t, result.Text, "éx")
			problem := result.Mappings.Validate(result.Text, "éx")
			assert.Assert(t, problem == nil, "%v", problem)
			mapped, fidelity := result.Mappings.GeneratedToOriginalPosition(2)
			assert.Equal(t, int(mapped), 2)
			assert.Equal(t, fidelity, spanmap.FidelityExact)
			assert.Equal(t, result.Diagnostics[0].Pos(), 2)
			assert.Equal(t, result.Diagnostics[0].End(), 3)
		})
	}
}

func TestRunnerRejectsUnsupportedPositionEncoding(t *testing.T) {
	t.Parallel()
	r := contentmapper.NewHost(t.Context(), &fakeSpawner{handler: unicodeMapper{encoding: "utf-32"}}, locale.Default)
	defer r.Close()
	mapper := &contentmapper.Mapper{Manifest: contentmapper.Manifest{Name: "invalid", Exec: []string{"mapper"}}}
	_, err := r.Transform(mapper, contentmapper.Request{FileName: "/a.vue", Content: "x"})
	assert.ErrorContains(t, err, "unsupported position encoding")
}

func TestRunnerRejectsInvalidDiagnosticSource(t *testing.T) {
	t.Parallel()
	for _, source := range []string{"", " ", "ts", "TS", "d.ts", "json", "typescript", "TypeScript", "tsc", "TSC"} {
		t.Run(source, func(t *testing.T) {
			t.Parallel()
			handler := unicodeMapper{encoding: contentmapper.PositionEncodingUTF8, source: &source}
			r := contentmapper.NewHost(t.Context(), &fakeSpawner{handler: handler}, locale.Default)
			defer r.Close()
			mapper := &contentmapper.Mapper{Manifest: contentmapper.Manifest{Name: "invalid", Exec: []string{"mapper"}}}
			_, err := r.Transform(mapper, contentmapper.Request{FileName: "/a.vue", Content: "x"})
			if strings.TrimSpace(source) == "" {
				assert.ErrorContains(t, err, "diagnostic source must not be empty")
			} else {
				assert.ErrorContains(t, err, "is reserved by TypeScript")
			}
		})
	}
}

func TestRunnerRejectsPositionsInsideUnicodeCharacters(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		encoding contentmapper.PositionEncoding
		content  string
	}{
		{encoding: contentmapper.PositionEncodingUTF8, content: "é"},
		{encoding: contentmapper.PositionEncodingUTF16, content: "😀"},
	} {
		t.Run(string(test.encoding), func(t *testing.T) {
			t.Parallel()
			r := contentmapper.NewHost(t.Context(), &fakeSpawner{handler: invalidDiagnosticMapper{encoding: test.encoding}}, locale.Default)
			defer r.Close()
			mapper := &contentmapper.Mapper{Manifest: contentmapper.Manifest{Name: string(test.encoding), Exec: []string{"mapper"}}}
			_, err := r.Transform(mapper, contentmapper.Request{FileName: "/a.vue", Content: test.content})
			assert.ErrorContains(t, err, "splits a Unicode code point")
		})
	}
}

func TestRunnerConsolidatesByIdentity(t *testing.T) {
	t.Parallel()
	var spawner fakeSpawner
	r := contentmapper.NewHost(t.Context(), &spawner, locale.Default)
	defer r.Close()

	// Two logically-separate mappers with the same identity share one process.
	vueA := &contentmapper.Mapper{Definition: contentmapper.Definition{Package: "a"}, Manifest: contentmapper.Manifest{Name: "vue", Version: "1.0.0", Exec: []string{"vue-mapper"}}}
	vueB := &contentmapper.Mapper{Definition: contentmapper.Definition{Package: "b"}, Manifest: contentmapper.Manifest{Name: "vue", Version: "1.0.0", Exec: []string{"vue-mapper"}}}
	svelte := &contentmapper.Mapper{Manifest: contentmapper.Manifest{Name: "svelte", Version: "2.0.0", Exec: []string{"svelte-mapper"}}}

	for _, m := range []*contentmapper.Mapper{vueA, vueB, vueA, svelte} {
		_, err := r.Transform(m, contentmapper.Request{FileName: "/x", Content: "y"})
		assert.NilError(t, err)
	}
	assert.Equal(t, spawner.spawns.Load(), int32(2), "expected one process per identity")
}

func TestRunnerLeaseLifecycle(t *testing.T) {
	t.Parallel()
	var spawner fakeSpawner
	r := contentmapper.NewHost(t.Context(), &spawner, locale.Default)
	defer r.Close()

	vueA := &contentmapper.Mapper{Definition: contentmapper.Definition{Package: "a"}, Manifest: contentmapper.Manifest{Name: "vue", Version: "1.0.0", Exec: []string{"vue-mapper"}}}
	vueB := &contentmapper.Mapper{Definition: contentmapper.Definition{Package: "b"}, Manifest: contentmapper.Manifest{Name: "vue", Version: "1.0.0", Exec: []string{"vue-mapper"}}}
	svelte := &contentmapper.Mapper{Manifest: contentmapper.Manifest{Name: "svelte", Version: "2.0.0", Exec: []string{"svelte-mapper"}}}

	releaseVueA := r.Acquire([]*contentmapper.Mapper{vueA, vueA})
	releaseVueB := r.Acquire([]*contentmapper.Mapper{vueB})
	releaseSvelte := r.Acquire([]*contentmapper.Mapper{svelte})
	for _, mapper := range []*contentmapper.Mapper{vueA, svelte} {
		_, err := r.Transform(mapper, contentmapper.Request{FileName: "/x", Content: "y"})
		assert.NilError(t, err)
	}
	assert.Equal(t, spawner.spawns.Load(), int32(2))

	releaseVueA()
	assert.Equal(t, spawner.closes.Load(), int32(0), "shared vue process should remain owned")
	releaseSvelte()
	assert.Equal(t, spawner.closes.Load(), int32(1), "final release should close the process")
	releaseVueB()
	releaseVueB()
	assert.Equal(t, spawner.closes.Load(), int32(2), "final vue owner should close once")

	releaseNew := r.Acquire([]*contentmapper.Mapper{vueA})
	_, err := r.Transform(vueA, contentmapper.Request{FileName: "/x", Content: "y"})
	assert.NilError(t, err)
	assert.Equal(t, spawner.spawns.Load(), int32(3), "reacquiring should spawn a fresh process lazily")
	releaseNew()
	assert.Equal(t, spawner.closes.Load(), int32(3))
}

// recordingMapper captures (as JSON) the options it receives on transform so a test can assert the host
// forwarded only the declared subset, in order.
type recordingMapper struct {
	mu             sync.Mutex
	received       string
	receivedLocale string
}

func (m *recordingMapper) HandleRequest(ctx context.Context, method string, params json.Value) (any, error) {
	switch method {
	case contentmapper.MethodInitialize:
		var p contentmapper.InitializeParams
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		m.mu.Lock()
		m.receivedLocale = p.Locale
		m.mu.Unlock()
		return contentmapper.InitializeResult{ProtocolVersion: contentmapper.ProtocolVersion, PositionEncoding: contentmapper.PositionEncodingUTF8, DiagnosticSource: "mapper"}, nil
	case contentmapper.MethodTransform:
		var p contentmapper.TransformParams
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		raw, err := json.Marshal(p.CompilerOptions)
		if err != nil {
			return nil, err
		}
		m.mu.Lock()
		m.received = string(raw)
		m.mu.Unlock()
		return contentmapper.TransformResult{Text: p.Content}, nil
	default:
		return nil, fmt.Errorf("unexpected method %s", method)
	}
}

func (m *recordingMapper) HandleNotification(ctx context.Context, method string, params json.Value) error {
	return nil
}

func TestRunnerForwardsDeclaredOptions(t *testing.T) {
	t.Parallel()
	mapper := &recordingMapper{}
	diagnosticLocale, ok := locale.Parse("cs-CZ")
	assert.Assert(t, ok)
	r := contentmapper.NewHost(t.Context(), &fakeSpawner{handler: mapper}, diagnosticLocale)
	defer r.Close()

	// target is declared and set (forwarded); jsx is declared but unset (omitted); strict is set but
	// undeclared (excluded).
	_, err := r.Transform(
		&contentmapper.Mapper{Manifest: contentmapper.Manifest{Name: "vue", Version: "1.0.0", Exec: []string{"vue-mapper"}, CompilerOptions: []string{"target", "jsx"}}},
		contentmapper.Request{
			FileName:        "/a.vue",
			Content:         "x",
			CompilerOptions: &core.CompilerOptions{Target: core.ScriptTargetES2020, Strict: core.TSTrue},
		},
	)
	assert.NilError(t, err)

	want, err := json.Marshal(core.ScriptTargetES2020)
	assert.NilError(t, err)
	mapper.mu.Lock()
	defer mapper.mu.Unlock()
	assert.Equal(t, mapper.received, fmt.Sprintf(`{"target":%s}`, want))
	assert.Equal(t, mapper.receivedLocale, "cs-CZ")
}
