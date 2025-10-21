package pnp

import (
	"errors"
	"fmt"
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

type FallbackPool map[string]*PackageDependency

type PackageDependency struct {
	Reference string
	Alias     [2]string
	IsAlias   bool
}

type PackageRegistryData map[string]map[string]PackageInformation

var _ json.UnmarshalerFrom = (*Manifest)(nil)

func (m *Manifest) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	if _, err := dec.ReadToken(); err != nil {
		return err
	}
	for dec.PeekKind() != jsontext.EndObject.Kind() {
		fieldName, err := dec.ReadValue()
		if err != nil {
			return err
		}

		switch string(fieldName) {
		case `"enableTopLevelFallback"`:
			if err := json.UnmarshalDecode(dec, &m.EnableTopLevelFallback); err != nil {
				return err
			}
		case `"ignorePatternData"`:
			if err := json.UnmarshalDecode(dec, &m.IgnorePatternData); err != nil {
				return err
			}
		case `"dependencyTreeRoots"`:
			if err := json.UnmarshalDecode(dec, &m.DependencyTreeRoots); err != nil {
				return err
			}
		case `"fallbackPool"`:
			if err := json.UnmarshalDecode(dec, &m.FallbackPool); err != nil {
				return err
			}
		case `"fallbackExclusionList"`:
			if err := json.UnmarshalDecode(dec, &m.FallbackExclusionList); err != nil {
				fmt.Println("error", err)
				return err
			}
		case `"packageRegistryData"`:
			if err := json.UnmarshalDecode(dec, &m.PackageRegistryData); err != nil {
				return err
			}
		default:
			if _, err := dec.ReadValue(); err != nil {
				return err
			}
		}
	}
	if _, err := dec.ReadToken(); err != nil {
		return err
	}

	return nil
}

var _ json.UnmarshalerFrom = (*PackageDependency)(nil)

func (p *PackageDependency) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	switch dec.PeekKind() {
	case '"': // string case
		var s string
		if err := json.UnmarshalDecode(dec, &s); err != nil {
			return err
		}
		p.Reference = s
		p.IsAlias = false
		return nil
	case '[': // array case
		if _, err := dec.ReadToken(); err != nil {
			return err
		}
		var arr []string
		for dec.PeekKind() != jsontext.EndArray.Kind() {
			var s string
			if err := json.UnmarshalDecode(dec, &s); err != nil {
				return err
			}
			arr = append(arr, s)
		}
		if _, err := dec.ReadToken(); err != nil {
			return err
		}
		if len(arr) != 2 {
			return errors.New("PackageDependency: array must have length 2")
		}
		p.IsAlias = true
		p.Alias = [2]string{arr[0], arr[1]}
		return nil
	default:
		return errors.New("PackageDependency: unsupported JSON shape")
	}
}

var _ json.UnmarshalerFrom = (*RegexDef)(nil)

func (r *RegexDef) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	if dec.PeekKind() == jsontext.Null.Kind() {
		if _, err := dec.ReadToken(); err != nil {
			return err
		}
		*r = RegexDef{}
		return nil
	}

	var s string
	if err := json.UnmarshalDecode(dec, &s); err != nil {
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

var _ json.UnmarshalerFrom = (*FallbackPool)(nil)

func (f *FallbackPool) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	if _, err := dec.ReadToken(); err != nil {
		return err
	}

	res := make(FallbackPool)

	for dec.PeekKind() != jsontext.EndArray.Kind() {
		if _, err := dec.ReadToken(); err != nil {
			return err
		}

		var key string
		if err := json.UnmarshalDecode(dec, &key); err != nil {
			return err
		}

		if dec.PeekKind() == jsontext.Null.Kind() {
			if _, err := dec.ReadToken(); err != nil {
				return err
			}
			res[key] = nil
		} else {
			var dep PackageDependency
			if err := json.UnmarshalDecode(dec, &dep); err != nil {
				return err
			}
			res[key] = &dep
		}

		if _, err := dec.ReadToken(); err != nil {
			return err
		}
	}

	if _, err := dec.ReadToken(); err != nil {
		return err
	}

	*f = res
	return nil
}

type FallbackExclusionList map[string][]string

var _ json.UnmarshalerFrom = (*FallbackExclusionList)(nil)

func (f *FallbackExclusionList) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	// start of array
	if _, err := dec.ReadToken(); err != nil {
		return err
	}

	res := make(FallbackExclusionList)

	for dec.PeekKind() != jsontext.EndArray.Kind() {
		// start of array
		if _, err := dec.ReadToken(); err != nil {
			return err
		}

		var key string
		if err := json.UnmarshalDecode(dec, &key); err != nil {
			return err
		}

		var arr []string
		if err := json.UnmarshalDecode(dec, &arr); err != nil {
			return err
		}

		res[key] = arr

		// end of array
		if _, err := dec.ReadToken(); err != nil {
			return err
		}
	}

	// end of array
	if _, err := dec.ReadToken(); err != nil {
		return err
	}

	*f = res
	return nil
}

var _ json.UnmarshalerFrom = (*PackageRegistryData)(nil)

func (m *PackageRegistryData) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	if _, err := dec.ReadToken(); err != nil {
		return err
	}
	res := make(PackageRegistryData)

	for dec.PeekKind() != jsontext.EndArray.Kind() {
		// start of array1
		if _, err := dec.ReadToken(); err != nil {
			return err
		}
		var key1 *string
		if err := json.UnmarshalDecode(dec, &key1); err != nil {
			return err
		}
		k1 := ""
		if key1 != nil {
			k1 = *key1
		}
		// start of array2
		if _, err := dec.ReadToken(); err != nil {
			return err
		}
		innerMap := make(map[string]PackageInformation)
		for dec.PeekKind() != jsontext.EndArray.Kind() {
			if _, err := dec.ReadToken(); err != nil {
				return err
			}

			var key2 *string
			if err := json.UnmarshalDecode(dec, &key2); err != nil {
				return err
			}
			k2 := ""
			if key2 != nil {
				k2 = *key2
			}

			var info PackageInformation
			if err := json.UnmarshalDecode(dec, &info); err != nil {
				return err
			}

			innerMap[k2] = info

			if _, err := dec.ReadToken(); err != nil {
				return err
			}
		}
		// end of array2
		if _, err := dec.ReadToken(); err != nil {
			return err
		}
		res[k1] = innerMap

		// end of array1
		if _, err := dec.ReadToken(); err != nil {
			return err
		}
	}
	if _, err := dec.ReadToken(); err != nil {
		return err
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
