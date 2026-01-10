package organizeimports

import (
	"cmp"
	"context"
	"math"
	"slices"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/change"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/scanner"
	"github.com/microsoft/typescript-go/internal/stringutil"
	"github.com/microsoft/typescript-go/internal/tspath"
)

var (
	caseInsensitiveOrganizeImportsComparer = []func(a, b string) int{getOrganizeImportsOrdinalStringComparer(true)}
	caseSensitiveOrganizeImportsComparer   = []func(a, b string) int{getOrganizeImportsOrdinalStringComparer(false)}
	organizeImportsComparers               = []func(a, b string) int{
		caseInsensitiveOrganizeImportsComparer[0],
		caseSensitiveOrganizeImportsComparer[0],
	}
)

type OrganizeImportsMode int

const (
	OrganizeImportsModeAll            OrganizeImportsMode = 0
	OrganizeImportsModeSortAndCombine OrganizeImportsMode = 1
	OrganizeImportsModeRemoveUnused   OrganizeImportsMode = 2
)

// OrganizeImports organizes imports by:
//  1. Removing unused imports
//  2. Coalescing imports from the same module
//  3. Sorting imports
func OrganizeImports(
	ctx context.Context,
	sourceFile *ast.SourceFile,
	changeTracker *change.Tracker,
	program *compiler.Program,
	preferences *lsutil.UserPreferences,
	mode OrganizeImportsMode,
) {
	shouldSort := mode == OrganizeImportsModeSortAndCombine || mode == OrganizeImportsModeAll
	shouldCombine := shouldSort
	shouldRemove := mode == OrganizeImportsModeRemoveUnused || mode == OrganizeImportsModeAll

	topLevelImportDecls := filterImportDeclarations(sourceFile.Statements.Nodes)
	topLevelImportGroupDecls := groupByNewlineContiguous(sourceFile, topLevelImportDecls)

	comparersToTest, typeOrdersToTest := getDetectionLists(preferences)
	defaultComparer := comparersToTest[0]

	moduleSpecifierComparer := defaultComparer
	namedImportComparer := defaultComparer
	typeOrder := lsutil.OrganizeImportsTypeOrderLast
	if preferences != nil && preferences.OrganizeImportsTypeOrder != lsutil.OrganizeImportsTypeOrderAuto {
		typeOrder = preferences.OrganizeImportsTypeOrder
	}

	if preferences == nil || preferences.OrganizeImportsIgnoreCase.IsUnknown() {
		result := detectModuleSpecifierCaseBySort(topLevelImportGroupDecls, comparersToTest)
		moduleSpecifierComparer = result.comparer
	}

	if typeOrder == lsutil.OrganizeImportsTypeOrderAuto || (preferences != nil && preferences.OrganizeImportsIgnoreCase.IsUnknown()) {
		namedImportSort := detectNamedImportOrganizationBySort(topLevelImportDecls, comparersToTest, typeOrdersToTest)
		if namedImportSort != nil {
			if namedImportComparer == nil || (preferences != nil && preferences.OrganizeImportsIgnoreCase.IsUnknown()) {
				namedImportComparer = namedImportSort.namedImportComparer
			}
			if typeOrder == lsutil.OrganizeImportsTypeOrderAuto {
				typeOrder = namedImportSort.typeOrder
			}
		}
	}

	comparer := comparerSettings{
		moduleSpecifierComparer: moduleSpecifierComparer,
		namedImportComparer:     namedImportComparer,
		typeOrder:               typeOrder,
	}

	for _, importGroupDecl := range topLevelImportGroupDecls {
		organizeImportsWorker(importGroupDecl, comparer, shouldSort, shouldCombine, shouldRemove, sourceFile, program, changeTracker, ctx)
	}

	// TODO: Handle exports when mode != RemoveUnused

	// TODO: Handle ambient modules
}

type comparerSettings struct {
	moduleSpecifierComparer func(a, b string) int
	namedImportComparer     func(a, b string) int
	typeOrder               lsutil.OrganizeImportsTypeOrder
}

