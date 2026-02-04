package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
)

func TestPackageJsonImportsMalformed(t *testing.T) {
	t.Parallel()
	// This test verifies that a malformed package.json imports mapping doesn't crash the server.
	// The mapping "./src/*ts" is missing a dot before "ts" (should be "./src/*.ts").
	// When importing "#/b.", the pattern substitution results in a module reference
	// without a valid TS extension, which previously caused a panic.
	const content = `// @Filename: /src/a.ts
import * as b from "#/b.";

b.foo();

// @Filename: /src/b.ts
export function foo() {}

// @Filename: /package.json
{
    "imports": {
        "#/*": {
            "types": "./src/*ts",
            "default": "./dist/*js"
        }
    }
}

// @Filename: /tsconfig.json
{
    "compilerOptions": {
        "module": "nodenext",
        "moduleResolution": "nodenext",
        "rootDir": "src",
        "outDir": "dist"
    },
    "include": ["src"]
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	// Request diagnostics for src/a.ts - this should NOT crash
	// The malformed package.json might produce some diagnostics, but that's fine
	f.VerifyDiagnostics(t, nil)
}
