// @allowJs: true
// @noEmit: true
// @esModuleInterop: true
// @module: commonjs
// @Filename: shared.vars.js
const foo = ['bar', 'baz'];

module.exports = {
    foo,
};

// @Filename: index.ts
import { foo } from "./shared.vars";

console.log(foo);