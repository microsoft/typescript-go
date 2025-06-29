//// [tests/cases/compiler/tsconfigNullOverrideMultipleFields.ts] ////

//// [tsconfig-base.json]
{
  "compilerOptions": {
    "types": ["node", "@types/jest"],
    "lib": ["es2020", "dom"],
    "typeRoots": ["./types", "./node_modules/@types"]
  }
}

//// [index.ts]
export {};

//// [index.js]
