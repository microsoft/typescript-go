//// [tests/cases/conformance/externalModules/rewriteRelativeImportExtensions/emitModuleCommonJS.ts] ////

//// [a.js]
{
  require("" + "./foo.ts");
  import("" + "./foo.ts");
  require("./foo.ts");
  import("./foo.ts");
}

//// [b.ts]
{
  import("" + "./foo.ts");
  import("./foo.ts");
}


//// [a.js]
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
{
    require(__rewriteRelativeImportExtension("" + "./foo.ts"));
    import(__rewriteRelativeImportExtension("" + "./foo.ts"));
    require("./foo.js");
    import("./foo.js");
}
//// [b.js]
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
{
    import(__rewriteRelativeImportExtension("" + "./foo.ts"));
    import("./foo.js");
}
