//// [tests/cases/compiler/declarationEmitObjectLiteralAccessorsJs1.ts] ////

//// [index.js]
// same type accessors
export const obj1 = {
  /**
   * my awesome getter (first in source order)
   * @returns {string}
   */
  get x() {
    return "";
  },
  /** 
   * my awesome setter (second in source order)
   * @param {string} a
   */
  set x(a) {},
};

// divergent accessors
export const obj2 = {
  /** 
   * my awesome getter
   * @returns {string}
   */
  get x() {
    return "";
  },
  /** 
   * my awesome setter
   * @param {number} a
   */
  set x(a) {},
};

export const obj3 = {
  /**
   * my awesome getter
   * @returns {string}
   */
  get x() {
    return "";
  },
};

export const obj4 = {
  /**
   * my awesome setter
   * @param {number} a
   */
  set x(a) {},
};


//// [index.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.obj4 = exports.obj3 = exports.obj2 = exports.obj1 = void 0;
exports.obj1 = {
    get x() {
        return "";
    },
    set x(a) { },
};
exports.obj2 = {
    get x() {
        return "";
    },
    set x(a) { },
};
exports.obj3 = {
    get x() {
        return "";
    },
};
exports.obj4 = {
    set x(a) { },
};
