package ls

// movetofile.go implements the core types and logic for the "Move to file"
// and "Move to a new file" refactorings.

import (
"slices"

"github.com/microsoft/typescript-go/internal/ast"
"github.com/microsoft/typescript-go/internal/astnav"
"github.com/microsoft/typescript-go/internal/checker"
"github.com/microsoft/typescript-go/internal/core"
)

// StatementRange represents a contiguous range of statements in a source file.
type StatementRange struct {
First     *ast.Node // The first statement in the range
AfterLast *ast.Node // The statement after the last in the range (may be nil)
}

// ToMove holds the set of statements to move and their ranges in the original file.
type ToMove struct {
// All statements that will be moved.
All []*ast.Node
// Ranges of those statements in the original file.
Ranges []StatementRange
}

// UsageInfo captures the symbol usage analysis needed to update imports/exports.
type UsageInfo struct {
// Symbols whose declarations are moved from the old file to the new file.
MovedSymbols map[*ast.Symbol]bool

// Symbols declared in the old file that must be imported by the new file. (May not already be exported.)
// Maps symbol -> isValidTypeOnlyUseSite
TargetFileImportsFromOldFile map[*ast.Symbol]bool

// Subset of movedSymbols that are still used elsewhere in the old file and must be imported back.
// Maps symbol -> isValidTypeOnlyUseSite
OldFileImportsFromTargetFile map[*ast.Symbol]bool

// Imports in old file needed by moved code.
// Maps symbol -> (isValidTypeOnlyUseSite, import declaration node)
OldImportsNeededByTargetFile map[*ast.Symbol]importInfo

// Subset of oldImportsNeededByTargetFile that will no longer be used in the old file.
UnusedImportsFromOldFile map[*ast.Symbol]bool
}

// importInfo holds info about an import symbol.
type importInfo struct {
isValidTypeOnly bool
declaration     *ast.Node // the import declaration node (ImportDeclaration, etc.)
}

// rangeToMove is an internal helper type for getRangeToMove.
type rangeToMove struct {
toMove    []*ast.Node
afterLast *ast.Node
}

// getRangeToMove finds the contiguous block of statements in the source file
// that overlap with the given span.
func getRangeToMove(file *ast.SourceFile, span core.TextRange) *rangeToMove {
statements := file.Statements.Nodes

// Find the first statement whose end is after span.Pos().
startNodeIndex := -1
for i, s := range statements {
if s.End() > span.Pos() {
startNodeIndex = i
break
}
}
if startNodeIndex == -1 {
return nil
}

// Find the last statement whose end is >= span.End(), starting from startNodeIndex.
endNodeIndex := -1
for i := startNodeIndex; i < len(statements); i++ {
if statements[i].End() >= span.End() {
endNodeIndex = i
break
}
}

// If the range ends before the start of the end statement, go back one.
if endNodeIndex != -1 && span.End() <= astnav.GetStartOfNode(statements[endNodeIndex], file, false) {
endNodeIndex--
}

var toMove []*ast.Node
var afterLast *ast.Node
if endNodeIndex == -1 {
toMove = statements[startNodeIndex:]
} else {
toMove = statements[startNodeIndex : endNodeIndex+1]
if endNodeIndex+1 < len(statements) {
afterLast = statements[endNodeIndex+1]
}
}

return &rangeToMove{toMove: toMove, afterLast: afterLast}
}

// getRangesWhere calls cb for each contiguous range of elements in arr for which pred returns true.
// cb receives (startIndex, afterEndIndex).
func getRangesWhere(arr []*ast.Node, pred func(*ast.Node) bool, cb func(start, afterEnd int)) {
start := -1
for i, elem := range arr {
if pred(elem) {
if start == -1 {
start = i
}
} else {
if start != -1 {
cb(start, i)
start = -1
}
}
}
if start != -1 {
cb(start, len(arr))
}
}

// getStatementsToMove determines which statements in the given span should be
// moved, filtering out imports and prologue directives.
func getStatementsToMove(file *ast.SourceFile, span core.TextRange) *ToMove {
r := getRangeToMove(file, span)
if r == nil {
return nil
}

var all []*ast.Node
var ranges []StatementRange

getRangesWhere(r.toMove, isAllowedStatementToMove, func(start, afterEnd int) {
for i := start; i < afterEnd; i++ {
all = append(all, r.toMove[i])
}
var afterLast *ast.Node
if afterEnd < len(r.toMove) {
afterLast = r.toMove[afterEnd]
} else {
afterLast = r.afterLast
}
ranges = append(ranges, StatementRange{
First:     r.toMove[start],
AfterLast: afterLast,
})
})

if len(all) == 0 {
return nil
}
return &ToMove{All: all, Ranges: ranges}
}

