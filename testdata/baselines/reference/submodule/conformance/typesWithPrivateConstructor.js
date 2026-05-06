//// [tests/cases/conformance/types/members/typesWithPrivateConstructor.ts] ////

//// [typesWithPrivateConstructor.ts]
class C {
    private constructor() { }
}

var c = new C(); // error C is private
var r: () => void = c.constructor;

class C2 {
    private constructor(x: number);
    private constructor(x: any) { }
}

var c2 = new C2(); // error C2 is private
var r2: (x: number) => void = c2.constructor;

//// [typesWithPrivateConstructor.js]
"use strict";
class C {
    constructor() { }
}
var c = new C(); // error C is private
var r = c.constructor;
class C2 {
    constructor(x) { }
}
var c2 = new C2(); // error C2 is private
var r2 = c2.constructor;


//// [typesWithPrivateConstructor.d.ts]
class C {
    private constructor();
}
var c: C;
var r: () => void;
class C2 {
    private constructor();
}
var c2: C2;
var r2: (x: number) => void;
