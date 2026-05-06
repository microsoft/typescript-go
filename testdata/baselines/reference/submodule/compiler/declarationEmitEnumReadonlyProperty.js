//// [tests/cases/compiler/declarationEmitEnumReadonlyProperty.ts] ////

//// [declarationEmitEnumReadonlyProperty.ts]
enum E {
    A = 'a',
    B = 'b'
}

class C {
    readonly type = E.A;
}

let x: E.A = new C().type;

//// [declarationEmitEnumReadonlyProperty.js]
"use strict";
var E;
(function (E) {
    E["A"] = "a";
    E["B"] = "b";
})(E || (E = {}));
class C {
    constructor() {
        this.type = E.A;
    }
}
let x = new C().type;


//// [declarationEmitEnumReadonlyProperty.d.ts]
enum E {
    A = "a",
    B = "b"
}
class C {
    readonly type = E.A;
}
let x: E.A;


//// [DtsFileErrors]


declarationEmitEnumReadonlyProperty.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitEnumReadonlyProperty.d.ts (1 errors) ====
    enum E {
    ~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        A = "a",
        B = "b"
    }
    class C {
        readonly type = E.A;
    }
    let x: E.A;
    