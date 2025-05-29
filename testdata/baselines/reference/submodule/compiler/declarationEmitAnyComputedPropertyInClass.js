//// [tests/cases/compiler/declarationEmitAnyComputedPropertyInClass.ts] ////

//// [ambient.d.ts]
declare module "abcdefgh";

//// [main.ts]
import Test from "abcdefgh";

export class C {
    [Test.someKey]() {};
}


//// [main.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.C = void 0;
const abcdefgh_1 = require("abcdefgh");
class C {
    [abcdefgh_1.default.someKey]() { }
    ;
}
exports.C = C;


//// [main.d.ts]
import Test from "abcdefgh";
export declare class C {
    [x: number]: () => void;
    [Test.someKey](): void;
}


//// [DtsFileErrors]


main.d.ts(4,5): error TS1165: A computed property name in an ambient context must refer to an expression whose type is a literal type or a 'unique symbol' type.


==== ambient.d.ts (0 errors) ====
    declare module "abcdefgh";
    
==== main.d.ts (1 errors) ====
    import Test from "abcdefgh";
    export declare class C {
        [x: number]: () => void;
        [Test.someKey](): void;
        ~~~~~~~~~~~~~~
!!! error TS1165: A computed property name in an ambient context must refer to an expression whose type is a literal type or a 'unique symbol' type.
    }
    