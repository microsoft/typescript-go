package ls

import (
	"context"
	"slices"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/locale"
	"github.com/microsoft/typescript-go/internal/ls/change"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/scanner"
)

const (
	fixNameUnusedIdentifier            = "unusedIdentifier"
	fixIdUnusedIdentifierPrefix        = "unusedIdentifier_prefix"
	fixIdUnusedIdentifierDelete        = "unusedIdentifier_delete"
	fixIdUnusedIdentifierDeleteImports = "unusedIdentifier_deleteImports"
	fixIdUnusedIdentifierInfer         = "unusedIdentifier_infer"
)

var unusedIdentifierErrorCodes = []int32{
	diagnostics.X_0_is_declared_but_its_value_is_never_read.Code(),
	diagnostics.X_0_is_declared_but_never_used.Code(),
	diagnostics.Property_0_is_declared_but_its_value_is_never_read.Code(),
	diagnostics.All_imports_in_import_declaration_are_unused.Code(),
	diagnostics.All_destructured_elements_are_unused.Code(),
	diagnostics.All_variables_are_unused.Code(),
	diagnostics.All_type_parameters_are_unused.Code(),
}

var UnusedIdentifierFixProvider = &CodeFixProvider{
	ErrorCodes:     unusedIdentifierErrorCodes,
	GetCodeActions: getUnusedIdentifierCodeActions,
	FixIds: []string{
		fixIdUnusedIdentifierPrefix,
		fixIdUnusedIdentifierDelete,
		fixIdUnusedIdentifierDeleteImports,
		fixIdUnusedIdentifierInfer,
	},
	GetAllCodeActions: getAllUnusedIdentifierCodeActions,
}

type unusedIdentifierFixKind int

const (
	unusedIdentifierFixKindDeleteNode unusedIdentifierFixKind = iota
	unusedIdentifierFixKindDeleteTypeParameters
	unusedIdentifierFixKindDeleteImportSpecifier
	unusedIdentifierFixKindDeleteDestructuring
	unusedIdentifierFixKindDeleteVariableStatement
	unusedIdentifierFixKindDeleteFunctionDeclaration
	unusedIdentifierFixKindDeleteDeclaration
	unusedIdentifierFixKindReplaceInferWithUnknown
	unusedIdentifierFixKindPrefixDeclaration
)

type unusedIdentifierFixInfo struct {
	kind        unusedIdentifierFixKind
	fixID       string
	description string
	node        *ast.Node
}

type unusedIdentifierFixer struct {
	ctx         context.Context
	ls          *LanguageService
	sourceFile  *ast.SourceFile
	program     *compiler.Program
	typeChecker *checker.Checker
	locale      locale.Locale
	fixAll      bool
}

func getUnusedIdentifierCodeActions(ctx context.Context, fixContext *CodeFixContext) ([]*CodeAction, error) {
	program := fixContext.Program
	typeChecker, done := program.GetTypeChecker(ctx)
	defer done()

	token := astnav.GetTokenAtPosition(fixContext.SourceFile, fixContext.Span.Pos())
	if token == nil {
		return nil, nil
	}

	fixer := &unusedIdentifierFixer{
		ctx:         ctx,
		ls:          fixContext.LS,
		sourceFile:  fixContext.SourceFile,
		program:     program,
		typeChecker: typeChecker,
		locale:      locale.FromContext(ctx),
	}

	var actions []*CodeAction
	for _, info := range fixer.getInfo(fixContext.ErrorCode, token) {
		changeTracker := fixer.createChangeTracker()
		fixer.writeChanges(changeTracker, info)

		edits := changeTracker.GetChanges()[fixer.sourceFile.FileName()]
		if len(edits) == 0 {
			continue
		}

		actions = append(actions, fixer.createCodeAction(info.description, info.fixID, edits))
	}
	return actions, nil
}

