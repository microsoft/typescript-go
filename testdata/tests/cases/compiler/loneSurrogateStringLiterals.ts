// @strict: true

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
