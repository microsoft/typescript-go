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
	declarations := l.getSourceDefinitionDeclarations(ctx, program, file, node)
	if len(declarations) == 0 {
		return l.provideDefinitionWorker(ctx, documentURI, position)
	}
	return l.createDefinitionLocations(originSelectionRange, clientSupportsLink, declarations, nil /*reference*/), nil
}

func (l *LanguageService) getSourceDefinitionDeclarations(
	ctx context.Context,
	program *compiler.Program,
	currentFile *ast.SourceFile,
	node *ast.Node,
) []*ast.Node {
	options := program.Options().Clone()
	options.NoDtsResolution = core.TSTrue
	resolver := module.NewResolver(program.Host(), options, program.GetGlobalTypingsCacheLocation(), "")

	var declarations []*ast.Node
	moduleSpecifier := findContainingModuleSpecifier(node)

	if moduleSpecifier != nil {
		moduleDeclarations := l.getSourceDefinitionDeclarationsForModuleSpecifier(
			resolver,
			program,
			currentFile,
			moduleSpecifier,
			node,
			getSourceDefinitionNamesForNode(node),
		)
		declarations = append(declarations, moduleDeclarations...)
		if shouldPreferModuleSpecifierResult(node, moduleSpecifier, moduleDeclarations) {
			return uniqueDeclarationNodes(declarations)
		}
	}

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

	for _, declaration := range definitionDeclarations {
		declarations = append(declarations, l.mapDeclarationToSourceDefinitions(resolver, program, currentFile, node, declaration)...)
	}

	declarations = uniqueDeclarationNodes(declarations)
	if len(getConcreteSourceDeclarations(declarations)) == 0 {
		return nil
	}
	return declarations
}

