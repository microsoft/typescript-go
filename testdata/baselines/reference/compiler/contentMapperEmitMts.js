//// [tests/cases/compiler/contentMapperEmitMts.ts] ////

//// [package.json]
{
    "name": "mapper",
    "version": "1.0.0",
    "tsContentMapper": {
        "exec": ["verbatim-mapper"],
        "extensions": { ".component": ".mts" },
        "supportsEmit": true
    }
}

//// [value.component]
export default 1;

//// [main.mts]
import value from "./value.component";
console.log(value);


//// [main.mjs]
import value from "./value.component.mjs";
console.log(value);
//// [value.component.mjs]
export default 1;
