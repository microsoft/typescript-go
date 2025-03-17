//// [tests/cases/compiler/declarationMapsOutFile2.ts] ////

//// [a.ts]
class Foo {
    doThing(x: {a: number}) {
        return {b: x.a};
    }
    static make() {
        return new Foo();
    }
}
//// [index.ts]
const c = new Foo();
c.doThing({a: 42});

let x = c.doThing({a: 12});


//// [index.js]
const c = new Foo();
c.doThing({ a: 42 });
let x = c.doThing({ a: 12 });
//// [a.js]
class Foo {
    doThing(x) {
        return { b: x.a };
    }
    static make() {
        return new Foo();
    }
}
