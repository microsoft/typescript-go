// @strict: true
// @noEmit: true

type Target = {
  prop: string;
};
function f1(): Target {
  let v;
  v = {
    prop: "foo",
    extra: 1,
  };
  return v;
}

interface Inner { x: number }
interface U { a: Inner[] }

declare const cond: boolean;

function f2(): U {
  let v;
  if (cond) {
    v = { a: [{ x: 1, extra: true }] }
  } else {
    v = { a: [{ x: 1, extra: true }] }
  }
  return v;
}

function f3(): U {
  let v;
  if (cond) {
    v = { a: [{ x: 1, extra: true }] }
  } else {
    v = { a: [{ x: 1, extra2: true }] }
  }
  return v;
}
