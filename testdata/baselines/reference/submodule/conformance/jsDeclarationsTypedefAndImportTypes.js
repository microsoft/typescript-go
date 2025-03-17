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
class Conn {
    constructor() { }
    item = 3;
    method() { }
}
module.exports = Conn;
//// [usage.js]
class Wrap {
    constructor(c) {
        this.connItem = c.item;
        this.another = "";
    }
}
module.exports = {
    Wrap
};
