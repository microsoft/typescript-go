package contentmapper_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/contentmapper"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/spanmap"
	"github.com/microsoft/typescript-go/internal/tspath"
	"gotest.tools/v3/assert"
)

func TestParseResultSupplementalFileExtensions(t *testing.T) {
	t.Parallel()
	mappings := spanmap.New(nil)
	result := contentmapper.Result{
		Mappings: mappings,
		Supplemental: []contentmapper.MappedResult{
			{VirtualExtension: ".js", Mappings: mappings},
			{VirtualExtension: ".jsx", Mappings: mappings},
			{VirtualExtension: ".ts", Mappings: mappings},
			{VirtualExtension: ".tsx", Mappings: mappings},
			{VirtualExtension: ".mts", Mappings: mappings},
			{VirtualExtension: ".cts", Mappings: mappings},
			{VirtualExtension: ".json", Mappings: mappings},
		},
	}
	files, err := contentmapper.ParseResult(
		ast.SourceFileParseOptions{FileName: "/component.astro", Path: "/component.astro"},
		"",
		&contentmapper.Mapper{Definition: contentmapper.Definition{Extensions: []string{".astro"}}, Manifest: contentmapper.Manifest{Name: "mapper", Extensions: map[string]string{".astro": ".ts"}}},
		result,
	)
	assert.NilError(t, err)

	expected := []struct {
		fileName   string
		scriptKind core.ScriptKind
	}{
		{"/component.astro.0.js", core.ScriptKindJS},
		{"/component.astro.1.jsx", core.ScriptKindJSX},
		{"/component.astro.2.ts", core.ScriptKindTS},
		{"/component.astro.3.tsx", core.ScriptKindTSX},
		{"/component.astro.4.mts", core.ScriptKindTS},
		{"/component.astro.5.cts", core.ScriptKindTS},
		{"/component.astro.6.json", core.ScriptKindJSON},
	}
	assert.Equal(t, len(files.Supplemental), len(expected))
	for i, expected := range expected {
		assert.Equal(t, files.Supplemental[i].FileName(), expected.fileName)
		assert.Equal(t, files.Supplemental[i].Path(), tspath.Path(expected.fileName))
		assert.Equal(t, files.Supplemental[i].ScriptKind, expected.scriptKind)
	}
}

func TestParseResultAllowsSupplementalModules(t *testing.T) {
	t.Parallel()
	mappings := spanmap.New(nil)
	files, err := contentmapper.ParseResult(
		ast.SourceFileParseOptions{FileName: "/component.astro", Path: "/component.astro"},
		"",
		&contentmapper.Mapper{Definition: contentmapper.Definition{Extensions: []string{".astro"}}, Manifest: contentmapper.Manifest{Name: "mapper", Extensions: map[string]string{".astro": ".ts"}}},
		contentmapper.Result{
			Text:     "export {};",
			Mappings: mappings,
			Supplemental: []contentmapper.MappedResult{{
				Text:             "export const value = 1;",
				VirtualExtension: ".mts",
				Mappings:         mappings,
			}},
		},
	)
	assert.NilError(t, err)
	assert.Assert(t, ast.IsExternalModule(files.Supplemental[0]))
}
