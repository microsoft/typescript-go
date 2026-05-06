//// [tests/cases/conformance/importAssertion/importAssertion1.ts] ////

//// [0.ts]
export const a = 1;
export const b = 2;

//// [1.ts]
import './0' with { type: "json" }
import { a, b } from './0' with { "type": "json" }
import * as foo from './0' with { type: "json" }
a;
b;
foo.a;
foo.b;

//// [2.ts]
import { a, b } from './0' with {}
import { a as c, b as d } from './0' with { a: "a", b: "b", c: "c" }
a;
b;
c;
d;

//// [3.ts]
const a = import('./0')
const b = import('./0', { with: { type: "json" } })
const c = import('./0', { with: { type: "json", ttype: "typo" } })
const d = import('./0', { with: {} })
const dd = import('./0', {})
declare function foo(): any;
const e = import('./0', foo())
const f = import()
const g = import('./0', {}, {})
const h = import('./0', { with: { type: "json" }},)



//// [0.js]
export const a = 1;
export const b = 2;
//// [1.js]
import './0' with { type: "json" };
import { a, b } from './0' with { "type": "json" };
import * as foo from './0' with { type: "json" };
a;
b;
foo.a;
foo.b;
//// [2.js]
import { a, b } from './0' with {};
import { a as c, b as d } from './0' with { a: "a", b: "b", c: "c" };
a;
b;
c;
d;
//// [3.js]
"use strict";
const a = import('./0');
const b = import('./0', { with: { type: "json" } });
const c = import('./0', { with: { type: "json", ttype: "typo" } });
const d = import('./0', { with: {} });
const dd = import('./0', {});
const e = import('./0', foo());
const f = import();
const g = import('./0', {}, {});
const h = import('./0', { with: { type: "json" } });


//// [0.d.ts]
export const a = 1;
export const b = 2;
//// [1.d.ts]
import './0';
//// [2.d.ts]
export {};
//// [3.d.ts]
const a: Promise<typeof import("./0")>;
const b: Promise<typeof import("./0")>;
const c: Promise<typeof import("./0")>;
const d: Promise<typeof import("./0")>;
const dd: Promise<typeof import("./0")>;
function foo(): any;
const e: Promise<typeof import("./0")>;
const f: Promise<any>;
const g: Promise<typeof import("./0")>;
const h: Promise<typeof import("./0")>;
