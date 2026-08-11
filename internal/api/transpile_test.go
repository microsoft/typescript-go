package api

import (
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/contentmapper"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/locale"
	"github.com/microsoft/typescript-go/internal/testutil/contentmappertest"
	"gotest.tools/v3/assert"
)

func TestTranspileOutputWithContentMapper(t *testing.T) {
	t.Parallel()
	host := contentmapper.NewHost(t.Context(), contentmappertest.NewSpawner(), locale.Default)
	defer func() { assert.NilError(t, host.Close()) }()

	output, err := transpileOutput(t.Context(), host, "export const version = #{target};\n", TranspileOptions{
		CompilerOptions: &core.CompilerOptions{Target: core.ScriptTargetES2020},
		FileName:        "component.box",
		ContentMapper: &TranspileContentMapperOptions{
			Manifest: ContentMapperManifest{
				Name:            contentmappertest.PackageName,
				Version:         "1.0.0",
				Exec:            []string{contentmappertest.TransformingMapper},
				CompilerOptions: contentmappertest.DeclaredOptions,
				SupportsEmit:    true,
				Extensions:      map[string]string{".box": ".ts"},
			},
			Options: map[string]any{"mode": "test"},
		},
	}, false)

	assert.NilError(t, err)
	assert.Assert(t, strings.Contains(output.OutputText, "export const version = 7;"), output.OutputText)
}
