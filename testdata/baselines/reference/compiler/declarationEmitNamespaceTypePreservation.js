//// [tests/cases/compiler/declarationEmitNamespaceTypePreservation.ts] ////

//// [declarationEmitNamespaceTypePreservation.ts]
export namespace Pull {
  export interface Halt<E> {
    readonly _tag: "Halt"
    readonly error: E
  }

  export type ExcludeHalt<E> = Exclude<E, Halt<unknown>>
}

export namespace Arr {
  export type NonEmptyReadonlyArray<A> = readonly [A, ...ReadonlyArray<A>]
}

// Namespace-qualified types should be preserved
export const excludeHalt = <E>(error: E): Pull.ExcludeHalt<E> => error as Pull.ExcludeHalt<E>
export const toNonEmpty = <A>(arr: Arr.NonEmptyReadonlyArray<A>): Arr.NonEmptyReadonlyArray<A> => arr
export const process = <E>(input: E): Pull.ExcludeHalt<E> | null => input as Pull.ExcludeHalt<E>


//// [declarationEmitNamespaceTypePreservation.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.process = exports.toNonEmpty = exports.excludeHalt = void 0;
// Namespace-qualified types should be preserved
const excludeHalt = (error) => error;
exports.excludeHalt = excludeHalt;
const toNonEmpty = (arr) => arr;
exports.toNonEmpty = toNonEmpty;
const process = (input) => input;
exports.process = process;


//// [declarationEmitNamespaceTypePreservation.d.ts]
export declare namespace Pull {
    interface Halt<E> {
        readonly _tag: "Halt";
        readonly error: E;
    }
    type ExcludeHalt<E> = Exclude<E, Halt<unknown>>;
}
export declare namespace Arr {
    type NonEmptyReadonlyArray<A> = readonly [A, ...ReadonlyArray<A>];
}
export declare const excludeHalt: <E>(error: E) => Pull.ExcludeHalt<E>;
export declare const toNonEmpty: <A>(arr: Arr.NonEmptyReadonlyArray<A>) => Arr.NonEmptyReadonlyArray<A>;
export declare const process: <E>(input: E) => Pull.ExcludeHalt<E> | null;
