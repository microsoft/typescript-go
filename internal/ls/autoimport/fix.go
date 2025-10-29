package autoimport

import (
	"context"
	"slices"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/debug"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/format"
	"github.com/microsoft/typescript-go/internal/ls/change"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/ls/organizeimports"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/modulespecifiers"
	"github.com/microsoft/typescript-go/internal/stringutil"
)

type ImportKind int

const (
	ImportKindNamed     ImportKind = 0
	ImportKindDefault   ImportKind = 1
	ImportKindNamespace ImportKind = 2
	ImportKindCommonJS  ImportKind = 3
)

type FixKind int

const (
	// Sorted with the preferred fix coming first.
	FixKindUseNamespace    FixKind = 0
	FixKindJsdocTypeImport FixKind = 1
	FixKindAddToExisting   FixKind = 2
	FixKindAddNew          FixKind = 3
	FixKindPromoteTypeOnly FixKind = 4
)

type AddAsTypeOnly int

const (
	// These should not be combined as bitflags, but are given powers of 2 values to
	// easily detect conflicts between `NotAllowed` and `Required` by giving them a unique sum.
	// They're also ordered in terms of increasing priority for a fix-all scenario (see
	// `reduceAddAsTypeOnlyValues`).
	AddAsTypeOnlyAllowed    AddAsTypeOnly = 1 << 0
	AddAsTypeOnlyRequired   AddAsTypeOnly = 1 << 1
	AddAsTypeOnlyNotAllowed AddAsTypeOnly = 1 << 2
)

type newImportBinding struct {
	kind          ImportKind
	propertyName  string
	name          string
	addAsTypeOnly AddAsTypeOnly
}

type Fix struct {
	Kind          FixKind       `json:"kind"`
	Name          string        `json:"name,omitzero"`
	ImportKind    ImportKind    `json:"importKind"`
	UseRequire    bool          `json:"useRequire,omitzero"`
	AddAsTypeOnly AddAsTypeOnly `json:"addAsTypeOnly"`

	// FixKindAddNew

	ModuleSpecifier string `json:"moduleSpecifier,omitzero"`

	// FixKindAddToExisting

	// ImportIndex is the index of the existing import in file.Imports()
	ImportIndex int `json:"importIndex"`
}

