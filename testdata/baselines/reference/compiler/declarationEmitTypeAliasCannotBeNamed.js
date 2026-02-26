//// [tests/cases/compiler/declarationEmitTypeAliasCannotBeNamed.ts] ////

//// [lib.ts]
interface Base { }

export type ExportBase = Base;

//// [main.ts]
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


//// [lib.js]
export {};
//// [main.js]
export default class C {
    async method(res) {
        if (Math.random() > 0.5) {
            res.clone();
        }
        else {
            return res.clone();
        }
    }
    ;
}


//// [lib.d.ts]
interface Base {
}
export type ExportBase = Base;
export {};
//// [main.d.ts]
import type { ExportBase } from './lib.js';
interface Clonable<T> {
    clone(): Clonable<T>;
}
export default class C {
    method(res: Clonable<{
        key: ExportBase;
    }>): Promise<Clonable<{
        key: ExportBase;
    }> | undefined>;
}
export {};
