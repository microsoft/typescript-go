package estransforms

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/printer"
	"github.com/microsoft/typescript-go/internal/transformers"
)

type classPrivateFieldsTransformer struct {
	transformers.Transformer
	classElementVisitor *ast.NodeVisitor // visits a member of a class
	currentClassElement *ast.ClassElement
}

func newClassPrivateFieldsTransformer(opts *transformers.TransformOptions) *transformers.Transformer {
	tx := &classPrivateFieldsTransformer{}
	tx.classElementVisitor = opts.Context.NewNodeVisitor(tx.visitClassElement)
	return tx.NewTransformer(tx.visit, opts.Context)
}

func (tx *classPrivateFieldsTransformer) visit(node *ast.Node) *ast.Node {
	// !!!
	if node == nil {
		return nil
	}
	if node.SubtreeFacts()&ast.SubtreeContainsClassFields == 0 {
		return node
	}
	switch node.Kind {
	case ast.KindSourceFile:
		return tx.visitSourceFile(node.AsSourceFile())
	case ast.KindPrivateIdentifier:
		return tx.visitPrivateIdentifier(node.AsPrivateIdentifier())
	case ast.KindExpressionStatement:
		return tx.visitExpressionStatement(node.AsExpressionStatement())
	case ast.KindClassDeclaration:
		return tx.visitClassDeclaration(node.AsClassDeclaration())
	case ast.KindPropertyAccessExpression:
		return tx.visitPropertyAccessExpression(node.AsPropertyAccessExpression())
	case ast.KindBinaryExpression:
		return tx.visitBinaryExpression(node.AsBinaryExpression())
	}
	return tx.Visitor().VisitEachChild(node)
}

func (tx *classPrivateFieldsTransformer) visitSourceFile(file *ast.SourceFile) *ast.Node {
	if file.IsDeclarationFile {
		return file.AsNode()
	}

	visited := tx.Visitor().VisitEachChild((file.AsNode()))
	tx.EmitContext().AddEmitHelper(visited.AsNode(), tx.EmitContext().ReadEmitHelpers()...)
	return visited
}

// If we visit a private name, this means it is an undeclared private name.
// Replace it with an empty identifier to indicate a problem with the code,
// unless we are in a statement position - otherwise this will not trigger
// a SyntaxError.
func (tx *classPrivateFieldsTransformer) visitPrivateIdentifier(node *ast.PrivateIdentifier) *ast.Node {
	if ast.IsStatement(node.Parent) {
		return node.AsNode()
	}

	// !!! TODO: private identifier should be replaced in other cases.
	// e.g. constructor { this.#foo = 3 } -> constructor { this. = 3 }
	if !ast.IsPropertyDeclaration(node.Parent) {
		return node.AsNode()
	}
	// !!! TODO: support class expression
	// e.g. array.push(class A { #foo = "hello" })
	classContainer := tx.EmitContext().GetClassContainer()
	if classContainer == nil {
		return node.AsNode()
	}

	emptyIdentifer := tx.Factory().NewIdentifier("")
	tx.EmitContext().SetOriginal(emptyIdentifer, node.AsNode())
	return emptyIdentifer
}

func (tx *classPrivateFieldsTransformer) visitExpressionStatement(node *ast.ExpressionStatement) *ast.Node {
	return tx.Factory().UpdateExpressionStatement(
		node,
		tx.Visitor().VisitNode(node.Expression),
	)
}

func (tx *classPrivateFieldsTransformer) visitClassDeclaration(node *ast.ClassDeclaration) *ast.Node {
	tx.EmitContext().StartClassLexicalEnvironment(node.AsNode())

	members := tx.transformClassMembers(node.AsNode())

	statements := make([]*ast.Statement, 0)

	if len(tx.EmitContext().GetPendingExpressions()) > 0 {
		statements = append(
			statements,
			tx.Factory().NewExpressionStatement(tx.Factory().InlineExpressions(tx.EmitContext().GetPendingExpressions())),
		)
	}
	classDecl := tx.Factory().UpdateClassDeclaration(
		node.AsClassDeclaration(),
		node.Modifiers(),
		node.Name(),
		nil, /*typeParameters*/
		node.HeritageClauses,
		members,
	)
	statements = append([]*ast.Statement{classDecl}, statements...)

	tx.EmitContext().EndClassLexicalEnvironment()

	return transformers.SingleOrMany(statements, tx.Factory())
}