func (f *Fix) Edits(
	ctx context.Context,
	file *ast.SourceFile,
	compilerOptions *core.CompilerOptions,
	formatOptions *format.FormatCodeSettings,
	converters *lsconv.Converters,
	preferences *lsutil.UserPreferences,
) ([]*lsproto.TextEdit, string) {
	tracker := change.NewTracker(ctx, compilerOptions, formatOptions, converters)
	switch f.Kind {
	case FixKindAddToExisting:
		if len(file.Imports()) <= f.ImportIndex {
			panic("import index out of range")
		}
		moduleSpecifier := file.Imports()[f.ImportIndex]
		importDecl := ast.TryGetImportFromModuleSpecifier(moduleSpecifier)
		if importDecl == nil {
			panic("expected import declaration")
		}
		var importClauseOrBindingPattern *ast.Node
		if importDecl.Kind == ast.KindImportDeclaration {
			importClauseOrBindingPattern = ast.GetImportClauseOfDeclaration(importDecl).AsNode()
			if importClauseOrBindingPattern == nil {
				panic("expected import clause")
			}
		} else if importDecl.Kind == ast.KindVariableDeclaration {
			importClauseOrBindingPattern = importDecl.Name().AsBindingPattern().AsNode()
		} else {
			panic("expected import declaration or variable declaration")
		}

		defaultImport := core.IfElse(f.ImportKind == ImportKindDefault, &newImportBinding{kind: ImportKindDefault, name: f.Name}, nil)
		namedImports := core.IfElse(f.ImportKind == ImportKindNamed, []*newImportBinding{{kind: ImportKindNamed, name: f.Name}}, nil)
		addToExistingImport(tracker, file, importClauseOrBindingPattern, defaultImport, namedImports, preferences)
		return tracker.GetChanges()[file.FileName()], diagnostics.Update_import_from_0.Format(f.ModuleSpecifier)
	case FixKindAddNew:
		var declarations []*ast.Statement
		defaultImport := core.IfElse(f.ImportKind == ImportKindDefault, &newImportBinding{name: f.Name}, nil)
		namedImports := core.IfElse(f.ImportKind == ImportKindNamed, []*newImportBinding{{name: f.Name}}, nil)
		var namespaceLikeImport *newImportBinding
		// qualification := f.qualification()
		// if f.ImportKind == ImportKindNamespace || f.ImportKind == ImportKindCommonJS {
		// 	namespaceLikeImport = &newImportBinding{kind: f.ImportKind, name: f.Name}
		// 	if qualification != nil && qualification.namespacePref != "" {
		// 		namespaceLikeImport.name = qualification.namespacePref
		// 	}
		// }

		if f.UseRequire {
			declarations = getNewRequires(tracker, f.ModuleSpecifier, defaultImport, namedImports, namespaceLikeImport, compilerOptions)
		} else {
			declarations = getNewImports(tracker, f.ModuleSpecifier, lsutil.GetQuotePreference(file, preferences), defaultImport, namedImports, namespaceLikeImport, compilerOptions, preferences)
		}

		insertImports(
			tracker,
			file,
			declarations,
			/*blankLineBetween*/ true,
			preferences,
		)
		// if qualification != nil {
		// 	addNamespaceQualifier(tracker, file, qualification)
		// }
		return tracker.GetChanges()[file.FileName()], diagnostics.Add_import_from_0.Format(f.ModuleSpecifier)
	default:
		panic("unimplemented fix edit")
	}
}

func addToExistingImport(
	ct *change.Tracker,
	file *ast.SourceFile,
	importClauseOrBindingPattern *ast.Node,
	defaultImport *newImportBinding,
	namedImports []*newImportBinding,
	preferences *lsutil.UserPreferences,
) {

	switch importClauseOrBindingPattern.Kind {
	case ast.KindObjectBindingPattern:
		bindingPattern := importClauseOrBindingPattern.AsBindingPattern()
		if defaultImport != nil {
			addElementToBindingPattern(ct, file, bindingPattern, defaultImport.name, "default")
		}
		for _, namedImport := range namedImports {
			addElementToBindingPattern(ct, file, bindingPattern, namedImport.name, "")
		}
		return
	case ast.KindImportClause:
		importClause := importClauseOrBindingPattern.AsImportClause()
		var existingSpecifiers []*ast.Node
		if importClause.NamedBindings != nil && importClause.NamedBindings.Kind == ast.KindNamedImports {
			existingSpecifiers = importClause.NamedBindings.Elements()
		}

		if defaultImport != nil {
			debug.Assert(importClause.Name() == nil, "Cannot add a default import to an import clause that already has one")
			ct.InsertNodeAt(file, core.TextPos(astnav.GetStartOfNode(importClause.AsNode(), file, false)), ct.NodeFactory.NewIdentifier(defaultImport.name), change.NodeOptions{Suffix: ", "})
		}

		if len(namedImports) > 0 {
			specifierComparer, isSorted := organizeimports.GetNamedImportSpecifierComparerWithDetection(importClause.Parent, file, preferences)
			newSpecifiers := core.Map(namedImports, func(namedImport *newImportBinding) *ast.Node {
				var identifier *ast.Node
				if namedImport.propertyName != "" {
					identifier = ct.NodeFactory.NewIdentifier(namedImport.propertyName).AsIdentifier().AsNode()
				}
				return ct.NodeFactory.NewImportSpecifier(
					false,
					identifier,
					ct.NodeFactory.NewIdentifier(namedImport.name),
				)
			})
			slices.SortFunc(newSpecifiers, specifierComparer)
			if len(existingSpecifiers) > 0 && isSorted != core.TSFalse {
				for _, spec := range newSpecifiers {
					insertionIndex := organizeimports.GetImportSpecifierInsertionIndex(existingSpecifiers, spec, specifierComparer)
					ct.InsertImportSpecifierAtIndex(file, spec, importClause.NamedBindings, insertionIndex)
				}
			} else if len(existingSpecifiers) > 0 && isSorted.IsTrue() {
				// Existing specifiers are sorted, so insert each new specifier at the correct position
				for _, spec := range newSpecifiers {
					insertionIndex := organizeimports.GetImportSpecifierInsertionIndex(existingSpecifiers, spec, specifierComparer)
					if insertionIndex >= len(existingSpecifiers) {
						// Insert at the end
						ct.InsertNodeInListAfter(file, existingSpecifiers[len(existingSpecifiers)-1], spec.AsNode(), existingSpecifiers)
					} else {
						// Insert before the element at insertionIndex
						ct.InsertNodeInListAfter(file, existingSpecifiers[insertionIndex], spec.AsNode(), existingSpecifiers)
					}
				}
			} else if len(existingSpecifiers) > 0 {
				// Existing specifiers may not be sorted, append to the end
				for _, spec := range newSpecifiers {
					ct.InsertNodeInListAfter(file, existingSpecifiers[len(existingSpecifiers)-1], spec.AsNode(), existingSpecifiers)
				}
			} else {
				if len(newSpecifiers) > 0 {
					namedImports := ct.NodeFactory.NewNamedImports(ct.NodeFactory.NewNodeList(newSpecifiers))
					if importClause.NamedBindings != nil {
						ct.ReplaceNode(file, importClause.NamedBindings, namedImports, nil)
					} else {
						if importClause.Name() == nil {
							panic("Import clause must have either named imports or a default import")
						}
						ct.InsertNodeAfter(file, importClause.Name(), namedImports)
					}
				}
			}
		}
	}
}