// isAllowedStatementToMove returns true if a statement can be moved to a new
// file. Pure imports and prologue directives are excluded.
func isAllowedStatementToMove(statement *ast.Node) bool {
return !isPureImport(statement) && !ast.IsPrologueDirective(statement)
}

// isPureImport returns true if a node is a pure import statement (not exported).
func isPureImport(node *ast.Node) bool {
if node == nil {
return false
}
switch node.Kind {
case ast.KindImportDeclaration:
return true
case ast.KindImportEqualsDeclaration:
return !ast.HasSyntacticModifier(node, ast.ModifierFlagsExport)
case ast.KindVariableStatement:
decls := node.AsVariableStatement().DeclarationList.Declarations.Nodes
return len(decls) > 0 && core.Every(decls, func(d *ast.Node) bool {
vd := d.AsVariableDeclaration()
return vd.Initializer != nil && ast.IsRequireCall(vd.Initializer, true)
})
default:
return false
}
}

// isInImport returns true if a declaration node is inside an import statement.
func isInImport(decl *ast.Node) bool {
switch decl.Kind {
case ast.KindImportEqualsDeclaration, ast.KindImportSpecifier, ast.KindImportClause, ast.KindNamespaceImport:
return true
case ast.KindVariableDeclaration:
return isVariableDeclarationInImport(decl)
case ast.KindBindingElement:
if decl.Parent != nil && decl.Parent.Parent != nil && ast.IsVariableDeclaration(decl.Parent.Parent) {
return isVariableDeclarationInImport(decl.Parent.Parent)
}
return false
default:
return false
}
}

// isVariableDeclarationInImport returns true if the variable declaration is a require() call at the top level.
func isVariableDeclarationInImport(decl *ast.Node) bool {
vd := decl.AsVariableDeclaration()
// parent = VariableDeclarationList, grandparent = VariableStatement, great-grandparent = SourceFile
if decl.Parent == nil || decl.Parent.Parent == nil || decl.Parent.Parent.Parent == nil {
return false
}
return ast.IsSourceFile(decl.Parent.Parent.Parent) &&
vd.Initializer != nil && ast.IsRequireCall(vd.Initializer, true)
}

// isTopLevelDeclaration checks if a node is a top-level declaration.
func isTopLevelDeclaration(node *ast.Node) bool {
if node == nil || node.Parent == nil {
return false
}
switch node.Kind {
case ast.KindFunctionDeclaration, ast.KindClassDeclaration, ast.KindModuleDeclaration,
ast.KindEnumDeclaration, ast.KindTypeAliasDeclaration, ast.KindInterfaceDeclaration,
ast.KindImportEqualsDeclaration:
return ast.IsSourceFile(node.Parent)
case ast.KindVariableDeclaration:
// parent = VariableDeclarationList, grandparent = VariableStatement, great-grandparent = SourceFile
return node.Parent != nil && node.Parent.Parent != nil &&
node.Parent.Parent.Parent != nil && ast.IsSourceFile(node.Parent.Parent.Parent)
}
return false
}

// forEachTopLevelDeclarationInStatement calls cb for each top-level declaration symbol in a statement.
func forEachTopLevelDeclarationInStatement(statement *ast.Node, cb func(*ast.Symbol)) {
switch statement.Kind {
case ast.KindFunctionDeclaration, ast.KindClassDeclaration, ast.KindModuleDeclaration,
ast.KindEnumDeclaration, ast.KindTypeAliasDeclaration, ast.KindInterfaceDeclaration,
ast.KindImportEqualsDeclaration:
if sym := statement.Symbol(); sym != nil {
cb(sym)
}
case ast.KindVariableStatement:
for _, decl := range statement.AsVariableStatement().DeclarationList.Declarations.Nodes {
forEachTopLevelDeclarationInBindingName(decl.AsVariableDeclaration().Name, cb)
}
}
}

func forEachTopLevelDeclarationInBindingName(name *ast.Node, cb func(*ast.Symbol)) {
if name == nil {
return
}
switch name.Kind {
case ast.KindIdentifier:
if name.Parent != nil {
if sym := name.Parent.Symbol(); sym != nil {
cb(sym)
}
}
case ast.KindArrayBindingPattern, ast.KindObjectBindingPattern:
for _, elem := range name.AsBindingPattern().Elements.Nodes {
if !ast.IsOmittedExpression(elem) {
be := elem.AsBindingElement()
forEachTopLevelDeclarationInBindingName(be.Name(), cb)
}
}
}
}

