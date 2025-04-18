//// [tests/cases/conformance/interfaces/declarationMerging/mergedInterfacesWithIndexers.ts] ////

//// [mergedInterfacesWithIndexers.ts]
// indexers should behave like other members when merging interface declarations

interface A {
    [x: number]: string;
}


interface A {
    [x: string]: { length: number };
}

var a: A;
var r = a[1];
var r2 = a['1'];
var r3 = a['hi'];

//// [mergedInterfacesWithIndexers.js]
var a;
var r = a[1];
var r2 = a['1'];
var r3 = a['hi'];
