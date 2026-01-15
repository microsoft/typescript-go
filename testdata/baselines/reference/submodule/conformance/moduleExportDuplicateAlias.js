//// [tests/cases/conformance/salsa/moduleExportDuplicateAlias.ts] ////

//// [moduleExportAliasDuplicateAlias.js]
exports.apply = undefined;
function a() { }
exports.apply()
exports.apply = a;
exports.apply()

//// [test.js]
const { apply } = require('./moduleExportAliasDuplicateAlias')
apply()


//// [moduleExportAliasDuplicateAlias.js]
"use strict";
exports.apply = undefined;
function a() { }
exports.apply();
exports.apply = a;
exports.apply();
//// [test.js]
"use strict";
const { apply } = require('./moduleExportAliasDuplicateAlias');
apply();


//// [moduleExportAliasDuplicateAlias.d.ts]
export declare var apply: undefined;
export declare var apply: undefined;
//// [test.d.ts]
export {};
