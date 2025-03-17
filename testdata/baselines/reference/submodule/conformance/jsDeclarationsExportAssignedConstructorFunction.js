//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsExportAssignedConstructorFunction.ts] ////

//// [jsDeclarationsExportAssignedConstructorFunction.js]
/** @constructor */
module.exports.MyClass = function() {
    this.x = 1
}
module.exports.MyClass.prototype = {
    a: function() {
    }
}


//// [jsDeclarationsExportAssignedConstructorFunction.js]
module.exports.MyClass = function () {
    this.x = 1;
};
module.exports.MyClass.prototype = {
    a: function () {
    }
};
