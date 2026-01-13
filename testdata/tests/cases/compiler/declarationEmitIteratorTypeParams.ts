// @declaration: true
// @strict: true
// @noTypesAndSymbols: true
// @target: esnext

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
