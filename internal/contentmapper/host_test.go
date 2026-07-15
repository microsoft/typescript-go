package contentmapper_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/microsoft/typescript-go/internal/contentmapper"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ipc"
	"github.com/microsoft/typescript-go/internal/json"
	"github.com/microsoft/typescript-go/internal/spanmap"
	"gotest.tools/v3/assert"
)

// fakeMapper is an in-process mapper that transforms content verbatim and reports one diagnostic.
type fakeMapper struct{}

func (fakeMapper) HandleRequest(ctx context.Context, method string, params json.Value) (any, error) {
	switch method {
	case contentmapper.MethodInitialize:
		return contentmapper.InitializeResult{ProtocolVersion: contentmapper.ProtocolVersion}, nil
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
				Length:      3,
				Code:        9999,
				Source:      "vue",
			}},
		}, nil
	default:
		return nil, fmt.Errorf("unexpected method %s", method)
	}
}

func (fakeMapper) HandleNotification(ctx context.Context, method string, params json.Value) error {
	return nil
}

// fakeSpawner serves each spawn request with an in-process mapper over a net.Pipe, counting spawns so
// tests can assert process consolidation. When handler is nil it serves a fakeMapper.
type fakeSpawner struct {
	spawns  atomic.Int32
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
	return client, nil
}

func TestRunnerTransform(t *testing.T) {
	t.Parallel()
	r := contentmapper.NewHost(t.Context(), &fakeSpawner{})
	defer r.Close()

	mapper := &contentmapper.Mapper{Manifest: contentmapper.Manifest{Name: "vue", Version: "1.0.0", Exec: []string{"vue-mapper"}}}
	result, err := r.Transform(mapper, contentmapper.Request{FileName: "/a.vue", Content: "export const x = 1;"})
	assert.NilError(t, err)
	assert.Equal(t, result.Text, "export const x = 1;")
	assert.Equal(t, result.ScriptKind, core.ScriptKindTS)
	assert.Assert(t, result.Mappings != nil)
	assert.Equal(t, len(result.Diagnostics), 1)
	assert.Equal(t, result.Diagnostics[0].Code(), int32(9999))
}

func TestRunnerConsolidatesByIdentity(t *testing.T) {
	t.Parallel()
	var spawner fakeSpawner
	r := contentmapper.NewHost(t.Context(), &spawner)
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

// recordingMapper captures (as JSON) the options it receives on transform so a test can assert the host
// forwarded only the declared subset, in order.
type recordingMapper struct {
	mu       sync.Mutex
	received string
}

func (m *recordingMapper) HandleRequest(ctx context.Context, method string, params json.Value) (any, error) {
	switch method {
	case contentmapper.MethodInitialize:
		return contentmapper.InitializeResult{ProtocolVersion: contentmapper.ProtocolVersion}, nil
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
	r := contentmapper.NewHost(t.Context(), &fakeSpawner{handler: mapper})
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
}
