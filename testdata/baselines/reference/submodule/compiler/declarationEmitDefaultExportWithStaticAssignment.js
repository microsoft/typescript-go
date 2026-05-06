//// [tests/cases/compiler/declarationEmitDefaultExportWithStaticAssignment.ts] ////

//// [foo.ts]
export class Foo {}

//// [index1.ts]
import {Foo} from './foo';
export default function Example() {}
Example.Foo = Foo

//// [index2.ts]
import {Foo} from './foo';
export {Foo};
export default function Example() {}
Example.Foo = Foo

//// [index3.ts]
export class Bar {}
export default function Example() {}

Example.Bar = Bar

//// [index4.ts]
function A() {  }

function B() { }

export function C() {
  return null;
}

C.A = A;
C.B = B;

//// [foo.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Foo = void 0;
class Foo {
}
exports.Foo = Foo;
//// [index1.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = Example;
const foo_1 = require("./foo");
function Example() { }
Example.Foo = foo_1.Foo;
//// [index2.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Foo = void 0;
exports.default = Example;
const foo_1 = require("./foo");
Object.defineProperty(exports, "Foo", { enumerable: true, get: function () { return foo_1.Foo; } });
function Example() { }
Example.Foo = foo_1.Foo;
//// [index3.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.Bar = void 0;
exports.default = Example;
class Bar {
}
exports.Bar = Bar;
function Example() { }
Example.Bar = Bar;
//// [index4.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.C = C;
function A() { }
function B() { }
function C() {
    return null;
}
C.A = A;
C.B = B;


//// [foo.d.ts]
export class Foo {
}
//// [index1.d.ts]
function Example(): void;
export default Example;
declare namespace Example {
    var Foo: typeof import("./foo").Foo;
}
//// [index2.d.ts]
import { Foo } from './foo';
export { Foo };
function Example(): void;
export default Example;
declare namespace Example {
    var Foo: typeof import("./foo").Foo;
}
//// [index3.d.ts]
export class Bar {
}
function Example(): void;
export default Example;
declare namespace Example {
    var Bar: typeof import("./index3").Bar;
}
//// [index4.d.ts]
export function C(): any;
export declare namespace C {
    var A: () => void;
}
export declare namespace C {
    var B: () => void;
}


//// [DtsFileErrors]


index1.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
index2.d.ts(3,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
index3.d.ts(3,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== foo.d.ts (0 errors) ====
    export class Foo {
    }
    
==== index1.d.ts (1 errors) ====
    function Example(): void;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export default Example;
    declare namespace Example {
        var Foo: typeof import("./foo").Foo;
    }
    
==== index2.d.ts (1 errors) ====
    import { Foo } from './foo';
    export { Foo };
    function Example(): void;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export default Example;
    declare namespace Example {
        var Foo: typeof import("./foo").Foo;
    }
    
==== index3.d.ts (1 errors) ====
    export class Bar {
    }
    function Example(): void;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export default Example;
    declare namespace Example {
        var Bar: typeof import("./index3").Bar;
    }
    
==== index4.d.ts (0 errors) ====
    export function C(): any;
    export declare namespace C {
        var A: () => void;
    }
    export declare namespace C {
        var B: () => void;
    }
    