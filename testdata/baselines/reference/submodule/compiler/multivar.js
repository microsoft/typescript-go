//// [tests/cases/compiler/multivar.ts] ////

//// [multivar.ts]
var a,b,c;
var x=1,y=2,z=3;

module m2 {

    export var a, b2: number = 10, b;
    var m1;
    var a2, b22: number = 10, b222;
    var m3;

    class C {
        constructor (public b) {
        }
    }

    export class C2 {
        constructor (public b) {
        }
    }
    var m;
    declare var d1, d2;
    var b2;

    declare var v1;
}

var d;
var a22, b22 = 10, c22 = 30, dn;
var nn;

declare var da1, da2;
var normalVar;
declare var dv1;
var xl;
var x3;
var z4;

function foo(a2) {
    var a = 10;
}


for (var i = 0; i < 30; i++) {
    i++;
}
var b5 = 10;

//// [multivar.js]
var a, b, c;
var x = 1, y = 2, z = 3;
var m2;
(function (m2) {
    m2.b2 = 10;
    var m1;
    var a2, b22 = 10, b222;
    var m3;
    class C {
        b;
        constructor(b) {
            this.b = b;
        }
    }
    class C2 {
        b;
        constructor(b) {
            this.b = b;
        }
    }
    m2.C2 = C2;
    var m;
    var b2;
})(m2 || (m2 = {}));
var d;
var a22, b22 = 10, c22 = 30, dn;
var nn;
var normalVar;
var xl;
var x3;
var z4;
function foo(a2) {
    var a = 10;
}
for (var i = 0; i < 30; i++) {
    i++;
}
var b5 = 10;
