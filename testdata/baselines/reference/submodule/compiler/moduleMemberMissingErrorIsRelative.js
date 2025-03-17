//// [tests/cases/compiler/moduleMemberMissingErrorIsRelative.ts] ////

//// [foo.ts]
export {};
//// [bar.ts]
import {nosuch} from './foo';

//// [bar.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [foo.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
