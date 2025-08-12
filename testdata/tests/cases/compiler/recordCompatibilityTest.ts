// @strict: true

function test() {
  const acceptsRecord = (record: Record<string, string>) => {};
  
  // This should still work (fresh object literal)
  acceptsRecord({});
  
  // This should also work (object with string properties)
  acceptsRecord({ a: "hello", b: "world" });
}