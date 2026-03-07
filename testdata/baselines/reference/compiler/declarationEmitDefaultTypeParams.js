//// [tests/cases/compiler/declarationEmitDefaultTypeParams.ts] ////

//// [declarationEmitDefaultTypeParams.ts]
export interface Effect<A, E = never, R = never> {
  readonly _A: A
  readonly _E: E
  readonly _R: R
}

export interface Deferred<A, E = never> {
  readonly _A: A
  readonly _E: E
}

// When return type uses defaults, they should be omitted
export const succeed = <A>(value: A): Effect<A> => ({ _A: value, _E: undefined as never, _R: undefined as never })
export const fail = <E>(error: E): Effect<never, E> => ({ _A: undefined as never, _E: error, _R: undefined as never })
export const makeDeferred = <A>(): Deferred<A> => ({ _A: undefined as never, _E: undefined as never })

// Many default params - only non-defaults should appear
export interface Channel<Out, Err = never, Done = void, In = unknown, InErr = unknown, InDone = unknown, Env = never> {
  readonly _Out: Out
}

export const fromString = (s: string): Channel<string> => ({ _Out: s })


//// [declarationEmitDefaultTypeParams.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.fromString = exports.makeDeferred = exports.fail = exports.succeed = void 0;
// When return type uses defaults, they should be omitted
const succeed = (value) => ({ _A: value, _E: undefined, _R: undefined });
exports.succeed = succeed;
const fail = (error) => ({ _A: undefined, _E: error, _R: undefined });
exports.fail = fail;
const makeDeferred = () => ({ _A: undefined, _E: undefined })
// Many default params - only non-defaults should appear
;
exports.makeDeferred = makeDeferred;
const fromString = (s) => ({ _Out: s });
exports.fromString = fromString;


//// [declarationEmitDefaultTypeParams.d.ts]
export interface Effect<A, E = never, R = never> {
    readonly _A: A;
    readonly _E: E;
    readonly _R: R;
}
export interface Deferred<A, E = never> {
    readonly _A: A;
    readonly _E: E;
}
export declare const succeed: <A>(value: A) => Effect<A>;
export declare const fail: <E>(error: E) => Effect<never, E>;
export declare const makeDeferred: <A>() => Deferred<A>;
export interface Channel<Out, Err = never, Done = void, In = unknown, InErr = unknown, InDone = unknown, Env = never> {
    readonly _Out: Out;
}
export declare const fromString: (s: string) => Channel<string>;
