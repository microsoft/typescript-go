//// [tests/cases/compiler/declarationEmitDetachedComment2.ts] ////

//// [test1.ts]
/*! Copyright 2015 MyCompany Inc. */

/**
 * Hello class
 */
class Hello {

}

//// [test2.ts]
/* A comment at the top of the file. */

/**
 * Hi class
 */
class Hi {

}

//// [test3.ts]
// A one-line comment at the top of the file.

/**
 * Hola class
 */
class Hola {

}


//// [test1.js]
"use strict";
/*! Copyright 2015 MyCompany Inc. */
class Hello {
}
//// [test2.js]
"use strict";
class Hi {
}
//// [test3.js]
"use strict";
class Hola {
}


//// [test1.d.ts]
/*! Copyright 2015 MyCompany Inc. */
class Hello {
}
//// [test2.d.ts]
class Hi {
}
//// [test3.d.ts]
class Hola {
}


//// [DtsFileErrors]


test1.d.ts(2,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
test2.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
test3.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== test1.d.ts (1 errors) ====
    /*! Copyright 2015 MyCompany Inc. */
    class Hello {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    }
    
==== test2.d.ts (1 errors) ====
    class Hi {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    }
    
==== test3.d.ts (1 errors) ====
    class Hola {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    }
    