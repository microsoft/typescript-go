package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestAutoFormattingOnPasting(t *testing.T) {
<<<<<<< HEAD
	t.Parallel()
	t.Skip()
=======
	fourslash.SkipIfFailing(t)
	t.Parallel()
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `module TestModule {
/**/
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "")
	f.Paste(t, " class TestClass{\nprivate   foo;\npublic testMethod( )\n{}\n}")
<<<<<<< HEAD
	f.VerifyCurrentFileContent(t, `module TestModule {
    class TestClass {
        private foo;
        public testMethod() { }
    }
}`)
=======
	f.VerifyCurrentFileContentIs(t, "module TestModule {\n    class TestClass {\n        private foo;\n        public testMethod() { }\n    }\n}")
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
}
