# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

TypeScript 7 is a native port of the TypeScript compiler and language server written in Go. This is a work-in-progress project that aims to achieve feature parity with TypeScript 5.8, with most development happening in the `internal` directory.

## Build System and Commands

This project uses `hereby` as the primary build tool. The following npm scripts are available:

```bash
# Core build and test commands
npm run build               # Build the tsgo binary using hereby
npm run build:watch         # Build with watch mode
npm run build:watch:debug   # Build with debug output in watch mode
npm run test                # Run all tests

# Package-specific commands
npm run api:build           # Build the API package
npm run extension:build     # Build the VS Code extension
npm run extension:watch     # Build extension in watch mode

# Development utilities
npm run node                # Run node with special conditions
npm run convertfourslash    # Convert fourslash test files
npm run updatefailing       # Update failing test baselines
npm run makemanual          # Create manual test files
```

Additional hereby commands (run with `npx hereby <command>`):
- `npx hereby format` - Format the code using dprint
- `npx hereby lint` - Run Go linters
- `npx hereby baseline-accept` - Update test baselines/snapshots

### Build Options

The build system supports several flags for debugging and performance analysis:
- `--debug` - Include debug symbols and disable optimizations
- `--race` - Enable Go race detector (use `TSGO_HEREBY_RACE=true`)
- Build tags: `release` (default), `noembed` for debugging

## Testing

### Running Specific Tests

```bash
# For pre-existing "submodule" tests from _submodules/TypeScript
go test -run='TestSubmodule/<test name>' ./internal/testrunner

# For new "local" tests in testdata/tests/cases
go test -run='TestLocal/<test name>' ./internal/testrunner

# Run tests with race detection (for debugging concurrency issues)
go test -race ./internal/testrunner

# Run tests for a specific package (e.g., parser, checker, etc.)
go test ./internal/parser
go test ./internal/checker
```

### Writing New Compiler Tests

New compiler tests go in `testdata/tests/cases/compiler/` and use TypeScript syntax with special comments:

```ts
// @target: esnext
// @module: preserve
// @moduleResolution: bundler
// @strict: true
// @checkJs: true

// @filename: fileA.ts
export interface Person {
    name: string;
    age: number;
}

// @filename: fileB.js
/** @import { Person } from "./fileA" */
function greet(person) {
    console.log(`Hello, ${person.name}!`);
}
```

**Always enable strict mode (`@strict: true`) unless testing non-strict behavior.**

Test outputs are generated in `testdata/baselines/local` and compared against `testdata/baselines/reference`. Use `npx hereby baseline-accept` to update baselines after test changes.

## Architecture

### Key Directories

- **`internal/`** - Main compiler and language server code (Go)
  - `compiler/` - Core compilation logic
  - `parser/` - TypeScript parsing
  - `checker/` - Type checking
  - `binder/` - Symbol binding
  - `scanner/` - Lexical analysis
  - `lsp/` - Language Server Protocol implementation
  - `testrunner/` - Test execution framework
- **`_extension/`** - VS Code extension (TypeScript/JavaScript)
- **`_packages/`** - NPM packages
  - `native-preview/` - Preview npm package
  - `api/` - TypeScript API bindings
- **`_submodules/TypeScript`** - Reference TypeScript implementation
- **`testdata/`** - Test cases and baselines

### Development Workflow

1. Write minimal test cases demonstrating the bug/feature
2. Run tests to verify failure (bugs) or expected behavior (features)
3. Accept generated baselines if needed
4. Implement the fix/feature in `internal/`
5. Re-run tests and accept new baselines
6. Ensure code is formatted (`npx hereby format`) and linted (`npx hereby lint`)

### Code Reference

The `internal/` directory is ported from `_submodules/TypeScript`. When implementing features or fixing bugs, search the TypeScript submodule for reference implementations.

## Language Server Integration

The LSP implementation is in `internal/lsp/`. Editor functionality requires integration testing with the language server, not just compiler tests.

## Go Development Notes

### Module Structure
- Uses Go 1.25+ with workspace mode
- Main module: `github.com/microsoft/typescript-go`
- Key dependencies: `regexp2`, `go-json-experiment/json`, `xxh3` for hashing
- Build tools included as Go tools: `moq`, `stringer`, `gofumpt`

### Performance Considerations
- UTF-8 string handling differs from TypeScript's UTF-16 approach
- Node positions use UTF-8 offsets, not UTF-16 (affects non-ASCII character positions)
- xxh3 hashing used for performance-critical operations
- Memory pooling and string interning used in hot paths

### Debugging
- Use `--debug` flag for builds to include debug symbols
- Race detector available with `--race` or `TSGO_HEREBY_RACE=true`
- LSP debugging can be enabled through VS Code extension settings

## Intentional Changes

See `CHANGES.md` for documented differences between this Go port and the original TypeScript compiler, including scanner, parser, and JSDoc handling changes.