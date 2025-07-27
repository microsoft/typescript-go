//// [tests/cases/compiler/declarationEmitFunctionKeywordProp.ts] ////

//// [declarationEmitFunctionKeywordProp.ts]
function foo() {}
foo.null = true;

function bar() {}
bar.async = true;
bar.normal = false;

function baz() {}
baz.class = true;
baz.normal = false;

//// [declarationEmitFunctionKeywordProp.js]
function foo() { }
foo.null = true;
function bar() { }
bar.async = true;
bar.normal = false;
function baz() { }
baz.class = true;
baz.normal = false;


//// [declarationEmitFunctionKeywordProp.d.ts]
declare function foo(): void;
declare function bar(): void;
declare function baz(): void;
declare namespace bar {
    var async: true;
    var normal: false;
}
declare namespace baz {
    var class_1: true;
    export { class_1 as class };
    var normal: false;
}
declare namespace foo {
    var null_1: true;
    export { null_1 as null };
}


!!!! File declarationEmitFunctionKeywordProp.d.ts differs from original emit in noCheck emit
//// [declarationEmitFunctionKeywordProp.d.ts]
--- Expected	The full check baseline
+++ Actual	with noCheck set
@@ -1,6 +1,10 @@
 declare function foo(): void;
 declare function bar(): void;
 declare function baz(): void;
+declare namespace foo {
+    var null_1: true;
+    export { null_1 as null };
+}
 declare namespace bar {
     var async: true;
     var normal: false;
@@ -9,8 +13,4 @@
     var class_1: true;
     export { class_1 as class };
     var normal: false;
-}
-declare namespace foo {
-    var null_1: true;
-    export { null_1 as null };
 }