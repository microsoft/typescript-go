//// [tests/cases/compiler/declarationEmitDestructuringWithOptionalBindingParameters.ts] ////

//// [declarationEmitDestructuringWithOptionalBindingParameters.ts]
function foo([x,y,z]?: [string, number, boolean]) {
}
function foo1( { x, y, z }?: { x: string; y: number; z: boolean }) {
}

//// [declarationEmitDestructuringWithOptionalBindingParameters.js]
"use strict";
function foo([x, y, z]) {
}
function foo1({ x, y, z }) {
}


//// [declarationEmitDestructuringWithOptionalBindingParameters.d.ts]
function foo([x, y, z]?: [string, number, boolean]): void;
function foo1({ x, y, z }?: {
    x: string;
    y: number;
    z: boolean;
}): void;