func (tx *classPrivateFieldsTransformer) visitClassElement(node *ast.ClassElement) *ast.Node {
	switch node.Kind {
	case ast.KindConstructor:
		return tx.setCurrentClassElementAnd(node, tx.visitConstructorDeclaration, node)
	// case ast.KindGetAccessor, ast.KindSetAccessor, ast.KindMethodDeclaration:
	// 	return tx.setCurrentClassElementAnd(node, tx.visitMethodOrAccessorDeclaration, node)
	case ast.KindPropertyDeclaration:
		return tx.setCurrentClassElementAnd(node, tx.visitPropertyDeclaration, node)
	// case ast.KindClassStaticBlockDeclaration:
	// 	return tx.setCurrentClassElementAnd(node, tx.visitClassStaticBlockDeclaration, node)
	// case ast.KindComputedPropertyName:
	// 	return tx.visitComputedPropertyName(node)
	// case ast.KindSemicolonClassElement:
	// 	return node
	default:
		return tx.Visitor().Visit(node)
	}
}

func (tx *classPrivateFieldsTransformer) visitConstructorDeclaration(node *ast.Node) *ast.Node {
	classContainer := tx.EmitContext().GetClassContainer()
	return tx.transformConstructor(node.AsConstructorDeclaration(), classContainer).AsNode()
}

func (tx *classPrivateFieldsTransformer) visitPropertyDeclaration(node *ast.Node) *ast.Node {
	// !!! TODO: supports static and auto accessor private fields
	if !ast.IsPrivateIdentifierClassElementDeclaration(node) || ast.IsStatic(node) || ast.IsAutoAccessorPropertyDeclaration(node) {
		return node
	}

	// If we are transforming private elements into WeakMap/WeakSet, we should elide the node.
	info := tx.EmitContext().GetPrivateIdentifierInfo(node.Name().AsPrivateIdentifier())

	if info == nil {
		panic("Undeclared private name for property declaration.")
	}

	// Leave invalid code untransformed
	if !info.IsValid() {
		return node.AsNode()
	}

	return nil
}

