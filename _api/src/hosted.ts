import { dirname } from "path";
import {
    API,
    type APIOptions,
} from "./sync/api.ts";
import type { FileSystemEntries } from "./types.ts";

function createFS(): APIOptions["fs"] {
    const files = new Map<string, string>([
        ["/project/tsconfig.json", `{}`],
        ["/project/index.ts", `import { hello } from "./hello";\nconsole.log(hello());`],
        ["/project/hello.ts", `export function hello() { return "Hello, world!"; }`],
    ]);

    function getAccessibleEntries(directoryName: string): FileSystemEntries {
        if (directoryName === "/") {
            return { files: [], directories: ["/project"] };
        }
        if (!directoryName.startsWith("/project")) {
            return { files: [], directories: [] };
        }
        return {
            files: Array.from(files.keys()).filter(f => dirname(f) === directoryName),
            directories: [],
        };
    }

    return {
        fileExists: path => files.has(path),
        readFile: path => files.get(path),
        directoryExists: () => true,
        getAccessibleEntries,
    };
}

const api = new API({
    tsserverPath: new URL("../../built/local/tsgo", import.meta.url).pathname,
    cwd: dirname(new URL("../../", import.meta.url).pathname),
    fs: createFS(),
});

const project = api.loadProject("/project/tsconfig.json");
const symbol = project.getSymbolAtPosition("index.ts", 47);
console.log(symbol);
console.log(symbol?.getType());
