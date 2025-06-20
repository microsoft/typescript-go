// @target: esnext
// @module: nodenext
// @checkJs: true
// @declaration: true

// @filename: a.ts
export class A<T> {}

// @filename: b.js
import { A } from './a.js';

export class B1 extends A {
    constructor() {
        super();
    }
}

export class B2 extends A<unknown, unknown> {
    constructor() {
        super();
    }
}