//// [tests/cases/compiler/importHelpersOutFile.ts] ////

//// [a.ts]
export class A { }

//// [b.ts]
import { A } from "./a";
export class B extends A { }

//// [c.ts]
import { A } from "./a";
export class C extends A { }

//// [tslib.d.ts]
export declare function __extends(d: Function, b: Function): void;
export declare function __assign(t: any, ...sources: any[]): any;
export declare function __decorate(decorators: Function[], target: any, key?: string | symbol, desc?: any): any;
export declare function __param(paramIndex: number, decorator: Function): Function;
export declare function __metadata(metadataKey: any, metadataValue: any): Function;
export declare function __awaiter(thisArg: any, _arguments: any, P: Function, generator: Function): any;


//// [a.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.A = void 0;
class A {
}
exports.A = A;
//// [b.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.B = void 0;
const a_1 = require("./a");
class B extends a_1.A {
}
exports.B = B;
//// [c.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.C = void 0;
const a_1 = require("./a");
class C extends a_1.A {
}
exports.C = C;
