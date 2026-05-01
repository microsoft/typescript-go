package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

// Test that auto-imports for JSX tags don't crash when React is type-imported.
// When both the JSX namespace (React) and the component need to be imported,
// getSymbolNamesToImport returns multiple names and the type-only promotion
// path should handle this gracefully instead of panicking.
func TestCodeFixPromoteTypeOnlyImportJsxTag(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @module: preserve
// @verbatimModuleSyntax: true
// @jsx: react
// @Filename: /react.ts
const React: any = {};
export default React;
// @Filename: /bar.tsx
import type React from "./react";

<Foo/**/ />;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "")
	// The fix should promote the type-only import of React to a regular import.
	// Only the promotion fix is returned; the missing "Foo" component is a
	// separate error handled by a different diagnostic.
	f.VerifyImportFixAtPosition(t, []string{
		// Promotion fix from the "cannot use as value because type-imported" error
		`import React from "./react";

<Foo />;`,
		// Auto-import fix from the "Cannot find name 'Foo'" error, which also
		// needs React for JSX
		`import type React from "./react";
import React from "./react";

<Foo />;`,
	}, nil /*preferences*/)
}
