//// [tests/cases/compiler/regularExpressionBackslashK.ts] ////

//// [regularExpressionBackslashK.ts]
// Test that \k without < is an error when named groups are present

// Valid: \k followed by <name> with named groups
const validBackref = /(?<foo>a)\k<foo>/;

// Invalid: \k not followed by < when named groups are present (even in non-Unicode mode)
const invalidK = /(?<foo>a)\k/;

// Invalid: \k followed by other chars when named groups are present
const invalidKWithText = /(?<bar>x)\kb/;

// Valid: \k without < is OK when there are NO named groups (identity escape)
const validIdentityEscape = /a\kb/;

// Invalid: \k without < in Unicode mode (regardless of named groups)
const invalidKUnicode = /a\kb/u;

// Edge cases

// Multiple named groups, valid backreferences
const multiGroup = /(?<a>x)(?<b>y)\k<a>\k<b>/;

// Named group in alternation with \k in different branch
const alternation = /(?<foo>a)|\k<foo>/;

// Named group with \k in lookahead
const lookahead = /(?<bar>b)(?=\k<bar>)/;

// Named group with bare \k - should error
const bareKWithGroups = /(?<x>.)(?<y>.)\k/;

// Bare \k at start when named group comes later - should still error
const bareKBeforeGroup = /\k(?<name>pattern)/;

// Identity escape \k is valid when no named groups at all
const noNamedGroups = /\ka\kb/;

// Unicode characters

// Unicode characters before named group
const unicodeBefore = /😀(?<foo>a)\k<foo>/;

// Unicode characters after named group
const unicodeAfter = /(?<bar>b)\k<bar>😀/;

// Unicode characters in between
const unicodeMiddle = /(?<x>.)😀\k<x>/;

// Unicode with bare \k - should error
const unicodeWithBareK = /😀(?<name>.)\k/;

// Unicode without named groups and \k - should be OK
const unicodeNoGroups = /😀\k😀/;


//// [regularExpressionBackslashK.js]
// Test that \k without < is an error when named groups are present
// Valid: \k followed by <name> with named groups
const validBackref = /(?<foo>a)\k<foo>/;
// Invalid: \k not followed by < when named groups are present (even in non-Unicode mode)
const invalidK = /(?<foo>a)\k/;
// Invalid: \k followed by other chars when named groups are present
const invalidKWithText = /(?<bar>x)\kb/;
// Valid: \k without < is OK when there are NO named groups (identity escape)
const validIdentityEscape = /a\kb/;
// Invalid: \k without < in Unicode mode (regardless of named groups)
const invalidKUnicode = /a\kb/u;
// Edge cases
// Multiple named groups, valid backreferences
const multiGroup = /(?<a>x)(?<b>y)\k<a>\k<b>/;
// Named group in alternation with \k in different branch
const alternation = /(?<foo>a)|\k<foo>/;
// Named group with \k in lookahead
const lookahead = /(?<bar>b)(?=\k<bar>)/;
// Named group with bare \k - should error
const bareKWithGroups = /(?<x>.)(?<y>.)\k/;
// Bare \k at start when named group comes later - should still error
const bareKBeforeGroup = /\k(?<name>pattern)/;
// Identity escape \k is valid when no named groups at all
const noNamedGroups = /\ka\kb/;
// Unicode characters
// Unicode characters before named group
const unicodeBefore = /😀(?<foo>a)\k<foo>/;
// Unicode characters after named group
const unicodeAfter = /(?<bar>b)\k<bar>😀/;
// Unicode characters in between
const unicodeMiddle = /(?<x>.)😀\k<x>/;
// Unicode with bare \k - should error
const unicodeWithBareK = /😀(?<name>.)\k/;
// Unicode without named groups and \k - should be OK
const unicodeNoGroups = /😀\k😀/;
