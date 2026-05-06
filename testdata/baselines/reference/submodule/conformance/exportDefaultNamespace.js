//// [tests/cases/conformance/declarationEmit/exportDefaultNamespace.ts] ////

//// [exportDefaultNamespace.ts]
export default function someFunc() {
    return 'hello!';
}

someFunc.someProp = 'yo';


//// [exportDefaultNamespace.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = someFunc;
function someFunc() {
    return 'hello!';
}
someFunc.someProp = 'yo';


//// [exportDefaultNamespace.d.ts]
function someFunc(): string;
export default someFunc;
declare namespace someFunc {
    var someProp: string;
}


//// [DtsFileErrors]


exportDefaultNamespace.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== exportDefaultNamespace.d.ts (1 errors) ====
    function someFunc(): string;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export default someFunc;
    declare namespace someFunc {
        var someProp: string;
    }
    