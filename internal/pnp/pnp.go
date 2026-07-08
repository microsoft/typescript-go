// Package pnp implements Yarn Plug'n'Play module resolution. A PnP project has
// no node_modules tree; instead a manifest (.pnp.data.json, or .pnp.cjs with the
// data embedded) maps every (package, importer) pair to the on-disk location of
// the dependency. This package loads that manifest and answers the one question
// the module resolver needs: given a bare specifier's package identifier and the
// directory it is imported from, where does that package live?
//
// The algorithm is a transcription of the Yarn PnP specification's
// RESOLVE_TO_UNQUALIFIED, following the reference resolver's actual behavior where
// it diverges from the spec (the trailing-slash longest-prefix locator loop). It
// deliberately mirrors esbuild's Go implementation so the two agree.
package pnp

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/microsoft/typescript-go/internal/json"
	"github.com/microsoft/typescript-go/internal/tspath"
)

// identAndRef is a dependency target: a reference (ident==""), an alias (both
// set), or an unfulfilled peer dependency (both ""). An absent entry is distinct
// from a present-but-null one, so callers check presence separately.
type identAndRef struct {
	ident     string
	reference string
}

type pkg struct {
	packageLocation     string
	packageDependencies map[string]identAndRef
}

type locatorByLocation struct {
	locator           identAndRef
	discardFromLookup bool
}

// Manifest is a compiled PnP manifest.
type Manifest struct {
	// manifestDir is the absolute, normalized directory the manifest lives in;
	// all packageLocation values are relative to it.
	manifestDir string

	enableTopLevelFallback bool
	fallbackPool           map[string]identAndRef
	fallbackExclusionList  map[string]map[string]bool
	ignorePattern          *regexp.Regexp

	// packageRegistry[ident][reference] -> package. The top-level package is
	// keyed by ("", "").
	packageRegistry map[string]map[string]pkg
	// packageLocatorsByLocation reverse-indexes a package location (verbatim,
	// as stored in the manifest, relative and trailing-slashed) to its locator.
	packageLocatorsByLocation map[string]locatorByLocation
}

// --- raw JSON shapes -------------------------------------------------------

// The manifest tables are arrays of pairs, e.g. packageRegistryData is
// [[ident|null, [[reference|null, packageObject], ...]], ...]. json.Value
// preserves the null-vs-string distinction the spec depends on.
type rawManifest struct {
	EnableTopLevelFallback bool            `json:"enableTopLevelFallback"`
	IgnorePatternData      string          `json:"ignorePatternData"`
	FallbackPool           [][2]json.Value `json:"fallbackPool"`
	FallbackExclusionList  [][2]json.Value `json:"fallbackExclusionList"`
	PackageRegistryData    [][2]json.Value `json:"packageRegistryData"`
}

type rawPackage struct {
	PackageLocation     string          `json:"packageLocation"`
	PackageDependencies [][2]json.Value `json:"packageDependencies"`
	DiscardFromLookup   bool            `json:"discardFromLookup"`
}

// stringOrNull decodes a JSON string or null; null yields ("", true), matching
// the spec's convention of encoding a null ident/reference as the empty string.
func stringOrNull(raw json.Value) (string, bool) {
	var s *string
	if err := json.Unmarshal(raw, &s); err != nil {
		return "", false
	}
	if s == nil {
		return "", true
	}
	return *s, true
}

// dependencyTarget decodes a packageDependencies value: a string reference, a
// [ident, reference] alias array, or null (unfulfilled peer). ok is false only
// on a malformed value.
func dependencyTarget(raw json.Value) (identAndRef, bool) {
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "null" {
		return identAndRef{}, true
	}
	if strings.HasPrefix(trimmed, "[") {
		var pair [2]string
		if err := json.Unmarshal(raw, &pair); err != nil {
			return identAndRef{}, false
		}
		return identAndRef{ident: pair[0], reference: pair[1]}, true
	}
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return identAndRef{}, false
	}
	return identAndRef{reference: s}, true
}

