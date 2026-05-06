//// [tests/cases/compiler/declarationEmitPropertyNumericStringKey.ts] ////

//// [declarationEmitPropertyNumericStringKey.ts]
// https://github.com/microsoft/TypeScript/issues/55292

const STATUS = {
    ["404"]: "not found",
} as const;

const hundredStr = "100";
const obj = { [hundredStr]: "foo" };

const hundredNum = 100;
const obj2 = { [hundredNum]: "bar" };


//// [declarationEmitPropertyNumericStringKey.js]
"use strict";
// https://github.com/microsoft/TypeScript/issues/55292
const STATUS = {
    ["404"]: "not found",
};
const hundredStr = "100";
const obj = { [hundredStr]: "foo" };
const hundredNum = 100;
const obj2 = { [hundredNum]: "bar" };


//// [declarationEmitPropertyNumericStringKey.d.ts]
const STATUS: {
    readonly ["404"]: "not found";
};
const hundredStr = "100";
const obj: {
    "100": string;
};
const hundredNum = 100;
const obj2: {
    100: string;
};


//// [DtsFileErrors]


declarationEmitPropertyNumericStringKey.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitPropertyNumericStringKey.d.ts (1 errors) ====
    const STATUS: {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        readonly ["404"]: "not found";
    };
    const hundredStr = "100";
    const obj: {
        "100": string;
    };
    const hundredNum = 100;
    const obj2: {
        100: string;
    };
    