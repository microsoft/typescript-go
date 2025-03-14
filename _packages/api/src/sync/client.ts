import { SyncRpcChannel } from "libsyncrpc";
import type { FileSystem } from "../base/fs.ts";

export interface ClientOptions {
    tsserverPath: string;
    cwd?: string;
    logFile?: string;
    fs?: FileSystem;
}

export class Client {
    private channel: SyncRpcChannel;
    private decoder = new TextDecoder();
    private encoder = new TextEncoder();

    constructor(options: ClientOptions) {
        this.channel = new SyncRpcChannel(options.tsserverPath, [
            "api",
            "-cwd",
            options.cwd ?? process.cwd(),
        ]);

        this.channel.requestSync(
            "configure",
            JSON.stringify({
                logFile: options.logFile,
                callbacks: Object.keys(options.fs ?? {}),
            }),
        );
    }

    registerCallback(method: string, callback: (payload: any) => any): void {
        this.channel.registerCallback(method, (_, arg) => {
            const result = callback(JSON.parse(this.decoder.decode(arg)));
            return JSON.stringify(result) ?? "";
        });
        this.channel.requestSync("registerCallback", method);
    }

    request(method: string, payload: any): any {
        const encodedPayload = JSON.stringify(payload);
        const result = this.channel.requestSync(method, encodedPayload);
        const decodedResult = JSON.parse(result);
        return decodedResult;
    }

    requestBinary(method: string, payload: any): Uint8Array {
        return this.channel.requestBinarySync(method, this.encoder.encode(JSON.stringify(payload)));
    }

    close(): void {
        this.channel.murderInColdBlood();
    }
}
