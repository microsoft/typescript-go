package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestFindFileReferences(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /a.ts
export const a = {};

// @Filename: /b.ts
import "./a";

// @Filename: /c.ts
import {} from "./a";

// @Filename: /d.ts
import { a } from "/a";
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineFindFileReferences(t, "/a.ts")
}

func TestFindFileReferencesNonModule(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /a.ts
const x = 1;

// @Filename: /b.ts
/// <reference path="./a.ts" />
const y = 2;
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineFindFileReferences(t, "/a.ts")
}
