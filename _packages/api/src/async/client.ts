import type { ChildProcess } from "node:child_process";
import type { Socket } from "node:net";
import {
    createMessageConnection,
    type MessageConnection,
    RequestType,
    SocketMessageReader,
    SocketMessageWriter,
    StreamMessageReader,
    StreamMessageWriter,
} from "vscode-jsonrpc/node";

export interface AsyncClientSocketOptions {
    /** Path to the Unix domain socket for API communication */
    pipePath: string;
}

export interface AsyncClientSpawnOptions {
    /** Path to the tsgo executable */
    tsserverPath: string;
    /** Current working directory */
    cwd?: string;
}

export type AsyncClientOptions = AsyncClientSocketOptions | AsyncClientSpawnOptions;

function isSpawnOptions(options: AsyncClientOptions): options is AsyncClientSpawnOptions {
    return "tsserverPath" in options;
}

/**
 * AsyncClient handles communication with the TypeScript API server
 * over STDIO (spawned process) or a Unix domain socket using JSON-RPC.
 */
export class AsyncClient {
    private socket: Socket | null = null;
    private process: ChildProcess | null = null;
    private connection: MessageConnection | null = null;
    private options: AsyncClientOptions;
    private connected = false;

    constructor(options: AsyncClientOptions) {
        this.options = options;
    }

    async connect(): Promise<void> {
        if (this.connected) return;

        if (isSpawnOptions(this.options)) {
            await this.connectViaSpawn(this.options);
        }
        else {
            await this.connectViaSocket(this.options);
        }
    }

    private async connectViaSpawn(options: AsyncClientSpawnOptions): Promise<void> {
        const { spawn } = await import("node:child_process");

        const args = [
            "--api",
            "--async",
            "-cwd",
            options.cwd ?? process.cwd(),
        ];

        this.process = spawn(options.tsserverPath, args, {
            stdio: ["pipe", "pipe", "inherit"],
        });

        const reader = new StreamMessageReader(this.process.stdout!);
        const writer = new StreamMessageWriter(this.process.stdin!);
        this.connection = createMessageConnection(reader, writer);
        this.connection.listen();
        this.connected = true;
    }

    private async connectViaSocket(options: AsyncClientSocketOptions): Promise<void> {
        const { createConnection } = await import("node:net");

        return new Promise((resolve, reject) => {
            this.socket = createConnection(options.pipePath, () => {
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
        if (this.process) {
            this.process.kill();
            this.process = null;
        }
        this.connected = false;
    }
}
