package organizeimports

import (
	"cmp"
	"context"
	"math"
	"slices"
	"strings"
	"unicode"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/locale"
	"github.com/microsoft/typescript-go/internal/ls/change"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/printer"
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
	kind lsproto.CodeActionKind,
) {
	shouldSort := kind == "source.organizeImports.sortAndCombine" || kind == lsproto.CodeActionKindSourceOrganizeImports
	shouldCombine := shouldSort
	shouldRemove := kind == "source.organizeImports.removeUnused" || kind == lsproto.CodeActionKindSourceOrganizeImports
	topLevelImportDecls := filterImportDeclarations(sourceFile.Statements.Nodes)
	topLevelImportGroupDecls := groupByNewlineContiguous(sourceFile, topLevelImportDecls)

	comparersToTest, typeOrdersToTest := getDetectionLists(preferences)
	defaultComparer := comparersToTest[0]

	moduleSpecifierComparer := defaultComparer
	namedImportComparer := defaultComparer
	typeOrder := lsutil.OrganizeImportsTypeOrderAuto
	if preferences != nil {
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

	if kind != "source.organizeImports.removeUnused" {
		topLevelExportGroupDecls := getTopLevelExportGroups(sourceFile)
		for _, exportGroupDecl := range topLevelExportGroupDecls {
			organizeExportsWorker(exportGroupDecl, comparer, sourceFile, changeTracker)
		}
	}

	for _, stmt := range sourceFile.Statements.Nodes {
		if !ast.IsAmbientModule(stmt.AsNode()) {
			continue
		}

		ambientModule := stmt.AsModuleDeclaration()
		if ambientModule.Body == nil {
			continue
		}

		moduleBody := ambientModule.Body.AsModuleBlock()

		ambientModuleImportDecls := filterImportDeclarations(moduleBody.Statements.Nodes)
		ambientModuleImportGroupDecls := groupByNewlineContiguous(sourceFile, ambientModuleImportDecls)

		for _, importGroupDecl := range ambientModuleImportGroupDecls {
			organizeImportsWorker(importGroupDecl, comparer, shouldSort, shouldCombine, shouldRemove, sourceFile, program, changeTracker, ctx)
		}

		if kind != "source.organizeImports.removeUnused" {
			var ambientModuleExportDecls []*ast.Statement
			for _, s := range moduleBody.Statements.Nodes {
				if s.Kind == ast.KindExportDeclaration {
					ambientModuleExportDecls = append(ambientModuleExportDecls, s)
				}
			}
			organizeExportsWorker(ambientModuleExportDecls, comparer, sourceFile, changeTracker)
		}
	}
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

	// Header comment preservation is handled via LeadingTriviaOptionExclude in the change tracker below

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
			coalesced := coalesceImportsWorker(importGroup, comparer.moduleSpecifierComparer, specifierComparer, sourceFile, changeTracker)
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

	if len(oldImportDecls) > 0 {
		if len(newImportDecls) == 0 {
			changeTracker.DeleteNodeRange(
				sourceFile,
				oldImportDecls[0].AsNode(),
				oldImportDecls[len(oldImportDecls)-1].AsNode(),
				change.LeadingTriviaOptionExclude, // Preserve header comment
				change.TrailingTriviaOptionInclude,
			)
		} else {
			for _, imp := range newImportDecls {
				changeTracker.SetEmitFlags(imp.AsNode(), printer.EFNoLeadingComments)
			}

			options := change.NodeOptions{
				LeadingTriviaOption:  change.LeadingTriviaOptionExclude, // Preserve header comment
				TrailingTriviaOption: change.TrailingTriviaOptionInclude,
				Suffix:               "\n",
			}

			newNodes := core.Map(newImportDecls, func(s *ast.Statement) *ast.Node { return s.AsNode() })
			changeTracker.ReplaceNodeWithNodes(sourceFile, oldImportDecls[0].AsNode(), newNodes, &options)

			if len(oldImportDecls) > 1 {
				for i := 1; i < len(oldImportDecls); i++ {
					changeTracker.Delete(sourceFile, oldImportDecls[i].AsNode())
				}
			}
		}
	}
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

	compilerOptions := program.Options()
	jsxNamespace := typeChecker.GetJsxNamespace(sourceFile.AsNode())
	jsxFragmentFactory := typeChecker.GetJsxFragmentFactory(sourceFile.AsNode())

	jsxElementsPresent := (sourceFile.AsNode().SubtreeFacts() & ast.SubtreeContainsJsx) != 0
	jsxModeNeedsExplicitImport := compilerOptions.Jsx == core.JsxEmitReact || compilerOptions.Jsx == core.JsxEmitReactNative

	factory := ast.NewNodeFactory(ast.NodeFactoryHooks{})
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

		if name != nil && !isDeclarationUsed(name.AsIdentifier(), jsxNamespace, jsxFragmentFactory, jsxElementsPresent, jsxModeNeedsExplicitImport, typeChecker, sourceFile) {
			name = nil
		}

		if namedBindings != nil {
			switch namedBindings.Kind {
			case ast.KindNamespaceImport:
				nsImport := namedBindings.AsNamespaceImport()
				if !isDeclarationUsed(nsImport.Name().AsIdentifier(), jsxNamespace, jsxFragmentFactory, jsxElementsPresent, jsxModeNeedsExplicitImport, typeChecker, sourceFile) {
					namedBindings = nil
				}
			case ast.KindNamedImports:
				namedImports := namedBindings.AsNamedImports()
				newElements := filterUsedImportSpecifiers(namedImports.Elements.Nodes, jsxNamespace, jsxFragmentFactory, jsxElementsPresent, jsxModeNeedsExplicitImport, typeChecker, sourceFile)
				if len(newElements) == 0 {
					namedBindings = nil
				} else if len(newElements) < len(namedImports.Elements.Nodes) {
					newList := factory.NewNodeList(newElements)
					namedBindings = factory.UpdateNamedImports(namedImports, newList).AsNode()
				}
			}
		}

		if name != nil || namedBindings != nil {
			importDeclNode := importDecl.AsImportDeclaration()
			newClause := factory.UpdateImportClause(clause, clause.PhaseModifier, name, namedBindings)
			newImportDecl := factory.UpdateImportDeclaration(
				importDeclNode,
				importDeclNode.Modifiers(),
				newClause.AsNode(),
				importDeclNode.ModuleSpecifier,
				importDeclNode.Attributes,
			)
			usedImports = append(usedImports, newImportDecl)
		} else {
			moduleSpecifier := importDecl.ModuleSpecifier()
			if hasModuleDeclarationMatchingSpecifier(sourceFile, moduleSpecifier) {
				if sourceFile.IsDeclarationFile {
					importDeclNode := importDecl.AsImportDeclaration()
					newImportDecl := factory.NewImportDeclaration(
						importDeclNode.Modifiers(),
						nil, // no import clause
						importDeclNode.ModuleSpecifier,
						nil, // no attributes
					)
					usedImports = append(usedImports, newImportDecl)
				} else {
					usedImports = append(usedImports, importDecl)
				}
			}
		}
	}

	return usedImports
}

