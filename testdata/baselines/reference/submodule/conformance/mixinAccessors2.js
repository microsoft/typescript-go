//// [tests/cases/conformance/classes/mixinAccessors2.ts] ////

//// [mixinAccessors2.ts]
function mixin<T extends { new (...args: any[]): {} }>(superclass: T) {
  return class extends superclass {
    accessor name = "";
  };
}

class BaseClass {
  accessor name = "";
}

class MyClass extends mixin(BaseClass) {
  accessor name = "";
}


//// [mixinAccessors2.js]
"use strict";
function mixin(superclass) {
    return class extends superclass {
        accessor name = "";
    };
}
class BaseClass {
    accessor name = "";
}
class MyClass extends mixin(BaseClass) {
    accessor name = "";
}


//// [mixinAccessors2.d.ts]
function mixin<T extends {
    new (...args: any[]): {};
}>(superclass: T): {
    new (...args: any[]): {
        get name(): string;
        set name(arg: string);
    };
} & T;
class BaseClass {
    accessor name: string;
}
const MyClass_base: {
    new (...args: any[]): {
        get name(): string;
        set name(arg: string);
    };
} & typeof BaseClass;
class MyClass extends MyClass_base {
    accessor name: string;
}


//// [DtsFileErrors]


mixinAccessors2.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== mixinAccessors2.d.ts (1 errors) ====
    function mixin<T extends {
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        new (...args: any[]): {};
    }>(superclass: T): {
        new (...args: any[]): {
            get name(): string;
            set name(arg: string);
        };
    } & T;
    class BaseClass {
        accessor name: string;
    }
    const MyClass_base: {
        new (...args: any[]): {
            get name(): string;
            set name(arg: string);
        };
    } & typeof BaseClass;
    class MyClass extends MyClass_base {
        accessor name: string;
    }
    