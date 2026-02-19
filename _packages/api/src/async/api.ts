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
    type Checker as BaseChecker,
    type DocumentIdentifier,
    type DocumentPosition,
    type NodeHandle as BaseNodeHandle,
    type Program as BaseProgram,
    type Project as BaseProject,
    resolveFileName,
    type Snapshot as BaseSnapshot,
    type Symbol as BaseSymbol,
    type Type as BaseType,
} from "../base/api.ts";
import { ObjectRegistry } from "../base/objectRegistry.ts";
import { SourceFileCache } from "../base/sourceFileCache.ts";
import {
    findDescendant,
    parseNodeHandle,
    readParseOptionsKey,
    readSourceFileHash,
    RemoteSourceFile,
} from "../node.ts";
import {
    createGetCanonicalFileName,
    toPath,
} from "../path.ts";
import type {
    ConfigResponse,
    InitializeResponse,
    ProjectResponse,
    SourceFileResponse,
    SymbolResponse,
    TypeResponse,
    UpdateSnapshotResponse,
} from "../proto.ts";
import type {
    LSPUpdateSnapshotParams,
    UpdateSnapshotParams,
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

/** Type alias for the snapshot-scoped object registry */
type SnapshotObjectRegistry = ObjectRegistry<Symbol, Type>;

export class API<FromLSP extends boolean = false> implements BaseAPI<true, FromLSP> {
    private client: Client;
    private sourceFileCache: SourceFileCache;
    private toPath: ((fileName: string) => Path) | undefined;
    private initialized: boolean = false;
    private activeSnapshots: Set<Snapshot> = new Set();
    private latestSnapshot: Snapshot | undefined;

    constructor(options: APIOptions | LSPConnectionOptions) {
        this.client = new Client(options);
        this.sourceFileCache = new SourceFileCache();
    }

    /**
     * Create an API instance from an existing LSP connection's API session.
     * Use this when connecting to an API pipe provided by an LSP server via custom/initializeAPISession.
     */
    static async fromLSPConnection(options: LSPConnectionOptions): Promise<API<true>> {
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
            this.initialized = true;
        }
    }

    async parseConfigFile(file: DocumentIdentifier): Promise<ConfigResponse> {
        await this.ensureInitialized();
        return this.client.apiRequest<ConfigResponse>("parseConfigFile", { file });
    }

    async updateSnapshot(params?: FromLSP extends true ? LSPUpdateSnapshotParams : UpdateSnapshotParams): Promise<Snapshot> {
        await this.ensureInitialized();

        const requestParams: UpdateSnapshotParams = params ?? {};
        if (requestParams.openProject) {
            requestParams.openProject = resolveFileName(requestParams.openProject);
        }

        const data = await this.client.apiRequest<UpdateSnapshotResponse>("updateSnapshot", requestParams);

        // Retain cached source files from previous snapshot for unchanged files
        if (this.latestSnapshot) {
            this.sourceFileCache.retainForSnapshot(data.snapshot, this.latestSnapshot.id, data.changes);
            if (this.latestSnapshot.isDisposed()) {
                this.sourceFileCache.releaseSnapshot(this.latestSnapshot.id);
            }
        }

        const snapshot = new Snapshot(
            data,
            this.client,
            this.sourceFileCache,
            this.toPath!,
            () => {
                this.activeSnapshots.delete(snapshot);
                if (snapshot !== this.latestSnapshot) {
                    this.sourceFileCache.releaseSnapshot(snapshot.id);
                }
            },
        );
        this.latestSnapshot = snapshot;
        this.activeSnapshots.add(snapshot);

        return snapshot;
    }

    async close(): Promise<void> {
        // Dispose all active snapshots
        for (const snapshot of [...this.activeSnapshots]) {
            await snapshot.dispose();
        }
        // Release the latest snapshot's cache refs if still held
        if (this.latestSnapshot) {
            this.sourceFileCache.releaseSnapshot(this.latestSnapshot.id);
            this.latestSnapshot = undefined;
        }
        await this.client.close();
        this.sourceFileCache.clear();
    }
}

export class Snapshot implements BaseSnapshot<true> {
    readonly id: string;
    private projectMap: Map<Path, Project>;
    private toPath: (fileName: string) => Path;
    private client: Client;
    private objectRegistry: SnapshotObjectRegistry;
    private disposed: boolean = false;
    private onDispose: () => void;

    constructor(
        data: UpdateSnapshotResponse,
        client: Client,
        sourceFileCache: SourceFileCache,
        toPath: (fileName: string) => Path,
        onDispose: () => void,
    ) {
        this.id = data.snapshot;
        this.client = client;
        this.toPath = toPath;
        this.onDispose = onDispose;

        this.objectRegistry = new ObjectRegistry<Symbol, Type>({
            createSymbol: symbolData => new Symbol(symbolData),
            createType: typeData => new Type(typeData),
        });

        // Create projects
        this.projectMap = new Map();
        for (const projData of data.projects) {
            const project = new Project(projData, this.id, client, this.objectRegistry, sourceFileCache, toPath);
            this.projectMap.set(toPath(projData.configFileName), project);
        }
    }

