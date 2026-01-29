/// <reference path="../node.ts" preserve="true" />
import { SymbolFlags } from "#symbolFlags";
import { TypeFlags } from "#typeFlags";
import type {
    Node,
    SourceFile,
} from "@typescript/ast";
import {
    type API as BaseAPI,
    ObjectRegistry,
    type Project as BaseProject,
    type Symbol as BaseSymbol,
    type Type as BaseType,
} from "../base/index.ts";
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

    private constructor(client: AsyncClient) {
        this.client = client;
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
        return new AsyncAPI(client);
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

    async getSymbolAtLocation(node: Node): Promise<AsyncSymbol | undefined> {
        this.ensureNotDisposed();
        const data = await this.client.apiRequest<SymbolResponse | null>("getSymbolAtLocation", {
            project: this.id,
            location: node.id,
        });
        return data ? this.objectRegistry.getSymbol(data) : undefined;
    }

    async getSymbolsAtLocations(nodes: readonly Node[]): Promise<(AsyncSymbol | undefined)[]> {
        this.ensureNotDisposed();
        const data = await this.client.apiRequest<(SymbolResponse | null)[]>("getSymbolsAtLocations", {
            project: this.id,
            locations: nodes.map(node => node.id),
        });
        return data.map(d => d ? this.objectRegistry.getSymbol(d) : undefined);
    }

    async getSymbolAtPosition(fileName: string, position: number): Promise<AsyncSymbol | undefined> {
        this.ensureNotDisposed();
        const data = await this.client.apiRequest<SymbolResponse | null>("getSymbolAtPosition", {
            project: this.id,
            fileName,
            position,
        });
        return data ? this.objectRegistry.getSymbol(data) : undefined;
    }

    async getSymbolsAtPositions(fileName: string, positions: readonly number[]): Promise<(AsyncSymbol | undefined)[]> {
        this.ensureNotDisposed();
        const data = await this.client.apiRequest<(SymbolResponse | null)[]>("getSymbolsAtPositions", {
            project: this.id,
            fileName,
            positions,
        });
        return data.map(d => d ? this.objectRegistry.getSymbol(d) : undefined);
    }

    async getTypeOfSymbol(symbol: AsyncSymbol): Promise<AsyncType | undefined> {
        this.ensureNotDisposed();
        const data = await this.client.apiRequest<TypeResponse | null>("getTypeOfSymbol", {
            project: this.id,
            symbol: symbol.ensureNotDisposed().id,
        });
        return data ? this.objectRegistry.getType(data) : undefined;
    }

    async getTypesOfSymbols(symbols: readonly AsyncSymbol[]): Promise<(AsyncType | undefined)[]> {
        this.ensureNotDisposed();
        const data = await this.client.apiRequest<(TypeResponse | null)[]>("getTypesOfSymbols", {
            project: this.id,
            symbols: symbols.map(s => s.ensureNotDisposed().id),
        });
        return data.map(d => d ? this.objectRegistry.getType(d) : undefined);
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
