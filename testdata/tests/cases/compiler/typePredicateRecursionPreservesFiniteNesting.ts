// @noEmit: true
// @noTypesAndSymbols: true
// @strict: true

interface Guard<T> {
    guard(): this is this & T;
}

namespace Left {
    export interface L2 { next: Guard<L3>; }
    export interface L3 { next: Guard<L4>; }
    export interface L4 { next: Guard<Leaf>; }
    export interface Leaf { id: string; }
}

namespace Right {
    export interface L2 { next: Guard<L3>; }
    export interface L3 { next: Guard<L4>; }
    export interface L4 { next: Guard<Leaf>; }
    export interface Leaf { id: number; }
}

declare const left: Guard<Left.L2>;
declare const right: Guard<Right.L2>;

const predicateComparison: typeof left.guard = right.guard;
const ordinaryComparison: Guard<Left.L2> = right;
