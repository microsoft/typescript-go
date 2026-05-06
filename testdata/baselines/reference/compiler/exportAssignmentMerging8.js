//// [tests/cases/compiler/exportAssignmentMerging8.ts] ////

//// [a.ts]
type SomeTypeAlias = { x: string };
class Foo {}
export = Foo;
export { SomeTypeAlias };
//// [b.ts]
import { SomeTypeAlias } from "./a";
const value: SomeTypeAlias = { x: "ok" };



//// [a.js]
"use strict";
class Foo {
}
module.exports = Foo;
//// [b.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const value = { x: "ok" };


//// [a.d.ts]
type SomeTypeAlias = {
    x: string;
};
class Foo {
}
export = Foo;
export { SomeTypeAlias };
//// [b.d.ts]
export {};


//// [DtsFileErrors]


a.d.ts(4,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== a.d.ts (1 errors) ====
    type SomeTypeAlias = {
        x: string;
    };
    class Foo {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    }
    export = Foo;
    export { SomeTypeAlias };
    
==== b.d.ts (0 errors) ====
    export {};
    