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
		ScriptKind: core.ScriptKindTS,
		Mappings:   mappings,
		Supplemental: []contentmapper.MappedResult{
			{ScriptKind: core.ScriptKindJS, Mappings: mappings},
			{ScriptKind: core.ScriptKindJSX, Mappings: mappings},
			{ScriptKind: core.ScriptKindTS, Mappings: mappings},
			{ScriptKind: core.ScriptKindTSX, Mappings: mappings},
			{ScriptKind: core.ScriptKindJSON, Mappings: mappings},
		},
	}
	files, err := contentmapper.ParseResult(
		ast.SourceFileParseOptions{FileName: "/component.astro", Path: "/component.astro"},
		"",
		&contentmapper.Mapper{Manifest: contentmapper.Manifest{Name: "mapper"}},
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
		{"/component.astro.4.json", core.ScriptKindJSON},
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
		&contentmapper.Mapper{Manifest: contentmapper.Manifest{Name: "mapper"}},
		contentmapper.Result{
			Text:       "export {};",
			ScriptKind: core.ScriptKindTS,
			Mappings:   mappings,
			Supplemental: []contentmapper.MappedResult{{
				Text:       "export const value = 1;",
				ScriptKind: core.ScriptKindTS,
				Mappings:   mappings,
			}},
		},
	)
	assert.NilError(t, err)
	assert.Assert(t, ast.IsExternalModule(files.Supplemental[0]))
}
