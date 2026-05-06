//// [tests/cases/compiler/declarationEmitFirstTypeArgumentGenericFunctionType.ts] ////

//// [declarationEmitFirstTypeArgumentGenericFunctionType.ts]
class X<A> {
}
var prop11: X< <Tany>() => Tany >; // spaces before the first type argument
var prop12: X<(<Tany>() => Tany)>; // spaces before the first type argument
function f1() { // Inferred return type
    return prop11;
}
function f2() { // Inferred return type
    return prop12;
}
function f3(): X< <Tany>() => Tany> { // written with space before type argument
    return prop11;
}
function f4(): X<(<Tany>() => Tany)> { // written type with parenthesis
    return prop12;
}
class Y<A, B> {
}
var prop2: Y<string[], <Tany>() => Tany>; // No space after second type argument
var prop2: Y<string[], <Tany>() => Tany>; // space after second type argument
var prop3: Y< <Tany>() => Tany, <Tany>() => Tany>; // space before first type argument
var prop4: Y<(<Tany>() => Tany), <Tany>() => Tany>; // parenthesized first type argument


//// [declarationEmitFirstTypeArgumentGenericFunctionType.js]
"use strict";
class X {
}
var prop11; // spaces before the first type argument
var prop12; // spaces before the first type argument
function f1() {
    return prop11;
}
function f2() {
    return prop12;
}
function f3() {
    return prop11;
}
function f4() {
    return prop12;
}
class Y {
}
var prop2; // No space after second type argument
var prop2; // space after second type argument
var prop3; // space before first type argument
var prop4; // parenthesized first type argument


//// [declarationEmitFirstTypeArgumentGenericFunctionType.d.ts]
class X<A> {
}
var prop11: X<<Tany>() => Tany>;
var prop12: X<(<Tany>() => Tany)>;
function f1(): X<<Tany>() => Tany>;
function f2(): X<<Tany>() => Tany>;
function f3(): X<<Tany>() => Tany>;
function f4(): X<(<Tany>() => Tany)>;
class Y<A, B> {
}
var prop2: Y<string[], <Tany>() => Tany>;
var prop2: Y<string[], <Tany>() => Tany>;
var prop3: Y<<Tany>() => Tany, <Tany>() => Tany>;
var prop4: Y<(<Tany>() => Tany), <Tany>() => Tany>;


//// [DtsFileErrors]


declarationEmitFirstTypeArgumentGenericFunctionType.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitFirstTypeArgumentGenericFunctionType.d.ts (1 errors) ====
    class X<A> {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    }
    var prop11: X<<Tany>() => Tany>;
    var prop12: X<(<Tany>() => Tany)>;
    function f1(): X<<Tany>() => Tany>;
    function f2(): X<<Tany>() => Tany>;
    function f3(): X<<Tany>() => Tany>;
    function f4(): X<(<Tany>() => Tany)>;
    class Y<A, B> {
    }
    var prop2: Y<string[], <Tany>() => Tany>;
    var prop2: Y<string[], <Tany>() => Tany>;
    var prop3: Y<<Tany>() => Tany, <Tany>() => Tany>;
    var prop4: Y<(<Tany>() => Tany), <Tany>() => Tany>;
    