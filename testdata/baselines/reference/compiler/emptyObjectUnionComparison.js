//// [tests/cases/compiler/emptyObjectUnionComparison.ts] ////

//// [emptyObjectUnionComparison.ts]
function f() {
  const v: {} = "lol";
  const acceptsRecord = (record: Record<string, string>) => {};
  acceptsRecord(v || {});
}

function g() {
  const v: unknown = "lol";
  const acceptsRecord = (record: Record<string, string>) => {};
  acceptsRecord(v || {});
}

//// [emptyObjectUnionComparison.js]
function f() {
    const v = "lol";
    const acceptsRecord = (record) => { };
    acceptsRecord(v || {});
}
function g() {
    const v = "lol";
    const acceptsRecord = (record) => { };
    acceptsRecord(v || {});
}
