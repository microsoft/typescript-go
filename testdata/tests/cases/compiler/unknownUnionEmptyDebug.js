// @strict: true
function f() {
    const emptyObj = {};
    const acceptsRecord = (record) => { };
    // This should fail if {} is not assignable to Record<string, string>
    acceptsRecord(emptyObj);
}
