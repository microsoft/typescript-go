//// [tests/cases/conformance/esDecorators/classDeclaration/esDecorators-classDeclaration-sourceMap.ts] ////

//// [esDecorators-classDeclaration-sourceMap.ts]
declare var dec: any;

@dec
@dec
class C {
    @dec
    @dec
    method() {}

    @dec
    @dec
    get x() { return 1; }

    @dec
    @dec
    set x(value: number) { }

    @dec
    @dec
    y = 1;

    @dec
    @dec
    accessor z = 1;

    @dec
    @dec
    static #method() {}

    @dec
    @dec
    static get #x() { return 1; }

    @dec
    @dec
    static set #x(value: number) { }

    @dec
    @dec
    static #y = 1;

    @dec
    @dec
    static accessor #z = 1;
}


//// [esDecorators-classDeclaration-sourceMap.js]
"use strict";
let C = (() => {
    let _classDecorators = [dec, dec];
    let _classDescriptor;
    let _classExtraInitializers = [];
    let _classThis;
    let _staticExtraInitializers = [];
    let _instanceExtraInitializers = [];
    let _static_private_method_decorators;
    let _static_private_method_descriptor;
    let _static_private_get_x_decorators;
    let _static_private_get_x_descriptor;
    let _static_private_set_x_decorators;
    let _static_private_set_x_descriptor;
    let _static_private_y_decorators;
    let _static_private_y_initializers = [];
    let _static_private_y_extraInitializers = [];
    let _static_private_z_decorators;
    let _static_private_z_initializers = [];
    let _static_private_z_extraInitializers = [];
    let _method_decorators;
    let _get_x_decorators;
    let _set_x_decorators;
    let _y_decorators;
    let _y_initializers = [];
    let _y_extraInitializers = [];
    let _z_decorators;
    let _z_initializers = [];
    let _z_extraInitializers = [];
    var C = class {
        static { _classThis = this; }
        static {
            const _metadata = typeof Symbol === "function" && Symbol.metadata ? Object.create(null) : void 0;
            _method_decorators = [dec, dec];
            _get_x_decorators = [dec, dec];
            _set_x_decorators = [dec, dec];
            _y_decorators = [dec, dec];
            _z_decorators = [dec, dec];
            _static_private_method_decorators = [dec, dec];
            _static_private_get_x_decorators = [dec, dec];
            _static_private_set_x_decorators = [dec, dec];
            _static_private_y_decorators = [dec, dec];
            _static_private_z_decorators = [dec, dec];
            __esDecorate(this, _static_private_method_descriptor = { value: __setFunctionName(function () { }, "#method") }, _static_private_method_decorators, { kind: "method", name: "#method", static: true, private: true, access: { has: obj => #method in obj, get: obj => obj.#method }, metadata: _metadata }, null, _staticExtraInitializers);
            __esDecorate(this, _static_private_get_x_descriptor = { get: __setFunctionName(function () { return 1; }, "#x", "get") }, _static_private_get_x_decorators, { kind: "getter", name: "#x", static: true, private: true, access: { has: obj => #x in obj, get: obj => obj.#x }, metadata: _metadata }, null, _staticExtraInitializers);
            __esDecorate(this, _static_private_set_x_descriptor = { set: __setFunctionName(function (value) { }, "#x", "set") }, _static_private_set_x_decorators, { kind: "setter", name: "#x", static: true, private: true, access: { has: obj => #x in obj, set: (obj, value) => { obj.#x = value; } }, metadata: _metadata }, null, _staticExtraInitializers);
            __esDecorate(this, null, _static_private_z_decorators, { kind: "accessor", name: "#z", static: true, private: true, access: { has: obj => #z in obj, get: obj => obj.#z, set: (obj, value) => { obj.#z = value; } }, metadata: _metadata }, _static_private_z_initializers, _static_private_z_extraInitializers);
            __esDecorate(this, null, _method_decorators, { kind: "method", name: "method", static: false, private: false, access: { has: obj => "method" in obj, get: obj => obj.method }, metadata: _metadata }, null, _instanceExtraInitializers);
            __esDecorate(this, null, _get_x_decorators, { kind: "getter", name: "x", static: false, private: false, access: { has: obj => "x" in obj, get: obj => obj.x }, metadata: _metadata }, null, _instanceExtraInitializers);
            __esDecorate(this, null, _set_x_decorators, { kind: "setter", name: "x", static: false, private: false, access: { has: obj => "x" in obj, set: (obj, value) => { obj.x = value; } }, metadata: _metadata }, null, _instanceExtraInitializers);
            __esDecorate(this, null, _z_decorators, { kind: "accessor", name: "z", static: false, private: false, access: { has: obj => "z" in obj, get: obj => obj.z, set: (obj, value) => { obj.z = value; } }, metadata: _metadata }, _z_initializers, _z_extraInitializers);
            __esDecorate(null, null, _static_private_y_decorators, { kind: "field", name: "#y", static: true, private: true, access: { has: obj => #y in obj, get: obj => obj.#y, set: (obj, value) => { obj.#y = value; } }, metadata: _metadata }, _static_private_y_initializers, _static_private_y_extraInitializers);
            __esDecorate(null, null, _y_decorators, { kind: "field", name: "y", static: false, private: false, access: { has: obj => "y" in obj, get: obj => obj.y, set: (obj, value) => { obj.y = value; } }, metadata: _metadata }, _y_initializers, _y_extraInitializers);
            __esDecorate(null, _classDescriptor = { value: _classThis }, _classDecorators, { kind: "class", name: _classThis.name, metadata: _metadata }, null, _classExtraInitializers);
            C = _classThis = _classDescriptor.value;
            if (_metadata) Object.defineProperty(_classThis, Symbol.metadata, { enumerable: true, configurable: true, writable: true, value: _metadata });
        }
        method() { }
        get x() { return 1; }
        set x(value) { }
        y = (__runInitializers(this, _instanceExtraInitializers), __runInitializers(this, _y_initializers, 1));
        accessor z = (__runInitializers(this, _y_extraInitializers), __runInitializers(this, _z_initializers, 1));
        static get #method() { return _static_private_method_descriptor.value; }
        static get #x() { return _static_private_get_x_descriptor.get.call(this); }
        static set #x(value) { return _static_private_set_x_descriptor.set.call(this, value); }
        static #y = (__runInitializers(_classThis, _staticExtraInitializers), __runInitializers(_classThis, _static_private_y_initializers, 1));
        static accessor #z = (__runInitializers(_classThis, _static_private_y_extraInitializers), __runInitializers(_classThis, _static_private_z_initializers, 1));
        constructor() {
            __runInitializers(this, _z_extraInitializers);
        }
        static {
            __runInitializers(_classThis, _static_private_z_extraInitializers);
            __runInitializers(_classThis, _classExtraInitializers);
        }
    };
    return C = _classThis;
})();
//# sourceMappingURL=esDecorators-classDeclaration-sourceMap.js.map

//// [esDecorators-classDeclaration-sourceMap.d.ts]
declare var dec: any;
declare class C {
    #private;
    method(): void;
    get x(): number;
    set x(value: number);
    y: number;
    accessor z: number;
}
//# sourceMappingURL=esDecorators-classDeclaration-sourceMap.d.ts.map