package lsproto

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type DocumentUri string // !!!

func (uri DocumentUri) FileName() string {
	if strings.HasPrefix(string(uri), "file://") {
		parsed := core.Must(url.Parse(string(uri)))
		if parsed.Host != "" {
			return "//" + parsed.Host + parsed.Path
		}
		return fixWindowsURIPath(parsed.Path)
	}

	// Leave all other URIs escaped so we can round-trip them.

	scheme, path, ok := strings.Cut(string(uri), ":")
	if !ok {
		panic(fmt.Sprintf("invalid URI: %s", uri))
	}

	authority := "ts-nul-authority"
	if rest, ok := strings.CutPrefix(path, "//"); ok {
		authority, path, ok = strings.Cut(rest, "/")
		if !ok {
			panic(fmt.Sprintf("invalid URI: %s", uri))
		}
	}

	return "^/" + scheme + "/" + authority + "/" + path
}

func (uri DocumentUri) Path(useCaseSensitiveFileNames bool) tspath.Path {
	fileName := uri.FileName()
	return tspath.ToPath(fileName, "", useCaseSensitiveFileNames)
}

func fixWindowsURIPath(path string) string {
	if rest, ok := strings.CutPrefix(path, "/"); ok {
		if volume, rest, ok := tspath.SplitVolumePath(rest); ok {
			return volume + rest
		}
	}
	return path
}

type URI string // !!!

type Method string

type Nullable[T any] struct {
	Value T
	Null  bool
}

func ToNullable[T any](v T) Nullable[T] {
	return Nullable[T]{Value: v}
}

func Null[T any]() Nullable[T] {
	return Nullable[T]{Null: true}
}

func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if n.Null {
		return []byte(`null`), nil
	}
	return json.Marshal(n.Value)
}

func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	*n = Nullable[T]{}
	if string(data) == `null` {
		n.Null = true
		return nil
	}
	return json.Unmarshal(data, &n.Value)
}

func unmarshalPtrTo[T any](data []byte) (*T, error) {
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %T: %w", (*T)(nil), err)
	}
	return &v, nil
}

func unmarshalAny(data []byte) (any, error) {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("failed to unmarshal any: %w", err)
	}
	return v, nil
}

func unmarshalEmpty(data []byte) (any, error) {
	if len(data) != 0 {
		return nil, fmt.Errorf("expected empty, got: %s", string(data))
	}
	return nil, nil
}

func assertOnlyOne(message string, values ...bool) {
	count := 0
	for _, v := range values {
		if v {
			count++
		}
	}
	if count != 1 {
		panic(message)
	}
}
