package pnp

import (
	"strings"
	"testing"
)

// A hand-written manifest mirroring the shape Yarn emits: a top-level package, a
// workspace package (SOFT link, on-disk location) that depends on an npm package
// (HARD link, a location inside a cache .zip), plus an aliased dependency and an
// unfulfilled peer.
const fixtureManifest = `{
  "__info": ["generated"],
  "dependencyTreeRoots": [{"name": "root", "reference": "workspace:."}],
  "enableTopLevelFallback": true,
  "ignorePatternData": "(^(?:\\.yarn\\/sdks(?:\\/(?!\\.{1,2}(?:\\/|$))(?:(?:(?!(?:^|\\/)\\.{1,2}(?:\\/|$)).)*?)|$))$)",
  "fallbackExclusionList": [["root", ["workspace:."]]],
  "fallbackPool": [],
  "packageRegistryData": [
    [null, [["workspace:.", {
      "packageLocation": "./",
      "packageDependencies": [["@t/util", "workspace:packages/util"]],
      "linkType": "SOFT"
    }]]],
    ["@t/util", [["workspace:packages/util", {
      "packageLocation": "./packages/util/",
      "packageDependencies": [
        ["@t/util", "workspace:packages/util"],
        ["lodash", "npm:4.17.21"],
        ["aliased", ["real-pkg", "npm:1.0.0"]],
        ["peer-dep", null]
      ],
      "linkType": "SOFT"
    }]]],
    ["lodash", [["npm:4.17.21", {
      "packageLocation": "./.yarn/cache/lodash-npm-4.17.21-abc.zip/node_modules/lodash/",
      "packageDependencies": [["lodash", "npm:4.17.21"]],
      "linkType": "HARD"
    }]]],
    ["real-pkg", [["npm:1.0.0", {
      "packageLocation": "./.yarn/cache/real-pkg-npm-1.0.0-def.zip/node_modules/real-pkg/",
      "packageDependencies": [["real-pkg", "npm:1.0.0"]],
      "linkType": "HARD"
    }]]]
  ]
}`

func loadFixture(t *testing.T) *Manifest {
	t.Helper()
	const manifestPath = "/proj/.pnp.data.json"
	m, ok := Load(manifestPath, func(p string) (string, bool) {
		if p == manifestPath {
			return fixtureManifest, true
		}
		return "", false
	})
	if !ok {
		t.Fatal("failed to load fixture manifest")
	}
	return m
}

func TestResolveWorkspaceDependency(t *testing.T) {
	t.Parallel()
	m := loadFixture(t)
	r := m.Resolve("@t/util", "/proj")
	if r.Status != Success {
		t.Fatalf("status = %d, want Success", r.Status)
	}
	if r.PackageDir != "/proj/packages/util" {
		t.Fatalf("PackageDir = %q", r.PackageDir)
	}
}

func TestResolveNpmDependencyInsideZip(t *testing.T) {
	t.Parallel()
	m := loadFixture(t)
	r := m.Resolve("lodash", "/proj/packages/util")
	if r.Status != Success {
		t.Fatalf("status = %d, want Success", r.Status)
	}
	want := "/proj/.yarn/cache/lodash-npm-4.17.21-abc.zip/node_modules/lodash"
	if r.PackageDir != want {
		t.Fatalf("PackageDir = %q, want %q", r.PackageDir, want)
	}
}

func TestResolveSubpathCarriedThrough(t *testing.T) {
	t.Parallel()
	m := loadFixture(t)
	r := m.Resolve("lodash/fp", "/proj/packages/util")
	if r.Status != Success || r.Subpath != "/fp" {
		t.Fatalf("status = %d subpath = %q", r.Status, r.Subpath)
	}
}

func TestResolveAliasedDependency(t *testing.T) {
	t.Parallel()
	m := loadFixture(t)
	// "aliased" points at ["real-pkg","npm:1.0.0"]; it must resolve to real-pkg's location.
	r := m.Resolve("aliased", "/proj/packages/util")
	if r.Status != Success {
		t.Fatalf("status = %d, want Success", r.Status)
	}
	if !strings.Contains(r.PackageDir, ".zip/") || r.PackageDir != "/proj/.yarn/cache/real-pkg-npm-1.0.0-def.zip/node_modules/real-pkg" {
		t.Fatalf("aliased PackageDir = %q", r.PackageDir)
	}
}

func TestResolveUnfulfilledPeer(t *testing.T) {
	t.Parallel()
	m := loadFixture(t)
	r := m.Resolve("peer-dep", "/proj/packages/util")
	if r.Status != UnfulfilledPeer {
		t.Fatalf("status = %d, want UnfulfilledPeer", r.Status)
	}
}

