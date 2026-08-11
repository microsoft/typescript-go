// @loadExternalPlugins: true

// @Filename: /tsconfig.json
{
    "compilerOptions": {
        "target": "es2020",
        "module": "nodenext",
        "moduleResolution": "nodenext"
    },
    "contentMappers": [
        { "package": "mapper", "extensions": [".component"] }
    ]
}

// @Filename: /node_modules/mapper/package.json
{
    "name": "mapper",
    "version": "1.0.0",
    "tsContentMapper": {
        "exec": ["verbatim-mapper"],
        "extensions": { ".component": ".mts" },
        "supportsEmit": true
    }
}

// @Filename: /value.component
export default 1;

// @Filename: /main.mts
import value from "./value.component";
console.log(value);
