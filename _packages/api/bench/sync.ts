import {
    API,
    type Project,
    type SourceFile,
} from "@typescript/api/sync";
import { Bench } from "tinybench";

const bench = new Bench({
    name: "Sync API",
    teardown: () => {
        api.close();
        api = undefined!;
        project = undefined!;
        file = undefined!;
    },
});

let api: API;
let project: Project;
let file: SourceFile;

const SMALL_STRING = "ping";
const LARGE_STRING = "a".repeat(1_000_000);
const SMALL_UINT8_ARRAY = new Uint8Array([1, 2, 3, 4, 5]);
const LARGE_UINT8_ARRAY = new Uint8Array(1_000_000);

bench
    .add("spawn API", () => {
        spawnAPI();
    })
    .add("echo (small string)", () => {
        api.echo(SMALL_STRING);
    }, { beforeAll: spawnAPI })
    .add("echo (large string)", () => {
        api.echo(LARGE_STRING);
    }, { beforeAll: spawnAPI })
    .add("echo (small Uint8Array)", () => {
        api.echoBinary(SMALL_UINT8_ARRAY);
    }, { beforeAll: spawnAPI })
    .add("echo (large Uint8Array)", () => {
        api.echoBinary(LARGE_UINT8_ARRAY);
    }, { beforeAll: spawnAPI })
    .add("load project", () => {
        loadProject();
    }, { beforeAll: spawnAPI })
    .add("transfer debug.ts", () => {
        getDebugTS();
    }, { beforeAll: all(spawnAPI, loadProject) })
    .add("transfer program.ts", () => {
        getProgramTS();
    }, { beforeAll: all(spawnAPI, loadProject) })
    .add("transfer checker.ts", () => {
        getCheckerTS();
    }, { beforeAll: all(spawnAPI, loadProject) })
    .add("materialize checker.ts", () => {
        file.forEachChild(function visit(node) {
            node.forEachChild(visit);
        });
    }, { beforeAll: all(spawnAPI, loadProject, getCheckerTS) });

await bench.run();
console.table(bench.table());

function spawnAPI() {
    api = new API({
        cwd: new URL("../../../", import.meta.url).pathname,
        tsserverPath: new URL("../../../built/local/tsgo", import.meta.url).pathname,
    });
}

function loadProject() {
    project = api.loadProject("_submodules/TypeScript/src/compiler/tsconfig.json");
}

function getDebugTS() {
    file = project.getSourceFile("debug.ts")!;
}

function getProgramTS() {
    file = project.getSourceFile("program.ts")!;
}

function getCheckerTS() {
    file = project.getSourceFile("checker.ts")!;
}

function all(...fns: (() => void)[]) {
    return () => {
        for (const fn of fns) {
            fn();
        }
    };
}
