//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsExportAssignedClassExpressionShadowing.ts] ////

//// [index.js]
class A {
    member = new Q();
}
class Q {
    x = 42;
}
module.exports = class Q {
    constructor() {
        this.x = new A();
    }
}
module.exports.Another = Q;


//// [index.js]
"use strict";
class A {
    member = new Q();
}
class Q {
    x = 42;
}
module.exports = class Q {
    constructor() {
        this.x = new A();
    }
};
module.exports.Another = Q;


//// [index.d.ts]
declare class Q {
    x: number;
}
declare const _default: {
    new (): Q;
};
export = _default;
export declare var Another: typeof Q;


!!!! File out/index.d.ts differs from original emit in noCheck emit
//// [index.d.ts]
--- Expected	The full check baseline
+++ Actual	with noCheck set
@@ -1,8 +1,13 @@
+declare class A {
+    member: Q;
+}
 declare class Q {
     x: number;
 }
 declare const _default: {
-    new (): Q;
+    new (): {
+        x: A;
+    };
 };
 export = _default;
 export declare var Another: typeof Q;