//// [tests/cases/compiler/es6ImportDefaultBindingFollowedWithNamespaceBinding1InEs5.ts] ////

//// [es6ImportDefaultBindingFollowedWithNamespaceBindingInEs5_0.ts]
var a = 10;
export default a;

//// [es6ImportDefaultBindingFollowedWithNamespaceBindingInEs5_1.ts]
import defaultBinding, * as nameSpaceBinding  from "./es6ImportDefaultBindingFollowedWithNamespaceBindingInEs5_0";
var x: number = defaultBinding;

//// [es6ImportDefaultBindingFollowedWithNamespaceBindingInEs5_0.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var a = 10;
exports.default = a;
//// [es6ImportDefaultBindingFollowedWithNamespaceBindingInEs5_1.js]
"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const es6ImportDefaultBindingFollowedWithNamespaceBindingInEs5_0_1 = __importDefault(require("./es6ImportDefaultBindingFollowedWithNamespaceBindingInEs5_0"));
var x = es6ImportDefaultBindingFollowedWithNamespaceBindingInEs5_0_1.default;


//// [es6ImportDefaultBindingFollowedWithNamespaceBindingInEs5_0.d.ts]
var a: number;
export default a;
//// [es6ImportDefaultBindingFollowedWithNamespaceBindingInEs5_1.d.ts]
export {};


//// [DtsFileErrors]


es6ImportDefaultBindingFollowedWithNamespaceBindingInEs5_0.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== es6ImportDefaultBindingFollowedWithNamespaceBindingInEs5_0.d.ts (1 errors) ====
    var a: number;
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export default a;
    
==== es6ImportDefaultBindingFollowedWithNamespaceBindingInEs5_1.d.ts (0 errors) ====
    export {};
    