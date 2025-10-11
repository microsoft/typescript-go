package pnp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

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

type UndeclaredDependency struct {
	Message        string
	Request        string
	DependencyName string
	IssuerLocator  PackageLocator
	IssuerPath     string
}

func (e *UndeclaredDependency) Error() string { return e.Message }

type MissingPeerDependency struct {
	Message         string
	Request         string
	DependencyName  string
	IssuerLocator   PackageLocator
	IssuerPath      string
	BrokenAncestors []PackageLocator
}

func (e *MissingPeerDependency) Error() string { return e.Message }

type FailedManifestHydration struct {
	Message      string
	ManifestPath string
	Err          error
}

func (e *FailedManifestHydration) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s\n\nOriginal error: %v", e.Message, e.Err)
	}
	return e.Message
}
func (e *FailedManifestHydration) Unwrap() error { return e.Err }

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
		return "", nil, &FailedManifestHydration{
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
	if abs, err := filepath.Abs(manifestPath); err == nil {
		m.ManifestPath = abs
		m.ManifestDir = filepath.Dir(abs)
	} else {
		return fmt.Errorf("assertion failed: %w", err)
	}

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
		return fmt.Errorf("assertion failed: should have a top-level name key")
	}
	top, ok := ranges[""]
	if !ok {
		return fmt.Errorf("assertion failed: should have a top-level range key")
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
		return Manifest{}, &FailedManifestHydration{
			Message: "We failed to read the content of the manifest.",
			Err:     err,
			ManifestPath: func() string {
				if abs, e := filepath.Abs(p); e == nil {
					return abs
				}
				return p
			}(),
		}
	}

	loc := rePNP.FindIndex(content)
	if loc == nil {
		return Manifest{}, &FailedManifestHydration{
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
			i++
			break
		} else if c == '\\' && !escaped {
			escaped = true
		} else {
			escaped = false
			jsonBuf = append(jsonBuf, c)
		}
	}

	var manifest Manifest
	if err := json.Unmarshal(jsonBuf, &manifest); err != nil {
		return Manifest{}, &FailedManifestHydration{
			Message:      "We failed to parse the PnP data payload as proper JSON; Did you manually edit the file?",
			ManifestPath: p,
			Err:          err,
		}
	}

	InitPNPManifest(&manifest, p)
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

func FindLocator(manifest *Manifest, absPath string) *PackageLocator {
	rel, err := filepath.Rel(manifest.ManifestDir, absPath)
	if err != nil {
		panic(fmt.Sprintf("Assertion failed: Provided path should be absolute but received %s", absPath))
	}

	if manifest.IgnorePatternData != nil {
		re, err := manifest.IgnorePatternData.compile()
		if err == nil {
			if re.MatchString(utils.NormalizePath(rel)) {
				return nil
			}
		}
	}

	normFull := utils.NormalizePath(absPath)
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
				return Resolution{}, &UndeclaredDependency{
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
			return Resolution{}, &UndeclaredDependency{
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
			return Resolution{}, &UndeclaredDependency{
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

		return Resolution{}, &MissingPeerDependency{
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
		return Resolution{}, fmt.Errorf("no host configured")
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
