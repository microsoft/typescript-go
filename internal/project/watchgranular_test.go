package project

import (
	"slices"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/collections"
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
