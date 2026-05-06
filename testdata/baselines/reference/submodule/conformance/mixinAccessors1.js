//// [tests/cases/conformance/classes/mixinAccessors1.ts] ////

//// [mixinAccessors1.ts]
// https://github.com/microsoft/TypeScript/issues/58790

function mixin<T extends { new (...args: any[]): {} }>(superclass: T) {
  return class extends superclass {
    get validationTarget(): HTMLElement {
      return document.createElement("input");
    }
  };
}

class BaseClass {
  get validationTarget(): HTMLElement {
    return document.createElement("div");
  }
}

class MyClass extends mixin(BaseClass) {
  get validationTarget(): HTMLElement {
    return document.createElement("select");
  }
}

//// [mixinAccessors1.js]
"use strict";
// https://github.com/microsoft/TypeScript/issues/58790
function mixin(superclass) {
    return class extends superclass {
        get validationTarget() {
            return document.createElement("input");
        }
    };
}
class BaseClass {
    get validationTarget() {
        return document.createElement("div");
    }
}
class MyClass extends mixin(BaseClass) {
    get validationTarget() {
        return document.createElement("select");
    }
}


//// [mixinAccessors1.d.ts]
function mixin<T extends {
    new (...args: any[]): {};
}>(superclass: T): {
    new (...args: any[]): {
        get validationTarget(): HTMLElement;
    };
} & T;
class BaseClass {
    get validationTarget(): HTMLElement;
}
const MyClass_base: {
    new (...args: any[]): {
        get validationTarget(): HTMLElement;
    };
} & typeof BaseClass;
class MyClass extends MyClass_base {
    get validationTarget(): HTMLElement;
}


//// [DtsFileErrors]


mixinAccessors1.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== mixinAccessors1.d.ts (1 errors) ====
    function mixin<T extends {
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        new (...args: any[]): {};
    }>(superclass: T): {
        new (...args: any[]): {
            get validationTarget(): HTMLElement;
        };
    } & T;
    class BaseClass {
        get validationTarget(): HTMLElement;
    }
    const MyClass_base: {
        new (...args: any[]): {
            get validationTarget(): HTMLElement;
        };
    } & typeof BaseClass;
    class MyClass extends MyClass_base {
        get validationTarget(): HTMLElement;
    }
    