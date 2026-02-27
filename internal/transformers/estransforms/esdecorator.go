package estransforms

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/printer"
	"github.com/microsoft/typescript-go/internal/scanner"
	"github.com/microsoft/typescript-go/internal/transformers"
)

// lexicalEntryKind discriminates the kind of lexical scope entry.
type lexicalEntryKind int

const (
	lexicalEntryKindClass lexicalEntryKind = iota
	lexicalEntryKindClassElement
	lexicalEntryKindName
	lexicalEntryKindOther
)

// lexicalEntry represents a single entry in the lexical scope stack used to track
// nested class declarations and their state during transformation.
type lexicalEntry struct {
	kind                    lexicalEntryKind
	next                    *lexicalEntry
	classInfoData           *classInfo
	savedPendingExpressions []*ast.Expression
	classThisData           *ast.IdentifierNode
	classSuperData          *ast.IdentifierNode
	depth                   int
}

// memberInfo stores decoration-related data for a single class element.
type memberInfo struct {
	memberDecoratorsName        *ast.IdentifierNode
	memberInitializersName      *ast.IdentifierNode
	memberExtraInitializersName *ast.IdentifierNode
	memberDescriptorName        *ast.IdentifierNode
}

// classInfo stores all transformation data for a single decorated class.
type classInfo struct {
	class                                 *ast.Node
	classDecoratorsName                   *ast.IdentifierNode
	classDescriptorName                   *ast.IdentifierNode
	classExtraInitializersName            *ast.IdentifierNode
	classThis                             *ast.IdentifierNode
	classSuper                            *ast.IdentifierNode
	metadataReference                     *ast.IdentifierNode
	memberInfos                           map[*ast.Node]*memberInfo
	instanceMethodExtraInitializersName   *ast.IdentifierNode
	staticMethodExtraInitializersName     *ast.IdentifierNode
	staticNonFieldDecorationStatements    []*ast.Statement
	nonStaticNonFieldDecorationStatements []*ast.Statement
	staticFieldDecorationStatements       []*ast.Statement
	nonStaticFieldDecorationStatements    []*ast.Statement
	hasStaticInitializers                 bool
	hasNonAmbientInstanceFields           bool
	hasStaticPrivateClassElements         bool
	pendingStaticInitializers             []*ast.Expression
	pendingInstanceInitializers           []*ast.Expression
}

type esDecoratorTransformer struct {
	transformers.Transformer
	compilerOptions    *core.CompilerOptions
	top                *lexicalEntry
	classInfoStack     *classInfo
	classThis          *ast.IdentifierNode
	classSuper         *ast.IdentifierNode
	pendingExpressions []*ast.Expression
	discardedVisitor   *ast.NodeVisitor
}

func newESDecoratorTransformer(opts *transformers.TransformOptions) *transformers.Transformer {
	tx := &esDecoratorTransformer{compilerOptions: opts.CompilerOptions}
	result := tx.NewTransformer(tx.visit, opts.Context)
	tx.discardedVisitor = tx.EmitContext().NewNodeVisitor(tx.discardedValueVisit)
	return result
}

// --- Lexical environment stack management ---

func (tx *esDecoratorTransformer) updateState() {
	tx.classInfoStack = nil
	tx.classThis = nil
	tx.classSuper = nil
	if tx.top == nil {
		return
	}
	switch tx.top.kind {
	case lexicalEntryKindClass:
		tx.classInfoStack = tx.top.classInfoData
	case lexicalEntryKindClassElement:
		if tx.top.next != nil {
			tx.classInfoStack = tx.top.next.classInfoData
		}
		tx.classThis = tx.top.classThisData
		tx.classSuper = tx.top.classSuperData
	case lexicalEntryKindName:
		// name -> class-element -> class -> ???
		if tx.top.next != nil && tx.top.next.next != nil && tx.top.next.next.next != nil {
			grandparent := tx.top.next.next.next
			if grandparent.kind == lexicalEntryKindClassElement {
				if grandparent.next != nil {
					tx.classInfoStack = grandparent.next.classInfoData
				}
				tx.classThis = grandparent.classThisData
				tx.classSuper = grandparent.classSuperData
			}
		}
	}
}

func (tx *esDecoratorTransformer) enterClass(ci *classInfo) {
	tx.top = &lexicalEntry{
		kind:                    lexicalEntryKindClass,
		next:                    tx.top,
		classInfoData:           ci,
		savedPendingExpressions: tx.pendingExpressions,
	}
	tx.pendingExpressions = nil
	tx.updateState()
}

func (tx *esDecoratorTransformer) exitClass() {
	tx.pendingExpressions = tx.top.savedPendingExpressions
	tx.top = tx.top.next
	tx.updateState()
}

func (tx *esDecoratorTransformer) enterClassElement(node *ast.Node) {
	entry := &lexicalEntry{
		kind: lexicalEntryKindClassElement,
		next: tx.top,
	}
	if ast.IsClassStaticBlockDeclaration(node) || (ast.IsPropertyDeclaration(node) && ast.HasStaticModifier(node)) {
		if tx.top != nil && tx.top.classInfoData != nil {
			entry.classThisData = tx.top.classInfoData.classThis
			entry.classSuperData = tx.top.classInfoData.classSuper
		}
	}
	tx.top = entry
	tx.updateState()
}

func (tx *esDecoratorTransformer) exitClassElement() {
	tx.top = tx.top.next
	tx.updateState()
}

func (tx *esDecoratorTransformer) enterName() {
	tx.top = &lexicalEntry{
		kind: lexicalEntryKindName,
		next: tx.top,
	}
	tx.updateState()
}

func (tx *esDecoratorTransformer) exitName() {
	tx.top = tx.top.next
	tx.updateState()
}

func (tx *esDecoratorTransformer) enterOther() {
	if tx.top != nil && tx.top.kind == lexicalEntryKindOther {
		tx.top.depth++
	} else {
		tx.top = &lexicalEntry{
			kind:                    lexicalEntryKindOther,
			next:                    tx.top,
			savedPendingExpressions: tx.pendingExpressions,
		}
		tx.pendingExpressions = nil
		tx.updateState()
	}
}

func (tx *esDecoratorTransformer) exitOther() {
	if tx.top.depth > 0 {
		tx.top.depth--
	} else {
		tx.pendingExpressions = tx.top.savedPendingExpressions
		tx.top = tx.top.next
		tx.updateState()
	}
}

// --- Visitor dispatch ---

func (tx *esDecoratorTransformer) shouldVisitNode(node *ast.Node) bool {
	return node.SubtreeFacts()&ast.SubtreeContainsDecorators != 0 ||
		(tx.classThis != nil && node.SubtreeFacts()&ast.SubtreeContainsLexicalThis != 0) ||
		(tx.classThis != nil && tx.classSuper != nil && node.SubtreeFacts()&ast.SubtreeContainsLexicalSuper != 0)
}

func (tx *esDecoratorTransformer) visit(node *ast.Node) *ast.Node {
	if !tx.shouldVisitNode(node) {
		return node
	}
	switch node.Kind {
	case ast.KindDecorator:
		return nil
	case ast.KindClassDeclaration:
		return tx.visitClassDeclaration(node.AsClassDeclaration())
	case ast.KindClassExpression:
		return tx.visitClassExpression(node.AsClassExpression())
	case ast.KindParameter:
		return tx.visitParameterDeclaration(node)
	case ast.KindBinaryExpression:
		return tx.visitBinaryExpression(node, false /*discarded*/)
	case ast.KindPrefixUnaryExpression, ast.KindPostfixUnaryExpression:
		return tx.visitPreOrPostfixUnaryExpression(node, false /*discarded*/)
	case ast.KindParenthesizedExpression:
		return tx.visitParenthesizedExpression(node, false /*discarded*/)
	case ast.KindPropertyAssignment:
		return tx.visitPropertyAssignment(node)
	case ast.KindVariableDeclaration:
		return tx.visitVariableDeclaration(node)
	case ast.KindBindingElement:
		return tx.visitBindingElement(node)
	case ast.KindExportAssignment:
		return tx.visitExportAssignment(node)
	case ast.KindThisKeyword:
		return tx.visitThisExpression(node)
	case ast.KindForStatement:
		return tx.visitForStatement(node)
	case ast.KindExpressionStatement:
		return tx.visitExpressionStatement(node)
	case ast.KindCallExpression:
		return tx.visitCallExpression(node)
	case ast.KindTaggedTemplateExpression:
		return tx.visitTaggedTemplateExpression(node)
	case ast.KindPropertyAccessExpression:
		return tx.visitPropertyAccessExpression(node)
	case ast.KindElementAccessExpression:
		return tx.visitElementAccessExpression(node)
	case ast.KindComputedPropertyName:
		return tx.visitComputedPropertyName(node)
	case ast.KindMethodDeclaration,
		ast.KindSetAccessor,
		ast.KindGetAccessor,
		ast.KindFunctionExpression,
		ast.KindFunctionDeclaration:
		tx.enterOther()
		result := tx.Visitor().VisitEachChild(node)
		tx.exitOther()
		return result
	default:
		return tx.Visitor().VisitEachChild(node)
	}
}

func (tx *esDecoratorTransformer) discardedValueVisit(node *ast.Node) *ast.Node {
	if !tx.shouldVisitNode(node) {
		return node
	}
	switch node.Kind {
	case ast.KindPrefixUnaryExpression, ast.KindPostfixUnaryExpression:
		return tx.visitPreOrPostfixUnaryExpression(node, true /*discarded*/)
	case ast.KindBinaryExpression:
		return tx.visitBinaryExpression(node, true /*discarded*/)
	case ast.KindParenthesizedExpression:
		return tx.visitParenthesizedExpression(node, true /*discarded*/)
	default:
		return tx.visit(node)
	}
}

func (tx *esDecoratorTransformer) classElementVisitor(node *ast.Node) *ast.Node {
	switch node.Kind {
	case ast.KindConstructor:
		return tx.visitConstructorDeclaration(node)
	case ast.KindMethodDeclaration:
		return tx.visitMethodDeclaration(node)
	case ast.KindGetAccessor:
		return tx.visitGetAccessorDeclaration(node)
	case ast.KindSetAccessor:
		return tx.visitSetAccessorDeclaration(node)
	case ast.KindPropertyDeclaration:
		return tx.visitPropertyDeclaration(node)
	case ast.KindClassStaticBlockDeclaration:
		return tx.visitClassStaticBlockDeclaration(node)
	default:
		return tx.visit(node)
	}
}

func (tx *esDecoratorTransformer) modifierVisitor(node *ast.Node) *ast.Node {
	if node.Kind == ast.KindDecorator {
		return nil
	}
	return node
}

// --- Helper utilities ---

func isDecoratedClassLike(node *ast.Node) bool {
	return ast.ClassOrConstructorParameterIsDecorated(false, node) ||
		ast.ChildIsDecorated(false, node, nil)
}

func moveRangePastDecorators(node *ast.Node) core.TextRange {
	var lastDecorator *ast.Node
	if ast.CanHaveModifiers(node) {
		nodes := node.ModifierNodes()
		if nodes != nil {
			lastDecorator = core.FindLast(nodes, ast.IsDecorator)
		}
	}
	if lastDecorator != nil && !ast.PositionIsSynthesized(lastDecorator.End()) {
		return core.NewTextRange(lastDecorator.End(), node.End())
	}
	return node.Loc
}

func moveRangePastModifiers(node *ast.Node) core.TextRange {
	if ast.IsPropertyDeclaration(node) || ast.IsMethodDeclaration(node) {
		return core.NewTextRange(node.Name().Pos(), node.End())
	}
	var lastModifier *ast.Node
	if ast.CanHaveModifiers(node) {
		lastModifier = core.LastOrNil(node.ModifierNodes())
	}
	if lastModifier != nil && !ast.PositionIsSynthesized(lastModifier.End()) {
		return core.NewTextRange(lastModifier.End(), node.End())
	}
	return moveRangePastDecorators(node)
}

func getHelperVariableName(node *ast.Node) string {
	declarationName := ""
	if node.Name() != nil && ast.IsIdentifier(node.Name()) {
		declarationName = node.Name().Text()
	} else if node.Name() != nil && ast.IsPrivateIdentifier(node.Name()) {
		text := node.Name().Text()
		if len(text) > 1 {
			declarationName = text[1:]
		}
	} else if node.Name() != nil && ast.IsStringLiteral(node.Name()) && scanner.IsIdentifierText(node.Name().Text(), core.LanguageVariantStandard) {
		declarationName = node.Name().Text()
	} else if ast.IsClassLike(node) {
		declarationName = "class"
	} else {
		declarationName = "member"
	}

	if ast.IsGetAccessorDeclaration(node) {
		declarationName = "get_" + declarationName
	}
	if ast.IsSetAccessorDeclaration(node) {
		declarationName = "set_" + declarationName
	}
	if node.Name() != nil && ast.IsPrivateIdentifier(node.Name()) {
		declarationName = "private_" + declarationName
	}
	if ast.IsStatic(node) {
		declarationName = "static_" + declarationName
	}
	return "_" + declarationName
}

func (tx *esDecoratorTransformer) createHelperVariable(node *ast.Node, suffix string) *ast.IdentifierNode {
	return tx.Factory().NewUniqueNameEx(
		getHelperVariableName(node)+"_"+suffix,
		printer.AutoGenerateOptions{Flags: printer.GeneratedIdentifierFlagsOptimistic | printer.GeneratedIdentifierFlagsReservedInNestedScopes},
	)
}

func (tx *esDecoratorTransformer) createLet(name *ast.IdentifierNode, initializer *ast.Expression) *ast.Statement {
	return tx.Factory().NewVariableStatement(
		nil,
		tx.Factory().NewVariableDeclarationList(
			ast.NodeFlagsLet,
			tx.Factory().NewNodeList([]*ast.Node{
				tx.Factory().NewVariableDeclaration(name, nil, nil, initializer),
			}),
		),
	)
}

// getAllDecoratorsOfClass returns the decorators for a class-like declaration (ES decorators).
func getAllDecoratorsOfClass(node *ast.Node) []*ast.Node {
	return node.Decorators()
}

// getAllDecoratorsOfClassElement returns the decorators for a class element (ES decorators).
func getAllDecoratorsOfClassElement(member *ast.Node, parent *ast.Node) []*ast.Node {
	switch member.Kind {
	case ast.KindGetAccessor, ast.KindSetAccessor, ast.KindMethodDeclaration:
		if member.Body() == nil {
			return nil
		}
		return member.Decorators()
	case ast.KindPropertyDeclaration:
		return member.Decorators()
	default:
		return nil
	}
}

