// @declaration: true

// @filename: index.ts
export interface Foo {}

// @filename: main.ts
/**
 * Some random docs not related to foo
 */
/* trigger */
import * as x from './index.js';
export const foo = 1;
