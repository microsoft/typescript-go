// @target: es2015
// @strict: true
// @noEmit: true

// repro #51138

type SubTup2FixedLength<T extends unknown[]> = T extends [
  ...infer B extends [any, any],
  any
]
  ? B
  : never;

type SubTup2FixedLengthTest2 = SubTup2FixedLength<[a: 0]>;

type SubTup2Variadic<T extends unknown[]> = T extends [
  ...infer B extends [any, any],
  ...any
]
  ? B
  : never;

type SubTup2VariadicTest3 = SubTup2Variadic<[a: 0, ...b: 1[]]>;
type SubTup2VariadicTest4 = SubTup2Variadic<[...a: 0[]]>;
type SubTup2VariadicTest5 = SubTup2Variadic<[a: 0, b: 1]>;

type SubTup2VariadicAndRest<T extends unknown[]> = T extends [
  ...infer B extends [any, any],
  ...(infer C)[]
]
  ? [...B, ...[C]]
  : never;

type SubTup2VariadicAndRestTest2 = SubTup2VariadicAndRest<[a: 0, ...b: 1[]]>;
type SubTup2VariadicAndRestTest3 = SubTup2VariadicAndRest<[...a: 0[]]>;
type SubTup2VariadicAndRestTest4 = SubTup2VariadicAndRest<[a: 0, b: 1]>;

type SubTup2TrailingVariadic<T extends unknown[]> = T extends [
  ...any,
  ...infer B extends [any, any],
]
  ? B
  : never;

type SubTup2TrailingVariadicTest3 = SubTup2TrailingVariadic<[...a: 0[], b: 1]>;
type SubTup2TrailingVariadicTest4 = SubTup2TrailingVariadic<[...a: 0[]]>;
type SubTup2TrailingVariadicTest5 = SubTup2TrailingVariadic<[b: 1, c: 2]>;

type SubTup2RestAndTrailingVariadic2<T extends unknown[]> = T extends [
  ...(infer C)[],
  ...infer B extends [any, any],
]
  ? [C, ...B]
  : never;

type SubTup2RestAndTrailingVariadic2Test2 = SubTup2RestAndTrailingVariadic2<[...a: 0[], b: 1]>;
type SubTup2RestAndTrailingVariadic2Test3 = SubTup2RestAndTrailingVariadic2<[...a: 0[]]>;
type SubTup2RestAndTrailingVariadic2Test4 = SubTup2RestAndTrailingVariadic2<[b: 1, c: 2]>;

type SubTup2VariadicWithLeadingFixedElements<T extends unknown[]> = T extends [
  any,
  ...infer B extends [any, any],
  ...any
]
  ? B
  : never;

type SubTup2VariadicWithLeadingFixedElementsTest3 = SubTup2VariadicWithLeadingFixedElements<[a: 0, b: 1, ...c: 2[]]>;
type SubTup2VariadicWithLeadingFixedElementsTest4 = SubTup2VariadicWithLeadingFixedElements<[a: 0, ...b: 1[]]>;
type SubTup2VariadicWithLeadingFixedElementsTest5 = SubTup2VariadicWithLeadingFixedElements<[...a: 0[]]>;

type SubTup2VariadicWithLeadingFixedElements2<T extends unknown[]> = T extends [
  any,
  any,
  ...infer B extends [any],
  ...any
]
  ? B
  : never;

type SubTup2VariadicWithLeadingFixedElements2Test = SubTup2VariadicWithLeadingFixedElements2<[a: 0, b: 1, c: 2, ...d: 3[]]>;
type SubTup2VariadicWithLeadingFixedElements2Test2 = SubTup2VariadicWithLeadingFixedElements2<[a: 0, b: 1, c: 2, d: 3, ...e: 4[]]>;
type SubTup2VariadicWithLeadingFixedElements2Test3 = SubTup2VariadicWithLeadingFixedElements2<[a: 0, b: 1, ...c: 2[]]>;
type SubTup2VariadicWithLeadingFixedElements2Test4 = SubTup2VariadicWithLeadingFixedElements2<[a: 0, ...b: 1[]]>;
type SubTup2VariadicWithLeadingFixedElements2Test5 = SubTup2VariadicWithLeadingFixedElements2<[...a: 0[]]>;

