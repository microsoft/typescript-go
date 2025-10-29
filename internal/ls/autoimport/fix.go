package autoimport

import (
	"context"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
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

type Fix struct {
	Kind       FixKind    `json:"kind"`
	ImportKind ImportKind `json:"importKind"`

	// FixKindAddNew

	ModuleSpecifier string `json:"moduleSpecifier,omitempty"`

	// FixKindAddToExisting

	// ImportIndex is the index of the existing import in file.Imports()
	ImportIndex int `json:"importIndex"`
}

func (f *Fix) Edits(ctx context.Context, file *ast.SourceFile) []*lsproto.TextEdit {
	return nil
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
