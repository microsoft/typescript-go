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
    LoadProjectResponse,
    ProjectChanges,
    ProjectResponse,
    SymbolResponse,
    TypeResponse,
} from "../proto.ts";
import { Client } from "./client.ts";

export { SymbolFlags, TypeFlags };
export type { DocumentIdentifier, DocumentPosition };
export { documentURIToFileName, fileNameToDocumentURI } from "../path.ts";

export interface APIOptions extends BaseAPIOptions {
    fs?: FileSystem;
}

/** Type alias for the sync object registry */
type SyncObjectRegistry = ObjectRegistry<Project, Symbol, Type>;

export abstract class DisposableObject {
    private disposed: boolean = false;
    protected objectRegistry: SyncObjectRegistry;
    abstract readonly id: string;

    constructor(objectRegistry: SyncObjectRegistry) {
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

export class API implements BaseAPI<false> {
    private client: Client;
    private objectRegistry: SyncObjectRegistry;
    private sourceFileCache: SourceFileCache;
    private useCaseSensitiveFileNames: boolean;
    private toPath: (fileName: string) => Path;

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

        this.objectRegistry = new ObjectRegistry<Project, Symbol, Type>(
            {
                createProject: data => new Project(this.client, this.objectRegistry, this.sourceFileCache, this.toPath, data),
                createSymbol: data => new Symbol(this.objectRegistry, data),
                createType: data => new Type(this.objectRegistry, data),
            },
            id => this.client.request("release", id),
        );
    }

    parseConfigFile(file: DocumentIdentifier | string): ConfigResponse {
        return this.client.request("parseConfigFile", { fileName: resolveFileName(file) });
    }

    loadProject(configFile: DocumentIdentifier | string): Project {
        const data = this.client.request("loadProject", { configFileName: resolveFileName(configFile) });
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
        this.sourceFileCache.clear();
    }
}

export class Project extends DisposableObject implements BaseProject<false> {
    private decoder = new TextDecoder();
    private client: Client;
    private sourceFileCache: SourceFileCache;
    private toPath: (fileName: string) => Path;
    private parseOptionsKey!: string;

    readonly id: string;
    configFileName!: string;
    compilerOptions!: Record<string, unknown>;
    rootFiles!: readonly string[];

    constructor(client: Client, objectRegistry: SyncObjectRegistry, sourceFileCache: SourceFileCache, toPath: (fileName: string) => Path, data: ProjectResponse) {
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

    reload(): ProjectChanges | undefined {
        this.ensureNotDisposed();
        const data: LoadProjectResponse = this.client.request("loadProject", { configFileName: this.configFileName });
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

    getSourceFile(file: DocumentIdentifier | string): SourceFile | undefined {
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
        const response: Uint8Array | undefined = this.client.requestBinary("getSourceFile", {
            project: this.id,
            fileName,
        });
        if (!response) {
            return undefined;
        }

        const sourceFile = new RemoteSourceFile(response, this.decoder) as unknown as SourceFile;
        this.sourceFileCache.set(path, sourceFile, parseCacheKey, this.id);

        return sourceFile;
    }

    getSymbolAtLocation(node: Node): Symbol | undefined;
    getSymbolAtLocation(nodes: readonly Node[]): (Symbol | undefined)[];
    getSymbolAtLocation(nodeOrNodes: Node | readonly Node[]): Symbol | (Symbol | undefined)[] | undefined {
        this.ensureNotDisposed();
        if (Array.isArray(nodeOrNodes)) {
            const data = this.client.request("getSymbolsAtLocations", { project: this.id, locations: nodeOrNodes.map(node => node.id) });
            return data.map((d: SymbolResponse | null) => d ? this.objectRegistry.getSymbol(d) : undefined);
        }
        const data = this.client.request("getSymbolAtLocation", { project: this.id, location: (nodeOrNodes as Node).id });
        return data ? this.objectRegistry.getSymbol(data) : undefined;
    }

    getSymbolAtPosition(file: DocumentIdentifier | string, position: number): Symbol | undefined;
    getSymbolAtPosition(file: DocumentIdentifier | string, positions: readonly number[]): (Symbol | undefined)[];
    getSymbolAtPosition(file: DocumentIdentifier | string, positionOrPositions: number | readonly number[]): Symbol | (Symbol | undefined)[] | undefined {
        this.ensureNotDisposed();
        const fileName = resolveFileName(file);
        if (typeof positionOrPositions === "number") {
            const data = this.client.request("getSymbolAtPosition", { project: this.id, fileName, position: positionOrPositions });
            return data ? this.objectRegistry.getSymbol(data) : undefined;
        }
        const data = this.client.request("getSymbolsAtPositions", { project: this.id, fileName, positions: positionOrPositions });
        return data.map((d: SymbolResponse | null) => d ? this.objectRegistry.getSymbol(d) : undefined);
    }

    getTypeOfSymbol(symbol: Symbol): Type | undefined;
    getTypeOfSymbol(symbols: readonly Symbol[]): (Type | undefined)[];
    getTypeOfSymbol(symbolOrSymbols: Symbol | readonly Symbol[]): Type | (Type | undefined)[] | undefined {
        this.ensureNotDisposed();
        if (Array.isArray(symbolOrSymbols)) {
            const data = this.client.request("getTypesOfSymbols", { project: this.id, symbols: symbolOrSymbols.map(symbol => symbol.ensureNotDisposed().id) });
            return data.map((d: TypeResponse | null) => d ? this.objectRegistry.getType(d) : undefined);
        }
        const data = this.client.request("getTypeOfSymbol", { project: this.id, symbol: (symbolOrSymbols as Symbol).ensureNotDisposed().id });
        return data ? this.objectRegistry.getType(data) : undefined;
    }

    /**
     * Resolve a name to a symbol at a given location.
     * @param name The name to resolve
     * @param meaning Symbol flags indicating what kind of symbol to look for
     * @param location Optional node or document position for location context
     * @param excludeGlobals Whether to exclude global symbols
     */
    resolveName(
        name: string,
        meaning: SymbolFlags,
        location?: Node | DocumentPosition,
        excludeGlobals?: boolean,
    ): Symbol | undefined {
        this.ensureNotDisposed();
        // Distinguish Node (has `id`) from DocumentPosition (has `document` and `position`)
        const isNode = location && "id" in location;
        const data = this.client.request("resolveName", {
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
        const sourceFile = project.getSourceFile(this.path);
        if (!sourceFile) {
            return undefined;
        }
        // Find the node at the stored position with matching kind and end
        return findDescendant(sourceFile, this.pos, this.end, this.kind);
    }
}

export class Symbol extends DisposableObject implements BaseSymbol<false> {
    readonly id: string;
    readonly name: string;
    readonly flags: SymbolFlags;
    readonly checkFlags: number;
    readonly declarations: readonly NodeHandle[];
    readonly valueDeclaration: NodeHandle | undefined;

    constructor(objectRegistry: SyncObjectRegistry, data: SymbolResponse) {
        super(objectRegistry);
        this.id = data.id;
        this.name = data.name;
        this.flags = data.flags;
        this.checkFlags = data.checkFlags;
        this.declarations = (data.declarations ?? []).map(d => new NodeHandle(d));
        this.valueDeclaration = data.valueDeclaration ? new NodeHandle(data.valueDeclaration) : undefined;
    }
}

export class Type extends DisposableObject implements BaseType<false> {
    readonly id: string;
    readonly flags: TypeFlags;

    constructor(objectRegistry: SyncObjectRegistry, data: TypeResponse) {
        super(objectRegistry);
        this.id = data.id;
        this.flags = data.flags;
    }
}
