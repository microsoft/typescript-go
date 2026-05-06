//// [tests/cases/compiler/typeInterfaceDeclarationsInBlockStatements1.ts] ////

//// [typeInterfaceDeclarationsInBlockStatements1.ts]
// https://github.com/microsoft/TypeScript/issues/60175

function f1() {
  if (true) type s = string;
  console.log("" as s);
}

function f2() {
  if (true) {
    type s = string;
  }
  console.log("" as s);
}

function f3() {
  if (true)
    interface s {
      length: number;
    }
  console.log("" as s);
}

function f4() {
  if (true) {
    interface s {
      length: number;
    }
  }
  console.log("" as s);
}


//// [typeInterfaceDeclarationsInBlockStatements1.js]
"use strict";
// https://github.com/microsoft/TypeScript/issues/60175
function f1() {
    if (true)
        ;
    console.log("");
}
function f2() {
    if (true) {
    }
    console.log("");
}
function f3() {
    if (true)
        ;
    console.log("");
}
function f4() {
    if (true) {
    }
    console.log("");
}


//// [typeInterfaceDeclarationsInBlockStatements1.d.ts]
function f1(): void;
function f2(): void;
function f3(): void;
function f4(): void;
