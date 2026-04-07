package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

// Adapted from TypeScript's fourslash/server/getFileReferences_deduplicate.ts.
// util.ts is referenced by index.ts, which is included in tsconfig.build.json and
// tsconfig.test.json. The reference will be returned from both projects' language
// services; this test ensures it gets deduplicated.
func TestGetFileReferences_deduplicate(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /home/src/workspaces/project/tsconfig.json
{ "files": [], "references": [{ "path": "tsconfig.build.json" }, { "path": "tsconfig.test.json" }] }
// @Filename: /home/src/workspaces/project/tsconfig.utils.json
{ "compilerOptions": { "lib": ["es5"], "rootDir": "src", "outDir": "dist/utils", "composite": true }, "files": ["util.ts"] }
// @Filename: /home/src/workspaces/project/tsconfig.build.json
{ "compilerOptions": { "lib": ["es5"], "rootDir": "src", "outDir": "dist/build", "composite": true }, "files": ["index.ts"], "references": [{ "path": "tsconfig.utils.json" }] }
// @Filename: /home/src/workspaces/project/index.ts
export * from "./util";
// @Filename: /home/src/workspaces/project/tsconfig.test.json
{ "compilerOptions": { "lib": ["es5"], "rootDir": "src", "outDir": "dist/test", "composite": true }, "files": ["test.ts", "index.ts"], "references": [{ "path": "tsconfig.utils.json" }] }
// @Filename: /home/src/workspaces/project/test.ts
import "./util";
// @Filename: /home/src/workspaces/project/util.ts
export {}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineFindFileReferences(t, "/home/src/workspaces/project/util.ts")
}
