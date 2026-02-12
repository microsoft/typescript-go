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
    /** Handle of a previous snapshot for computing changes */
    previousSnapshot?: string;
}

/**
 * Response from updateSnapshot.
 */
export interface UpdateSnapshotResponse {
    /** Handle for the newly created snapshot */
    snapshot: string;
    /** List of projects in the snapshot */
    projects: ProjectResponse[];
    /** Changes relative to previousSnapshot, if provided */
    changes?: SnapshotChanges;
}

/**
 * Changes between two snapshots for cache invalidation.
 */
export interface SnapshotChanges {
    /** List of project handles that no longer exist */
    removedProjects?: string[];
    /** Map of project handles to their file changes */
    projectChanges?: Record<string, ProjectChanges>;
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

/**
 * Describes file changes within a project.
 */
export interface ProjectChanges {
    /** List of file paths whose content has changed */
    changedFiles?: string[];
    /** List of file paths that no longer exist in the project */
    removedFiles?: string[];
}
