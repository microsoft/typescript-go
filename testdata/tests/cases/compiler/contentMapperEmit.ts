// @loadExternalPlugins: true

// @Filename: /tsconfig.json
{
    "compilerOptions": {
        "target": "es2020",
        "module": "esnext",
        "moduleResolution": "bundler"
    },
    "contentMappers": [
        { "package": "mapper", "extensions": [".mdx"] }
    ]
}

// @Filename: /node_modules/mapper/package.json
{
    "name": "mapper",
    "version": "1.0.0",
    "tsContentMapper": {
        "exec": ["verbatim-mapper"],
        "extensions": { ".mdx": ".ts" },
        "supportsEmit": true
    }
}

// @Filename: /component.mdx
export default function Component() {
    return "content";
}

// @Filename: /unsafe.mdx/index.ts
export const unsafe = true;

// @Filename: /main.ts
import Component from "./component.mdx";
import { unsafe } from "./unsafe.mdx";
const dynamicPath = "./component.mdx";
console.log(Component(), unsafe, import(dynamicPath));
