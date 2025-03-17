//// [tests/cases/compiler/constructorPropertyJs.ts] ////

//// [a.js]
class C {
    /**
     * @param {any} a
     */
    foo(a) {
        this.constructor = a;
    }
}


//// [a.js]
class C {
    foo(a) {
        this.constructor = a;
    }
}
