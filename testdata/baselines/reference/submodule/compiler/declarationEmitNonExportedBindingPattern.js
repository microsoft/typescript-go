//// [tests/cases/compiler/declarationEmitNonExportedBindingPattern.ts] ////

//// [test.ts]
function getFoo() {
  return { foo: { test: 42 } }
}

const { foo } = getFoo()

export type AliasType = typeof foo

const { foo: renamed } = getFoo()

export type AliasType2 = typeof renamed

function getNested() {
  return { a: { b: { c: 'd' } } }
}

const { a: { b: { c } } } = getNested()

export type AliasType3 = typeof c


//// [test.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
function getFoo() {
    return { foo: { test: 42 } };
}
const { foo } = getFoo();
const { foo: renamed } = getFoo();
function getNested() {
    return { a: { b: { c: 'd' } } };
}
const { a: { b: { c } } } = getNested();


//// [test.d.ts]
const foo: {
    test: number;
};
export type AliasType = typeof foo;
const renamed: {
    test: number;
};
export type AliasType2 = typeof renamed;
const c: string;
export type AliasType3 = typeof c;
export {};


//// [DtsFileErrors]


test.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== test.d.ts (1 errors) ====
    const foo: {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        test: number;
    };
    export type AliasType = typeof foo;
    const renamed: {
        test: number;
    };
    export type AliasType2 = typeof renamed;
    const c: string;
    export type AliasType3 = typeof c;
    export {};
    