package estransforms

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/printer"
	"github.com/microsoft/typescript-go/internal/scanner"
	"github.com/microsoft/typescript-go/internal/transformers"
)

type ClassFacts int32

const (
	ClassFactsClassWasDecorated = 1 << iota
	ClassFactsNeedsClassConstructorReference
	ClassFactsNeedsClassSuperReference
	ClassFactsNeedsSubstitutionForThisInClassStaticField
	ClassFactsWillHoistInitializersToConstructor
	ClassFactsNone ClassFacts = 0
)

type PrivateEnvironment struct {
	// used for prefixing generated variable names
	ClassName *ast.Node
	// used for brand check on private methods
	WeakSetName *ast.Node

	// A mapping of generated private names to information needed for transformation.
	GeneratedIdentifiers map[*ast.Node]printer.PrivateIdentifierInfo
	// A mapping of private names to information needed for transformation.
	Identifiers map[string]printer.PrivateIdentifierInfo
}

type classLexicalEnvrionment struct {
	facts ClassFacts
	// Used for brand checks on static members, and `this` references in static initializers
	classConstructor *ast.IdentifierNode
	classThis        *ast.IdentifierNode
	// Used for `super` references in static initializers.
	superClassReference *ast.IdentifierNode

	privateEnvrionment *PrivateEnvironment
}

type ClassPropertySubstitutionFlags int

const (
	// Enables substitutions for class expressions with static fields which have initializers that reference the class name.
	ClassPropertySubstitutionFlagsClassAliases = 1 << iota
	// Enables substitutions for class expressions with static fields which have initializers that reference the 'this' or 'super'.
	ClassPropertySubstitutionFlagsClassStaticThisOrSuperReference
	ClassPropertySubstitutionFlagsNone = 0
)

type classFieldsTransformer struct {
	transformers.Transformer
	discardedValueVisitor                               *ast.NodeVisitor // visits expressions whose values would be discarded at runtime
	classDeclarationInNewClassLexicalEnvironmentVisitor *ast.NodeVisitor // visits class declaration in new class lexical environment
	classExpressionInNewClassLexicalEnvironmentVisitor  *ast.NodeVisitor // visits class expression in new class lexical environment
	classElementVisitor                                 *ast.NodeVisitor // visits a member of a class

	parentNode  *ast.Node // used for ancestor tracking via pushNode/popNode to detect expression identifiers
	currentNode *ast.Node // used for ancestor tracking via pushNode/popNode to detect expression identifiers

	currentClassContainer *ast.ClassLikeDeclaration
	currentClassElement   *ast.ClassElement

	classScopeStack core.Stack[*classLexicalEnvrionment]

	enabledSubstitutions ClassPropertySubstitutionFlags
	classAlias           map[*ast.Node]*ast.Node
	// Tracks what computed name expressions originating from elided names must be inlined
	// at the next execution site, in document order
	pendingExpressions []*ast.Expression
}

func newClassFieldsTransformer(emitContext *printer.EmitContext) *transformers.Transformer {
	tx := &classFieldsTransformer{}
	tx.discardedValueVisitor = emitContext.NewNodeVisitor(tx.visitDiscardedValue)
	tx.classDeclarationInNewClassLexicalEnvironmentVisitor = emitContext.NewNodeVisitor(tx.visitClassDeclarationInNewClassLexicalEnvironment)
	tx.classExpressionInNewClassLexicalEnvironmentVisitor = emitContext.NewNodeVisitor(tx.visitClassExpressionInNewClassLexicalEnvironment)
	tx.classElementVisitor = emitContext.NewNodeVisitor(tx.visitClassElement)
	tx.enabledSubstitutions = ClassPropertySubstitutionFlagsNone
	return tx.NewTransformer(tx.visit, emitContext)
}

func (tx *classFieldsTransformer) GetPendingExpressions() []*ast.Expression {
	if tx.pendingExpressions == nil {
		tx.pendingExpressions = []*ast.Expression{}
	}
	return tx.pendingExpressions
}

func (tx *classFieldsTransformer) SetPendingExpressions(expressions []*ast.Expression) {
	tx.pendingExpressions = expressions
}

func (tx *classFieldsTransformer) AppendPendingExpressions(expressions ...*ast.Expression) {
	tx.pendingExpressions = append(tx.GetPendingExpressions(), expressions...)
}

func (tx *classFieldsTransformer) PrependPendingExpressions(expressions ...*ast.Expression) {
	tx.pendingExpressions = append(expressions, tx.GetPendingExpressions()...)
}

func (tx *classFieldsTransformer) visit(node *ast.Node) *ast.Node {
	if node == nil {
		return nil
	}
	if node.SubtreeFacts()&(ast.SubtreeContainsClassFields|ast.SubtreeContainsLexicalThisOrSuper) == 0 {
		return node
	}
	switch node.Kind {
	case ast.KindSourceFile:
		return tx.visitSourceFile(node.AsSourceFile())
	case ast.KindClassDeclaration:
		return tx.visitClassDeclartion(node.AsClassDeclaration())
	case ast.KindClassExpression:
		return tx.visitClassExpression(node.AsClassExpression())
	case ast.KindClassStaticBlockDeclaration, ast.KindPropertyDeclaration:
		panic("Use `classElementVisitor` instead.")
	case ast.KindPropertyAssignment:
		return tx.visitPropertyAssignment(node.AsPropertyAssignment())
	case ast.KindVariableStatement:
		return tx.visitVariableStatement(node.AsVariableStatement())
	case ast.KindVariableDeclaration:
		return tx.visitVariableDeclaration(node.AsVariableDeclaration())
	case ast.KindParameter:
		return tx.visitParameterDeclaration(node.AsParameterDeclaration())
	case ast.KindBindingElement:
		return tx.visitBindingElement(node.AsBindingElement())
	case ast.KindExportAssignment:
		return tx.visitExportAssignment(node.AsExportAssignment())
	case ast.KindPrivateIdentifier:
		return tx.visitPrivateIdentifier(node.AsPrivateIdentifier())
	case ast.KindPropertyAccessExpression:
		return tx.visitPropertyAccessExpression(node.AsPropertyAccessExpression())
	case ast.KindElementAccessExpression:
		return tx.visitElementAccessExpression(node.AsElementAccessExpression())
	case ast.KindPrefixUnaryExpression:
		return tx.visitPrefixUnaryExpression(node.AsPrefixUnaryExpression(), false /*resultIsDiscarded*/)
	case ast.KindPostfixUnaryExpression:
		return tx.visitPostfixUnaryExpression(node.AsPostfixUnaryExpression(), false /*resultIsDiscarded*/)
	case ast.KindBinaryExpression:
		return tx.visitBinaryExpression(node.AsBinaryExpression(), false /*resultIsDiscarded*/)
	case ast.KindParenthesizedExpression:
		return tx.visitParenthesizedExpression(node.AsParenthesizedExpression(), false /*resultIsDiscarded*/)
	case ast.KindCallExpression:
		return tx.visitCallExpression(node.AsCallExpression())
	case ast.KindExpressionStatement:
		return tx.visitExpressionStatement(node.AsExpressionStatement())
	case ast.KindTaggedTemplateExpression:
		return tx.visitTaggedTemplateExpression(node.AsTaggedTemplateExpression())
	case ast.KindForStatement:
		return tx.visitForStatement(node.AsForStatement())
	case ast.KindThisKeyword:
		return tx.visitThisExpression(node.AsKeywordExpression())
		// !!! other kinds
	case ast.KindFunctionDeclaration, ast.KindFunctionExpression:
		// If we are descending into a new scope, clear the current class element
		return tx.setCurrentClassElementAnd(
			nil, /*classElement*/
			tx.Visitor().VisitEachChild,
			node,
		)
	case ast.KindConstructor, ast.KindMethodDeclaration, ast.KindGetAccessor, ast.KindSetAccessor:
		// If we are descending into a class element, set the class element
		return tx.setCurrentClassElementAnd(
			node,
			tx.Visitor().VisitEachChild,
			node,
		)
	}
	return tx.Visitor().VisitEachChild(node)
}

