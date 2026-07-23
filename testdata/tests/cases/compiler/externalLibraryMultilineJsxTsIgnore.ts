// @strict: true
// @skipLibCheck: true
// @jsx: preserve
// @module: esnext
// @moduleResolution: bundler
// @noEmit: true

// @filename: /node_modules/pkg/package.json
{
    "name": "pkg",
    "version": "1.0.0",
    "types": "./exports.d.ts"
}

// @filename: /node_modules/pkg/exports.d.ts
export { x, y } from "./src/ui";

// @filename: /node_modules/pkg/src/ui.tsx
declare global {
    namespace JSX {
        interface Element {}
        interface IntrinsicAttributes {}
        interface ElementChildrenAttribute {
            children: {};
        }
        interface IntrinsicElements {
            div: {};
        }
    }
}

declare const Component: (props: {}) => JSX.Element;

export const x = (
    // @ts-ignore
    <Component
        invalid=""
    >
        <div />
    </Component>
);

export const y = (
    <div>
        {/* @ts-ignore */}
        <Component
            invalid=""
        >
            <div />
        </Component>
    </div>
);

// @filename: /index.ts
import { x, y } from "pkg";

x;
y;
