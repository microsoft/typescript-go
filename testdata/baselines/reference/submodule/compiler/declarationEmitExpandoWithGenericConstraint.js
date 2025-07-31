//// [tests/cases/compiler/declarationEmitExpandoWithGenericConstraint.ts] ////

//// [declarationEmitExpandoWithGenericConstraint.ts]
export interface Point {
    readonly x: number;
    readonly y: number;
}

export interface Rect<p extends Point> {
    readonly a: p;
    readonly b: p;
}

export const Point = (x: number, y: number): Point => ({ x, y });
export const Rect = <p extends Point>(a: p, b: p): Rect<p> => ({ a, b });

Point.zero = (): Point => Point(0, 0);

//// [declarationEmitExpandoWithGenericConstraint.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Rect = exports.Point = void 0;
const Point = (x, y) => ({ x, y });
exports.Point = Point;
const Rect = (a, b) => ({ a, b });
exports.Rect = Rect;
exports.Point.zero = () => (0, exports.Point)(0, 0);


//// [declarationEmitExpandoWithGenericConstraint.d.ts]
export interface Point {
    readonly x: number;
    readonly y: number;
}
export interface Rect<p extends Point> {
    readonly a: p;
    readonly b: p;
}
export declare const Point: {
    (x: number, y: number): Point;
    zero: () => Point;
};
export declare const Rect: <p extends Point>(a: p, b: p) => Rect<p>;
declare namespace Point {
    const zero: () => Point;
}


//// [DtsFileErrors]


declarationEmitExpandoWithGenericConstraint.d.ts(1,18): error TS2451: Cannot redeclare block-scoped variable 'Point'.
declarationEmitExpandoWithGenericConstraint.d.ts(9,22): error TS2395: Individual declarations in merged declaration 'Point' must be all exported or all local.
declarationEmitExpandoWithGenericConstraint.d.ts(9,22): error TS2451: Cannot redeclare block-scoped variable 'Point'.
declarationEmitExpandoWithGenericConstraint.d.ts(14,19): error TS2395: Individual declarations in merged declaration 'Point' must be all exported or all local.
declarationEmitExpandoWithGenericConstraint.d.ts(14,19): error TS2451: Cannot redeclare block-scoped variable 'Point'.


==== declarationEmitExpandoWithGenericConstraint.d.ts (5 errors) ====
    export interface Point {
                     ~~~~~
!!! error TS2451: Cannot redeclare block-scoped variable 'Point'.
        readonly x: number;
        readonly y: number;
    }
    export interface Rect<p extends Point> {
        readonly a: p;
        readonly b: p;
    }
    export declare const Point: {
                         ~~~~~
!!! error TS2395: Individual declarations in merged declaration 'Point' must be all exported or all local.
                         ~~~~~
!!! error TS2451: Cannot redeclare block-scoped variable 'Point'.
        (x: number, y: number): Point;
        zero: () => Point;
    };
    export declare const Rect: <p extends Point>(a: p, b: p) => Rect<p>;
    declare namespace Point {
                      ~~~~~
!!! error TS2395: Individual declarations in merged declaration 'Point' must be all exported or all local.
                      ~~~~~
!!! error TS2451: Cannot redeclare block-scoped variable 'Point'.
        const zero: () => Point;
    }
    