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
    constructor() {
        this.item = 3;
    }
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
class Conn {
    constructor();
    item: number;
    method(): void;
}
export = Conn;
//// [usage.d.ts]
/**
 * @typedef {import("./conn")} Conn
 */
export type Conn = import("./conn");
class Wrap {
    connItem: number;
    /** @type {import("./conn").Whatever} */
    another: import("./conn").Whatever;
    /**
     * @param {Conn} c
     */
    constructor(c: Conn);
}
const _default: {
    Wrap: typeof Wrap;
};
export = _default;


//// [DtsFileErrors]


out/conn.d.ts(5,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
out/usage.d.ts(5,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== out/conn.d.ts (1 errors) ====
    /**
     * @typedef {string | number} Whatever
     */
    export type Whatever = string | number;
    class Conn {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        constructor();
        item: number;
        method(): void;
    }
    export = Conn;
    
==== out/usage.d.ts (1 errors) ====
    /**
     * @typedef {import("./conn")} Conn
     */
    export type Conn = import("./conn");
    class Wrap {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        connItem: number;
        /** @type {import("./conn").Whatever} */
        another: import("./conn").Whatever;
        /**
         * @param {Conn} c
         */
        constructor(c: Conn);
    }
    const _default: {
        Wrap: typeof Wrap;
    };
    export = _default;
    