//// [tests/cases/conformance/jsdoc/instantiateTemplateTagTypeParameterOnVariableStatement.ts] ////

//// [instantiateTemplateTagTypeParameterOnVariableStatement.js]
/**
 * @template T
 * @param {T} a
 * @returns {(b: T) => T}
 */
const seq = a => b => b;

const text1 = "hello";
const text2 = "world";

/** @type {string} */
var text3 = seq(text1)(text2);




//// [instantiateTemplateTagTypeParameterOnVariableStatement.d.ts]
/**
 * @template T
 * @param {T} a
 * @returns {(b: T) => T}
 */
const seq: <T>(a: T) => (b: T) => T;
const text1 = "hello";
const text2 = "world";
/** @type {string} */
var text3: string;


//// [DtsFileErrors]


instantiateTemplateTagTypeParameterOnVariableStatement.d.ts(6,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== instantiateTemplateTagTypeParameterOnVariableStatement.d.ts (1 errors) ====
    /**
     * @template T
     * @param {T} a
     * @returns {(b: T) => T}
     */
    const seq: <T>(a: T) => (b: T) => T;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    const text1 = "hello";
    const text2 = "world";
    /** @type {string} */
    var text3: string;
    