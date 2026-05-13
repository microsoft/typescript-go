import { lsp } from "monaco-editor";

/**
 * Creates a MonacoLspClient connected to the Go WASM LSP worker.
 *
 * Uses Monaco's built-in createTransportToWorker which sends/receives
 * raw JSON-RPC Message objects via postMessage. The worker buffers
 * messages in its stdin queue, so even if MonacoLspClient sends
 * `initialize` before the Go runtime starts, it gets processed
 * once WASM finishes loading.
 */
export function createPlaygroundLspClient(worker: Worker) {
    const transport = lsp.createTransportToWorker(worker);
    const client = new lsp.MonacoLspClient(transport);
    return { client, transport };
}
