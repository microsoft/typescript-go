//// [tests/cases/compiler/internalAliasInterfaceInsideTopLevelModuleWithExport.ts] ////

//// [internalAliasInterfaceInsideTopLevelModuleWithExport.ts]
export namespace a {
    export interface I {
    }
}

export import b = a.I;
export var x: b;


//// [internalAliasInterfaceInsideTopLevelModuleWithExport.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.x = void 0;


//// [internalAliasInterfaceInsideTopLevelModuleWithExport.d.ts]
export namespace a {
    interface I {
    }
}
export import b = a.I;
export var x: b;