    getProjects(): readonly Project[] {
        this.ensureNotDisposed();
        return [...this.projectMap.values()];
    }

    getProject(configFileName: string): Project | undefined {
        this.ensureNotDisposed();
        return this.projectMap.get(this.toPath(configFileName));
    }

    async getDefaultProjectForFile(file: DocumentIdentifier): Promise<Project | undefined> {
        this.ensureNotDisposed();
        const data = await this.client.apiRequest<ProjectResponse | null>("getDefaultProjectForFile", {
            snapshot: this.id,
            file,
        });
        if (!data) return undefined;
        return this.projectMap.get(this.toPath(data.configFileName));
    }

    [globalThis.Symbol.dispose](): void {
        this.dispose();
    }

    async dispose(): Promise<void> {
        if (this.disposed) return;
        this.disposed = true;
        this.objectRegistry.clear();
        this.onDispose();
        await this.client.apiRequest("release", { handle: this.id });
    }

    isDisposed(): boolean {
        return this.disposed;
    }

    private ensureNotDisposed(): void {
        if (this.disposed) {
            throw new Error("Snapshot is disposed");
        }
    }
}

export class Project implements BaseProject<true> {
    readonly id: string;
    readonly configFileName: string;
    readonly compilerOptions: Record<string, unknown>;
    readonly rootFiles: readonly string[];

    readonly program: Program;
    readonly checker: Checker;

    constructor(
        data: ProjectResponse,
        snapshotId: string,
        client: Client,
        objectRegistry: SnapshotObjectRegistry,
        sourceFileCache: SourceFileCache,
        toPath: (fileName: string) => Path,
    ) {
        this.id = data.id;
        this.configFileName = data.configFileName;
        this.compilerOptions = data.compilerOptions;
        this.rootFiles = data.rootFiles;
        this.program = new Program(
            snapshotId,
            this.id,
            client,
            sourceFileCache,
            toPath,
        );
        this.checker = new Checker(
            snapshotId,
            this.id,
            client,
            objectRegistry,
        );
    }
}

export class Program implements BaseProgram<true> {
    private snapshotId: string;
    private projectId: string;
    private client: Client;
    private sourceFileCache: SourceFileCache;
    private toPath: (fileName: string) => Path;
    private decoder = new TextDecoder();

    constructor(
        snapshotId: string,
        projectId: string,
        client: Client,
        sourceFileCache: SourceFileCache,
        toPath: (fileName: string) => Path,
    ) {
        this.snapshotId = snapshotId;
        this.projectId = projectId;
        this.client = client;
        this.sourceFileCache = sourceFileCache;
        this.toPath = toPath;
    }

    async getSourceFile(file: DocumentIdentifier): Promise<SourceFile | undefined> {
        const fileName = resolveFileName(file);
        const path = this.toPath(fileName);

        // Check if we already have a retained cache entry for this (snapshot, project) pair
        const retained = this.sourceFileCache.getRetained(path, this.snapshotId, this.projectId);
        if (retained) {
            return retained;
        }

        // Fetch from server
        const response = await this.client.apiRequest<SourceFileResponse | undefined>("getSourceFile", {
            snapshot: this.snapshotId,
            project: this.projectId,
            file,
        });
        if (!response?.data) {
            return undefined;
        }

        // Decode base64 to Uint8Array
        const binaryData = Uint8Array.from(atob(response.data), c => c.charCodeAt(0));
        const view = new DataView(binaryData.buffer, binaryData.byteOffset, binaryData.byteLength);
        const contentHash = readSourceFileHash(view);
        const parseOptionsKey = readParseOptionsKey(view);

        // Create a new RemoteSourceFile and cache it (set returns existing if hash matches)
        const sourceFile = new RemoteSourceFile(binaryData, this.decoder) as unknown as SourceFile;
        return this.sourceFileCache.set(path, sourceFile, parseOptionsKey, contentHash, this.snapshotId, this.projectId);
    }
}

export class Checker implements BaseChecker<true> {
    private snapshotId: string;
    private projectId: string;
    private client: Client;
    private objectRegistry: SnapshotObjectRegistry;

    constructor(
        snapshotId: string,
        projectId: string,
        client: Client,
        objectRegistry: SnapshotObjectRegistry,
    ) {
        this.snapshotId = snapshotId;
        this.projectId = projectId;
        this.client = client;
        this.objectRegistry = objectRegistry;
    }

