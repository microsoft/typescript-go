//// [tests/cases/compiler/es6ExportClause.ts] ////

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
class c {
}
var m;
(function (m) {
    m.x = 10;
})(m || (m = {}));
var x = 10;
export { c };
export { c as c2 };
export { m as instantiatedModule };
export { x };


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
    