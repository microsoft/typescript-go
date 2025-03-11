import { SyncRpcChannel } from "libsyncrpc";
import type { FileSystemEntries } from "../types.ts";

export interface ClientOptions {
    tsserverPath: string;
    cwd?: string;
    logFile?: string;
    fs?: {
        directoryExists?: (directoryName: string) => boolean | undefined;
        fileExists?: (fileName: string) => boolean | undefined;
        getAccessibleEntries?: (directoryName: string) => FileSystemEntries | undefined;
        readFile?: (fileName: string) => string | null | undefined;
        realpath?: (path: string) => string | undefined;
    };
}

export class Client {
    private channel: SyncRpcChannel;

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
            const result = callback(JSON.parse(arg));
            return JSON.stringify(result) ?? "";
        });
        this.channel.requestSync("registerCallback", method);
    }

    request(method: string, payload: any): any {
        console.time("encode payload");
        const encodedPayload = JSON.stringify(payload);
        console.timeEnd("encode payload");
        console.time("request");
        const result = this.channel.requestSync(method, encodedPayload);
        console.timeEnd("request");
        console.time("decode result");
        const decodedResult = JSON.parse(result);
        console.timeEnd("decode result");
        return decodedResult;
    }

    requestBinary(method: string, payload: any): Buffer {
        return this.channel.requestBinarySync(method, JSON.stringify(payload));
    }

    close(): void {
        this.channel.murderInColdBlood();
    }
}
