//// [tests/cases/compiler/packageJsonImportsWildcardMatchesTsExtension.ts] ////

//// [package.json]
{
    "type": "module",
    "imports": {
        "#/*.omg": "./src/*"
    }
}

//// [foo.ts]
export function hello() {
    return "world";
}

//// [index.ts]
import { hello } from "#/foo.ts.omg";

hello();


//// [foo.js]
export function hello() {
    return "world";
}
//// [index.js]
import { hello } from "#/foo.ts.omg";
hello();