func (tx *classFieldsTransformer) visitSourceFile(file *ast.SourceFile) *ast.Node {
	if file.IsDeclarationFile {
		return file.AsNode()
	}

	visited := tx.Visitor().VisitEachChild((file.AsNode()))
	tx.EmitContext().AddEmitHelper(visited.AsNode(), tx.EmitContext().ReadEmitHelpers()...)
	return visited
}

func (tx *classFieldsTransformer) visitClassDeclartion(node *ast.ClassDeclaration) *ast.Node {
	return tx.visitInNewClassLexicalEnviornment(node.AsNode(), func(n *ast.Node) *ast.Node {
		return tx.classDeclarationInNewClassLexicalEnvironmentVisitor.Visit(n)
	})
}

func (tx *classFieldsTransformer) visitClassExpression(node *ast.ClassExpression) *ast.Node {
	return tx.visitInNewClassLexicalEnviornment(node.AsNode(), func(n *ast.Node) *ast.Node {
		return tx.classExpressionInNewClassLexicalEnvironmentVisitor.Visit(n)
	})
}

func (tx *classFieldsTransformer) visitPropertyAssignment(node *ast.PropertyAssignment) *ast.Node {
	// !!!
	return node.AsNode()
}

func (tx *classFieldsTransformer) visitVariableStatement(node *ast.VariableStatement) *ast.Node {
	// !!!
	return node.AsNode()
}

func (tx *classFieldsTransformer) visitVariableDeclaration(node *ast.VariableDeclaration) *ast.Node {
	// !!!
	return node.AsNode()
}

func (tx *classFieldsTransformer) visitParameterDeclaration(node *ast.ParameterDeclaration) *ast.Node {
	// !!!
	return node.AsNode()
}

func (tx *classFieldsTransformer) visitBindingElement(node *ast.BindingElement) *ast.Node {
	// !!!
	return node.AsNode()
}

func (tx *classFieldsTransformer) visitExportAssignment(node *ast.ExportAssignment) *ast.Node {
	// !!!
	return node.AsNode()
}

// If we visit a private name, this means it is an undeclared private name.
// Replace it with an empty identifier to indicate a problem with the code,
// unless we are in a statement position - otherwise this will not trigger
// a SyntaxError.
func (tx *classFieldsTransformer) visitPrivateIdentifier(node *ast.PrivateIdentifier) *ast.Node {
	// !!!
	if ast.IsStatement(node.Parent) {
		return node.AsNode()
	}

	// !!! replace it with an empty identifier
	return node.AsNode()
}

