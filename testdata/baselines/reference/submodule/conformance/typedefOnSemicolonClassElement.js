//// [tests/cases/conformance/jsdoc/typedefOnSemicolonClassElement.ts] ////

//// [typedefOnSemicolonClassElement.js]
export class Preferences {
  /** @typedef {string} A */
  ;
  /** @type {A} */
  a = 'ok'
}


//// [typedefOnSemicolonClassElement.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Preferences = void 0;
class Preferences {
    ;
    a = 'ok';
}
exports.Preferences = Preferences;
