//// [tests/cases/compiler/unknownUnionEmpty.ts] ////

//// [unknownUnionEmpty.ts]
function f() {
  const v: unknown = "lol";
  const acceptsRecord = (record: Record<string, string>) => {};
  acceptsRecord(v || {});
}

//// [unknownUnionEmpty.js]
function f() {
    const v = "lol";
    const acceptsRecord = (record) => { };
    acceptsRecord(v || {});
}
