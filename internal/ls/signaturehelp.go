package ls

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/printer"
	"github.com/microsoft/typescript-go/internal/scanner"
)

type invocationKind = int

const (
	invocationKindCall invocationKind = iota
	invocationKindTypeArgs
	invocationKinContextual
)

type callInvocation struct {
	kind invocationKind
	node *ast.Node
}

type typeArgsInvocation struct {
	kind   invocationKind
	called *ast.Identifier
}

type contextualInvocation struct {
	kind      invocationKind
	signature *checker.Signature
	node      *ast.Node // Just for enclosingDeclaration for printing types
	symbol    *ast.Symbol
}

type invocation struct {
	kind                 invocationKind
	callInvocation       callInvocation
	typeArgsInvocation   typeArgsInvocation
	contextualInvocation contextualInvocation
}

type symbolDisplayPart struct {
	// Text of an item describing the symbol.
	text string
	// The symbol's kind (such as 'className' or 'parameterName' or plain 'text').
	kind string
}

// Signature help information for a single parameter
type SignatureHelpParameter struct {
	name          string
	documentation []symbolDisplayPart
	displayParts  []symbolDisplayPart
	isOptional    bool
	isRest        bool
}

type jSDocTagInfo struct {
	name string
	text []symbolDisplayPart
}
type SignatureHelpTriggerReason struct {
	Invoked        *SignatureHelpInvokedReason
	CharacterTyped *SignatureHelpCharacterTypedReason
	Retriggered    *SignatureHelpRetriggeredReason
}

type SignatureHelpTriggerCharacter string

const (
	CommaTriggerCharacter      SignatureHelpTriggerCharacter = ","
	OpenParenTriggerCharacter  SignatureHelpTriggerCharacter = "("
	LessThanTriggerCharacter   SignatureHelpTriggerCharacter = "<"
	CloseParenTriggerCharacter string                        = ")"
)

type SignatureHelpRetriggerCharacter struct {
	CommaTrigger      SignatureHelpTriggerCharacter
	OpenParenTrigger  SignatureHelpTriggerCharacter
	LessThanTrigger   SignatureHelpTriggerCharacter
	CloseParenTrigger string
}

// Signals that the user manually requested signature help.
// The language service will unconditionally attempt to provide a result.
type SignatureHelpInvokedReason struct {
	Kind             string
	TriggerCharacter *string
}

// Signals that the signature help request came from a user typing a character.
// Depending on the character and the syntactic context, the request may or may not be served a result.
type SignatureHelpCharacterTypedReason struct {
	Kind string
	// Character that was responsible for triggering signature help.
	TriggerCharacter *SignatureHelpTriggerCharacter
}

// Signals that this signature help request came from typing a character or moving the cursor.
// This should only occur if a signature help session was already active and the editor needs to see if it should adjust.
// The language service will unconditionally attempt to provide a result.
// `triggerCharacter` can be `undefined` for a retrigger caused by a cursor move.
type SignatureHelpRetriggeredReason struct {
	Kind string
	// Character that was responsible for triggering signature help.
	TriggerCharacter *SignatureHelpRetriggerCharacter
}

func (l *LanguageService) GetSignatureHelpItems(fileName string, position int, triggerReason *SignatureHelpTriggerReason) *SignatureHelpItems {
	program, sourceFile := l.getProgramAndFile(fileName)
	typeChecker := program.GetTypeChecker()

	// Decide whether to show signature help
	startingToken := astnav.FindPrecedingToken(sourceFile, position)
	if startingToken == nil {
		// We are at the beginning of the file
		return nil
	}

	// Only need to be careful if the user typed a character and signature help wasn't showing.
	onlyUseSyntacticOwners := triggerReason != nil && triggerReason.CharacterTyped != nil //&&  == "characterTyped"

	// Bail out quickly in the middle of a string or comment, don't provide signature help unless the user explicitly requested it.
	if onlyUseSyntacticOwners && IsInString(sourceFile, position, startingToken) { // isInComment(sourceFile, position) needs formatting implemented
		return nil
	}

	isManuallyInvoked := triggerReason != nil && triggerReason.Invoked != nil //invoked.kind == "invoked"
	argumentInfo := getContainingArgumentInfo(startingToken, sourceFile, program.GetTypeChecker(), isManuallyInvoked, position)
	if argumentInfo == nil {
		return nil
	}

	// cancellationToken.throwIfCancellationRequested();

	// Extra syntactic and semantic filtering of signature help
	candidateInfo := getCandidateOrTypeInfo(argumentInfo, typeChecker, sourceFile, startingToken, onlyUseSyntacticOwners)
	// cancellationToken.throwIfCancellationRequested();

	// if (!candidateInfo) { tbd
	// 	// We didn't have any sig help items produced by the TS compiler.  If this is a JS
	// 	// file, then see if we can figure out anything better.
	// 	return isSourceFileJS(sourceFile) ? createJSSignatureHelpItems(argumentInfo, program, cancellationToken) : undefined;
	// }

	// return typeChecker.runWithCancellationToken(cancellationToken, typeChecker =>
	if candidateInfo.candidateInfo.kind == candidateKind {
		return createSignatureHelpItems(candidateInfo.candidateInfo.candidates, candidateInfo.candidateInfo.resolvedSignature, argumentInfo, sourceFile, typeChecker, onlyUseSyntacticOwners)
	}
	return createTypeHelpItems(candidateInfo.typeInfo.symbol, argumentInfo, sourceFile, typeChecker)
}

func createTypeHelpItems(symbol *ast.Symbol, argumentInfo *argumentListInfo, sourceFile *ast.SourceFile, c *checker.Checker) *SignatureHelpItems {
	typeParameters := c.GetLocalTypeParametersOfClassOrInterfaceOrTypeAlias(symbol)
	if typeParameters == nil {
		return nil
	}
	items := getTypeHelpItem(symbol, typeParameters, sourceFile, c)
	return &SignatureHelpItems{
		Signatures:      []signatureHelpItem{items},
		ActiveSignature: 0,
		ActiveParameter: argumentInfo.argumentIndex,
	}
}

func getTypeHelpItem(symbol *ast.Symbol, typeParameter []*checker.Type, sourceFile *ast.SourceFile, c *checker.Checker) signatureHelpItem {
	printer := printer.NewPrinter(
		printer.PrinterOptions{
			NewLine: core.NewLineKindLF,
		},
		printer.PrintHandlers{},
		nil)
	var parameters []signatureHelpParameter = []signatureHelpParameter{}
	for _, typeParam := range typeParameter {
		parameters = append(parameters, createSignatureHelpParameterForTypeParameter(typeParam.Symbol(), sourceFile, c, printer))
	}
	return signatureHelpItem{
		Label:         symbol.Name,
		Documentation: nil,
		Parameters:    &parameters,
		IsVariadic:    false,
	}
}