func organizeImportsWorker(
	oldImportDecls []*ast.Statement,
	comparer comparerSettings,
	shouldSort bool,
	shouldCombine bool,
	shouldRemove bool,
	sourceFile *ast.SourceFile,
	program *compiler.Program,
	changeTracker *change.Tracker,
	ctx context.Context,
) {
	if len(oldImportDecls) == 0 {
		return
	}

	// TODO: Set EmitFlags.NoLeadingComments on first import

	processedImports := oldImportDecls
	if shouldRemove {
		processedImports = removeUnusedImports(processedImports, sourceFile, program, ctx)
	}

	var newImportDecls []*ast.Statement
	if shouldCombine {
		grouped := groupByModuleSpecifier(processedImports)
		if shouldSort {
			slices.SortFunc(grouped, func(a, b []*ast.Statement) int {
				if len(a) == 0 || len(b) == 0 {
					return 0
				}
				return compareModuleSpecifiersWorker(
					a[0].ModuleSpecifier(),
					b[0].ModuleSpecifier(),
					comparer.moduleSpecifierComparer,
				)
			})
		}

		specifierComparer := GetNamedImportSpecifierComparer(
			&lsutil.UserPreferences{OrganizeImportsTypeOrder: comparer.typeOrder},
			comparer.namedImportComparer,
		)

		for _, importGroup := range grouped {
			coalesced := coalesceImportsWorker(importGroup, comparer.moduleSpecifierComparer, specifierComparer, sourceFile)
			newImportDecls = append(newImportDecls, coalesced...)
		}
	} else {
		newImportDecls = processedImports
	}

	if shouldSort && !shouldCombine {
		slices.SortFunc(newImportDecls, func(a, b *ast.Statement) int {
			return CompareImportsOrRequireStatements(a, b, comparer.moduleSpecifierComparer)
		})
	}

	// TODO: Apply changes using changeTracker
	_ = newImportDecls
}

func filterImportDeclarations(statements []*ast.Statement) []*ast.Statement {
	var result []*ast.Statement
	for _, stmt := range statements {
		if stmt.Kind == ast.KindImportDeclaration {
			result = append(result, stmt)
		}
	}
	return result
}

func groupByModuleSpecifier(imports []*ast.Statement) [][]*ast.Statement {
	groups := make(map[string][]*ast.Statement)
	var order []string

	for _, imp := range imports {
		specifier := getExternalModuleName(imp.ModuleSpecifier())
		if _, exists := groups[specifier]; !exists {
			order = append(order, specifier)
		}
		groups[specifier] = append(groups[specifier], imp)
	}

	result := make([][]*ast.Statement, 0, len(order))
	for _, key := range order {
		result = append(result, groups[key])
	}
	return result
}

