package estransforms

import (
	"slices"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/printer"
	"github.com/microsoft/typescript-go/internal/transformers"
)

type contextFlags int

const (
	contextFlagsNone        contextFlags = 0
	contextFlagsNonTopLevel contextFlags = 1 << iota
	contextFlagsHasLexicalThis
)

type asyncTransformer struct {
	transformers.Transformer
	compilerOptions *core.CompilerOptions
	emitResolver    printer.EmitResolver

	contextFlags contextFlags

	enclosingFunctionParameterNames *collections.Set[string]
	capturedSuperProperties         *collections.Set[string]
	hasSuperElementAccess           bool
	lexicalArgumentsBinding         *ast.IdentifierNode
	usedLexicalArguments            bool

	asyncBodyVisitor *ast.NodeVisitor
}

func (tx *asyncTransformer) visit(node *ast.Node) *ast.Node {
	if node.Kind == ast.KindAsyncKeyword {
		// ES2017 async modifier should be elided for targets < ES2017
		return nil
	}
	if node.SubtreeFacts()&(ast.SubtreeContainsAnyAwait|ast.SubtreeContainsAwait) == 0 {
		if tx.lexicalArgumentsBinding != nil {
			return tx.argumentsVisitor(node)
		}
		return node
	}
	switch node.Kind {
	case ast.KindSourceFile:
		return tx.visitSourceFile(node.AsSourceFile())
	case ast.KindAwaitExpression:
		return tx.visitAwaitExpression(node.AsAwaitExpression())
	case ast.KindMethodDeclaration:
		return tx.doWithContext(contextFlagsNonTopLevel|contextFlagsHasLexicalThis, tx.visitMethodDeclaration, node)
	case ast.KindFunctionDeclaration:
		return tx.doWithContext(contextFlagsNonTopLevel|contextFlagsHasLexicalThis, tx.visitFunctionDeclaration, node)
	case ast.KindFunctionExpression:
		return tx.doWithContext(contextFlagsNonTopLevel|contextFlagsHasLexicalThis, tx.visitFunctionExpression, node)
	case ast.KindArrowFunction:
		return tx.doWithContext(contextFlagsNonTopLevel, tx.visitArrowFunction, node)
	case ast.KindPropertyAccessExpression:
		if tx.capturedSuperProperties != nil && node.Expression().Kind == ast.KindSuperKeyword {
			tx.capturedSuperProperties.Add(node.Name().Text())
		}
		return tx.Visitor().VisitEachChild(node)
	case ast.KindElementAccessExpression:
		if tx.capturedSuperProperties != nil && node.Expression().Kind == ast.KindSuperKeyword {
			tx.hasSuperElementAccess = true
		}
		return tx.Visitor().VisitEachChild(node)
	case ast.KindGetAccessor:
		return tx.doWithContext(contextFlagsNonTopLevel|contextFlagsHasLexicalThis, tx.visitGetAccessorDeclaration, node)
	case ast.KindSetAccessor:
		return tx.doWithContext(contextFlagsNonTopLevel|contextFlagsHasLexicalThis, tx.visitSetAccessorDeclaration, node)
	case ast.KindConstructor:
		return tx.doWithContext(contextFlagsNonTopLevel|contextFlagsHasLexicalThis, tx.visitConstructorDeclaration, node)
	case ast.KindClassDeclaration, ast.KindClassExpression:
		return tx.doWithContext(contextFlagsNonTopLevel|contextFlagsHasLexicalThis, tx.visitDefault, node)
	default:
		return tx.Visitor().VisitEachChild(node)
	}
}

func (tx *asyncTransformer) visitSourceFile(node *ast.SourceFile) *ast.Node {
	if node.IsDeclarationFile {
		return node.AsNode()
	}

	tx.setContextFlag(contextFlagsNonTopLevel, false)
	tx.setContextFlag(contextFlagsHasLexicalThis, !isEffectiveStrictModeSourceFile(node, tx.compilerOptions))
	visited := tx.Visitor().VisitEachChild(node.AsNode())
	tx.EmitContext().AddEmitHelper(visited, tx.EmitContext().ReadEmitHelpers()...)
	return visited
}

