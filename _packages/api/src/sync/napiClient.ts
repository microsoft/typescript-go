/**
 * NAPI-based client for the TypeScript Go API.
 *
 * This client loads the TypeScript compiler directly as a native Node.js
 * addon (NAPI module) instead of spawning a child process. All communication
 * happens in-process via direct function calls, avoiding IPC overhead.
 *
 * The NAPI client provides the same interface as the sync IPC client,
 * so it can be used as a drop-in replacement.
 */

import { createRequire } from "node:module";
import path from "node:path";
import { fileURLToPath } from "node:url";

interface NapiModule {
    createSession(cwd: string, defaultLibraryPath?: string, callbacks?: string): void;
    request(method: string, payload: string): string;
    requestBinary(method: string, payload: Uint8Array): Uint8Array;
    close(): void;
}

/**
 * Resolve the path to the tsgo.node native addon.
 *
 * Checks in order:
 *   1. Explicit path from options.napiModulePath
 *   2. Repository's built/local/tsgo.node
 *   3. Platform-specific npm package location
 */
function resolveNapiModulePath(explicitPath?: string): string {
    if (explicitPath) return explicitPath;

    const __dirname = path.dirname(fileURLToPath(import.meta.url));
    const normalizedDirname = __dirname.replace(/\\/g, "/");

    // Check if we're in the repository
    if (normalizedDirname.endsWith("/_packages/api/src/sync") || normalizedDirname.endsWith("/_packages/api/dist/sync")) {
        const repoRoot = path.resolve(__dirname, "..", "..", "..", "..");
        const repoPath = path.join(repoRoot, "built", "local", "tsgo.node");
        return repoPath;
    }

    // Fallback: look in the same directory as the module
    return path.join(__dirname, "..", "tsgo.node");
}

export interface NapiClientOptions {
    /** Current working directory */
    cwd?: string;
    /** Explicit path to the tsgo.node native addon */
    napiModulePath?: string;
    /**
     * Path to the directory containing bundled lib.d.ts files.
     * Required for noembed builds. If not provided, defaults to the
     * directory containing the .node addon.
     */
    defaultLibraryPath?: string;
}

export class NapiClient {
    private module: NapiModule;
    private encoder = new TextEncoder();

    constructor(options: NapiClientOptions = {}) {
        const cwd = options.cwd ?? process.cwd();
        const modulePath = path.resolve(resolveNapiModulePath(options.napiModulePath));
        const defaultLibraryPath = options.defaultLibraryPath ?? path.dirname(modulePath);

        // Use createRequire to load native .node modules (ESM can't load them directly)
        const require = createRequire(import.meta.url);
        this.module = require(modulePath) as NapiModule;

        this.module.createSession(cwd, defaultLibraryPath);
    }

    apiRequest<T>(method: string, params?: unknown): T {
        const encodedPayload = JSON.stringify(params);
        const result = this.module.request(method, encodedPayload);
        if (result.length) {
            return JSON.parse(result) as T;
        }
        return undefined as unknown as T;
    }

    apiRequestBinary(method: string, params?: unknown): Uint8Array | undefined {
        const result = this.module.requestBinary(method, this.encoder.encode(JSON.stringify(params)));
        if (result.length === 0) return undefined;
        return result;
    }

    echo(payload: string): string {
        return this.module.request("echo", payload);
    }

    echoBinary(payload: Uint8Array): Uint8Array {
        return this.module.requestBinary("echo", payload);
    }

    close(): void {
        this.module.close();
    }
}
