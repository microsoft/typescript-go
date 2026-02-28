//// [tests/cases/compiler/classFieldsNamedEvaluationDestructuringAssignment.ts] ////

//// [classFieldsNamedEvaluationDestructuringAssignment.ts]
// Bug 2: Named evaluation missing in destructuring assignment elements.
// Anonymous class expressions used as default values in destructuring
// assignments should receive their inferred name.

let x: any;

// Array destructuring assignment with anonymous class default
[x = class { static #y = 1; }] = [];

// Object destructuring assignment (shorthand) with anonymous class default
({ x = class { static #z = 2; } } = {} as any);

// Object destructuring assignment (property) with anonymous class default
({ y: x = class { static #w = 3; } } = {} as any);


//// [classFieldsNamedEvaluationDestructuringAssignment.js]
"use strict";
// Bug 2: Named evaluation missing in destructuring assignment elements.
// Anonymous class expressions used as default values in destructuring
// assignments should receive their inferred name.
var _a, _x_y, _b, _z, _c, _x_w;
let x;
// Array destructuring assignment with anonymous class default
[x = (_a = class {
        },
        _x_y = { value: 1 },
        _a)] = [];
// Object destructuring assignment (shorthand) with anonymous class default
({ x = (_b = class {
        },
        _z = { value: 2 },
        _b) } = {});
// Object destructuring assignment (property) with anonymous class default
({ y: x = (_c = class {
        },
        _x_w = { value: 3 },
        _c) } = {});
