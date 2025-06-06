package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestAutoImportsWithRootDirsAndRootedPath01(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		// Without embedding, we'd need to read all of the lib files out from disk into the MapFS.
		// Just skip this for now.
		t.Skip("bundled files are not embedded")
	}
	t.Skip()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /dir/foo.ts
 export function foo() {}
// @Filename: /dir/bar.ts
 /*$*/
// @Filename: /dir/tsconfig.json
{
    "compilerOptions": {
        "module": "amd",
        "moduleResolution": "classic",
        "rootDirs": ["D:/"]
    }
}`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.GoToMarker(t, "$")
	f.VerifyCompletions(t, nil, nil)
}