func addElementToBindingPattern(
	ct *change.Tracker,
	file *ast.SourceFile,
	bindingPattern *ast.BindingPattern,
	name string,
	propertyName string,
) {
	element := ct.NodeFactory.NewBindingElement(nil, nil, ct.NodeFactory.NewIdentifier(name), core.IfElse(propertyName == "", nil, ct.NodeFactory.NewIdentifier(propertyName)))
	if len(bindingPattern.Elements.Nodes) > 0 {
		ct.InsertNodeInListAfter(file, bindingPattern.Elements.Nodes[len(bindingPattern.Elements.Nodes)-1], element, bindingPattern.Elements.Nodes)
	} else {
		ct.ReplaceNode(file, bindingPattern.AsNode(), ct.NodeFactory.NewBindingPattern(ast.KindObjectBindingPattern, ct.AsNodeFactory().NewNodeList([]*ast.Node{element})), nil)
	}
}

func getNewImports(
	ct *change.Tracker,
	moduleSpecifier string,
	quotePreference lsutil.QuotePreference,
	defaultImport *newImportBinding,
	namedImports []*newImportBinding,
	namespaceLikeImport *newImportBinding, // { importKind: ImportKind.CommonJS | ImportKind.Namespace; }
	compilerOptions *core.CompilerOptions,
	preferences *lsutil.UserPreferences,
) []*ast.Statement {
	moduleSpecifierStringLiteral := ct.NodeFactory.NewStringLiteral(moduleSpecifier)
	if quotePreference == lsutil.QuotePreferenceSingle {
		moduleSpecifierStringLiteral.AsStringLiteral().TokenFlags |= ast.TokenFlagsSingleQuote
	}
	var statements []*ast.Statement // []AnyImportSyntax
	if defaultImport != nil || len(namedImports) > 0 {
		// `verbatimModuleSyntax` should prefer top-level `import type` -
		// even though it's not an error, it would add unnecessary runtime emit.
		topLevelTypeOnly := (defaultImport == nil || needsTypeOnly(defaultImport.addAsTypeOnly)) &&
			core.Every(namedImports, func(i *newImportBinding) bool { return needsTypeOnly(i.addAsTypeOnly) }) ||
			(compilerOptions.VerbatimModuleSyntax.IsTrue() || preferences.PreferTypeOnlyAutoImports) &&
				defaultImport != nil && defaultImport.addAsTypeOnly != AddAsTypeOnlyNotAllowed && !core.Some(namedImports, func(i *newImportBinding) bool { return i.addAsTypeOnly == AddAsTypeOnlyNotAllowed })

		var defaultImportNode *ast.Node
		if defaultImport != nil {
			defaultImportNode = ct.NodeFactory.NewIdentifier(defaultImport.name)
		}

		statements = append(statements, makeImport(ct, defaultImportNode, core.Map(namedImports, func(namedImport *newImportBinding) *ast.Node {
			var namedImportPropertyName *ast.Node
			if namedImport.propertyName != "" {
				namedImportPropertyName = ct.NodeFactory.NewIdentifier(namedImport.propertyName)
			}
			return ct.NodeFactory.NewImportSpecifier(
				!topLevelTypeOnly && shouldUseTypeOnly(namedImport.addAsTypeOnly, preferences),
				namedImportPropertyName,
				ct.NodeFactory.NewIdentifier(namedImport.name),
			)
		}), moduleSpecifierStringLiteral, topLevelTypeOnly))
	}

	if namespaceLikeImport != nil {
		var declaration *ast.Statement
		if namespaceLikeImport.kind == ImportKindCommonJS {
			declaration = ct.NodeFactory.NewImportEqualsDeclaration(
				/*modifiers*/ nil,
				shouldUseTypeOnly(namespaceLikeImport.addAsTypeOnly, preferences),
				ct.NodeFactory.NewIdentifier(namespaceLikeImport.name),
				ct.NodeFactory.NewExternalModuleReference(moduleSpecifierStringLiteral),
			)
		} else {
			declaration = ct.NodeFactory.NewImportDeclaration(
				/*modifiers*/ nil,
				ct.NodeFactory.NewImportClause(
					/*phaseModifier*/ core.IfElse(shouldUseTypeOnly(namespaceLikeImport.addAsTypeOnly, preferences), ast.KindTypeKeyword, ast.KindUnknown),
					/*name*/ nil,
					ct.NodeFactory.NewNamespaceImport(ct.NodeFactory.NewIdentifier(namespaceLikeImport.name)),
				),
				moduleSpecifierStringLiteral,
				/*attributes*/ nil,
			)
		}
		statements = append(statements, declaration)
	}
	if len(statements) == 0 {
		panic("No statements to insert for new imports")
	}
	return statements
}

