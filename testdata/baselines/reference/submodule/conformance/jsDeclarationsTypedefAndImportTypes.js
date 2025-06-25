//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsTypedefAndImportTypes.ts] ////

//// [conn.js]
/**
 * @typedef {string | number} Whatever
 */

class Conn {
    constructor() {}
    item = 3;
    method() {}
}

module.exports = Conn;

//// [usage.js]
/**
 * @typedef {import("./conn")} Conn
 */

class Wrap {
    /**
     * @param {Conn} c
     */
    constructor(c) {
        this.connItem = c.item;
        /** @type {import("./conn").Whatever} */
        this.another = "";
    }
}

module.exports = {
    Wrap
};


//// [conn.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
/**
 * @typedef {string | number} Whatever
 */
class Conn {
    constructor() { }
    item = 3;
    method() { }
}
export = Conn;
module.exports = Conn;
//// [usage.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
/**
 * @typedef {import("./conn")} Conn
 */
class Wrap {
    /**
     * @param {Conn} c
     */
    constructor(c) {
        this.connItem = c.item;
        /** @type {import("./conn").Whatever} */
        this.another = "";
    }
}
module.exports = {
    Wrap
};


//// [conn.d.ts]
export type Whatever = string | number;
export = Conn;
//// [usage.d.ts]
export type Conn = import("./conn");
/**
 * @typedef {import("./conn")} Conn
 */
declare class Wrap {
    /**
     * @param {Conn} c
     */
    constructor(c: Conn);
}
declare const _default: {
    Wrap: typeof Wrap;
};
export = _default;


//// [DtsFileErrors]


out/conn.d.ts(2,1): error TS2309: An export assignment cannot be used in a module with other exported elements.
out/conn.d.ts(2,10): error TS2304: Cannot find name 'Conn'.
out/usage.d.ts(1,20): error TS1340: Module './conn' does not refer to a type, but is used as a type here. Did you mean 'typeof import('./conn')'?
out/usage.d.ts(14,1): error TS2309: An export assignment cannot be used in a module with other exported elements.


==== out/conn.d.ts (2 errors) ====
    export type Whatever = string | number;
    export = Conn;
    ~~~~~~~~~~~~~~
!!! error TS2309: An export assignment cannot be used in a module with other exported elements.
             ~~~~
!!! error TS2304: Cannot find name 'Conn'.
    
==== out/usage.d.ts (2 errors) ====
    export type Conn = import("./conn");
                       ~~~~~~~~~~~~~~~~
!!! error TS1340: Module './conn' does not refer to a type, but is used as a type here. Did you mean 'typeof import('./conn')'?
    /**
     * @typedef {import("./conn")} Conn
     */
    declare class Wrap {
        /**
         * @param {Conn} c
         */
        constructor(c: Conn);
    }
    declare const _default: {
        Wrap: typeof Wrap;
    };
    export = _default;
    ~~~~~~~~~~~~~~~~~~
!!! error TS2309: An export assignment cannot be used in a module with other exported elements.
    