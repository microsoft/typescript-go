// @loadExternalPlugins: true
// @noEmit: true

// @Filename: /tsconfig.json
{
    "compilerOptions": {
        "target": "es2020",
        "module": "esnext",
        "moduleResolution": "bundler",
        "strict": true
    },
    "contentMappers": [
        { "package": "mapper", "extensions": [".vue"] }
    ]
}

// @Filename: /node_modules/mapper/package.json
{
    "name": "mapper",
    "version": "1.0.0",
    "tsContentMapper": { "exec": ["supplemental-globals-mapper"], "extensions": { ".vue": ".ts" } }
}

// @Filename: /a.vue
a

// @Filename: /b.vue
b

// @Filename: /extra.d.ts
interface Extra { extra: boolean }
