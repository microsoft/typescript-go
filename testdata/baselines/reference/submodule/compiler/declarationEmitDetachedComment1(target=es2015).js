//// [tests/cases/compiler/declarationEmitDetachedComment1.ts] ////

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
/**
 * Hello class
 */
class Hello {
}
//// [test2.js]
"use strict";
/* A comment at the top of the file. */
/**
 * Hi class
 */
class Hi {
}
//// [test3.js]
"use strict";
// A one-line comment at the top of the file.
/**
 * Hola class
 */
class Hola {
}


//// [test1.d.ts]
/*! Copyright 2015 MyCompany Inc. */
/**
 * Hello class
 */
class Hello {
}
//// [test2.d.ts]
/**
 * Hi class
 */
class Hi {
}
//// [test3.d.ts]
/**
 * Hola class
 */
class Hola {
}


//// [DtsFileErrors]


test1.d.ts(5,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
test2.d.ts(4,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
test3.d.ts(4,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== test1.d.ts (1 errors) ====
    /*! Copyright 2015 MyCompany Inc. */
    /**
     * Hello class
     */
    class Hello {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    }
    
==== test2.d.ts (1 errors) ====
    /**
     * Hi class
     */
    class Hi {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    }
    
==== test3.d.ts (1 errors) ====
    /**
     * Hola class
     */
    class Hola {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    }
    