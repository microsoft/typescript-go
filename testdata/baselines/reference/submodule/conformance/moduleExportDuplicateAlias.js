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
export var apply: typeof a;
function a(): void;
export var apply: typeof a;
//// [test.d.ts]
export {};
