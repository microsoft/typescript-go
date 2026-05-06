//// [tests/cases/compiler/declarationEmitResolveTypesIfNotReusable.ts] ////

//// [decl.ts]
const u = "X";
type A = { a: { b : "value of b", notNecessary: typeof u }}
const a = { a: "value of a", notNecessary: u } as const


export const o1 = (o: A['a']['b']) => {}

export const o2 = (o: (typeof a)['a']) => {}
export const o3 = (o:  typeof a['a']) => {}

export const o4 = (o: keyof (A['a'])) => {}

//// [main.ts]
import * as d  from './decl'

export const f = {...d}

//// [decl.js]
const u = "X";
const a = { a: "value of a", notNecessary: u };
export const o1 = (o) => { };
export const o2 = (o) => { };
export const o3 = (o) => { };
export const o4 = (o) => { };
//// [main.js]
import * as d from './decl';
export const f = { ...d };


//// [decl.d.ts]
const u = "X";
type A = {
    a: {
        b: "value of b";
        notNecessary: typeof u;
    };
};
const a: {
    readonly a: "value of a";
    readonly notNecessary: "X";
};
export const o1: (o: A['a']['b']) => void;
export const o2: (o: (typeof a)['a']) => void;
export const o3: (o: (typeof a)['a']) => void;
export const o4: (o: keyof (A['a'])) => void;
export {};
//// [main.d.ts]
export const f: {
    o1: (o: "value of b") => void;
    o2: (o: ({
        readonly a: "value of a";
        readonly notNecessary: "X";
    })['a']) => void;
    o3: (o: "value of a") => void;
    o4: (o: keyof ({
        b: "value of b";
        notNecessary: "X";
    })) => void;
};


//// [DtsFileErrors]


decl.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== decl.d.ts (1 errors) ====
    const u = "X";
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    type A = {
        a: {
            b: "value of b";
            notNecessary: typeof u;
        };
    };
    const a: {
        readonly a: "value of a";
        readonly notNecessary: "X";
    };
    export const o1: (o: A['a']['b']) => void;
    export const o2: (o: (typeof a)['a']) => void;
    export const o3: (o: (typeof a)['a']) => void;
    export const o4: (o: keyof (A['a'])) => void;
    export {};
    
==== main.d.ts (0 errors) ====
    export const f: {
        o1: (o: "value of b") => void;
        o2: (o: ({
            readonly a: "value of a";
            readonly notNecessary: "X";
        })['a']) => void;
        o3: (o: "value of a") => void;
        o4: (o: keyof ({
            b: "value of b";
            notNecessary: "X";
        })) => void;
    };
    