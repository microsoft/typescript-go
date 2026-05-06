//// [tests/cases/conformance/jsdoc/seeTag2.ts] ////

//// [seeTag2.ts]
/** @see {} empty*/
const a = ""

/** @see {aaaaaa} unknown name*/
const b = ""

/** @see {?????} invalid */
const c = ""

/** @see c without brace */
const d = ""

/** @see ?????? wowwwwww*/
const e = ""

/** @see {}*/
const f = ""

/** @see */
const g = ""


//// [seeTag2.js]
"use strict";
/** @see {} empty*/
const a = "";
/** @see {aaaaaa} unknown name*/
const b = "";
/** @see {?????} invalid */
const c = "";
/** @see c without brace */
const d = "";
/** @see ?????? wowwwwww*/
const e = "";
/** @see {}*/
const f = "";
/** @see */
const g = "";


//// [seeTag2.d.ts]
/** @see {} empty*/
const a = "";
/** @see {aaaaaa} unknown name*/
const b = "";
/** @see {?????} invalid */
const c = "";
/** @see c without brace */
const d = "";
/** @see ?????? wowwwwww*/
const e = "";
/** @see {}*/
const f = "";
/** @see */
const g = "";


//// [DtsFileErrors]


seeTag2.d.ts(2,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== seeTag2.d.ts (1 errors) ====
    /** @see {} empty*/
    const a = "";
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    /** @see {aaaaaa} unknown name*/
    const b = "";
    /** @see {?????} invalid */
    const c = "";
    /** @see c without brace */
    const d = "";
    /** @see ?????? wowwwwww*/
    const e = "";
    /** @see {}*/
    const f = "";
    /** @see */
    const g = "";
    