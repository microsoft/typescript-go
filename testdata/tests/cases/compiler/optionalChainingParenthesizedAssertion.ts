// @target: es2020
// Tests that parenthesized type assertions around optional chains
// correctly preserve parentheses after type erasure.
// This matches tsc behavior - see TypeScript#50148.

declare let a: undefined | { b: string[], c(): void, d: { e: number } };

// Parenthesized assertion breaks the optional chain - parens preserved
let r1 = (a?.b as any).length;
let r2 = (a?.b as string[]).length;
let r3 = (<any>a?.b).length;
let r4 = (a?.b!).length;
let r5 = (a?.b satisfies any).length;

// Regular optional chain - no parens needed
let r6 = a?.b.length;

// Call expression after assertion - parens preserved
(a?.c as any)();

// Optional chain after assertion - parens removed
(a?.c as any)?.();
