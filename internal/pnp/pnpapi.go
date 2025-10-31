package pnp

/*
 * Yarn Plug'n'Play (generally referred to as Yarn PnP) is the default installation strategy in modern releases of Yarn.
 * Yarn PnP generates a single Node.js loader file in place of the typical node_modules folder.
 * This loader file, named .pnp.cjs, contains all information about your project's dependency tree, informing your tools as to
 * the location of the packages on the disk and letting them know how to resolve require and import calls.
 *
 * The full specification is available at https://yarnpkg.com/advanced/pnp-spec
 */

import (
	"errors"
	"fmt"
	"strings"

	"github.com/microsoft/typescript-go/internal/tspath"
)

type PnpApi struct {
	fs       PnpApiFS
	url      string
	manifest *PnpManifestData
}

// FS abstraction used by the PnpApi to access the file system
// We can't use the vfs.FS interface because it creates an import cycle: core -> pnp -> vfs -> core
type PnpApiFS interface {
	FileExists(path string) bool
	ReadFile(path string) (contents string, ok bool)
}

func (p *PnpApi) RefreshManifest() error {
	var newData *PnpManifestData
	var err error

	if p.manifest == nil {
		newData, err = p.findClosestPnpManifest()
	} else {
		newData, err = parseManifestFromPath(p.fs, p.manifest.dirPath)
	}

	if err != nil {
		return err
	}

	p.manifest = newData
	return nil
}

func (p *PnpApi) ResolveToUnqualified(specifier string, parentPath string) (string, error) {
	if p.manifest == nil {
		panic("ResolveToUnqualified called with no PnP manifest available")
	}

	ident, modulePath, err := p.ParseBareIdentifier(specifier)
	if err != nil {
		// Skipping resolution
		return "", nil
	}

	parentLocator, err := p.FindLocator(parentPath)
	if err != nil || parentLocator == nil {
		// Skipping resolution
		return "", nil
	}

	parentPkg := p.GetPackage(parentLocator)

	var referenceOrAlias *PackageDependency
	for _, dep := range parentPkg.PackageDependencies {
		if dep.Ident == ident {
			referenceOrAlias = &dep
			break
		}
	}

	// If not found, try fallback if enabled
	if referenceOrAlias == nil {
		if p.manifest.enableTopLevelFallback {
			excluded := false
			if exclusion, ok := p.manifest.fallbackExclusionMap[parentLocator.Name]; ok {
				for _, entry := range exclusion.Entries {
					if entry == parentLocator.Reference {
						excluded = true
						break
					}
				}
			}
			if !excluded {
				fallback := p.ResolveViaFallback(ident)
				if fallback != nil {
					referenceOrAlias = fallback
				}
			}
		}
	}

	// undeclared dependency
	if referenceOrAlias == nil {
		if parentLocator.Name == "" {
			return "", fmt.Errorf("Your application tried to access %s, but it isn't declared in your dependencies; this makes the require call ambiguous and unsound.\n\nRequired package: %s\nRequired by: %s", ident, ident, parentPath)
		}
		return "", fmt.Errorf("%s tried to access %s, but it isn't declared in your dependencies; this makes the require call ambiguous and unsound.\n\nRequired package: %s\nRequired by: %s", parentLocator.Name, ident, ident, parentPath)
	}

	// unfulfilled peer dependency
	if !referenceOrAlias.IsAlias() && referenceOrAlias.Reference == "" {
		if parentLocator.Name == "" {
			return "", fmt.Errorf("Your application tried to access %s (a peer dependency); this isn't allowed as there is no ancestor to satisfy the requirement. Use a devDependency if needed.\n\nRequired package: %s\nRequired by: %s", ident, ident, parentPath)
		}
		return "", fmt.Errorf("%s tried to access %s (a peer dependency) but it isn't provided by its ancestors/your application; this makes the require call ambiguous and unsound.\n\nRequired package: %s\nRequired by: %s", parentLocator.Name, ident, ident, parentPath)
	}

	var dependencyPkg *PackageInfo
	if referenceOrAlias.IsAlias() {
		dependencyPkg = p.GetPackage(&Locator{Name: referenceOrAlias.AliasName, Reference: referenceOrAlias.Reference})
	} else {
		dependencyPkg = p.GetPackage(&Locator{Name: referenceOrAlias.Ident, Reference: referenceOrAlias.Reference})
	}

	return tspath.ResolvePath(p.manifest.dirPath, dependencyPkg.PackageLocation, modulePath), nil
}

func (p *PnpApi) findClosestPnpManifest() (*PnpManifestData, error) {
	directoryPath := tspath.GetNormalizedAbsolutePath(p.url, "/")

	for {
		pnpPath := tspath.CombinePaths(directoryPath, ".pnp.cjs")
		if p.fs.FileExists(pnpPath) {
			return parseManifestFromPath(p.fs, directoryPath)
		}

		if tspath.IsDiskPathRoot(directoryPath) {
			return nil, errors.New("no PnP manifest found")
		}

		directoryPath = tspath.GetDirectoryPath(directoryPath)
	}
}

func (p *PnpApi) GetPackage(locator *Locator) *PackageInfo {
	packageRegistryMap := p.manifest.packageRegistryMap
	packageInfo, ok := packageRegistryMap[locator.Name][locator.Reference]
	if !ok {
		panic(locator.Name + " should have an entry in the package registry")
	}

	return packageInfo
}

