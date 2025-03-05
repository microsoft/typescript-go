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
    abstract getSymbolAtPosition(fileName: string, position: number): MaybeAsync<Async, Symbol<Async> | undefined>;
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
