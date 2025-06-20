// @target: esnext
// @module: nodenext
// @checkJs: true
// @declaration: true

// @filename: a.ts
export class A<T> {}

// @filename: b.js
import { A } from './a.js';

/** @extends {A} */
export class B1 extends A {
    constructor() {
        super();
    }
}

/** @extends {A<unknown, unknown>} */
export class B2 extends A {
    constructor() {
        super();
    }
}