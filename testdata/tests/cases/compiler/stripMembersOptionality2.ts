// @strict: true
// @exactOptionalPropertyTypes: true, false
// @noEmit: true

// https://github.com/microsoft/TypeScript/issues/63291

type WithObject1 = Required<{ a?: string | undefined }>;
const obj1: WithObject1 = { a: undefined };
type WithArray1 = Required<[(string | undefined)?]>;
const tup1: WithArray1 = [undefined];
