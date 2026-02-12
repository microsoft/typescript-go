/**
 * Shared utilities for the TypeScript API client.
 */

import type { FileSystem } from "./fs.ts";
import type { SyncTransport } from "./syncChannel.ts";

export type { SyncTransport };

export interface ClientSocketOptions {
    /** Path to the Unix domain socket or Windows named pipe for API communication */
    pipe: string;
}

export interface ClientSpawnOptions {
    /** Path to the tsgo executable */
    tsserverPath: string;
    /** Current working directory */
    cwd?: string;
    /** Virtual filesystem callbacks */
    fs?: FileSystem;
    /** Transport mechanism: "stdio" (default on Unix), "pipe" (default on Windows), "fifo" */
    transport?: SyncTransport;
}

export type ClientOptions = ClientSocketOptions | ClientSpawnOptions;

export function isSpawnOptions(options: ClientOptions): options is ClientSpawnOptions {
    return "tsserverPath" in options;
}

export interface LSPConnectionOptions extends ClientSocketOptions {
}

export interface APIOptions extends ClientSpawnOptions {
}
