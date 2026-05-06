//// [tests/cases/compiler/internalAliasEnum.ts] ////

//// [internalAliasEnum.ts]
namespace a {
    export enum weekend {
        Friday,
        Saturday,
        Sunday
    }
}

namespace c {
    import b = a.weekend;
    export var bVal: b = b.Sunday;
}


//// [internalAliasEnum.js]
"use strict";
var a;
(function (a) {
    let weekend;
    (function (weekend) {
        weekend[weekend["Friday"] = 0] = "Friday";
        weekend[weekend["Saturday"] = 1] = "Saturday";
        weekend[weekend["Sunday"] = 2] = "Sunday";
    })(weekend = a.weekend || (a.weekend = {}));
})(a || (a = {}));
var c;
(function (c) {
    var b = a.weekend;
    c.bVal = b.Sunday;
})(c || (c = {}));


//// [internalAliasEnum.d.ts]
namespace a {
    enum weekend {
        Friday = 0,
        Saturday = 1,
        Sunday = 2
    }
}
namespace c {
    import b = a.weekend;
    var bVal: b;
}


//// [DtsFileErrors]


internalAliasEnum.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== internalAliasEnum.d.ts (1 errors) ====
    namespace a {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        enum weekend {
            Friday = 0,
            Saturday = 1,
            Sunday = 2
        }
    }
    namespace c {
        import b = a.weekend;
        var bVal: b;
    }
    