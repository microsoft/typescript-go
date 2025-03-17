//// [tests/cases/compiler/moduleAugmentationOfAlias.ts] ////

//// [a.ts]
interface I {}
export default I;

//// [b.ts]
export {};
declare module './a' {
    export default interface I { x: number; }
}

//// [c.ts]
import I from "./a";
function f(i: I) {
    i.x;
}


//// [c.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
function f(i) {
    i.x;
}
//// [b.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
