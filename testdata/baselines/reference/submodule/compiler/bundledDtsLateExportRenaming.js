//// [tests/cases/compiler/bundledDtsLateExportRenaming.ts] ////

//// [index.ts]
export * from "./nested";

//// [base.ts]
import { B } from "./shared";

export function f() {
    return new B();
}

//// [derived.ts]
import { f } from "./base";

export function g() {
    return f();
}

//// [index.ts]
export * from "./base";

export * from "./derived";
export * from "./shared";

//// [shared.ts]
export class B {}


//// [index.js]
"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __exportStar = (this && this.__exportStar) || function(m, exports) {
    for (var p in m) if (p !== "default" && !Object.prototype.hasOwnProperty.call(exports, p)) __createBinding(exports, m, p);
};
Object.defineProperty(exports, "__esModule", { value: true });
__exportStar(require("./nested"), exports);
//// [index.js]
"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __exportStar = (this && this.__exportStar) || function(m, exports) {
    for (var p in m) if (p !== "default" && !Object.prototype.hasOwnProperty.call(exports, p)) __createBinding(exports, m, p);
};
Object.defineProperty(exports, "__esModule", { value: true });
__exportStar(require("./base"), exports);
__exportStar(require("./derived"), exports);
__exportStar(require("./shared"), exports);
//// [derived.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.g = g;
const base_1 = require("./base");
function g() {
    return (0, base_1.f)();
}
//// [base.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.f = f;
const shared_1 = require("./shared");
function f() {
    return new shared_1.B();
}
//// [shared.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.B = void 0;
class B {
}
exports.B = B;
