package transpile_test

import (
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/contentmapper"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/locale"
	"github.com/microsoft/typescript-go/internal/testutil/contentmappertest"
	"github.com/microsoft/typescript-go/internal/transpile"
	"gotest.tools/v3/assert"
)

func TestTranspileModuleWithContentMapper(t *testing.T) {
	t.Parallel()
	host := contentmapper.NewHost(t.Context(), contentmappertest.NewSpawner(), locale.Default)
	defer func() { assert.NilError(t, host.Close()) }()

	output := transpile.TranspileModule(t.Context(), "export const version = #{target};\n", transpile.Options{
		CompilerOptions: &core.CompilerOptions{Target: core.ScriptTargetES2020, SourceMap: core.TSTrue},
		FileName:        "component.box",
		ContentMapper: &transpile.ContentMapperOptions{
			Manifest: contentmapper.Manifest{
				Name:            contentmappertest.PackageName,
				Version:         "1.0.0",
				Exec:            []string{contentmappertest.TransformingMapper},
				CompilerOptions: contentmappertest.DeclaredOptions,
				SupportsEmit:    true,
				Extensions:      map[string]string{".box": ".ts"},
			},
			Options: map[string]any{"mode": "test"},
		},
	}, host)

	assert.Assert(t, output != nil)
	assert.Assert(t, strings.Contains(output.OutputText, "export const version = 7;"), output.OutputText)
	assert.Equal(t, output.SourceMapText, "")
}

func TestTranspileDeclarationWithDynamicConfigContentMapper(t *testing.T) {
	t.Parallel()
	lifecycle := &contentmappertest.ProjectLifecycle{}
	host := contentmapper.NewHost(t.Context(), contentmappertest.NewSpawnerWithProjectLifecycle(lifecycle), locale.Default)
	defer func() { assert.NilError(t, host.Close()) }()

	output := transpile.TranspileDeclaration(t.Context(), "export const value = 1;\n", transpile.Options{
		FileName: "component.box",
		ContentMapper: &transpile.ContentMapperOptions{
			Manifest: contentmapper.Manifest{
				Name:          contentmappertest.PackageName,
				Version:       "1.0.0",
				Exec:          []string{contentmappertest.DynamicVerbatimMapper},
				DynamicConfig: true,
				Extensions:    map[string]string{".box": ".ts"},
			},
			Options: map[string]any{"mode": "test"},
		},
	}, host)

	assert.Assert(t, output != nil)
	assert.Assert(t, strings.Contains(output.OutputText, "export declare const value = 1;"), output.OutputText)
	assert.Equal(t, lifecycle.Opens.Load(), int32(1))
	assert.Equal(t, lifecycle.Closes.Load(), int32(1))
}
