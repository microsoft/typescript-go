package ls

import (
	"strconv"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/locale"
	"github.com/microsoft/typescript-go/internal/ls/autoimport"
	"github.com/microsoft/typescript-go/internal/ls/change"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/nodebuilder"
)

type preserveOptionalFlags int

const (
	preserveOptionalFlagsMethod preserveOptionalFlags = 1 << iota
	preserveOptionalFlagsProperty
	preserveOptionalFlagsAll = preserveOptionalFlagsMethod | preserveOptionalFlagsProperty
)

type missingMemberFixer struct {
	ChangeTracker *change.Tracker
	Checker       *checker.Checker
	Program       *compiler.Program
	Preferences   lsutil.UserPreferences
	ImportAdder   autoimport.ImportAdder
	Locale        locale.Locale
}

func newMissingMemberFixer(changeTracker *change.Tracker, program *compiler.Program, typeChecker *checker.Checker, preferences lsutil.UserPreferences, importAdder autoimport.ImportAdder, locale locale.Locale) *missingMemberFixer {
	return &missingMemberFixer{
		ChangeTracker: changeTracker,
		Checker:       typeChecker,
		Program:       program,
		Preferences:   preferences,
		ImportAdder:   importAdder,
		Locale:        locale,
	}
}

func (f *missingMemberFixer) createNodeBuilder() (*checker.NodeBuilder, map[*ast.IdentifierNode]*ast.Symbol) {
	idToSymbol := make(map[*ast.IdentifierNode]*ast.Symbol)
	nodeBuilder := checker.NewNodeBuilderEx(f.Checker, f.ChangeTracker.EmitContext, idToSymbol)
	return nodeBuilder, idToSymbol
}