func getAllUnusedIdentifierCodeActions(ctx context.Context, fixContext *CodeFixContext) (*CombinedCodeActions, error) {
	program := fixContext.Program
	typeChecker, done := program.GetTypeChecker(ctx)
	defer done()

	fixer := &unusedIdentifierFixer{
		ctx:         ctx,
		ls:          fixContext.LS,
		sourceFile:  fixContext.SourceFile,
		program:     program,
		typeChecker: typeChecker,
		locale:      locale.FromContext(ctx),
		fixAll:      true,
	}

	changeTracker := fixer.createChangeTracker()
	for _, diag := range getAllDiagnostics(ctx, program, fixContext.SourceFile) {
		if containsErrorCode(unusedIdentifierErrorCodes, diag.Code()) {
			token := astnav.GetTokenAtPosition(fixContext.SourceFile, diag.Loc().Pos())
			if token == nil {
				continue
			}
			for _, info := range fixer.getInfo(diag.Code(), token) {
				if info.fixID == fixContext.FixID {
					fixer.writeChanges(changeTracker, info)
				}
			}
		}
	}

	changes := changeTracker.GetChanges()[fixContext.SourceFile.FileName()]
	if len(changes) == 0 {
		return nil, nil
	}

	return &CombinedCodeActions{
		Description: getFixAllDescriptionMessage(fixContext.FixID).Localize(locale.FromContext(ctx)),
		Changes:     changes,
	}, nil
}

func getFixAllDescriptionMessage(fixID string) *diagnostics.Message {
	switch fixID {
	case fixIdUnusedIdentifierDeleteImports:
		return diagnostics.Delete_all_unused_imports
	case fixIdUnusedIdentifierInfer:
		return diagnostics.Replace_all_unused_infer_with_unknown
	case fixIdUnusedIdentifierPrefix:
		return diagnostics.Prefix_all_unused_declarations_with_where_possible
	default:
		return diagnostics.Delete_all_unused_declarations
	}
}