func (tx *asyncTransformer) argumentsVisitor(node *ast.Node) *ast.Node {
	switch node.Kind {
	case ast.KindFunctionExpression,
		ast.KindFunctionDeclaration,
		ast.KindMethodDeclaration,
		ast.KindGetAccessor,
		ast.KindSetAccessor,
		ast.KindConstructor:
		return node
	case ast.KindParameter,
		ast.KindBindingElement,
		ast.KindVariableDeclaration:
		// fall through to visitEachChild
	case ast.KindIdentifier:
		if tx.lexicalArgumentsBinding != nil && tx.emitResolver != nil && tx.emitResolver.IsArgumentsLocalBinding(node) {
			tx.usedLexicalArguments = true
			return tx.lexicalArgumentsBinding
		}
	}
	return tx.Visitor().VisitEachChild(node)
}

func (tx *asyncTransformer) visitAsyncBodyNode(node *ast.Node) *ast.Node {
	if isNodeWithPossibleHoistedDeclaration(node) {
		switch node.Kind {
		case ast.KindVariableStatement:
			return tx.visitVariableStatementInAsyncBody(node)
		case ast.KindForStatement:
			return tx.visitForStatementInAsyncBody(node.AsForStatement())
		case ast.KindForInStatement:
			return tx.visitForInStatementInAsyncBody(node.AsForInOrOfStatement())
		case ast.KindForOfStatement:
			return tx.visitForOfStatementInAsyncBody(node.AsForInOrOfStatement())
		case ast.KindCatchClause:
			return tx.visitCatchClauseInAsyncBody(node.AsCatchClause())
		case ast.KindBlock,
			ast.KindSwitchStatement,
			ast.KindCaseBlock,
			ast.KindCaseClause,
			ast.KindDefaultClause,
			ast.KindTryStatement,
			ast.KindDoStatement,
			ast.KindWhileStatement,
			ast.KindIfStatement,
			ast.KindWithStatement,
			ast.KindLabeledStatement:
			return tx.asyncBodyVisitor.VisitEachChild(node)
		}
	}
	return tx.visit(node)
}

func (tx *asyncTransformer) setContextFlag(flag contextFlags, val bool) {
	if val {
		tx.contextFlags |= flag
	} else {
		tx.contextFlags &= ^flag
	}
}

func (tx *asyncTransformer) inContext(flags contextFlags) bool {
	return tx.contextFlags&flags != 0
}

func (tx *asyncTransformer) inTopLevelContext() bool {
	return !tx.inContext(contextFlagsNonTopLevel)
}

func (tx *asyncTransformer) inHasLexicalThisContext() bool {
	return tx.inContext(contextFlagsHasLexicalThis)
}

func (tx *asyncTransformer) doWithContext(flags contextFlags, cb func(*ast.Node) *ast.Node, node *ast.Node) *ast.Node {
	contextFlagsToSet := flags & ^tx.contextFlags
	if contextFlagsToSet != 0 {
		tx.setContextFlag(contextFlagsToSet, true)
		result := cb(node)
		tx.setContextFlag(contextFlagsToSet, false)
		return result
	}
	return cb(node)
}

func (tx *asyncTransformer) visitDefault(node *ast.Node) *ast.Node {
	return tx.Visitor().VisitEachChild(node)
}

func (tx *asyncTransformer) visitAwaitExpression(node *ast.AwaitExpression) *ast.Node {
	// do not downlevel a top-level await as it is module syntax...
	if tx.inTopLevelContext() {
		return tx.Visitor().VisitEachChild(node.AsNode())
	}
	yieldExpr := tx.Factory().NewYieldExpression(
		nil, /*asteriskToken*/
		tx.Visitor().VisitNode(node.Expression),
	)
	yieldExpr.Loc = node.Loc
	tx.EmitContext().SetOriginal(yieldExpr, node.AsNode())
	return yieldExpr
}

func (tx *asyncTransformer) visitConstructorDeclaration(node *ast.Node) *ast.Node {
	decl := node.AsConstructorDeclaration()
	savedLexicalArgumentsBinding := tx.lexicalArgumentsBinding
	tx.lexicalArgumentsBinding = nil
	updated := tx.Factory().UpdateConstructorDeclaration(
		decl,
		tx.Visitor().VisitModifiers(decl.Modifiers()),
		nil, /*typeParameters*/
		tx.Visitor().VisitNodes(decl.Parameters),
		nil, /*returnType*/
		nil, /*fullSignature*/
		tx.transformMethodBody(node),
	)
	tx.lexicalArgumentsBinding = savedLexicalArgumentsBinding
	return updated
}

