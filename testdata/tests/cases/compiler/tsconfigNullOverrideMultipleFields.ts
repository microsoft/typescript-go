// @target: esnext
// @module: preserve
// @moduleResolution: bundler

// Tests that null values in tsconfig.json properly override extended values for different field types

// @filename: tsconfig-base.json

{
  "compilerOptions": {
    "types": ["node", "@types/jest"],
    "lib": ["es2020", "dom"],
    "typeRoots": ["./types", "./node_modules/@types"]
  }
}

// @filename: tsconfig.json

{
  "compilerOptions": {
    "types": null,
    "lib": null,
    "typeRoots": null
  },
  "extends": "./tsconfig-base.json"
}

// @filename: index.ts

export {};