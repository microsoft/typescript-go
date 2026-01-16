//// [tests/cases/conformance/classes/members/privateNames/privateNamesInNestedClasses-2.ts] ////

//// [privateNamesInNestedClasses-2.ts]
class A {
    static #x = 5;
    constructor () {
        class B {
            #x = 5;
            constructor() {
                class C {
                    constructor() {
                        A.#x // error
                    }
                }
            }
        }
    }
}


//// [privateNamesInNestedClasses-2.js]
"use strict";
class A {
    static #x = 5;
    constructor() {
        var _B_x;
        class B {
            constructor() {
                _B_x.set(this, 5);
                class C {
                    constructor() {
                        A.#x; // error
                    }
                }
            }
        }
        _B_x = new WeakMap();
    }
}
