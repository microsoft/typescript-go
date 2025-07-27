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
export declare class Foo {
}
//// [index1.d.ts]
import { Foo } from './foo';
declare function Example(): void;
declare namespace Example {
    var Foo: typeof Foo;
}
export default Example;
//// [index2.d.ts]
import { Foo } from './foo';
export { Foo };
declare function Example(): void;
declare namespace Example {
    var Foo: typeof Foo;
}
export default Example;
//// [index3.d.ts]
export declare class Bar {
}
declare function Example(): void;
declare namespace Example {
    var Bar: typeof Bar;
}
export default Example;
//// [index4.d.ts]
declare function A(): void;
declare function B(): void;
export declare function C(): any;
export declare namespace C {
    var A: typeof A;
    var B: typeof B;
}


//// [DtsFileErrors]


index1.d.ts(4,9): error TS2502: 'Foo' is referenced directly or indirectly in its own type annotation.
index2.d.ts(5,9): error TS2502: 'Foo' is referenced directly or indirectly in its own type annotation.
index3.d.ts(5,9): error TS2502: 'Bar' is referenced directly or indirectly in its own type annotation.
index4.d.ts(5,9): error TS2502: 'A' is referenced directly or indirectly in its own type annotation.
index4.d.ts(6,9): error TS2502: 'B' is referenced directly or indirectly in its own type annotation.


==== foo.d.ts (0 errors) ====
    export declare class Foo {
    }
    
==== index1.d.ts (1 errors) ====
    import { Foo } from './foo';
    declare function Example(): void;
    declare namespace Example {
        var Foo: typeof Foo;
            ~~~
!!! error TS2502: 'Foo' is referenced directly or indirectly in its own type annotation.
    }
    export default Example;
    
==== index2.d.ts (1 errors) ====
    import { Foo } from './foo';
    export { Foo };
    declare function Example(): void;
    declare namespace Example {
        var Foo: typeof Foo;
            ~~~
!!! error TS2502: 'Foo' is referenced directly or indirectly in its own type annotation.
    }
    export default Example;
    
==== index3.d.ts (1 errors) ====
    export declare class Bar {
    }
    declare function Example(): void;
    declare namespace Example {
        var Bar: typeof Bar;
            ~~~
!!! error TS2502: 'Bar' is referenced directly or indirectly in its own type annotation.
    }
    export default Example;
    
==== index4.d.ts (2 errors) ====
    declare function A(): void;
    declare function B(): void;
    export declare function C(): any;
    export declare namespace C {
        var A: typeof A;
            ~
!!! error TS2502: 'A' is referenced directly or indirectly in its own type annotation.
        var B: typeof B;
            ~
!!! error TS2502: 'B' is referenced directly or indirectly in its own type annotation.
    }
    