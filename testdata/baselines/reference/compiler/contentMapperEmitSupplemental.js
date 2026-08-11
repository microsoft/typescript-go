//// [tests/cases/compiler/contentMapperEmitSupplemental.ts] ////

//// [package.json]
{
    "name": "mapper",
    "version": "1.0.0",
    "tsContentMapper": {
        "exec": ["supplemental-module-mapper"],
        "extensions": { ".mdx": ".ts" },
        "supportsEmit": true
    }
}

//// [component.mdx]
export default 1;

//// [component.mdx.0.js]
export const privateValue = "wrong";
//// [component.mdx.js]
export default 1;
