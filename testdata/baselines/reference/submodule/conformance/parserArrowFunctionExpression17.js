//// [tests/cases/conformance/parser/ecmascript5/ArrowFunctionExpressions/parserArrowFunctionExpression17.ts] ////

//// [fileJs.js]
a ? b : (c) : d => e // Not legal JS; "Unexpected token ':'" at last colon

//// [fileTs.ts]
a ? b : (c) : d => e


//// [fileTs.js]
a ? b : (c) => e;
//// [fileJs.js]
a ? b : (c) => e;
