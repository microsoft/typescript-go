/**
 * Pure JS replacement for @typescript/libsyncrpc.
 *
 * Spawns a child process and communicates with it synchronously over
 * stdin/stdout pipes using the same MessagePack-based tuple protocol:
 *   [MessageType (u8), method (bin), payload (bin)]
 *
 * Synchronous I/O is achieved by calling fs.readSync / fs.writeSync
 * directly on the pipe file descriptors obtained from the spawned
 * ChildProcess.
 */

import {
    type ChildProcess,
    spawn,
} from "node:child_process";
import {
    readSync,
    writeSync,
} from "node:fs";

// ── MessagePack format constants ────────────────────────────────────
const MSGPACK_FIXARRAY3 = 0x93; // 3-element fixarray
const MSGPACK_BIN8 = 0xc4;
const MSGPACK_BIN16 = 0xc5;
const MSGPACK_BIN32 = 0xc6;
const MSGPACK_U8 = 0xcc; // uint8 marker

// ── MessageType constants (matches Go / Rust / TS definitions) ──────
// Sent by channel (parent → child)
const MSG_REQUEST = 1;
const MSG_CALL_RESPONSE = 2;
const MSG_CALL_ERROR = 3;
// Sent by child (child → parent)
const MSG_RESPONSE = 4;
const MSG_ERROR = 5;
const MSG_CALL = 6;

// Pre-allocated buffer used by Atomics.wait for tiny sleeps when a
// non-blocking fd returns EAGAIN.
const sleepBuf = new Int32Array(new SharedArrayBuffer(4));

// ── Global cleanup tracking ─────────────────────────────────────────
// Track all live child processes so they can be killed on process exit.
// This mimics the auto-cleanup behavior of the native libsyncrpc module,
// whose Rust/C++ destructors would kill children automatically.
const liveChildren = new Set<ChildProcess>();

process.on("exit", () => {
    for (const child of liveChildren) {
        try {
            child.kill();
        }
        catch {
            // swallow – process may already be dead
        }
    }
    liveChildren.clear();
});

/**
 * SyncRpcChannel – drop-in replacement for the native libsyncrpc class.
 *
 * API surface intentionally matches the original:
 *   - constructor(exe, args)
 *   - requestSync(method, payload): string
 *   - requestBinarySync(method, payload): Uint8Array
 *   - registerCallback(name, cb)
 *   - close()
 */
export class SyncRpcChannel {
    private child: ChildProcess;
    private readFd: number;
    private writeFd: number;
    private callbacks = new Map<string, (name: string, payload: string) => string>();
    private encoder = new TextEncoder();

    // Reusable 1-byte buffer for readByte()
    private oneByte = Buffer.allocUnsafe(1);
    // Reusable header buffer for reading size fields (up to 4 bytes)
    private headerBuf = Buffer.allocUnsafe(4);

    constructor(exe: string, args: Array<string>) {
        this.child = spawn(exe, args, {
            stdio: ["pipe", "pipe", "inherit"],
        });

        const stdout = this.child.stdout!;
        const stdin = this.child.stdin!;

        // Obtain the underlying OS file descriptors from the libuv pipe
        // handles. On POSIX this is the real fd; on Windows _handle.fd is -1
        // (not supported without native code).
        // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
        this.readFd = (stdout as any)._handle.fd as number;
        // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
        this.writeFd = (stdin as any)._handle.fd as number;

        if (typeof this.readFd !== "number" || this.readFd < 0 || typeof this.writeFd !== "number" || this.writeFd < 0) {
            throw new Error(
                "SyncRpcChannel: could not obtain pipe file descriptors. " +
                    "This implementation requires POSIX (Linux / macOS).",
            );
        }

        // Track for auto-cleanup on process exit.
        liveChildren.add(this.child);

        // Set the pipe handles to blocking mode. Under node --test's
        // process isolation, pipes are created in non-blocking mode
        // (for the IPC channel). This causes readSync/writeSync to get
        // EAGAIN, requiring costly 1ms sleeps per retry. Setting
        // blocking mode ensures readSync blocks properly until data
        // arrives, matching the behavior of the native libsyncrpc.
        // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access, @typescript-eslint/no-unsafe-call
        (stdout as any)._handle.setBlocking(true);
        // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access, @typescript-eslint/no-unsafe-call
        (stdin as any)._handle.setBlocking(true);

        // Prevent Node's event-loop from reading stdout or keeping the
        // process alive – we will use fs.readSync exclusively.
        stdout.pause();
        (stdout as any).unref();
        (stdin as any).unref();
        this.child.unref();
    }

