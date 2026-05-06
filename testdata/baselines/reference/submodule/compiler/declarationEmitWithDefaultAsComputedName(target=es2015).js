//// [tests/cases/compiler/declarationEmitWithDefaultAsComputedName.ts] ////

//// [other.ts]
type Experiment<Name> = {
    name: Name;
};
declare const createExperiment: <Name extends string>(
    options: Experiment<Name>
) => Experiment<Name>;
export default createExperiment({
    name: "foo"
});

//// [main.ts]
import other from "./other";
export const obj = {
    [other.name]: 1,
};

//// [other.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = createExperiment({
    name: "foo"
});
//// [main.js]
"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.obj = void 0;
const other_1 = __importDefault(require("./other"));
exports.obj = {
    [other_1.default.name]: 1,
};


//// [other.d.ts]
type Experiment<Name> = {
    name: Name;
};
const _default: Experiment<"foo">;
export default _default;
//// [main.d.ts]
export const obj: {
    foo: number;
};


//// [DtsFileErrors]


other.d.ts(4,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== other.d.ts (1 errors) ====
    type Experiment<Name> = {
        name: Name;
    };
    const _default: Experiment<"foo">;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export default _default;
    
==== main.d.ts (0 errors) ====
    export const obj: {
        foo: number;
    };
    