func (f *missingMemberFixer) createMemberFromSymbol(symbol *ast.Symbol, enclosingDeclaration *ast.Node, sourceFile *ast.SourceFile, body *ast.FunctionBody, preserveOptional preserveOptionalFlags) []*ast.Node {
	declarations := symbol.Declarations
	declaration := core.FirstOrNil(declarations)

	quotePreference := lsutil.GetQuotePreference(sourceFile, f.Preferences)
	ambient := enclosingDeclaration.Flags&ast.NodeFlagsAmbient != 0
	optional := symbol.Flags&ast.SymbolFlagsOptional != 0
	kind := ast.KindPropertySignature
	if declaration != nil {
		kind = declaration.Kind
	}
	declarationName := createDeclarationName(f.ChangeTracker.NodeFactory, f.Checker, symbol, declaration)
	modifiers := f.createModifiers(symbol, declaration)

	flags := nodebuilder.FlagsNoTruncation
	if quotePreference == lsutil.QuotePreferenceSingle {
		flags |= nodebuilder.FlagsUseSingleQuotesForStringLiteralType
	}

	t := f.Checker.GetTypeOfSymbolAtLocation(symbol, enclosingDeclaration)
	nodeBuilder, idToSymbol := f.createNodeBuilder()
	var nodes []*ast.Node

	switch kind {
	case ast.KindPropertySignature, ast.KindPropertyDeclaration:
		typeNode := f.createTypeNode(t, enclosingDeclaration, flags, nodeBuilder, idToSymbol)
		var questionToken *ast.TokenNode
		if optional && preserveOptional&preserveOptionalFlagsProperty != 0 {
			questionToken = f.ChangeTracker.NodeFactory.NewToken(ast.KindQuestionToken)
		}
		return append(nodes, f.ChangeTracker.NodeFactory.NewPropertyDeclaration(modifiers, createPropertyName(f.ChangeTracker.NodeFactory, declarationName, quotePreference), questionToken, typeNode, nil /*initializer*/))

	case ast.KindGetAccessor, ast.KindSetAccessor:
		accessors := ast.GetAllAccessorDeclarations(symbol.Declarations, declaration)
		var orderedAccessors []*ast.Node
		if accessors.SecondAccessor == nil {
			orderedAccessors = append(orderedAccessors, accessors.FirstAccessor)
		} else {
			orderedAccessors = append(orderedAccessors, accessors.FirstAccessor, accessors.SecondAccessor)
		}

		for _, accessor := range orderedAccessors {
			if ast.IsGetAccessorDeclaration(accessor) {
				nodes = append(nodes,
					f.ChangeTracker.NodeFactory.NewGetAccessorDeclaration(
						modifiers, createPropertyName(f.ChangeTracker.NodeFactory, declarationName, quotePreference),
						nil /*typeParameters*/, nil /*parameters*/, f.createTypeNode(t, enclosingDeclaration, flags, nodeBuilder, idToSymbol), nil /*fullSignature*/, f.createBody(body, ambient, quotePreference)),
				)
			}

			if ast.IsSetAccessorDeclaration(accessor) {
				parameter := checker.GetSetAccessorValueParameter(accessor)
				if parameter == nil {
					panic("Expected set accessor to have a parameter.")
				}

				nodes = append(nodes, f.ChangeTracker.NodeFactory.NewSetAccessorDeclaration(
					modifiers, createPropertyName(f.ChangeTracker.NodeFactory, declarationName, quotePreference),
					nil /*typeParameters*/, createDummyParameters(f.ChangeTracker.NodeFactory, 1, []string{parameter.Name().Text()}, []*ast.TypeNode{f.createTypeNode(t, enclosingDeclaration, flags, nodeBuilder, idToSymbol)}, 1, ast.IsInJSFile(enclosingDeclaration)),
					nil /*type*/, nil /*fullSignature*/, f.createBody(body, ambient, quotePreference)),
				)
			}
		}
		return nodes

	case ast.KindMethodSignature, ast.KindMethodDeclaration:
		signatures := f.Checker.GetSignaturesOfType(t, checker.SignatureKindCall)
		if len(signatures) == 0 {
			return nil
		}

		if len(declarations) == 1 {
			method := f.createSignatureDeclarationFromSignature(core.FirstOrNil(signatures), ast.KindMethodDeclaration, sourceFile, enclosingDeclaration, f.createBody(body, ambient, quotePreference), modifiers, declarationName, optional)
			if method != nil {
				nodes = append(nodes, method)
			}
			return nodes
		}

		for _, signature := range signatures {
			if signature.Declaration() != nil && signature.Declaration().Flags&ast.NodeFlagsAmbient != 0 {
				continue
			}

			method := f.createSignatureDeclarationFromSignature(signature, ast.KindMethodDeclaration, sourceFile, enclosingDeclaration, nil, modifiers, declarationName, optional)
			if method != nil {
				nodes = append(nodes, method)
			}
		}

		if ambient {
			return nodes
		}

		if len(declarations) > len(signatures) {
			signature := f.Checker.GetSignatureFromDeclaration(core.LastOrNil(declarations))
			method := f.createSignatureDeclarationFromSignature(signature, ast.KindMethodDeclaration, sourceFile, enclosingDeclaration, f.createBody(body, ambient, quotePreference), modifiers, declarationName, optional)
			if method != nil {
				nodes = append(nodes, method)
			}
		} else {
			method := f.createSignatureDeclarationFromSignatures(signatures, declarationName, optional, modifiers, quotePreference, body, enclosingDeclaration)
			if method != nil {
				nodes = append(nodes, method)
			}
		}

		return nodes
	}
	return nil
}

func (f *missingMemberFixer) createTypeNode(t *checker.Type, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, nodeBuilder *checker.NodeBuilder, idToSymbol map[*ast.IdentifierNode]*ast.Symbol) *ast.TypeNode {
	return f.importTypeNode(nodeBuilder.TypeToTypeNode(t, enclosingDeclaration, flags, nodebuilder.InternalFlagsNone, nil /*tracker*/), idToSymbol)
}

