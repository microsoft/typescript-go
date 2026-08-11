//// [tests/cases/compiler/contentMapperEmitNoEmit.ts] ////

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

//// [main.ts]
import Component from "./component.mdx";
console.log(Component());


//// [main.js]
import Component from "./component.mdx";
console.log(Component());