func (tx *classPrivateFieldsTransformer) visitBinaryExpression(node *ast.BinaryExpression) *ast.Node {
	// !!! destructuring assignment
	// e.g. ({ x: obj.#x } = ...)

	if ast.IsAssignmentExpression(node.AsNode(), false /*excludeCompoundAssignment*/) {
		// 13.15.2 RS: Evaluation
		//   AssignmentExpression : LeftHandSideExpression `=` AssignmentExpression
		//     1. If |LeftHandSideExpression| is neither an |ObjectLiteral| nor an |ArrayLiteral|, then
		//        a. Let _lref_ be ? Evaluation of |LeftHandSideExpression|.
		//        b. If IsAnonymousFunctionDefinition(|AssignmentExpression|) and IsIdentifierRef of |LeftHandSideExpression| are both *true*, then
		//           i. Let _rval_ be ? NamedEvaluation of |AssignmentExpression| with argument _lref_.[[ReferencedName]].
		//     ...
		//
		//   AssignmentExpression : LeftHandSideExpression `&&=` AssignmentExpression
		//     ...
		//     5. If IsAnonymousFunctionDefinition(|AssignmentExpression|) is *true* and IsIdentifierRef of |LeftHandSideExpression| is *true*, then
		//        a. Let _rval_ be ? NamedEvaluation of |AssignmentExpression| with argument _lref_.[[ReferencedName]].
		//     ...
		//
		//   AssignmentExpression : LeftHandSideExpression `||=` AssignmentExpression
		//     ...
		//     5. If IsAnonymousFunctionDefinition(|AssignmentExpression|) is *true* and IsIdentifierRef of |LeftHandSideExpression| is *true*, then
		//        a. Let _rval_ be ? NamedEvaluation of |AssignmentExpression| with argument _lref_.[[ReferencedName]].
		//     ...
		//
		//   AssignmentExpression : LeftHandSideExpression `??=` AssignmentExpression
		//     ...
		//     4. If IsAnonymousFunctionDefinition(|AssignmentExpression|) is *true* and IsIdentifierRef of |LeftHandSideExpression| is *true*, then
		//        a. Let _rval_ be ? NamedEvaluation of |AssignmentExpression| with argument _lref_.[[ReferencedName]].
		//     ...

		left := ast.SkipOuterExpressions(node.Left, ast.OEKPartiallyEmittedExpressions|ast.OEKParentheses)
		if isPrivateIdentifierPropertyAccessExpression(left) {
			// obj.#x = ...
			info := tx.EmitContext().GetPrivateIdentifierInfo(left.Name().AsPrivateIdentifier())
			if info != nil {
				assignment := tx.NewPrivateIdentifierAssignment(info, left.Expression(), node.Right, node.OperatorToken)
				tx.EmitContext().SetOriginal(assignment, node.AsNode())
				assignment.Loc = node.Loc
				return assignment
			}
		}
	}

	return tx.Visitor().VisitEachChild(node.AsNode())
}

func (tx *classPrivateFieldsTransformer) setCurrentClassElementAnd(classElement *ast.ClassElement, visitor func(arg *ast.Node) *ast.Node, arg *ast.Node) *ast.Node {
	if classElement != tx.currentClassElement {
		savedCurrentClassElement := tx.currentClassElement
		tx.currentClassElement = classElement
		result := visitor(arg)
		tx.currentClassElement = savedCurrentClassElement
		return result
	}
	return visitor(arg)
}

func (tx *classPrivateFieldsTransformer) transformClassMembers(node *ast.ClassLikeDeclaration) *ast.NodeList {
	for _, member := range node.Members() {
		if ast.IsPrivateIdentifierClassElementDeclaration(member) {
			// !!! TODO: supports static and auto accessor private identifier class elements
			if ast.IsPropertyDeclaration(member) && !ast.IsStatic(member) && !ast.IsAutoAccessorPropertyDeclaration(member) {
				tx.EmitContext().AddPrivateIdentifierToEnvironment(member, member.Name().AsPrivateIdentifier())
			}
		}
	}

	members := tx.classElementVisitor.VisitNodes(node.MemberList())

	var syntheticConstructor *ast.ConstructorDeclaration
	if !core.Some(members.Nodes, ast.IsConstructorDeclaration) {
		syntheticConstructor = tx.transformConstructor(nil /*constructor*/, node)
	}

	if syntheticConstructor != nil {
		classThisAssignmentBlock := core.Find(members.Nodes, func(member *ast.Node) bool {
			return isClassThisAssignmentBlock(tx.EmitContext(), member)
		})
		classNamedEvaluationHelperBlock := core.Find(members.Nodes, func(member *ast.Node) bool {
			return isClassNamedEvaluationHelperBlock(tx.EmitContext(), member)
		})
		membersArray := make([]*ast.Node, 0)
		if classThisAssignmentBlock != nil {
			membersArray = append(membersArray, classThisAssignmentBlock)
		}
		if classNamedEvaluationHelperBlock != nil {
			membersArray = append(membersArray, classNamedEvaluationHelperBlock)
		}
		if syntheticConstructor != nil {
			membersArray = append(membersArray, syntheticConstructor.AsNode())
		}
		remainingMembers := members.Nodes
		if classThisAssignmentBlock != nil || classNamedEvaluationHelperBlock != nil {
			remainingMembers = core.Filter(remainingMembers, func(member *ast.Node) bool {
				return member != classThisAssignmentBlock && member != classNamedEvaluationHelperBlock
			})
		}
		membersArray = append(membersArray, remainingMembers...)
		members = tx.Factory().NewNodeList(membersArray)
		members.Loc = node.MemberList().Loc
	}

	return members
}