func (f *missingMemberFixer) createModifiers(symbol *ast.Symbol, declaration *ast.Node) *ast.ModifierList {
	modifierFlags := ast.ModifierFlagsNone
	if declaration != nil {
		effective := checker.GetDeclarationModifierFlagsFromSymbol(symbol)
		modifierFlags = effective & ast.ModifierFlagsStatic
		if effective&ast.ModifierFlagsPublic != 0 {
			modifierFlags |= ast.ModifierFlagsPublic
		} else if effective&ast.ModifierFlagsProtected != 0 {
			modifierFlags |= ast.ModifierFlagsProtected
		}
		if ast.IsAutoAccessorPropertyDeclaration(declaration) {
			modifierFlags |= ast.ModifierFlagsAccessor
		}
	}
	if f.shouldAddOverrideKeyword(declaration) {
		modifierFlags |= ast.ModifierFlagsOverride
	}
	if modifierFlags == ast.ModifierFlagsNone {
		return nil
	}
	return f.ChangeTracker.NodeFactory.NewModifierList(ast.CreateModifiersFromModifierFlags(modifierFlags, f.ChangeTracker.NodeFactory.NewModifier))
}

func (f *missingMemberFixer) shouldAddOverrideKeyword(declaration *ast.Node) bool {
	return f.Program.Options().NoImplicitOverride.IsTrue() && ast.HasAbstractModifier(declaration)
}

func (f *missingMemberFixer) createSignatureDeclarationFromSignature(signature *checker.Signature, kind ast.Kind, sourceFile *ast.SourceFile, enclosingDeclaration *ast.Node, body *ast.FunctionBody, modifiers *ast.ModifierList, name *ast.PropertyName, optional bool) *ast.Node {
	flags := nodebuilder.FlagsNoTruncation | nodebuilder.FlagsSuppressAnyReturnType | nodebuilder.FlagsAllowEmptyTuple
	if lsutil.GetQuotePreference(sourceFile, f.Preferences) == lsutil.QuotePreferenceSingle {
		flags |= nodebuilder.FlagsUseSingleQuotesForStringLiteralType
	}

	nodeBuilder, idToSymbol := f.createNodeBuilder()
	signatureDeclaration := nodeBuilder.SignatureToSignatureDeclaration(signature, kind, enclosingDeclaration, flags, nodebuilder.InternalFlagsAllowUnresolvedNames, nil /*tracker*/)
	if signatureDeclaration == nil {
		return nil
	}

	isJS := ast.IsInJSFile(enclosingDeclaration)
	parameters := signatureDeclaration.ParameterList()
	typeParameters := core.IfElse(isJS, nil, signatureDeclaration.TypeParameterList())
	typeNode := core.IfElse(isJS, nil, signatureDeclaration.Type())

	if typeParameters != nil && len(typeParameters.Nodes) > 0 {
		nodes := make([]*ast.Node, 0, len(typeParameters.Nodes))
		for _, tp := range typeParameters.Nodes {
			if tp == nil {
				continue
			}

			if ast.IsTypeParameterDeclaration(tp) {
				typeParameter := tp.AsTypeParameterDeclaration()

				constraint := typeParameter.Constraint
				if constraint != nil {
					constraint = f.importTypeNode(constraint, idToSymbol)
				}

				defaultType := typeParameter.DefaultType
				if defaultType != nil {
					defaultType = f.importTypeNode(defaultType, idToSymbol)
				}

				nodes = append(nodes,
					f.ChangeTracker.NodeFactory.UpdateTypeParameterDeclaration(typeParameter, typeParameter.Modifiers(), typeParameter.Name(), constraint, typeParameter.Expression, defaultType))
			} else {
				nodes = append(nodes, tp)
			}
		}
		typeParameters = f.ChangeTracker.NodeFactory.NewNodeList(nodes)
	}

	if parameters != nil {
		nodes := make([]*ast.Node, 0, len(parameters.Nodes))
		for _, p := range parameters.Nodes {
			if p == nil {
				continue
			}

			parameter := p.AsParameterDeclaration()
			parameterTypeNode := parameter.Type
			if parameterTypeNode != nil {
				parameterTypeNode = f.importTypeNode(parameterTypeNode, idToSymbol)
			}

			nodes = append(nodes,
				f.ChangeTracker.NodeFactory.UpdateParameterDeclaration(parameter, parameter.Modifiers(), parameter.DotDotDotToken, parameter.Name(), core.IfElse(isJS, nil, parameter.QuestionToken), parameterTypeNode, parameter.Initializer))
		}
		parameters = f.ChangeTracker.NodeFactory.NewNodeList(nodes)
	}

	if typeNode != nil {
		typeNode = f.importTypeNode(typeNode, idToSymbol)
	}

	var questionToken *ast.TokenNode
	if optional {
		questionToken = f.ChangeTracker.NodeFactory.NewToken(ast.KindQuestionToken)
	}

	switch kind {
	case ast.KindFunctionExpression:
		fn := signatureDeclaration.AsFunctionExpression()
		return f.ChangeTracker.NodeFactory.UpdateFunctionExpression(fn, modifiers, fn.AsteriskToken, core.IfElse(name != nil && ast.IsIdentifier(name), name, nil), typeParameters, parameters, typeNode, fn.FullSignature, core.OrElse(body, fn.Body))

	case ast.KindArrowFunction:
		fn := signatureDeclaration.AsArrowFunction()
		return f.ChangeTracker.NodeFactory.UpdateArrowFunction(fn, modifiers, typeParameters, parameters, typeNode, fn.FullSignature, fn.EqualsGreaterThanToken, core.OrElse(body, fn.Body))

	case ast.KindMethodDeclaration:
		method := signatureDeclaration.AsMethodDeclaration()
		methodName := core.IfElse(name == nil, f.ChangeTracker.NodeFactory.NewIdentifier(""), createPropertyName(f.ChangeTracker.NodeFactory, name, lsutil.GetQuotePreference(sourceFile, f.Preferences)))
		return f.ChangeTracker.NodeFactory.UpdateMethodDeclaration(method, modifiers, method.AsteriskToken, methodName, questionToken, typeParameters, parameters, typeNode, method.FullSignature, body)

	case ast.KindFunctionDeclaration:
		fn := signatureDeclaration.AsFunctionDeclaration()
		return f.ChangeTracker.NodeFactory.UpdateFunctionDeclaration(fn, modifiers, fn.AsteriskToken, core.IfElse(name != nil && ast.IsIdentifier(name), name, nil), typeParameters, parameters, typeNode, fn.FullSignature, core.OrElse(body, fn.Body))
	}

	return nil
}

