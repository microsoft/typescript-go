import { SyncRpcChannel } from "libsyncrpc";

export interface ClientOptions {
    tsserverPath: string;
    cwd?: string;
}

export class Client {
    private channel: SyncRpcChannel;

    constructor(options: ClientOptions) {
        this.channel = new SyncRpcChannel(options.tsserverPath, [
            "api",
            "-cwd",
            options.cwd ?? process.cwd(),
        ]);
    }

    registerCallback(method: string, callback: (payload: any) => any): void {
        this.channel.registerCallback(method, (_, arg) => {
            const result = callback(JSON.parse(arg));
            return JSON.stringify(result) ?? "";
        });
        this.channel.requestSync("registerCallback", method);
    }

    request(method: string, payload: any): any {
        return JSON.parse(this.channel.requestSync(method, JSON.stringify(payload)));
    }

    close(): void {
        this.channel.murderInColdBlood();
    }
}
