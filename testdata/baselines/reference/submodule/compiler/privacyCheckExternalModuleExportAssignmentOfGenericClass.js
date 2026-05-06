//// [tests/cases/compiler/privacyCheckExternalModuleExportAssignmentOfGenericClass.ts] ////

//// [privacyCheckExternalModuleExportAssignmentOfGenericClass_0.ts]
export = Foo;
class Foo<A> {
    constructor(public a: A) { }
}

//// [privacyCheckExternalModuleExportAssignmentOfGenericClass_1.ts]
import Foo = require("./privacyCheckExternalModuleExportAssignmentOfGenericClass_0");
export = Bar;
interface Bar {
    foo: Foo<number>;
}

//// [privacyCheckExternalModuleExportAssignmentOfGenericClass_0.js]
"use strict";
class Foo {
    constructor(a) {
        this.a = a;
    }
}
module.exports = Foo;
//// [privacyCheckExternalModuleExportAssignmentOfGenericClass_1.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });


//// [privacyCheckExternalModuleExportAssignmentOfGenericClass_0.d.ts]
export = Foo;
class Foo<A> {
    a: A;
    constructor(a: A);
}
//// [privacyCheckExternalModuleExportAssignmentOfGenericClass_1.d.ts]
import Foo = require("./privacyCheckExternalModuleExportAssignmentOfGenericClass_0");
export = Bar;
interface Bar {
    foo: Foo<number>;
}


//// [DtsFileErrors]


privacyCheckExternalModuleExportAssignmentOfGenericClass_0.d.ts(2,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== privacyCheckExternalModuleExportAssignmentOfGenericClass_1.d.ts (0 errors) ====
    import Foo = require("./privacyCheckExternalModuleExportAssignmentOfGenericClass_0");
    export = Bar;
    interface Bar {
        foo: Foo<number>;
    }
    
==== privacyCheckExternalModuleExportAssignmentOfGenericClass_0.d.ts (1 errors) ====
    export = Foo;
    class Foo<A> {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        a: A;
        constructor(a: A);
    }
    