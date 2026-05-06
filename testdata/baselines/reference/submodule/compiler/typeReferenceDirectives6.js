//// [tests/cases/compiler/typeReferenceDirectives6.ts] ////

//// [ref.d.ts]
declare let $: { x: number }
    
//// [index.d.ts]
interface $ { x }


//// [app.ts]
/// <reference path="./ref.d.ts"/>
/// <reference types="lib"/>

let x: $;
let y = () => x



//// [app.js]
"use strict";
/// <reference path="./ref.d.ts"/>
/// <reference types="lib"/>
let x;
let y = () => x;


//// [app.d.ts]
let x: $;
let y: () => $;


//// [DtsFileErrors]


/app.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== /app.d.ts (1 errors) ====
    let x: $;
    ~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    let y: () => $;
    
==== /ref.d.ts (0 errors) ====
    declare let $: { x: number }
        
==== /types/lib/index.d.ts (0 errors) ====
    interface $ { x }
    
    