func (tx *asyncTransformer) visitMethodDeclaration(node *ast.Node) *ast.Node {
	decl := node.AsMethodDeclaration()
	functionFlags := getFunctionFlags(node)
	savedLexicalArgumentsBinding := tx.lexicalArgumentsBinding
	tx.lexicalArgumentsBinding = nil

	var parameters *ast.NodeList
	var body *ast.Node
	if functionFlags&checker.FunctionFlagsAsync != 0 {
		parameters = tx.transformAsyncFunctionParameterList(node)
		body = tx.transformAsyncFunctionBody(node, parameters)
	} else {
		parameters = tx.Visitor().VisitNodes(decl.Parameters)
		body = tx.transformMethodBody(node)
	}

	updated := tx.Factory().UpdateMethodDeclaration(
		decl,
		tx.Visitor().VisitModifiers(decl.Modifiers()),
		decl.AsteriskToken,
		decl.Name(),
		nil, /*postfixToken*/
		nil, /*typeParameters*/
		parameters,
		nil, /*returnType*/
		nil, /*fullSignature*/
		body,
	)
	tx.lexicalArgumentsBinding = savedLexicalArgumentsBinding
	return updated
}

func (tx *asyncTransformer) visitGetAccessorDeclaration(node *ast.Node) *ast.Node {
	decl := node.AsGetAccessorDeclaration()
	savedLexicalArgumentsBinding := tx.lexicalArgumentsBinding
	tx.lexicalArgumentsBinding = nil
	updated := tx.Factory().UpdateGetAccessorDeclaration(
		decl,
		tx.Visitor().VisitModifiers(decl.Modifiers()),
		decl.Name(),
		nil, /*typeParameters*/
		tx.Visitor().VisitNodes(decl.Parameters),
		nil, /*returnType*/
		nil, /*fullSignature*/
		tx.transformMethodBody(node),
	)
	tx.lexicalArgumentsBinding = savedLexicalArgumentsBinding
	return updated
}

func (tx *asyncTransformer) visitSetAccessorDeclaration(node *ast.Node) *ast.Node {
	decl := node.AsSetAccessorDeclaration()
	savedLexicalArgumentsBinding := tx.lexicalArgumentsBinding
	tx.lexicalArgumentsBinding = nil
	updated := tx.Factory().UpdateSetAccessorDeclaration(
		decl,
		tx.Visitor().VisitModifiers(decl.Modifiers()),
		decl.Name(),
		nil, /*typeParameters*/
		tx.Visitor().VisitNodes(decl.Parameters),
		nil, /*returnType*/
		nil, /*fullSignature*/
		tx.transformMethodBody(node),
	)
	tx.lexicalArgumentsBinding = savedLexicalArgumentsBinding
	return updated
}

func (tx *asyncTransformer) visitFunctionDeclaration(node *ast.Node) *ast.Node {
	decl := node.AsFunctionDeclaration()
	functionFlags := getFunctionFlags(node)
	savedLexicalArgumentsBinding := tx.lexicalArgumentsBinding
	tx.lexicalArgumentsBinding = nil

	var parameters *ast.NodeList
	var body *ast.Node
	if functionFlags&checker.FunctionFlagsAsync != 0 {
		parameters = tx.transformAsyncFunctionParameterList(node)
		body = tx.transformAsyncFunctionBody(node, parameters)
	} else {
		parameters = tx.Visitor().VisitNodes(decl.Parameters)
		body = tx.Visitor().VisitNode(decl.Body)
	}

	updated := tx.Factory().UpdateFunctionDeclaration(
		decl,
		tx.Visitor().VisitModifiers(decl.Modifiers()),
		decl.AsteriskToken,
		tx.Visitor().VisitNode(decl.Name()),
		nil, /*typeParameters*/
		parameters,
		nil, /*returnType*/
		nil, /*fullSignature*/
		body,
	)
	tx.lexicalArgumentsBinding = savedLexicalArgumentsBinding
	return updated
}

func (tx *asyncTransformer) visitFunctionExpression(node *ast.Node) *ast.Node {
	decl := node.AsFunctionExpression()
	functionFlags := getFunctionFlags(node)
	savedLexicalArgumentsBinding := tx.lexicalArgumentsBinding
	tx.lexicalArgumentsBinding = nil

	var parameters *ast.NodeList
	var body *ast.Node
	if functionFlags&checker.FunctionFlagsAsync != 0 {
		parameters = tx.transformAsyncFunctionParameterList(node)
		body = tx.transformAsyncFunctionBody(node, parameters)
	} else {
		parameters = tx.Visitor().VisitNodes(decl.Parameters)
		body = tx.Visitor().VisitNode(decl.Body)
	}

	updated := tx.Factory().UpdateFunctionExpression(
		decl,
		tx.Visitor().VisitModifiers(decl.Modifiers()),
		decl.AsteriskToken,
		tx.Visitor().VisitNode(decl.Name()),
		nil, /*typeParameters*/
		parameters,
		nil, /*returnType*/
		nil, /*fullSignature*/
		body,
	)
	tx.lexicalArgumentsBinding = savedLexicalArgumentsBinding
	return updated
}

