package compiler

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/microsoft/typescript-go/internal/compiler/diagnostics"
)

type ProgramOptions struct {
	RootPath       string
	Host           CompilerHost
	Options        *CompilerOptions
	SingleThreaded bool
}

type Program struct {
	host                        CompilerHost
	options                     *CompilerOptions
	rootPath                    string
	files                       []*SourceFile
	filesByPath                 map[string]*SourceFile
	nodeModules                 map[string]*SourceFile
	checker                     *Checker
	usesUriStyleNodeCoreModules Tristate
	currentNodeModulesDepth     int
}

var extensions = []string{".ts", ".tsx"}

func NewProgram(options ProgramOptions) *Program {
	p := &Program{}
	p.options = options.Options
	if p.options == nil {
		p.options = &CompilerOptions{}
	}
	p.host = options.Host
	if p.host == nil {
		p.host = NewCompilerHost(p.options, options.SingleThreaded)
	}
	rootPath := options.RootPath
	if rootPath == "" {
		rootPath = "."
	}
	p.rootPath = p.host.AbsFileName(rootPath)
	fileInfos := p.host.ReadDirectory(rootPath, extensions)
	// Sort files by descending file size
	slices.SortFunc(fileInfos, func(a FileInfo, b FileInfo) int {
		return int(b.Size) - int(a.Size)
	})
	p.parseSourceFiles(fileInfos)
	return p
}

func (p *Program) SourceFiles() []*SourceFile { return p.files }
func (p *Program) Options() *CompilerOptions  { return p.options }
func (p *Program) Host() CompilerHost         { return p.host }

func (p *Program) parseSourceFiles(fileInfos []FileInfo) {
	p.files = make([]*SourceFile, len(fileInfos))[:len(fileInfos)]
	for i := range fileInfos {
		p.host.RunTask(func() {
			fileName := fileInfos[i].Name
			text, _ := p.host.ReadFile(fileName)
			sourceFile := ParseSourceFile(fileName, text, getEmitScriptTarget(p.options))
			sourceFile.path, _ = filepath.Abs(fileName)
			p.collectExternalModuleReferences(sourceFile)
			p.files[i] = sourceFile
		})
	}
	p.host.WaitForTasks()
	p.filesByPath = make(map[string]*SourceFile)
	for _, file := range p.files {
		p.filesByPath[file.path] = file
	}
}

func (p *Program) bindSourceFiles() {
	for _, file := range p.files {
		if !file.isBound {
			p.host.RunTask(func() {
				bindSourceFile(file, p.options)
			})
		}
	}
	p.host.WaitForTasks()
}

func (p *Program) getResolvedModule(currentSourceFile *SourceFile, moduleReference string) *SourceFile {
	directory := filepath.Dir(currentSourceFile.path)
	if isExternalModuleNameRelative(moduleReference) {
		return p.findSourceFile(filepath.Join(directory, moduleReference))
	}
	return p.findNodeModule(moduleReference)
}

func (p *Program) findSourceFile(candidate string) *SourceFile {
	extensionless := removeFileExtension(candidate)
	for _, ext := range []string{ExtensionTs, ExtensionTsx, ExtensionDts} {
		path := extensionless + ext
		if result, ok := p.filesByPath[path]; ok {
			return result
		}
	}
	return nil
}

func (p *Program) findNodeModule(moduleReference string) *SourceFile {
	if p.nodeModules == nil {
		p.nodeModules = make(map[string]*SourceFile)
	}
	if sourceFile, ok := p.nodeModules[moduleReference]; ok {
		return sourceFile
	}
	sourceFile := p.tryLoadNodeModule(filepath.Join(p.rootPath, "node_modules", moduleReference))
	if sourceFile == nil {
		sourceFile = p.tryLoadNodeModule(filepath.Join(p.rootPath, "node_modules/@types", moduleReference))
	}
	p.nodeModules[moduleReference] = sourceFile
	return sourceFile
}

