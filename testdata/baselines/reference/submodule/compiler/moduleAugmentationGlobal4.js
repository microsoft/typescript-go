//// [tests/cases/compiler/moduleAugmentationGlobal4.ts] ////

//// [f1.ts]
declare global {
    interface Something {x}
}
export {};
//// [f2.ts]
declare global {
    interface Something {y}
}
export {};
//// [f3.ts]
import "./f1";
import "./f2";



//// [f1.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [f2.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [f3.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
require("./f1");
require("./f2");


//// [f1.d.ts]
global {
    interface Something {
        x: any;
    }
}
export {};
//// [f2.d.ts]
global {
    interface Something {
        y: any;
    }
}
export {};
//// [f3.d.ts]
import "./f1";
import "./f2";


//// [DtsFileErrors]


f1.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
f2.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== f1.d.ts (1 errors) ====
    global {
    ~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        interface Something {
            x: any;
        }
    }
    export {};
    
==== f2.d.ts (1 errors) ====
    global {
    ~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        interface Something {
            y: any;
        }
    }
    export {};
    
==== f3.d.ts (0 errors) ====
    import "./f1";
    import "./f2";
    