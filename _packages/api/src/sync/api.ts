import type {
    Node,
    NodeArray,
    Statement,
    SyntaxKind,
} from "@typescript/ast";
import {
    type API as BaseAPI,
    type APIOptions as BaseAPIOptions,
    Project as BaseProject,
    RemoteSourceFile as BaseRemoteSourceFile,
    Symbol as BaseSymbol,
    Type as BaseType,
} from "../base/api.ts";
import type {
    ConfigResponse,
    GetSymbolAtPositionParams,
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

    echo(message: string): string {
        return this.client.echo(message);
    }

    echoBinary(message: Uint8Array): Uint8Array {
        return this.client.echoBinary(message);
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

    getSourceFile(fileName: string): SourceFile | undefined {
        const data = this.client.requestBinary("getSourceFile", { project: this.id, fileName });
        return data ? new RemoteSourceFile(this.client, this, data) as unknown as SourceFile : undefined;
    }

    getSymbolAtPosition(requests: readonly GetSymbolAtPositionParams[]): (Symbol | undefined)[];
    getSymbolAtPosition(fileName: string, position: number | number[]): Symbol | undefined;
    getSymbolAtPosition(...params: [fileName: string, position: number | number[]] | [readonly GetSymbolAtPositionParams[]]): Symbol | undefined | (Symbol | undefined)[] {
        if (params.length === 2) {
            if (typeof params[1] === "number") {
                const data = this.client.request("getSymbolAtPosition", { project: this.id, fileName: params[0], position: params[1] });
                return data ? new Symbol(this.client, this, data) : undefined;
            }
            const data = this.client.request("getSymbolAtPositions", { project: this.id, fileName: params[0], positions: params[1] });
            return data.map((d: SymbolResponse | null) => d ? new Symbol(this.client, this, d) : undefined);
        }
        const data = this.client.request("getSymbolAtPosition", params[0].map(({ fileName, position }) => ({ project: this.id, fileName, position })));
        return data.map((d: SymbolResponse | null) => d ? new Symbol(this.client, this, d) : undefined);
    }
}

export interface SourceFile extends Node {
    kind: SyntaxKind.SourceFile;
    statements: NodeArray<Statement>;
}

class RemoteSourceFile extends BaseRemoteSourceFile {
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

    constructor(client: Client, project: Project, data: SymbolResponse) {
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
    constructor(client: Client, data: TypeResponse) {
        super(data);
        this.client = client;
    }
}
