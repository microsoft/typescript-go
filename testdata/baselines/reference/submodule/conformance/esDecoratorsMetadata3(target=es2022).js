//// [tests/cases/conformance/esDecorators/metadata/esDecoratorsMetadata3.ts] ////

//// [foo.ts]
function appendMeta(key: string, value: string) {
    return (_, context) => {
        const existing = context.metadata[key] ?? [];
        context.metadata[key] = [...existing, value];
    };
}

@appendMeta('a', 'x')
class C {
}

@appendMeta('a', 'z')
class D extends C {
}

C[Symbol.metadata].a; // ['x']
D[Symbol.metadata].a; // ['x', 'z']


//// [foo.js]
"use strict";
function appendMeta(key, value) {
    return (_, context) => {
        const existing = context.metadata[key] ?? [];
        context.metadata[key] = [...existing, value];
    };
}
let C = (() => {
    let _classDecorators = [appendMeta('a', 'x')];
    let _classDescriptor;
    let _classExtraInitializers = [];
    let _classThis;
    var C = class {
        static { _classThis = this; }
        static {
            const _metadata = typeof Symbol === "function" && Symbol.metadata ? Object.create(null) : void 0;
            __esDecorate(null, _classDescriptor = { value: _classThis }, _classDecorators, { kind: "class", name: _classThis.name, metadata: _metadata }, null, _classExtraInitializers);
            C = _classThis = _classDescriptor.value;
            if (_metadata) Object.defineProperty(_classThis, Symbol.metadata, { enumerable: true, configurable: true, writable: true, value: _metadata });
            __runInitializers(_classThis, _classExtraInitializers);
        }
    };
    return C = _classThis;
})();
let D = (() => {
    let _classDecorators = [appendMeta('a', 'z')];
    let _classDescriptor;
    let _classExtraInitializers = [];
    let _classThis;
    let _classSuper = C;
    var D = class extends _classSuper {
        static { _classThis = this; }
        static {
            const _metadata = typeof Symbol === "function" && Symbol.metadata ? Object.create(_classSuper[Symbol.metadata] ?? null) : void 0;
            __esDecorate(null, _classDescriptor = { value: _classThis }, _classDecorators, { kind: "class", name: _classThis.name, metadata: _metadata }, null, _classExtraInitializers);
            D = _classThis = _classDescriptor.value;
            if (_metadata) Object.defineProperty(_classThis, Symbol.metadata, { enumerable: true, configurable: true, writable: true, value: _metadata });
            __runInitializers(_classThis, _classExtraInitializers);
        }
    };
    return D = _classThis;
})();
C[Symbol.metadata].a; // ['x']
D[Symbol.metadata].a; // ['x', 'z']