func (f *missingMemberFixer) createSignatureDeclarationFromSignatures(signatures []*checker.Signature, name *ast.PropertyName, optional bool, modifiers *ast.ModifierList, quotePreference lsutil.QuotePreference, body *ast.FunctionBody, enclosingDeclaration *ast.Node) *ast.Node {
	if len(signatures) == 0 {
		return nil
	}

	nodeBuilder, idToSymbol := f.createNodeBuilder()
	maxArgsSignature := signatures[0]
	minArgumentCount := signatures[0].MinArgumentCount()

	hasRestParameter := false
	for _, signature := range signatures {
		minArgumentCount = min(minArgumentCount, signature.MinArgumentCount())
		if signature.HasRestParameter() {
			hasRestParameter = true
		}
		if len(signature.Parameters()) >= len(maxArgsSignature.Parameters()) && (!signature.HasRestParameter() || maxArgsSignature.HasRestParameter()) {
			maxArgsSignature = signature
		}
	}

	maxNonRestArgs := len(maxArgsSignature.Parameters()) - core.IfElse(maxArgsSignature.HasRestParameter(), 1, 0)
	parameterNames := make([]string, 0, len(maxArgsSignature.Parameters()))
	for _, symbol := range maxArgsSignature.Parameters() {
		parameterNames = append(parameterNames, symbol.Name)
	}
	parameters := createDummyParameters(f.ChangeTracker.NodeFactory, maxNonRestArgs, parameterNames, nil /*types*/, minArgumentCount, ast.IsInJSFile(enclosingDeclaration))

	if hasRestParameter {
		restParameterName := "rest"
		if maxNonRestArgs < len(parameterNames) && parameterNames[maxNonRestArgs] != "" {
			restParameterName = parameterNames[maxNonRestArgs]
		}

		var questionToken *ast.QuestionToken
		if maxNonRestArgs >= minArgumentCount {
			questionToken = f.ChangeTracker.NodeFactory.NewToken(ast.KindQuestionToken)
		}

		parameters.Nodes = append(parameters.Nodes, f.ChangeTracker.NodeFactory.NewParameterDeclaration(
			nil /*modifiers*/, f.ChangeTracker.NodeFactory.NewToken(ast.KindDotDotDotToken),
			f.ChangeTracker.NodeFactory.NewIdentifier(restParameterName), questionToken,
			f.ChangeTracker.NodeFactory.NewArrayTypeNode(f.ChangeTracker.NodeFactory.NewKeywordTypeNode(ast.KindUnknownKeyword)), nil, /*initializer*/
		))
	}

	methodName := core.IfElse(name == nil, f.ChangeTracker.NodeFactory.NewIdentifier(""), createPropertyName(f.ChangeTracker.NodeFactory, name, quotePreference))

	return f.ChangeTracker.NodeFactory.NewMethodDeclaration(
		modifiers, nil /*asteriskToken*/, methodName, core.IfElse(optional, f.ChangeTracker.NodeFactory.NewToken(ast.KindQuestionToken), nil),
		nil /*typeParameters*/, parameters, f.getReturnTypeFromSignatures(signatures, enclosingDeclaration, nodeBuilder, idToSymbol),
		nil /*fullSignature*/, f.createBody(body, false /*ambient*/, quotePreference),
	)
}

