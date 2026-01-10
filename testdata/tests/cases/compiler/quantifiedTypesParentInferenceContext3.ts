declare const f: <A>(x: { a: A, f: [<T> ((a: A) => void)] }) => void
f({ a: 0, f: [a => { a satisfies number; }] })
