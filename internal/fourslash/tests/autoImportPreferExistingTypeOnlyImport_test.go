package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

// Regression test: auto-import should add a type to an existing `import type`
// declaration instead of adding an inline `type` modifier to a value import.
// https://github.com/microsoft/typescript-go/issues/3029
func TestAutoImportPreferExistingTypeOnlyImport(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /node_modules/typebox/index.d.ts
export declare const Type: { Object: Function; String: Function };
export type Static<T> = T;
// @Filename: /main.ts
import { Type } from "typebox";
import type { } from "typebox";

const object = Type.Object({});

type ObjectType = Static/**/;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "")
	f.VerifyImportFixAtPosition(t, []string{
		`import { Type } from "typebox";
import type { Static } from "typebox";

const object = Type.Object({});

type ObjectType = Static;`,
	}, nil /*preferences*/)
}
