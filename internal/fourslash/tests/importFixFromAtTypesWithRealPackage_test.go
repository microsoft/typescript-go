package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestImportFixFromAtTypesWithRealPackage(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// Simulate a project where both `myLib` (JS-only package) and `@types/myLib` (type declarations) are installed.
	// The import fix should suggest importing from "myLib", not "@types/myLib".
	const content = `// @Filename: /node_modules/myLib/package.json
{"name":"myLib","version":"1.0.0","main":"index.js"}
// @Filename: /node_modules/myLib/index.js
module.exports = {};
// @Filename: /node_modules/@types/myLib/index.d.ts
export function f1(): void;
export var v1: number;
// @Filename: /index.ts
[|f1/*0*/();|]`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyImportFixModuleSpecifiers(t, "0", []string{"myLib"}, nil /*preferences*/)
}

func TestImportFixFromAtTypesWithRealPackageExports(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// Like the above test, but the real package has an exports field pointing to JS files.
	// This is the React 19 scenario: react has exports but no .d.ts, @types/react provides types.
	const content = `// @Filename: /node_modules/myLib/package.json
{"name":"myLib","version":"1.0.0","exports":{".":{"default":"./index.js"}}}
// @Filename: /node_modules/myLib/index.js
module.exports = {};
// @Filename: /node_modules/@types/myLib/index.d.ts
export function f1(): void;
export var v1: number;
// @Filename: /index.ts
[|f1/*0*/();|]`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyImportFixModuleSpecifiers(t, "0", []string{"myLib"}, nil /*preferences*/)
}
