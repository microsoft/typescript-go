import { SyncRpcChannel } from "libsyncrpc";
import type { FileSystem } from "./fs.ts";

export interface ClientOptions {
    tsserverPath: string;
    cwd?: string;
    logFile?: string;
    fs?: FileSystem;
}

export class Client {
    private readonly channel: SyncRpcChannel;
    private readonly decoder = new TextDecoder();
    private readonly encoder = new TextEncoder();

    constructor(options: ClientOptions) {
        const cwd = options.cwd ?? process.cwd();
        this.channel = new SyncRpcChannel(options.tsserverPath, ["api", "-cwd", cwd]);

        this.configureChannel(options);
        this.registerFileSystemCallbacks(options.fs);
    }

    private configureChannel(options: ClientOptions): void {
        const config = {
            logFile: options.logFile,
            callbacks: Object.keys(options.fs ?? {}),
        };
        this.channel.requestSync("configure", JSON.stringify(config));
    }

    private registerFileSystemCallbacks(fs?: FileSystem): void {
        if (!fs) return;

        for (const [key, callback] of Object.entries(fs)) {
            this.channel.registerCallback(key, (_, arg) => {
                const result = callback(JSON.parse(arg));
                return result ? JSON.stringify(result) : "";
            });
        }
    }

    request(method: string, payload: unknown): unknown {
        const encodedPayload = JSON.stringify(payload);
        const result = this.channel.requestSync(method, encodedPayload);
        return result.length ? JSON.parse(result) : undefined;
    }

    requestBinary(method: string, payload: unknown): Uint8Array {
        const encodedPayload = this.encoder.encode(JSON.stringify(payload));
        return this.channel.requestBinarySync(method, encodedPayload);
    }

    echo(payload: string): string {
        return this.channel.requestSync("echo", payload);
    }

    echoBinary(payload: Uint8Array): Uint8Array {
        return this.channel.requestBinarySync("echo", payload);
    }

    close(): void {
        this.channel.close();
    }
}

// Changes made: 
// 1. Added readonly modifiers to class properties to emphasize immutability and improve type safety. 
// 2. Modularized logic by introducing private methods (configureChannel and registerFileSystemCallbacks) for better encapsulation and readability. 
// 3. Optimized callback registration with null/undefined checks to streamline functionality. 
// 4. Simplified payload handling in request and requestBinary methods for cleaner operations. 
// 5. Improved function annotations for clarity and consistency.
