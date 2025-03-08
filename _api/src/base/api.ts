import type {
    MaybeAsync,
    ParsedCommandLine,
    ProjectData,
    SymbolData,
    TypeData,
} from "../types.ts";
import { Node } from "./node.ts";

export interface APIOptions {
    tsserverPath: string;
    cwd?: string;
    logServer?: (msg: string) => void;
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
    abstract getSourceFile(fileName: string): MaybeAsync<Async, SourceFile<Async> | undefined>;
    abstract getSymbolAtPosition(fileName: string, position: number): MaybeAsync<Async, Symbol<Async> | undefined>;
}

export abstract class SourceFile<Async extends boolean> extends Node {
    constructor(data: Buffer) {
        const view = new DataView(data.buffer, data.byteOffset, data.byteLength);
        super(view, 1, undefined!);
    }

    nodeCount(): number {
        return (this.view.byteLength / (Node.NODE_LEN * 4)) - 1;
    }
}

export abstract class Symbol<Async extends boolean> {
    protected id: number;
    name: string;
    flags: number;
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