func createSignatureHelpItems(candidates *[]*checker.Signature, resolvedSignature *checker.Signature, argumentInfo *argumentListInfo, sourceFile *ast.SourceFile, c *checker.Checker, useFullPrefic bool) *SignatureHelpItems {
	enclosingDeclaration := getEnclosingDeclarationFromInvocation(argumentInfo.invocation)
	if enclosingDeclaration == nil {
		return nil
	}
	var callTargetSymbol *ast.Symbol
	if argumentInfo.invocation.kind == invocationKinContextual {
		callTargetSymbol = argumentInfo.invocation.contextualInvocation.symbol
	} else {
		callTargetSymbol = c.GetSymbolAtLocation(getExpressionFromInvocation(argumentInfo))
		if callTargetSymbol == nil && useFullPrefic && resolvedSignature.Declaration() != nil {
			callTargetSymbol = resolvedSignature.Declaration().Symbol()
		}
	}

	callTargetDisplayParts := []string{}
	if callTargetSymbol != nil {
		callTargetDisplayParts = append(callTargetDisplayParts, c.SymbolToString(callTargetSymbol))
	}
	var items [][]signatureHelpItem
	for _, candidateSignature := range *candidates {
		items = append(items, getSignatureHelpItem(candidateSignature, argumentInfo.isTypeParameterList, callTargetDisplayParts, enclosingDeclaration, argumentInfo.argumentCount, sourceFile, c))
	}

	selectedItemIndex := 0
	itemSeen := 0
	for i := 0; i < len(items); i++ {
		item := items[i]
		if (*candidates)[i] == resolvedSignature {
			selectedItemIndex = itemSeen
			if len(item) > 1 {
				count := 0
				for _, j := range item {
					if j.IsVariadic || len(*j.Parameters) >= argumentInfo.argumentCount {
						selectedItemIndex = itemSeen + count
						break
					}
					count++
				}
			}
		}
		itemSeen = itemSeen + len(item)
	}

	//Debug.assert(selectedItemIndex !== -1)
	var flattenedSignatures = []signatureHelpItem{}
	for _, item := range items {
		flattenedSignatures = append(flattenedSignatures, item...)
	}
	help := &SignatureHelpItems{
		Signatures:      flattenedSignatures,
		ActiveSignature: selectedItemIndex,
		ActiveParameter: argumentInfo.argumentIndex,
	}
	selected := help.Signatures[selectedItemIndex]
	if selected.IsVariadic {
		firstRest := core.FindIndex(*selected.Parameters, func(p signatureHelpParameter) bool {
			return p.isRest
		})
		if -1 < firstRest && firstRest < len(*selected.Parameters)-1 {
			// We don't have any code to get this correct; instead, don't highlight a current parameter AT ALL
			help.ActiveParameter = len(*selected.Parameters)
		} else {
			if help.ActiveParameter > len(*selected.Parameters)-1 {
				help.ActiveParameter = len(*selected.Parameters) - 1
			}
		}
	}
	return help
}

func getSignatureHelpItem(candidate *checker.Signature, isTypeParameterList bool, callTargetSymbol []string, enclosingDeclaration *ast.Node, activeArgumentIndex int, sourceFile *ast.SourceFile, c *checker.Checker) []signatureHelpItem {
	var infos []*signatureHelpItemInfo
	if isTypeParameterList {
		infos = itemInfoForTypeParameters(candidate, c, enclosingDeclaration, sourceFile)
	} else {
		infos = itemInfoForParameters(candidate, c, enclosingDeclaration, sourceFile)
	}

	suffixDisplayParts := returnTypeToDisplayParts(candidate, c)

	result := []signatureHelpItem{}
	for _, info := range infos {
		display := strings.Join(info.displayParts, "") + suffixDisplayParts
		result = append(result, signatureHelpItem{
			Label:           strings.Join(callTargetSymbol, "") + display,
			Documentation:   nil,
			Parameters:      info.parameters,
			IsVariadic:      info.isVariadic,
			ActiveParameter: activeArgumentIndex,
		})
	}
	return result
}

func returnTypeToDisplayParts(candidateSignature *checker.Signature, c *checker.Checker) string {
	var returnType string
	returnType = ": " //returnType.WriteString(":")
	//returnType.WriteString(" ")
	predicate := c.GetTypePredicateOfSignature(candidateSignature)
	if predicate != nil {
		returnType += c.TypePredicateToString(predicate) //returnType.WriteString(c.TypePredicateToString(predicate))
	} else {
		returnType += c.TypeToString(c.GetReturnTypeOfSignature(candidateSignature)) //returnType.WriteString(c.TypeToString(c.GetReturnTypeOfSignature(candidateSignature)))
	}
	return returnType
}

func itemInfoForTypeParameters(candidateSignature *checker.Signature, c *checker.Checker, enclosingDeclaration *ast.Node, sourceFile *ast.SourceFile) []*signatureHelpItemInfo {
	printer := printer.NewPrinter(
		printer.PrinterOptions{
			NewLine: core.NewLineKindLF,
		},
		printer.PrintHandlers{},
		nil)

	var typeParameters []*checker.Type
	if candidateSignature.Target() != nil {
		typeParameters = candidateSignature.Target().TypeParameters()
	} else {
		typeParameters = candidateSignature.TypeParameters()
	}
	var getTypeParameters []signatureHelpParameter = []signatureHelpParameter{}
	for _, typeParameter := range typeParameters {
		getTypeParameters = append(getTypeParameters, createSignatureHelpParameterForTypeParameter(typeParameter.Symbol(), sourceFile, c, printer))
	}
	var typeParameterDisplayParts = []string{}
	for _, typeParameter := range getTypeParameters {
		typeParameterDisplayParts = append(typeParameterDisplayParts, typeParameter.Label)
	}
	thisParameter := []signatureHelpParameter{} //tbd
	if candidateSignature.ThisParameter() != nil {
		thisParameter = []signatureHelpParameter{createSignatureHelpParameterForParameter(candidateSignature.ThisParameter(), printer, sourceFile, c)}
	}

	// Creating display parts for parameters. For example, <T>(a: string, b: number)
	lists := c.GetExpandedParameters(candidateSignature, false)

	var result []*signatureHelpItemInfo
	parameterLabels := []string{}
	for _, parameterList := range lists {
		parameters := thisParameter
		for _, param := range parameterList {
			parameter := createSignatureHelpParameterForTypeParameter(param, sourceFile, c, printer)
			parameterLabels = append(parameterLabels, parameter.Label)
			parameters = append(parameters, parameter)
		}

		result = append(result, &signatureHelpItemInfo{
			isVariadic:   false,
			parameters:   &getTypeParameters,
			displayParts: []string{scanner.TokenToString(ast.KindLessThanToken), strings.Join(typeParameterDisplayParts, ","), scanner.TokenToString(ast.KindGreaterThanToken), scanner.TokenToString(ast.KindOpenParenToken), strings.Join(parameterLabels, ", "), scanner.TokenToString(ast.KindCloseParenToken)},
		})
	}
	return result
}

