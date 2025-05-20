// @noEmit: true
// @filename: /project/subdir/eslint.config.js
export { default } from '../eslint.config.js';

// @filename: /project/eslint.config.js
export default { rules: {} };

// @filename: /project/main.ts
// This is a regular TypeScript file that will be compiled to JS
// The import should produce an error because of the circular reference
import { default as config } from './subdir/eslint.config.js';
console.log(config);

// This test case reproduces a panic in the module resolution system
// when dealing with circular references in module imports.
// The issue occurs in the Pattern.Matches method when StarIndex is 0,
// causing a slice bounds out of range [1:0] panic.