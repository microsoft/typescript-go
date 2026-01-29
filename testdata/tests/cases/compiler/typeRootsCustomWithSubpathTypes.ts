// @strict: true
// @noImplicitReferences: true
// @module: nodenext
// @moduleResolution: nodenext
// @types: testpkg/globals
// @typeRoots: /node_modules
// @traceResolution: true
// @noEmit: true

// @filename: /node_modules/testpkg/package.json
{
    "name": "testpkg",
    "version": "1.0.0",
    "exports": {
        ".": {
            "import": {
                "types": "./index.d.ts"
            }
        },
        "./globals": {
            "types": "./globals.d.ts"
        }
    }
}

// @filename: /node_modules/testpkg/index.d.ts
export interface TestPkg {
    name: string;
}

// @filename: /node_modules/testpkg/globals.d.ts
// This file provides global declarations, similar to vitest/globals
declare global {
    var testGlobal: string;
    function testFunc(): void;
}
export {}

// @filename: /index.ts
// When typeRoots is set to a custom directory (not ending in @types),
// type reference directives with subpaths like "testpkg/globals" should
// still resolve correctly by checking if the typeRoot directory exists
// (not the candidate path), then loading from the candidate.

// If resolution works correctly, testGlobal and testFunc should be available
const x: string = testGlobal;
testFunc();