func (p *Program) tryLoadNodeModule(modulePath string) *SourceFile {
	if packageJson, ok := p.host.ReadFile(filepath.Join(modulePath, "package.json")); ok {
		var jsonMap map[string]any
		if json.Unmarshal([]byte(packageJson), &jsonMap) == nil {
			typesValue := jsonMap["types"]
			if typesValue == nil {
				typesValue = jsonMap["typings"]
			}
			if fileName, ok := typesValue.(string); ok {
				path := filepath.Join(modulePath, fileName)
				return p.filesByPath[path]
			}
		}
	}
	return nil
}

func (p *Program) GetSyntacticDiagnostics(sourceFile *SourceFile) []*Diagnostic {
	return p.getDiagnosticsHelper(sourceFile, p.getSyntaticDiagnosticsForFile)
}

func (p *Program) GetBindDiagnostics(sourceFile *SourceFile) []*Diagnostic {
	p.bindSourceFiles()
	return p.getDiagnosticsHelper(sourceFile, p.getBindDiagnosticsForFile)
}

func (p *Program) GetSemanticDiagnostics(sourceFile *SourceFile) []*Diagnostic {
	return p.getDiagnosticsHelper(sourceFile, p.getSemanticDiagnosticsForFile)
}

func (p *Program) GetGlobalDiagnostics() []*Diagnostic {
	return sortAndDeduplicateDiagnostics(p.getTypeChecker().GetGlobalDiagnostics())
}

func (p *Program) getTypeChecker() *Checker {
	if p.checker == nil {
		p.checker = NewChecker(p)
	}
	return p.checker
}

func (p *Program) getSyntaticDiagnosticsForFile(sourceFile *SourceFile) []*Diagnostic {
	return sourceFile.diagnostics
}

func (p *Program) getBindDiagnosticsForFile(sourceFile *SourceFile) []*Diagnostic {
	return sourceFile.bindDiagnostics
}

func (p *Program) getSemanticDiagnosticsForFile(sourceFile *SourceFile) []*Diagnostic {
	return p.getTypeChecker().GetDiagnostics(sourceFile)
}

func (p *Program) getDiagnosticsHelper(sourceFile *SourceFile, getDiagnostics func(*SourceFile) []*Diagnostic) []*Diagnostic {
	if sourceFile != nil {
		return sortAndDeduplicateDiagnostics(getDiagnostics(sourceFile))
	}
	var result []*Diagnostic
	for _, file := range p.files {
		result = append(result, getDiagnostics(file)...)
	}
	return sortAndDeduplicateDiagnostics(result)
}

func (p *Program) PrintTypeAliases() {
	for _, file := range p.files {
		if filepath.Base(file.fileName) == "main.ts" {
			file.AsNode().ForEachChild(p.printTypeAlias)
		}
	}
}

func (p *Program) printTypeAlias(node *Node) bool {
	if isTypeAliasDeclaration(node) {
		fmt.Println(p.getTypeChecker().typeAliasToString(node.AsTypeAliasDeclaration()))
	}
	return node.ForEachChild(p.printTypeAlias)
}

