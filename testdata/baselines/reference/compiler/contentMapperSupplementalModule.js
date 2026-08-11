//// [tests/cases/compiler/contentMapperSupplementalModule.ts] ////

//// [package.json]
{
    "name": "mapper",
    "version": "1.0.0",
    "tsContentMapper": { "exec": ["supplemental-module-mapper"], "extensions": { ".vue": ".ts" } }
}

//// [component.vue]
export default 1;



//// [component.vue.0.d.ts]
export declare const privateValue: number;
//// [component.d.vue.ts]
/// <reference path="./component.vue.0.d.ts" />
declare const _default = 1;
export default _default;
