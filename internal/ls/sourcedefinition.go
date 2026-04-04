package ls

import (
	"context"
	"math"
	"slices"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/binder"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/module"
	"github.com/microsoft/typescript-go/internal/modulespecifiers"
	"github.com/microsoft/typescript-go/internal/parser"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

func (l *LanguageService) ProvideSourceDefinition(
	ctx context.Context,
	documentURI lsproto.DocumentUri,
	position lsproto.Position,
) (lsproto.DefinitionResponse, error) {
	caps := lsproto.GetClientCapabilities(ctx)
	clientSupportsLink := caps.TextDocument.Definition.LinkSupport

	program, file := l.getProgramAndFile(documentURI)
	pos := int(l.converters.LineAndCharacterToPosition(file, position))
	node := astnav.GetTouchingPropertyName(file, pos)
	if node.Kind == ast.KindSourceFile {
		return lsproto.LocationOrLocationsOrDefinitionLinksOrNull{}, nil
	}

	originSelectionRange := l.createLspRangeFromNode(node, file)
	declarations := l.newSourceDefResolver(ctx, program, file, node).resolve(node)
	if len(declarations) == 0 {
		return l.provideDefinitionWorker(ctx, documentURI, position)
	}
	return l.createDefinitionLocations(originSelectionRange, clientSupportsLink, declarations, nil /*reference*/), nil
}

// sourceDefResolver resolves source definitions by mapping .d.ts declarations
// to their implementation files (.js/.ts). It holds values derived from the
// Program upfront so that inner methods never access the Program directly,
// preventing accidental use of Program methods on externally-parsed source files.
type sourceDefResolver struct {
	ls            *LanguageService
	resolver      *module.Resolver
	options       *core.CompilerOptions
	fs            vfs.FS
	getSourceFile func(string) *ast.SourceFile

	// moduleSpecifier is the containing module specifier node for the
	// original cursor position, if one exists. Used to resolve the
	// implementation file from the import statement.
	moduleSpecifier *ast.Node
	// specifierMode is the pre-computed resolution mode for moduleSpecifier,
	// derived from Program.GetModeForUsageLocation with the current file.
	specifierMode core.ResolutionMode
	// resolveFrom is the file name to resolve module names relative to
	// (always currentFile.FileName()).
	resolveFrom string

	// definitionDeclarations are the checker-derived declarations for the
	// node at the cursor position. Computed during construction so the
	// checker is not held for the duration of the resolution.
	definitionDeclarations []*ast.Node
}

func (l *LanguageService) newSourceDefResolver(
	ctx context.Context,
	program *compiler.Program,
	currentFile *ast.SourceFile,
	node *ast.Node,
) *sourceDefResolver {
	options := program.Options().Clone()
	options.NoDtsResolution = core.TSTrue

	moduleSpecifier := findContainingModuleSpecifier(node)

	r := &sourceDefResolver{
		ls:              l,
		resolver:        module.NewResolver(program.Host(), options, program.GetGlobalTypingsCacheLocation(), ""),
		options:         options,
		fs:              program.Host().FS(),
		getSourceFile:   program.GetSourceFile,
		moduleSpecifier: moduleSpecifier,
		resolveFrom:     currentFile.FileName(),
	}
	if moduleSpecifier != nil {
		r.specifierMode = program.GetModeForUsageLocation(currentFile, moduleSpecifier)
	}
	r.definitionDeclarations = getDefinitionDeclarationsFromChecker(ctx, program, currentFile, node)
	return r
}

func (r *sourceDefResolver) resolve(node *ast.Node) []*ast.Node {
	var declarations []*ast.Node

	if r.moduleSpecifier != nil {
		moduleDeclarations := r.resolveFromModuleSpecifier(
			r.moduleSpecifier,
			node,
			getSourceDefinitionNamesForNode(node),
		)
		declarations = append(declarations, moduleDeclarations...)
		if shouldPreferModuleSpecifierResult(node, r.moduleSpecifier, moduleDeclarations) {
			return uniqueDeclarationNodes(declarations)
		}
	}

	for _, declaration := range r.definitionDeclarations {
		declarations = append(declarations, r.mapDeclarationToSource(node, declaration)...)
	}

	declarations = uniqueDeclarationNodes(declarations)
	if len(getConcreteSourceDeclarations(declarations)) == 0 {
		// Fallback for property access on imported values where the property
		// has no checker declarations (e.g. mapped type properties). Trace
		// the parent expression back to its import to find the implementation file.
		if node.Parent != nil && ast.IsAccessExpression(node.Parent) && node.Parent.Name() == node {
			if fallback := r.resolvePropertyViaParentImport(node); len(fallback) != 0 {
				return fallback
			}
		}
		return nil
	}
	return declarations
}

func getDefinitionDeclarationsFromChecker(
	ctx context.Context,
	program *compiler.Program,
	currentFile *ast.SourceFile,
	node *ast.Node,
) []*ast.Node {
	c, done := program.GetTypeCheckerForFile(ctx, currentFile)
	defer done()

	definitionDeclarations := getDeclarationsFromLocation(c, node)
	if len(definitionDeclarations) == 0 && node.Parent != nil && ast.IsAccessExpression(node.Parent) && node.Parent.Name() == node {
		if left := node.Parent.Expression(); left != nil {
			if prop := c.GetPropertyOfType(c.GetTypeAtLocation(left), node.Text()); prop != nil {
				definitionDeclarations = prop.Declarations
			}
		}
	}
	if calledDeclaration := tryGetSignatureDeclaration(c, node); calledDeclaration != nil {
		nonFunctionDeclarations := core.Filter(definitionDeclarations, func(node *ast.Node) bool { return !ast.IsFunctionLike(node) })
		definitionDeclarations = append(nonFunctionDeclarations, calledDeclaration)
	}
	return definitionDeclarations
}

func (r *sourceDefResolver) resolveFromModuleSpecifier(
	moduleSpecifier *ast.Node,
	originalNode *ast.Node,
	names []string,
) []*ast.Node {
	implementationFile := resolveImplementationFromModuleName(r.resolver, moduleSpecifier.Text(), r.resolveFrom, r.specifierMode)
	if implementationFile == "" {
		return nil
	}

	sourceFile := r.getOrParseSourceFile(implementationFile)
	if sourceFile == nil {
		return nil
	}
	if originalNode == moduleSpecifier {
		return getSourceDefinitionEntryDeclarations(sourceFile)
	}
	if isDefaultImportName(originalNode) {
		// For default imports, only search for "default" declarations to avoid
		// matching unrelated declarations with the same identifier name.
		defaultDeclarations := r.findDeclarationsInFile(implementationFile, []string{"default"}, &collections.Set[string]{})
		if len(defaultDeclarations) != 0 {
			return filterPreferredSourceDeclarations(originalNode, defaultDeclarations)
		}
		return getSourceDefinitionEntryDeclarations(sourceFile)
	}

	declarations := r.findDeclarationsInFile(implementationFile, names, &collections.Set[string]{})
	if len(declarations) != 0 {
		return filterPreferredSourceDeclarations(originalNode, declarations)
	}
	return getSourceDefinitionEntryDeclarations(sourceFile)
}

func isDefaultImportName(node *ast.Node) bool {
	if node == nil || node.Parent == nil || !ast.IsImportClause(node.Parent) || node.Parent.Name() != node || node.Parent.Parent == nil {
		return false
	}
	return ast.IsDefaultImport(node.Parent.Parent)
}

func shouldPreferModuleSpecifierResult(node *ast.Node, moduleSpecifier *ast.Node, declarations []*ast.Node) bool {
	if moduleSpecifier == nil || len(declarations) == 0 {
		return false
	}
	if node == moduleSpecifier {
		return true
	}
	if ast.IsPartOfTypeNode(node) || ast.IsPartOfTypeOnlyImportOrExportDeclaration(node) {
		return len(getConcreteSourceDeclarations(declarations)) != 0
	}
	return true
}

func getSourceDefinitionEntryNode(sourceFile *ast.SourceFile) *ast.Node {
	if len(sourceFile.Statements.Nodes) != 0 {
		return sourceFile.Statements.Nodes[0].AsNode()
	}
	return sourceFile.AsNode()
}

func getSourceDefinitionEntryDeclarations(sourceFile *ast.SourceFile) []*ast.Node {
	return []*ast.Node{getSourceDefinitionEntryNode(sourceFile)}
}

func (r *sourceDefResolver) mapDeclarationToSource(
	originalNode *ast.Node,
	declaration *ast.Node,
) []*ast.Node {
	file, startPos := getFileAndStartPosFromDeclaration(declaration)
	fileName := file.FileName()

	if mapped := r.ls.tryGetSourcePosition(fileName, startPos); mapped != nil {
		if sourceFile := r.getOrParseSourceFile(mapped.FileName); sourceFile != nil {
			return []*ast.Node{findClosestDeclarationNode(sourceFile, mapped.Pos)}
		}
	}

	if !tspath.IsDeclarationFileName(fileName) {
		return []*ast.Node{declaration}
	}

	implementationFile := r.resolveImplementationFileForDeclaration(declaration)
	if implementationFile == "" {
		return nil
	}

	sourceFile := r.getOrParseSourceFile(implementationFile)
	if sourceFile == nil {
		return nil
	}

	names := getCandidateSourceDeclarationNames(originalNode, declaration)
	declarations := r.findDeclarationsInFile(implementationFile, names, &collections.Set[string]{})
	if len(declarations) != 0 {
		return filterPreferredSourceDeclarations(originalNode, declarations)
	}
	if len(names) != 0 {
		return nil
	}
	return getSourceDefinitionEntryDeclarations(sourceFile)
}

func (r *sourceDefResolver) resolveImplementationFileForDeclaration(
	declaration *ast.Node,
) string {
	if r.moduleSpecifier != nil {
		if implementationFile := resolveImplementationFromModuleName(r.resolver, r.moduleSpecifier.Text(), r.resolveFrom, r.specifierMode); implementationFile != "" {
			return implementationFile
		}
	}

	dtsFileName := ast.GetSourceFileOfNode(declaration).FileName()

	preferredMode := inferImpliedNodeFormat(r.resolver, dtsFileName)
	if r.moduleSpecifier != nil {
		preferredMode = r.specifierMode
	}
	return r.findImplementationFileFromDtsFileName(dtsFileName, preferredMode)
}

func (r *sourceDefResolver) findImplementationFileFromDtsFileName(
	dtsFileName string,
	preferredMode core.ResolutionMode,
) string {
	if jsExt := module.TryGetJSExtensionForFile(dtsFileName, r.options); jsExt != "" {
		candidate := tspath.ChangeExtension(dtsFileName, jsExt)
		if r.fs.FileExists(candidate) {
			return candidate
		}
	}

	parts := modulespecifiers.GetNodeModulePathParts(dtsFileName)
	if parts == nil {
		return ""
	}

	// Ensure the file only contains one /node_modules/ segment. If there's more
	// than one, the package name extraction may be incorrect, so bail out.
	if strings.LastIndex(dtsFileName, "/node_modules/") != parts.TopLevelNodeModulesIndex {
		return ""
	}

	packageNamePathPart := dtsFileName[parts.TopLevelPackageNameIndex+1 : parts.PackageRootIndex]
	packageName := module.GetPackageNameFromTypesPackageName(module.UnmangleScopedPackageName(packageNamePathPart))
	if packageName == "" {
		return ""
	}

	pathToFileInPackage := dtsFileName[parts.PackageRootIndex+1:]
	tryResolvePackageSubpath := func() string {
		if pathToFileInPackage == "" {
			return ""
		}
		specifier := packageName + "/" + tspath.RemoveFileExtension(pathToFileInPackage)
		return resolveImplementationFromModuleName(r.resolver, specifier, r.resolveFrom, preferredMode)
	}

	tryPackageRootFirst := pathToFileInPackage == "index.d.ts" || strings.HasSuffix(pathToFileInPackage, "/index.d.ts")

	if !tryPackageRootFirst {
		if implementationFile := tryResolvePackageSubpath(); implementationFile != "" {
			return implementationFile
		}
	}
	if implementationFile := resolveImplementationFromModuleName(r.resolver, packageName, r.resolveFrom, preferredMode); implementationFile != "" {
		return implementationFile
	}
	if !tryPackageRootFirst {
		return ""
	}
	return tryResolvePackageSubpath()
}

func resolveImplementationFromModuleName(
	resolver *module.Resolver,
	moduleName string,
	resolveFromFile string,
	preferredMode core.ResolutionMode,
) string {
	modes := []core.ResolutionMode{preferredMode}
	if preferredMode != core.ModuleKindESNext {
		modes = append(modes, core.ModuleKindESNext)
	}
	if preferredMode != core.ModuleKindCommonJS {
		modes = append(modes, core.ModuleKindCommonJS)
	}

	for _, mode := range modes {
		resolved, _ := resolver.ResolveModuleName(moduleName, resolveFromFile, mode, nil)
		if resolved != nil && resolved.IsResolved() && !tspath.IsDeclarationFileName(resolved.ResolvedFileName) {
			return resolved.ResolvedFileName
		}
	}
	return ""
}

func (r *sourceDefResolver) getOrParseSourceFile(fileName string) *ast.SourceFile {
	if sourceFile := r.getSourceFile(fileName); sourceFile != nil {
		return sourceFile
	}
	text, ok := r.ls.ReadFile(fileName)
	if !ok {
		return nil
	}
	sourceFile := parser.ParseSourceFile(
		ast.SourceFileParseOptions{FileName: fileName, Path: r.ls.toPath(fileName)},
		text,
		core.GetScriptKindFromFileName(fileName),
	)
	binder.BindSourceFile(sourceFile)
	return sourceFile
}

// inferImpliedNodeFormat determines the module format for a source file that may not be
// in the program, using the file extension and nearest package.json "type" field.
func inferImpliedNodeFormat(resolver *module.Resolver, fileName string) core.ResolutionMode {
	var packageJsonType string
	if scope := resolver.GetPackageScopeForPath(tspath.GetDirectoryPath(fileName)); scope.Exists() {
		if value, ok := scope.Contents.Type.GetValue(); ok {
			packageJsonType = value
		}
	}
	return ast.GetImpliedNodeFormatForFile(fileName, packageJsonType)
}

func findContainingModuleSpecifier(node *ast.Node) *ast.Node {
	for current := node; current != nil; current = current.Parent {
		if ast.IsAnyImportOrReExport(current) || ast.IsRequireCall(current, true /*requireStringLiteralLikeArgument*/) || ast.IsImportCall(current) {
			if moduleSpecifier := ast.GetExternalModuleName(current); moduleSpecifier != nil && ast.IsStringLiteralLike(moduleSpecifier) {
				return moduleSpecifier
			}
		}
	}
	return nil
}

// resolvePropertyViaParentImport is a fallback for property access on imported
// values where the property itself has no checker declarations (e.g. mapped type
// properties). It walks the access chain to find the root identifier, looks up
// its import declaration, and searches the implementation file for the property name.
func (r *sourceDefResolver) resolvePropertyViaParentImport(node *ast.Node) []*ast.Node {
	// Walk left in the access chain to find the root identifier.
	expr := node.Parent.Expression()
	for expr != nil && ast.IsAccessExpression(expr) {
		expr = expr.Expression()
	}
	if expr == nil || !ast.IsIdentifier(expr) {
		return nil
	}

	// Find the import declaration for this identifier in the current file.
	currentFile := ast.GetSourceFileOfNode(node)
	moduleSpecifier := findImportModuleSpecifierForName(currentFile, expr.Text())
	if moduleSpecifier == nil || !ast.IsStringLiteralLike(moduleSpecifier) {
		return nil
	}

	// Resolve the module to an implementation file.
	preferredMode := inferImpliedNodeFormat(r.resolver, r.resolveFrom)
	implementationFile := resolveImplementationFromModuleName(r.resolver, moduleSpecifier.Text(), r.resolveFrom, preferredMode)
	if implementationFile == "" {
		return nil
	}

	// Search the implementation file for the property name.
	propertyName := node.Text()
	declarations := r.findDeclarationsInFile(implementationFile, []string{propertyName}, &collections.Set[string]{})
	if len(declarations) != 0 {
		return filterPreferredSourceDeclarations(node, declarations)
	}
	return nil
}

// findImportModuleSpecifierForName searches the source file for an import
// declaration that imports the given identifier name and returns its module specifier.
func findImportModuleSpecifierForName(sourceFile *ast.SourceFile, name string) *ast.Node {
	for _, stmt := range sourceFile.Statements.Nodes {
		if !ast.IsImportDeclaration(stmt) {
			continue
		}
		importDecl := stmt.AsImportDeclaration()
		if importDecl.ImportClause == nil {
			continue
		}
		clause := importDecl.ImportClause.AsImportClause()
		// Default import: import name from "pkg"
		if clause.Name() != nil && clause.Name().Text() == name {
			return importDecl.ModuleSpecifier
		}
		if clause.NamedBindings == nil {
			continue
		}
		if ast.IsNamespaceImport(clause.NamedBindings) {
			// import * as name from "pkg"
			ns := clause.NamedBindings.AsNamespaceImport()
			if ns.Name() != nil && ns.Name().Text() == name {
				return importDecl.ModuleSpecifier
			}
		} else if ast.IsNamedImports(clause.NamedBindings) {
			// import { name } from "pkg" or import { orig as name } from "pkg"
			for _, spec := range clause.NamedBindings.AsNamedImports().Elements.Nodes {
				if spec.AsImportSpecifier().Name().Text() == name {
					return importDecl.ModuleSpecifier
				}
			}
		}
	}
	return nil
}

func getSourceDefinitionNamesForNode(node *ast.Node) []string {
	names := getCandidateSourceDeclarationNames(node, nil)
	if isDefaultImportName(node) {
		names = append(names, "default")
	}
	// For aliased imports (import { original as alias }) and re-exports
	// (export { original as alias }), include the original export name
	// so we search the .js file for the right declaration.
	if node != nil && node.Parent != nil {
		if ast.IsImportSpecifier(node.Parent) || ast.IsExportSpecifier(node.Parent) {
			if propName := node.Parent.PropertyName(); propName != nil {
				names = append(names, propName.Text())
			}
		}
	}
	return core.Deduplicate(core.Filter(names, func(name string) bool { return name != "" }))
}

func (r *sourceDefResolver) findDeclarationsInFile(
	fileName string,
	names []string,
	seen *collections.Set[string],
) []*ast.Node {
	if fileName == "" || len(names) == 0 {
		return nil
	}
	if !seen.AddIfAbsent(fileName) {
		return nil
	}

	sourceFile := r.getOrParseSourceFile(fileName)
	if sourceFile == nil {
		return nil
	}

	declarations := findDeclarationNodesByName(sourceFile, names)
	if len(declarations) != 0 && len(getConcreteSourceDeclarations(declarations)) != 0 {
		return declarations
	}

	var forwarded []*ast.Node
	for _, forwardedFile := range r.getForwardedImplementationFiles(sourceFile) {
		forwarded = append(forwarded, r.findDeclarationsInFile(forwardedFile, names, seen)...)
	}
	if len(forwarded) != 0 {
		if len(getConcreteSourceDeclarations(forwarded)) != 0 {
			return uniqueDeclarationNodes(forwarded)
		}
		return uniqueDeclarationNodes(append(slices.Clip(declarations), forwarded...))
	}
	return declarations
}

func (r *sourceDefResolver) getForwardedImplementationFiles(sourceFile *ast.SourceFile) []string {
	preferredMode := inferImpliedNodeFormat(r.resolver, sourceFile.FileName())

	var files []string
	for _, imp := range sourceFile.Imports() {
		moduleName := imp.Text()
		if implementationFile := resolveImplementationFromModuleName(r.resolver, moduleName, sourceFile.FileName(), preferredMode); implementationFile != "" {
			files = append(files, implementationFile)
		}
	}
	return core.Deduplicate(files)
}

func getCandidateSourceDeclarationNames(originalNode *ast.Node, declaration *ast.Node) []string {
	var names []string
	if declaration != nil {
		if name := ast.GetNameOfDeclaration(declaration); name != nil {
			if text := ast.GetTextOfPropertyName(name); text != "" {
				names = append(names, text)
			}
		}
		if declaration.Kind == ast.KindExportAssignment {
			names = append(names, "default")
		}
		if (ast.IsFunctionDeclaration(declaration) || ast.IsClassDeclaration(declaration)) && declaration.ModifierFlags()&ast.ModifierFlagsExportDefault == ast.ModifierFlagsExportDefault {
			names = append(names, "default")
		}
		// For aliased import/export specifiers (e.g. import { original as alias }),
		// also include the original name (propertyName) so we search for it in .js.
		if ast.IsImportSpecifier(declaration) || ast.IsExportSpecifier(declaration) {
			if propName := declaration.PropertyName(); propName != nil {
				names = append(names, propName.Text())
			}
		}
	}
	if originalNode != nil {
		switch {
		case ast.IsIdentifier(originalNode), ast.IsPrivateIdentifier(originalNode):
			names = append(names, originalNode.Text())
		case ast.IsStringLiteralLike(originalNode):
			text := originalNode.Text()
			if text != "" && !strings.ContainsRune(text, '/') && !tspath.IsExternalModuleNameRelative(text) {
				names = append(names, text)
			}
		}
	}
	return names
}

func findDeclarationNodesByName(sourceFile *ast.SourceFile, names []string) []*ast.Node {
	names = core.Deduplicate(core.Filter(names, func(name string) bool { return name != "" }))
	if len(names) == 0 {
		return nil
	}

	var wanted collections.Set[string]
	wantDefault := false
	for _, name := range names {
		if name == "default" {
			wantDefault = true
			continue
		}
		wanted.Add(name)
	}

	type candidate struct {
		node  *ast.Node
		depth int
	}
	var candidates []candidate
	minDepth := math.MaxInt

	var visit ast.Visitor
	visit = func(node *ast.Node) bool {
		matched := false
		if name := ast.GetNameOfDeclaration(node); name != nil {
			if text := ast.GetTextOfPropertyName(name); text != "" {
				if wanted.Has(text) {
					matched = true
				}
			}
		}
		if wantDefault && node.Kind == ast.KindExportAssignment {
			matched = true
		}
		if wantDefault && (ast.IsFunctionDeclaration(node) || ast.IsClassDeclaration(node)) && node.ModifierFlags()&ast.ModifierFlagsExportDefault == ast.ModifierFlagsExportDefault {
			matched = true
		}
		if matched {
			depth := getContainerDepth(node)
			candidates = append(candidates, candidate{node: node, depth: depth})
			if depth < minDepth {
				minDepth = depth
			}
		}
		return node.ForEachChild(visit)
	}
	sourceFile.AsNode().ForEachChild(visit)

	// Only keep declarations at the shallowest depth, like getTopMostDeclarationNamesInFile.
	var declarations []*ast.Node
	for _, c := range candidates {
		if c.depth == minDepth {
			declarations = append(declarations, c.node)
		}
	}
	return uniqueDeclarationNodes(declarations)
}

// getContainerDepth counts the number of container nodes above a declaration,
// matching the behavior of getDepth in getTopMostDeclarationNamesInFile.
func getContainerDepth(node *ast.Node) int {
	depth := 0
	current := node
	for current != nil {
		current = getContainerNode(current)
		depth++
	}
	return depth
}

func filterPreferredSourceDeclarations(originalNode *ast.Node, declarations []*ast.Node) []*ast.Node {
	if len(declarations) <= 1 || originalNode == nil {
		return declarations
	}
	if preferred := getPropertyLikeSourceDeclarations(originalNode, declarations); len(preferred) != 0 {
		return preferred
	}
	if preferred := getConcreteSourceDeclarations(declarations); len(preferred) != 0 {
		return preferred
	}
	return declarations
}

func getPropertyLikeSourceDeclarations(originalNode *ast.Node, declarations []*ast.Node) []*ast.Node {
	if originalNode.Parent == nil || !ast.IsAccessExpression(originalNode.Parent) || originalNode.Parent.Name() != originalNode {
		return nil
	}
	return core.Filter(declarations, func(node *ast.Node) bool {
		switch node.Kind {
		case ast.KindPropertyAssignment,
			ast.KindShorthandPropertyAssignment,
			ast.KindPropertyDeclaration,
			ast.KindPropertySignature,
			ast.KindMethodDeclaration,
			ast.KindMethodSignature,
			ast.KindGetAccessor,
			ast.KindSetAccessor,
			ast.KindEnumMember:
			return true
		default:
			return false
		}
	})
}

func getConcreteSourceDeclarations(declarations []*ast.Node) []*ast.Node {
	return core.Filter(declarations, isConcreteSourceDeclaration)
}

func isConcreteSourceDeclaration(node *ast.Node) bool {
	if !ast.IsDeclaration(node) || node.Kind == ast.KindExportAssignment || node.Kind == ast.KindJSExportAssignment {
		return false
	}
	if (ast.IsBinaryExpression(node) || ast.IsCallExpression(node)) && ast.GetAssignmentDeclarationKind(node) != ast.JSDeclarationKindNone {
		return false
	}
	switch node.Kind {
	case ast.KindParameter,
		ast.KindTypeParameter,
		ast.KindBindingElement,
		ast.KindImportClause,
		ast.KindImportSpecifier,
		ast.KindNamespaceImport,
		ast.KindExportSpecifier,
		ast.KindPropertyAccessExpression,
		ast.KindElementAccessExpression,
		ast.KindCommonJSExport:
		return false
	default:
		return true
	}
}

func uniqueDeclarationNodes(nodes []*ast.Node) []*ast.Node {
	type declarationKey struct {
		fileName string
		loc      core.TextRange
	}
	var seen collections.Set[declarationKey]
	result := make([]*ast.Node, 0, len(nodes))
	for _, node := range nodes {
		if node == nil {
			continue
		}
		fileName := ast.GetSourceFileOfNode(node).FileName()
		key := declarationKey{fileName: fileName, loc: node.Loc}
		if !seen.AddIfAbsent(key) {
			continue
		}
		result = append(result, node)
	}
	return result
}

func findClosestDeclarationNode(sourceFile *ast.SourceFile, pos int) *ast.Node {
	node := astnav.GetTouchingPropertyName(sourceFile, pos)
	for current := node; current != nil; current = current.Parent {
		if ast.IsDeclaration(current) || current.Kind == ast.KindExportAssignment {
			return current
		}
	}
	return getSourceDefinitionEntryNode(sourceFile)
}
