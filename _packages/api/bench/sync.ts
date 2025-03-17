import type { GetSymbolAtPositionParams } from "@typescript/api/base/proto";
import {
    API,
    type Project,
    type SourceFile,
} from "@typescript/api/sync";
import { SyntaxKind } from "@typescript/ast";
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
    .add("materialize program.ts", () => {
        file.forEachChild(function visit(node) {
            node.forEachChild(visit);
        });
    }, { beforeAll: all(spawnAPI, loadProject, getProgramTS) })
    .add("materialize checker.ts", () => {
        file.forEachChild(function visit(node) {
            node.forEachChild(visit);
        });
    }, { beforeAll: all(spawnAPI, loadProject, getCheckerTS) })
    .add("getSymbolAtPosition - one location", () => {
        project.getSymbolAtPosition("program.ts", 8895);
    }, { beforeAll: all(spawnAPI, loadProject) })
    .add("getSymbolAtPosition - all identifiers", () => {
        file.forEachChild(function visit(node) {
            if (node.kind === SyntaxKind.Identifier) {
                project.getSymbolAtPosition("program.ts", node.pos);
            }
            node.forEachChild(visit);
        });
    }, { beforeAll: all(spawnAPI, loadProject, getProgramTS) })
    .add("getSymbolAtPosition - all identifiers (batched)", () => {
        const positions: GetSymbolAtPositionParams[] = [];
        file.forEachChild(function visit(node) {
            if (node.kind === SyntaxKind.Identifier) {
                positions.push({ fileName: "program.ts", position: node.pos });
            }
            node.forEachChild(visit);
        });
        project.getSymbolAtPosition(positions);
    }, { beforeAll: all(spawnAPI, loadProject, getProgramTS) })
    .add("getSymbolAtPosition - all identifiers (batched 2)", () => {
        const positions: number[] = [];
        file.forEachChild(function visit(node) {
            if (node.kind === SyntaxKind.Identifier) {
                positions.push(node.pos);
            }
            node.forEachChild(visit);
        });
        project.getSymbolAtPosition("program.ts", positions);
    }, { beforeAll: all(spawnAPI, loadProject, getProgramTS) });

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
