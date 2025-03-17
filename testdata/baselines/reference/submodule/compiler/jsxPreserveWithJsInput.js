//// [tests/cases/compiler/jsxPreserveWithJsInput.ts] ////

//// [a.js]
var elemA = 42;

//// [b.jsx]
var elemB = <b>{"test"}</b>;

//// [c.js]
var elemC = <c>{42}</c>;

//// [d.ts]
var elemD = 42;

//// [e.tsx]
var elemE = <e>{true}</e>;


//// [e.jsx]
var elemE = <e>{true}</e>;
//// [d.js]
var elemD = 42;
//// [c.js]
var elemC = <c>{42}</c>;
//// [b.jsx]
var elemB = <b>{"test"}</b>;
//// [a.js]
var elemA = 42;