func (tx *asyncTransformer) visitArrowFunction(node *ast.Node) *ast.Node {
	decl := node.AsArrowFunction()
	functionFlags := getFunctionFlags(node)

	var parameters *ast.NodeList
	var body *ast.Node
	if functionFlags&checker.FunctionFlagsAsync != 0 {
		parameters = tx.transformAsyncFunctionParameterList(node)
		body = tx.transformAsyncFunctionBody(node, parameters)
	} else {
		parameters = tx.Visitor().VisitNodes(decl.Parameters)
		body = tx.Visitor().VisitNode(decl.Body)
	}

	return tx.Factory().UpdateArrowFunction(
		decl,
		tx.Visitor().VisitModifiers(decl.Modifiers()),
		nil, /*typeParameters*/
		parameters,
		nil, /*returnType*/
		nil, /*fullSignature*/
		decl.EqualsGreaterThanToken,
		body,
	)
}

func (tx *asyncTransformer) transformMethodBody(node *ast.Node) *ast.Node {
	savedCapturedSuperProperties := tx.capturedSuperProperties
	savedHasSuperElementAccess := tx.hasSuperElementAccess
	tx.capturedSuperProperties = &collections.Set[string]{}
	tx.hasSuperElementAccess = false

	updated := tx.Visitor().VisitNode(node.Body())

	// Minor optimization, emit `_super` helper to capture `super` access in an arrow.
	// This step isn't needed if we eventually transform this to ES5.
	// !!! super property access/assignment check flags not yet available on resolver
	// emitSuperHelpers :=
	// 	tx.emitResolver != nil &&
	// 	(resolver.hasNodeCheckFlag(node, NodeCheckFlagsMethodWithSuperPropertyAssignmentInAsync) ||
	// 	 resolver.hasNodeCheckFlag(node, NodeCheckFlagsMethodWithSuperPropertyAccessInAsync)) &&
	// 	(getFunctionFlags(tx.getOriginalIfFunctionLike(node)) & checker.FunctionFlagsAsyncGenerator) != checker.FunctionFlagsAsyncGenerator

	tx.capturedSuperProperties = savedCapturedSuperProperties
	tx.hasSuperElementAccess = savedHasSuperElementAccess
	return updated
}

func (tx *asyncTransformer) transformAsyncFunctionParameterList(node *ast.Node) *ast.NodeList {
	if isSimpleParameterList(node.Parameters()) {
		return tx.Visitor().VisitNodes(node.ParameterList())
	}

	var newParameters []*ast.Node
	for _, parameter := range node.Parameters() {
		param := parameter.AsParameterDeclaration()
		if param.Initializer != nil || param.DotDotDotToken != nil {
			// for an arrow function, capture the remaining arguments in a rest parameter.
			if node.Kind == ast.KindArrowFunction {
				restParameter := tx.Factory().NewParameterDeclaration(
					nil,
					tx.Factory().NewToken(ast.KindDotDotDotToken),
					tx.Factory().NewUniqueNameEx("args", printer.AutoGenerateOptions{Flags: printer.GeneratedIdentifierFlagsReservedInNestedScopes}),
					nil,
					nil,
					nil,
				)
				newParameters = append(newParameters, restParameter)
			}
			break
		}
		// for arrow functions we capture fixed parameters to forward to `__awaiter`. For all other functions
		// we add fixed parameters to preserve the function's `length` property.
		newParameter := tx.Factory().NewParameterDeclaration(
			nil,
			nil,
			tx.Factory().NewGeneratedNameForNodeEx(param.Name(), printer.AutoGenerateOptions{Flags: printer.GeneratedIdentifierFlagsReservedInNestedScopes}),
			nil,
			nil,
			nil,
		)
		newParameters = append(newParameters, newParameter)
	}
	newParametersArray := tx.Factory().NewNodeList(newParameters)
	newParametersArray.Loc = node.ParameterList().Loc
	return newParametersArray
}

