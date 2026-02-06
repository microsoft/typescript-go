/// <reference path="../node.ts" preserve="true" />
import { SymbolFlags } from "#symbolFlags";
import { TypeFlags } from "#typeFlags";
import type {
    Node,
    Path,
    SourceFile,
    SyntaxKind,
} from "@typescript/ast";
import {
    type API as BaseAPI,
    type DocumentIdentifier,
    type DocumentPosition,
    type NodeHandle as BaseNodeHandle,
    type Project as BaseProject,
    resolveFileName,
    type Symbol as BaseSymbol,
    type Type as BaseType,
} from "../base/api.ts";
import { ObjectRegistry } from "../base/objectRegistry.ts";
import { SourceFileCache } from "../base/sourceFileCache.ts";
import {
    findDescendant,
    parseNodeHandle,
    RemoteSourceFile,
} from "../node.ts";
import {
    createGetCanonicalFileName,
    toPath,
} from "../path.ts";
import type {
    AdoptLSPStateResponse,
    ConfigResponse,
    InitializeResponse,
    LoadProjectResponse,
    ProjectChanges,
    ProjectResponse,
    SourceFileResponse,
    SymbolResponse,
    TypeResponse,
} from "../proto.ts";
import {
    Client,
    type ClientSocketOptions,
    type ClientSpawnOptions,
} from "./client.ts";

export { SymbolFlags, TypeFlags };
export type { DocumentIdentifier, DocumentPosition };
export { documentURIToFileName, fileNameToDocumentURI } from "../path.ts";

export interface LSPConnectionOptions extends ClientSocketOptions {
}

export interface APIOptions extends ClientSpawnOptions {
}

/** Type alias for the async object registry */
type AsyncObjectRegistry = ObjectRegistry<Project, Symbol, Type>;

export abstract class DisposableObject {
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

export class API implements BaseAPI<true> {
    private client: Client;
    private objectRegistry: AsyncObjectRegistry;
    private sourceFileCache: SourceFileCache;
    private toPath: ((fileName: string) => Path) | undefined;
    private initialized: boolean = false;

    constructor(options: APIOptions | LSPConnectionOptions) {
        this.client = new Client(options);
        this.sourceFileCache = new SourceFileCache();
        this.objectRegistry = new ObjectRegistry<Project, Symbol, Type>(
            {
                createProject: data => new Project(this.client, this.objectRegistry, this.sourceFileCache, this.toPath!, data),
                createSymbol: data => new Symbol(this.objectRegistry, data),
                createType: data => new Type(this.objectRegistry, data),
            },
            id => {
                // Currently fire-and-forget, may need to track failure here
                this.client.apiRequest("release", id).catch(() => {});
            },
        );
    }

    /**
     * Create an API instance from an existing LSP connection's API session.
     * Use this when connecting to an API pipe provided by an LSP server via custom/initializeAPISession.
     */
    static async fromLSPConnection(options: LSPConnectionOptions): Promise<API> {
        const api = new API(options);
        await api.ensureInitialized();
        return api;
    }

    private async ensureInitialized(): Promise<void> {
        if (!this.initialized) {
            const response = await this.client.apiRequest<InitializeResponse>("initialize", null);
            const getCanonicalFileName = createGetCanonicalFileName(response.useCaseSensitiveFileNames);
            const currentDirectory = response.currentDirectory;
            this.toPath = (fileName: string) => toPath(fileName, currentDirectory, getCanonicalFileName) as Path;
        }
    }

    async parseConfigFile(file: DocumentIdentifier | string): Promise<ConfigResponse> {
        await this.ensureInitialized();
        return this.client.apiRequest<ConfigResponse>("parseConfigFile", { fileName: resolveFileName(file) });
    }

