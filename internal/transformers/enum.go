package transformers

// !!! Unscoped enum member references across merged enum declarations are not yet supported (e.g `enum E {A}; enum E {B=A}`)

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/binder"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/jsnum"
	"github.com/microsoft/typescript-go/internal/printer"
	"github.com/microsoft/typescript-go/internal/scanner"
)

type EnumTransformer struct {
	Transformer
	compilerOptions                     *core.CompilerOptions
	parent                              *ast.Node
	current                             *ast.Node
	currentSourceFile                   *ast.Node
	currentLexicalScope                 *ast.Node // SourceFile | Block | ModuleBlock | CaseBlock
	currentScopeFirstDeclarationsOfName map[string]*ast.Node
	nameResolver                        *binder.NameResolver
	referencedEnumMembers               core.Set[*ast.Node]
	isVisitingEnumBody                  bool
}

func NewEnumTransformer(emitContext *printer.EmitContext, compilerOptions *core.CompilerOptions) *Transformer {
	tx := &EnumTransformer{compilerOptions: compilerOptions}
	return tx.newTransformer(tx.visit, emitContext)
}

// Visits each node in the AST
func (tx *EnumTransformer) visit(node *ast.Node) *ast.Node {
	savedLexicalScope := tx.currentLexicalScope
	savedCurrentScopeFirstDeclarationsOfName := tx.currentScopeFirstDeclarationsOfName
	savedParent := tx.parent
	tx.parent = tx.current
	tx.current = node

	switch node.Kind {
	case ast.KindSourceFile:
		tx.currentLexicalScope = node
		tx.currentSourceFile = node
		tx.currentScopeFirstDeclarationsOfName = nil
	case ast.KindCaseBlock, ast.KindModuleBlock, ast.KindBlock:
		tx.currentLexicalScope = node
		tx.currentScopeFirstDeclarationsOfName = nil
	case ast.KindFunctionDeclaration, ast.KindClassDeclaration, ast.KindModuleDeclaration:
		tx.recordEmittedDeclarationInScope(node)
	}

	switch node.Kind {
	case ast.KindEnumDeclaration:
		node = tx.visitEnumDeclaration(node.AsEnumDeclaration())
	case ast.KindIdentifier:
		node = tx.visitIdentifier(node)
	default:
		node = tx.VisitEachChild(node)
	}

	if tx.currentLexicalScope != savedLexicalScope {
		// only reset the first declaration for a name if we are exiting the scope in which it was declared
		tx.currentScopeFirstDeclarationsOfName = savedCurrentScopeFirstDeclarationsOfName
	}
	tx.currentLexicalScope = savedLexicalScope
	tx.current = tx.parent
	tx.parent = savedParent
	return node
}

// Records that a declaration was emitted in the current scope, if it was the first declaration for the provided symbol.
func (tx *EnumTransformer) recordEmittedDeclarationInScope(node *ast.Node) {
	name := node.Name()
	if name != nil && ast.IsIdentifier(name) {
		if tx.currentScopeFirstDeclarationsOfName == nil {
			tx.currentScopeFirstDeclarationsOfName = make(map[string]*ast.Node)
		}
		text := name.Text()
		if _, found := tx.currentScopeFirstDeclarationsOfName[text]; !found {
			tx.currentScopeFirstDeclarationsOfName[text] = node
		}
	}
}

// Determines whether a declaration is the first declaration with the same name emitted in the current scope.
func (tx *EnumTransformer) isFirstEmittedDeclarationInScope(node *ast.Node) bool {
	name := node.Name()
	if name != nil && ast.IsIdentifier(name) {
		text := name.Text()
		if firstDeclaration, found := tx.currentScopeFirstDeclarationsOfName[text]; found {
			return firstDeclaration == node
		}
	}
	return false
}

