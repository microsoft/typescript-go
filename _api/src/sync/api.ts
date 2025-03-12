import type { SourceFile as SourceFileNode } from "../ast/ast.ts";
import {
    type API as BaseAPI,
    type APIOptions as BaseAPIOptions,
    Project as BaseProject,
    SourceFile as BaseSourceFile,
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
}

export class API implements BaseAPI<false> {
    private client: Client;
    constructor(options: APIOptions) {
        this.client = new Client(options);
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

    getSourceFile(fileName: string): SourceFileNode | undefined {
        const data = this.client.requestBinary("getSourceFile", { project: this.id, fileName });
        return data ? new SourceFile(this.client, this, data) as unknown as SourceFileNode : undefined;
    }

    getSymbolAtPosition(requests: readonly { fileName: string; position: number; }[]): (Symbol | undefined)[];
    getSymbolAtPosition(fileName: string, position: number): Symbol | undefined;
    getSymbolAtPosition(...params: [fileName: string, position: number] | [readonly { fileName: string; position: number; }[]]): Symbol | undefined | (Symbol | undefined)[] {
        if (params.length === 2) {
            const data = this.client.request("getSymbolAtPosition", { project: this.id, fileName: params[0], position: params[1] });
            return data ? new Symbol(this.client, this, data) : undefined;
        }
        else {
            const data = this.client.request("getSymbolAtPosition", params[0].map(({ fileName, position }) => ({ project: this.id, fileName, position })));
            return data.map((d: SymbolData | null) => d ? new Symbol(this.client, this, d) : undefined);
        }
    }
}

export class SourceFile extends BaseSourceFile {
    private client: Client;
    private project: Project;
    constructor(client: Client, project: Project, data: Uint8Array) {
        super(data);
        this.client = client;
        this.project = project;
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
        const data = this.client.request("getTypeOfSymbol", { project: this.project.id, symbol: this.id });
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
