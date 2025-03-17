//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsClassMethod.ts] ////

//// [jsDeclarationsClassMethod.js]
function C1() {
    /**
     * A comment prop
     * @param {number} x
     * @param {number} y
     * @returns {number}
     */
    this.prop = function (x, y) {
        return x + y;
    }
}

/**
 * A comment method
 * @param {number} x
 * @param {number} y
 * @returns {number}
 */
C1.prototype.method = function (x, y) {
    return x + y;
}

/**
 * A comment staticProp
 * @param {number} x
 * @param {number} y
 * @returns {number}
 */
C1.staticProp = function (x, y) {
    return x + y;
}

class C2 {
    /**
     * A comment method1
     * @param {number} x
     * @param {number} y
     * @returns {number}
     */
    method1(x, y) {
        return x + y;
    }
}

/**
 * A comment method2
 * @param {number} x
 * @param {number} y
 * @returns {number}
 */
C2.prototype.method2 = function (x, y) {
    return x + y;
}

/**
 * A comment staticProp
 * @param {number} x
 * @param {number} y
 * @returns {number}
 */
C2.staticProp = function (x, y) {
    return x + y;
}


//// [jsDeclarationsClassMethod.js]
function C1() {
    this.prop = function (x, y) {
        return x + y;
    };
}
C1.prototype.method = function (x, y) {
    return x + y;
};
C1.staticProp = function (x, y) {
    return x + y;
};
class C2 {
    method1(x, y) {
        return x + y;
    }
}
C2.prototype.method2 = function (x, y) {
    return x + y;
};
C2.staticProp = function (x, y) {
    return x + y;
};