func (p *PnpApi) FindLocator(parentPath string) (*Locator, error) {
	relativePath := tspath.GetRelativePathFromDirectory(p.manifest.dirPath, parentPath,
		tspath.ComparePathsOptions{UseCaseSensitiveFileNames: true})

	if p.manifest.ignorePatternData != nil {
		match, err := p.manifest.ignorePatternData.MatchString(relativePath)
		if err != nil {
			return nil, err
		}

		if match {
			return nil, nil
		}
	}

	var relativePathWithDot string
	if strings.HasPrefix(relativePath, "../") {
		relativePathWithDot = relativePath
	} else {
		relativePathWithDot = "./" + relativePath
	}

	var bestLength int
	var bestLocator *Locator
	pathSegments := strings.Split(relativePathWithDot, "/")
	currentTrie := p.manifest.packageRegistryTrie

	// Go down the trie, looking for the latest defined packageInfo that matches the path
	for index, segment := range pathSegments {
		currentTrie = currentTrie.childrenPathSegments[segment]

		if currentTrie == nil || currentTrie.childrenPathSegments == nil {
			break
		}

		if currentTrie.packageData != nil && index >= bestLength {
			bestLength = index
			bestLocator = &Locator{Name: currentTrie.packageData.ident, Reference: currentTrie.packageData.reference}
		}
	}

	if bestLocator == nil {
		return nil, fmt.Errorf("no package found for path %s", relativePath)
	}

	return bestLocator, nil
}

func (p *PnpApi) ResolveViaFallback(name string) *PackageDependency {
	topLevelPkg := p.GetPackage(&Locator{Name: "", Reference: ""})

	if topLevelPkg != nil {
		for _, dep := range topLevelPkg.PackageDependencies {
			if dep.Ident == name {
				return &dep
			}
		}
	}

	for _, dep := range p.manifest.fallbackPool {
		if dep[0] == name {
			return &PackageDependency{
				Ident:     dep[0],
				Reference: dep[1],
				AliasName: "",
			}
		}
	}

	return nil
}

func (p *PnpApi) ParseBareIdentifier(specifier string) (ident string, modulePath string, err error) {
	if len(specifier) == 0 {
		return "", "", fmt.Errorf("Empty specifier: %s", specifier)
	}

	firstSlash := strings.Index(specifier, "/")

	if specifier[0] == '@' {
		if firstSlash == -1 {
			return "", "", fmt.Errorf("Invalid specifier: %s", specifier)
		}

		secondSlash := strings.Index(specifier[firstSlash+1:], "/")

		if secondSlash == -1 {
			ident = specifier
		} else {
			ident = specifier[:firstSlash+1+secondSlash]
		}
	} else {
		firstSlash := strings.Index(specifier, "/")

		if firstSlash == -1 {
			ident = specifier
		} else {
			ident = specifier[:firstSlash]
		}
	}

	modulePath = specifier[len(ident):]

	return ident, modulePath, nil
}

func (p *PnpApi) GetPnpTypeRoots(currentDirectory string) []string {
	if p.manifest == nil {
		return []string{}
	}

	currentDirectory = tspath.NormalizePath(currentDirectory)

	currentPackage, err := p.FindLocator(currentDirectory)
	if err != nil {
		return []string{}
	}

	if currentPackage == nil {
		return []string{}
	}

	packageDependencies := p.GetPackage(currentPackage).PackageDependencies

	typeRoots := []string{}
	for _, dep := range packageDependencies {
		if strings.HasPrefix(dep.Ident, "@types/") && dep.Reference != "" {
			packageInfo := p.GetPackage(&Locator{Name: dep.Ident, Reference: dep.Reference})
			typeRoots = append(typeRoots, tspath.GetDirectoryPath(
				tspath.ResolvePath(p.manifest.dirPath, packageInfo.PackageLocation),
			))
		}
	}

	return typeRoots
}

func (p *PnpApi) IsImportable(fromFileName string, toFileName string) bool {
	fromLocator, errFromLocator := p.FindLocator(fromFileName)
	toLocator, errToLocator := p.FindLocator(toFileName)

	if fromLocator == nil || toLocator == nil || errFromLocator != nil || errToLocator != nil {
		return false
	}

	fromInfo := p.GetPackage(fromLocator)
	for _, dep := range fromInfo.PackageDependencies {
		if dep.Reference == toLocator.Reference {
			if dep.IsAlias() && dep.AliasName == toLocator.Name {
				return true
			}

			if dep.Ident == toLocator.Name {
				return true
			}
		}
	}

	return false
}

func (p *PnpApi) GetPackageLocationAbsolutePath(packageInfo *PackageInfo) string {
	if packageInfo == nil {
		return ""
	}

	packageLocation := packageInfo.PackageLocation
	return tspath.ResolvePath(p.manifest.dirPath, packageLocation)
}

func (p *PnpApi) IsInPnpModule(fromFileName string, toFileName string) bool {
	fromLocator, _ := p.FindLocator(fromFileName)
	toLocator, _ := p.FindLocator(toFileName)
	// The targeted filename is in a pnp module different from the requesting filename
	return fromLocator != nil && toLocator != nil && fromLocator.Name != toLocator.Name
}

func (p *PnpApi) AppendPnpTypeRoots(nmTypes []string, baseDir string, nmFromConfig bool) ([]string, bool) {
	pnpTypes := p.GetPnpTypeRoots(baseDir)

	if len(nmTypes) > 0 {
		return append(nmTypes, pnpTypes...), nmFromConfig
	}

	if len(pnpTypes) > 0 {
		return pnpTypes, false
	}

	return nil, false
}
