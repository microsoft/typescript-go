//// [tests/cases/compiler/constDeclarations.ts] ////

//// [constDeclarations.ts]
// No error
const c1 = false;
const c2: number = 23;
const c3 = 0, c4 :string = "", c5 = null;


for(const c4 = 0; c4 < 9; ) { break; }


for(const c5 = 0, c6 = 0; c5 < c6; ) { break; }

//// [constDeclarations.js]
"use strict";
// No error
const c1 = false;
const c2 = 23;
const c3 = 0, c4 = "", c5 = null;
for (const c4 = 0; c4 < 9;) {
    break;
}
for (const c5 = 0, c6 = 0; c5 < c6;) {
    break;
}


//// [constDeclarations.d.ts]
const c1 = false;
const c2: number;
const c3 = 0, c4: string, c5: any;


//// [DtsFileErrors]


constDeclarations.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== constDeclarations.d.ts (1 errors) ====
    const c1 = false;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    const c2: number;
    const c3 = 0, c4: string, c5: any;
    