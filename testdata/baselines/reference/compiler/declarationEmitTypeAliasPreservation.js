//// [tests/cases/compiler/declarationEmitTypeAliasPreservation.ts] ////

//// [declarationEmitTypeAliasPreservation.ts]
export type NonEmptyArray<A> = [A, ...Array<A>]
export type NonEmptyReadonlyArray<A> = readonly [A, ...ReadonlyArray<A>]

// All of these should preserve the type alias in declaration emit
export const make = <A>(...elements: NonEmptyArray<A>): NonEmptyArray<A> => elements
export const fromArray = <A>(arr: Array<A>): Array<A> => arr
export const fromReadonly = <A>(arr: ReadonlyArray<A>): ReadonlyArray<A> => arr
export const first = <A>(arr: NonEmptyReadonlyArray<A>): A => arr[0]
export const allocate = <A>(n: number): Array<A | undefined> => new Array(n)


//// [declarationEmitTypeAliasPreservation.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.allocate = exports.first = exports.fromReadonly = exports.fromArray = exports.make = void 0;
// All of these should preserve the type alias in declaration emit
const make = (...elements) => elements;
exports.make = make;
const fromArray = (arr) => arr;
exports.fromArray = fromArray;
const fromReadonly = (arr) => arr;
exports.fromReadonly = fromReadonly;
const first = (arr) => arr[0];
exports.first = first;
const allocate = (n) => new Array(n);
exports.allocate = allocate;


//// [declarationEmitTypeAliasPreservation.d.ts]
export type NonEmptyArray<A> = [A, ...Array<A>];
export type NonEmptyReadonlyArray<A> = readonly [A, ...ReadonlyArray<A>];
export declare const make: <A>(...elements: NonEmptyArray<A>) => NonEmptyArray<A>;
export declare const fromArray: <A>(arr: Array<A>) => Array<A>;
export declare const fromReadonly: <A>(arr: ReadonlyArray<A>) => ReadonlyArray<A>;
export declare const first: <A>(arr: NonEmptyReadonlyArray<A>) => A;
export declare const allocate: <A>(n: number) => Array<A | undefined>;
