//// [tests/cases/conformance/dynamicImport/importCallExpressionIncorrect1.ts] ////

//// [0.ts]
export function foo() { return "foo"; }

//// [1.ts]
import
import { foo } from './0';


//// [1.js]
import ;
//// [0.js]
export function foo() { return "foo"; }
