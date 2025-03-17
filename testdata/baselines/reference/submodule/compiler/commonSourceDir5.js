//// [tests/cases/compiler/commonSourceDir5.ts] ////

//// [bar.ts]
import {z} from "./foo";
export var x = z + z;

//// [foo.ts]
import {pi} from "B:/baz";
export var i = Math.sqrt(-1);
export var z = pi * pi;

//// [baz.ts]
import {x} from "A:/bar";
import {i} from "A:/foo";
export var pi = Math.PI;
export var y = x * i;

//// [bar.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.x = void 0;
const foo_1 = require("./foo");
exports.x = foo_1.z + foo_1.z;
//// [foo.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.z = exports.i = void 0;
const baz_1 = require("B:/baz");
exports.i = Math.sqrt(-1);
exports.z = baz_1.pi * baz_1.pi;
//// [baz.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.y = exports.pi = void 0;
const bar_1 = require("A:/bar");
const foo_1 = require("A:/foo");
exports.pi = Math.PI;
exports.y = bar_1.x * foo_1.i;
