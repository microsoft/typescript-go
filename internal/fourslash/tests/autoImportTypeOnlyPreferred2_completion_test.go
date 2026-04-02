package fourslash_test

import (
"testing"

"github.com/microsoft/typescript-go/internal/fourslash"
"github.com/microsoft/typescript-go/internal/lsp/lsproto"
"github.com/microsoft/typescript-go/internal/testutil"
)

func TestAutoImportTypeOnlyPreferred2_Completion(t *testing.T) {
t.Parallel()
defer testutil.RecoverAndFail(t, "Panic on fourslash test")
const content = `// @module: node18
// @Filename: /mod.ts
export interface ComponentType {}
export interface ComponentProps {}
export declare function useState<T>(initialState: T): [T, (newState: T) => void];
export declare function useEffect(callback: () => void, deps: any[]): void;
// @Filename: /main.ts
import type { ComponentType } from "./mod.js";
import { useState } from "./mod.js";

export function Component({ prop } : { prop: ComponentType }) {
    const codeIsUnimportant = useState(1);
    /*1*/
}`
f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
defer done()
f.VerifyApplyCodeActionFromCompletion(t, new("1"), &fourslash.ApplyCodeActionFromCompletionOptions{
Name:        "useEffect",
Source:      "./mod",
Description: "Add import from \"./mod.js\"",
AutoImportFix: &lsproto.AutoImportFix{
ModuleSpecifier: "./mod.js",
},
NewFileContent: new(`import type { ComponentType } from "./mod.js";
import { useEffect, useState } from "./mod.js";

export function Component({ prop } : { prop: ComponentType }) {
    const codeIsUnimportant = useState(1);
    
}`),
})
}