func (tx *asyncTransformer) transformAsyncFunctionBody(node *ast.Node, outerParameters *ast.NodeList) *ast.Node {
	innerParameters := (*ast.NodeList)(nil)
	if !isSimpleParameterList(node.Parameters()) {
		innerParameters = tx.Visitor().VisitNodes(node.ParameterList())
	}
	tx.EmitContext().StartVariableEnvironment()

	isArrow := node.Kind == ast.KindArrowFunction
	savedLexicalArgumentsBinding := tx.lexicalArgumentsBinding
	captureLexicalArguments := tx.lexicalArgumentsBinding == nil
	if captureLexicalArguments {
		tx.lexicalArgumentsBinding = tx.Factory().NewUniqueName("arguments")
		tx.usedLexicalArguments = false
	}

	var argumentsExpression *ast.Expression
	if innerParameters != nil {
		if isArrow {
			// `node` does not have a simple parameter list, so `outerParameters` refers to placeholders
			// forwarded to `innerParameters`.
			var parameterBindings []*ast.Node
			outerLen := len(outerParameters.Nodes)
			for i, param := range node.Parameters() {
				if i >= outerLen {
					break
				}
				originalParameter := param.AsParameterDeclaration()
				outerParameter := outerParameters.Nodes[i].AsParameterDeclaration()
				if originalParameter.Initializer != nil || originalParameter.DotDotDotToken != nil {
					parameterBindings = append(parameterBindings, tx.Factory().NewSpreadElement(outerParameter.Name()))
					break
				}
				parameterBindings = append(parameterBindings, outerParameter.Name())
			}
			argumentsExpression = tx.Factory().NewArrayLiteralExpression(tx.Factory().NewNodeList(parameterBindings), false)
		} else {
			argumentsExpression = tx.Factory().NewIdentifier("arguments")
		}
	}

	// An async function is emit as an outer function that calls an inner
	// generator function. To preserve lexical bindings, we pass the current
	// `this` and `arguments` objects to `__awaiter`. The generator function
	// passed to `__awaiter` is executed inside of the callback to the
	// promise constructor.

	savedEnclosingFunctionParameterNames := tx.enclosingFunctionParameterNames
	tx.enclosingFunctionParameterNames = &collections.Set[string]{}
	for _, parameter := range node.Parameters() {
		tx.recordDeclarationName(parameter, tx.enclosingFunctionParameterNames)
	}

	savedCapturedSuperProperties := tx.capturedSuperProperties
	savedHasSuperElementAccess := tx.hasSuperElementAccess
	if !isArrow {
		tx.capturedSuperProperties = &collections.Set[string]{}
		tx.hasSuperElementAccess = false
	}

	hasLexicalThis := tx.inHasLexicalThisContext()

	asyncBody := tx.transformAsyncFunctionBodyWorker(node.Body())
	asyncBody = tx.Factory().UpdateBlock(
		asyncBody.AsBlock(),
		tx.EmitContext().EndAndMergeVariableEnvironmentList(asyncBody.StatementList()),
	)

	var result *ast.Node
	if !isArrow {
		var statements []*ast.Node
		statements = append(statements, tx.Factory().NewReturnStatement(
			tx.Factory().NewAwaiterHelper(
				hasLexicalThis,
				argumentsExpression,
				innerParameters,
				asyncBody,
			),
		))

		// !!! super property access/assignment helpers
		// (requires onEmitNode/onSubstituteNode support)

		if captureLexicalArguments && tx.usedLexicalArguments {
			prologue, rest := tx.Factory().SplitStandardPrologue(statements)
			statements = append(prologue, append([]*ast.Node{tx.createCaptureArgumentsStatement()}, rest...)...)
		}

		block := tx.Factory().NewBlock(tx.Factory().NewNodeList(statements), true)
		block.Loc = node.Body().Loc

		// !!! super element access helpers
		// (requires onEmitNode/onSubstituteNode support)

		result = block
	} else {
		result = tx.Factory().NewAwaiterHelper(
			hasLexicalThis,
			argumentsExpression,
			innerParameters,
			asyncBody,
		)

		if captureLexicalArguments && tx.usedLexicalArguments {
			block := tx.convertToFunctionBlock(result)
			result = tx.Factory().UpdateBlock(
				block.AsBlock(),
				tx.EmitContext().MergeEnvironmentList(block.StatementList(), []*ast.Node{tx.createCaptureArgumentsStatement()}),
			)
		}
	}

	tx.enclosingFunctionParameterNames = savedEnclosingFunctionParameterNames
	if !isArrow {
		tx.capturedSuperProperties = savedCapturedSuperProperties
		tx.hasSuperElementAccess = savedHasSuperElementAccess
		tx.lexicalArgumentsBinding = savedLexicalArgumentsBinding
	}
	return result
}

