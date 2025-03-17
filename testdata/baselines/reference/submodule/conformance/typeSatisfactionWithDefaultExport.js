//// [tests/cases/conformance/expressions/typeSatisfaction/typeSatisfactionWithDefaultExport.ts] ////

//// [a.ts]
interface Foo {
    a: number;
}
export default {} satisfies Foo;

//// [b.ts]
interface Foo {
    a: number;
}
export default { a: 1 } satisfies Foo;


//// [b.js]
export default { a: 1 };
//// [a.js]
export default {};
