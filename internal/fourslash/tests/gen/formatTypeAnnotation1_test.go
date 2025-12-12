package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFormatTypeAnnotation1(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `function foo(x: number, y?: string): number {}
interface Foo {
    x: number;
    y?: number;
}`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.SetFormatOption(t, "insertSpaceBeforeTypeAnnotation", true)
	f.FormatDocument(t, "")
	f.VerifyCurrentFileContent(t, `function foo(x : number, y ?: string) : number { }
interface Foo {
    x : number;
    y ?: number;
}`)
}