// --- transformDecorator transforms a decorator expression ---

func (tx *esDecoratorTransformer) transformDecorator(decorator *ast.Node) *ast.Expression {
	expression := tx.Visitor().VisitNode(decorator.AsDecorator().Expression)
	tx.EmitContext().SetEmitFlags(expression, printer.EFNoComments)

	// preserve the 'this' binding for an access expression
	innerExpression := ast.SkipOuterExpressions(expression, ast.OEKAll)
	if ast.IsAccessExpression(innerExpression) {
		target, thisArg := tx.createCallBinding(expression)
		bindCall := tx.Factory().NewFunctionBindCall(target, thisArg, nil)
		return tx.Factory().RestoreOuterExpressions(expression, bindCall, ast.OEKAll)
	}
	return expression
}

// createCallBinding is a simplified version of the factory's createCallBinding.
func (tx *esDecoratorTransformer) createCallBinding(expression *ast.Expression) (target *ast.Expression, thisArg *ast.Expression) {
	callee := ast.SkipOuterExpressions(expression, ast.OEKAll)
	if ast.IsSuperProperty(callee) {
		return callee, tx.Factory().NewThisExpression()
	}
	if callee.Kind == ast.KindSuperKeyword {
		return callee, tx.Factory().NewThisExpression()
	}
	if ast.IsPropertyAccessExpression(callee) {
		pa := callee.AsPropertyAccessExpression()
		if tx.shouldBeCapturedInTempVariable(pa.Expression) {
			thisArg = tx.Factory().NewTempVariable()
			assign := tx.Factory().NewAssignmentExpression(thisArg, pa.Expression)
			assign.Loc = pa.Expression.Loc
			target = tx.Factory().NewPropertyAccessExpression(assign, nil, pa.Name(), ast.NodeFlagsNone)
			target.Loc = callee.Loc
			return target, thisArg
		}
		return callee, pa.Expression
	}
	if ast.IsElementAccessExpression(callee) {
		ea := callee.AsElementAccessExpression()
		if tx.shouldBeCapturedInTempVariable(ea.Expression) {
			thisArg = tx.Factory().NewTempVariable()
			assign := tx.Factory().NewAssignmentExpression(thisArg, ea.Expression)
			assign.Loc = ea.Expression.Loc
			target = tx.Factory().NewElementAccessExpression(assign, nil, ea.ArgumentExpression, ast.NodeFlagsNone)
			target.Loc = callee.Loc
			return target, thisArg
		}
		return callee, ea.Expression
	}
	return expression, tx.Factory().NewVoidZeroExpression()
}

func (tx *esDecoratorTransformer) shouldBeCapturedInTempVariable(node *ast.Expression) bool {
	// Capture identifiers that are not simple (i.e., not just a name or "this")
	if ast.IsIdentifier(node) || node.Kind == ast.KindThisKeyword {
		return false
	}
	return true
}

func (tx *esDecoratorTransformer) transformAllDecoratorsOfDeclaration(decorators []*ast.Node) []*ast.Expression {
	if len(decorators) == 0 {
		return nil
	}
	result := make([]*ast.Expression, 0, len(decorators))
	for _, d := range decorators {
		result = append(result, tx.transformDecorator(d))
	}
	return result
}

// --- createClassInfo: creates bookkeeping struct for class transformation ---

func (tx *esDecoratorTransformer) createClassInfo(node *ast.Node) *classInfo {
	f := tx.Factory()
	ci := &classInfo{
		class: node,
		metadataReference: f.NewUniqueNameEx("_metadata", printer.AutoGenerateOptions{
			Flags: printer.GeneratedIdentifierFlagsOptimistic | printer.GeneratedIdentifierFlagsFileLevel,
		}),
	}

	// If the class itself is decorated, create a _classThis binding
	if ast.HasDecorators(node) && ast.NodeCanBeDecorated(false, node, nil, nil) {
		needsUniqueClassThis := core.Some(node.Members(), func(member *ast.Node) bool {
			return (ast.IsPrivateIdentifierClassElementDeclaration(member) || ast.IsAutoAccessorPropertyDeclaration(member)) && ast.HasStaticModifier(member)
		})
		var flags printer.GeneratedIdentifierFlags = printer.GeneratedIdentifierFlagsOptimistic | printer.GeneratedIdentifierFlagsFileLevel
		if needsUniqueClassThis {
			flags = printer.GeneratedIdentifierFlagsOptimistic | printer.GeneratedIdentifierFlagsReservedInNestedScopes
		}
		ci.classThis = f.NewUniqueNameEx("_classThis", printer.AutoGenerateOptions{Flags: flags})
	}

	for _, member := range node.Members() {
		if ast.IsMethodOrAccessor(member) && ast.NodeOrChildIsDecorated(false, member, node, nil) {
			if ast.HasStaticModifier(member) {
				if ci.staticMethodExtraInitializersName == nil {
					ci.staticMethodExtraInitializersName = f.NewUniqueNameEx("_staticExtraInitializers", printer.AutoGenerateOptions{
						Flags: printer.GeneratedIdentifierFlagsOptimistic | printer.GeneratedIdentifierFlagsFileLevel,
					})
					var renamedClassThis *ast.Node
					if ci.classThis != nil {
						renamedClassThis = ci.classThis
					} else {
						renamedClassThis = f.NewThisExpression()
					}
					initializer := f.NewRunInitializersHelper(renamedClassThis, ci.staticMethodExtraInitializersName, nil)
					nameRange := node.Name()
					if nameRange != nil {
						tx.EmitContext().SetSourceMapRange(initializer, nameRange.Loc)
					} else {
						tx.EmitContext().SetSourceMapRange(initializer, moveRangePastDecorators(node))
					}
					ci.pendingStaticInitializers = append(ci.pendingStaticInitializers, initializer)
				}
			} else {
				if ci.instanceMethodExtraInitializersName == nil {
					ci.instanceMethodExtraInitializersName = f.NewUniqueNameEx("_instanceExtraInitializers", printer.AutoGenerateOptions{
						Flags: printer.GeneratedIdentifierFlagsOptimistic | printer.GeneratedIdentifierFlagsFileLevel,
					})
					initializer := f.NewRunInitializersHelper(f.NewThisExpression(), ci.instanceMethodExtraInitializersName, nil)
					nameRange := node.Name()
					if nameRange != nil {
						tx.EmitContext().SetSourceMapRange(initializer, nameRange.Loc)
					} else {
						tx.EmitContext().SetSourceMapRange(initializer, moveRangePastDecorators(node))
					}
					ci.pendingInstanceInitializers = append(ci.pendingInstanceInitializers, initializer)
				}
			}
		}

		if ast.IsClassStaticBlockDeclaration(member) {
			if !isClassNamedEvaluationHelperBlock(tx.EmitContext(), member) {
				ci.hasStaticInitializers = true
			}
		} else if ast.IsPropertyDeclaration(member) {
			if ast.HasStaticModifier(member) {
				ci.hasStaticInitializers = ci.hasStaticInitializers || member.Initializer() != nil || ast.HasDecorators(member)
			} else {
				ci.hasNonAmbientInstanceFields = ci.hasNonAmbientInstanceFields || !ast.HasSyntacticModifier(member, ast.ModifierFlagsAmbient)
			}
		}

		if (ast.IsPrivateIdentifierClassElementDeclaration(member) || ast.IsAutoAccessorPropertyDeclaration(member)) && ast.HasStaticModifier(member) {
			ci.hasStaticPrivateClassElements = true
		}

		// exit early if possible
		if ci.staticMethodExtraInitializersName != nil &&
			ci.instanceMethodExtraInitializersName != nil &&
			ci.hasStaticInitializers &&
			ci.hasNonAmbientInstanceFields &&
			ci.hasStaticPrivateClassElements {
			break
		}
	}

	return ci
}

// --- transformClassLike: the main transformation workhorse ---

