//// [tests/cases/compiler/funcdecl.ts] ////

//// [funcdecl.ts]
function simpleFunc() {
    return "this is my simple func";
}
var simpleFuncVar = simpleFunc;

function anotherFuncNoReturn() {
}
var anotherFuncNoReturnVar = anotherFuncNoReturn;

function withReturn() : string{
    return "Hello";
}
var withReturnVar = withReturn;

function withParams(a : string) : string{
    return a;
}
var withparamsVar = withParams;

function withMultiParams(a : number, b, c: Object) {
    return a;
}
var withMultiParamsVar = withMultiParams;

function withOptionalParams(a?: string) {
}
var withOptionalParamsVar = withOptionalParams;

function withInitializedParams(a: string, b0, b = 30, c = "string value") {
}
var withInitializedParamsVar = withInitializedParams;

function withOptionalInitializedParams(a: string, c: string = "hello string") {
}
var withOptionalInitializedParamsVar = withOptionalInitializedParams;

function withRestParams(a: string, ... myRestParameter : number[]) {
    return myRestParameter;
}
var withRestParamsVar = withRestParams;

function overload1(n: number) : string;
function overload1(s: string) : string;
function overload1(ns: any) {
    return ns.toString();
}
var withOverloadSignature = overload1;

function f(n: () => void) { }

namespace m2 {
    export function foo(n: () => void ) {
    }

}

m2.foo(() =>  {

    var b = 30;
    return b;
});


declare function fooAmbient(n: number): string;

declare function overloadAmbient(n: number): string;
declare function overloadAmbient(s: string): string;

var f2 = () => {
    return "string";
}

//// [funcdecl.js]
"use strict";
function simpleFunc() {
    return "this is my simple func";
}
var simpleFuncVar = simpleFunc;
function anotherFuncNoReturn() {
}
var anotherFuncNoReturnVar = anotherFuncNoReturn;
function withReturn() {
    return "Hello";
}
var withReturnVar = withReturn;
function withParams(a) {
    return a;
}
var withparamsVar = withParams;
function withMultiParams(a, b, c) {
    return a;
}
var withMultiParamsVar = withMultiParams;
function withOptionalParams(a) {
}
var withOptionalParamsVar = withOptionalParams;
function withInitializedParams(a, b0, b = 30, c = "string value") {
}
var withInitializedParamsVar = withInitializedParams;
function withOptionalInitializedParams(a, c = "hello string") {
}
var withOptionalInitializedParamsVar = withOptionalInitializedParams;
function withRestParams(a, ...myRestParameter) {
    return myRestParameter;
}
var withRestParamsVar = withRestParams;
function overload1(ns) {
    return ns.toString();
}
var withOverloadSignature = overload1;
function f(n) { }
var m2;
(function (m2) {
    function foo(n) {
    }
    m2.foo = foo;
})(m2 || (m2 = {}));
m2.foo(() => {
    var b = 30;
    return b;
});
var f2 = () => {
    return "string";
};


//// [funcdecl.d.ts]
function simpleFunc(): string;
var simpleFuncVar: typeof simpleFunc;
function anotherFuncNoReturn(): void;
var anotherFuncNoReturnVar: typeof anotherFuncNoReturn;
function withReturn(): string;
var withReturnVar: typeof withReturn;
function withParams(a: string): string;
var withparamsVar: typeof withParams;
function withMultiParams(a: number, b: any, c: Object): number;
var withMultiParamsVar: typeof withMultiParams;
function withOptionalParams(a?: string): void;
var withOptionalParamsVar: typeof withOptionalParams;
function withInitializedParams(a: string, b0: any, b?: number, c?: string): void;
var withInitializedParamsVar: typeof withInitializedParams;
function withOptionalInitializedParams(a: string, c?: string): void;
var withOptionalInitializedParamsVar: typeof withOptionalInitializedParams;
function withRestParams(a: string, ...myRestParameter: number[]): number[];
var withRestParamsVar: typeof withRestParams;
function overload1(n: number): string;
function overload1(s: string): string;
var withOverloadSignature: typeof overload1;
function f(n: () => void): void;
namespace m2 {
    function foo(n: () => void): void;
}
function fooAmbient(n: number): string;
function overloadAmbient(n: number): string;
function overloadAmbient(s: string): string;
var f2: () => string;


//// [DtsFileErrors]


funcdecl.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== funcdecl.d.ts (1 errors) ====
    function simpleFunc(): string;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    var simpleFuncVar: typeof simpleFunc;
    function anotherFuncNoReturn(): void;
    var anotherFuncNoReturnVar: typeof anotherFuncNoReturn;
    function withReturn(): string;
    var withReturnVar: typeof withReturn;
    function withParams(a: string): string;
    var withparamsVar: typeof withParams;
    function withMultiParams(a: number, b: any, c: Object): number;
    var withMultiParamsVar: typeof withMultiParams;
    function withOptionalParams(a?: string): void;
    var withOptionalParamsVar: typeof withOptionalParams;
    function withInitializedParams(a: string, b0: any, b?: number, c?: string): void;
    var withInitializedParamsVar: typeof withInitializedParams;
    function withOptionalInitializedParams(a: string, c?: string): void;
    var withOptionalInitializedParamsVar: typeof withOptionalInitializedParams;
    function withRestParams(a: string, ...myRestParameter: number[]): number[];
    var withRestParamsVar: typeof withRestParams;
    function overload1(n: number): string;
    function overload1(s: string): string;
    var withOverloadSignature: typeof overload1;
    function f(n: () => void): void;
    namespace m2 {
        function foo(n: () => void): void;
    }
    function fooAmbient(n: number): string;
    function overloadAmbient(n: number): string;
    function overloadAmbient(s: string): string;
    var f2: () => string;
    