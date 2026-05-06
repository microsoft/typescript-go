//// [tests/cases/compiler/jsDeclarationsRequireImportForms.ts] ////

//// [obj.js]
class Obj {
    constructor() {
        this.x = 12;
    }
}
module.exports.Obj = Obj
//// [index.js]
const {Obj, Obj: Other} = require("./obj");

class Container {
    constructor() {
        this.usage = new Obj();
        /** @type {Other} */
        this.usage2 = new Other();
    }
}

module.exports = Container;

//// [obj.js]
"use strict";
class Obj {
    constructor() {
        this.x = 12;
    }
}
module.exports.Obj = Obj;
//// [index.js]
"use strict";
const { Obj, Obj: Other } = require("./obj");
class Container {
    constructor() {
        this.usage = new Obj();
        /** @type {Other} */
        this.usage2 = new Other();
    }
}
module.exports = Container;


//// [obj.d.ts]
class Obj {
    x: number;
    constructor();
}
export { Obj };
//// [index.d.ts]
import { Obj, Obj as Other } from "./obj";
class Container {
    usage: Obj;
    /** @type {Other} */
    usage2: Other;
    constructor();
}
export = Container;


//// [DtsFileErrors]


out/index.d.ts(2,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
out/obj.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== out/index.d.ts (1 errors) ====
    import { Obj, Obj as Other } from "./obj";
    class Container {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        usage: Obj;
        /** @type {Other} */
        usage2: Other;
        constructor();
    }
    export = Container;
    
==== out/obj.d.ts (1 errors) ====
    class Obj {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        x: number;
        constructor();
    }
    export { Obj };
    