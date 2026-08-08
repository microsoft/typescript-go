import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";
import getExePath from "./getExePath.js";

let optimizedExePath;

export function optimizeBin() {
    if (optimizedExePath !== undefined) {
        return optimizedExePath;
    }

    try {
        optimizedExePath = tryOptimizeBin();
    }
    catch {
        // Keep the JS wrapper in place when the platform package is unavailable.
    }
    return optimizedExePath;
}

function tryOptimizeBin() {
    const dirname = path.dirname(fileURLToPath(import.meta.url));
    const normalizedDirname = dirname.replace(/\\/g, "/");
    if (normalizedDirname.endsWith("/_packages/native-preview/lib")) {
        return undefined;
    }

    const packageDir = path.resolve(dirname, "..");
    const exe = getExePath();
    const binDir = path.join(packageDir, "bin");
    const binPath = path.join(binDir, process.platform === "win32" ? "tsgo.exe" : "tsgo");
    const tempBinPath = path.join(binDir, `.${path.basename(binPath)}.tmp`);
    const relativeExe = path.relative(binDir, exe);

    try {
        fs.rmSync(tempBinPath, { force: true });
        fs.symlinkSync(relativeExe, tempBinPath);
        fs.renameSync(tempBinPath, binPath);

        if (process.platform === "win32") {
            patchWindowsShims(packageDir, binPath);
        }
    }
    catch {
        // Keep the JS wrapper in place when the package manager does not allow symlinks.
    }
    try {
        fs.rmSync(tempBinPath, { force: true });
    }
    catch {
        // ignore cleanup failures
    }

    return exe;
}

function patchWindowsShims(packageDir, exePath) {
    for (const shimDir of getWindowsShimDirs(packageDir)) {
        patchWindowsShim(path.join(shimDir, "tsgo"), exePath, writeShellShim);
        patchWindowsShim(path.join(shimDir, "tsgo.cmd"), exePath, writeCmdShim);
        patchWindowsShim(path.join(shimDir, "tsgo.ps1"), exePath, writePowerShellShim);
    }
}

function getWindowsShimDirs(packageDir) {
    const dirs = [];
    const parent = path.dirname(packageDir);
    if (path.basename(parent).startsWith("@") && path.basename(path.dirname(parent)) === "node_modules") {
        dirs.push(path.join(path.dirname(parent), ".bin"));
    }
    else if (path.basename(parent) === "node_modules") {
        dirs.push(path.join(parent, ".bin"));
    }

    if (process.env.npm_config_global === "true" && process.env.npm_config_prefix) {
        dirs.push(process.env.npm_config_prefix);
    }

    return [...new Set(dirs)];
}

function patchWindowsShim(shimPath, exePath, write) {
    if (fs.existsSync(shimPath)) {
        write(shimPath, exePath);
    }
}

function writeCmdShim(shimPath, exePath) {
    const relativeExe = path.relative(path.dirname(shimPath), exePath).replaceAll("\\", "\\\\");
    fs.writeFileSync(
        shimPath,
        `@ECHO off\r\nSETLOCAL\r\nSET "_prog=%~dp0${relativeExe}"\r\nendLocal & "%_prog%" %*\r\n`,
    );
}

function writePowerShellShim(shimPath, exePath) {
    const relativeExe = path.relative(path.dirname(shimPath), exePath).replaceAll("\\", "/");
    fs.writeFileSync(
        shimPath,
        `#!/usr/bin/env pwsh\n$basedir=Split-Path $MyInvocation.MyCommand.Definition -Parent\n$exe=Join-Path $basedir '${relativeExe}'\n& $exe @args\nexit $LASTEXITCODE\n`,
    );
}

function writeShellShim(shimPath, exePath) {
    const relativeExe = path.relative(path.dirname(shimPath), exePath).replaceAll("\\", "/");
    fs.writeFileSync(
        shimPath,
        `#!/bin/sh\nbasedir=$(dirname "$(echo "$0" | sed -e 's,\\\\,/,g')")\ncase \`uname\` in\n    *CYGWIN*|*MINGW*|*MSYS*) basedir=\`cygpath -w "$basedir"\`;;\nesac\nexec "$basedir/${relativeExe}" "$@"\n`,
    );
}

optimizeBin();