    getSymbolAtLocation(node: Node): Promise<Symbol | undefined>;
    getSymbolAtLocation(nodes: readonly Node[]): Promise<(Symbol | undefined)[]>;
    async getSymbolAtLocation(nodeOrNodes: Node | readonly Node[]): Promise<Symbol | (Symbol | undefined)[] | undefined> {
        if (Array.isArray(nodeOrNodes)) {
            const data = await this.client.apiRequest<(SymbolResponse | null)[]>("getSymbolsAtLocations", {
                snapshot: this.snapshotId,
                project: this.projectId,
                locations: nodeOrNodes.map(node => node.id),
            });
            return data.map(d => d ? this.objectRegistry.getSymbol(d) : undefined);
        }
        const data = await this.client.apiRequest<SymbolResponse | null>("getSymbolAtLocation", {
            snapshot: this.snapshotId,
            project: this.projectId,
            location: (nodeOrNodes as Node).id,
        });
        return data ? this.objectRegistry.getSymbol(data) : undefined;
    }

    getSymbolAtPosition(file: DocumentIdentifier, position: number): Promise<Symbol | undefined>;
    getSymbolAtPosition(file: DocumentIdentifier, positions: readonly number[]): Promise<(Symbol | undefined)[]>;
    async getSymbolAtPosition(file: DocumentIdentifier, positionOrPositions: number | readonly number[]): Promise<Symbol | (Symbol | undefined)[] | undefined> {
        if (typeof positionOrPositions === "number") {
            const data = await this.client.apiRequest<SymbolResponse | null>("getSymbolAtPosition", {
                snapshot: this.snapshotId,
                project: this.projectId,
                file,
                position: positionOrPositions,
            });
            return data ? this.objectRegistry.getSymbol(data) : undefined;
        }
        const data = await this.client.apiRequest<(SymbolResponse | null)[]>("getSymbolsAtPositions", {
            snapshot: this.snapshotId,
            project: this.projectId,
            file,
            positions: positionOrPositions,
        });
        return data.map(d => d ? this.objectRegistry.getSymbol(d) : undefined);
    }

    getTypeOfSymbol(symbol: Symbol): Promise<Type | undefined>;
    getTypeOfSymbol(symbols: readonly Symbol[]): Promise<(Type | undefined)[]>;
    async getTypeOfSymbol(symbolOrSymbols: Symbol | readonly Symbol[]): Promise<Type | (Type | undefined)[] | undefined> {
        if (Array.isArray(symbolOrSymbols)) {
            const data = await this.client.apiRequest<(TypeResponse | null)[]>("getTypesOfSymbols", {
                snapshot: this.snapshotId,
                project: this.projectId,
                symbols: symbolOrSymbols.map(s => s.id),
            });
            return data.map(d => d ? this.objectRegistry.getType(d) : undefined);
        }
        const data = await this.client.apiRequest<TypeResponse | null>("getTypeOfSymbol", {
            snapshot: this.snapshotId,
            project: this.projectId,
            symbol: (symbolOrSymbols as Symbol).id,
        });
        return data ? this.objectRegistry.getType(data) : undefined;
    }

    async resolveName(
        name: string,
        meaning: SymbolFlags,
        location?: Node | DocumentPosition,
        excludeGlobals?: boolean,
    ): Promise<Symbol | undefined> {
        // Distinguish Node (has `id`) from DocumentPosition (has `document` and `position`)
        const isNode = location && "id" in location;
        const data = await this.client.apiRequest<SymbolResponse | null>("resolveName", {
            snapshot: this.snapshotId,
            project: this.projectId,
            name,
            meaning,
            location: isNode ? (location as Node).id : undefined,
            file: !isNode && location ? (location as DocumentPosition).document : undefined,
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
        const sourceFile = await project.program.getSourceFile(this.path);
        if (!sourceFile) {
            return undefined;
        }
        // Find the node at the stored position with matching kind and end
        return findDescendant(sourceFile, this.pos, this.end, this.kind);
    }
}

export class Symbol implements BaseSymbol<true> {
    readonly id: string;
    readonly name: string;
    readonly flags: SymbolFlags;
    readonly checkFlags: number;
    readonly declarations: readonly NodeHandle[];
    readonly valueDeclaration: NodeHandle | undefined;

    constructor(data: SymbolResponse) {
        this.id = data.id;
        this.name = data.name;
        this.flags = data.flags;
        this.checkFlags = data.checkFlags;
        this.declarations = (data.declarations ?? []).map(d => new NodeHandle(d));
        this.valueDeclaration = data.valueDeclaration ? new NodeHandle(data.valueDeclaration) : undefined;
    }
}

export class Type implements BaseType<true> {
    readonly id: string;
    readonly flags: TypeFlags;

    constructor(data: TypeResponse) {
        this.id = data.id;
        this.flags = data.flags;
    }
}
