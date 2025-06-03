package format

import (
	"slices"

	"github.com/microsoft/typescript-go/internal/ast"
)

func getAllRules() []ruleSpec {
	allTokens := make([]ast.Kind, 0, ast.KindLastToken-ast.KindFirstToken+1)
	for token := ast.KindFirstToken; token <= ast.KindLastToken; token++ {
		allTokens = append(allTokens, token)
	}

	anyTokenExcept := func(tokens ...ast.Kind) tokenRange {
		newTokens := make([]ast.Kind, 0, ast.KindLastToken-ast.KindFirstToken+1)
		for token := ast.KindFirstToken; token <= ast.KindLastToken; token++ {
			if slices.Contains(tokens, token) {
				continue
			}
			newTokens = append(newTokens, token)
		}
		return tokenRange{
			isSpecific: false,
			tokens:     newTokens,
		}
	}

	anyToken := tokenRange{
		isSpecific: false,
		tokens:     allTokens,
	}

	anyTokenIncludingMultilineComments := tokenRangeFromEx(allTokens, ast.KindMultiLineCommentTrivia)
	anyTokenIncludingEOF := tokenRangeFromEx(allTokens, ast.KindEndOfFile)
	keywords := tokenRangeFromRange(ast.KindFirstKeyword, ast.KindLastKeyword)
	binaryOperators := tokenRangeFromRange(ast.KindFirstBinaryOperator, ast.KindLastBinaryOperator)
	binaryKeywordOperators := []ast.Kind{
		ast.KindInKeyword,
		ast.KindInstanceOfKeyword,
		ast.KindOfKeyword,
		ast.KindAsKeyword,
		ast.KindIsKeyword,
		ast.KindSatisfiesKeyword,
	}
	unaryPrefixOperators := []ast.Kind{ast.KindPlusPlusToken, ast.KindMinusToken, ast.KindTildeToken, ast.KindExclamationToken}
	unaryPrefixExpressions := []ast.Kind{
		ast.KindNumericLiteral,
		ast.KindBigIntLiteral,
		ast.KindIdentifier,
		ast.KindOpenParenToken,
		ast.KindOpenBracketToken,
		ast.KindOpenBraceToken,
		ast.KindThisKeyword,
		ast.KindNewKeyword,
	}
	unaryPreincrementExpressions := []ast.Kind{ast.KindIdentifier, ast.KindOpenParenToken, ast.KindThisKeyword, ast.KindNewKeyword}
	unaryPostincrementExpressions := []ast.Kind{ast.KindIdentifier, ast.KindCloseParenToken, ast.KindCloseBracketToken, ast.KindNewKeyword}
	unaryPredecrementExpressions := []ast.Kind{ast.KindIdentifier, ast.KindOpenParenToken, ast.KindThisKeyword, ast.KindNewKeyword}
	unaryPostdecrementExpressions := []ast.Kind{ast.KindIdentifier, ast.KindCloseParenToken, ast.KindCloseBracketToken, ast.KindNewKeyword}
	comments := []ast.Kind{ast.KindSingleLineCommentTrivia, ast.KindMultiLineCommentTrivia}
	typeKeywords := []ast.Kind{
		ast.KindAnyKeyword,
		ast.KindAssertsKeyword,
		ast.KindBigIntKeyword,
		ast.KindBooleanKeyword,
		ast.KindFalseKeyword,
		ast.KindInferKeyword,
		ast.KindKeyOfKeyword,
		ast.KindNeverKeyword,
		ast.KindNullKeyword,
		ast.KindNumberKeyword,
		ast.KindObjectKeyword,
		ast.KindReadonlyKeyword,
		ast.KindStringKeyword,
		ast.KindSymbolKeyword,
		ast.KindTypeOfKeyword,
		ast.KindTrueKeyword,
		ast.KindVoidKeyword,
		ast.KindUndefinedKeyword,
		ast.KindUniqueKeyword,
		ast.KindUnknownKeyword,
	}
	typeNames := append([]ast.Kind{ast.KindIdentifier}, typeKeywords...)

	// Place a space before open brace in a function declaration
	// TypeScript: Function can have return types, which can be made of tons of different token kinds
	functionOpenBraceLeftTokenRange := anyTokenIncludingMultilineComments

	// Place a space before open brace in a TypeScript declaration that has braces as children (class, module, enum, etc)
	typeScriptOpenBraceLeftTokenRange := tokenRangeFrom(ast.KindIdentifier, ast.KindGreaterThanToken, ast.KindMultiLineCommentTrivia, ast.KindClassKeyword, ast.KindExportKeyword, ast.KindImportKeyword)

	// Place a space before open brace in a control flow construct
	controlOpenBraceLeftTokenRange := tokenRangeFrom(ast.KindCloseParenToken, ast.KindMultiLineCommentTrivia, ast.KindDoKeyword, ast.KindTryKeyword, ast.KindFinallyKeyword, ast.KindElseKeyword, ast.KindCatchKeyword)

	// These rules are higher in priority than user-configurable
	highPriorityCommonRules := []ruleSpec{
		// Leave comments alone
		rule("IgnoreBeforeComment", anyToken, comments, anyContext, RuleActionStopProcessingSpaceActions),
		rule("IgnoreAfterLineComment", ast.KindSingleLineCommentTrivia, anyToken, anyContext, RuleActionStopProcessingSpaceActions),

		rule("NotSpaceBeforeColon", anyToken, ast.KindColonToken, []ContextPredicate{isNonJsxSameLineTokenContext, isNotBinaryOpContext, isNotTypeAnnotationContext}, RuleActionDeleteSpace),
		rule("SpaceAfterColon", ast.KindColonToken, anyToken, []ContextPredicate{isNonJsxSameLineTokenContext, isNotBinaryOpContext, isNextTokenParentNotJsxNamespacedName}, RuleActionInsertSpace),
		rule("NoSpaceBeforeQuestionMark", anyToken, ast.KindQuestionToken, []ContextPredicate{isNonJsxSameLineTokenContext, isNotBinaryOpContext, isNotTypeAnnotationContext}, RuleActionDeleteSpace),
		// insert space after '?' only when it is used in conditional operator
		rule("SpaceAfterQuestionMarkInConditionalOperator", ast.KindQuestionToken, anyToken, []ContextPredicate{isNonJsxSameLineTokenContext, isConditionalOperatorContext}, RuleActionInsertSpace),

		// in other cases there should be no space between '?' and next token
		rule("NoSpaceAfterQuestionMark", ast.KindQuestionToken, anyToken, []ContextPredicate{isNonJsxSameLineTokenContext, isNonOptionalPropertyContext}, RuleActionDeleteSpace),

		rule("NoSpaceBeforeDot", anyToken, []ast.Kind{ast.KindDotToken, ast.KindQuestionDotToken}, []ContextPredicate{isNonJsxSameLineTokenContext, isNotPropertyAccessOnIntegerLiteral}, RuleActionDeleteSpace),
		rule("NoSpaceAfterDot", []ast.Kind{ast.KindDotToken, ast.KindQuestionDotToken}, anyToken, []ContextPredicate{isNonJsxSameLineTokenContext}, RuleActionDeleteSpace),

		rule("NoSpaceBetweenImportParenInImportType", ast.KindImportKeyword, ast.KindOpenParenToken, []ContextPredicate{isNonJsxSameLineTokenContext, isImportTypeContext}, RuleActionDeleteSpace),

		// Special handling of unary operators.
		// Prefix operators generally shouldn't have a space between
		// them and their target unary expression.
		rule("NoSpaceAfterUnaryPrefixOperator", unaryPrefixOperators, unaryPrefixExpressions, []ContextPredicate{isNonJsxSameLineTokenContext, isNotBinaryOpContext}, RuleActionDeleteSpace),
		rule("NoSpaceAfterUnaryPreincrementOperator", ast.KindPlusPlusToken, unaryPreincrementExpressions, []ContextPredicate{isNonJsxSameLineTokenContext}, RuleActionDeleteSpace),
		rule("NoSpaceAfterUnaryPredecrementOperator", ast.KindMinusMinusToken, unaryPredecrementExpressions, []ContextPredicate{isNonJsxSameLineTokenContext}, RuleActionDeleteSpace),
		rule("NoSpaceBeforeUnaryPostincrementOperator", unaryPostincrementExpressions, ast.KindPlusPlusToken, []ContextPredicate{isNonJsxSameLineTokenContext, isNotStatementConditionContext}, RuleActionDeleteSpace),
		rule("NoSpaceBeforeUnaryPostdecrementOperator", unaryPostdecrementExpressions, ast.KindMinusMinusToken, []ContextPredicate{isNonJsxSameLineTokenContext, isNotStatementConditionContext}, RuleActionDeleteSpace),

		// More unary operator special-casing.
		// DevDiv 181814: Be careful when removing leading whitespace
		// around unary operators.  Examples:
		//      1 - -2  --X--> 1--2
		//      a + ++b --X--> a+++b
		rule("SpaceAfterPostincrementWhenFollowedByAdd", ast.KindPlusPlusToken, ast.KindPlusToken, []ContextPredicate{isNonJsxSameLineTokenContext, isBinaryOpContext}, RuleActionInsertSpace),
		rule("SpaceAfterAddWhenFollowedByUnaryPlus", ast.KindPlusToken, ast.KindPlusToken, []ContextPredicate{isNonJsxSameLineTokenContext, isBinaryOpContext}, RuleActionInsertSpace),
		rule("SpaceAfterAddWhenFollowedByPreincrement", ast.KindPlusToken, ast.KindPlusPlusToken, []ContextPredicate{isNonJsxSameLineTokenContext, isBinaryOpContext}, RuleActionInsertSpace),
		rule("SpaceAfterPostdecrementWhenFollowedBySubtract", ast.KindMinusMinusToken, ast.KindMinusToken, []ContextPredicate{isNonJsxSameLineTokenContext, isBinaryOpContext}, RuleActionInsertSpace),
		rule("SpaceAfterSubtractWhenFollowedByUnaryMinus", ast.KindMinusToken, ast.KindMinusToken, []ContextPredicate{isNonJsxSameLineTokenContext, isBinaryOpContext}, RuleActionInsertSpace),
		rule("SpaceAfterSubtractWhenFollowedByPredecrement", ast.KindMinusToken, ast.KindMinusMinusToken, []ContextPredicate{isNonJsxSameLineTokenContext, isBinaryOpContext}, RuleActionInsertSpace),

		rule("NoSpaceAfterCloseBrace", ast.KindCloseBraceToken, []ast.Kind{ast.KindCommaToken, ast.KindSemicolonToken}, []ContextPredicate{isNonJsxSameLineTokenContext}, RuleActionDeleteSpace),
		// For functions and control block place } on a new line []ast.Kind{multi-line rule}
		rule("NewLineBeforeCloseBraceInBlockContext", anyTokenIncludingMultilineComments, ast.KindCloseBraceToken, []ContextPredicate{isMultilineBlockContext}, RuleActionInsertNewLine),

		// Space/new line after }.
		rule("SpaceAfterCloseBrace", ast.KindCloseBraceToken, anyTokenExcept(ast.KindCloseParenToken), []ContextPredicate{isNonJsxSameLineTokenContext, isAfterCodeBlockContext}, RuleActionInsertSpace),
		// Special case for (}, else) and (}, while) since else & while tokens are not part of the tree which makes SpaceAfterCloseBrace rule not applied
		// Also should not apply to })
		rule("SpaceBetweenCloseBraceAndElse", ast.KindCloseBraceToken, ast.KindElseKeyword, []ContextPredicate{isNonJsxSameLineTokenContext}, RuleActionInsertSpace),
		rule("SpaceBetweenCloseBraceAndWhile", ast.KindCloseBraceToken, ast.KindWhileKeyword, []ContextPredicate{isNonJsxSameLineTokenContext}, RuleActionInsertSpace),
		rule("NoSpaceBetweenEmptyBraceBrackets", ast.KindOpenBraceToken, ast.KindCloseBraceToken, []ContextPredicate{isNonJsxSameLineTokenContext, isObjectContext}, RuleActionDeleteSpace),

		// Add a space after control dec context if the next character is an open bracket ex: 'if (false)[]ast.Kind{a, b} = []ast.Kind{1, 2};' -> 'if (false) []ast.Kind{a, b} = []ast.Kind{1, 2};'
		rule("SpaceAfterConditionalClosingParen", ast.KindCloseParenToken, ast.KindOpenBracketToken, []ContextPredicate{isControlDeclContext}, RuleActionInsertSpace),

		rule("NoSpaceBetweenFunctionKeywordAndStar", ast.KindFunctionKeyword, ast.KindAsteriskToken, []ContextPredicate{isFunctionDeclarationOrFunctionExpressionContext}, RuleActionDeleteSpace),
		rule("SpaceAfterStarInGeneratorDeclaration", ast.KindAsteriskToken, ast.KindIdentifier, []ContextPredicate{isFunctionDeclarationOrFunctionExpressionContext}, RuleActionInsertSpace),

		rule("SpaceAfterFunctionInFuncDecl", ast.KindFunctionKeyword, anyToken, []ContextPredicate{isFunctionDeclContext}, RuleActionInsertSpace),
		// Insert new line after { and before } in multi-line contexts.
		rule("NewLineAfterOpenBraceInBlockContext", ast.KindOpenBraceToken, anyToken, []ContextPredicate{isMultilineBlockContext}, RuleActionInsertNewLine),

		// For get/set members, we check for (identifier,identifier) since get/set don't have tokens and they are represented as just an identifier token.
		// Though, we do extra check on the context to make sure we are dealing with get/set node. Example:
		//      get x() {}
		//      set x(val) {}
		rule("SpaceAfterGetSetInMember", []ast.Kind{ast.KindGetKeyword, ast.KindSetKeyword}, ast.KindIdentifier, []ContextPredicate{isFunctionDeclContext}, RuleActionInsertSpace),

		rule("NoSpaceBetweenYieldKeywordAndStar", ast.KindYieldKeyword, ast.KindAsteriskToken, []ContextPredicate{isNonJsxSameLineTokenContext, isYieldOrYieldStarWithOperand}, RuleActionDeleteSpace),
		rule("SpaceBetweenYieldOrYieldStarAndOperand", []ast.Kind{ast.KindYieldKeyword, ast.KindAsteriskToken}, anyToken, []ContextPredicate{isNonJsxSameLineTokenContext, isYieldOrYieldStarWithOperand}, RuleActionInsertSpace),

		rule("NoSpaceBetweenReturnAndSemicolon", ast.KindReturnKeyword, ast.KindSemicolonToken, []ContextPredicate{isNonJsxSameLineTokenContext}, RuleActionDeleteSpace),
		rule("SpaceAfterCertainKeywords", []ast.Kind{ast.KindVarKeyword, ast.KindThrowKeyword, ast.KindNewKeyword, ast.KindDeleteKeyword, ast.KindReturnKeyword, ast.KindTypeOfKeyword, ast.KindAwaitKeyword}, anyToken, []ContextPredicate{isNonJsxSameLineTokenContext}, RuleActionInsertSpace),
		rule("SpaceAfterLetConstInVariableDeclaration", []ast.Kind{ast.KindLetKeyword, ast.KindConstKeyword}, anyToken, []ContextPredicate{isNonJsxSameLineTokenContext, isStartOfVariableDeclarationList}, RuleActionInsertSpace),
		rule("NoSpaceBeforeOpenParenInFuncCall", anyToken, ast.KindOpenParenToken, []ContextPredicate{isNonJsxSameLineTokenContext, isFunctionCallOrNewContext, isPreviousTokenNotComma}, RuleActionDeleteSpace),

		// Special case for binary operators (that are keywords). For these we have to add a space and shouldn't follow any user options.
		rule("SpaceBeforeBinaryKeywordOperator", anyToken, binaryKeywordOperators, []ContextPredicate{isNonJsxSameLineTokenContext, isBinaryOpContext}, RuleActionInsertSpace),
		rule("SpaceAfterBinaryKeywordOperator", binaryKeywordOperators, anyToken, []ContextPredicate{isNonJsxSameLineTokenContext, isBinaryOpContext}, RuleActionInsertSpace),

		rule("SpaceAfterVoidOperator", ast.KindVoidKeyword, anyToken, []ContextPredicate{isNonJsxSameLineTokenContext, isVoidOpContext}, RuleActionInsertSpace),

		// Async-await
		rule("SpaceBetweenAsyncAndOpenParen", ast.KindAsyncKeyword, ast.KindOpenParenToken, []ContextPredicate{isArrowFunctionContext, isNonJsxSameLineTokenContext}, RuleActionInsertSpace),
		rule("SpaceBetweenAsyncAndFunctionKeyword", ast.KindAsyncKeyword, []ast.Kind{ast.KindFunctionKeyword, ast.KindIdentifier}, []ContextPredicate{isNonJsxSameLineTokenContext}, RuleActionInsertSpace),

		// Template string
		rule("NoSpaceBetweenTagAndTemplateString", []ast.Kind{ast.KindIdentifier, ast.KindCloseParenToken}, []ast.Kind{ast.KindNoSubstitutionTemplateLiteral, ast.KindTemplateHead}, []ContextPredicate{isNonJsxSameLineTokenContext}, RuleActionDeleteSpace),

		// JSX opening elements
		rule("SpaceBeforeJsxAttribute", anyToken, ast.KindIdentifier, []ContextPredicate{isNextTokenParentJsxAttribute, isNonJsxSameLineTokenContext}, RuleActionInsertSpace),
		rule("SpaceBeforeSlashInJsxOpeningElement", anyToken, ast.KindSlashToken, []ContextPredicate{isJsxSelfClosingElementContext, isNonJsxSameLineTokenContext}, RuleActionInsertSpace),
		rule("NoSpaceBeforeGreaterThanTokenInJsxOpeningElement", ast.KindSlashToken, ast.KindGreaterThanToken, []ContextPredicate{isJsxSelfClosingElementContext, isNonJsxSameLineTokenContext}, RuleActionDeleteSpace),
		rule("NoSpaceBeforeEqualInJsxAttribute", anyToken, ast.KindEqualsToken, []ContextPredicate{isJsxAttributeContext, isNonJsxSameLineTokenContext}, RuleActionDeleteSpace),
		rule("NoSpaceAfterEqualInJsxAttribute", ast.KindEqualsToken, anyToken, []ContextPredicate{isJsxAttributeContext, isNonJsxSameLineTokenContext}, RuleActionDeleteSpace),
		rule("NoSpaceBeforeJsxNamespaceColon", ast.KindIdentifier, ast.KindColonToken, []ContextPredicate{isNextTokenParentJsxNamespacedName}, RuleActionDeleteSpace),
		rule("NoSpaceAfterJsxNamespaceColon", ast.KindColonToken, ast.KindIdentifier, []ContextPredicate{isNextTokenParentJsxNamespacedName}, RuleActionDeleteSpace),

		// TypeScript-specific rules
		// Use of module as a function call. e.g.: import m2 = module("m2");
		rule("NoSpaceAfterModuleImport", []ast.Kind{ast.KindModuleKeyword, ast.KindRequireKeyword}, ast.KindOpenParenToken, []ContextPredicate{isNonJsxSameLineTokenContext}, RuleActionDeleteSpace),
		// Add a space around certain TypeScript keywords
		rule(
			"SpaceAfterCertainTypeScriptKeywords",
			[]ast.Kind{
				ast.KindAbstractKeyword,
				ast.KindAccessorKeyword,
				ast.KindClassKeyword,
				ast.KindDeclareKeyword,
				ast.KindDefaultKeyword,
				ast.KindEnumKeyword,
				ast.KindExportKeyword,
				ast.KindExtendsKeyword,
				ast.KindGetKeyword,
				ast.KindImplementsKeyword,
				ast.KindImportKeyword,
				ast.KindInterfaceKeyword,
				ast.KindModuleKeyword,
				ast.KindNamespaceKeyword,
				ast.KindPrivateKeyword,
				ast.KindPublicKeyword,
				ast.KindProtectedKeyword,
				ast.KindReadonlyKeyword,
				ast.KindSetKeyword,
				ast.KindStaticKeyword,
				ast.KindTypeKeyword,
				ast.KindFromKeyword,
				ast.KindKeyOfKeyword,
				ast.KindInferKeyword,
			},
			anyToken,
			[]ContextPredicate{isNonJsxSameLineTokenContext},
			RuleActionInsertSpace,
		),
		rule(
			"SpaceBeforeCertainTypeScriptKeywords",
			anyToken,
			[]ast.Kind{ast.KindExtendsKeyword, ast.KindImplementsKeyword, ast.KindFromKeyword},
			[]ContextPredicate{isNonJsxSameLineTokenContext},
			RuleActionInsertSpace,
		),
		// Treat string literals in module names as identifiers, and add a space between the literal and the opening Brace braces, e.g.: module "m2" {
		rule("SpaceAfterModuleName", ast.KindStringLiteral, ast.KindOpenBraceToken, []ContextPredicate{isModuleDeclContext}, RuleActionInsertSpace),

		// Lambda expressions
		rule("SpaceBeforeArrow", anyToken, ast.KindEqualsGreaterThanToken, []ContextPredicate{isNonJsxSameLineTokenContext}, RuleActionInsertSpace),
		rule("SpaceAfterArrow", ast.KindEqualsGreaterThanToken, anyToken, []ContextPredicate{isNonJsxSameLineTokenContext}, RuleActionInsertSpace),

		// Optional parameters and let args
		rule("NoSpaceAfterEllipsis", ast.KindDotDotDotToken, ast.KindIdentifier, []ContextPredicate{isNonJsxSameLineTokenContext}, RuleActionDeleteSpace),
		rule("NoSpaceAfterOptionalParameters", ast.KindQuestionToken, []ast.Kind{ast.KindCloseParenToken, ast.KindCommaToken}, []ContextPredicate{isNonJsxSameLineTokenContext, isNotBinaryOpContext}, RuleActionDeleteSpace),

		// Remove spaces in empty interface literals. e.g.: x: {}
		rule("NoSpaceBetweenEmptyInterfaceBraceBrackets", ast.KindOpenBraceToken, ast.KindCloseBraceToken, []ContextPredicate{isNonJsxSameLineTokenContext, isObjectTypeContext}, RuleActionDeleteSpace),

		// generics and type assertions
		rule("NoSpaceBeforeOpenAngularBracket", typeNames, ast.KindLessThanToken, []ContextPredicate{isNonJsxSameLineTokenContext, isTypeArgumentOrParameterOrAssertionContext}, RuleActionDeleteSpace),
		rule("NoSpaceBetweenCloseParenAndAngularBracket", ast.KindCloseParenToken, ast.KindLessThanToken, []ContextPredicate{isNonJsxSameLineTokenContext, isTypeArgumentOrParameterOrAssertionContext}, RuleActionDeleteSpace),
		rule("NoSpaceAfterOpenAngularBracket", ast.KindLessThanToken, anyToken, []ContextPredicate{isNonJsxSameLineTokenContext, isTypeArgumentOrParameterOrAssertionContext}, RuleActionDeleteSpace),
		rule("NoSpaceBeforeCloseAngularBracket", anyToken, ast.KindGreaterThanToken, []ContextPredicate{isNonJsxSameLineTokenContext, isTypeArgumentOrParameterOrAssertionContext}, RuleActionDeleteSpace),
		rule("NoSpaceAfterCloseAngularBracket", ast.KindGreaterThanToken, []ast.Kind{ast.KindOpenParenToken, ast.KindOpenBracketToken, ast.KindGreaterThanToken, ast.KindCommaToken}, []ContextPredicate{
			isNonJsxSameLineTokenContext,
			isTypeArgumentOrParameterOrAssertionContext,
			isNotFunctionDeclContext, /*To prevent an interference with the SpaceBeforeOpenParenInFuncDecl rule*/
			isNonTypeAssertionContext,
		}, RuleActionDeleteSpace),

		// decorators
		rule("SpaceBeforeAt", []ast.Kind{ast.KindCloseParenToken, ast.KindIdentifier}, ast.KindAtToken, []ContextPredicate{isNonJsxSameLineTokenContext}, RuleActionInsertSpace),
		rule("NoSpaceAfterAt", ast.KindAtToken, anyToken, []ContextPredicate{isNonJsxSameLineTokenContext}, RuleActionDeleteSpace),
		// Insert space after @ in decorator
		rule(
			"SpaceAfterDecorator",
			anyToken,
			[]ast.Kind{
				ast.KindAbstractKeyword,
				ast.KindIdentifier,
				ast.KindExportKeyword,
				ast.KindDefaultKeyword,
				ast.KindClassKeyword,
				ast.KindStaticKeyword,
				ast.KindPublicKeyword,
				ast.KindPrivateKeyword,
				ast.KindProtectedKeyword,
				ast.KindGetKeyword,
				ast.KindSetKeyword,
				ast.KindOpenBracketToken,
				ast.KindAsteriskToken,
			},
			[]ContextPredicate{isEndOfDecoratorContextOnSameLine},
			RuleActionInsertSpace,
		),

		rule("NoSpaceBeforeNonNullAssertionOperator", anyToken, ast.KindExclamationToken, []ContextPredicate{isNonJsxSameLineTokenContext, isNonNullAssertionContext}, RuleActionDeleteSpace),
		rule("NoSpaceAfterNewKeywordOnConstructorSignature", ast.KindNewKeyword, ast.KindOpenParenToken, []ContextPredicate{isNonJsxSameLineTokenContext, isConstructorSignatureContext}, RuleActionDeleteSpace),
		rule("SpaceLessThanAndNonJSXTypeAnnotation", ast.KindLessThanToken, ast.KindLessThanToken, []ContextPredicate{isNonJsxSameLineTokenContext}, RuleActionInsertSpace),
	}

	// These rules are applied after high priority
	userConfigurableRules := []ruleSpec{
		// Treat constructor as an identifier in a function declaration, and remove spaces between constructor and following left parentheses
		rule("SpaceAfterConstructor", ast.KindConstructorKeyword, ast.KindOpenParenToken, []ContextPredicate{isOptionEnabled(insertSpaceAfterConstructorOption), isNonJsxSameLineTokenContext}, RuleActionInsertSpace),
		rule("NoSpaceAfterConstructor", ast.KindConstructorKeyword, ast.KindOpenParenToken, []ContextPredicate{isOptionDisabledOrUndefined(insertSpaceAfterConstructorOption), isNonJsxSameLineTokenContext}, RuleActionDeleteSpace),

		rule("SpaceAfterComma", ast.KindCommaToken, anyToken, []ContextPredicate{isOptionEnabled(insertSpaceAfterCommaDelimiterOption), isNonJsxSameLineTokenContext, isNonJsxElementOrFragmentContext, isNextTokenNotCloseBracket, isNextTokenNotCloseParen}, RuleActionInsertSpace),
		rule("NoSpaceAfterComma", ast.KindCommaToken, anyToken, []ContextPredicate{isOptionDisabledOrUndefined(insertSpaceAfterCommaDelimiterOption), isNonJsxSameLineTokenContext, isNonJsxElementOrFragmentContext}, RuleActionDeleteSpace),

		// Insert space after function keyword for anonymous functions
		rule("SpaceAfterAnonymousFunctionKeyword", []ast.Kind{ast.KindFunctionKeyword, ast.KindAsteriskToken}, ast.KindOpenParenToken, []ContextPredicate{isOptionEnabled(insertSpaceAfterFunctionKeywordForAnonymousFunctionsOption), isFunctionDeclContext}, RuleActionInsertSpace),
		rule("NoSpaceAfterAnonymousFunctionKeyword", []ast.Kind{ast.KindFunctionKeyword, ast.KindAsteriskToken}, ast.KindOpenParenToken, []ContextPredicate{isOptionDisabledOrUndefined(insertSpaceAfterFunctionKeywordForAnonymousFunctionsOption), isFunctionDeclContext}, RuleActionDeleteSpace),

		// Insert space after keywords in control flow statements
		rule("SpaceAfterKeywordInControl", keywords, ast.KindOpenParenToken, []ContextPredicate{isOptionEnabled(insertSpaceAfterKeywordsInControlFlowStatementsOption), isControlDeclContext}, RuleActionInsertSpace),
		rule("NoSpaceAfterKeywordInControl", keywords, ast.KindOpenParenToken, []ContextPredicate{isOptionDisabledOrUndefined(insertSpaceAfterKeywordsInControlFlowStatementsOption), isControlDeclContext}, RuleActionDeleteSpace),

		// Insert space after opening and before closing nonempty parenthesis
		rule("SpaceAfterOpenParen", ast.KindOpenParenToken, anyToken, []ContextPredicate{isOptionEnabled(insertSpaceAfterOpeningAndBeforeClosingNonemptyParenthesisOption), isNonJsxSameLineTokenContext}, RuleActionInsertSpace),
		rule("SpaceBeforeCloseParen", anyToken, ast.KindCloseParenToken, []ContextPredicate{isOptionEnabled(insertSpaceAfterOpeningAndBeforeClosingNonemptyParenthesisOption), isNonJsxSameLineTokenContext}, RuleActionInsertSpace),
		rule("SpaceBetweenOpenParens", ast.KindOpenParenToken, ast.KindOpenParenToken, []ContextPredicate{isOptionEnabled(insertSpaceAfterOpeningAndBeforeClosingNonemptyParenthesisOption), isNonJsxSameLineTokenContext}, RuleActionInsertSpace),
		rule("NoSpaceBetweenParens", ast.KindOpenParenToken, ast.KindCloseParenToken, []ContextPredicate{isNonJsxSameLineTokenContext}, RuleActionDeleteSpace),
		rule("NoSpaceAfterOpenParen", ast.KindOpenParenToken, anyToken, []ContextPredicate{isOptionDisabledOrUndefined(insertSpaceAfterOpeningAndBeforeClosingNonemptyParenthesisOption), isNonJsxSameLineTokenContext}, RuleActionDeleteSpace),
		rule("NoSpaceBeforeCloseParen", anyToken, ast.KindCloseParenToken, []ContextPredicate{isOptionDisabledOrUndefined(insertSpaceAfterOpeningAndBeforeClosingNonemptyParenthesisOption), isNonJsxSameLineTokenContext}, RuleActionDeleteSpace),

		// Insert space after opening and before closing nonempty brackets
		rule("SpaceAfterOpenBracket", ast.KindOpenBracketToken, anyToken, []ContextPredicate{isOptionEnabled(insertSpaceAfterOpeningAndBeforeClosingNonemptyBracketsOption), isNonJsxSameLineTokenContext}, RuleActionInsertSpace),
		rule("SpaceBeforeCloseBracket", anyToken, ast.KindCloseBracketToken, []ContextPredicate{isOptionEnabled(insertSpaceAfterOpeningAndBeforeClosingNonemptyBracketsOption), isNonJsxSameLineTokenContext}, RuleActionInsertSpace),
		rule("NoSpaceBetweenBrackets", ast.KindOpenBracketToken, ast.KindCloseBracketToken, []ContextPredicate{isNonJsxSameLineTokenContext}, RuleActionDeleteSpace),
		rule("NoSpaceAfterOpenBracket", ast.KindOpenBracketToken, anyToken, []ContextPredicate{isOptionDisabledOrUndefined(insertSpaceAfterOpeningAndBeforeClosingNonemptyBracketsOption), isNonJsxSameLineTokenContext}, RuleActionDeleteSpace),
		rule("NoSpaceBeforeCloseBracket", anyToken, ast.KindCloseBracketToken, []ContextPredicate{isOptionDisabledOrUndefined(insertSpaceAfterOpeningAndBeforeClosingNonemptyBracketsOption), isNonJsxSameLineTokenContext}, RuleActionDeleteSpace),

		// Insert a space after { and before } in single-line contexts, but remove space from empty object literals {}.
		rule("SpaceAfterOpenBrace", ast.KindOpenBraceToken, anyToken, []ContextPredicate{isOptionEnabledOrUndefined(insertSpaceAfterOpeningAndBeforeClosingNonemptyBracesOption), isBraceWrappedContext}, RuleActionInsertSpace),
		rule("SpaceBeforeCloseBrace", anyToken, ast.KindCloseBraceToken, []ContextPredicate{isOptionEnabledOrUndefined(insertSpaceAfterOpeningAndBeforeClosingNonemptyBracesOption), isBraceWrappedContext}, RuleActionInsertSpace),
		rule("NoSpaceBetweenEmptyBraceBrackets", ast.KindOpenBraceToken, ast.KindCloseBraceToken, []ContextPredicate{isNonJsxSameLineTokenContext, isObjectContext}, RuleActionDeleteSpace),
		rule("NoSpaceAfterOpenBrace", ast.KindOpenBraceToken, anyToken, []ContextPredicate{isOptionDisabled(insertSpaceAfterOpeningAndBeforeClosingNonemptyBracesOption), isNonJsxSameLineTokenContext}, RuleActionDeleteSpace),
		rule("NoSpaceBeforeCloseBrace", anyToken, ast.KindCloseBraceToken, []ContextPredicate{isOptionDisabled(insertSpaceAfterOpeningAndBeforeClosingNonemptyBracesOption), isNonJsxSameLineTokenContext}, RuleActionDeleteSpace),

		// Insert a space after opening and before closing empty brace brackets
		rule("SpaceBetweenEmptyBraceBrackets", ast.KindOpenBraceToken, ast.KindCloseBraceToken, []ContextPredicate{isOptionEnabled(insertSpaceAfterOpeningAndBeforeClosingEmptyBracesOption)}, RuleActionInsertSpace),
		rule("NoSpaceBetweenEmptyBraceBrackets", ast.KindOpenBraceToken, ast.KindCloseBraceToken, []ContextPredicate{isOptionDisabled(insertSpaceAfterOpeningAndBeforeClosingEmptyBracesOption), isNonJsxSameLineTokenContext}, RuleActionDeleteSpace),

		// Insert space after opening and before closing template string braces
		rule("SpaceAfterTemplateHeadAndMiddle", []ast.Kind{ast.KindTemplateHead, ast.KindTemplateMiddle}, anyToken, []ContextPredicate{isOptionEnabled(insertSpaceAfterOpeningAndBeforeClosingTemplateStringBracesOption), isNonJsxTextContext}, RuleActionInsertSpace, RuleFlagsCanDeleteNewLines),
		rule("SpaceBeforeTemplateMiddleAndTail", anyToken, []ast.Kind{ast.KindTemplateMiddle, ast.KindTemplateTail}, []ContextPredicate{isOptionEnabled(insertSpaceAfterOpeningAndBeforeClosingTemplateStringBracesOption), isNonJsxSameLineTokenContext}, RuleActionInsertSpace),
		rule("NoSpaceAfterTemplateHeadAndMiddle", []ast.Kind{ast.KindTemplateHead, ast.KindTemplateMiddle}, anyToken, []ContextPredicate{isOptionDisabledOrUndefined(insertSpaceAfterOpeningAndBeforeClosingTemplateStringBracesOption), isNonJsxTextContext}, RuleActionDeleteSpace, RuleFlagsCanDeleteNewLines),
		rule("NoSpaceBeforeTemplateMiddleAndTail", anyToken, []ast.Kind{ast.KindTemplateMiddle, ast.KindTemplateTail}, []ContextPredicate{isOptionDisabledOrUndefined(insertSpaceAfterOpeningAndBeforeClosingTemplateStringBracesOption), isNonJsxSameLineTokenContext}, RuleActionDeleteSpace),

		// No space after { and before } in JSX expression
		rule("SpaceAfterOpenBraceInJsxExpression", ast.KindOpenBraceToken, anyToken, []ContextPredicate{isOptionEnabled(insertSpaceAfterOpeningAndBeforeClosingJsxExpressionBracesOption), isNonJsxSameLineTokenContext, isJsxExpressionContext}, RuleActionInsertSpace),
		rule("SpaceBeforeCloseBraceInJsxExpression", anyToken, ast.KindCloseBraceToken, []ContextPredicate{isOptionEnabled(insertSpaceAfterOpeningAndBeforeClosingJsxExpressionBracesOption), isNonJsxSameLineTokenContext, isJsxExpressionContext}, RuleActionInsertSpace),
		rule("NoSpaceAfterOpenBraceInJsxExpression", ast.KindOpenBraceToken, anyToken, []ContextPredicate{isOptionDisabledOrUndefined(insertSpaceAfterOpeningAndBeforeClosingJsxExpressionBracesOption), isNonJsxSameLineTokenContext, isJsxExpressionContext}, RuleActionDeleteSpace),
		rule("NoSpaceBeforeCloseBraceInJsxExpression", anyToken, ast.KindCloseBraceToken, []ContextPredicate{isOptionDisabledOrUndefined(insertSpaceAfterOpeningAndBeforeClosingJsxExpressionBracesOption), isNonJsxSameLineTokenContext, isJsxExpressionContext}, RuleActionDeleteSpace),

		// Insert space after semicolon in for statement
		rule("SpaceAfterSemicolonInFor", ast.KindSemicolonToken, anyToken, []ContextPredicate{isOptionEnabled(insertSpaceAfterSemicolonInForStatementsOption), isNonJsxSameLineTokenContext, isForContext}, RuleActionInsertSpace),
		rule("NoSpaceAfterSemicolonInFor", ast.KindSemicolonToken, anyToken, []ContextPredicate{isOptionDisabledOrUndefined(insertSpaceAfterSemicolonInForStatementsOption), isNonJsxSameLineTokenContext, isForContext}, RuleActionDeleteSpace),

		// Insert space before and after binary operators
		rule("SpaceBeforeBinaryOperator", anyToken, binaryOperators, []ContextPredicate{isOptionEnabled(insertSpaceBeforeAndAfterBinaryOperatorsOption), isNonJsxSameLineTokenContext, isBinaryOpContext}, RuleActionInsertSpace),
		rule("SpaceAfterBinaryOperator", binaryOperators, anyToken, []ContextPredicate{isOptionEnabled(insertSpaceBeforeAndAfterBinaryOperatorsOption), isNonJsxSameLineTokenContext, isBinaryOpContext}, RuleActionInsertSpace),
		rule("NoSpaceBeforeBinaryOperator", anyToken, binaryOperators, []ContextPredicate{isOptionDisabledOrUndefined(insertSpaceBeforeAndAfterBinaryOperatorsOption), isNonJsxSameLineTokenContext, isBinaryOpContext}, RuleActionDeleteSpace),
		rule("NoSpaceAfterBinaryOperator", binaryOperators, anyToken, []ContextPredicate{isOptionDisabledOrUndefined(insertSpaceBeforeAndAfterBinaryOperatorsOption), isNonJsxSameLineTokenContext, isBinaryOpContext}, RuleActionDeleteSpace),

		rule("SpaceBeforeOpenParenInFuncDecl", anyToken, ast.KindOpenParenToken, []ContextPredicate{isOptionEnabled(insertSpaceBeforeFunctionParenthesisOption), isNonJsxSameLineTokenContext, isFunctionDeclContext}, RuleActionInsertSpace),
		rule("NoSpaceBeforeOpenParenInFuncDecl", anyToken, ast.KindOpenParenToken, []ContextPredicate{isOptionDisabledOrUndefined(insertSpaceBeforeFunctionParenthesisOption), isNonJsxSameLineTokenContext, isFunctionDeclContext}, RuleActionDeleteSpace),

		// Open Brace braces after control block
		rule("NewLineBeforeOpenBraceInControl", controlOpenBraceLeftTokenRange, ast.KindOpenBraceToken, []ContextPredicate{isOptionEnabled(placeOpenBraceOnNewLineForControlBlocksOption), isControlDeclContext, isBeforeMultilineBlockContext}, RuleActionInsertNewLine, RuleFlagsCanDeleteNewLines),

		// Open Brace braces after function
		// TypeScript: Function can have return types, which can be made of tons of different token kinds
		rule("NewLineBeforeOpenBraceInFunction", functionOpenBraceLeftTokenRange, ast.KindOpenBraceToken, []ContextPredicate{isOptionEnabled(placeOpenBraceOnNewLineForFunctionsOption), isFunctionDeclContext, isBeforeMultilineBlockContext}, RuleActionInsertNewLine, RuleFlagsCanDeleteNewLines),
		// Open Brace braces after TypeScript module/class/interface
		rule("NewLineBeforeOpenBraceInTypeScriptDeclWithBlock", typeScriptOpenBraceLeftTokenRange, ast.KindOpenBraceToken, []ContextPredicate{isOptionEnabled(placeOpenBraceOnNewLineForFunctionsOption), isTypeScriptDeclWithBlockContext, isBeforeMultilineBlockContext}, RuleActionInsertNewLine, RuleFlagsCanDeleteNewLines),

		rule("SpaceAfterTypeAssertion", ast.KindGreaterThanToken, anyToken, []ContextPredicate{isOptionEnabled(insertSpaceAfterTypeAssertionOption), isNonJsxSameLineTokenContext, isTypeAssertionContext}, RuleActionInsertSpace),
		rule("NoSpaceAfterTypeAssertion", ast.KindGreaterThanToken, anyToken, []ContextPredicate{isOptionDisabledOrUndefined(insertSpaceAfterTypeAssertionOption), isNonJsxSameLineTokenContext, isTypeAssertionContext}, RuleActionDeleteSpace),

		rule("SpaceBeforeTypeAnnotation", anyToken, []ast.Kind{ast.KindQuestionToken, ast.KindColonToken}, []ContextPredicate{isOptionEnabled(insertSpaceBeforeTypeAnnotationOption), isNonJsxSameLineTokenContext, isTypeAnnotationContext}, RuleActionInsertSpace),
		rule("NoSpaceBeforeTypeAnnotation", anyToken, []ast.Kind{ast.KindQuestionToken, ast.KindColonToken}, []ContextPredicate{isOptionDisabledOrUndefined(insertSpaceBeforeTypeAnnotationOption), isNonJsxSameLineTokenContext, isTypeAnnotationContext}, RuleActionDeleteSpace),

		rule("NoOptionalSemicolon", ast.KindSemicolonToken, anyTokenIncludingEOF, []ContextPredicate{optionEquals(semicolonOption, SemicolonPreferenceRemove), isSemicolonDeletionContext}, RuleActionDeleteToken),
		rule("OptionalSemicolon", anyToken, anyTokenIncludingEOF, []ContextPredicate{optionEquals(semicolonOption, SemicolonPreferenceInsert), isSemicolonInsertionContext}, RuleActionInsertTrailingSemicolon),
	}

	// These rules are lower in priority than user-configurable. Rules earlier in this list have priority over rules later in the list.
	lowPriorityCommonRules := []ruleSpec{
		// Space after keyword but not before ; or : or ?
		rule("NoSpaceBeforeSemicolon", anyToken, ast.KindSemicolonToken, []ContextPredicate{isNonJsxSameLineTokenContext}, RuleActionDeleteSpace),

		rule("SpaceBeforeOpenBraceInControl", controlOpenBraceLeftTokenRange, ast.KindOpenBraceToken, []ContextPredicate{isOptionDisabledOrUndefinedOrTokensOnSameLine(placeOpenBraceOnNewLineForControlBlocksOption), isControlDeclContext, isNotFormatOnEnter, isSameLineTokenOrBeforeBlockContext}, RuleActionInsertSpace, RuleFlagsCanDeleteNewLines),
		rule("SpaceBeforeOpenBraceInFunction", functionOpenBraceLeftTokenRange, ast.KindOpenBraceToken, []ContextPredicate{isOptionDisabledOrUndefinedOrTokensOnSameLine(placeOpenBraceOnNewLineForFunctionsOption), isFunctionDeclContext, isBeforeBlockContext, isNotFormatOnEnter, isSameLineTokenOrBeforeBlockContext}, RuleActionInsertSpace, RuleFlagsCanDeleteNewLines),
		rule("SpaceBeforeOpenBraceInTypeScriptDeclWithBlock", typeScriptOpenBraceLeftTokenRange, ast.KindOpenBraceToken, []ContextPredicate{isOptionDisabledOrUndefinedOrTokensOnSameLine(placeOpenBraceOnNewLineForFunctionsOption), isTypeScriptDeclWithBlockContext, isNotFormatOnEnter, isSameLineTokenOrBeforeBlockContext}, RuleActionInsertSpace, RuleFlagsCanDeleteNewLines),

		rule("NoSpaceBeforeComma", anyToken, ast.KindCommaToken, []ContextPredicate{isNonJsxSameLineTokenContext}, RuleActionDeleteSpace),

		// No space before and after indexer `x[]ast.Kind{}`
		rule("NoSpaceBeforeOpenBracket", anyTokenExcept(ast.KindAsyncKeyword, ast.KindCaseKeyword), ast.KindOpenBracketToken, []ContextPredicate{isNonJsxSameLineTokenContext}, RuleActionDeleteSpace),
		rule("NoSpaceAfterCloseBracket", ast.KindCloseBracketToken, anyToken, []ContextPredicate{isNonJsxSameLineTokenContext, isNotBeforeBlockInFunctionDeclarationContext}, RuleActionDeleteSpace),
		rule("SpaceAfterSemicolon", ast.KindSemicolonToken, anyToken, []ContextPredicate{isNonJsxSameLineTokenContext}, RuleActionInsertSpace),

		// Remove extra space between for and await
		rule("SpaceBetweenForAndAwaitKeyword", ast.KindForKeyword, ast.KindAwaitKeyword, []ContextPredicate{isNonJsxSameLineTokenContext}, RuleActionInsertSpace),

		// Remove extra spaces between ... and type name in tuple spread
		rule("SpaceBetweenDotDotDotAndTypeName", ast.KindDotDotDotToken, typeNames, []ContextPredicate{isNonJsxSameLineTokenContext}, RuleActionDeleteSpace),

		// Add a space between statements. All keywords except (do,else,case) has open/close parens after them.
		// So, we have a rule to add a space for []ast.Kind{),Any}, []ast.Kind{do,Any}, []ast.Kind{else,Any}, and []ast.Kind{case,Any}
		rule(
			"SpaceBetweenStatements",
			[]ast.Kind{ast.KindCloseParenToken, ast.KindDoKeyword, ast.KindElseKeyword, ast.KindCaseKeyword},
			anyToken,
			[]ContextPredicate{isNonJsxSameLineTokenContext, isNonJsxElementOrFragmentContext, isNotForContext},
			RuleActionInsertSpace,
		),
		// This low-pri rule takes care of "try {", "catch {" and "finally {" in case the rule SpaceBeforeOpenBraceInControl didn't execute on FormatOnEnter.
		rule("SpaceAfterTryCatchFinally", []ast.Kind{ast.KindTryKeyword, ast.KindCatchKeyword, ast.KindFinallyKeyword}, ast.KindOpenBraceToken, []ContextPredicate{isNonJsxSameLineTokenContext}, RuleActionInsertSpace),
	}

	result := make([]ruleSpec, 0, len(highPriorityCommonRules)+len(userConfigurableRules)+len(lowPriorityCommonRules))
	result = append(result, highPriorityCommonRules...)
	result = append(result, userConfigurableRules...)
	result = append(result, lowPriorityCommonRules...)
	return result
}

func tokenRangeFrom(tokens ...ast.Kind) tokenRange {
	return tokenRange{
		isSpecific: true,
		tokens:     tokens,
	}
}

func tokenRangeFromEx(prefix []ast.Kind, tokens ...ast.Kind) tokenRange {
	tokens = append(prefix, tokens...)
	return tokenRange{
		isSpecific: true,
		tokens:     tokens,
	}
}

func tokenRangeFromRange(start ast.Kind, end ast.Kind) tokenRange {
	tokens := make([]ast.Kind, 0, end-start+1)
	for token := start; token <= end; token++ {
		tokens = append(tokens, token)
	}

	return tokenRangeFrom(tokens...)
}
