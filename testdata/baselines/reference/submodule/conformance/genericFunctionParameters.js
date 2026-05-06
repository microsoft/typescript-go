//// [tests/cases/conformance/types/typeRelationships/typeInference/genericFunctionParameters.ts] ////

//// [genericFunctionParameters.ts]
declare function f1<T>(cb: <S>(x: S) => T): T;
declare function f2<T>(cb: <S extends number>(x: S) => T): T;
declare function f3<T>(cb: <S extends Array<S>>(x: S) => T): T;

let x1 = f1(x => x);  // {}
let x2 = f2(x => x);  // number
let x3 = f3(x => x);  // Array<any>

// Repro from #19345

declare const s: <R>(go: <S>(ops: { init(): S; }) => R) => R;
const x = s(a => a.init());  // x is any, should have been {}


//// [genericFunctionParameters.js]
"use strict";
let x1 = f1(x => x); // {}
let x2 = f2(x => x); // number
let x3 = f3(x => x); // Array<any>
const x = s(a => a.init()); // x is any, should have been {}


//// [genericFunctionParameters.d.ts]
function f1<T>(cb: <S>(x: S) => T): T;
function f2<T>(cb: <S extends number>(x: S) => T): T;
function f3<T>(cb: <S extends Array<S>>(x: S) => T): T;
let x1: unknown;
let x2: number;
let x3: any[][];
const s: <R>(go: <S>(ops: {
    init(): S;
}) => R) => R;
const x: unknown;


//// [DtsFileErrors]


genericFunctionParameters.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== genericFunctionParameters.d.ts (1 errors) ====
    function f1<T>(cb: <S>(x: S) => T): T;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    function f2<T>(cb: <S extends number>(x: S) => T): T;
    function f3<T>(cb: <S extends Array<S>>(x: S) => T): T;
    let x1: unknown;
    let x2: number;
    let x3: any[][];
    const s: <R>(go: <S>(ops: {
        init(): S;
    }) => R) => R;
    const x: unknown;
    