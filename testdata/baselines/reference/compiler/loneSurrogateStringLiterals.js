//// [tests/cases/compiler/loneSurrogateStringLiterals.ts] ////

//// [loneSurrogateStringLiterals.ts]
// Lone surrogates should be distinct string literal types
const highSurrogate: "\uD800" = "\uD800"; // ok
const lowSurrogate: "\uDC00" = "\uDC00"; // ok

// These should be errors - different surrogates are not assignable to each other
const highToLow: "\uD800" = "\uDC00"; // error
const lowToHigh: "\uDC00" = "\uD800"; // error

// Different high surrogates should also be distinct
const high1: "\uD800" = "\uD801"; // error
const high2: "\uD801" = "\uD800"; // error

// Different low surrogates should also be distinct
const low1: "\uDC00" = "\uDC01"; // error
const low2: "\uDC01" = "\uDC00"; // error

// Extended Unicode escape syntax should also work
const extHigh: "\u{D800}" = "\u{D800}"; // ok
const extLow: "\u{DC00}" = "\u{DC00}"; // ok
const extHighToLow: "\u{D800}" = "\u{DC00}"; // error
const extLowToHigh: "\u{DC00}" = "\u{D800}"; // error

// Mixed syntax should also be equivalent
const mixedHigh: "\uD800" = "\u{D800}"; // ok
const mixedLow: "\u{DC00}" = "\uDC00"; // ok
const mixedError1: "\uD800" = "\u{DC00}"; // error
const mixedError2: "\u{D800}" = "\uDC00"; // error


//// [loneSurrogateStringLiterals.js]
"use strict";
// Lone surrogates should be distinct string literal types
const highSurrogate = "\uD800"; // ok
const lowSurrogate = "\uDC00"; // ok
// These should be errors - different surrogates are not assignable to each other
const highToLow = "\uDC00"; // error
const lowToHigh = "\uD800"; // error
// Different high surrogates should also be distinct
const high1 = "\uD801"; // error
const high2 = "\uD800"; // error
// Different low surrogates should also be distinct
const low1 = "\uDC01"; // error
const low2 = "\uDC00"; // error
// Extended Unicode escape syntax should also work
const extHigh = "\u{D800}"; // ok
const extLow = "\u{DC00}"; // ok
const extHighToLow = "\u{DC00}"; // error
const extLowToHigh = "\u{D800}"; // error
// Mixed syntax should also be equivalent
const mixedHigh = "\u{D800}"; // ok
const mixedLow = "\uDC00"; // ok
const mixedError1 = "\u{DC00}"; // error
const mixedError2 = "\uDC00"; // error
