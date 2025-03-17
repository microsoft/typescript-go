//// [tests/cases/compiler/isolatedDeclarationErrorsDefault.ts] ////

//// [a.ts]
export default 1 + 1;


//// [b.ts]
export default { foo: 1 + 1 };

//// [c.ts]
export default [{ foo: 1 + 1 }];

//// [d.ts]
export default [{ foo: 1 + 1 }] as const;

//// [e.ts]
export default [{ foo: 1 + 1 }] as const;

//// [f.ts]
const a = { foo: 1 };
export default a;

//// [f.js]
const a = { foo: 1 };
export default a;
//// [e.js]
export default [{ foo: 1 + 1 }];
//// [d.js]
export default [{ foo: 1 + 1 }];
//// [c.js]
export default [{ foo: 1 + 1 }];
//// [b.js]
export default { foo: 1 + 1 };
//// [a.js]
export default 1 + 1;