// forEachReference calls onReference for every non-declaration identifier reference in node.
func forEachReference(node *ast.Node, c *checker.Checker, onReference func(sym *ast.Symbol, isValidTypeOnly bool)) {
var visit func(n *ast.Node) bool
visit = func(n *ast.Node) bool {
if ast.IsIdentifier(n) && !ast.IsDeclarationName(n) {
sym := c.GetSymbolAtLocation(n)
if sym != nil {
onReference(sym, ast.IsValidTypeOnlyAliasUseSite(n))
}
}
n.ForEachChild(visit)
return false
}
node.ForEachChild(visit)
}

// getUsageInfo analyzes symbol usage across the moved statements.
func getUsageInfo(oldFile *ast.SourceFile, toMove []*ast.Node, c *checker.Checker) *UsageInfo {
movedSymbols := make(map[*ast.Symbol]bool)
oldImportsNeededByTargetFile := make(map[*ast.Symbol]importInfo)
targetFileImportsFromOldFile := make(map[*ast.Symbol]bool)

// Find all symbols declared in the moved statements.
for _, stmt := range toMove {
forEachTopLevelDeclarationInStatement(stmt, func(sym *ast.Symbol) {
movedSymbols[sym] = true
})
}

toMoveSet := make(map[*ast.Node]bool, len(toMove))
for _, s := range toMove {
toMoveSet[s] = true
}

// For each identifier reference in moved statements, classify it.
for _, stmt := range toMove {
forEachReference(stmt, c, func(sym *ast.Symbol, isValidTypeOnly bool) {
if sym == nil || len(sym.Declarations) == 0 {
return
}
// Skip symbols that are being moved.
if movedSymbols[sym] {
return
}
// Check if this symbol is from an import.
importedDecl := findDeclarationInImport(sym)
if importedDecl != nil {
// This symbol is imported in the old file - new file needs it too.
if existing, ok := oldImportsNeededByTargetFile[sym]; ok {
if !isValidTypeOnly {
existing.isValidTypeOnly = false
oldImportsNeededByTargetFile[sym] = existing
}
} else {
oldImportsNeededByTargetFile[sym] = importInfo{
isValidTypeOnly: isValidTypeOnly,
declaration:     importedDecl,
}
}
} else if isFromFile(sym, oldFile) {
// Symbol is declared in the old file (not an import, not moved).
// The new file needs to import it.
prev := targetFileImportsFromOldFile[sym]
targetFileImportsFromOldFile[sym] = prev || isValidTypeOnly
}
})
}

// Find which imports are still used in the remaining old file.
unusedImportsFromOldFile := make(map[*ast.Symbol]bool)
for sym := range oldImportsNeededByTargetFile {
unusedImportsFromOldFile[sym] = true
}

oldFileImportsFromTargetFile := make(map[*ast.Symbol]bool)

for _, stmt := range oldFile.Statements.Nodes {
if toMoveSet[stmt] {
continue
}
forEachReference(stmt, c, func(sym *ast.Symbol, isValidTypeOnly bool) {
// If this symbol is a moved one, old file needs to import it from new file.
if movedSymbols[sym] {
prev := oldFileImportsFromTargetFile[sym]
oldFileImportsFromTargetFile[sym] = prev || isValidTypeOnly
}
// This symbol is still used in old file, so it's not unused.
delete(unusedImportsFromOldFile, sym)
})
}

return &UsageInfo{
MovedSymbols:                movedSymbols,
TargetFileImportsFromOldFile: targetFileImportsFromOldFile,
OldFileImportsFromTargetFile: oldFileImportsFromTargetFile,
OldImportsNeededByTargetFile: oldImportsNeededByTargetFile,
UnusedImportsFromOldFile:    unusedImportsFromOldFile,
}
}

// findDeclarationInImport returns the import declaration node if any of the symbol's
// declarations is inside an import statement.
func findDeclarationInImport(sym *ast.Symbol) *ast.Node {
for _, decl := range sym.Declarations {
if isInImport(decl) {
// Walk up to find the top-level import statement.
n := decl
for n != nil && !ast.IsSourceFile(n) {
if ast.IsImportDeclaration(n) || ast.IsImportEqualsDeclaration(n) || ast.IsVariableStatement(n) {
return n
}
n = n.Parent
}
return decl
}
}
return nil
}

// isFromFile returns true if the symbol is declared in the given source file.
func isFromFile(sym *ast.Symbol, file *ast.SourceFile) bool {
for _, decl := range sym.Declarations {
if getSourceFileOfDeclaration(decl) == file.AsNode() {
return true
}
}
return false
}

// getSourceFileOfDeclaration walks up the parent chain to find the SourceFile.
func getSourceFileOfDeclaration(decl *ast.Node) *ast.Node {
n := decl
for n != nil {
if ast.IsSourceFile(n) {
return n
}
n = n.Parent
}
return nil
}

// containsNode returns true if the nodes slice contains n.
func containsNode(nodes []*ast.Node, n *ast.Node) bool {
return slices.Contains(nodes, n)
}
