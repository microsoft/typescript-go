package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestWhiteSpaceTrimming2(t *testing.T) {
<<<<<<< HEAD
	t.Parallel()
	t.Skip()
=======
	fourslash.SkipIfFailing(t)
	t.Parallel()
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `let noSubTemplate = ` + "`" + `/*    /*1*/` + "`" + `;
let templateHead = ` + "`" + `/*    /*2*/${1 + 2}` + "`" + `;
let templateMiddle = ` + "`" + `/*    ${1 + 2    /*3*/}` + "`" + `;
let templateTail = ` + "`" + `/*    ${1 + 2}    /*4*/` + "`" + `;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "1")
	f.Insert(t, "\n")
	f.GoToMarker(t, "2")
	f.Insert(t, "\n")
	f.GoToMarker(t, "3")
	f.Insert(t, "\n")
	f.GoToMarker(t, "4")
	f.Insert(t, "\n")
<<<<<<< HEAD
	f.VerifyCurrentFileContent(t, "let noSubTemplate = `/*    \n`;\nlet templateHead = `/*    \n${1 + 2}`;\nlet templateMiddle = `/*    ${1 + 2\n    }`;\nlet templateTail = `/*    ${1 + 2}    \n`;")
=======
	f.VerifyCurrentFileContentIs(t, "let noSubTemplate = `/*    \n`;\nlet templateHead = `/*    \n${1 + 2}`;\nlet templateMiddle = `/*    ${1 + 2\n    }`;\nlet templateTail = `/*    ${1 + 2}    \n`;")
>>>>>>> 20bf4fc90d3d38016f07fda1fb972eedc715bb02
}
