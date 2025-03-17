//// [tests/cases/compiler/importEqualsError45874.ts] ////

//// [globals.ts]
namespace globals {
  export type Foo = {};
  export const Bar = {};
}
import Foo = globals.toString.Blah;

//// [index.ts]
const Foo = {};


//// [index.js]
const Foo = {};
//// [globals.js]
var globals;
(function (globals) {
    globals.Bar = {};
})(globals || (globals = {}));