func (tx *esDecoratorTransformer) transformClassLike(node *ast.Node) *ast.Expression {
	f := tx.Factory()
	ec := tx.EmitContext()

	ec.StartVariableEnvironment()

	// When a class has class decorators, if it doesn't have an assigned name, give it one of "".
	if !classHasDeclaredOrExplicitlyAssignedName(ec, node) && ast.ClassOrConstructorParameterIsDecorated(false, node) {
		node = injectClassNamedEvaluationHelperBlockIfMissing(ec, node, f.NewStringLiteral("", 0), nil)
	}

	classReference := f.GetLocalNameEx(node, printer.AssignedNameOptions{})
	ci := tx.createClassInfo(node)
	classDefinitionStatements := []*ast.Statement{}
	var leadingBlockStatements []*ast.Statement
	var trailingBlockStatements []*ast.Statement
	var syntheticConstructor *ast.Node
	var heritageClauses *ast.NodeList

	// 1. Class decorators are evaluated outside the private name scope of the class
	classDecorators := tx.transformAllDecoratorsOfDeclaration(getAllDecoratorsOfClass(node))
	if len(classDecorators) > 0 {
		ci.classDecoratorsName = f.NewUniqueNameEx("_classDecorators", printer.AutoGenerateOptions{
			Flags: printer.GeneratedIdentifierFlagsOptimistic | printer.GeneratedIdentifierFlagsFileLevel,
		})
		ci.classDescriptorName = f.NewUniqueNameEx("_classDescriptor", printer.AutoGenerateOptions{
			Flags: printer.GeneratedIdentifierFlagsOptimistic | printer.GeneratedIdentifierFlagsFileLevel,
		})
		ci.classExtraInitializersName = f.NewUniqueNameEx("_classExtraInitializers", printer.AutoGenerateOptions{
			Flags: printer.GeneratedIdentifierFlagsOptimistic | printer.GeneratedIdentifierFlagsFileLevel,
		})

		decoratorsArray := f.NewArrayLiteralExpression(
			f.NewNodeList(expressionsToNodes(classDecorators)),
			false,
		)
		classDefinitionStatements = append(classDefinitionStatements,
			tx.createLet(ci.classDecoratorsName, decoratorsArray),
			tx.createLet(ci.classDescriptorName, nil),
			tx.createLet(ci.classExtraInitializersName, f.NewArrayLiteralExpression(f.NewNodeList(nil), false)),
			tx.createLet(ci.classThis, nil),
		)
	}

	// Rewrite super in static initializers
	extendsClause := ast.GetHeritageClause(node, ast.KindExtendsKeyword)
	var extendsElement *ast.Node
	if extendsClause != nil {
		hc := extendsClause.AsHeritageClause()
		if hc.Types != nil && len(hc.Types.Nodes) > 0 {
			extendsElement = hc.Types.Nodes[0]
		}
	}
	var extendsExpression *ast.Expression
	if extendsElement != nil {
		extendsExpression = tx.Visitor().VisitNode(extendsElement.AsExpressionWithTypeArguments().Expression)
	}

	if extendsExpression != nil {
		ci.classSuper = f.NewUniqueNameEx("_classSuper", printer.AutoGenerateOptions{
			Flags: printer.GeneratedIdentifierFlagsOptimistic | printer.GeneratedIdentifierFlagsFileLevel,
		})

		// Ensure we do not give the class or function an assigned name
		unwrapped := ast.SkipOuterExpressions(extendsExpression, ast.OEKAll)
		safeExtendsExpression := extendsExpression
		if (ast.IsClassExpression(unwrapped) && unwrapped.Name() == nil) ||
			(ast.IsFunctionExpression(unwrapped) && unwrapped.Name() == nil) ||
			ast.IsArrowFunction(unwrapped) {
			safeExtendsExpression = f.NewCommaExpression(
				f.NewNumericLiteral("0", 0),
				extendsExpression,
			)
		}
		classDefinitionStatements = append(classDefinitionStatements, tx.createLet(ci.classSuper, safeExtendsExpression))

		updatedExtendsElement := f.UpdateExpressionWithTypeArguments(extendsElement.AsExpressionWithTypeArguments(), ci.classSuper, nil)
		updatedExtendsClause := f.UpdateHeritageClause(extendsClause.AsHeritageClause(), f.NewNodeList([]*ast.Node{updatedExtendsElement}))
		heritageClauses = f.NewNodeList([]*ast.Node{updatedExtendsClause})
	}

	var renamedClassThis *ast.Node
	if ci.classThis != nil {
		renamedClassThis = ci.classThis
	} else {
		renamedClassThis = f.NewThisExpression()
	}

	// Visit each member
	tx.enterClass(ci)

	leadingBlockStatements = append(leadingBlockStatements, tx.createMetadata(ci.metadataReference, ci.classSuper))

	// Visit non-constructor members first, then constructor
	// We visit in two passes to ensure instance field initializers are resolved before the constructor.
	nonConstructorVisitor := tx.EmitContext().NewNodeVisitor(func(n *ast.Node) *ast.Node {
		if ast.IsConstructorDeclaration(n) {
			return n // skip constructors in pass 1
		}
		return tx.classElementVisitor(n)
	})
	members := nonConstructorVisitor.VisitNodes(node.MemberList())
	constructorVisitor := tx.EmitContext().NewNodeVisitor(func(n *ast.Node) *ast.Node {
		if ast.IsConstructorDeclaration(n) {
			return tx.classElementVisitor(n)
		}
		return n
	})
	members = constructorVisitor.VisitNodes(members)

	// Handle pending expressions (computed property names and decorator evaluations)
	if len(tx.pendingExpressions) > 0 {
		var outerThis *ast.IdentifierNode
		for _, expr := range tx.pendingExpressions {
			// If a pending expression contains lexical `this`, capture it
			if expr.SubtreeFacts()&ast.SubtreeContainsLexicalThis != 0 {
				var thisVisitor *ast.NodeVisitor
				thisVisitor = ec.NewNodeVisitor(func(n *ast.Node) *ast.Node {
					if n.SubtreeFacts()&ast.SubtreeContainsLexicalThis == 0 && n.Kind != ast.KindThisKeyword {
						return n
					}
					if n.Kind == ast.KindThisKeyword {
						if outerThis == nil {
							outerThis = f.NewUniqueNameEx("_outerThis", printer.AutoGenerateOptions{
								Flags: printer.GeneratedIdentifierFlagsOptimistic,
							})
							classDefinitionStatements = append(
								[]*ast.Statement{tx.createLet(outerThis, f.NewThisExpression())},
								classDefinitionStatements...,
							)
						}
						return outerThis
					}
					return thisVisitor.VisitEachChild(n)
				})
				expr = thisVisitor.VisitNode(expr)
			}
			statement := f.NewExpressionStatement(expr)
			leadingBlockStatements = append(leadingBlockStatements, statement)
		}
		tx.pendingExpressions = nil
	}
	tx.exitClass()

	// If there are instance initializers but no constructor, synthesize one
	if len(ci.pendingInstanceInitializers) > 0 && ast.GetFirstConstructorWithBody(node) == nil {
		initializerStatements := tx.prepareConstructor(ci)
		if len(initializerStatements) > 0 {
			isDerivedClass := extendsClause != nil
			constructorStatements := []*ast.Statement{}
			if isDerivedClass {
				spreadArguments := f.NewSpreadElement(f.NewIdentifier("arguments"))
				superCall := f.NewCallExpression(f.NewKeywordExpression(ast.KindSuperKeyword), nil, nil, f.NewNodeList([]*ast.Expression{spreadArguments}), ast.NodeFlagsNone)
				constructorStatements = append(constructorStatements, f.NewExpressionStatement(superCall))
			}
			constructorStatements = append(constructorStatements, initializerStatements...)
			constructorBody := f.NewBlock(f.NewNodeList(statementsToNodes(constructorStatements)), true)
			syntheticConstructor = f.NewConstructorDeclaration(nil, nil, f.NewNodeList(nil), nil, nil, constructorBody)
		}
	}

	// Used in class definition steps 5,7,11
	if ci.staticMethodExtraInitializersName != nil {
		classDefinitionStatements = append(classDefinitionStatements,
			tx.createLet(ci.staticMethodExtraInitializersName, f.NewArrayLiteralExpression(f.NewNodeList(nil), false)),
		)
	}

	// Used in class definition steps 6,8, and construction
	if ci.instanceMethodExtraInitializersName != nil {
		classDefinitionStatements = append(classDefinitionStatements,
			tx.createLet(ci.instanceMethodExtraInitializersName, f.NewArrayLiteralExpression(f.NewNodeList(nil), false)),
		)
	}

	// Emit member info variable declarations
	// The reference implementation emits static member vars first, then non-static
	if ci.memberInfos != nil {
		for _, member := range node.Members() {
			if !ast.IsStatic(member) {
				continue
			}
			mi, ok := ci.memberInfos[member]
			if !ok {
				continue
			}
			classDefinitionStatements = append(classDefinitionStatements, tx.createLet(mi.memberDecoratorsName, nil))
			if mi.memberInitializersName != nil {
				classDefinitionStatements = append(classDefinitionStatements, tx.createLet(mi.memberInitializersName, f.NewArrayLiteralExpression(f.NewNodeList(nil), false)))
			}
			if mi.memberExtraInitializersName != nil {
				classDefinitionStatements = append(classDefinitionStatements, tx.createLet(mi.memberExtraInitializersName, f.NewArrayLiteralExpression(f.NewNodeList(nil), false)))
			}
			if mi.memberDescriptorName != nil {
				classDefinitionStatements = append(classDefinitionStatements, tx.createLet(mi.memberDescriptorName, nil))
			}
		}
		for _, member := range node.Members() {
			if ast.IsStatic(member) {
				continue
			}
			mi, ok := ci.memberInfos[member]
			if !ok {
				continue
			}
			classDefinitionStatements = append(classDefinitionStatements, tx.createLet(mi.memberDecoratorsName, nil))
			if mi.memberInitializersName != nil {
				classDefinitionStatements = append(classDefinitionStatements, tx.createLet(mi.memberInitializersName, f.NewArrayLiteralExpression(f.NewNodeList(nil), false)))
			}
			if mi.memberExtraInitializersName != nil {
				classDefinitionStatements = append(classDefinitionStatements, tx.createLet(mi.memberExtraInitializersName, f.NewArrayLiteralExpression(f.NewNodeList(nil), false)))
			}
			if mi.memberDescriptorName != nil {
				classDefinitionStatements = append(classDefinitionStatements, tx.createLet(mi.memberDescriptorName, nil))
			}
		}
	}

	// 5. Static non-field element decorators are applied
	leadingBlockStatements = append(leadingBlockStatements, ci.staticNonFieldDecorationStatements...)

	// 6. Non-static non-field element decorators are applied
	leadingBlockStatements = append(leadingBlockStatements, ci.nonStaticNonFieldDecorationStatements...)

	// 7. Static field element decorators are applied
	leadingBlockStatements = append(leadingBlockStatements, ci.staticFieldDecorationStatements...)

	// 8. Non-static field element decorators are applied
	leadingBlockStatements = append(leadingBlockStatements, ci.nonStaticFieldDecorationStatements...)

	// 9. Class decorators are applied
	// 10. Class binding is initialized
	if ci.classDescriptorName != nil && ci.classDecoratorsName != nil && ci.classExtraInitializersName != nil && ci.classThis != nil {
		valueProperty := f.NewPropertyAssignment(nil, f.NewIdentifier("value"), nil, nil, renamedClassThis)
		classDescriptor := f.NewObjectLiteralExpression(f.NewNodeList([]*ast.Node{valueProperty}), false)
		classDescriptorAssignment := f.NewAssignmentExpression(ci.classDescriptorName, classDescriptor)
		classNameReference := f.NewPropertyAccessExpression(renamedClassThis, nil, f.NewIdentifier("name"), ast.NodeFlagsNone)

		contextObj := tx.createESDecorateClassContext(classNameReference, ci.metadataReference)

		esDecorateHelper := f.NewESDecorateHelper(
			f.NewToken(ast.KindNullKeyword),
			classDescriptorAssignment,
			ci.classDecoratorsName,
			contextObj,
			f.NewToken(ast.KindNullKeyword),
			ci.classExtraInitializersName,
		)
		esDecorateStatement := f.NewExpressionStatement(esDecorateHelper)
		ec.SetSourceMapRange(esDecorateStatement, moveRangePastDecorators(node))
		leadingBlockStatements = append(leadingBlockStatements, esDecorateStatement)

		// C = _classThis = _classDescriptor.value
		classDescriptorValueRef := f.NewPropertyAccessExpression(ci.classDescriptorName, nil, f.NewIdentifier("value"), ast.NodeFlagsNone)
		classThisAssignment := f.NewAssignmentExpression(ci.classThis, classDescriptorValueRef)
		classReferenceAssignment := f.NewAssignmentExpression(classReference, classThisAssignment)
		leadingBlockStatements = append(leadingBlockStatements, f.NewExpressionStatement(classReferenceAssignment))
	}

	// Symbol.metadata assignment
	leadingBlockStatements = append(leadingBlockStatements, tx.createSymbolMetadata(renamedClassThis, ci.metadataReference))

	// 11. Static extra initializers
	// 12. Static fields are initialized
	if len(ci.pendingStaticInitializers) > 0 {
		for _, initializer := range ci.pendingStaticInitializers {
			initializerStatement := f.NewExpressionStatement(initializer)
			ec.SetSourceMapRange(initializerStatement, ec.SourceMapRange(initializer))
			trailingBlockStatements = append(trailingBlockStatements, initializerStatement)
		}
		ci.pendingStaticInitializers = nil
	}

	// 13. Class extra initializers
	if ci.classExtraInitializersName != nil {
		runClassInitializersHelper := f.NewRunInitializersHelper(renamedClassThis, ci.classExtraInitializersName, nil)
		runClassInitializersStatement := f.NewExpressionStatement(runClassInitializersHelper)
		if node.Name() != nil {
			ec.SetSourceMapRange(runClassInitializersStatement, node.Name().Loc)
		} else {
			ec.SetSourceMapRange(runClassInitializersStatement, moveRangePastDecorators(node))
		}
		trailingBlockStatements = append(trailingBlockStatements, runClassInitializersStatement)
	}

	// If there are no static initializers, combine leading and trailing
	if len(leadingBlockStatements) > 0 && len(trailingBlockStatements) > 0 && !ci.hasStaticInitializers {
		leadingBlockStatements = append(leadingBlockStatements, trailingBlockStatements...)
		trailingBlockStatements = nil
	}

	// Create leading static block
	var leadingStaticBlock *ast.Node
	if len(leadingBlockStatements) > 0 {
		leadingStaticBlock = f.NewClassStaticBlockDeclaration(
			nil,
			f.NewBlock(f.NewNodeList(statementsToNodes(leadingBlockStatements)), true),
		)
	}

	// Create trailing static block
	var trailingStaticBlock *ast.Node
	if len(trailingBlockStatements) > 0 {
		trailingStaticBlock = f.NewClassStaticBlockDeclaration(
			nil,
			f.NewBlock(f.NewNodeList(statementsToNodes(trailingBlockStatements)), true),
		)
	}

	// Assemble new members list
	if leadingStaticBlock != nil || syntheticConstructor != nil || trailingStaticBlock != nil {
		newMembers := make([]*ast.Node, 0, len(members.Nodes)+3)

		// Find the existing NamedEvaluation helper block index
		existingNamedEvaluationHelperBlockIndex := -1
		for i, m := range members.Nodes {
			if isClassNamedEvaluationHelperBlock(ec, m) {
				existingNamedEvaluationHelperBlockIndex = i
				break
			}
		}

		if leadingStaticBlock != nil {
			newMembers = append(newMembers, members.Nodes[:existingNamedEvaluationHelperBlockIndex+1]...)
			newMembers = append(newMembers, leadingStaticBlock)
			newMembers = append(newMembers, members.Nodes[existingNamedEvaluationHelperBlockIndex+1:]...)
		} else {
			newMembers = append(newMembers, members.Nodes...)
		}

		if syntheticConstructor != nil {
			newMembers = append(newMembers, syntheticConstructor)
		}

		if trailingStaticBlock != nil {
			newMembers = append(newMembers, trailingStaticBlock)
		}

		membersList := f.NewNodeList(newMembers)
		membersList.Loc = members.Loc
		members = membersList
	}

	lexicalEnvironment := ec.EndVariableEnvironment()

	var classExpression *ast.Node
	if len(classDecorators) > 0 {
		classExpression = f.NewClassExpression(nil, nil, nil, heritageClauses, members)
		ec.SetOriginal(classExpression, node)
		if ci.classThis != nil {
			classExpression = injectClassThisAssignmentIfMissing(ec, f, classExpression, ci.classThis)
		}

		classReferenceDeclaration := f.NewVariableDeclaration(classReference, nil, nil, classExpression)
		classReferenceVarDeclList := f.NewVariableDeclarationList(ast.NodeFlagsNone, f.NewNodeList([]*ast.Node{classReferenceDeclaration}))
		var returnExpr *ast.Expression
		if ci.classThis != nil {
			returnExpr = f.NewAssignmentExpression(classReference, ci.classThis)
		} else {
			returnExpr = classReference
		}
		classDefinitionStatements = append(classDefinitionStatements,
			f.NewVariableStatement(nil, classReferenceVarDeclList),
			f.NewReturnStatement(returnExpr),
		)
	} else {
		classExpression = f.NewClassExpression(nil, node.Name(), nil, heritageClauses, members)
		ec.SetOriginal(classExpression, node)
		classDefinitionStatements = append(classDefinitionStatements, f.NewReturnStatement(classExpression))
	}

	mergedStatements := ec.MergeEnvironment(classDefinitionStatements, lexicalEnvironment)
	return f.NewImmediatelyInvokedArrowFunction(mergedStatements)
}

// --- visitClassDeclaration ---

func (tx *esDecoratorTransformer) visitClassDeclaration(node *ast.ClassDeclaration) *ast.Node {
	if isDecoratedClassLike(node.AsNode()) {
		f := tx.Factory()
		ec := tx.EmitContext()
		statements := []*ast.Statement{}

		originalClass := ec.MostOriginal(node.AsNode())
		if !ast.IsClassLike(originalClass) {
			originalClass = node.AsNode()
		}
		var className *ast.Expression
		if originalClass.Name() != nil {
			className = f.NewStringLiteralFromNode(originalClass.Name())
		} else {
			className = f.NewStringLiteral("default", 0)
		}

		isExport := ast.HasSyntacticModifier(node.AsNode(), ast.ModifierFlagsExport)
		isDefault := ast.HasSyntacticModifier(node.AsNode(), ast.ModifierFlagsDefault)

		classNode := node.AsNode()
		if node.Name() == nil {
			classNode = injectClassNamedEvaluationHelperBlockIfMissing(ec, classNode, className, nil)
		}

		if isExport && isDefault {
			iife := tx.transformClassLike(classNode)
			if classNode.Name() != nil {
				varDecl := f.NewVariableDeclaration(f.GetLocalName(classNode), nil, nil, iife)
				ec.SetOriginal(varDecl, classNode)
				varDecls := f.NewVariableDeclarationList(ast.NodeFlagsLet, f.NewNodeList([]*ast.Node{varDecl}))
				varStatement := f.NewVariableStatement(nil, varDecls)
				statements = append(statements, varStatement)

				exportStatement := tx.createExportDefault(f.GetDeclarationName(classNode))
				ec.SetOriginal(exportStatement, classNode)
				ec.AssignCommentRange(exportStatement, classNode)
				ec.SetSourceMapRange(exportStatement, moveRangePastDecorators(classNode))
				statements = append(statements, exportStatement)
			} else {
				exportStatement := tx.createExportDefault(iife)
				ec.SetOriginal(exportStatement, classNode)
				ec.AssignCommentRange(exportStatement, classNode)
				ec.SetSourceMapRange(exportStatement, moveRangePastDecorators(classNode))
				statements = append(statements, exportStatement)
			}
		} else {
			iife := tx.transformClassLike(classNode)
			modifierVisitor := ec.NewNodeVisitor(func(n *ast.Node) *ast.Node {
				if isExport && n.Kind == ast.KindExportKeyword {
					return nil
				}
				return tx.modifierVisitor(n)
			})
			modifiers := modifierVisitor.VisitModifiers(classNode.Modifiers())

			declName := f.GetLocalNameEx(classNode, printer.AssignedNameOptions{AllowSourceMaps: true})
			varDecl := f.NewVariableDeclaration(declName, nil, nil, iife)
			ec.SetOriginal(varDecl, classNode)
			varDecls := f.NewVariableDeclarationList(ast.NodeFlagsLet, f.NewNodeList([]*ast.Node{varDecl}))
			varStatement := f.NewVariableStatement(modifiers, varDecls)
			ec.SetOriginal(varStatement, classNode)
			ec.AssignCommentRange(varStatement, classNode)
			statements = append(statements, varStatement)

			if isExport {
				exportStatement := tx.createExternalModuleExport(declName)
				ec.SetOriginal(exportStatement, classNode)
				statements = append(statements, exportStatement)
			}
		}

		if len(statements) == 1 {
			return statements[0]
		}
		return transformers.SingleOrMany(statementsToNodes(statements), f)
	}

	// Non-decorated class
	modVisitor := tx.EmitContext().NewNodeVisitor(tx.modifierVisitor)
	modifiers := modVisitor.VisitModifiers(node.Modifiers())
	hcVisitor := tx.EmitContext().NewNodeVisitor(tx.visit)
	heritageClauses := hcVisitor.VisitNodes(node.HeritageClauses)
	tx.enterClass(nil)
	ceVisitor := tx.EmitContext().NewNodeVisitor(tx.classElementVisitor)
	members := ceVisitor.VisitNodes(node.Members)
	tx.exitClass()
	return tx.Factory().UpdateClassDeclaration(node, modifiers, node.Name(), nil, heritageClauses, members)
}

