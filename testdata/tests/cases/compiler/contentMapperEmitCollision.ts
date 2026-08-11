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
export const fromMdx = true;

// @Filename: /component.mdx.ts
export const fromTs = true;
