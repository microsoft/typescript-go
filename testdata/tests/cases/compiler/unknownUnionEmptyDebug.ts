// @strict: true

function f() {
  const emptyObj = {};
  const acceptsRecord = (record: Record<string, string>) => {};
  
  // This should fail if {} is not assignable to Record<string, string>
  acceptsRecord(emptyObj);
}