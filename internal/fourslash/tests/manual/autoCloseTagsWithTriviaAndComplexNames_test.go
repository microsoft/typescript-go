package fourslash

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestAutoCloseTagsWithTriviaAndComplexNames(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// Using separate files for each example to avoid unclosed JSX tags affecting other tests.
	const content = `// @noLib: true

// @Filename: /0.tsx
// JSDoc
const x = <
	/** hello world! */
	div /** hello world! */
	>/*0*/;

// @Filename: /1.tsx
// Single-line comments
const x =
	<
	// hello world!
	div // hello world!
	>/*1*/;

// @Filename: /2.tsx
// Namespaced tag
const x =
	<ns:sometag>/*2*/

// @Filename: /3.tsx
// Namespace with single-line comments
const x = <
	// pre-ns	
	ns
	// pre-colon
	:
	// post-colon
	sometag
	// post-id
	>/*3*/;

// @Filename: /4.tsx
// UppercaseComponent-named tag
const x = <SomeComponent>/*4*/
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineClosingTags(t)
}
