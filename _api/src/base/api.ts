import { SymbolFlags } from "#symbolFlags";
import type { SourceFile as SourceFileNode } from "../ast/ast.ts";
import { RemoteNode } from "../ast/node.ts";
import type {
    MaybeAsync,
    ParsedCommandLine,
    ProjectData,
    SymbolData,
    TypeData,
} from "../types.ts";

export interface APIOptions {
    tsserverPath: string;
    cwd?: string;
    logFile?: string;
}

export interface API<Async extends boolean> {
    parseConfigFile(fileName: string): MaybeAsync<Async, ParsedCommandLine>;
    loadProject(configFileName: string): MaybeAsync<Async, Project<Async>>;
}

export abstract class Project<Async extends boolean> {
    configFileName!: string;
    compilerOptions!: Record<string, unknown>;
    rootFiles!: readonly string[];

    constructor(data: ProjectData) {
        this.loadData(data);
    }

    loadData(data: ProjectData): void {
        this.configFileName = data.configFileName;
        this.compilerOptions = data.compilerOptions;
        this.rootFiles = data.rootFiles;
    }

    abstract reload(): MaybeAsync<Async, void>;
    abstract getSourceFile(fileName: string): MaybeAsync<Async, SourceFileNode | undefined>;
    abstract getSymbolAtPosition(fileName: string, position: number): MaybeAsync<Async, Symbol<Async> | undefined>;
}

export abstract class SourceFile extends RemoteNode {
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

    constructor(data: SymbolData) {
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

    constructor(data: TypeData) {
        this.id = data.id;
        this.flags = data.flags;
    }
}
