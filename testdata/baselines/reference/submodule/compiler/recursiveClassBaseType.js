//// [tests/cases/compiler/recursiveClassBaseType.ts] ////

//// [recursiveClassBaseType.ts]
// Repro from #44281

declare const p: <T>(fn: () => T) => T;

declare const Base: <T>(val: T) => { new(): T };

class C extends Base({ x: p<C[]>(() => []) }) { }

// Repro from #44359

abstract class Base1 {
    abstract root(): Derived1;
}

class Derived1 extends class extends Base1 {
    root() {
        return undefined as any;
    }
}
{ }


//// [recursiveClassBaseType.js]
"use strict";
// Repro from #44281
class C extends Base({ x: p(() => []) }) {
}
// Repro from #44359
class Base1 {
}
class Derived1 extends class extends Base1 {
    root() {
        return undefined;
    }
} {
}


//// [recursiveClassBaseType.d.ts]
const p: <T>(fn: () => T) => T;
const Base: <T>(val: T) => {
    new (): T;
};
const C_base: new () => {
    x: C[];
};
class C extends C_base {
}
abstract class Base1 {
    abstract root(): Derived1;
}
const Derived1_base: {
    new (): {
        root(): any;
    };
};
class Derived1 extends Derived1_base {
}


//// [DtsFileErrors]


recursiveClassBaseType.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== recursiveClassBaseType.d.ts (1 errors) ====
    const p: <T>(fn: () => T) => T;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    const Base: <T>(val: T) => {
        new (): T;
    };
    const C_base: new () => {
        x: C[];
    };
    class C extends C_base {
    }
    abstract class Base1 {
        abstract root(): Derived1;
    }
    const Derived1_base: {
        new (): {
            root(): any;
        };
    };
    class Derived1 extends Derived1_base {
    }
    