// --- visitClassExpression ---

func (tx *esDecoratorTransformer) visitClassExpression(node *ast.ClassExpression) *ast.Node {
	if isDecoratedClassLike(node.AsNode()) {
		iife := tx.transformClassLike(node.AsNode())
		tx.EmitContext().SetOriginal(iife, node.AsNode())
		return iife
	}

	modVisitor := tx.EmitContext().NewNodeVisitor(tx.modifierVisitor)
	modifiers := modVisitor.VisitModifiers(node.Modifiers())
	hcVisitor := tx.EmitContext().NewNodeVisitor(tx.visit)
	heritageClauses := hcVisitor.VisitNodes(node.HeritageClauses)
	tx.enterClass(nil)
	ceVisitor := tx.EmitContext().NewNodeVisitor(tx.classElementVisitor)
	members := ceVisitor.VisitNodes(node.Members)
	tx.exitClass()
	return tx.Factory().UpdateClassExpression(node, modifiers, node.Name(), nil, heritageClauses, members)
}

// --- prepareConstructor ---

func (tx *esDecoratorTransformer) prepareConstructor(ci *classInfo) []*ast.Statement {
	if len(ci.pendingInstanceInitializers) == 0 {
		return nil
	}
	f := tx.Factory()
	statements := []*ast.Statement{
		f.NewExpressionStatement(f.InlineExpressions(ci.pendingInstanceInitializers)),
	}
	ci.pendingInstanceInitializers = nil
	return statements
}

// --- visitConstructorDeclaration ---

func (tx *esDecoratorTransformer) visitConstructorDeclaration(node *ast.Node) *ast.Node {
	tx.enterClassElement(node)
	modVisitor := tx.EmitContext().NewNodeVisitor(tx.modifierVisitor)
	modifiers := modVisitor.VisitModifiers(node.Modifiers())
	paramVisitor := tx.EmitContext().NewNodeVisitor(tx.visit)
	parameters := paramVisitor.VisitNodes(node.ParameterList())

	var body *ast.Node
	ctor := node.AsConstructorDeclaration()
	if ctor.Body != nil && tx.classInfoStack != nil {
		initializerStatements := tx.prepareConstructor(tx.classInfoStack)
		if len(initializerStatements) > 0 {
			stmts := []*ast.Statement{}
			prologue, rest := tx.Factory().SplitStandardPrologue(ctor.Body.AsBlock().Statements.Nodes)
			stmts = append(stmts, prologue...)

			superStatementIndices := findSuperStatementIndexPath(rest, 0)
			if len(superStatementIndices) > 0 {
				tx.transformConstructorBodyWorker(&stmts, rest, 0, superStatementIndices, 0, initializerStatements)
			} else {
				stmts = append(stmts, initializerStatements...)
				stmtVisitor := tx.EmitContext().NewNodeVisitor(tx.visit)
				visited, _ := stmtVisitor.VisitSlice(nodesToStatements(rest))
				stmts = append(stmts, visited...)
			}

			body = tx.Factory().NewBlock(tx.Factory().NewNodeList(statementsToNodes(stmts)), true)
			tx.EmitContext().SetOriginal(body, ctor.Body.AsNode())
			body.Loc = ctor.Body.Loc
		}
	}

	if body == nil {
		body = tx.Visitor().VisitNode(ctor.Body.AsNode())
	}
	tx.exitClassElement()
	return tx.Factory().UpdateConstructorDeclaration(ctor, modifiers, nil, parameters, nil, nil, body)
}

func (tx *esDecoratorTransformer) transformConstructorBodyWorker(statementsOut *[]*ast.Statement, statementsIn []*ast.Statement, statementOffset int, superPath []int, superPathDepth int, initializerStatements []*ast.Statement) {
	superStatementIndex := superPath[superPathDepth]
	// Visit statements before super
	if superStatementIndex > statementOffset {
		stmtVisitor := tx.EmitContext().NewNodeVisitor(tx.visit)
		for _, s := range statementsIn[statementOffset:superStatementIndex] {
			*statementsOut = append(*statementsOut, stmtVisitor.VisitNode(s))
		}
	}

	superStatement := statementsIn[superStatementIndex]
	if ast.IsTryStatement(superStatement) {
		// Recurse into try block
		tryBlockNode := superStatement.AsTryStatement().TryBlock
		tryBlock := tryBlockNode.AsBlock()
		tryBlockStatements := []*ast.Statement{}
		tx.transformConstructorBodyWorker(&tryBlockStatements, tryBlock.Statements.Nodes, 0, superPath, superPathDepth+1, initializerStatements)

		newTryBlock := tx.Factory().NewBlock(tx.Factory().NewNodeList(statementsToNodes(tryBlockStatements)), true)
		newTryBlock.Loc = tryBlockNode.Loc

		var catchClause *ast.Node
		if superStatement.AsTryStatement().CatchClause != nil {
			catchClause = tx.Visitor().VisitNode(superStatement.AsTryStatement().CatchClause)
		}
		var finallyBlock *ast.Node
		if superStatement.AsTryStatement().FinallyBlock != nil {
			finallyBlock = tx.Visitor().VisitNode(superStatement.AsTryStatement().FinallyBlock)
		}
		updated := tx.Factory().UpdateTryStatement(superStatement.AsTryStatement(), newTryBlock, catchClause, finallyBlock)
		*statementsOut = append(*statementsOut, updated)
	} else {
		stmtVisitor := tx.EmitContext().NewNodeVisitor(tx.visit)
		*statementsOut = append(*statementsOut, stmtVisitor.VisitNode(superStatement))
		*statementsOut = append(*statementsOut, initializerStatements...)
	}

	// Visit statements after super
	if superStatementIndex+1 < len(statementsIn) {
		stmtVisitor := tx.EmitContext().NewNodeVisitor(tx.visit)
		for _, s := range statementsIn[superStatementIndex+1:] {
			*statementsOut = append(*statementsOut, stmtVisitor.VisitNode(s))
		}
	}
}

// --- finishClassElement ---

func (tx *esDecoratorTransformer) finishClassElement(updated *ast.Node, original *ast.Node) *ast.Node {
	if updated != original {
		tx.EmitContext().AssignCommentRange(updated, original)
		tx.EmitContext().SetSourceMapRange(updated, moveRangePastDecorators(original))
	}
	return updated
}

// --- partialTransformClassElement ---

type partialResult struct {
	modifiers             *ast.ModifierList
	referencedName        *ast.Expression
	name                  *ast.Node
	initializersName      *ast.IdentifierNode
	extraInitializersName *ast.IdentifierNode
	descriptorName        *ast.IdentifierNode
	thisArg               *ast.IdentifierNode
}

type createDescriptorFunc func(member *ast.Node, modifiers *ast.ModifierList) *ast.Expression

func (tx *esDecoratorTransformer) partialTransformClassElement(member *ast.Node, ci *classInfo, createDescriptor createDescriptorFunc) partialResult {
	f := tx.Factory()
	ec := tx.EmitContext()

	if ci == nil {
		modVisitor := ec.NewNodeVisitor(tx.modifierVisitor)
		modifiers := modVisitor.VisitModifiers(member.Modifiers())
		tx.enterName()
		name := tx.visitPropertyName(member.Name())
		tx.exitName()
		return partialResult{modifiers: modifiers, name: name}
	}

	// Collect decorators for this member
	memberDecorators := tx.transformAllDecoratorsOfDeclaration(getAllDecoratorsOfClassElement(member, ci.class))
	modVisitor := ec.NewNodeVisitor(tx.modifierVisitor)
	modifiers := modVisitor.VisitModifiers(member.Modifiers())

	var result partialResult
	result.modifiers = modifiers

	if len(memberDecorators) > 0 {
		memberDecoratorsName := tx.createHelperVariable(member, "decorators")
		memberDecoratorsArray := f.NewArrayLiteralExpression(
			f.NewNodeList(expressionsToNodes(memberDecorators)),
			false,
		)
		memberDecoratorsAssignment := f.NewAssignmentExpression(memberDecoratorsName, memberDecoratorsArray)
		mi := &memberInfo{memberDecoratorsName: memberDecoratorsName}
		if ci.memberInfos == nil {
			ci.memberInfos = make(map[*ast.Node]*memberInfo)
		}
		ci.memberInfos[member] = mi
		tx.pendingExpressions = append(tx.pendingExpressions, memberDecoratorsAssignment)

		// Determine which decoration statement bucket to use
		var statements *[]*ast.Statement
		if ast.IsMethodOrAccessor(member) || ast.IsAutoAccessorPropertyDeclaration(member) {
			if ast.IsStatic(member) {
				statements = &ci.staticNonFieldDecorationStatements
			} else {
				statements = &ci.nonStaticNonFieldDecorationStatements
			}
		} else if ast.IsPropertyDeclaration(member) && !ast.IsAutoAccessorPropertyDeclaration(member) {
			if ast.IsStatic(member) {
				statements = &ci.staticFieldDecorationStatements
			} else {
				statements = &ci.nonStaticFieldDecorationStatements
			}
		}

		// Determine decorator kind
		kind := ""
		switch {
		case ast.IsGetAccessorDeclaration(member):
			kind = "getter"
		case ast.IsSetAccessorDeclaration(member):
			kind = "setter"
		case ast.IsMethodDeclaration(member):
			kind = "method"
		case ast.IsAutoAccessorPropertyDeclaration(member):
			kind = "accessor"
		case ast.IsPropertyDeclaration(member):
			kind = "field"
		}

		// Determine the property name for the context
		var propertyNameComputed bool
		var propertyNameExpr *ast.Expression
		if member.Name() != nil && (ast.IsIdentifier(member.Name()) || ast.IsPrivateIdentifier(member.Name())) {
			propertyNameComputed = false
			propertyNameExpr = member.Name()
		} else if member.Name() != nil && ast.IsPropertyNameLiteral(member.Name()) {
			propertyNameComputed = true
			propertyNameExpr = f.NewStringLiteralFromNode(member.Name())
		} else if member.Name() != nil && ast.IsComputedPropertyName(member.Name()) {
			cpn := member.Name().AsComputedPropertyName()
			if ast.IsPropertyNameLiteral(cpn.Expression) && !ast.IsIdentifier(cpn.Expression) {
				propertyNameComputed = true
				propertyNameExpr = f.NewStringLiteralFromNode(cpn.Expression)
			} else {
				tx.enterName()
				result.referencedName, result.name = tx.visitReferencedPropertyName(member.Name())
				tx.exitName()
				propertyNameComputed = true
				propertyNameExpr = result.referencedName
			}
		}

		contextObj := tx.createESDecorateElementContext(
			kind,
			propertyNameComputed,
			propertyNameExpr,
			ast.IsStatic(member),
			member.Name() != nil && ast.IsPrivateIdentifier(member.Name()),
			ast.IsPropertyDeclaration(member) || ast.IsGetAccessorDeclaration(member) || ast.IsMethodDeclaration(member),
			ast.IsPropertyDeclaration(member) || ast.IsSetAccessorDeclaration(member),
			ci.metadataReference,
		)

		if ast.IsMethodOrAccessor(member) {
			methodExtraInitializersName := ci.instanceMethodExtraInitializersName
			if ast.IsStatic(member) {
				methodExtraInitializersName = ci.staticMethodExtraInitializersName
			}

			var descriptorArg *ast.Expression
			if ast.IsPrivateIdentifierClassElementDeclaration(member) && createDescriptor != nil {
				// For private members, extract the method/accessor body into a descriptor object.
				// Filter modifiers to only keep async.
				asyncMods := tx.filterAsyncModifier(modifiers)
				descriptor := createDescriptor(member, asyncMods)
				mi.memberDescriptorName = tx.createHelperVariable(member, "descriptor")
				result.descriptorName = mi.memberDescriptorName
				descriptorArg = f.NewAssignmentExpression(mi.memberDescriptorName, descriptor)
			} else {
				descriptorArg = f.NewToken(ast.KindNullKeyword)
			}

			esDecorateExpr := f.NewESDecorateHelper(
				f.NewThisExpression(),
				descriptorArg,
				memberDecoratorsName,
				contextObj,
				f.NewToken(ast.KindNullKeyword),
				methodExtraInitializersName,
			)
			esDecorateStatement := f.NewExpressionStatement(esDecorateExpr)
			ec.SetSourceMapRange(esDecorateStatement, moveRangePastDecorators(member))
			if statements != nil {
				*statements = append(*statements, esDecorateStatement)
			}
		} else if ast.IsPropertyDeclaration(member) {
			mi.memberInitializersName = tx.createHelperVariable(member, "initializers")
			mi.memberExtraInitializersName = tx.createHelperVariable(member, "extraInitializers")
			result.initializersName = mi.memberInitializersName
			result.extraInitializersName = mi.memberExtraInitializersName
			if ast.IsStatic(member) {
				result.thisArg = ci.classThis
			}

			ctorArg := (*ast.Node)(nil)
			if ast.IsAutoAccessorPropertyDeclaration(member) {
				ctorArg = f.NewThisExpression()
			} else {
				ctorArg = f.NewToken(ast.KindNullKeyword)
			}

			esDecorateExpr := f.NewESDecorateHelper(
				ctorArg,
				f.NewToken(ast.KindNullKeyword),
				memberDecoratorsName,
				contextObj,
				mi.memberInitializersName,
				mi.memberExtraInitializersName,
			)
			esDecorateStatement := f.NewExpressionStatement(esDecorateExpr)
			ec.SetSourceMapRange(esDecorateStatement, moveRangePastDecorators(member))
			if statements != nil {
				*statements = append(*statements, esDecorateStatement)
			}
		}
	}

	if result.name == nil {
		tx.enterName()
		result.name = tx.visitPropertyName(member.Name())
		tx.exitName()
	}

	if modifiers == nil && (ast.IsMethodDeclaration(member) || ast.IsPropertyDeclaration(member)) {
		ec.SetEmitFlags(result.name, printer.EFNoLeadingComments)
	}

	return result
}

