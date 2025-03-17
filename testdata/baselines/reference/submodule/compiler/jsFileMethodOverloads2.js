//// [tests/cases/compiler/jsFileMethodOverloads2.ts] ////

//// [jsFileMethodOverloads2.js]
// Also works if all @overload tags are combined in one comment.
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
   *
   * @overload
   * @param {Example<string>} this
   * @returns {'string'}
   *
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
   *
   * @overload
   * @returns {T}
   *
   * @param {(y: T) => unknown} [fn]
   * @returns {unknown}
   */
  transform(fn) {
    return fn ? fn(this.value) : this.value;
  }
}


//// [jsFileMethodOverloads2.js]
class Example {
    constructor(value) {
        this.value = value;
    }
    getTypeName() {
        return typeof this.value;
    }
    transform(fn) {
        return fn ? fn(this.value) : this.value;
    }
}
