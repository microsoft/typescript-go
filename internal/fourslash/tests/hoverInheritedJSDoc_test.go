package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestHoverInheritedJSDocInterface(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
interface foo {
  /** base jsdoc */
  bar(k: string): number;
  /** other jsdoc */
  other: 24;
}

interface bar extends foo {
  bar(k: string | symbol): number | 99;
}

declare const f: foo;
declare const b: bar;

f./*1*/bar;
f./*2*/other;
b./*3*/bar;
b./*4*/other;
`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.VerifyBaselineHover(t)
}

func TestHoverInheritedJSDocClass(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
declare class thing {
  /** doc comment */
  method(s: string): void;
}

declare class potato extends thing {
  method(s: "1234"): void;
}

declare const t: thing;
declare const p: potato;

t./*1*/method;
p./*2*/method;
`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	// t.method should show "doc comment"
	f.VerifyQuickInfoAt(t, "1", "(method) thing.method(s: string): void", "doc comment")
	// p.method should inherit "doc comment" from thing
	f.VerifyQuickInfoAt(t, "2", "(method) potato.method(s: \"1234\"): void", "doc comment")
}
