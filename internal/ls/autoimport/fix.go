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
	"github.com/microsoft/typescript-go/internal/format"
	"github.com/microsoft/typescript-go/internal/ls/change"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/modulespecifiers"
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

type newImportBinding struct {
	kind         ImportKind
	propertyName string
	name         string
}

type Fix struct {
	Kind       FixKind    `json:"kind"`
	Name       string     `json:"name,omitempty"`
	ImportKind ImportKind `json:"importKind"`

	// FixKindAddNew

	ModuleSpecifier string `json:"moduleSpecifier,omitempty"`

	// FixKindAddToExisting

	// ImportIndex is the index of the existing import in file.Imports()
	ImportIndex int `json:"importIndex"`
}

func (f *Fix) Edits(ctx context.Context, file *ast.SourceFile, compilerOptions *core.CompilerOptions, formatOptions *format.FormatCodeSettings, converters *lsconv.Converters) []*lsproto.TextEdit {
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
		addToExistingImport(tracker, file, importClauseOrBindingPattern, defaultImport, namedImports)
		return tracker.GetChanges()[file.FileName()]
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
		namedBindings := importClause.NamedBindings
		if namedBindings == nil || namedBindings.Kind != ast.KindNamedImports {
			panic("expected named imports")
		}
		if defaultImport != nil {
			debug.Assert(importClause.Name() == nil, "Cannot add a default import to an import clause that already has one")
			ct.InsertNodeAt(file, core.TextPos(astnav.GetStartOfNode(importClause.AsNode(), file, false)), ct.NodeFactory.NewIdentifier(defaultImport.name), change.NodeOptions{Suffix: ", "})
		}

		if len(namedImports) > 0 {
			specifierComparer, isSorted := ls.getNamedImportSpecifierComparerWithDetection(importClause.Parent, file)
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

	moduleSpecifier := GetModuleSpecifier(fromFile, export, userPreferences, program, program.Options())
	if moduleSpecifier == "" {
		return nil
	}
	return []*Fix{}
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
