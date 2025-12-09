package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestDestructuredInterfaceJSDoc(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
interface FooBar {
    /** foo comment */
    foo: number;
    /** bar comment */
    bar: string;
    /** baz comment */
    baz: string;
}

declare const fubar: FooBar;

const {/*1*/foo, /*2*/bar, /*3*/baz: /*4*/biz} = fubar;
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyQuickInfoAt(t, "1", "const foo: number", "foo comment")
	f.VerifyQuickInfoAt(t, "2", "const bar: string", "bar comment")
	f.VerifyQuickInfoAt(t, "3", "(property) FooBar.baz: string", "baz comment")
	f.VerifyQuickInfoAt(t, "4", "const biz: string", "baz comment")
}

func TestDestructuredInterfaceJSDocWithRename(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
interface FooBar {
    /** foo comment */
    foo: number;
    /** bar comment */
    bar: string;
}

declare const fubar: FooBar;

const {foo: /*1*/myFoo, bar: /*2*/myBar} = fubar;
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyQuickInfoAt(t, "1", "const myFoo: number", "foo comment")
	f.VerifyQuickInfoAt(t, "2", "const myBar: string", "bar comment")
}
