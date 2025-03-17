//// [tests/cases/compiler/isolatedModulesDontElideReExportStar.ts] ////

//// [a.ts]
export type T = number;

//// [b.ts]
export * from "./a";


//// [b.js]
export * from "./a";
//// [a.js]
export {};
