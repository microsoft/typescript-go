//// [tests/cases/compiler/fixCrashAliasLookupForDefauledImport.ts] ////

//// [input.ts]
export type Foo<T = string> = {};

//// [usage.ts]
import {Foo} from "./input";

function bar<T>(element: Foo) {
    return 1;
}

bar(1 as Foo<number>);


//// [usage.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
function bar(element) {
    return 1;
}
bar(1);
//// [input.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
