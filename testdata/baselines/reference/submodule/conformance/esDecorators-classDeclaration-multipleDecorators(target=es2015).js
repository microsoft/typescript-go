//// [tests/cases/conformance/esDecorators/classDeclaration/esDecorators-classDeclaration-multipleDecorators.ts] ////

//// [esDecorators-classDeclaration-multipleDecorators.ts]
declare let dec1: any, dec2: any;

@dec1
@dec2
class C {
}


//// [esDecorators-classDeclaration-multipleDecorators.js]
"use strict";
let C = (() => {
    var _a;
    let _classDecorators = [dec1, dec2];
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
            _a = _classThis = _classDescriptor.value;
            if (_metadata) Object.defineProperty(_classThis, Symbol.metadata, { enumerable: true, configurable: true, writable: true, value: _metadata });
            __runInitializers(_classThis, _classExtraInitializers);
        })(),
        _a);
    return C = _classThis;
})();
