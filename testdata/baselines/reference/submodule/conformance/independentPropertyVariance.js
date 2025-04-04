//// [tests/cases/conformance/types/typeRelationships/comparable/independentPropertyVariance.ts] ////

//// [independentPropertyVariance.ts]
// Verify that properties can vary independently in comparable relationship

declare const x: { a: 1, b: string };
declare const y: { a: number, b: 'a' };

x === y;


//// [independentPropertyVariance.js]
x === y;
