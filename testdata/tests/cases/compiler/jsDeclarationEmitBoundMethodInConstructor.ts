// @allowJs: true
// @checkJs: true
// @declaration: true
// @emitDeclarationOnly: true
// @outDir: out
// @stableTypeOrdering: true

// @filename: main.js
export class C {
    constructor() {
        this.foo = this.foo.bind(this);
    }
    foo() {}
}