    /**
     * Adopt the latest state from the LSP server.
     * Only meaningful when connected to an LSP server via `fromLSPConnection`.
     *
     * This method invalidates cached source files based on the server's diff response,
     * which indicates which files have changed or been removed since the last snapshot.
     */
    async adoptLSPState(): Promise<void> {
        await this.ensureInitialized();
        const response = await this.client.apiRequest<AdoptLSPStateResponse>("adoptLSPState");

        // Handleremoved projects - release their cached files
        if (response.removedProjects) {
            for (const projectId of response.removedProjects) {
                this.sourceFileCache.releaseProject(projectId);
            }
        }

        // Handle file changes within projects
        if (response.projectChanges) {
            for (const [_projectId, changes] of Object.entries(response.projectChanges)) {
                // Mark changed files as dirty (they'll be re-fetched on next access)
                if (changes.changedFiles) {
                    this.sourceFileCache.remove(changes.changedFiles as Path[]);
                }
                // Remove deleted files from cache entirely
                if (changes.removedFiles) {
                    this.sourceFileCache.remove(changes.removedFiles as Path[]);
                }
            }
        }
    }

    async loadProject(configFile: DocumentIdentifier | string): Promise<Project> {
        await this.ensureInitialized();
        const data = await this.client.apiRequest<ProjectResponse>("loadProject", { configFileName: resolveFileName(configFile) });
        return this.objectRegistry.getProject(data);
    }

    async getDefaultProjectForFile(file: DocumentIdentifier | string): Promise<Project | undefined> {
        await this.ensureInitialized();
        const data = await this.client.apiRequest<ProjectResponse | null>("getDefaultProjectForFile", { fileName: resolveFileName(file) });
        return data ? this.objectRegistry.getProject(data) : undefined;
    }

    async close(): Promise<void> {
        await this.client.close();
        this.objectRegistry.clear();
        this.sourceFileCache.clear();
    }
}

export class Project extends DisposableObject implements BaseProject<true> {
    private client: Client;
    private sourceFileCache: SourceFileCache;
    private decoder = new TextDecoder();

    readonly id: string;
    readonly toPath: (fileName: string) => Path;
    configFileName!: string;
    compilerOptions!: Record<string, unknown>;
    rootFiles!: readonly string[];
    parseOptionsKey!: string;

    constructor(client: Client, objectRegistry: AsyncObjectRegistry, sourceFileCache: SourceFileCache, toPath: (fileName: string) => Path, data: ProjectResponse) {
        super(objectRegistry);
        this.id = data.id;
        this.client = client;
        this.sourceFileCache = sourceFileCache;
        this.toPath = toPath;
        this.loadData(data);
    }

    loadData(data: ProjectResponse): void {
        this.configFileName = data.configFileName;
        this.compilerOptions = data.compilerOptions;
        this.rootFiles = data.rootFiles;
        this.parseOptionsKey = data.parseOptionsKey;
    }

    async reload(): Promise<ProjectChanges | undefined> {
        this.ensureNotDisposed();
        const data = await this.client.apiRequest<LoadProjectResponse>("loadProject", { configFileName: this.configFileName });
        // !!! TODO: handle parseOptionsKey change effect on source file cache
        this.loadData(data);

        // Handle cache invalidation based on changes
        if (data.changes) {
            if (data.changes.changedFiles) {
                this.sourceFileCache.remove(data.changes.changedFiles as Path[]);
            }
            if (data.changes.removedFiles) {
                this.sourceFileCache.remove(data.changes.removedFiles as Path[]);
            }
        }

        return data.changes;
    }

    async getSourceFile(file: DocumentIdentifier | string): Promise<SourceFile | undefined> {
        this.ensureNotDisposed();
        const fileName = resolveFileName(file);
        const path = this.toPath(fileName);
        const parseCacheKey = this.parseOptionsKey;

        // Check cache first
        const cached = this.sourceFileCache.get(path, parseCacheKey);
        if (cached) {
            this.sourceFileCache.retain(path, this.id);
            return cached;
        }

        // Fetch from server
        const response = await this.client.apiRequest<SourceFileResponse | undefined>("getSourceFile", {
            project: this.id,
            fileName,
        });
        if (!response?.data) {
            return undefined;
        }

        // Decode base64 to Uint8Array and create RemoteSourceFile
        const binaryData = Uint8Array.from(atob(response.data), c => c.charCodeAt(0));
        const sourceFile = new RemoteSourceFile(binaryData, this.decoder) as unknown as SourceFile;
        this.sourceFileCache.set(path, sourceFile, parseCacheKey, this.id);

        return sourceFile;
    }

