//// [tests/cases/compiler/signatureHelpTokenCacheMismatch.ts] ////

//// [signatureHelpTokenCacheMismatch.ts]
// Test case for token cache mismatch issue in signature help
// This should reproduce the issue when AST structure changes
declare const array: number[];
array?.at(0);

//// [signatureHelpTokenCacheMismatch.js]
array?.at(0);