type SubTup2VariadicWithLeadingFixedElements3<T extends unknown[]> = T extends [
  any,
  ...infer B extends [any, any],
  ...(infer C)[]
]
  ? [B, C]
  : never;

type SubTup2VariadicWithLeadingFixedElements3Test = SubTup2VariadicWithLeadingFixedElements3<[a: 0, b: 1, c: 2, d: 3, ...e: 4[]]>;
type SubTup2VariadicWithLeadingFixedElements3Test2 = SubTup2VariadicWithLeadingFixedElements3<[a: 0, b: 1, c: 2, ...d: 3[]]>;
type SubTup2VariadicWithLeadingFixedElements3Test3 = SubTup2VariadicWithLeadingFixedElements3<[a: 0, b: 1, ...c: 2[]]>;
type SubTup2VariadicWithLeadingFixedElements3Test4 = SubTup2VariadicWithLeadingFixedElements3<[a: 0, ...b: 1[]]>;
type SubTup2VariadicWithLeadingFixedElements3Test5 = SubTup2VariadicWithLeadingFixedElements3<[...a: 0[]]>;
type SubTup2VariadicWithLeadingFixedElements3Test6 = SubTup2VariadicWithLeadingFixedElements3<[a: 0, b: 1, c: 2]>;

type SubTup2VariadicWithTrailingOptionalElement<T extends unknown[]> = T extends [
  ...infer B extends [0, 1?],
  ...any
]
  ? B
  : never;

type SubTup2VariadicWithTrailingOptionalElementTest = SubTup2VariadicWithTrailingOptionalElement<[a: 0]>;
type SubTup2VariadicWithTrailingOptionalElementTest2 = SubTup2VariadicWithTrailingOptionalElement<[a: 0, b: 1]>;
type SubTup2VariadicWithTrailingOptionalElementTest3 = SubTup2VariadicWithTrailingOptionalElement<[a: 0, b: 1, c: 2]>;
type SubTup2VariadicWithTrailingOptionalElementTest4 = SubTup2VariadicWithTrailingOptionalElement<[a: 0, ...b: 1[]]>;
type SubTup2VariadicWithTrailingOptionalElementTest5 = SubTup2VariadicWithTrailingOptionalElement<[...a: 0[]]>;

type SubTup2VariadicWithLeadingOptionalElement<T extends unknown[]> = T extends [
  ...infer B extends [0?, 1],
  ...any
]
  ? B
  : never;

type SubTup2VariadicWithLeadingOptionalElementTest = SubTup2VariadicWithLeadingOptionalElement<[b: 1]>;
type SubTup2VariadicWithLeadingOptionalElementTest2 = SubTup2VariadicWithLeadingOptionalElement<[a: 0, b: 1]>;
type SubTup2VariadicWithLeadingOptionalElementTest3 = SubTup2VariadicWithLeadingOptionalElement<[a: 0, b: 1, c: 2]>;
type SubTup2VariadicWithLeadingOptionalElementTest4 = SubTup2VariadicWithLeadingOptionalElement<[a: 0, ...b: 1[]]>;
type SubTup2VariadicWithLeadingOptionalElementTest5 = SubTup2VariadicWithLeadingOptionalElement<[...a: 1[]]>;

type SubTup2TrailingVariadicWithTrailingFixedElements<T extends unknown[]> = T extends [
  ...any,
  ...infer B extends [any, any],
  any,
]
  ? B
  : never;

type SubTup2TrailingVariadicWithTrailingFixedElementsTest3 = SubTup2TrailingVariadicWithTrailingFixedElements<[...a: 0[], b: 1, c: 2]>;
type SubTup2TrailingVariadicWithTrailingFixedElementsTest4 = SubTup2TrailingVariadicWithTrailingFixedElements<[...a: 0[], b: 1]>;
type SubTup2TrailingVariadicWithTrailingFixedElementsTest5 = SubTup2TrailingVariadicWithTrailingFixedElements<[...a: 0[]]>;

type SubTup2TrailingVariadicWithTrailingFixedElements2<T extends unknown[]> = T extends [
  ...any,
  ...infer B extends [any],
  any,
  any,
]
  ? B
  : never;

