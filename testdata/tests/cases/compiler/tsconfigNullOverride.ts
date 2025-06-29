// @target: esnext
// @module: preserve
// @moduleResolution: bundler

// Tests that null values in tsconfig.json properly override extended values

// @filename: tsconfig-base.json

{
  "compilerOptions": {
    "types": []
  }
}

// @filename: tsconfig.json

{
  "compilerOptions": {
    "types": null
  },
  "extends": "./tsconfig-base.json"
}

// @filename: index.ts

export {};