func (tx *classFieldsTransformer) visitPropertyAccessExpression(node *ast.PropertyAccessExpression) *ast.Node {
	if ast.IsPrivateIdentifier(node.Name()) {
		privateIdentifierInfo := tx.accessPrivateIdentifier(node.Name().AsPrivateIdentifier())
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
	// !!!
	return tx.Visitor().VisitEachChild(node.AsNode())
}

func (tx *classFieldsTransformer) visitElementAccessExpression(node *ast.ElementAccessExpression) *ast.Node {
	// !!!
	return node.AsNode()
}

func (tx *classFieldsTransformer) visitPrefixUnaryExpression(node *ast.PrefixUnaryExpression, resultIsDiscarded bool) *ast.Node {
	// !!!
	return node.AsNode()
}

func (tx *classFieldsTransformer) visitPostfixUnaryExpression(node *ast.PostfixUnaryExpression, resultIsDiscarded bool) *ast.Node {
	// !!!
	return node.AsNode()
}

func (tx *classFieldsTransformer) visitBinaryExpression(node *ast.BinaryExpression, resultIsDiscarded bool) *ast.Node {
	// !!!
	return node.AsNode()
}

func (tx *classFieldsTransformer) visitParenthesizedExpression(node *ast.ParenthesizedExpression, resultIsDiscarded bool) *ast.Node {
	// !!!
	return node.AsNode()
}

func (tx *classFieldsTransformer) visitCallExpression(node *ast.CallExpression) *ast.Node {
	// !!!
	return node.AsNode()
}

func (tx *classFieldsTransformer) visitExpressionStatement(node *ast.ExpressionStatement) *ast.Node {
	return tx.Factory().UpdateExpressionStatement(node, tx.discardedValueVisitor.VisitNode(node.Expression))
}

func (tx *classFieldsTransformer) visitTaggedTemplateExpression(node *ast.TaggedTemplateExpression) *ast.Node {
	// !!!
	return node.AsNode()
}

func (tx *classFieldsTransformer) visitForStatement(node *ast.ForStatement) *ast.Node {
	// !!!
	return node.AsNode()
}

func (tx *classFieldsTransformer) visitThisExpression(node *ast.KeywordExpression) *ast.Node {
	// !!!
	return node.AsNode()
}

// Pushes a new child node onto the ancestor tracking stack, returning the grandparent node to be restored later via `popNode`.
func (tx *classFieldsTransformer) pushNode(node *ast.Node) (grandparentNode *ast.Node) {
	grandparentNode = tx.parentNode
	tx.parentNode = tx.currentNode
	tx.currentNode = node
	return
}

// Pops the last child node off the ancestor tracking stack, restoring the grandparent node.
func (tx *classFieldsTransformer) popNode(grandparentNode *ast.Node) {
	tx.currentNode = tx.parentNode
	tx.parentNode = grandparentNode
}

func (tx *classFieldsTransformer) visitDiscardedValue(node *ast.Node) *ast.Node {
	grandparentNode := tx.pushNode(node)
	defer tx.popNode(grandparentNode)

	return tx.visitNoStack(node, true /*resultIsDiscarded*/)
}

func (tx *classFieldsTransformer) visitNoStack(node *ast.Node, resultIsDiscarded bool) *ast.Node {
	switch node.Kind {
	case ast.KindSourceFile:
		node = tx.visitSourceFile(node.AsSourceFile())
	case ast.KindPrefixUnaryExpression:
		node = tx.visitPrefixUnaryExpression(node.AsPrefixUnaryExpression(), resultIsDiscarded)
	case ast.KindPostfixUnaryExpression:
		node = tx.visitPostfixUnaryExpression(node.AsPostfixUnaryExpression(), resultIsDiscarded)
	case ast.KindBinaryExpression:
		node = tx.visitBinaryExpression(node.AsBinaryExpression(), resultIsDiscarded)
	// !!!
	// case ast.KindCommaListExpression:
	// 	node = tx.visitCommaListExpression(node.AsCommaListExpression(), resultIsDiscarded)
	case ast.KindParenthesizedExpression:
		node = tx.visitParenthesizedExpression(node.AsParenthesizedExpression(), resultIsDiscarded)
	default:
		node = tx.Visitor().Visit(node)
	}

	return node
}

func (tx *classFieldsTransformer) visitInNewClassLexicalEnviornment(node *ast.Node, visitor func(node *ast.Node) *ast.Node) *ast.Node {
	savedCurrentClassContainer := tx.currentClassContainer
	savedPendingExpressions := tx.GetPendingExpressions()
	tx.currentClassContainer = node
	tx.SetPendingExpressions(nil)
	tx.startClassLexicalEnvironment()

	name := ast.GetNameOfDeclaration(node.AsNode())
	if name != nil && ast.IsIdentifier(name) {
		tx.getPrivateEnvironment().ClassName = name
	} else {
		assignedName := tx.EmitContext().AssignedName(node)
		if assignedName != nil {
			if ast.IsStringLiteral(assignedName) {
				// If the class name was assigned from a string literal based on an Identifier, use the Identifier
				// as the prefix.
				if textSourceNode := tx.EmitContext().TextSource(assignedName); textSourceNode != nil && ast.IsIdentifier(textSourceNode) {
					tx.getPrivateEnvironment().ClassName = textSourceNode
				} else if
				// If the class name was assigned from a string literal that is a valid identifier, create an
				// identifier from it.
				scanner.IsIdentifierText(assignedName.Text(), core.LanguageVariantStandard) {
					prefixName := tx.Factory().NewIdentifier(assignedName.Text())
					tx.getPrivateEnvironment().ClassName = prefixName
				}
			}
		}
	}

	privateInstanceMethodsAndAccessors := getPrivateInstanceMethodsAndAccessors(node.AsNode())
	if len(privateInstanceMethodsAndAccessors) > 0 {
		// !!! tx.getPrivateEnvironment().WeakSetName = tx.NewHoistedVariableForClass("instances", privateInstanceMethodsAndAccessors[0])
		tx.getPrivateEnvironment().WeakSetName = tx.Factory().NewUniqueName("instance")
	}

	facts := getClassFacts(tx.EmitContext(), node.AsNode())
	tx.getClassLexicalEnvironment().facts = facts

	if facts&ClassFactsNeedsSubstitutionForThisInClassStaticField != 0 {
		tx.enableSubstitutionForClassStaticThisOrSuperReference()
	}

	result := visitor(node)

	tx.endClassLexicalEnvironment()
	tx.currentClassContainer = savedCurrentClassContainer
	tx.SetPendingExpressions(savedPendingExpressions)

	return result
}

func (tx *classFieldsTransformer) visitClassDeclarationInNewClassLexicalEnvironment(node *ast.Node) *ast.Node {
	facts := tx.getClassLexicalEnvironment().facts
	// If a class has private static fields, or a static field has a `this` or `super` reference,
	// then we need to allocate a temp variable to hold on to that reference.
	var pendingClassReferenceAssignment *ast.Expression
	if facts&ClassFactsNeedsClassConstructorReference != 0 {
		// If we aren't transforming class static blocks, then we can't reuse `_classThis` since in
		// `class C { ... static { _classThis = ... } }; _classThis = C` the outer assignment would occur *after*
		// class static blocks evaluate and would overwrite the replacement constructor produced by class
		// decorators.

		// If we are transforming class static blocks, then we can reuse `_classThis` since the assignment
		// will be evaluated *before* the transformed static blocks are evaluated and thus won't overwrite
		// the replacement constructor.
		classThis := tx.EmitContext().ClassThis(node)
		if classThis != nil {
			tx.getClassLexicalEnvironment().classConstructor = classThis
			pendingClassReferenceAssignment = tx.Factory().NewAssignmentExpression(classThis, tx.Factory().GetInternalName(node))
		} else {
			temp := tx.Factory().NewTempVariableEx(printer.AutoGenerateOptions{
				Flags: printer.GeneratedIdentifierFlagsReservedInNestedScopes,
			})
			tx.EmitContext().AddVariableDeclaration(temp)
			tx.getClassLexicalEnvironment().classConstructor = tx.Factory().DeepCloneNode(temp)
			pendingClassReferenceAssignment = tx.Factory().NewAssignmentExpression(temp, tx.Factory().GetInternalName(node))
		}
	}

	if tx.EmitContext().ClassThis(node) != nil {
		tx.getClassLexicalEnvironment().classThis = tx.EmitContext().ClassThis(node)
	}

	// !!! const isClassWithConstructorReference = tx.EmitResolver().HasNodeCheckFlag(node, NodeCheckFlagsisClassWithConstructorReference)
	isClassWithConstructorReference := false
	isExport := ast.HasSyntacticModifier(node, ast.ModifierFlagsExport)
	isDefault := ast.HasSyntacticModifier(node, ast.ModifierFlagsDefault)
	modifiers := transformers.ExtractModifiers(tx.EmitContext(), node.Modifiers(), ast.ModifierFlagsAll)
	// !!! const heritageClauses = visitNodes(node.heritageClauses, heritageClauseVisitor, isHeritageClause);
	var heritageClauses *ast.HeritageClauseList
	members, prologue := tx.transformClassMembers(node)

	var statements []*ast.Statement = make([]*ast.Statement, 0)
	if pendingClassReferenceAssignment != nil {
		tx.PrependPendingExpressions(pendingClassReferenceAssignment)
	}

	// Write any pending expressions from elided or moved computed property names
	if len(tx.GetPendingExpressions()) > 0 {
		statements = append(statements, tx.Factory().NewExpressionStatement(tx.Factory().InlineExpressions(tx.GetPendingExpressions())))
	}

	// !!! if (shouldTransformInitializersUsingSet || shouldTransformPrivateElementsOrClassStaticBlocks || getInternalEmitFlags(node) & InternalEmitFlags.TransformPrivateStaticElements) {
	// Emit static property assignment. Because classDeclaration is lexically evaluated,
	// it is safe to emit static property assignment after classDeclaration
	// From ES6 specification:
	//      HasLexicalDeclaration (N) : Determines if the argument identifier has a binding in this environment record that was created using
	//                                  a lexical declaration such as a LexicalDeclaration or a ClassDeclaration.
	staticProperties := ast.GetStaticPropertiesAndClassStaticBlock(node)
	if len(staticProperties) > 0 {
		statements = tx.addPropertyOrClassStaticBlockStatements(statements, staticProperties, tx.Factory().GetInternalName(node))
	}

	if len(statements) > 0 && isExport && isDefault {
		// exportOrDefaultModifiers := transformers.ExtractModifiers(tx.EmitContext(), node.Modifiers(), ast.ModifierFlagsExportDefault)
		statements = append(statements, tx.Factory().NewExportAssignment(
			nil,   /*modifiers*/
			false, /*isExportEquals*/
			nil,   /*typeNode*/
			tx.Factory().GetLocalNameEx(node, printer.AssignedNameOptions{
				AllowComments:   false,
				AllowSourceMaps: true,
			})))
	}

	alias := tx.getClassLexicalEnvironment().classConstructor
	if isClassWithConstructorReference && alias != nil {
		tx.enableSubstitutionForClassAliases()
		tx.setClassAlias(node, alias)
	}

	classDecl := tx.Factory().UpdateClassDeclaration(
		node.AsClassDeclaration(),
		modifiers,
		node.Name(),
		nil, /*typeParameters*/
		heritageClauses,
		members,
	)
	statements = append([]*ast.Statement{classDecl}, statements...)

	if prologue != nil {
		statements = append([]*ast.Statement{
			tx.Factory().NewExpressionStatement(prologue),
		}, statements...)
	}

	return transformers.SingleOrMany(statements, tx.Factory())
}

func (tx *classFieldsTransformer) visitClassExpressionInNewClassLexicalEnvironment(node *ast.Node) *ast.Node {
	// !!!
	return node
}

func (tx *classFieldsTransformer) visitClassElement(node *ast.ClassElement) *ast.Node {
	switch node.Kind {
	case ast.KindConstructor:
		return tx.setCurrentClassElementAnd(node, tx.visitConstructorDeclaration, node)
	case ast.KindGetAccessor, ast.KindSetAccessor, ast.KindMethodDeclaration:
		return tx.setCurrentClassElementAnd(node, tx.visitMethodOrAccessorDeclaration, node)
	case ast.KindPropertyDeclaration:
		return tx.setCurrentClassElementAnd(node, tx.visitPropertyDeclaration, node)
	case ast.KindClassStaticBlockDeclaration:
		return tx.setCurrentClassElementAnd(node, tx.visitClassStaticBlockDeclaration, node)
	case ast.KindComputedPropertyName:
		return tx.visitComputedPropertyName(node)
	case ast.KindSemicolonClassElement:
		return node
	default:
		if ast.IsModifierLike(node) {
			return tx.Factory().NewSyntaxList(
				transformers.ExtractModifiers(tx.EmitContext(), node.Modifiers(), ast.ModifierFlagsAll).Nodes,
			)
		}
		return tx.Visitor().Visit(node)
	}
}

func (tx *classFieldsTransformer) visitConstructorDeclaration(node *ast.Node) *ast.Node {
	if tx.currentClassContainer != nil {
		return tx.transformConstructor(node.AsConstructorDeclaration(), tx.currentClassContainer).AsNode()
	}
	return tx.Visitor().VisitEachChild(node)
}

func (tx *classFieldsTransformer) visitMethodOrAccessorDeclaration(node *ast.Node) *ast.Node {
	if !ast.IsPrivateIdentifierClassElementDeclaration(node) || !tx.shouldTransformClassElementToWeakMap(node) {
		return tx.classElementVisitor.VisitEachChild(node)
	}

	// !!!
	return tx.classElementVisitor.VisitEachChild(node)
}

func (tx *classFieldsTransformer) visitPropertyDeclaration(node *ast.Node) *ast.Node {
	// !!!
	// If this is an auto-accessor, we defer to `transformAutoAccessor`. That function
	// will in turn call `transformFieldInitializer` as needed.
	return tx.transformFieldInitializer(node.AsPropertyDeclaration())
}

func (tx *classFieldsTransformer) visitClassStaticBlockDeclaration(node *ast.Node) *ast.Node {
	// !!!
	tx.EmitContext().StartVariableEnvironment()
	return tx.Visitor().VisitEachChild(node)
}

func (tx *classFieldsTransformer) visitComputedPropertyName(node *ast.Node) *ast.Node {
	expression := node.Expression()
	return tx.Factory().UpdateComputedPropertyName(node.AsComputedPropertyName(), tx.injectPendingExpressions(expression))
}

func (tx *classFieldsTransformer) injectPendingExpressions(expression *ast.Node) *ast.Node {
	if len(tx.GetPendingExpressions()) > 0 {
		if ast.IsParenthesizedExpression(expression) {
			tx.AppendPendingExpressions(expression.Expression())
			expression = tx.Factory().UpdateParenthesizedExpression(expression.AsParenthesizedExpression(), tx.Factory().InlineExpressions(tx.pendingExpressions))
		} else {
			tx.AppendPendingExpressions(expression)
			expression = tx.Factory().InlineExpressions(tx.pendingExpressions)
		}
		tx.SetPendingExpressions(nil)
	}
	return expression
}

func (tx *classFieldsTransformer) startClassLexicalEnvironment() {
	tx.classScopeStack.Push(&classLexicalEnvrionment{})
}

func (tx *classFieldsTransformer) endClassLexicalEnvironment() {
	// !!!
	tx.classScopeStack.Pop()
}

func (tx *classFieldsTransformer) getClassLexicalEnvironment() *classLexicalEnvrionment {
	return tx.classScopeStack.Peek()
}

func (tx *classFieldsTransformer) getPrivateEnvironment() *PrivateEnvironment {
	classScope := tx.classScopeStack.Peek()
	if classScope.privateEnvrionment == nil {
		classScope.privateEnvrionment = &PrivateEnvironment{}
	}
	return classScope.privateEnvrionment
}

func (tx *classFieldsTransformer) transformClassMembers(node *ast.ClassLikeDeclaration) (*ast.NodeList, *ast.Expression) {

	// !!! Declare private name
	// !!! if (shouldTransformPrivateElementsOrClassStaticBlocks || shouldTransformPrivateStaticElementsInFile) {
	for _, member := range node.Members() {
		if ast.IsPrivateIdentifierClassElementDeclaration(member) {
			// !!! shouldTransformClassElementToWeakMap?
			if tx.shouldTransformClassElementToWeakMap(node) {
				tx.addPrivateIdentifierToEnvironment(member, member.Name().AsPrivateIdentifier())
			} else {
				// !!!
			}
		}
	}

	members := tx.classElementVisitor.VisitNodes(node.MemberList())

	var syntheticConstructor *ast.ConstructorDeclaration
	if !core.Some(members.Nodes, ast.IsConstructorDeclaration) {
		syntheticConstructor = tx.transformConstructor(nil /*constructor*/, node)
	}

	var prologue *ast.Expression

	// !!!
	// If there are pending expressions create a class static block in which to evaluate them, but only if
	// class static blocks are not also being transformed. This block will be injected at the top of the class
	// to ensure that expressions from computed property names are evaluated before any other static
	// initializers.
	var syntheticStaticBlock *ast.ClassStaticBlockDeclaration

	if syntheticConstructor != nil || syntheticStaticBlock != nil {
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
		if syntheticStaticBlock != nil {
			membersArray = append(membersArray, syntheticStaticBlock.AsNode())
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

	return members, prologue
}

func (tx *classFieldsTransformer) transformConstructor(constructor *ast.ConstructorDeclaration, container *ast.ClassLikeDeclaration) *ast.ConstructorDeclaration {
	if constructor != nil {
		constructor = tx.Visitor().VisitNode(constructor.AsNode()).AsConstructorDeclaration()
	}
	if tx.getClassLexicalEnvironment().facts&ClassFactsWillHoistInitializersToConstructor == 0 {
		return constructor
	}

	extendsClauseElement := ast.GetEffectiveBaseTypeNode(container)
	isDerivedClass := extendsClauseElement != nil && ast.SkipOuterExpressions(extendsClauseElement.Expression(), ast.OEKAll).Kind != ast.KindNullKeyword
	body := tx.transformConstructorBody(container, constructor, isDerivedClass)
	if body == nil {
		return constructor
	}

	var parameters *ast.ParameterList
	if constructor != nil {
		parameters = tx.EmitContext().VisitParameters(constructor.ParameterList(), tx.Visitor())
	}

	// !!! handled by runtimesyntax
	if constructor != nil {
		return constructor
		// return tx.Factory().UpdateConstructorDeclaration(
		// 	constructor,
		// 	nil,                        /*modifiers*/
		// 	constructor.TypeParameters, /*typeParameters*/
		// 	parameters,
		// 	constructor.Type,          /*returnType*/
		// 	constructor.FullSignature, /*fullSignature*/
		// 	body,
		// ).AsConstructorDeclaration()
	}

	result := tx.Factory().NewConstructorDeclaration(nil /*modifiers*/, nil /*typeParameters*/, parameters, nil /*returnType*/, nil /*fullSignature*/, body).AsConstructorDeclaration()
	if container != nil {
		result.Loc = container.Loc
	}
	tx.EmitContext().AddEmitFlags(result.AsNode(), printer.EFStartOnNewLine)

	return result
}

func (tx *classFieldsTransformer) transformConstructorBody(node *ast.ClassLikeDeclaration, constructor *ast.ConstructorDeclaration, isDerivedClass bool) *ast.Node {
	instanceProperties := getProperties(node, false /*requireInitializer*/, false /*isStatic*/)
	properties := instanceProperties
	// !!! if (useDefineForClassFields) {
	properties = core.Filter(properties, func(property *ast.Node) bool {
		return property.Initializer() != nil || ast.IsPrivateIdentifier(property.Name()) || ast.HasAccessorModifier(property)
	})
	// }

	privateMethodsAndAccessors := getPrivateInstanceMethodsAndAccessors(node)
	needsConstructorBody := len(properties) > 0 || len(privateMethodsAndAccessors) > 0

	// Only generate synthetic constructor when there are property initializers to move.
	if constructor == nil && !needsConstructorBody {
		return tx.EmitContext().VisitFunctionBody(nil /*node*/, tx.Visitor())
	}

	/// !!! tx.EmitContext().ResumeVariableEnvironment()
	tx.EmitContext().StartVariableEnvironment()

	// needsSyntheticConstructor := constructor == nil && isDerivedClass
	// statementOffset := 0
	statements := make([]*ast.Statement, 0)

	// Add the property initializers. Transforms this:
	//
	//  public x = 1;
	//
	// Into this:
	//
	//  constructor() {
	//      this.x = 1;
	//  }
	//
	initializerStatements := make([]*ast.Statement, 0)
	receiver := tx.Factory().NewThisExpression() // createThis

	// private methods can be called in property initializers, they should execute first.
	initializerStatements = tx.addInstanceMethodStatements(initializerStatements, privateMethodsAndAccessors, receiver)
	if constructor != nil {
		parameterProperties := core.Filter(instanceProperties, func(property *ast.Node) bool {
			return ast.IsParameterPropertyDeclaration(tx.EmitContext().MostOriginal(property), constructor.AsNode())
		})
		nonParameterProperties := core.Filter(instanceProperties, func(property *ast.Node) bool {
			return !ast.IsParameterPropertyDeclaration(tx.EmitContext().MostOriginal(property), constructor.AsNode())
		})
		initializerStatements = tx.addPropertyOrClassStaticBlockStatements(initializerStatements, parameterProperties, receiver)
		initializerStatements = tx.addPropertyOrClassStaticBlockStatements(initializerStatements, nonParameterProperties, receiver)
	} else {
		initializerStatements = tx.addPropertyOrClassStaticBlockStatements(initializerStatements, properties, receiver)
	}

	// !!! handled by runtimesyntax
	// if constructor != nil && constructor.Body != nil {
	// 	_, rest := tx.Factory().SplitStandardPrologue(constructor.Body.Statements())
	// 	superPath := transformers.FindSuperStatementIndexPath(rest, 0)
	// 	if len(superPath) > 0 {
	// 		statements = tx.transformConstructorBodyWorker(rest, superPath, initializerStatements)
	// 	} else {
	// 		statements = append(statements, initializerStatements...)
	// 		statements = append(statements, core.FirstResult(tx.Visitor().VisitSlice(rest))...)
	// 	}
	// } else {
	// 	// !!!
	// 	// if (needsSyntheticConstructor) {
	// 	// Add a synthetic `super` call:
	// 	//
	// 	//  super(...arguments);
	// 	//

	// }
	statements = append(statements, initializerStatements...)

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

func (tx *classFieldsTransformer) transformConstructorBodyWorker(statementsIn []*ast.Statement, superPath []int, initializerStatements []*ast.Statement) []*ast.Statement {
	var statementsOut []*ast.Statement
	superStatementIndex := superPath[0]
	// superStatement := statementsIn[superStatementIndex]

	// visit up to the statement containing `super`
	statementsOut = append(statementsOut, core.FirstResult(tx.Visitor().VisitSlice(statementsIn[:superStatementIndex]))...)
	return statementsOut
}

func (tx *classFieldsTransformer) ClassAlias(node *ast.Node) *ast.Node {
	return tx.classAlias[node]
}

func (tx *classFieldsTransformer) setClassAlias(node *ast.Node, alias *ast.Node) {
	if tx.classAlias == nil {
		tx.classAlias = make(map[*ast.Node]*ast.Node)
	}
	tx.classAlias[node] = alias
}

func (tx *classFieldsTransformer) clearClassAlias() {
	tx.classAlias = nil
}

func (tx *classFieldsTransformer) transformFieldInitializer(node *ast.PropertyDeclaration) *ast.Node {
	if ast.IsPrivateIdentifierClassElementDeclaration(node.AsNode()) {
		return tx.transformPrivateFieldInitializer(node)
	} else {
		return tx.transformPublicFieldInitializer(node)
	}
}

func (tx *classFieldsTransformer) transformPrivateFieldInitializer(node *ast.PropertyDeclaration) *ast.Node {
	if tx.shouldTransformClassElementToWeakMap(node.AsNode()) {
		// If we are transforming private elements into WeakMap/WeakSet, we should elide the node.
		info := tx.accessPrivateIdentifier(node.Name().AsPrivateIdentifier())

		// Leave invalid code untransformed
		if info == nil || info.IsValid() {
			return node.AsNode()
		}

		return nil
	}

	// !!!
	// If we encounter a valid private static field and we're not transforming class static blocks, initialize it

	if isNamedEvaluationAnd(tx.EmitContext(), node.AsNode(), func(v *anonymousFunctionDefinition) bool {
		return isAnonymousClassNeedingAssignedName(tx.EmitContext(), v)
	}) {
		node = transformNamedEvaluation(tx.EmitContext(), node.AsNode(), false /*ignoreEmptyStringLiteral*/, "" /*assignedName*/).AsPropertyDeclaration()
	}

	return tx.Factory().UpdatePropertyDeclaration(
		node,
		node.Modifiers(),
		node.Name(),
		nil, /*postfixToken*/
		nil, /*typeNode*/
		node.Initializer,
	)
}

func (tx *classFieldsTransformer) transformPublicFieldInitializer(node *ast.PropertyDeclaration) *ast.Node {
	// !!! auto accessor
	if !ast.IsAutoAccessorPropertyDeclaration(node.AsNode()) {
		// Create a temporary variable to store a computed property name (if necessary).
		// If it's not inlineable, then we emit an expression after the class which assigns
		// the property name to the temporary variable.
		// expr := tx.getPropertyNameExpressionIfNeeded(
		// 	node.Name(),
		// 	node.Initializer != nil /*shouldHoist*/,
		// )
		// if expr != nil {
		// 	tx.pendingExpressions = append(tx.pendingExpressions, flattenCommaList(expr)...)
		// }

		// !!! isStatic(node)

		return nil
	}

	return tx.Factory().UpdatePropertyDeclaration(
		node,
		node.Modifiers(),
		node.Name(),
		nil, /*postfixToken*/
		nil, /*typeNode*/
		node.Initializer,
	)
}

func (tx *classFieldsTransformer) accessPrivateIdentifier(name *ast.PrivateIdentifier) printer.PrivateIdentifierInfo {
	// !!!
	var info printer.PrivateIdentifierInfo
	tx.classScopeStack.FindFromTop(func(item *classLexicalEnvrionment) bool {
		privateEnv := item.privateEnvrionment
		if privateEnv == nil {
			return false
		}
		info = getPrivateIdentifierInfo(tx.EmitContext(), privateEnv, name)
		return info != nil
	})

	if info == nil || info.Kind() == printer.PrivateIdentifierKindUntransformed {
		return nil
	}

	return info
}

// Generates assignment statements for property initializers.
func (tx *classFieldsTransformer) addPropertyOrClassStaticBlockStatements(statements []*ast.Statement, properties []*ast.Node, receiver *ast.LeftHandSideExpression) []*ast.Statement {
	for _, property := range properties {
		// !!!
		// if (isStatic(property) && !shouldTransformPrivateElementsOrClassStaticBlocks) {
		// 		continue;
		// }

		statement := tx.transformPropertyOrClassStaticBlock(property, receiver)
		if statement == nil {
			continue
		}
		statements = append(statements, statement)
	}

	return statements
}

func (tx *classFieldsTransformer) addPrivateIdentifierToEnvironment(node *ast.Node, name *ast.PrivateIdentifier) {
	previousInfo := getPrivateIdentifierInfo(tx.EmitContext(), tx.getPrivateEnvironment(), name)
	isStatic := ast.HasStaticModifier(node)
	isValid := IsReservedPrivateName(tx.EmitContext(), name) && previousInfo == nil
	if ast.IsAutoAccessorPropertyDeclaration(node) {
		// !!!
	} else if ast.IsPropertyDeclaration(node) {
		if isStatic {
			// !!!
		} else {
			className := tx.getPrivateEnvironment().ClassName
			weakMapName := tx.Factory().NewGeneratedNameForNodeEx(name.AsNode(), printer.AutoGenerateOptions{
				Prefix: "_" + className.Text() + "_",
			})
			tx.EmitContext().AddVariableDeclaration(weakMapName)

			tx.setPrivateIdentifierInfo(name, printer.NewPrivateIdentifierInstanceFieldInfo(
				weakMapName.AsIdentifier(),
				isValid,
			))

			tx.PrependPendingExpressions(tx.Factory().NewAssignmentExpression(
				weakMapName,
				tx.Factory().NewNewExpression(
					tx.Factory().NewIdentifier("WeakMap"),
					nil, /*typeArguments*/
					&ast.NodeList{},
				),
			))

		}
	} else if ast.IsMethodDeclaration(node) {
		// !!!
	} else if ast.IsGetAccessorDeclaration(node) {
		// !!!
	} else if ast.IsSetAccessorDeclaration(node) {
		// !!!
	}
}

func (tx *classFieldsTransformer) addInstanceMethodStatements(statements []*ast.Statement, methods []*ast.Node, receiver *ast.LeftHandSideExpression) []*ast.Statement {
	if len(methods) == 0 {
		return statements
	}

	weakSetName := tx.getPrivateEnvironment().WeakSetName
	return append(statements, tx.Factory().NewExpressionStatement(
		tx.Factory().NewCallExpression(
			tx.Factory().NewPropertyAccessExpression(weakSetName, nil /*questionDotToken*/, tx.Factory().NewIdentifier("add"), ast.NodeFlagsNone),
			nil, /*questionDotToken*/
			nil, /*typeArguments*/
			tx.Factory().NewNodeList([]*ast.Expression{receiver}),
			ast.NodeFlagsNone,
		),
	))
}

func (tx *classFieldsTransformer) transformPropertyOrClassStaticBlock(property *ast.Node, receiver *ast.LeftHandSideExpression) *ast.Statement {
	var expression *ast.Expression
	if ast.IsClassStaticBlockDeclaration(property) {
		expression = tx.setCurrentClassElementAnd(property, func(arg *ast.Node) *ast.Node {
			return tx.transformClassStaticBlockDeclaration(arg.AsClassStaticBlockDeclaration())
		}, property)
	} else {
		expression = tx.transformProperty(property.AsPropertyDeclaration(), receiver)
	}
	if expression == nil {
		return nil
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

	return statement
}

func (tx *classFieldsTransformer) setCurrentClassElementAnd(classElement *ast.ClassElement, visitor func(arg *ast.Node) *ast.Node, arg *ast.Node) *ast.Node {
	if classElement != tx.currentClassElement {
		savedCurrentClassElement := tx.currentClassElement
		tx.currentClassElement = classElement
		result := visitor(arg)
		tx.currentClassElement = savedCurrentClassElement
		return result
	}
	return visitor(arg)
}

func (tx *classFieldsTransformer) transformClassStaticBlockDeclaration(node *ast.ClassStaticBlockDeclaration) *ast.Node {
	// !!!
	return node.AsNode()
}

// Transforms a property initializer into an assignment statement.
func (tx *classFieldsTransformer) transformProperty(property *ast.PropertyDeclaration, receiver *ast.LeftHandSideExpression) *ast.Node {
	savedCurrentClassElement := tx.currentClassElement
	transformed := tx.transformPropertyWorker(property, receiver)
	if transformed != nil &&
		// !!! lexicalEnvironment?.data?.facts &&
		ast.HasStaticModifier(property.AsNode()) {
		// capture the lexical environment for the member
		tx.EmitContext().SetOriginal(transformed, property.AsNode())
		tx.EmitContext().SetSourceMapRange(transformed, tx.EmitContext().SourceMapRange(property.Name()))
		// !!! lexicalEnvironmentMap.set(getOriginalNode(property), lexicalEnvironment);
	}
	tx.currentClassElement = savedCurrentClassElement
	return transformed
}

func (tx *classFieldsTransformer) transformPropertyWorker(property *ast.PropertyDeclaration, receiver *ast.LeftHandSideExpression) *ast.Expression {
	// !!! const emitAssignment = !useDefineForClassFields;
	emitAssignment := true

	if isNamedEvaluationAnd(tx.EmitContext(), property.AsNode(), func(v *anonymousFunctionDefinition) bool {
		return isAnonymousClassNeedingAssignedName(tx.EmitContext(), v)
	}) {
		property = transformNamedEvaluation(tx.EmitContext(), property.AsNode(), false /*ignoreEmptyStringLiteral*/, "" /*assignedName*/).AsPropertyDeclaration()
	}

	propertyName := property.Name()
	if ast.HasAccessorModifier(property.AsNode()) {
		propertyName = tx.Factory().NewGeneratedPrivateNameForNode(property.Name())
	} else if ast.IsComputedPropertyName(property.Name()) && transformers.IsSimpleInlineableExpression(property.Name().Expression()) {
		propertyName = tx.Factory().UpdateComputedPropertyName(property.Name().AsComputedPropertyName(), tx.Factory().NewGeneratedNameForNode(property.Name()))
	}

	if ast.HasStaticModifier(property.AsNode()) {
		tx.currentClassElement = property.AsNode()
	}

	if ast.IsPrivateIdentifier(propertyName) && tx.shouldTransformClassElementToWeakMap(property.AsNode()) {
		privateIdentifierInfo := tx.accessPrivateIdentifier(propertyName.AsPrivateIdentifier())
		if privateIdentifierInfo != nil {
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
			case *printer.PrivateIdentifierStaticFieldInfo:
				{
					// `createPrivateStaticFieldInitializer` in Strada
					initializer := tx.Visitor().VisitNode(property.Initializer)
					if initializer == nil {
						initializer = tx.Factory().NewVoidZeroExpression()
					}
					return tx.Factory().NewAssignmentExpression(
						info.VariableName.AsNode(),
						tx.Factory().NewObjectLiteralExpression(
							tx.Factory().NewNodeList([]*ast.Expression{
								tx.Factory().NewPropertyAssignment(
									nil, /*modifiers*/
									tx.Factory().NewIdentifier("value"),
									nil, /*postfixToken*/
									nil, /*typeNode*/
									initializer,
								),
							}),
							false,
						),
					)
				}
			default:
				return nil
			}
		} else {
			panic("Undeclared private name for property declaration.")
		}
	}

	if ast.IsPrivateIdentifier(propertyName) || ast.HasAccessorModifier(property.AsNode()) && property.Initializer == nil {
		return nil
	}

	propertyOriginalNode := tx.EmitContext().MostOriginal(property.AsNode())
	if ast.HasSyntacticModifier(propertyOriginalNode, ast.ModifierFlagsAbstract) {
		return nil
	}

	// !!! initializer
	initializer := tx.Visitor().VisitNode(property.Initializer)
	if ast.IsParameterPropertyDeclaration(propertyOriginalNode, propertyOriginalNode.Parent) && ast.IsIdentifier(propertyName) {
		// A parameter-property declaration always overrides the initializer. The only time a parameter-property
		// declaration *should* have an initializer is when decorators have added initializers that need to run before
		// any other initializer
		localName := tx.Factory().DeepCloneNode(propertyName)
		if initializer != nil {
			// unwrap `(__runInitializers(this, _instanceExtraInitializers), void 0)`
			if ast.IsParenthesizedExpression(initializer) &&
				ast.IsCommaExpression(initializer.Expression()) &&
				tx.EmitContext().IsCallToHelper(initializer.Expression(), "___runInitializers") &&
				ast.IsVoidExpression(initializer.Expression().AsBinaryExpression().Right) &&
				ast.IsNumericLiteral(initializer.Expression().AsBinaryExpression().Right.Expression()) {
				initializer = initializer.Expression().AsBinaryExpression().Left
			}
			initializer = tx.Factory().InlineExpressions([]*ast.Expression{initializer, localName})
		} else {
			initializer = localName
		}
		tx.EmitContext().SetEmitFlags(propertyName, printer.EFNoComments|printer.EFNoSourceMap)
		tx.EmitContext().SetSourceMapRange(localName, propertyOriginalNode.Name().Loc)
		tx.EmitContext().SetEmitFlags(localName, printer.EFNoComments)
	} else {
		if initializer == nil {
			initializer = tx.Factory().NewVoidZeroExpression()
		}
	}

	if emitAssignment || ast.IsPrivateIdentifier(propertyName) {
		memberAccess := tx.NewMemberAccessForPropertyName(receiver, propertyName)
		tx.EmitContext().AddEmitFlags(memberAccess, printer.EFNoLeadingComments)
		expression := tx.Factory().NewAssignmentExpression(memberAccess, initializer)
		return expression
	} else {
		name := propertyName
		if ast.IsComputedPropertyName(propertyName) {
			name = propertyName.Expression()
		} else if ast.IsIdentifier(propertyName) {
			name = tx.Factory().NewStringLiteral(propertyName.Text())
		}
		// !!! factory.createPropertyDescriptor
		// const descriptor = factory.createPropertyDescriptor({ value: initializer, configurable: true, writable: true, enumerable: true });
		// return factory.createObjectDefinePropertyCall(receiver, name, descriptor);
		return tx.Factory().NewAssignmentExpression(receiver, name)
	}

}

func (tx *classFieldsTransformer) shouldTransformClassElementToWeakMap(node *ast.Node) bool {
	// !!!
	return true
}

func (tx *classFieldsTransformer) NewMemberAccessForPropertyName(target *ast.Expression, memberName *ast.PropertyName) *ast.Expression {
	var expression *ast.Expression
	if ast.IsComputedPropertyName(memberName) {
		expression = tx.Factory().NewElementAccessExpression(target, nil, memberName.Expression(), ast.NodeFlagsNone)
	} else {
		if ast.IsMemberName(memberName) {
			expression = tx.Factory().NewPropertyAccessExpression(target, nil, memberName, ast.NodeFlagsNone)
		} else {
			expression = tx.Factory().NewElementAccessExpression(target, nil, memberName, ast.NodeFlagsNone)
		}
		tx.EmitContext().AddEmitFlags(expression, printer.EFNoNestedSourceMaps)
	}
	expression.Loc = memberName.Loc

	return expression
}

func (tx *classFieldsTransformer) NewPrivateIdentifierAccess(info printer.PrivateIdentifierInfo, receiver *ast.Expression) *ast.Expression {
	receiver = tx.Visitor().VisitNode(receiver)
	// !!! ensureDynamicThisIfNeeded(node)
	return tx.NewPrivateIdentifierAccessHelper(info, receiver)
}

func (tx *classFieldsTransformer) NewPrivateIdentifierAccessHelper(info printer.PrivateIdentifierInfo, receiver *ast.Expression) *ast.Expression {
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

func (tx *classFieldsTransformer) enableSubstitutionForClassAliases() {
	if tx.enabledSubstitutions&ClassPropertySubstitutionFlagsClassAliases == 0 {
		tx.enabledSubstitutions |= ClassPropertySubstitutionFlagsClassAliases

		// We need to enable substitutions for identifiers. This allows us to
		// substitute class names inside of a class declaration.
		// context.enableSubstitution(SyntaxKind.Identifier);

		// Keep track of class aliases.
		tx.clearClassAlias()
	}
}

func (tx *classFieldsTransformer) enableSubstitutionForClassStaticThisOrSuperReference() {
	if tx.enabledSubstitutions&ClassPropertySubstitutionFlagsClassStaticThisOrSuperReference == 0 {
		tx.enabledSubstitutions |= ClassPropertySubstitutionFlagsClassStaticThisOrSuperReference

		// substitute `this` in a static field initializer

		// !!!
		// context.enableSubstitution(SyntaxKind.ThisKeyword)
		// context.enableEmitNotification(SyntaxKind.FunctionDeclaration)
		// context.enableEmitNotification(SyntaxKind.FunctionExpression)
		// context.enableEmitNotification(SyntaxKind.Constructor)
		// context.enableEmitNotification(SyntaxKind.GetAccessor)
		// context.enableEmitNotification(SyntaxKind.SetAccessor)
		// context.enableEmitNotification(SyntaxKind.MethodDeclaration)
		// context.enableEmitNotification(SyntaxKind.PropertyDeclaration)
		// context.enableEmitNotification(SyntaxKind.ComputedPropertyName)
	}
}

func (tx *classFieldsTransformer) getPrivateIdentifierInfo(name *ast.PrivateIdentifier) printer.PrivateIdentifierInfo {
	privateEnv := tx.getPrivateEnvironment()
	if transformers.IsGeneratedIdentifier(tx.EmitContext(), name.AsNode()) {
		if privateEnv.GeneratedIdentifiers == nil {
			return nil
		} else {
			return privateEnv.GeneratedIdentifiers[tx.EmitContext().GetNodeForGeneratedName(name.AsNode())]
		}
	} else {
		if privateEnv.Identifiers == nil {
			return nil
		}
		return privateEnv.Identifiers[name.Text]
	}
}

func (tx *classFieldsTransformer) setPrivateIdentifierInfo(name *ast.PrivateIdentifier, info printer.PrivateIdentifierInfo) {
	privateEnv := tx.getPrivateEnvironment()
	if transformers.IsGeneratedIdentifier(tx.EmitContext(), name.AsNode()) {
		if privateEnv.GeneratedIdentifiers == nil {
			privateEnv.GeneratedIdentifiers = make(map[*ast.Node]printer.PrivateIdentifierInfo)
		}
		privateEnv.GeneratedIdentifiers[tx.EmitContext().GetNodeForGeneratedName(name.AsNode())] = info
	} else {
		if privateEnv.Identifiers == nil {
			privateEnv.Identifiers = make(map[string]printer.PrivateIdentifierInfo)
		}
		privateEnv.Identifiers[name.Text] = info
	}
}

/**
 * Gets a value indicating whether a class element is a private instance method or accessor.
 */
func isNonStaticMethodOrAccessorWithPrivateName(classMemberNode *ast.Node) bool {
	return !ast.IsStatic(classMemberNode) && (ast.IsMethodOrAccessor(classMemberNode) || ast.IsAutoAccessorPropertyDeclaration(classMemberNode)) && ast.IsPrivateIdentifier(classMemberNode.Name())
}

func getPrivateInstanceMethodsAndAccessors(node *ast.ClassLikeDeclaration) []*ast.Node {
	return core.Filter(node.Members(), isNonStaticMethodOrAccessorWithPrivateName)
}

func getClassFacts(emitContext *printer.EmitContext, node *ast.ClassLikeDeclaration) ClassFacts {
	facts := ClassFactsNone

	// !!! class decorated
	// if (isClassLike(original) && classOrConstructorParameterIsDecorated(legacyDecorators, original)) {
	//     facts |= ClassFacts.ClassWasDecorated;
	// }

	if classHasClassThisAssignment(emitContext, node) || classHasExplicitlyAssignedName(emitContext, node) {
		facts |= ClassFactsNeedsClassConstructorReference
	}

	containsPublicInstanceFields := false
	containsInitializedPublicInstanceFields := false
	containsInstancePrivateElements := false
	containsInstanceAutoAccessors := false

	for _, member := range node.Members() {
		if ast.IsStatic(member) {
			if member.Name() != nil && (ast.IsPrivateIdentifier(member.Name()) || ast.IsAutoAccessorPropertyDeclaration(member)) {
				facts |= ClassFactsNeedsClassConstructorReference
			} else if
			// !!! shouldTransformAutoAccessors &&
			member.Name() == nil &&
				emitContext.ClassThis(member) == nil &&
				ast.IsAutoAccessorPropertyDeclaration(member) {
				facts |= ClassFactsNeedsClassConstructorReference
			}

			if ast.IsPropertyDeclaration(member) || ast.IsClassStaticBlockDeclaration(member) {
				if member.SubtreeFacts()&ast.SubtreeContainsLexicalThis != 0 {
					facts |= ClassFactsNeedsSubstitutionForThisInClassStaticField
					// !!!
					// if (!(facts & ClassFacts.ClassWasDecorated)) {
					facts |= ClassFactsNeedsClassConstructorReference
					// }
				}
				if member.SubtreeFacts()&ast.SubtreeContainsLexicalSuper != 0 {
					// !!!
					// if (!(facts & ClassFacts.ClassWasDecorated)) {
					facts |= ClassFactsNeedsClassConstructorReference | ClassFactsNeedsClassSuperReference
					// }
				}
			}
		} else if !ast.HasAccessorModifier(emitContext.MostOriginal(member)) {
			if ast.IsAutoAccessorPropertyDeclaration(member) {
				containsInstanceAutoAccessors = true
				if ast.IsPrivateIdentifierClassElementDeclaration(member) {
					containsInstancePrivateElements = true
				}
			} else if ast.IsPrivateIdentifierClassElementDeclaration(member) {
				containsInstancePrivateElements = true
				// !!!
				// if (resolver.hasNodeCheckFlag(member, NodeCheckFlags.ContainsConstructorReference)) {
				// 		facts |= ClassFacts.NeedsClassConstructorReference;
				// }
			} else if ast.IsPropertyDeclaration(member) {
				containsPublicInstanceFields = true
				if member.Initializer() != nil {
					containsInitializedPublicInstanceFields = true
				}
			}
		}
	}

	// !!! shouldTransformInitializersUsingSet
	// !!! shouldTransformAutoAccessors === Ternary.True
	willHoistInitializersToConstructor := containsPublicInstanceFields ||
		containsInitializedPublicInstanceFields ||
		containsInstancePrivateElements ||
		containsInstanceAutoAccessors

	if willHoistInitializersToConstructor {
		facts |= ClassFactsWillHoistInitializersToConstructor
	}

	return facts

}

func isAnonymousClassNeedingAssignedName(emitContext *printer.EmitContext, node *anonymousFunctionDefinition) bool {
	if ast.IsClassExpression(node) && node.Name() != nil {
		staticPropertiesOrClassStaticBlocks := ast.GetStaticPropertiesAndClassStaticBlock(node)
		if core.Some(staticPropertiesOrClassStaticBlocks, func(item *ast.Node) bool {
			return isClassNamedEvaluationHelperBlock(emitContext, item)
		}) {
			return false
		}

		return true
		// !!!
		// const hasTransformableStatics = (shouldTransformPrivateElementsOrClassStaticBlocks ||
		// !!(getInternalEmitFlags(node) && InternalEmitFlags.TransformPrivateStaticElements)) &&
		// some(staticPropertiesOrClassStaticBlocks, node =>
		// 		isClassStaticBlockDeclaration(node) ||
		// 		isPrivateIdentifierClassElementDeclaration(node) ||
		// 		shouldTransformInitializers && isInitializedProperty(node));
		// return hasTransformableStatics;
	}

	return false
}

// Gets all the static or all the instance property declarations of a class
func getProperties(node *ast.ClassLikeDeclaration, requireInitializer bool, isStatic bool) []*ast.Node {
	return core.Filter(node.Members(), func(member *ast.Node) bool {
		return member.Kind == ast.KindPropertyDeclaration && (member.Initializer() != nil || !requireInitializer) && ast.HasStaticModifier(member) == isStatic
	})
}

func IsReservedPrivateName(emitContext *printer.EmitContext, node *ast.PrivateIdentifier) bool {
	return !transformers.IsGeneratedIdentifier(emitContext, node.AsNode()) && node.Text == "#constructor"
}

func getPrivateIdentifierInfo(emitContext *printer.EmitContext, privateEnv *PrivateEnvironment, name *ast.PrivateIdentifier) printer.PrivateIdentifierInfo {
	if transformers.IsGeneratedIdentifier(emitContext, name.AsNode()) {
		if privateEnv.GeneratedIdentifiers == nil {
			return nil
		} else {
			return privateEnv.GeneratedIdentifiers[emitContext.GetNodeForGeneratedName(name.AsNode())]
		}
	} else {
		if privateEnv.Identifiers == nil {
			return nil
		}
		return privateEnv.Identifiers[name.Text]
	}
}
