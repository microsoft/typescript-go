package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestSemanticClassificationAlias(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /a.ts
export type x = number;
export class y {};
// @Filename: /b.ts
import { /*0*/x, /*1*/y } from "./a";
const v: /*2*/x = /*3*/y;`
	f := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	f.GoToFile(t, "/b.ts")
	f.VerifySemanticTokens(t, []fourslash.SemanticToken{{Type: "variable.declaration.readonly", Text: "v"}, {Type: "type", Text: "x"}, {Type: "class", Text: "y"}})
}
