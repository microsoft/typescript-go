//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsFunctionClassesCjsExportAssignment.ts] ////

//// [timer.js]
/**
 * @param {number} timeout
 */
function Timer(timeout) {
    this.timeout = timeout;
}
module.exports = Timer;
//// [hook.js]
/**
 * @typedef {(arg: import("./context")) => void} HookHandler
 */
/**
 * @param {HookHandler} handle
 */
function Hook(handle) {
    this.handle = handle;
}
module.exports = Hook;

//// [context.js]
/**
 * Imports
 *
 * @typedef {import("./timer")} Timer
 * @typedef {import("./hook")} Hook
 * @typedef {import("./hook").HookHandler} HookHandler
 */

/**
 * Input type definition
 *
 * @typedef {Object} Input
 * @prop {Timer} timer
 * @prop {Hook} hook
 */
 
/**
 * State type definition
 *
 * @typedef {Object} State
 * @prop {Timer} timer
 * @prop {Hook} hook
 */

/**
 * New `Context`
 *
 * @class
 * @param {Input} input
 */

function Context(input) {
    if (!(this instanceof Context)) {
      return new Context(input)
    }
    this.state = this.construct(input);
}
Context.prototype = {
    /**
     * @param {Input} input
     * @param {HookHandler=} handle
     * @returns {State}
     */
    construct(input, handle = () => void 0) {
        return input;
    }
}
module.exports = Context;


//// [context.js]
function Context(input) {
    if (!(this instanceof Context)) {
        return new Context(input);
    }
    this.state = this.construct(input);
}
Context.prototype = {
    construct(input, handle = () => void 0) {
        return input;
    }
};
module.exports = Context;
//// [hook.js]
function Hook(handle) {
    this.handle = handle;
}
module.exports = Hook;
//// [timer.js]
function Timer(timeout) {
    this.timeout = timeout;
}
module.exports = Timer;
