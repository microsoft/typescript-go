// @loadExternalPlugins: true

// @Filename: /tsconfig.json
{
    "compilerOptions": {
        "target": "es2020",
        "module": "esnext",
        "moduleResolution": "bundler",
        "jsx": "preserve"
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
        "exec": ["tsx-verbatim-mapper"],
        "extensions": { ".mdx": ".tsx" },
        "supportsEmit": true
    }
}

// @Filename: /component.mdx
export default function Component() {
    return <div />;
}

// @Filename: /main.ts
import Component from "./component.mdx";
const dynamicPath = "./component.mdx";
console.log(Component(), import(dynamicPath));
