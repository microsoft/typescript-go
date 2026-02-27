//// [tests/cases/conformance/esDecorators/metadata/esDecoratorsMetadata5.ts] ////

//// [foo.ts]
declare var metadata: any;
class C {
    @metadata m() {}
}


//// [foo.js]
"use strict";
let C = (() => {
    let _instanceExtraInitializers = [];
    let _m_decorators;
    return class C {
        static {
            const _metadata = typeof Symbol === "function" && Symbol.metadata ? Object.create(null) : void 0;
            _m_decorators = [metadata];
            __esDecorate(this, null, _m_decorators, { kind: "method", name: "m", static: false, private: false, access: { has: obj => "m" in obj, get: obj => obj.m }, metadata: _metadata }, null, _instanceExtraInitializers);
            if (_metadata) Object.defineProperty(this, Symbol.metadata, { enumerable: true, configurable: true, writable: true, value: _metadata });
        }
        m() { }
        constructor() {
            __runInitializers(this, _instanceExtraInitializers);
        }
    };
})();
