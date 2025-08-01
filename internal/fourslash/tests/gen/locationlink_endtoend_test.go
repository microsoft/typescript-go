package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestLocationLinkEndToEnd(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: types.ts
export interface [|/*def*/Person|] {
    name: string;
    age: number;
}

export function createPerson(name: string, age: number): Person {
    return { name, age };
}

// @Filename: main.ts
import { /*usage*/Person, createPerson } from "./types";

const john: Person = createPerson("John", 30);`

	// Test with LinkSupport enabled
	linkSupport := true
	capabilities := &lsproto.ClientCapabilities{
		TextDocument: &lsproto.TextDocumentClientCapabilities{
			Definition: &lsproto.DefinitionClientCapabilities{
				LinkSupport: &linkSupport,
			},
		},
	}

	f := fourslash.NewFourslash(t, capabilities, content)
	
	// Verify that going to definition from "usage" marker works
	// and returns LocationLink format when client supports it
	f.VerifyBaselineGoToDefinition(t, "usage")
}