    // ── Public API ──────────────────────────────────────────────────

    /**
     * Send a request and synchronously wait for the response (string).
     * Handles Call (callback) messages from the child inline.
     */
    requestSync(method: string, payload: string): string {
        const result = this.requestBytesSync(method, Buffer.from(payload, "utf-8"));
        return result.toString("utf-8");
    }

    /**
     * Send a request and synchronously wait for the response (binary).
     * Handles Call (callback) messages from the child inline.
     */
    requestBinarySync(method: string, payload: Uint8Array): Uint8Array {
        return this.requestBytesSync(method, payload instanceof Buffer ? payload : Buffer.from(payload));
    }

    /** Register a string→string callback that the child may invoke. */
    registerCallback(name: string, callback: (name: string, payload: string) => string): void {
        this.callbacks.set(name, callback);
    }

    /** Kill the child process and release resources. */
    close(): void {
        try {
            liveChildren.delete(this.child);
            // Destroy the stdio streams so that their pipe handles are closed
            // and no longer prevent the event loop from draining.
            this.child.stdout?.destroy();
            this.child.stdin?.destroy();
            this.child.kill();
            this.readFd = -1;
            this.writeFd = -1;
        }
        catch {
            // swallow – process may already be dead
        }
    }

    // ── Core request loop ───────────────────────────────────────────

    private requestBytesSync(method: string, payload: Buffer): Buffer {
        const methodBuf = Buffer.from(method, "utf-8");
        this.writeTuple(MSG_REQUEST, methodBuf, payload);

        for (;;) {
            const [type, name, data] = this.readTuple();

            switch (type) {
                case MSG_RESPONSE: {
                    const rName = name.toString("utf-8");
                    if (rName !== method) {
                        throw new Error(
                            `name mismatch for response: expected \`${method}\`, got \`${rName}\``,
                        );
                    }
                    return data;
                }
                case MSG_ERROR: {
                    const eName = name.toString("utf-8");
                    if (eName === method) {
                        throw new Error(data.toString("utf-8"));
                    }
                    throw new Error(
                        `name mismatch for response: expected \`${method}\`, got \`${eName}\``,
                    );
                }
                case MSG_CALL: {
                    this.handleCall(name.toString("utf-8"), data);
                    break;
                }
                default:
                    throw new Error(`Invalid message type from child: ${type}`);
            }
        }
    }

    // ── Callback handling ───────────────────────────────────────────

    private handleCall(name: string, payload: Buffer): void {
        const cb = this.callbacks.get(name);
        if (!cb) {
            const errMsg = `unknown callback: \`${name}\`. Please make sure to register it on the JavaScript side before invoking it.`;
            this.writeTuple(
                MSG_CALL_ERROR,
                Buffer.from(name, "utf-8"),
                Buffer.from(errMsg, "utf-8"),
            );
            throw new Error(`no callback named \`${name}\` found`);
        }

        try {
            const result = cb(name, payload.toString("utf-8"));
            this.writeTuple(
                MSG_CALL_RESPONSE,
                Buffer.from(name, "utf-8"),
                Buffer.from(result, "utf-8"),
            );
        }
        catch (e: any) {
            const errMsg = String(e?.message ?? e).trim();
            this.writeTuple(
                MSG_CALL_ERROR,
                Buffer.from(name, "utf-8"),
                Buffer.from(errMsg, "utf-8"),
            );
            throw new Error(`Error calling callback \`${name}\`: ${errMsg}`);
        }
    }

    // ── MessagePack tuple write ─────────────────────────────────────

    private writeTuple(type: number, name: Buffer, payload: Buffer | Uint8Array): void {
        // [0x93] [type as fixint] [bin name] [bin payload]
        this.writeAllBuf(Buffer.from([MSGPACK_FIXARRAY3, type]));
        this.writeBin(name);
        this.writeBin(payload);
    }

