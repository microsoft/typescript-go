//// [tests/cases/compiler/stylexOptionsExportEqualsMerging.ts] ////

//// [package.json]
{
    "name": "@stylexjs/babel-plugin",
    "version": "0.16.2",
    "main": "./lib/index.js",
    "types": "./lib/index.d.ts"
}

//// [shared.d.ts]
export interface SharedOptions {
    answer: number;
}

//// [index.d.ts]
import type { SharedOptions } from "./shared";

export type Options = SharedOptions;

declare function stylexPlugin(): void;

declare const exported: {
    (): void;
    withOptions(options: Partial<SharedOptions>): [typeof stylexPlugin, Partial<SharedOptions>];
};

export = exported;

//// [stylex.ts]
import { Options } from "@stylexjs/babel-plugin";

const example: Options = { answer: 42 };
void example;


//// [stylex.js]
const example = { answer: 42 };
void example;
export {};
