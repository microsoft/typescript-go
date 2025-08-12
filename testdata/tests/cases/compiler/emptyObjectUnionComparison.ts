// @strict: true

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