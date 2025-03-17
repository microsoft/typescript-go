//// [tests/cases/compiler/blockScopedNamespaceDifferentFile.ts] ////

//// [test.ts]
namespace C {
    export class Name {
        static funcData = A.AA.func();
        static someConst = A.AA.foo;

        constructor(parameters) {}
    }
}

//// [typings.d.ts]
declare namespace A {
    namespace AA {
        function func(): number;
        const foo = "";
    }
}


//// [test.js]
var C;
(function (C) {
    class Name {
        static funcData = A.AA.func();
        static someConst = A.AA.foo;
        constructor(parameters) { }
    }
    C.Name = Name;
})(C || (C = {}));
