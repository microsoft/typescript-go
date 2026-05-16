//// [tests/cases/compiler/typeofSwitchEmptyStringCase.ts] ////

//// [typeofSwitchEmptyStringCase.ts]
function f(x: string | number) {
  switch (typeof x) {
    case "":
    case "string":
      x.charAt(0);
      break;
  }
}


//// [typeofSwitchEmptyStringCase.js]
"use strict";
function f(x) {
    switch (typeof x) {
        case "":
        case "string":
            x.charAt(0);
            break;
    }
}
