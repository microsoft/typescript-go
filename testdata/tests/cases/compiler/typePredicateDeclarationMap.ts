// @declaration: true
// @declarationMap: true

// @filename: predicateExport.ts
export function createPredicate() {
  return (_item: unknown): _item is boolean => {
    return true;
  };
}

// @filename: predicateImport.ts  
import { createPredicate } from './predicateExport';
export const predicate = createPredicate();