func (f *unusedIdentifierFixer) getInfo(errorCode int32, token *ast.Node) []unusedIdentifierFixInfo {
	if ast.IsJSDocTemplateTag(token) {
		return []unusedIdentifierFixInfo{
			{
				kind:        unusedIdentifierFixKindDeleteNode,
				fixID:       fixIdUnusedIdentifierDelete,
				description: diagnostics.Remove_template_tag.Localize(f.locale),
				node:        token,
			},
		}
	}

	if token.Kind == ast.KindLessThanToken {
		return []unusedIdentifierFixInfo{
			{
				kind:        unusedIdentifierFixKindDeleteTypeParameters,
				fixID:       fixIdUnusedIdentifierDelete,
				description: diagnostics.Remove_type_parameters.Localize(f.locale),
				node:        token,
			},
		}
	}

	if importDecl := tryGetImportDeclaration(token); importDecl != nil {
		return []unusedIdentifierFixInfo{
			{
				kind:        unusedIdentifierFixKindDeleteNode,
				fixID:       fixIdUnusedIdentifierDeleteImports,
				description: diagnostics.Remove_import_from_0.Localize(f.locale, ast.GetExternalModuleName(importDecl.AsNode()).Text()),
				node:        importDecl.AsNode(),
			},
		}
	}

	if isImport(token) {
		return []unusedIdentifierFixInfo{
			{
				kind:        unusedIdentifierFixKindDeleteImportSpecifier,
				fixID:       fixIdUnusedIdentifierDeleteImports,
				description: diagnostics.Remove_unused_declaration_for_Colon_0.Localize(f.locale, token.Text()),
				node:        token,
			},
		}
	}

	if ast.IsObjectBindingPattern(token.Parent) || ast.IsArrayBindingPattern(token.Parent) {
		bindingPattern := token.Parent
		parent := bindingPattern.Parent
		if f.fixAll {
			var node *ast.Node
			switch {
			case ast.IsVariableDeclaration(parent):
				if parent.AsVariableDeclaration().Initializer != nil {
					return nil
				}
				node = parent
			case ast.IsParameterDeclaration(parent):
				if f.isNotProvidedArguments(parent.AsParameterDeclaration()) {
					node = core.IfElse(ast.IsArrayBindingPattern(bindingPattern), bindingPattern, parent)
				} else {
					return nil
				}
			default:
				return nil
			}
			return []unusedIdentifierFixInfo{{kind: unusedIdentifierFixKindDeleteNode, fixID: fixIdUnusedIdentifierDelete, node: node}}
		}

		var description string
		if ast.IsParameterDeclaration(parent) {
			elements := bindingPattern.AsBindingPattern().Elements.Nodes
			msg := core.IfElse(len(elements) > 1, diagnostics.Remove_unused_declarations_for_Colon_0, diagnostics.Remove_unused_declaration_for_Colon_0)
			names := make([]string, len(elements))
			for i, e := range elements {
				names[i] = e.Name().Text()
			}
			description = msg.Localize(f.locale, strings.Join(names, ", "))
		} else {
			description = diagnostics.Remove_unused_destructuring_declaration.Localize(f.locale)
		}

		return []unusedIdentifierFixInfo{
			{
				kind:        unusedIdentifierFixKindDeleteDestructuring,
				fixID:       fixIdUnusedIdentifierDelete,
				description: description,
				node:        bindingPattern,
			},
		}
	}

	if canDeleteVariableStatement(f.sourceFile, token) {
		return []unusedIdentifierFixInfo{
			{
				kind:        unusedIdentifierFixKindDeleteVariableStatement,
				fixID:       fixIdUnusedIdentifierDelete,
				description: diagnostics.Remove_variable_statement.Localize(f.locale),
				node:        token.Parent,
			},
		}
	}

	if ast.IsIdentifier(token) && ast.IsFunctionDeclaration(token.Parent) {
		return []unusedIdentifierFixInfo{
			{
				kind:        unusedIdentifierFixKindDeleteFunctionDeclaration,
				fixID:       fixIdUnusedIdentifierDelete,
				description: diagnostics.Remove_unused_declaration_for_Colon_0.Localize(f.locale, token.Text()),
				node:        token.Parent,
			},
		}
	}

	var infos []unusedIdentifierFixInfo
	if token.Kind == ast.KindInferKeyword {
		name := token.Parent.AsInferTypeNode().TypeParameter.AsTypeParameterDeclaration().Name().Text()
		infos = append(infos, unusedIdentifierFixInfo{
			kind:        unusedIdentifierFixKindReplaceInferWithUnknown,
			fixID:       fixIdUnusedIdentifierInfer,
			description: diagnostics.Replace_infer_0_with_unknown.Localize(f.locale, name),
			node:        token,
		})
	} else if token.Parent != nil {
		var name string
		if ast.IsComputedPropertyName(token.Parent) {
			name = scanner.GetSourceTextOfNodeFromSourceFile(f.sourceFile, token.Parent, false /*includeTrivia*/)
		} else if ast.IsPropertyNameLiteral(token) {
			name = scanner.GetSourceTextOfNodeFromSourceFile(f.sourceFile, token, false /*includeTrivia*/)
		} else {
			name = token.Text()
		}
		infos = append(infos, unusedIdentifierFixInfo{
			kind:        unusedIdentifierFixKindDeleteDeclaration,
			fixID:       fixIdUnusedIdentifierDelete,
			description: diagnostics.Remove_unused_declaration_for_Colon_0.Localize(f.locale, name),
			node:        token,
		})
	}

	if errorCode != diagnostics.Property_0_is_declared_but_its_value_is_never_read.Code() {
		prefixToken := token
		if prefixToken.Kind == ast.KindInferKeyword && ast.IsInferTypeNode(prefixToken.Parent) {
			prefixToken = prefixToken.Parent.AsInferTypeNode().TypeParameter.AsTypeParameterDeclaration().Name()
		}
		if ast.IsIdentifier(prefixToken) && canPrefix(prefixToken) {
			infos = append(infos, unusedIdentifierFixInfo{
				kind:        unusedIdentifierFixKindPrefixDeclaration,
				fixID:       fixIdUnusedIdentifierPrefix,
				description: diagnostics.Prefix_0_with_an_underscore.Localize(f.locale, prefixToken.Text()),
				node:        prefixToken,
			})
		}
	}

	return infos
}

