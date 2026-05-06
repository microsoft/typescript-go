//// [tests/cases/compiler/parameterDestructuringObjectLiteral.ts] ////

//// [parameterDestructuringObjectLiteral.ts]
// Repro from #22644

const fn1 = (options: { headers?: {} }) => { };
fn1({ headers: { foo: 1 } });

const fn2 = ({ headers = {} }) => { };
fn2({ headers: { foo: 1 } });


//// [parameterDestructuringObjectLiteral.js]
"use strict";
// Repro from #22644
const fn1 = (options) => { };
fn1({ headers: { foo: 1 } });
const fn2 = ({ headers = {} }) => { };
fn2({ headers: { foo: 1 } });


//// [parameterDestructuringObjectLiteral.d.ts]
const fn1: (options: {
    headers?: {};
}) => void;
const fn2: ({ headers }: {
    headers?: {} | undefined;
}) => void;


//// [DtsFileErrors]


parameterDestructuringObjectLiteral.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== parameterDestructuringObjectLiteral.d.ts (1 errors) ====
    const fn1: (options: {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        headers?: {};
    }) => void;
    const fn2: ({ headers }: {
        headers?: {} | undefined;
    }) => void;
    