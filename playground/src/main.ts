import * as monaco from "monaco-editor";
import { createPlaygroundLspClient } from "./lspClient";

// ---------------------------------------------------------------------------
// Disable Monaco's built-in TypeScript language features entirely.
// MonacoLspClient provides all features via the tsgo LSP server.
// ---------------------------------------------------------------------------
const disabledModeConfig: Record<keyof monaco.typescript.ModeConfiguration, false> = {
    completionItems: false,
    hovers: false,
    documentSymbols: false,
    definitions: false,
    references: false,
    documentHighlights: false,
    rename: false,
    diagnostics: false,
    documentRangeFormattingEdits: false,
    signatureHelp: false,
    onTypeFormattingEdits: false,
    codeActions: false,
    inlayHints: false,
};

monaco.typescript.typescriptDefaults.setModeConfiguration(disabledModeConfig);
monaco.typescript.javascriptDefaults.setModeConfiguration(disabledModeConfig);

// ---------------------------------------------------------------------------
// Monaco environment - only the base editor worker is needed.
// MonacoLspClient handles all TypeScript features via LSP.
// ---------------------------------------------------------------------------
self.MonacoEnvironment = {
    getWorker() {
        return new Worker(
            new URL("monaco-editor/esm/vs/editor/editor.worker.js", import.meta.url),
            { type: "module" },
        );
    },
};

// ---------------------------------------------------------------------------
// Constants
// ---------------------------------------------------------------------------
const INITIAL_CODE = `// Welcome to the typescript-go Playground!
// Powered by typescript-go compiled to WebAssembly.

interface Person {
  name: string;
  age: number;
}

function greet(person: Person): string {
  return \`Hello, \${person.name}! You are \${person.age} years old.\`;
}

const alice: Person = { name: "Alice", age: 30 };
console.log(greet(alice));

// Try introducing an error - uncomment the line below:
// const bob: Person = { name: "Bob" };
`;

// ---------------------------------------------------------------------------
// UI references
// ---------------------------------------------------------------------------
const statusEl = document.getElementById("status")!;

function setStatus(text: string, kind: "loading" | "ready" | "error" = "loading") {
    statusEl.textContent = text;
    statusEl.className = kind;
}

// ---------------------------------------------------------------------------
// Create a Monaco editor instance with a `file:///` URI.
// Otherwise, Monaco sends over something like `inmemory://model/1` which the LSP server cannot find.
// ---------------------------------------------------------------------------
const model = monaco.editor.createModel(
    INITIAL_CODE,
    "typescript",
    monaco.Uri.parse("file:///playground/index.ts"),
);

const editor = monaco.editor.create(document.getElementById("editor")!, {
    model: model,
    theme: "vs-dark",
    automaticLayout: true,
    minimap: { enabled: false },
    fontSize: 14,
    padding: { top: 16 },
});

// ---------------------------------------------------------------------------
// Start LSP via MonacoLspClient
// ---------------------------------------------------------------------------
async function initLsp() {
    setStatus("Loading WASM + initializing LSP...");

    const worker = new Worker("/lspWorker.js");
    // MonacoLspClient sends initialize automatically and registers
    // all language features (diagnostics, hover, completions, etc.)
    // Messages are buffered in the worker until Go WASM finishes loading.
    createPlaygroundLspClient(worker);

    setStatus("Ready", "ready");
}

await initLsp().catch(err => {
    setStatus(`Error: ${err.message || err}`, "error");
    console.error("LSP init error:", err);
});