// --- Class member visitors ---

func (tx *esDecoratorTransformer) visitMethodDeclaration(node *ast.Node) *ast.Node {
	tx.enterClassElement(node)
	result := tx.partialTransformClassElement(node, tx.classInfoStack, tx.createMethodDescriptorObject)
	if result.descriptorName != nil {
		tx.exitClassElement()
		return tx.finishClassElement(tx.createMethodDescriptorForwarder(result.modifiers, result.name, result.descriptorName), node)
	}
	paramVisitor := tx.EmitContext().NewNodeVisitor(tx.visit)
	parameters := paramVisitor.VisitNodes(node.ParameterList())
	body := tx.Visitor().VisitNode(node.Body())
	tx.exitClassElement()
	method := node.AsMethodDeclaration()
	return tx.finishClassElement(
		tx.Factory().UpdateMethodDeclaration(method, result.modifiers, method.AsteriskToken, result.name, nil, nil, parameters, nil, nil, body),
		node,
	)
}

func (tx *esDecoratorTransformer) visitGetAccessorDeclaration(node *ast.Node) *ast.Node {
	tx.enterClassElement(node)
	result := tx.partialTransformClassElement(node, tx.classInfoStack, tx.createGetAccessorDescriptorObject)
	if result.descriptorName != nil {
		tx.exitClassElement()
		return tx.finishClassElement(tx.createGetAccessorDescriptorForwarder(result.modifiers, result.name, result.descriptorName), node)
	}
	paramVisitor := tx.EmitContext().NewNodeVisitor(tx.visit)
	parameters := paramVisitor.VisitNodes(node.ParameterList())
	body := tx.Visitor().VisitNode(node.Body())
	tx.exitClassElement()
	accessor := node.AsGetAccessorDeclaration()
	return tx.finishClassElement(
		tx.Factory().UpdateGetAccessorDeclaration(accessor, result.modifiers, result.name, nil, parameters, nil, nil, body),
		node,
	)
}

func (tx *esDecoratorTransformer) visitSetAccessorDeclaration(node *ast.Node) *ast.Node {
	tx.enterClassElement(node)
	result := tx.partialTransformClassElement(node, tx.classInfoStack, tx.createSetAccessorDescriptorObject)
	if result.descriptorName != nil {
		tx.exitClassElement()
		return tx.finishClassElement(tx.createSetAccessorDescriptorForwarder(result.modifiers, result.name, result.descriptorName), node)
	}
	paramVisitor := tx.EmitContext().NewNodeVisitor(tx.visit)
	parameters := paramVisitor.VisitNodes(node.ParameterList())
	body := tx.Visitor().VisitNode(node.Body())
	tx.exitClassElement()
	accessor := node.AsSetAccessorDeclaration()
	return tx.finishClassElement(
		tx.Factory().UpdateSetAccessorDeclaration(accessor, result.modifiers, result.name, nil, parameters, nil, nil, body),
		node,
	)
}

func (tx *esDecoratorTransformer) visitPropertyDeclaration(node *ast.Node) *ast.Node {
	if isNamedEvaluation(tx.EmitContext(), node) && isAnonymousClassNeedingAssignedName(node.Initializer()) {
		node = transformNamedEvaluation(tx.EmitContext(), node, false, "")
	}

	tx.enterClassElement(node)
	f := tx.Factory()
	ec := tx.EmitContext()

	result := tx.partialTransformClassElement(node, tx.classInfoStack, nil)

	ec.StartVariableEnvironment()

	initializer := tx.Visitor().VisitNode(node.Initializer())
	if result.initializersName != nil {
		var thisArg *ast.Node
		if result.thisArg != nil {
			thisArg = result.thisArg
		} else {
			thisArg = f.NewThisExpression()
		}
		if initializer == nil {
			initializer = f.NewVoidZeroExpression()
		}
		initializer = f.NewRunInitializersHelper(thisArg, result.initializersName, initializer)
	}

	if ast.IsStatic(node) && tx.classInfoStack != nil && initializer != nil {
		tx.classInfoStack.hasStaticInitializers = true
	}

	declarations := ec.EndVariableEnvironment()
	if len(declarations) > 0 {
		stmts := make([]*ast.Statement, len(declarations)+1)
		copy(stmts, declarations)
		stmts[len(declarations)] = f.NewReturnStatement(initializer)
		initializer = f.NewImmediatelyInvokedArrowFunction(stmts)
	}

	if tx.classInfoStack != nil {
		if ast.IsStatic(node) {
			initializer = tx.injectPendingInitializers(tx.classInfoStack, true, initializer)
			if result.extraInitializersName != nil {
				var thisArg *ast.Node
				if tx.classInfoStack.classThis != nil {
					thisArg = tx.classInfoStack.classThis
				} else {
					thisArg = f.NewThisExpression()
				}
				tx.classInfoStack.pendingStaticInitializers = append(tx.classInfoStack.pendingStaticInitializers,
					f.NewRunInitializersHelper(thisArg, result.extraInitializersName, nil),
				)
			}
		} else {
			initializer = tx.injectPendingInitializers(tx.classInfoStack, false, initializer)
			if result.extraInitializersName != nil {
				tx.classInfoStack.pendingInstanceInitializers = append(tx.classInfoStack.pendingInstanceInitializers,
					f.NewRunInitializersHelper(f.NewThisExpression(), result.extraInitializersName, nil),
				)
			}
		}
	}

	tx.exitClassElement()

	prop := node.AsPropertyDeclaration()
	return tx.finishClassElement(
		f.UpdatePropertyDeclaration(prop, result.modifiers, result.name, nil, nil, initializer),
		node,
	)
}

func (tx *esDecoratorTransformer) visitClassStaticBlockDeclaration(node *ast.Node) *ast.Node {
	tx.enterClassElement(node)
	f := tx.Factory()

	var result *ast.Node
	if isClassNamedEvaluationHelperBlock(tx.EmitContext(), node) {
		result = tx.Visitor().VisitEachChild(node)
		// Transfer AssignedName metadata to the new node so isClassNamedEvaluationHelperBlock
		// can still find it after visiting (visiting may create a new node when this->_classThis)
		if assignedName := tx.EmitContext().AssignedName(node); assignedName != nil && result != node {
			tx.EmitContext().SetAssignedName(result, assignedName)
		}
	} else if isClassThisAssignmentBlock(tx.EmitContext(), node) {
		savedClassThis := tx.classThis
		tx.classThis = nil
		result = tx.Visitor().VisitEachChild(node)
		tx.classThis = savedClassThis
	} else {
		result = tx.Visitor().VisitEachChild(node)
		if tx.classInfoStack != nil {
			tx.classInfoStack.hasStaticInitializers = true
			if len(tx.classInfoStack.pendingStaticInitializers) > 0 {
				stmts := []*ast.Statement{}
				for _, init := range tx.classInfoStack.pendingStaticInitializers {
					initStmt := f.NewExpressionStatement(init)
					tx.EmitContext().SetSourceMapRange(initStmt, tx.EmitContext().SourceMapRange(init))
					stmts = append(stmts, initStmt)
				}
				body := f.NewBlock(f.NewNodeList(statementsToNodes(stmts)), true)
				staticBlock := f.NewClassStaticBlockDeclaration(nil, body)
				tx.classInfoStack.pendingStaticInitializers = nil
				// Return both the new static block and the original
				tx.exitClassElement()
				return transformers.SingleOrMany([]*ast.Node{staticBlock, result}, tx.Factory())
			}
		}
	}

	tx.exitClassElement()
	return result
}

// --- Expression visitors ---

func (tx *esDecoratorTransformer) visitThisExpression(node *ast.Node) *ast.Node {
	if tx.classThis != nil {
		return tx.classThis
	}
	return node
}

func (tx *esDecoratorTransformer) visitCallExpression(node *ast.Node) *ast.Node {
	call := node.AsCallExpression()
	if ast.IsSuperProperty(call.Expression) && tx.classThis != nil {
		expression := tx.Visitor().VisitNode(call.Expression)
		argVisitor := tx.EmitContext().NewNodeVisitor(tx.visit)
		argumentsList := argVisitor.VisitNodes(call.Arguments)
		invocation := tx.Factory().NewFunctionCallCall(expression, tx.classThis, argumentsList.Nodes)
		tx.EmitContext().SetOriginal(invocation, node)
		invocation.Loc = node.Loc
		return invocation
	}
	return tx.Visitor().VisitEachChild(node)
}

func (tx *esDecoratorTransformer) visitTaggedTemplateExpression(node *ast.Node) *ast.Node {
	tte := node.AsTaggedTemplateExpression()
	if ast.IsSuperProperty(tte.Tag) && tx.classThis != nil {
		tag := tx.Visitor().VisitNode(tte.Tag)
		boundTag := tx.Factory().NewFunctionBindCall(tag, tx.classThis, []*ast.Expression{})
		tx.EmitContext().SetOriginal(boundTag, node)
		boundTag.Loc = node.Loc
		template := tx.Visitor().VisitNode(tte.Template)
		return tx.Factory().UpdateTaggedTemplateExpression(tte, boundTag, nil, nil, template)
	}
	return tx.Visitor().VisitEachChild(node)
}

func (tx *esDecoratorTransformer) visitPropertyAccessExpression(node *ast.Node) *ast.Node {
	pa := node.AsPropertyAccessExpression()
	if ast.IsSuperProperty(node) && ast.IsIdentifier(pa.Name()) && tx.classThis != nil && tx.classSuper != nil {
		propertyName := tx.Factory().NewStringLiteralFromNode(pa.Name())
		superProperty := tx.Factory().NewReflectGetCall(tx.classSuper, propertyName, tx.classThis)
		tx.EmitContext().SetOriginal(superProperty, pa.Expression)
		superProperty.Loc = pa.Expression.Loc
		return superProperty
	}
	return tx.Visitor().VisitEachChild(node)
}

func (tx *esDecoratorTransformer) visitElementAccessExpression(node *ast.Node) *ast.Node {
	ea := node.AsElementAccessExpression()
	if ast.IsSuperProperty(node) && tx.classThis != nil && tx.classSuper != nil {
		propertyName := tx.Visitor().VisitNode(ea.ArgumentExpression)
		superProperty := tx.Factory().NewReflectGetCall(tx.classSuper, propertyName, tx.classThis)
		tx.EmitContext().SetOriginal(superProperty, ea.Expression)
		superProperty.Loc = ea.Expression.Loc
		return superProperty
	}
	return tx.Visitor().VisitEachChild(node)
}

// --- Simple pass-through visitors that handle NamedEvaluation ---

func (tx *esDecoratorTransformer) visitParameterDeclaration(node *ast.Node) *ast.Node {
	if isNamedEvaluation(tx.EmitContext(), node) && isAnonymousClassNeedingAssignedName(node.Initializer()) {
		node = transformNamedEvaluation(tx.EmitContext(), node, false, "")
	}
	return tx.Visitor().VisitEachChild(node)
}

func (tx *esDecoratorTransformer) visitBinaryExpression(node *ast.Node, discarded bool) *ast.Node {
	f := tx.Factory()
	ec := tx.EmitContext()
	bin := node.AsBinaryExpression()

	if ast.IsDestructuringAssignment(node) {
		left := tx.visitAssignmentPattern(bin.Left)
		right := tx.Visitor().VisitNode(bin.Right)
		return f.UpdateBinaryExpression(bin, nil, left, nil, bin.OperatorToken, right)
	}

	if ast.IsAssignmentExpression(node, false) {
		if isNamedEvaluation(ec, node) && isAnonymousClassNeedingAssignedName(bin.Right) {
			node = transformNamedEvaluation(ec, node, false, "")
			return tx.Visitor().VisitEachChild(node)
		}

		if ast.IsSuperProperty(bin.Left) && tx.classThis != nil && tx.classSuper != nil {
			var setterName *ast.Expression
			if ast.IsElementAccessExpression(bin.Left) {
				setterName = tx.Visitor().VisitNode(bin.Left.AsElementAccessExpression().ArgumentExpression)
			} else if ast.IsPropertyAccessExpression(bin.Left) && ast.IsIdentifier(bin.Left.AsPropertyAccessExpression().Name()) {
				setterName = f.NewStringLiteralFromNode(bin.Left.AsPropertyAccessExpression().Name())
			}
			if setterName != nil {
				expression := tx.Visitor().VisitNode(bin.Right)
				if isCompoundAssignment(bin.OperatorToken.Kind) {
					getterName := setterName
					if !transformers.IsSimpleInlineableExpression(setterName) {
						getterName = f.NewTempVariable()
						ec.AddVariableDeclaration(getterName)
						setterName = f.NewAssignmentExpression(getterName, setterName)
					}
					superPropertyGet := f.NewReflectGetCall(tx.classSuper, getterName, tx.classThis)
					ec.SetOriginal(superPropertyGet, bin.Left)
					superPropertyGet.Loc = bin.Left.Loc
					expression = f.AsNodeFactory().NewBinaryExpression(
						nil,
						superPropertyGet,
						nil,
						f.NewToken(getNonAssignmentOperatorForCompoundAssignment(bin.OperatorToken.Kind)),
						expression,
					)
					expression.Loc = node.Loc
				}
				var temp *ast.Expression
				if !discarded {
					temp = f.NewTempVariable()
					ec.AddVariableDeclaration(temp)
				}
				if temp != nil {
					expression = f.NewAssignmentExpression(temp, expression)
					expression.Loc = node.Loc
				}
				expression = f.NewReflectSetCall(tx.classSuper, setterName, expression, tx.classThis)
				ec.SetOriginal(expression, node)
				expression.Loc = node.Loc
				if temp != nil {
					expression = f.NewCommaExpression(expression, temp)
					expression.Loc = node.Loc
				}
				return expression
			}
		}
	}

	if bin.OperatorToken.Kind == ast.KindCommaToken {
		left := tx.discardedVisitor.VisitNode(bin.Left)
		var right *ast.Node
		if discarded {
			right = tx.discardedVisitor.VisitNode(bin.Right)
		} else {
			right = tx.Visitor().VisitNode(bin.Right)
		}
		return f.UpdateBinaryExpression(bin, nil, left, nil, bin.OperatorToken, right)
	}

	return tx.Visitor().VisitEachChild(node)
}

