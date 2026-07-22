module A {
  export enum ExportedEnum {
    APP_PAGE = 1,
    PAGES = 2,
  }
}

import ImportedEnum = A.ExportedEnum;

// Test 1: Exported enum directly
declare const type1: A.ExportedEnum;
const kind1 = type1 === A.ExportedEnum.APP_PAGE || A.ExportedEnum.PAGES ? 'left' : 'right';

// Test 2: Imported enum alias
declare const type2: ImportedEnum;
const kind2 = type2 === ImportedEnum.APP_PAGE || ImportedEnum.PAGES ? 'left' : 'right';

// Test 3: Normal enum (baseline verification)
enum NormalEnum {
  A = 1,
  B = 2,
}

declare const type3: NormalEnum;
const kind3 = type3 === NormalEnum.A || NormalEnum.B ? 'left' : 'right';