func (f *unusedIdentifierFixer) writeChanges(changeTracker *change.Tracker, info unusedIdentifierFixInfo) {
	switch info.kind {
	case unusedIdentifierFixKindDeleteNode:
		changeTracker.Delete(f.sourceFile, info.node)
	case unusedIdentifierFixKindDeleteTypeParameters:
		f.deleteTypeParameters(changeTracker, info.node)
	case unusedIdentifierFixKindDeleteImportSpecifier:
		f.deleteImportSpecifier(changeTracker, info.node)
	case unusedIdentifierFixKindDeleteDestructuring:
		f.deleteDestructuring(changeTracker, info.node)
	case unusedIdentifierFixKindDeleteVariableStatement:
		f.deleteVariableStatement(changeTracker, info.node)
	case unusedIdentifierFixKindDeleteFunctionDeclaration:
		if f.fixAll {
			changeTracker.Delete(f.sourceFile, info.node)
		} else {
			for _, decl := range info.node.AsFunctionDeclaration().Symbol.Declarations {
				changeTracker.Delete(f.sourceFile, decl)
			}
		}
	case unusedIdentifierFixKindDeleteDeclaration:
		f.deleteDeclaration(changeTracker, info.node)
	case unusedIdentifierFixKindReplaceInferWithUnknown:
		f.replaceInferWithUnknown(changeTracker, info.node)
	case unusedIdentifierFixKindPrefixDeclaration:
		f.prefixDeclaration(changeTracker, info.node)
	}
}

func (f *unusedIdentifierFixer) deleteTypeParameters(changeTracker *change.Tracker, node *ast.Node) {
	typeParams := node.Parent.TypeParameters()
	if len(typeParams) == 0 {
		return
	}
	first := core.FirstOrNil(typeParams)
	last := core.LastOrNil(typeParams)
	openAngle := astnav.FindPrecedingToken(f.sourceFile, first.Pos())
	closeAngle := astnav.GetTokenAtPosition(f.sourceFile, last.End())
	if openAngle == nil || closeAngle == nil {
		return
	}
	if openAngle.Kind == ast.KindLessThanToken && closeAngle.Kind == ast.KindGreaterThanToken {
		changeTracker.DeleteNodeRange(f.sourceFile, openAngle, closeAngle, change.LeadingTriviaOptionExclude, change.TrailingTriviaOptionExclude)
	}
}

func (f *unusedIdentifierFixer) deleteImportSpecifier(changeTracker *change.Tracker, token *ast.Node) {
	node := token.Parent
	if ast.IsImportClause(token.Parent) {
		node = token
	}
	changeTracker.Delete(f.sourceFile, node)
}

func (f *unusedIdentifierFixer) deleteDestructuring(changeTracker *change.Tracker, node *ast.Node) {
	parent := node.Parent
	if ast.IsParameterDeclaration(parent) {
		for _, elem := range node.AsBindingPattern().Elements.Nodes {
			changeTracker.Delete(f.sourceFile, elem)
		}
		return
	}
	if ast.IsVariableDeclaration(parent) && parent.AsVariableDeclaration().Initializer != nil && ast.IsCallLikeExpression(parent.AsVariableDeclaration().Initializer) {
		if ast.IsVariableDeclarationList(parent.Parent) && len(parent.Parent.AsVariableDeclarationList().Declarations.Nodes) > 1 {
			variableStatement := parent.Parent.AsVariableDeclarationList().Parent
			pos := astnav.GetStartOfNode(variableStatement, f.sourceFile, false /*includeJSDoc*/)
			end := variableStatement.End()
			options := change.NodeOptions{
				Prefix: f.program.Options().NewLine.GetNewLineCharacter() + f.sourceFile.Text()[getPrecedingNonSpaceCharacterPosition(f.sourceFile.Text(), pos-1):pos],
				Suffix: core.IfElse(lsutil.ProbablyUsesSemicolons(f.sourceFile), ";", ""),
			}
			changeTracker.Delete(f.sourceFile, parent)
			changeTracker.InsertNodeAt(f.sourceFile, core.TextPos(end), parent.AsVariableDeclaration().Initializer, options)
		} else {
			changeTracker.ReplaceNode(f.sourceFile, parent.Parent, parent.AsVariableDeclaration().Initializer, nil /*options*/)
		}
	} else {
		changeTracker.Delete(f.sourceFile, parent)
	}
}