func removeUnusedImports(oldImports []*ast.Statement, sourceFile *ast.SourceFile, program *compiler.Program, ctx context.Context) []*ast.Statement {
	typeChecker, done := program.GetTypeCheckerForFile(ctx, sourceFile)
	defer done()

	// TODO: Get JSX namespace and fragment factory
	jsxNamespace := typeChecker.GetJsxNamespace(sourceFile.AsNode())
	_ = jsxNamespace

	// TODO: Check if JSX elements are present
	jsxElementsPresent := false
	_ = jsxElementsPresent

	usedImports := make([]*ast.Statement, 0, len(oldImports))

	for _, importDecl := range oldImports {
		importClause := importDecl.AsImportDeclaration().ImportClause
		if importClause == nil {
			usedImports = append(usedImports, importDecl)
			continue
		}

		clause := importClause.AsImportClause()
		name := clause.Name()
		namedBindings := clause.NamedBindings

		if name != nil && !isDeclarationUsed(name.AsIdentifier(), typeChecker, sourceFile) {
			name = nil
		}

		if namedBindings != nil {
			if namedBindings.Kind == ast.KindNamespaceImport {
				nsImport := namedBindings.AsNamespaceImport()
				if !isDeclarationUsed(nsImport.Name().AsIdentifier(), typeChecker, sourceFile) {
					namedBindings = nil
				}
			} else if namedBindings.Kind == ast.KindNamedImports {
				namedImports := namedBindings.AsNamedImports()
				newElements := filterUsedImportSpecifiers(namedImports.Elements.Nodes, typeChecker, sourceFile)
				if len(newElements) == 0 {
					namedBindings = nil
				} else if len(newElements) < len(namedImports.Elements.Nodes) {
					// TODO: Create updated named imports
					_ = newElements
				}
			}
		}

		if name != nil || namedBindings != nil {
			// TODO: Create updated import declaration
			usedImports = append(usedImports, importDecl)
		} else {
			moduleSpecifier := importDecl.ModuleSpecifier()
			if hasModuleDeclarationMatchingSpecifier(sourceFile, moduleSpecifier) {
				if sourceFile.IsDeclarationFile {
					// TODO: Create import without clause
					usedImports = append(usedImports, importDecl)
				} else {
					usedImports = append(usedImports, importDecl)
				}
			}
		}
	}

	return usedImports
}

func isDeclarationUsed(identifier *ast.Identifier, typeChecker *checker.Checker, sourceFile *ast.SourceFile) bool {
	// TODO: Check for JSX factory usage
	// TODO: Implement isSymbolReferencedInFile from FindAllReferences
	_ = identifier
	_ = typeChecker
	_ = sourceFile
	return true
}

func filterUsedImportSpecifiers(elements []*ast.Statement, typeChecker *checker.Checker, sourceFile *ast.SourceFile) []*ast.Statement {
	var result []*ast.Statement
	for _, elem := range elements {
		spec := elem.AsImportSpecifier()
		if isDeclarationUsed(spec.Name().AsIdentifier(), typeChecker, sourceFile) {
			result = append(result, elem)
		}
	}
	return result
}

func hasModuleDeclarationMatchingSpecifier(sourceFile *ast.SourceFile, moduleSpecifier *ast.Expression) bool {
	if moduleSpecifier == nil || !ast.IsStringLiteral(moduleSpecifier.AsNode()) {
		return false
	}
	moduleSpecifierText := moduleSpecifier.Text()

	// TODO: Check source file module augmentations
	_ = moduleSpecifierText
	return false
}

func coalesceImportsWorker(
	importDecls []*ast.Statement,
	comparer func(a, b string) int,
	specifierComparer func(s1, s2 *ast.Node) int,
	sourceFile *ast.SourceFile,
) []*ast.Statement {
	if len(importDecls) == 0 {
		return importDecls
	}

	// TODO: Group by attributes
	categorized := getCategorizedImports(importDecls)

	coalescedImports := make([]*ast.Statement, 0)

	if categorized.importWithoutClause != nil {
		coalescedImports = append(coalescedImports, categorized.importWithoutClause)
	}

	for _, group := range []importGroup{categorized.regularImports, categorized.typeOnlyImports} {
		if group.isEmpty() {
			continue
		}

		slices.SortFunc(group.namespaceImports, func(a, b *ast.Statement) int {
			n1 := a.AsImportDeclaration().ImportClause.AsImportClause().NamedBindings.AsNamespaceImport().Name()
			n2 := b.AsImportDeclaration().ImportClause.AsImportClause().NamedBindings.AsNamespaceImport().Name()
			return comparer(n1.Text(), n2.Text())
		})

		for _, nsImport := range group.namespaceImports {
			coalescedImports = append(coalescedImports, nsImport)
		}

		// TODO: Handle default and named imports combination
		_ = group.defaultImports
		_ = group.namedImports
		_ = specifierComparer
	}

	return coalescedImports
}

type categorizedImports struct {
	importWithoutClause *ast.Statement
	typeOnlyImports     importGroup
	regularImports      importGroup
}

