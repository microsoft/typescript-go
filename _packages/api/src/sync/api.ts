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
    type APIOptions as BaseAPIOptions,
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
import type { FileSystem } from "../fs.ts";
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
    SymbolResponse,
    TypeResponse,
    UpdateSnapshotResponse,
} from "../proto.ts";
import type { UpdateSnapshotParams } from "../proto.ts";
import { Client } from "./client.ts";

export { SymbolFlags, TypeFlags };
export type { DocumentIdentifier, DocumentPosition };
export { documentURIToFileName, fileNameToDocumentURI } from "../path.ts";

export interface APIOptions extends BaseAPIOptions {
    fs?: FileSystem;
}

/** Type alias for the snapshot-scoped object registry */
type SnapshotObjectRegistry = ObjectRegistry<Symbol, Type>;

export class API implements BaseAPI<false> {
    /** @internal */
    readonly client: Client;
    private sourceFileCache: SourceFileCache;
    private useCaseSensitiveFileNames: boolean;
    private toPath: (fileName: string) => Path;
    private activeSnapshots: Set<Snapshot> = new Set();

    constructor(options: APIOptions) {
        this.client = new Client(options);
        this.sourceFileCache = new SourceFileCache();

        // Initialize and get file system settings
        const initResponse: InitializeResponse = this.client.request("initialize", null);
        this.useCaseSensitiveFileNames = initResponse.useCaseSensitiveFileNames;

        // Create the toPath function using the server's current directory and case sensitivity
        const getCanonicalFileName = createGetCanonicalFileName(this.useCaseSensitiveFileNames);
        const currentDirectory = initResponse.currentDirectory;
        this.toPath = (fileName: string) => toPath(fileName, currentDirectory, getCanonicalFileName) as Path;
    }

    parseConfigFile(file: DocumentIdentifier | string): ConfigResponse {
        return this.client.request("parseConfigFile", { fileName: resolveFileName(file) });
    }