func (p *Program) collectExternalModuleReferences(file *SourceFile) {
	if file.moduleReferencesProcessed {
		return
	}
	file.moduleReferencesProcessed = true
	// !!!
	// If we are importing helpers, we need to add a synthetic reference to resolve the
	// helpers library. (A JavaScript file without `externalModuleIndicator` set might be
	// a CommonJS module; `commonJsModuleIndicator` doesn't get set until the binder has
	// run. We synthesize a helpers import for it just in case; it will never be used if
	// the binder doesn't find and set a `commonJsModuleIndicator`.)
	// if (isJavaScriptFile || (!file.isDeclarationFile && (getIsolatedModules(options) || isExternalModule(file)))) {
	// 	if (options.importHelpers) {
	// 		// synthesize 'import "tslib"' declaration
	// 		imports = [createSyntheticImport(externalHelpersModuleNameText, file)];
	// 	}
	// 	const jsxImport = getJSXRuntimeImport(getJSXImplicitImportBase(options, file), options);
	// 	if (jsxImport) {
	// 		// synthesize `import "base/jsx-runtime"` declaration
	// 		(imports ||= []).push(createSyntheticImport(jsxImport, file));
	// 	}
	// }
	for _, node := range file.statements {
		p.collectModuleReferences(file, node, false /*inAmbientModule*/)
	}
	// if ((file.flags & NodeFlags.PossiblyContainsDynamicImport) || isJavaScriptFile) {
	// 	collectDynamicImportOrRequireOrJsDocImportCalls(file);
	// }
	// function collectDynamicImportOrRequireOrJsDocImportCalls(file: SourceFile) {
	// 	const r = /import|require/g;
	// 	while (r.exec(file.text) !== null) { // eslint-disable-line no-restricted-syntax
	// 		const node = getNodeAtPosition(file, r.lastIndex);
	// 		if (isJavaScriptFile && isRequireCall(node, /*requireStringLiteralLikeArgument*/ true)) {
	// 			setParentRecursive(node, /*incremental*/ false); // we need parent data on imports before the program is fully bound, so we ensure it's set here
	// 			imports = append(imports, node.arguments[0]);
	// 		}
	// 		// we have to check the argument list has length of at least 1. We will still have to process these even though we have parsing error.
	// 		else if (isImportCall(node) && node.arguments.length >= 1 && isStringLiteralLike(node.arguments[0])) {
	// 			setParentRecursive(node, /*incremental*/ false); // we need parent data on imports before the program is fully bound, so we ensure it's set here
	// 			imports = append(imports, node.arguments[0]);
	// 		}
	// 		else if (isLiteralImportTypeNode(node)) {
	// 			setParentRecursive(node, /*incremental*/ false); // we need parent data on imports before the program is fully bound, so we ensure it's set here
	// 			imports = append(imports, node.argument.literal);
	// 		}
	// 		else if (isJavaScriptFile && isJSDocImportTag(node)) {
	// 			const moduleNameExpr = getExternalModuleName(node);
	// 			if (moduleNameExpr && isStringLiteral(moduleNameExpr) && moduleNameExpr.text) {
	// 				setParentRecursive(node, /*incremental*/ false);
	// 				imports = append(imports, moduleNameExpr);
	// 			}
	// 		}
	// 	}
	// }
	// /** Returns a token if position is in [start-of-leading-trivia, end), includes JSDoc only in JS files */
	// function getNodeAtPosition(sourceFile: SourceFile, position: number): Node {
	// 	let current: Node = sourceFile;
	// 	const getContainingChild = (child: Node) => {
	// 		if (child.pos <= position && (position < child.end || (position === child.end && (child.kind === SyntaxKind.EndOfFileToken)))) {
	// 			return child;
	// 		}
	// 	};
	// 	while (true) {
	// 		const child = isJavaScriptFile && hasJSDocNodes(current) && forEach(current.jsDoc, getContainingChild) || forEachChild(current, getContainingChild);
	// 		if (!child) {
	// 			return current;
	// 		}
	// 		current = child;
	// 	}
	// }
}

var unprefixedNodeCoreModules = map[string]bool{
	"assert":              true,
	"assert/strict":       true,
	"async_hooks":         true,
	"buffer":              true,
	"child_process":       true,
	"cluster":             true,
	"console":             true,
	"constants":           true,
	"crypto":              true,
	"dgram":               true,
	"diagnostics_channel": true,
	"dns":                 true,
	"dns/promises":        true,
	"domain":              true,
	"events":              true,
	"fs":                  true,
	"fs/promises":         true,
	"http":                true,
	"http2":               true,
	"https":               true,
	"inspector":           true,
	"inspector/promises":  true,
	"module":              true,
	"net":                 true,
	"os":                  true,
	"path":                true,
	"path/posix":          true,
	"path/win32":          true,
	"perf_hooks":          true,
	"process":             true,
	"punycode":            true,
	"querystring":         true,
	"readline":            true,
	"readline/promises":   true,
	"repl":                true,
	"stream":              true,
	"stream/consumers":    true,
	"stream/promises":     true,
	"stream/web":          true,
	"string_decoder":      true,
	"sys":                 true,
	"test/mock_loader":    true,
	"timers":              true,
	"timers/promises":     true,
	"tls":                 true,
	"trace_events":        true,
	"tty":                 true,
	"url":                 true,
	"util":                true,
	"util/types":          true,
	"v8":                  true,
	"vm":                  true,
	"wasi":                true,
	"worker_threads":      true,
	"zlib":                true,
}

