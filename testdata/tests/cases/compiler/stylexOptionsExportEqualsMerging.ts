// @skipLibCheck: true
// @filename: /node_modules/@stylexjs/babel-plugin/package.json
{
    "name": "@stylexjs/babel-plugin",
    "version": "0.16.2",
    "main": "./lib/index.js",
    "types": "./lib/index.d.ts"
}

// @filename: /node_modules/@stylexjs/babel-plugin/lib/shared.d.ts
export interface SharedOptions {
    answer: number;
}

// @filename: /node_modules/@stylexjs/babel-plugin/lib/index.d.ts
import type { SharedOptions } from "./shared";

export type Options = SharedOptions;

declare function stylexPlugin(): void;

declare const exported: {
    (): void;
    withOptions(options: Partial<SharedOptions>): [typeof stylexPlugin, Partial<SharedOptions>];
};

export = exported;

// @filename: /stylex.ts
import { Options } from "@stylexjs/babel-plugin";

const example: Options = { answer: 42 };
void example;