type importGroup struct {
	defaultImports   []*ast.Statement
	namespaceImports []*ast.Statement
	namedImports     []*ast.Statement
}

func (g importGroup) isEmpty() bool {
	return len(g.defaultImports) == 0 && len(g.namespaceImports) == 0 && len(g.namedImports) == 0
}

func getCategorizedImports(importDecls []*ast.Statement) categorizedImports {
	var importWithoutClause *ast.Statement
	var typeOnlyImports, regularImports importGroup

	for _, importDecl := range importDecls {
		if importDecl.ImportClause == nil {
			if importWithoutClause == nil {
				importWithoutClause = importDecl
			}
			continue
		}

		clause := importDecl.AsImportDeclaration().ImportClause.AsImportClause()
		group := &regularImports
		if clause.IsTypeOnly() {
			group = &typeOnlyImports
		}

		name := clause.Name()
		namedBindings := clause.NamedBindings

		if name != nil {
			group.defaultImports = append(group.defaultImports, importDecl)
		}

		if namedBindings != nil {
			if namedBindings.Kind == ast.KindNamespaceImport {
				group.namespaceImports = append(group.namespaceImports, importDecl)
			} else if namedBindings.Kind == ast.KindNamedImports {
				group.namedImports = append(group.namedImports, importDecl)
			}
		}
	}

	return categorizedImports{
		importWithoutClause: importWithoutClause,
		typeOnlyImports:     typeOnlyImports,
		regularImports:      regularImports,
	}
}

func groupByNewlineContiguous(sourceFile *ast.SourceFile, decls []*ast.Statement) [][]*ast.Statement {
	s := scanner.NewScanner()
	var groups [][]*ast.Statement
	var currentGroup []*ast.Statement

	for _, decl := range decls {
		if len(currentGroup) > 0 && isNewGroup(sourceFile, decl, s) {
			groups = append(groups, currentGroup)
			currentGroup = nil
		}
		currentGroup = append(currentGroup, decl)
	}

	if len(currentGroup) > 0 {
		groups = append(groups, currentGroup)
	}

	return groups
}

func isNewGroup(sourceFile *ast.SourceFile, decl *ast.Statement, s *scanner.Scanner) bool {
	startPos := scanner.SkipTrivia(sourceFile.Text(), decl.Pos())
	fullStart := decl.Pos()
	triviaLen := startPos - fullStart
	s.SetText(sourceFile.Text()[fullStart:startPos])

	numberOfNewLines := 0
	for s.TokenStart() < triviaLen {
		tokenKind := s.Scan()
		if tokenKind == ast.KindNewLineTrivia {
			numberOfNewLines++
			if numberOfNewLines >= 2 {
				return true
			}
		}
	}

	return false
}

func getDetectionLists(preferences *lsutil.UserPreferences) (comparersToTest []func(a, b string) int, typeOrdersToTest []lsutil.OrganizeImportsTypeOrder) {
	if preferences != nil && !preferences.OrganizeImportsIgnoreCase.IsUnknown() {
		ignoreCase := preferences.OrganizeImportsIgnoreCase.IsTrue()
		comparersToTest = []func(a, b string) int{getOrganizeImportsOrdinalStringComparer(ignoreCase)}
	} else {
		comparersToTest = []func(a, b string) int{
			getOrganizeImportsOrdinalStringComparer(true),
			getOrganizeImportsOrdinalStringComparer(false),
		}
	}

	if preferences != nil && preferences.OrganizeImportsTypeOrder != lsutil.OrganizeImportsTypeOrderAuto {
		typeOrdersToTest = []lsutil.OrganizeImportsTypeOrder{preferences.OrganizeImportsTypeOrder}
	} else {
		typeOrdersToTest = []lsutil.OrganizeImportsTypeOrder{
			lsutil.OrganizeImportsTypeOrderLast,
			lsutil.OrganizeImportsTypeOrderInline,
			lsutil.OrganizeImportsTypeOrderFirst,
		}
	}

	return comparersToTest, typeOrdersToTest
}

