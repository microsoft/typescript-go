import {
    API,
    type Project,
    type Snapshot,
} from "@typescript/api/napi";
import {
    type Node,
    type SourceFile,
    SyntaxKind,
} from "@typescript/ast";
import {
    existsSync,
    writeFileSync,
} from "node:fs";
import inspector from "node:inspector";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { parseArgs } from "node:util";
import { Bench } from "tinybench";
import { RemoteSourceFile } from "../../src/node/node.ts";

const isMain = process.argv[1] === fileURLToPath(import.meta.url);
if (isMain) {
    const { values } = parseArgs({
        options: {
            filter: { type: "string" },
            singleIteration: { type: "boolean", default: false },
            cpuprofile: { type: "boolean", default: false },
        },
    });
    runBenchmarks(values);
}

export function runBenchmarks(options?: { filter?: string; singleIteration?: boolean; cpuprofile?: boolean; }) {
    const { filter, singleIteration, cpuprofile } = options ?? {};
    const repoRoot = fileURLToPath(new URL("../../../../", import.meta.url).toString());
    if (!existsSync(path.join(repoRoot, "_submodules/TypeScript/src/compiler"))) {
        console.warn("Warning: TypeScript submodule is not cloned; skipping benchmarks.");
        return;
    }

    const napiModulePath = path.join(repoRoot, "built", "local", "tsgo.node");
    if (!existsSync(napiModulePath)) {
        console.warn("Warning: tsgo.node not found; run `npx hereby tsgo-napi` first. Skipping NAPI benchmarks.");
        return;
    }

    const bench = new Bench({
        name: "NAPI API",
        teardown,
        iterations: 10,
        warmupIterations: 4,
        ...singleIteration ? {
            iterations: 1,
            warmup: false,
            time: 0,
        } : undefined,
    });

    let api: API;
    let snapshot: Snapshot;
    let project: Project;
    let file: SourceFile;

    const programIdentifierCount = (() => {
        spawnAPI();
        loadSnapshot();
        getProgramTS();
        let count = 0;
        file!.forEachChild(function visit(node) {
            if (node.kind === SyntaxKind.Identifier) {
                count++;
            }
            node.forEachChild(visit);
        });
        teardown();
        return count;
    })();

    const isAsync = false;

    bench
        .add("spawn API", () => {
            spawnAPI();
        }, { async: isAsync })
        .add("load snapshot", () => {
            loadSnapshot();
        }, { async: isAsync, beforeAll: spawnAPI })
        .add("transfer debug.ts", () => {
            getDebugTS();
        }, { async: isAsync, beforeAll: all(spawnAPI, loadSnapshot), beforeEach: clearSourceFileCache })
        .add("transfer program.ts", () => {
            getProgramTS();
        }, { async: isAsync, beforeAll: all(spawnAPI, loadSnapshot), beforeEach: clearSourceFileCache })
        .add("transfer checker.ts", () => {
            getCheckerTS();
        }, { async: isAsync, beforeAll: all(spawnAPI, loadSnapshot), beforeEach: clearSourceFileCache })
        .add("materialize program.ts", () => {
            const { view, _decoder } = file as unknown as RemoteSourceFile;
            new RemoteSourceFile(new Uint8Array(view.buffer, view.byteOffset, view.byteLength), _decoder).forEachChild(function visit(node) {
                node.forEachChild(visit);
            });
        }, { async: isAsync, beforeAll: all(spawnAPI, loadSnapshot, getProgramTS) })
        .add("materialize checker.ts", () => {
            const { view, _decoder } = file as unknown as RemoteSourceFile;
            new RemoteSourceFile(new Uint8Array(view.buffer, view.byteOffset, view.byteLength), _decoder).forEachChild(function visit(node) {
                node.forEachChild(visit);
            });
        }, { async: isAsync, beforeAll: all(spawnAPI, loadSnapshot, getCheckerTS) })
        .add("getSymbolAtPosition - one location", () => {
            project.checker.getSymbolAtPosition("program.ts", 8895);
        }, { async: isAsync, beforeAll: all(spawnAPI, loadSnapshot, createChecker) })
        .add(`getSymbolAtPosition - ${programIdentifierCount} identifiers`, () => {
            for (const node of collectIdentifiers(file)) {
                project.checker.getSymbolAtPosition("program.ts", node.pos);
            }
        }, { async: isAsync, beforeAll: all(spawnAPI, loadSnapshot, createChecker, getProgramTS) })
        .add(`getSymbolAtPosition - ${programIdentifierCount} identifiers (batched)`, () => {
            const positions = collectIdentifiers(file).map(node => node.pos);
            project.checker.getSymbolAtPosition("program.ts", positions);
        }, { async: isAsync, beforeAll: all(spawnAPI, loadSnapshot, createChecker, getProgramTS) })
        .add(`getSymbolAtLocation - ${programIdentifierCount} identifiers`, () => {
            for (const node of collectIdentifiers(file)) {
                project.checker.getSymbolAtLocation(node);
            }
        }, { async: isAsync, beforeAll: all(spawnAPI, loadSnapshot, createChecker, getProgramTS) })
        .add(`getSymbolAtLocation - ${programIdentifierCount} identifiers (batched)`, () => {
            const nodes = collectIdentifiers(file);
            project.checker.getSymbolAtLocation(nodes);
        }, { async: isAsync, beforeAll: all(spawnAPI, loadSnapshot, createChecker, getProgramTS) });

    if (filter) {
        const pattern = filter.toLowerCase();
        for (const task of [...bench.tasks]) {
            if (!task.name.toLowerCase().includes(pattern)) {
                bench.remove(task.name);
            }
        }
    }

    let session: inspector.Session | undefined;
    if (cpuprofile) {
        session = new inspector.Session();
        session.connect();
        session.post("Profiler.enable");
        session.post("Profiler.start");
    }

    bench.runSync();

    if (session) {
        session.post("Profiler.stop", (err, { profile }) => {
            if (err) throw err;
            const outPath = `bench-${Date.now()}.cpuprofile`;
            writeFileSync(outPath, JSON.stringify(profile));
            console.log(`CPU profile written to ${outPath}`);
        });
        session.disconnect();
    }
    console.table(bench.table());

    function collectIdentifiers(sourceFile: SourceFile): Node[] {
        const nodes: Node[] = [];
        sourceFile.forEachChild(function visit(node) {
            if (node.kind === SyntaxKind.Identifier) {
                nodes.push(node);
            }
            node.forEachChild(visit);
        });
        return nodes;
    }

    function spawnAPI() {
        api = new API({
            cwd: repoRoot,
            napiModulePath,
        });
    }

    function loadSnapshot() {
        snapshot = api.updateSnapshot({ openProject: "_submodules/TypeScript/src/compiler/tsconfig.json" });
        project = snapshot.getProjects()[0];
    }

    function createChecker() {
        project.checker.getSymbolAtPosition("core.ts", 0);
    }

    function getDebugTS() {
        file = (project.program.getSourceFile("debug.ts"))!;
    }

    function getProgramTS() {
        file = (project.program.getSourceFile("program.ts"))!;
    }

    function getCheckerTS() {
        file = (project.program.getSourceFile("checker.ts"))!;
    }

    function clearSourceFileCache() {
        api.clearSourceFileCache();
    }

    function teardown() {
        api?.close();
        api = undefined!;
        snapshot = undefined!;
        project = undefined!;
        file = undefined!;
    }

    function all(...fns: (() => void)[]) {
        return () => {
            for (const fn of fns) {
                fn();
            }
        };
    }
}
