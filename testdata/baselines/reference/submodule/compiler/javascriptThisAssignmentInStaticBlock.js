//// [tests/cases/compiler/javascriptThisAssignmentInStaticBlock.ts] ////

//// [a.js]
class Thing {
    static {
        this.doSomething = () => {};
    }
}

Thing.doSomething();

// GH#46468
class ElementsArray extends Array {
    static {
        const superisArray = super.isArray;
        const customIsArray = (arg)=> superisArray(arg);
        this.isArray = customIsArray;
    }
}

ElementsArray.isArray(new ElementsArray());

//// [a.js]
"use strict";
var _a, _b, _c;
class Thing {
}
_a = Thing;
(() => {
    _a.doSomething = () => { };
})();
Thing.doSomething();
// GH#46468
class ElementsArray extends (_c = Array) {
}
_b = ElementsArray;
(() => {
    const superisArray = Reflect.get(_c, "isArray", _b);
    const customIsArray = (arg) => superisArray(arg);
    _b.isArray = customIsArray;
})();
ElementsArray.isArray(new ElementsArray());


//// [a.d.ts]
class Thing {
    static doSomething: () => void;
}
class ElementsArray extends Array {
    static isArray: (arg: any) => arg is any[];
}


//// [DtsFileErrors]


/out/a.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== /out/a.d.ts (1 errors) ====
    class Thing {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        static doSomething: () => void;
    }
    class ElementsArray extends Array {
        static isArray: (arg: any) => arg is any[];
    }
    