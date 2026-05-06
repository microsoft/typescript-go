//// [tests/cases/compiler/declarationEmitSubpathImportsReexport.ts] ////

//// [package.json]
{
  "name": "package-b",
  "type": "module",
  "exports": {
    ".": "./index.js"
  }
}

//// [index.js]
export {};

//// [index.d.ts]
export interface B {
	b: "b";
}

//// [package.json]
{
  "name": "package-a",
  "type": "module",
  "imports": {
    "#re_export": "./src/re_export.ts"
  },
  "exports": {
    ".": "./dist/index.js"
  }
}


//// [re_export.ts]
import type { B } from "package-b";
declare function foo(): Promise<B>
export const re = { foo };

//// [index.ts]
import { re } from "#re_export";
const { foo } = re;
export { foo };




//// [re_export.js]
export const re = { foo };
//// [index.js]
import { re } from "#re_export";
const { foo } = re;
export { foo };


//// [re_export.d.ts]
import type { B } from "package-b";
function foo(): Promise<B>;
export const re: {
    foo: typeof foo;
};
export {};
//// [index.d.ts]
const foo: () => Promise<import("package-b").B>;
export { foo };


//// [DtsFileErrors]


/packages/a/dist/index.d.ts(1,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
/packages/a/dist/re_export.d.ts(2,1): error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.


==== /packages/a/tsconfig.json (0 errors) ====
    {
      "compilerOptions": {
        "module": "nodenext",
        "outDir": "dist",
        "rootDir": "src",
        "declaration": true,
      },
      "include": ["src/**/*.ts"]
    }
    
==== /packages/a/dist/re_export.d.ts (1 errors) ====
    import type { B } from "package-b";
    function foo(): Promise<B>;
    ~~~~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export const re: {
        foo: typeof foo;
    };
    export {};
    
==== /packages/a/dist/index.d.ts (1 errors) ====
    const foo: () => Promise<import("package-b").B>;
    ~~~~~
!!! error TS1046: Top-level declarations in .d.ts files must start with either a 'declare' or 'export' modifier.
    export { foo };
    
==== /packages/b/package.json (0 errors) ====
    {
      "name": "package-b",
      "type": "module",
      "exports": {
        ".": "./index.js"
      }
    }
    
==== /packages/b/index.d.ts (0 errors) ====
    export interface B {
    	b: "b";
    }
    
==== /packages/a/package.json (0 errors) ====
    {
      "name": "package-a",
      "type": "module",
      "imports": {
        "#re_export": "./src/re_export.ts"
      },
      "exports": {
        ".": "./dist/index.js"
      }
    }
    
    