import fs from "node:fs";
import module from "node:module";
import path from "node:path";
import { fileURLToPath } from "node:url";

export default function getExePath() {
    // Allow users to bypass resolution entirely. Useful in monorepos where
    // import.meta.resolve is slow across many separate tsgo invocations.
    const envPath = process.env.TSGO_BINARY;
    if (envPath) {
        let envStat;
        try {
            envStat = fs.statSync(envPath);
        }
        catch {
            throw new Error("TSGO_BINARY points to non-existent path: " + envPath);
        }
        if (!envStat.isFile()) {
            throw new Error("TSGO_BINARY does not point to a file: " + envPath);
        }
        if (process.platform === "win32" && envPath.length >= 248) {
            return "\\\\?\\" + envPath;
        }
        return envPath;
    }

    const __dirname = path.dirname(fileURLToPath(import.meta.url));
    const normalizedDirname = __dirname.replace(/\\/g, "/");

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
        // We're running from an installed package.
        // Fast path: check the sibling package at the standard npm layout first.
        // This avoids calling import.meta.resolve, which can be slow in large
        // monorepos where many tsgo invocations happen as separate processes.
        const platformPackageName = "@typescript/" + expectedPackage;
        const siblingDir = path.resolve(__dirname, "..", "..", expectedPackage, "lib");
        const siblingExe = path.join(siblingDir, process.platform === "win32" ? "tsgo.exe" : "tsgo");
        if (fs.existsSync(siblingExe)) {
            exeDir = siblingDir;
        }
        else {
            try {
                if (typeof import.meta.resolve === "undefined") {
                    // v16.20.1
                    const require = module.createRequire(import.meta.url);
                    const packageJson = require.resolve(platformPackageName + "/package.json");
                    exeDir = path.join(path.dirname(packageJson), "lib");
                }
                else {
                    // v20.6.0, v18.19.0
                    const packageJson = import.meta.resolve(platformPackageName + "/package.json");
                    const packageJsonPath = fileURLToPath(packageJson);
                    exeDir = path.join(path.dirname(packageJsonPath), "lib");
                }
            }
            catch (e) {
                throw new Error("Unable to resolve " + platformPackageName + ". Either your platform is unsupported, or you are missing the package on disk.");
            }
        }
    }

    let exe = path.join(exeDir, "tsgo");
    if (process.platform === "win32") {
        exe += ".exe";
        if (exe.length >= 248) {
            exe = "\\\\?\\" + exe;
        }
    }

    if (!fs.existsSync(exe)) {
        throw new Error("Executable not found: " + exe);
    }

    return exe;
}