func itemInfoForParameters(candidateSignature *checker.Signature, c *checker.Checker, enclosingDeclaratipn *ast.Node, sourceFile *ast.SourceFile) []*signatureHelpItemInfo {
	printer := printer.NewPrinter(
		printer.PrinterOptions{
			NewLine: core.NewLineKindLF,
		},
		printer.PrintHandlers{},
		nil)

	var getTypeParameters []signatureHelpParameter = []signatureHelpParameter{}
	if candidateSignature.TypeParameters() != nil && len(candidateSignature.TypeParameters()) != 0 {
		for _, typeParameter := range candidateSignature.TypeParameters() {
			getTypeParameters = append(getTypeParameters, createSignatureHelpParameterForTypeParameter(typeParameter.Symbol(), sourceFile, c, printer))
		}
	}
	typeParameterDisplayParts := []string{}
	for _, typeParameter := range getTypeParameters {
		typeParameterDisplayParts = append(typeParameterDisplayParts, typeParameter.Label)
	}
	typeParameterParts := []string{}
	if len(typeParameterDisplayParts) != 0 {
		typeParameterParts = []string{scanner.TokenToString(ast.KindLessThanToken), strings.Join(typeParameterDisplayParts, ", "), scanner.TokenToString(ast.KindGreaterThanToken)}
	}

	// This part is only for parameters. For example, (a: string, b: number)
	lists := c.GetExpandedParameters(candidateSignature, false)

	isVariadic := func(parameterList []*ast.Symbol) bool {
		if !c.HasEffectiveRestParameter(candidateSignature) {
			return false
		}
		if len(lists) == 1 {
			return true
		}
		return len(parameterList) != 0 && parameterList[len(parameterList)-1] != nil && (parameterList[len(parameterList)-1].CheckFlags&ast.CheckFlagsRestParameter != 0)
	}

	var result []*signatureHelpItemInfo

	for _, parameterList := range lists {
		var parameters []signatureHelpParameter
		parameterLabels := []string{}
		for _, param := range parameterList {
			parameter := createSignatureHelpParameterForParameter(param, printer, sourceFile, c)
			parameterLabels = append(parameterLabels, parameter.Label)
			parameters = append(parameters, parameter)
		}

		result = append(result, &signatureHelpItemInfo{
			isVariadic:   isVariadic(parameterList),
			parameters:   &parameters,
			displayParts: append(typeParameterParts, scanner.TokenToString(ast.KindOpenParenToken), strings.Join(parameterLabels, ", "), scanner.TokenToString(ast.KindCloseParenToken)),
		})
	}

	return result
}

func getEffectiveParameterDeclaration(symbol *ast.Symbol) *ast.Node {
	parameterDeclaration := ast.GetDeclarationOfKind(symbol, ast.KindParameter)
	if parameterDeclaration != nil {
		return parameterDeclaration
	}
	if symbol.Flags&ast.SymbolFlagsTransient == 0 {
		return ast.GetDeclarationOfKind(symbol, ast.KindJSDocParameterTag)
	}
	return nil
}

func createSignatureHelpParameterForParameter(parameter *ast.Symbol, p *printer.Printer, sourceFile *ast.SourceFile, c *checker.Checker) signatureHelpParameter {
	display := p.Emit(tempSymbolToParameterDeclaration(parameter, c), sourceFile)
	isOptional := parameter.CheckFlags&ast.CheckFlagsOptionalParameter != 0 //any extra checks needed?
	isRest := parameter.CheckFlags&ast.CheckFlagsRestParameter != 0
	return signatureHelpParameter{
		Label:         display,
		Documentation: nil,
		isRest:        isRest,
		isOptional:    isOptional,
	}
}

func tempSymbolToParameterDeclaration(parameterSymbol *ast.Symbol, c *checker.Checker) *ast.Node {
	factory := printer.NewEmitContext().Factory
	parameterDeclaration := getEffectiveParameterDeclaration(parameterSymbol)
	parameterType := c.GetTypeOfSymbol(parameterSymbol)
	parameterTypeNode := createTypeNode(parameterType, parameterSymbol, parameterDeclaration, factory)
	isRest := parameterDeclaration != nil && checker.IsRestParameter(parameterDeclaration) || parameterSymbol.CheckFlags&ast.CheckFlagsRestParameter != 0
	var dotDotDotToken *ast.Node
	if isRest {
		dotDotDotToken = factory.NewToken(ast.KindDotDotDotToken)
	}
	var name *ast.Node
	if parameterDeclaration == nil || parameterDeclaration.Name() == nil {
		name = factory.NewIdentifier(parameterSymbol.Name)
	} else {
		name = factory.DeepCloneNode(parameterDeclaration.Name())
	}
	isOptional := parameterDeclaration != nil && parameterSymbol.CheckFlags&ast.CheckFlagsOptionalParameter != 0 // || c.GetEmitResolver(nil, true).IsOptionalParameter(parameterDeclaration)
	var questionToken *ast.Node
	if isOptional {
		questionToken = factory.NewToken(ast.KindQuestionToken)
	}

	parameterNode := factory.NewParameterDeclaration(
		nil,
		dotDotDotToken,
		name,
		questionToken,
		parameterTypeNode,
		/*initializer*/ nil,
	)
	return parameterNode
}

func createTypeNode(t *checker.Type, symbol *ast.Symbol, parameterDeclaration *ast.Node, factory *ast.NodeFactory) *ast.Node {
	if t == nil {
		return factory.NewKeywordTypeNode(ast.KindAnyKeyword)
	}
	if t.Flags()&checker.TypeFlagsString != 0 {
		return factory.NewKeywordTypeNode(ast.KindStringKeyword)
	}
	if t.Flags()&checker.TypeFlagsNumber != 0 {
		return factory.NewKeywordTypeNode(ast.KindNumberKeyword)
	}
	if t.Flags()&checker.TypeFlagsAny != 0 {
		return factory.NewKeywordTypeNode(ast.KindAnyKeyword)
	}
	if t.Flags()&checker.TypeFlagsUnknown != 0 {
		return factory.NewKeywordTypeNode(ast.KindUnknownKeyword)
	}
	if t.Flags()&checker.TypeFlagsBigInt != 0 {
		return factory.NewKeywordTypeNode(ast.KindBigIntKeyword)
	}
	if t.Flags()&checker.TypeFlagsBoolean != 0 {
		return factory.NewKeywordTypeNode(ast.KindBooleanKeyword)
	}
	if t.Flags()&checker.TypeFlagsVoid != 0 {
		return factory.NewKeywordTypeNode(ast.KindVoidKeyword)
	}
	if t.Flags()&checker.TypeFlagsUndefined != 0 {
		return factory.NewKeywordTypeNode(ast.KindUndefinedKeyword)
	}
	if t.Flags()&checker.TypeFlagsNull != 0 {
		return factory.NewKeywordTypeNode(ast.KindNullKeyword)
	}
	if t.Flags()&checker.TypeFlagsNever != 0 {
		return factory.NewKeywordTypeNode(ast.KindNeverKeyword)
	}
	if t.Flags()&checker.TypeFlagsESSymbol != 0 {
		return factory.NewKeywordTypeNode(ast.KindSymbolKeyword)
	}
	if t.Flags()&checker.TypeFlagsNonPrimitive != 0 {
		return factory.NewKeywordTypeNode(ast.KindObjectKeyword)
	}
	if t.Flags()&checker.TypeFlagsObject != 0 {
		// parameterDeclaration.Kind = ast.KindTypeReference
		// return parameterDeclaration
		return factory.NewKeywordTypeNode(ast.KindObjectKeyword)
	}
	if t.Flags()&checker.TypeFlagsTypeParameter != 0 {
		return factory.NewKeywordTypeNode(ast.KindTypeParameter)
	}
	if t.Flags()&checker.TypeFlagsUnion != 0 || t.Flags()&checker.TypeFlagsIntersection != 0 {
		return factory.NewUnionTypeNode(nil) // come back
	}
	if t.Flags()&checker.TypeFlagsIndex != 0 {
		return factory.NewKeywordTypeNode(ast.KindIndexSignature)
	}
	if t.Flags()&checker.TypeFlagsTemplateLiteral != 0 {
		return factory.NewKeywordTypeNode(ast.KindTemplateLiteralType)
	}
	if t.Flags()&checker.TypeFlagsIndexedAccess != 0 {
		return factory.NewKeywordTypeNode(ast.KindIndexedAccessType)
	}
	if t.Flags()&checker.TypeFlagsConditional != 0 {
		return factory.NewKeywordTypeNode(ast.KindConditionalType)
	}
	return nil
}

