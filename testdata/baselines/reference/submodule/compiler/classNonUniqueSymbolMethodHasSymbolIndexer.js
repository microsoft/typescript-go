//// [tests/cases/compiler/classNonUniqueSymbolMethodHasSymbolIndexer.ts] ////

//// [classNonUniqueSymbolMethodHasSymbolIndexer.ts]
declare const a: symbol;
export class A {
    [a]() { return 1 };
}
declare const e1: A[typeof a]; // no error, `A` has `symbol` index

type Constructor = new (...args: any[]) => {};
declare function Mix<T extends Constructor>(classish: T): T & (new (...args: any[]) => {mixed: true});

export const Mixer = Mix(class {
    [a]() { return 1 };
});


//// [classNonUniqueSymbolMethodHasSymbolIndexer.js]
export class A {
    [a]() { return 1; }
    ;
}
export const Mixer = Mix(class {
    [a]() { return 1; }
    ;
});


//// [classNonUniqueSymbolMethodHasSymbolIndexer.d.ts]
const a: symbol;
export class A {
    [a]: () => number;
}
export const Mixer: {
    new (): {
        [x: symbol]: () => number;
    };
} & (new (...args: any[]) => {
    mixed: true;
});
export {};


//// [DtsFileErrors]


classNonUniqueSymbolMethodHasSymbolIndexer.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== classNonUniqueSymbolMethodHasSymbolIndexer.d.ts (1 errors) ====
    const a: symbol;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export class A {
        [a]: () => number;
    }
    export const Mixer: {
        new (): {
            [x: symbol]: () => number;
        };
    } & (new (...args: any[]) => {
        mixed: true;
    });
    export {};
    