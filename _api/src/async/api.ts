import type { ChildProcessWithoutNullStreams } from "child_process";
import type { ParsedCommandLine } from "typescript";
import type { APIOptions } from "../base/api.ts";
import {
    type API as BaseAPI,
    Project as BaseProject,
    Symbol as BaseSymbol,
    Type as BaseType,
} from "../base/api.ts";
import type {
    ProjectResponse,
    SymbolResponse,
    TypeResponse,
} from "../base/proto.ts";
import { Client } from "./client.ts";
import { startLSPServer } from "./lsp.ts";

export class API implements BaseAPI<true> {
    private server: ChildProcessWithoutNullStreams;
    private client: Client;

    constructor(options: APIOptions) {
        this.client = new Client(this.server = startLSPServer(options.tsserverPath, options.cwd ?? process.cwd(), options.logServer));
    }

    async parseConfigFile(fileName: string): Promise<ParsedCommandLine> {
        await this.ensureInitialized();
        return this.client.request("parseConfigFile", { fileName });
    }

    async loadProject(configFileName: string): Promise<Project> {
        await this.ensureInitialized();
        const data = await this.client.loadProject(configFileName);
        return new Project(this.client, data);
    }

    async ensureInitialized(): Promise<void> {
        if (!this.client.initialized) {
            await this.client.initialize();
        }
    }

    async close(): Promise<void> {
        await this.client.shutdown();
        this.client.exit();
        this.server.kill();
    }
}

export class Project extends BaseProject<true> {
    private client: Client;

    constructor(client: Client, data: ProjectResponse) {
        super(data);
        this.client = client;
    }

    async reload(): Promise<void> {
        this.loadData(await this.client.request("loadProject", { configFileName: this.configFileName }));
    }

    async getSymbolAtPosition(fileName: string, position: number): Promise<Symbol | undefined> {
        const data = await this.client.request("getSymbolAtPosition", { project: this.configFileName, fileName, position });
        return data ? new Symbol(this.client, this, data) : undefined;
    }
}

export class Symbol extends BaseSymbol<true> {
    private client: Client;
    private project: Project;

    constructor(client: Client, project: Project, data: SymbolResponse) {
        super(data);
        this.client = client;
        this.project = project;
    }

    async getType(): Promise<Type | undefined> {
        const data = await this.client.request("getTypeOfSymbol", { project: this.project.configFileName, symbol: this.id });
        return data ? new Type(this.client, data) : undefined;
    }
}

export class Type extends BaseType<true> {
    private client: Client;

    constructor(client: Client, data: TypeResponse) {
        super(data);
        this.client = client;
    }
}
