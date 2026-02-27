//// [tests/cases/conformance/classes/propertyMemberDeclarations/staticAutoAccessorsWithDecorators.ts] ////

//// [staticAutoAccessorsWithDecorators.ts]
// https://github.com/microsoft/TypeScript/issues/53752

class A {
    // uses class reference
    @((t, c) => {})
    static accessor x = 1;

    // uses 'this'
    @((t, c) => {})
    accessor y = 2;
}


//// [staticAutoAccessorsWithDecorators.js]
"use strict";
// https://github.com/microsoft/TypeScript/issues/53752
let A = (() => {
    let _static_x_decorators;
    let _static_x_initializers = [];
    let _static_x_extraInitializers = [];
    let _y_decorators;
    let _y_initializers = [];
    let _y_extraInitializers = [];
    return class A {
        static {
            const _metadata = typeof Symbol === "function" && Symbol.metadata ? Object.create(null) : void 0;
            _static_x_decorators = [((t, c) => { })];
            _y_decorators = [((t, c) => { })];
            __esDecorate(this, null, _static_x_decorators, { kind: "accessor", name: "x", static: true, private: false, access: { has: obj => "x" in obj, get: obj => obj.x, set: (obj, value) => { obj.x = value; } }, metadata: _metadata }, _static_x_initializers, _static_x_extraInitializers);
            __esDecorate(this, null, _y_decorators, { kind: "accessor", name: "y", static: false, private: false, access: { has: obj => "y" in obj, get: obj => obj.y, set: (obj, value) => { obj.y = value; } }, metadata: _metadata }, _y_initializers, _y_extraInitializers);
            if (_metadata) Object.defineProperty(this, Symbol.metadata, { enumerable: true, configurable: true, writable: true, value: _metadata });
        }
        // uses class reference
        static accessor x = __runInitializers(this, _static_x_initializers, 1);
        // uses 'this'
        accessor y = __runInitializers(this, _y_initializers, 2);
        constructor() {
            __runInitializers(this, _y_extraInitializers);
        }
        static {
            __runInitializers(this, _static_x_extraInitializers);
        }
    };
})();
