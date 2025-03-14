export interface FileSystemEntries {
    files: string[];
    directories: string[];
}

export interface FileSystem {
    directoryExists?: (directoryName: string) => boolean | undefined;
    fileExists?: (fileName: string) => boolean | undefined;
    getAccessibleEntries?: (directoryName: string) => FileSystemEntries | undefined;
    readFile?: (fileName: string) => string | null | undefined;
    realpath?: (path: string) => string | undefined;
}
