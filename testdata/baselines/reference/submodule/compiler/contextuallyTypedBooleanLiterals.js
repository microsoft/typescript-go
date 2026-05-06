//// [tests/cases/compiler/contextuallyTypedBooleanLiterals.ts] ////

//// [contextuallyTypedBooleanLiterals.ts]
// Repro from #48363

type Box<T> = {
    get: () => T,
    set: (value: T) => void
}

declare function box<T>(value: T): Box<T>;

const bn1 = box(0);  // Box<number>
const bn2: Box<number> = box(0);  // Ok

const bb1 = box(false);  // Box<boolean>
const bb2: Box<boolean> = box(false);  // Error, box<false> not assignable to Box<boolean>

// Repro from #48150

interface Observable<T>
{
  (): T;
  (value: T): any;
}

declare function observable<T>(value: T): Observable<T>;

const x: Observable<boolean> = observable(false);


//// [contextuallyTypedBooleanLiterals.js]
"use strict";
// Repro from #48363
const bn1 = box(0); // Box<number>
const bn2 = box(0); // Ok
const bb1 = box(false); // Box<boolean>
const bb2 = box(false); // Error, box<false> not assignable to Box<boolean>
const x = observable(false);


//// [contextuallyTypedBooleanLiterals.d.ts]
type Box<T> = {
    get: () => T;
    set: (value: T) => void;
};
function box<T>(value: T): Box<T>;
const bn1: Box<number>;
const bn2: Box<number>;
const bb1: Box<boolean>;
const bb2: Box<boolean>;
interface Observable<T> {
    (): T;
    (value: T): any;
}
function observable<T>(value: T): Observable<T>;
const x: Observable<boolean>;


//// [DtsFileErrors]


contextuallyTypedBooleanLiterals.d.ts(5,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== contextuallyTypedBooleanLiterals.d.ts (1 errors) ====
    type Box<T> = {
        get: () => T;
        set: (value: T) => void;
    };
    function box<T>(value: T): Box<T>;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    const bn1: Box<number>;
    const bn2: Box<number>;
    const bb1: Box<boolean>;
    const bb2: Box<boolean>;
    interface Observable<T> {
        (): T;
        (value: T): any;
    }
    function observable<T>(value: T): Observable<T>;
    const x: Observable<boolean>;
    