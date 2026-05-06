//// [tests/cases/compiler/commentsClass.ts] ////

//// [commentsClass.ts]
/** This is class c2 without constuctor*/
class c2 {
} // trailing comment1
var i2 = new c2();
var i2_c = c2;
class c3 {
    /** Constructor comment*/
    constructor() {
    } // trailing comment of constructor
} /* trailing comment 2 */
var i3 = new c3();
var i3_c = c3;
/** Class comment*/
class c4 {
    /** Constructor comment*/
    constructor() {
    } /* trailing comment of constructor 2*/
}
var i4 = new c4();
var i4_c = c4;
/** Class with statics*/
class c5 {
    static s1: number;
}
var i5 = new c5();
var i5_c = c5;

/// class with statics and constructor
class c6 { /// class with statics and constructor2
    /// s1 comment
    static s1: number; /// s1 comment2
    /// constructor comment
    constructor() { /// constructor comment2
    }
}
var i6 = new c6();
var i6_c = c6;

// class with statics and constructor
class c7 {
    // s1 comment
    static s1: number;
    // constructor comment
    constructor() {
    }
}
var i7 = new c7();
var i7_c = c7;

/** class with statics and constructor
 */
class c8 {
    /** s1 comment */
    static s1: number; /** s1 comment2 */
    /** constructor comment 
    */
    constructor() {
        /** constructor comment2 
        */
    }
}
var i8 = new c8();
var i8_c = c8;

class c9 {
    constructor() {
        /// This is some detached comment

        // should emit this leading comment of } too
    }
}


//// [commentsClass.js]
"use strict";
/** This is class c2 without constuctor*/
class c2 {
} // trailing comment1
var i2 = new c2();
var i2_c = c2;
class c3 {
    /** Constructor comment*/
    constructor() {
    } // trailing comment of constructor
} /* trailing comment 2 */
var i3 = new c3();
var i3_c = c3;
/** Class comment*/
class c4 {
    /** Constructor comment*/
    constructor() {
    } /* trailing comment of constructor 2*/
}
var i4 = new c4();
var i4_c = c4;
/** Class with statics*/
class c5 {
}
var i5 = new c5();
var i5_c = c5;
/// class with statics and constructor
class c6 {
    /// constructor comment
    constructor() {
    }
}
var i6 = new c6();
var i6_c = c6;
// class with statics and constructor
class c7 {
    // constructor comment
    constructor() {
    }
}
var i7 = new c7();
var i7_c = c7;
/** class with statics and constructor
 */
class c8 {
    /** constructor comment
    */
    constructor() {
        /** constructor comment2
        */
    }
}
var i8 = new c8();
var i8_c = c8;
class c9 {
    constructor() {
        /// This is some detached comment
        // should emit this leading comment of } too
    }
}


//// [commentsClass.d.ts]
/** This is class c2 without constuctor*/
class c2 {
}
var i2: c2;
var i2_c: typeof c2;
class c3 {
    /** Constructor comment*/
    constructor();
}
var i3: c3;
var i3_c: typeof c3;
/** Class comment*/
class c4 {
    /** Constructor comment*/
    constructor();
}
var i4: c4;
var i4_c: typeof c4;
/** Class with statics*/
class c5 {
    static s1: number;
}
var i5: c5;
var i5_c: typeof c5;
class c6 {
    static s1: number;
    constructor();
}
var i6: c6;
var i6_c: typeof c6;
class c7 {
    static s1: number;
    constructor();
}
var i7: c7;
var i7_c: typeof c7;
/** class with statics and constructor
 */
class c8 {
    /** s1 comment */
    static s1: number; /** s1 comment2 */
    /** constructor comment
    */
    constructor();
}
var i8: c8;
var i8_c: typeof c8;
class c9 {
    constructor();
}


//// [DtsFileErrors]


commentsClass.d.ts(2,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== commentsClass.d.ts (1 errors) ====
    /** This is class c2 without constuctor*/
    class c2 {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    }
    var i2: c2;
    var i2_c: typeof c2;
    class c3 {
        /** Constructor comment*/
        constructor();
    }
    var i3: c3;
    var i3_c: typeof c3;
    /** Class comment*/
    class c4 {
        /** Constructor comment*/
        constructor();
    }
    var i4: c4;
    var i4_c: typeof c4;
    /** Class with statics*/
    class c5 {
        static s1: number;
    }
    var i5: c5;
    var i5_c: typeof c5;
    class c6 {
        static s1: number;
        constructor();
    }
    var i6: c6;
    var i6_c: typeof c6;
    class c7 {
        static s1: number;
        constructor();
    }
    var i7: c7;
    var i7_c: typeof c7;
    /** class with statics and constructor
     */
    class c8 {
        /** s1 comment */
        static s1: number; /** s1 comment2 */
        /** constructor comment
        */
        constructor();
    }
    var i8: c8;
    var i8_c: typeof c8;
    class c9 {
        constructor();
    }
    