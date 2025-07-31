//// [tests/cases/compiler/declarationEmitModuleWithScopeMarker.ts] ////

//// [declarationEmitModuleWithScopeMarker.ts]
declare module "bar" {
    var before: typeof func;

    export function normal(): void;

    export default function func(): typeof func;

    var after: typeof func;

    export {}
}


//// [declarationEmitModuleWithScopeMarker.js]


//// [declarationEmitModuleWithScopeMarker.d.ts]
declare module "bar" {
    var before: typeof func;
    export function normal(): void;
    function func(): typeof func;
    var after: typeof func;
    export {};
    export default func;
}
export default func;


//// [DtsFileErrors]


declarationEmitModuleWithScopeMarker.d.ts(9,16): error TS2304: Cannot find name 'func'.


==== declarationEmitModuleWithScopeMarker.d.ts (1 errors) ====
    declare module "bar" {
        var before: typeof func;
        export function normal(): void;
        function func(): typeof func;
        var after: typeof func;
        export {};
        export default func;
    }
    export default func;
                   ~~~~
!!! error TS2304: Cannot find name 'func'.
    