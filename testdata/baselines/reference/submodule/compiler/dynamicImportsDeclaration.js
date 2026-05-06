//// [tests/cases/compiler/dynamicImportsDeclaration.ts] ////

//// [case0.ts]
export default 0;

//// [case1.ts]
export default 1;

//// [caseFallback.ts]
export default 'fallback';

//// [index.ts]
export const mod = await (async () => {
  const x: number = 0;
  switch (x) {
    case 0:
      return await import("./case0.js");
    case 1:
      return await import("./case1.js");
    default:
      return await import("./caseFallback.js");
  }
})();

//// [case0.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = 0;
//// [case1.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = 1;
//// [caseFallback.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.default = 'fallback';
//// [index.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.mod = void 0;
exports.mod = await (async () => {
    const x = 0;
    switch (x) {
        case 0:
            return await import("./case0.js");
        case 1:
            return await import("./case1.js");
        default:
            return await import("./caseFallback.js");
    }
})();


//// [case0.d.ts]
const _default = 0;
export default _default;
//// [case1.d.ts]
const _default = 1;
export default _default;
//// [caseFallback.d.ts]
const _default = "fallback";
export default _default;
//// [index.d.ts]
export const mod: {
    default: typeof import("./case0.js");
} | {
    default: typeof import("./case1.js");
} | {
    default: typeof import("./caseFallback.js");
};


//// [DtsFileErrors]


/case0.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
/case1.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
/caseFallback.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== /case0.d.ts (1 errors) ====
    const _default = 0;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export default _default;
    
==== /case1.d.ts (1 errors) ====
    const _default = 1;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export default _default;
    
==== /caseFallback.d.ts (1 errors) ====
    const _default = "fallback";
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export default _default;
    
==== /index.d.ts (0 errors) ====
    export const mod: {
        default: typeof import("./case0.js");
    } | {
        default: typeof import("./case1.js");
    } | {
        default: typeof import("./caseFallback.js");
    };
    