    updateSnapshot(params?: UpdateSnapshotParams): Snapshot {
        const requestParams: UpdateSnapshotParams = params ?? {};
        if (requestParams.openProject) {
            requestParams.openProject = resolveFileName(requestParams.openProject);
        }

        const data: UpdateSnapshotResponse = this.client.request("updateSnapshot", requestParams);
        const snapshot = new Snapshot(
            data,
            this.client,
            this.sourceFileCache,
            this.toPath,
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

    echo(message: string): string {
        return this.client.echo(message);
    }

    echoBinary(message: Uint8Array): Uint8Array {
        return this.client.echoBinary(message);
    }

    close(): void {
        // Dispose all active snapshots
        for (const snapshot of [...this.activeSnapshots]) {
            snapshot.dispose();
        }
        this.client.close();
        this.sourceFileCache.clear();
    }
}

export class Snapshot implements BaseSnapshot<false> {
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

    getDefaultProjectForFile(file: DocumentIdentifier | string): Project | undefined {
        this.ensureNotDisposed();
        const data: ProjectResponse | null = this.client.request("getDefaultProjectForFile", {
            snapshot: this.id,
            fileName: resolveFileName(file),
        });
        if (!data) return undefined;
        return this.projectMap.get(data.id);
    }

    [globalThis.Symbol.dispose](): void {
        this.dispose();
    }

    dispose(): void {
        if (this.disposed) return;
        this.disposed = true;
        this.objectRegistry.clear();
        this.onDispose();
        this.client.request("release", { handle: this.id });
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

export class Project implements BaseProject<false> {
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

export class Program implements BaseProgram<false> {
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

    getSourceFile(file: DocumentIdentifier | string): SourceFile | undefined {
        const fileName = resolveFileName(file);
        const path = this.toPath(fileName);

        // Check cache first
        const cached = this.sourceFileCache.get(path, this.parseOptionsKey);
        if (cached) {
            this.sourceFileCache.retain(path, this.projectId);
            return cached;
        }

        // Fetch from server
        const response: Uint8Array | undefined = this.client.requestBinary("getSourceFile", {
            snapshot: this.snapshotId,
            project: this.projectId,
            fileName,
        });
        if (!response) {
            return undefined;
        }

        const sourceFile = new RemoteSourceFile(response, this.decoder) as unknown as SourceFile;
        this.sourceFileCache.set(path, sourceFile, this.parseOptionsKey, this.projectId);

        return sourceFile;
    }
}

export class Checker implements BaseChecker<false> {
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

    getSymbolAtLocation(node: Node): Symbol | undefined;
    getSymbolAtLocation(nodes: readonly Node[]): (Symbol | undefined)[];
    getSymbolAtLocation(nodeOrNodes: Node | readonly Node[]): Symbol | (Symbol | undefined)[] | undefined {
        if (Array.isArray(nodeOrNodes)) {
            const data = this.client.request("getSymbolsAtLocations", { snapshot: this.snapshotId, project: this.projectId, locations: nodeOrNodes.map(node => node.id) });
            return data.map((d: SymbolResponse | null) => d ? this.objectRegistry.getSymbol(d) : undefined);
        }
        const data = this.client.request("getSymbolAtLocation", { snapshot: this.snapshotId, project: this.projectId, location: (nodeOrNodes as Node).id });
        return data ? this.objectRegistry.getSymbol(data) : undefined;
    }

    getSymbolAtPosition(file: DocumentIdentifier | string, position: number): Symbol | undefined;
    getSymbolAtPosition(file: DocumentIdentifier | string, positions: readonly number[]): (Symbol | undefined)[];
    getSymbolAtPosition(file: DocumentIdentifier | string, positionOrPositions: number | readonly number[]): Symbol | (Symbol | undefined)[] | undefined {
        const fileName = resolveFileName(file);
        if (typeof positionOrPositions === "number") {
            const data = this.client.request("getSymbolAtPosition", { snapshot: this.snapshotId, project: this.projectId, fileName, position: positionOrPositions });
            return data ? this.objectRegistry.getSymbol(data) : undefined;
        }
        const data = this.client.request("getSymbolsAtPositions", { snapshot: this.snapshotId, project: this.projectId, fileName, positions: positionOrPositions });
        return data.map((d: SymbolResponse | null) => d ? this.objectRegistry.getSymbol(d) : undefined);
    }

    getTypeOfSymbol(symbol: Symbol): Type | undefined;
    getTypeOfSymbol(symbols: readonly Symbol[]): (Type | undefined)[];
    getTypeOfSymbol(symbolOrSymbols: Symbol | readonly Symbol[]): Type | (Type | undefined)[] | undefined {
        if (Array.isArray(symbolOrSymbols)) {
            const data = this.client.request("getTypesOfSymbols", { snapshot: this.snapshotId, project: this.projectId, symbols: symbolOrSymbols.map(s => s.id) });
            return data.map((d: TypeResponse | null) => d ? this.objectRegistry.getType(d) : undefined);
        }
        const data = this.client.request("getTypeOfSymbol", { snapshot: this.snapshotId, project: this.projectId, symbol: (symbolOrSymbols as Symbol).id });
        return data ? this.objectRegistry.getType(data) : undefined;
    }

    resolveName(
        name: string,
        meaning: SymbolFlags,
        location?: Node | DocumentPosition,
        excludeGlobals?: boolean,
    ): Symbol | undefined {
        // Distinguish Node (has `id`) from DocumentPosition (has `document` and `position`)
        const isNode = location && "id" in location;
        const data = this.client.request("resolveName", {
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

export class NodeHandle implements BaseNodeHandle<false> {
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
    resolve(project: Project): Node | undefined {
        const sourceFile = project.program.getSourceFile(this.path);
        if (!sourceFile) {
            return undefined;
        }
        // Find the node at the stored position with matching kind and end
        return findDescendant(sourceFile, this.pos, this.end, this.kind);
    }
}

export class Symbol implements BaseSymbol<false> {
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

export class Type implements BaseType<false> {
    readonly id: string;
    readonly flags: TypeFlags;

    constructor(data: TypeResponse) {
        this.id = data.id;
        this.flags = data.flags;
    }
}
