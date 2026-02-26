// @declaration: true
// @filename: lib.ts
interface Base { }

export type ExportBase = Base;

// @filename: main.ts
import type { ExportBase } from './lib.js'

interface Clonable<T> {
    clone(): Clonable<T>
}

export default class C {
    async method(res: Clonable<{ key: ExportBase }>) {
        if (Math.random() > 0.5) {
            res.clone();
        } else {
            return res.clone();
        }
    };
}
