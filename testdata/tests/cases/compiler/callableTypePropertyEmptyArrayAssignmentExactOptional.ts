// @strict: true
// @exactOptionalPropertyTypes: true
// @noEmit: true

// https://github.com/microsoft/typescript-go/issues/3976
// Regression: under exactOptionalPropertyTypes the helper must use getWriteTypeOfSymbol
// (not getTypeOfSymbol) so the internal `missing` type is stripped before being
// returned as the initializer type of the assignment declaration.

const fn1: {
    (): void;
    items?: string[];
} = () => undefined;
fn1.items = [];  // No error: items is declared as string[]

const obj: { tags?: string[] } = {};
obj.tags = [];  // No error: tags is declared as string[]

// Without a type annotation TS7008 should still fire.
const fn2 = () => undefined;
fn2.labels = [];  // Error: Member 'labels' implicitly has an 'any[]' type.