func (tx *classPrivateFieldsTransformer) transformConstructor(constructor *ast.ConstructorDeclaration, container *ast.ClassLikeDeclaration) *ast.ConstructorDeclaration {
	if constructor != nil {
		constructor = tx.Visitor().VisitNode(constructor.AsNode()).AsConstructorDeclaration()
	}
	if tx.EmitContext().GetClassFacts()&printer.ClassFactsWillHoistInitializersToConstructor == 0 {
		return constructor
	}

	extendsClauseElement := ast.GetEffectiveBaseTypeNode(container)
	isDerivedClass := extendsClauseElement != nil && ast.SkipOuterExpressions(extendsClauseElement.Expression(), ast.OEKAll).Kind != ast.KindNullKeyword
	var parameters *ast.ParameterList
	if constructor == nil {
		parameters = tx.EmitContext().VisitParameters(nil, tx.Visitor())
	} else {
		parameters = tx.EmitContext().VisitParameters(constructor.ParameterList(), tx.Visitor())
	}
	body := tx.transformConstructorBody(container, constructor, isDerivedClass)
	if body == nil {
		return constructor
	}

	if constructor != nil {
		return tx.Factory().UpdateConstructorDeclaration(
			constructor,
			constructor.Modifiers(),    /*modifiers*/
			constructor.TypeParameters, /*typeParameters*/
			parameters,
			constructor.Type,          /*returnType*/
			constructor.FullSignature, /*fullSignature*/
			body,
		).AsConstructorDeclaration()
	}

	result := tx.Factory().NewConstructorDeclaration(nil /*modifiers*/, nil /*typeParameters*/, nil /*parameters*/, nil /*returnType*/, nil /*fullSignature*/, body).AsConstructorDeclaration()
	if container != nil {
		result.Loc = container.Loc
	}
	tx.EmitContext().AddEmitFlags(result.AsNode(), printer.EFStartOnNewLine)

	return result
}

func (tx *classPrivateFieldsTransformer) transformConstructorBody(node *ast.ClassLikeDeclaration, constructor *ast.ConstructorDeclaration, isDerivedClass bool) *ast.Node {
	properties := getPropertiesNeedingInitialization(node)

	// privateMethodsAndAccessors := getPrivateInstanceMethodsAndAccessors(node)
	// !!! needsConstructorBody := len(properties) > 0 || len(privateMethodsAndAccessors) > 0
	needsConstructorBody := len(properties) > 0

	// Only generate synthetic constructor when there are property initializers to move.
	if constructor == nil && !needsConstructorBody {
		return tx.EmitContext().VisitFunctionBody(nil /*node*/, tx.Visitor())
	}

	/// !!! tx.EmitContext().ResumeVariableEnvironment()

	// needsSyntheticConstructor := constructor == nil && isDerivedClass
	// statementOffset := 0
	statements := make([]*ast.Statement, 0)

	// Add the property initializers. Transforms this:
	//
	//  private x = 1;
	//
	// Into this:
	//
	//  constructor() {
	//      this.x = 1;
	//  }
	//

	receiver := tx.Factory().NewThisExpression() // createThis

	// private methods can be called in property initializers, they should execute first.
	// !!!
	// initializerStatements = tx.addInstanceMethodStatements(initializerStatements, privateMethodsAndAccessors, receiver)
	initializerStatements := tx.transformPropertyStatements(properties, receiver)

	statements = append(statements, initializerStatements...)
	if constructor != nil && constructor.Body != nil {
		statements = append(statements, tx.Visitor().VisitNodes(constructor.Body.StatementList()).Nodes...)
	}

	statements = tx.EmitContext().EndAndMergeVariableEnvironment(statements)

	if len(statements) == 0 && constructor == nil {
		return nil
	}

	statementsNodeArray := tx.Factory().NewNodeList(statements)
	if constructor != nil && constructor.Body != nil && len(constructor.Body.Statements()) > 0 {
		statementsNodeArray.Loc = constructor.Body.StatementList().Loc
	} else {
		statementsNodeArray.Loc = node.MemberList().Loc
	}

	multiline := len(statements) > 0
	if constructor != nil && constructor.Body != nil && len(constructor.Body.Statements()) >= len(statements) {
		multiline = constructor.Body.AsBlock().Multiline
	}

	blockNode := tx.Factory().NewBlock(statementsNodeArray, multiline)
	if constructor != nil {
		if constructor.Body != nil {
			blockNode.Loc = constructor.Body.Loc
		} else {
			blockNode.Loc = constructor.Loc
		}
	}
	return blockNode
}

