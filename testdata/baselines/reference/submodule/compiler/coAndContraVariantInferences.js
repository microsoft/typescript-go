//// [tests/cases/compiler/coAndContraVariantInferences.ts] ////

//// [coAndContraVariantInferences.ts]
type A = { kind: 'a' };
type B = { kind: 'b' };

declare const a: A;
declare const b: B;

declare function fab(arg: A | B): void;

declare function foo<T>(x: { kind: T }, f: (arg: { kind: T }) => void): void;

foo(a, fab);
foo(b, fab);

// Repro from #45603

interface Action<TName extends string,TPayload> {
    name: TName,
    payload: TPayload
}

const actionA = { payload: 'any-string' } as Action<'ACTION_A', string>;
const actionB = { payload: true } as Action<'ACTION_B', boolean>;

function call<TName extends string,TPayload>(
  action: Action<TName,TPayload>,
  fn: (action: Action<TName,TPayload>)=> any,
) {
  fn(action);
}

const printFn = (action: typeof actionA | typeof actionB)=> console.log(action);

call(actionA, printFn);
call(actionB, printFn);


//// [coAndContraVariantInferences.js]
"use strict";
foo(a, fab);
foo(b, fab);
const actionA = { payload: 'any-string' };
const actionB = { payload: true };
function call(action, fn) {
    fn(action);
}
const printFn = (action) => console.log(action);
call(actionA, printFn);
call(actionB, printFn);


//// [coAndContraVariantInferences.d.ts]
type A = {
    kind: 'a';
};
type B = {
    kind: 'b';
};
const a: A;
const b: B;
function fab(arg: A | B): void;
function foo<T>(x: {
    kind: T;
}, f: (arg: {
    kind: T;
}) => void): void;
interface Action<TName extends string, TPayload> {
    name: TName;
    payload: TPayload;
}
const actionA: Action<'ACTION_A', string>;
const actionB: Action<'ACTION_B', boolean>;
function call<TName extends string, TPayload>(action: Action<TName, TPayload>, fn: (action: Action<TName, TPayload>) => any): void;
const printFn: (action: typeof actionA | typeof actionB) => void;


//// [DtsFileErrors]


coAndContraVariantInferences.d.ts(7,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== coAndContraVariantInferences.d.ts (1 errors) ====
    type A = {
        kind: 'a';
    };
    type B = {
        kind: 'b';
    };
    const a: A;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    const b: B;
    function fab(arg: A | B): void;
    function foo<T>(x: {
        kind: T;
    }, f: (arg: {
        kind: T;
    }) => void): void;
    interface Action<TName extends string, TPayload> {
        name: TName;
        payload: TPayload;
    }
    const actionA: Action<'ACTION_A', string>;
    const actionB: Action<'ACTION_B', boolean>;
    function call<TName extends string, TPayload>(action: Action<TName, TPayload>, fn: (action: Action<TName, TPayload>) => any): void;
    const printFn: (action: typeof actionA | typeof actionB) => void;
    