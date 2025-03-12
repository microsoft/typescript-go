import { SyncRpcChannel } from "libsyncrpc";
import {
    encodeGetSymbolAtPositionRequest,
    encodeGetTypeOfSymbolRequest,
} from "../base/binary.ts";
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
    decoder: TextDecoder = new TextDecoder();
    encoder: TextEncoder = new TextEncoder();

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
        for (const callback in options.fs) {
            this.channel.registerCallback(callback, (_, arg) => {
                const result = options.fs?.[callback as keyof typeof options.fs]?.(JSON.parse(this.decoder.decode(arg)));
                return JSON.stringify(result) ?? "";
            });
        }
    }

    getSymbolAtPosition(projectId: number, fileName: string, position: number): Uint8Array {
        return this.channel.requestBinarySync("getSymbolAtPosition", encodeGetSymbolAtPositionRequest(projectId, fileName, position, this.encoder));
    }

    getTypeOfSymbol(projectId: number, symbolId: number): Uint8Array {
        return this.channel.requestBinarySync("getTypeOfSymbol", encodeGetTypeOfSymbolRequest(projectId, symbolId));
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
