//// [tests/cases/compiler/declarationEmitExpressionInExtends.ts] ////

//// [declarationEmitExpressionInExtends.ts]
var x: {
    new<T>(s: any): Q;
}

class Q {
    s: string;    
}

class B extends x<string> {    
}

var q: B;
q.s;

//// [declarationEmitExpressionInExtends.js]
"use strict";
var x;
class Q {
}
class B extends x {
}
var q;
q.s;


//// [declarationEmitExpressionInExtends.d.ts]
var x: {
    new <T>(s: any): Q;
};
class Q {
    s: string;
}
class B extends x<string> {
}
var q: B;
