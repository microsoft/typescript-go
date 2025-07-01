// @checkJs: true
// @filename: crash-on-neo-async.js
// This test reproduces a crash that occurs when parsing files like webpack/node_modules/neo-async/async.js
// The crash happens in GetAssignmentDeclarationKind when checking IsIdentifier(bin.Left.Name())
// where bin.Left.Name() returns nil for ElementAccessExpression
// Pattern that causes the crash - element access assignment
var obj = {};
var prop = 'test';
obj[prop] = function () {
    // This assignment with element access should not crash
    return 42;
};
