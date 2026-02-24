//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsExportAssignedClassExpressionAnonymousWithSub.ts] ////

//// [index.js]
module.exports = class {
    /**
     * @param {number} p
     */
    constructor(p) {
        this.t = 12 + p;
    }
}
module.exports.Sub = class {
    constructor() {
        this.instance = new module.exports(10);
    }
}


//// [index.js]
"use strict";
module.exports = class {
    /**
     * @param {number} p
     */
    constructor(p) {
        this.t = 12 + p;
    }
};
module.exports.Sub = class {
    constructor() {
        this.instance = new module.exports(10);
    }
};


//// [index.d.ts]
declare const _default: any;
export = _default;
export declare var Sub: any;


!!!! File out/index.d.ts differs from original emit in noCheck emit
//// [index.d.ts]
--- Expected	The full check baseline
+++ Actual	with noCheck set
@@ -1,3 +1,13 @@
-declare const _default: any;
+declare const _default: {
+    new (p: number): {
+        t: number;
+    };
+};
 export = _default;
-export declare var Sub: any;
+export declare var Sub: {
+    new (): {
+        instance: {
+            t: number;
+        };
+    };
+};