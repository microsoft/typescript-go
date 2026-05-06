//// [tests/cases/compiler/declarationEmitForModuleImportingModuleAugmentationRetainsImport.ts] ////

//// [child1.ts]
import { ParentThing } from './parent';

declare module './parent' {
    interface ParentThing {
        add: (a: number, b: number) => number;
    }
}

export function child1(prototype: ParentThing) {
    prototype.add = (a: number, b: number) => a + b;
}

//// [parent.ts]
import { child1 } from './child1'; // this import should still exist in some form in the output, since it augments this module

export class ParentThing implements ParentThing {}

child1(ParentThing.prototype);

//// [parent.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ParentThing = void 0;
const child1_1 = require("./child1"); // this import should still exist in some form in the output, since it augments this module
class ParentThing {
}
exports.ParentThing = ParentThing;
(0, child1_1.child1)(ParentThing.prototype);
//// [child1.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.child1 = child1;
function child1(prototype) {
    prototype.add = (a, b) => a + b;
}


//// [parent.d.ts]
import './child1';
export class ParentThing implements ParentThing {
}
//// [child1.d.ts]
import { ParentThing } from './parent';
module './parent' {
    interface ParentThing {
        add: (a: number, b: number) => number;
    }
}
export function child1(prototype: ParentThing): void;


//// [DtsFileErrors]


child1.d.ts(2,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== child1.d.ts (1 errors) ====
    import { ParentThing } from './parent';
    module './parent' {
    ~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        interface ParentThing {
            add: (a: number, b: number) => number;
        }
    }
    export function child1(prototype: ParentThing): void;
    
==== parent.d.ts (0 errors) ====
    import './child1';
    export class ParentThing implements ParentThing {
    }
    