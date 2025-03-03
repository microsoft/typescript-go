import type { ChildProcessWithoutNullStreams } from "node:child_process";
import { Client } from "./client.ts";
import { startServer } from "./server.ts";
import type {
    ParsedCommandLine,
    ProjectData,
    SymbolData,
} from "./types.ts";

export interface APIOptions {
    tsserverPath: string;
    cwd?: string;
    logServer?: (msg: string) => void;
}

export class API {
    private server: ChildProcessWithoutNullStreams;
    private client: Client;
    constructor(options: APIOptions) {
        this.client = new Client(this.server = startServer(options.tsserverPath, options.cwd ?? process.cwd(), options.logServer));
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

export class Project {
    private client: Client;
    configFileName!: string;
    commandLine!: ParsedCommandLine;

    constructor(client: Client, data: ProjectData) {
        this.client = client;
        this.loadData(data);
    }

    private loadData(data: ProjectData) {
        this.configFileName = data.configFileName;
        this.commandLine = data.commandLine;
    }

    async reload(): Promise<void> {
        this.loadData(await this.client.request("loadProject", { configFileName: this.configFileName }));
    }

    async getSymbolAtPosition(fileName: string, position: number): Promise<Symbol | undefined> {
        const data = await this.client.request("getSymbolAtPosition", { project: this.configFileName, fileName, position });
        return data ? new Symbol(this.client, this, data) : undefined;
    }
}

export class Symbol {
    private client: Client;
    private project: Project;

    private projectVersion: number;
    name: string;
    flags: number;
    checkFlags: number;

    constructor(client: Client, project: Project, data: SymbolData) {
        this.client = client;
        this.project = project;
        this.projectVersion = data.projectVersion;
        this.name = data.name;
        this.flags = data.flags;
        this.checkFlags = data.checkFlags;
    }
}
