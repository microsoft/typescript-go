//// [tests/cases/conformance/types/spread/spreadDuplicate.ts] ////

//// [spreadDuplicate.ts]
// Repro from #44438

declare let a: { a: string };
declare let b: { a?: string };
declare let c: { a: string | undefined };
declare let d: { a?: string | undefined };

declare let t: boolean;

let a1 = { a: 123, ...a };  // string (Error)
let b1 = { a: 123, ...b };  // string | number
let c1 = { a: 123, ...c };  // string | undefined (Error)
let d1 = { a: 123, ...d };  // string | number

let a2 = { a: 123, ...(t ? a : {}) };  // string | number
let b2 = { a: 123, ...(t ? b : {}) };  // string | number
let c2 = { a: 123, ...(t ? c : {}) };  // string | number
let d2 = { a: 123, ...(t ? d : {}) };  // string | number


//// [spreadDuplicate.js]
"use strict";
// Repro from #44438
let a1 = Object.assign({ a: 123 }, a); // string (Error)
let b1 = Object.assign({ a: 123 }, b); // string | number
let c1 = Object.assign({ a: 123 }, c); // string | undefined (Error)
let d1 = Object.assign({ a: 123 }, d); // string | number
let a2 = Object.assign({ a: 123 }, (t ? a : {})); // string | number
let b2 = Object.assign({ a: 123 }, (t ? b : {})); // string | number
let c2 = Object.assign({ a: 123 }, (t ? c : {})); // string | number
let d2 = Object.assign({ a: 123 }, (t ? d : {})); // string | number


//// [spreadDuplicate.d.ts]
let a: {
    a: string;
};
let b: {
    a?: string;
};
let c: {
    a: string | undefined;
};
let d: {
    a?: string | undefined;
};
let t: boolean;
let a1: {
    a: string;
};
let b1: {
    a: string | number;
};
let c1: {
    a: string | undefined;
};
let d1: {
    a: string | number;
};
let a2: {
    a: string | number;
};
let b2: {
    a: string | number;
};
let c2: {
    a: string | number;
};
let d2: {
    a: string | number;
};
