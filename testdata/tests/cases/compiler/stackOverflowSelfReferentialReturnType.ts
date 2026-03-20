// @strict: true

// Regression test for https://github.com/microsoft/TypeScript/issues/63273
// Self-referential type involving ReturnType and typeof should not cause stack overflow.
function clone(): <T>(obj: T) => T extends any ? ReturnType<typeof clone>[0] : 0;
