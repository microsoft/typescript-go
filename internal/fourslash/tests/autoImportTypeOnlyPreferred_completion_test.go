package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

// Tests that auto-import from completions correctly picks the value import
// (not the type-only import) when two imports from the same module exist.
func TestAutoImportTypeOnlyPreferredCompletion(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @module: esnext
// @Filename: /mod.ts
export interface Dispatch {}
export declare function useReducer(): void;
export declare function useState(): void;
// @Filename: /main.ts
import type { Dispatch } from "./mod";
import { useReducer } from "./mod";

useState/**/`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyApplyCodeActionFromCompletion(t, new(""), &fourslash.ApplyCodeActionFromCompletionOptions{
		Name:        "useState",
		Source:      "./mod",
		Description: "Update import from \"./mod\"",
		NewFileContent: new(`import type { Dispatch } from "./mod";
import { useReducer, useState } from "./mod";

useState`),
	})
}

// Tests that auto-import from completions correctly picks the type-only import
// (not the value import) when adding a type and two imports from the same module exist.
func TestAutoImportTypeOnlyPreferredCompletionReversed(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @module: esnext
// @Filename: /mod.ts
export interface ComponentType {}
export interface ComponentProps {}
export declare function useState(): void;
// @Filename: /main.ts
import { useState } from "./mod";
import type { ComponentType } from "./mod";

type _ = ComponentProps/**/;`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyApplyCodeActionFromCompletion(t, new(""), &fourslash.ApplyCodeActionFromCompletionOptions{
		Name:        "ComponentProps",
		Source:      "./mod",
		Description: "Update import from \"./mod\"",
		NewFileContent: new(`import { useState } from "./mod";
import type { ComponentProps, ComponentType } from "./mod";

type _ = ComponentProps;`),
	})
}

// Tests the fix for the specific divergence between isTypeOnlyLocation and
// IsValidTypeOnlyAliasUseSite. In an incomplete type argument position
// (e.g., `myFunc<MyClass` without closing `>`), the parser treats it as a
// binary expression. isPossiblyTypeArgumentPosition returns true (making
// isTypeOnlyLocation true), but IsValidTypeOnlyAliasUseSite correctly
// returns false since the location is an expression node. This ensures a
// class (which has both type and value meaning) is added to the value import.
func TestAutoImportTypeOnlyPreferredCompletionTypeArgPosition(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @module: esnext
// @Filename: /mod.ts
export class MyClass {}
export declare function myFunc<T>(): T;
// @Filename: /main.ts
import type { } from "./mod";
import { myFunc } from "./mod";

myFunc<MyClass/**/`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.VerifyApplyCodeActionFromCompletion(t, new(""), &fourslash.ApplyCodeActionFromCompletionOptions{
		Name:        "MyClass",
		Source:      "./mod",
		Description: "Update import from \"./mod\"",
		NewFileContent: new(`import type { } from "./mod";
import { MyClass, myFunc } from "./mod";

myFunc<MyClass`),
	})
}

