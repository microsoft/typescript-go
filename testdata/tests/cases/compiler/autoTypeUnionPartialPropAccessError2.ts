// @strict: true
// @allowJs: true
// @checkJs: true
// @noEmit: true

// @filename: index.js

// https://github.com/microsoft/typescript-go/issues/3639

export class C {
    constructor() {
        const error = { message: 'props.message' };
        const shouldEmitError = Math.random() < 0.5;

        this.state = {
            ctx: shouldEmitError ? { ...error, name: 'A' } : { name: 'B' },
        }
    }

    method() {
        const { ctx: { message, name } } = this.state;
    }
}
