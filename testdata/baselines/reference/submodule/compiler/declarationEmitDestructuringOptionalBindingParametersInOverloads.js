//// [tests/cases/compiler/declarationEmitDestructuringOptionalBindingParametersInOverloads.ts] ////

//// [declarationEmitDestructuringOptionalBindingParametersInOverloads.ts]
function foo([x, y, z] ?: [string, number, boolean]);
function foo(...rest: any[]) {
}

function foo2( { x, y, z }?: { x: string; y: number; z: boolean });
function foo2(...rest: any[]) {

}

//// [declarationEmitDestructuringOptionalBindingParametersInOverloads.js]
"use strict";
function foo(...rest) {
}
function foo2(...rest) {
}


//// [declarationEmitDestructuringOptionalBindingParametersInOverloads.d.ts]
function foo([x, y, z]?: [string, number, boolean]): any;
function foo2({ x, y, z }?: {
    x: string;
    y: number;
    z: boolean;
}): any;


//// [DtsFileErrors]


declarationEmitDestructuringOptionalBindingParametersInOverloads.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitDestructuringOptionalBindingParametersInOverloads.d.ts (1 errors) ====
    function foo([x, y, z]?: [string, number, boolean]): any;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    function foo2({ x, y, z }?: {
        x: string;
        y: number;
        z: boolean;
    }): any;
    