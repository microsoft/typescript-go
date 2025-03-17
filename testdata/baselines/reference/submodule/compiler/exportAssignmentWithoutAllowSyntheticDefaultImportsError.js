//// [tests/cases/compiler/exportAssignmentWithoutAllowSyntheticDefaultImportsError.ts] ////

//// [bar.ts]
export = bar;
function bar() {}

//// [foo.ts]
import bar from './bar';

//// [foo.js]
export {};
//// [bar.js]
function bar() { }
export {};
