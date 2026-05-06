//// [tests/cases/compiler/declareFileExportAssignment.ts] ////

//// [declareFileExportAssignment.ts]
namespace m2 {
    export interface connectModule {
        (res, req, next): void;
    }
    export interface connectExport {
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

//// [declareFileExportAssignment.js]
"use strict";
var m2;
module.exports = m2;


//// [declareFileExportAssignment.d.ts]
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


declareFileExportAssignment.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declareFileExportAssignment.d.ts (1 errors) ====
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
    