func (f *unusedIdentifierFixer) deleteVariableStatement(changeTracker *change.Tracker, node *ast.Node) {
	if ast.IsVariableStatement(node.Parent) {
		node = node.Parent
	}
	changeTracker.Delete(f.sourceFile, node)
}

func (f *unusedIdentifierFixer) replaceInferWithUnknown(changeTracker *change.Tracker, token *ast.Node) {
	changeTracker.ReplaceNode(f.sourceFile, token.Parent, changeTracker.NewKeywordTypeNode(ast.KindUnknownKeyword), nil /*options*/)
}

func (f *unusedIdentifierFixer) prefixDeclaration(changeTracker *change.Tracker, token *ast.Node) {
	changeTracker.ReplaceNode(f.sourceFile, token, changeTracker.NewIdentifier("_"+token.Text()), nil /*options*/)
	if ast.IsParameterDeclaration(token.Parent) {
		fn := token.Parent.Parent
		for _, tag := range getAllJSDocTags(fn) {
			if ast.IsJSDocParameterTag(tag) && ast.IsIdentifier(tag.Name()) && tag.Name().Text() == token.Text() {
				changeTracker.ReplaceNode(f.sourceFile, tag.Name(), changeTracker.NewIdentifier("_"+tag.Name().Text()), nil /*options*/)
			}
		}
	}
}

func (f *unusedIdentifierFixer) deleteDeclaration(changeTracker *change.Tracker, token *ast.Node) {
	parent := token.Parent
	if parent == nil {
		return
	}
	if ast.IsParameterDeclaration(parent) {
		param := parent.AsParameterDeclaration()
		if f.mayDeleteParameter(param) {
			modifiers := param.Modifiers()
			if modifiers != nil && len(modifiers.Nodes) > 0 {
				name := param.Name()
				referencedParameterName := ast.IsIdentifier(name) && isSymbolReferencedInFile(name.AsIdentifier(), f.typeChecker, f.sourceFile, f.sourceFile.AsNode())
				if ast.IsBindingPattern(name) || referencedParameterName {
					for _, modifier := range modifiers.Nodes {
						if ast.IsModifier(modifier) {
							changeTracker.DeleteModifier(f.sourceFile, modifier)
						}
					}
				} else if param.Initializer == nil && f.isNotProvidedArguments(param) {
					changeTracker.Delete(f.sourceFile, param.AsNode())
				}
			} else if param.Initializer == nil && f.isNotProvidedArguments(param) {
				changeTracker.Delete(f.sourceFile, param.AsNode())
			}
		}
	} else {
		if f.fixAll && ast.IsIdentifier(token) && isSymbolReferencedInFile(token.AsIdentifier(), f.typeChecker, f.sourceFile, f.sourceFile.AsNode()) {
			return
		}
		node := parent
		if ast.IsImportClause(parent) {
			node = token
		} else if ast.IsComputedPropertyName(parent) {
			node = parent.Parent
		}

		arrayBindingElement := f.fixAll && ast.IsBindingElement(parent) && ast.IsArrayBindingPattern(parent.Parent) && ast.IsVariableDeclaration(parent.Parent.Parent)
		if arrayBindingElement || node == nil || node == f.sourceFile.AsNode() {
			return
		}

		changeTracker.Delete(f.sourceFile, node)
	}

	if f.fixAll {
		return
	}
	if ast.IsIdentifier(token) {
		for _, ref := range f.findPossibleReferencesForNode(token) {
			if ast.IsPropertyAccessExpression(ref.Parent) && ref.Parent.Name() == ref {
				ref = ref.Parent
			}
			if mayDeleteExpression(ref) {
				changeTracker.DeleteNode(f.sourceFile, ref.Parent.Parent, change.LeadingTriviaOptionExclude, change.TrailingTriviaOptionInclude)
			}
		}
	}
}

