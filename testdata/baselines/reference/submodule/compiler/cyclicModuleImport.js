//// [tests/cases/compiler/cyclicModuleImport.ts] ////

//// [cyclicModuleImport.ts]
declare module "SubModule" {
    import MainModule = require('MainModule');
    class SubModule {
        public static StaticVar: number;
        public InstanceVar: number;
        public main: MainModule;
        constructor();
    }
    export = SubModule;
}
declare module "MainModule" {
    import SubModule = require('SubModule');
    class MainModule {
        public SubModule: SubModule;
        constructor();
    }
    export = MainModule;
}


//// [cyclicModuleImport.js]
"use strict";


//// [cyclicModuleImport.d.ts]
module "SubModule" {
    import MainModule = require('MainModule');
    class SubModule {
        static StaticVar: number;
        InstanceVar: number;
        main: MainModule;
        constructor();
    }
    export = SubModule;
}
module "MainModule" {
    import SubModule = require('SubModule');
    class MainModule {
        SubModule: SubModule;
        constructor();
    }
    export = MainModule;
}


//// [DtsFileErrors]


cyclicModuleImport.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== cyclicModuleImport.d.ts (1 errors) ====
    module "SubModule" {
    ~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        import MainModule = require('MainModule');
        class SubModule {
            static StaticVar: number;
            InstanceVar: number;
            main: MainModule;
            constructor();
        }
        export = SubModule;
    }
    module "MainModule" {
        import SubModule = require('SubModule');
        class MainModule {
            SubModule: SubModule;
            constructor();
        }
        export = MainModule;
    }
    