func getNewRequires(
	changeTracker *change.Tracker,
	moduleSpecifier string,
	defaultImport *newImportBinding,
	namedImports []*newImportBinding,
	namespaceLikeImport *newImportBinding,
	compilerOptions *core.CompilerOptions,
) []*ast.Statement {
	quotedModuleSpecifier := changeTracker.NodeFactory.NewStringLiteral(moduleSpecifier)
	var statements []*ast.Statement

	// const { default: foo, bar, etc } = require('./mod');
	if defaultImport != nil || len(namedImports) > 0 {
		bindingElements := []*ast.Node{}
		for _, namedImport := range namedImports {
			var propertyName *ast.Node
			if namedImport.propertyName != "" {
				propertyName = changeTracker.NodeFactory.NewIdentifier(namedImport.propertyName)
			}
			bindingElements = append(bindingElements, changeTracker.NodeFactory.NewBindingElement(
				/*dotDotDotToken*/ nil,
				propertyName,
				changeTracker.NodeFactory.NewIdentifier(namedImport.name),
				/*initializer*/ nil,
			))
		}
		if defaultImport != nil {
			bindingElements = append([]*ast.Node{
				changeTracker.NodeFactory.NewBindingElement(
					/*dotDotDotToken*/ nil,
					changeTracker.NodeFactory.NewIdentifier("default"),
					changeTracker.NodeFactory.NewIdentifier(defaultImport.name),
					/*initializer*/ nil,
				),
			}, bindingElements...)
		}
		declaration := createConstEqualsRequireDeclaration(
			changeTracker,
			changeTracker.NodeFactory.NewBindingPattern(
				ast.KindObjectBindingPattern,
				changeTracker.NodeFactory.NewNodeList(bindingElements),
			),
			quotedModuleSpecifier,
		)
		statements = append(statements, declaration)
	}

	// const foo = require('./mod');
	if namespaceLikeImport != nil {
		declaration := createConstEqualsRequireDeclaration(
			changeTracker,
			changeTracker.NodeFactory.NewIdentifier(namespaceLikeImport.name),
			quotedModuleSpecifier,
		)
		statements = append(statements, declaration)
	}

	debug.AssertIsDefined(statements)
	return statements
}

