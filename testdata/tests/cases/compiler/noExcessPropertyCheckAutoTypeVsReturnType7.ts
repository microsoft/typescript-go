// @strict: true
// @noEmit: true

interface Inner {
  x: number;
}

interface U {
  a: Inner[];
}

declare const cond: boolean;

function f(): U[] {
  let v;
  if (cond) {
    v = [{ a: [{ x: 1, extra: true }] }];
  } else {
    const tmp = [];
    tmp.push({ a: [{ x: 1, extra: true }] });
    v = tmp;
  }
  return v;
}
