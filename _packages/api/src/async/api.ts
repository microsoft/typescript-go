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
    SnapshotChanges,
    SourceFileResponse,
    SymbolResponse,
    TypeResponse,
    UpdateSnapshotResponse,
} from "../proto.ts";
import type { UpdateSnapshotParams } from "../proto.ts";
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

export class API implements BaseAPI<true> {
    private client: Client;
    private sourceFileCache: SourceFileCache;
    private toPath: ((fileName: string) => Path) | undefined;
    private initialized: boolean = false;
    private activeSnapshots: Set<Snapshot> = new Set();

    constructor(options: APIOptions | LSPConnectionOptions) {
        this.client = new Client(options);
        this.sourceFileCache = new SourceFileCache();
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
            this.initialized = true;
        }
    }

    async parseConfigFile(file: DocumentIdentifier | string): Promise<ConfigResponse> {
        await this.ensureInitialized();
        return this.client.apiRequest<ConfigResponse>("parseConfigFile", { fileName: resolveFileName(file) });
    }

    async updateSnapshot(params?: UpdateSnapshotParams): Promise<Snapshot> {
        await this.ensureInitialized();

        const requestParams: UpdateSnapshotParams = params ?? {};
        if (requestParams.openProject) {
            requestParams.openProject = resolveFileName(requestParams.openProject);
        }

        const data = await this.client.apiRequest<UpdateSnapshotResponse>("updateSnapshot", requestParams);
        const snapshot = new Snapshot(
            data,
            this.client,
            this.sourceFileCache,
            this.toPath!,
            () => this.activeSnapshots.delete(snapshot),
        );
        this.activeSnapshots.add(snapshot);

        // Apply cache invalidation if changes are available
        if (data.changes) {
            this.applySnapshotChanges(data.changes);
        }

        return snapshot;
    }

    private applySnapshotChanges(changes: SnapshotChanges): void {
        // Handle removed projects - release their cached files
        if (changes.removedProjects) {
            for (const projectId of changes.removedProjects) {
                this.sourceFileCache.releaseProject(projectId);
            }
        }

        // Handle file changes within projects
        if (changes.projectChanges) {
            for (const [_projectId, projectChanges] of Object.entries(changes.projectChanges)) {
                if (projectChanges.changedFiles) {
                    this.sourceFileCache.remove(projectChanges.changedFiles as Path[]);
                }
                if (projectChanges.removedFiles) {
                    this.sourceFileCache.remove(projectChanges.removedFiles as Path[]);
                }
            }
        }
    }

    async close(): Promise<void> {
        // Dispose all active snapshots
        for (const snapshot of [...this.activeSnapshots]) {
            await snapshot.dispose();
        }
        await this.client.close();
        this.sourceFileCache.clear();
    }
}

export class Snapshot implements BaseSnapshot<true> {
    readonly id: string;
    readonly projects: readonly Project[];
    readonly changes: SnapshotChanges | undefined;
    private projectMap: Map<string, Project>;
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
        this.changes = data.changes;
        this.client = client;
        this.onDispose = onDispose;

        this.objectRegistry = new ObjectRegistry<Symbol, Type>({
            createSymbol: symbolData => new Symbol(symbolData),
            createType: typeData => new Type(typeData),
        });

        // Create projects
        this.projectMap = new Map();
        const projects: Project[] = [];
        for (const projData of data.projects) {
            const project = new Project(projData, this.id, client, this.objectRegistry, sourceFileCache, toPath);
            this.projectMap.set(projData.id, project);
            projects.push(project);
        }
        this.projects = projects;
    }

    getProject(id: string): Project | undefined {
        this.ensureNotDisposed();
        return this.projectMap.get(id);
    }

    async getDefaultProjectForFile(file: DocumentIdentifier | string): Promise<Project | undefined> {
        this.ensureNotDisposed();
        const data = await this.client.apiRequest<ProjectResponse | null>("getDefaultProjectForFile", {
            snapshot: this.id,
            fileName: resolveFileName(file),
        });
        if (!data) return undefined;
        return this.projectMap.get(data.id);
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
            data.parseOptionsKey,
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
    private parseOptionsKey: string;
    private decoder = new TextDecoder();

    constructor(
        snapshotId: string,
        projectId: string,
        client: Client,
        sourceFileCache: SourceFileCache,
        toPath: (fileName: string) => Path,
        parseOptionsKey: string,
    ) {
        this.snapshotId = snapshotId;
        this.projectId = projectId;
        this.client = client;
        this.sourceFileCache = sourceFileCache;
        this.toPath = toPath;
        this.parseOptionsKey = parseOptionsKey;
    }

    async getSourceFile(file: DocumentIdentifier | string): Promise<SourceFile | undefined> {
        const fileName = resolveFileName(file);
        const path = this.toPath(fileName);

        // Check cache first
        const cached = this.sourceFileCache.get(path, this.parseOptionsKey);
        if (cached) {
            this.sourceFileCache.retain(path, this.projectId);
            return cached;
        }

        // Fetch from server
        const response = await this.client.apiRequest<SourceFileResponse | undefined>("getSourceFile", {
            snapshot: this.snapshotId,
            project: this.projectId,
            fileName,
        });
        if (!response?.data) {
            return undefined;
        }

        // Decode base64 to Uint8Array and create RemoteSourceFile
        const binaryData = Uint8Array.from(atob(response.data), c => c.charCodeAt(0));
        const sourceFile = new RemoteSourceFile(binaryData, this.decoder) as unknown as SourceFile;
        this.sourceFileCache.set(path, sourceFile, this.parseOptionsKey, this.projectId);

        return sourceFile;
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

    getSymbolAtPosition(file: DocumentIdentifier | string, position: number): Promise<Symbol | undefined>;
    getSymbolAtPosition(file: DocumentIdentifier | string, positions: readonly number[]): Promise<(Symbol | undefined)[]>;
    async getSymbolAtPosition(file: DocumentIdentifier | string, positionOrPositions: number | readonly number[]): Promise<Symbol | (Symbol | undefined)[] | undefined> {
        const fileName = resolveFileName(file);
        if (typeof positionOrPositions === "number") {
            const data = await this.client.apiRequest<SymbolResponse | null>("getSymbolAtPosition", {
                snapshot: this.snapshotId,
                project: this.projectId,
                fileName,
                position: positionOrPositions,
            });
            return data ? this.objectRegistry.getSymbol(data) : undefined;
        }
        const data = await this.client.apiRequest<(SymbolResponse | null)[]>("getSymbolsAtPositions", {
            snapshot: this.snapshotId,
            project: this.projectId,
            fileName,
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
