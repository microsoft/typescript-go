//// [tests/cases/conformance/classes/members/privateNames/privateNameFieldClassExpression.ts] ////

//// [privateNameFieldClassExpression.ts]
class B {
    #foo = class {
        constructor() {
            console.log("hello");
        }
        static test = 123;
    };
    #foo2 = class Foo {
        static otherClass = 123;
    };
}




//// [privateNameFieldClassExpression.js]
var _B_foo, _B_foo2;
class B {
    constructor() {
        _B_foo.set(this, class {
            constructor() {
                console.log("hello");
            }
            static test = 123;
        });
        _B_foo2.set(this, class Foo {
            static otherClass = 123;
        });
    }
}
_B_foo = new WeakMap(), _B_foo2 = new WeakMap();