func (l *LanguageService) getSourceDefinitionDeclarationsForModuleSpecifier(
	resolver *module.Resolver,
	program *compiler.Program,
	currentFile *ast.SourceFile,
	moduleSpecifier *ast.Node,
	originalNode *ast.Node,
	names []string,
) []*ast.Node {
	implementationFile := l.resolveImplementationFileForModuleSpecifier(resolver, program, currentFile, moduleSpecifier)
	if implementationFile == "" {
		return nil
	}

	sourceFile := l.getOrParseSourceFile(program, implementationFile)
	if sourceFile == nil {
		return nil
	}
	if originalNode == moduleSpecifier {
		return getSourceDefinitionEntryDeclarations(sourceFile)
	}
	if isDefaultImportName(originalNode) {
		// For default imports, only search for "default" declarations to avoid
		// matching unrelated declarations with the same identifier name.
		defaultDeclarations := l.findSourceDefinitionDeclarationsInFile(resolver, program, implementationFile, []string{"default"}, &collections.Set[string]{})
		if len(defaultDeclarations) != 0 {
			return filterPreferredSourceDeclarations(originalNode, defaultDeclarations)
		}
		return getSourceDefinitionEntryDeclarations(sourceFile)
	}

	declarations := l.findSourceDefinitionDeclarationsInFile(resolver, program, implementationFile, names, &collections.Set[string]{})
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

func (l *LanguageService) mapDeclarationToSourceDefinitions(
	resolver *module.Resolver,
	program *compiler.Program,
	currentFile *ast.SourceFile,
	originalNode *ast.Node,
	declaration *ast.Node,
) []*ast.Node {
	file, startPos := getFileAndStartPosFromDeclaration(declaration)
	fileName := file.FileName()

	if mapped := l.tryGetSourcePosition(fileName, startPos); mapped != nil {
		if sourceFile := l.getOrParseSourceFile(program, mapped.FileName); sourceFile != nil {
			return []*ast.Node{findClosestDeclarationNode(sourceFile, mapped.Pos)}
		}
	}

	if !tspath.IsDeclarationFileName(fileName) {
		return []*ast.Node{declaration}
	}

	implementationFile := l.resolveImplementationFileForDeclaration(resolver, program, currentFile, originalNode, declaration)
	if implementationFile == "" {
		return nil
	}

	sourceFile := l.getOrParseSourceFile(program, implementationFile)
	if sourceFile == nil {
		return nil
	}

	names := getCandidateSourceDeclarationNames(originalNode, declaration)
	declarations := l.findSourceDefinitionDeclarationsInFile(resolver, program, implementationFile, names, &collections.Set[string]{})
	if len(declarations) != 0 {
		return filterPreferredSourceDeclarations(originalNode, declarations)
	}
	if len(names) != 0 {
		return nil
	}
	return getSourceDefinitionEntryDeclarations(sourceFile)
}

func (l *LanguageService) resolveImplementationFileForDeclaration(
	resolver *module.Resolver,
	program *compiler.Program,
	currentFile *ast.SourceFile,
	originalNode *ast.Node,
	declaration *ast.Node,
) string {
	originalModuleSpecifier := findContainingModuleSpecifier(originalNode)
	if originalModuleSpecifier != nil {
		if implementationFile := l.resolveImplementationFileForModuleSpecifier(resolver, program, currentFile, originalModuleSpecifier); implementationFile != "" {
			return implementationFile
		}
	}

	dtsFileName := ast.GetSourceFileOfNode(declaration).FileName()

	preferredMode := inferImpliedNodeFormat(resolver, dtsFileName)
	if originalModuleSpecifier != nil {
		preferredMode = program.GetModeForUsageLocation(currentFile, originalModuleSpecifier)
	}
	return l.findImplementationFileFromDtsFileName(resolver, program, dtsFileName, currentFile.FileName(), preferredMode)
}

func (l *LanguageService) resolveImplementationFileForModuleSpecifier(
	resolver *module.Resolver,
	program *compiler.Program,
	currentFile *ast.SourceFile,
	moduleSpecifier *ast.Node,
) string {
	mode := program.GetModeForUsageLocation(currentFile, moduleSpecifier)
	return resolveImplementationFromModuleName(resolver, moduleSpecifier.Text(), currentFile.FileName(), mode)
}

func (l *LanguageService) findImplementationFileFromDtsFileName(
	resolver *module.Resolver,
	program *compiler.Program,
	dtsFileName string,
	resolveFromFile string,
	preferredMode core.ResolutionMode,
) string {
	options := program.Options()

	if jsExt := module.TryGetJSExtensionForFile(dtsFileName, options); jsExt != "" {
		candidate := tspath.ChangeExtension(dtsFileName, jsExt)
		if program.Host().FS().FileExists(candidate) {
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
		return resolveImplementationFromModuleName(resolver, specifier, resolveFromFile, preferredMode)
	}

	tryPackageRootFirst := pathToFileInPackage == "index.d.ts" || strings.HasSuffix(pathToFileInPackage, "/index.d.ts")

	if !tryPackageRootFirst {
		if implementationFile := tryResolvePackageSubpath(); implementationFile != "" {
			return implementationFile
		}
	}
	if implementationFile := resolveImplementationFromModuleName(resolver, packageName, resolveFromFile, preferredMode); implementationFile != "" {
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

func (l *LanguageService) getOrParseSourceFile(program *compiler.Program, fileName string) *ast.SourceFile {
	if sourceFile := program.GetSourceFile(fileName); sourceFile != nil {
		return sourceFile
	}
	text, ok := l.ReadFile(fileName)
	if !ok {
		return nil
	}
	sourceFile := parser.ParseSourceFile(
		ast.SourceFileParseOptions{FileName: fileName, Path: l.toPath(fileName)},
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

func getSourceDefinitionNamesForNode(node *ast.Node) []string {
	names := getCandidateSourceDeclarationNames(node, nil)
	if isDefaultImportName(node) {
		names = append(names, "default")
	}
	return core.Deduplicate(core.Filter(names, func(name string) bool { return name != "" }))
}

func (l *LanguageService) findSourceDefinitionDeclarationsInFile(
	resolver *module.Resolver,
	program *compiler.Program,
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

	sourceFile := l.getOrParseSourceFile(program, fileName)
	if sourceFile == nil {
		return nil
	}

	declarations := findDeclarationNodesByName(sourceFile, names)
	if len(declarations) != 0 && len(getConcreteSourceDeclarations(declarations)) != 0 {
		return declarations
	}

	var forwarded []*ast.Node
	for _, forwardedFile := range l.getForwardedImplementationFiles(resolver, program, sourceFile) {
		forwarded = append(forwarded, l.findSourceDefinitionDeclarationsInFile(resolver, program, forwardedFile, names, seen)...)
	}
	if len(forwarded) != 0 {
		if len(getConcreteSourceDeclarations(forwarded)) != 0 {
			return uniqueDeclarationNodes(forwarded)
		}
		return uniqueDeclarationNodes(append(slices.Clip(declarations), forwarded...))
	}
	return declarations
}

func (l *LanguageService) getForwardedImplementationFiles(resolver *module.Resolver, program *compiler.Program, sourceFile *ast.SourceFile) []string {
	preferredMode := inferImpliedNodeFormat(resolver, sourceFile.FileName())

	var files []string
	for _, imp := range sourceFile.Imports() {
		moduleName := imp.Text()
		if implementationFile := resolveImplementationFromModuleName(resolver, moduleName, sourceFile.FileName(), preferredMode); implementationFile != "" {
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