func isDeclarationUsed(
	identifier *ast.Identifier,
	jsxNamespace string,
	jsxFragmentFactory string,
	jsxElementsPresent bool,
	jsxModeNeedsExplicitImport bool,
	typeChecker *checker.Checker,
	sourceFile *ast.SourceFile,
) bool {
	if jsxElementsPresent && jsxModeNeedsExplicitImport {
		identifierText := identifier.Text
		if identifierText == jsxNamespace {
			return true
		}
		if jsxFragmentFactory != "" && identifierText == jsxFragmentFactory {
			return true
		}
	}

	symbol := typeChecker.GetSymbolAtLocation(identifier.AsNode())
	if symbol == nil {
		return true
	}

	return isSymbolReferencedInFile(identifier, symbol, typeChecker, sourceFile)
}

func isSymbolReferencedInFile(
	definition *ast.Identifier,
	symbol *ast.Symbol,
	typeChecker *checker.Checker,
	sourceFile *ast.SourceFile,
) bool {
	identifierText := definition.Text
	for _, token := range getPossibleSymbolReferenceNodes(sourceFile, identifierText, sourceFile.AsNode()) {
		if !ast.IsIdentifier(token) {
			continue
		}
		id := token.AsIdentifier()
		if id == definition || id.Text != identifierText {
			continue
		}
		refSymbol := typeChecker.GetSymbolAtLocation(token)
		if refSymbol == symbol {
			return true
		}
		if token.Parent != nil && token.Parent.Kind == ast.KindShorthandPropertyAssignment {
			shorthandSymbol := typeChecker.GetShorthandAssignmentValueSymbol(token.Parent)
			if shorthandSymbol == symbol {
				return true
			}
		}
		if token.Parent != nil && ast.IsExportSpecifier(token.Parent) {
			localSymbol := getLocalSymbolForExportSpecifier(token, refSymbol, token.Parent.AsExportSpecifier(), typeChecker)
			if localSymbol == symbol {
				return true
			}
		}
	}
	return false
}

func getPossibleSymbolReferenceNodes(sourceFile *ast.SourceFile, symbolName string, container *ast.Node) []*ast.Node {
	return core.MapNonNil(getPossibleSymbolReferencePositions(sourceFile, symbolName, container), func(pos int) *ast.Node {
		if referenceLocation := astnav.GetTouchingPropertyName(sourceFile, pos); referenceLocation != sourceFile.AsNode() {
			return referenceLocation
		}
		return nil
	})
}