func (tx *esDecoratorTransformer) visitPropertyAssignment(node *ast.Node) *ast.Node {
	if isNamedEvaluation(tx.EmitContext(), node) && isAnonymousClassNeedingAssignedName(node.Initializer()) {
		node = transformNamedEvaluation(tx.EmitContext(), node, false, "")
	}
	return tx.Visitor().VisitEachChild(node)
}

func (tx *esDecoratorTransformer) visitVariableDeclaration(node *ast.Node) *ast.Node {
	if isNamedEvaluation(tx.EmitContext(), node) && isAnonymousClassNeedingAssignedName(node.Initializer()) {
		node = transformNamedEvaluation(tx.EmitContext(), node, false, "")
	}
	return tx.Visitor().VisitEachChild(node)
}

func (tx *esDecoratorTransformer) visitBindingElement(node *ast.Node) *ast.Node {
	if isNamedEvaluation(tx.EmitContext(), node) && isAnonymousClassNeedingAssignedName(node.Initializer()) {
		node = transformNamedEvaluation(tx.EmitContext(), node, false, "")
	}
	return tx.Visitor().VisitEachChild(node)
}

func (tx *esDecoratorTransformer) visitExportAssignment(node *ast.Node) *ast.Node {
	if isNamedEvaluation(tx.EmitContext(), node) && isAnonymousClassNeedingAssignedName(node.Expression()) {
		node = transformNamedEvaluation(tx.EmitContext(), node, false, "")
	}
	return tx.Visitor().VisitEachChild(node)
}

func (tx *esDecoratorTransformer) visitForStatement(node *ast.Node) *ast.Node {
	f := tx.Factory()
	forStmt := node.AsForStatement()
	return f.UpdateForStatement(
		forStmt,
		tx.discardedVisitor.VisitNode(forStmt.Initializer),
		tx.Visitor().VisitNode(forStmt.Condition),
		tx.discardedVisitor.VisitNode(forStmt.Incrementor),
		tx.Visitor().VisitNode(forStmt.Statement),
	)
}

func (tx *esDecoratorTransformer) visitExpressionStatement(node *ast.Node) *ast.Node {
	return tx.discardedVisitor.VisitEachChild(node)
}

func (tx *esDecoratorTransformer) visitPreOrPostfixUnaryExpression(node *ast.Node, discarded bool) *ast.Node {
	f := tx.Factory()
	ec := tx.EmitContext()

	var operator ast.Kind
	var operandNode *ast.Node
	if ast.IsPrefixUnaryExpression(node) {
		operator = node.AsPrefixUnaryExpression().Operator
		operandNode = node.AsPrefixUnaryExpression().Operand
	} else {
		operator = node.AsPostfixUnaryExpression().Operator
		operandNode = node.AsPostfixUnaryExpression().Operand
	}

	if operator == ast.KindPlusPlusToken || operator == ast.KindMinusMinusToken {
		operand := ast.SkipParentheses(operandNode)
		if ast.IsSuperProperty(operand) && tx.classThis != nil && tx.classSuper != nil {
			var setterName *ast.Expression
			if ast.IsElementAccessExpression(operand) {
				setterName = tx.Visitor().VisitNode(operand.AsElementAccessExpression().ArgumentExpression)
			} else if ast.IsPropertyAccessExpression(operand) && ast.IsIdentifier(operand.AsPropertyAccessExpression().Name()) {
				setterName = f.NewStringLiteralFromNode(operand.AsPropertyAccessExpression().Name())
			}
			if setterName != nil {
				getterName := setterName
				if !transformers.IsSimpleInlineableExpression(setterName) {
					getterName = f.NewTempVariable()
					ec.AddVariableDeclaration(getterName)
					setterName = f.NewAssignmentExpression(getterName, setterName)
				}

				expression := f.NewReflectGetCall(tx.classSuper, getterName, tx.classThis)
				ec.SetOriginal(expression, node)
				expression.Loc = node.Loc

				var temp *ast.Expression
				if !discarded {
					temp = f.NewTempVariable()
					ec.AddVariableDeclaration(temp)
				}

				expression = expandPreOrPostfixIncrementOrDecrementExpression(f, node, expression, func(name *ast.IdentifierNode) {
					ec.AddVariableDeclaration(name)
				}, temp)

				expression = f.NewReflectSetCall(tx.classSuper, setterName, expression, tx.classThis)
				ec.SetOriginal(expression, node)
				expression.Loc = node.Loc

				if temp != nil {
					expression = f.NewCommaExpression(expression, temp)
					expression.Loc = node.Loc
				}

				return expression
			}
		}
	}

	return tx.Visitor().VisitEachChild(node)
}

func (tx *esDecoratorTransformer) visitParenthesizedExpression(node *ast.Node, discarded bool) *ast.Node {
	f := tx.Factory()
	pe := node.AsParenthesizedExpression()
	var expression *ast.Node
	if discarded {
		expression = tx.discardedVisitor.VisitNode(pe.Expression)
	} else {
		expression = tx.Visitor().VisitNode(pe.Expression)
	}
	return f.UpdateParenthesizedExpression(pe, expression)
}

func (tx *esDecoratorTransformer) visitComputedPropertyName(node *ast.Node) *ast.Node {
	cpn := node.AsComputedPropertyName()
	expression := tx.Visitor().VisitNode(cpn.Expression)
	if !transformers.IsSimpleInlineableExpression(expression) {
		expression = tx.injectPendingExpressions(expression)
	}
	return tx.Factory().UpdateComputedPropertyName(cpn, expression)
}

// --- Property name visitors ---

func (tx *esDecoratorTransformer) visitPropertyName(node *ast.Node) *ast.Node {
	if ast.IsComputedPropertyName(node) {
		return tx.visitComputedPropertyName(node)
	}
	return tx.Visitor().VisitNode(node)
}

func (tx *esDecoratorTransformer) visitReferencedPropertyName(node *ast.Node) (referencedName *ast.Expression, name *ast.Node) {
	if ast.IsPropertyNameLiteral(node) || ast.IsPrivateIdentifier(node) {
		referencedName = tx.Factory().NewStringLiteralFromNode(node)
		name = tx.Visitor().VisitNode(node)
		return
	}

	cpn := node.AsComputedPropertyName()
	if ast.IsPropertyNameLiteral(cpn.Expression) && !ast.IsIdentifier(cpn.Expression) {
		referencedName = tx.Factory().NewStringLiteralFromNode(cpn.Expression)
		name = tx.Visitor().VisitNode(node)
		return
	}

	referencedName = tx.Factory().NewGeneratedNameForNode(node)
	tx.EmitContext().AddVariableDeclaration(referencedName)

	key := tx.Factory().NewPropKeyHelper(tx.Visitor().VisitNode(cpn.Expression))
	assignment := tx.Factory().NewAssignmentExpression(referencedName, key)
	updatedName := tx.Factory().UpdateComputedPropertyName(cpn, tx.injectPendingExpressions(assignment))
	return referencedName, updatedName
}

// --- Pending expressions injection ---

func (tx *esDecoratorTransformer) injectPendingExpressions(expression *ast.Expression) *ast.Expression {
	f := tx.Factory()
	if len(tx.pendingExpressions) > 0 {
		if ast.IsParenthesizedExpression(expression) {
			pe := expression.AsParenthesizedExpression()
			exprs := make([]*ast.Expression, len(tx.pendingExpressions)+1)
			copy(exprs, tx.pendingExpressions)
			exprs[len(tx.pendingExpressions)] = pe.Expression
			expression = f.UpdateParenthesizedExpression(pe, f.InlineExpressions(exprs))
		} else {
			exprs := make([]*ast.Expression, len(tx.pendingExpressions)+1)
			copy(exprs, tx.pendingExpressions)
			exprs[len(tx.pendingExpressions)] = expression
			expression = f.InlineExpressions(exprs)
		}
		tx.pendingExpressions = nil
	}
	return expression
}

func (tx *esDecoratorTransformer) injectPendingInitializers(ci *classInfo, isStatic bool, expression *ast.Expression) *ast.Expression {
	f := tx.Factory()
	var pending []*ast.Expression
	if isStatic {
		pending = ci.pendingStaticInitializers
	} else {
		pending = ci.pendingInstanceInitializers
	}

	if len(pending) > 0 {
		if expression != nil {
			if ast.IsParenthesizedExpression(expression) {
				pe := expression.AsParenthesizedExpression()
				exprs := make([]*ast.Expression, len(pending)+1)
				copy(exprs, pending)
				exprs[len(pending)] = pe.Expression
				expression = f.UpdateParenthesizedExpression(pe, f.InlineExpressions(exprs))
			} else {
				exprs := make([]*ast.Expression, len(pending)+1)
				copy(exprs, pending)
				exprs[len(pending)] = expression
				expression = f.InlineExpressions(exprs)
			}
		} else {
			expression = f.InlineExpressions(pending)
		}
		if isStatic {
			ci.pendingStaticInitializers = nil
		} else {
			ci.pendingInstanceInitializers = nil
		}
	}
	return expression
}

// --- Context object creation for __esDecorate ---

func (tx *esDecoratorTransformer) createESDecorateClassContext(nameExpr *ast.Expression, metadata *ast.IdentifierNode) *ast.Expression {
	f := tx.Factory()
	props := []*ast.Node{
		f.NewPropertyAssignment(nil, f.NewIdentifier("kind"), nil, nil, f.NewStringLiteral("class", 0)),
		f.NewPropertyAssignment(nil, f.NewIdentifier("name"), nil, nil, nameExpr),
		f.NewPropertyAssignment(nil, f.NewIdentifier("metadata"), nil, nil, metadata),
	}
	return f.NewObjectLiteralExpression(f.NewNodeList(props), false)
}

func (tx *esDecoratorTransformer) createESDecorateElementContext(
	kind string,
	nameComputed bool,
	nameExpr *ast.Expression,
	isStatic bool,
	isPrivate bool,
	hasGet bool,
	hasSet bool,
	metadata *ast.IdentifierNode,
) *ast.Expression {
	f := tx.Factory()

	// Build the name value for the context's "name" property
	var nameValue *ast.Expression
	if nameComputed {
		nameValue = nameExpr
	} else if nameExpr != nil && ast.IsPrivateIdentifier(nameExpr) {
		nameValue = f.NewStringLiteralFromNode(nameExpr)
	} else if nameExpr != nil && ast.IsIdentifier(nameExpr) {
		nameValue = f.NewStringLiteralFromNode(nameExpr)
	} else {
		nameValue = nameExpr
	}

	// Build the access object with has/get/set arrow functions
	accessObj := tx.createESDecorateClassElementAccessObject(nameComputed, nameExpr, hasGet, hasSet)

	var staticExpr *ast.Node
	if isStatic {
		staticExpr = f.NewTrueExpression()
	} else {
		staticExpr = f.NewFalseExpression()
	}

	var privateExpr *ast.Node
	if isPrivate {
		privateExpr = f.NewTrueExpression()
	} else {
		privateExpr = f.NewFalseExpression()
	}

	props := []*ast.Node{
		f.NewPropertyAssignment(nil, f.NewIdentifier("kind"), nil, nil, f.NewStringLiteral(kind, 0)),
		f.NewPropertyAssignment(nil, f.NewIdentifier("name"), nil, nil, nameValue),
		f.NewPropertyAssignment(nil, f.NewIdentifier("static"), nil, nil, staticExpr),
		f.NewPropertyAssignment(nil, f.NewIdentifier("private"), nil, nil, privateExpr),
		f.NewPropertyAssignment(nil, f.NewIdentifier("access"), nil, nil, accessObj),
		f.NewPropertyAssignment(nil, f.NewIdentifier("metadata"), nil, nil, metadata),
	}
	return f.NewObjectLiteralExpression(f.NewNodeList(props), false)
}

// createESDecorateClassElementAccessObject creates the "access" object for a class element decorator context.
// Per the spec (15.7.3 CreateDecoratorAccessObject):
//   - has: obj => name in obj
//   - get: obj => obj.name  (or obj => obj[name] for computed)
//   - set: (obj, value) => { obj.name = value; }  (or obj[name] for computed)
func (tx *esDecoratorTransformer) createESDecorateClassElementAccessObject(
	nameComputed bool,
	nameExpr *ast.Expression,
	hasGet bool,
	hasSet bool,
) *ast.Expression {
	f := tx.Factory()
	accessProps := []*ast.Node{}

	// "has" method: obj => name in obj
	accessProps = append(accessProps, tx.createESDecorateClassElementAccessHasMethod(nameComputed, nameExpr))

	// "get" method: obj => obj.name or obj => obj[name]
	if hasGet {
		accessProps = append(accessProps, tx.createESDecorateClassElementAccessGetMethod(nameComputed, nameExpr))
	}

	// "set" method: (obj, value) => { obj.name = value; } or (obj, value) => { obj[name] = value; }
	if hasSet {
		accessProps = append(accessProps, tx.createESDecorateClassElementAccessSetMethod(nameComputed, nameExpr))
	}

	return f.NewObjectLiteralExpression(f.NewNodeList(accessProps), false)
}

func (tx *esDecoratorTransformer) createESDecorateClassElementAccessHasMethod(
	nameComputed bool,
	nameExpr *ast.Expression,
) *ast.Node {
	f := tx.Factory()

	// The property name for the "in" expression
	var propertyName *ast.Expression
	if nameComputed {
		propertyName = nameExpr
	} else if nameExpr != nil && ast.IsIdentifier(nameExpr) {
		propertyName = f.NewStringLiteralFromNode(nameExpr)
	} else {
		propertyName = nameExpr
	}

	objParam := f.NewParameterDeclaration(nil, nil, f.NewIdentifier("obj"), nil, nil, nil)
	inExpr := f.NewBinaryExpression(nil, propertyName, nil, f.NewToken(ast.KindInKeyword), f.NewIdentifier("obj"))

	arrow := f.NewArrowFunction(
		nil, nil,
		f.NewNodeList([]*ast.Node{objParam}),
		nil, nil,
		f.NewToken(ast.KindEqualsGreaterThanToken),
		inExpr,
	)

	return f.NewPropertyAssignment(nil, f.NewIdentifier("has"), nil, nil, arrow)
}