    getSymbolAtLocation(node: Node): Promise<Symbol | undefined>;
    getSymbolAtLocation(nodes: readonly Node[]): Promise<(Symbol | undefined)[]>;
    async getSymbolAtLocation(nodeOrNodes: Node | readonly Node[]): Promise<Symbol | (Symbol | undefined)[] | undefined> {
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

    getSymbolAtPosition(file: DocumentIdentifier | string, position: number): Promise<Symbol | undefined>;
    getSymbolAtPosition(file: DocumentIdentifier | string, positions: readonly number[]): Promise<(Symbol | undefined)[]>;
    async getSymbolAtPosition(file: DocumentIdentifier | string, positionOrPositions: number | readonly number[]): Promise<Symbol | (Symbol | undefined)[] | undefined> {
        this.ensureNotDisposed();
        const fileName = resolveFileName(file);
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

    getTypeOfSymbol(symbol: Symbol): Promise<Type | undefined>;
    getTypeOfSymbol(symbols: readonly Symbol[]): Promise<(Type | undefined)[]>;
    async getTypeOfSymbol(symbolOrSymbols: Symbol | readonly Symbol[]): Promise<Type | (Type | undefined)[] | undefined> {
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
            symbol: (symbolOrSymbols as Symbol).ensureNotDisposed().id,
        });
        return data ? this.objectRegistry.getType(data) : undefined;
    }

    /**
     * Resolve a name to a symbol at a given location.
     * @param name The name to resolve
     * @param meaning Symbol flags indicating what kind of symbol to look for
     * @param location Optional node or document position for location context
     * @param excludeGlobals Whether to exclude global symbols
     */
    async resolveName(
        name: string,
        meaning: SymbolFlags,
        location?: Node | DocumentPosition,
        excludeGlobals?: boolean,
    ): Promise<Symbol | undefined> {
        this.ensureNotDisposed();
        // Distinguish Node (has `id`) from DocumentPosition (has `document` and `position`)
        const isNode = location && "id" in location;
        const data = await this.client.apiRequest<SymbolResponse | null>("resolveName", {
            project: this.id,
            name,
            meaning,
            location: isNode ? (location as Node).id : undefined,
            fileName: !isNode && location ? resolveFileName((location as DocumentPosition).document) : undefined,
            position: !isNode && location ? (location as DocumentPosition).position : undefined,
            excludeGlobals,
        });
        return data ? this.objectRegistry.getSymbol(data) : undefined;
    }
}

export class NodeHandle implements BaseNodeHandle<true> {
    readonly kind: SyntaxKind;
    readonly pos: number;
    readonly end: number;
    readonly path: Path;

    constructor(handle: string) {
        const parsed = parseNodeHandle(handle);
        this.pos = parsed.pos;
        this.end = parsed.end;
        this.kind = parsed.kind;
        this.path = parsed.path;
    }

    /**
     * Resolve this handle to the actual AST node by fetching the source file
     * from the given project and finding the node at the stored position.
     */
    async resolve(project: Project): Promise<Node | undefined> {
        const sourceFile = await project.getSourceFile(this.path);
        if (!sourceFile) {
            return undefined;
        }
        // Find the node at the stored position with matching kind and end
        return findDescendant(sourceFile, this.pos, this.end, this.kind);
    }
}

export class Symbol extends DisposableObject implements BaseSymbol<true> {
    readonly id: string;
    readonly name: string;
    readonly flags: SymbolFlags;
    readonly checkFlags: number;
    readonly declarations: readonly NodeHandle[];
    readonly valueDeclaration: NodeHandle | undefined;

    constructor(objectRegistry: AsyncObjectRegistry, data: SymbolResponse) {
        super(objectRegistry);
        this.id = data.id;
        this.name = data.name;
        this.flags = data.flags;
        this.checkFlags = data.checkFlags;
        this.declarations = (data.declarations ?? []).map(d => new NodeHandle(d));
        this.valueDeclaration = data.valueDeclaration ? new NodeHandle(data.valueDeclaration) : undefined;
    }
}

export class Type extends DisposableObject implements BaseType<true> {
    readonly id: string;
    readonly flags: TypeFlags;

    constructor(objectRegistry: AsyncObjectRegistry, data: TypeResponse) {
        super(objectRegistry);
        this.id = data.id;
        this.flags = data.flags;
    }
}