// Gets an expression that represents a property name.
func (tx *EnumTransformer) getExpressionForPropertyName(member *ast.EnumMember) *ast.Expression {
	name := member.Name()
	switch name.Kind {
	case ast.KindPrivateIdentifier:
		return tx.Factory.NewIdentifier("")
	case ast.KindComputedPropertyName:
		n := name.AsComputedPropertyName()
		// enums don't support computed properties so we always generate the 'expression' part of the name as-is.
		return tx.VisitNode(n.Expression)
	case ast.KindIdentifier:
		return tx.EmitContext.NewStringLiteralFromNode(name)
	case ast.KindStringLiteral:
		return tx.Factory.NewStringLiteral(name.AsStringLiteral().Text)
	case ast.KindNumericLiteral:
		return tx.Factory.NewNumericLiteral(name.AsNumericLiteral().Text)
	default:
		return name
	}
}

// Gets a variable name that represents a property name.
func (tx *EnumTransformer) getVariableNameForPropertyName(member *ast.EnumMember) *ast.Expression {
	name := member.Name()

	switch name.Kind {
	case ast.KindIdentifier:
		temp := copyIdentifier(tx.EmitContext, name)
		tx.EmitContext.SetEmitFlags(temp, printer.EFNoComments)
		return temp

	case ast.KindStringLiteral:
		if scanner.IsIdentifierText(name.AsStringLiteral().Text, tx.compilerOptions.GetEmitScriptTarget()) {
			return tx.Factory.NewIdentifier(name.AsStringLiteral().Text)
		}
		if scanner.IsIdentifierText(name.AsStringLiteral().Text, core.ScriptTargetLatest) {
			return tx.EmitContext.NewTempVariable(printer.AutoGenerateOptions{})
		}
	}

	return nil
}

// Gets the declaration name used inside of a namespace or enum.
func (tx *EnumTransformer) getNamespaceParameterName(node *ast.EnumDeclarationNode) *ast.IdentifierNode {
	name := tx.EmitContext.NewGeneratedNameForNode(node, printer.AutoGenerateOptions{})
	tx.EmitContext.SetSourceMapRange(name, node.Name().Loc)
	return name
}

// Gets the expression used to refer to a namespace or enum within the body of its declaration.
func (tx *EnumTransformer) getNamespaceContainerName(node *ast.EnumDeclarationNode) *ast.IdentifierNode {
	return tx.EmitContext.NewGeneratedNameForNode(node, printer.AutoGenerateOptions{})
}

func (tx *EnumTransformer) addVarForEnumDeclaration(statements []*ast.Statement, node *ast.EnumDeclaration) ([]*ast.Statement, bool) {
	tx.recordEmittedDeclarationInScope(node.AsNode())
	if !tx.isFirstEmittedDeclarationInScope(node.AsNode()) {
		return statements, false
	}

	// var name;
	name := getLocalName(tx.EmitContext, node.AsNode(), false /*allowComments*/, true /*allowSourceMaps*/, false /*ignoreAssignedName*/)
	varDecl := tx.Factory.NewVariableDeclaration(name, nil, nil, nil)
	varFlags := core.IfElse(tx.currentLexicalScope == tx.currentSourceFile, ast.NodeFlagsNone, ast.NodeFlagsLet)
	varDecls := tx.Factory.NewVariableDeclarationList(varFlags, tx.Factory.NewNodeList([]*ast.Node{varDecl}))
	varModifiers := extractModifiers(tx.EmitContext, node.Modifiers(), ast.ModifierFlagsExport)
	varStatement := tx.Factory.NewVariableStatement(varModifiers, varDecls)

	tx.EmitContext.SetOriginal(varDecl, node.AsNode())
	// !!! synthetic comments
	tx.EmitContext.SetOriginal(varStatement, node.AsNode())

	// Adjust the source map emit to match the old emitter.
	tx.EmitContext.SetSourceMapRange(varDecls, node.Loc)

	// Trailing comments for enum declaration should be emitted after the function closure
	// instead of the variable statement:
	//
	//     /** Leading comment*/
	//     enum E {
	//         A
	//     } // trailing comment
	//
	// Should emit:
	//
	//     /** Leading comment*/
	//     var E;
	//     (function (E) {
	//         E[E["A"] = 0] = "A";
	//     })(E || (E = {})); // trailing comment
	//
	tx.EmitContext.SetCommentRange(varStatement, node.Loc)
	tx.EmitContext.AddEmitFlags(varStatement, printer.EFNoTrailingComments)
	statements = append(statements, varStatement)

	return statements, true
}

