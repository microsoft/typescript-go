//// [tests/cases/conformance/classes/members/privateNames/privateNameDuplicateField.ts] ////

//// [privateNameDuplicateField.ts]
function Field() {

    // Error
    class A_Field_Field {
        #foo = "foo";
        #foo = "foo";
    }

    // Error
    class A_Field_Method {
        #foo = "foo";
        #foo() { }
    }

    // Error
    class A_Field_Getter {
        #foo = "foo";
        get #foo() { return ""}
    }

    // Error
    class A_Field_Setter {
        #foo = "foo";
        set #foo(value: string) { }
    }

    // Error
    class A_Field_StaticField {
        #foo = "foo";
        static #foo = "foo";
    }

    // Error
    class A_Field_StaticMethod {
        #foo = "foo";
        static #foo() { }
    }

    // Error
    class A_Field_StaticGetter {
        #foo = "foo";
        static get #foo() { return ""}
    }

    // Error
    class A_Field_StaticSetter {
        #foo = "foo";
        static set #foo(value: string) { }
    }
}

function Method() {
    // Error
    class A_Method_Field {
        #foo() { }
        #foo = "foo";
    }

    // Error
    class A_Method_Method {
        #foo() { }
        #foo() { }
    }

    // Error
    class A_Method_Getter {
        #foo() { }
        get #foo() { return ""}
    }

    // Error
    class A_Method_Setter {
        #foo() { }
        set #foo(value: string) { }
    }

    // Error
    class A_Method_StaticField {
        #foo() { }
        static #foo = "foo";
    }

    // Error
    class A_Method_StaticMethod {
        #foo() { }
        static #foo() { }
    }

    // Error
    class A_Method_StaticGetter {
        #foo() { }
        static get #foo() { return ""}
    }

    // Error
    class A_Method_StaticSetter {
        #foo() { }
        static set #foo(value: string) { }
    }
}


function Getter() {
    // Error
    class A_Getter_Field {
        get #foo() { return ""}
        #foo = "foo";
    }

    // Error
    class A_Getter_Method {
        get #foo() { return ""}
        #foo() { }
    }

    // Error
    class A_Getter_Getter {
        get #foo() { return ""}
        get #foo() { return ""}
    }

    //OK
    class A_Getter_Setter {
        get #foo() { return ""}
        set #foo(value: string) { }
    }

    // Error
    class A_Getter_StaticField {
        get #foo() { return ""}
        static #foo() { }
    }

    // Error
    class A_Getter_StaticMethod {
        get #foo() { return ""}
        static #foo() { }
    }

    // Error
    class A_Getter_StaticGetter {
        get #foo() { return ""}
        static get #foo() { return ""}
    }

    // Error
    class A_Getter_StaticSetter {
        get #foo() { return ""}
        static set #foo(value: string) { }
    }
}

function Setter() {
    // Error
    class A_Setter_Field {
        set #foo(value: string) { }
        #foo = "foo";
    }

    // Error
    class A_Setter_Method {
        set #foo(value: string) { }
        #foo() { }
    }

    // OK
    class A_Setter_Getter {
        set #foo(value: string) { }
        get #foo() { return ""}
    }

    // Error
    class A_Setter_Setter {
        set #foo(value: string) { }
        set #foo(value: string) { }
    }

    // Error
    class A_Setter_StaticField {
        set #foo(value: string) { }
        static #foo = "foo";
    }

    // Error
    class A_Setter_StaticMethod {
        set #foo(value: string) { }
        static #foo() { }
    }

    // Error
    class A_Setter_StaticGetter {
        set #foo(value: string) { }
        static get #foo() { return ""}
    }

    // Error
    class A_Setter_StaticSetter {
        set #foo(value: string) { }
        static set #foo(value: string) { }
    }
}

function StaticField() {
    // Error
    class A_StaticField_Field {
        static #foo = "foo";
        #foo = "foo";
    }

    // Error
    class A_StaticField_Method {
        static #foo = "foo";
        #foo() { }
    }

    // Error
    class A_StaticField_Getter {
        static #foo = "foo";
        get #foo() { return ""}
    }

    // Error
    class A_StaticField_Setter {
        static #foo = "foo";
        set #foo(value: string) { }
    }

    // Error
    class A_StaticField_StaticField {
        static #foo = "foo";
        static #foo = "foo";
    }

    // Error
    class A_StaticField_StaticMethod {
        static #foo = "foo";
        static #foo() { }
    }

    // Error
    class A_StaticField_StaticGetter {
        static #foo = "foo";
        static get #foo() { return ""}
    }

    // Error
    class A_StaticField_StaticSetter {
        static #foo = "foo";
        static set #foo(value: string) { }
    }
}

