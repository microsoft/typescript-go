//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsFunctionKeywordProp.ts] ////

//// [source.js]
function foo() {}
foo.null = true;

function bar() {}
bar.async = true;
bar.normal = false;

function baz() {}
baz.class = true;
baz.normal = false;

//// [source.js]
function foo() { }
foo.null = true;
function bar() { }
bar.async = true;
bar.normal = false;
function baz() { }
baz.class = true;
baz.normal = false;


//// [source.d.ts]
declare function foo(): void;
declare function bar(): void;
declare function baz(): void;
declare namespace foo {
    const null_1: true;
    export { null_1 as null };
}
declare namespace bar {
    const async: true;
    const normal: false;
}
declare namespace baz {
    const class_1: true;
    export { class_1 as class };
    const normal: false;
}


!!!! File out/source.d.ts differs from original emit in noCheck emit
//// [source.d.ts]
--- Expected	The full check baseline
+++ Actual	with noCheck set
@@ -1,16 +1,16 @@
 declare function foo(): void;
 declare function bar(): void;
 declare function baz(): void;
+declare namespace baz {
+    const class_1: true;
+    export { class_1 as class };
+    const normal: false;
+}
 declare namespace foo {
     const null_1: true;
     export { null_1 as null };
 }
 declare namespace bar {
     const async: true;
-    const normal: false;
-}
-declare namespace baz {
-    const class_1: true;
-    export { class_1 as class };
     const normal: false;
 }