// Load parses and compiles a PnP manifest. manifestPath is the absolute,
// normalized path of a .pnp.data.json sidecar or the .pnp.cjs the data is inlined
// into (Yarn's default); readFile reads it.
func Load(manifestPath string, readFile func(string) (string, bool)) (*Manifest, bool) {
	contents, ok := readFile(manifestPath)
	if !ok {
		return nil, false
	}
	if !strings.HasSuffix(manifestPath, ".json") {
		// A .pnp.cjs/.pnp.js embeds the data as a JS string literal; extract it.
		contents, ok = extractInlinedManifest(contents)
		if !ok {
			return nil, false
		}
	}
	var raw rawManifest
	if err := json.Unmarshal([]byte(contents), &raw); err != nil {
		return nil, false
	}

	m := &Manifest{
		manifestDir:               tspath.GetDirectoryPath(manifestPath),
		enableTopLevelFallback:    raw.EnableTopLevelFallback,
		fallbackPool:              map[string]identAndRef{},
		fallbackExclusionList:     map[string]map[string]bool{},
		packageRegistry:           map[string]map[string]pkg{},
		packageLocatorsByLocation: map[string]locatorByLocation{},
	}

	if raw.IgnorePatternData != "" {
		m.ignorePattern = compileIgnorePattern(raw.IgnorePatternData)
	}

	for _, entry := range raw.FallbackExclusionList {
		ident, ok := stringOrNull(entry[0])
		if !ok {
			continue
		}
		var refs []string
		if err := json.Unmarshal(entry[1], &refs); err != nil {
			continue
		}
		set := make(map[string]bool, len(refs))
		for _, r := range refs {
			set[r] = true
		}
		m.fallbackExclusionList[ident] = set
	}

	for _, entry := range raw.FallbackPool {
		ident, ok := stringOrNull(entry[0])
		if !ok {
			continue
		}
		if target, ok := dependencyTarget(entry[1]); ok {
			m.fallbackPool[ident] = target
		}
	}

	for _, entry := range raw.PackageRegistryData {
		ident, ok := stringOrNull(entry[0])
		if !ok {
			continue
		}
		var refList [][2]json.Value
		if err := json.Unmarshal(entry[1], &refList); err != nil {
			continue
		}
		refs := m.packageRegistry[ident]
		if refs == nil {
			refs = map[string]pkg{}
			m.packageRegistry[ident] = refs
		}
		for _, refEntry := range refList {
			reference, ok := stringOrNull(refEntry[0])
			if !ok {
				continue
			}
			var rp rawPackage
			if err := json.Unmarshal(refEntry[1], &rp); err != nil {
				continue
			}
			deps := make(map[string]identAndRef, len(rp.PackageDependencies))
			for _, dep := range rp.PackageDependencies {
				var depIdent string
				if err := json.Unmarshal(dep[0], &depIdent); err != nil {
					continue
				}
				if target, ok := dependencyTarget(dep[1]); ok {
					deps[depIdent] = target
				}
			}
			refs[reference] = pkg{
				packageLocation:     rp.PackageLocation,
				packageDependencies: deps,
			}

			// Replicate Yarn's hydrateRuntimeState: first locator at a location
			// wins unless discarded; a later non-discarded entry overwrites a
			// discarded one; the discard flag is AND-ed across entries.
			if existing, ok := m.packageLocatorsByLocation[rp.PackageLocation]; !ok {
				m.packageLocatorsByLocation[rp.PackageLocation] = locatorByLocation{
					locator:           identAndRef{ident: ident, reference: reference},
					discardFromLookup: rp.DiscardFromLookup,
				}
			} else {
				existing.discardFromLookup = existing.discardFromLookup && rp.DiscardFromLookup
				if !rp.DiscardFromLookup {
					existing.locator = identAndRef{ident: ident, reference: reference}
				}
				m.packageLocatorsByLocation[rp.PackageLocation] = existing
			}
		}
	}

	return m, true
}

