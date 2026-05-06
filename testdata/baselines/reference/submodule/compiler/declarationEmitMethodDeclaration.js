//// [tests/cases/compiler/declarationEmitMethodDeclaration.ts] ////

//// [a.js]
export default {
    methods: {
        foo() { }
    }
}




//// [a.d.ts]
const _default: {
    methods: {
        foo(): void;
    };
};
export default _default;


//// [DtsFileErrors]


/a.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== /a.d.ts (1 errors) ====
    const _default: {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        methods: {
            foo(): void;
        };
    };
    export default _default;
    