func createSignatureHelpParameterForTypeParameter(parameter *ast.Symbol, sourceFile *ast.SourceFile, c *checker.Checker, p *printer.Printer) signatureHelpParameter {
	parameterDeclaration := getEffectiveParameterDeclaration(parameter)
	var display string
	if parameterDeclaration == nil {
		display = ast.SymbolName(parameter)
	} else {
		display = p.Emit(getEffectiveParameterDeclaration(parameter), sourceFile)
	}
	return signatureHelpParameter{
		Label:         display,
		Documentation: nil,
		isRest:        false,
		isOptional:    false,
		//displayParts:  []string{param.Text()},
	}
}

type SignatureHelpItems struct {
	Signatures      []signatureHelpItem
	ActiveSignature int
	ActiveParameter int
}

// Represents the signature of something callable. A signature
// can have a label, like a function-name, a doc-comment, and
// a set of parameters.
type signatureHelpItem struct {
	// The Label of this signature. Will be shown in
	// the UI.
	Label string

	// The human-readable doc-comment of this signature. Will be shown
	// in the UI but can be omitted.
	Documentation *string

	// The Parameters of this signature.
	Parameters *[]signatureHelpParameter

	// The index of the active parameter.
	ActiveParameter int

	// Needed only here, not in lsp
	IsVariadic bool
}

type signatureHelpItemInfo struct { //later change the name to sigHelpDisplayParts
	isVariadic   bool
	parameters   *[]signatureHelpParameter
	displayParts []string
}

// Represents a parameter of a callable-signature. A parameter can
// have a label and a doc-comment.
type signatureHelpParameter struct {
	Label         string
	Documentation *string
	isRest        bool
	isOptional    bool
}

func getEnclosingDeclarationFromInvocation(invocation *invocation) *ast.Node {
	if invocation.kind == invocationKindCall {
		return invocation.callInvocation.node
	} else if invocation.kind == invocationKindTypeArgs {
		return invocation.typeArgsInvocation.called.AsNode()
	} else {
		return invocation.contextualInvocation.node
	}
}

func getExpressionFromInvocation(argumentInfo *argumentListInfo) *ast.Node {
	if argumentInfo.invocation.kind == invocationKindCall {
		return checker.GetInvokedExpression(argumentInfo.invocation.callInvocation.node)
	}
	return argumentInfo.invocation.typeArgsInvocation.called.AsNode()
}

type candidateOrTypeKind = int

const (
	candidateKind candidateOrTypeKind = iota
	typeKind
)

type candidateInfo struct {
	kind              candidateOrTypeKind
	candidates        *[]*checker.Signature
	resolvedSignature *checker.Signature
}

type typeInfo struct {
	kind   candidateOrTypeKind
	symbol *ast.Symbol
}

type CandidateOrTypeInfo struct {
	candidateInfo
	typeInfo
}

func getCandidateOrTypeInfo(info *argumentListInfo, c *checker.Checker, sourceFile *ast.SourceFile, startingToken *ast.Node, onlyUseSyntacticOwners bool) *CandidateOrTypeInfo {
	switch info.invocation.kind {
	case invocationKindCall:
		if onlyUseSyntacticOwners && !isSyntacticOwner(startingToken, info.invocation.callInvocation.node, sourceFile) {
			return nil
		}
		var candidates *[]*checker.Signature = &[]*checker.Signature{} //check
		resolvedSignature := c.GetResolvedSignatureForSignatureHelp(info.invocation.callInvocation.node, candidates, info.argumentCount)
		// if candidates == nil { //check yes can check if exists and if len = 0
		// 	return nil
		// }
		return &CandidateOrTypeInfo{
			candidateInfo: candidateInfo{
				kind:              candidateKind,
				candidates:        candidates,
				resolvedSignature: resolvedSignature,
			},
		}
	case invocationKindTypeArgs:
		called := info.invocation.typeArgsInvocation.called.AsNode()
		container := called
		if ast.IsIdentifier(called) {
			container = called.Parent
		}
		if onlyUseSyntacticOwners && !containsPrecedingToken(startingToken, sourceFile, container) {
			return nil
		}
		candidates := getPossibleGenericSignatures(called, info.argumentCount, c)
		if len(candidates) != 0 {
			return &CandidateOrTypeInfo{
				candidateInfo: candidateInfo{
					kind:              candidateKind,
					candidates:        &candidates,
					resolvedSignature: candidates[0],
				},
			}
		}
		symbol := c.GetSymbolAtLocation(called)
		return &CandidateOrTypeInfo{
			typeInfo: typeInfo{
				kind:   typeKind,
				symbol: symbol,
			},
		}
	case invocationKinContextual:
		return &CandidateOrTypeInfo{
			candidateInfo: candidateInfo{
				kind:              candidateKind,
				candidates:        &[]*checker.Signature{info.invocation.contextualInvocation.signature},
				resolvedSignature: info.invocation.contextualInvocation.signature,
			},
		}
	default:
		return nil //return Debug.assertNever(invocation);
	}
}

func isSyntacticOwner(startingToken *ast.Node, node *ast.Node, sourceFile *ast.SourceFile) bool { //tbd
	if !checker.IsCallOrNewExpression(node) {
		return false
	}
	invocationChildren := getChildren(node, sourceFile)
	switch startingToken.Kind {
	case ast.KindOpenParenToken:
		return containsNode(invocationChildren, startingToken)
	case ast.KindCommaToken:
		return containsNode(invocationChildren, startingToken) //tbd does this need to check only for argument list?
		// const containingList = findContainingList(startingToken);
		// return !!containingList && contains(invocationChildren, containingList);
	case ast.KindLessThanToken:
		return containsPrecedingToken(startingToken, sourceFile, node.AsCallExpression().Expression)
	default:
		return false
	}
}

