// @module: nodenext
// @moduleResolution: nodenext
// @noEmit: true

// @filename: /node_modules/pkg/package.json
{ "name": "pkg", "version": "1.0.0", "types": "index.d.ts" }

// @filename: /node_modules/pkg/index.d.ts
export declare const main: string;

// @filename: /node_modules/pkg/sub/package.json
{ "name": "pkg/sub", "version": "1.0.0", "types": "index.d.ts" }

// @filename: /node_modules/pkg/sub/index.d.ts
export declare const sub: number;

// @filename: /index.mts
import { main } from "pkg";
import { sub } from "pkg/sub";