    private writeBin(data: Buffer | Uint8Array): void {
        const len = data.length;
        if (len < 0x100) {
            this.headerBuf[0] = MSGPACK_BIN8;
            this.headerBuf[1] = len;
            this.writeAllBuf(this.headerBuf, 2);
        }
        else if (len < 0x10000) {
            this.headerBuf[0] = MSGPACK_BIN16;
            this.headerBuf[1] = (len >>> 8) & 0xff;
            this.headerBuf[2] = len & 0xff;
            this.writeAllBuf(this.headerBuf, 3);
        }
        else {
            const hdr = Buffer.allocUnsafe(5);
            hdr[0] = MSGPACK_BIN32;
            hdr.writeUInt32BE(len, 1);
            this.writeAllBuf(hdr, 5);
        }
        if (len > 0) {
            this.writeAllBuf(data);
        }
    }

    // ── MessagePack tuple read ──────────────────────────────────────

    private readTuple(): [type: number, name: Buffer, payload: Buffer] {
        // Fixed 3-element array marker
        const marker = this.readByte();
        if (marker !== MSGPACK_FIXARRAY3) {
            throw new Error(
                `Expected fixed 3-element array (0x93), received: 0x${marker.toString(16)}`,
            );
        }

        // Message type – positive fixint or uint8
        const tb = this.readByte();
        let msgType: number;
        if (tb <= 0x7f) {
            msgType = tb;
        }
        else if (tb === MSGPACK_U8) {
            msgType = this.readByte();
        }
        else {
            throw new Error(
                `Expected positive fixint or uint8 marker, received: 0x${tb.toString(16)}`,
            );
        }

        const name = this.readBin();
        const payload = this.readBin();
        return [msgType, name, payload];
    }

    private readBin(): Buffer {
        const marker = this.readByte();
        let size: number;
        switch (marker) {
            case MSGPACK_BIN8:
                size = this.readByte();
                break;
            case MSGPACK_BIN16:
                this.readExactInto(this.headerBuf, 2);
                size = (this.headerBuf[0] << 8) | this.headerBuf[1];
                break;
            case MSGPACK_BIN32:
                this.readExactInto(this.headerBuf, 4);
                size = this.headerBuf.readUInt32BE(0);
                break;
            default:
                throw new Error(
                    `Expected binary data (0xc4-0xc6), received: 0x${marker.toString(16)}`,
                );
        }
        if (size === 0) return Buffer.alloc(0);
        return this.readExact(size);
    }

    // ── Low-level synchronous I/O ───────────────────────────────────

    private readByte(): number {
        this.readExactInto(this.oneByte, 1);
        return this.oneByte[0];
    }

    private readExact(length: number): Buffer {
        // Use Buffer.alloc (not allocUnsafe) so the buffer has its own backing
        // ArrayBuffer at byteOffset 0.  This is critical because callers such as
        // RemoteSourceFile create DataView/Uint8Array over buffer.buffer with
        // absolute offsets.  Buffer.allocUnsafe returns slices of a shared pool
        // whose byteOffset is non-zero, corrupting those downstream views.
        const buf = Buffer.alloc(length);
        this.readExactInto(buf, length);
        return buf;
    }

    /**
     * Synchronously read exactly `length` bytes into `buffer`.
     * Retries on EAGAIN (non-blocking pipe) with a tiny sleep.
     */
    private readExactInto(buffer: Buffer, length: number): void {
        let pos = 0;
        while (pos < length) {
            try {
                const n = readSync(this.readFd, buffer, pos, length - pos, null);
                if (n === 0) {
                    throw new Error("Unexpected EOF while reading from child process");
                }
                pos += n;
            }
            catch (e: any) {
                if (e.code === "EAGAIN" || e.code === "EWOULDBLOCK") {
                    // Pipe is non-blocking; yield briefly and retry.
                    Atomics.wait(sleepBuf, 0, 0, 1);
                    continue;
                }
                throw e;
            }
        }
    }

    /**
     * Synchronously write all bytes from `data` (up to `length`).
     * Retries on EAGAIN.
     */
    private writeAllBuf(data: Buffer | Uint8Array, length?: number): void {
        const total = length ?? data.length;
        let pos = 0;
        while (pos < total) {
            try {
                const n = writeSync(this.writeFd, data, pos, total - pos);
                pos += n;
            }
            catch (e: any) {
                if (e.code === "EAGAIN" || e.code === "EWOULDBLOCK") {
                    Atomics.wait(sleepBuf, 0, 0, 1);
                    continue;
                }
                throw e;
            }
        }
    }
}
