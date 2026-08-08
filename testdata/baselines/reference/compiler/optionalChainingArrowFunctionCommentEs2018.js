//// [tests/cases/compiler/optionalChainingArrowFunctionCommentEs2018.ts] ////

//// [optionalChainingArrowFunctionCommentEs2018.ts]
const thing = { nested: { condition: true } };
const wat = () =>
    // explanatory comment
    thing?.nested?.condition ? "pass" : "fail";


//// [optionalChainingArrowFunctionCommentEs2018.js]
"use strict";
const thing = { nested: { condition: true } };
const wat = () => { var _a; 
// explanatory comment
return ((_a = thing === null || thing === void 0 ? void 0 : thing.nested) === null || _a === void 0 ? void 0 : _a.condition) ? "pass" : "fail"; };
