//// [tests/cases/compiler/declFileAmbientExternalModuleWithSingleExportedModule.ts] ////

//// [declFileAmbientExternalModuleWithSingleExportedModule_0.ts]
declare module "SubModule" {
    export namespace m {
        export namespace m3 {
            interface c {
            }
        }
    }
}

//// [declFileAmbientExternalModuleWithSingleExportedModule_1.ts]
///<reference path='declFileAmbientExternalModuleWithSingleExportedModule_0.ts' preserve="true" />
import SubModule = require('SubModule');
export var x: SubModule.m.m3.c;



//// [declFileAmbientExternalModuleWithSingleExportedModule_0.js]
"use strict";
//// [declFileAmbientExternalModuleWithSingleExportedModule_1.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.x = void 0;


//// [declFileAmbientExternalModuleWithSingleExportedModule_0.d.ts]
module "SubModule" {
    namespace m {
        namespace m3 {
            interface c {
            }
        }
    }
}
//// [declFileAmbientExternalModuleWithSingleExportedModule_1.d.ts]
/// <reference path="declFileAmbientExternalModuleWithSingleExportedModule_0.d.ts" preserve="true" />
import SubModule = require('SubModule');
export var x: SubModule.m.m3.c;


//// [DtsFileErrors]


declFileAmbientExternalModuleWithSingleExportedModule_0.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileAmbientExternalModuleWithSingleExportedModule_1.d.ts (0 errors) ====
    /// <reference path="declFileAmbientExternalModuleWithSingleExportedModule_0.d.ts" preserve="true" />
    import SubModule = require('SubModule');
    export var x: SubModule.m.m3.c;
    
==== declFileAmbientExternalModuleWithSingleExportedModule_0.d.ts (1 errors) ====
    module "SubModule" {
    ~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        namespace m {
            namespace m3 {
                interface c {
                }
            }
        }
    }
    