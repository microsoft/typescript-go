// @allowJs: true
// @esModuleInterop: true
// @noEmit: true
// @filename: shared.vars.js
const foo = ['bar', 'baz'];

module.exports = {
    foo,
};
// @filename: index.ts
import { foo } from "./shared.vars";

console.log(foo);
