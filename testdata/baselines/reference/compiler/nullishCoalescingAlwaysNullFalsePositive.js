//// [tests/cases/compiler/nullishCoalescingAlwaysNullFalsePositive.ts] ////

//// [nullishCoalescingAlwaysNullFalsePositive.ts]
// Repro for https://github.com/microsoft/TypeScript/issues/63642
// TS2871 false positive: conditional with ?? null in true branch over spread-of-union

const make = (s: string | null) =>
  s === null
    ? { paymentIntentId: null, transferGroup: null }
    : { paymentIntentId: s as string | null, transferGroup: s as string | null }
const decorated = [null, 'pi_1'].map(s => ({ tag: s, ...make(s) }))
const localTxMap = new Map<string, string>()
const groupMap = new Map<string, string>()
// Should not produce TS2871 on the parenthesized conditional
const rows = decorated.map(({ paymentIntentId, transferGroup }) => ({
  localTransactionId: (paymentIntentId ? localTxMap.get(paymentIntentId) ?? null : null)
    ?? (transferGroup ? groupMap.get(transferGroup) ?? null : null),
}))
console.log(rows)

// Additional cases: ?? semantics
declare function maybeString(): string | undefined;

// x ?? null: can be non-null when x returns a string — NOT always nullish
const a = (maybeString() ?? null) ?? "fallback";

// null ?? null: always nullish (should still error on the left side of ??)
// (no assertion here, just ensuring no false positive from above)


//// [nullishCoalescingAlwaysNullFalsePositive.js]
"use strict";
// Repro for https://github.com/microsoft/TypeScript/issues/63642
// TS2871 false positive: conditional with ?? null in true branch over spread-of-union
const make = (s) => s === null
    ? { paymentIntentId: null, transferGroup: null }
    : { paymentIntentId: s, transferGroup: s };
const decorated = [null, 'pi_1'].map(s => ({ tag: s, ...make(s) }));
const localTxMap = new Map();
const groupMap = new Map();
// Should not produce TS2871 on the parenthesized conditional
const rows = decorated.map(({ paymentIntentId, transferGroup }) => ({
    localTransactionId: (paymentIntentId ? localTxMap.get(paymentIntentId) ?? null : null)
        ?? (transferGroup ? groupMap.get(transferGroup) ?? null : null),
}));
console.log(rows);
// x ?? null: can be non-null when x returns a string — NOT always nullish
const a = (maybeString() ?? null) ?? "fallback";
// null ?? null: always nullish (should still error on the left side of ??)
// (no assertion here, just ensuring no false positive from above)
