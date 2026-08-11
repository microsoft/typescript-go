//// [tests/cases/compiler/rewriteRelativeImportExtensionsDynamicImportCjs.ts] ////

//// [index.cts]
import { getSpecifier } from "./specifier.cjs";

import(getSpecifier());

//// [specifier.cts]
export function getSpecifier() {
    return "./target.ts";
}

//// [target.ts]
export {};


//// [specifier.cjs]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.getSpecifier = getSpecifier;
function getSpecifier() {
    return "./target.ts";
}
//// [index.cjs]
"use strict";
var __rewriteRelativeImportExtension = (this && this.__rewriteRelativeImportExtension) || function (path, preserveJsx, extraExtensions) {
    if (typeof path === "string" && /^\.\.?\//.test(path)) {
        if (extraExtensions) {
            for (var extension in extraExtensions) {
                var outputExtension = extraExtensions[extension];
                if (path.length > extension.length && path.slice(-extension.length) === extension) return path + outputExtension;
            }
        }
        return path.replace(/\.(tsx)$|((?:\.d)?)((?:\.[^./]+?)?)\.([cm]?)ts$/i, function (m, tsx, d, ext, cm) {
            return tsx ? preserveJsx ? ".jsx" : ".js" : d && (!ext || !cm) ? m : (d + ext + "." + cm.toLowerCase() + "js");
        });
    }
    return path;
};
Object.defineProperty(exports, "__esModule", { value: true });
const specifier_cjs_1 = require("./specifier.cjs");
import(__rewriteRelativeImportExtension((0, specifier_cjs_1.getSpecifier)()));
//// [target.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
