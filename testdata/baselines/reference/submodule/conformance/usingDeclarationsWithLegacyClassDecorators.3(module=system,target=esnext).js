//// [tests/cases/conformance/statements/VariableStatements/usingDeclarations/usingDeclarationsWithLegacyClassDecorators.3.ts] ////

//// [usingDeclarationsWithLegacyClassDecorators.3.ts]
export {};

declare var dec: any;

using before = null;

@dec
export default class C {
}


//// [usingDeclarationsWithLegacyClassDecorators.3.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
using before = null;
@dec
class C {
}
exports.default = C;