type namedImportSortResult struct {
	namedImportComparer func(a, b string) int
	typeOrder           lsutil.OrganizeImportsTypeOrder
	isSorted            bool
}

func detectNamedImportOrganizationBySort(
	originalGroups []*ast.Statement,
	comparersToTest []func(a, b string) int,
	typesToTest []lsutil.OrganizeImportsTypeOrder,
) *namedImportSortResult {
	var bothNamedImports bool
	var importDeclsWithNamed []*ast.Statement

	for _, imp := range originalGroups {
		if imp.ImportClause == nil {
			continue
		}
		clause := imp.AsImportDeclaration().ImportClause.AsImportClause()
		if clause.NamedBindings == nil || clause.NamedBindings.Kind != ast.KindNamedImports {
			continue
		}
		namedImports := clause.NamedBindings.AsNamedImports()
		if len(namedImports.Elements.Nodes) == 0 {
			continue
		}

		if !bothNamedImports {
			hasTypeOnly := false
			hasRegular := false
			for _, elem := range namedImports.Elements.Nodes {
				if elem.IsTypeOnly() {
					hasTypeOnly = true
				} else {
					hasRegular = true
				}
			}
			if hasTypeOnly && hasRegular {
				bothNamedImports = true
			}
		}

		importDeclsWithNamed = append(importDeclsWithNamed, imp)
	}

	if len(importDeclsWithNamed) == 0 {
		return nil
	}

	namedImportsByDecl := make([][]*ast.Statement, 0, len(importDeclsWithNamed))
	for _, imp := range importDeclsWithNamed {
		clause := imp.AsImportDeclaration().ImportClause.AsImportClause()
		namedImports := clause.NamedBindings.AsNamedImports()
		namedImportsByDecl = append(namedImportsByDecl, namedImports.Elements.Nodes)
	}

	if !bothNamedImports || len(typesToTest) == 0 {
		namesList := make([][]string, len(namedImportsByDecl))
		for i, imports := range namedImportsByDecl {
			names := make([]string, len(imports))
			for j, imp := range imports {
				names[j] = imp.Name().Text()
			}
			namesList[i] = names
		}
		sortState := detectCaseSensitivityBySort(namesList, comparersToTest)
		typeOrder := lsutil.OrganizeImportsTypeOrderLast
		if len(typesToTest) == 1 {
			typeOrder = typesToTest[0]
		}
		return &namedImportSortResult{
			namedImportComparer: sortState.comparer,
			typeOrder:           typeOrder,
			isSorted:            sortState.isSorted,
		}
	}

	bestDiff := map[lsutil.OrganizeImportsTypeOrder]int{
		lsutil.OrganizeImportsTypeOrderFirst:  math.MaxInt,
		lsutil.OrganizeImportsTypeOrderLast:   math.MaxInt,
		lsutil.OrganizeImportsTypeOrderInline: math.MaxInt,
	}
	bestComparer := map[lsutil.OrganizeImportsTypeOrder]func(a, b string) int{
		lsutil.OrganizeImportsTypeOrderFirst:  comparersToTest[0],
		lsutil.OrganizeImportsTypeOrderLast:   comparersToTest[0],
		lsutil.OrganizeImportsTypeOrderInline: comparersToTest[0],
	}

	for _, curComparer := range comparersToTest {
		currDiff := map[lsutil.OrganizeImportsTypeOrder]int{
			lsutil.OrganizeImportsTypeOrderFirst:  0,
			lsutil.OrganizeImportsTypeOrderLast:   0,
			lsutil.OrganizeImportsTypeOrderInline: 0,
		}

		for _, importDecl := range namedImportsByDecl {
			for _, typeOrder := range typesToTest {
				prefs := &lsutil.UserPreferences{OrganizeImportsTypeOrder: typeOrder}
				diff := measureSortedness(importDecl, func(n1, n2 *ast.Node) int {
					return compareImportOrExportSpecifiers(n1, n2, curComparer, prefs)
				})
				currDiff[typeOrder] = currDiff[typeOrder] + diff
			}
		}

		for _, typeOrder := range typesToTest {
			if currDiff[typeOrder] < bestDiff[typeOrder] {
				bestDiff[typeOrder] = currDiff[typeOrder]
				bestComparer[typeOrder] = curComparer
			}
		}
	}

	for _, bestTypeOrder := range typesToTest {
		isBest := true
		for _, testTypeOrder := range typesToTest {
			if bestDiff[testTypeOrder] < bestDiff[bestTypeOrder] {
				isBest = false
				break
			}
		}
		if isBest {
			return &namedImportSortResult{
				namedImportComparer: bestComparer[bestTypeOrder],
				typeOrder:           bestTypeOrder,
				isSorted:            bestDiff[bestTypeOrder] == 0,
			}
		}
	}

	return &namedImportSortResult{
		namedImportComparer: bestComparer[lsutil.OrganizeImportsTypeOrderLast],
		typeOrder:           lsutil.OrganizeImportsTypeOrderLast,
		isSorted:            bestDiff[lsutil.OrganizeImportsTypeOrderLast] == 0,
	}
}