func createConstEqualsRequireDeclaration(changeTracker *change.Tracker, name *ast.Node, quotedModuleSpecifier *ast.Node) *ast.Statement {
	return changeTracker.NodeFactory.NewVariableStatement(
		/*modifiers*/ nil,
		changeTracker.NodeFactory.NewVariableDeclarationList(
			ast.NodeFlagsConst,
			changeTracker.NodeFactory.NewNodeList([]*ast.Node{
				changeTracker.NodeFactory.NewVariableDeclaration(
					name,
					/*exclamationToken*/ nil,
					/*type*/ nil,
					changeTracker.NodeFactory.NewCallExpression(
						changeTracker.NodeFactory.NewIdentifier("require"),
						/*questionDotToken*/ nil,
						/*typeArguments*/ nil,
						changeTracker.NodeFactory.NewNodeList([]*ast.Node{quotedModuleSpecifier}),
						ast.NodeFlagsNone,
					),
				),
			}),
		),
	)
}

func insertImports(ct *change.Tracker, sourceFile *ast.SourceFile, imports []*ast.Statement, blankLineBetween bool, preferences *lsutil.UserPreferences) {
	var existingImportStatements []*ast.Statement

	if imports[0].Kind == ast.KindVariableStatement {
		existingImportStatements = core.Filter(sourceFile.Statements.Nodes, ast.IsRequireVariableStatement)
	} else {
		existingImportStatements = core.Filter(sourceFile.Statements.Nodes, ast.IsAnyImportSyntax)
	}
	comparer, isSorted := organizeimports.GetOrganizeImportsStringComparerWithDetection(existingImportStatements, preferences)
	sortedNewImports := slices.Clone(imports)
	slices.SortFunc(sortedNewImports, func(a, b *ast.Statement) int {
		return organizeimports.CompareImportsOrRequireStatements(a, b, comparer)
	})
	// !!! FutureSourceFile
	// if !isFullSourceFile(sourceFile) {
	//     for _, newImport := range sortedNewImports {
	//         // Insert one at a time to send correct original source file for accurate text reuse
	//         // when some imports are cloned from existing ones in other files.
	//         ct.insertStatementsInNewFile(sourceFile.fileName, []*ast.Node{newImport}, ast.GetSourceFileOfNode(getOriginalNode(newImport)))
	//     }
	// return;
	// }

	if len(existingImportStatements) > 0 && isSorted {
		// Existing imports are sorted, insert each new import at the correct position
		for _, newImport := range sortedNewImports {
			insertionIndex := organizeimports.GetImportDeclarationInsertIndex(existingImportStatements, newImport, func(a, b *ast.Statement) stringutil.Comparison {
				return organizeimports.CompareImportsOrRequireStatements(a, b, comparer)
			})
			if insertionIndex == 0 {
				// If the first import is top-of-file, insert after the leading comment which is likely the header
				ct.InsertNodeAt(sourceFile, core.TextPos(astnav.GetStartOfNode(existingImportStatements[0], sourceFile, false)), newImport.AsNode(), change.NodeOptions{})
			} else {
				prevImport := existingImportStatements[insertionIndex-1]
				ct.InsertNodeAfter(sourceFile, prevImport.AsNode(), newImport.AsNode())
			}
		}
	} else if len(existingImportStatements) > 0 {
		ct.InsertNodesAfter(sourceFile, existingImportStatements[len(existingImportStatements)-1], sortedNewImports)
	} else {
		ct.InsertAtTopOfFile(sourceFile, sortedNewImports, blankLineBetween)
	}
}

