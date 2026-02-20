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

type asyncContextFlags int

const (
	asyncContextNone        asyncContextFlags = 0
	asyncContextNonTopLevel asyncContextFlags = 1 << iota
	asyncContextHasLexicalThis
)

type lexicalArgumentsInfo struct {
	binding *ast.IdentifierNode
	used    bool
}

type asyncTransformer struct {
	transformers.Transformer
	compilerOptions *core.CompilerOptions

	contextFlags asyncContextFlags

	enclosingFunctionParameterNames *collections.Set[string]
	capturedSuperProperties         *collections.Set[string]
	hasSuperElementAccess           bool
	hasSuperPropertyAssignment      bool
	superBinding                    *ast.IdentifierNode
	superIndexBinding               *ast.IdentifierNode
	lexicalArguments                lexicalArgumentsInfo

	asyncBodyVisitor *ast.NodeVisitor
}

func (tx *asyncTransformer) visit(node *ast.Node) *ast.Node {
	if node.Kind == ast.KindAsyncKeyword {
		// ES2017 async modifier should be elided for targets < ES2017
		return nil
	}
	if node.SubtreeFacts()&(ast.SubtreeContainsAnyAwait|ast.SubtreeContainsAwait) == 0 {
		if tx.capturedSuperProperties != nil || tx.lexicalArguments.binding != nil {
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
		return tx.doWithContext(asyncContextNonTopLevel|asyncContextHasLexicalThis, tx.visitMethodDeclaration, node)
	case ast.KindFunctionDeclaration:
		return tx.doWithContext(asyncContextNonTopLevel|asyncContextHasLexicalThis, tx.visitFunctionDeclaration, node)
	case ast.KindFunctionExpression:
		return tx.doWithContext(asyncContextNonTopLevel|asyncContextHasLexicalThis, tx.visitFunctionExpression, node)
	case ast.KindArrowFunction:
		return tx.doWithContext(asyncContextNonTopLevel, tx.visitArrowFunction, node)
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
	case ast.KindBinaryExpression:
		if tx.capturedSuperProperties != nil && ast.IsAssignmentOperator(node.AsBinaryExpression().OperatorToken.Kind) && assignmentTargetContainsSuperProperty(node.AsBinaryExpression().Left) {
			tx.hasSuperPropertyAssignment = true
		}
		return tx.Visitor().VisitEachChild(node)
	case ast.KindPrefixUnaryExpression:
		if tx.capturedSuperProperties != nil && isUpdateExpression(node) && assignmentTargetContainsSuperProperty(node.AsPrefixUnaryExpression().Operand) {
			tx.hasSuperPropertyAssignment = true
		}
		return tx.Visitor().VisitEachChild(node)
	case ast.KindPostfixUnaryExpression:
		if tx.capturedSuperProperties != nil && isUpdateExpression(node) && assignmentTargetContainsSuperProperty(node.AsPostfixUnaryExpression().Operand) {
			tx.hasSuperPropertyAssignment = true
		}
		return tx.Visitor().VisitEachChild(node)
	case ast.KindGetAccessor:
		return tx.doWithContext(asyncContextNonTopLevel|asyncContextHasLexicalThis, tx.visitGetAccessorDeclaration, node)
	case ast.KindSetAccessor:
		return tx.doWithContext(asyncContextNonTopLevel|asyncContextHasLexicalThis, tx.visitSetAccessorDeclaration, node)
	case ast.KindConstructor:
		return tx.doWithContext(asyncContextNonTopLevel|asyncContextHasLexicalThis, tx.visitConstructorDeclaration, node)
	case ast.KindClassDeclaration, ast.KindClassExpression:
		return tx.doWithContext(asyncContextNonTopLevel|asyncContextHasLexicalThis, tx.visitDefault, node)
	default:
		return tx.Visitor().VisitEachChild(node)
	}
}

func (tx *asyncTransformer) visitSourceFile(node *ast.SourceFile) *ast.Node {
	if node.IsDeclarationFile {
		return node.AsNode()
	}

	tx.setContextFlag(asyncContextNonTopLevel, false)
	tx.setContextFlag(asyncContextHasLexicalThis, !isEffectiveStrictModeSourceFile(node, tx.compilerOptions))
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
		if tx.lexicalArguments.binding != nil && isArgumentsIdentifier(node) {
			tx.lexicalArguments.used = true
			return tx.lexicalArguments.binding
		}
	case ast.KindPropertyAccessExpression:
		if tx.capturedSuperProperties != nil && node.Expression().Kind == ast.KindSuperKeyword {
			tx.capturedSuperProperties.Add(node.Name().Text())
		}
	case ast.KindElementAccessExpression:
		if tx.capturedSuperProperties != nil && node.Expression().Kind == ast.KindSuperKeyword {
			tx.hasSuperElementAccess = true
		}
	case ast.KindBinaryExpression:
		if tx.capturedSuperProperties != nil && ast.IsAssignmentOperator(node.AsBinaryExpression().OperatorToken.Kind) && assignmentTargetContainsSuperProperty(node.AsBinaryExpression().Left) {
			tx.hasSuperPropertyAssignment = true
		}
	case ast.KindPrefixUnaryExpression:
		if tx.capturedSuperProperties != nil && isUpdateExpression(node) && assignmentTargetContainsSuperProperty(node.AsPrefixUnaryExpression().Operand) {
			tx.hasSuperPropertyAssignment = true
		}
	case ast.KindPostfixUnaryExpression:
		if tx.capturedSuperProperties != nil && isUpdateExpression(node) && assignmentTargetContainsSuperProperty(node.AsPostfixUnaryExpression().Operand) {
			tx.hasSuperPropertyAssignment = true
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

func (tx *asyncTransformer) setContextFlag(flag asyncContextFlags, val bool) {
	if val {
		tx.contextFlags |= flag
	} else {
		tx.contextFlags &= ^flag
	}
}

func (tx *asyncTransformer) inContext(flags asyncContextFlags) bool {
	return tx.contextFlags&flags != 0
}

func (tx *asyncTransformer) inTopLevelContext() bool {
	return !tx.inContext(asyncContextNonTopLevel)
}

func (tx *asyncTransformer) inHasLexicalThisContext() bool {
	return tx.inContext(asyncContextHasLexicalThis)
}

func (tx *asyncTransformer) doWithContext(flags asyncContextFlags, cb func(*ast.Node) *ast.Node, node *ast.Node) *ast.Node {
	flagsToSet := flags & ^tx.contextFlags
	if flagsToSet != 0 {
		tx.setContextFlag(flagsToSet, true)
		result := cb(node)
		tx.setContextFlag(flagsToSet, false)
		return result
	}
	return cb(node)
}

func (tx *asyncTransformer) visitDefault(node *ast.Node) *ast.Node {
	return tx.Visitor().VisitEachChild(node)
}

// visitAwaitExpression visits an AwaitExpression node.
//
// This function will be called any time a ES2017 await expression is encountered.
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
	savedLexicalArguments := tx.lexicalArguments
	tx.lexicalArguments = lexicalArgumentsInfo{}
	updated := tx.Factory().UpdateConstructorDeclaration(
		decl,
		tx.Visitor().VisitModifiers(decl.Modifiers()),
		nil, /*typeParameters*/
		tx.Visitor().VisitNodes(decl.Parameters),
		nil, /*returnType*/
		nil, /*fullSignature*/
		tx.transformMethodBody(node),
	)
	tx.lexicalArguments = savedLexicalArguments
	return updated
}

// visitMethodDeclaration visits a MethodDeclaration node.
//
// This function will be called when one of the following conditions are met:
// - The node is marked as async
func (tx *asyncTransformer) visitMethodDeclaration(node *ast.Node) *ast.Node {
	decl := node.AsMethodDeclaration()
	functionFlags := getFunctionFlags(node)
	savedLexicalArguments := tx.lexicalArguments
	tx.lexicalArguments = lexicalArgumentsInfo{}

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
	tx.lexicalArguments = savedLexicalArguments
	return updated
}

func (tx *asyncTransformer) visitGetAccessorDeclaration(node *ast.Node) *ast.Node {
	decl := node.AsGetAccessorDeclaration()
	savedLexicalArguments := tx.lexicalArguments
	tx.lexicalArguments = lexicalArgumentsInfo{}
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
	tx.lexicalArguments = savedLexicalArguments
	return updated
}

func (tx *asyncTransformer) visitSetAccessorDeclaration(node *ast.Node) *ast.Node {
	decl := node.AsSetAccessorDeclaration()
	savedLexicalArguments := tx.lexicalArguments
	tx.lexicalArguments = lexicalArgumentsInfo{}
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
	tx.lexicalArguments = savedLexicalArguments
	return updated
}

// visitFunctionDeclaration visits a FunctionDeclaration node.
//
// This function will be called when one of the following conditions are met:
// - The node is marked async
func (tx *asyncTransformer) visitFunctionDeclaration(node *ast.Node) *ast.Node {
	decl := node.AsFunctionDeclaration()
	functionFlags := getFunctionFlags(node)
	savedLexicalArguments := tx.lexicalArguments
	tx.lexicalArguments = lexicalArgumentsInfo{}

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
	tx.lexicalArguments = savedLexicalArguments
	return updated
}

// visitFunctionExpression visits a FunctionExpression node.
//
// This function will be called when one of the following conditions are met:
// - The node is marked async
func (tx *asyncTransformer) visitFunctionExpression(node *ast.Node) *ast.Node {
	decl := node.AsFunctionExpression()
	functionFlags := getFunctionFlags(node)
	savedLexicalArguments := tx.lexicalArguments
	tx.lexicalArguments = lexicalArgumentsInfo{}

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
	tx.lexicalArguments = savedLexicalArguments
	return updated
}

// visitArrowFunction visits an ArrowFunction.
//
// This function will be called when one of the following conditions are met:
// - The node is marked async
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
	savedHasSuperPropertyAssignment := tx.hasSuperPropertyAssignment
	savedSuperBinding := tx.superBinding
	savedSuperIndexBinding := tx.superIndexBinding
	tx.capturedSuperProperties = &collections.Set[string]{}
	tx.hasSuperElementAccess = false
	tx.hasSuperPropertyAssignment = false
	tx.superBinding = tx.Factory().NewUniqueNameEx("_super", printer.AutoGenerateOptions{Flags: printer.GeneratedIdentifierFlagsOptimistic | printer.GeneratedIdentifierFlagsFileLevel})
	tx.superIndexBinding = tx.Factory().NewUniqueNameEx("_superIndex", printer.AutoGenerateOptions{Flags: printer.GeneratedIdentifierFlagsOptimistic | printer.GeneratedIdentifierFlagsFileLevel})

	updated := tx.Visitor().VisitNode(node.Body())

	// Minor optimization, emit `_super` helper to capture `super` access in an arrow.
	// This step isn't needed if we eventually transform this to ES5.
	emitSuperHelpers := (tx.capturedSuperProperties.Len() > 0 || tx.hasSuperElementAccess) &&
		(getFunctionFlags(tx.getOriginalIfFunctionLike(node))&checker.FunctionFlagsAsyncGenerator) != checker.FunctionFlagsAsyncGenerator

	if emitSuperHelpers {
		block := updated.AsBlock()
		statements := append([]*ast.Node{}, block.Statements.Nodes...)
		var helpers []*ast.Node
		multiLine := block.Multiline
		if tx.hasSuperElementAccess {
			helpers = append(helpers, tx.createSuperIndexVariableStatement())
			multiLine = true
		}
		if tx.capturedSuperProperties.Len() > 0 {
			helpers = append(helpers, tx.createSuperAccessVariableStatement())
		}
		prologue, rest := tx.Factory().SplitStandardPrologue(statements)
		statements = slices.Clone(prologue)
		statements = append(statements, helpers...)
		statements = append(statements, rest...)
		if multiLine != block.Multiline {
			newBlock := tx.Factory().NewBlock(tx.Factory().NewNodeList(statements), multiLine)
			newBlock.Loc = updated.Loc
			updated = newBlock
		} else {
			updated = tx.Factory().UpdateBlock(block, tx.Factory().NewNodeList(statements))
		}
	}

	tx.capturedSuperProperties = savedCapturedSuperProperties
	tx.hasSuperElementAccess = savedHasSuperElementAccess
	tx.hasSuperPropertyAssignment = savedHasSuperPropertyAssignment
	tx.superBinding = savedSuperBinding
	tx.superIndexBinding = savedSuperIndexBinding
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
			// for any other function/method this isn't necessary as we can just use `arguments`.
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
	savedLexicalArguments := tx.lexicalArguments
	captureLexicalArguments := tx.lexicalArguments.binding == nil
	if captureLexicalArguments {
		tx.lexicalArguments = lexicalArgumentsInfo{
			binding: tx.Factory().NewUniqueName("arguments"),
		}
	}

	var argumentsExpression *ast.Expression
	if innerParameters != nil {
		if isArrow {
			// `node` does not have a simple parameter list, so `outerParameters` refers to placeholders that are
			// forwarded to `innerParameters`, matching how they are introduced in `transformAsyncFunctionParameterList`.
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
	savedHasSuperPropertyAssignment := tx.hasSuperPropertyAssignment
	savedSuperBinding := tx.superBinding
	savedSuperIndexBinding := tx.superIndexBinding
	if !isArrow {
		tx.capturedSuperProperties = &collections.Set[string]{}
		tx.hasSuperElementAccess = false
		tx.hasSuperPropertyAssignment = false
		tx.superBinding = tx.Factory().NewUniqueNameEx("_super", printer.AutoGenerateOptions{Flags: printer.GeneratedIdentifierFlagsOptimistic | printer.GeneratedIdentifierFlagsFileLevel})
		tx.superIndexBinding = tx.Factory().NewUniqueNameEx("_superIndex", printer.AutoGenerateOptions{Flags: printer.GeneratedIdentifierFlagsOptimistic | printer.GeneratedIdentifierFlagsFileLevel})
	}

	hasLexicalThis := tx.inHasLexicalThisContext()

	asyncBody := tx.transformAsyncFunctionBodyWorker(node.Body())
	asyncBody = tx.Factory().UpdateBlock(
		asyncBody.AsBlock(),
		tx.EmitContext().EndAndMergeVariableEnvironmentList(asyncBody.StatementList()),
	)

	// Substitute super property accesses with _super/_superIndex helpers
	emitSuperHelpers := tx.capturedSuperProperties != nil &&
		(tx.capturedSuperProperties.Len() > 0 || tx.hasSuperElementAccess)
	if emitSuperHelpers {
		asyncBody = tx.substituteSuperAccessesInBody(asyncBody)
	}

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

		// Minor optimization, emit `_super` helper to capture `super` access in an arrow.
		// This step isn't needed if we eventually transform this to ES5.
		if emitSuperHelpers {
			var superHelpers []*ast.Node
			if tx.hasSuperElementAccess {
				superHelpers = append(superHelpers, tx.createSuperIndexVariableStatement())
			}
			if tx.capturedSuperProperties.Len() > 0 {
				superHelpers = append(superHelpers, tx.createSuperAccessVariableStatement())
			}
			prologue, rest := tx.Factory().SplitStandardPrologue(statements)
			statements = slices.Clone(prologue)
			statements = append(statements, superHelpers...)
			statements = append(statements, rest...)
		}

		if captureLexicalArguments && tx.lexicalArguments.used {
			prologue, rest := tx.Factory().SplitStandardPrologue(statements)
			statements = slices.Clone(prologue)
			statements = append(statements, tx.createCaptureArgumentsStatement())
			statements = append(statements, rest...)
		}

		block := tx.Factory().NewBlock(tx.Factory().NewNodeList(statements), true)
		block.Loc = node.Body().Loc

		result = block
	} else {
		result = tx.Factory().NewAwaiterHelper(
			hasLexicalThis,
			argumentsExpression,
			innerParameters,
			asyncBody,
		)

		if captureLexicalArguments && tx.lexicalArguments.used {
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
		tx.hasSuperPropertyAssignment = savedHasSuperPropertyAssignment
		tx.superBinding = savedSuperBinding
		tx.superIndexBinding = savedSuperIndexBinding
		tx.lexicalArguments = savedLexicalArguments
	} else if captureLexicalArguments && !tx.lexicalArguments.used {
		// If we created a new binding but it wasn't used, restore the previous state.
		// If it was used, keep the binding alive so sibling arrows can reuse it
		// (the `var` declaration hoists to the enclosing function scope).
		tx.lexicalArguments = savedLexicalArguments
	} else if captureLexicalArguments {
		// Keep the binding but clear the used flag so siblings don't re-emit the capture statement.
		tx.lexicalArguments.used = false
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
		tx.lexicalArguments.binding,
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

// isSuperProperty checks if a node is super.x or super[x].
func isSuperProperty(node *ast.Node) bool {
	return (ast.IsPropertyAccessExpression(node) || ast.IsElementAccessExpression(node)) &&
		node.Expression().Kind == ast.KindSuperKeyword
}

// assignmentTargetContainsSuperProperty checks top-down whether an assignment target
// expression contains a super property or element access (super.x or super[x]).
// This avoids relying on parent pointers (IsAssignmentTarget) which may not be set
// on synthesized AST nodes from prior transforms.
func assignmentTargetContainsSuperProperty(node *ast.Node) bool {
	switch node.Kind {
	case ast.KindPropertyAccessExpression, ast.KindElementAccessExpression:
		return node.Expression().Kind == ast.KindSuperKeyword
	case ast.KindParenthesizedExpression:
		return assignmentTargetContainsSuperProperty(node.AsParenthesizedExpression().Expression)
	case ast.KindArrayLiteralExpression:
		return slices.ContainsFunc(node.AsArrayLiteralExpression().Elements.Nodes, assignmentTargetContainsSuperProperty)
	case ast.KindObjectLiteralExpression:
		for _, prop := range node.AsObjectLiteralExpression().Properties.Nodes {
			switch prop.Kind {
			case ast.KindPropertyAssignment:
				if assignmentTargetContainsSuperProperty(prop.AsPropertyAssignment().Initializer) {
					return true
				}
			case ast.KindShorthandPropertyAssignment:
				if assignmentTargetContainsSuperProperty(prop.AsShorthandPropertyAssignment().Name()) {
					return true
				}
			case ast.KindSpreadAssignment:
				if assignmentTargetContainsSuperProperty(prop.AsSpreadAssignment().Expression) {
					return true
				}
			}
		}
	case ast.KindSpreadElement:
		return assignmentTargetContainsSuperProperty(node.AsSpreadElement().Expression)
	}
	return false
}

// isUpdateExpression checks if a prefix/postfix unary expression is ++ or --.
func isUpdateExpression(node *ast.Node) bool {
	if ast.IsPrefixUnaryExpression(node) {
		op := node.AsPrefixUnaryExpression().Operator
		return op == ast.KindPlusPlusToken || op == ast.KindMinusMinusToken
	}
	if ast.IsPostfixUnaryExpression(node) {
		op := node.AsPostfixUnaryExpression().Operator
		return op == ast.KindPlusPlusToken || op == ast.KindMinusMinusToken
	}
	return false
}

// substituteSuperAccessesInBody walks the async body and replaces super property/element
// accesses with _super/_superIndex references. This is necessary because the async body
// ends up inside a generator function where `super` is not valid.
func (tx *asyncTransformer) substituteSuperAccessesInBody(body *ast.Node) *ast.Node {
	var visitor *ast.NodeVisitor
	var doVisit func(node *ast.Node) *ast.Node
	doVisit = func(node *ast.Node) *ast.Node {
		switch node.Kind {
		case ast.KindCallExpression:
			call := node.AsCallExpression()
			if isSuperProperty(call.Expression) {
				return tx.substituteCallExpressionWithSuperAccess(call, visitor)
			}
			return visitor.VisitEachChild(node)
		case ast.KindPropertyAccessExpression:
			if node.Expression().Kind == ast.KindSuperKeyword {
				// super.x → _super.x
				return tx.Factory().NewPropertyAccessExpression(
					tx.superBinding, nil, node.Name(), ast.NodeFlagsNone,
				)
			}
			return visitor.VisitEachChild(node)
		case ast.KindElementAccessExpression:
			if node.Expression().Kind == ast.KindSuperKeyword {
				// super[x] → _superIndex(x) or _superIndex(x).value
				return tx.createSuperElementAccessInAsyncMethod(
					node.AsElementAccessExpression().ArgumentExpression,
				)
			}
			return visitor.VisitEachChild(node)
		// Don't recurse into non-arrow function scopes or classes
		case ast.KindFunctionExpression, ast.KindFunctionDeclaration,
			ast.KindMethodDeclaration, ast.KindGetAccessor, ast.KindSetAccessor,
			ast.KindConstructor, ast.KindClassDeclaration, ast.KindClassExpression:
			return node
		default:
			return visitor.VisitEachChild(node)
		}
	}
	visitor = tx.EmitContext().NewNodeVisitor(doVisit)
	return visitor.VisitNode(body)
}

// substituteCallExpressionWithSuperAccess handles super.x(args) and super[x](args).
func (tx *asyncTransformer) substituteCallExpressionWithSuperAccess(call *ast.CallExpression, visitor *ast.NodeVisitor) *ast.Node {
	expression := call.Expression
	var target *ast.Node

	if ast.IsPropertyAccessExpression(expression) {
		// super.x(args) → _super.x.call(this, args)
		target = tx.Factory().NewPropertyAccessExpression(
			tx.superBinding, nil,
			expression.AsPropertyAccessExpression().Name(), ast.NodeFlagsNone,
		)
	} else if ast.IsElementAccessExpression(expression) {
		// super[x](args) → _superIndex(x).call(this, args) or _superIndex(x).value.call(this, args)
		target = tx.createSuperElementAccessInAsyncMethod(
			expression.AsElementAccessExpression().ArgumentExpression,
		)
	} else {
		return visitor.VisitEachChild(call.AsNode())
	}

	callTarget := tx.Factory().NewPropertyAccessExpression(
		target, nil,
		tx.Factory().NewIdentifier("call"), ast.NodeFlagsNone,
	)

	var allArgs []*ast.Node
	allArgs = append(allArgs, tx.Factory().NewThisExpression())
	if call.Arguments != nil {
		visitedArgs := visitor.VisitNodes(call.Arguments)
		if visitedArgs != nil {
			allArgs = append(allArgs, visitedArgs.Nodes...)
		}
	}

	result := tx.Factory().NewCallExpression(
		callTarget, nil, nil,
		tx.Factory().NewNodeList(allArgs), ast.NodeFlagsNone,
	)
	result.Loc = call.Loc
	return result
}

// createSuperElementAccessInAsyncMethod creates _superIndex(x) or _superIndex(x).value.
func (tx *asyncTransformer) createSuperElementAccessInAsyncMethod(argumentExpression *ast.Node) *ast.Node {
	superIndexCall := tx.Factory().NewCallExpression(
		tx.superIndexBinding, nil, nil,
		tx.Factory().NewNodeList([]*ast.Node{argumentExpression}),
		ast.NodeFlagsNone,
	)
	if tx.hasSuperPropertyAssignment {
		return tx.Factory().NewPropertyAccessExpression(
			superIndexCall, nil,
			tx.Factory().NewIdentifier("value"), ast.NodeFlagsNone,
		)
	}
	return superIndexCall
}

// createSuperAccessVariableStatement creates a variable named `_super` with accessor
// properties for the given property names.
//
// Create a variable declaration with a getter/setter (if binding) definition for each name:
//
//	const _super = Object.create(null, {
//	    x: { get: () => super.x },                           // read-only
//	    x: { get: () => super.x, set: (v) => super.x = v }, // read-write
//	});
func (tx *asyncTransformer) createSuperAccessVariableStatement() *ast.Node {
	f := tx.Factory()
	var accessors []*ast.Node

	var sortedNames []string
	for name := range tx.capturedSuperProperties.Keys() {
		sortedNames = append(sortedNames, name)
	}
	slices.Sort(sortedNames)

	for _, name := range sortedNames {
		var descriptorProperties []*ast.Node

		// getter: get: () => super.name
		getterBody := f.NewPropertyAccessExpression(
			f.NewKeywordExpression(ast.KindSuperKeyword), nil,
			f.NewIdentifier(name), ast.NodeFlagsNone,
		)
		getterArrow := f.NewArrowFunction(
			nil, nil,
			f.NewNodeList([]*ast.Node{}),
			nil, nil,
			f.NewToken(ast.KindEqualsGreaterThanToken),
			getterBody,
		)
		getter := f.NewPropertyAssignment(nil, f.NewIdentifier("get"), nil, nil, getterArrow)
		descriptorProperties = append(descriptorProperties, getter)

		if tx.hasSuperPropertyAssignment {
			// setter: set: v => super.name = v
			vParam := f.NewParameterDeclaration(nil, nil, f.NewIdentifier("v"), nil, nil, nil)
			superProp := f.NewPropertyAccessExpression(
				f.NewKeywordExpression(ast.KindSuperKeyword), nil,
				f.NewIdentifier(name), ast.NodeFlagsNone,
			)
			assignExpr := f.NewAssignmentExpression(superProp, f.NewIdentifier("v"))
			setterArrow := f.NewArrowFunction(
				nil, nil,
				f.NewNodeList([]*ast.Node{vParam}),
				nil, nil,
				f.NewToken(ast.KindEqualsGreaterThanToken),
				assignExpr,
			)
			setter := f.NewPropertyAssignment(nil, f.NewIdentifier("set"), nil, nil, setterArrow)
			descriptorProperties = append(descriptorProperties, setter)
		}

		descriptor := f.NewObjectLiteralExpression(f.NewNodeList(descriptorProperties), false)
		accessor := f.NewPropertyAssignment(nil, f.NewIdentifier(name), nil, nil, descriptor)
		accessors = append(accessors, accessor)
	}

	descriptorsObject := f.NewObjectLiteralExpression(f.NewNodeList(accessors), true)

	objectCreateCall := f.NewCallExpression(
		f.NewPropertyAccessExpression(
			f.NewIdentifier("Object"), nil,
			f.NewIdentifier("create"), ast.NodeFlagsNone,
		), nil, nil,
		f.NewNodeList([]*ast.Node{
			f.NewKeywordExpression(ast.KindNullKeyword),
			descriptorsObject,
		}),
		ast.NodeFlagsNone,
	)

	decl := f.NewVariableDeclaration(tx.superBinding, nil, nil, objectCreateCall)
	declList := f.NewVariableDeclarationList(ast.NodeFlagsConst, f.NewNodeList([]*ast.Node{decl}))
	return f.NewVariableStatement(nil, declList)
}

// createSuperIndexVariableStatement creates the _superIndex helper variable.
func (tx *asyncTransformer) createSuperIndexVariableStatement() *ast.Node {
	if tx.hasSuperPropertyAssignment {
		return tx.createAdvancedSuperIndexVariableStatement()
	}
	return tx.createSimpleSuperIndexVariableStatement()
}

// createSimpleSuperIndexVariableStatement creates: const _superIndex = name => super[name];
func (tx *asyncTransformer) createSimpleSuperIndexVariableStatement() *ast.Node {
	f := tx.Factory()
	nameParam := f.NewParameterDeclaration(nil, nil, f.NewIdentifier("name"), nil, nil, nil)
	superElementAccess := f.NewElementAccessExpression(
		f.NewKeywordExpression(ast.KindSuperKeyword), nil,
		f.NewIdentifier("name"), ast.NodeFlagsNone,
	)
	arrow := f.NewArrowFunction(
		nil, nil,
		f.NewNodeList([]*ast.Node{nameParam}),
		nil, nil,
		f.NewToken(ast.KindEqualsGreaterThanToken),
		superElementAccess,
	)
	decl := f.NewVariableDeclaration(tx.superIndexBinding, nil, nil, arrow)
	declList := f.NewVariableDeclarationList(ast.NodeFlagsConst, f.NewNodeList([]*ast.Node{decl}))
	return f.NewVariableStatement(nil, declList)
}

// createAdvancedSuperIndexVariableStatement creates:
//
//	const _superIndex = (function (geti, seti) {
//	    const cache = Object.create(null);
//	    return name => cache[name] || (cache[name] = { get value() { return geti(name); }, set value(v) { seti(name, v); } });
//	})(name => super[name], (name, value) => super[name] = value);
func (tx *asyncTransformer) createAdvancedSuperIndexVariableStatement() *ast.Node {
	f := tx.Factory()

	// const cache = Object.create(null)
	objectCreateNull := f.NewCallExpression(
		f.NewPropertyAccessExpression(f.NewIdentifier("Object"), nil, f.NewIdentifier("create"), ast.NodeFlagsNone),
		nil, nil,
		f.NewNodeList([]*ast.Node{f.NewKeywordExpression(ast.KindNullKeyword)}),
		ast.NodeFlagsNone,
	)
	cacheDecl := f.NewVariableDeclaration(f.NewIdentifier("cache"), nil, nil, objectCreateNull)
	cacheDeclList := f.NewVariableDeclarationList(ast.NodeFlagsConst, f.NewNodeList([]*ast.Node{cacheDecl}))
	cacheStmt := f.NewVariableStatement(nil, cacheDeclList)

	// geti(name)
	getiCall := f.NewCallExpression(f.NewIdentifier("geti"), nil, nil, f.NewNodeList([]*ast.Node{f.NewIdentifier("name")}), ast.NodeFlagsNone)
	// seti(name, v)
	setiCall := f.NewCallExpression(f.NewIdentifier("seti"), nil, nil, f.NewNodeList([]*ast.Node{f.NewIdentifier("name"), f.NewIdentifier("v")}), ast.NodeFlagsNone)

	// { get value() { return geti(name); }, set value(v) { seti(name, v); } }
	getterBody := f.NewBlock(f.NewNodeList([]*ast.Node{f.NewReturnStatement(getiCall)}), false)
	getAccessor := f.NewGetAccessorDeclaration(nil, f.NewIdentifier("value"), nil, f.NewNodeList([]*ast.Node{}), nil, nil, getterBody)

	setterVParam := f.NewParameterDeclaration(nil, nil, f.NewIdentifier("v"), nil, nil, nil)
	setterBody := f.NewBlock(f.NewNodeList([]*ast.Node{f.NewExpressionStatement(setiCall)}), false)
	setAccessor := f.NewSetAccessorDeclaration(nil, f.NewIdentifier("value"), nil, f.NewNodeList([]*ast.Node{setterVParam}), nil, nil, setterBody)

	descriptor := f.NewObjectLiteralExpression(f.NewNodeList([]*ast.Node{getAccessor, setAccessor}), false)

	// cache[name] || (cache[name] = descriptor)
	cacheAccess1 := f.NewElementAccessExpression(f.NewIdentifier("cache"), nil, f.NewIdentifier("name"), ast.NodeFlagsNone)
	cacheAccess2 := f.NewElementAccessExpression(f.NewIdentifier("cache"), nil, f.NewIdentifier("name"), ast.NodeFlagsNone)
	cacheAssign := f.NewParenthesizedExpression(f.NewAssignmentExpression(cacheAccess2, descriptor))
	orExpr := f.NewBinaryExpression(nil, cacheAccess1, nil, f.NewToken(ast.KindBarBarToken), cacheAssign)

	// name => cache[name] || (cache[name] = descriptor)
	innerNameParam := f.NewParameterDeclaration(nil, nil, f.NewIdentifier("name"), nil, nil, nil)
	innerArrow := f.NewArrowFunction(nil, nil, f.NewNodeList([]*ast.Node{innerNameParam}), nil, nil, f.NewToken(ast.KindEqualsGreaterThanToken), orExpr)

	// return innerArrow
	returnStmt := f.NewReturnStatement(innerArrow)

	// function(geti, seti) { const cache = ...; return ...; }
	funcBody := f.NewBlock(f.NewNodeList([]*ast.Node{cacheStmt, returnStmt}), true)
	getiParam := f.NewParameterDeclaration(nil, nil, f.NewIdentifier("geti"), nil, nil, nil)
	setiParam := f.NewParameterDeclaration(nil, nil, f.NewIdentifier("seti"), nil, nil, nil)
	outerFunc := f.NewFunctionExpression(nil, nil, nil, nil, f.NewNodeList([]*ast.Node{getiParam, setiParam}), nil, nil, funcBody)

	// Getter arg: name => super[name]
	getterArgParam := f.NewParameterDeclaration(nil, nil, f.NewIdentifier("name"), nil, nil, nil)
	getterArgBody := f.NewElementAccessExpression(f.NewKeywordExpression(ast.KindSuperKeyword), nil, f.NewIdentifier("name"), ast.NodeFlagsNone)
	getterArg := f.NewArrowFunction(nil, nil, f.NewNodeList([]*ast.Node{getterArgParam}), nil, nil, f.NewToken(ast.KindEqualsGreaterThanToken), getterArgBody)

	// Setter arg: (name, value) => super[name] = value
	setterArgNameParam := f.NewParameterDeclaration(nil, nil, f.NewIdentifier("name"), nil, nil, nil)
	setterArgValueParam := f.NewParameterDeclaration(nil, nil, f.NewIdentifier("value"), nil, nil, nil)
	setterArgSuperAccess := f.NewElementAccessExpression(f.NewKeywordExpression(ast.KindSuperKeyword), nil, f.NewIdentifier("name"), ast.NodeFlagsNone)
	setterArgAssign := f.NewAssignmentExpression(setterArgSuperAccess, f.NewIdentifier("value"))
	setterArg := f.NewArrowFunction(nil, nil, f.NewNodeList([]*ast.Node{setterArgNameParam, setterArgValueParam}), nil, nil, f.NewToken(ast.KindEqualsGreaterThanToken), setterArgAssign)

	// IIFE: (function(geti, seti) { ... })(getterArg, setterArg)
	iife := f.NewCallExpression(
		f.NewParenthesizedExpression(outerFunc), nil, nil,
		f.NewNodeList([]*ast.Node{getterArg, setterArg}),
		ast.NodeFlagsNone,
	)

	decl := f.NewVariableDeclaration(tx.superIndexBinding, nil, nil, iife)
	declList := f.NewVariableDeclarationList(ast.NodeFlagsConst, f.NewNodeList([]*ast.Node{decl}))
	return f.NewVariableStatement(nil, declList)
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

// isArgumentsIdentifier checks whether an identifier refers to the `arguments` object.
// Since we always assume strict mode, a simple name check suffices.
func isArgumentsIdentifier(node *ast.Node) bool {
	if node.Text() != "arguments" {
		return false
	}
	parent := node.Parent
	if parent == nil {
		return false
	}
	// Exclude property name positions like obj.arguments or { arguments: ... }
	if (ast.IsPropertyAccessExpression(parent) || ast.IsPropertyAssignment(parent)) && parent.Name() == node {
		return false
	}
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
	result := tx.NewTransformer(tx.visit, opts.Context)
	tx.asyncBodyVisitor = tx.EmitContext().NewNodeVisitor(tx.visitAsyncBodyNode)
	return result
}
