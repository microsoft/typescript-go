import type {
    Path,
    SourceFile,
} from "@typescript/ast";

/**
 * A cached source file entry.
 */
export interface CachedSourceFile {
    /** The cached source file object */
    file: SourceFile;
    /** The parse cache key that was used to create this file */
    parseCacheKey: string;
    /** Set of project IDs that are retaining this cache entry */
    retainingProjects: Set<string>;
}

/**
 * Client-side cache for source files.
 *
 * This cache stores source files by their path and allows multiple projects
 * to share cached entries when their parse options are compatible (same parseCacheKey).
 *
 * Lifecycle:
 * - Files are added when fetched from the server
 * - Projects retain cache entries to keep them alive
 * - Files are removed when the server reports changes (e.g., via adoptLSPState)
 * - Files are removed when no projects retain them
 */
export class SourceFileCache {
    private cache: Map<Path, CachedSourceFile> = new Map();

    /**
     * Get a cached source file if it exists, is not dirty, and has a compatible parse key.
     * @param path The file path to look up
     * @param parseCacheKey The expected parse cache key for compatibility check
     * @returns The cached source file if found and compatible, undefined otherwise
     */
    get(path: Path, parseCacheKey: string): SourceFile | undefined {
        const entry = this.cache.get(path);
        if (!entry) {
            return undefined;
        }
        if (entry.parseCacheKey !== parseCacheKey) {
            return undefined;
        }
        return entry.file;
    }

    /**
     * Get the raw cache entry for a path (for inspection/debugging).
     */
    getEntry(path: Path): CachedSourceFile | undefined {
        return this.cache.get(path);
    }

    /**
     * Store a source file in the cache.
     * @param path The file path
     * @param file The source file to cache
     * @param parseCacheKey The parse cache key used to create this file
     * @param projectId The ID of the project that fetched this file (will retain it)
     */
    set(path: Path, file: SourceFile, parseCacheKey: string, projectId: string): void {
        const existing = this.cache.get(path);
        if (existing) {
            // Update existing entry
            existing.file = file;
            existing.parseCacheKey = parseCacheKey;
            existing.retainingProjects.add(projectId);
        }
        else {
            // Create new entry
            this.cache.set(path, {
                file,
                parseCacheKey,
                retainingProjects: new Set([projectId]),
            });
        }
    }

    /**
     * Add a project as a retainer of a cache entry.
     * @param path The file path
     * @param projectId The project ID to add as a retainer
     */
    retain(path: Path, projectId: string): void {
        const entry = this.cache.get(path);
        if (entry) {
            entry.retainingProjects.add(projectId);
        }
    }

    /**
     * Remove a project as a retainer of a cache entry.
     * If no projects remain, the entry is removed from the cache.
     * @param path The file path
     * @param projectId The project ID to remove
     */
    release(path: Path, projectId: string): void {
        const entry = this.cache.get(path);
        if (entry) {
            entry.retainingProjects.delete(projectId);
            if (entry.retainingProjects.size === 0) {
                this.cache.delete(path);
            }
        }
    }

    /**
     * Release all cache entries retained by a project.
     * @param projectId The project ID whose entries should be released
     */
    releaseProject(projectId: string): void {
        for (const [path, entry] of this.cache) {
            entry.retainingProjects.delete(projectId);
            if (entry.retainingProjects.size === 0) {
                this.cache.delete(path);
            }
        }
    }

    /**
     * Remove specific paths from the cache.
     * @param paths The paths to remove
     */
    remove(paths: Iterable<Path>): void {
        for (const path of paths) {
            this.cache.delete(path);
        }
    }

    /**
     * Clear all entries from the cache.
     */
    clear(): void {
        this.cache.clear();
    }

    /**
     * Get the number of entries in the cache.
     */
    get size(): number {
        return this.cache.size;
    }

    /**
     * Check if a path is in the cache (regardless of dirty state).
     */
    has(path: Path): boolean {
        return this.cache.has(path);
    }
}
