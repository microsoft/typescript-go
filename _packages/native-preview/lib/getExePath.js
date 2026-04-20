import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

let cachedExePath;

function findPeerPackageLib(startDir, packageName) {
    let dir = startDir;
    while (true) {
        const candidate = path.join(dir, "node_modules", packageName, "lib");
        if (fs.existsSync(candidate)) {
            return candidate;
        }
        const parent = path.dirname(dir);
        if (parent === dir) {
            return undefined;
        }
        dir = parent;
    }
}

export default function getExePath() {
    if (cachedExePath !== undefined) {
        return cachedExePath;
    }

    const __dirname = path.dirname(fileURLToPath(import.meta.url));
    const normalizedDirname = __dirname.replace(/\/g, "/");

    let exeDir;

    const expectedPackage = "native-preview-" + process.platform + "-" + process.arch;

    if (normalizedDirname.endsWith("/_packages/native-preview/lib")) {
        // We're running directly from source in the repo.
        exeDir = path.resolve(__dirname, "..", "..", "..", "built", "local");
    }
    else if (normalizedDirname.endsWith("/built/npm/native-preview/lib")) {
        // We're running from the built output.
        exeDir = path.resolve(__dirname, "..", "..", expectedPackage, "lib");
    }
    else {
        // We're running from an installed package. Use plain FS lookups instead of
        // import.meta.resolve/require.resolve, which can be slow in large monorepos
        // because they trigger full package resolution. Peer dependencies are always
        // placed in a parent node_modules directory, so a simple directory walk suffices.
        const platformPackageName = "@typescript/" + expectedPackage;
        exeDir = findPeerPackageLib(path.resolve(__dirname, ".."), platformPackageName);

        if (exeDir === undefined) {
            throw new Error("Unable to resolve " + platformPackageName + ". Either your platform is unsupported, or you are missing the package on disk.");
        }
    }

    let exe = path.join(exeDir, "tsgo");
    if (process.platform === "win32") {
        exe += ".exe";
        if (exe.length >= 248) {
            exe = "\\?\\" + exe;
        }
    }

    if (!fs.existsSync(exe)) {
        throw new Error("Executable not found: " + exe);
    }

    cachedExePath = exe;
    return exe;
}
