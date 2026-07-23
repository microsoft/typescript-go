// @allowJs: true
// @checkJs: true
// @declaration: true
// @emitDeclarationOnly: true
// @outDir: ./out

// @filename: lib.js
/**
 * @typedef {Object} Foo
 * @property {boolean} bool Whether `.bool` is true or not
 */
export class C {
    /** @returns {Foo} */
    getFoo() { return { bool: false }; }
}

// @filename: main.js
import { C } from './lib.js';

export class Main {
    constructor() {
        this.c = new C();
    }

    getFoo() { return { ...this.c.getFoo() }; }
}
