//// [tests/cases/compiler/argumentsPropertyNameInJsMode1.ts] ////

//// [a.js]
const foo = {
   f1: (params) => { }
}

function f2(x) {
  foo.f1({ x, arguments: [] });
}

f2(1, 2, 3);


//// [a.js]
const foo = {
    f1: (params) => { }
};
function f2(x) {
    foo.f1({ x, arguments: [] });
}
f2(1, 2, 3);
