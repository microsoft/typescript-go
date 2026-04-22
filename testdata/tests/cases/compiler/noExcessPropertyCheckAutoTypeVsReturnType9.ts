// @strict: true
// @noEmit: true

interface Inner { x: number }
interface U { a: Inner[] }

declare const cond: boolean;

function f1(): U {
  let v;
  if (cond) {
    v = { a: [{ x: 1, extra: true }] }
  } else {
    v = { a: [{ x: 1, extra: true }] }
  }
  return v;
}

function f2(): U {
  let v;
  if (cond) {
    v = { a: [{ x: 1, extra: true }] }
  } else {
    v = { a: [{ x: 1, extra2: true }] }
  }
  return v;
}
