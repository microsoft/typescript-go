// tsgo LSP Web Worker
// Loads the Go WASM binary and bridges stdin/stdout for LSP communication
// with the main thread via postMessage.

// ---------------------------------------------------------------------------
// Virtual filesystem
// ---------------------------------------------------------------------------
const virtualFiles = new Map();
const virtualDirs = new Set(["/", "/playground", "/home", "/tmp"]);

// ---------------------------------------------------------------------------
// Stdin buffer (LSP messages from main thread -> Go WASM)
// ---------------------------------------------------------------------------
const stdinChunks = [];
let pendingReadCb = null;

// ---------------------------------------------------------------------------
// Stdout buffer (Go WASM -> main thread)
// ---------------------------------------------------------------------------
let stdoutBuf = new Uint8Array(0);

// ---------------------------------------------------------------------------
// File descriptor management (for virtual file reads)
// ---------------------------------------------------------------------------
let nextFd = 10;
const openFds = new Map(); // fd -> { data: Uint8Array, pos: number }

// ---------------------------------------------------------------------------
// Utility functions
// ---------------------------------------------------------------------------
function enosys() {
    const err = new Error("not implemented");
    err.code = "ENOSYS";
    return err;
}

function enoent(op, path) {
    const err = new Error(`ENOENT: no such file or directory, ${op} '${path}'`);
    err.code = "ENOENT";
    return err;
}

function concat(a, b) {
    const c = new Uint8Array(a.length + b.length);
    c.set(a);
    c.set(b, a.length);
    return c;
}

function makeStat(isDir, size) {
    const now = Date.now();
    return {
        dev: 0,
        ino: 0,
        mode: isDir ? 0o40755 : 0o100644,
        nlink: 1,
        uid: 0,
        gid: 0,
        rdev: 0,
        size: size || 0,
        blksize: 4096,
        blocks: 0,
        atimeMs: now,
        mtimeMs: now,
        ctimeMs: now,
        birthtimeMs: now,
        atime: new Date(now),
        mtime: new Date(now),
        ctime: new Date(now),
        birthtime: new Date(now),
        isDirectory() {
            return isDir;
        },
        isFile() {
            return !isDir;
        },
        isSymbolicLink() {
            return false;
        },
        isBlockDevice() {
            return false;
        },
        isCharacterDevice() {
            return false;
        },
        isFIFO() {
            return false;
        },
        isSocket() {
            return false;
        },
    };
}

// ---------------------------------------------------------------------------
// Stdin helpers
// ---------------------------------------------------------------------------
function drainStdin() {
    if (!pendingReadCb || stdinChunks.length === 0) return;
    const { buffer, offset, length, callback } = pendingReadCb;
    pendingReadCb = null;

    const chunk = stdinChunks[0];
    const n = Math.min(length, chunk.length);
    for (let i = 0; i < n; i++) buffer[offset + i] = chunk[i];

    if (n < chunk.length) {
        stdinChunks[0] = chunk.subarray(n);
    }
    else {
        stdinChunks.shift();
    }
    callback(null, n);
}

// ---------------------------------------------------------------------------
// Stdout parser — extracts LSP messages from Content-Length framed stream
// and posts them as raw JSON-RPC Message objects via postMessage
// (compatible with @hediet/json-rpc-browser's worker transport)
// ---------------------------------------------------------------------------
function parseStdoutMessages() {
    const decoder = new TextDecoder();
    const encoder = new TextEncoder();
    while (true) {
        const str = decoder.decode(stdoutBuf);
        const idx = str.indexOf("\r\n\r\n");
        if (idx === -1) break;

        const headerStr = str.substring(0, idx);
        const match = headerStr.match(/Content-Length:\s*(\d+)/i);
        if (!match) break;

        const contentLength = parseInt(match[1], 10);
        const headerBytesLen = encoder.encode(str.substring(0, idx + 4)).length;
        const totalLength = headerBytesLen + contentLength;

        if (stdoutBuf.length < totalLength) break;

        const body = decoder.decode(
            stdoutBuf.subarray(headerBytesLen, totalLength),
        );
        stdoutBuf = stdoutBuf.subarray(totalLength);

        // Post raw JSON-RPC Message object
        self.postMessage(JSON.parse(body));
    }
}

