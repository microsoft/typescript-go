//// [tests/cases/compiler/declarationEmitBindingPatterns.ts] ////

//// [declarationEmitBindingPatterns.ts]
const k = ({x: z = 'y'}) => { }

var a;
function f({} = a, [] = a, { p: {} = a} = a) {
}

//// [declarationEmitBindingPatterns.js]
"use strict";
const k = ({ x: z = 'y' }) => { };
var a;
function f({} = a, [] = a, { p: {} = a } = a) {
}


//// [declarationEmitBindingPatterns.d.ts]
const k: ({ x: z }: {
    x?: string;
}) => void;
var a: any;
function f({}?: any, []?: any, { p: {} }?: any): void;


//// [DtsFileErrors]


declarationEmitBindingPatterns.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitBindingPatterns.d.ts (1 errors) ====
    const k: ({ x: z }: {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        x?: string;
    }) => void;
    var a: any;
    function f({}?: any, []?: any, { p: {} }?: any): void;
    