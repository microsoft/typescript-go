// @module: node20
// @filename: a.ts
export class A<T> {}

// @filename: b.js
// export class B2 extends A {
//     constructor() {
//         super();
//     }
// }

export class B2 extends A<unknown, unknown> {
    constructor() {
        super();
    }
}