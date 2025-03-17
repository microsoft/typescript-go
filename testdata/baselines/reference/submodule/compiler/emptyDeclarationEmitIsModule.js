//// [tests/cases/compiler/emptyDeclarationEmitIsModule.ts] ////

//// [module.ts]
import * as i from "./index";
class Foo {}
//// [index.ts]
import {} from "./module";
export interface Bar {
    x: string
}

//// [module.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
class Foo {
}
//// [index.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
