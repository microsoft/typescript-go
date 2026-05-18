// @strict: true
// @noEmit: true

// https://github.com/microsoft/typescript-go/issues/3976

// Assigning [] to a property of a callable-type variable with a declared type
// annotation should not produce TS7008.

const fn1: {
    (): void;
    items?: string[];
} = () => undefined;
fn1.items = [];  // No error: items is declared as string[] | undefined

const fn2: {
    (): void;
    counts?: number[];
} = () => undefined;
fn2.counts = [];  // No error: counts is declared as number[] | undefined

// A variable with a simple object type annotation also benefits from this.
const obj: { tags?: string[] } = {};
obj.tags = [];  // No error: tags is declared as string[] | undefined

// Without a type annotation, the empty array assignment should still produce TS7008.
const fn3 = () => undefined;
fn3.labels = [];  // Error: Member 'labels' implicitly has an 'any[]' type.
