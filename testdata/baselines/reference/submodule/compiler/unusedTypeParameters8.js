//// [tests/cases/compiler/unusedTypeParameters8.ts] ////

//// [a.ts]
class C<T> { }

//// [b.ts]
interface C<T> { }

//// [b.js]
//// [a.js]
class C {
}
