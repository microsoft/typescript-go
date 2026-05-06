//// [tests/cases/compiler/declarationImportTypeAliasInferredAndEmittable.ts] ////

//// [foo.ts]
class Conn {
    constructor() { }
    item = 3;
    method() { }
}

export = Conn;
//// [usage.ts]
type Conn = import("./foo");
declare var x: Conn;

export class Wrap {
    connItem: number;
    constructor(c = x) {
        this.connItem = c.item;
    }
}


//// [foo.js]
"use strict";
class Conn {
    constructor() {
        this.item = 3;
    }
    method() { }
}
module.exports = Conn;
//// [usage.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Wrap = void 0;
class Wrap {
    constructor(c = x) {
        this.connItem = c.item;
    }
}
exports.Wrap = Wrap;


//// [foo.d.ts]
class Conn {
    constructor();
    item: number;
    method(): void;
}
export = Conn;
//// [usage.d.ts]
export class Wrap {
    connItem: number;
    constructor(c?: import("./foo"));
}


//// [DtsFileErrors]


foo.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== foo.d.ts (1 errors) ====
    class Conn {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        constructor();
        item: number;
        method(): void;
    }
    export = Conn;
    
==== usage.d.ts (0 errors) ====
    export class Wrap {
        connItem: number;
        constructor(c?: import("./foo"));
    }
    