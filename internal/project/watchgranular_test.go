package project

import (
	"context"
	"io"
	"slices"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project/logging"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

// TestGranularResolutionGlobsAreNonRecursive verifies that granular watch mode
// emits non-recursive `<dir>/*` globs for resolution lookup directories, rather
// than recursive `<dir>/**/*` globs. The latter collapse shallow directories
// (the workspace root, node_modules) back into broad recursive watches, which
// is exactly what granular mode exists to avoid.
func TestGranularResolutionGlobsAreNonRecursive(t *testing.T) {
	t.Parallel()

	const root = "/home/user/proj"
	fs := vfstest.FromMap(map[string]string{
		root + "/package.json":            "{}",
		root + "/src/a.ts":                "",
		root + "/node_modules/foo/foo.ts": "",
	}, true)

	mapper := createResolutionLookupGlobMapper(root, root+"/lib", root, true /*useCaseSensitiveFileNames*/, true /*granularWatches*/)

	paths := &collections.SyncSet[tspath.Path]{}
	for _, p := range []string{
		root + "/package.json",            // directly in workspace root
		root + "/src/a.ts",                // nested workspace dir
		root + "/node_modules/events.ts",  // bare/core-module probe -> node_modules itself
		root + "/node_modules/foo/foo.ts", // nested node_modules dir
	} {
		paths.Add(tspath.Path(p))
	}

	result := mapper(paths, fs)

	assert.Assert(t, len(result.patternsInsideWorkspace) > 0, "expected some watch globs")
	for _, g := range result.patternsInsideWorkspace {
		assert.Assert(t, !strings.Contains(g, "**"), "granular glob should be non-recursive, got %q", g)
		assert.Assert(t, strings.HasSuffix(g, "/*"), "granular glob should target a directory, got %q", g)
	}

	// The workspace root and node_modules directories must appear only as
	// non-recursive directory watches.
	patterns := result.patternsInsideWorkspace
	assert.Assert(t, slices.Contains(patterns, root+"/*"), "expected non-recursive workspace-root watch, got %v", patterns)
	assert.Assert(t, slices.Contains(patterns, root+"/node_modules/*"), "expected non-recursive node_modules watch, got %v", patterns)
}

// TestAutoImportWatchGlobsGranular verifies that the auto-import node_modules
// watcher emits non-recursive `<nm>/*` globs in granular mode and recursive
// `<nm>/**/*` globs in broad mode.
func TestAutoImportWatchGlobsGranular(t *testing.T) {
	t.Parallel()

	dirs := map[tspath.Path]string{
		"/proj/node_modules":     "/proj/node_modules",
		"/proj/pkg/node_modules": "/proj/pkg/node_modules",
	}

	granular := autoImportWatchGlobs(dirs, true)
	assert.DeepEqual(t, granular.patternsInsideWorkspace, []string{
		"/proj/node_modules/*",
		"/proj/pkg/node_modules/*",
	})
	for _, g := range granular.patternsInsideWorkspace {
		assert.Assert(t, !strings.Contains(g, "**"), "granular auto-import glob should be non-recursive, got %q", g)
	}

	broad := autoImportWatchGlobs(dirs, false)
	assert.DeepEqual(t, broad.patternsInsideWorkspace, []string{
		"/proj/node_modules/**/*",
		"/proj/pkg/node_modules/**/*",
	})
}

// TestGranularWildcardIncludeStaysRecursive verifies that a recursive wildcard
// `include` (e.g. `./vs/**/*.ts`) still produces a recursive `<dir>/**/*` watch
// in granular mode. Downgrading it to a non-recursive `<dir>/*` watch would miss
// files added in subdirectories of the wildcard root, so the program would never
// learn about new root files there.
func TestGranularWildcardIncludeStaysRecursive(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	const root = "/home/user/proj"
	files := map[string]any{
		root + "/tsconfig.json":     `{ "include": ["./vs/**/*.ts"] }`,
		root + "/vs/main.ts":        "export const x = 1;",
		root + "/vs/editor/edit.ts": "export const y = 2;",
	}
	fs := bundled.WrapFS(vfstest.FromMap(files, true /*useCaseSensitiveFileNames*/))

	client := &benchClientMock{}
	session := NewSession(&SessionInit{
		BackgroundCtx: context.Background(),
		FS:            fs,
		Client:        client,
		Logger:        logging.NewLogger(io.Discard),
		Options: &SessionOptions{
			CurrentDirectory:   root,
			DefaultLibraryPath: bundled.LibPath(),
			PositionEncoding:   lsproto.PositionEncodingKindUTF8,
			WatchEnabled:       true,
			GranularWatches:    true,
		},
	})

	session.DidOpenFile(context.Background(), lsproto.DocumentUri("file://"+root+"/vs/main.ts"), 1, files[root+"/vs/main.ts"].(string), lsproto.LanguageKindTypeScript)
	session.WaitForBackgroundTasks()

	var patterns []string
	for _, w := range client.watchers() {
		if w.GlobPattern.Pattern != nil {
			patterns = append(patterns, *w.GlobPattern.Pattern)
		}
	}

	wantRecursive := root + "/vs/**/*"
	assert.Assert(t, slices.Contains(patterns, wantRecursive),
		"expected recursive watch %q for recursive wildcard include, got %v", wantRecursive, patterns)
	// The non-recursive form must not be the only coverage of the wildcard root.
	assert.Assert(t, !slices.Contains(patterns, root+"/vs/*") || slices.Contains(patterns, wantRecursive),
		"recursive wildcard include must not be downgraded to a non-recursive watch, got %v", patterns)
}
