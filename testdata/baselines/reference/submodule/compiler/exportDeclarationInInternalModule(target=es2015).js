//// [tests/cases/compiler/exportDeclarationInInternalModule.ts] ////

//// [exportDeclarationInInternalModule.ts]
class Bbb {
}

class Aaa extends Bbb { }

namespace Aaa {
    export class SomeType { }
}

namespace Bbb {
    export class SomeType { }

    export * from Aaa;      // this line causes the nullref
}

var a: Bbb.SomeType;


//// [exportDeclarationInInternalModule.js]
"use strict";
class Bbb {
}
class Aaa extends Bbb {
}
(function (Aaa) {
    class SomeType {
    }
    Aaa.SomeType = SomeType;
})(Aaa || (Aaa = {}));
(function (Bbb) {
    class SomeType {
    }
    Bbb.SomeType = SomeType;
})(Bbb || (Bbb = {}));
var a;


//// [exportDeclarationInInternalModule.d.ts]
class Bbb {
}
class Aaa extends Bbb {
}
namespace Aaa {
    class SomeType {
    }
}
namespace Bbb {
    export class SomeType {
    }
    export * from Aaa;
}
var a: Bbb.SomeType;
