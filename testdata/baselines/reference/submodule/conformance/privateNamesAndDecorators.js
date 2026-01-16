//// [tests/cases/conformance/classes/members/privateNames/privateNamesAndDecorators.ts] ////

//// [privateNamesAndDecorators.ts]
declare function dec<T>(target: T): T;

class A {
    @dec                // Error
    #foo = 1;
    @dec                // Error
    #bar(): void { }
}


//// [privateNamesAndDecorators.js]
var _A_foo;
class A {
    constructor() {
        _A_foo.set(this, 1);
    }
    #bar() { }
}
_A_foo = new WeakMap();