func containsPrecedingToken(startingToken *ast.Node, sourceFile *ast.SourceFile, container *ast.Node) bool {
	pos := startingToken.Pos()
	// There's a possibility that `startingToken.parent` contains only `startingToken` and
	// missing nodes, none of which are valid to be returned by `findPrecedingToken`. In that
	// case, the preceding token we want is actually higher up the tree—almost definitely the
	// next parent, but theoretically the situation with missing nodes might be happening on
	// multiple nested levels.
	currentParent := startingToken.Parent
	for currentParent != nil {
		precedingToken := astnav.FindPrecedingToken(sourceFile, pos)
		if precedingToken != nil {
			return RangeContainsRange(container.Loc, precedingToken.Loc)
		}
		currentParent = currentParent.Parent
	}
	// return Debug.fail("Could not find preceding token");
	return false
}
func getContainingArgumentInfo(node *ast.Node, sourceFile *ast.SourceFile, checker *checker.Checker, isManuallyInvoked bool, position int) *argumentListInfo {
	for n := node; !ast.IsSourceFile(n) && (isManuallyInvoked || !ast.IsBlock(n)); n = n.Parent {
		// If the node is not a subspan of its parent, this is a big problem.
		// There have been crashes that might be caused by this violation.
		//Debug.assert(rangeContainsRange(n.parent, n), "Not a subspan", () => `Child: ${Debug.formatSyntaxKind(n.kind)}, parent: ${Debug.formatSyntaxKind(n.parent.kind)}`);
		argumentInfo := getImmediatelyContainingArgumentOrContextualParameterInfo(n, position, sourceFile, checker)
		if argumentInfo != nil {
			return argumentInfo
		}
	}
	return nil
}

func getImmediatelyContainingArgumentOrContextualParameterInfo(node *ast.Node, position int, sourceFile *ast.SourceFile, checker *checker.Checker) *argumentListInfo {
	result := tryGetParameterInfo(node, sourceFile, checker)
	if result == nil {
		return getImmediatelyContainingArgumentInfo(node, position, sourceFile, checker)
	}
	return result
}

type argumentListInfo struct {
	isTypeParameterList bool
	invocation          *invocation
	argumentsRange      core.TextRange
	argumentIndex       int
	/** argumentCount is the *apparent* number of arguments. */
	argumentCount int
}

// Returns relevant information for the argument list and the current argument if we are
// in the argument of an invocation; returns undefined otherwise.
func getImmediatelyContainingArgumentInfo(node *ast.Node, position int, sourceFile *ast.SourceFile, c *checker.Checker) *argumentListInfo {
	parent := node.Parent
	if checker.IsCallOrNewExpression(parent) {
		// There are 3 cases to handle:
		//   1. The token introduces a list, and should begin a signature help session
		//   2. The token is either not associated with a list, or ends a list, so the session should end
		//   3. The token is buried inside a list, and should give signature help
		//
		// The following are examples of each:
		//
		//    Case 1:
		//          foo<#T, U>(#a, b)    -> The token introduces a list, and should begin a signature help session
		//    Case 2:
		//          fo#o<T, U>#(a, b)#   -> The token is either not associated with a list, or ends a list, so the session should end
		//    Case 3:
		//          foo<T#, U#>(a#, #b#) -> The token is buried inside a list, and should give signature help
		// Find out if 'node' is an argument, a type argument, or neither
		// const info = getArgumentOrParameterListInfo(node, position, sourceFile, checker);
		list, argumentIndex, argumentCount, argumentSpan := getArgumentOrParameterListInfo(node, sourceFile, c)
		isTypeParameterList := false
		if parent.TypeArgumentList() != nil {
			if parent.TypeArgumentList().Pos() == list.Pos() {
				isTypeParameterList = true
			}
		}
		return &argumentListInfo{
			isTypeParameterList: isTypeParameterList,
			invocation:          &invocation{kind: invocationKindCall, callInvocation: callInvocation{kind: invocationKindCall, node: parent}},
			argumentsRange:      argumentSpan,
			argumentIndex:       argumentIndex,
			argumentCount:       argumentCount,
		}
	} else if isNoSubstitutionTemplateLiteral(node) && isTaggedTemplateExpression(parent) {
		// Check if we're actually inside the template;
		// otherwise we'll fall out and return undefined.
		if isInsideTemplateLiteral(node, position, sourceFile) {
			return getArgumentListInfoForTemplate(parent.AsTaggedTemplateExpression(), 0, sourceFile)
		}
		return nil
	} else if isTemplateHead(node) && parent.Parent.Kind == ast.KindTaggedTemplateExpression {
		templateExpression := parent.AsTemplateExpression()
		tagExpression := templateExpression.Parent.AsTaggedTemplateExpression()

		argumentIndex := 1
		if isInsideTemplateLiteral(node, position, sourceFile) {
			argumentIndex = 0
		}
		return getArgumentListInfoForTemplate(tagExpression, argumentIndex, sourceFile)
	} else if ast.IsTemplateSpan(parent) && isTaggedTemplateExpression(parent.Parent.Parent) {
		templateSpan := parent
		tagExpression := parent.Parent.Parent

		// If we're just after a template tail, don't show signature help.
		if isTemplateTail(node) && !isInsideTemplateLiteral(node, position, sourceFile) {
			return nil
		}

		spanIndex := checker.IndexOfNode(templateSpan.Parent.AsTemplateExpression().TemplateSpans.Nodes, templateSpan)
		argumentIndex := getArgumentIndexForTemplatePiece(spanIndex, templateSpan, position, sourceFile)

		return getArgumentListInfoForTemplate(tagExpression.AsTaggedTemplateExpression(), argumentIndex, sourceFile)
	} else if checker.IsJsxOpeningLikeElement(parent) {
		// Provide a signature help for JSX opening element or JSX self-closing element.
		// This is not guarantee that JSX tag-name is resolved into stateless function component. (that is done in "getSignatureHelpItems")
		// i.e
		//      export function MainButton(props: ButtonProps, context: any): JSX.Element { ... }
		//      <MainButton /*signatureHelp*/
		attributeSpanStart := parent.AsJsxOpeningElement().Attributes.Loc.Pos()
		attributeSpanEnd := scanner.SkipTrivia(sourceFile.Text(), parent.AsJsxOpeningElement().Attributes.End())
		return &argumentListInfo{
			isTypeParameterList: false,
			invocation:          &invocation{kind: invocationKindCall, callInvocation: callInvocation{kind: invocationKindCall, node: parent}},
			argumentsRange:      core.NewTextRange(attributeSpanStart, attributeSpanEnd-attributeSpanStart),
			argumentIndex:       0,
			argumentCount:       1,
		}
	} else {
		typeArgInfo := getPossibleTypeArgumentsInfo(node, sourceFile)
		if typeArgInfo != nil {
			called := typeArgInfo.called
			nTypeArguments := typeArgInfo.nTypeArguments
			invoc := typeArgsInvocation{kind: invocationKindTypeArgs, called: called}
			argumentRange := core.NewTextRange(called.Loc.Pos(), node.End())
			return &argumentListInfo{
				isTypeParameterList: true,
				invocation: &invocation{
					kind:               invocationKindTypeArgs,
					typeArgsInvocation: invoc,
				},
				argumentsRange: argumentRange,
				argumentIndex:  nTypeArguments,
				argumentCount:  nTypeArguments + 1,
			}
		}
	}
	return nil
}