func getOrganizeImportsOrdinalStringComparer(ignoreCase bool) func(a, b string) int {
	if ignoreCase {
		return stringutil.CompareStringsCaseInsensitiveEslintCompatible
	}
	return stringutil.CompareStringsCaseSensitive
}

// getModuleSpecifierExpression returns the module specifier expression from an import/require statement
func getModuleSpecifierExpression(declaration *ast.Statement) *ast.Expression {
	switch declaration.Kind {
	case ast.KindImportEqualsDeclaration:
		importEquals := declaration.AsImportEqualsDeclaration()
		if importEquals.ModuleReference.Kind == ast.KindExternalModuleReference {
			return importEquals.ModuleReference.Expression()
		}
		return nil
	case ast.KindImportDeclaration:
		return declaration.ModuleSpecifier()
	case ast.KindVariableStatement:
		// For require statements: const x = require('...')
		variableStatement := declaration.AsVariableStatement()
		declarations := variableStatement.DeclarationList.AsVariableDeclarationList().Declarations.Nodes
		if len(declarations) > 0 {
			decl := declarations[0]
			initializer := decl.Initializer()
			if initializer != nil && initializer.Kind == ast.KindCallExpression {
				callExpr := initializer.AsCallExpression()
				if len(callExpr.Arguments.Nodes) > 0 {
					return callExpr.Arguments.Nodes[0]
				}
			}
		}
		return nil
	default:
		return nil
	}
}

func getExternalModuleName(specifier *ast.Expression) string {
	if specifier != nil && ast.IsStringLiteralLike(specifier.AsNode()) {
		return specifier.Text()
	}
	return ""
}

// compareModuleSpecifiersWorker compares two module specifiers
func compareModuleSpecifiersWorker(m1 *ast.Expression, m2 *ast.Expression, comparer func(a, b string) int) int {
	name1 := getExternalModuleName(m1)
	name2 := getExternalModuleName(m2)
	if cmp := core.CompareBooleans(name1 == "", name2 == ""); cmp != 0 {
		return cmp
	}
	if cmp := core.CompareBooleans(tspath.IsExternalModuleNameRelative(name1), tspath.IsExternalModuleNameRelative(name2)); cmp != 0 {
		return cmp
	}
	return comparer(name1, name2)
}

// compareImportKind returns comparison order based on import kind
func compareImportKind(s1 *ast.Statement, s2 *ast.Statement) int {
	return cmp.Compare(getImportKindOrder(s1), getImportKindOrder(s2))
}

