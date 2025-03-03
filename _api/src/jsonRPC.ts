import type { ChildProcessWithoutNullStreams } from "child_process";
import { JSONRPCClient } from "json-rpc-2.0";
import { LSPError } from "./errors.ts";

export function createJSONRPCClient(serverProcess: ChildProcessWithoutNullStreams, onError: (error: LSPError) => void): JSONRPCClient {
    const client = new JSONRPCClient(jsonRPCRequest => {
        const content = JSON.stringify(jsonRPCRequest);
        const contentLength = Buffer.byteLength(content, "utf8");
        const header = `Content-Length: ${contentLength}\r\n\r\n`;
        serverProcess.stdin.write(header + content, "utf8");
    });

    // Buffer for incoming data
    let buffer = "";
    let contentLength: number | null = null;

    serverProcess.stdout.on("data", (data: Buffer) => {
        buffer += data.toString();

        // Process all complete messages in the buffer
        while (true) {
            // If we don't have a content length yet, try to parse the header
            if (contentLength === null) {
                const headerEnd = buffer.indexOf("\r\n\r\n");
                if (headerEnd === -1) {
                    // Not enough data to read the header
                    break;
                }

                const header = buffer.substring(0, headerEnd);
                const match = header.match(/Content-Length:\s*(\d+)/i);
                if (!match) {
                    onError(
                        new LSPError(
                            "Invalid message header: Content-Length not found",
                            "MessageFormat",
                        ),
                    );
                    buffer = "";
                    return;
                }

                contentLength = parseInt(match[1], 10);
                buffer = buffer.substring(headerEnd + 4); // Remove the header
            }

            // If we have a content length and enough data, process the message
            if (contentLength !== null && buffer.length >= contentLength) {
                const message = buffer.substring(0, contentLength);
                buffer = buffer.substring(contentLength); // Remove the processed message
                contentLength = null;

                try {
                    const json = JSON.parse(message);
                    client.receive(json);
                }
                catch (error) {
                    onError(
                        new LSPError(
                            `Error parsing message: ${error instanceof Error ? error.message : error}`,
                            "Parse",
                            { raw: message },
                        ),
                    );
                }
            }
            else {
                // Not enough data to process a message
                break;
            }
        }
    });

    return client;
}