// spanIndex is either the index for a given template span.
// This does not give appropriate results for a NoSubstitutionTemplateLiteral
func getArgumentIndexForTemplatePiece(spanIndex int, node *ast.Node, position int, sourceFile *ast.SourceFile) int {
	// Because the TemplateStringsArray is the first argument, we have to offset each substitution expression by 1.
	// There are three cases we can encounter:
	//      1. We are precisely in the template literal (argIndex = 0).
	//      2. We are in or to the right of the substitution expression (argIndex = spanIndex + 1).
	//      3. We are directly to the right of the template literal, but because we look for the token on the left,
	//          not enough to put us in the substitution expression; we should consider ourselves part of
	//          the *next* span's expression by offsetting the index (argIndex = (spanIndex + 1) + 1).
	//
	// Example: f  `# abcd $#{#  1 + 1#  }# efghi ${ #"#hello"#  }  #  `
	//              ^       ^ ^       ^   ^          ^ ^      ^     ^
	// Case:        1       1 3       2   1          3 2      2     1
	//Debug.assert(position >= node.getStart(), "Assumed 'position' could not occur before node.");
	if ast.IsTemplateLiteralToken(node) {
		if isInsideTemplateLiteral(node, position, sourceFile) {
			return 0
		}
		return spanIndex + 2
	}
	return spanIndex + 1
}

func getAdjustedNode(node *ast.Node) *ast.Node {
	switch node.Kind {
	case ast.KindOpenParenToken, ast.KindCommaToken:
		return node
	default:
		ast.FindAncestor(node.Parent, func(n *ast.Node) bool {
			if ast.IsParameter(n) {
				return true
			} else if ast.IsBindingElement(n) || ast.IsObjectBindingPattern(n) || ast.IsArrayBindingPattern(n) {
				return false
			}
			return false
		})
	}
	return nil
}

type contextualSignatureLocationInfo struct {
	contextualType *checker.Type
	argumentIndex  int
	argumentCount  int
	argumentsSpan  core.TextRange
}

func getSpreadElementCount(node *ast.SpreadElement, c *checker.Checker) int {
	spreadType := c.GetTypeAtLocation(node.Expression)
	if checker.IsTupleType(spreadType) {
		tupleType := spreadType.Target().AsTupleType() //spreadType.AsTupleType() //.Target()
		if tupleType == nil {
			return 0
		}
		elementFlags := tupleType.ElementFlags()
		fixedLength := tupleType.FixedLength()
		if fixedLength == 0 {
			return 0
		}

		firstOptionalIndex := core.FindIndex(elementFlags, func(f checker.ElementFlags) bool {
			return (f&checker.ElementFlagsRequired == 0)
		})
		if firstOptionalIndex < 0 {
			return fixedLength
		}
		return firstOptionalIndex
	}
	return 0
}

func getArgumentIndex(node *ast.Node, arguments *ast.NodeList, sourceFile *ast.SourceFile, c *checker.Checker) int {
	return getArgumentIndexOrCount(getTokenFromNodeList(arguments, node.Parent, sourceFile), node, c)
}

func getArgumentCount(node *ast.Node, arguments *ast.NodeList, sourceFile *ast.SourceFile, c *checker.Checker) int {
	// The argument count for a list is normally the number of non-comma children it has.
	// For example, if you have "Foo(a,b)" then there will be three children of the arg
	// list 'a' '<comma>' 'b'. So, in this case the arg count will be 2. However, there
	// is a small subtlety. If you have "Foo(a,)", then the child list will just have
	// 'a' '<comma>'. So, in the case where the last child is a comma, we increase the

	// arg count by one to compensate.
	// if node.Kind == ast.KindCommaToken {
	// 	return len(arguments.Nodes) + 1
	// }
	// return len(arguments.Nodes) //tbd

	return getArgumentIndexOrCount(getTokenFromNodeList(arguments, node.Parent, sourceFile), nil, c)
}

func getArgumentIndexOrCount(arguments []*ast.Node, node *ast.Node, c *checker.Checker) int {
	argumentIndex := 0
	skipComma := false
	for _, arg := range arguments {
		if node != nil && arg == node {
			if !skipComma && arg.Kind == ast.KindCommaToken {
				argumentIndex++
			}
			return argumentIndex
		}
		if ast.IsSpreadElement(arg) {
			argumentIndex = argumentIndex + getSpreadElementCount(arg.AsSpreadElement(), c)
			skipComma = true
			continue
		}
		if arg.Kind != ast.KindCommaToken {
			argumentIndex++
			skipComma = true
			continue
		}
		if skipComma {
			skipComma = false
			continue
		}
		argumentIndex++
	}
	if node != nil {
		return argumentIndex
	}
	// The argument count for a list is normally the number of non-comma children it has.
	// For example, if you have "Foo(a,b)" then there will be three children of the arg
	// list 'a' '<comma>' 'b'. So, in this case the arg count will be 2. However, there
	// is a small subtlety. If you have "Foo(a,)", then the child list will just have
	// 'a' '<comma>'. So, in the case where the last child is a comma, we increase the
	// arg count by one to compensate.
	argumentCount := argumentIndex
	if len(arguments) > 0 && arguments[len(arguments)-1].Kind == ast.KindCommaToken {
		argumentCount = argumentIndex + 1
	}
	return argumentCount
}

func getArgumentOrParameterListInfo(node *ast.Node, sourceFile *ast.SourceFile, c *checker.Checker) (*ast.NodeList, int, int, core.TextRange) {
	arguments, argumentIndex := getArgumentOrParameterListAndIndex(node, sourceFile, c)
	argumentCount := getArgumentCount(node, arguments, sourceFile, c)
	argumentSpan := getApplicableSpanForArguments(arguments, node, sourceFile)
	return arguments, argumentIndex, argumentCount, argumentSpan
}

func getApplicableSpanForArguments(argumentList *ast.NodeList, node *ast.Node, sourceFile *ast.SourceFile) core.TextRange { //tbd
	// We use full start and skip trivia on the end because we want to include trivia on
	// both sides. For example,
	//
	//    foo(   /*comment */     a, b, c      /*comment*/     )
	//        |                                               |
	//
	// The applicable span is from the first bar to the second bar (inclusive,
	// but not including parentheses)
	if argumentList == nil && node != nil {
		// If the user has just opened a list, and there are no arguments.
		// For example, foo(    )
		//                  |  |
		return core.NewTextRange(node.End(), scanner.SkipTrivia(sourceFile.Text(), node.End()))
	}
	applicableSpanStart := argumentList.Pos()
	applicableSpanEnd := scanner.SkipTrivia(sourceFile.Text(), argumentList.End())
	return core.NewTextRange(applicableSpanStart, applicableSpanEnd)
}

