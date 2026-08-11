//// [tests/cases/compiler/contentMapperEmitTsx.ts] ////

//// [package.json]
{
    "name": "mapper",
    "version": "1.0.0",
    "tsContentMapper": {
        "exec": ["tsx-verbatim-mapper"],
        "extensions": { ".mdx": ".tsx" },
        "supportsEmit": true
    }
}

//// [component.mdx]
export default function Component() {
    return <div />;
}

//// [main.ts]
import Component from "./component.mdx";
const dynamicPath = "./component.mdx";
console.log(Component(), import(dynamicPath));


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
import Component from "./component.mdx.jsx";
const dynamicPath = "./component.mdx";
console.log(Component(), import(__rewriteRelativeImportExtension(dynamicPath, true, { ".mdx": ".jsx" })));
//// [component.mdx.jsx]
export default function Component() {
    return <div />;
}
