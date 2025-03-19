import { SymbolFlags } from "#symbolFlags";
import type { SourceFile } from "@typescript/ast";
import { Client } from "./client.ts";
import type { FileSystem } from "./fs.ts";
import { RemoteSourceFile } from "./node.ts";
import type {
    ConfigResponse,
    GetSymbolAtPositionParams,
    ProjectResponse,
    SymbolResponse,
    TypeResponse,
} from "./proto.ts";

export { SymbolFlags };

export interface APIOptions {
    tsserverPath: string;
    cwd?: string;
    logFile?: string;
    fs?: FileSystem;
}

export class API {
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

export class Project {
    private decoder = new TextDecoder();
    private client: Client;

    id: number;
    configFileName!: string;
    compilerOptions!: Record<string, unknown>;
    rootFiles!: readonly string[];

    constructor(client: Client, data: ProjectResponse) {
        this.id = data.id;
        this.client = client;
        this.loadData(data);
    }

    loadData(data: ProjectResponse): void {
        this.configFileName = data.configFileName;
        this.compilerOptions = data.compilerOptions;
        this.rootFiles = data.rootFiles;
    }

    reload(): void {
        this.loadData(this.client.request("loadProject", { configFileName: this.configFileName }));
    }

    getSourceFile(fileName: string): SourceFile | undefined {
        const data = this.client.requestBinary("getSourceFile", { project: this.id, fileName });
        return data ? new RemoteSourceFile(data, this.decoder) as unknown as SourceFile : undefined;
    }

    getSymbolAtPosition(requests: readonly GetSymbolAtPositionParams[]): (Symbol | undefined)[];
    getSymbolAtPosition(fileName: string, position: number | number[]): Symbol | undefined;
    getSymbolAtPosition(...params: [fileName: string, position: number | number[]] | [readonly GetSymbolAtPositionParams[]]): Symbol | undefined | (Symbol | undefined)[] {
        if (params.length === 2) {
            if (typeof params[1] === "number") {
                const data = this.client.request("getSymbolAtPosition", { project: this.id, fileName: params[0], position: params[1] });
                return data ? new Symbol(this.client, data) : undefined;
            }
            const data = this.client.request("getSymbolAtPositions", { project: this.id, fileName: params[0], positions: params[1] });
            return data.map((d: SymbolResponse | null) => d ? new Symbol(this.client, d) : undefined);
        }
        const data = this.client.request("getSymbolAtPosition", params[0].map(({ fileName, position }) => ({ project: this.id, fileName, position })));
        return data.map((d: SymbolResponse | null) => d ? new Symbol(this.client, d) : undefined);
    }
}

export class Symbol {
    private client: Client;
    private id: number;
    name: string;
    flags: SymbolFlags;
    checkFlags: number;

    constructor(client: Client, data: SymbolResponse) {
        this.client = client;
        this.id = data.id;
        this.name = data.name;
        this.flags = data.flags;
        this.checkFlags = data.checkFlags;
    }
}

export class Type {
    private client: Client;
    private id: number;
    flags: number;
    constructor(client: Client, data: TypeResponse) {
        this.client = client;
        this.id = data.id;
        this.flags = data.flags;
    }
}
