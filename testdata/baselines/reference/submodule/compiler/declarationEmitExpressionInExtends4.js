//// [tests/cases/compiler/declarationEmitExpressionInExtends4.ts] ////

//// [declarationEmitExpressionInExtends4.ts]
function getSomething() {
    return class D { }
}

class C extends getSomething()<number, string> {

}

class C2 extends SomeUndefinedFunction()<number, string> {

}


class C3 extends SomeUndefinedFunction {

}

//// [declarationEmitExpressionInExtends4.js]
"use strict";
function getSomething() {
    return class D {
    };
}
class C extends getSomething() {
}
class C2 extends SomeUndefinedFunction() {
}
class C3 extends SomeUndefinedFunction {
}


//// [declarationEmitExpressionInExtends4.d.ts]
function getSomething(): {
    new (): {};
};
const C_base: {
    new (): {};
};
class C extends C_base<number, string> {
}
const C2_base: any;
class C2 extends C2_base<number, string> {
}
class C3 extends SomeUndefinedFunction {
}
