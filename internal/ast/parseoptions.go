package ast

import (
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type SourceFileParseOptions struct {
	FileName string
	Path     tspath.Path
}

var isFileForcedToBeModuleByFormatExtensions = []string{tspath.ExtensionCjs, tspath.ExtensionCts, tspath.ExtensionMjs, tspath.ExtensionMts}

func isFileForcedToBeModuleByFormat(fileName string, options *core.CompilerOptions, metadata SourceFileMetaData) bool {
	// Excludes declaration files - they still require an explicit `export {}` or the like
	// for back compat purposes. The only non-declaration files _not_ forced to be a module are `.js` files
	// that aren't esm-mode (meaning not in a `type: module` scope).
	if GetImpliedNodeFormatForEmitWorker(fileName, options.GetEmitModuleKind(), metadata) == core.ModuleKindESNext || tspath.FileExtensionIsOneOf(fileName, isFileForcedToBeModuleByFormatExtensions) {
		return true
	}
	return false
}

// GetExternalModuleIndicator computes the full external module indicator for a file,
// combining syntax-based detection with compiler-options-dependent checks (JSX, forced module format).
// This should be called by the Program (or similar host) and NOT stored on the SourceFile.
func GetExternalModuleIndicator(file *SourceFile, options *core.CompilerOptions, metadata SourceFileMetaData) *Node {
	if file.ScriptKind == core.ScriptKindJSON {
		return nil
	}

	// Syntax-based indicator is pre-computed by the parser and stored on the SourceFile.
	// We use the cached field here instead of re-running IsFileProbablyExternalModule
	// to avoid reading Flags, which may be concurrently written by the binder on shared files.
	if file.SyntacticExternalModuleIndicator != nil {
		return file.SyntacticExternalModuleIndicator
	}

	if file.IsDeclarationFile {
		return nil
	}

	fileName := file.FileName()
	if tspath.IsDeclarationFileName(fileName) {
		return nil
	}

	switch options.GetEmitModuleDetectionKind() {
	case core.ModuleDetectionKindForce:
		// All non-declaration files are modules, declaration files still do the usual isFileProbablyExternalModule
		return file.AsNode()
	case core.ModuleDetectionKindLegacy:
		// Files are modules if they have imports, exports, or import.meta
		return nil
	case core.ModuleDetectionKindAuto:
		// If jsx is react-jsx or react-jsxdev then jsx tags force module-ness
		if options.Jsx == core.JsxEmitReactJSX || options.Jsx == core.JsxEmitReactJSXDev {
			if node := IsFileModuleFromUsingJSXTag(file); node != nil {
				return node
			}
		}
		// If module is nodenext or node16, all esm format files are modules
		if isFileForcedToBeModuleByFormat(fileName, options, metadata) {
			return file.AsNode()
		}
		return nil
	default:
		return nil
	}
}

// IsFileProbablyExternalModule checks if a file has explicit module syntax
// (import/export declarations or import.meta). This is the syntax-only check,
// independent of compiler options.
func IsFileProbablyExternalModule(sourceFile *SourceFile) *Node {
	for _, statement := range sourceFile.Statements.Nodes {
		if IsAnExternalModuleIndicatorNode(statement) {
			return statement
		}
	}
	return getImportMetaIfNecessary(sourceFile)
}

func IsAnExternalModuleIndicatorNode(node *Node) bool {
	return HasSyntacticModifier(node, ModifierFlagsExport) ||
		IsImportEqualsDeclaration(node) && IsExternalModuleReference(node.AsImportEqualsDeclaration().ModuleReference) ||
		IsImportDeclaration(node) || IsExportAssignment(node) || IsExportDeclaration(node)
}

func getImportMetaIfNecessary(sourceFile *SourceFile) *Node {
	if sourceFile.AsNode().Flags&NodeFlagsPossiblyContainsImportMeta != 0 {
		return findChildNode(sourceFile.AsNode(), IsImportMeta)
	}
	return nil
}

func findChildNode(root *Node, check func(*Node) bool) *Node {
	var result *Node
	var visit func(*Node) bool
	visit = func(node *Node) bool {
		if check(node) {
			result = node
			return true
		}
		return node.ForEachChild(visit)
	}
	visit(root)
	return result
}

func IsFileModuleFromUsingJSXTag(file *SourceFile) *Node {
	return walkTreeForJSXTags(file.AsNode())
}

// This is a somewhat unavoidable full tree walk to locate a JSX tag - `import.meta` requires the same,
// but we avoid that walk (or parts of it) if at all possible using the `PossiblyContainsImportMeta` node flag.
// Unfortunately, there's no `NodeFlag` space to do the same for JSX.
func walkTreeForJSXTags(node *Node) *Node {
	var found *Node

	var visitor func(node *Node) bool
	visitor = func(node *Node) bool {
		if found != nil {
			return true
		}
		if node.SubtreeFacts()&SubtreeContainsJsx == 0 {
			return false
		}
		if IsJsxOpeningElement(node) || IsJsxFragment(node) {
			found = node
			return true
		}
		return node.ForEachChild(visitor)
	}
	visitor(node)

	return found
}
