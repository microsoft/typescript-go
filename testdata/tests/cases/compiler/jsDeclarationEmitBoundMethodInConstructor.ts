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

export class D {
    constructor() {
        this.bar = 1;
    }
    static bar() {}
}

export class E {
    static init() {
        this.baz = 1;
    }
    baz() {}
}