func getPossibleSymbolReferencePositions(sourceFile *ast.SourceFile, symbolName string, container *ast.Node) []int {
	positions := []int{}

	if symbolName == "" {
		return positions
	}

	text := sourceFile.Text()
	sourceLength := len(text)
	symbolNameLength := len(symbolName)

	if container == nil {
		container = sourceFile.AsNode()
	}

	searchStart := max(container.Pos(), 0)
	endPos := container.End()
	if endPos < 0 || endPos > sourceLength {
		endPos = sourceLength
	}

	position := strings.Index(text[searchStart:], symbolName)
	if position >= 0 {
		position += searchStart
	}

	for position >= 0 {
		if position > endPos {
			break
		}

		endPosition := position + symbolNameLength

		if (position == 0 || !scanner.IsIdentifierPart(rune(text[position-1]))) &&
			(endPosition >= sourceLength || !scanner.IsIdentifierPart(rune(text[endPosition]))) {
			positions = append(positions, position)
		}

		nextStart := position + symbolNameLength + 1
		if nextStart >= sourceLength {
			break
		}
		foundIndex := strings.Index(text[nextStart:], symbolName)
		if foundIndex == -1 {
			break
		}
		position = nextStart + foundIndex
	}

	return positions
}

func getLocalSymbolForExportSpecifier(
	node *ast.Node,
	refSymbol *ast.Symbol,
	exportSpecifier *ast.ExportSpecifier,
	typeChecker *checker.Checker,
) *ast.Symbol {
	if exportSpecifier.PropertyName != nil {
		return refSymbol
	}
	return typeChecker.GetExportSpecifierLocalTargetSymbol(exportSpecifier.AsNode())
}

