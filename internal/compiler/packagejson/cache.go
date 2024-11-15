package packagejson

import (
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler/diagnostics"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/semver"
	"github.com/microsoft/typescript-go/internal/tspath"
)

var typeScriptVersion = semver.MustParse(core.Version)

type PackageJson struct {
	Fields
	versionPaths *VersionPaths
}

func (p *PackageJson) GetVersionPaths(collectTraces bool) (value *VersionPaths, traces []string) {
	if p.versionPaths == nil {
		if p.Fields.TypesVersions.Type == JSONValueTypeNotPresent {
			p.versionPaths = &VersionPaths{}
			if collectTraces {
				traces = append(traces, diagnostics.Format(diagnostics.X_package_json_does_not_have_a_0_field, "typesVersions"))
			}
			return nil, traces
		}
		if p.Fields.TypesVersions.Type != JSONValueTypeObject {
			p.versionPaths = &VersionPaths{}
			if collectTraces {
				traces = append(traces, diagnostics.Format(diagnostics.Expected_type_of_0_field_in_package_json_to_be_1_got_2, "typesVersions", "object", p.Fields.TypesVersions.Type.String()))
			}
			return nil, traces
		}

		if collectTraces {
			traces = append(traces, diagnostics.Format(diagnostics.X_package_json_has_a_typesVersions_field_with_version_specific_path_mappings, "typesVersions"))
		}

		for key, value := range p.Fields.TypesVersions.AsObject().Entries() {
			keyRange, ok := semver.TryParseVersionRange(key)
			if ok {
				if collectTraces {
					traces = append(traces, diagnostics.Format(diagnostics.X_package_json_has_a_typesVersions_entry_0_that_is_not_a_valid_semver_range, key))
				}
				continue
			}
			if keyRange.Test(&typeScriptVersion) {
				if value.Type != JSONValueTypeObject {
					if collectTraces {
						traces = append(traces, diagnostics.Format(diagnostics.Expected_type_of_0_field_in_package_json_to_be_1_got_2, "typesVersions['"+key+"']", "object", value.Type.String()))
					}
					p.versionPaths = &VersionPaths{}
					return nil, traces
				}
				p.versionPaths = &VersionPaths{
					Version:   key,
					pathsJSON: value.AsObject(),
				}
				return p.versionPaths, traces
			}
		}

		p.versionPaths = &VersionPaths{}
		if collectTraces {
			traces = append(traces, diagnostics.Format(diagnostics.X_package_json_does_not_have_a_typesVersions_entry_that_matches_version_0, core.VersionMajorMinor))
		}
		return nil, traces
	}
	if p.versionPaths.IsValid() {
		return p.versionPaths, traces
	}
	return nil, traces
}

type VersionPaths struct {
	Version   string
	pathsJSON *collections.OrderedMap[string, *JSONValue]
	paths     map[string][]string
}

func (v *VersionPaths) IsValid() bool {
	return v.pathsJSON != nil
}

func (v *VersionPaths) GetPaths() map[string][]string {
	if !v.IsValid() {
		return nil
	}
	if v.paths != nil {
		return v.paths
	}
	v.paths = make(map[string][]string, v.pathsJSON.Size())
	for key, value := range v.pathsJSON.Entries() {
		if value.Type != JSONValueTypeArray {
			continue
		}
		v.paths[key] = make([]string, len(value.AsArray()))
		for i, path := range value.AsArray() {
			if path.Type != JSONValueTypeString {
				continue
			}
			v.paths[key][i] = path.Value.(string)
		}
	}
	return v.paths
}

type InfoCacheEntry struct {
	PackageDirectory string
	DirectoryExists  bool
	Contents         *PackageJson
}

func (p *InfoCacheEntry) Exists() bool {
	return p != nil && p.Contents != nil
}

type InfoCache struct {
	IsReadonly                bool
	cache                     map[tspath.Path]InfoCacheEntry
	currentDirectory          string
	useCaseSensitiveFileNames bool
}

func NewInfoCache(currentDirectory string, useCaseSensitiveFileNames bool) *InfoCache {
	return &InfoCache{
		currentDirectory:          currentDirectory,
		useCaseSensitiveFileNames: useCaseSensitiveFileNames,
	}
}

func (p *InfoCache) Get(packageJsonPath string) *InfoCacheEntry {
	key := tspath.ToPath(packageJsonPath, p.currentDirectory, p.useCaseSensitiveFileNames)
	if p.cache == nil {
		p.cache = make(map[tspath.Path]InfoCacheEntry)
		return nil
	}
	entry, ok := p.cache[key]
	if !ok {
		return nil
	}
	return &entry
}

func (p *InfoCache) Set(packageJsonPath string, info *InfoCacheEntry) {
	key := tspath.ToPath(packageJsonPath, p.currentDirectory, p.useCaseSensitiveFileNames)
	if p.cache == nil {
		p.cache = make(map[tspath.Path]InfoCacheEntry)
	}
	p.cache[key] = *info
}