func TestResolveNotADependency(t *testing.T) {
	t.Parallel()
	m := loadFixture(t)
	// lodash is a dep of @t/util but not of the top-level package, and there is
	// no fallback pool entry, so resolving it from the root fails.
	r := m.Resolve("lodash", "/proj")
	if r.Status == Success {
		t.Fatalf("unexpectedly resolved lodash from root to %q", r.PackageDir)
	}
}

func TestResolveIgnoredImporterSkips(t *testing.T) {
	t.Parallel()
	m := loadFixture(t)
	// A .yarn/sdks path matches ignorePatternData -> Skipped (fall back to classic).
	r := m.Resolve("lodash", "/proj/.yarn/sdks/typescript")
	if r.Status != Skipped {
		t.Fatalf("status = %d, want Skipped", r.Status)
	}
}

func TestResolveGovernedNotFoundDistinctFromSkipped(t *testing.T) {
	t.Parallel()
	m := loadFixture(t)
	// The root workspace IS governed by PnP, but "lodash" is not one of its
	// declared dependencies and there is no fallback pool entry: the status is
	// NotFound, NOT Skipped. The resolver relies on this distinction — a governed
	// NotFound must stop the classic node_modules walk (so a phantom dependency in
	// a stray node_modules is not resolved), whereas Skipped falls through to it.
	if r := m.Resolve("lodash", "/proj"); r.Status != NotFound {
		t.Fatalf("undeclared dep from governed importer: status = %d, want NotFound", r.Status)
	}
	// An unfulfilled peer is also governed (not Skipped).
	if r := m.Resolve("peer-dep", "/proj/packages/util"); r.Status != UnfulfilledPeer {
		t.Fatalf("unfulfilled peer: status = %d, want UnfulfilledPeer", r.Status)
	}
	// A non-absolute importer directory cannot be located against the absolute
	// manifest: Skipped (fall back to the classic walk), never a false match.
	if r := m.Resolve("lodash", "packages/util"); r.Status != Skipped {
		t.Fatalf("relative importer: status = %d, want Skipped", r.Status)
	}
}

