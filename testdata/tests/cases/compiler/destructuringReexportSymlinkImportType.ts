// Test case for destructuring re-exports using type from symlinked node-modules
// Expected: import() types should use package names instead of relative paths
// @declaration: true

// @Filename: /real-packages/package-b/package.json
{
  "name": "package-b",
  "main": "./index.js",
  "types": "./index.d.ts"
}

// @Filename: /real-packages/package-b/index.d.ts
export interface B {
  value: string;
}

// @Filename: /project/src/types.ts
import type { B } from "package-b";
export type { B };

// @Filename: /project/src/main.ts
import type { B } from "./types";

export function useB(param: B): B {
  return param;
}

// @link: /real-packages/package-b -> /project/node_modules/package-b