func (tx *EnumTransformer) visitEnumDeclaration(node *ast.EnumDeclaration) *ast.Node {
	statements := []*ast.Statement{}

	// If needed, we should emit a variable declaration for the enum:
	//  var name;
	statements, varAdded := tx.addVarForEnumDeclaration(statements, node)

	// If we emit a leading variable declaration, we should not emit leading comments for the enum body, but we should
	// still emit the comments if we are emitting to a System module.
	emitFlags := printer.EFNone
	if varAdded && (tx.compilerOptions.GetEmitModuleKind() != core.ModuleKindSystem || tx.currentLexicalScope != tx.currentSourceFile) {
		emitFlags |= printer.EFNoLeadingComments
	}

	// `parameterName` is the declaration name used inside of the enum.
	parameterName := tx.getNamespaceParameterName(node.AsNode())

	// `containerName` is the expression used inside of the enum for assignments.
	containerName := tx.getNamespaceContainerName(node.AsNode())

	// `exportName` is the expression used within this node's container for any exported references.
	exportName := getDeclarationName(tx.EmitContext, node.AsNode(), false /*allowComments*/, true /*allowSourceMaps*/)

	//  x || (x = {})
	//  exports.x || (exports.x = {})
	enumArg := tx.Factory.NewBinaryExpression(
		exportName,
		tx.Factory.NewToken(ast.KindBarBarToken),
		tx.Factory.NewBinaryExpression(
			exportName,
			tx.Factory.NewToken(ast.KindEqualsToken),
			tx.Factory.NewObjectLiteralExpression(tx.Factory.NewNodeList([]*ast.Node{}), false),
		),
	)

	// !!! handle export of namespace

	// (function (name) { ... })(name || (name = {}))
	enumBody := tx.transformEnumBody(node, containerName)
	enumParam := tx.Factory.NewParameterDeclaration(nil, nil, parameterName, nil, nil, nil)
	enumFunc := tx.Factory.NewFunctionExpression(nil, nil, nil, nil, tx.Factory.NewNodeList([]*ast.Node{enumParam}), nil, enumBody)
	enumCall := tx.Factory.NewCallExpression(enumFunc, nil, nil, tx.Factory.NewNodeList([]*ast.Node{enumArg}), ast.NodeFlagsNone)
	enumStatement := tx.Factory.NewExpressionStatement(enumCall)
	tx.EmitContext.SetOriginal(enumStatement, node.AsNode())
	tx.EmitContext.SetCommentAndSourceMapRanges(enumStatement, node.Loc)
	tx.EmitContext.AddEmitFlags(enumStatement, emitFlags)
	statements = append(statements, enumStatement)
	return tx.Factory.NewSyntaxList(statements)
}

// Transforms the body of an enum declaration.
func (tx *EnumTransformer) transformEnumBody(node *ast.EnumDeclaration, enumName *ast.IdentifierNode) *ast.BlockNode {
	// visit the children of `node` in advance to capture any references to enum members
	savedIsVisitingEnumBody := tx.isVisitingEnumBody
	tx.isVisitingEnumBody = true
	node = node.VisitEachChild(&tx.NodeVisitor).AsEnumDeclaration()
	tx.isVisitingEnumBody = savedIsVisitingEnumBody

	statements := []*ast.Statement{}
	if len(node.Members.Nodes) > 0 {
		tx.EmitContext.StartVarEnvironment()

		var autoValue jsnum.Number
		var autoVar *ast.IdentifierNode
		var useAutoVar bool
		for i := range len(node.Members.Nodes) {
			//  E[E["A"] = 0] = "A";
			statements = tx.transformEnumMember(
				statements,
				node.Members.Nodes,
				i,
				enumName,
				&autoValue,
				&autoVar,
				&useAutoVar,
			)
			autoValue++
		}

		if autoVar != nil {
			tx.EmitContext.HoistVariable(autoVar)
		}

		statements = core.Concatenate(tx.EmitContext.EndVarEnvironment(), statements)
	}

	statementList := tx.Factory.NewNodeList(statements)
	statementList.Loc = node.Members.Loc
	return tx.Factory.NewBlock(statementList, true /*multiline*/)
}

