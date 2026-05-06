//// [tests/cases/conformance/classes/mixinAccessors4.ts] ////

//// [mixinAccessors4.ts]
// https://github.com/microsoft/TypeScript/issues/44938

class A {
  constructor(...args: any[]) {}
  get myName(): string {
    return "A";
  }
}

function Mixin<T extends typeof A>(Super: T) {
  return class B extends Super {
    get myName(): string {
      return "B";
    }
  };
}

class C extends Mixin(A) {
  get myName(): string {
    return "C";
  }
}


//// [mixinAccessors4.js]
"use strict";
// https://github.com/microsoft/TypeScript/issues/44938
class A {
    constructor(...args) { }
    get myName() {
        return "A";
    }
}
function Mixin(Super) {
    return class B extends Super {
        get myName() {
            return "B";
        }
    };
}
class C extends Mixin(A) {
    get myName() {
        return "C";
    }
}


//// [mixinAccessors4.d.ts]
class A {
    constructor(...args: any[]);
    get myName(): string;
}
function Mixin<T extends typeof A>(Super: T): {
    new (...args: any[]): {
        get myName(): string;
    };
} & T;
const C_base: {
    new (...args: any[]): {
        get myName(): string;
    };
} & typeof A;
class C extends C_base {
    get myName(): string;
}


//// [DtsFileErrors]


mixinAccessors4.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== mixinAccessors4.d.ts (1 errors) ====
    class A {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        constructor(...args: any[]);
        get myName(): string;
    }
    function Mixin<T extends typeof A>(Super: T): {
        new (...args: any[]): {
            get myName(): string;
        };
    } & T;
    const C_base: {
        new (...args: any[]): {
            get myName(): string;
        };
    } & typeof A;
    class C extends C_base {
        get myName(): string;
    }
    