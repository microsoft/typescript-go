package pnp

import (
	"bytes"
	"encoding/json"
	"errors"
	"regexp"

	"github.com/microsoft/typescript-go/internal/module/pnp/utils"
)

type Trie[T any] = utils.Trie[T]

type RegexDef struct {
	Source string `json:"source"`
	Flags  string `json:"flags,omitempty"`
}

type Manifest struct {
	ManifestDir  string               `json:"-"`
	ManifestPath string               `json:"-"`
	LocationTrie Trie[PackageLocator] `json:"-"`

	EnableTopLevelFallback bool      `json:"enableTopLevelFallback"`
	IgnorePatternData      *RegexDef `json:"ignorePatternData,omitempty"`

	DependencyTreeRoots []PackageLocator `json:"dependencyTreeRoots"`

	FallbackPool FallbackPool `json:"fallbackPool"`

	FallbackExclusionList FallbackExclusionList `json:"fallbackExclusionList"`

	PackageRegistryData PackageRegistryData `json:"packageRegistryData"`
}

type PackageLocator struct {
	Name      string `json:"name"`
	Reference string `json:"reference"`
}

type PackageInformation struct {
	PackageLocation string `json:"packageLocation"`

	DiscardFromLookup bool `json:"discardFromLookup"`

	PackageDependencies FallbackPool `json:"packageDependencies"`
}

type PackageDependency struct {
	Reference string
	Alias     [2]string
	IsAlias   bool
}

func (p *PackageDependency) UnmarshalJSON(data []byte) error {
	// case 1: string
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		p.Reference = s
		p.IsAlias = false
		return nil
	}

	// case 2: [string, string]
	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		if len(arr) != 2 {
			return errors.New("PackageDependency: array must have length 2")
		}
		p.IsAlias = true
		p.Alias = [2]string{arr[0], arr[1]}
		return nil
	}

	return errors.New("PackageDependency: unsupported JSON shape")
}

type FallbackPool map[string]*PackageDependency

func (r *RegexDef) UnmarshalJSON(b []byte) error {
	if bytes.Equal(bytes.TrimSpace(b), []byte("null")) {
		*r = RegexDef{}
		return nil
	}
	if len(b) > 0 && b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		r.Source = s
		r.Flags = ""
		return nil
	}
	type alias RegexDef
	var a alias
	if err := json.Unmarshal(b, &a); err != nil {
		return err
	}
	*r = RegexDef(a)
	return nil
}

func (m *FallbackPool) UnmarshalJSON(data []byte) error {
	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}
	res := make(FallbackPool, len(items))
	for _, it := range items {
		var pair []json.RawMessage
		if err := json.Unmarshal(it, &pair); err != nil {
			return err
		}
		if len(pair) != 2 {
			return errors.New("fallbackPool: each item must have 2 elements")
		}
		var key string
		if err := json.Unmarshal(pair[0], &key); err != nil {
			return err
		}

		raw := bytes.TrimSpace(pair[1])
		if bytes.Equal(raw, []byte("null")) {
			res[key] = nil
			continue
		}

		var dep PackageDependency
		if err := json.Unmarshal(pair[1], &dep); err != nil {
			return err
		}
		res[key] = &dep
	}
	*m = res
	return nil
}

type FallbackExclusionList map[string][]string

func (m *FallbackExclusionList) UnmarshalJSON(data []byte) error {
	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}
	res := make(FallbackExclusionList, len(items))
	for _, it := range items {
		var pair []json.RawMessage
		if err := json.Unmarshal(it, &pair); err != nil {
			return err
		}
		if len(pair) != 2 {
			return errors.New("fallbackExclusionList: each item must have 2 elements")
		}
		var key string
		if err := json.Unmarshal(pair[0], &key); err != nil {
			return err
		}
		var arr []string
		if err := json.Unmarshal(pair[1], &arr); err != nil {
			return err
		}
		res[key] = arr
	}
	*m = res
	return nil
}

type PackageRegistryData map[string]map[string]PackageInformation

func (m *PackageRegistryData) UnmarshalJSON(data []byte) error {
	var items []json.RawMessage
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}
	res := make(PackageRegistryData, len(items))
	for _, it := range items {
		var pair []json.RawMessage
		if err := json.Unmarshal(it, &pair); err != nil {
			return err
		}
		if len(pair) != 2 {
			return errors.New("packageRegistryData: each item must have 2 elements")
		}

		var key1 *string
		if err := json.Unmarshal(pair[0], &key1); err != nil {
			return err
		}
		k1 := ""
		if key1 != nil {
			k1 = *key1
		}

		var inner []json.RawMessage
		if err := json.Unmarshal(pair[1], &inner); err != nil {
			return err
		}
		innerMap := make(map[string]PackageInformation, len(inner))
		for _, e := range inner {
			var p2 []json.RawMessage
			if err := json.Unmarshal(e, &p2); err != nil {
				return err
			}
			if len(p2) != 2 {
				return errors.New("packageRegistryData: inner item must have 2 elements")
			}
			var key2 *string
			if err := json.Unmarshal(p2[0], &key2); err != nil {
				return err
			}
			k2 := ""
			if key2 != nil {
				k2 = *key2
			}
			var info PackageInformation
			if err := json.Unmarshal(p2[1], &info); err != nil {
				return err
			}
			innerMap[k2] = info
		}
		res[k1] = innerMap
	}
	*m = res
	return nil
}

func (r *RegexDef) compile() (*regexp.Regexp, error) {
	return regexp.Compile(r.Source)
}
