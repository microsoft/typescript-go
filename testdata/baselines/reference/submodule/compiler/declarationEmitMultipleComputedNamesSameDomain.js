//// [tests/cases/compiler/declarationEmitMultipleComputedNamesSameDomain.ts] ////

//// [declarationEmitMultipleComputedNamesSameDomain.ts]
declare const x: string;
declare const y: "y";

export class Test {
    [x] = 10;
    [y] = 10;
}

//// [declarationEmitMultipleComputedNamesSameDomain.js]
var _a, _b;
export class Test {
    constructor() {
        this[_a] = 10;
        this[_b] = 10;
    }
}
_a = x, _b = y;


//// [declarationEmitMultipleComputedNamesSameDomain.d.ts]
const x: string;
const y: "y";
export class Test {
    [x]: number;
    [y]: number;
}
export {};


//// [DtsFileErrors]


declarationEmitMultipleComputedNamesSameDomain.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitMultipleComputedNamesSameDomain.d.ts (1 errors) ====
    const x: string;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    const y: "y";
    export class Test {
        [x]: number;
        [y]: number;
    }
    export {};
    