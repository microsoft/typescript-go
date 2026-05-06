//// [tests/cases/compiler/jsFileMethodOverloads.ts] ////

//// [jsFileMethodOverloads.js]
/**
 * @template T
 */
 class Example {
  /**
   * @param {T} value 
   */
  constructor(value) {
    this.value = value;
  }

  /**
   * @overload
   * @param {Example<number>} this
   * @returns {'number'}
   */
  /**
   * @overload
   * @param {Example<string>} this
   * @returns {'string'}
   */
  /**
   * @returns {string}
   */
  getTypeName() {
    return typeof this.value;
  }

  /**
   * @template U
   * @overload
   * @param {(y: T) => U} fn
   * @returns {U}
   */
  /**
   * @overload
   * @returns {T}
   */
  /**
   * @param {(y: T) => unknown} [fn]
   * @returns {unknown}
   */
  transform(fn) {
    return fn ? fn(this.value) : this.value;
  }
}


//// [jsFileMethodOverloads.js]
"use strict";
/**
 * @template T
 */
class Example {
    /**
     * @param {T} value
     */
    constructor(value) {
        this.value = value;
    }
    /**
     * @overload
     * @param {Example<number>} this
     * @returns {'number'}
     */
    /**
     * @overload
     * @param {Example<string>} this
     * @returns {'string'}
     */
    /**
     * @returns {string}
     */
    getTypeName() {
        return typeof this.value;
    }
    /**
     * @template U
     * @overload
     * @param {(y: T) => U} fn
     * @returns {U}
     */
    /**
     * @overload
     * @returns {T}
     */
    /**
     * @param {(y: T) => unknown} [fn]
     * @returns {unknown}
     */
    transform(fn) {
        return fn ? fn(this.value) : this.value;
    }
}


//// [jsFileMethodOverloads.d.ts]
/**
 * @template T
 */
class Example<T> {
    value: T;
    /**
     * @param {T} value
     */
    constructor(value: T);
    /**
     * @overload
     * @param {Example<number>} this
     * @returns {'number'}
     */
    /**
     * @overload
     * @param {Example<string>} this
     * @returns {'string'}
     */
    /**
     * @returns {string}
     */
    getTypeName(this: Example<number>): 'number';
    /**
     * @overload
     * @param {Example<number>} this
     * @returns {'number'}
     */
    /**
     * @overload
     * @param {Example<string>} this
     * @returns {'string'}
     */
    /**
     * @returns {string}
     */
    getTypeName(this: Example<string>): 'string';
    /**
     * @template U
     * @overload
     * @param {(y: T) => U} fn
     * @returns {U}
     */
    /**
     * @overload
     * @returns {T}
     */
    /**
     * @param {(y: T) => unknown} [fn]
     * @returns {unknown}
     */
    transform<U>(fn: (y: T) => U): U;
    /**
     * @template U
     * @overload
     * @param {(y: T) => U} fn
     * @returns {U}
     */
    /**
     * @overload
     * @returns {T}
     */
    /**
     * @param {(y: T) => unknown} [fn]
     * @returns {unknown}
     */
    transform(): T;
}


//// [DtsFileErrors]


dist/jsFileMethodOverloads.d.ts(4,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== dist/jsFileMethodOverloads.d.ts (1 errors) ====
    /**
     * @template T
     */
    class Example<T> {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        value: T;
        /**
         * @param {T} value
         */
        constructor(value: T);
        /**
         * @overload
         * @param {Example<number>} this
         * @returns {'number'}
         */
        /**
         * @overload
         * @param {Example<string>} this
         * @returns {'string'}
         */
        /**
         * @returns {string}
         */
        getTypeName(this: Example<number>): 'number';
        /**
         * @overload
         * @param {Example<number>} this
         * @returns {'number'}
         */
        /**
         * @overload
         * @param {Example<string>} this
         * @returns {'string'}
         */
        /**
         * @returns {string}
         */
        getTypeName(this: Example<string>): 'string';
        /**
         * @template U
         * @overload
         * @param {(y: T) => U} fn
         * @returns {U}
         */
        /**
         * @overload
         * @returns {T}
         */
        /**
         * @param {(y: T) => unknown} [fn]
         * @returns {unknown}
         */
        transform<U>(fn: (y: T) => U): U;
        /**
         * @template U
         * @overload
         * @param {(y: T) => U} fn
         * @returns {U}
         */
        /**
         * @overload
         * @returns {T}
         */
        /**
         * @param {(y: T) => unknown} [fn]
         * @returns {unknown}
         */
        transform(): T;
    }
    