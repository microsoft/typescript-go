package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestCodeLensOnFunctionAcrossProjects1(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @filename: ./a/tsconfig.json
{
  "compilerOptions": {
	"composite": true,
	"declaration": true,
	"declarationMaps": true,
	"outDir": "./dist",
	"rootDir": "src"
  },
  "include": ["./src"]
}

// @filename: ./a/src/foo.ts
export function aaa() {}
aaa();

// @filename: ./b/tsconfig.json
{
  "compilerOptions": {
	"composite": true,
	"declaration": true,
	"declarationMaps": true,
	"outDir": "./dist",
	"rootDir": "src"
  },
  "references": [{ "path": "../a" }],
  "include": ["./src"]
}

// @filename: ./b/src/bar.ts
import * as foo from '../../a/dist/foo.js';
foo.aaa();
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	f.VerifyBaselineCodeLens(t, &lsutil.UserPreferences{
		CodeLens: lsutil.CodeLensUserPreferences{
			ReferencesCodeLensEnabled: true,
		},
	})
}
