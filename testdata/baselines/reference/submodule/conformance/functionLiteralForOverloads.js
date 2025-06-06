//// [tests/cases/conformance/types/specifyingTypes/typeLiterals/functionLiteralForOverloads.ts] ////

//// [functionLiteralForOverloads.ts]
// basic uses of function literals with overloads

var f: {
    (x: string): string;
    (x: number): number;
} = (x) => x;

var f2: {
    <T>(x: string): string;
    <T>(x: number): number;
} = (x) => x;

var f3: {
    <T>(x: T): string;
    <T>(x: T): number;
} = (x) => x;

var f4: {
    <T>(x: string): T;
    <T>(x: number): T;
} = (x) => x;

//// [functionLiteralForOverloads.js]
// basic uses of function literals with overloads
var f = (x) => x;
var f2 = (x) => x;
var f3 = (x) => x;
var f4 = (x) => x;
