// @strict: true
// @declaration: true
// @emitDeclarationOnly: true
// @skipLibCheck: true
// @moduleResolution: node16
// @module: node16

// Test that skipLibCheck suppresses TS2742 errors for internal library types
// when generating declarations in pnpm workspaces where package exports block
// direct access to internal files.

// @Filename: /node_modules/.pnpm/@base-lib+react@1.0.0/node_modules/@base-lib/react/package.json
{
  "name": "@base-lib/react",
  "version": "1.0.0",
  "exports": {
    ".": "./index.d.ts"
  }
}

// @Filename: /node_modules/.pnpm/@base-lib+react@1.0.0/node_modules/@base-lib/react/index.d.ts
export { Tooltip } from "./esm/Tooltip";

// @Filename: /node_modules/.pnpm/@base-lib+react@1.0.0/node_modules/@base-lib/react/esm/Tooltip.d.ts
import { InternalUtil } from "./utils/internal";
export interface TooltipProps { content: string; }
export declare const Tooltip: TooltipProps & { __internal: InternalUtil };

// @Filename: /node_modules/.pnpm/@base-lib+react@1.0.0/node_modules/@base-lib/react/esm/utils/internal.d.ts
// Internal utility type - should not be referenceable due to exports field
export interface InternalUtil { __private: never; }

// @Filename: /node_modules/@base-lib/react/package.json
{
  "name": "@base-lib/react",
  "version": "1.0.0"
}

// @Filename: /node_modules/@base-lib/react/index.d.ts
export * from "../../.pnpm/@base-lib+react@1.0.0/node_modules/@base-lib/react/index.d.ts";

// @Filename: /src/component.ts
import { Tooltip } from "@base-lib/react";

// This function returns a value whose inferred type includes Tooltip with its internal types.
// Without syntacticNodeBuilder, tsgo will:
// 1. Analyze the function body to infer the return type
// 2. Encounter Tooltip type which has __internal: InternalUtil
// 3. Try to serialize InternalUtil type
// 4. Fail to generate clean specifier (blocked by exports)
// 5. Fall back to relative path through .pnpm directory
// 6. With skipLibCheck, should suppress error instead of reporting TS2742
export function createComponent() {
  return {
    data: Tooltip
  };
}

