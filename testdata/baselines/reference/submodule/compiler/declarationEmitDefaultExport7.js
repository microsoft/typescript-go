//// [tests/cases/compiler/declarationEmitDefaultExport7.ts] ////

//// [declarationEmitDefaultExport7.ts]
class A {}
export default new A();


//// [declarationEmitDefaultExport7.js]
class A {
}
export default new A();


//// [declarationEmitDefaultExport7.d.ts]
class A {
}
const _default: A;
export default _default;


//// [DtsFileErrors]


declarationEmitDefaultExport7.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declarationEmitDefaultExport7.d.ts (1 errors) ====
    class A {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    }
    const _default: A;
    export default _default;
    