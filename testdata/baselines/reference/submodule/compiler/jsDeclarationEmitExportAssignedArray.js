//// [tests/cases/compiler/jsDeclarationEmitExportAssignedArray.ts] ////

//// [file.js]
module.exports = [{ name: 'other', displayName: 'Other', defaultEnabled: true }];



//// [file.d.ts]
const _default: {
    name: string;
    displayName: string;
    defaultEnabled: boolean;
}[];
export = _default;


//// [DtsFileErrors]


file.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== file.d.ts (1 errors) ====
    const _default: {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        name: string;
        displayName: string;
        defaultEnabled: boolean;
    }[];
    export = _default;
    