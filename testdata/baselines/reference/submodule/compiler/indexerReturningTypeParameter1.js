//// [tests/cases/compiler/indexerReturningTypeParameter1.ts] ////

//// [indexerReturningTypeParameter1.ts]
interface f {
    groupBy<T>(): { [key: string]: T[]; };
}
var a: f;
var r = a.groupBy();

class c {
    groupBy<T>(): { [key: string]: T[]; } {
        return null;
    }
}
var a2: c;
var r2 = a2.groupBy();

//// [indexerReturningTypeParameter1.js]
var a;
var r = a.groupBy();
class c {
    groupBy() {
        return null;
    }
}
var a2;
var r2 = a2.groupBy();
