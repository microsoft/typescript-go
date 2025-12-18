import type { Socket } from "node:net";
import {
    createMessageConnection,
    type MessageConnection,
    RequestType,
    SocketMessageReader,
    SocketMessageWriter,
} from "vscode-jsonrpc/node";

export interface AsyncClientOptions {
    pipePath: string;
}

/**
 * AsyncClient handles communication with the TypeScript API server
 * over a Unix domain socket using JSON-RPC.
 */
export class AsyncClient {
    private socket: Socket | null = null;
    private connection: MessageConnection | null = null;
    private options: AsyncClientOptions;
    private connected = false;

    constructor(options: AsyncClientOptions) {
        this.options = options;
    }

    async connect(): Promise<void> {
        if (this.connected) return;

        const { createConnection } = await import("node:net");

        return new Promise((resolve, reject) => {
            this.socket = createConnection(this.options.pipePath, () => {
                const reader = new SocketMessageReader(this.socket!);
                const writer = new SocketMessageWriter(this.socket!);
                this.connection = createMessageConnection(reader, writer);
                this.connection.listen();
                this.connected = true;
                resolve();
            });

            this.socket.on("error", error => {
                reject(new Error(`Socket error: ${error.message}`));
            });
        });
    }

    async apiRequest<T>(method: string, params?: unknown): Promise<T> {
        if (!this.connected) {
            await this.connect();
        }
        if (!this.connection) {
            throw new Error("Connection not established");
        }

        const requestType = new RequestType<unknown, T, void>(method);
        return this.connection.sendRequest(requestType, params);
    }

    async close(): Promise<void> {
        if (this.connection) {
            this.connection.dispose();
            this.connection = null;
        }
        if (this.socket) {
            this.socket.destroy();
            this.socket = null;
        }
        this.connected = false;
    }
}
