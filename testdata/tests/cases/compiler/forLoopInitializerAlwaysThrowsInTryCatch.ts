// @strict: true
// @noEmit: true

// Repro from github.com/microsoft/TypeScript/issues/63092
// A for-loop initializer that always throws, wrapped in try/catch,
// should not cause a stack overflow in control flow analysis.

function test1(v: 0 | 1 | 2) {
  try {
    for (
      (function () {
        throw new Error("");
      })();
      v;
      v++
    ) {}
  } catch (e) {}
  v;
}

function test2() {
  try {
    for (
      (function () { throw "1"; })();
      (function* () { throw "2"; })();
      (function* () { throw "3"; })()
    ) {}
  } catch (e) {}
}

function test3() {
  let x = 0;
  while (true) {
    try {
      for (
        (function () { throw new Error(""); })();
        true;
        x++
      ) {
        break;
        continue;
      }
    } catch (e) {}
  }
}

function test4() {
  let x = 0;
  outer: for (; x < 10; x++) {
    try {
      inner: for (
        (function () { throw new Error(""); })();
        true;
        x++
      ) {
        continue inner;
      }
    } catch (e) {}
  }
}


