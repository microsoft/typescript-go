package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestGoToDefinitionGetterReturnsCallableInterface(t *testing.T) {
	t.Parallel()

	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: type.d.ts
export interface Disposable {
  dispose(): void;
}

export interface [|Event|]<T> {
  [|(|]
    listener: (e: T) => any,
    thisArgs?: any,
    disposables?: Disposable[]
  [|)|]: Disposable;
}

export interface TextDocumentChangeEvent<T> {
  document: T;
}

export declare class TextDocuments<
  T extends {
    uri: string;
  }
> {
  [|get onDidChangeContent()|]: Event<TextDocumentChangeEvent<T>>;
}

export interface TextDocument {
  uri: string;
}

// @Filename: index.ts
import { TextDocument, TextDocuments } from "./type";

var documents: TextDocuments<TextDocument>;

documents!.[|onDid/*1*/ChangeContent|]()
`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyBaselineGoToDefinition(t, true, "1")
}
