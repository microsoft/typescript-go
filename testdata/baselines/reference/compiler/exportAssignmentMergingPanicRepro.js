//// [tests/cases/compiler/exportAssignmentMergingPanicRepro.ts] ////

//// [a.ts]
type Constructor<T = {}> = new (...args: any[]) => T;
interface Printable {
    print(): void;
}
function Mixin<TBase extends Constructor>(Base: TBase) {
    return class extends Base implements Printable {
        print() {}
    };
}
class CoreBase {
    id: number = 0;
    static Printable: Printable = { print() {} };
}
const Mixed = Mixin(CoreBase);
export = Mixed;
export { Printable };
//// [b.ts]
import Mixed = require("./a");
import { Printable } from "./a";
class App extends Mixed {
    doPrint(p: Printable) {
        p.print();
    }
}


//// [a.js]
"use strict";
function Mixin(Base) {
    return class extends Base {
        print() { }
    };
}
class CoreBase {
    id = 0;
    static Printable = { print() { } };
}
const Mixed = Mixin(CoreBase);
module.exports = Mixed;
//// [b.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const Mixed = require("./a");
class App extends Mixed {
    doPrint(p) {
        p.print();
    }
}


//// [a.d.ts]
interface Printable {
    print(): void;
}
declare class CoreBase {
    id: number;
    static Printable: Printable;
}
declare const Mixed: {
    new (...args: any[]): {
        print(): void;
    };
} & typeof CoreBase;
export = Mixed;
export { Printable };
//// [b.d.ts]
export {};
