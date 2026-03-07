//// [tests/cases/compiler/declarationEmitIteratorTypeParams.ts] ////

//// [declarationEmitIteratorTypeParams.ts]
// Iterator with two type params should not get third param
export const fromIterator = <A, L>(
  iterator: () => Iterator<A, L>
): Array<A> => {
  const result: Array<A> = []
  let next = iterator().next()
  while (!next.done) {
    result.push(next.value as A)
    next = iterator().next()
  }
  return result
}

// Generator function
export function* counted(): Generator<number, string> {
  yield 1
  yield 2
  yield 3
  return "done"
}

// Iterable with return type
export const fromIterable = <A, L>(iterable: Iterable<A, L>): Array<A> => Array.from(iterable as Iterable<A>)

// AsyncIterator
export const fromAsyncIterator = <A, L>(
  iterator: () => AsyncIterator<A, L>
): Promise<Array<A>> => Promise.resolve([])


//// [declarationEmitIteratorTypeParams.js]
// Iterator with two type params should not get third param
export const fromIterator = (iterator) => {
    const result = [];
    let next = iterator().next();
    while (!next.done) {
        result.push(next.value);
        next = iterator().next();
    }
    return result;
};
// Generator function
export function* counted() {
    yield 1;
    yield 2;
    yield 3;
    return "done";
}
// Iterable with return type
export const fromIterable = (iterable) => Array.from(iterable);
// AsyncIterator
export const fromAsyncIterator = (iterator) => Promise.resolve([]);


//// [declarationEmitIteratorTypeParams.d.ts]
export declare const fromIterator: <A, L>(iterator: () => Iterator<A, L>) => Array<A>;
export declare function counted(): Generator<number, string>;
export declare const fromIterable: <A, L>(iterable: Iterable<A, L>) => Array<A>;
export declare const fromAsyncIterator: <A, L>(iterator: () => AsyncIterator<A, L>) => Promise<Array<A>>;
