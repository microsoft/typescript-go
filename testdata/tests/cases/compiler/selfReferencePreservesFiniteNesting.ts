// @noEmit: true
// @noTypesAndSymbols: true
// @strict: true

interface RecursiveBox<T> {
    aNext?: RecursiveBox<T>;
    zValue: T;
}

namespace Left {
    export interface Outer {
        box: RecursiveBox<Inner>;
    }
    export interface Inner {
        box: RecursiveBox<Mid>;
    }
    export interface Mid {
        box: RecursiveBox<Leaf>;
    }
    export interface Leaf {
        id: string;
    }
}

namespace Right {
    export interface Outer {
        box: RecursiveBox<Inner>;
    }
    export interface Inner {
        box: RecursiveBox<Mid>;
    }
    export interface Mid {
        box: RecursiveBox<Leaf>;
    }
    export interface Leaf {
        id: number;
    }
}

function test(left: Left.Outer, right: Right.Outer) {
    left = right;
}
