//// [tests/cases/compiler/unknownUnionEmptyDebug.ts] ////

//// [unknownUnionEmptyDebug.ts]
function f() {
  const emptyObj = {};
  const acceptsRecord = (record: Record<string, string>) => {};
  
  // This should fail if {} is not assignable to Record<string, string>
  acceptsRecord(emptyObj);
}

//// [unknownUnionEmptyDebug.js]
function f() {
    const emptyObj = {};
    const acceptsRecord = (record) => { };
    // This should fail if {} is not assignable to Record<string, string>
    acceptsRecord(emptyObj);
}
