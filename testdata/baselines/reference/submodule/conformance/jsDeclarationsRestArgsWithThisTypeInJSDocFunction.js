//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsRestArgsWithThisTypeInJSDocFunction.ts] ////

//// [bug38550.js]
export class Clazz {
  /**
   * @param {function(this:Object, ...*):*} functionDeclaration
   */
  method(functionDeclaration) {}
}


//// [bug38550.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Clazz = void 0;
class Clazz {
    method(functionDeclaration) { }
}
exports.Clazz = Clazz;
