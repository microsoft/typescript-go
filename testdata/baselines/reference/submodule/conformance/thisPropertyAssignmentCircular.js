//// [tests/cases/conformance/salsa/thisPropertyAssignmentCircular.ts] ////

//// [thisPropertyAssignmentCircular.js]
export class Foo {
    constructor() {
        this.foo = "Hello";
    }
    slicey() {
        this.foo = this.foo.slice();
    }
    m() {
        this.foo
    }
}

/** @class */
function C() {
    this.x = 0;
    this.x = function() { this.x.toString(); }
}


//// [thisPropertyAssignmentCircular.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Foo = void 0;
class Foo {
    constructor() {
        this.foo = "Hello";
    }
    slicey() {
        this.foo = this.foo.slice();
    }
    m() {
        this.foo;
    }
}
exports.Foo = Foo;
function C() {
    this.x = 0;
    this.x = function () { this.x.toString(); };
}