// getImportKindOrder returns the sort order for different import kinds:
// 1. Side-effect imports
// 2. Type-only imports
// 3. Namespace imports
// 4. Default imports
// 5. Named imports
// 6. ImportEqualsDeclarations
// 7. Require variable statements
func getImportKindOrder(s1 *ast.Statement) int {
	switch s1.Kind {
	case ast.KindImportDeclaration:
		importDecl := s1.AsImportDeclaration()
		if importDecl.ImportClause == nil {
			return 0 // Side-effect import
		}
		importClause := importDecl.ImportClause.AsImportClause()
		if importClause.IsTypeOnly() {
			return 1 // Type-only import
		}
		if importClause.NamedBindings != nil && importClause.NamedBindings.Kind == ast.KindNamespaceImport {
			return 2 // Namespace import
		}
		if importClause.Name() != nil {
			return 3 // Default import
		}
		return 4 // Named imports
	case ast.KindImportEqualsDeclaration:
		return 5
	case ast.KindVariableStatement:
		return 6 // Require statement
	default:
		return 7
	}
}

// CompareImportsOrRequireStatements compares two import or require statements for sorting
func CompareImportsOrRequireStatements(s1 *ast.Statement, s2 *ast.Statement, comparer func(a, b string) int) int {
	if cmp := compareModuleSpecifiersWorker(getModuleSpecifierExpression(s1), getModuleSpecifierExpression(s2), comparer); cmp != 0 {
		return cmp
	}
	return compareImportKind(s1, s2)
}

// compareImportOrExportSpecifiers compares two import or export specifiers
func compareImportOrExportSpecifiers(s1 *ast.Node, s2 *ast.Node, comparer func(a, b string) int, preferences *lsutil.UserPreferences) int {
	typeOrder := lsutil.OrganizeImportsTypeOrderLast
	if preferences != nil {
		typeOrder = preferences.OrganizeImportsTypeOrder
	}

	s1Name := s1.Name().Text()
	s2Name := s2.Name().Text()

	switch typeOrder {
	case lsutil.OrganizeImportsTypeOrderFirst:
		if cmp := core.CompareBooleans(s2.IsTypeOnly(), s1.IsTypeOnly()); cmp != 0 {
			return cmp
		}
		return comparer(s1Name, s2Name)
	case lsutil.OrganizeImportsTypeOrderInline:
		return comparer(s1Name, s2Name)
	default: // OrganizeImportsTypeOrderLast
		if cmp := core.CompareBooleans(s1.IsTypeOnly(), s2.IsTypeOnly()); cmp != 0 {
			return cmp
		}
		return comparer(s1Name, s2Name)
	}
}

// GetNamedImportSpecifierComparer returns a comparer function for import/export specifiers
func GetNamedImportSpecifierComparer(preferences *lsutil.UserPreferences, comparer func(a, b string) int) func(s1, s2 *ast.Node) int {
	if comparer == nil {
		ignoreCase := false
		if preferences != nil && !preferences.OrganizeImportsIgnoreCase.IsUnknown() {
			ignoreCase = preferences.OrganizeImportsIgnoreCase.IsTrue()
		}
		comparer = getOrganizeImportsOrdinalStringComparer(ignoreCase)
	}
	return func(s1, s2 *ast.Node) int {
		return compareImportOrExportSpecifiers(s1, s2, comparer, preferences)
	}
}

// GetImportSpecifierInsertionIndex finds the insertion index for a new import specifier
func GetImportSpecifierInsertionIndex(sortedImports []*ast.Node, newImport *ast.Node, comparer func(s1, s2 *ast.Node) int) int {
	return core.FirstResult(core.BinarySearchUniqueFunc(sortedImports, func(mid int, value *ast.Node) int {
		return comparer(value, newImport)
	}))
}

// GetImportDeclarationInsertIndex finds the insertion index for a new import declaration
func GetImportDeclarationInsertIndex(sortedImports []*ast.Statement, newImport *ast.Statement, comparer func(a, b *ast.Statement) int) int {
	return core.FirstResult(core.BinarySearchUniqueFunc(sortedImports, func(mid int, value *ast.Statement) int {
		return comparer(value, newImport)
	}))
}

