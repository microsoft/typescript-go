package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestOrganizeImportsAttributes4(t *testing.T) {
	fourslash.SkipIfFailing(t)
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `import { A } from "./a" assert { foo: "foo", bar: "bar" };
import { B } from "./a" assert { bar: "bar", foo: "foo" };
import { D } from "./a" assert { bar: "foo", foo: "bar" };
import { E } from "./a" assert { foo: 'bar', bar: "foo" };
import { C } from "./a" assert { foo: "bar", bar: "foo" };
import { F } from "./a" assert { foo: "42" };
import { Y } from "./a" assert { foo: 42 };
import { Z } from "./a" assert { foo: "42" };

export type G = A | B | C | D | E | F | Y | Z;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyOrganizeImports(t, `import { A, B } from "./a" assert { foo: "foo", bar: "bar" };
import { C, D, E } from "./a" assert { bar: "foo", foo: "bar" };
import { F, Z } from "./a" assert { foo: "42" };
import { Y } from "./a" assert { foo: 42 };

export type G = A | B | C | D | E | F | Y | Z;`, lsproto.CodeActionKindSourceOrganizeImports, nil)
}
