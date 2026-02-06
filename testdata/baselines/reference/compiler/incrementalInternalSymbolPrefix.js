//// [tests/cases/compiler/incrementalInternalSymbolPrefix.ts] ////

//// [incrementalInternalSymbolPrefix.ts]
// This test triggers the "X is specified more than once, so this usage will be overwritten" diagnostic
// which contains internal symbol names like __@iterator@. When marshalling build info, the internal
// symbol name prefix (\xFE) caused invalid UTF-8 errors.
// See: https://github.com/microsoft/typescript-go/issues/1531

const items = [1, 2, 3];

// The spread ...items overwrites the [Symbol.iterator] property specified above it
const obj = {
    length: items.length,
    [Symbol.iterator]: function* () {
        for (const item of items) yield item;
    },
    ...items,  // This spread overwrites 'length' and '[Symbol.iterator]'
};


//// [incrementalInternalSymbolPrefix.js]
"use strict";
// This test triggers the "X is specified more than once, so this usage will be overwritten" diagnostic
// which contains internal symbol names like __@iterator@. When marshalling build info, the internal
// symbol name prefix (\xFE) caused invalid UTF-8 errors.
// See: https://github.com/microsoft/typescript-go/issues/1531
const items = [1, 2, 3];
// The spread ...items overwrites the [Symbol.iterator] property specified above it
const obj = {
    length: items.length,
    [Symbol.iterator]: function* () {
        for (const item of items)
            yield item;
    },
    ...items, // This spread overwrites 'length' and '[Symbol.iterator]'
};


//// [incrementalInternalSymbolPrefix.d.ts]
declare const items: number[];
declare const obj: {
    length: number;
    toString(): string;
    toLocaleString(): string;
    toLocaleString(locales: string | string[], options?: (Intl.NumberFormatOptions & Intl.DateTimeFormatOptions) | undefined): string;
    pop(): number | undefined;
    push(...items: number[]): number;
    concat(...items: ConcatArray<number>[]): number[];
    concat(...items: (number | ConcatArray<number>)[]): number[];
    join(separator?: string | undefined): string;
    reverse(): number[];
    shift(): number | undefined;
    slice(start?: number | undefined, end?: number | undefined): number[];
    sort(compareFn?: ((a: number, b: number) => number) | undefined): number[];
    splice(start: number, deleteCount?: number | undefined): number[];
    splice(start: number, deleteCount: number, ...items: number[]): number[];
    unshift(...items: number[]): number;
    indexOf(searchElement: number, fromIndex?: number | undefined): number;
    lastIndexOf(searchElement: number, fromIndex?: number | undefined): number;
    every<S extends number>(predicate: (value: number, index: number, array: number[]) => value is S, thisArg?: any): this is S[];
    every(predicate: (value: number, index: number, array: number[]) => unknown, thisArg?: any): boolean;
    some(predicate: (value: number, index: number, array: number[]) => unknown, thisArg?: any): boolean;
    forEach(callbackfn: (value: number, index: number, array: number[]) => void, thisArg?: any): void;
    map<U>(callbackfn: (value: number, index: number, array: number[]) => U, thisArg?: any): U[];
    filter<S extends number>(predicate: (value: number, index: number, array: number[]) => value is S, thisArg?: any): S[];
    filter(predicate: (value: number, index: number, array: number[]) => unknown, thisArg?: any): number[];
    reduce(callbackfn: (previousValue: number, currentValue: number, currentIndex: number, array: number[]) => number): number;
    reduce(callbackfn: (previousValue: number, currentValue: number, currentIndex: number, array: number[]) => number, initialValue: number): number;
    reduce<U>(callbackfn: (previousValue: U, currentValue: number, currentIndex: number, array: number[]) => U, initialValue: U): U;
    reduceRight(callbackfn: (previousValue: number, currentValue: number, currentIndex: number, array: number[]) => number): number;
    reduceRight(callbackfn: (previousValue: number, currentValue: number, currentIndex: number, array: number[]) => number, initialValue: number): number;
    reduceRight<U>(callbackfn: (previousValue: U, currentValue: number, currentIndex: number, array: number[]) => U, initialValue: U): U;
    find<S extends number>(predicate: (value: number, index: number, obj: number[]) => value is S, thisArg?: any): S | undefined;
    find(predicate: (value: number, index: number, obj: number[]) => unknown, thisArg?: any): number | undefined;
    findIndex(predicate: (value: number, index: number, obj: number[]) => unknown, thisArg?: any): number;
    fill(value: number, start?: number | undefined, end?: number | undefined): number[];
    copyWithin(target: number, start: number, end?: number | undefined): number[];
    [Symbol.iterator](): ArrayIterator<number>;
    entries(): ArrayIterator<[number, number]>;
    keys(): ArrayIterator<number>;
    values(): ArrayIterator<number>;
    [Symbol.unscopables]: {
        [x: number]: boolean | undefined;
        length?: boolean | undefined;
        toString?: boolean | undefined;
        toLocaleString?: boolean | undefined;
        pop?: boolean | undefined;
        push?: boolean | undefined;
        concat?: boolean | undefined;
        join?: boolean | undefined;
        reverse?: boolean | undefined;
        shift?: boolean | undefined;
        slice?: boolean | undefined;
        sort?: boolean | undefined;
        splice?: boolean | undefined;
        unshift?: boolean | undefined;
        indexOf?: boolean | undefined;
        lastIndexOf?: boolean | undefined;
        every?: boolean | undefined;
        some?: boolean | undefined;
        forEach?: boolean | undefined;
        map?: boolean | undefined;
        filter?: boolean | undefined;
        reduce?: boolean | undefined;
        reduceRight?: boolean | undefined;
        find?: boolean | undefined;
        findIndex?: boolean | undefined;
        fill?: boolean | undefined;
        copyWithin?: boolean | undefined;
        [Symbol.iterator]?: boolean | undefined;
        entries?: boolean | undefined;
        keys?: boolean | undefined;
        values?: boolean | undefined;
        readonly [Symbol.unscopables]?: boolean | undefined;
        includes?: boolean | undefined;
        flatMap?: boolean | undefined;
        flat?: boolean | undefined;
        at?: boolean | undefined;
        findLast?: boolean | undefined;
        findLastIndex?: boolean | undefined;
        toReversed?: boolean | undefined;
        toSorted?: boolean | undefined;
        toSpliced?: boolean | undefined;
        with?: boolean | undefined;
    };
    includes(searchElement: number, fromIndex?: number | undefined): boolean;
    flatMap<U, This = undefined>(callback: (this: This, value: number, index: number, array: number[]) => U | readonly U[], thisArg?: This | undefined): U[];
    flat<A, D extends number = 1>(this: A, depth?: D | undefined): FlatArray<A, D>[];
    at(index: number): number | undefined;
    findLast<S extends number>(predicate: (value: number, index: number, array: number[]) => value is S, thisArg?: any): S | undefined;
    findLast(predicate: (value: number, index: number, array: number[]) => unknown, thisArg?: any): number | undefined;
    findLastIndex(predicate: (value: number, index: number, array: number[]) => unknown, thisArg?: any): number;
    toReversed(): number[];
    toSorted(compareFn?: ((a: number, b: number) => number) | undefined): number[];
    toSpliced(start: number, deleteCount: number, ...items: number[]): number[];
    toSpliced(start: number, deleteCount?: number | undefined): number[];
    with(index: number, value: number): number[];
};
