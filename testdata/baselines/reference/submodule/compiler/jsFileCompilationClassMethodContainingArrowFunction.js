//// [tests/cases/compiler/jsFileCompilationClassMethodContainingArrowFunction.ts] ////

//// [a.js]
class c {
    method(a) {
        let x = a => this.method(a);
    }
}


//// [a.js]
class c {
    method(a) {
        let x = a => this.method(a);
    }
}
