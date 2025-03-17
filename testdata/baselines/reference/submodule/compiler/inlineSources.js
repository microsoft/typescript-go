//// [tests/cases/compiler/inlineSources.ts] ////

//// [a.ts]
var a = 0;
console.log(a);

//// [b.ts]
var b = 0;
console.log(b);


//// [b.js]
var b = 0;
console.log(b);
//// [a.js]
var a = 0;
console.log(a);
