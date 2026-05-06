//// [tests/cases/compiler/letDeclarations.ts] ////

//// [letDeclarations.ts]
let l1;
let l2: number;
let l3, l4, l5 :string, l6;

let l7 = false;
let l8: number = 23;
let l9 = 0, l10 :string = "", l11 = null;

for(let l11 in {}) { }

for(let l12 = 0; l12 < 9; l12++) { }


//// [letDeclarations.js]
"use strict";
let l1;
let l2;
let l3, l4, l5, l6;
let l7 = false;
let l8 = 23;
let l9 = 0, l10 = "", l11 = null;
for (let l11 in {}) { }
for (let l12 = 0; l12 < 9; l12++) { }


//// [letDeclarations.d.ts]
let l1: any;
let l2: number;
let l3: any, l4: any, l5: string, l6: any;
let l7: boolean;
let l8: number;
let l9: number, l10: string, l11: any;


//// [DtsFileErrors]


letDeclarations.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== letDeclarations.d.ts (1 errors) ====
    let l1: any;
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    let l2: number;
    let l3: any, l4: any, l5: string, l6: any;
    let l7: boolean;
    let l8: number;
    let l9: number, l10: string, l11: any;
    