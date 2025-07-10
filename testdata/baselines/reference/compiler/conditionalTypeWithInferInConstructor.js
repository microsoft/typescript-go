//// [tests/cases/compiler/conditionalTypeWithInferInConstructor.ts] ////

//// [conditionalTypeWithInferInConstructor.ts]
// This is the exact case from issue #1379
type ExtractReturn<T> = T extends { new(): infer R } ? R : never;

//// [conditionalTypeWithInferInConstructor.js]


//// [conditionalTypeWithInferInConstructor.d.ts]
// This is the exact case from issue #1379
type ExtractReturn<T> = T extends {
    new ();
} ? R : never;


//// [DtsFileErrors]


conditionalTypeWithInferInConstructor.d.ts(4,5): error TS2304: Cannot find name 'R'.


==== conditionalTypeWithInferInConstructor.d.ts (1 errors) ====
    // This is the exact case from issue #1379
    type ExtractReturn<T> = T extends {
        new ();
    } ? R : never;
        ~
!!! error TS2304: Cannot find name 'R'.
    