// Transforms an enum member into a statement. It is expected that `members` have already been visited.
func (tx *EnumTransformer) transformEnumMember(
	statements []*ast.Statement,
	members []*ast.EnumMemberNode,
	index int,
	enumName *ast.IdentifierNode,
	autoValue *jsnum.Number,
	autoVar **ast.IdentifierNode,
	useAutoVar *bool,
) []*ast.Statement {
	memberNode := members[index]
	member := memberNode.AsEnumMember()

	savedParent := tx.parent
	tx.parent = tx.current
	tx.current = memberNode

	var useConditionalReverseMapping bool
	var useExplicitReverseMapping bool

	//  E[E["A"] = 0] = "A";
	//      ^^^         ^^^
	name := tx.getExpressionForPropertyName(member)

	var expression *ast.Expression
	if member.Initializer == nil {
		// Enum members without an initializer are auto-numbered. We will use constant values if there was no preceding
		// initialized member, or if the preceding initialized member was a numeric literal.
		if *useAutoVar {
			// If you are using an auto-numbered member following a non-numeric literal, we assume the previous member
			// produced a valid numeric value. This assumption is intended to be validated by the type checker prior to
			// emit.
			//  E[E["A"] = ++auto] = "A";
			//             ^^^^^^
			expression = tx.Factory.NewPrefixUnaryExpression(ast.KindPlusPlusToken, *autoVar)
			useExplicitReverseMapping = true
		} else {
			// If the preceding auto value is a finite number, we can emit a numeric literal for the member initializer:
			//  E[E["A"] = 0] = "A";
			//             ^
			// If not, we cannot emit a valid numeric literal for the member initializer and emit `void 0` instead:
			//  E["A"] = void 0;
			//           ^^^^^^
			expression = constantExpression(*autoValue, tx.Factory)
			if expression != nil {
				useExplicitReverseMapping = true
			} else {
				expression = tx.Factory.NewVoidExpression(tx.Factory.NewNumericLiteral("0"))
			}
		}
	} else {
		// Enum members with an initializer may restore auto-numbering if the initializer is a numeric literal. If we
		// cannot syntactically determine the initializer value and the following enum member is auto-numbered, we will
		// use an `auto` variable to perform the remaining auto-numbering at runtime.
		var hasNumericInitializer, hasStringInitializer bool
		switch value := constantValue(member.Initializer).(type) {
		case jsnum.Number:
			hasNumericInitializer = true
			*autoValue = value
		case string:
			hasStringInitializer = true
			*autoValue = jsnum.NaN()
		default:
			*autoValue = jsnum.NaN()
		}

		nextIsAuto := index+1 < len(members) && members[index+1].AsEnumMember().Initializer == nil
		useExplicitReverseMapping = hasNumericInitializer || !hasStringInitializer && nextIsAuto
		useConditionalReverseMapping = !hasNumericInitializer && !hasStringInitializer && !nextIsAuto
		if *useAutoVar = nextIsAuto && !hasNumericInitializer && !hasStringInitializer; *useAutoVar {
			//  E[E["A"] = auto = x] = "A";
			//             ^^^^^^^^
			if *autoVar == nil {
				*autoVar = tx.EmitContext.NewUniqueName("auto", printer.AutoGenerateOptions{Flags: printer.GeneratedIdentifierFlagsOptimistic})
			}
			expression = tx.Factory.NewBinaryExpression(*autoVar, tx.Factory.NewToken(ast.KindEqualsToken), member.Initializer)
		} else {
			//  E[E["A"] = x] = "A";
			//             ^
			expression = member.Initializer
		}
	}

	// To support references to other enum members in the same declaration, we introduce a local binding for the enum member
	var local *ast.IdentifierNode
	if tx.referencedEnumMembers.Has(memberNode) || useConditionalReverseMapping {
		//  E[E["A"] = A = ++auto] = "A";
		//             ^^^-------
		local = tx.getVariableNameForPropertyName(member)
		if local != nil {
			tx.EmitContext.HoistVariable(local)
			expression = tx.Factory.NewBinaryExpression(local, tx.Factory.NewToken(ast.KindEqualsToken), expression)
		}
	}

	// Define the enum member property:
	//  E[E["A"] = ++auto] = "A";
	//    ^^^^^^^^--_____
	expression = tx.Factory.NewBinaryExpression(tx.Factory.NewElementAccessExpression(enumName, nil, name, ast.NodeFlagsNone), tx.Factory.NewToken(ast.KindEqualsToken), expression)

	// If this is syntactically a numeric literal initializer, or is auto numbered, then we unconditionally define the
	// reverse mapping for the enum member.
	if useExplicitReverseMapping {
		//  E[E["A"] = A = ++auto] = "A";
		//  ^^-------------------^^^^^^^
		expression = tx.Factory.NewBinaryExpression(tx.Factory.NewElementAccessExpression(enumName, nil, expression, ast.NodeFlagsNone), tx.Factory.NewToken(ast.KindEqualsToken), name)
	}

	memberStatement := tx.Factory.NewExpressionStatement(expression)
	tx.EmitContext.SetCommentAndSourceMapRanges(expression, member.Loc)
	tx.EmitContext.SetCommentAndSourceMapRanges(memberStatement, member.Loc)
	statements = append(statements, memberStatement)

	// If this is not auto numbered and is not syntactically a string or numeric literal initializer, then we
	// conditionally define the reverse mapping for the enum member.
	if useConditionalReverseMapping {
		//  E["A"] = A = x;
		//  if (typeof A !== "string") E[A] = "A";
		//  ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

		// If we cannot refer to a local variable for the enum member, then we need to look up the value when evaluating
		// the condition.
		condition := local
		if condition == nil {
			condition = tx.Factory.NewElementAccessExpression(enumName, nil, name, ast.NodeFlagsNone)
		}

		ifStatement := tx.Factory.NewIfStatement(
			tx.Factory.NewBinaryExpression(
				tx.Factory.NewTypeOfExpression(condition),
				tx.Factory.NewToken(ast.KindExclamationEqualsEqualsToken),
				tx.Factory.NewStringLiteral("string"),
			),
			tx.Factory.NewExpressionStatement(
				tx.Factory.NewBinaryExpression(
					tx.Factory.NewElementAccessExpression(enumName, nil,
						condition,
						ast.NodeFlagsNone,
					),
					tx.Factory.NewToken(ast.KindEqualsToken),
					name,
				),
			),
			nil,
		)

		tx.EmitContext.AddEmitFlags(ifStatement, printer.EFSingleLine)
		statements = append(statements, ifStatement)
	}

	tx.current = tx.parent
	tx.parent = savedParent
	return statements
}

func (tx *EnumTransformer) visitIdentifier(node *ast.IdentifierNode) *ast.Node {
	if tx.isVisitingEnumBody && isIdentifierReference(node, tx.parent) {
		if symbol := tx.resolveName(node, node.Text(), ast.SymbolFlagsEnumMember); symbol != nil {
			if enumMember := ast.GetDeclarationOfKind(symbol, ast.KindEnumMember); enumMember != nil {
				tx.referencedEnumMembers.Add(enumMember)
			}
		}
	}
	return node
}

func (tx *EnumTransformer) resolveName(location *ast.Node, name string, meaning ast.SymbolFlags) *ast.Symbol {
	if tx.nameResolver == nil {
		tx.nameResolver = &binder.NameResolver{CompilerOptions: tx.compilerOptions}
	}
	location = tx.EmitContext.MostOriginal(location)
	return tx.nameResolver.Resolve(location, name, meaning, nil /*nameNotFoundMessage*/, false /*isUse*/, true /*excludeGlobals*/)
}
