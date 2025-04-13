/// <reference path="./node.ts" preserve="true" />
import { SymbolFlags, TypeFlags } from "#symbolFlags";
import type { Node, SourceFile } from "@typescript/ast";
import { Client } from "./client.ts";
import type { FileSystem } from "./fs.ts";
import { RemoteSourceFile } from "./node.ts";
import { ObjectRegistry } from "./objectRegistry.ts";
import type {
    ConfigResponse,
    ProjectResponse,
    SymbolResponse,
    TypeResponse,
} from "./proto.ts";

export { SymbolFlags, TypeFlags };

export interface APIOptions {
    tsserverPath: string;
    cwd?: string;
    logFile?: string;
    fs?: FileSystem;
}

export class API {
    private readonly client: Client;
    private readonly objectRegistry: ObjectRegistry;

    constructor(options: APIOptions) {
        this.client = new Client(options);
        this.objectRegistry = new ObjectRegistry(this.client);
    }

    parseConfigFile(fileName: string): ConfigResponse {
        return this.client.request("parseConfigFile", { fileName });
    }

    loadProject(configFileName: string): Project {
        const data = this.client.request("loadProject", { configFileName });
        return this.objectRegistry.getProject(data);
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

export class DisposableObject {
    private disposed = false;
    protected readonly objectRegistry: ObjectRegistry;

    constructor(objectRegistry: ObjectRegistry) {
        this.objectRegistry = objectRegistry;
    }

    [globalThis.Symbol.dispose](): void {
        this.dispose();
    }

    dispose(): void {
        if (!this.disposed) {
            this.objectRegistry.release(this);
            this.disposed = true;
        }
    }

    isDisposed(): boolean {
        return this.disposed;
    }

    protected ensureNotDisposed(): void {
        if (this.disposed) {
            throw new Error(`${this.constructor.name} is disposed`);
        }
    }
}

export class Project extends DisposableObject {
    private readonly decoder = new TextDecoder();
    private readonly client: Client;

    readonly id: string;
    configFileName!: string;
    compilerOptions!: Record<string, unknown>;
    rootFiles!: readonly string[];

    constructor(client: Client, objectRegistry: ObjectRegistry, data: ProjectResponse) {
        super(objectRegistry);
        this.client = client;
        this.id = data.id;
        this.loadData(data);
    }

    private loadData(data: ProjectResponse): void {
        this.configFileName = data.configFileName;
        this.compilerOptions = data.compilerOptions;
        this.rootFiles = data.rootFiles;
    }

    reload(): void {
        this.ensureNotDisposed();
        this.loadData(
            this.client.request("loadProject", { configFileName: this.configFileName })
        );
    }

    getSourceFile(fileName: string): SourceFile | undefined {
        this.ensureNotDisposed();
        const data = this.client.requestBinary("getSourceFile", {
            project: this.id,
            fileName,
        });
        return data ? (new RemoteSourceFile(data, this.decoder) as SourceFile) : undefined;
    }

    getSymbolAtLocation(
        nodeOrNodes: Node | readonly Node[]
    ): Symbol | (Symbol | undefined)[] | undefined {
        this.ensureNotDisposed();
        if (Array.isArray(nodeOrNodes)) {
            return this.client
                .request("getSymbolsAtLocations", {
                    project: this.id,
                    locations: nodeOrNodes.map((node) => node.id),
                })
                .map((d: SymbolResponse | null) =>
                    d ? this.objectRegistry.getSymbol(d) : undefined
                );
        }

        const data = this.client.request("getSymbolAtLocation", {
            project: this.id,
            location: nodeOrNodes.id,
        });
        return data ? this.objectRegistry.getSymbol(data) : undefined;
    }

    getSymbolAtPosition(
        fileName: string,
        positionOrPositions: number | readonly number[]
    ): Symbol | (Symbol | undefined)[] | undefined {
        this.ensureNotDisposed();
        if (typeof positionOrPositions === "number") {
            const data = this.client.request("getSymbolAtPosition", {
                project: this.id,
                fileName,
                position: positionOrPositions,
            });
            return data ? this.objectRegistry.getSymbol(data) : undefined;
        }

        const data = this.client.request("getSymbolsAtPositions", {
            project: this.id,
            fileName,
            positions: positionOrPositions,
        });
        return data.map((d: SymbolResponse | null) =>
            d ? this.objectRegistry.getSymbol(d) : undefined
        );
    }

    getTypeOfSymbol(
        symbolOrSymbols: Symbol | readonly Symbol[]
    ): Type | (Type | undefined)[] | undefined {
        this.ensureNotDisposed();
        if (Array.isArray(symbolOrSymbols)) {
            return this.client
                .request("getTypesOfSymbols", {
                    project: this.id,
                    symbols: symbolOrSymbols.map((symbol) =>
                        symbol.ensureNotDisposed().id
                    ),
                })
                .map((d: TypeResponse | null) =>
                    d ? this.objectRegistry.getType(d) : undefined
                );
        }

        const data = this.client.request("getTypeOfSymbol", {
            project: this.id,
            symbol: symbolOrSymbols.ensureNotDisposed().id,
        });
        return data ? this.objectRegistry.getType(data) : undefined;
    }
}

export class Symbol extends DisposableObject {
    readonly id: string;
    readonly name: string;
    readonly flags: SymbolFlags;
    readonly checkFlags: number;

    constructor(client: Client, objectRegistry: ObjectRegistry, data: SymbolResponse) {
        super(objectRegistry);
        this.id = data.id;
        this.name = data.name;
        this.flags = data.flags;
        this.checkFlags = data.checkFlags;
    }
}

export class Type extends DisposableObject {
    readonly id: string;
    readonly flags: TypeFlags;

    constructor(client: Client, objectRegistry: ObjectRegistry, data: TypeResponse) {
        super(objectRegistry);
        this.id = data.id;
        this.flags = data.flags;
    }
}

// Changes made: 
// 1. Replaced mutable fields with readonly where applicable to emphasize immutability. 
// 2. Modularized logic with private methods (loadData) for better encapsulation. 
// 3. Ensured consistent null/undefined checks with concise expressions. 
// 4. Improved function annotations for type safety and clarity. 
// 5. Streamlined dispose method to avoid redundant calls.
