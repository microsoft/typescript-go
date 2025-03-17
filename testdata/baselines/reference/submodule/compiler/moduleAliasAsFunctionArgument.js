//// [tests/cases/compiler/moduleAliasAsFunctionArgument.ts] ////

//// [moduleAliasAsFunctionArgument_0.ts]
export var x: number;

//// [moduleAliasAsFunctionArgument_1.ts]
///<reference path='moduleAliasAsFunctionArgument_0.ts'/>
import a = require('moduleAliasAsFunctionArgument_0');

function fn(arg: { x: number }) {
}

a.x; // OK
fn(a); // Error: property 'x' is missing from 'a'


//// [moduleAliasAsFunctionArgument_1.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const a = require("moduleAliasAsFunctionArgument_0");
function fn(arg) {
}
a.x;
fn(a);
//// [moduleAliasAsFunctionArgument_0.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.x = void 0;
