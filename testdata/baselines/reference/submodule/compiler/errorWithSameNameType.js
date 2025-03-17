//// [tests/cases/compiler/errorWithSameNameType.ts] ////

//// [a.ts]
export interface F {
    foo1: number
}

//// [b.ts]
export interface F {
    foo2: number
}

//// [c.ts]
import * as A from './a'
import * as B from './b'

let a: A.F
let b: B.F

if (a === b) {

}

a = b


//// [c.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
let a;
let b;
if (a === b) {
}
a = b;
//// [b.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
