package ls

import (
	"context"
	"slices"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
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
	return l.provideSourceDefinition(ctx, documentURI, position)
}

func (l *LanguageService) provideSourceDefinition(
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
	declarations := l.getSourceDefinitionDeclarations(ctx, program, file, node, pos)
	if len(declarations) == 0 {
		return lsproto.LocationOrLocationsOrDefinitionLinksOrNull{}, nil
	}
	return l.createDefinitionLocations(originSelectionRange, clientSupportsLink, declarations, nil /*reference*/), nil
}

func (l *LanguageService) getSourceDefinitionDeclarations(
	ctx context.Context,
	program *compiler.Program,
	currentFile *ast.SourceFile,
	node *ast.Node,
	position int,
) []*ast.Node {
	var declarations []*ast.Node
	moduleSpecifier := findContainingModuleSpecifier(node)

	if moduleSpecifier != nil {
		moduleDeclarations := l.getSourceDefinitionDeclarationsForModuleSpecifier(
			program,
			currentFile,
			moduleSpecifier,
			node,
			getSourceDefinitionNamesForNode(node, moduleSpecifier),
		)
		declarations = append(declarations, moduleDeclarations...)
		if len(moduleDeclarations) != 0 && (node == moduleSpecifier || isImportOrExportName(node)) {
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
		declarations = append(declarations, l.mapDeclarationToSourceDefinitions(program, currentFile, node, declaration)...)
	}

	reference := getReferenceAtPosition(currentFile, position, program)
	if reference != nil && len(declarations) == 0 {
		if sourceFile := l.getOrParseSourceFile(program, reference.fileName); sourceFile != nil && !tspath.IsDeclarationFileName(reference.fileName) {
			if len(sourceFile.Statements.Nodes) != 0 {
				declarations = append(declarations, sourceFile.Statements.Nodes[0].AsNode())
			} else {
				declarations = append(declarations, sourceFile.AsNode())
			}
		}
	}

	declarations = uniqueDeclarationNodes(declarations)
	if filtered := filterImportLikeDeclarations(currentFile, declarations); len(filtered) != 0 {
		declarations = filtered
	}
	if moduleSpecifier != nil && isImportOrExportName(node) && (len(declarations) == 0 || core.Every(declarations, func(declaration *ast.Node) bool {
		return ast.GetSourceFileOfNode(declaration).FileName() == currentFile.FileName() && isImportLikeDeclaration(declaration)
	})) {
		if fallback := l.getSourceDefinitionDeclarationsForModuleSpecifier(program, currentFile, moduleSpecifier, node, nil); len(fallback) != 0 {
			return uniqueDeclarationNodes(fallback)
		}
	}
	return declarations
}

func (l *LanguageService) getSourceDefinitionDeclarationsForModuleSpecifier(
	program *compiler.Program,
	currentFile *ast.SourceFile,
	moduleSpecifier *ast.Node,
	originalNode *ast.Node,
	names []string,
) []*ast.Node {
	implementationFile := l.resolveImplementationFileForModuleSpecifier(program, currentFile, moduleSpecifier)
	if implementationFile == "" {
		return nil
	}

	sourceFile := l.getOrParseSourceFile(program, implementationFile)
	if sourceFile == nil {
		return nil
	}
	if originalNode == moduleSpecifier || originalNode != nil && originalNode.Parent != nil && ast.IsImportClause(originalNode.Parent) && originalNode.Parent.Name() == originalNode {
		if len(sourceFile.Statements.Nodes) != 0 {
			return []*ast.Node{sourceFile.Statements.Nodes[0].AsNode()}
		}
		return []*ast.Node{sourceFile.AsNode()}
	}

	declarations := l.findSourceDefinitionDeclarationsInFile(program, implementationFile, names, map[string]struct{}{})
	if len(declarations) != 0 {
		return declarations
	}
	if len(names) != 0 && !isImportOrExportName(originalNode) {
		return nil
	}
	if len(sourceFile.Statements.Nodes) != 0 {
		return []*ast.Node{sourceFile.Statements.Nodes[0].AsNode()}
	}
	return []*ast.Node{sourceFile.AsNode()}
}

func isImportOrExportName(node *ast.Node) bool {
	if node == nil {
		return false
	}
	if ast.IsImportOrExportSpecifier(node) {
		return true
	}
	switch node.Kind {
	case ast.KindImportClause,
		ast.KindNamespaceImport,
		ast.KindNamedImports,
		ast.KindNamedExports:
		return true
	}
	if node.Parent == nil {
		return false
	}
	return (node.Parent.Kind == ast.KindImportClause || node.Parent.Kind == ast.KindNamespaceImport || ast.IsImportOrExportSpecifier(node.Parent)) && node.Parent.Name() == node
}

func (l *LanguageService) mapDeclarationToSourceDefinitions(
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

	implementationFile := l.resolveImplementationFileForDeclaration(program, currentFile, originalNode, declaration)
	if implementationFile == "" {
		return nil
	}

	sourceFile := l.getOrParseSourceFile(program, implementationFile)
	if sourceFile == nil {
		return nil
	}

	names := getCandidateSourceDeclarationNames(originalNode, declaration)
	declarations := l.findSourceDefinitionDeclarationsInFile(program, implementationFile, names, map[string]struct{}{})
	if len(declarations) != 0 {
		return filterPreferredSourceDeclarations(originalNode, declarations)
	}
	if len(names) != 0 {
		return nil
	}
	if len(sourceFile.Statements.Nodes) != 0 {
		return []*ast.Node{sourceFile.Statements.Nodes[0].AsNode()}
	}
	return []*ast.Node{sourceFile.AsNode()}
}

func (l *LanguageService) resolveImplementationFileForDeclaration(
	program *compiler.Program,
	currentFile *ast.SourceFile,
	originalNode *ast.Node,
	declaration *ast.Node,
) string {
	if moduleSpecifier := findContainingModuleSpecifier(originalNode); moduleSpecifier != nil {
		if implementationFile := l.resolveImplementationFileForModuleSpecifier(program, currentFile, moduleSpecifier); implementationFile != "" {
			return implementationFile
		}
	}
	if moduleSpecifier := findContainingModuleSpecifier(declaration); moduleSpecifier != nil {
		if implementationFile := l.resolveImplementationFileForModuleSpecifier(program, currentFile, moduleSpecifier); implementationFile != "" {
			return implementationFile
		}
	}

	_, startPos := getFileAndStartPosFromDeclaration(declaration)
	preferredMode := core.ModuleKindESNext
	if moduleSpecifier := findContainingModuleSpecifier(originalNode); moduleSpecifier != nil {
		preferredMode = program.GetModeForUsageLocation(currentFile, moduleSpecifier)
	}
	return l.findImplementationFileFromDtsFileName(program, ast.GetSourceFileOfNode(declaration).FileName(), currentFile.FileName(), preferredMode, startPos)
}

func (l *LanguageService) resolveImplementationFileForModuleSpecifier(
	program *compiler.Program,
	currentFile *ast.SourceFile,
	moduleSpecifier *ast.Node,
) string {
	options := *program.Options()
	options.NoDtsResolution = core.TSTrue
	resolver := module.NewResolver(program.Host(), &options, program.GetGlobalTypingsCacheLocation(), "")
	mode := program.GetModeForUsageLocation(currentFile, moduleSpecifier)

	if implementationFile := resolveImplementationFromModuleName(resolver, moduleSpecifier.Text(), currentFile.FileName(), mode); implementationFile != "" {
		return implementationFile
	}

	if resolved := program.GetResolvedModuleFromModuleSpecifier(currentFile, moduleSpecifier); resolved != nil && resolved.IsResolved() {
		if !tspath.IsDeclarationFileName(resolved.ResolvedFileName) {
			return resolved.ResolvedFileName
		}
		if implementationFile := l.findImplementationFileFromDtsFileName(program, resolved.ResolvedFileName, currentFile.FileName(), mode, core.TextPos(moduleSpecifier.Pos())); implementationFile != "" {
			return implementationFile
		}
	}

	return ""
}

func (l *LanguageService) findImplementationFileFromDtsFileName(
	program *compiler.Program,
	dtsFileName string,
	resolveFromFile string,
	preferredMode core.ResolutionMode,
	_ core.TextPos,
) string {
	options := *program.Options()
	options.NoDtsResolution = core.TSTrue
	resolver := module.NewResolver(program.Host(), &options, program.GetGlobalTypingsCacheLocation(), "")

	if jsExt := module.TryGetJSExtensionForFile(dtsFileName, &options); jsExt != "" {
		candidate := tspath.ChangeExtension(dtsFileName, jsExt)
		if program.Host().FS().FileExists(candidate) {
			return candidate
		}
	}

	parts := modulespecifiers.GetNodeModulePathParts(dtsFileName)
	if parts == nil {
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
	if tryPackageRootFirst {
		if implementationFile := tryResolvePackageSubpath(); implementationFile != "" {
			return implementationFile
		}
	}

	return ""
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
	return parser.ParseSourceFile(
		ast.SourceFileParseOptions{FileName: fileName, Path: l.toPath(fileName)},
		text,
		core.GetScriptKindFromFileName(fileName),
	)
}

func findContainingModuleSpecifier(node *ast.Node) *ast.Node {
	for current := node; current != nil; current = current.Parent {
		switch {
		case ast.IsImportDeclaration(current),
			ast.IsExportDeclaration(current),
			ast.IsJSImportDeclaration(current),
			ast.IsImportEqualsDeclaration(current),
			ast.IsRequireCall(current, true /*requireStringLiteralLikeArgument*/),
			ast.IsImportCall(current):
			if moduleSpecifier := ast.GetExternalModuleName(current); moduleSpecifier != nil && ast.IsStringLiteralLike(moduleSpecifier) {
				return moduleSpecifier
			}
		}
	}
	return nil
}

func getImportNamesForModuleSpecifier(moduleSpecifier *ast.Node) []string {
	var names []string
	for current := moduleSpecifier.Parent; current != nil; current = current.Parent {
		if !ast.IsImportDeclaration(current) {
			continue
		}
		if importClause := current.ImportClause(); importClause != nil {
			if importClause.Name() != nil {
				names = append(names, "default", importClause.Name().Text())
			}
			if namedBindings := importClause.AsImportClause().NamedBindings; namedBindings != nil {
				if ast.IsNamedImports(namedBindings) {
					for _, element := range namedBindings.AsNamedImports().Elements.Nodes {
						if propertyName := element.PropertyName(); propertyName != nil {
							names = append(names, propertyName.Text())
						} else if name := element.Name(); name != nil {
							names = append(names, name.Text())
						}
					}
				}
			}
		}
		break
	}
	return names
}

func getSourceDefinitionNamesForNode(node *ast.Node, moduleSpecifier *ast.Node) []string {
	if node == nil {
		return nil
	}

	names := getCandidateSourceDeclarationNames(node, nil)
	if node == moduleSpecifier {
		names = append(names, getImportNamesForModuleSpecifier(moduleSpecifier)...)
	} else if node.Parent != nil && ast.IsImportClause(node.Parent) && node.Parent.Name() == node {
		names = append(names, "default")
	}
	return core.Deduplicate(core.Filter(names, func(name string) bool { return name != "" }))
}

func (l *LanguageService) findSourceDefinitionDeclarationsInFile(
	program *compiler.Program,
	fileName string,
	names []string,
	seen map[string]struct{},
) []*ast.Node {
	if fileName == "" || len(names) == 0 {
		return nil
	}
	if _, ok := seen[fileName]; ok {
		return nil
	}
	seen[fileName] = struct{}{}

	sourceFile := l.getOrParseSourceFile(program, fileName)
	if sourceFile == nil {
		return nil
	}

	declarations := findDeclarationNodesByName(sourceFile, names)
	if len(declarations) != 0 {
		return declarations
	}

	var forwarded []*ast.Node
	for _, forwardedFile := range l.getForwardedImplementationFiles(program, sourceFile) {
		forwarded = append(forwarded, l.findSourceDefinitionDeclarationsInFile(program, forwardedFile, names, seen)...)
	}
	return uniqueDeclarationNodes(forwarded)
}

func (l *LanguageService) getForwardedImplementationFiles(program *compiler.Program, sourceFile *ast.SourceFile) []string {
	if sourceFile == nil {
		return nil
	}

	options := *program.Options()
	options.NoDtsResolution = core.TSTrue
	resolver := module.NewResolver(program.Host(), &options, program.GetGlobalTypingsCacheLocation(), "")
	type forwardedFile struct {
		fileName string
		pos      int
	}
	var files []forwardedFile
	addForwardedFile := func(moduleName string, pos int) {
		if moduleName == "" {
			return
		}
		if implementationFile := resolveImplementationFromModuleName(resolver, moduleName, sourceFile.FileName(), core.ModuleKindCommonJS); implementationFile != "" {
			files = append(files, forwardedFile{fileName: implementationFile, pos: pos})
		}
	}

	var visit ast.Visitor
	visit = func(node *ast.Node) bool {
		if node == nil {
			return false
		}
		switch {
		case ast.IsImportDeclaration(node),
			ast.IsExportDeclaration(node),
			ast.IsJSImportDeclaration(node),
			ast.IsImportEqualsDeclaration(node),
			ast.IsRequireCall(node, true /*requireStringLiteralLikeArgument*/),
			ast.IsImportCall(node):
			if moduleSpecifier := ast.GetExternalModuleName(node); moduleSpecifier != nil && ast.IsStringLiteralLike(moduleSpecifier) {
				addForwardedFile(moduleSpecifier.Text(), moduleSpecifier.Pos())
			}
		}
		return node.ForEachChild(visit)
	}
	sourceFile.AsNode().ForEachChild(visit)
	slices.SortFunc(files, func(a, b forwardedFile) int {
		switch {
		case a.pos < b.pos:
			return -1
		case a.pos > b.pos:
			return 1
		default:
			return 0
		}
	})
	return core.Deduplicate(core.Map(files, func(file forwardedFile) string { return file.fileName }))
}

func filterImportLikeDeclarations(currentFile *ast.SourceFile, declarations []*ast.Node) []*ast.Node {
	if currentFile == nil || len(declarations) <= 1 {
		return declarations
	}
	currentFileName := currentFile.FileName()
	if !core.Some(declarations, func(node *ast.Node) bool {
		return ast.GetSourceFileOfNode(node).FileName() != currentFileName
	}) {
		return declarations
	}
	filtered := core.Filter(declarations, func(node *ast.Node) bool {
		return ast.GetSourceFileOfNode(node).FileName() != currentFileName || !isImportLikeDeclaration(node)
	})
	if len(filtered) == 0 {
		return declarations
	}
	return filtered
}

func isImportLikeDeclaration(node *ast.Node) bool {
	for current := node; current != nil; current = current.Parent {
		switch {
		case ast.IsImportSpecifier(current),
			ast.IsImportClause(current),
			ast.IsNamespaceImport(current),
			ast.IsImportEqualsDeclaration(current),
			ast.IsImportDeclaration(current),
			ast.IsJSImportDeclaration(current):
			return true
		case ast.IsExportSpecifier(current):
			return current.Parent != nil && current.Parent.Parent != nil && current.Parent.Parent.ModuleSpecifier() != nil
		case ast.IsExportDeclaration(current) && current.ModuleSpecifier() != nil:
			return true
		case current.Kind == ast.KindSourceFile:
			return false
		}
	}
	return false
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

	var declarations []*ast.Node
	bestDepth := -1
	addMatch := func(nameNode *ast.Node, declaration *ast.Node) {
		depth := getDeclarationSearchDepth(nameNode)
		switch {
		case bestDepth == -1 || depth < bestDepth:
			declarations = []*ast.Node{declaration}
			bestDepth = depth
		case depth == bestDepth:
			declarations = append(declarations, declaration)
		}
	}

	wanted := make(map[string]struct{}, len(names))
	wantDefault := false
	for _, name := range names {
		if name == "default" {
			wantDefault = true
			continue
		}
		wanted[name] = struct{}{}
		for _, candidate := range getPossibleSymbolReferenceNodes(sourceFile, name, nil /*container*/) {
			if declaration := ast.GetDeclarationFromName(candidate); declaration != nil {
				addMatch(candidate, declaration)
			}
		}
	}
	var visit ast.Visitor
	visit = func(node *ast.Node) bool {
		if node == nil {
			return false
		}
		if name := ast.GetNameOfDeclaration(node); name != nil {
			if text := ast.GetTextOfPropertyName(name); text != "" {
				if _, ok := wanted[text]; ok {
					addMatch(name, node)
				}
			}
		}
		if wantDefault && node.Kind == ast.KindExportAssignment {
			addMatch(node, node)
		}
		return node.ForEachChild(visit)
	}
	sourceFile.AsNode().ForEachChild(visit)
	return uniqueDeclarationNodes(declarations)
}

func getDeclarationSearchDepth(node *ast.Node) int {
	depth := 0
	for current := node; current != nil; current = current.Parent {
		switch {
		case current.Kind == ast.KindSourceFile,
			current.Kind == ast.KindBlock,
			current.Kind == ast.KindModuleBlock,
			current.Kind == ast.KindCaseBlock,
			current.Kind == ast.KindExportAssignment,
			current.Kind == ast.KindJSExportAssignment,
			ast.IsFunctionLike(current),
			ast.IsClassLike(current):
			depth++
		}
	}
	return depth
}

func filterPreferredSourceDeclarations(originalNode *ast.Node, declarations []*ast.Node) []*ast.Node {
	if len(declarations) <= 1 || originalNode == nil {
		return declarations
	}
	if originalNode.Parent != nil && ast.IsAccessExpression(originalNode.Parent) && originalNode.Parent.Name() == originalNode {
		preferred := core.Filter(declarations, func(node *ast.Node) bool {
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
		if len(preferred) != 0 {
			return preferred
		}
	}

	preferred := core.Filter(declarations, func(node *ast.Node) bool {
		if !ast.IsDeclaration(node) || node.Kind == ast.KindExportAssignment || node.Kind == ast.KindJSExportAssignment {
			return false
		}
		if ast.IsBinaryExpression(node) && ast.GetAssignmentDeclarationKind(node) != ast.JSDeclarationKindNone {
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
	})
	if len(preferred) != 0 {
		return preferred
	}
	return declarations
}

func uniqueDeclarationNodes(nodes []*ast.Node) []*ast.Node {
	type declarationKey struct {
		fileName string
		pos      int
		end      int
	}
	seen := make(map[declarationKey]struct{}, len(nodes))
	result := make([]*ast.Node, 0, len(nodes))
	for _, node := range nodes {
		if node == nil {
			continue
		}
		fileName := ast.GetSourceFileOfNode(node).FileName()
		key := declarationKey{fileName: fileName, pos: node.Pos(), end: node.End()}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
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
	if len(sourceFile.Statements.Nodes) != 0 {
		return sourceFile.Statements.Nodes[0].AsNode()
	}
	return sourceFile.AsNode()
}
