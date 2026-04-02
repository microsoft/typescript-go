package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/testutil"
)

// Tests that auto-import from completions correctly picks the matching import
// declaration when two imports from the same module exist (one type-only, one value).
// This was fixed in #3244 by making tryAddToExistingImport prefer imports whose
// type-only-ness matches the addAsTypeOnly requirement.
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
