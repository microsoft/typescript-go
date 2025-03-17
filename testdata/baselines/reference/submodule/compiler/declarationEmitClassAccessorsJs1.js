//// [tests/cases/compiler/declarationEmitClassAccessorsJs1.ts] ////

//// [index.js]
// https://github.com/microsoft/TypeScript/issues/58167

export class VFile {
  /**
   * @returns {string}
   */
  get path() {
    return ''
  }

  /**
   * @param {URL | string} path
   */
  set path(path) {
  }
}


//// [index.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.VFile = void 0;
class VFile {
    get path() {
        return '';
    }
    set path(path) {
    }
}
exports.VFile = VFile;
