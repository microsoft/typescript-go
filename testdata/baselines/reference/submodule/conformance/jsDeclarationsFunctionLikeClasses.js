//// [tests/cases/conformance/jsdoc/declarations/jsDeclarationsFunctionLikeClasses.ts] ////

//// [source.js]
/**
 * @param {number} x
 * @param {number} y
 */
export function Point(x, y) {
    if (!(this instanceof Point)) {
        return new Point(x, y);
    }
    this.x = x;
    this.y = y;
}

//// [referencer.js]
import {Point} from "./source";

/**
 * @param {Point} p
 */
export function magnitude(p) {
    return Math.sqrt(p.x ** 2 + p.y ** 2);
}


//// [referencer.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.magnitude = magnitude;
const source_1 = require("./source");
function magnitude(p) {
    return Math.sqrt(p.x ** 2 + p.y ** 2);
}
//// [source.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Point = Point;
function Point(x, y) {
    if (!(this instanceof Point)) {
        return new Point(x, y);
    }
    this.x = x;
    this.y = y;
}
