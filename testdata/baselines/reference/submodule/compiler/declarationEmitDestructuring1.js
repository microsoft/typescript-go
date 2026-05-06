//// [tests/cases/compiler/declarationEmitDestructuring1.ts] ////

//// [declarationEmitDestructuring1.ts]
function foo([a, b, c]: [string, string, string]): void { }
function far([a, [b], [[c]]]: [number, boolean[], string[][]]): void { }
function bar({a1, b1, c1}: { a1: number, b1: boolean, c1: string }): void { }
function baz({a2, b2: {b1, c1}}: { a2: number, b2: { b1: boolean, c1: string } }): void { } 


//// [declarationEmitDestructuring1.js]
"use strict";
function foo([a, b, c]) { }
function far([a, [b], [[c]]]) { }
function bar({ a1, b1, c1 }) { }
function baz({ a2, b2: { b1, c1 } }) { }


//// [declarationEmitDestructuring1.d.ts]
function foo([a, b, c]: [string, string, string]): void;
function far([a, [b], [[c]]]: [number, boolean[], string[][]]): void;
function bar({ a1, b1, c1 }: {
    a1: number;
    b1: boolean;
    c1: string;
}): void;
function baz({ a2, b2: { b1, c1 } }: {
    a2: number;
    b2: {
        b1: boolean;
        c1: string;
    };
}): void;


//// [DtsFileErrors]


declarationEmitDestructuring1.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitDestructuring1.d.ts (1 errors) ====
    function foo([a, b, c]: [string, string, string]): void;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    function far([a, [b], [[c]]]: [number, boolean[], string[][]]): void;
    function bar({ a1, b1, c1 }: {
        a1: number;
        b1: boolean;
        c1: string;
    }): void;
    function baz({ a2, b2: { b1, c1 } }: {
        a2: number;
        b2: {
            b1: boolean;
            c1: string;
        };
    }): void;
    