package ls

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/core"
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

// Represents a single signature to show in signature help.
// The id is used for subsequent calls into the language service to ask questions about the
// signature help item in the context of any documents that have been updated.  i.e. after
// an edit has happened, while signature help is still active, the host can ask important
// questions like 'what parameter is the user currently contained within?'.
// type signatureHelpItem struct {
// 	isVariadic            bool
// 	prefixDisplayParts    []symbolDisplayPart
// 	suffixDisplayParts    []symbolDisplayPart
// 	separatorDisplayParts []symbolDisplayPart
// 	parameters            []SignatureHelpParameter
// 	documentation         []symbolDisplayPart
// 	tags                  jSDocTagInfo
// }

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
	startingToken := astnav.FindTokenOnLeftOfPosition(sourceFile, position)
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
	argumentInfo := getContainingArgumentInfo(startingToken, sourceFile, program.GetTypeChecker(), isManuallyInvoked)
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
		return createSignatureHelpItems(candidateInfo.candidateInfo.candidates, candidateInfo.candidateInfo.resolvedSignature, *argumentInfo, sourceFile, typeChecker, onlyUseSyntacticOwners)
	}
	return createTypeHelpItems(candidateInfo.typeInfo.symbol, argumentInfo, typeChecker)
}

func createTypeHelpItems(symbol *ast.Symbol, argumentInfo *argumentListInfo, c *checker.Checker) *SignatureHelpItems {
	typeParameters := c.GetLocalTypeParametersOfClassOrInterfaceOrTypeAlias(symbol)
	if typeParameters == nil {
		return nil
	}
	items := getTypeHelpItems(symbol, typeParameters)
	return &SignatureHelpItems{
		Signatures:      []signatureHelpItem{items},
		ActiveSignature: 0,
		ActiveParameter: argumentInfo.argumentIndex,
	}
}

func getTypeHelpItems(symbol *ast.Symbol, typeParameter []*checker.Type) signatureHelpItem {
	parameters := createSignatureHelpParameterForTypeParameter(typeParameter)
	// const documentation = symbol.getDocumentationComment(checker);
	// const tags = symbol.getJsDocTags(checker);
	return signatureHelpItem{
		Label:         symbol.Name,
		Documentation: nil,
		Parameters:    &parameters,
		IsVariadic:    false,
	}
}

