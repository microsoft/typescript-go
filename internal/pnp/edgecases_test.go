package pnp

import (
	"strings"
	"testing"
)

// A manifest exercising the top-level fallback: enableTopLevelFallback is on, the
// top-level package declares "shared", and app-b is on the fallbackExclusionList.
// app-a (not excluded) can reach "shared" through the fallback even though it does
// not declare it; app-b (excluded) cannot.
const fallbackManifest = `{
  "dependencyTreeRoots": [{"name": "root", "reference": "workspace:."}],
  "enableTopLevelFallback": true,
  "fallbackExclusionList": [["app-b", ["workspace:packages/app-b"]]],
  "fallbackPool": [["pooled", "npm:9.0.0"]],
  "packageRegistryData": [
    [null, [[null, {
      "packageLocation": "./",
      "packageDependencies": [
        ["app-a", "workspace:packages/app-a"],
        ["app-b", "workspace:packages/app-b"],
        ["shared", "npm:1.0.0"]
      ],
      "linkType": "SOFT"
    }]]],
    ["app-a", [["workspace:packages/app-a", {
      "packageLocation": "./packages/app-a/",
      "packageDependencies": [["app-a", "workspace:packages/app-a"]],
      "linkType": "SOFT"
    }]]],
    ["app-b", [["workspace:packages/app-b", {
      "packageLocation": "./packages/app-b/",
      "packageDependencies": [["app-b", "workspace:packages/app-b"]],
      "linkType": "SOFT"
    }]]],
    ["shared", [["npm:1.0.0", {
      "packageLocation": "./.yarn/cache/shared-npm-1.0.0-abc.zip/node_modules/shared/",
      "packageDependencies": [["shared", "npm:1.0.0"]],
      "linkType": "HARD"
    }]]],
    ["pooled", [["npm:9.0.0", {
      "packageLocation": "./.yarn/cache/pooled-npm-9.0.0-def.zip/node_modules/pooled/",
      "packageDependencies": [["pooled", "npm:9.0.0"]],
      "linkType": "HARD"
    }]]]
  ]
}`

func loadFallback(t *testing.T) *Manifest {
	t.Helper()
	const manifestPath = "/proj/.pnp.data.json"
	m, ok := Load(manifestPath, func(p string) (string, bool) {
		if p == manifestPath {
			return fallbackManifest, true
		}
		return "", false
	})
	if !ok {
		t.Fatal("failed to load fallback manifest")
	}
	return m
}

func TestResolveTopLevelFallback(t *testing.T) {
	t.Parallel()
	m := loadFallback(t)
	// app-a does not declare "shared", but the top-level package does and app-a is
	// not on the exclusion list: the fallback resolves it.
	if r := m.Resolve("shared", "/proj/packages/app-a"); r.Status != Success ||
		!strings.Contains(r.PackageDir, "shared-npm-1.0.0-abc.zip") {
		t.Fatalf("app-a fallback to shared: status %d dir %q", r.Status, r.PackageDir)
	}
	// The fallback pool (distinct from the top-level deps) also satisfies an
	// undeclared specifier.
	if r := m.Resolve("pooled", "/proj/packages/app-a"); r.Status != Success ||
		!strings.Contains(r.PackageDir, "pooled-npm-9.0.0-def.zip") {
		t.Fatalf("app-a fallback to pooled: status %d dir %q", r.Status, r.PackageDir)
	}
}

func TestResolveFallbackExclusion(t *testing.T) {
	t.Parallel()
	m := loadFallback(t)
	// app-b IS on the fallbackExclusionList, so the top-level fallback does not
	// apply: an undeclared specifier is NotFound, not silently resolved.
	if r := m.Resolve("shared", "/proj/packages/app-b"); r.Status != NotFound {
		t.Fatalf("app-b (excluded) resolving shared: status %d, want NotFound", r.Status)
	}
}

func TestResolveTopLevelDoesNotFallBack(t *testing.T) {
	t.Parallel()
	m := loadFallback(t)
	// The synthetic top-level (root) locator never uses fallbacks — it is the
	// fallback source. "pooled" lives only in the fallback pool, so a non-root
	// package reaches it but the root does not (an undeclared specifier from the
	// root is NotFound, matching Yarn's `issuerLocator.name !== null` guard).
	if r := m.Resolve("pooled", "/proj/packages/app-a"); r.Status != Success {
		t.Fatalf("app-a resolving pooled: status %d, want Success", r.Status)
	}
	if r := m.Resolve("pooled", "/proj"); r.Status != NotFound {
		t.Fatalf("root resolving pooled via pool: status %d, want NotFound (root must not fall back)", r.Status)
	}
}

func TestResolveBareIdentifierParse(t *testing.T) {
	t.Parallel()
	m := loadFixture(t)
	cases := []struct {
		specifier string
		wantFail  bool
	}{
		{"@scope", true},          // scoped name with no second segment is invalid
		{"@scope/pkg/sub", false}, // scoped name + subpath parses (then resolves/NotFound)
		{"pkg/sub/deep", false},   // unscoped name + multi-segment subpath parses
	}
	for _, c := range cases {
		r := m.Resolve(c.specifier, "/proj/packages/util")
		if c.wantFail && r.Status != Failed {
			t.Errorf("Resolve(%q) status = %d, want Failed", c.specifier, r.Status)
		}
		if !c.wantFail && r.Status == Failed {
			t.Errorf("Resolve(%q) status = Failed, want a parse success", c.specifier)
		}
	}
}

func TestResolveCarriesScopedSubpath(t *testing.T) {
	t.Parallel()
	m := loadFixture(t)
	// The subpath after a scoped package name is carried through unresolved-file
	// resolution as the Subpath (here the package @t/util resolves and "/deep/x"
	// rides along).
	r := m.Resolve("@t/util/deep/x", "/proj")
	if r.Status != Success || r.Subpath != "/deep/x" {
		t.Fatalf("@t/util/deep/x: status %d subpath %q", r.Status, r.Subpath)
	}
}

func TestLoadRejectsMalformed(t *testing.T) {
	t.Parallel()
	read := func(contents string) func(string) (string, bool) {
		return func(p string) (string, bool) {
			if p == "/proj/.pnp.data.json" || p == "/proj/.pnp.cjs" {
				return contents, true
			}
			return "", false
		}
	}
	// Malformed JSON sidecar.
	if _, ok := Load("/proj/.pnp.data.json", read(`{not json`)); ok {
		t.Error("Load accepted malformed JSON")
	}
	// A .pnp.cjs with no RAW_RUNTIME_STATE marker.
	if _, ok := Load("/proj/.pnp.cjs", read("module.exports = {};\n")); ok {
		t.Error("Load accepted a .cjs with no manifest marker")
	}
	// The marker present but NOT an assignment (e.g. only the JSON.parse use):
	// requiring "=" keeps a stray textual mention from being mis-extracted.
	if _, ok := Load("/proj/.pnp.cjs", read("hydrate(JSON.parse(RAW_RUNTIME_STATE));\n")); ok {
		t.Error("Load accepted a RAW_RUNTIME_STATE use with no assignment")
	}
	// A missing file.
	if _, ok := Load("/proj/.pnp.data.json", func(string) (string, bool) { return "", false }); ok {
		t.Error("Load accepted a missing manifest file")
	}
}