func (tx *classPrivateFieldsTransformer) transformPropertyStatements(properties []*ast.Node, receiver *ast.LeftHandSideExpression) []*ast.Node {
	var statements []*ast.Node
	for _, property := range properties {
		var expression *ast.Node
		if !ast.IsStatic(property) {
			expression = tx.transformProperty(property.AsPropertyDeclaration(), receiver)
		}

		if expression == nil {
			continue
		}

		statement := tx.Factory().NewExpressionStatement(expression)
		tx.EmitContext().SetOriginal(statement, property)
		tx.EmitContext().AddEmitFlags(statement, tx.EmitContext().EmitFlags(property)&printer.EFNoComments)
		tx.EmitContext().SetCommentRange(statement, property.Loc)

		propertyOriginalNode := tx.EmitContext().MostOriginal(property)
		if ast.IsParameter(propertyOriginalNode) {
			// replicate comment and source map behavior from the ts transform for parameter properties.
			tx.EmitContext().SetSourceMapRange(statement, propertyOriginalNode.Loc)
			tx.EmitContext().RemoveAllComments(statement)
		} else {
			tx.EmitContext().SetSourceMapRange(statement, core.NewTextRange(property.Name().Pos(), property.End()))
		}

		// // `setOriginalNode` *copies* the `emitNode` from `property`, so now both
		// // `statement` and `expression` have a copy of the synthesized comments.
		// // Drop the comments from expression to avoid printing them twice.
		// setSyntheticLeadingComments(expression, undefined);
		// setSyntheticTrailingComments(expression, undefined);

		// If the property was originally an auto-accessor, don't emit comments here since they will be attached to
		// the synthezized getter.
		if ast.HasAccessorModifier(propertyOriginalNode) {
			tx.EmitContext().AddEmitFlags(statement, printer.EFNoComments)
		}

		statements = append(statements, statement)
	}

	return statements
}

// Transforms a property initializer into an assignment statement.
func (tx *classPrivateFieldsTransformer) transformProperty(property *ast.PropertyDeclaration, receiver *ast.LeftHandSideExpression) *ast.Node {
	savedCurrentClassElement := tx.currentClassElement
	transformed := tx.transformPropertyWorker(property, receiver)
	if transformed != nil &&
		ast.HasStaticModifier(property.AsNode()) {
		tx.EmitContext().SetOriginal(transformed, property.AsNode())
		tx.EmitContext().SetSourceMapRange(transformed, tx.EmitContext().SourceMapRange(property.Name()))
	}
	tx.currentClassElement = savedCurrentClassElement
	return transformed
}

