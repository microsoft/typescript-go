package declarations

import (
	"fmt"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler/diagnostics"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/jsnum"
	"github.com/microsoft/typescript-go/internal/nodebuilder"
	"github.com/microsoft/typescript-go/internal/printer"
	"github.com/microsoft/typescript-go/internal/transformers"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type ReferencedFilePair struct {
	file *ast.SourceFile
	ref  *ast.FileReference
}

type OutputPaths interface {
	DeclarationFilePath() string
	JsFilePath() string
}

// Used to be passed in the TransformationContext, which is now just an EmitContext
type DeclarationEmitHost interface {
	GetCurrentDirectory() string
	UseCaseSensitiveFileNames() bool
	GetSourceFileFromReference(origin *ast.SourceFile, ref *ast.FileReference) *ast.SourceFile

	GetOutputPathsFor(file *ast.SourceFile, forceDtsPaths bool) OutputPaths
	GetResolutionModeOverride(node *ast.Node) core.ResolutionMode
	GetEffectiveDeclarationFlags(node *ast.Node, flags ast.ModifierFlags) ast.ModifierFlags
}

type DeclarationTransformer struct {
	transformers.Transformer
	host            DeclarationEmitHost
	compilerOptions *core.CompilerOptions
	diagnostics     []*ast.Diagnostic
	tracker         *SymbolTrackerImpl
	state           *SymbolTrackerSharedState
	resolver        printer.EmitResolver

	declarationFilePath string
	declarationMapPath  string

	isBundledEmit                    bool
	needsDeclare                     bool
	needsScopeFixMarker              bool
	resultHasScopeMarker             bool
	enclosingDeclaration             *ast.Node
	getSymbolAccessibilityDiagnostic GetSymbolAccessibilityDiagnostic
	resultHasExternalModuleIndicator bool
	suppressNewDiagnosticContexts    bool
	lateStatementReplacementMap      map[ast.NodeId]*ast.Node
	rawReferencedFiles               []ReferencedFilePair
	rawTypeReferenceDirectives       []*ast.FileReference
	rawLibReferenceDirectives        []*ast.FileReference
}

func NewDeclarationTransformer(host DeclarationEmitHost, resolver printer.EmitResolver, context *printer.EmitContext, compilerOptions *core.CompilerOptions, declarationFilePath string, declarationMapPath string) *DeclarationTransformer {
	shared := &SymbolTrackerSharedState{isolatedDeclarations: compilerOptions.IsolatedDeclarations.IsTrue(), resolver: resolver}
	tracker := NewSymbolTracker(resolver, shared)
	// TODO: Use new host GetOutputPathsFor method instead of passing in entrypoint paths (which will also better support bundled emit)
	tx := &DeclarationTransformer{compilerOptions: compilerOptions, tracker: tracker, state: shared, declarationFilePath: declarationFilePath, declarationMapPath: declarationMapPath, host: host}
	tx.NewTransformer(tx.visit, context)
	return tx
}

func (tx *DeclarationTransformer) GetDiagnostics() []*ast.Diagnostic {
	return tx.diagnostics
}

const declarationEmitNodeBuilderFlags = nodebuilder.FlagsMultilineObjectLiterals |
	nodebuilder.FlagsWriteClassExpressionAsTypeLiteral |
	nodebuilder.FlagsUseTypeOfFunction |
	nodebuilder.FlagsUseStructuralFallback |
	nodebuilder.FlagsAllowEmptyTuple |
	nodebuilder.FlagsGenerateNamesForShadowedTypeParams |
	nodebuilder.FlagsNoTruncation

const declarationEmitInternalNodeBuilderFlags = nodebuilder.InternalFlagsAllowUnresolvedNames

// functions as both `visitDeclarationStatements` and `transformRoot`, utilitzing SyntaxList nodes
func (tx *DeclarationTransformer) visit(node *ast.Node) *ast.Node {
	// !!! TODO: Bundle support?
	switch node.Kind {
	case ast.KindSourceFile:
		return tx.visitSourceFile(node.AsSourceFile())
	default:
		return tx.visitDeclarationStatements(node)
	}
}

func throwDiagnostic(result printer.SymbolAccessibilityResult) *SymbolAccessibilityDiagnostic {
	panic("Diagnostic emitted without context")
}

func (tx *DeclarationTransformer) visitSourceFile(node *ast.SourceFile) *ast.Node {
	if node.IsDeclarationFile {
		return node.AsNode()
	}

	tx.isBundledEmit = false
	tx.needsDeclare = true
	tx.needsScopeFixMarker = false
	tx.resultHasScopeMarker = false
	tx.enclosingDeclaration = node.AsNode()
	tx.getSymbolAccessibilityDiagnostic = throwDiagnostic
	tx.resultHasExternalModuleIndicator = false
	tx.suppressNewDiagnosticContexts = false
	tx.state.lateMarkedStatements = make([]*ast.Node, 0)
	tx.lateStatementReplacementMap = make(map[ast.NodeId]*ast.Node)
	tx.rawReferencedFiles = make([]ReferencedFilePair, 0)
	tx.rawTypeReferenceDirectives = make([]*ast.FileReference, 0)
	tx.rawLibReferenceDirectives = make([]*ast.FileReference, 0)
	tx.state.currentSourceFile = node
	tx.collectFileReferences(node)
	updated := tx.transformSourceFile(node)
	tx.state.currentSourceFile = nil
	return updated
}

func (tx *DeclarationTransformer) collectFileReferences(sourceFile *ast.SourceFile) {
	tx.rawReferencedFiles = append(tx.rawReferencedFiles, core.Map(sourceFile.ReferencedFiles, func(ref *ast.FileReference) ReferencedFilePair { return ReferencedFilePair{file: sourceFile, ref: ref} })...)
	tx.rawTypeReferenceDirectives = append(tx.rawTypeReferenceDirectives, sourceFile.TypeReferenceDirectives...)
	tx.rawLibReferenceDirectives = append(tx.rawLibReferenceDirectives, sourceFile.LibReferenceDirectives...)
}

func (tx *DeclarationTransformer) transformSourceFile(node *ast.SourceFile) *ast.Node {
	var combinedStatements *ast.StatementList
	if ast.IsSourceFileJS(node) {
		// !!! TODO: JS declaration emit support
		combinedStatements = tx.Factory().NewNodeList([]*ast.Node{})
	} else {
		statements := tx.Visitor().VisitNodes(node.Statements)
		combinedStatements = tx.transformAndReplaceLatePaintedStatements(statements)
		combinedStatements.Loc = statements.Loc // setTextRange
		if ast.IsExternalModule(node) && (!tx.resultHasExternalModuleIndicator || (tx.needsScopeFixMarker && !tx.resultHasScopeMarker)) {
			marker := createEmptyExports(tx.Factory())
			newList := append(combinedStatements.Nodes, marker)
			withMarker := tx.Factory().NewNodeList(newList)
			withMarker.Loc = combinedStatements.Loc
			combinedStatements = withMarker
		}
	}
	outputFilePath := tspath.GetDirectoryPath(tspath.NormalizeSlashes(tx.declarationFilePath))
	result := tx.Factory().UpdateSourceFile(node, combinedStatements)
	result.AsSourceFile().LibReferenceDirectives = tx.getLibReferences()
	result.AsSourceFile().TypeReferenceDirectives = tx.getTypeReferences()
	result.AsSourceFile().HasNoDefaultLib = node.HasNoDefaultLib
	result.AsSourceFile().IsDeclarationFile = true
	result.AsSourceFile().ReferencedFiles = tx.getReferencedFiles(outputFilePath)
	return result.AsNode()
}

func createEmptyExports(factory *ast.NodeFactory) *ast.Node {
	return factory.NewExportDeclaration(nil /*isTypeOnly*/, false, factory.NewNamedExports(factory.NewNodeList([]*ast.Node{})), nil, nil)
}

func (tx *DeclarationTransformer) transformAndReplaceLatePaintedStatements(statements *ast.StatementList) *ast.StatementList {
	// This is a `while` loop because `handleSymbolAccessibilityError` can see additional import aliases marked as visible during
	// error handling which must now be included in the output and themselves checked for errors.
	// For example:
	// ```
	// module A {
	//   export module Q {}
	//   import B = Q;
	//   import C = B;
	//   export import D = C;
	// }
	// ```
	// In such a scenario, only Q and D are initially visible, but we don't consider imports as private names - instead we say they if they are referenced they must
	// be recorded. So while checking D's visibility we mark C as visible, then we must check C which in turn marks B, completing the chain of
	// dependent imports and allowing a valid declaration file output. Today, this dependent alias marking only happens for internal import aliases.
	for true {
		if len(tx.state.lateMarkedStatements) == 0 {
			break
		}

		next := tx.state.lateMarkedStatements[0]
		tx.state.lateMarkedStatements = tx.state.lateMarkedStatements[1:]

		priorNeedsDeclare := tx.needsDeclare
		tx.needsDeclare = next.Parent != nil && ast.IsSourceFile(next.Parent) && !(ast.IsExternalModule(next.Parent.AsSourceFile()) && tx.isBundledEmit)

		result := tx.transformTopLevelDeclaration(next)

		tx.needsDeclare = priorNeedsDeclare
		original := tx.EmitContext().MostOriginal(next)
		id := ast.GetNodeId(original)
		tx.lateStatementReplacementMap[id] = result
	}

	// And lastly, we need to get the final form of all those indetermine import declarations from before and add them to the output list
	// (and remove them from the set to examine for outter declarations)
	results := make([]*ast.Node, 0, len(statements.Nodes))
	for _, statement := range statements.Nodes {
		if !isLateVisibilityPaintedStatement(statement) {
			results = append(results, statement)
			continue
		}
		original := tx.EmitContext().MostOriginal(statement)
		id := ast.GetNodeId(original)
		replacement, ok := tx.lateStatementReplacementMap[id]
		if !ok || replacement == nil {
			continue // not replaced, elide
		}
		if replacement.Kind == ast.KindSyntaxList {
			if !tx.needsScopeFixMarker || !tx.resultHasExternalModuleIndicator {
				for _, elem := range replacement.AsSyntaxList().Children {
					if needsScopeMarker(elem) {
						tx.needsScopeFixMarker = true
					}
					if ast.IsSourceFile(statement.Parent) && ast.IsExternalModuleIndicator(replacement) {
						tx.resultHasExternalModuleIndicator = true
					}
				}
			}
			results = append(results, replacement.AsSyntaxList().Children...)
		} else {
			if needsScopeMarker(replacement) {
				tx.needsScopeFixMarker = true
			}
			if ast.IsSourceFile(statement.Parent) && ast.IsExternalModuleIndicator(replacement) {
				tx.resultHasExternalModuleIndicator = true
			}
			results = append(results, replacement)
		}
	}

	return tx.Factory().NewNodeList(results)
}

func (tx *DeclarationTransformer) getReferencedFiles(outputFilePath string) (results []*ast.FileReference) {
	// Handle path rewrites for triple slash ref comments
	for _, pair := range tx.rawReferencedFiles {
		sourceFile := pair.file
		ref := pair.ref

		if !ref.Preserve {
			continue
		}

		file := tx.host.GetSourceFileFromReference(sourceFile, ref)
		if file == nil {
			continue
		}

		var declFileName string
		if file.IsDeclarationFile {
			declFileName = file.FileName()
		} else {
			// !!! bundled emit support, omit bundled refs
			// if (tx.isBundledEmit && contains((node as Bundle).sourceFiles, file)) continue
			paths := tx.host.GetOutputPathsFor(file, true)
			// Try to use output path for referenced file, or output js path if that doesn't exist, or the input path if all else fails
			declFileName = paths.DeclarationFilePath()
			if len(declFileName) == 0 {
				declFileName = paths.JsFilePath()
			}
			if len(declFileName) == 0 {
				declFileName = file.FileName()
			}
		}
		// Should only be missing if the source file is missing a fileName (at which point we can't name a reference to it anyway)
		// TODO: Shouldn't this be a crash or assert instead of a silent continue?
		if len(declFileName) == 0 {
			continue
		}

		fileName := tspath.GetRelativePathToDirectoryOrUrl(
			outputFilePath,
			declFileName,
			false, // TODO: Probably unsafe to assume this isn't a URL, but that's what strada does
			tspath.ComparePathsOptions{
				CurrentDirectory:          tx.host.GetCurrentDirectory(),
				UseCaseSensitiveFileNames: tx.host.UseCaseSensitiveFileNames(),
			},
		)

		results = append(results, &ast.FileReference{
			TextRange:      core.NewTextRange(-1, -1),
			FileName:       fileName,
			ResolutionMode: ref.ResolutionMode,
			Preserve:       ref.Preserve,
		})
	}
	return results
}

func (tx *DeclarationTransformer) getLibReferences() (result []*ast.FileReference) {
	// clone retained references
	for _, ref := range tx.rawLibReferenceDirectives {
		if !ref.Preserve {
			continue
		}
		result = append(result, &ast.FileReference{
			TextRange:      core.NewTextRange(-1, -1),
			FileName:       ref.FileName,
			ResolutionMode: ref.ResolutionMode,
			Preserve:       ref.Preserve,
		})
	}
	return result
}

func (tx *DeclarationTransformer) getTypeReferences() (result []*ast.FileReference) {
	// clone retained references
	for _, ref := range tx.rawTypeReferenceDirectives {
		if !ref.Preserve {
			continue
		}
		result = append(result, &ast.FileReference{
			TextRange:      core.NewTextRange(-1, -1),
			FileName:       ref.FileName,
			ResolutionMode: ref.ResolutionMode,
			Preserve:       ref.Preserve,
		})
	}
	return result
}

func (tx *DeclarationTransformer) visitDeclarationStatements(input *ast.Node) *ast.Node {
	if !isPreservedDeclarationStatement(input) {
		// return nil for unmatched kinds to omit them from the tree
		return nil
	}
	// !!! TODO: stripInternal support?
	// if (shouldStripInternal(input)) return nil
	switch input.Kind {
	case ast.KindExportDeclaration:
		if ast.IsSourceFile(input.Parent) {
			tx.resultHasExternalModuleIndicator = true
		}
		tx.resultHasScopeMarker = true
		// Rewrite external module names if necessary
		return tx.Factory().UpdateExportDeclaration(
			input.AsExportDeclaration(),
			input.Modifiers(),
			input.IsTypeOnly(),
			input.AsExportDeclaration().ExportClause,
			tx.rewriteModuleSpecifier(input, input.AsExportDeclaration().ModuleSpecifier),
			tx.tryGetResolutionModeOverride(input.AsExportDeclaration().Attributes),
		)
	case ast.KindExportAssignment:
		if ast.IsSourceFile(input.Parent) {
			tx.resultHasExternalModuleIndicator = true
		}
		tx.resultHasScopeMarker = true
		if input.AsExportAssignment().Expression.Kind == ast.KindIdentifier {
			return input
		}
		// expression is non-identifier, create _default typed variable to reference
		newId := tx.EmitContext().NewUniqueName("_default", printer.AutoGenerateOptions{Flags: printer.GeneratedIdentifierFlagsOptimistic})
		tx.getSymbolAccessibilityDiagnostic = func(_ printer.SymbolAccessibilityResult) *SymbolAccessibilityDiagnostic {
			return &SymbolAccessibilityDiagnostic{
				diagnosticMessage: diagnostics.Default_export_of_the_module_has_or_is_using_private_name_0,
				errorNode:         input,
			}
		}
		tx.tracker.PushErrorFallbackNode(input)
		type_ := tx.ensureType(input, false)
		varDecl := tx.Factory().NewVariableDeclaration(newId, nil, type_, nil)
		tx.tracker.PopErrorFallbackNode()
		var modList *ast.ModifierList
		if tx.needsDeclare {
			modList = tx.Factory().NewModifierList([]*ast.Node{tx.Factory().NewModifier(ast.KindDeclareKeyword)})
		} else {
			modList = tx.Factory().NewModifierList([]*ast.Node{})
		}
		statement := tx.Factory().NewVariableStatement(modList, tx.Factory().NewVariableDeclarationList(ast.NodeFlagsConst, tx.Factory().NewNodeList([]*ast.Node{varDecl})))

		assignment := tx.Factory().UpdateExportAssignment(input.AsExportAssignment(), input.Modifiers(), newId)
		// Remove coments from the export declaration and copy them onto the synthetic _default declaration
		tx.preserveJsDoc(statement, input)
		tx.removeAllComments(assignment)
		tx.Factory().NewSyntaxList([]*ast.Node{statement, assignment})
		return nil
	default:
		result := tx.transformTopLevelDeclaration(input)
		// Don't actually transform yet; just leave as original node - will be elided/swapped by late pass
		original := tx.EmitContext().MostOriginal(input)
		id := ast.GetNodeId(original)
		tx.lateStatementReplacementMap[id] = result
		return input
	}
}

func (tx *DeclarationTransformer) rewriteModuleSpecifier(parent *ast.Node, input *ast.Node) *ast.Node {
	if input == nil {
		return nil
	}
	tx.resultHasExternalModuleIndicator = tx.resultHasExternalModuleIndicator || (parent.Kind != ast.KindModuleDeclaration && parent.Kind != ast.KindImportType)
	if ast.IsStringLiteralLike(input) {
		if tx.isBundledEmit {
			// !!! TODO: support bundled emit specifier rewriting
		}
	}
	return input
}

func (tx *DeclarationTransformer) tryGetResolutionModeOverride(node *ast.Node) *ast.Node {
	if node == nil {
		return node
	}
	mode := tx.host.GetResolutionModeOverride(node)
	if mode != core.ResolutionModeNone {
		return node
	}
	return nil
}

func (tx *DeclarationTransformer) preserveJsDoc(updated *ast.Node, original *ast.Node) {
	// !!! TODO: JSDoc comment support
	// if (hasJSDocNodes(updated) && hasJSDocNodes(original)) {
	// 	updated.jsDoc = original.jsDoc;
	// }
	// return setCommentRange(updated, getCommentRange(original));
}

func (tx *DeclarationTransformer) removeAllComments(node *ast.Node) {
	tx.EmitContext().AddEmitFlags(node, printer.EFNoComments)
	// !!! TODO: Also remove synthetic trailing/leading comments added by transforms
	// emitNode.leadingComments = undefined;
	// emitNode.trailingComments = undefined;
}

func (tx *DeclarationTransformer) ensureType(node *ast.Node, ignorePrivate bool) *ast.Node {
	if !ignorePrivate && tx.host.GetEffectiveDeclarationFlags(node, ast.ModifierFlagsPrivate) != 0 {
		// Private nodes emit no types (except private parameter properties, whose parameter types are actually visible)
		return nil
	}

	if tx.shouldPrintWithInitializer(node) {
		// Literal const declarations will have an initializer ensured rather than a type
		return nil
	}

	// Should be removed createTypeOfDeclaration will actually now reuse the existing annotation so there is no real need to duplicate type walking
	// Left in for now to minimize diff during syntactic type node builder refactor
	if !ast.IsExportAssignment(node) && !ast.IsBindingElement(node) && node.Type() != nil && (!ast.IsParameter(node) || !tx.resolver.RequiresAddingImplicitUndefined(node, nil, tx.enclosingDeclaration)) {
		// return visitNode(node.type, visitDeclarationSubtree, isTypeNode); // !!! TODO: visitDeclarationSubtree for ensuring syntactic completeness and vis errors
		return node
	}

	oldErrorNameNode := tx.state.errorNameNode
	tx.state.errorNameNode = node.Name()
	var oldDiag GetSymbolAccessibilityDiagnostic
	if !tx.suppressNewDiagnosticContexts {
		oldDiag = tx.getSymbolAccessibilityDiagnostic
		if canProduceDiagnostics(node) {
			tx.state.getSymbolAccessibilityDiagnostic = createGetSymbolAccessibilityDiagnosticForNode(node)
		}
	}
	var typeNode *ast.Node

	if hasInferredType(node) {
		typeNode = tx.resolver.CreateTypeOfDeclaration(tx.EmitContext(), node, tx.enclosingDeclaration, declarationEmitNodeBuilderFlags, declarationEmitInternalNodeBuilderFlags, tx.tracker)
	} else if ast.IsFunctionLike(node) {
		typeNode = tx.resolver.CreateReturnTypeOfSignatureDeclaration(tx.EmitContext(), node, tx.enclosingDeclaration, declarationEmitNodeBuilderFlags, declarationEmitInternalNodeBuilderFlags, tx.tracker)
	} else {
		// Debug.assertNever(node); // !!!
	}

	tx.state.errorNameNode = oldErrorNameNode
	if !tx.suppressNewDiagnosticContexts {
		tx.state.getSymbolAccessibilityDiagnostic = oldDiag
	}
	if typeNode == nil {
		return tx.Factory().NewKeywordTypeNode(ast.KindAnyKeyword)
	}
	return typeNode
}

func (tx *DeclarationTransformer) shouldPrintWithInitializer(node *ast.Node) bool {
	return canHaveLiteralInitializer(tx.host, node) && node.Initializer() != nil && tx.resolver.IsLiteralConstDeclaration(tx.EmitContext().MostOriginal(node))
}

func (tx *DeclarationTransformer) checkEntityNameVisibility(entityName *ast.Node, enclosingDeclaration *ast.Node) {
	visibilityResult := tx.resolver.IsEntityNameVisible(entityName, enclosingDeclaration)
	tx.tracker.handleSymbolAccessibilityError(visibilityResult)
}

// Transforms the direct child of a source file into zero or more replacement statements
func (tx *DeclarationTransformer) transformTopLevelDeclaration(input *ast.Node) *ast.Node {
	if len(tx.state.lateMarkedStatements) > 0 {
		// Remove duplicates of the current statement from the deferred work queue (this was done via orderedRemoveItem in strada - why? to ensure the same backing array? microop?)
		tx.state.lateMarkedStatements = core.Filter(tx.state.lateMarkedStatements, func(node *ast.Node) bool { return node != input })
	}
	// !!! TODO: stripInternal support?
	// if (shouldStripInternal(input)) return;
	if input.Kind == ast.KindImportEqualsDeclaration {
		return tx.transformImportEqualsDeclaration(input.AsImportEqualsDeclaration())
	}
	if input.Kind == ast.KindImportDeclaration {
		return tx.transformImportDeclaration(input.AsImportDeclaration())
	}
	if ast.IsDeclaration(input) && isDeclarationAndNotVisible(tx.EmitContext(), tx.resolver, input) {
		return nil
	}

	// !!! TODO: JSDoc support
	// if (isJSDocImportTag(input)) return;

	// Elide implementation signatures from overload sets
	if ast.IsFunctionLike(input) && tx.resolver.IsImplementationOfOverload(input) {
		return nil
	}
	previousEnclosingDeclaration := tx.enclosingDeclaration
	if isEnclosingDeclaration(input) {
		tx.enclosingDeclaration = input
	}

	canProdiceDiagnostic := canProduceDiagnostics(input)
	oldDiag := tx.state.getSymbolAccessibilityDiagnostic
	if canProdiceDiagnostic {
		tx.state.getSymbolAccessibilityDiagnostic = createGetSymbolAccessibilityDiagnosticForNode(input)
	}
	previousNeedsDeclare := tx.needsDeclare

	var result *ast.Node
	switch input.Kind {
	// !!!
	case ast.KindVariableStatement:
		result = tx.transformVariableStatement(input.AsVariableStatement())
	case ast.KindEnumDeclaration:
		result = tx.Factory().UpdateEnumDeclaration(
			input.AsEnumDeclaration(),
			tx.Factory().NewModifierList(tx.ensureModifiers(input)),
			input.Name(),
			tx.Factory().NewNodeList(core.MapNonNil(input.AsEnumDeclaration().Members.Nodes, func(m *ast.Node) *ast.Node {
				// !!! TODO: stripInternal support?
				// if (shouldStripInternal(m)) return;

				// !!! TODO: isolatedDeclarations support
				// if (
				// 	isolatedDeclarations && m.initializer && enumValue?.hasExternalReferences &&
				// 	// This will be its own compiler error instead, so don't report.
				// 	!isComputedPropertyName(m.name)
				// ) {
				// 	context.addDiagnostic(createDiagnosticForNode(m, Diagnostics.Enum_member_initializers_must_be_computable_without_references_to_external_symbols_with_isolatedDeclarations));
				// }

				// Rewrite enum values to their constants, if available
				enumValue := tx.resolver.GetEnumMemberValue(m)
				var newInitializer *ast.Node
				switch value := enumValue.Value.(type) {
				case jsnum.Number:
					if value >= 0 {
						newInitializer = tx.Factory().NewNumericLiteral(value.String())
					} else {
						newInitializer = tx.Factory().NewPrefixUnaryExpression(
							ast.KindMinusToken,
							tx.Factory().NewNumericLiteral((-value).String()),
						)
					}
				case string:
					newInitializer = tx.Factory().NewStringLiteral(value)
				default:
					// nil
					newInitializer = nil
				}
				result := tx.Factory().UpdateEnumMember(m.AsEnumMember(), m.Name(), newInitializer)
				tx.preserveJsDoc(result, m)
				return result
			})),
		)
	default:
		// Anything left unhandled is an error, so this should be unreachable
		panic(fmt.Sprintf("Unhandled top-level node in declaration emit: %q", input.Kind))
	}

	if isEnclosingDeclaration(input) {
		tx.enclosingDeclaration = previousEnclosingDeclaration
	}
	if canProdiceDiagnostic {
		tx.state.getSymbolAccessibilityDiagnostic = oldDiag
	}
	if input.Kind == ast.KindModuleDeclaration {
		tx.needsDeclare = previousNeedsDeclare
	}
	if result == input {
		return input
	}
	tx.state.errorNameNode = nil
	return result
}

func (tx *DeclarationTransformer) transformVariableStatement(node *ast.VariableStatement) *ast.Node {
	return nil // !!!
}

func (tx *DeclarationTransformer) ensureModifiers(node *ast.Node) []*ast.Node {
	currentFlags := tx.host.GetEffectiveDeclarationFlags(node, ast.ModifierFlagsAll)
	newFlags := tx.ensureModifierFlags(node)
	if currentFlags == newFlags {
		// Elide decorators
		return core.Filter(node.Modifiers().Nodes, ast.IsModifier)
	}
	return ast.CreateModifiersFromModifierFlags(newFlags, tx.Factory().NewModifier)
}

func (tx *DeclarationTransformer) ensureModifierFlags(node *ast.Node) ast.ModifierFlags {
	mask := ast.ModifierFlagsAll ^ (ast.ModifierFlagsPublic | ast.ModifierFlagsAsync | ast.ModifierFlagsOverride) // No async and override modifiers in declaration files
	additions := ast.ModifierFlagsNone
	if tx.needsDeclare && !isAlwaysType(node) {
		additions = ast.ModifierFlagsAmbient
	}
	parentIsFile := node.Parent.Kind == ast.KindSourceFile
	if !parentIsFile || (tx.isBundledEmit && parentIsFile && ast.IsExternalModule(node.Parent.AsSourceFile())) {
		mask ^= ast.ModifierFlagsAmbient
		additions = ast.ModifierFlagsNone
	}
	return maskModifierFlagsEx(tx.host, node, mask, additions)
}

func (tx *DeclarationTransformer) transformImportEqualsDeclaration(decl *ast.ImportEqualsDeclaration) *ast.Node {
	if !tx.resolver.IsDeclarationVisible(decl.AsNode()) {
		return nil
	}
	if decl.ModuleReference.Kind == ast.KindExternalModuleReference {
		// Rewrite external module names if necessary
		specifier := ast.GetExternalModuleImportEqualsDeclarationExpression(decl.AsNode())
		return tx.Factory().UpdateImportEqualsDeclaration(
			decl,
			decl.Modifiers(),
			decl.IsTypeOnly,
			decl.Name(),
			tx.Factory().UpdateExternalModuleReference(decl.ModuleReference.AsExternalModuleReference(), tx.rewriteModuleSpecifier(decl.AsNode(), specifier)),
		)
	} else {
		oldDiag := tx.getSymbolAccessibilityDiagnostic
		tx.getSymbolAccessibilityDiagnostic = createGetSymbolAccessibilityDiagnosticForNode(decl.AsNode())
		tx.checkEntityNameVisibility(decl.ModuleReference, tx.enclosingDeclaration)
		tx.getSymbolAccessibilityDiagnostic = oldDiag
		return decl.AsNode()
	}
}

func (tx *DeclarationTransformer) transformImportDeclaration(decl *ast.ImportDeclaration) *ast.Node {
	if decl.ImportClause == nil {
		// import "mod" - possibly needed for side effects? (global interface patches, module augmentations, etc)
		return tx.Factory().UpdateImportDeclaration(
			decl,
			decl.Modifiers(),
			decl.ImportClause,
			tx.rewriteModuleSpecifier(decl.AsNode(), decl.ModuleSpecifier),
			tx.tryGetResolutionModeOverride(decl.Attributes),
		)
	}
	// The `importClause` visibility corresponds to the default's visibility.
	var visibleDefaultBinding *ast.Node
	if decl.ImportClause != nil && decl.ImportClause.Name() != nil && tx.resolver.IsDeclarationVisible(decl.ImportClause) {
		visibleDefaultBinding = decl.ImportClause.Name()
	}
	if decl.ImportClause.AsImportClause().NamedBindings == nil {
		// No named bindings (either namespace or list), meaning the import is just default or should be elided
		if visibleDefaultBinding == nil {
			return nil
		}
		return tx.Factory().UpdateImportDeclaration(
			decl,
			decl.Modifiers(),
			tx.Factory().UpdateImportClause(
				decl.ImportClause.AsImportClause(),
				decl.ImportClause.AsImportClause().IsTypeOnly,
				visibleDefaultBinding,
				/*namedBindings*/ nil,
			),
			tx.rewriteModuleSpecifier(decl.AsNode(), decl.ModuleSpecifier),
			tx.tryGetResolutionModeOverride(decl.Attributes),
		)
	}
	if decl.ImportClause.AsImportClause().NamedBindings.Kind == ast.KindNamespaceImport {
		// Namespace import (optionally with visible default)
		var namedBindings *ast.Node
		if tx.resolver.IsDeclarationVisible(decl.ImportClause.AsImportClause().NamedBindings) {
			namedBindings = decl.ImportClause.AsImportClause().NamedBindings
		}
		if visibleDefaultBinding == nil && namedBindings == nil {
			return nil
		}
		return tx.Factory().UpdateImportDeclaration(
			decl,
			decl.Modifiers(),
			tx.Factory().UpdateImportClause(
				decl.ImportClause.AsImportClause(),
				decl.ImportClause.AsImportClause().IsTypeOnly,
				visibleDefaultBinding,
				namedBindings,
			),
			tx.rewriteModuleSpecifier(decl.AsNode(), decl.ModuleSpecifier),
			tx.tryGetResolutionModeOverride(decl.Attributes),
		)
	}
	// Named imports (optionally with visible default)
	bindingList := core.Filter(
		decl.ImportClause.AsImportClause().NamedBindings.AsNamedImports().Elements.Nodes,
		func(b *ast.Node) bool {
			return tx.resolver.IsDeclarationVisible(b)
		},
	)
	if len(bindingList) > 0 || visibleDefaultBinding != nil {
		var namedImports *ast.Node
		if len(bindingList) > 0 {
			namedImports = tx.Factory().UpdateNamedImports(
				decl.ImportClause.AsImportClause().NamedBindings.AsNamedImports(),
				tx.Factory().NewNodeList(bindingList),
			)
		}
		return tx.Factory().UpdateImportDeclaration(
			decl,
			decl.Modifiers(),
			tx.Factory().UpdateImportClause(
				decl.ImportClause.AsImportClause(),
				decl.ImportClause.AsImportClause().IsTypeOnly,
				visibleDefaultBinding,
				namedImports,
			),
			tx.rewriteModuleSpecifier(decl.AsNode(), decl.ModuleSpecifier),
			tx.tryGetResolutionModeOverride(decl.Attributes),
		)
	}
	// Augmentation of export depends on import
	if tx.resolver.IsImportRequiredByAugmentation(decl) {
		// IsolatedDeclarations support
		// if (isolatedDeclarations) {
		// 	context.addDiagnostic(createDiagnosticForNode(decl, Diagnostics.Declaration_emit_for_this_file_requires_preserving_this_import_for_augmentations_This_is_not_supported_with_isolatedDeclarations));
		// }
		return tx.Factory().UpdateImportDeclaration(
			decl,
			decl.Modifiers(),
			/*importClause*/ nil,
			tx.rewriteModuleSpecifier(decl.AsNode(), decl.ModuleSpecifier),
			tx.tryGetResolutionModeOverride(decl.Attributes),
		)
	}
	// Nothing visible
	return nil
}