func getArgumentOrParameterListAndIndex(node *ast.Node, sourceFile *ast.SourceFile, c *checker.Checker) (*ast.NodeList, int) {
	if node.Kind == ast.KindLessThanToken || node.Kind == ast.KindOpenParenToken {
		// Find the list that starts right *after* the < or ( token.
		// If the user has just opened a list, consider this item 0.
		list := getChildListThatStartsWithOpenerToken(node.Parent, node, sourceFile)
		return list, 0
	} else {
		// findListItemInfo can return undefined if we are not in parent's argument list
		// or type argument list. This includes cases where the cursor is:
		//   - To the right of the closing parenthesis, non-substitution template, or template tail.
		//   - Between the type arguments and the arguments (greater than token)
		//   - On the target of the call (parent.func)
		//   - On the 'new' keyword in a 'new' expression
		// list := findContainingList(node)
		// return list && { list, argumentIndex: getArgumentIndex(checker, list, node) };
		var arguments *ast.NodeList
		if ast.IsCallExpression(node.Parent) {
			arguments = node.Parent.AsCallExpression().Arguments
		} else if ast.IsNewExpression(node.Parent) {
			arguments = node.Parent.AsNewExpression().Arguments
		}
		if arguments == nil {
			return nil, 0
		}
		// Find the index of the argument that contains the node.
		argumentIndex := getArgumentIndex(node, arguments, sourceFile, c)
		return arguments, argumentIndex
	}
}

func getChildListThatStartsWithOpenerToken(parent *ast.Node, openerToken *ast.Node, sourceFile *ast.SourceFile) *ast.NodeList { //tbd check for the implementation
	if ast.IsCallExpression(parent) {
		return parent.AsCallExpression().Arguments
	} else if ast.IsNewExpression(parent) {
		return parent.AsNewExpression().Arguments
	}
	return nil
}

func tryGetParameterInfo(startingToken *ast.Node, sourceFile *ast.SourceFile, c *checker.Checker) *argumentListInfo {
	node := getAdjustedNode(startingToken)
	if node == nil {
		return nil
	}
	info := getContextualSignatureLocationInfo(node, sourceFile, c)
	if info == nil {
		return nil
	}

	// for optional function condition
	nonNullableContextualType := c.GetNonNullableType(info.contextualType)
	if nonNullableContextualType == nil {
		return nil
	}

	symbol := nonNullableContextualType.Symbol()
	if symbol == nil {
		return nil
	}

	signatures := c.GetSignaturesOfType(nonNullableContextualType, checker.SignatureKindCall)
	if signatures == nil || signatures[len(signatures)-1] == nil {
		return nil
	}
	signature := signatures[len(signatures)-1]

	contextualInvocation := contextualInvocation{
		kind:      invocationKinContextual,
		signature: signature,
		node:      startingToken,
		symbol:    chooseBetterSymbol(symbol),
	}
	return &argumentListInfo{
		isTypeParameterList: false,
		invocation:          &invocation{kind: invocationKinContextual, contextualInvocation: contextualInvocation},
		argumentsRange:      info.argumentsSpan,
		argumentIndex:       info.argumentIndex,
		argumentCount:       info.argumentCount,
	}
}

func chooseBetterSymbol(s *ast.Symbol) *ast.Symbol {
	if s.Name == ast.InternalSymbolNameType {
		for _, d := range s.Declarations {
			if ast.IsFunctionTypeNode(d) && checker.CanHaveSymbol(d.Parent) {
				return d.Parent.Symbol()
			}
		}
	}
	return s
}
func getContextualSignatureLocationInfo(node *ast.Node, sourceFile *ast.SourceFile, c *checker.Checker) *contextualSignatureLocationInfo {
	parent := node.Parent
	switch parent.Kind {
	case ast.KindParenthesizedExpression, ast.KindMethodDeclaration, ast.KindFunctionExpression, ast.KindArrowFunction:
		_, argumentIndex, argumentCount, argumentSpan := getArgumentOrParameterListInfo(node, sourceFile, c)

		var contextualType *checker.Type
		if ast.IsMethodDeclaration(parent) {
			contextualType = c.GetContextualTypeForObjectLiteralElement(parent, checker.ContextFlagsNone)
		} else {
			contextualType = c.GetContextualType(parent, checker.ContextFlagsNone)
		}
		if contextualType != nil {
			return &contextualSignatureLocationInfo{
				contextualType: contextualType,
				argumentIndex:  argumentIndex,
				argumentCount:  argumentCount,
				argumentsSpan:  argumentSpan,
			}
		}
		return nil
	case ast.KindBinaryExpression:
		highestBinary := getHighestBinary(parent.AsBinaryExpression())
		contextualType := c.GetContextualType(highestBinary.AsNode(), checker.ContextFlagsNone)
		var argumentIndex int
		if node.Kind == ast.KindOpenParenToken {
			argumentIndex = 0
		} else {
			argumentIndex = countBinaryExpressionParameters(parent.AsBinaryExpression()) - 1
			argumentCount := countBinaryExpressionParameters(highestBinary)
			if contextualType != nil {
				return &contextualSignatureLocationInfo{
					contextualType: contextualType,
					argumentIndex:  argumentIndex,
					argumentCount:  argumentCount,
					argumentsSpan:  core.NewTextRange(parent.Pos(), parent.End()), //node.getStart(sourceFile), (endNode || node).getEnd() used for textSpan
				}
			}
			return nil
		}
	}
	return nil
}

func getHighestBinary(b *ast.BinaryExpression) *ast.BinaryExpression {
	if ast.IsBinaryExpression(b.Parent) {
		return getHighestBinary(b.Parent.AsBinaryExpression())
	}
	return b
}

func countBinaryExpressionParameters(b *ast.BinaryExpression) int {
	if ast.IsBinaryExpression(b.Left) {
		return countBinaryExpressionParameters(b.Left.AsBinaryExpression()) + 1
	}
	return 2
}

func getChildren(node *ast.Node, sourceFile *ast.SourceFile) []*ast.Node {
	if node == nil {
		return nil
	}
	var children []*ast.Node
	current := node
	left := node.Pos()
	scanner := scanner.GetScannerForSourceFile(sourceFile, left)
	for left < current.End() {
		token := scanner.Token()
		tokenFullStart := scanner.TokenFullStart()
		tokenEnd := scanner.TokenEnd()
		children = append(children, sourceFile.GetOrCreateToken(token, tokenFullStart, tokenEnd, current))
		left = tokenEnd
		scanner.Scan()
	}
	return children
}

