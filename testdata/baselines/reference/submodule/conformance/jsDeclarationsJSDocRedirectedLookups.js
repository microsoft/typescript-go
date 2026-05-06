//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsJSDocRedirectedLookups.ts] ////

//// [index.js]
// these are recognized as TS concepts by the checker
/** @type {String} */const a = "";
/** @type {Number} */const b = 0;
/** @type {Boolean} */const c = true;
/** @type {Void} */const d = undefined;
/** @type {Undefined} */const e = undefined;
/** @type {Null} */const f = null;

/** @type {Function} */const g = () => void 0;
/** @type {function} */const h = () => void 0;
/** @type {array} */const i = [];
/** @type {promise} */const j = Promise.resolve(0);
/** @type {Object<string, string>} */const k = {x: "x"};


// these are not recognized as anything and should just be lookup failures
// ignore the errors to try to ensure they're emitted as `any` in declaration emit
// @ts-ignore
/** @type {class} */const l = true;
// @ts-ignore
/** @type {bool} */const m = true;
// @ts-ignore
/** @type {int} */const n = true;
// @ts-ignore
/** @type {float} */const o = true;
// @ts-ignore
/** @type {integer} */const p = true;

// or, in the case of `event` likely erroneously refers to the type of the global Event object
/** @type {event} */const q = undefined;

//// [index.js]
"use strict";
// these are recognized as TS concepts by the checker
/** @type {String} */ const a = "";
/** @type {Number} */ const b = 0;
/** @type {Boolean} */ const c = true;
/** @type {Void} */ const d = undefined;
/** @type {Undefined} */ const e = undefined;
/** @type {Null} */ const f = null;
/** @type {Function} */ const g = () => void 0;
/** @type {function} */ const h = () => void 0;
/** @type {array} */ const i = [];
/** @type {promise} */ const j = Promise.resolve(0);
/** @type {Object<string, string>} */ const k = { x: "x" };
// these are not recognized as anything and should just be lookup failures
// ignore the errors to try to ensure they're emitted as `any` in declaration emit
// @ts-ignore
/** @type {class} */ const l = true;
// @ts-ignore
/** @type {bool} */ const m = true;
// @ts-ignore
/** @type {int} */ const n = true;
// @ts-ignore
/** @type {float} */ const o = true;
// @ts-ignore
/** @type {integer} */ const p = true;
// or, in the case of `event` likely erroneously refers to the type of the global Event object
/** @type {event} */ const q = undefined;


//// [index.d.ts]
/** @type {String} */ const a: string;
/** @type {Number} */ const b: number;
/** @type {Boolean} */ const c: boolean;
/** @type {Void} */ const d: void;
/** @type {Undefined} */ const e: undefined;
/** @type {Null} */ const f: null;
/** @type {Function} */ const g: Function;
/** @type {function} */ const h: Function;
/** @type {array} */ const i: any[];
/** @type {promise} */ const j: Promise<any>;
/** @type {Object<string, string>} */ const k: Record<string, string>;
/** @type {class} */ const l: class;
/** @type {bool} */ const m: bool;
/** @type {int} */ const n: int;
/** @type {float} */ const o: float;
/** @type {integer} */ const p: integer;
/** @type {event} */ const q: event;
