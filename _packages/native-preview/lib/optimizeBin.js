import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";
import getExePath from "./getExePath.js";

function optimizeBin() {
    if (process.platform === "win32") {
        return;
    }

    const dirname = path.dirname(fileURLToPath(import.meta.url));
    const normalizedDirname = dirname.replace(/\\/g, "/");
    if (normalizedDirname.endsWith("/_packages/native-preview/lib")) {
        return;
    }

    const exe = getExePath();
    const binDir = path.resolve(dirname, "..", "bin");
    const binPath = path.join(binDir, "tsgo");
    const tempBinPath = path.join(binDir, ".tsgo.tmp");
    const relativeExe = path.relative(binDir, exe);

    fs.rmSync(tempBinPath, { force: true });
    try {
        fs.symlinkSync(relativeExe, tempBinPath);
        fs.renameSync(tempBinPath, binPath);
    }
    finally {
        fs.rmSync(tempBinPath, { force: true });
    }
}

try {
    optimizeBin();
}
catch {
    // Keep the JS wrapper in place when the platform package is unavailable or
    // the package manager does not allow symlinks.
}
