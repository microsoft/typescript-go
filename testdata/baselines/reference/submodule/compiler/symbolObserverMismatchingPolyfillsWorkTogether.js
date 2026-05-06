//// [tests/cases/compiler/symbolObserverMismatchingPolyfillsWorkTogether.ts] ////

//// [symbolObserverMismatchingPolyfillsWorkTogether.ts]
interface SymbolConstructor {
    readonly observer: symbol;
}
interface SymbolConstructor {
    readonly observer: unique symbol;
}

const obj = {
    [Symbol.observer]: 0
};

//// [symbolObserverMismatchingPolyfillsWorkTogether.js]
"use strict";
const obj = {
    [Symbol.observer]: 0
};


//// [symbolObserverMismatchingPolyfillsWorkTogether.d.ts]
interface SymbolConstructor {
    readonly observer: symbol;
}
interface SymbolConstructor {
    readonly observer: unique symbol;
}
const obj: {
    [Symbol.observer]: number;
};


//// [DtsFileErrors]


symbolObserverMismatchingPolyfillsWorkTogether.d.ts(7,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== symbolObserverMismatchingPolyfillsWorkTogether.d.ts (1 errors) ====
    interface SymbolConstructor {
        readonly observer: symbol;
    }
    interface SymbolConstructor {
        readonly observer: unique symbol;
    }
    const obj: {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        [Symbol.observer]: number;
    };
    