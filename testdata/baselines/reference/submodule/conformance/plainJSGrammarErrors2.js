//// [tests/cases/conformance/salsa/plainJSGrammarErrors2.ts] ////

//// [plainJSGrammarErrors2.js]

//// [a.js]
export default 1;

//// [b.js]
/**
 * @deprecated
 */
export { default as A } from "./a";


//// [b.js]
export { default as A } from "./a";
//// [a.js]
export default 1;
//// [plainJSGrammarErrors2.js]