func createSignatureHelpItems(candidates *[]*checker.Signature, resolvedSignature *checker.Signature, argumentInfo argumentListInfo, sourceFile *ast.SourceFile, c *checker.Checker, useFullPrefic bool) *SignatureHelpItems {
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
	// var parameters []signatureHelpParameterInformation
	var items []signatureHelpItem
	for _, candidate := range *candidates {
		items = append(items, getSignatureHelpItem(candidate, argumentInfo.isTypeParameterList, callTargetSymbol))
	}

	selectedItemIndex := 0
	//itemSeen := 0 this needs items to be a 2d array
	for i := 0; i < len(items); i++ {
		if (*candidates)[i] == resolvedSignature {
			selectedItemIndex = i
			// selectedItemIndex = itemsSeen;
			// if (item.length > 1) {
			//     // check to see if any items in the list better match than the first one, as the checker isn't filtering the nested lists
			//     // (those come from tuple parameter expansion)
			//     let count = 0;
			//     for (const i of item) {
			//         if (i.isVariadic || i.parameters.length >= argumentCount) {
			//             selectedItemIndex = itemsSeen + count;
			//             break;
			//         }
			//         count++;
			//     }
			// }
		}
	}
	help := &SignatureHelpItems{
		Signatures:      items,
		ActiveSignature: selectedItemIndex,
		ActiveParameter: argumentInfo.argumentIndex,
	}
	selected := items[selectedItemIndex]
	if selected.IsVariadic {
		firstRest := core.FindIndex(*selected.Parameters, func(p signatureHelpParameters) bool {
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

func getSignatureHelpItem(candidate *checker.Signature, isTypeParameterList bool, callTargetSymbol *ast.Symbol) signatureHelpItem {
	var parameters []signatureHelpParameters
	if isTypeParameterList {
		var typeParameters []*checker.Type
		if candidate.Target() != nil {
			typeParameters = candidate.Target().TypeParameters()
		} else {
			typeParameters = candidate.TypeParameters()
		}
		parameters = createSignatureHelpParameterForTypeParameter(typeParameters)
	} else {
		parameters = createSignatureHelpParameterInformation(candidate.Parameters())
	}
	return signatureHelpItem{
		Label:           callTargetSymbol.Name,
		Documentation:   nil,
		Parameters:      &parameters,
		ActiveParameter: -1, //how to set this?
	}
}

// func itemInfoForParameter(candidateSignature *checker.Signature, c *checker.Checker, enclosingDeclaratipn *ast.Node, sourceFile *ast.SourceFile) signatureHelpItem {

// }
func createSignatureHelpParameterInformation(symbol []*ast.Symbol) []signatureHelpParameters {
	var parameterInfo []signatureHelpParameters
	for _, s := range symbol {
		parameterInfo = append(parameterInfo, signatureHelpParameters{
			Label:         s.Name,
			Documentation: nil,
		})
	}
	return parameterInfo
}

func createSignatureHelpParameterForTypeParameter(typeParameters []*checker.Type) []signatureHelpParameters {
	var parameterInfo []signatureHelpParameters
	for _, s := range typeParameters {
		parameterInfo = append(parameterInfo, signatureHelpParameters{
			Label:         s.Symbol().Name,
			Documentation: nil,
		})
	}
	return parameterInfo
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
	Parameters *[]signatureHelpParameters

	// The index of the active parameter.
	ActiveParameter int

	// Needed only here, not in lsp
	IsVariadic bool
}

// Represents a parameter of a callable-signature. A parameter can
// have a label and a doc-comment.
type signatureHelpParameters struct {
	Label         string
	Documentation *string
	isRest        bool
}

// const defaultMaximumTruncationLength = 160;
// func symbolToDisplayParts(symbol *ast.Symbol,c *checker.Checker) {
// 	sym := c.SymbolToString(symbol)
// 	absoluteMaximumLength := defaultMaximumTruncationLength * 10

// //come back
// }
func getEnclosingDeclarationFromInvocation(invocation invocation) *ast.Node {
	if invocation.kind == invocationKindCall {
		return invocation.callInvocation.node
	} else if invocation.kind == invocationKindTypeArgs {
		return invocation.typeArgsInvocation.called.AsNode()
	} else {
		return invocation.contextualInvocation.node
	}
}

// func createJSSignatureHelpItems(argumentInfo *argumentListInfo, p *compiler.Program, c *checker.Checker) *signatureHelpItems {
// 	if argumentInfo.invocation.kind == invocationKinContextual {
// 		return nil
// 	}
// 	// See if we can find some symbol with the call expression name that has call signatures.
// 	expression := getExpressionFromInvocation(argumentInfo)
// 	name := ""
// 	if ast.IsPropertyAccessExpression(expression) {
// 		name = expression.Name().Text()
// 	}
// 	if name != "" {
// 		sourceFiles := p.GetSourceFiles()
// 		for _, sourceFile := range sourceFiles {

// 		}
// 	}
// 	return nil

// }

func getExpressionFromInvocation(argumentInfo argumentListInfo) *ast.Node {
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
		resolvedSignature := c.GetResolvedSignatureForSignatureHelp(info.invocation.callInvocation.node, candidates)
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
	case ast.KindCommaToken:
		return containsNode(invocationChildren, startingToken) //tbd does this need to check only for argument list?
		// const containingList = findContainingList(startingToken);
		// return !!containingList && contains(invocationChildren, containingList);
	case ast.KindLessThanToken:
		return containsPrecedingToken(startingToken, sourceFile, node.AsCallExpression().Expression)
	default:
		return false
	}
	return false
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
func getContainingArgumentInfo(node *ast.Node, sourceFile *ast.SourceFile, checker *checker.Checker, isManuallyInvoked bool) *argumentListInfo {
	for n := node; !ast.IsSourceFile(n) && (isManuallyInvoked || !ast.IsBlock(n)); n = n.Parent {
		// If the node is not a subspan of its parent, this is a big problem.
		// There have been crashes that might be caused by this violation.
		//Debug.assert(rangeContainsRange(n.parent, n), "Not a subspan", () => `Child: ${Debug.formatSyntaxKind(n.kind)}, parent: ${Debug.formatSyntaxKind(n.parent.kind)}`);
		argumentInfo := getImmediatelyContainingArgumentOrContextualParameterInfo(n, sourceFile, checker)
		if argumentInfo != nil {
			return argumentInfo
		}
	}
	return nil
}

func getImmediatelyContainingArgumentOrContextualParameterInfo(node *ast.Node, sourceFile *ast.SourceFile, checker *checker.Checker) *argumentListInfo {
	result := tryGetParameterInfo(node, sourceFile, checker)
	if result == nil {
		return getImmediatelyContainingArgumentInfo(node, sourceFile, checker)
	}
	return result
}

type argumentListInfo struct {
	isTypeParameterList bool
	invocation          invocation
	argumentsRange      core.TextRange
	argumentIndex       int
	/** argumentCount is the *apparent* number of arguments. */
	argumentCount int
}

// Returns relevant information for the argument list and the current argument if we are
// in the argument of an invocation; returns undefined otherwise.
func getImmediatelyContainingArgumentInfo(node *ast.Node, sourceFile *ast.SourceFile, checker *checker.Checker) *argumentListInfo {
	parent := node.Parent
	if ast.IsCallExpression(parent) || ast.IsNewExpression(parent) {
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
		list, argumentIndex, argumentCount, argumentSpan := getArgumentOrParameterListInfo(node, sourceFile, checker)
		if list == nil {
			return nil
		}
		isTypeParameterList := parent.TypeArguments() != nil && parent.TypeArguments()[0].Pos() == list.Loc.Pos()
		return &argumentListInfo{
			isTypeParameterList: isTypeParameterList,
			invocation:          invocation{kind: invocationKindCall, callInvocation: callInvocation{kind: invocationKindCall, node: parent}},
			argumentsRange:      argumentSpan,
			argumentIndex:       argumentIndex,
			argumentCount:       argumentCount,
		}
	}
	return nil

}
func getAdjustedNode(node *ast.Node) *ast.Node {
	switch node.Kind {
	case ast.KindOpenParenToken:
	case ast.KindCommaToken:
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

func getArgumentOrParameterListInfo(node *ast.Node, sourceFile *ast.SourceFile, checker *checker.Checker) (*ast.NodeList, int, int, core.TextRange) {
	list, argumentIndex := getArgumentOrParameterListAndIndex(node, sourceFile, checker)
	if list == nil {
		return nil, 0, 0, core.TextRange{}
	}
	argumentCount := len(list.Nodes)
	argumentSpan := getApplicableSpanForArguments(list, sourceFile)
	return list, argumentIndex, argumentCount, argumentSpan
}

func getApplicableSpanForArguments(argumentList *ast.NodeList, sourceFile *ast.SourceFile) core.TextRange { //tbd
	// We use full start and skip trivia on the end because we want to include trivia on
	// both sides. For example,
	//
	//    foo(   /*comment */     a, b, c      /*comment*/     )
	//        |                                               |
	//
	// The applicable span is from the first bar to the second bar (inclusive,
	// but not including parentheses)
	applicableSpanStart := argumentList.Loc.Pos() //fullstart - .pos but getStart is getPosOfToken
	applicableSpanEnd := scanner.SkipTrivia(sourceFile.Text, argumentList.Nodes[len(argumentList.Nodes)-1].End())
	return core.NewTextRange(applicableSpanStart, applicableSpanEnd)
}
func getArgumentOrParameterListAndIndex(node *ast.Node, sourceFile *ast.SourceFile, checker *checker.Checker) (*ast.NodeList, int) {
	if node.Kind == ast.KindLessThanToken || node.Kind == ast.KindOpenParenToken {
		// Find the list that starts right *after* the < or ( token.
		// If the user has just opened a list, consider this item 0.
		return getChildListThatStartsWithOpenerToken(node.Parent, node, sourceFile), 0
	} else { //tbd
		// findListItemInfo can return undefined if we are not in parent's argument list
		// or type argument list. This includes cases where the cursor is:
		//   - To the right of the closing parenthesis, non-substitution template, or template tail.
		//   - Between the type arguments and the arguments (greater than token)
		//   - On the target of the call (parent.func)
		//   - On the 'new' keyword in a 'new' expression
		// list := findContainingList(node)
		// return list && { list, argumentIndex: getArgumentIndex(checker, list, node) };
	}
	return nil, 0
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
	signatures := nonNullableContextualType.AsStructuredType().CallSignatures()
	signature := signatures[len(signatures)-1]
	symbol := nonNullableContextualType.Symbol()
	if symbol == nil {
		return nil
	}

	contextualInvocation := contextualInvocation{
		kind:      invocationKinContextual,
		signature: signature,
		node:      startingToken,
		symbol:    nonNullableContextualType.Symbol(), //check
	}
	return &argumentListInfo{
		isTypeParameterList: false,
		invocation:          invocation{kind: invocationKinContextual, contextualInvocation: contextualInvocation},
		argumentsRange:      info.argumentsSpan,
		argumentIndex:       info.argumentIndex,
		argumentCount:       info.argumentCount,
	}
}

func getContextualSignatureLocationInfo(node *ast.Node, sourceFile *ast.SourceFile, c *checker.Checker) *contextualSignatureLocationInfo {
	parent := node.Parent
	switch parent.Kind {
	case ast.KindParenthesizedExpression:
	case ast.KindMethodDeclaration:
	case ast.KindFunctionDeclaration:
	case ast.KindArrowFunction:
		list, argumentIndex, argumentCount, argumentSpan := getArgumentOrParameterListInfo(node, sourceFile, c)
		if list == nil {
			return nil
		}
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

func containsNode(nodes []*ast.Node, node *ast.Node) bool {
	for _, n := range nodes {
		if n == node {
			return true
		}
	}
	return false
}