// compileIgnorePattern strips the negative-lookahead sequences Yarn emits (Go's
// RE2 rejects them; they only exclude "."/".." path segments that never occur
// here) and compiles the rest, or returns nil if it still won't compile.
func compileIgnorePattern(src string) *regexp.Regexp {
	for _, la := range []string{
		`(?!\.)`,
		`(?!(?:^|\/)\.)`,
		`(?!\.{1,2}(?:\/|$))`,
		`(?!(?:^|\/)\.{1,2}(?:\/|$))`,
	} {
		src = strings.ReplaceAll(src, la, "")
	}
	reg, err := regexp.Compile(src)
	if err != nil {
		return nil
	}
	return reg
}

// extractInlinedManifest pulls the JSON payload out of a .pnp.cjs. Yarn embeds it
// as `const RAW_RUNTIME_STATE = '<json>';` — a single-quoted JS string literal
// whose newlines are backslash line-continuations — and passes it to
// hydrateRuntimeState(JSON.parse(RAW_RUNTIME_STATE)). This finds that literal and
// unescapes it back to JSON without a JS parser. Returns false if the marker or a
// closing quote is absent.
func extractInlinedManifest(contents string) (string, bool) {
	const marker = "RAW_RUNTIME_STATE"
	_, rest, found := strings.Cut(contents, marker)
	if !found {
		return "", false
	}
	// Require the assignment operator (so an earlier textual mention of
	// RAW_RUNTIME_STATE — a comment, or the JSON.parse(RAW_RUNTIME_STATE) use —
	// cannot be mistaken for the declaration), tolerating any spacing/newlines
	// around it (`RAW_RUNTIME_STATE =`, `RAW_RUNTIME_STATE=`, `... =\n'`). Then
	// find the opening quote of the string literal.
	rest = strings.TrimLeft(rest, " \t\r\n")
	rest, found = strings.CutPrefix(rest, "=")
	if !found {
		return "", false
	}
	q := strings.IndexAny(rest, "'\"`")
	if q < 0 {
		return "", false
	}
	quote := rest[q]
	body := rest[q+1:]
	var sb strings.Builder
	for j := 0; j < len(body); j++ {
		c := body[j]
		if c == quote {
			return sb.String(), true
		}
		if c != '\\' || j+1 >= len(body) {
			sb.WriteByte(c)
			continue
		}
		n := body[j+1]
		j++
		switch n {
		case '\n': // line continuation: emit nothing
		case '\r':
			if j+1 < len(body) && body[j+1] == '\n' {
				j++
			}
		case 'n':
			sb.WriteByte('\n')
		case 't':
			sb.WriteByte('\t')
		case 'r':
			sb.WriteByte('\r')
		case 'b':
			sb.WriteByte('\b')
		case 'f':
			sb.WriteByte('\f')
		case 'v':
			sb.WriteByte('\v')
		case '0':
			sb.WriteByte(0)
		default:
			// \\, \', \", \`, and any other escape: the following character
			// verbatim (JS treats an unknown escape as the character itself).
			sb.WriteByte(n)
		}
	}
	return "", false
}

