//// [tests/cases/compiler/es6ExportClauseInEs5.ts] ////

//// [server.ts]
class c {
}
interface i {
}
namespace m {
    export var x = 10;
}
var x = 10;
namespace uninstantiated {
}
export { c };
export { c as c2 };
export { i, m as instantiatedModule };
export { uninstantiated };
export { x };

//// [server.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.x = exports.instantiatedModule = exports.c2 = exports.c = void 0;
class c {
}
exports.c = c;
exports.c2 = c;
var m;
(function (m) {
    m.x = 10;
})(m || (exports.instantiatedModule = m = {}));
var x = 10;
exports.x = x;


//// [server.d.ts]
class c {
}
interface i {
}
namespace m {
    var x: number;
}
var x: number;
namespace uninstantiated {
}
export { c };
export { c as c2 };
export { i, m as instantiatedModule };
export { uninstantiated };
export { x };


//// [DtsFileErrors]


server.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== server.d.ts (1 errors) ====
    class c {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    }
    interface i {
    }
    namespace m {
        var x: number;
    }
    var x: number;
    namespace uninstantiated {
    }
    export { c };
    export { c as c2 };
    export { i, m as instantiatedModule };
    export { uninstantiated };
    export { x };
    