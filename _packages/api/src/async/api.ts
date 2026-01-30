/// <reference path="../node.ts" preserve="true" />
import { SymbolFlags } from "#symbolFlags";
import { TypeFlags } from "#typeFlags";
import type {
    Node,
    SourceFile,
} from "@typescript/ast";
import {
    type API as BaseAPI,
    type Project as BaseProject,
    type Symbol as BaseSymbol,
    type Type as BaseType,
} from "../base/api.ts";
import { ObjectRegistry } from "../base/objectRegistry.ts";
import { RemoteSourceFile } from "../node.ts";
import type {
    ConfigResponse,
    ProjectResponse,
    SourceFileResponse,
    SymbolResponse,
    TypeResponse,
} from "../proto.ts";
import { AsyncClient } from "./client.ts";

export { SymbolFlags, TypeFlags };

export interface LSPConnectionOptions {
    /** Path to the Unix domain socket for API communication */
    pipePath: string;
}

export interface AsyncAPIOptions {
    /** Path to the tsgo executable */
    tsserverPath: string;
    /** Current working directory */
    cwd?: string;
}

/** Type alias for the async object registry */
export type AsyncObjectRegistry = ObjectRegistry<AsyncProject, AsyncSymbol, AsyncType>;

export abstract class AsyncDisposableObject {
    private disposed: boolean = false;
    protected objectRegistry: AsyncObjectRegistry;
    abstract readonly id: string;

    constructor(objectRegistry: AsyncObjectRegistry) {
        this.objectRegistry = objectRegistry;
    }
    [globalThis.Symbol.dispose](): void {
        this.objectRegistry.release(this);
        this.disposed = true;
    }
    dispose(): void {
        this[globalThis.Symbol.dispose]();
    }
    isDisposed(): boolean {
        return this.disposed;
    }
    ensureNotDisposed(): this {
        if (this.disposed) {
            throw new Error(`${this.constructor.name} is disposed`);
        }
        return this;
    }
}

export class AsyncAPI implements BaseAPI<true> {
    private client: AsyncClient;
    private objectRegistry: AsyncObjectRegistry;

    /**
     * Create an AsyncAPI instance by spawning a new tsgo process.
     */
    constructor(options: AsyncAPIOptions) {
        this.client = new AsyncClient(options);
        // Create registry with factories - fire-and-forget release for async
        this.objectRegistry = new ObjectRegistry<AsyncProject, AsyncSymbol, AsyncType>(
            {
                createProject: data => new AsyncProject(this.client, this.objectRegistry, data),
                createSymbol: data => new AsyncSymbol(this.objectRegistry, data),
                createType: data => new AsyncType(this.objectRegistry, data),
            },
            id => {
                this.client.apiRequest("release", id).catch(() => {});
            },
        );
    }

    /**
     * Create an AsyncAPI instance from an existing LSP connection's API session.
     * Use this when connecting to an API pipe provided by an LSP server via $/initializeAPISession.
     */
    static fromLSPConnection(options: LSPConnectionOptions): AsyncAPI {
        const client = new AsyncClient(options);
        const api = Object.create(AsyncAPI.prototype) as AsyncAPI;
        api.client = client;
        api.objectRegistry = new ObjectRegistry<AsyncProject, AsyncSymbol, AsyncType>(
            {
                createProject: data => new AsyncProject(client, api.objectRegistry, data),
                createSymbol: data => new AsyncSymbol(api.objectRegistry, data),
                createType: data => new AsyncType(api.objectRegistry, data),
            },
            id => {
                client.apiRequest("release", id).catch(() => {});
            },
        );
        return api;
    }

    async parseConfigFile(fileName: string): Promise<ConfigResponse> {
        return this.client.apiRequest<ConfigResponse>("parseConfigFile", { fileName });
    }

    async loadProject(configFileName: string): Promise<AsyncProject> {
        const data = await this.client.apiRequest<ProjectResponse>("loadProject", { configFileName });
        return this.objectRegistry.getProject(data);
    }

    async getDefaultProjectForFile(fileName: string): Promise<AsyncProject | undefined> {
        const data = await this.client.apiRequest<ProjectResponse | null>("getDefaultProjectForFile", { fileName });
        return data ? this.objectRegistry.getProject(data) : undefined;
    }

    async close(): Promise<void> {
        await this.client.close();
        this.objectRegistry.clear();
    }
}

export class AsyncProject extends AsyncDisposableObject implements BaseProject<true> {
    private client: AsyncClient;
    private decoder = new TextDecoder();

    readonly id: string;
    configFileName!: string;
    compilerOptions!: Record<string, unknown>;
    rootFiles!: readonly string[];

    constructor(client: AsyncClient, objectRegistry: AsyncObjectRegistry, data: ProjectResponse) {
        super(objectRegistry);
        this.id = data.id;
        this.client = client;
        this.loadData(data);
    }