var exclusivelyPrefixedNodeCoreModules = map[string]bool{
	"node:sea":            true,
	"node:sqlite":         true,
	"node:test":           true,
	"node:test/reporters": true,
}

func (p *Program) collectModuleReferences(file *SourceFile, node *Statement, inAmbientModule bool) {
	if isAnyImportOrReExport(node) {
		moduleNameExpr := getExternalModuleName(node)
		// TypeScript 1.0 spec (April 2014): 12.1.6
		// An ExternalImportDeclaration in an AmbientExternalModuleDeclaration may reference other external modules
		// only through top - level external module names. Relative external module names are not permitted.
		if moduleNameExpr != nil && isStringLiteral(moduleNameExpr) {
			moduleName := moduleNameExpr.AsStringLiteral().text
			if moduleName != "" && (!inAmbientModule || !isExternalModuleNameRelative(moduleName)) {
				setParentInChildren(node) // we need parent data on imports before the program is fully bound, so we ensure it's set here
				file.imports = append(file.imports, moduleNameExpr)
				if file.usesUriStyleNodeCoreModules != TSTrue && p.currentNodeModulesDepth == 0 && !file.isDeclarationFile {
					if strings.HasPrefix(moduleName, "node:") && !exclusivelyPrefixedNodeCoreModules[moduleName] {
						// Presence of `node:` prefix takes precedence over unprefixed node core modules
						file.usesUriStyleNodeCoreModules = TSTrue
					} else if file.usesUriStyleNodeCoreModules == TSUnknown && unprefixedNodeCoreModules[moduleName] {
						// Avoid `unprefixedNodeCoreModules.has` for every import
						file.usesUriStyleNodeCoreModules = TSFalse
					}
				}
			}
		}
		return
	}
	if isModuleDeclaration(node) && isAmbientModule(node) && (inAmbientModule || hasSyntacticModifier(node, ModifierFlagsAmbient) || file.isDeclarationFile) {
		setParentInChildren(node)
		nameText := getTextOfIdentifierOrLiteral(node.AsModuleDeclaration().name)
		// Ambient module declarations can be interpreted as augmentations for some existing external modules.
		// This will happen in two cases:
		// - if current file is external module then module augmentation is a ambient module declaration defined in the top level scope
		// - if current file is not external module then module augmentation is an ambient module declaration with non-relative module name
		//   immediately nested in top level ambient module declaration .
		if isExternalModule(file) || (inAmbientModule && !isExternalModuleNameRelative(nameText)) {
			file.moduleAugmentations = append(file.moduleAugmentations, node.AsModuleDeclaration().name)
		} else if !inAmbientModule {
			if file.isDeclarationFile {
				// for global .d.ts files record name of ambient module
				file.ambientModuleNames = append(file.ambientModuleNames, nameText)
			}
			// An AmbientExternalModuleDeclaration declares an external module.
			// This type of declaration is permitted only in the global module.
			// The StringLiteral must specify a top - level external module name.
			// Relative external module names are not permitted
			// NOTE: body of ambient module is always a module block, if it exists
			if node.AsModuleDeclaration().body != nil {
				for _, statement := range node.AsModuleDeclaration().body.AsModuleBlock().statements {
					p.collectModuleReferences(file, statement, true /*inAmbientModule*/)
				}
			}
		}
	}
}

type DiagnosticsFormattingOptions struct {
	CurrentDirectory     string
	NewLine              string
	GetCanonicalFileName func(fileName string) string
}

