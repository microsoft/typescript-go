export type MaybeAsync<Async extends boolean, T> = Async extends true ? Promise<T> : T;

export interface FileSystemEntries {
    files: string[];
    directories: string[];
}
