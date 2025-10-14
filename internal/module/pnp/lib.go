// BSD 2-Clause License
//
// Copyright (c) 2016-present, Yarn Contributors.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
//  1. Redistributions of source code must retain the above copyright notice, this
//     list of conditions and the following disclaimer.
//
//  2. Redistributions in binary form must reproduce the above copyright notice,
//     this list of conditions and the following disclaimer in the documentation
//     and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
package pnp

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/go-json-experiment/json"
	"github.com/microsoft/typescript-go/internal/module/pnp/utils"
)

type ResolutionKind int

const (
	ResolutionSkipped ResolutionKind = iota
	ResolutionResolved
)

type ResolutionHost struct {
	FindPNPManifest func(start string) (*Manifest, error)
}

type ResolutionConfig struct {
	Host ResolutionHost
}

type Resolution struct {
	Kind       ResolutionKind
	Path       string
	ModulePath *string
}

type UndeclaredDependencyError struct {
	Message        string
	Request        string
	DependencyName string
	IssuerLocator  PackageLocator
	IssuerPath     string
}

func (e *UndeclaredDependencyError) Error() string { return e.Message }

type MissingPeerDependencyError struct {
	Message         string
	Request         string
	DependencyName  string
	IssuerLocator   PackageLocator
	IssuerPath      string
	BrokenAncestors []PackageLocator
}

func (e *MissingPeerDependencyError) Error() string { return e.Message }

type FailedManifestHydrationError struct {
	Message      string
	ManifestPath string
	Err          error
}

func (e *FailedManifestHydrationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s\n\nOriginal error: %v", e.Message, e.Err)
	}
	return e.Message
}
func (e *FailedManifestHydrationError) Unwrap() error { return e.Err }

func parseScopedPackageName(spec string) (pkg string, sub *string, ok bool) {
	parts := strings.SplitN(spec, "/", 3)
	if len(parts) < 2 {
		return "", nil, false
	}
	scope, name := parts[0], parts[1]
	if scope == "" || name == "" {
		return "", nil, false
	}
	pkg = scope + "/" + name
	if len(parts) == 3 {
		s := parts[2]
		sub = &s
	}
	return pkg, sub, true
}

func parseGlobalPackageName(spec string) (pkg string, sub *string, ok bool) {
	parts := strings.SplitN(spec, "/", 2)
	if len(parts) == 0 || parts[0] == "" {
		return "", nil, false
	}
	pkg = parts[0]
	if len(parts) == 2 {
		s := parts[1]
		sub = &s
	}
	return pkg, sub, true
}

func ParseBareIdentifier(spec string) (string, *string, error) {
	var (
		pkg string
		sub *string
		ok  bool
	)
	if strings.HasPrefix(spec, "@") {
		pkg, sub, ok = parseScopedPackageName(spec)
	} else {
		pkg, sub, ok = parseGlobalPackageName(spec)
	}
	if !ok {
		return "", nil, &FailedManifestHydrationError{
			Message:      "Invalid specifier",
			ManifestPath: spec,
		}
	}
	return pkg, sub, nil
}