    loadData(data: ProjectResponse): void {
        this.configFileName = data.configFileName;
        this.compilerOptions = data.compilerOptions;
        this.rootFiles = data.rootFiles;
    }

    async reload(): Promise<void> {
        this.ensureNotDisposed();
        const data = await this.client.apiRequest<ProjectResponse>("loadProject", { configFileName: this.configFileName });
        this.loadData(data);
    }

    async getSourceFile(fileName: string): Promise<SourceFile | undefined> {
        this.ensureNotDisposed();
        const response = await this.client.apiRequest<SourceFileResponse | null>("getSourceFile", {
            project: this.id,
            fileName,
        });
        if (!response?.data) {
            return undefined;
        }
        // Decode base64 to Uint8Array
        const binaryData = Uint8Array.from(atob(response.data), c => c.charCodeAt(0));
        return new RemoteSourceFile(binaryData, this.decoder) as unknown as SourceFile;
    }

    getSymbolAtLocation(node: Node): Promise<AsyncSymbol | undefined>;
    getSymbolAtLocation(nodes: readonly Node[]): Promise<(AsyncSymbol | undefined)[]>;
    async getSymbolAtLocation(nodeOrNodes: Node | readonly Node[]): Promise<AsyncSymbol | (AsyncSymbol | undefined)[] | undefined> {
        this.ensureNotDisposed();
        if (Array.isArray(nodeOrNodes)) {
            const data = await this.client.apiRequest<(SymbolResponse | null)[]>("getSymbolsAtLocations", {
                project: this.id,
                locations: nodeOrNodes.map(node => node.id),
            });
            return data.map(d => d ? this.objectRegistry.getSymbol(d) : undefined);
        }
        const data = await this.client.apiRequest<SymbolResponse | null>("getSymbolAtLocation", {
            project: this.id,
            location: (nodeOrNodes as Node).id,
        });
        return data ? this.objectRegistry.getSymbol(data) : undefined;
    }

    getSymbolAtPosition(fileName: string, position: number): Promise<AsyncSymbol | undefined>;
    getSymbolAtPosition(fileName: string, positions: readonly number[]): Promise<(AsyncSymbol | undefined)[]>;
    async getSymbolAtPosition(fileName: string, positionOrPositions: number | readonly number[]): Promise<AsyncSymbol | (AsyncSymbol | undefined)[] | undefined> {
        this.ensureNotDisposed();
        if (typeof positionOrPositions === "number") {
            const data = await this.client.apiRequest<SymbolResponse | null>("getSymbolAtPosition", {
                project: this.id,
                fileName,
                position: positionOrPositions,
            });
            return data ? this.objectRegistry.getSymbol(data) : undefined;
        }
        const data = await this.client.apiRequest<(SymbolResponse | null)[]>("getSymbolsAtPositions", {
            project: this.id,
            fileName,
            positions: positionOrPositions,
        });
        return data.map(d => d ? this.objectRegistry.getSymbol(d) : undefined);
    }

    getTypeOfSymbol(symbol: AsyncSymbol): Promise<AsyncType | undefined>;
    getTypeOfSymbol(symbols: readonly AsyncSymbol[]): Promise<(AsyncType | undefined)[]>;
    async getTypeOfSymbol(symbolOrSymbols: AsyncSymbol | readonly AsyncSymbol[]): Promise<AsyncType | (AsyncType | undefined)[] | undefined> {
        this.ensureNotDisposed();
        if (Array.isArray(symbolOrSymbols)) {
            const data = await this.client.apiRequest<(TypeResponse | null)[]>("getTypesOfSymbols", {
                project: this.id,
                symbols: symbolOrSymbols.map(s => s.ensureNotDisposed().id),
            });
            return data.map(d => d ? this.objectRegistry.getType(d) : undefined);
        }
        const data = await this.client.apiRequest<TypeResponse | null>("getTypeOfSymbol", {
            project: this.id,
            symbol: (symbolOrSymbols as AsyncSymbol).ensureNotDisposed().id,
        });
        return data ? this.objectRegistry.getType(data) : undefined;
    }
}

export class AsyncSymbol extends AsyncDisposableObject implements BaseSymbol<true> {
    readonly id: string;
    readonly name: string;
    readonly flags: SymbolFlags;
    readonly checkFlags: number;

    constructor(objectRegistry: AsyncObjectRegistry, data: SymbolResponse) {
        super(objectRegistry);
        this.id = data.id;
        this.name = data.name;
        this.flags = data.flags;
        this.checkFlags = data.checkFlags;
    }
}

export class AsyncType extends AsyncDisposableObject implements BaseType<true> {
    readonly id: string;
    readonly flags: TypeFlags;

    constructor(objectRegistry: AsyncObjectRegistry, data: TypeResponse) {
        super(objectRegistry);
        this.id = data.id;
        this.flags = data.flags;
    }
}
