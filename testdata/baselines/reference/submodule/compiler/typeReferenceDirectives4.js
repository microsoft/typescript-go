//// [tests/cases/compiler/typeReferenceDirectives4.ts] ////

//// [ref.d.ts]
interface $ { x }

//// [index.d.ts]
declare let $: { x: number }


//// [app.ts]
/// <reference path="./ref.d.ts"/>
/// <reference types="lib" preserve="true" />

let x: $;
let y = () => x

//// [app.js]
"use strict";
/// <reference path="./ref.d.ts"/>
/// <reference types="lib" preserve="true" />
let x;
let y = () => x;


//// [app.d.ts]
/// <reference types="lib" preserve="true" />
let x: $;
let y: () => $;


//// [DtsFileErrors]


/app.d.ts(2,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
/app.d.ts(2,8): error TS2749: '$' refers to a value, but is being used as a type here. Did you mean 'typeof $'?
/app.d.ts(3,14): error TS2749: '$' refers to a value, but is being used as a type here. Did you mean 'typeof $'?


==== /app.d.ts (3 errors) ====
    /// <reference types="lib" preserve="true" />
    let x: $;
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
           ~
!!! error TS2749: '$' refers to a value, but is being used as a type here. Did you mean 'typeof $'?
    let y: () => $;
                 ~
!!! error TS2749: '$' refers to a value, but is being used as a type here. Did you mean 'typeof $'?
    
==== /ref.d.ts (0 errors) ====
    interface $ { x }
    
==== /types/lib/index.d.ts (0 errors) ====
    declare let $: { x: number }
    
    