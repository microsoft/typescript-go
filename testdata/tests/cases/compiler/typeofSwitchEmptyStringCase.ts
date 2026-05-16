// @strict: true

function f(x: string | number) {
  switch (typeof x) {
    case "":
    case "string":
      x.charAt(0);
      break;
  }
}
