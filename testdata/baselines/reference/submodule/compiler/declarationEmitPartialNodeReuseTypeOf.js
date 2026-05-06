//// [tests/cases/compiler/declarationEmitPartialNodeReuseTypeOf.ts] ////

//// [a.ts]
export const nImported = "nImported"
export const nNotImported = "nNotImported"
const nPrivate = "private"
export const o = (p1: typeof nImported, p2: typeof nNotImported, p3: typeof nPrivate) => null! as { foo: typeof nImported, bar: typeof nPrivate, baz: typeof nNotImported }

//// [b.ts]
import { o, nImported } from "./a";
export const g = o
console.log(nImported);

//// [c.ts]
import * as a from "./a";
export const g = a.o


//// [a.js]
export const nImported = "nImported";
export const nNotImported = "nNotImported";
const nPrivate = "private";
export const o = (p1, p2, p3) => null;
//// [b.js]
import { o, nImported } from "./a";
export const g = o;
console.log(nImported);
//// [c.js]
import * as a from "./a";
export const g = a.o;


//// [a.d.ts]
export const nImported = "nImported";
export const nNotImported = "nNotImported";
const nPrivate = "private";
export const o: (p1: typeof nImported, p2: typeof nNotImported, p3: typeof nPrivate) => {
    foo: typeof nImported;
    bar: typeof nPrivate;
    baz: typeof nNotImported;
};
export {};
//// [b.d.ts]
import { nImported } from "./a";
export const g: (p1: typeof nImported, p2: typeof import("./a").nNotImported, p3: "private") => {
    foo: typeof nImported;
    bar: "private";
    baz: typeof import("./a").nNotImported;
};
//// [c.d.ts]
import * as a from "./a";
export const g: (p1: typeof a.nImported, p2: typeof a.nNotImported, p3: "private") => {
    foo: typeof a.nImported;
    bar: "private";
    baz: typeof a.nNotImported;
};


//// [DtsFileErrors]


a.d.ts(3,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== a.d.ts (1 errors) ====
    export const nImported = "nImported";
    export const nNotImported = "nNotImported";
    const nPrivate = "private";
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export const o: (p1: typeof nImported, p2: typeof nNotImported, p3: typeof nPrivate) => {
        foo: typeof nImported;
        bar: typeof nPrivate;
        baz: typeof nNotImported;
    };
    export {};
    
==== b.d.ts (0 errors) ====
    import { nImported } from "./a";
    export const g: (p1: typeof nImported, p2: typeof import("./a").nNotImported, p3: "private") => {
        foo: typeof nImported;
        bar: "private";
        baz: typeof import("./a").nNotImported;
    };
    
==== c.d.ts (0 errors) ====
    import * as a from "./a";
    export const g: (p1: typeof a.nImported, p2: typeof a.nNotImported, p3: "private") => {
        foo: typeof a.nImported;
        bar: "private";
        baz: typeof a.nNotImported;
    };
    