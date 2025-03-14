import type { JSONRPCClient } from "json-rpc-2.0";
import type { ChildProcessByStdio } from "node:child_process";
import type {
    Readable,
    Writable,
} from "node:stream";
import type {
    ConfigResponse,
    ProjectResponse,
    SymbolResponse,
} from "../base/proto.ts";
import {
    LSPError,
    LSPErrorCollection,
} from "./errors.ts";
import { createJSONRPCClient } from "./lsp.ts";

export class Client {
    initialized = false;
    jsonRPCClient: JSONRPCClient;
    errorStack: LSPError[] = [];

    constructor(serverProcess: ChildProcessByStdio<Writable, Readable, null>) {
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

    async parseConfigFile(configFileName: string): Promise<ConfigResponse> {
        return this.request("parseConfigFile", {
            configFileName,
        });
    }

    async loadProject(configFileName: string): Promise<ProjectResponse> {
        return this.request("loadProject", {
            configFileName,
        });
    }

    async getSymbolAtPosition(fileName: string, position: number): Promise<SymbolResponse | undefined> {
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
