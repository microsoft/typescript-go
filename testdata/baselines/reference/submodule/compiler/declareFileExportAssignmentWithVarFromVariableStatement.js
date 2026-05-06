//// [tests/cases/compiler/declareFileExportAssignmentWithVarFromVariableStatement.ts] ////

//// [declareFileExportAssignmentWithVarFromVariableStatement.ts]
namespace m2 {
    export interface connectModule {
        (res, req, next): void;
    }
    export interface connectExport {
        use: (mod: connectModule) => connectExport;
        listen: (port: number) => void;
    }

}

var x = 10, m2: {
    (): m2.connectExport;
    test1: m2.connectModule;
    test2(): m2.connectModule;
};

export = m2;

//// [declareFileExportAssignmentWithVarFromVariableStatement.js]
"use strict";
var x = 10, m2;
module.exports = m2;


//// [declareFileExportAssignmentWithVarFromVariableStatement.d.ts]
namespace m2 {
    interface connectModule {
        (res: any, req: any, next: any): void;
    }
    interface connectExport {
        use: (mod: connectModule) => connectExport;
        listen: (port: number) => void;
    }
}
var m2: {
    (): m2.connectExport;
    test1: m2.connectModule;
    test2(): m2.connectModule;
};
export = m2;


//// [DtsFileErrors]


declareFileExportAssignmentWithVarFromVariableStatement.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declareFileExportAssignmentWithVarFromVariableStatement.d.ts (1 errors) ====
    namespace m2 {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        interface connectModule {
            (res: any, req: any, next: any): void;
        }
        interface connectExport {
            use: (mod: connectModule) => connectExport;
            listen: (port: number) => void;
        }
    }
    var m2: {
        (): m2.connectExport;
        test1: m2.connectModule;
        test2(): m2.connectModule;
    };
    export = m2;
    