func (tx *esDecoratorTransformer) createESDecorateClassElementAccessGetMethod(
	nameComputed bool,
	nameExpr *ast.Expression,
) *ast.Node {
	f := tx.Factory()

	var accessor *ast.Expression
	if nameComputed {
		accessor = f.NewElementAccessExpression(f.NewIdentifier("obj"), nil, nameExpr, ast.NodeFlagsNone)
	} else {
		accessor = f.NewPropertyAccessExpression(f.NewIdentifier("obj"), nil, nameExpr, ast.NodeFlagsNone)
	}

	objParam := f.NewParameterDeclaration(nil, nil, f.NewIdentifier("obj"), nil, nil, nil)

	arrow := f.NewArrowFunction(
		nil, nil,
		f.NewNodeList([]*ast.Node{objParam}),
		nil, nil,
		f.NewToken(ast.KindEqualsGreaterThanToken),
		accessor,
	)

	return f.NewPropertyAssignment(nil, f.NewIdentifier("get"), nil, nil, arrow)
}

func (tx *esDecoratorTransformer) createESDecorateClassElementAccessSetMethod(
	nameComputed bool,
	nameExpr *ast.Expression,
) *ast.Node {
	f := tx.Factory()

	var accessor *ast.Expression
	if nameComputed {
		accessor = f.NewElementAccessExpression(f.NewIdentifier("obj"), nil, nameExpr, ast.NodeFlagsNone)
	} else {
		accessor = f.NewPropertyAccessExpression(f.NewIdentifier("obj"), nil, nameExpr, ast.NodeFlagsNone)
	}

	assignment := f.NewAssignmentExpression(accessor, f.NewIdentifier("value"))
	stmt := f.NewExpressionStatement(assignment)
	body := f.NewBlock(f.NewNodeList([]*ast.Node{stmt}), false)

	objParam := f.NewParameterDeclaration(nil, nil, f.NewIdentifier("obj"), nil, nil, nil)
	valueParam := f.NewParameterDeclaration(nil, nil, f.NewIdentifier("value"), nil, nil, nil)

	arrow := f.NewArrowFunction(
		nil, nil,
		f.NewNodeList([]*ast.Node{objParam, valueParam}),
		nil, nil,
		f.NewToken(ast.KindEqualsGreaterThanToken),
		body,
	)

	return f.NewPropertyAssignment(nil, f.NewIdentifier("set"), nil, nil, arrow)
}

// --- Private member descriptor helpers ---

// filterAsyncModifier returns a modifier list containing only the async modifier, if present.
func (tx *esDecoratorTransformer) filterAsyncModifier(modifiers *ast.ModifierList) *ast.ModifierList {
	if modifiers == nil {
		return nil
	}
	var asyncMods []*ast.Node
	for _, mod := range modifiers.Nodes {
		if mod.Kind == ast.KindAsyncKeyword {
			asyncMods = append(asyncMods, mod)
		}
	}
	if len(asyncMods) == 0 {
		return nil
	}
	return tx.Factory().NewModifierList(asyncMods)
}

// filterStaticModifier returns a modifier list containing only the static modifier, if present.
func (tx *esDecoratorTransformer) filterStaticModifier(modifiers *ast.ModifierList) *ast.ModifierList {
	if modifiers == nil {
		return nil
	}
	var staticMods []*ast.Node
	for _, mod := range modifiers.Nodes {
		if mod.Kind == ast.KindStaticKeyword {
			staticMods = append(staticMods, mod)
		}
	}
	if len(staticMods) == 0 {
		return nil
	}
	return tx.Factory().NewModifierList(staticMods)
}

// createDescriptorMethod creates a property assignment that wraps a method body
// into a FunctionExpression with __setFunctionName applied.
// kind is "value", "get", or "set".
func (tx *esDecoratorTransformer) createDescriptorMethod(
	original *ast.Node,
	name *ast.Node, // PrivateIdentifier
	modifiers *ast.ModifierList,
	asteriskToken *ast.TokenNode,
	kind string,
	parameters *ast.NodeList,
	body *ast.Node,
) *ast.Node {
	f := tx.Factory()
	ec := tx.EmitContext()

	if body == nil {
		body = f.NewBlock(f.NewNodeList([]*ast.Node{}), false)
	}

	funcExpr := f.NewFunctionExpression(
		modifiers,
		asteriskToken,
		nil, // name
		nil, // typeParameters
		parameters,
		nil, // type
		nil, // fullSignature
		body,
	)
	ec.SetOriginal(funcExpr, original)
	ec.SetSourceMapRange(funcExpr, moveRangePastDecorators(original))
	ec.SetEmitFlags(funcExpr, printer.EFNoComments)

	var prefix string
	if kind == "get" || kind == "set" {
		prefix = kind
	}
	functionName := f.NewStringLiteralFromNode(name)
	namedFunction := f.NewSetFunctionNameHelper(funcExpr, functionName, prefix)

	method := f.NewPropertyAssignment(nil, f.NewIdentifier(kind), nil, nil, namedFunction)
	ec.SetOriginal(method, original)
	ec.SetSourceMapRange(method, moveRangePastDecorators(original))
	ec.SetEmitFlags(method, printer.EFNoComments)
	return method
}

// createMethodDescriptorObject creates { value: __setFunctionName(function(...) { body }, "#name") }
func (tx *esDecoratorTransformer) createMethodDescriptorObject(member *ast.Node, modifiers *ast.ModifierList) *ast.Expression {
	f := tx.Factory()
	paramVisitor := tx.EmitContext().NewNodeVisitor(tx.visit)
	parameters := paramVisitor.VisitNodes(member.ParameterList())
	body := tx.Visitor().VisitNode(member.Body())
	method := member.AsMethodDeclaration()
	return f.NewObjectLiteralExpression(
		f.NewNodeList([]*ast.Node{
			tx.createDescriptorMethod(member, member.Name(), modifiers, method.AsteriskToken, "value", parameters, body),
		}),
		false,
	)
}

// createGetAccessorDescriptorObject creates { get: __setFunctionName(function() { body }, "#name", "get") }
func (tx *esDecoratorTransformer) createGetAccessorDescriptorObject(member *ast.Node, modifiers *ast.ModifierList) *ast.Expression {
	f := tx.Factory()
	body := tx.Visitor().VisitNode(member.Body())
	return f.NewObjectLiteralExpression(
		f.NewNodeList([]*ast.Node{
			tx.createDescriptorMethod(member, member.Name(), modifiers, nil, "get", f.NewNodeList([]*ast.Node{}), body),
		}),
		false,
	)
}

// createSetAccessorDescriptorObject creates { set: __setFunctionName(function(value) { body }, "#name", "set") }
func (tx *esDecoratorTransformer) createSetAccessorDescriptorObject(member *ast.Node, modifiers *ast.ModifierList) *ast.Expression {
	f := tx.Factory()
	paramVisitor := tx.EmitContext().NewNodeVisitor(tx.visit)
	parameters := paramVisitor.VisitNodes(member.ParameterList())
	body := tx.Visitor().VisitNode(member.Body())
	return f.NewObjectLiteralExpression(
		f.NewNodeList([]*ast.Node{
			tx.createDescriptorMethod(member, member.Name(), modifiers, nil, "set", parameters, body),
		}),
		false,
	)
}

// createMethodDescriptorForwarder creates:
//
//	get #name() { return _descriptor.value; }
func (tx *esDecoratorTransformer) createMethodDescriptorForwarder(modifiers *ast.ModifierList, name *ast.Node, descriptorName *ast.IdentifierNode) *ast.Node {
	f := tx.Factory()
	staticOnly := tx.filterStaticModifier(modifiers)
	return f.NewGetAccessorDeclaration(
		staticOnly,
		name,
		nil, // typeParameters
		f.NewNodeList([]*ast.Node{}),
		nil, // type
		nil, // fullSignature
		f.NewBlock(f.NewNodeList([]*ast.Node{
			f.NewReturnStatement(
				f.NewPropertyAccessExpression(descriptorName, nil, f.NewIdentifier("value"), ast.NodeFlagsNone),
			),
		}), false),
	)
}

// createGetAccessorDescriptorForwarder creates:
//
//	get #name() { return _descriptor.get.call(this); }
func (tx *esDecoratorTransformer) createGetAccessorDescriptorForwarder(modifiers *ast.ModifierList, name *ast.Node, descriptorName *ast.IdentifierNode) *ast.Node {
	f := tx.Factory()
	staticOnly := tx.filterStaticModifier(modifiers)
	return f.NewGetAccessorDeclaration(
		staticOnly,
		name,
		nil, // typeParameters
		f.NewNodeList([]*ast.Node{}),
		nil, // type
		nil, // fullSignature
		f.NewBlock(f.NewNodeList([]*ast.Node{
			f.NewReturnStatement(
				f.NewFunctionCallCall(
					f.NewPropertyAccessExpression(descriptorName, nil, f.NewIdentifier("get"), ast.NodeFlagsNone),
					f.NewThisExpression(),
					nil,
				),
			),
		}), false),
	)
}

// createSetAccessorDescriptorForwarder creates:
//
//	set #name(value) { return _descriptor.set.call(this, value); }
func (tx *esDecoratorTransformer) createSetAccessorDescriptorForwarder(modifiers *ast.ModifierList, name *ast.Node, descriptorName *ast.IdentifierNode) *ast.Node {
	f := tx.Factory()
	staticOnly := tx.filterStaticModifier(modifiers)
	return f.NewSetAccessorDeclaration(
		staticOnly,
		name,
		nil, // typeParameters
		f.NewNodeList([]*ast.Node{
			f.NewParameterDeclaration(nil, nil, f.NewIdentifier("value"), nil, nil, nil),
		}),
		nil, // type
		nil, // fullSignature
		f.NewBlock(f.NewNodeList([]*ast.Node{
			f.NewReturnStatement(
				f.NewFunctionCallCall(
					f.NewPropertyAccessExpression(descriptorName, nil, f.NewIdentifier("set"), ast.NodeFlagsNone),
					f.NewThisExpression(),
					[]*ast.Node{f.NewIdentifier("value")},
				),
			),
		}), false),
	)
}

// --- Metadata helpers ---

func (tx *esDecoratorTransformer) createMetadata(name *ast.IdentifierNode, classSuper *ast.IdentifierNode) *ast.Statement {
	f := tx.Factory()

	var superMetadata *ast.Expression
	if classSuper != nil {
		superMetadata = tx.createSymbolMetadataReference(classSuper)
	} else {
		superMetadata = f.NewToken(ast.KindNullKeyword)
	}

	objectCreate := f.NewCallExpression(
		f.NewPropertyAccessExpression(f.NewIdentifier("Object"), nil, f.NewIdentifier("create"), ast.NodeFlagsNone),
		nil, nil,
		f.NewNodeList([]*ast.Expression{superMetadata}),
		ast.NodeFlagsNone,
	)

	symbolCheck := f.NewLogicalANDExpression(
		f.NewTypeCheck(f.NewIdentifier("Symbol"), "function"),
		f.NewPropertyAccessExpression(f.NewIdentifier("Symbol"), nil, f.NewIdentifier("metadata"), ast.NodeFlagsNone),
	)

	conditional := f.NewConditionalExpression(
		symbolCheck,
		f.NewToken(ast.KindQuestionToken),
		objectCreate,
		f.NewToken(ast.KindColonToken),
		f.NewVoidZeroExpression(),
	)

	varDecl := f.NewVariableDeclaration(name, nil, nil, conditional)
	varDeclList := f.NewVariableDeclarationList(ast.NodeFlagsConst, f.NewNodeList([]*ast.Node{varDecl}))
	return f.NewVariableStatement(nil, varDeclList)
}

func (tx *esDecoratorTransformer) createSymbolMetadata(target *ast.Expression, value *ast.IdentifierNode) *ast.Statement {
	f := tx.Factory()

	// Object.defineProperty(target, Symbol.metadata, { configurable: true, writable: true, enumerable: true, value })
	symbolMetadata := f.NewPropertyAccessExpression(f.NewIdentifier("Symbol"), nil, f.NewIdentifier("metadata"), ast.NodeFlagsNone)

	descriptorProps := []*ast.Node{
		f.NewPropertyAssignment(nil, f.NewIdentifier("enumerable"), nil, nil, f.NewTrueExpression()),
		f.NewPropertyAssignment(nil, f.NewIdentifier("configurable"), nil, nil, f.NewTrueExpression()),
		f.NewPropertyAssignment(nil, f.NewIdentifier("writable"), nil, nil, f.NewTrueExpression()),
		f.NewPropertyAssignment(nil, f.NewIdentifier("value"), nil, nil, value),
	}
	descriptor := f.NewObjectLiteralExpression(f.NewNodeList(descriptorProps), false)

	defineProperty := f.NewCallExpression(
		f.NewPropertyAccessExpression(f.NewIdentifier("Object"), nil, f.NewIdentifier("defineProperty"), ast.NodeFlagsNone),
		nil, nil,
		f.NewNodeList([]*ast.Expression{target, symbolMetadata, descriptor}),
		ast.NodeFlagsNone,
	)

	ifStatement := f.NewIfStatement(value, f.NewExpressionStatement(defineProperty), nil)
	tx.EmitContext().SetEmitFlags(ifStatement, printer.EFSingleLine)
	return ifStatement
}

func (tx *esDecoratorTransformer) createSymbolMetadataReference(classSuper *ast.IdentifierNode) *ast.Expression {
	f := tx.Factory()
	symbolMetadata := f.NewPropertyAccessExpression(f.NewIdentifier("Symbol"), nil, f.NewIdentifier("metadata"), ast.NodeFlagsNone)
	elementAccess := f.NewElementAccessExpression(classSuper, nil, symbolMetadata, ast.NodeFlagsNone)
	return f.NewBinaryExpression(nil, elementAccess, nil, f.NewToken(ast.KindQuestionQuestionToken), f.NewToken(ast.KindNullKeyword))
}

// --- Export helpers ---

func (tx *esDecoratorTransformer) createExportDefault(expression *ast.Expression) *ast.Statement {
	f := tx.Factory()
	return f.NewExportAssignment(nil, false, nil, expression)
}

