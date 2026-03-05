//// [tests/cases/conformance/statements/VariableStatements/usingDeclarations/usingDeclarationsWithESClassDecorators.1.ts] ////

//// [usingDeclarationsWithESClassDecorators.1.ts]
export {};

declare var dec: any;

using before = null;

@dec
class C {
}


//// [usingDeclarationsWithESClassDecorators.1.js]
"use strict";
var before, C;
Object.defineProperty(exports, "__esModule", { value: true });
const env_1 = { stack: [], error: void 0, hasError: false };
try {
    before = __addDisposableResource(env_1, null, false);
    C = (() => {
        var _a;
        let _classDecorators = [dec];
        let _classDescriptor;
        let _classExtraInitializers = [];
        let _classThis;
        var C = (_a = class {
            },
            _classThis = _a,
            __setFunctionName(_a, "C"),
            (() => {
                const _metadata = typeof Symbol === "function" && Symbol.metadata ? Object.create(null) : void 0;
                __esDecorate(null, _classDescriptor = { value: _classThis }, _classDecorators, { kind: "class", name: _classThis.name, metadata: _metadata }, null, _classExtraInitializers);
                C = _classThis = _classDescriptor.value;
                if (_metadata) Object.defineProperty(_classThis, Symbol.metadata, { enumerable: true, configurable: true, writable: true, value: _metadata });
                __runInitializers(_classThis, _classExtraInitializers);
            })(),
            _a);
        return C = _classThis;
    })();
}
catch (e_1) {
    env_1.error = e_1;
    env_1.hasError = true;
}
finally {
    __disposeResources(env_1);
}