func (tx *asyncTransformer) transformAsyncFunctionBodyWorker(body *ast.Node) *ast.Node {
	if ast.IsBlock(body) {
		return tx.Factory().UpdateBlock(
			body.AsBlock(),
			tx.asyncBodyVisitor.VisitNodes(body.StatementList()),
		)
	}
	// Convert expression body to block body with return statement
	visited := tx.asyncBodyVisitor.VisitNode(body)
	ret := tx.Factory().NewReturnStatement(visited)
	ret.Loc = body.Loc
	list := tx.Factory().NewNodeList([]*ast.Node{ret})
	list.Loc = body.Loc
	block := tx.Factory().NewBlock(list, false /*multiLine*/)
	block.Loc = body.Loc
	return block
}

func (tx *asyncTransformer) createCaptureArgumentsStatement() *ast.Node {
	variable := tx.Factory().NewVariableDeclaration(
		tx.lexicalArgumentsBinding,
		nil,
		nil,
		tx.Factory().NewIdentifier("arguments"),
	)
	declList := tx.Factory().NewVariableDeclarationList(ast.NodeFlagsNone, tx.Factory().NewNodeList([]*ast.Node{variable}))
	statement := tx.Factory().NewVariableStatement(nil, declList)
	tx.EmitContext().AddEmitFlags(statement, printer.EFStartOnNewLine|printer.EFCustomPrologue)
	return statement
}

func (tx *asyncTransformer) convertToFunctionBlock(node *ast.Node) *ast.Node {
	if ast.IsBlock(node) {
		return node
	}
	ret := tx.Factory().NewReturnStatement(node)
	ret.Loc = node.Loc
	tx.EmitContext().SetOriginal(ret, node)
	list := tx.Factory().NewNodeList([]*ast.Node{ret})
	list.Loc = node.Loc
	block := tx.Factory().NewBlock(list, true)
	block.Loc = node.Loc
	return block
}

func (tx *asyncTransformer) recordDeclarationName(node *ast.Node, names *collections.Set[string]) {
	name := node.Name()
	if name == nil {
		return
	}
	if ast.IsIdentifier(name) {
		names.Add(name.Text())
	} else if ast.IsBindingPattern(name) {
		for _, element := range name.AsBindingPattern().Elements.Nodes {
			if !ast.IsOmittedExpression(element) {
				tx.recordDeclarationName(element, names)
			}
		}
	}
}

func (tx *asyncTransformer) visitCatchClauseInAsyncBody(node *ast.CatchClause) *ast.Node {
	catchClauseNames := &collections.Set[string]{}
	if node.VariableDeclaration != nil {
		tx.recordDeclarationName(node.VariableDeclaration, catchClauseNames)
	}

	// names declared in a catch variable are block scoped
	var catchClauseUnshadowedNames *collections.Set[string]
	for escapedName := range catchClauseNames.Keys() {
		if tx.enclosingFunctionParameterNames != nil && tx.enclosingFunctionParameterNames.Has(escapedName) {
			if catchClauseUnshadowedNames == nil {
				catchClauseUnshadowedNames = tx.enclosingFunctionParameterNames.Clone()
			}
			catchClauseUnshadowedNames.Delete(escapedName)
		}
	}

	if catchClauseUnshadowedNames != nil {
		savedEnclosingFunctionParameterNames := tx.enclosingFunctionParameterNames
		tx.enclosingFunctionParameterNames = catchClauseUnshadowedNames
		result := tx.asyncBodyVisitor.VisitEachChild(node.AsNode())
		tx.enclosingFunctionParameterNames = savedEnclosingFunctionParameterNames
		return result
	}
	return tx.asyncBodyVisitor.VisitEachChild(node.AsNode())
}

func (tx *asyncTransformer) visitVariableStatementInAsyncBody(node *ast.Node) *ast.Node {
	declList := node.AsVariableStatement().DeclarationList
	if tx.isVariableDeclarationListWithCollidingName(declList) {
		expression := tx.visitVariableDeclarationListWithCollidingNames(declList.AsVariableDeclarationList(), false)
		if expression != nil {
			return tx.Factory().NewExpressionStatement(expression)
		}
		return nil
	}
	return tx.Visitor().VisitEachChild(node)
}

