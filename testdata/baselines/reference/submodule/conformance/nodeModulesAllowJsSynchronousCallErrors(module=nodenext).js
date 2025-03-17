//// [tests/cases/conformance/node/allowJs/nodeModulesAllowJsSynchronousCallErrors.ts] ////

//// [index.js]
// cjs format file
import {h} from "../index.js";
import mod = require("../index.js");
import {f as _f} from "./index.js";
import mod2 = require("./index.js");
export async function f() {
    const mod3 = await import ("../index.js");
    const mod4 = await import ("./index.js");
    h();
}
//// [index.js]
// esm format file
import {h as _h} from "./index.js";
import mod = require("./index.js");
import {f} from "./subfolder/index.js";
import mod2 = require("./subfolder/index.js");
export async function h() {
    const mod3 = await import ("./index.js");
    const mod4 = await import ("./subfolder/index.js");
    f();
}
//// [package.json]
{
    "name": "package",
    "private": true,
    "type": "module"
}
//// [package.json]
{
    "type": "commonjs"
}

//// [index.js]
import { createRequire as _createRequire } from "module";
const __require = _createRequire(import.meta.url);
import { h as _h } from "./index.js";
const mod = __require("./index.js");
import { f } from "./subfolder/index.js";
const mod2 = __require("./subfolder/index.js");
export async function h() {
    const mod3 = await import("./index.js");
    const mod4 = await import("./subfolder/index.js");
    f();
}
//// [index.js]
import { createRequire as _createRequire } from "module";
const __require = _createRequire(import.meta.url);
import { h } from "../index.js";
const mod = __require("../index.js");
import { f as _f } from "./index.js";
const mod2 = __require("./index.js");
export async function f() {
    const mod3 = await import("../index.js");
    const mod4 = await import("./index.js");
    h();
}
