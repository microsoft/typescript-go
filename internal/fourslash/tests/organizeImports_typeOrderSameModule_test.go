package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestOrganizeImports_typeOrderSameModule(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// Test that import type statements are sorted before regular import statements
	// when importing from the same module. Type-only imports should come before
	// value imports according to getImportKindOrder.
	const content = `import { foo } from 'package';
import type { Foo } from 'package';

console.log(foo, Foo);`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyOrganizeImports(t,
		`import type { Foo } from 'package';
import { foo } from 'package';

console.log(foo, Foo);`,
		lsproto.CodeActionKindSourceOrganizeImports,
		nil,
	)
}

func TestOrganizeImports_typeOrderSameModuleReversed(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	// Test that when import type is already first, it stays first
	const content = `import type { Foo } from 'package';
import { foo } from 'package';

console.log(foo, Foo);`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyOrganizeImports(t,
		`import type { Foo } from 'package';
import { foo } from 'package';

console.log(foo, Foo);`,
		lsproto.CodeActionKindSourceOrganizeImports,
		nil,
	)
}
