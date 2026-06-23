package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestCodeFixMissingTypeAnnotationOnExportsUsingTypeFromUnimportedFile(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `
// @module: preserve
// @verbatimModuleSyntax: true
// @declaration: true
// @isolatedDeclarations: true

// @Filename: /types.ts
export type Thing = {};

// @Filename: /funcs.ts
import type { Thing } from "./types";

export function makeThing(): Thing {
    throw {};
}

// @Filename: /exporter.ts
import { makeThing } from "./funcs";

export let thing/**/ = makeThing();`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	f.GoToMarker(t, "")
	f.VerifyCodeFix(t, fourslash.VerifyCodeFixOptions{
		Description: "Add annotation of type 'Thing'",
		NewFileContent: `import { makeThing } from "./funcs";
import type { Thing } from "./types";

export let thing: Thing = makeThing();`,
	})
}
