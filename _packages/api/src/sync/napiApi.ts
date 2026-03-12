/**
 * NAPI-backed TypeScript API.
 *
 * This module provides the same API surface as `@typescript/api/sync`,
 * but uses a native Node.js addon (NAPI) instead of spawning a child
 * process. This eliminates IPC overhead and provides better performance
 * for in-process use cases.
 *
 * Usage:
 *   import { API } from "@typescript/api/napi";
 *   const api = new API({ napiModulePath: "./tsgo.node" });
 */

export { NapiClient, type NapiClientOptions } from "./napiClient.ts";

// Re-export all types from the sync API
export { Checker, type ClientSocketOptions, type ClientSpawnOptions, type DocumentIdentifier, type DocumentPosition, ElementFlags, Emitter, type LSPConnectionOptions, NodeHandle, ObjectFlags, Program, Project, SignatureFlags, SignatureKind, Snapshot, SymbolFlags, TypeFlags, TypePredicateKind } from "./api.ts";

export type { APIOptions, AssertsIdentifierTypePredicate, AssertsThisTypePredicate, ConditionalType, IdentifierTypePredicate, IndexedAccessType, IndexInfo, IndexType, InterfaceType, IntersectionType, LiteralType, ObjectType, StringMappingType, SubstitutionType, TemplateLiteralType, ThisTypePredicate, TupleType, Type, TypeParameter, TypePredicate, TypePredicateBase, TypeReference, UnionOrIntersectionType, UnionType } from "./api.ts";

export { documentURIToFileName, fileNameToDocumentURI } from "../path.ts";

import type {
    Path,
    SourceFile,
} from "@typescript/ast";
import { RemoteSourceFile } from "../node/node.ts";
import {
    readParseOptionsKey,
    readSourceFileHash,
} from "../node/node.ts";
import { ObjectRegistry } from "../objectRegistry.ts";
import {
    createGetCanonicalFileName,
    toPath,
} from "../path.ts";
import type {
    ConfigResponse,
    DocumentIdentifier,
    InitializeResponse,
    ProjectResponse,
    UpdateSnapshotParams,
    UpdateSnapshotResponse,
} from "../proto.ts";
import { resolveFileName } from "../proto.ts";
import { SourceFileCache } from "../sourceFileCache.ts";
import {
    Checker,
    Emitter,
    NodeHandle,
    Program,
    Project,
    Snapshot,
} from "./api.ts";
import {
    NapiClient,
    type NapiClientOptions,
} from "./napiClient.ts";

/**
 * NAPI-backed TypeScript API.
 *
 * Drop-in replacement for the sync `API` class that uses a native Node.js
 * addon instead of IPC. Provides the same interface but with better performance
 * for in-process use cases.
 */
export class API {
    private client: NapiClient;
    private sourceFileCache: SourceFileCache;
    private toPath: ((fileName: string) => Path) | undefined;
    private initialized: boolean = false;
    private activeSnapshots: Set<Snapshot> = new Set();
    private latestSnapshot: Snapshot | undefined;

    constructor(options: NapiClientOptions) {
        this.client = new NapiClient(options);
        this.sourceFileCache = new SourceFileCache();
    }

    private ensureInitialized(): void {
        if (!this.initialized) {
            const response = this.client.apiRequest<InitializeResponse>("initialize", null);
            const getCanonicalFileName = createGetCanonicalFileName(response.useCaseSensitiveFileNames);
            const currentDirectory = response.currentDirectory;
            this.toPath = (fileName: string) => toPath(fileName, currentDirectory, getCanonicalFileName) as Path;
            this.initialized = true;
        }
    }

    parseConfigFile(file: DocumentIdentifier): ConfigResponse {
        this.ensureInitialized();
        return this.client.apiRequest<ConfigResponse>("parseConfigFile", { file });
    }

    updateSnapshot(params?: UpdateSnapshotParams): Snapshot {
        this.ensureInitialized();

        const requestParams: UpdateSnapshotParams = params ?? {};
        if (requestParams.openProject) {
            requestParams.openProject = resolveFileName(requestParams.openProject);
        }

        const data = this.client.apiRequest<UpdateSnapshotResponse>("updateSnapshot", requestParams);

        // Retain cached source files from previous snapshot for unchanged files
        if (this.latestSnapshot) {
            this.sourceFileCache.retainForSnapshot(data.snapshot, this.latestSnapshot.id, data.changes);
            if (this.latestSnapshot.isDisposed()) {
                this.sourceFileCache.releaseSnapshot(this.latestSnapshot.id);
            }
        }

        // The NapiClient is structurally compatible with the sync Client
        const client = this.client as any;
        const snapshot = new Snapshot(
            data,
            client,
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

    close(): void {
        // Dispose all active snapshots
        for (const snapshot of [...this.activeSnapshots]) {
            snapshot.dispose();
        }
        // Release the latest snapshot's cache refs if still held
        if (this.latestSnapshot) {
            this.sourceFileCache.releaseSnapshot(this.latestSnapshot.id);
            this.latestSnapshot = undefined;
        }
        this.client.close();
        this.sourceFileCache.clear();
    }

    clearSourceFileCache(): void {
        this.sourceFileCache.clear();
    }
}
