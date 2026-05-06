//// [tests/cases/compiler/declarationEmitFBoundedTypeParams.ts] ////

//// [declarationEmitFBoundedTypeParams.ts]
// Repro from #6040

function append<a, b extends a>(result: a[], value: b): a[] {
    result.push(value);
    return result;
}


//// [declarationEmitFBoundedTypeParams.js]
"use strict";
// Repro from #6040
function append(result, value) {
    result.push(value);
    return result;
}


//// [declarationEmitFBoundedTypeParams.d.ts]
function append<a, b extends a>(result: a[], value: b): a[];


//// [DtsFileErrors]


declarationEmitFBoundedTypeParams.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitFBoundedTypeParams.d.ts (1 errors) ====
    function append<a, b extends a>(result: a[], value: b): a[];
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    