// DerefVirtualPath rewrites a Yarn PnP virtual path
// (<prefix>/__virtual__/<hash>/<n>/<suffix>) to its real backing path by applying
// the ".." operation n times to the prefix, so file access reaches the real cache
// archive. A path with no virtual segment is returned unchanged. Ported from
// esbuild's ParseYarnPnPVirtualPath; forward-slash paths only (tsgo normalizes).
func DerefVirtualPath(path string) string {
	i := 0
	for {
		start := i
		slash := strings.IndexByte(path[i:], '/')
		if slash == -1 {
			return path
		}
		i += slash + 1
		segment := path[start : i-1]
		if segment != "__virtual__" && segment != "$$virtual" {
			continue
		}
		slash = strings.IndexByte(path[i:], '/')
		if slash == -1 {
			return path
		}
		j := i + slash + 1
		var count, suffix string
		if s := strings.IndexByte(path[j:], '/'); s != -1 {
			count = path[j : j+s]
			suffix = path[j+s:]
		} else {
			count = path[j:]
		}
		n, err := strconv.ParseInt(count, 10, 64)
		if err != nil {
			continue
		}
		prefix := path[:start]
		for n > 0 && strings.HasSuffix(prefix, "/") {
			s := strings.LastIndexByte(prefix[:len(prefix)-1], '/')
			if s == -1 {
				break
			}
			prefix = prefix[:s+1]
			n--
		}
		if suffix == "" && strings.IndexByte(prefix, '/') != strings.LastIndexByte(prefix, '/') {
			prefix = prefix[:len(prefix)-1]
		} else if prefix == "" {
			prefix = "."
		} else if strings.HasPrefix(suffix, "/") {
			suffix = suffix[1:]
		}
		return prefix + suffix
	}
}

// parseBareIdentifier splits a specifier into its package identifier and the
// subpath within the package. A scoped name without a "/" is invalid.
func parseBareIdentifier(specifier string) (ident string, modulePath string, ok bool) {
	slash := strings.IndexByte(specifier, '/')
	if strings.HasPrefix(specifier, "@") {
		if slash == -1 {
			return "", "", false
		}
		if slash2 := strings.IndexByte(specifier[slash+1:], '/'); slash2 != -1 {
			ident = specifier[:slash+1+slash2]
		} else {
			ident = specifier
		}
	} else if slash != -1 {
		ident = specifier[:slash]
	} else {
		ident = specifier
	}
	modulePath = specifier[len(ident):]
	return ident, modulePath, true
}

// Status is the outcome of a PnP resolution.
type Status uint8

const (
	// Skipped: the importer is not governed by PnP (ignored path, or no locator
	// found). The caller should fall back to classic node_modules resolution.
	Skipped Status = iota
	// Success: the package was located.
	Success
	// NotFound: the specifier is not a declared dependency of the importer.
	NotFound
	// UnfulfilledPeer: the dependency is a null-valued (unfulfilled) peer.
	UnfulfilledPeer
	// Failed: a manifest invariant was violated (should not happen).
	Failed
)

// Result of a resolution.
type Result struct {
	Status Status
	// PackageDir is the absolute, normalized directory the resolved package
	// lives in (may be inside a Yarn cache .zip). Valid only when Status is
	// Success.
	PackageDir string
	// Subpath is the portion of the specifier after the package identifier
	// (including a leading "/" when present).
	Subpath string
}

