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
    /**
     * A comment prop
     * @param {number} x
     * @param {number} y
     * @returns {number}
     */
    this.prop = function (x, y) {
        return x + y;
    };
}
/**
 * A comment method
 * @param {number} x
 * @param {number} y
 * @returns {number}
 */
C1.prototype.method = function (x, y) {
    return x + y;
};
/**
 * A comment staticProp
 * @param {number} x
 * @param {number} y
 * @returns {number}
 */
C1.staticProp = function (x, y) {
    return x + y;
};
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
};
/**
 * A comment staticProp
 * @param {number} x
 * @param {number} y
 * @returns {number}
 */
C2.staticProp = function (x, y) {
    return x + y;
};


//// [jsDeclarationsClassMethod.d.ts]
declare function C1(): void;
declare class C2 {
    /**
     * A comment method1
     * @param {number} x
     * @param {number} y
     * @returns {number}
     */
    method1(x: number, y: number): number;
}
declare namespace C2 {
    const staticProp: (x: any, y: any) => any;
}
declare namespace C1 {
    const staticProp: (x: any, y: any) => any;
}


!!!! File out/jsDeclarationsClassMethod.d.ts differs from original emit in noCheck emit
//// [jsDeclarationsClassMethod.d.ts]
--- Expected	The full check baseline
+++ Actual	with noCheck set
@@ -8,9 +8,9 @@
      */
     method1(x: number, y: number): number;
 }
-declare namespace C2 {
+declare namespace C1 {
     const staticProp: (x: any, y: any) => any;
 }
-declare namespace C1 {
+declare namespace C2 {
     const staticProp: (x: any, y: any) => any;
 }