func (tx *classPrivateFieldsTransformer) transformPropertyWorker(property *ast.PropertyDeclaration, receiver *ast.LeftHandSideExpression) *ast.Expression {
	propertyName := property.Name()
	if ast.HasAccessorModifier(property.AsNode()) {
		propertyName = tx.Factory().NewGeneratedPrivateNameForNode(property.Name())
	} else if ast.IsComputedPropertyName(property.Name()) && transformers.IsSimpleInlineableExpression(property.Name().Expression()) {
		propertyName = tx.Factory().UpdateComputedPropertyName(property.Name().AsComputedPropertyName(), tx.Factory().NewGeneratedNameForNode(property.Name()))
	}

	if ast.IsPrivateIdentifier(propertyName) {
		privateIdentifierInfo := tx.EmitContext().GetPrivateIdentifierInfo(propertyName.AsPrivateIdentifier())
		if privateIdentifierInfo == nil {
			panic("Undeclared private name for property declaration.")
		}

		switch info := privateIdentifierInfo.(type) {
		case *printer.PrivateIdentifierInstanceFieldInfo:
			{
				// `createPrivateInstanceFieldInitializer` in Strada
				initializer := tx.Visitor().VisitNode(property.Initializer)
				if initializer == nil {
					initializer = tx.Factory().NewVoidZeroExpression()
				}
				arguments := []*ast.Expression{
					receiver,
					initializer,
				}
				return tx.Factory().NewCallExpression(
					tx.Factory().NewPropertyAccessExpression(
						info.BrandCheckIdentifier().AsNode(),
						nil, /*questionDotToken*/
						tx.Factory().NewIdentifier("set"),
						ast.NodeFlagsNone,
					),
					nil, /*questionDotToken*/
					nil, /*typeArguments*/
					tx.Factory().NewNodeList(arguments),
					ast.NodeFlagsNone,
				)
			}
		}
	}
	return nil
}

func (tx *classPrivateFieldsTransformer) visitPropertyAccessExpression(node *ast.PropertyAccessExpression) *ast.Node {
	if ast.IsPrivateIdentifier(node.Name()) {
		privateIdentifierInfo := tx.EmitContext().GetPrivateIdentifierInfo(node.Name().AsPrivateIdentifier())
		if privateIdentifierInfo != nil {
			privateIdentifierAccess := tx.NewPrivateIdentifierAccess(privateIdentifierInfo, node.Expression)
			tx.EmitContext().SetOriginal(
				privateIdentifierAccess,
				node.AsNode(),
			)
			privateIdentifierAccess.Loc = node.Loc
			return privateIdentifierAccess
		}
	}
	return tx.Visitor().VisitEachChild(node.AsNode())
}

func (tx *classPrivateFieldsTransformer) NewPrivateIdentifierAccess(info printer.PrivateIdentifierInfo, receiver *ast.Expression) *ast.Expression {
	receiver = tx.Visitor().VisitNode(receiver)
	return tx.NewPrivateIdentifierAccessHelper(info, receiver)
}

func (tx *classPrivateFieldsTransformer) NewPrivateIdentifierAccessHelper(info printer.PrivateIdentifierInfo, receiver *ast.Expression) *ast.Expression {
	tx.EmitContext().SetCommentRange(receiver, core.NewTextRange(-1, receiver.End()))

	switch v := info.(type) {
	case *printer.PrivateIdentifierAccessorInfo:
		return tx.Factory().NewClassPrivateFieldGetHelper(
			receiver,
			v.BrandCheckIdentifier(),
			v.Kind(),
			v.GetterName,
		)
	case *printer.PrivateIdentifierMethodInfo:
		return tx.Factory().NewClassPrivateFieldGetHelper(
			receiver,
			v.BrandCheckIdentifier(),
			v.Kind(),
			v.MethodName,
		)
	case *printer.PrivateIdentifierInstanceFieldInfo:
		return tx.Factory().NewClassPrivateFieldGetHelper(
			receiver,
			v.BrandCheckIdentifier(),
			v.Kind(),
			nil, /*f*/
		)
	case *printer.PrivateIdentifierStaticFieldInfo:
		return tx.Factory().NewClassPrivateFieldGetHelper(
			receiver,
			v.BrandCheckIdentifier(),
			v.Kind(),
			v.VariableName,
		)
	case *printer.PrivateIdentifierUntransformedInfo:
		panic("Access helpers should not be created for untransformed private elements")
	default:
		panic("Unknown private identifier info")
	}
}