type SubTup2TrailingVariadicWithTrailingFixedElements2Test = SubTup2TrailingVariadicWithTrailingFixedElements2<[...a: 0[], b: 1, c: 2, d: 3]>;
type SubTup2TrailingVariadicWithTrailingFixedElements2Test2 = SubTup2TrailingVariadicWithTrailingFixedElements2<[...a: 0[], b: 1, c: 2, d: 3, e: 4]>;
type SubTup2TrailingVariadicWithTrailingFixedElements2Test3 = SubTup2TrailingVariadicWithTrailingFixedElements2<[...a: 0[], b: 1, c: 2]>;
type SubTup2TrailingVariadicWithTrailingFixedElements2Test4 = SubTup2TrailingVariadicWithTrailingFixedElements2<[...a: 0[], b: 1]>;
type SubTup2TrailingVariadicWithTrailingFixedElements2Test5 = SubTup2TrailingVariadicWithTrailingFixedElements2<[...a: 0[]]>;

type SubTup2TrailingVariadicWithTrailingFixedElements3<T extends unknown[]> = T extends [
  ...(infer C)[],
  ...infer B extends [any, any],
  any,
]
  ? [C, B]
  : never;

type SubTup2TrailingVariadicWithTrailingFixedElements3Test = SubTup2TrailingVariadicWithTrailingFixedElements3<[...a: 0[], b: 1, c: 2, d: 3, e: 4]>;
type SubTup2TrailingVariadicWithTrailingFixedElements3Test2 = SubTup2TrailingVariadicWithTrailingFixedElements3<[...a: 0[], b: 1, c: 2, d: 3]>;
type SubTup2TrailingVariadicWithTrailingFixedElements3Test3 = SubTup2TrailingVariadicWithTrailingFixedElements3<[...a: 0[], b: 1, c: 2]>;
type SubTup2TrailingVariadicWithTrailingFixedElements3Test4 = SubTup2TrailingVariadicWithTrailingFixedElements3<[...a: 0[], b: 1]>;
type SubTup2TrailingVariadicWithTrailingFixedElements3Test5 = SubTup2TrailingVariadicWithTrailingFixedElements3<[...a: 0[]]>;
type SubTup2TrailingVariadicWithTrailingFixedElements3Test6 = SubTup2TrailingVariadicWithTrailingFixedElements3<[b: 1, c: 2, d: 3]>;

type SubTup2TrailingVariadicWithLeadingOptionalElement<T extends unknown[]> = T extends [
  ...any,
  ...infer B extends [0?, 1],
]
  ? B
  : never;

type SubTup2TrailingVariadicWithLeadingOptionalElementTest = SubTup2TrailingVariadicWithLeadingOptionalElement<[b: 1]>;
type SubTup2TrailingVariadicWithLeadingOptionalElementTest2 = SubTup2TrailingVariadicWithLeadingOptionalElement<[a: 0, b: 1]>;
type SubTup2TrailingVariadicWithLeadingOptionalElementTest3 = SubTup2TrailingVariadicWithLeadingOptionalElement<[a: 2, b: 0, c: 1]>;
type SubTup2TrailingVariadicWithLeadingOptionalElementTest4 = SubTup2TrailingVariadicWithLeadingOptionalElement<[...a: 0[], b: 1]>;
type SubTup2TrailingVariadicWithLeadingOptionalElementTest5 = SubTup2TrailingVariadicWithLeadingOptionalElement<[...a: 1[]]>;

type SubTup2TrailingVariadicWithTrailingOptionalElement<T extends unknown[]> = T extends [
  ...any,
  ...infer B extends [0, 1?],
]
  ? B
  : never;

type SubTup2TrailingVariadicWithTrailingOptionalElementTest = SubTup2TrailingVariadicWithTrailingOptionalElement<[a: 0]>;
type SubTup2TrailingVariadicWithTrailingOptionalElementTest2 = SubTup2TrailingVariadicWithTrailingOptionalElement<[a: 0, b: 1]>;
type SubTup2TrailingVariadicWithTrailingOptionalElementTest3 = SubTup2TrailingVariadicWithTrailingOptionalElement<[a: 2, b: 0, c: 1]>;
type SubTup2TrailingVariadicWithTrailingOptionalElementTest4 = SubTup2TrailingVariadicWithTrailingOptionalElement<[...a: 0[]]>;
type SubTup2TrailingVariadicWithTrailingOptionalElementTest5 = SubTup2TrailingVariadicWithTrailingOptionalElement<[...a: 0[], b: 1]>;
