//// [tests/cases/compiler/destructuringReexportSymlinkImportType.ts] ////

//// [package.json]
{
  "name": "package-b",
  "main": "./index.js",
  "types": "./index.d.ts"
}

//// [index.d.ts]
export interface B {
  value: string;
}

//// [types.ts]
import type { B } from "package-b";
export type { B };

//// [main.ts]
import type { B } from "./types";

export function useB(param: B): B {
  return param;
}


//// [types.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
//// [main.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.useB = useB;
function useB(param) {
    return param;
}


//// [types.d.ts]
import type { B } from "package-b";
export type { B };
//// [main.d.ts]
import type { B } from "./types";
export declare function useB(param: B): B;
