# LSP Server Communication Example

This is a minimal VS Code extension demonstrating how to use the `@typescript/api` package to communicate with the TypeScript Go language server.

## Prerequisites

- The TypeScript Native Preview extension must be installed and active in VS Code
- A TypeScript project must be open

## How It Works

1. The extension registers a command "Print AST of Current File"
2. When executed, it:
   - Gets an API connection pipe path from the TypeScript Native Preview extension via `typescript.native-preview.initializeAPIConnection`
   - Connects using `AsyncAPI.fromLSPConnection()` from `@typescript/api`
   - Gets the project for the current file with `api.getDefaultProjectForFile()`
   - Gets the AST with `project.getSourceFile()`
   - Walks the tree using `node.forEachChild()` and prints to an output channel

## Development

```bash
# Install dependencies
npm install

# Compile (uses esbuild to bundle)
npm run compile

# Watch mode
npm run watch
```

## Testing

1. Use the "Launch VS Code extension + API example" launch configuration from the repo root
2. In the Extension Development Host, ensure the TypeScript Native Preview extension activates
3. Open a TypeScript file
4. Run "Print AST of Current File" from the Command Palette (Cmd+Shift+P)
5. Check the "TypeScript AST" output channel for results

## Key Code

```typescript
import { AsyncAPI } from "@typescript/api/async";
import { SyntaxKind } from "@typescript/ast";

// Get pipe path from TypeScript Native Preview extension
const pipePath = await vscode.commands.executeCommand<string>(
    "typescript.native-preview.initializeAPIConnection"
);

// Connect to the API
const api = AsyncAPI.fromLSPConnection({ pipePath });

// Get project and source file
const project = await api.getDefaultProjectForFile(fileName);
const sourceFile = await project.getSourceFile(fileName);

// Walk the AST
sourceFile.forEachChild(function visit(node) {
    console.log(SyntaxKind[node.kind]);
    node.forEachChild(visit);
});

await api.close();
```
