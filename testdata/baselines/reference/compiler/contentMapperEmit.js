//// [tests/cases/compiler/contentMapperEmit.ts] ////

//// [package.json]
{
    "name": "mapper",
    "version": "1.0.0",
    "tsContentMapper": {
        "exec": ["verbatim-mapper"],
        "extensions": { ".mdx": ".ts" },
        "supportsEmit": true
    }
}

//// [component.mdx]
export default function Component() {
    return "content";
}

//// [index.ts]
export const unsafe = true;

//// [main.ts]
import Component from "./component.mdx";
import { unsafe } from "./unsafe.mdx";
const dynamicPath = "./component.mdx";
console.log(Component(), unsafe, import(dynamicPath));


//// [index.js]
export const unsafe = true;
//// [main.js]
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
import Component from "./component.mdx.js";
import { unsafe } from "./unsafe.mdx.js";
const dynamicPath = "./component.mdx";
console.log(Component(), unsafe, import(__rewriteRelativeImportExtension(dynamicPath, false, { ".mdx": ".js" })));
//// [component.mdx.js]
export default function Component() {
    return "content";
}
