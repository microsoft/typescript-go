#!/usr/bin/env node

// tsgo: TypeScript compiler executable entry point.
// Uses process.execve for better performance when available, falling back to execFileSync.

import getExePath from "#getExePath";
import { execFileSync } from "node:child_process";

const exe = getExePath();

if (process.platform !== "win32" && typeof process.execve === "function") {
    // > v22.15.0
    try {
        process.execve(exe, [exe, ...process.argv.slice(2)]);
    }
    catch {
        // may not be available, ignore the error and fallback
    }
}

try {
    execFileSync(exe, process.argv.slice(2), { stdio: "inherit" });
}
catch (e) {
    if (e.status) {
        process.exitCode = e.status;
    }
    else {
        throw e;
    }
}
