import type { APIOptions as BaseAPIOptions } from "../base/api.ts";
import {
    type API as BaseAPI,
    Project as BaseProject,
    Symbol as BaseSymbol,
    Type as BaseType,
} from "../base/api.ts";
import type {
    FileSystemEntries,
    ParsedCommandLine,
    ProjectData,
    SymbolData,
    TypeData,
} from "../types.ts";
import { Client } from "./client.ts";

export interface APIOptions extends BaseAPIOptions {
    fs?: {
        directoryExists?: (directoryName: string) => boolean | undefined;
        fileExists?: (fileName: string) => boolean | undefined;
        getAccessibleEntries?: (directoryName: string) => FileSystemEntries | undefined;
        getEntries?: (directoryName: string) => FileSystemEntries | undefined;
        readFile?: (fileName: string) => string | null | undefined;
        realpath?: (path: string) => string | undefined;
    };
}

export class API implements BaseAPI<false> {
    private client: Client;
    constructor(options: APIOptions) {
        this.client = new Client(options);
        if (options.fs) {
            if (options.fs) {
                if (options.fs.directoryExists) {
                    this.client.registerCallback("directoryExists", options.fs.directoryExists);
                }
                if (options.fs.fileExists) {
                    this.client.registerCallback("fileExists", options.fs.fileExists);
                }
                if (options.fs.getAccessibleEntries) {
                    this.client.registerCallback("getAccessibleEntries", options.fs.getAccessibleEntries);
                }
                if (options.fs.getEntries) {
                    this.client.registerCallback("getEntries", options.fs.getEntries);
                }
                if (options.fs.readFile) {
                    this.client.registerCallback("readFile", options.fs.readFile);
                }
                if (options.fs.realpath) {
                    this.client.registerCallback("realpath", options.fs.realpath);
                }
            }
        }
    }

    parseConfigFile(fileName: string): ParsedCommandLine {
        return this.client.request("parseConfigFile", { fileName });
    }

    loadProject(configFileName: string): Project {
        const data = this.client.request("loadProject", { configFileName });
        return new Project(this.client, data);
    }

    close(): void {
        this.client.close();
    }
}

export class Project extends BaseProject<false> {
    private client: Client;

    constructor(client: Client, data: ProjectData) {
        super(data);
        this.client = client;
    }

    reload(): void {
        this.loadData(this.client.request("loadProject", { configFileName: this.configFileName }));
    }

    getSymbolAtPosition(fileName: string, position: number): Symbol | undefined {
        const data = this.client.request("getSymbolAtPosition", { project: this.configFileName, fileName, position });
        return data ? new Symbol(this.client, this, data) : undefined;
    }
}

export class Symbol extends BaseSymbol<false> {
    private client: Client;
    private project: Project;

    constructor(client: Client, project: Project, data: SymbolData) {
        super(data);
        this.client = client;
        this.project = project;
    }

    getType(): Type | undefined {
        const data = this.client.request("getTypeOfSymbol", { project: this.project.configFileName, symbol: this.id });
        return data ? new Type(this.client, data) : undefined;
    }
}

export class Type extends BaseType<false> {
    private client: Client;
    constructor(client: Client, data: TypeData) {
        super(data);
        this.client = client;
    }
}
