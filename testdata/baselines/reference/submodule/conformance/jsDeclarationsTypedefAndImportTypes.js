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
/**
 * @typedef {string | number} Whatever
 */
class Conn {
    constructor() { }
    item = 3;
    method() { }
}
module.exports = Conn;
//// [usage.js]
"use strict";
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
/**
 * @typedef {string | number} Whatever
 */
export type Whatever = string | number;
export = Conn;
//// [usage.d.ts]
/**
 * @typedef {import("./conn")} Conn
 */
export type Conn = import("./conn");
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


out/conn.d.ts(5,10): error TS2304: Cannot find name 'Conn'.


==== out/conn.d.ts (1 errors) ====
    /**
     * @typedef {string | number} Whatever
     */
    export type Whatever = string | number;
    export = Conn;
             ~~~~
!!! error TS2304: Cannot find name 'Conn'.
    
==== out/usage.d.ts (0 errors) ====
    /**
     * @typedef {import("./conn")} Conn
     */
    export type Conn = import("./conn");
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
    