//// [tests/cases/conformance/declarationEmit/anonymousClassAccessorsDeclarationEmit1.ts] ////

//// [anonymousClassAccessorsDeclarationEmit1.ts]
export abstract class Base {
  accessor a = 1;
}

export function middle(Super = Base) {
  abstract class Middle extends Super {}
  return Middle;
}

class A {
  constructor(...args: any[]) {}
}

export function Mixin<T extends typeof A>(Super: T) {
  return class B extends Super {
    get myName(): string {
      return "B";
    }
    set myName(arg: string) {}
  };
}


//// [anonymousClassAccessorsDeclarationEmit1.js]
export class Base {
    accessor a = 1;
}
export function middle(Super = Base) {
    class Middle extends Super {
    }
    return Middle;
}
class A {
    constructor(...args) { }
}
export function Mixin(Super) {
    return class B extends Super {
        get myName() {
            return "B";
        }
        set myName(arg) { }
    };
}


//// [anonymousClassAccessorsDeclarationEmit1.d.ts]
export abstract class Base {
    accessor a: number;
}
export function middle(Super?: typeof Base): abstract new () => {
    get a(): number;
    set a(arg: number);
};
class A {
    constructor(...args: any[]);
}
export function Mixin<T extends typeof A>(Super: T): {
    new (...args: any[]): {
        get myName(): string;
        set myName(arg: string);
    };
} & T;
export {};


//// [DtsFileErrors]


anonymousClassAccessorsDeclarationEmit1.d.ts(8,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== anonymousClassAccessorsDeclarationEmit1.d.ts (1 errors) ====
    export abstract class Base {
        accessor a: number;
    }
    export function middle(Super?: typeof Base): abstract new () => {
        get a(): number;
        set a(arg: number);
    };
    class A {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        constructor(...args: any[]);
    }
    export function Mixin<T extends typeof A>(Super: T): {
        new (...args: any[]): {
            get myName(): string;
            set myName(arg: string);
        };
    } & T;
    export {};
    