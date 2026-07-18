// @strict: true
// @noEmit: true

// https://github.com/microsoft/TypeScript/issues/63627

function arr1(a: [number, ...NoInfer<[]>]) {}
arr1([1]);

function arr2(a: [number, ...[]]) {}
arr2([1]);

function fun1(a: (x: number, ...o: NoInfer<[]>) => void) {}
fun1(x => {});

function fun2(a: (x: number, ...o: []) => void) {}
fun2(x => {});

function func1<A extends unknown[]>(
  args: A,
  fn: (x: number, ...args: NoInfer<A>) => void,
) {}

func1([] as const, x => {});

function func2<A extends unknown[]>(
  args: A,
  fn: (...args: NoInfer<A>) => void,
) {}

function foo(u: number, v: number) {
  func2([u, v] as const, () => {});
  func2([u, v] as const, (x: number) => {});
  func2([u, v] as const, (x: number, y: number) => {});
}
