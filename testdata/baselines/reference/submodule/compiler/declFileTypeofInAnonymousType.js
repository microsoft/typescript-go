//// [tests/cases/compiler/declFileTypeofInAnonymousType.ts] ////

//// [declFileTypeofInAnonymousType.ts]
namespace m1 {
    export class c {
    }
    export enum e {
        weekday,
        weekend,
        holiday
    }
}
var a: { c: m1.c; };
var b = {
    c: m1.c,
    m1: m1
};
var c = { m1: m1 };
var d = {
    m: { mod: m1 },
    mc: { cl: m1.c },
    me: { en: m1.e },
    mh: m1.e.holiday
};

//// [declFileTypeofInAnonymousType.js]
"use strict";
var m1;
(function (m1) {
    class c {
    }
    m1.c = c;
    let e;
    (function (e) {
        e[e["weekday"] = 0] = "weekday";
        e[e["weekend"] = 1] = "weekend";
        e[e["holiday"] = 2] = "holiday";
    })(e = m1.e || (m1.e = {}));
})(m1 || (m1 = {}));
var a;
var b = {
    c: m1.c,
    m1: m1
};
var c = { m1: m1 };
var d = {
    m: { mod: m1 },
    mc: { cl: m1.c },
    me: { en: m1.e },
    mh: m1.e.holiday
};


//// [declFileTypeofInAnonymousType.d.ts]
namespace m1 {
    class c {
    }
    enum e {
        weekday = 0,
        weekend = 1,
        holiday = 2
    }
}
var a: {
    c: m1.c;
};
var b: {
    c: typeof m1.c;
    m1: typeof m1;
};
var c: {
    m1: typeof m1;
};
var d: {
    m: {
        mod: typeof m1;
    };
    mc: {
        cl: typeof m1.c;
    };
    me: {
        en: typeof m1.e;
    };
    mh: m1.e;
};


//// [DtsFileErrors]


declFileTypeofInAnonymousType.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== declFileTypeofInAnonymousType.d.ts (1 errors) ====
    namespace m1 {
    ~~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        class c {
        }
        enum e {
            weekday = 0,
            weekend = 1,
            holiday = 2
        }
    }
    var a: {
        c: m1.c;
    };
    var b: {
        c: typeof m1.c;
        m1: typeof m1;
    };
    var c: {
        m1: typeof m1;
    };
    var d: {
        m: {
            mod: typeof m1;
        };
        mc: {
            cl: typeof m1.c;
        };
        me: {
            en: typeof m1.e;
        };
        mh: m1.e;
    };
    