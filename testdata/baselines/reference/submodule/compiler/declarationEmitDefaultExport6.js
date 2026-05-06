//// [tests/cases/compiler/declarationEmitDefaultExport6.ts] ////

//// [declarationEmitDefaultExport6.ts]
export class A {}
export default new A();


//// [declarationEmitDefaultExport6.js]
export class A {
}
export default new A();


//// [declarationEmitDefaultExport6.d.ts]
export class A {
}
const _default: A;
export default _default;


//// [DtsFileErrors]


declarationEmitDefaultExport6.d.ts(3,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitDefaultExport6.d.ts (1 errors) ====
    export class A {
    }
    const _default: A;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export default _default;
    