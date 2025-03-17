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

//// [c.js]
export const bar = (p) => { };
//// [b.js]
export {};
//// [a.js]
export {};
