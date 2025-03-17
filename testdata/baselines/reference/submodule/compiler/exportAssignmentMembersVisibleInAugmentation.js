//// [tests/cases/compiler/exportAssignmentMembersVisibleInAugmentation.ts] ////

//// [index.d.ts]
export = foo;
declare namespace foo {
    export type T = number;
}

//// [a.ts]
import * as foo from "foo";
declare module "foo" {
    export function f(): T; // OK
}

//// [b.ts]
import * as foo from "foo";
declare module "foo" {
    export function g(): foo.T; // OK
}


//// [b.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