// Resolve locates the package named by specifier as imported from importerDir
// (an absolute, normalized directory). It answers RESOLVE_TO_UNQUALIFIED: the
// package directory, before file/extension/exports resolution runs on it.
func (m *Manifest) Resolve(specifier string, importerDir string) Result {
	ident, modulePath, ok := parseBareIdentifier(specifier)
	if !ok {
		return Result{Status: Failed}
	}

	parentLocator, ok := m.findLocator(importerDir)
	if !ok {
		return Result{Status: Skipped}
	}

	parentPkg, ok := m.getPackage(parentLocator.ident, parentLocator.reference)
	if !ok {
		return Result{Status: Failed}
	}

	target, ok := parentPkg.packageDependencies[ident]
	if !ok || target.reference == "" {
		// Not listed (or null-valued): try the fallback set. The synthetic
		// top-level locator (empty ident, name===null in Yarn) never uses
		// fallbacks — it IS the fallback source, so a specifier the root does not
		// declare is genuinely NotFound rather than pulled from the pool. Workspaces
		// are kept out separately via the fallbackExclusionList. (Mirrors berry
		// makeApi resolveToUnqualified: `if (issuerLocator.name !== null)` guarding
		// the fallback block, plus the exclusion check.)
		if m.enableTopLevelFallback && parentLocator.ident != "" {
			if set := m.fallbackExclusionList[parentLocator.ident]; !set[parentLocator.reference] {
				if fallback, found := m.resolveViaFallback(ident); found && fallback.reference != "" {
					target = fallback
					ok = true
				}
			}
		}
	}

	if !ok {
		return Result{Status: NotFound}
	}
	if target.reference == "" {
		return Result{Status: UnfulfilledPeer}
	}

	var dependencyPkg pkg
	if target.ident != "" {
		// Aliased dependency.
		if dependencyPkg, ok = m.getPackage(target.ident, target.reference); !ok {
			return Result{Status: Failed}
		}
	} else {
		if dependencyPkg, ok = m.getPackage(ident, target.reference); !ok {
			return Result{Status: Failed}
		}
	}

	// path.resolve(manifestDir, packageLocation, modulePath). packageLocation is
	// relative (e.g. "./packages/util/" or "./.yarn/cache/x.zip/node_modules/x/"),
	// and for a peer-dependency-virtualized package it is a __virtual__/<hash>/<n>
	// path. The location is returned as-is (virtual paths are NOT dereferenced
	// here): the package registry's locator table is keyed by these raw locations,
	// so a file resolved under a virtual path must keep that path for findLocator to
	// identify its owning package when resolving its own imports. The zip overlay
	// dereferences virtual paths when it actually reads from the filesystem.
	// Canonicalize to a directory path with no trailing separator; the resolver
	// joins file names onto it.
	pkgDir := tspath.RemoveTrailingDirectorySeparator(tspath.NormalizePath(tspath.CombinePaths(m.manifestDir, dependencyPkg.packageLocation)))
	return Result{
		Status:     Success,
		PackageDir: pkgDir,
		Subpath:    modulePath,
	}
}

// findLocator finds the package that owns importerDir, walking up path prefixes.
func (m *Manifest) findLocator(importerDir string) (identAndRef, bool) {
	// A non-absolute importer directory cannot be located against the manifest's
	// absolute locations; treat it as not PnP-governed rather than risk a relative
	// path accidentally prefix-matching a locator.
	if !tspath.IsRootedDiskPath(importerDir) {
		return identAndRef{}, false
	}
	// Relative path from the manifest directory to the importer, forward-slashed.
	rel := tspath.GetRelativePathFromDirectory(m.manifestDir, importerDir, tspath.ComparePathsOptions{
		UseCaseSensitiveFileNames: true,
		CurrentDirectory:          m.manifestDir,
	})
	rel = tspath.NormalizeSlashes(rel)
	rel = strings.TrimPrefix(rel, "./")

	if m.ignorePattern != nil && m.ignorePattern.MatchString(rel) {
		return identAndRef{}, false
	}

	// Locations in the manifest are trailing-slashed and start with ./ or ../.
	if !strings.HasSuffix(rel, "/") {
		rel += "/"
	}
	if !strings.HasPrefix(rel, "./") && !strings.HasPrefix(rel, "../") {
		rel = "./" + rel
	}

	// Longest-prefix loop, per the reference implementation (not the spec's
	// hypothetical algorithm).
	for {
		entry, ok := m.packageLocatorsByLocation[rel]
		if !ok || entry.discardFromLookup {
			idx := strings.LastIndexByte(rel[:len(rel)-1], '/')
			if idx < 0 {
				break
			}
			rel = rel[:idx+1]
			if rel == "" {
				break
			}
			continue
		}
		return entry.locator, true
	}
	return identAndRef{}, false
}

func (m *Manifest) resolveViaFallback(ident string) (identAndRef, bool) {
	topLevel, ok := m.getPackage("", "")
	if !ok {
		return identAndRef{}, false
	}
	if target, found := topLevel.packageDependencies[ident]; found {
		return target, true
	}
	target, found := m.fallbackPool[ident]
	return target, found
}

func (m *Manifest) getPackage(ident string, reference string) (pkg, bool) {
	if refs, ok := m.packageRegistry[ident]; ok {
		if p, ok := refs[reference]; ok {
			return p, true
		}
	}
	return pkg{}, false
}
