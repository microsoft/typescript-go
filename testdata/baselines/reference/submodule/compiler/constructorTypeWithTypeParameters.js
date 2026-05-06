//// [tests/cases/compiler/constructorTypeWithTypeParameters.ts] ////

//// [constructorTypeWithTypeParameters.ts]
declare var X: {
    new <T>(): number;
}
declare var Y: {
    new (): number;
}
var anotherVar: new <T>() => number;

//// [constructorTypeWithTypeParameters.js]
"use strict";
var anotherVar;


//// [constructorTypeWithTypeParameters.d.ts]
var X: {
    new <T>(): number;
};
var Y: {
    new (): number;
};
var anotherVar: new <T>() => number;


//// [DtsFileErrors]


constructorTypeWithTypeParameters.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== constructorTypeWithTypeParameters.d.ts (1 errors) ====
    var X: {
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        new <T>(): number;
    };
    var Y: {
        new (): number;
    };
    var anotherVar: new <T>() => number;
    