func filterUsedImportSpecifiers(
	elements []*ast.Statement,
	jsxNamespace string,
	jsxFragmentFactory string,
	jsxElementsPresent bool,
	jsxModeNeedsExplicitImport bool,
	typeChecker *checker.Checker,
	sourceFile *ast.SourceFile,
) []*ast.Statement {
	var result []*ast.Statement
	for _, elem := range elements {
		spec := elem.AsImportSpecifier()
		if isDeclarationUsed(spec.Name().AsIdentifier(), jsxNamespace, jsxFragmentFactory, jsxElementsPresent, jsxModeNeedsExplicitImport, typeChecker, sourceFile) {
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

	for _, moduleName := range sourceFile.ModuleAugmentations {
		if ast.IsStringLiteral(moduleName) && moduleName.Text() == moduleSpecifierText {
			return true
		}
	}

	return false
}

func coalesceImportsWorker(
	importDecls []*ast.Statement,
	comparer func(a, b string) int,
	specifierComparer func(s1, s2 *ast.Node) int,
	sourceFile *ast.SourceFile,
	changeTracker *change.Tracker,
) []*ast.Statement {
	if len(importDecls) == 0 {
		return importDecls
	}

	importGroupsByAttributes := make(map[string][]*ast.Statement)
	var attributeKeys []string

	for _, importDecl := range importDecls {
		key := getImportAttributesKey(importDecl.AsImportDeclaration().Attributes)
		if _, exists := importGroupsByAttributes[key]; !exists {
			attributeKeys = append(attributeKeys, key)
		}
		importGroupsByAttributes[key] = append(importGroupsByAttributes[key], importDecl)
	}

	coalescedImports := make([]*ast.Statement, 0)

	for _, attributeKey := range attributeKeys {
		importGroupSameAttrs := importGroupsByAttributes[attributeKey]
		categorized := getCategorizedImports(importGroupSameAttrs)

		if categorized.importWithoutClause != nil {
			coalescedImports = append(coalescedImports, categorized.importWithoutClause)
		}

		factory := ast.NewNodeFactory(ast.NodeFactoryHooks{})

		for _, group := range []importGroup{categorized.regularImports, categorized.typeOnlyImports} {
			if group.isEmpty() {
				continue
			}

			isTypeOnly := &group == &categorized.typeOnlyImports

			if !isTypeOnly && len(group.defaultImports) == 1 && len(group.namespaceImports) == 1 && len(group.namedImports) == 0 {
				defaultImport := group.defaultImports[0]
				namespaceImport := group.namespaceImports[0]

				defaultClause := defaultImport.AsImportDeclaration().ImportClause.AsImportClause()
				namespaceBindings := namespaceImport.AsImportDeclaration().ImportClause.AsImportClause().NamedBindings

				newClause := factory.UpdateImportClause(defaultClause, defaultClause.PhaseModifier, defaultClause.Name(), namespaceBindings)
				defaultDeclNode := defaultImport.AsImportDeclaration()
				newImportDecl := factory.UpdateImportDeclaration(
					defaultDeclNode,
					defaultDeclNode.Modifiers(),
					newClause,
					defaultDeclNode.ModuleSpecifier,
					defaultDeclNode.Attributes,
				)
				coalescedImports = append(coalescedImports, newImportDecl)
				continue
			}

			slices.SortFunc(group.namespaceImports, func(a, b *ast.Statement) int {
				n1 := a.AsImportDeclaration().ImportClause.AsImportClause().NamedBindings.AsNamespaceImport().Name()
				n2 := b.AsImportDeclaration().ImportClause.AsImportClause().NamedBindings.AsNamespaceImport().Name()
				return comparer(n1.Text(), n2.Text())
			})

			for _, nsImport := range group.namespaceImports {
				nsImportDecl := nsImport.AsImportDeclaration()
				clause := nsImportDecl.ImportClause.AsImportClause()
				newClause := factory.UpdateImportClause(clause, clause.PhaseModifier, nil, clause.NamedBindings)
				newImportDecl := factory.UpdateImportDeclaration(
					nsImportDecl,
					nsImportDecl.Modifiers(),
					newClause,
					nsImportDecl.ModuleSpecifier,
					nsImportDecl.Attributes,
				)
				coalescedImports = append(coalescedImports, newImportDecl)
			}

			var firstDefaultImport *ast.Statement
			var firstNamedImport *ast.Statement

			if len(group.defaultImports) > 0 {
				firstDefaultImport = group.defaultImports[0]
			}
			if len(group.namedImports) > 0 {
				firstNamedImport = group.namedImports[0]
			}

			importDecl := firstDefaultImport
			if importDecl == nil {
				importDecl = firstNamedImport
			}
			if importDecl == nil {
				continue
			}

			var newDefaultImport *ast.IdentifierNode
			var newImportSpecifiers []*ast.Node

			if len(group.defaultImports) == 1 {
				newDefaultImport = group.defaultImports[0].AsImportDeclaration().ImportClause.AsImportClause().Name()
			} else {
				for _, defaultImport := range group.defaultImports {
					defaultClause := defaultImport.AsImportDeclaration().ImportClause.AsImportClause()
					defaultName := defaultClause.Name()
					propertyName := factory.NewIdentifier("default")
					importSpec := factory.NewImportSpecifier(false, propertyName, defaultName)
					newImportSpecifiers = append(newImportSpecifiers, importSpec)
				}
			}

			newImportSpecifiers = append(newImportSpecifiers, getNewImportSpecifiers(group.namedImports, factory)...)
			slices.SortStableFunc(newImportSpecifiers, specifierComparer)

			var newNamedImports *ast.NamedImportBindings
			if len(newImportSpecifiers) == 0 {
				if newDefaultImport != nil {
					newNamedImports = nil
				} else {
					newNamedImports = factory.NewNamedImports(factory.NewNodeList(nil))
				}
			} else {
				sortedList := factory.NewNodeList(newImportSpecifiers)
				if firstNamedImport != nil {
					firstNamedBindings := firstNamedImport.AsImportDeclaration().ImportClause.AsImportClause().NamedBindings.AsNamedImports()
					originalElements := firstNamedBindings.Elements
					if originalElements.HasTrailingComma() {
						sortedList.Loc = originalElements.Loc
					}
					newNamedImports = factory.UpdateNamedImports(firstNamedBindings, sortedList).AsNode()
				} else {
					newNamedImports = factory.NewNamedImports(sortedList)
				}
			}

			if sourceFile != nil && newNamedImports != nil && firstNamedImport != nil {
				firstNamedBindings := firstNamedImport.AsImportDeclaration().ImportClause.AsImportClause().NamedBindings
				if !ast.NodeIsSynthesized(firstNamedBindings.AsNode()) && !rangeIsOnSingleLine(firstNamedBindings.Loc, sourceFile) {
					changeTracker.SetEmitFlags(newNamedImports.AsNode(), printer.EFMultiLine)
				}
			}

			if isTypeOnly && newDefaultImport != nil && newNamedImports != nil {
				importDeclNode := importDecl.AsImportDeclaration()

				defaultClause := factory.NewImportClause(importDeclNode.ImportClause.AsImportClause().PhaseModifier, newDefaultImport, nil)
				defaultImportDecl := factory.UpdateImportDeclaration(
					importDeclNode,
					importDeclNode.Modifiers(),
					defaultClause,
					importDeclNode.ModuleSpecifier,
					importDeclNode.Attributes,
				)
				coalescedImports = append(coalescedImports, defaultImportDecl)

				namedDeclNode := firstNamedImport
				if namedDeclNode == nil {
					namedDeclNode = importDecl
				}
				namedImportDeclNode := namedDeclNode.AsImportDeclaration()
				namedClause := factory.NewImportClause(namedImportDeclNode.ImportClause.AsImportClause().PhaseModifier, nil, newNamedImports)
				namedImportDecl := factory.UpdateImportDeclaration(
					namedImportDeclNode,
					namedImportDeclNode.Modifiers(),
					namedClause,
					namedImportDeclNode.ModuleSpecifier,
					namedImportDeclNode.Attributes,
				)
				coalescedImports = append(coalescedImports, namedImportDecl)
			} else {
				importDeclNode := importDecl.AsImportDeclaration()
				clauseNode := importDeclNode.ImportClause.AsImportClause()
				newClause := factory.UpdateImportClause(clauseNode, clauseNode.PhaseModifier, newDefaultImport, newNamedImports)
				newImportDecl := factory.UpdateImportDeclaration(
					importDeclNode,
					importDeclNode.Modifiers(),
					newClause,
					importDeclNode.ModuleSpecifier,
					importDeclNode.Attributes,
				)
				coalescedImports = append(coalescedImports, newImportDecl)
			}
		}
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

func getImportAttributesKey(attributes *ast.ImportAttributesNode) string {
	if attributes == nil {
		return ""
	}

	importAttrs := attributes.AsImportAttributes()
	var key strings.Builder
	key.WriteString(importAttrs.Token.String())
	key.WriteString(" ")

	attrNodes := make([]*ast.Node, len(importAttrs.Attributes.Nodes))
	copy(attrNodes, importAttrs.Attributes.Nodes)
	slices.SortFunc(attrNodes, func(a, b *ast.Node) int {
		aName := a.AsImportAttribute().Name().Text()
		bName := b.AsImportAttribute().Name().Text()
		return stringutil.CompareStringsCaseSensitive(aName, bName)
	})

	for _, attrNode := range attrNodes {
		attr := attrNode.AsImportAttribute()
		key.WriteString(attr.Name().Text())
		key.WriteString(":")
		if ast.IsStringLiteralLike(attr.Value.AsNode()) {
			key.WriteString(`"`)
			key.WriteString(attr.Value.Text())
			key.WriteString(`"`)
		} else {
			key.WriteString(attr.Value.AsNode().Text())
		}
		key.WriteString(" ")
	}

	return key.String()
}

func getCategorizedImports(importDecls []*ast.Statement) categorizedImports {
	var importWithoutClause *ast.Statement
	var typeOnlyImports, regularImports importGroup

	for _, importDecl := range importDecls {
		if importDecl.AsImportDeclaration().ImportClause == nil {
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
			switch namedBindings.Kind {
			case ast.KindNamespaceImport:
				group.namespaceImports = append(group.namespaceImports, importDecl)
			case ast.KindNamedImports:
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

func rangeIsOnSingleLine(r core.TextRange, sourceFile *ast.SourceFile) bool {
	if r.Pos() < 0 || r.End() < 0 {
		return true
	}
	startLine, _ := scanner.GetECMALineAndCharacterOfPosition(sourceFile, r.Pos())
	endLine, _ := scanner.GetECMALineAndCharacterOfPosition(sourceFile, r.End())
	return startLine == endLine
}

func getNewImportSpecifiers(namedImports []*ast.Statement, factory *ast.NodeFactory) []*ast.Node {
	var result []*ast.Node

	for _, namedImport := range namedImports {
		elements := tryGetNamedBindingElements(namedImport)
		if elements == nil {
			continue
		}

		for _, elem := range elements {
			spec := elem.AsImportSpecifier()

			if spec.PropertyName != nil && spec.Name() != nil {
				propertyText := spec.PropertyName.Text()
				nameText := spec.Name().Text()

				if propertyText == nameText {
					normalized := factory.UpdateImportSpecifier(spec, spec.IsTypeOnly, nil, spec.Name())
					result = append(result, normalized)
					continue
				}
			}

			result = append(result, elem)
		}
	}

	return result
}

func tryGetNamedBindingElements(namedImport *ast.Statement) []*ast.Statement {
	if namedImport.Kind != ast.KindImportDeclaration {
		return nil
	}

	importDecl := namedImport.AsImportDeclaration()
	if importDecl.ImportClause == nil {
		return nil
	}

	clause := importDecl.ImportClause.AsImportClause()
	namedBindings := clause.NamedBindings

	if namedBindings != nil && namedBindings.Kind == ast.KindNamedImports {
		namedImportsNode := namedBindings.AsNamedImports()
		return namedImportsNode.Elements.Nodes
	}

	return nil
}

func groupByNewlineContiguous(sourceFile *ast.SourceFile, decls []*ast.Statement) [][]*ast.Statement {
	s := scanner.NewScanner()
	s.SetSkipTrivia(false) // Must not skip trivia to detect newlines
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
	fullStart := decl.Pos()
	if fullStart < 0 {
		return false
	}

	text := sourceFile.Text()
	textLen := len(text)

	if fullStart >= textLen {
		return false
	}

	startPos := scanner.SkipTrivia(text, fullStart)
	if startPos <= fullStart {
		return false
	}

	triviaLen := startPos - fullStart
	s.SetText(text[fullStart:startPos])

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
		comparersToTest = []func(a, b string) int{getOrganizeImportsStringComparer(preferences, ignoreCase)}
	} else {
		comparersToTest = []func(a, b string) int{
			getOrganizeImportsStringComparer(preferences, true),
			getOrganizeImportsStringComparer(preferences, false),
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
		if imp.AsImportDeclaration().ImportClause == nil {
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

func getOrganizeImportsUnicodeStringComparer(ignoreCase bool, preferences *lsutil.UserPreferences) func(a, b string) int {
	resolvedLocale := getOrganizeImportsLocale(preferences)

	caseFirst := lsutil.OrganizeImportsCaseFirstFalse
	numeric := false
	accents := true

	if preferences != nil {
		caseFirst = preferences.OrganizeImportsCaseFirst
		numeric = preferences.OrganizeImportsNumericCollation
		accents = preferences.OrganizeImportsAccentCollation
	}

	tag, _ := language.Parse(resolvedLocale)

	var opts []collate.Option

	if numeric {
		opts = append(opts, collate.Numeric)
	}

	looseOpts := append([]collate.Option{}, opts...)
	looseOpts = append(looseOpts, collate.Loose)
	looseCollator := collate.New(tag, looseOpts...)

	if !ignoreCase {
		caseInsensitiveOpts := append([]collate.Option{}, opts...)
		caseInsensitiveOpts = append(caseInsensitiveOpts, collate.IgnoreCase)
		caseInsensitiveCollator := collate.New(tag, caseInsensitiveOpts...)

		fullCollator := collate.New(tag, opts...)

		return func(a, b string) int {
			var primaryCmp int
			if !accents {
				primaryCmp = looseCollator.CompareString(a, b)
			} else {
				primaryCmp = caseInsensitiveCollator.CompareString(a, b)
			}
			if primaryCmp != 0 {
				return primaryCmp
			}

			aRunes := []rune(a)
			bRunes := []rune(b)
			minLen := min(len(aRunes), len(bRunes))

			for i := range minLen {
				aUpper := unicode.IsUpper(aRunes[i])
				bUpper := unicode.IsUpper(bRunes[i])
				if aUpper != bUpper {
					switch caseFirst {
					case lsutil.OrganizeImportsCaseFirstUpper:
						if aUpper {
							return -1
						}
						return 1
					case lsutil.OrganizeImportsCaseFirstLower:
						if !aUpper {
							return -1
						}
						return 1
					default:
						if aUpper {
							return 1
						}
						return -1
					}
				}
			}

			if !accents {
				if len(aRunes) != len(bRunes) {
					return len(aRunes) - len(bRunes)
				}
				return 0
			}

			return fullCollator.CompareString(a, b)
		}
	}

	if ignoreCase {
		opts = append(opts, collate.IgnoreCase)
		if !accents {
			opts = append(opts, collate.Loose)
		}
	}

	collator := collate.New(tag, opts...)

	return func(a, b string) int {
		return collator.CompareString(a, b)
	}
}

func getOrganizeImportsLocale(preferences *lsutil.UserPreferences) string {
	localeStr := "en"
	if preferences != nil && preferences.OrganizeImportsLocale != "" {
		localeStr = preferences.OrganizeImportsLocale
	}

	if localeStr == "auto" {
		if locale.Default != (locale.Locale{}) {
			tag := language.Tag(locale.Default)
			return tag.String()
		}
		return "en"
	}

	if locale, ok := locale.Parse(localeStr); ok {
		tag := language.Tag(locale)
		return tag.String()
	}

	return "en"
}

func getOrganizeImportsStringComparer(preferences *lsutil.UserPreferences, ignoreCase bool) func(a, b string) int {
	collation := lsutil.OrganizeImportsCollationOrdinal
	if preferences != nil {
		collation = preferences.OrganizeImportsCollation
	}

	if collation == lsutil.OrganizeImportsCollationUnicode {
		return getOrganizeImportsUnicodeStringComparer(ignoreCase, preferences)
	}
	return getOrganizeImportsOrdinalStringComparer(ignoreCase)
}

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

func CompareImportsOrRequireStatements(s1 *ast.Statement, s2 *ast.Statement, comparer func(a, b string) int) int {
	if cmp := compareModuleSpecifiersWorker(getModuleSpecifierExpression(s1), getModuleSpecifierExpression(s2), comparer); cmp != 0 {
		return cmp
	}
	return compareImportKind(s1, s2)
}

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

// compareExportSpecifiers compares two export specifiers considering type order.
func compareExportSpecifiers(s1 *ast.Node, s2 *ast.Node, comparer func(a, b string) int, typeOrder lsutil.OrganizeImportsTypeOrder) int {
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
	default: // OrganizeImportsTypeOrderLast or Auto (defaults to Last)
		if cmp := core.CompareBooleans(s1.IsTypeOnly(), s2.IsTypeOnly()); cmp != 0 {
			return cmp
		}
		return comparer(s1Name, s2Name)
	}
}

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

func GetImportSpecifierInsertionIndex(sortedImports []*ast.Node, newImport *ast.Node, comparer func(s1, s2 *ast.Node) int) int {
	return core.FirstResult(core.BinarySearchUniqueFunc(sortedImports, func(mid int, value *ast.Node) int {
		return comparer(value, newImport)
	}))
}

func GetImportDeclarationInsertIndex(sortedImports []*ast.Statement, newImport *ast.Statement, comparer func(a, b *ast.Statement) int) int {
	return core.FirstResult(core.BinarySearchUniqueFunc(sortedImports, func(mid int, value *ast.Statement) int {
		return comparer(value, newImport)
	}))
}

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

func getTopLevelExportGroups(sourceFile *ast.SourceFile) [][]*ast.Statement {
	var topLevelExportGroups [][]*ast.Statement
	statements := sourceFile.Statements.Nodes
	statementsLen := len(statements)

	i := 0
	groupIndex := 0
	for i < statementsLen {
		if statements[i].Kind == ast.KindExportDeclaration {
			if groupIndex >= len(topLevelExportGroups) {
				topLevelExportGroups = append(topLevelExportGroups, []*ast.Statement{})
			}
			exportDecl := statements[i].AsExportDeclaration()
			if exportDecl.ModuleSpecifier != nil {
				topLevelExportGroups[groupIndex] = append(topLevelExportGroups[groupIndex], statements[i])
				i++
			} else {
				for i < statementsLen && statements[i].Kind == ast.KindExportDeclaration {
					topLevelExportGroups[groupIndex] = append(topLevelExportGroups[groupIndex], statements[i])
					i++
				}
				groupIndex++
			}
		} else {
			i++
			if groupIndex < len(topLevelExportGroups) && len(topLevelExportGroups[groupIndex]) > 0 {
				groupIndex++
			}
		}
	}

	var result [][]*ast.Statement
	for _, exportGroup := range topLevelExportGroups {
		subGroups := groupByNewlineContiguous(sourceFile, exportGroup)
		result = append(result, subGroups...)
	}

	return result
}

func organizeExportsWorker(
	oldExportDecls []*ast.Statement,
	comparer comparerSettings,
	sourceFile *ast.SourceFile,
	changeTracker *change.Tracker,
) {
	if len(oldExportDecls) == 0 {
		return
	}

	specifierComparerFunc := func(s1, s2 *ast.Node) int {
		return compareExportSpecifiers(s1, s2, comparer.namedImportComparer, comparer.typeOrder)
	}

	newExportDecls := coalesceExportsWorker(oldExportDecls, specifierComparerFunc, comparer.moduleSpecifierComparer, sourceFile, changeTracker)

	if len(oldExportDecls) > 0 {
		if len(newExportDecls) == 0 {
			changeTracker.DeleteNodeRange(
				sourceFile,
				oldExportDecls[0].AsNode(),
				oldExportDecls[len(oldExportDecls)-1].AsNode(),
				change.LeadingTriviaOptionExclude,
				change.TrailingTriviaOptionInclude,
			)
		} else {
			options := change.NodeOptions{
				LeadingTriviaOption:  change.LeadingTriviaOptionExclude,
				TrailingTriviaOption: change.TrailingTriviaOptionInclude,
				Suffix:               "\n",
			}

			newNodes := core.Map(newExportDecls, func(s *ast.Statement) *ast.Node { return s.AsNode() })
			changeTracker.ReplaceNodeWithNodes(sourceFile, oldExportDecls[0].AsNode(), newNodes, &options)

			if len(oldExportDecls) > 1 {
				for i := 1; i < len(oldExportDecls); i++ {
					changeTracker.Delete(sourceFile, oldExportDecls[i].AsNode())
				}
			}
		}
	}
}

func coalesceExportsWorker(
	exportGroup []*ast.Statement,
	specifierComparer func(s1, s2 *ast.Node) int,
	moduleSpecifierComparer func(a, b string) int,
	sourceFile *ast.SourceFile,
	changeTracker *change.Tracker,
) []*ast.Statement {
	if len(exportGroup) == 0 {
		return exportGroup
	}

	exportsByModuleSpecifier := make(map[string][]*ast.Statement)
	var moduleSpecifierOrder []string

	for _, exportDecl := range exportGroup {
		export := exportDecl.AsExportDeclaration()
		var moduleSpecifier string
		if export.ModuleSpecifier != nil {
			moduleSpecifier = export.ModuleSpecifier.Text()
		}
		if _, exists := exportsByModuleSpecifier[moduleSpecifier]; !exists {
			moduleSpecifierOrder = append(moduleSpecifierOrder, moduleSpecifier)
		}
		exportsByModuleSpecifier[moduleSpecifier] = append(exportsByModuleSpecifier[moduleSpecifier], exportDecl)
	}

	slices.SortStableFunc(moduleSpecifierOrder, func(a, b string) int {
		if a == "" && b != "" {
			return 1
		}
		if a != "" && b == "" {
			return -1
		}
		return moduleSpecifierComparer(a, b)
	})

	var coalescedExports []*ast.Statement
	factory := ast.NewNodeFactory(ast.NodeFactoryHooks{})

	for _, moduleSpecifier := range moduleSpecifierOrder {
		group := exportsByModuleSpecifier[moduleSpecifier]

		categorized := getCategorizedExports(group)

		if categorized.exportWithoutClause != nil {
			coalescedExports = append(coalescedExports, categorized.exportWithoutClause)
		}

		for _, subGroup := range [][]*ast.Statement{categorized.namedExports, categorized.typeOnlyExports} {
			if len(subGroup) == 0 {
				continue
			}

			var newExportSpecifiers []*ast.Node
			for _, exportDecl := range subGroup {
				exportClause := exportDecl.AsExportDeclaration().ExportClause
				if exportClause != nil && exportClause.Kind == ast.KindNamedExports {
					namedExports := exportClause.AsNamedExports()
					newExportSpecifiers = append(newExportSpecifiers, namedExports.Elements.Nodes...)
				}
			}

			slices.SortStableFunc(newExportSpecifiers, specifierComparer)

			exportDecl := subGroup[0].AsExportDeclaration()

			var updatedExportClause *ast.NamedExportBindings
			if exportDecl.ExportClause != nil {
				if exportDecl.ExportClause.Kind == ast.KindNamedExports {
					namedExports := exportDecl.ExportClause.AsNamedExports()
					sortedList := factory.NewNodeList(newExportSpecifiers)
					updatedExportClause = factory.UpdateNamedExports(namedExports, sortedList)

					if sourceFile != nil && !ast.NodeIsSynthesized(namedExports.AsNode()) && !rangeIsOnSingleLine(namedExports.Loc, sourceFile) {
						changeTracker.SetEmitFlags(updatedExportClause.AsNode(), printer.EFMultiLine)
					}
				} else {
					updatedExportClause = exportDecl.ExportClause
				}
			}

			newExportDecl := factory.UpdateExportDeclaration(
				exportDecl,
				exportDecl.Modifiers(),
				exportDecl.IsTypeOnly,
				updatedExportClause,
				exportDecl.ModuleSpecifier,
				exportDecl.Attributes,
			)
			coalescedExports = append(coalescedExports, newExportDecl)
		}
	}

	return coalescedExports
}

type categorizedExports struct {
	exportWithoutClause *ast.Statement
	namedExports        []*ast.Statement
	typeOnlyExports     []*ast.Statement
}

func getCategorizedExports(exportGroup []*ast.Statement) categorizedExports {
	var exportWithoutClause *ast.Statement
	var namedExports, typeOnlyExports []*ast.Statement

	for _, exportDecl := range exportGroup {
		export := exportDecl.AsExportDeclaration()
		if export.ExportClause == nil {
			if exportWithoutClause == nil {
				exportWithoutClause = exportDecl
			}
		} else if export.IsTypeOnly {
			typeOnlyExports = append(typeOnlyExports, exportDecl)
		} else {
			namedExports = append(namedExports, exportDecl)
		}
	}

	return categorizedExports{
		exportWithoutClause: exportWithoutClause,
		namedExports:        namedExports,
		typeOnlyExports:     typeOnlyExports,
	}
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

func GetNamedImportSpecifierComparerWithDetection(importDecl *ast.Node, sourceFile *ast.SourceFile, preferences *lsutil.UserPreferences) (specifierComparer func(s1, s2 *ast.Node) int, isSorted core.Tristate) {
	comparersToTest, typeOrdersToTest := getDetectionLists(preferences)

	var importStmt *ast.Statement
	if importDecl.Kind == ast.KindImportDeclaration {
		importStmt = importDecl
	}

	specifierComparer = GetNamedImportSpecifierComparer(preferences, comparersToTest[0])
	isSorted = core.TSUnknown

	if (preferences == nil || preferences.OrganizeImportsIgnoreCase.IsUnknown() || preferences.OrganizeImportsTypeOrder == lsutil.OrganizeImportsTypeOrderAuto) && importStmt != nil {
		detectFromDecl := detectNamedImportOrganizationBySort([]*ast.Statement{importStmt}, comparersToTest, typeOrdersToTest)
		if detectFromDecl != nil {
			isSorted = core.BoolToTristate(detectFromDecl.isSorted)
			specifierComparer = GetNamedImportSpecifierComparer(
				&lsutil.UserPreferences{OrganizeImportsTypeOrder: detectFromDecl.typeOrder},
				detectFromDecl.namedImportComparer,
			)
		} else if sourceFile != nil {
			allImports := filterImportDeclarations(sourceFile.Statements.Nodes)
			detectFromFile := detectNamedImportOrganizationBySort(allImports, comparersToTest, typeOrdersToTest)
			if detectFromFile != nil {
				isSorted = core.BoolToTristate(detectFromFile.isSorted)
				specifierComparer = GetNamedImportSpecifierComparer(
					&lsutil.UserPreferences{OrganizeImportsTypeOrder: detectFromFile.typeOrder},
					detectFromFile.namedImportComparer,
				)
			}
		}
	}

	return specifierComparer, isSorted
}
