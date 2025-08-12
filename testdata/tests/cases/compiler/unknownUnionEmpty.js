// @strict: true
// @strictNullChecks: true
function f() {
    const v = "lol";
    const acceptsRecord = (record) => { };
    acceptsRecord(v || {});
}
