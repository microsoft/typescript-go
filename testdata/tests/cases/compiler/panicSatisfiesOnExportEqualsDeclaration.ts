// @allowJs: true
// @checkJs: true
// @noEmit: true
// @module: commonjs
// @Filename: types.d.ts
export interface SupportVersionTraceMap {
    zlib?: any;
    'node:zlib'?: any;
}

// @Filename: panicSatisfiesOnExportEqualsDeclaration.js
const zlib = {};
const READ = Symbol('read');

/**
 * @satisfies {import('./types').SupportVersionTraceMap}
 */
module.exports = {
    zlib: zlib,
    'node:zlib': {
        ...zlib,
        [READ]: { supported: ['14.13.1', '12.20.0'] },
    },
};