func (tx *classPrivateFieldsTransformer) NewPrivateIdentifierAssignment(info printer.PrivateIdentifierInfo, receiver *ast.Expression, right *ast.Expression, operatorToken *ast.TokenNode) *ast.Expression {
	receiver = tx.Visitor().VisitNode(receiver)
	right = tx.Visitor().VisitNode(right)

	tx.EmitContext().SetCommentRange(receiver, core.NewTextRange(-1, receiver.End()))

	// e.g. #field += 1
	if checker.IsCompoundAssignment(operatorToken.Kind) {
		readExpression, initializeExpression := tx.NewCopiableReceiverExpr(receiver)
		if initializeExpression != nil {
			receiver = initializeExpression
		} else {
			receiver = readExpression
		}
		right = tx.Factory().NewBinaryExpression(
			nil, /*modifiers*/
			tx.NewPrivateIdentifierAccessHelper(info, readExpression),
			nil, /*typeNode*/
			getNonAssignmentOperatorForCompoundAssignment(tx.EmitContext(), operatorToken),
			right,
		)
	}

	switch v := info.(type) {
	case *printer.PrivateIdentifierAccessorInfo:
		return tx.Factory().NewClassPrivateFieldSetHelper(
			receiver,
			v.BrandCheckIdentifier(),
			right,
			v.Kind(),
			v.SetterName,
		)
	case *printer.PrivateIdentifierMethodInfo:
		return tx.Factory().NewClassPrivateFieldSetHelper(
			receiver,
			v.BrandCheckIdentifier(),
			right,
			v.Kind(),
			nil, /*f*/
		)
	case *printer.PrivateIdentifierInstanceFieldInfo:
		return tx.Factory().NewClassPrivateFieldSetHelper(
			receiver,
			v.BrandCheckIdentifier(),
			right,
			v.Kind(),
			nil, /*f*/
		)
	case *printer.PrivateIdentifierStaticFieldInfo:
		return tx.Factory().NewClassPrivateFieldSetHelper(
			receiver,
			v.BrandCheckIdentifier(),
			right,
			v.Kind(),
			v.VariableName,
		)
	case *printer.PrivateIdentifierUntransformedInfo:
		panic("Access helpers should not be created for untransformed private elements")
	default:
		panic("Unknown private element type")
	}
}

func (tx *classPrivateFieldsTransformer) NewCopiableReceiverExpr(receiver *ast.Expression) (*ast.Expression, *ast.Expression) {
	var clone *ast.Expression
	if ast.NodeIsSynthesized(receiver) {
		clone = receiver
	} else {
		clone = tx.Factory().DeepCloneNode(receiver)
	}
	if transformers.IsSimpleInlineableExpression(receiver) {
		return clone, nil
	}
	readExpression := tx.Factory().NewTempVariable()
	initializeExpression := tx.Factory().NewAssignmentExpression(readExpression, clone)
	return readExpression, initializeExpression
}

func isPrivateIdentifierPropertyAccessExpression(node *ast.Node) bool {
	return ast.IsPropertyAccessExpression(node) && ast.IsPrivateIdentifier(node.Name())
}

func getPropertiesNeedingInitialization(node *ast.ClassLikeDeclaration) []*ast.Node {
	var properties []*ast.Node

	for _, member := range node.Members() {
		if !ast.IsPropertyDeclaration(member) || ast.IsStatic(member) {
			continue
		}

		if member.Initializer() != nil ||
			ast.IsPrivateIdentifier(member.Name()) ||
			ast.HasAccessorModifier(member) {
			properties = append(properties, member)
		}
	}

	return properties
}
