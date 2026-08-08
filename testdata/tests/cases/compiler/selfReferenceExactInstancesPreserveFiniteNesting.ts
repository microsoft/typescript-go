// @noEmit: true
// @noTypesAndSymbols: true
// @strict: true

interface Box<T> {
    left2(): Box<Left.L2>;
    left3(): Box<Left.L3>;
    leftLeaf(): Box<Left.Leaf>;
    right2(): Box<Right.L2>;
    right3(): Box<Right.L3>;
    rightLeaf(): Box<Right.Leaf>;
    value: T;
}

namespace Left {
    export interface Outer {
        box: Box<L2>;
    }
    export interface L2 {
        box: Box<L3>;
    }
    export interface L3 {
        box: Box<Leaf>;
    }
    export interface Leaf {
        id: string;
    }
}

namespace Right {
    export interface Outer {
        box: Box<L2>;
    }
    export interface L2 {
        box: Box<L3>;
    }
    export interface L3 {
        box: Box<Leaf>;
    }
    export interface Leaf {
        id: number;
    }
}

declare const box: Box<unknown>;
box.left2();
box.left3();
box.leftLeaf();
box.right2();
box.right3();
box.rightLeaf();

function test(left: Left.Outer, right: Right.Outer) {
    left = right;
}