func (tx *esDecoratorTransformer) createExternalModuleExport(name *ast.IdentifierNode) *ast.Statement {
	f := tx.Factory()
	specifier := f.NewExportSpecifier(false, nil, name)
	namedExports := f.NewNamedExports(f.NewNodeList([]*ast.Node{specifier}))
	return f.NewExportDeclaration(nil, false, namedExports, nil, nil)
}

// --- Utility: isAnonymousClassNeedingAssignedName ---

func isAnonymousClassNeedingAssignedName(node *ast.Node) bool {
	if node == nil {
		return false
	}
	return ast.IsClassExpression(node) && node.Name() == nil && isDecoratedClassLike(node)
}

// --- injectClassThisAssignmentIfMissing ---

func classHasClassThisAssignment(ec *printer.EmitContext, node *ast.Node) bool {
	classThisExpr := ec.ClassThis(node)
	if classThisExpr == nil {
		return false
	}
	return core.Some(node.Members(), func(m *ast.Node) bool {
		return isClassThisAssignmentBlock(ec, m)
	})
}

func injectClassThisAssignmentIfMissing(ec *printer.EmitContext, f *printer.NodeFactory, node *ast.Node, classThis *ast.IdentifierNode) *ast.Node {
	if classHasClassThisAssignment(ec, node) {
		return node
	}

	// Create: static { _classThis = this; }
	expression := f.NewAssignmentExpression(classThis, f.NewThisExpression())
	statement := f.NewExpressionStatement(expression)
	body := f.NewBlock(f.NewNodeList([]*ast.Node{statement}), false)
	staticBlock := f.NewClassStaticBlockDeclaration(nil, body)
	ec.SetClassThis(staticBlock, classThis)

	if node.Name() != nil {
		ec.SetSourceMapRange(statement, node.Name().Loc)
	}

	newMembers := make([]*ast.Node, 0, 1+len(node.Members()))
	newMembers = append(newMembers, staticBlock)
	newMembers = append(newMembers, node.Members()...)
	membersList := f.NewNodeList(newMembers)
	membersList.Loc = node.MemberList().Loc

	var updatedNode *ast.Node
	if ast.IsClassDeclaration(node) {
		cd := node.AsClassDeclaration()
		updatedNode = f.UpdateClassDeclaration(cd, cd.Modifiers(), cd.Name(), nil, cd.HeritageClauses, membersList)
	} else {
		ce := node.AsClassExpression()
		updatedNode = f.UpdateClassExpression(ce, ce.Modifiers(), ce.Name(), nil, ce.HeritageClauses, membersList)
	}
	ec.SetClassThis(updatedNode, classThis)
	return updatedNode
}

// --- findSuperStatementIndexPath (duplicated from tstransforms) ---

func findSuperStatementIndexPath(statements []*ast.Statement, start int) []int {
	for i := start; i < len(statements); i++ {
		stmt := statements[i]
		if ast.IsExpressionStatement(stmt) {
			expr := stmt.Expression()
			if ast.IsCallExpression(expr) && expr.AsCallExpression().Expression.Kind == ast.KindSuperKeyword {
				return []int{i}
			}
		}
		if ast.IsTryStatement(stmt) {
			tryBlock := stmt.AsTryStatement().TryBlock.AsBlock()
			path := findSuperStatementIndexPath(tryBlock.Statements.Nodes, 0)
			if len(path) > 0 {
				return append([]int{i}, path...)
			}
		}
	}
	return nil
}

// --- Slice conversion helpers ---

func expressionsToNodes(exprs []*ast.Expression) []*ast.Node {
	nodes := make([]*ast.Node, len(exprs))
	for i, e := range exprs {
		nodes[i] = e
	}
	return nodes
}

func statementsToNodes(stmts []*ast.Statement) []*ast.Node {
	nodes := make([]*ast.Node, len(stmts))
	for i, s := range stmts {
		nodes[i] = s
	}
	return nodes
}

func nodesToStatements(nodes []*ast.Node) []*ast.Statement {
	stmts := make([]*ast.Statement, len(nodes))
	for i, n := range nodes {
		stmts[i] = n
	}
	return stmts
}

// --- Super property destructuring assignment visitors ---

func (tx *esDecoratorTransformer) visitDestructuringAssignmentTarget(node *ast.Node) *ast.Node {
	if ast.IsObjectLiteralExpression(node) || ast.IsArrayLiteralExpression(node) {
		return tx.visitAssignmentPattern(node)
	}

	if ast.IsSuperProperty(node) && tx.classThis != nil && tx.classSuper != nil {
		f := tx.Factory()
		ec := tx.EmitContext()
		var propertyName *ast.Expression
		if ast.IsElementAccessExpression(node) {
			propertyName = tx.Visitor().VisitNode(node.AsElementAccessExpression().ArgumentExpression)
		} else if ast.IsPropertyAccessExpression(node) && ast.IsIdentifier(node.AsPropertyAccessExpression().Name()) {
			propertyName = f.NewStringLiteralFromNode(node.AsPropertyAccessExpression().Name())
		}
		if propertyName != nil {
			expression := createAssignmentTargetWrapper(f, tx.classSuper, propertyName, tx.classThis)
			ec.SetOriginal(expression, node)
			expression.Loc = node.Loc
			return expression
		}
	}

	return tx.Visitor().VisitEachChild(node)
}

func (tx *esDecoratorTransformer) visitAssignmentPattern(node *ast.Node) *ast.Node {
	f := tx.Factory()
	ec := tx.EmitContext()
	if ast.IsArrayLiteralExpression(node) {
		ale := node.AsArrayLiteralExpression()
		arrayVisitor := ec.NewNodeVisitor(tx.visitArrayAssignmentElement)
		elements := arrayVisitor.VisitNodes(ale.Elements)
		return f.UpdateArrayLiteralExpression(ale, elements)
	}
	ole := node.AsObjectLiteralExpression()
	objVisitor := ec.NewNodeVisitor(tx.visitObjectAssignmentElement)
	properties := objVisitor.VisitNodes(ole.Properties)
	return f.UpdateObjectLiteralExpression(ole, properties)
}

func (tx *esDecoratorTransformer) visitArrayAssignmentElement(node *ast.Node) *ast.Node {
	if ast.IsSpreadElement(node) {
		return tx.visitAssignmentRestElement(node)
	}
	if !ast.IsOmittedExpression(node) {
		return tx.visitAssignmentElement(node)
	}
	return tx.Visitor().VisitEachChild(node)
}

func (tx *esDecoratorTransformer) visitAssignmentElement(node *ast.Node) *ast.Node {
	if ast.IsAssignmentExpression(node, true /*excludeCompoundAssignment*/) {
		f := tx.Factory()
		bin := node.AsBinaryExpression()
		if isNamedEvaluation(tx.EmitContext(), node) && isAnonymousClassNeedingAssignedName(bin.Right) {
			node = transformNamedEvaluation(tx.EmitContext(), node, false, "")
			bin = node.AsBinaryExpression()
		}
		assignmentTarget := tx.visitDestructuringAssignmentTarget(bin.Left)
		initializer := tx.Visitor().VisitNode(bin.Right)
		return f.UpdateBinaryExpression(bin, nil, assignmentTarget, nil, bin.OperatorToken, initializer)
	}
	return tx.visitDestructuringAssignmentTarget(node)
}

func (tx *esDecoratorTransformer) visitAssignmentRestElement(node *ast.Node) *ast.Node {
	se := node.AsSpreadElement()
	if ast.IsLeftHandSideExpression(se.Expression) {
		f := tx.Factory()
		expression := tx.visitDestructuringAssignmentTarget(se.Expression)
		return f.UpdateSpreadElement(se, expression)
	}
	return tx.Visitor().VisitEachChild(node)
}

func (tx *esDecoratorTransformer) visitObjectAssignmentElement(node *ast.Node) *ast.Node {
	if ast.IsSpreadAssignment(node) {
		return tx.visitAssignmentRestProperty(node)
	}
	if ast.IsShorthandPropertyAssignment(node) {
		return tx.visitShorthandAssignmentProperty(node)
	}
	if ast.IsPropertyAssignment(node) {
		return tx.visitAssignmentPropertyNode(node)
	}
	return tx.Visitor().VisitEachChild(node)
}

func (tx *esDecoratorTransformer) visitAssignmentPropertyNode(node *ast.Node) *ast.Node {
	f := tx.Factory()
	pa := node.AsPropertyAssignment()
	name := tx.Visitor().VisitNode(pa.Name())
	if ast.IsAssignmentExpression(pa.Initializer, true /*excludeCompoundAssignment*/) {
		assignmentElement := tx.visitAssignmentElement(pa.Initializer)
		return f.UpdatePropertyAssignment(pa, nil, name, nil, nil, assignmentElement)
	}
	if ast.IsLeftHandSideExpression(pa.Initializer) {
		assignmentElement := tx.visitDestructuringAssignmentTarget(pa.Initializer)
		return f.UpdatePropertyAssignment(pa, nil, name, nil, nil, assignmentElement)
	}
	return tx.Visitor().VisitEachChild(node)
}

func (tx *esDecoratorTransformer) visitShorthandAssignmentProperty(node *ast.Node) *ast.Node {
	if isNamedEvaluation(tx.EmitContext(), node) && isAnonymousClassNeedingAssignedName(node.AsShorthandPropertyAssignment().ObjectAssignmentInitializer) {
		node = transformNamedEvaluation(tx.EmitContext(), node, false, "")
	}
	return tx.Visitor().VisitEachChild(node)
}

func (tx *esDecoratorTransformer) visitAssignmentRestProperty(node *ast.Node) *ast.Node {
	sa := node.AsSpreadAssignment()
	if ast.IsLeftHandSideExpression(sa.Expression) {
		f := tx.Factory()
		expression := tx.visitDestructuringAssignmentTarget(sa.Expression)
		return f.UpdateSpreadAssignment(sa, expression)
	}
	return tx.Visitor().VisitEachChild(node)
}

// --- Utility: getNonAssignmentOperatorForCompoundAssignment ---

func getNonAssignmentOperatorForCompoundAssignment(kind ast.Kind) ast.Kind {
	switch kind {
	case ast.KindPlusEqualsToken:
		return ast.KindPlusToken
	case ast.KindMinusEqualsToken:
		return ast.KindMinusToken
	case ast.KindAsteriskEqualsToken:
		return ast.KindAsteriskToken
	case ast.KindAsteriskAsteriskEqualsToken:
		return ast.KindAsteriskAsteriskToken
	case ast.KindSlashEqualsToken:
		return ast.KindSlashToken
	case ast.KindPercentEqualsToken:
		return ast.KindPercentToken
	case ast.KindLessThanLessThanEqualsToken:
		return ast.KindLessThanLessThanToken
	case ast.KindGreaterThanGreaterThanEqualsToken:
		return ast.KindGreaterThanGreaterThanToken
	case ast.KindGreaterThanGreaterThanGreaterThanEqualsToken:
		return ast.KindGreaterThanGreaterThanGreaterThanToken
	case ast.KindAmpersandEqualsToken:
		return ast.KindAmpersandToken
	case ast.KindBarEqualsToken:
		return ast.KindBarToken
	case ast.KindCaretEqualsToken:
		return ast.KindCaretToken
	case ast.KindBarBarEqualsToken:
		return ast.KindBarBarToken
	case ast.KindAmpersandAmpersandEqualsToken:
		return ast.KindAmpersandAmpersandToken
	case ast.KindQuestionQuestionEqualsToken:
		return ast.KindQuestionQuestionToken
	default:
		return ast.KindUnknown
	}
}

// --- Utility: isCompoundAssignment ---

func isCompoundAssignment(kind ast.Kind) bool {
	return kind >= ast.KindFirstCompoundAssignment && kind <= ast.KindLastCompoundAssignment
}

// --- Utility: expandPreOrPostfixIncrementOrDecrementExpression ---

func expandPreOrPostfixIncrementOrDecrementExpression(
	f *printer.NodeFactory,
	node *ast.Node,
	expression *ast.Expression,
	recordTempVariable func(*ast.IdentifierNode),
	resultVariable *ast.Expression,
) *ast.Expression {
	var operator ast.Kind
	if ast.IsPrefixUnaryExpression(node) {
		operator = node.AsPrefixUnaryExpression().Operator
	} else {
		operator = node.AsPostfixUnaryExpression().Operator
	}

	temp := f.NewTempVariable()
	recordTempVariable(temp)
	expression = f.NewAssignmentExpression(temp, expression)
	expression.Loc = node.Loc

	var operation *ast.Expression
	if ast.IsPrefixUnaryExpression(node) {
		operation = f.AsNodeFactory().NewPrefixUnaryExpression(operator, temp)
	} else {
		operation = f.AsNodeFactory().NewPostfixUnaryExpression(temp, operator)
	}
	operation.Loc = node.Loc

	if resultVariable != nil {
		operation = f.NewAssignmentExpression(resultVariable, operation)
		operation.Loc = node.Loc
	}

	expression = f.NewCommaExpression(expression, operation)
	expression.Loc = node.Loc

	if ast.IsPostfixUnaryExpression(node) {
		expression = f.NewCommaExpression(expression, temp)
		expression.Loc = node.Loc
	}

	return expression
}

// --- Utility: createAssignmentTargetWrapper ---
// Creates: ({ set value(_p) { Reflect.set(target, key, _p, receiver) } }).value

func createAssignmentTargetWrapper(f *printer.NodeFactory, target *ast.Expression, propertyKey *ast.Expression, receiver *ast.Expression) *ast.Expression {
	paramName := f.NewTempVariable()
	reflectSetCall := f.NewReflectSetCall(target, propertyKey, paramName, receiver)
	statement := f.NewExpressionStatement(reflectSetCall)
	body := f.AsNodeFactory().NewBlock(f.AsNodeFactory().NewNodeList([]*ast.Node{statement}), false)
	param := f.AsNodeFactory().NewParameterDeclaration(nil, nil, paramName, nil, nil, nil)
	setter := f.AsNodeFactory().NewSetAccessorDeclaration(nil, f.AsNodeFactory().NewIdentifier("value"), nil, f.AsNodeFactory().NewNodeList([]*ast.Node{param}), nil, nil, body)
	obj := f.AsNodeFactory().NewObjectLiteralExpression(f.AsNodeFactory().NewNodeList([]*ast.Node{setter}), false)
	paren := f.AsNodeFactory().NewParenthesizedExpression(obj)
	return f.AsNodeFactory().NewPropertyAccessExpression(paren, nil, f.AsNodeFactory().NewIdentifier("value"), ast.NodeFlagsNone)
}
