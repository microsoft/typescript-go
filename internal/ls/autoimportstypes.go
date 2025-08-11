package ls

import (
	"fmt"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/checker"
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

type ExportKind int

const (
	ExportKindNamed        ExportKind = 0
	ExportKindDefault      ExportKind = 1
	ExportKindExportEquals ExportKind = 2
	ExportKindUMD          ExportKind = 3
	ExportKindModule       ExportKind = 4
)

func (k ExportKind) String() string {
	switch k {
	case ExportKindNamed:
		return "Named"
	case ExportKindDefault:
		return "Default"
	case ExportKindExportEquals:
		return "ExportEquals"
	case ExportKindUMD:
		return "UMD"
	case ExportKindModule:
		return "Module"
	}
	panic(fmt.Sprintf("unexpected export kind: %d", k))
}

type ImportFix interface {
	Kind() ImportFixKind
	Base() *ImportFixBase
}

type ImportFixKind int

const (
	// Sorted with the preferred fix coming first.
	ImportFixKindUseNamespace    ImportFixKind = 0
	ImportFixKindJsdocTypeImport ImportFixKind = 1
	ImportFixKindAddToExisting   ImportFixKind = 2
	ImportFixKindAddNew          ImportFixKind = 3
	ImportFixKindPromoteTypeOnly ImportFixKind = 4
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

type ImportFixBase struct {
	isReExport          *bool
	exportInfo          *SymbolExportInfo // !!! | FutureSymbolExportInfo | undefined
	moduleSpecifierKind modulespecifiers.ResultKind
	moduleSpecifier     string
}

type Qualification struct {
	usagePosition   lsproto.Position
	namespacePrefix string
}

type FixUseNamespaceImport struct {
	ImportFixBase
	Qualification
}

func (f *FixUseNamespaceImport) Kind() ImportFixKind {
	return ImportFixKindUseNamespace
}

func (f *FixUseNamespaceImport) Base() *ImportFixBase {
	return &f.ImportFixBase
}

func getUseNamespaceImport(
	moduleSpecifier string,
	moduleSpecifierKind modulespecifiers.ResultKind,
	namespacePrefix string,
	usagePosition lsproto.Position,
) *FixUseNamespaceImport {
	return &FixUseNamespaceImport{
		ImportFixBase: ImportFixBase{
			moduleSpecifierKind: moduleSpecifierKind,
			moduleSpecifier:     moduleSpecifier,
		},
		Qualification: Qualification{
			usagePosition:   usagePosition,
			namespacePrefix: namespacePrefix,
		},
	}
}

type FixAddJsdocTypeImport struct {
	ImportFixBase
	usagePosition *lsproto.Position
}

func (f *FixAddJsdocTypeImport) Kind() ImportFixKind {
	return ImportFixKindJsdocTypeImport
}

func (f *FixAddJsdocTypeImport) Base() *ImportFixBase {
	return &f.ImportFixBase
}

func getAddJsdocTypeImport(
	moduleSpecifier string,
	moduleSpecifierKind modulespecifiers.ResultKind,
	usagePosition *lsproto.Position,
	exportInfo *SymbolExportInfo,
	isReExport *bool,
) *FixAddJsdocTypeImport {
	return &FixAddJsdocTypeImport{
		ImportFixBase: ImportFixBase{
			isReExport:          isReExport,
			exportInfo:          exportInfo,
			moduleSpecifierKind: moduleSpecifierKind,
			moduleSpecifier:     moduleSpecifier,
		},
		usagePosition: usagePosition,
	}
}

type FixAddToExistingImport struct {
	ImportFixBase
	importClauseOrBindingPattern *ast.Node  // ImportClause | ObjectBindingPattern
	importKind                   ImportKind // ImportKindDefault | ImportKindNamed
	addAsTypeOnly                AddAsTypeOnly
	propertyName                 string
}

func (f *FixAddToExistingImport) Kind() ImportFixKind {
	return ImportFixKindAddToExisting
}

func (f *FixAddToExistingImport) Base() *ImportFixBase {
	return &f.ImportFixBase
}

func getAddToExistingImport(
	importClauseOrBindingPattern *ast.Node,
	importKind ImportKind,
	moduleSpecifier string,
	moduleSpecifierKind modulespecifiers.ResultKind,
	addAsTypeOnly AddAsTypeOnly,
) *FixAddToExistingImport {
	return &FixAddToExistingImport{
		ImportFixBase: ImportFixBase{
			moduleSpecifierKind: moduleSpecifierKind,
			moduleSpecifier:     moduleSpecifier,
		},
		importClauseOrBindingPattern: importClauseOrBindingPattern,
		importKind:                   importKind,
		addAsTypeOnly:                addAsTypeOnly,
	}
}

type FixAddNewImport struct {
	ImportFixBase
	*Qualification
	importKind    ImportKind
	addAsTypeOnly AddAsTypeOnly
	propertyName  string
	useRequire    bool
}

func (f *FixAddNewImport) Kind() ImportFixKind {
	return ImportFixKindAddNew
}

func (f *FixAddNewImport) Base() *ImportFixBase {
	return &f.ImportFixBase
}

func getNewAddNewImport(
	moduleSpecifier string,
	moduleSpecifierKind modulespecifiers.ResultKind,
	importKind ImportKind,
	useRequire bool,
	addAsTypeOnly AddAsTypeOnly,
	exportInfo *SymbolExportInfo, // !!! | FutureSymbolExportInfo
	isReExport *bool,
	qualification *Qualification,
) *FixAddNewImport {
	return &FixAddNewImport{
		ImportFixBase: ImportFixBase{
			isReExport:          isReExport,
			exportInfo:          exportInfo,
			moduleSpecifierKind: modulespecifiers.ResultKindNone,
			moduleSpecifier:     moduleSpecifier,
		},
		// Qualification: qualification,
		importKind:    importKind,
		addAsTypeOnly: addAsTypeOnly,
		useRequire:    useRequire,
	}
}

type FixPromoteTypeOnlyImport struct {
	ImportFixBase
	typeOnlyAliasDeclaration *ast.Declaration // TypeOnlyAliasDeclaration
}

func (f *FixPromoteTypeOnlyImport) Kind() ImportFixKind {
	return ImportFixKindPromoteTypeOnly
}

func (f *FixPromoteTypeOnlyImport) Base() *ImportFixBase {
	return &f.ImportFixBase
}

/** Information needed to augment an existing import declaration. */
// rename all fixes to say fix at end
// rename to AddToExistingImportInfo
type FixAddToExistingImportInfo struct {
	declaration *ast.Declaration
	importKind  ImportKind
	targetFlags ast.SymbolFlags
	symbol      *ast.Symbol
}

func (info *FixAddToExistingImportInfo) getNewImportFromExistingSpecifier(
	isValidTypeOnlyUseSite bool,
	useRequire bool,
	ch *checker.Checker,
	compilerOptions *core.CompilerOptions,
) *FixAddNewImport {
	moduleSpecifier := checker.TryGetModuleSpecifierFromDeclaration(info.declaration)
	if moduleSpecifier == nil || moduleSpecifier.Text() == "" {
		return nil
	}
	addAsTypeOnly := AddAsTypeOnlyNotAllowed
	if !useRequire {
		addAsTypeOnly = getAddAsTypeOnly(isValidTypeOnlyUseSite, info.symbol, info.targetFlags, ch, compilerOptions)
	}
	return getNewAddNewImport(
		moduleSpecifier.Text(),
		modulespecifiers.ResultKindNone,
		info.importKind,
		useRequire,
		addAsTypeOnly,
		nil, // exportInfo
		nil, // isReExport
		nil, // qualification
	)
}