// ---------------------------------------------------------------------------
// Handle messages from main thread
// Receives raw JSON-RPC Message objects from createTransportToWorker
// ---------------------------------------------------------------------------
self.onmessage = event => {
    const data = event.data;
    // Everything is a raw JSON-RPC message from the transport
    const bodyStr = JSON.stringify(data);
    const body = new TextEncoder().encode(bodyStr);
    const header = new TextEncoder().encode(
        `Content-Length: ${body.length}\r\n\r\n`,
    );
    stdinChunks.push(concat(header, body));
    drainStdin();
};

// ---------------------------------------------------------------------------
// Set up global `fs` BEFORE loading wasm_exec.js
// (wasm_exec.js checks `if (!globalThis.fs)` and skips if already set)
// ---------------------------------------------------------------------------
globalThis.fs = {
  constants: {
    O_WRONLY: -1,
    O_RDWR: -1,
    O_CREAT: -1,
    O_TRUNC: -1,
    O_APPEND: -1,
    O_EXCL: -1,
    O_DIRECTORY: -1,
  },

  writeSync(fd, buf) {
    if (fd === 1) {
      // stdout -> capture LSP output
      stdoutBuf = concat(stdoutBuf, buf);
      parseStdoutMessages();
      return buf.length;
    }
    if (fd === 2) {
      // stderr -> console
      const text = new TextDecoder().decode(buf);
      if (text.trim()) console.warn("[tsgo]", text.trimEnd());
      return buf.length;
    }
    console.warn("writeSync to unknown fd", fd);
    return buf.length;
  },

  write(fd, buf, offset, length, position, callback) {
    if (fd === 1 || fd === 2) {
      const data = buf.subarray(offset, offset + length);
      const n = this.writeSync(fd, data);
      callback(null, n);
      return;
    }
    callback(enosys());
  },

  read(fd, buffer, offset, length, position, callback) {
    if (fd === 0) {
      // stdin -> LSP input from main thread
      if (stdinChunks.length > 0) {
        const chunk = stdinChunks[0];
        const n = Math.min(length, chunk.length);
        for (let i = 0; i < n; i++) buffer[offset + i] = chunk[i];
        if (n < chunk.length) {
          stdinChunks[0] = chunk.subarray(n);
        } else {
          stdinChunks.shift();
        }
        callback(null, n);
      } else {
        // No data yet — store callback; will fire when data arrives via postMessage
        pendingReadCb = { buffer, offset, length, callback };
      }
      return;
    }

    // Virtual file reads
    const file = openFds.get(fd);
    if (file) {
      const pos =
        position !== null && position !== undefined ? position : file.pos;
      const n = Math.min(length, file.data.length - pos);
      if (n <= 0) {
        callback(null, 0);
        return;
      }
      for (let i = 0; i < n; i++) buffer[offset + i] = file.data[pos + i];
      file.pos = pos + n;
      callback(null, n);
      return;
    }

    callback(enosys());
  },

  open(path, flags, mode, callback) {
    if (virtualFiles.has(path)) {
      const fd = nextFd++;
      const content = virtualFiles.get(path);
      openFds.set(fd, { data: new TextEncoder().encode(content), pos: 0 });
      callback(null, fd);
      return;
    }
    callback(enoent("open", path));
  },

  close(fd, callback) {
    openFds.delete(fd);
    callback(null);
  },

  stat(path, callback) {
    if (virtualDirs.has(path)) {
      callback(null, makeStat(true));
      return;
    }
    if (virtualFiles.has(path)) {
      const content = virtualFiles.get(path);
      callback(
        null,
        makeStat(false, new TextEncoder().encode(content).length)
      );
      return;
    }
    callback(enoent("stat", path));
  },

  lstat(path, callback) {
    this.stat(path, callback);
  },

  fstat(fd, callback) {
    const file = openFds.get(fd);
    if (file) {
      callback(null, makeStat(false, file.data.length));
      return;
    }
    if (fd <= 2) {
      // stdin/stdout/stderr
      callback(null, makeStat(false, 0));
      return;
    }
    callback(enosys());
  },

  readdir(path, callback) {
    const prefix = path.endsWith("/") ? path : path + "/";
    const entries = new Set();

    for (const filePath of virtualFiles.keys()) {
      if (filePath.startsWith(prefix)) {
        const rest = filePath.substring(prefix.length);
        const name = rest.split("/")[0];
        if (name) entries.add(name);
      }
    }
    for (const dirPath of virtualDirs) {
      if (dirPath.startsWith(prefix) && dirPath !== path) {
        const rest = dirPath.substring(prefix.length);
        const name = rest.split("/")[0];
        if (name) entries.add(name);
      }
    }

    if (entries.size > 0 || virtualDirs.has(path)) {
      callback(null, [...entries]);
    } else {
      callback(enoent("readdir", path));
    }
  },

  chmod(path, mode, callback) {
    callback(null);
  },
  chown(path, uid, gid, callback) {
    callback(null);
  },
  fchmod(fd, mode, callback) {
    callback(null);
  },
  fchown(fd, uid, gid, callback) {
    callback(null);
  },
  fsync(fd, callback) {
    callback(null);
  },
  ftruncate(fd, length, callback) {
    callback(null);
  },
  lchown(path, uid, gid, callback) {
    callback(null);
  },
  link(path, link, callback) {
    callback(enosys());
  },
  mkdir(path, perm, callback) {
    virtualDirs.add(path);
    callback(null);
  },
  readlink(path, callback) {
    callback(enosys());
  },
  rename(from, to, callback) {
    callback(enosys());
  },
  rmdir(path, callback) {
    virtualDirs.delete(path);
    callback(null);
  },
  symlink(path, link, callback) {
    callback(enosys());
  },
  truncate(path, length, callback) {
    callback(null);
  },
  unlink(path, callback) {
    virtualFiles.delete(path);
    callback(null);
  },
  utimes(path, atime, mtime, callback) {
    callback(null);
  },
};

