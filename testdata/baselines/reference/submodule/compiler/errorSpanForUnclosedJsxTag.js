//// [tests/cases/compiler/errorSpanForUnclosedJsxTag.tsx] ////

//// [errorSpanForUnclosedJsxTag.tsx]
declare const React: any

let Foo = {
  Bar() {}
}

let Baz = () => {}

let x = <    Foo.Bar >Hello

let y = <   Baz >Hello

//// [errorSpanForUnclosedJsxTag.js]
let Foo = {
    Bar() { }
};
let Baz = () => { };
let x = React.createElement(Foo.Bar, null, "Hell let y = ", React.createElement(Baz, null, "Hello"));
