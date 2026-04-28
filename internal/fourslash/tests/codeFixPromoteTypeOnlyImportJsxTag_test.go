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
// https://github.com/microsoft/typescript-go/issues/1234
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
	// The main goal is that this doesn't panic. The fix should promote
	// the type-only import of React to a regular import.
	f.VerifyImportFixAtPosition(t, []string{
		`import type React from "./react";
import React from "./react";

<Foo />;`,
		`import React from "./react";

<Foo />;`,
		`import type React from "./react";
import React from "./react";

<Foo />;`,
	}, nil /*preferences*/)
}
