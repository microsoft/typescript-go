package ls_test

import (
	"context"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil/projecttestutil"
	"gotest.tools/v3/assert"
)

func TestInternalAliasGoToDefinition(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	// Source file content
	sourceContent := `export function helperFunction() {
    return "helper result";
}

export const helperConstant = 42;`

	// Declaration file content
	declContent := `export declare function helperFunction(): string;
export declare const helperConstant = 42;
//# sourceMappingURL=utils.d.ts.map`

	// Source map content
	sourceMapContent := `{"version":3,"file":"utils.d.ts","sourceRoot":"","sources":["../../../src/internal/helpers/utils.ts"],"names":[],"mappings":"AAAA,wBAAgB,cAAc,IAAI,MAAM,CAEvC;AAED,eAAO,MAAM,cAAc,KAAK,CAAC"}`

	// Main file content
	mainContent := `import { helperFunction } from "@internal/helpers/utils";

const result = helperFunction();`

	// tsconfig.json with path mapping
	tsconfigContent := `{
  "compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@internal/*": ["dist/internal/*", "src/internal/*"]
    },
    "declaration": true,
    "declarationMap": true,
    "outDir": "dist"
  }
}`

	// Set up test files
	files := map[string]any{
		"/tsconfig.json":                        tsconfigContent,
		"/src/internal/helpers/utils.ts":        sourceContent,
		"/dist/internal/helpers/utils.d.ts":     declContent,
		"/dist/internal/helpers/utils.d.ts.map": sourceMapContent,
		"/main.ts":                              mainContent,
	}

	session, _ := projecttestutil.Setup(files)

	ctx := projecttestutil.WithRequestID(context.Background())
	session.DidOpenFile(ctx, "file:///main.ts", 1, mainContent, lsproto.LanguageKindTypeScript)

	languageService, err := session.GetLanguageService(ctx, "file:///main.ts")
	assert.NilError(t, err)

	uri := lsproto.DocumentUri("file:///main.ts")
	lspPosition := lsproto.Position{Line: 0, Character: 9}

	definition, err := languageService.ProvideDefinition(ctx, uri, lspPosition)
	assert.NilError(t, err)

	if definition.Locations != nil {
		assert.Assert(t, len(*definition.Locations) == 1, "Expected 1 definition location, got %d", len(*definition.Locations))

		location := (*definition.Locations)[0]
		expectedURI := ls.FileNameToDocumentURI("/src/internal/helpers/utils.ts")
		actualURI := location.Uri

		assert.Equal(t, string(expectedURI), string(actualURI), "Should resolve to source .ts file, not .d.ts file")
	} else if definition.Location != nil {
		expectedURI := ls.FileNameToDocumentURI("/src/internal/helpers/utils.ts")
		assert.Equal(t, string(expectedURI), string(definition.Location.Uri), "Should resolve to source .ts file, not .d.ts file")
	} else {
		t.Fatal("No definition found - expected to find definition")
	}
}
