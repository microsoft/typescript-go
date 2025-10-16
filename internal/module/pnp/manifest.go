package pnp

import (
	"bytes"
	"errors"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"

	"github.com/microsoft/typescript-go/internal/tspath"
)

type LocationTrie[T any] struct {
	inner *Trie[T]
}

type RegexDef struct {
	Source string `json:"source"`
	reg    *regexp2.Regexp
}

type Manifest struct {
	ManifestDir  string                       `json:"-"`
	ManifestPath string                       `json:"-"`
	LocationTrie LocationTrie[PackageLocator] `json:"-"`

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
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	r.Source = s
	reg, err := regexp2.Compile(s, regexp2.ECMAScript)
	if err != nil {
		return err
	}
	r.reg = reg
	return nil
}

func (m *FallbackPool) UnmarshalJSON(data []byte) error {
	var items []jsontext.Value
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}
	res := make(FallbackPool, len(items))
	for _, it := range items {
		var pair []jsontext.Value
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
	var items []jsontext.Value
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}
	res := make(FallbackExclusionList, len(items))
	for _, it := range items {
		var pair []jsontext.Value
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
	var items []jsontext.Value
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}
	res := make(PackageRegistryData, len(items))
	for _, it := range items {
		var pair []jsontext.Value
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

		var inner []jsontext.Value
		if err := json.Unmarshal(pair[1], &inner); err != nil {
			return err
		}
		innerMap := make(map[string]PackageInformation, len(inner))
		for _, e := range inner {
			var p2 []jsontext.Value
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

func NewLocationTrie[T any]() *LocationTrie[T] {
	return &LocationTrie[T]{inner: New[T]()}
}

func (t *LocationTrie[T]) key(key string) string {
	p := tspath.NormalizePath(key)

	if !strings.HasSuffix(p, "/") {
		return p + "/"
	}

	return p
}

func (t *LocationTrie[T]) GetAncestorValue(p string) (*T, bool) {
	v, ok := t.inner.GetAncestorValue(t.key(p))
	return &v, ok
}

func (t *LocationTrie[T]) Insert(p string, v T) {
	t.inner.Set(t.key(p), v)
}
