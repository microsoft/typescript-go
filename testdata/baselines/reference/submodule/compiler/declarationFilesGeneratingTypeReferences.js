//// [tests/cases/compiler/declarationFilesGeneratingTypeReferences.ts] ////

//// [index.d.ts]
interface JQuery {

}

//// [app.ts]
/// <reference types="jquery" preserve="true" />
namespace Test {
    export var x: JQuery;
}


//// [app.js]
var Test;
(function (Test) {
})(Test || (Test = {}));