func (tx *asyncTransformer) visitForStatementInAsyncBody(node *ast.ForStatement) *ast.Node {
	initializer := node.Initializer
	var visitedInitializer *ast.Node
	if initializer != nil && tx.isVariableDeclarationListWithCollidingName(initializer) {
		visitedInitializer = tx.visitVariableDeclarationListWithCollidingNames(initializer.AsVariableDeclarationList(), false)
	} else {
		visitedInitializer = tx.Visitor().VisitNode(node.Initializer)
	}

	return tx.Factory().UpdateForStatement(
		node,
		visitedInitializer,
		tx.Visitor().VisitNode(node.Condition),
		tx.Visitor().VisitNode(node.Incrementor),
		tx.asyncBodyVisitor.VisitEmbeddedStatement(node.Statement),
	)
}

func (tx *asyncTransformer) visitForInStatementInAsyncBody(node *ast.ForInOrOfStatement) *ast.Node {
	var visitedInitializer *ast.Node
	if tx.isVariableDeclarationListWithCollidingName(node.Initializer) {
		visitedInitializer = tx.visitVariableDeclarationListWithCollidingNames(node.Initializer.AsVariableDeclarationList(), true)
	} else {
		visitedInitializer = tx.Visitor().VisitNode(node.Initializer)
	}

	return tx.Factory().UpdateForInOrOfStatement(
		node,
		nil, /*awaitModifier*/
		visitedInitializer,
		tx.Visitor().VisitNode(node.Expression),
		tx.asyncBodyVisitor.VisitEmbeddedStatement(node.Statement),
	)
}

func (tx *asyncTransformer) visitForOfStatementInAsyncBody(node *ast.ForInOrOfStatement) *ast.Node {
	var visitedInitializer *ast.Node
	if tx.isVariableDeclarationListWithCollidingName(node.Initializer) {
		visitedInitializer = tx.visitVariableDeclarationListWithCollidingNames(node.Initializer.AsVariableDeclarationList(), true)
	} else {
		visitedInitializer = tx.Visitor().VisitNode(node.Initializer)
	}

	return tx.Factory().UpdateForInOrOfStatement(
		node,
		tx.Visitor().VisitNode(node.AwaitModifier),
		visitedInitializer,
		tx.Visitor().VisitNode(node.Expression),
		tx.asyncBodyVisitor.VisitEmbeddedStatement(node.Statement),
	)
}

func (tx *asyncTransformer) isVariableDeclarationListWithCollidingName(node *ast.Node) bool {
	return node != nil &&
		ast.IsVariableDeclarationList(node) &&
		node.Flags&ast.NodeFlagsBlockScoped == 0 &&
		tx.hasCollidingDeclarations(node.AsVariableDeclarationList())
}

func (tx *asyncTransformer) hasCollidingDeclarations(node *ast.VariableDeclarationList) bool {
	return slices.ContainsFunc(node.Declarations.Nodes, tx.collidesWithParameterName)
}

func (tx *asyncTransformer) collidesWithParameterName(node *ast.Node) bool {
	name := node.Name()
	if name == nil {
		return false
	}
	if ast.IsIdentifier(name) {
		return tx.enclosingFunctionParameterNames != nil && tx.enclosingFunctionParameterNames.Has(name.Text())
	}
	if ast.IsBindingPattern(name) {
		for _, element := range name.AsBindingPattern().Elements.Nodes {
			if !ast.IsOmittedExpression(element) && tx.collidesWithParameterName(element) {
				return true
			}
		}
	}
	return false
}

func (tx *asyncTransformer) visitVariableDeclarationListWithCollidingNames(node *ast.VariableDeclarationList, hasReceiver bool) *ast.Node {
	tx.hoistVariableDeclarationList(node)

	var variables []*ast.Node
	for _, decl := range node.Declarations.Nodes {
		if decl.AsVariableDeclaration().Initializer != nil {
			variables = append(variables, decl)
		}
	}

	if len(variables) == 0 {
		if hasReceiver {
			name := node.Declarations.Nodes[0].Name()
			var target *ast.Node
			if ast.IsBindingPattern(name) {
				target = transformers.ConvertBindingPatternToAssignmentPattern(tx.EmitContext(), name.AsBindingPattern())
			} else {
				target = name
			}
			return tx.Visitor().VisitNode(target)
		}
		return nil
	}

	var expressions []*ast.Node
	for _, variable := range variables {
		expressions = append(expressions, tx.transformInitializedVariable(variable.AsVariableDeclaration()))
	}
	return tx.Factory().InlineExpressions(expressions)
}