const (
	foregroundColorEscapeGrey   = "\u001b[90m"
	foregroundColorEscapeRed    = "\u001b[91m"
	foregroundColorEscapeYellow = "\u001b[93m"
	foregroundColorEscapeBlue   = "\u001b[94m"
	foregroundColorEscapeCyan   = "\u001b[96m"
)

const (
	gutterStyleSequence = "\u001b[7m"
	gutterSeparator     = " "
	resetEscapeSequence = "\u001b[0m"
	ellipsis            = "..."
)

func FormatDiagnosticsWithColorAndContext(diags []*Diagnostic, formatOpts *DiagnosticsFormattingOptions) string {
	if len(diags) == 0 {
		return ""
	}

	var output strings.Builder

	for _, diagnostic := range diags {
		if diagnostic.file != nil {
			file := diagnostic.file
			pos := diagnostic.loc.Pos()
			writeLocation(&output, file, pos, formatOpts, writeWithStyleAndReset)
			output.WriteString(" - ")
		}

		writeWithStyleAndReset(&output, diagnostic.Category().Name(), getCategoryFormat(diagnostic.Category()))
		fmt.Fprintf(&output, "%s TS%d: %s", foregroundColorEscapeGrey, diagnostic.Code(), resetEscapeSequence)
		WriteFlattenedDiagnosticMessage(&output, diagnostic, formatOpts.NewLine, 0 /*indent*/)

		if diagnostic.File() != nil && diagnostic.Code() != diagnostics.File_appears_to_be_binary.Code() {
			output.WriteString(formatOpts.NewLine)
			writeCodeSnippet(&output, diagnostic.File(), diagnostic.Pos(), diagnostic.Length(), getCategoryFormat(diagnostic.Category()), formatOpts)
		}

		if (diagnostic.RelatedInformation() != nil) && (len(diagnostic.RelatedInformation()) > 0) {
			output.WriteString(formatOpts.NewLine)
			for _, relatedInformation := range diagnostic.RelatedInformation() {
				file := relatedInformation.File()
				if file != nil {
					output.WriteString(formatOpts.NewLine)
					pos := relatedInformation.Pos()
					writeLocation(&output, file, pos, formatOpts, writeWithStyleAndReset)
					writeCodeSnippet(&output, file, pos, relatedInformation.Length(), foregroundColorEscapeCyan, formatOpts)
				}

				output.WriteString(formatOpts.NewLine)
				WriteFlattenedDiagnosticMessage(&output, relatedInformation, formatOpts.NewLine, 0 /*indent*/)
			}
		}
	}

	return output.String()
}

