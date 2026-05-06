//// [tests/cases/compiler/jsdocCommentDefaultExport.ts] ////

//// [exportDefaultObject.ts]
/** Object comment */
export default {
    fn() {}
}

//// [exportDefaultFunction.ts]
/** Function comment */
export default function() {
    return 42;
}

//// [exportDefaultClass.ts]
/** Class comment */
export default class {
    method() {}
}

//// [exportDefaultLiteral.ts]
/** Literal comment */
export default 42;

//// [exportDefaultNull.ts]
/** Null comment */
export default null;


//// [exportDefaultObject.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
/** Object comment */
exports.default = {
    fn() { }
};
//// [exportDefaultFunction.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = default_1;
/** Function comment */
function default_1() {
    return 42;
}
//// [exportDefaultClass.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
/** Class comment */
class default_1 {
    method() { }
}
exports.default = default_1;
//// [exportDefaultLiteral.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
/** Literal comment */
exports.default = 42;
//// [exportDefaultNull.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
/** Null comment */
exports.default = null;


//// [exportDefaultObject.d.ts]
/** Object comment */
const _default: {
    fn(): void;
};
export default _default;
//// [exportDefaultFunction.d.ts]
/** Function comment */
export default function (): number;
//// [exportDefaultClass.d.ts]
/** Class comment */
export default class {
    method(): void;
}
//// [exportDefaultLiteral.d.ts]
/** Literal comment */
const _default = 42;
export default _default;
//// [exportDefaultNull.d.ts]
/** Null comment */
const _default: null;
export default _default;


//// [DtsFileErrors]


exportDefaultLiteral.d.ts(2,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
exportDefaultNull.d.ts(2,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
exportDefaultObject.d.ts(2,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== exportDefaultObject.d.ts (1 errors) ====
    /** Object comment */
    const _default: {
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
        fn(): void;
    };
    export default _default;
    
==== exportDefaultFunction.d.ts (0 errors) ====
    /** Function comment */
    export default function (): number;
    
==== exportDefaultClass.d.ts (0 errors) ====
    /** Class comment */
    export default class {
        method(): void;
    }
    
==== exportDefaultLiteral.d.ts (1 errors) ====
    /** Literal comment */
    const _default = 42;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export default _default;
    
==== exportDefaultNull.d.ts (1 errors) ====
    /** Null comment */
    const _default: null;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export default _default;
    