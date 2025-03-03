import type { JSONRPCClient } from "json-rpc-2.0";
import type { ChildProcessWithoutNullStreams } from "node:child_process";
import {
    LSPError,
    LSPErrorCollection,
} from "./errors.ts";
import { createJSONRPCClient } from "./jsonRPC.ts";
import type {
    ParsedCommandLine,
    ProjectData,
    SymbolData,
} from "./types.ts";

export class Client {
    initialized = false;
    jsonRPCClient: JSONRPCClient;
    errorStack: LSPError[] = [];

    constructor(serverProcess: ChildProcessWithoutNullStreams) {
        this.jsonRPCClient = createJSONRPCClient(serverProcess, error => {
            this.errorStack.push(error);
        });
    }

    async initialize(): Promise<void> {
        await this.jsonRPCClient.request("initialize", {
            rootUri: null,
            capabilities: {},
        });
        this.initialized = true;
        this.jsonRPCClient.notify("initialized", {});
        this.flushErrors();
    }

    async parseConfigFile(configFileName: string): Promise<ParsedCommandLine> {
        return this.request("parseConfigFile", {
            configFileName,
        });
    }

    async loadProject(configFileName: string): Promise<ProjectData> {
        return this.request("loadProject", {
            configFileName,
        });
    }

    async getSymbolAtPosition(fileName: string, position: number): Promise<SymbolData | undefined> {
        return this.request("getSymbolAtPosition", {
            fileName,
            position,
        });
    }

    async shutdown(): Promise<void> {
        await this.jsonRPCClient.request("shutdown", {});
        this.jsonRPCClient.notify("exit", {});
        this.flushErrors();
    }

    exit(): void {
        this.jsonRPCClient.notify("exit", {});
    }

    async request(method: string, params: any): Promise<any> {
        const result = await this.jsonRPCClient.request(`@ts/${method}`, params);
        this.flushErrors();
        return result;
    }

    private flushErrors() {
        if (this.errorStack.length > 1) {
            throw new LSPErrorCollection(this.errorStack);
        }
        if (this.errorStack.length === 1) {
            throw this.errorStack[0];
        }
    }
}
