//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsExportAssignedClassInstance2.ts] ////

//// [index.js]
class Foo {
    static stat = 10;
    member = 10;
}

module.exports = new Foo();

//// [index.js]
"use strict";
class Foo {
    constructor() {
        this.member = 10;
    }
}
Foo.stat = 10;
module.exports = new Foo();


//// [index.d.ts]
class Foo {
    static stat: number;
    member: number;
}
const _default: Foo;
export = _default;


//// [DtsFileErrors]


out/index.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== out/index.d.ts (1 errors) ====
    class Foo {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        static stat: number;
        member: number;
    }
    const _default: Foo;
    export = _default;
    