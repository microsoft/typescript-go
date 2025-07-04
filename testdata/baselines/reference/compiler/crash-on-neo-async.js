//// [tests/cases/compiler/crash-on-neo-async.ts] ////

//// [crash-on-neo-async.js]
// This test reproduces a crash that occurs when parsing files like webpack/node_modules/neo-async/async.js
// The crash happens in GetAssignmentDeclarationKind when checking IsIdentifier(bin.Left.Name())
// where bin.Left.Name() returns nil for ElementAccessExpression
// Pattern that causes the crash - element access assignment
var obj = {};
var prop = 'test';
// This assignment with element access should not crash
// It previously crashed because ElementAccessExpression.Name() returns nil
// and IsIdentifier was called on that nil value
obj[prop] = function () {
    return 42;
};
// Property access assignment should work fine (has a valid Name())
obj.prop2 = function () {
    return 43;
};
// Nested element access assignment
obj['nested'][prop] = function () {
    return 44;
};


//// [crash-on-neo-async.js]
// This test reproduces a crash that occurs when parsing files like webpack/node_modules/neo-async/async.js
// The crash happens in GetAssignmentDeclarationKind when checking IsIdentifier(bin.Left.Name())
// where bin.Left.Name() returns nil for ElementAccessExpression
// Pattern that causes the crash - element access assignment
var obj = {};
var prop = 'test';
// This assignment with element access should not crash
// It previously crashed because ElementAccessExpression.Name() returns nil
// and IsIdentifier was called on that nil value
obj[prop] = function () {
    return 42;
};
// Property access assignment should work fine (has a valid Name())
obj.prop2 = function () {
    return 43;
};
// Nested element access assignment
obj['nested'][prop] = function () {
    return 44;
};