func (f *unusedIdentifierFixer) mayDeleteParameter(param *ast.ParameterDeclaration) bool {
	parent := param.Parent
	if parent == nil {
		return false
	}
	switch parent.Kind {
	case ast.KindMethodDeclaration, ast.KindConstructor:
		index := slices.Index(parent.ParameterList().Nodes, param.AsNode())
		if index < 0 {
			return true
		}
		referent := parent
		if ast.IsMethodDeclaration(referent) {
			referent = parent.Name()
		}
		for _, node := range f.findReferencesForNode(referent) {
			if node.Kind == ast.KindSuperKeyword && ast.IsCallExpression(node.Parent) && len(node.Parent.AsCallExpression().Arguments.Nodes) > index {
				return false
			}
			if ast.IsPropertyAccessExpression(node.Parent) && node.Parent.AsPropertyAccessExpression().Expression.Kind == ast.KindSuperKeyword &&
				ast.IsCallExpression(node.Parent.Parent) && len(node.Parent.Parent.AsCallExpression().Arguments.Nodes) > index {
				return false
			}
			if (ast.IsMethodDeclaration(node.Parent) || ast.IsMethodSignatureDeclaration(node.Parent)) && node.Parent != param.Parent && len(node.Parent.Parameters()) > index {
				return false
			}
		}
		return true
	case ast.KindFunctionDeclaration:
		if parent.Name() != nil && ast.IsIdentifier(parent.Name()) && f.isCallbackLike(parent.Name()) {
			return f.isLastParameter(parent, param)
		}
		return true
	case ast.KindFunctionExpression, ast.KindArrowFunction:
		return f.isLastParameter(parent, param)
	case ast.KindSetAccessor:
		return false
	case ast.KindGetAccessor:
		return true
	default:
		return false
	}
}

func mayDeleteExpression(node *ast.Node) bool {
	parent := node.Parent
	if parent == nil {
		return false
	}

	isWriteExpression := ast.IsBinaryExpression(parent) && parent.AsBinaryExpression().Left == node ||
		(ast.IsPostfixUnaryExpression(parent) || ast.IsPrefixUnaryExpression(parent)) && parent.Expression() == node
	return isWriteExpression && ast.IsExpressionStatement(parent.Parent)
}

func (f *unusedIdentifierFixer) isLastParameter(function *ast.Node, param *ast.ParameterDeclaration) bool {
	params := function.ParameterList().Nodes
	index := slices.Index(params, param.AsNode())
	if index < 0 {
		return false
	}
	if f.fixAll {
		for _, p := range params[index+1:] {
			sym := p.Symbol()
			if sym == nil || f.typeChecker.IsReferenced(sym) {
				return false
			}
		}
		return true
	}
	return index == len(params)-1
}

func (f *unusedIdentifierFixer) isCallbackLike(node *ast.Node) bool {
	for _, candidate := range f.findPossibleReferencesForNode(node) {
		called := ast.ClimbPastPropertyAccess(candidate)
		call := called.Parent
		if call == nil {
			continue
		}
		if ast.IsCallExpression(call) && call.AsCallExpression().Expression == called && slices.Contains(call.AsCallExpression().Arguments.Nodes, candidate) {
			return true
		}
	}
	return false
}

