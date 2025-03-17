//// [tests/cases/compiler/declarationEmitOutFileBundlePaths.ts] ////

//// [versions.static.js]
export default {
    "@a/b": "1.0.0",
    "@a/c": "1.2.3"
};
//// [index.js]
import versions from './versions.static.js';

export {
    versions
};


//// [index.js]
"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.versions = void 0;
const versions_static_js_1 = __importDefault(require("./versions.static.js"));
exports.versions = versions_static_js_1.default;
//// [versions.static.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = {
    "@a/b": "1.0.0",
    "@a/c": "1.2.3"
};