// ---------------------------------------------------------------------------
// Set up global `process` BEFORE loading wasm_exec.js
// ---------------------------------------------------------------------------
globalThis.process = {
    getuid() {
        return -1;
    },
    getgid() {
        return -1;
    },
    geteuid() {
        return -1;
    },
    getegid() {
        return -1;
    },
    getgroups() {
        return [];
    },
    pid: 1,
    ppid: 0,
    umask() {
        return 0;
    },
    cwd() {
        return "/playground";
    },
    chdir() {},
    env: { HOME: "/home", TMPDIR: "/tmp" },
};

// ---------------------------------------------------------------------------
// Set up global `path` BEFORE loading wasm_exec.js
// ---------------------------------------------------------------------------
globalThis.path = {
    resolve(...segments) {
        return segments.join("/");
    },
};

// ---------------------------------------------------------------------------
// Load wasm_exec.js (sets globalThis.Go)
// Since we already set fs/process/path, wasm_exec.js will skip its stubs.
// ---------------------------------------------------------------------------
importScripts("/wasm_exec.js");

// ---------------------------------------------------------------------------
// Start the Go WASM LSP server
// ---------------------------------------------------------------------------
async function start() {
    try {
        const go = new Go();
        go.argv = ["tsgo", "--lsp", "--stdio"];
        go.env = { HOME: "/home", TMPDIR: "/tmp" };

        const result = await WebAssembly.instantiateStreaming(
            fetch("/tsgo.wasm"),
            go.importObject,
        );

        console.log("[tsgo] WASM loaded, starting LSP server...");

        // Run Go program — blocks until the Go program exits
        await go.run(result.instance);

        console.log("[tsgo] Go WASM process exited");
    }
    catch (err) {
        console.error("WASM startup error:", err);
    }
}

start();
