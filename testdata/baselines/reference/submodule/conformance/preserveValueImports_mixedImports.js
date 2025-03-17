//// [tests/cases/conformance/externalModules/typeOnly/preserveValueImports_mixedImports.ts] ////

//// [exports.ts]
export function Component() {}
export interface ComponentProps {}

//// [index.ts]
import { Component, ComponentProps } from "./exports.js";

//// [index.fixed.ts]
import { Component, type ComponentProps } from "./exports.js";


//// [index.fixed.js]
import "./exports.js";
//// [index.js]
export {};
//// [exports.js]
export function Component() { }
