//// [tests/cases/conformance/types/stringLiteral/stringLiteralTypesOverloads03.ts] ////

//// [stringLiteralTypesOverloads03.ts]
interface Base {
    x: string;
    y: number;
}

interface HelloOrWorld extends Base {
    p1: boolean;
}

interface JustHello extends Base {
    p2: boolean;
}

interface JustWorld extends Base {
    p3: boolean;
}

let hello: "hello";
let world: "world";
let helloOrWorld: "hello" | "world";

function f(p: "hello"): JustHello;
function f(p: "hello" | "world"): HelloOrWorld;
function f(p: "world"): JustWorld;
function f(p: string): Base;
function f(...args: any[]): any {
    return undefined;
}

let fResult1 = f(hello);
let fResult2 = f(world);
let fResult3 = f(helloOrWorld);

function g(p: string): Base;
function g(p: "hello"): JustHello;
function g(p: "hello" | "world"): HelloOrWorld;
function g(p: "world"): JustWorld;
function g(...args: any[]): any {
    return undefined;
}

let gResult1 = g(hello);
let gResult2 = g(world);
let gResult3 = g(helloOrWorld);

//// [stringLiteralTypesOverloads03.js]
"use strict";
let hello;
let world;
let helloOrWorld;
function f(...args) {
    return undefined;
}
let fResult1 = f(hello);
let fResult2 = f(world);
let fResult3 = f(helloOrWorld);
function g(...args) {
    return undefined;
}
let gResult1 = g(hello);
let gResult2 = g(world);
let gResult3 = g(helloOrWorld);


//// [stringLiteralTypesOverloads03.d.ts]
interface Base {
    x: string;
    y: number;
}
interface HelloOrWorld extends Base {
    p1: boolean;
}
interface JustHello extends Base {
    p2: boolean;
}
interface JustWorld extends Base {
    p3: boolean;
}
let hello: "hello";
let world: "world";
let helloOrWorld: "hello" | "world";
function f(p: "hello"): JustHello;
function f(p: "hello" | "world"): HelloOrWorld;
function f(p: "world"): JustWorld;
function f(p: string): Base;
let fResult1: JustHello;
let fResult2: JustWorld;
let fResult3: HelloOrWorld;
function g(p: string): Base;
function g(p: "hello"): JustHello;
function g(p: "hello" | "world"): HelloOrWorld;
function g(p: "world"): JustWorld;
let gResult1: JustHello;
let gResult2: JustWorld;
let gResult3: Base;
