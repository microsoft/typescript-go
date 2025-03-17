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


//// [test.js]
const { apply } = require('./moduleExportAliasDuplicateAlias');
apply();
