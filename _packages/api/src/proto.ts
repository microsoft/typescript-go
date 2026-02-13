/**
 * Response from the initialize method.
 */
export interface InitializeResponse {
    /** Whether the host file system is case-sensitive */
    useCaseSensitiveFileNames: boolean;
    /** The server's current working directory */
    currentDirectory: string;
}

export interface ConfigResponse {
    options: Record<string, unknown>;
    fileNames: string[];
}

/**
 * Parameters for updateSnapshot.
 */
export interface UpdateSnapshotParams {
    /** Path to a tsconfig.json file to open in the new snapshot */
    openProject?: string;
}

/**
 * Changes to source files within a single project.
 */
export interface ProjectFileChanges {
    /** Source file paths whose content changed */
    changedFiles?: string[];
    /** Source file paths removed from the project's program */
    deletedFiles?: string[];
}

/**
 * Changes between two consecutive snapshots, reported per-project.
 */
export interface SnapshotChanges {
    /** Project handles mapped to their file changes. Projects not listed are unchanged. */
    changedProjects?: Record<string, ProjectFileChanges>;
    /** Project handles that were removed from the snapshot */
    removedProjects?: string[];
}

/**
 * Response from updateSnapshot.
 */
export interface UpdateSnapshotResponse {
    /** Handle for the newly created snapshot */
    snapshot: string;
    /** List of projects in the snapshot */
    projects: ProjectResponse[];
    /** Changes from the previous snapshot (absent for the first snapshot) */
    changes?: SnapshotChanges;
}

export interface ProjectResponse {
    id: string;
    configFileName: string;
    compilerOptions: Record<string, unknown>;
    rootFiles: string[];
    /** Encodes the source-file-independent parse options for this project (for cache keying) */
    parseOptionsKey: string;
}

export interface SourceFileResponse {
    /** Base64-encoded binary AST data */
    data: string;
}

export interface SymbolResponse {
    id: string;
    name: string;
    flags: number;
    checkFlags: number;
    declarations?: string[];
    valueDeclaration?: string;
}

export interface TypeResponse {
    id: string;
    flags: number;
}
