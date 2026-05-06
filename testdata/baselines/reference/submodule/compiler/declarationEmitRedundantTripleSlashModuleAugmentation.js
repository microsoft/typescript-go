//// [tests/cases/compiler/declarationEmitRedundantTripleSlashModuleAugmentation.ts] ////

//// [index.d.ts]
declare module "foo" {
    export interface Original {}
}

//// [augmentation.ts]
export interface FooOptions {}
declare module "foo" {
    export interface Augmentation {}
}

//// [index.ts]
import { Original, Augmentation } from "foo";
import type { FooOptions } from "./augmentation";
export interface _ {
    original: Original;
    augmentation: Augmentation;
    options: FooOptions;
}




//// [augmentation.d.ts]
export interface FooOptions {
}
module "foo" {
    interface Augmentation {
    }
}
//// [index.d.ts]
import { Original, Augmentation } from "foo";
import type { FooOptions } from "./augmentation";
export interface _ {
    original: Original;
    augmentation: Augmentation;
    options: FooOptions;
}


//// [DtsFileErrors]


/augmentation.d.ts(3,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== /node_modules/foo/index.d.ts (0 errors) ====
    declare module "foo" {
        export interface Original {}
    }
    
==== /augmentation.d.ts (1 errors) ====
    export interface FooOptions {
    }
    module "foo" {
    ~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        interface Augmentation {
        }
    }
    
==== /index.d.ts (0 errors) ====
    import { Original, Augmentation } from "foo";
    import type { FooOptions } from "./augmentation";
    export interface _ {
        original: Original;
        augmentation: Augmentation;
        options: FooOptions;
    }
    