func (f *unusedIdentifierFixer) isNotProvidedArguments(node *ast.ParameterDeclaration) bool {
	index := slices.Index(node.Parent.ParameterList().Nodes, node.AsNode())
	if index < 0 {
		return false
	}
	name := node.Parent.Name()
	if name == nil || !ast.IsIdentifier(name) {
		return true
	}
	for _, ref := range f.findPossibleReferencesForNode(name) {
		called := ast.ClimbPastPropertyAccess(ref)
		call := called.Parent
		if call == nil {
			return false
		}
		if ast.IsCallExpression(call) && call.AsCallExpression().Expression == called {
			if len(call.AsCallExpression().Arguments.Nodes) > index {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func (f *unusedIdentifierFixer) findPossibleReferencesForNode(node *ast.Node) []*ast.Node {
	var refs []*ast.Node
	symbol := f.typeChecker.GetSymbolAtLocation(node)
	if symbol == nil {
		return refs
	}
	for _, candidate := range getPossibleSymbolReferenceNodes(f.sourceFile, symbol.Name, nil /*container*/) {
		if ast.IsIdentifier(candidate) {
			if candidate == node {
				continue
			}
			if candidate.Text() == node.Text() {
				ref := f.typeChecker.GetSymbolAtLocation(candidate)
				if ref == nil {
					continue
				}
				if slices.Contains(f.typeChecker.GetRootSymbols(ref), symbol) {
					refs = append(refs, candidate)
				}
			}
		}
	}
	return refs
}

func (f *unusedIdentifierFixer) findReferencesForNode(node *ast.Node) []*ast.Node {
	var refs []*ast.Node
	options := refOptions{
		use: referenceUseReferences,
	}
	for _, entry := range f.ls.getReferencedSymbolsForNode(f.ctx, node.Pos(), node, f.program, f.program.GetSourceFiles(), options) {
		for _, ref := range entry.references {
			if ref.node == nil {
				continue
			}
			refs = append(refs, ref.node)
		}
	}
	return refs
}

func (f *unusedIdentifierFixer) createChangeTracker() *change.Tracker {
	return change.NewTracker(f.ctx, f.program.Options(), f.ls.FormatOptions(), f.ls.converters)
}

func (f *unusedIdentifierFixer) createCodeAction(description string, fixID string, changes []*lsproto.TextEdit) *CodeAction {
	return &CodeAction{
		Description:       description,
		Changes:           changes,
		FixName:           fixNameUnusedIdentifier,
		FixID:             fixID,
		FixAllDescription: getFixAllDescriptionMessage(fixID).Localize(f.locale),
	}
}

func canDeleteVariableStatement(sourceFile *ast.SourceFile, token *ast.Node) bool {
	return ast.IsVariableDeclarationList(token.Parent) && astnav.GetTokenAtPosition(sourceFile, token.Parent.Pos()) == token
}

func tryGetImportDeclaration(token *ast.Node) *ast.ImportDeclaration {
	if token.Kind == ast.KindImportKeyword && token.Parent != nil && ast.IsImportDeclaration(token.Parent) {
		return token.Parent.AsImportDeclaration()
	}
	return nil
}

func isImport(token *ast.Node) bool {
	return token.Kind == ast.KindImportKeyword || ast.IsIdentifier(token) && (ast.IsImportSpecifier(token.Parent) || ast.IsImportClause(token.Parent))
}

func canPrefix(token *ast.Node) bool {
	parent := token.Parent
	if parent == nil {
		return false
	}
	switch parent.Kind {
	case ast.KindParameter, ast.KindTypeParameter:
		return true
	case ast.KindVariableDeclaration:
		if parent.Parent == nil || parent.Parent.Parent == nil {
			return false
		}
		if ast.IsForOfStatement(parent.Parent.Parent) || ast.IsForInStatement(parent.Parent.Parent) {
			return true
		}
	}
	return false
}

func getAllJSDocTags(node *ast.Node) []*ast.Node {
	if node.Flags&ast.NodeFlagsJSDoc == 0 {
		for current := node; current != nil; current = ast.GetNextJSDocCommentLocation(current) {
			jsdocs := current.JSDoc(nil)
			if len(jsdocs) == 0 {
				continue
			}
			lastJSDoc := jsdocs[len(jsdocs)-1].AsJSDoc()
			if lastJSDoc.Tags != nil {
				return lastJSDoc.Tags.Nodes
			}
		}
	}
	return nil
}