func TestLoadInlinedCjsTightSpacing(t *testing.T) {
	t.Parallel()
	// Yarn's spacing around the assignment is not guaranteed; the extractor must
	// tolerate `RAW_RUNTIME_STATE=` (no spaces) as well as `RAW_RUNTIME_STATE =`.
	jsEscaped := strings.ReplaceAll(fixtureManifest, `\`, `\\`)
	jsEscaped = strings.ReplaceAll(jsEscaped, "\n", "\\\n")
	cjs := "const RAW_RUNTIME_STATE='" + jsEscaped + "';\n"
	const manifestPath = "/proj/.pnp.cjs"
	m, ok := Load(manifestPath, func(p string) (string, bool) {
		if p == manifestPath {
			return cjs, true
		}
		return "", false
	})
	if !ok {
		t.Fatal("failed to load tightly-spaced inlined .pnp.cjs")
	}
	if r := m.Resolve("@t/util", "/proj"); r.Status != Success {
		t.Fatalf("@t/util from tight-spacing manifest: status %d", r.Status)
	}
}

// A manifest whose dependency locations are Yarn __virtual__ paths (the shape
// emitted for packages with peer dependencies): a workspace lib virtualized into
// a plain directory (count 1) and an npm package virtualized into a cache .zip
// (count 0).
const virtualFixtureManifest = `{
  "dependencyTreeRoots": [{"name": "root", "reference": "workspace:."}],
  "enableTopLevelFallback": true,
  "fallbackExclusionList": [],
  "fallbackPool": [],
  "packageRegistryData": [
    [null, [["workspace:.", {
      "packageLocation": "./",
      "packageDependencies": [
        ["@w/ui", "virtual:aaa#workspace:packages/ui"],
        ["react-dom", "virtual:bbb#npm:18.3.1"]
      ],
      "linkType": "SOFT"
    }]]],
    ["@w/ui", [["virtual:aaa#workspace:packages/ui", {
      "packageLocation": "./.yarn/__virtual__/w-ui-virtual-aaa/1/packages/ui/",
      "packageDependencies": [["@w/ui", "virtual:aaa#workspace:packages/ui"]],
      "linkType": "SOFT"
    }]]],
    ["react-dom", [["virtual:bbb#npm:18.3.1", {
      "packageLocation": "./.yarn/__virtual__/react-dom-virtual-bbb/0/cache/react-dom-npm-18.3.1-abc.zip/node_modules/react-dom/",
      "packageDependencies": [
        ["react-dom", "virtual:bbb#npm:18.3.1"],
        ["scheduler", "npm:0.23.0"]
      ],
      "linkType": "HARD"
    }]]],
    ["scheduler", [["npm:0.23.0", {
      "packageLocation": "./.yarn/cache/scheduler-npm-0.23.0-xyz.zip/node_modules/scheduler/",
      "packageDependencies": [["scheduler", "npm:0.23.0"]],
      "linkType": "HARD"
    }]]]
  ]
}`

func TestResolveKeepsVirtualLocation(t *testing.T) {
	t.Parallel()
	const manifestPath = "/proj/.pnp.data.json"
	m, ok := Load(manifestPath, func(p string) (string, bool) {
		if p == manifestPath {
			return virtualFixtureManifest, true
		}
		return "", false
	})
	if !ok {
		t.Fatal("failed to load virtual fixture manifest")
	}
	// A virtualized package's PackageDir is returned as the __virtual__ path, NOT
	// dereferenced: the locator table is keyed by these raw locations, so a file
	// resolved under the virtual path must keep it for findLocator to identify its
	// owning package when resolving the package's own imports. The zip overlay
	// dereferences the path only when it actually reads from the filesystem.
	wantUI := "/proj/.yarn/__virtual__/w-ui-virtual-aaa/1/packages/ui"
	if r := m.Resolve("@w/ui", "/proj"); r.Status != Success || r.PackageDir != wantUI {
		t.Fatalf("@w/ui: status %d dir %q, want Success %q", r.Status, r.PackageDir, wantUI)
	}
	wantRD := "/proj/.yarn/__virtual__/react-dom-virtual-bbb/0/cache/react-dom-npm-18.3.1-abc.zip/node_modules/react-dom"
	if r := m.Resolve("react-dom", "/proj"); r.Status != Success || r.PackageDir != wantRD {
		t.Fatalf("react-dom: status %d dir %q, want %q", r.Status, r.PackageDir, wantRD)
	}
	// findLocator round-trips the virtual location back to the owning package, so
	// an import made from inside the virtualized package resolves against that
	// package's own dependency map (not the root's).
	locator, found := m.findLocator(wantUI)
	if !found || locator.ident != "@w/ui" {
		t.Fatalf("findLocator(%q) = %+v found=%v, want @w/ui", wantUI, locator, found)
	}
}

// TestResolveTransitiveDepFromVirtualizedPackage is the direct regression for the
// bug where dereferencing PackageDir broke findLocator: "scheduler" is a
// dependency of react-dom but NOT of the root. A file inside react-dom's
// virtualized location importing "scheduler" must resolve through react-dom's own
// dependency map. If the virtual location were collapsed to its real .zip path,
// findLocator would miss react-dom's (virtual-keyed) locator, fall through to the
// root — which does not declare scheduler — and the import would be NotFound.
func TestResolveTransitiveDepFromVirtualizedPackage(t *testing.T) {
	t.Parallel()
	const manifestPath = "/proj/.pnp.data.json"
	m, ok := Load(manifestPath, func(p string) (string, bool) {
		if p == manifestPath {
			return virtualFixtureManifest, true
		}
		return "", false
	})
	if !ok {
		t.Fatal("failed to load virtual fixture manifest")
	}
	// Sanity: scheduler is not resolvable from the root (it is not a root dep and
	// the root does not fall back), so a successful resolution below can only come
	// from react-dom's own dependency map.
	if r := m.Resolve("scheduler", "/proj"); r.Status == Success {
		t.Fatalf("scheduler unexpectedly resolvable from root: %q", r.PackageDir)
	}
	// Model the real resolver flow: react-dom's resolved PackageDir becomes the
	// importer directory for react-dom's own imports. Feeding that returned
	// location straight back into Resolve is what exposes the bug — if PackageDir
	// were dereferenced to the real .zip path, the importer would no longer be in
	// the virtual locator space and findLocator would miss react-dom.
	rd := m.Resolve("react-dom", "/proj")
	if rd.Status != Success {
		t.Fatalf("react-dom from root: status %d, want Success", rd.Status)
	}
	r := m.Resolve("scheduler", rd.PackageDir)
	if r.Status != Success {
		t.Fatalf("scheduler from inside react-dom (importer %q): status %d, want Success (findLocator must match the virtual locator)", rd.PackageDir, r.Status)
	}
	if !strings.Contains(r.PackageDir, "scheduler-npm-0.23.0-xyz.zip") {
		t.Fatalf("scheduler PackageDir = %q, want the scheduler cache zip", r.PackageDir)
	}
}

func TestFindWalksUp(t *testing.T) {
	t.Parallel()
	// Only .pnp.cjs exists (Yarn's inlining default): Find returns it.
	got := Find("/proj/packages/app/src", func(p string) bool { return p == "/proj/.pnp.cjs" })
	if got != "/proj/.pnp.cjs" {
		t.Fatalf("Find = %q, want /proj/.pnp.cjs", got)
	}
	// When both exist, the .pnp.data.json sidecar is preferred.
	got = Find("/proj/packages/app/src", func(p string) bool {
		return p == "/proj/.pnp.cjs" || p == "/proj/.pnp.data.json"
	})
	if got != "/proj/.pnp.data.json" {
		t.Fatalf("Find = %q, want /proj/.pnp.data.json", got)
	}
	if Find("/elsewhere", func(string) bool { return false }) != "" {
		t.Fatal("Find should return empty when no manifest exists")
	}
}

func TestLoadInlinedCjs(t *testing.T) {
	t.Parallel()
	// Emulate a .pnp.cjs: the manifest JSON embedded as a backslash-continued,
	// single-quoted RAW_RUNTIME_STATE literal (backslashes doubled, newlines
	// continued), the shape Yarn writes by default (pnpEnableInlining).
	jsEscaped := strings.ReplaceAll(fixtureManifest, `\`, `\\`)
	jsEscaped = strings.ReplaceAll(jsEscaped, "\n", "\\\n")
	cjs := "/* prologue */\nconst RAW_RUNTIME_STATE =\n'" + jsEscaped +
		"'\n;\nfunction $$SETUP_STATE(h){return h(JSON.parse(RAW_RUNTIME_STATE));}\n"

	const manifestPath = "/proj/.pnp.cjs"
	m, ok := Load(manifestPath, func(p string) (string, bool) {
		if p == manifestPath {
			return cjs, true
		}
		return "", false
	})
	if !ok {
		t.Fatal("failed to load inlined .pnp.cjs")
	}
	if r := m.Resolve("@t/util", "/proj"); r.Status != Success || r.PackageDir != "/proj/packages/util" {
		t.Fatalf("@t/util from inlined manifest: status %d dir %q", r.Status, r.PackageDir)
	}
}

func TestDerefVirtualPath(t *testing.T) {
	t.Parallel()
	cases := []struct{ in, want string }{
		// n=0: no ".." applied, the __virtual__/<hash>/0 segment is removed.
		{
			"/proj/.yarn/__virtual__/react-dom-virtual-7b/0/cache/react-dom.zip/node_modules/react-dom/index.js",
			"/proj/.yarn/cache/react-dom.zip/node_modules/react-dom/index.js",
		},
		// n=1 into a plain directory (a workspace lib with a peer dep): pops the
		// ".yarn" segment before the triple. This is the non-.zip virtual path that
		// the overlay must dereference and delegate. Mirrors the real Yarn location
		// ./.yarn/__virtual__/@w-ui-virtual-<hash>/1/packages/ui/.
		{
			"/proj/.yarn/__virtual__/w-ui-virtual-aaa/1/packages/ui/src/index.ts",
			"/proj/packages/ui/src/index.ts",
		},
		// $$virtual: the legacy (Yarn <3) spelling is handled identically.
		{
			"/proj/.yarn/$$virtual/w-ui-virtual-aaa/1/packages/ui/src/index.ts",
			"/proj/packages/ui/src/index.ts",
		},
		// n=2 pops two parent segments.
		{
			"/proj/a/b/__virtual__/h/2/x/y.ts",
			"/proj/x/y.ts",
		},
		// Non-numeric count: NOT a Yarn virtual path, returned unchanged. This is
		// the guard that keeps a real directory literally named "__virtual__" from
		// being corrupted.
		{
			"/proj/src/__virtual__/fixtures/index.ts",
			"/proj/src/__virtual__/fixtures/index.ts",
		},
		// __virtual__ as the last segment (no hash/count following): unchanged.
		{
			"/proj/src/__virtual__",
			"/proj/src/__virtual__",
		},
		// Ported from Yarn berry VirtualFS.test.ts "should ignore non-hash virtual
		// components": a __virtual__ dir with a child but no count segment after it
		// is left unchanged.
		{
			"/proj/__virtual__/package.json",
			"/proj/__virtual__/package.json",
		},
		// berry "should map numbered virtual components (1, no file)": the package
		// dir itself at depth 1.
		{
			"/proj/__virtual__/12345/1",
			"/",
		},
		// berry "should preserve dots when mapping" and "empty strings".
		{".", "."},
		{"", ""},
		// Count larger than the available parents: the pop clamps at the root
		// rather than underflowing.
		{
			"/__virtual__/h/9/x.ts",
			"/x.ts",
		},
		// Count with no suffix after it (the package dir itself): the trailing
		// separator is dropped rather than left dangling.
		{
			"/root/.yarn/__virtual__/h/0",
			"/root/.yarn",
		},
		// no virtual segment: unchanged.
		{"/proj/.yarn/cache/react.zip/node_modules/react/index.js", "/proj/.yarn/cache/react.zip/node_modules/react/index.js"},
		{"/proj/packages/app/src/index.ts", "/proj/packages/app/src/index.ts"},
	}
	for _, c := range cases {
		if got := DerefVirtualPath(c.in); got != c.want {
			t.Errorf("DerefVirtualPath(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
