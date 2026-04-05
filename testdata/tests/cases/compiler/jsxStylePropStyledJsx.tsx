// @target: es2015
// @jsx: react
// @strict: true
// @module: commonjs
// @noEmit: true

// https://github.com/microsoft/typescript-go/issues/TODO
// Validates that the styled-jsx `jsx` and `global` boolean props are accepted on
// <style> elements when the types are properly declared via module augmentation.

/// <reference path="/.lib/react16.d.ts" />

import * as React from 'react';

// Module augmentation as provided by styled-jsx (global.d.ts)
declare module 'react' {
    interface StyleHTMLAttributes<T> extends HTMLAttributes<T> {
        jsx?: boolean;
        global?: boolean;
    }
}

// Should compile without errors - jsx is in the augmented StyleHTMLAttributes
const a = <style jsx>{`h1 { color: red; }`}</style>;

// Should compile without errors - global is in the augmented StyleHTMLAttributes
const b = <style global>{`h1 { color: red; }`}</style>;

// Should compile without errors - both jsx and global
const c = <style jsx global>{`h1 { color: red; }`}</style>;

// Should still error - unknown is not in the type
// @ts-expect-error
const d = <style unknown>{`h1 { color: red; }`}</style>;
