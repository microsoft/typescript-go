//// [tests/cases/conformance/classes/mixinAccessors3.ts] ////

//// [mixinAccessors3.ts]
function mixin<T extends { new (...args: any[]): {} }>(superclass: T) {
  return class extends superclass {
    get name() {
      return "";
    }
  };
}

class BaseClass {
  set name(v: string) {}
}

// error
class MyClass extends mixin(BaseClass) { 
  get name() {
    return "";
  }
}


//// [mixinAccessors3.js]
"use strict";
function mixin(superclass) {
    return class extends superclass {
        get name() {
            return "";
        }
    };
}
class BaseClass {
    set name(v) { }
}
// error
class MyClass extends mixin(BaseClass) {
    get name() {
        return "";
    }
}


//// [mixinAccessors3.d.ts]
function mixin<T extends {
    new (...args: any[]): {};
}>(superclass: T): {
    new (...args: any[]): {
        get name(): string;
    };
} & T;
class BaseClass {
    set name(v: string);
}
const MyClass_base: {
    new (...args: any[]): {
        get name(): string;
    };
} & typeof BaseClass;
class MyClass extends MyClass_base {
    get name(): string;
}
