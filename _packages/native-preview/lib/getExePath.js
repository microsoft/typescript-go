import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

export default function getExePath() {
    const __dirname = path.dirname(fileURLToPath(import.meta.url));
    const normalizedDirname = __dirname.replace(/\\/g, "/");

    let exeDir;
    let isDev = false;

    const expectedPackage = "native-preview-" + process.platform + "-" + process.arch;

    if (normalizedDirname.endsWith("/_packages/native-preview/lib")) {
        // We're running directly from source in the repo.
        exeDir = path.resolve(__dirname, "..", "..", "..", "built", "local");
        isDev = true;
    }
    else {
        // Peer dependency is a sibling directory; resolve relative to this file.
        exeDir = path.resolve(__dirname, "..", "..", expectedPackage, "lib");
    }

    let exe = path.join(exeDir, "tsgo");
    if (process.platform === "win32") {
        exe += ".exe";
        if (exe.length >= 248) {
            exe = "\\\\?\\" + exe;
        }
    }

    if (!fs.existsSync(exe)) {
        if (isDev) {
            throw new Error("tsgo executable not found at " + exe + ". Run 'npx hereby build' to build it.");
        }
        throw new Error([
            "Could not find the tsgo executable for your platform (" + process.platform + "-" + process.arch + ").",
            "The package @typescript/" + expectedPackage + " may not be installed; try reinstalling your dependencies.",
        ].join(" "));
    }

    return exe;
}
