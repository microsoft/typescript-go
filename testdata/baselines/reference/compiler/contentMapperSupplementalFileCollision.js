//// [tests/cases/compiler/contentMapperSupplementalFileCollision.ts] ////

//// [package.json]
{
    "name": "mapper",
    "version": "1.0.0",
    "tsContentMapper": { "exec": ["supplemental-mapper"], "extensions": { ".astro": ".ts" } }
}

//// [component.astro]

//// [component.astro.0.ts]
export const existing = true;


//// [component.astro.0.js]
export const existing = true;