func (f *missingMemberFixer) getReturnTypeFromSignatures(signatures []*checker.Signature, enclosingDeclaration *ast.Node, nodeBuilder *checker.NodeBuilder, idToSymbol map[*ast.IdentifierNode]*ast.Symbol) *ast.TypeNode {
	if len(signatures) == 0 {
		return nil
	}

	returnTypes := make([]*checker.Type, 0, len(signatures))
	for _, signature := range signatures {
		returnTypes = append(returnTypes, f.Checker.GetReturnTypeOfSignature(signature))
	}

	unionType := f.Checker.GetUnionType(returnTypes)
	return f.importTypeNode(nodeBuilder.TypeToTypeNode(unionType, enclosingDeclaration, nodebuilder.FlagsNoTruncation, nodebuilder.InternalFlagsAllowUnresolvedNames, nil /*typeArguments*/), idToSymbol)
}

func (f *missingMemberFixer) importTypeNode(typeNode *ast.TypeNode, idToSymbol map[*ast.IdentifierNode]*ast.Symbol) *ast.TypeNode {
	if typeNode == nil || f.ImportAdder == nil {
		return typeNode
	}

	importedTypeNode, symbols := autoimport.TryGetAutoImportableReferenceFromTypeNode(typeNode, idToSymbol)
	if importedTypeNode != nil {
		for _, symbol := range symbols {
			f.ImportAdder.AddImportFromExportedSymbol(symbol, true /*isValidTypeOnlyUseSite*/)
		}
		return importedTypeNode
	}

	seen := make(map[*ast.Symbol]bool)
	for _, symbol := range idToSymbol {
		if symbol == nil || seen[symbol] {
			continue
		}
		seen[symbol] = true
		f.ImportAdder.AddImportFromExportedSymbol(symbol, true /*isValidTypeOnlyUseSite*/)
	}
	return typeNode
}

func (f *missingMemberFixer) createIndexSignatureDeclarationFromType(classDeclaration *ast.Node, implementedType *checker.Type, keyType *checker.Type) *ast.Node {
	indexInfo := f.Checker.GetIndexInfoOfType(implementedType, keyType)
	if indexInfo == nil {
		return nil
	}

	builder := checker.NewNodeBuilder(f.Checker, f.ChangeTracker.EmitContext)
	return builder.IndexInfoToIndexSignatureDeclaration(indexInfo, classDeclaration, nodebuilder.FlagsNone, nodebuilder.InternalFlagsNone, nil)
}