function StaticMethod() {
    // Error
    class A_StaticMethod_Field {
        static #foo() { }
        #foo = "foo";
    }

    // Error
    class A_StaticMethod_Method {
        static #foo() { }
        #foo() { }
    }

    // Error
    class A_StaticMethod_Getter {
        static #foo() { }
        get #foo() { return ""}
    }

    // Error
    class A_StaticMethod_Setter {
        static #foo() { }
        set #foo(value: string) { }
    }

    // Error
    class A_StaticMethod_StaticField {
        static #foo() { }
        static #foo = "foo";
    }

    // Error
    class A_StaticMethod_StaticMethod {
        static #foo() { }
        static #foo() { }
    }

    // Error
    class A_StaticMethod_StaticGetter {
        static #foo() { }
        static get #foo() { return ""}
    }

    // Error
    class A_StaticMethod_StaticSetter {
        static #foo() { }
        static set #foo(value: string) { }
    }
}

function StaticGetter() {

    // Error
    class A_StaticGetter_Field {
        static get #foo() { return ""}
        #foo = "foo";
    }

    // Error
    class A_StaticGetter_Method {
        static get #foo() { return ""}
        #foo() { }
    }

    // Error
    class A_StaticGetter_Getter {
        static get #foo() { return ""}
        get #foo() { return ""}
    }

    // Error
    class A_StaticGetter_Setter {
        static get #foo() { return ""}
        set #foo(value: string) { }
    }

    // Error
    class A_StaticGetter_StaticField {
        static get #foo() { return ""}
        static #foo() { }
    }

    // Error
    class A_StaticGetter_StaticMethod {
        static get #foo() { return ""}
        static #foo() { }
    }

    // Error
    class A_StaticGetter_StaticGetter {
        static get #foo() { return ""}
        static get #foo() { return ""}
    }
    // OK
    class A_StaticGetter_StaticSetter {
        static get #foo() { return ""}
        static set #foo(value: string) { }
    }
}

function StaticSetter() {
    // Error
    class A_StaticSetter_Field {
        static set #foo(value: string) { }
        #foo = "foo";
    }

    // Error
    class A_StaticSetter_Method {
        static set #foo(value: string) { }
        #foo() { }
    }


    // Error
    class A_StaticSetter_Getter {
        static set #foo(value: string) { }
        get #foo() { return ""}
    }

    // Error
    class A_StaticSetter_Setter {
        static set #foo(value: string) { }
        set #foo(value: string) { }
    }

    // Error
    class A_StaticSetter_StaticField {
        static set #foo(value: string) { }
        static #foo = "foo";
    }

    // Error
    class A_StaticSetter_StaticMethod {
        static set #foo(value: string) { }
        static #foo() { }
    }

    // OK
    class A_StaticSetter_StaticGetter {
        static set #foo(value: string) { }
        static get #foo() { return ""}
    }

    // Error
    class A_StaticSetter_StaticSetter {
        static set #foo(value: string) { }
        static set #foo(value: string) { }
    }
}


