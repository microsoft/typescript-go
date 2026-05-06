//// [tests/cases/compiler/declFileTypeAnnotationBuiltInType.ts] ////

//// [declFileTypeAnnotationBuiltInType.ts]
// string
function foo(): string {
    return "";
}
function foo2() {
    return "";
}

// number
function foo3(): number {
    return 10;
}
function foo4() {
    return 10;
}

// boolean
function foo5(): boolean {
    return true;
}
function foo6() {
    return false;
}

// void
function foo7(): void {
    return;
}
function foo8() {
    return;
}

// any
function foo9(): any {
    return undefined;
}
function foo10() {
    return undefined;
}

//// [declFileTypeAnnotationBuiltInType.js]
"use strict";
// string
function foo() {
    return "";
}
function foo2() {
    return "";
}
// number
function foo3() {
    return 10;
}
function foo4() {
    return 10;
}
// boolean
function foo5() {
    return true;
}
function foo6() {
    return false;
}
// void
function foo7() {
    return;
}
function foo8() {
    return;
}
// any
function foo9() {
    return undefined;
}
function foo10() {
    return undefined;
}


//// [declFileTypeAnnotationBuiltInType.d.ts]
function foo(): string;
function foo2(): string;
function foo3(): number;
function foo4(): number;
function foo5(): boolean;
function foo6(): boolean;
function foo7(): void;
function foo8(): void;
function foo9(): any;
function foo10(): any;


//// [DtsFileErrors]


declFileTypeAnnotationBuiltInType.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileTypeAnnotationBuiltInType.d.ts (1 errors) ====
    function foo(): string;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    function foo2(): string;
    function foo3(): number;
    function foo4(): number;
    function foo5(): boolean;
    function foo6(): boolean;
    function foo7(): void;
    function foo8(): void;
    function foo9(): any;
    function foo10(): any;
    