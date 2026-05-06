//// [tests/cases/conformance/declarationEmit/leaveOptionalParameterAsWritten.ts] ////

//// [a.ts]
export interface Foo {}

//// [b.ts]
import * as a from "./a";
declare global {
  namespace teams {
    export namespace calling {
      export import Foo = a.Foo;
    }
  }
}

//// [c.ts]
type Foo = teams.calling.Foo;
export const bar = (p?: Foo) => {}



//// [a.d.ts]
export interface Foo {
}
//// [b.d.ts]
import * as a from "./a";
global {
    namespace teams {
        namespace calling {
            export import Foo = a.Foo;
        }
    }
}
//// [c.d.ts]
type Foo = teams.calling.Foo;
export const bar: (p?: Foo) => void;
export {};


//// [DtsFileErrors]


dist/b.d.ts(2,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== dist/a.d.ts (0 errors) ====
    export interface Foo {
    }
    
==== dist/b.d.ts (1 errors) ====
    import * as a from "./a";
    global {
    ~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        namespace teams {
            namespace calling {
                export import Foo = a.Foo;
            }
        }
    }
    
==== dist/c.d.ts (0 errors) ====
    type Foo = teams.calling.Foo;
    export const bar: (p?: Foo) => void;
    export {};
    