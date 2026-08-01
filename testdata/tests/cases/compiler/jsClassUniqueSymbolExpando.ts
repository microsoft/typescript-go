// @allowJs: true
// @checkJs: true
// @noEmit: true
// @lib: es2022

// @filename: index.js
const key = Symbol("key");

class Container {
    constructor() {
        this[key] = 1;
    }

    read() {
        return this[key];
    }
}

const value = new Container().read();
