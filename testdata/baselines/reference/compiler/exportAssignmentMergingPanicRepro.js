//// [tests/cases/compiler/exportAssignmentMergingPanicRepro.ts] ////

//// [a.d.ts]
declare class Lib {
    value: number;
}
declare namespace Lib {
    class Component {
        render(): void;
    }
}
export = Lib;
//// [b.ts]
import { Component } from "./a";
class App extends Component {
    render() { }
}


//// [b.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const a_1 = require("./a");
class App extends a_1.Component {
    render() { }
}


//// [b.d.ts]
export {};