// GetOrganizeImportsStringComparerWithDetection detects the string comparer to use based on existing imports
func GetOrganizeImportsStringComparerWithDetection(originalImportDecls []*ast.Statement, preferences *lsutil.UserPreferences) (comparer func(a, b string) int, isSorted bool) {
	result := detectModuleSpecifierCaseBySort([][]*ast.Statement{originalImportDecls}, getComparers(preferences))
	return result.comparer, result.isSorted
}

func getComparers(preferences *lsutil.UserPreferences) []func(a string, b string) int {
	if preferences != nil {
		switch preferences.OrganizeImportsIgnoreCase {
		case core.TSTrue:
			return caseInsensitiveOrganizeImportsComparer
		case core.TSFalse:
			return caseSensitiveOrganizeImportsComparer
		}
	}

	return organizeImportsComparers
}

type caseSensitivityDetectionResult struct {
	comparer func(a, b string) int
	isSorted bool
}

func detectModuleSpecifierCaseBySort(importDeclsByGroup [][]*ast.Statement, comparersToTest []func(a, b string) int) caseSensitivityDetectionResult {
	moduleSpecifiersByGroup := make([][]string, 0, len(importDeclsByGroup))
	for _, importGroup := range importDeclsByGroup {
		moduleNames := make([]string, 0, len(importGroup))
		for _, decl := range importGroup {
			if expr := getModuleSpecifierExpression(decl); expr != nil {
				moduleNames = append(moduleNames, getExternalModuleName(expr))
			} else {
				moduleNames = append(moduleNames, "")
			}
		}
		moduleSpecifiersByGroup = append(moduleSpecifiersByGroup, moduleNames)
	}
	return detectCaseSensitivityBySort(moduleSpecifiersByGroup, comparersToTest)
}

func detectCaseSensitivityBySort(originalGroups [][]string, comparersToTest []func(a, b string) int) caseSensitivityDetectionResult {
	var bestComparer func(a, b string) int
	bestDiff := math.MaxInt

	for _, curComparer := range comparersToTest {
		diffOfCurrentComparer := 0

		for _, listToSort := range originalGroups {
			if len(listToSort) <= 1 {
				continue
			}
			diff := measureSortedness(listToSort, curComparer)
			diffOfCurrentComparer += diff
		}

		if diffOfCurrentComparer < bestDiff {
			bestDiff = diffOfCurrentComparer
			bestComparer = curComparer
		}
	}

	if bestComparer == nil && len(comparersToTest) > 0 {
		bestComparer = comparersToTest[0]
	}

	return caseSensitivityDetectionResult{
		comparer: bestComparer,
		isSorted: bestDiff == 0,
	}
}

func measureSortedness[T any](arr []T, comparer func(a, b T) int) int {
	i := 0
	for j := range len(arr) - 1 {
		if comparer(arr[j], arr[j+1]) > 0 {
			i++
		}
	}
	return i
}

// GetNamedImportSpecifierComparerWithDetection detects the appropriate comparer for named imports
func GetNamedImportSpecifierComparerWithDetection(importDecl *ast.Node, sourceFile *ast.SourceFile, preferences *lsutil.UserPreferences) (specifierComparer func(s1, s2 *ast.Node) int, isSorted core.Tristate) {
	specifierComparer = GetNamedImportSpecifierComparer(preferences, getComparers(preferences)[0])
	// Try to detect from the current import declaration
	if (preferences == nil || preferences.OrganizeImportsIgnoreCase.IsUnknown() || preferences.OrganizeImportsTypeOrder == lsutil.OrganizeImportsTypeOrderLast) &&
		importDecl.Kind == ast.KindImportDeclaration {
		// For now, just return the default comparer
		// Full detection logic would require porting detectNamedImportOrganizationBySort
		return specifierComparer, core.TSUnknown
	}

	return specifierComparer, core.TSUnknown
}