func (tx *asyncTransformer) hoistVariableDeclarationList(node *ast.VariableDeclarationList) {
	for _, decl := range node.Declarations.Nodes {
		tx.hoistVariable(decl)
	}
}

func (tx *asyncTransformer) hoistVariable(node *ast.Node) {
	name := node.Name()
	if name == nil {
		return
	}
	if ast.IsIdentifier(name) {
		tx.EmitContext().AddVariableDeclaration(name)
	} else if ast.IsBindingPattern(name) {
		for _, element := range name.AsBindingPattern().Elements.Nodes {
			if !ast.IsOmittedExpression(element) {
				tx.hoistVariable(element)
			}
		}
	}
}

func (tx *asyncTransformer) transformInitializedVariable(node *ast.VariableDeclaration) *ast.Node {
	var target *ast.Node
	if ast.IsBindingPattern(node.Name()) {
		target = transformers.ConvertBindingPatternToAssignmentPattern(tx.EmitContext(), node.Name().AsBindingPattern())
	} else {
		target = node.Name()
	}
	converted := tx.Factory().NewAssignmentExpression(target, node.Initializer)
	tx.EmitContext().SetSourceMapRange(converted, node.Loc)
	return tx.Visitor().VisitNode(converted)
}

func (tx *asyncTransformer) getOriginalIfFunctionLike(node *ast.Node) *ast.Node {
	original := tx.EmitContext().MostOriginal(node)
	if original != nil && ast.IsFunctionLikeDeclaration(original) {
		return original
	}
	return node
}

// isEffectiveStrictModeSourceFile checks if the source file is in strict mode.
// alwaysStrict is always true in the Go port.
func isEffectiveStrictModeSourceFile(_ *ast.SourceFile, _ *core.CompilerOptions) bool {
	return true
}

// isSimpleParameterList checks if every parameter has no initializer and an Identifier name.
func isSimpleParameterList(params []*ast.Node) bool {
	for _, param := range params {
		p := param.AsParameterDeclaration()
		if p.Initializer != nil || !ast.IsIdentifier(p.Name()) {
			return false
		}
	}
	return true
}

// getFunctionFlags returns the function flags for a node.
func getFunctionFlags(node *ast.Node) checker.FunctionFlags {
	if node == nil {
		return checker.FunctionFlagsInvalid
	}
	data := node.BodyData()
	if data == nil {
		return checker.FunctionFlagsInvalid
	}
	flags := checker.FunctionFlagsNormal
	switch node.Kind {
	case ast.KindFunctionDeclaration, ast.KindFunctionExpression, ast.KindMethodDeclaration:
		if data.AsteriskToken != nil {
			flags |= checker.FunctionFlagsGenerator
		}
		fallthrough
	case ast.KindArrowFunction:
		if ast.HasSyntacticModifier(node, ast.ModifierFlagsAsync) {
			flags |= checker.FunctionFlagsAsync
		}
	}
	if data.Body == nil {
		flags |= checker.FunctionFlagsInvalid
	}
	return flags
}

// isNodeWithPossibleHoistedDeclaration checks if a node could contain hoisted declarations.
func isNodeWithPossibleHoistedDeclaration(node *ast.Node) bool {
	switch node.Kind {
	case ast.KindBlock,
		ast.KindVariableStatement,
		ast.KindWithStatement,
		ast.KindIfStatement,
		ast.KindSwitchStatement,
		ast.KindCaseBlock,
		ast.KindCaseClause,
		ast.KindDefaultClause,
		ast.KindLabeledStatement,
		ast.KindForStatement,
		ast.KindForInStatement,
		ast.KindForOfStatement,
		ast.KindDoStatement,
		ast.KindWhileStatement,
		ast.KindTryStatement,
		ast.KindCatchClause:
		return true
	}
	return false
}

func newAsyncTransformer(opts *transformers.TransformOptions) *transformers.Transformer {
	tx := &asyncTransformer{
		compilerOptions: opts.CompilerOptions,
	}
	if opts.EmitResolver != nil {
		tx.emitResolver = opts.EmitResolver
	}
	result := tx.NewTransformer(tx.visit, opts.Context)
	tx.asyncBodyVisitor = tx.EmitContext().NewNodeVisitor(tx.visitAsyncBodyNode)
	return result
}
