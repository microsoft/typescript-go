//// [tests/cases/conformance/esDecorators/metadata/esDecoratorsMetadata4.ts] ////

//// [foo.ts]
const PRIVATE_METADATA = new WeakMap();

function meta(key: string, value: string) {
    return (_, context) => {
        let metadata = PRIVATE_METADATA.get(context.metadata);

        if (!metadata) {
            metadata = {};
            PRIVATE_METADATA.set(context.metadata, metadata);
        }

        metadata[key] = value;
    };
}

@meta('a', 'x')
class C {
    @meta('b', 'y')
    m() { }
}

PRIVATE_METADATA.get(C[Symbol.metadata]).a; // 'x'
PRIVATE_METADATA.get(C[Symbol.metadata]).b; // 'y'


//// [foo.js]
"use strict";
const PRIVATE_METADATA = new WeakMap();
function meta(key, value) {
    return (_, context) => {
        let metadata = PRIVATE_METADATA.get(context.metadata);
        if (!metadata) {
            metadata = {};
            PRIVATE_METADATA.set(context.metadata, metadata);
        }
        metadata[key] = value;
    };
}
let C = (() => {
    let _classDecorators = [meta('a', 'x')];
    let _classDescriptor;
    let _classExtraInitializers = [];
    let _classThis;
    let _instanceExtraInitializers = [];
    let _m_decorators;
    var C = class {
        static { _classThis = this; }
        static {
            const _metadata = typeof Symbol === "function" && Symbol.metadata ? Object.create(null) : void 0;
            _m_decorators = [meta('b', 'y')];
            __esDecorate(this, null, _m_decorators, { kind: "method", name: "m", static: false, private: false, access: { has: obj => "m" in obj, get: obj => obj.m }, metadata: _metadata }, null, _instanceExtraInitializers);
            __esDecorate(null, _classDescriptor = { value: _classThis }, _classDecorators, { kind: "class", name: _classThis.name, metadata: _metadata }, null, _classExtraInitializers);
            C = _classThis = _classDescriptor.value;
            if (_metadata) Object.defineProperty(_classThis, Symbol.metadata, { enumerable: true, configurable: true, writable: true, value: _metadata });
            __runInitializers(_classThis, _classExtraInitializers);
        }
        m() { }
        constructor() {
            __runInitializers(this, _instanceExtraInitializers);
        }
    };
    return C = _classThis;
})();
PRIVATE_METADATA.get(C[Symbol.metadata]).a; // 'x'
PRIVATE_METADATA.get(C[Symbol.metadata]).b; // 'y'