func getTokenFromNodeList(nodeList *ast.NodeList, nodeListParent *ast.Node, sourceFile *ast.SourceFile) []*ast.Node {
	if nodeList == nil || nodeListParent == nil {
		return nil
	}
	left := nodeList.Pos()
	nodeListIndex := 0
	var tokens []*ast.Node
	for left < nodeList.End() {
		if len(nodeList.Nodes) > nodeListIndex && left == nodeList.Nodes[nodeListIndex].Pos() {
			tokens = append(tokens, nodeList.Nodes[nodeListIndex])
			left = nodeList.Nodes[nodeListIndex].End()
			nodeListIndex++
		} else {
			scanner := scanner.GetScannerForSourceFile(sourceFile, left)
			token := scanner.Token()
			tokenFullStart := scanner.TokenFullStart()
			tokenEnd := scanner.TokenEnd()
			tokens = append(tokens, sourceFile.GetOrCreateToken(token, tokenFullStart, tokenEnd, nodeListParent))
			left = tokenEnd
		}
	}
	return tokens
}

func containsNode(nodes []*ast.Node, node *ast.Node) bool {
	for i := 0; i < len(nodes); i++ {
		if nodes[i] == node {
			return true
		}
	}
	return false
}

func getArgumentListInfoForTemplate(tagExpression *ast.TaggedTemplateExpression, argumentIndex int, sourceFile *ast.SourceFile) *argumentListInfo {
	// argumentCount is either 1 or (numSpans + 1) to account for the template strings array argument.
	argumentCount := 1
	if !isNoSubstitutionTemplateLiteral(tagExpression.Template) {
		argumentCount = len(tagExpression.Template.AsTemplateExpression().TemplateSpans.Nodes) + 1
	}
	// if (argumentIndex !== 0) {
	//     Debug.assertLessThan(argumentIndex, argumentCount);
	// }
	return &argumentListInfo{
		isTypeParameterList: false,
		invocation:          &invocation{kind: invocationKindCall, callInvocation: callInvocation{kind: invocationKindCall, node: tagExpression.AsNode()}},
		argumentIndex:       argumentIndex,
		argumentCount:       argumentCount,
		argumentsRange:      getApplicableRangeForTaggedTemplate(tagExpression, sourceFile),
	}
}

func getApplicableRangeForTaggedTemplate(taggedTemplate *ast.TaggedTemplateExpression, sourceFile *ast.SourceFile) core.TextRange {
	template := taggedTemplate.Template
	applicableSpanStart := scanner.GetTokenPosOfNode(template, sourceFile, false)
	applicableSpanEnd := template.End()

	// We need to adjust the end position for the case where the template does not have a tail.
	// Otherwise, we will not show signature help past the expression.
	// For example,
	//
	//      ` ${ 1 + 1 foo(10)
	//       |       |
	// This is because a Missing node has no width. However, what we actually want is to include trivia
	// leading up to the next token in case the user is about to type in a TemplateMiddle or TemplateTail.
	if template.Kind == ast.KindTemplateExpression {
		templateSpans := template.AsTemplateExpression().TemplateSpans
		lastSpan := templateSpans.Nodes[len(templateSpans.Nodes)-1]
		if lastSpan.AsTemplateSpan().Literal.GetFullWidth() == 0 {
			applicableSpanEnd = scanner.SkipTrivia(sourceFile.Text(), applicableSpanEnd)
		}
	}

	return core.NewTextRange(applicableSpanStart, applicableSpanEnd-applicableSpanStart)
}

type possibleTypeArgumentInfo struct {
	called         *ast.Identifier
	nTypeArguments int
}

// Get info for an expression like `f <` that may be the start of type arguments.
func getPossibleTypeArgumentsInfo(tokenIn *ast.Node, sourceFile *ast.SourceFile) *possibleTypeArgumentInfo {
	// This is a rare case, but one that saves on a _lot_ of work if true - if the source file has _no_ `<` character,
	// then there obviously can't be any type arguments - no expensive brace-matching backwards scanning required
	if strings.LastIndex(sourceFile.Text(), "<") == -1 { // (sourceFile.text.lastIndexOf("<", tokenIn ? tokenIn.pos : sourceFile.text.length) === -1)
		return nil
	}
	var token *ast.Node
	// This function determines if the node could be type argument position
	// Since during editing, when type argument list is not complete,
	// the tree could be of any shape depending on the tokens parsed before current node,
	// scanning of the previous identifier followed by "<" before current node would give us better result
	// Note that we also balance out the already provided type arguments, arrays, object literals while doing so
	remainingLessThanTokens := 0
	nTypeArguments := 0
	for token != nil {
		switch token.Kind {
		case ast.KindLessThanToken:
			// Found the beginning of the generic argument expression
			token := astnav.FindPrecedingToken(sourceFile, token.Loc.Pos())
			if token != nil && token.Kind == ast.KindQuestionDotToken {
				token = astnav.FindPrecedingToken(sourceFile, token.Loc.Pos())
			}
			if token == nil && ast.IsIdentifier(token) {
				return nil
			}
			// if (!remainingLessThanTokens) {
			// 	return isDeclarationName(token) ? undefined : { called: token, nTypeArguments };
			// }
			remainingLessThanTokens--
			break
		case ast.KindGreaterThanGreaterThanGreaterThanToken:
			remainingLessThanTokens = +3
			break
		case ast.KindGreaterThanGreaterThanToken:
			remainingLessThanTokens = +2
			break
		case ast.KindGreaterThanToken:
			remainingLessThanTokens++
			break
		case ast.KindCloseBraceToken:
			// This can be object type, skip until we find the matching open brace token
			// Skip until the matching open brace token
			token := findPrecedingMatchingToken(token, ast.KindOpenBraceToken, sourceFile)
			if token == nil {
				return nil
			}
			break
		case ast.KindCloseParenToken:
			// This can be object type, skip until we find the matching open brace token
			// Skip until the matching open brace token
			token := findPrecedingMatchingToken(token, ast.KindOpenParenToken, sourceFile)
			if token == nil {
				return nil
			}
			break
		case ast.KindCloseBracketToken:
			// This can be object type, skip until we find the matching open brace token
			// Skip until the matching open brace token
			token := findPrecedingMatchingToken(token, ast.KindOpenBracketToken, sourceFile)
			if token == nil {
				return nil
			}
			break

		// Valid tokens in a type name. Skip.
		case ast.KindCommaToken:
			nTypeArguments++
			break
		case ast.KindEqualsGreaterThanToken:
		case ast.KindIdentifier:
		case ast.KindStringLiteral:
		case ast.KindNumericLiteral:
		case ast.KindBigIntLiteral:
		case ast.KindTrueKeyword:
		case ast.KindFalseKeyword:
		case ast.KindTypeOfKeyword:
		case ast.KindExtendsKeyword:
		case ast.KindKeyOfKeyword:
		case ast.KindDotToken:
		case ast.KindBarToken:
		case ast.KindQuestionToken:
		case ast.KindColonToken:
			break
		default:
			if ast.IsTypeNode(token) {
				break
			}

			// Invalid token in type
			return nil
		}
		token = astnav.FindPrecedingToken(sourceFile, token.Loc.Pos())
	}
	return nil
}
