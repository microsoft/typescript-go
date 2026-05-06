//// [tests/cases/compiler/uniqueSymbolPropertyDeclarationEmit.ts] ////

//// [test.ts]
import Op from './op';
import { Po } from './po';

export default function foo() {
  return {
    [Op.or]: [],
    [Po.ro]: {}
  };
}

//// [op.ts]
declare const Op: {
  readonly or: unique symbol;
};

export default Op;

//// [po.d.ts]
export declare const Po: {
  readonly ro: unique symbol;
};


//// [op.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = Op;
//// [test.js]
"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = foo;
const op_1 = __importDefault(require("./op"));
const po_1 = require("./po");
function foo() {
    return {
        [op_1.default.or]: [],
        [po_1.Po.ro]: {}
    };
}


//// [op.d.ts]
const Op: {
    readonly or: unique symbol;
};
export default Op;
//// [test.d.ts]
import Op from './op';
import { Po } from './po';
export default function foo(): {
    [Op.or]: any[];
    [Po.ro]: {};
};


//// [DtsFileErrors]


op.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== test.d.ts (0 errors) ====
    import Op from './op';
    import { Po } from './po';
    export default function foo(): {
        [Op.or]: any[];
        [Po.ro]: {};
    };
    
==== op.d.ts (1 errors) ====
    const Op: {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        readonly or: unique symbol;
    };
    export default Op;
    
==== po.d.ts (0 errors) ====
    export declare const Po: {
      readonly ro: unique symbol;
    };
    