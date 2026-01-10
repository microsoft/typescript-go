declare const f: <A>(a: A, f: [<T> ((a: A) => void)]) => void
f(0, [a => {}])