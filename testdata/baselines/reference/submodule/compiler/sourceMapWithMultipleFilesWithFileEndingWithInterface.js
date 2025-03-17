//// [tests/cases/compiler/sourceMapWithMultipleFilesWithFileEndingWithInterface.ts] ////

//// [a.ts]
module M {
    export var X = 1;
}
interface Navigator {
    getGamepads(func?: any): any;
    webkitGetGamepads(func?: any): any
    msGetGamepads(func?: any): any;
    webkitGamepads(func?: any): any;
}

//// [b.ts]
module m1 {
    export class c1 {
    }
}


//// [b.js]
var m1;
(function (m1) {
    class c1 {
    }
    m1.c1 = c1;
})(m1 || (m1 = {}));
//// [a.js]
var M;
(function (M) {
    M.X = 1;
})(M || (M = {}));
