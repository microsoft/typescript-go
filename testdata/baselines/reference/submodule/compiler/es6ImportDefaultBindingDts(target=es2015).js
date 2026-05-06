//// [tests/cases/compiler/es6ImportDefaultBindingDts.ts] ////

//// [server.ts]
class c { }
export default c;

//// [client.ts]
import defaultBinding from "./server";
export var x = new defaultBinding();
import defaultBinding2 from "./server"; // elide this import since defaultBinding2 is not used


//// [server.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
class c {
}
exports.default = c;
//// [client.js]
"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.x = void 0;
const server_1 = __importDefault(require("./server"));
exports.x = new server_1.default();


//// [server.d.ts]
class c {
}
export default c;
//// [client.d.ts]
import defaultBinding from "./server";
export var x: defaultBinding;


//// [DtsFileErrors]


server.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== server.d.ts (1 errors) ====
    class c {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    }
    export default c;
    
==== client.d.ts (0 errors) ====
    import defaultBinding from "./server";
    export var x: defaultBinding;
    