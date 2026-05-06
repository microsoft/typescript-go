//// [tests/cases/compiler/declInput.ts] ////

//// [declInput.ts]
interface bar {

}

class bar {
  public f() { return ''; }
  public g() { return {a: <bar>null, b: undefined, c: void 4 }; }
  public h(x = 4, y = null, z = '') { x++; }
}


//// [declInput.js]
"use strict";
class bar {
    f() { return ''; }
    g() { return { a: null, b: undefined, c: void 4 }; }
    h(x = 4, y = null, z = '') { x++; }
}


//// [declInput.d.ts]
interface bar {
}
class bar {
    f(): string;
    g(): {
        a: bar;
        b: any;
        c: any;
    };
    h(x?: number, y?: any, z?: string): void;
}


//// [DtsFileErrors]


declInput.d.ts(3,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declInput.d.ts (1 errors) ====
    interface bar {
    }
    class bar {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        f(): string;
        g(): {
            a: bar;
            b: any;
            c: any;
        };
        h(x?: number, y?: any, z?: string): void;
    }
    