func makeImport(ct *change.Tracker, defaultImport *ast.IdentifierNode, namedImports []*ast.Node, moduleSpecifier *ast.Expression, isTypeOnly bool) *ast.Statement {
	var newNamedImports *ast.Node
	if len(namedImports) > 0 {
		newNamedImports = ct.NodeFactory.NewNamedImports(ct.NodeFactory.NewNodeList(namedImports))
	}
	var importClause *ast.Node
	if defaultImport != nil || newNamedImports != nil {
		importClause = ct.NodeFactory.NewImportClause(core.IfElse(isTypeOnly, ast.KindTypeKeyword, ast.KindUnknown), defaultImport, newNamedImports)
	}
	return ct.NodeFactory.NewImportDeclaration( /*modifiers*/ nil, importClause, moduleSpecifier, nil /*attributes*/)
}

func GetFixes(
	ctx context.Context,
	export *RawExport,
	fromFile *ast.SourceFile,
	program *compiler.Program,
	userPreferences modulespecifiers.UserPreferences,
) []*Fix {
	ch, done := program.GetTypeChecker(ctx)
	defer done()

	existingImports := getExistingImports(fromFile, ch)
	// !!! tryUseExistingNamespaceImport
	if fix := tryAddToExistingImport(export, fromFile, existingImports, program); fix != nil {
		return []*Fix{fix}
	}

	// !!! getNewImportFromExistingSpecifier - even worth it?

	moduleSpecifier := GetModuleSpecifier(fromFile, export, userPreferences, program, program.Options())
	if moduleSpecifier == "" || modulespecifiers.ContainsNodeModules(moduleSpecifier) {
		return nil
	}
	importKind := getImportKind(fromFile, export, program)
	// !!! JSDoc type import, add as type only
	return []*Fix{
		{
			Kind:            FixKindAddNew,
			ImportKind:      importKind,
			ModuleSpecifier: moduleSpecifier,
			Name:            export.Name,
		},
	}
}