func writeCodeSnippet(writer *strings.Builder, sourceFile *SourceFile, start int, length int, squiggleColor string, formatOpts *DiagnosticsFormattingOptions) {
	firstLine, firstLineChar := GetLineAndCharacterOfPosition(sourceFile, start)
	lastLine, lastLineChar := GetLineAndCharacterOfPosition(sourceFile, start+length)

	lastLineOfFile, _ := GetLineAndCharacterOfPosition(sourceFile, len(sourceFile.text))

	hasMoreThanFiveLines := lastLine-firstLine >= 4
	gutterWidth := len(strconv.Itoa(lastLineOfFile + 1))

	for i := firstLine; i <= lastLine; i++ {
		writer.WriteString(formatOpts.NewLine)

		// If the error spans over 5 lines, we'll only show the first 2 and last 2 lines,
		// so we'll skip ahead to the second-to-last line.
		if hasMoreThanFiveLines && firstLine+1 < i && i < lastLine-1 {
			writer.WriteString(gutterStyleSequence)
			fmt.Fprintf(writer, "%*s", gutterWidth, ellipsis)
			writer.WriteString(resetEscapeSequence)
			writer.WriteString(gutterSeparator)
			writer.WriteString(formatOpts.NewLine)
			i = lastLine - 1
		}

		lineStart := GetPositionOfLineAndCharacter(sourceFile, i, 0)
		lineEnd := sourceFile.loc.end
		if i < lastLineOfFile {
			lineEnd = GetPositionOfLineAndCharacter(sourceFile, i+1, 0)
		}
		lineContent := strings.TrimRightFunc(sourceFile.text[lineStart:lineEnd], unicode.IsSpace) // trim from end
		lineContent = strings.ReplaceAll(lineContent, "\t", " ")                                  // convert tabs to single spaces

		// Output the gutter and the actual contents of the line.
		writer.WriteString(gutterStyleSequence)
		fmt.Fprintf(writer, "%*d", gutterWidth, i+1)
		writer.WriteString(resetEscapeSequence)
		writer.WriteString(gutterSeparator)
		writer.WriteString(lineContent)
		writer.WriteString(formatOpts.NewLine)

		// Output the gutter and the error span for the line using tildes.
		writer.WriteString(gutterStyleSequence)
		fmt.Fprintf(writer, "%*s", gutterWidth, "")
		writer.WriteString(resetEscapeSequence)
		writer.WriteString(gutterSeparator)
		writer.WriteString(squiggleColor)
		if i == firstLine {
			// If we're on the last line, then limit it to the last character of the last line.
			// Otherwise, we'll just squiggle the rest of the line, giving 'slice' no end position.
			lastCharForLine := ifElse(i == lastLine, lastLineChar, len(lineContent))

			// Fill with spaces until the first character,
			// then squiggle the remainder of the line.
			writer.WriteString(strings.Repeat(" ", firstLineChar))
			writer.WriteString(strings.Repeat("~", lastCharForLine-firstLineChar))
		} else if i == lastLine {
			// Squiggle until the final character.
			writer.WriteString(strings.Repeat("~", lastLineChar))
		} else {
			// Squiggle the entire line.
			writer.WriteString(strings.Repeat("~", len(lineContent)))
		}

		writer.WriteString(resetEscapeSequence)
	}
}

func WriteFlattenedDiagnosticMessage(writer *strings.Builder, diagnostic *Diagnostic, newline string, level int) {
	writer.WriteString(diagnostic.Message())

	for _, chain := range diagnostic.messageChain {
		flattenDiagnosticMessageChain(writer, chain, newline, level+1)
	}
}

func flattenDiagnosticMessageChain(writer *strings.Builder, chain *MessageChain, newline string, level int) {
	writer.WriteString(newline)
	for i := 0; i < level; i++ {
		writer.WriteString("  ")
	}

	writer.WriteString(chain.message)
	for _, child := range chain.messageChain {
		flattenDiagnosticMessageChain(writer, child, newline, level+1)
	}
}

func getCategoryFormat(category diagnostics.Category) string {
	switch category {
	case diagnostics.CategoryError:
		return foregroundColorEscapeRed
	case diagnostics.CategoryWarning:
		return foregroundColorEscapeYellow
	case diagnostics.CategorySuggestion:
		return foregroundColorEscapeGrey
	case diagnostics.CategoryMessage:
		return foregroundColorEscapeBlue
	}
	panic("Unhandled diagnostic category")
}

type FormattedWriter func(output *strings.Builder, text string, formatStyle string)

func writeWithStyleAndReset(output *strings.Builder, text string, formatStyle string) {
	output.WriteString(formatStyle)
	output.WriteString(text)
	output.WriteString(resetEscapeSequence)
}

func writeLocation(output *strings.Builder, file *SourceFile, pos int, formatOpts *DiagnosticsFormattingOptions, writeWithStyleAndReset FormattedWriter) {
	firstLine, firstChar := GetLineAndCharacterOfPosition(file, pos)
	var relativeFileName string
	if formatOpts != nil {
		relativeFileName = ConvertToRelativePath(file.path, formatOpts.CurrentDirectory, formatOpts.GetCanonicalFileName)
	} else {
		relativeFileName = file.path
	}

	writeWithStyleAndReset(output, relativeFileName, foregroundColorEscapeCyan)
	output.WriteRune(':')
	writeWithStyleAndReset(output, strconv.Itoa(firstLine+1), foregroundColorEscapeYellow)
	output.WriteRune(':')
	writeWithStyleAndReset(output, strconv.Itoa(firstChar+1), foregroundColorEscapeYellow)
}
