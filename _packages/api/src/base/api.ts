import { SymbolFlags } from "#symbolFlags";
import type { SourceFile as SourceFileNode } from "@typescript/ast";
import { RemoteNode } from "./node.ts";
import type {
    ConfigResponse,
    ProjectResponse,
    SymbolResponse,
    TypeResponse,
} from "./proto.ts";

export type MaybeAsync<Async extends boolean, T> = Async extends true ? Promise<T> : T;

export interface APIOptions {
    tsserverPath: string;
    cwd?: string;
    logFile?: string;
}

export interface API<Async extends boolean> {
    parseConfigFile(fileName: string): MaybeAsync<Async, ConfigResponse>;
    loadProject(configFileName: string): MaybeAsync<Async, Project<Async>>;
}

export abstract class Project<Async extends boolean> {
    id: number;
    configFileName!: string;
    compilerOptions!: Record<string, unknown>;
    rootFiles!: readonly string[];

    constructor(data: ProjectResponse) {
        this.id = data.id;
        this.loadData(data);
    }

    loadData(data: ProjectResponse): void {
        this.configFileName = data.configFileName;
        this.compilerOptions = data.compilerOptions;
        this.rootFiles = data.rootFiles;
    }

    abstract reload(): MaybeAsync<Async, void>;
    abstract getSourceFile(fileName: string): MaybeAsync<Async, SourceFileNode | undefined>;
    abstract getSymbolAtPosition(fileName: string, position: number): MaybeAsync<Async, Symbol<Async> | undefined>;
}

export abstract class RemoteSourceFile extends RemoteNode {
    constructor(data: Uint8Array) {
        const view = new DataView(data.buffer, data.byteOffset, data.byteLength);
        super(view, 1, undefined!);
    }
}

export { SymbolFlags };

export abstract class Symbol<Async extends boolean> {
    protected id: number;
    name: string;
    flags: SymbolFlags;
    checkFlags: number;

    constructor(data: SymbolResponse) {
        this.id = data.id;
        this.name = data.name;
        this.flags = data.flags;
        this.checkFlags = data.checkFlags;
    }

    abstract getType(): MaybeAsync<Async, Type<Async> | undefined>;
}

export abstract class Type<Async extends boolean> {
    protected id: number;
    flags: number;

    constructor(data: TypeResponse) {
        this.id = data.id;
        this.flags = data.flags;
    }
}
