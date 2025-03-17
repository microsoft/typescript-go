//// [tests/cases/conformance/jsdoc/jsdocTemplateTagNameResolution.ts] ////

//// [file.js]
/**
 * @template T
 * @template {keyof T} K
 * @typedef {T[K]} Foo
 */

const x = { a: 1 };

/** @type {Foo<typeof x, "a">} */
const y = "a";

//// [file.js]
const x = { a: 1 };
const y = "a";
