package module

import (
	"fmt"
	"math/bits"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type ResolutionHost interface {
	FS() vfs.FS
	GetCurrentDirectory() string
	Trace(msg string)
}

type ModeAwareCacheKey struct {
	name string
	mode core.ResolutionMode
}

type ResolvedProjectReference struct {
	CommandLine core.ParsedOptions
	SourceFile  *ast.SourceFile
	References  []*ResolvedProjectReference
}

type NodeResolutionFeatures int32

const (
	NodeResolutionFeaturesNone    NodeResolutionFeatures = 0
	NodeResolutionFeaturesImports NodeResolutionFeatures = 1 << iota
	NodeResolutionFeaturesSelfName
	NodeResolutionFeaturesExports
	NodeResolutionFeaturesExportsPatternTrailers

	NodeResolutionFeaturesAll             = NodeResolutionFeaturesImports | NodeResolutionFeaturesSelfName | NodeResolutionFeaturesExports | NodeResolutionFeaturesExportsPatternTrailers
	NodeResolutionFeaturesNode16Default   = NodeResolutionFeaturesImports | NodeResolutionFeaturesSelfName | NodeResolutionFeaturesExports | NodeResolutionFeaturesExportsPatternTrailers
	NodeResolutionFeaturesNodeNextDefault = NodeResolutionFeaturesAll
	NodeResolutionFeaturesBundlerDefault  = NodeResolutionFeaturesImports | NodeResolutionFeaturesSelfName | NodeResolutionFeaturesExports | NodeResolutionFeaturesExportsPatternTrailers
)

type PackageId struct {
	Name             string
	SubModuleName    string
	Version          string
	PeerDependencies string
}

func (p *PackageId) String() string {
	return fmt.Sprintf("%s@%s%s", p.Name, p.Version, p.PeerDependencies)
}

func (p *PackageId) PackageName() string {
	if p.SubModuleName != "" {
		return p.Name + "/" + p.SubModuleName
	}
	return p.Name
}

type LookupLocations struct {
	FailedLookupLocations []string
	AffectingLocations    []string
	ResolutionDiagnostics []ast.Diagnostic
}

type ResolvedModule struct {
	ResolvedFileName         string
	OriginalPath             string
	Extension                string
	ResolvedUsingTsExtension bool
	PackageId                PackageId
	IsExternalLibraryImport  bool
}

func (r *ResolvedModule) IsResolved() bool {
	return r.ResolvedFileName != ""
}

type ResolvedTypeReferenceDirective struct {
	Primary                 bool
	ResolvedFileName        string
	OriginalPath            string
	PackageId               PackageId
	IsExternalLibraryImport bool
}

func (r *ResolvedTypeReferenceDirective) IsResolved() bool {
	return r.ResolvedFileName != ""
}

type Extensions int32

const (
	extensionsTypeScript Extensions = 1 << iota
	extensionsJavaScript
	extensionsDeclaration
	extensionsJson

	extensionsImplementationFiles = extensionsTypeScript | extensionsJavaScript
)

func (e Extensions) String() string {
	result := make([]string, 0, bits.OnesCount(uint(e)))
	if e&extensionsTypeScript != 0 {
		result = append(result, "TypeScript")
	}
	if e&extensionsJavaScript != 0 {
		result = append(result, "JavaScript")
	}
	if e&extensionsDeclaration != 0 {
		result = append(result, "Declaration")
	}
	if e&extensionsJson != 0 {
		result = append(result, "JSON")
	}
	return strings.Join(result, ", ")
}

func (e Extensions) Array() []string {
	result := []string{}
	if e&extensionsTypeScript != 0 {
		result = append(result, tspath.ExtensionTs, tspath.ExtensionTsx)
	}
	if e&extensionsJavaScript != 0 {
		result = append(result, tspath.ExtensionJs, tspath.ExtensionJsx)
	}
	if e&extensionsDeclaration != 0 {
		result = append(result, tspath.ExtensionDts)
	}
	if e&extensionsJson != 0 {
		result = append(result, tspath.ExtensionJson)
	}
	return result
}

// ResolvedModule with an explicitly provided `extension` property.
// Prefer this over `ResolvedModule`.
// If changing this, remember to change `moduleResolutionIsEqualTo`.
type resolvedModuleFull struct {
	ResolvedModule
	// This is a file name with preserved original casing, not a normalized `Path`
	originalPath string
	// Extension of resolvedFileName. This must match what's at the end of resolvedFileName.
	// This is optional for backwards-compatibility, but will be added if not provided.
	extension string
	packageId PackageId
}

type ResolvedModuleWithFailedLookupLocations struct {
	ResolvedModule        *resolvedModuleFull
	failedLookupLocations []string
	affectingLocations    []string
	resolutionDiagnostics []ast.Diagnostic
	// Used to issue a better diagnostic when an unresolvable module may
	// have been resolvable under different module resolution settings.
	alternateResult string
}
