package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormatSpaceAfterTemplateHeadAndMiddle(t *testing.T) {
	t.Parallel()
	t.Skip()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `const a1 = ` + "`" + `${1}${1}` + "`" + `;
const a2 = ` + "`" + `
    ${1}${1}
` + "`" + `;
const a3 = ` + "`" + `


    ${1}${1}
` + "`" + `;
const a4 = ` + "`" + `

    ${1}${1}

` + "`" + `;
const a5 = ` + "`" + `text ${1} text ${1} text` + "`" + `;
const a6 = ` + "`" + `
    text ${1}
    text ${1}
    text
` + "`" + `;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.SetFormatOption(t, "insertSpaceAfterOpeningAndBeforeClosingTemplateStringBraces", true)
	f.FormatDocument(t, "")
	f.VerifyCurrentFileContent(t, "const a1 = `${ 1 }${ 1 }`;\n"+"const a2 = `\n"+`    ${ 1 }${ 1 }
`+"`;\n"+"const a3 = `\n"+`
`+`
`+`    ${ 1 }${ 1 }
`+"`;\n"+"const a4 = `\n"+`
`+`    ${ 1 }${ 1 }
`+`
`+"`;\n"+"const a5 = `text ${ 1 } text ${ 1 } text`;\n"+"const a6 = `\n"+`    text ${ 1 }
`+`    text ${ 1 }
`+`    text
`+"`;")
}