func (f *missingMemberFixer) createBody(body *ast.FunctionBody, ambient bool, quotePreference lsutil.QuotePreference) *ast.FunctionBody {
	if ambient {
		return nil
	}
	body = f.ChangeTracker.NodeFactory.DeepCloneNode(body)
	if body == nil {
		return f.createStubbedMethodBody(quotePreference)
	}
	return body
}

func (f *missingMemberFixer) createStubbedMethodBody(quotePreference lsutil.QuotePreference) *ast.FunctionBody {
	tokenFlags := ast.TokenFlagsNone
	if quotePreference == lsutil.QuotePreferenceSingle {
		tokenFlags = ast.TokenFlagsSingleQuote
	}

	return f.ChangeTracker.NodeFactory.NewBlock(f.ChangeTracker.NodeFactory.NewNodeList([]*ast.Node{
		f.ChangeTracker.NodeFactory.NewThrowStatement(
			f.ChangeTracker.NodeFactory.NewNewExpression(
				f.ChangeTracker.NodeFactory.NewIdentifier("Error"), nil /*typeArguments*/, f.ChangeTracker.NodeFactory.NewNodeList([]*ast.Node{
					f.ChangeTracker.NodeFactory.NewStringLiteral(diagnostics.Method_not_implemented.Localize(f.Locale), tokenFlags),
				}),
			),
		),
	}), true /*multiLine*/)
}

func createDummyParameters(factory *ast.NodeFactory, argCount int, names []string, types []*ast.TypeNode, minArgumentCount int, inJS bool) *ast.ParameterList {
	parameters := make([]*ast.Node, 0, argCount)
	parameterNameCounts := make(map[string]int)

	for i := range argCount {
		parameterName := ""
		if i < len(names) && names[i] != "" {
			parameterName = names[i]
		} else {
			parameterName = "arg" + strconv.Itoa(i)
		}

		count := parameterNameCounts[parameterName]
		parameterNameCounts[parameterName] = count + 1

		if count > 0 {
			parameterName += strconv.Itoa(count)
		}

		var questionToken *ast.QuestionToken
		if i >= minArgumentCount {
			questionToken = factory.NewToken(ast.KindQuestionToken)
		}

		var typeNode *ast.TypeNode
		if inJS {
			typeNode = nil
		} else if i < len(types) && types[i] != nil {
			typeNode = types[i]
		} else {
			typeNode = factory.NewKeywordTypeNode(ast.KindUnknownKeyword)
		}
		parameters = append(parameters,
			factory.NewParameterDeclaration(nil /*modifiers*/, nil /*dotDotDotToken*/, factory.NewIdentifier(parameterName), questionToken, typeNode, nil /*initializer*/))
	}
	return factory.NewNodeList(parameters)
}

func createDeclarationName(factory *ast.NodeFactory, typeChecker *checker.Checker, symbol *ast.Symbol, declaration *ast.Node) *ast.PropertyName {
	if symbol != nil && symbol.CheckFlags&ast.CheckFlagsMapped != 0 {
		nameType := typeChecker.GetNameTypeOfSymbol(symbol)
		if nameType != nil && checker.IsTypeUsableAsPropertyName(nameType) {
			return factory.NewIdentifier(checker.GetPropertyNameFromType(nameType))
		}
	}
	if declaration != nil && declaration.Name() != nil {
		return declaration.Name().Clone(factory)
	}
	if symbol != nil {
		return factory.NewIdentifier(symbol.Name)
	}
	return nil
}

func createPropertyName(factory *ast.NodeFactory, node *ast.Node, quotePreference lsutil.QuotePreference) *ast.PropertyName {
	if ast.IsIdentifier(node) && node.Text() == "constructor" {
		tokenFlags := ast.TokenFlagsNone
		if quotePreference == lsutil.QuotePreferenceSingle {
			tokenFlags = ast.TokenFlagsSingleQuote
		}
		return factory.NewComputedPropertyName(factory.NewStringLiteral(node.Text(), tokenFlags))
	}
	return factory.DeepCloneNode(node)
}
