import {
    type ChildProcessWithoutNullStreams,
    spawn,
} from "child_process";
import { LSPError } from "./errors.ts";

export function startServer(executable: string, cwd: string, log?: (msg: string) => void): ChildProcessWithoutNullStreams {
    const server = spawn(executable, ["api", "-cwd", cwd], {
        detached: true,
    });

    server.unref();

    server.on("error", error => {
        throw new LSPError(`Server process error: ${error.message}`, "Server");
    });

    if (log) {
        server.stderr.on("data", data => {
            log(data.toString());
        });
    }

    server.once("exit", code => {
        if (code !== 0 && code !== null) {
            throw new LSPError(`Server exited with code ${code}`, "Server", { code });
        }
    });

    return server;
}
