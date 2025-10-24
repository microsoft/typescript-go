//// [tests/cases/compiler/declarationEmitPnpmWorkspaceSkipLibCheck.ts] ////

//// [package.json]
{
  "name": "@base-lib/react",
  "version": "1.0.0",
  "exports": {
    ".": "./index.d.ts"
  }
}

//// [index.d.ts]
export { Tooltip } from "./esm/Tooltip";

//// [Tooltip.d.ts]
import { InternalUtil } from "./utils/internal";
export interface TooltipProps { content: string; }
export declare const Tooltip: TooltipProps & { __internal: InternalUtil };

//// [internal.d.ts]
// Internal utility type - should not be referenceable due to exports field
export interface InternalUtil { __private: never; }

//// [package.json]
{
  "name": "@base-lib/react",
  "version": "1.0.0"
}

//// [index.d.ts]
export * from "../../.pnpm/@base-lib+react@1.0.0/node_modules/@base-lib/react/index.d.ts";

//// [component.ts]
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





//// [component.d.ts]
export declare function createComponent(): {
    data: {
        __internal;
    };
};
