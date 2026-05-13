import { defineConfig } from "vite";

export default defineConfig({
    server: {
        headers: {
            // Required for SharedArrayBuffer and cross-origin isolation
            "Cross-Origin-Opener-Policy": "same-origin",
            "Cross-Origin-Embedder-Policy": "require-corp",
        },
    },
    preview: {
        headers: {
            "Cross-Origin-Opener-Policy": "same-origin",
            "Cross-Origin-Embedder-Policy": "require-corp",
        },
    },
    build: {
        target: "es2024",
    },
});
