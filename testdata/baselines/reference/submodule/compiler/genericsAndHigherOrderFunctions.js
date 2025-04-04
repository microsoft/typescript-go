//// [tests/cases/compiler/genericsAndHigherOrderFunctions.ts] ////

//// [genericsAndHigherOrderFunctions.ts]
// no errors expected

var combine: <T, S>(f: (_: T) => S) =>
    <U>(g: (_: U) => T) =>
    (x: U) => S

    = <T, S>(f: (_: T) => S) =>
        <U>(g: (_: U) => T) =>
            (x: U) => f(g(x))

var foo: <K, N>(g: (x: K) => N) =>
    (h: <M>(_: (_: K) => (_: M) => M) => (_: M) => M) =>
    <R>(f: (_: N) => (_: R) => R) => (_: R) => R

    = <K, N>(g: (x: K) => N) =>
        (h: <M>(_: (_: K) => (_: M) => M) => (_: M) => M) =>
            <R>(f: (_: N) => (_: R) => R) => h(combine(f)(g))

//// [genericsAndHigherOrderFunctions.js]
// no errors expected
var combine = (f) => (g) => (x) => f(g(x));
var foo = (g) => (h) => (f) => h(combine(f)(g));
