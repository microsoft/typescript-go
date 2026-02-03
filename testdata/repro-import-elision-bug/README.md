# Minimal Reproduction for Import Elision Bug

This directory contains a minimal test case that reproduces the bug where `tsgo` 
incorrectly emits a `require()` statement for type-only imports.

## Bug Description
When an import statement contains:
1. A `type` import specifier (e.g., `type ValueData`)
2. A regular import specifier that is ONLY used in type positions (e.g., `Value`)

tsgo incorrectly emits `require()` for the module, even though no runtime import is needed.

## Test Files
- `provider.ts` - Exports both a class (`Value`) and a type (`ValueData`)
- `broken.ts` - Imports both with and without `type` keyword - **BUG: emits require()**
- `working.ts` - Imports only `Value` without `type` keyword - **Works correctly**

## How to Reproduce

```bash
cd testdata/repro-import-elision-bug

# Compile the broken case
tsgo --module commonjs --target es2020 broken.ts
cat broken.js
# BUG: Contains `require("./provider");`

# Compile the working case  
tsgo --module commonjs --target es2020 working.ts
cat working.js
# CORRECT: No require() statement

# Compare with TypeScript compiler
tsc --module commonjs --target es2020 broken.ts
cat broken.js
# CORRECT: TypeScript does NOT emit require()
```

## Expected Behavior
Neither `broken.ts` nor `working.ts` should emit `require("./provider")` because 
`Value` is only used in type annotations, not at runtime.

## Actual Behavior
- `working.ts`: Correctly does NOT emit `require()` ✓
- `broken.ts`: Incorrectly DOES emit `require()` ✗

The presence of the `type ValueData` import causes the bug to manifest.