//// [privateNameDuplicateField.js]
"use strict";
function Field() {
    var _A_Field_Field_foo, _A_Field_Field_foo_1, _A_Field_Method_foo, _A_Field_Getter_foo, _A_Field_Setter_foo, _A_Field_StaticField_foo, _A_Field_StaticMethod_foo, _A_Field_StaticGetter_foo, _A_Field_StaticSetter_foo;
    // Error
    class A_Field_Field {
        constructor() {
            _A_Field_Field_foo_1.set(this, "foo");
            _A_Field_Field_foo_1.set(this, "foo");
        }
        #foo = "foo";
        #foo = "foo";
    }
    _A_Field_Field_foo = new WeakMap(), _A_Field_Field_foo_1 = new WeakMap();
    // Error
    class A_Field_Method {
        constructor() {
            _A_Field_Method_foo.set(this, "foo");
        }
        #foo() { }
    }
    _A_Field_Method_foo = new WeakMap();
    // Error
    class A_Field_Getter {
        constructor() {
            _A_Field_Getter_foo.set(this, "foo");
        }
        get #foo() { return ""; }
    }
    _A_Field_Getter_foo = new WeakMap();
    // Error
    class A_Field_Setter {
        constructor() {
            _A_Field_Setter_foo.set(this, "foo");
        }
        set #foo(value) { }
    }
    _A_Field_Setter_foo = new WeakMap();
    // Error
    class A_Field_StaticField {
        constructor() {
            _A_Field_StaticField_foo.set(this, "foo");
        }
        static #foo = "foo";
    }
    _A_Field_StaticField_foo = new WeakMap();
    // Error
    class A_Field_StaticMethod {
        constructor() {
            _A_Field_StaticMethod_foo.set(this, "foo");
        }
        static #foo() { }
    }
    _A_Field_StaticMethod_foo = new WeakMap();
    // Error
    class A_Field_StaticGetter {
        constructor() {
            _A_Field_StaticGetter_foo.set(this, "foo");
        }
        static get #foo() { return ""; }
    }
    _A_Field_StaticGetter_foo = new WeakMap();
    // Error
    class A_Field_StaticSetter {
        constructor() {
            _A_Field_StaticSetter_foo.set(this, "foo");
        }
        static set #foo(value) { }
    }
    _A_Field_StaticSetter_foo = new WeakMap();
}
function Method() {
    var _A_Method_Field_foo;
    // Error
    class A_Method_Field {
        constructor() {
            _A_Method_Field_foo.set(this, "foo");
        }
        #foo() { }
    }
    _A_Method_Field_foo = new WeakMap();
    // Error
    class A_Method_Method {
        #foo() { }
        #foo() { }
    }
    // Error
    class A_Method_Getter {
        #foo() { }
        get #foo() { return ""; }
    }
    // Error
    class A_Method_Setter {
        #foo() { }
        set #foo(value) { }
    }
    // Error
    class A_Method_StaticField {
        #foo() { }
        static #foo = "foo";
    }
    // Error
    class A_Method_StaticMethod {
        #foo() { }
        static #foo() { }
    }
    // Error
    class A_Method_StaticGetter {
        #foo() { }
        static get #foo() { return ""; }
    }
    // Error
    class A_Method_StaticSetter {
        #foo() { }
        static set #foo(value) { }
    }
}
function Getter() {
    var _A_Getter_Field_foo;
    // Error
    class A_Getter_Field {
        constructor() {
            _A_Getter_Field_foo.set(this, "foo");
        }
        get #foo() { return ""; }
    }
    _A_Getter_Field_foo = new WeakMap();
    // Error
    class A_Getter_Method {
        get #foo() { return ""; }
        #foo() { }
    }
    // Error
    class A_Getter_Getter {
        get #foo() { return ""; }
        get #foo() { return ""; }
    }
    //OK
    class A_Getter_Setter {
        get #foo() { return ""; }
        set #foo(value) { }
    }
    // Error
    class A_Getter_StaticField {
        get #foo() { return ""; }
        static #foo() { }
    }
    // Error
    class A_Getter_StaticMethod {
        get #foo() { return ""; }
        static #foo() { }
    }
    // Error
    class A_Getter_StaticGetter {
        get #foo() { return ""; }
        static get #foo() { return ""; }
    }
    // Error
    class A_Getter_StaticSetter {
        get #foo() { return ""; }
        static set #foo(value) { }
    }
}
function Setter() {
    var _A_Setter_Field_foo;
    // Error
    class A_Setter_Field {
        constructor() {
            _A_Setter_Field_foo.set(this, "foo");
        }
        set #foo(value) { }
    }
    _A_Setter_Field_foo = new WeakMap();
    // Error
    class A_Setter_Method {
        set #foo(value) { }
        #foo() { }
    }
    // OK
    class A_Setter_Getter {
        set #foo(value) { }
        get #foo() { return ""; }
    }
    // Error
    class A_Setter_Setter {
        set #foo(value) { }
        set #foo(value) { }
    }
    // Error
    class A_Setter_StaticField {
        set #foo(value) { }
        static #foo = "foo";
    }
    // Error
    class A_Setter_StaticMethod {
        set #foo(value) { }
        static #foo() { }
    }
    // Error
    class A_Setter_StaticGetter {
        set #foo(value) { }
        static get #foo() { return ""; }
    }
    // Error
    class A_Setter_StaticSetter {
        set #foo(value) { }
        static set #foo(value) { }
    }
}
function StaticField() {
    var _A_StaticField_Field_foo;
    // Error
    class A_StaticField_Field {
        constructor() {
            _A_StaticField_Field_foo.set(this, "foo");
        }
        static #foo = "foo";
    }
    _A_StaticField_Field_foo = new WeakMap();
    // Error
    class A_StaticField_Method {
        static #foo = "foo";
        #foo() { }
    }
    // Error
    class A_StaticField_Getter {
        static #foo = "foo";
        get #foo() { return ""; }
    }
    // Error
    class A_StaticField_Setter {
        static #foo = "foo";
        set #foo(value) { }
    }
    // Error
    class A_StaticField_StaticField {
        static #foo = "foo";
        static #foo = "foo";
    }
    // Error
    class A_StaticField_StaticMethod {
        static #foo = "foo";
        static #foo() { }
    }
    // Error
    class A_StaticField_StaticGetter {
        static #foo = "foo";
        static get #foo() { return ""; }
    }
    // Error
    class A_StaticField_StaticSetter {
        static #foo = "foo";
        static set #foo(value) { }
    }
}
function StaticMethod() {
    var _A_StaticMethod_Field_foo;
    // Error
    class A_StaticMethod_Field {
        constructor() {
            _A_StaticMethod_Field_foo.set(this, "foo");
        }
        static #foo() { }
    }
    _A_StaticMethod_Field_foo = new WeakMap();
    // Error
    class A_StaticMethod_Method {
        static #foo() { }
        #foo() { }
    }
    // Error
    class A_StaticMethod_Getter {
        static #foo() { }
        get #foo() { return ""; }
    }
    // Error
    class A_StaticMethod_Setter {
        static #foo() { }
        set #foo(value) { }
    }
    // Error
    class A_StaticMethod_StaticField {
        static #foo() { }
        static #foo = "foo";
    }
    // Error
    class A_StaticMethod_StaticMethod {
        static #foo() { }
        static #foo() { }
    }
    // Error
    class A_StaticMethod_StaticGetter {
        static #foo() { }
        static get #foo() { return ""; }
    }
    // Error
    class A_StaticMethod_StaticSetter {
        static #foo() { }
        static set #foo(value) { }
    }
}
function StaticGetter() {
    var _A_StaticGetter_Field_foo;
    // Error
    class A_StaticGetter_Field {
        constructor() {
            _A_StaticGetter_Field_foo.set(this, "foo");
        }
        static get #foo() { return ""; }
    }
    _A_StaticGetter_Field_foo = new WeakMap();
    // Error
    class A_StaticGetter_Method {
        static get #foo() { return ""; }
        #foo() { }
    }
    // Error
    class A_StaticGetter_Getter {
        static get #foo() { return ""; }
        get #foo() { return ""; }
    }
    // Error
    class A_StaticGetter_Setter {
        static get #foo() { return ""; }
        set #foo(value) { }
    }
    // Error
    class A_StaticGetter_StaticField {
        static get #foo() { return ""; }
        static #foo() { }
    }
    // Error
    class A_StaticGetter_StaticMethod {
        static get #foo() { return ""; }
        static #foo() { }
    }
    // Error
    class A_StaticGetter_StaticGetter {
        static get #foo() { return ""; }
        static get #foo() { return ""; }
    }
    // OK
    class A_StaticGetter_StaticSetter {
        static get #foo() { return ""; }
        static set #foo(value) { }
    }
}
function StaticSetter() {
    var _A_StaticSetter_Field_foo;
    // Error
    class A_StaticSetter_Field {
        constructor() {
            _A_StaticSetter_Field_foo.set(this, "foo");
        }
        static set #foo(value) { }
    }
    _A_StaticSetter_Field_foo = new WeakMap();
    // Error
    class A_StaticSetter_Method {
        static set #foo(value) { }
        #foo() { }
    }
    // Error
    class A_StaticSetter_Getter {
        static set #foo(value) { }
        get #foo() { return ""; }
    }
    // Error
    class A_StaticSetter_Setter {
        static set #foo(value) { }
        set #foo(value) { }
    }
    // Error
    class A_StaticSetter_StaticField {
        static set #foo(value) { }
        static #foo = "foo";
    }
    // Error
    class A_StaticSetter_StaticMethod {
        static set #foo(value) { }
        static #foo() { }
    }
    // OK
    class A_StaticSetter_StaticGetter {
        static set #foo(value) { }
        static get #foo() { return ""; }
    }
    // Error
    class A_StaticSetter_StaticSetter {
        static set #foo(value) { }
        static set #foo(value) { }
    }
}
