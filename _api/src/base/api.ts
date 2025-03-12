import { SymbolFlags } from "#symbolFlags";
import type { SourceFile as SourceFileNode } from "../ast/ast.ts";
import { RemoteNode } from "../ast/node.ts";
import type { MaybeAsync } from "../types.ts";
import type {
    ConfigResponse,
    ProjectResponse,
} from "./proto.ts";

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
    abstract getSymbolAtPosition(fileName: string, position: number): MaybeAsync<Async, Symbol | undefined>;
}

export abstract class SourceFile extends RemoteNode {
    constructor(data: Uint8Array) {
        const view = new DataView(data.buffer, data.byteOffset, data.byteLength);
        super(view, 1, undefined!);
    }
}

export { SymbolFlags };

export abstract class Symbol {
    private data: Uint8Array;
    private view: DataView;
    private decoder: TextDecoder;

    constructor(data: Uint8Array, decoder: TextDecoder) {
        this.data = data;
        this.view = new DataView(data.buffer, data.byteOffset, data.byteLength);
        this.decoder = decoder;
    }

    get id(): number {
        return this.view.getUint32(0, true);
    }

    get flags(): number {
        return this.view.getUint32(4, true);
    }

    get checkFlags(): number {
        return this.view.getUint32(8, true);
    }

    get name(): string {
        return this.decoder.decode(this.data.subarray(12));
    }
}

export abstract class Type<Async extends boolean> {
    private view: DataView;

    constructor(data: Uint8Array) {
        this.view = new DataView(data.buffer, data.byteOffset, data.byteLength);
    }

    get id(): number {
        return this.view.getUint32(0, true);
    }

    get flags(): number {
        return this.view.getUint32(4, true);
    }
}
