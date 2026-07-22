package tsctests

import (
	"context"
	"io"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/microsoft/typescript-go/internal/contentmapper"
	"github.com/microsoft/typescript-go/internal/execute"
	"github.com/microsoft/typescript-go/internal/fswatch"
	"github.com/microsoft/typescript-go/internal/testutil/contentmappertest"
	"gotest.tools/v3/assert"
)

type recordingContentMapperSystem struct {
	*TestSys
	spawner *recordingContentMapperSpawner
}

func (s *recordingContentMapperSystem) Spawn(command []string, dir string) (io.ReadWriteCloser, error) {
	return s.spawner.Spawn(command, dir)
}

type recordingContentMapperSpawner struct {
	inner  contentmapper.Spawner
	spawns atomic.Int32
	closes atomic.Int32
}

func (s *recordingContentMapperSpawner) Spawn(command []string, dir string) (io.ReadWriteCloser, error) {
	process, err := s.inner.Spawn(command, dir)
	if err != nil {
		return nil, err
	}
	s.spawns.Add(1)
	return &recordingContentMapperProcess{ReadWriteCloser: process, closes: &s.closes}, nil
}

type recordingContentMapperProcess struct {
	io.ReadWriteCloser
	closes *atomic.Int32
	once   sync.Once
}

func (p *recordingContentMapperProcess) Close() error {
	var err error
	p.once.Do(func() {
		p.closes.Add(1)
		err = p.ReadWriteCloser.Close()
	})
	return err
}

func TestContentMapperWatchLifecycle(t *testing.T) {
	t.Parallel()
	for _, test := range []struct {
		name string
		args []string
	}{
		{name: "watch", args: []string{"--watch", "--loadExternalPlugins"}},
		{name: "build watch", args: []string{"--build", "--watch", "--loadExternalPlugins"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			const configFileName = "/home/src/workspaces/project/tsconfig.json"
			input := &tscInput{files: FileMap{
				configFileName: `{
					"compilerOptions": { "composite": true },
					"contentMappers": [{ "package": "mapper-a", "extensions": [".vue"] }]
				}`,
				"/home/src/workspaces/project/app.vue":                            `export const app = 1;`,
				"/home/src/workspaces/project/node_modules/mapper-a/package.json": contentmappertest.PackageJSON(contentmappertest.VerbatimMapper),
				"/home/src/workspaces/project/node_modules/mapper-b/package.json": strings.Replace(contentmappertest.PackageJSON(contentmappertest.VerbatimMapper), `"version": "1.0.0"`, `"version": "2.0.0"`, 1),
			}}
			testSys := newTestSys(input, false)
			spawner := &recordingContentMapperSpawner{inner: contentmappertest.NewSpawner()}
			sys := &recordingContentMapperSystem{TestSys: testSys, spawner: spawner}
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			result := execute.CommandLine(ctx, sys, test.args, testSys)
			assert.Assert(t, result.Watcher != nil)
			assert.Equal(t, spawner.spawns.Load(), int32(1))
			assert.Equal(t, spawner.closes.Load(), int32(0))

			testSys.writeFileNoError(configFileName, `{
				"compilerOptions": { "composite": true },
				"contentMappers": [{ "package": "mapper-b", "extensions": [".vue"] }]
			}`)
			testSys.mockWatchBackend.SendEvents([]fswatch.Event{{Kind: fswatch.EventUpdate, Path: configFileName}})
			result.Watcher.DoCycle()

			assert.Equal(t, spawner.spawns.Load(), int32(2))
			assert.Equal(t, spawner.closes.Load(), int32(1))

			testSys.writeFileNoError(configFileName, `{ "compilerOptions": { "composite": true } }`)
			testSys.mockWatchBackend.SendEvents([]fswatch.Event{{Kind: fswatch.EventUpdate, Path: configFileName}})
			result.Watcher.DoCycle()

			assert.Equal(t, spawner.closes.Load(), int32(2))
		})
	}
}

func TestContentMapperBuildWatchSharedLifecycle(t *testing.T) {
	t.Parallel()
	const mapperConfig = `{
		"compilerOptions": { "composite": true },
		"contentMappers": [{ "package": "mapper", "extensions": [".vue"] }]
	}`
	input := &tscInput{files: FileMap{
		"/home/src/workspaces/project/tsconfig.json": `{
			"files": [],
			"references": [{ "path": "a" }, { "path": "b" }]
		}`,
		"/home/src/workspaces/project/a/tsconfig.json":                  mapperConfig,
		"/home/src/workspaces/project/a/app.vue":                        `export const a = 1;`,
		"/home/src/workspaces/project/b/tsconfig.json":                  mapperConfig,
		"/home/src/workspaces/project/b/app.vue":                        `export const b = 1;`,
		"/home/src/workspaces/project/node_modules/mapper/package.json": contentmappertest.PackageJSON(contentmappertest.VerbatimMapper),
	}}
	testSys := newTestSys(input, false)
	spawner := &recordingContentMapperSpawner{inner: contentmappertest.NewSpawner()}
	sys := &recordingContentMapperSystem{TestSys: testSys, spawner: spawner}
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	result := execute.CommandLine(ctx, sys, []string{"--build", "--watch", "--loadExternalPlugins"}, testSys)
	assert.Assert(t, result.Watcher != nil)
	assert.Equal(t, spawner.spawns.Load(), int32(1))
	assert.Equal(t, spawner.closes.Load(), int32(0))

	for _, project := range []string{"a", "b"} {
		configFileName := "/home/src/workspaces/project/" + project + "/tsconfig.json"
		testSys.writeFileNoError(configFileName, `{ "compilerOptions": { "composite": true } }`)
		testSys.mockWatchBackend.SendEvents([]fswatch.Event{{Kind: fswatch.EventUpdate, Path: configFileName}})
		result.Watcher.DoCycle()
		if project == "a" {
			assert.Equal(t, spawner.closes.Load(), int32(0))
		} else {
			assert.Equal(t, spawner.closes.Load(), int32(1))
		}
	}
}