func tryAddToExistingImport(
	export *RawExport,
	fromFile *ast.SourceFile,
	existingImports collections.MultiMap[ModuleID, existingImport],
	program *compiler.Program,
) *Fix {
	matchingDeclarations := existingImports.Get(export.ModuleID)
	if len(matchingDeclarations) == 0 {
		return nil
	}

	// Can't use an es6 import for a type in JS.
	if ast.IsSourceFileJS(fromFile) && export.Flags&ast.SymbolFlagsValue == 0 && !core.Every(matchingDeclarations, func(i existingImport) bool {
		return ast.IsJSDocImportTag(i.node)
	}) {
		return nil
	}

	importKind := getImportKind(fromFile, export, program)
	if importKind == ImportKindCommonJS || importKind == ImportKindNamespace {
		return nil
	}

	for _, existingImport := range matchingDeclarations {
		if existingImport.node.Kind == ast.KindImportEqualsDeclaration {
			continue
		}

		if existingImport.node.Kind == ast.KindVariableDeclaration {
			if (importKind == ImportKindNamed || importKind == ImportKindDefault) && existingImport.node.Name().Kind == ast.KindObjectBindingPattern {
				return &Fix{
					Kind:            FixKindAddToExisting,
					Name:            export.Name,
					ImportKind:      importKind,
					ImportIndex:     existingImport.index,
					ModuleSpecifier: existingImport.moduleSpecifier,
				}
			}
			continue
		}

		importClause := ast.GetImportClauseOfDeclaration(existingImport.node)
		if importClause == nil || !ast.IsStringLiteralLike(existingImport.node.ModuleSpecifier()) {
			continue
		}

		namedBindings := importClause.NamedBindings
		// A type-only import may not have both a default and named imports, so the only way a name can
		// be added to an existing type-only import is adding a named import to existing named bindings.
		if importClause.IsTypeOnly() && !(importKind == ImportKindNamed && namedBindings != nil) {
			continue
		}

		// Cannot add a named import to a declaration that has a namespace import
		if importKind == ImportKindNamed && namedBindings != nil && namedBindings.Kind == ast.KindNamespaceImport {
			continue
		}

		return &Fix{
			Kind:            FixKindAddToExisting,
			Name:            export.Name,
			ImportKind:      importKind,
			ImportIndex:     existingImport.index,
			ModuleSpecifier: existingImport.moduleSpecifier,
		}
	}

	return nil
}

func getImportKind(importingFile *ast.SourceFile, export *RawExport, program *compiler.Program) ImportKind {
	if program.Options().VerbatimModuleSyntax.IsTrue() && program.GetEmitModuleFormatOfFile(importingFile) == core.ModuleKindCommonJS {
		return ImportKindCommonJS
	}
	switch export.Syntax {
	case ExportSyntaxDefaultModifier, ExportSyntaxDefaultDeclaration:
		return ImportKindDefault
	case ExportSyntaxNamed, ExportSyntaxModifier:
		return ImportKindNamed
	case ExportSyntaxEquals:
		return ImportKindDefault
	default:
		panic("unhandled export syntax kind: " + export.Syntax.String())
	}
}

type existingImport struct {
	node            *ast.Node
	moduleSpecifier string
	index           int
}

func getExistingImports(file *ast.SourceFile, ch *checker.Checker) collections.MultiMap[ModuleID, existingImport] {
	result := collections.MultiMap[ModuleID, existingImport]{}
	for i, moduleSpecifier := range file.Imports() {
		node := ast.TryGetImportFromModuleSpecifier(moduleSpecifier)
		if node == nil {
			panic("error: did not expect node kind " + moduleSpecifier.Kind.String())
		} else if ast.IsVariableDeclarationInitializedToRequire(node.Parent) {
			if moduleSymbol := ch.ResolveExternalModuleName(moduleSpecifier); moduleSymbol != nil {
				result.Add(getModuleIDOfModuleSymbol(moduleSymbol), existingImport{node: node.Parent, moduleSpecifier: moduleSpecifier.Text(), index: i})
			}
		} else if node.Kind == ast.KindImportDeclaration || node.Kind == ast.KindImportEqualsDeclaration || node.Kind == ast.KindJSDocImportTag {
			if moduleSymbol := ch.GetSymbolAtLocation(moduleSpecifier); moduleSymbol != nil {
				result.Add(getModuleIDOfModuleSymbol(moduleSymbol), existingImport{node: node, moduleSpecifier: moduleSpecifier.Text(), index: i})
			}
		}
	}
	return result
}

func needsTypeOnly(addAsTypeOnly AddAsTypeOnly) bool {
	return addAsTypeOnly == AddAsTypeOnlyRequired
}

func shouldUseTypeOnly(addAsTypeOnly AddAsTypeOnly, preferences *lsutil.UserPreferences) bool {
	return needsTypeOnly(addAsTypeOnly) || addAsTypeOnly != AddAsTypeOnlyNotAllowed && preferences.PreferTypeOnlyAutoImports
}
