//// [tests/cases/conformance/importAttributes/importAttributes8.ts] ////

//// [a.ts]
export default {
    a: "a",
    b: "b",
}

//// [b.ts]
import a from "./a" with { a: "a", "b": "b" }; // ok


//// [b.js]
export {};
//// [a.js]
export default {
    a: "a",
    b: "b",
};