func FindClosestPNPManifestPath(start string) (string, bool) {
	dir := filepath.Clean(start)

	for {
		candidate := filepath.Join(dir, ".pnp.cjs")
		if st, err := os.Stat(candidate); err == nil && !st.IsDir() {
			return candidate, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", false
}

var rePNP = regexp.MustCompile(`(?s)(const[\ \r\n]+RAW_RUNTIME_STATE[\ \r\n]*=[\ \r\n]*|hydrateRuntimeState\(JSON\.parse\()'`)

func InitPNPManifest(m *Manifest, manifestPath string) error {
	m.ManifestPath = manifestPath
	m.ManifestDir = filepath.Dir(manifestPath)

	m.LocationTrie = *utils.NewTrie[PackageLocator]()

	for name, ranges := range m.PackageRegistryData {
		for reference, info := range ranges {
			joined := filepath.Join(m.ManifestDir, info.PackageLocation)

			if strings.HasSuffix(info.PackageLocation, "/") {
				joined = joined + "/"
			}

			norm := utils.NormalizePath(joined)

			info.PackageLocation = norm
			ranges[reference] = info

			if !info.DiscardFromLookup {
				m.LocationTrie.Insert(
					info.PackageLocation,
					PackageLocator{Name: name, Reference: reference},
				)
			}
		}
	}

	ranges, ok := m.PackageRegistryData[""]
	if !ok {
		return errors.New("assertion failed: should have a top-level name key")
	}
	top, ok := ranges[""]
	if !ok {
		return errors.New("assertion failed: should have a top-level range key")
	}

	if m.FallbackPool == nil {
		m.FallbackPool = make(FallbackPool)
	}
	for depName, dep := range top.PackageDependencies {
		if _, exists := m.FallbackPool[depName]; !exists {
			m.FallbackPool[depName] = dep
		}
	}

	return nil
}

func LoadPNPManifest(p string) (Manifest, error) {
	content, err := os.ReadFile(p)
	if err != nil {
		return Manifest{}, &FailedManifestHydrationError{
			Message:      "We failed to read the content of the manifest.",
			Err:          err,
			ManifestPath: p,
		}
	}

	loc := rePNP.FindIndex(content)
	if loc == nil {
		return Manifest{}, &FailedManifestHydrationError{
			Message:      "We failed to locate the PnP data payload inside its manifest file. Did you manually edit the file?",
			ManifestPath: p,
		}
	}

	i := loc[1]
	escaped := false
	jsonBuf := make([]byte, 0)

	for ; i < len(content); i++ {
		c := content[i]

		if c == '\'' && !escaped {
			break
		} else if c == '\\' && !escaped {
			escaped = true
		} else {
			escaped = false
			jsonBuf = append(jsonBuf, c)
		}
	}

	var manifest Manifest
	if err = json.Unmarshal(jsonBuf, &manifest); err != nil {
		return Manifest{}, &FailedManifestHydrationError{
			Message:      "We failed to parse the PnP data payload as proper JSON; Did you manually edit the file?",
			ManifestPath: p,
			Err:          err,
		}
	}

	err = InitPNPManifest(&manifest, p)
	if err != nil {
		return Manifest{}, &FailedManifestHydrationError{
			Message:      "We failed to init the PnP manifest",
			ManifestPath: p,
			Err:          err,
		}
	}
	return manifest, nil
}

func FindPNPManifest(parent string) (*Manifest, error) {
	path, ok := FindClosestPNPManifestPath(parent)
	if !ok {
		return nil, nil
	}
	manifest, err := LoadPNPManifest(path)
	if err != nil {
		return nil, err
	}
	return &manifest, nil
}

func IsDependencyTreeRoot(m *Manifest, loc *PackageLocator) bool {
	return slices.Contains(m.DependencyTreeRoots, PackageLocator{
		Name:      loc.Name,
		Reference: loc.Reference,
	})
}

func FindLocator(manifest *Manifest, path string) *PackageLocator {
	rel, err := filepath.Rel(manifest.ManifestDir, path)
	if err != nil {
		panic("Assertion failed: Provided path should be absolute but received" + path)
	}

	if manifest.IgnorePatternData != nil {
		re, err := manifest.IgnorePatternData.compile()
		if err == nil {
			if re.MatchString(utils.NormalizePath(rel)) {
				return nil
			}
		}
	}

	normFull := utils.NormalizePath(path)
	if v, ok := manifest.LocationTrie.GetAncestorValue(normFull); ok && v != nil {
		return v
	}
	return nil
}

func GetPackage(manifest *Manifest, locator *PackageLocator) (*PackageInformation, error) {
	refs, ok := manifest.PackageRegistryData[locator.Name]
	if !ok {
		return nil, fmt.Errorf("should have an entry in the package registry for %s", locator.Name)
	}
	info, ok := refs[locator.Reference]
	if !ok {
		return nil, fmt.Errorf("should have an entry in the package registry for %s", locator.Reference)
	}
	return &info, nil
}

func IsExcludedFromFallback(manifest *Manifest, locator *PackageLocator) bool {
	if refs, ok := manifest.FallbackExclusionList[locator.Name]; ok {
		for _, r := range refs {
			if r == locator.Reference {
				return true
			}
		}
	}
	return false
}

func FindBrokenPeerDependencies(specifier string, parent *PackageLocator) []PackageLocator {
	return []PackageLocator{}
}

func viaSuffix(specifier string, ident string) string {
	if ident != specifier {
		return fmt.Sprintf(" (via \"%s\")", ident)
	}
	return ""
}

func ResolveToUnqualifiedViaManifest(
	manifest *Manifest,
	specifier string,
	parentPath string,
) (Resolution, error) {
	ident, modulePath, err := ParseBareIdentifier(specifier)
	if err != nil {
		return Resolution{}, err
	}

	parentLocator := FindLocator(manifest, parentPath)
	if parentLocator == nil {
		return Resolution{Kind: ResolutionSkipped}, nil
	}

	parentPkg, err := GetPackage(manifest, parentLocator)
	if err != nil {
		return Resolution{}, err
	}

	var refOrAlias *PackageDependency

	if dep, ok := parentPkg.PackageDependencies[ident]; ok && dep != nil {
		refOrAlias = dep
	}

	if refOrAlias == nil && manifest.EnableTopLevelFallback && !IsExcludedFromFallback(manifest, parentLocator) {
		if dep, ok := manifest.FallbackPool[ident]; ok {
			refOrAlias = dep
		}
	}

	if refOrAlias == nil {
		if IsNodeJSBuiltin(specifier) {
			if IsDependencyTreeRoot(manifest, parentLocator) {
				msg := fmt.Sprintf(
					"Your application tried to access %s. While this module is usually interpreted as a Node builtin, your resolver is running inside a non-Node resolution context where such builtins are ignored. Since %s isn't otherwise declared in your dependencies, this makes the require call ambiguous and unsound.\n\nRequired package: %s%s\nRequired by: %s",
					ident, ident, ident, viaSuffix(specifier, ident), parentPath,
				)
				return Resolution{}, &UndeclaredDependencyError{
					Message:        msg,
					Request:        specifier,
					DependencyName: ident,
					IssuerLocator:  *parentLocator,
					IssuerPath:     parentPath,
				}
			}
			msg := fmt.Sprintf(
				"%s tried to access %s. While this module is usually interpreted as a Node builtin, your resolver is running inside a non-Node resolution context where such builtins are ignored. Since %s isn't otherwise declared in %s's dependencies, this makes the require call ambiguous and unsound.\n\nRequired package: %s%s\nRequired by: %s",
				parentLocator.Name, ident, ident, parentLocator.Name, ident, viaSuffix(specifier, ident), parentPath,
			)
			return Resolution{}, &UndeclaredDependencyError{
				Message:        msg,
				Request:        specifier,
				DependencyName: ident,
				IssuerLocator:  *parentLocator,
				IssuerPath:     parentPath,
			}
		}

		if IsDependencyTreeRoot(manifest, parentLocator) {
			msg := fmt.Sprintf(
				"Your application tried to access %s, but it isn't declared in your dependencies; this makes the require call ambiguous and unsound.\n\nRequired package: %s%s\nRequired by: %s",
				ident, ident, viaSuffix(specifier, ident), parentPath,
			)
			return Resolution{}, &UndeclaredDependencyError{
				Message:        msg,
				Request:        specifier,
				DependencyName: ident,
				IssuerLocator:  *parentLocator,
				IssuerPath:     parentPath,
			}
		}

		brokenAncestors := FindBrokenPeerDependencies(specifier, parentLocator)
		allBrokenAreRoots := len(brokenAncestors) > 0
		if allBrokenAreRoots {
			for _, brokenAncestor := range brokenAncestors {
				if !IsDependencyTreeRoot(manifest, &brokenAncestor) {
					allBrokenAreRoots = false
					break
				}
			}
		}

		var msg string
		if len(brokenAncestors) > 0 && allBrokenAreRoots {
			msg = fmt.Sprintf(
				"%s tried to access %s (a peer dependency) but it isn't provided by your application; this makes the require call ambiguous and unsound.\n\nRequired package: %s%s\nRequired by: %s@%s (via %s)",
				parentLocator.Name, ident, ident, viaSuffix(specifier, ident),
				parentLocator.Name, parentLocator.Reference, parentPath,
			)
		} else {
			msg = fmt.Sprintf(
				"%s tried to access %s (a peer dependency) but it isn't provided by its ancestors; this makes the require call ambiguous and unsound.\n\nRequired package: %s%s\nRequired by: %s@%s (via %s)",
				parentLocator.Name, ident, ident, viaSuffix(specifier, ident),
				parentLocator.Name, parentLocator.Reference, parentPath,
			)
		}

		return Resolution{}, &MissingPeerDependencyError{
			Message:         msg,
			Request:         specifier,
			DependencyName:  ident,
			IssuerLocator:   *parentLocator,
			IssuerPath:      parentPath,
			BrokenAncestors: brokenAncestors,
		}
	}

	var target PackageLocator
	if refOrAlias.IsAlias {
		target = PackageLocator{
			Name:      refOrAlias.Alias[0],
			Reference: refOrAlias.Alias[1],
		}
	} else {
		target = PackageLocator{
			Name:      ident,
			Reference: refOrAlias.Reference,
		}
	}
	depPkg, err := GetPackage(manifest, &target)
	if err != nil {
		return Resolution{}, err
	}

	return Resolution{
		Kind:       ResolutionResolved,
		Path:       depPkg.PackageLocation,
		ModulePath: modulePath,
	}, nil
}

func ResolveToUnqualified(specifier, parentPath string, cfg *ResolutionConfig) (Resolution, error) {
	if cfg == nil || cfg.Host.FindPNPManifest == nil {
		return Resolution{}, errors.New("no host configured")
	}
	m, err := cfg.Host.FindPNPManifest(parentPath)
	if err != nil {
		return Resolution{}, err
	}
	if m == nil {
		return Resolution{Kind: ResolutionSkipped}, nil
	}
	return ResolveToUnqualifiedViaManifest(m, specifier, parentPath)
}
