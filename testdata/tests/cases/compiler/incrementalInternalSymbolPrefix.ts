// @strict: true
// @target: esnext
// @incremental: true
// @declaration: true
// @outDir: ./out
// @tsBuildInfoFile: ./out/test.tsbuildinfo
// @lib: esnext

// This test triggers the "X is specified more than once, so this usage will be overwritten" diagnostic
// which contains internal symbol names like __@iterator@. When marshalling build info, the internal
// symbol name prefix (\xFE) caused invalid UTF-8 errors.
// See: https://github.com/microsoft/typescript-go/issues/1531

const items = [1, 2, 3];

// The spread ...items overwrites the [Symbol.iterator] property specified above it
const obj = {
    length: items.length,
    [Symbol.iterator]: function* () {
        for (const item of items) yield item;
    },
    ...items,  // This spread overwrites 'length' and '[Symbol.iterator]'
};
