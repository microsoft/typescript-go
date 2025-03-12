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
    ConfigResponse,
    ProjectResponse,
    SymbolResponse,
    TypeResponse,
} from "../base/proto.ts";
import { Client } from "./client.ts";

export interface APIOptions extends BaseAPIOptions {
}

export class API implements BaseAPI<false> {
    private client: Client;
    constructor(options: APIOptions) {
        this.client = new Client(options);
    }

    parseConfigFile(fileName: string): ConfigResponse {
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

    constructor(client: Client, data: ProjectResponse) {
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
            const data = this.client.getSymbolAtPosition(this.id, params[0], params[1]);
            return data.length ? new Symbol(this.client, data) : undefined;
        }
        else {
            // const data = this.client.request("getSymbolAtPosition", params[0].map(({ fileName, position }) => ({ project: this.id, fileName, position })));
            // return data.map((d: SymbolResponse | null) => d ? new Symbol(this.client, this, d) : undefined);
        }
    }

    getTypeOfSymbol(symbol: Symbol): Type | undefined {
        const data = this.client.getTypeOfSymbol(this.id, symbol.id);
        return data ? new Type(this.client, data) : undefined;
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

export class Symbol extends BaseSymbol {
    private client: Client;

    constructor(client: Client, data: Uint8Array) {
        super(data, client.decoder);
        this.client = client;
    }
}

export class Type extends BaseType<false> {
    private client: Client;
    constructor(client: Client, data: Uint8Array) {
        super(data);
        this.client = client;
    }
}
