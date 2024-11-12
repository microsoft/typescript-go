// Package printer exports a Printer for pretty-printing TS ASTs and writer interfaces and implementations for using them
// Intended ultimate usage:
//
//		func nodeToInlineStr(node *Node) {
//	   // Reuse singleton single-line writer (TODO: thread safety?)
//		  printer := printer.New({ writer: printer.SingleLineTextWriter, stripComments: true })
//		  printer.printNode(node)
//		  return printer.getText()
//		}
//
// // or
//
//		func nodeToStr(node *Node, options CompilerOptions) {
//	   // create new writer shared for the entire printing operation
//		  printer := printer.New({ writer: printer.NewTextWriter(options.newLine) })
//		  printer.printNode(node)
//		  return printer.getText()
//		}
package printer

import "github.com/microsoft/typescript-go/internal/compiler"

// Prints a node into a string - creates its' own text writer to facilitate this - prefer emitter.PrintNode where an emitter is available
func PrintNode(node *compiler.Node) string {
	writer := NewTextWriter("\n")
	// printNode(node, writer)
	return writer.getText()
}

// func printNode(node *compiler.Node, writer EmitTextWriter) {
// 	// TODO: Port emitter.ts's emit code more faithfully
// 	switch node.Kind {
// 	case compiler.ast.KindUnknown, compiler.ast.KindEndOfFile, compiler.ast.KindConflictMarkerTrivia, compiler.ast.KindNonTextFileMarkerTrivia:
// 		panic("non-printing node kind")
// 	case compiler.ast.KindNumericLiteral:
// 	case compiler.ast.KindBigintLiteral:
// 	case compiler.ast.KindStringLiteral:
// 	case compiler.ast.KindJsxText:
// 	case compiler.ast.KindJsxTextAllWhiteSpaces:
// 	case compiler.ast.KindRegularExpressionLiteral:
// 	case compiler.ast.KindNoSubstitutionTemplateLiteral:
// 	case compiler.ast.KindTemplateHead:
// 	case compiler.ast.KindTemplateMiddle:
// 	case compiler.ast.KindTemplateTail:
// 	case compiler.ast.KindOpenBraceToken:
// 	case compiler.ast.KindCloseBraceToken:
// 	case compiler.ast.KindOpenParenToken:
// 	case compiler.ast.KindCloseParenToken:
// 	case compiler.ast.KindOpenBracketToken:
// 	case compiler.ast.KindCloseBracketToken:
// 	case compiler.ast.KindDotToken:
// 	case compiler.ast.KindDotDotDotToken:
// 	case compiler.ast.KindSemicolonToken:
// 	case compiler.ast.KindCommaToken:
// 	case compiler.ast.KindQuestionDotToken:
// 	case compiler.ast.KindLessThanToken:
// 	case compiler.ast.KindLessThanSlashToken:
// 	case compiler.ast.KindGreaterThanToken:
// 	case compiler.ast.KindLessThanEqualsToken:
// 	case compiler.ast.KindGreaterThanEqualsToken:
// 	case compiler.ast.KindEqualsEqualsToken:
// 	case compiler.ast.KindExclamationEqualsToken:
// 	case compiler.ast.KindEqualsEqualsEqualsToken:
// 	case compiler.ast.KindExclamationEqualsEqualsToken:
// 	case compiler.ast.KindEqualsGreaterThanToken:
// 	case compiler.ast.KindPlusToken:
// 	case compiler.ast.KindMinusToken:
// 	case compiler.ast.KindAsteriskToken:
// 	case compiler.ast.KindAsteriskAsteriskToken:
// 	case compiler.ast.KindSlashToken:
// 	case compiler.ast.KindPercentToken:
// 	case compiler.ast.KindPlusPlusToken:
// 	case compiler.ast.KindMinusMinusToken:
// 	case compiler.ast.KindLessThanLessThanToken:
// 	case compiler.ast.KindGreaterThanGreaterThanToken:
// 	case compiler.ast.KindGreaterThanGreaterThanGreaterThanToken:
// 	case compiler.ast.KindAmpersandToken:
// 	case compiler.ast.KindBarToken:
// 	case compiler.ast.KindCaretToken:
// 	case compiler.ast.KindExclamationToken:
// 	case compiler.ast.KindTildeToken:
// 	case compiler.ast.KindAmpersandAmpersandToken:
// 	case compiler.ast.KindBarBarToken:
// 	case compiler.ast.KindQuestionToken:
// 	case compiler.ast.KindColonToken:
// 	case compiler.ast.KindAtToken:
// 	case compiler.ast.KindQuestionQuestionToken:
// 	/** Only the JSDoc scanner produces BacktickToken. The normal scanner produces NoSubstitutionTemplateLiteral and related kinds. */
// 	case compiler.ast.KindBacktickToken:
// 	/** Only the JSDoc scanner produces HashToken. The normal scanner produces PrivateIdentifier. */
// 	case compiler.ast.KindHashToken:
// 	// Assignments
// 	case compiler.ast.KindEqualsToken:
// 	case compiler.ast.KindPlusEqualsToken:
// 	case compiler.ast.KindMinusEqualsToken:
// 	case compiler.ast.KindAsteriskEqualsToken:
// 	case compiler.ast.KindAsteriskAsteriskEqualsToken:
// 	case compiler.ast.KindSlashEqualsToken:
// 	case compiler.ast.KindPercentEqualsToken:
// 	case compiler.ast.KindLessThanLessThanEqualsToken:
// 	case compiler.ast.KindGreaterThanGreaterThanEqualsToken:
// 	case compiler.ast.KindGreaterThanGreaterThanGreaterThanEqualsToken:
// 	case compiler.ast.KindAmpersandEqualsToken:
// 	case compiler.ast.KindBarEqualsToken:
// 	case compiler.ast.KindBarBarEqualsToken:
// 	case compiler.ast.KindAmpersandAmpersandEqualsToken:
// 	case compiler.ast.KindQuestionQuestionEqualsToken:
// 	case compiler.ast.KindCaretEqualsToken:
// 	// Identifiers and PrivateIdentifier
// 	case compiler.ast.KindIdentifier:
// 	case compiler.ast.KindPrivateIdentifier:
// 	case compiler.ast.KindJSDocCommentTextToken:
// 	// Reserved words
// 	case compiler.ast.KindBreakKeyword:
// 	case compiler.ast.KindCaseKeyword:
// 	case compiler.ast.KindCatchKeyword:
// 	case compiler.ast.KindClassKeyword:
// 	case compiler.ast.KindConstKeyword:
// 	case compiler.ast.KindContinueKeyword:
// 	case compiler.ast.KindDebuggerKeyword:
// 	case compiler.ast.KindDefaultKeyword:
// 	case compiler.ast.KindDeleteKeyword:
// 	case compiler.ast.KindDoKeyword:
// 	case compiler.ast.KindElseKeyword:
// 	case compiler.ast.KindEnumKeyword:
// 	case compiler.ast.KindExportKeyword:
// 	case compiler.ast.KindExtendsKeyword:
// 	case compiler.ast.KindFalseKeyword:
// 	case compiler.ast.KindFinallyKeyword:
// 	case compiler.ast.KindForKeyword:
// 	case compiler.ast.KindFunctionKeyword:
// 	case compiler.ast.KindIfKeyword:
// 	case compiler.ast.KindImportKeyword:
// 	case compiler.ast.KindInKeyword:
// 	case compiler.ast.KindInstanceOfKeyword:
// 	case compiler.ast.KindNewKeyword:
// 	case compiler.ast.KindNullKeyword:
// 	case compiler.ast.KindReturnKeyword:
// 	case compiler.ast.KindSuperKeyword:
// 	case compiler.ast.KindSwitchKeyword:
// 	case compiler.ast.KindThisKeyword:
// 	case compiler.ast.KindThrowKeyword:
// 	case compiler.ast.KindTrueKeyword:
// 	case compiler.ast.KindTryKeyword:
// 	case compiler.ast.KindTypeOfKeyword:
// 	case compiler.ast.KindVarKeyword:
// 	case compiler.ast.KindVoidKeyword:
// 	case compiler.ast.KindWhileKeyword:
// 	case compiler.ast.KindWithKeyword:
// 	case compiler.ast.KindImplementsKeyword:
// 	case compiler.ast.KindInterfaceKeyword:
// 	case compiler.ast.KindLetKeyword:
// 	case compiler.ast.KindPackageKeyword:
// 	case compiler.ast.KindPrivateKeyword:
// 	case compiler.ast.KindProtectedKeyword:
// 	case compiler.ast.KindPublicKeyword:
// 	case compiler.ast.KindStaticKeyword:
// 	case compiler.ast.KindYieldKeyword:
// 	case compiler.ast.KindAbstractKeyword:
// 	case compiler.ast.KindAccessorKeyword:
// 	case compiler.ast.KindAsKeyword:
// 	case compiler.ast.KindAssertsKeyword:
// 	case compiler.ast.KindAssertKeyword:
// 	case compiler.ast.KindAnyKeyword:
// 	case compiler.ast.KindAsyncKeyword:
// 	case compiler.ast.KindAwaitKeyword:
// 	case compiler.ast.KindBooleanKeyword:
// 	case compiler.ast.KindConstructorKeyword:
// 	case compiler.ast.KindDeclareKeyword:
// 	case compiler.ast.KindGetKeyword:
// 	case compiler.ast.KindImmediateKeyword:
// 	case compiler.ast.KindInferKeyword:
// 	case compiler.ast.KindIntrinsicKeyword:
// 	case compiler.ast.KindIsKeyword:
// 	case compiler.ast.KindKeyOfKeyword:
// 	case compiler.ast.KindModuleKeyword:
// 	case compiler.ast.KindNamespaceKeyword:
// 	case compiler.ast.KindNeverKeyword:
// 	case compiler.ast.KindOutKeyword:
// 	case compiler.ast.KindReadonlyKeyword:
// 	case compiler.ast.KindRequireKeyword:
// 	case compiler.ast.KindNumberKeyword:
// 	case compiler.ast.KindObjectKeyword:
// 	case compiler.ast.KindSatisfiesKeyword:
// 	case compiler.ast.KindSetKeyword:
// 	case compiler.ast.KindStringKeyword:
// 	case compiler.ast.KindSymbolKeyword:
// 	case compiler.ast.KindTypeKeyword:
// 	case compiler.ast.KindUndefinedKeyword:
// 	case compiler.ast.KindUniqueKeyword:
// 	case compiler.ast.KindUnknownKeyword:
// 	case compiler.ast.KindUsingKeyword:
// 	case compiler.ast.KindFromKeyword:
// 	case compiler.ast.KindGlobalKeyword:
// 	case compiler.ast.KindBigIntKeyword:
// 	case compiler.ast.KindOverrideKeyword:
// 	case compiler.ast.KindOfKeyword: // LastKeyword and LastToken and LastContextualKeyword
// 	case compiler.ast.KindQualifiedName:
// 	case compiler.ast.KindComputedPropertyName:
// 	case compiler.ast.KindModifierList:
// 	case compiler.ast.KindTypeParameterList:
// 	case compiler.ast.KindTypeArgumentList:
// 	case compiler.ast.KindTypeParameter:
// 	case compiler.ast.KindParameter:
// 	case compiler.ast.KindDecorator:
// 	case compiler.ast.KindPropertySignature:
// 	case compiler.ast.KindPropertyDeclaration:
// 	case compiler.ast.KindMethodSignature:
// 	case compiler.ast.KindMethodDeclaration:
// 	case compiler.ast.KindClassStaticBlockDeclaration:
// 	case compiler.ast.KindConstructor:
// 	case compiler.ast.KindGetAccessor:
// 	case compiler.ast.KindSetAccessor:
// 	case compiler.ast.KindCallSignature:
// 	case compiler.ast.KindConstructSignature:
// 	case compiler.ast.KindIndexSignature:
// 	case compiler.ast.KindTypePredicate:
// 	case compiler.ast.KindTypeReference:
// 	case compiler.ast.KindFunctionType:
// 	case compiler.ast.KindConstructorType:
// 	case compiler.ast.KindTypeQuery:
// 	case compiler.ast.KindTypeLiteral:
// 	case compiler.ast.KindArrayType:
// 	case compiler.ast.KindTupleType:
// 	case compiler.ast.KindOptionalType:
// 	case compiler.ast.KindRestType:
// 	case compiler.ast.KindUnionType:
// 	case compiler.ast.KindIntersectionType:
// 	case compiler.ast.KindConditionalType:
// 	case compiler.ast.KindInferType:
// 	case compiler.ast.KindParenthesizedType:
// 	case compiler.ast.KindThisType:
// 	case compiler.ast.KindTypeOperator:
// 	case compiler.ast.KindIndexedAccessType:
// 	case compiler.ast.KindMappedType:
// 	case compiler.ast.KindLiteralType:
// 	case compiler.ast.KindNamedTupleMember:
// 	case compiler.ast.KindTemplateLiteralType:
// 	case compiler.ast.KindTemplateLiteralTypeSpan:
// 	case compiler.ast.KindImportType:
// 	case compiler.ast.KindObjectBindingPattern:
// 	case compiler.ast.KindArrayBindingPattern:
// 	case compiler.ast.KindBindingElement:
// 	case compiler.ast.KindArrayLiteralExpression:
// 	case compiler.ast.KindObjectLiteralExpression:
// 	case compiler.ast.KindPropertyAccessExpression:
// 	case compiler.ast.KindElementAccessExpression:
// 	case compiler.ast.KindCallExpression:
// 	case compiler.ast.KindNewExpression:
// 	case compiler.ast.KindTaggedTemplateExpression:
// 	case compiler.ast.KindTypeAssertionExpression:
// 	case compiler.ast.KindParenthesizedExpression:
// 	case compiler.ast.KindFunctionExpression:
// 	case compiler.ast.KindArrowFunction:
// 	case compiler.ast.KindDeleteExpression:
// 	case compiler.ast.KindTypeOfExpression:
// 	case compiler.ast.KindVoidExpression:
// 	case compiler.ast.KindAwaitExpression:
// 	case compiler.ast.KindPrefixUnaryExpression:
// 	case compiler.ast.KindPostfixUnaryExpression:
// 	case compiler.ast.KindBinaryExpression:
// 	case compiler.ast.KindConditionalExpression:
// 	case compiler.ast.KindTemplateExpression:
// 	case compiler.ast.KindYieldExpression:
// 	case compiler.ast.KindSpreadElement:
// 	case compiler.ast.KindClassExpression:
// 	case compiler.ast.KindOmittedExpression:
// 	case compiler.ast.KindExpressionWithTypeArguments:
// 	case compiler.ast.KindAsExpression:
// 	case compiler.ast.KindNonNullExpression:
// 	case compiler.ast.KindMetaProperty:
// 	case compiler.ast.KindSyntheticExpression:
// 	case compiler.ast.KindSatisfiesExpression:
// 	case compiler.ast.KindTemplateSpan:
// 	case compiler.ast.KindSemicolonClassElement:
// 	case compiler.ast.KindBlock:
// 	case compiler.ast.KindEmptyStatement:
// 	case compiler.ast.KindVariableStatement:
// 	case compiler.ast.KindExpressionStatement:
// 	case compiler.ast.KindIfStatement:
// 	case compiler.ast.KindDoStatement:
// 	case compiler.ast.KindWhileStatement:
// 	case compiler.ast.KindForStatement:
// 	case compiler.ast.KindForInStatement:
// 	case compiler.ast.KindForOfStatement:
// 	case compiler.ast.KindContinueStatement:
// 	case compiler.ast.KindBreakStatement:
// 	case compiler.ast.KindReturnStatement:
// 	case compiler.ast.KindWithStatement:
// 	case compiler.ast.KindSwitchStatement:
// 	case compiler.ast.KindLabeledStatement:
// 	case compiler.ast.KindThrowStatement:
// 	case compiler.ast.KindTryStatement:
// 	case compiler.ast.KindDebuggerStatement:
// 	case compiler.ast.KindVariableDeclaration:
// 	case compiler.ast.KindVariableDeclarationList:
// 	case compiler.ast.KindFunctionDeclaration:
// 	case compiler.ast.KindClassDeclaration:
// 	case compiler.ast.KindInterfaceDeclaration:
// 	case compiler.ast.KindTypeAliasDeclaration:
// 	case compiler.ast.KindEnumDeclaration:
// 	case compiler.ast.KindModuleDeclaration:
// 	case compiler.ast.KindModuleBlock:
// 	case compiler.ast.KindCaseBlock:
// 	case compiler.ast.KindNamespaceExportDeclaration:
// 	case compiler.ast.KindImportEqualsDeclaration:
// 	case compiler.ast.KindImportDeclaration:
// 	case compiler.ast.KindImportClause:
// 	case compiler.ast.KindNamespaceImport:
// 	case compiler.ast.KindNamedImports:
// 	case compiler.ast.KindImportSpecifier:
// 	case compiler.ast.KindExportAssignment:
// 	case compiler.ast.KindExportDeclaration:
// 	case compiler.ast.KindNamedExports:
// 	case compiler.ast.KindNamespaceExport:
// 	case compiler.ast.KindExportSpecifier:
// 	case compiler.ast.KindMissingDeclaration:
// 	case compiler.ast.KindExternalModuleReference:
// 	case compiler.ast.KindJsxElement:
// 	case compiler.ast.KindJsxSelfClosingElement:
// 	case compiler.ast.KindJsxOpeningElement:
// 	case compiler.ast.KindJsxClosingElement:
// 	case compiler.ast.KindJsxFragment:
// 	case compiler.ast.KindJsxOpeningFragment:
// 	case compiler.ast.KindJsxClosingFragment:
// 	case compiler.ast.KindJsxAttribute:
// 	case compiler.ast.KindJsxAttributes:
// 	case compiler.ast.KindJsxSpreadAttribute:
// 	case compiler.ast.KindJsxExpression:
// 	case compiler.ast.KindJsxNamespacedName:
// 	case compiler.ast.KindCaseClause:
// 	case compiler.ast.KindDefaultClause:
// 	case compiler.ast.KindHeritageClause:
// 	case compiler.ast.KindCatchClause:
// 	case compiler.ast.KindImportAttributes:
// 	case compiler.ast.KindImportAttribute:
// 	case compiler.ast.KindPropertyAssignment:
// 	case compiler.ast.KindShorthandPropertyAssignment:
// 	case compiler.ast.KindSpreadAssignment:
// 	case compiler.ast.KindEnumMember:
// 	case compiler.ast.KindSourceFile:
// 		// for i, statement := range node.AsSourceFile().Statements {
// 		// 	if i != 0 {
// 		// 		writer.writeLine()
// 		// 	}
// 		// 	printNode(statement, writer)
// 		// }
// 	case compiler.ast.KindBundle:
// 	case compiler.ast.KindJSDocTypeExpression:
// 	case compiler.ast.KindJSDocNameReference:
// 	case compiler.ast.KindJSDocMemberName: // C#p
// 	case compiler.ast.KindJSDocAllType: // The * type
// 	case compiler.ast.KindJSDocUnknownType: // The ? type
// 	case compiler.ast.KindJSDocNullableType:
// 	case compiler.ast.KindJSDocNonNullableType:
// 	case compiler.ast.KindJSDocOptionalType:
// 	case compiler.ast.KindJSDocFunctionType:
// 	case compiler.ast.KindJSDocVariadicType:
// 	case compiler.ast.KindJSDocNamepathType: // https://jsdoc.app/about-namepaths.html
// 	case compiler.ast.KindJSDoc:
// 	case compiler.ast.KindJSDocText:
// 	case compiler.ast.KindJSDocTypeLiteral:
// 	case compiler.ast.KindJSDocSignature:
// 	case compiler.ast.KindJSDocLink:
// 	case compiler.ast.KindJSDocLinkCode:
// 	case compiler.ast.KindJSDocLinkPlain:
// 	case compiler.ast.KindJSDocTag:
// 	case compiler.ast.KindJSDocAugmentsTag:
// 	case compiler.ast.KindJSDocImplementsTag:
// 	case compiler.ast.KindJSDocAuthorTag:
// 	case compiler.ast.KindJSDocDeprecatedTag:
// 	case compiler.ast.KindJSDocImmediateTag:
// 	case compiler.ast.KindJSDocClassTag:
// 	case compiler.ast.KindJSDocPublicTag:
// 	case compiler.ast.KindJSDocPrivateTag:
// 	case compiler.ast.KindJSDocProtectedTag:
// 	case compiler.ast.KindJSDocReadonlyTag:
// 	case compiler.ast.KindJSDocOverrideTag:
// 	case compiler.ast.KindJSDocCallbackTag:
// 	case compiler.ast.KindJSDocOverloadTag:
// 	case compiler.ast.KindJSDocEnumTag:
// 	case compiler.ast.KindJSDocParameterTag:
// 	case compiler.ast.KindJSDocReturnTag:
// 	case compiler.ast.KindJSDocThisTag:
// 	case compiler.ast.KindJSDocTypeTag:
// 	case compiler.ast.KindJSDocTemplateTag:
// 	case compiler.ast.KindJSDocTypedefTag:
// 	case compiler.ast.KindJSDocSeeTag:
// 	case compiler.ast.KindJSDocPropertyTag:
// 	case compiler.ast.KindJSDocThrowsTag:
// 	case compiler.ast.KindJSDocSatisfiesTag:
// 	case compiler.ast.KindJSDocImportTag:
// 	case compiler.ast.KindSyntaxList:
// 	case compiler.ast.KindNotEmittedStatement:
// 	case compiler.ast.KindPartiallyEmittedExpression:
// 	case compiler.ast.KindCommaListExpression:
// 	case compiler.ast.KindSyntheticReferenceExpression:
// 	default:
// 		panic("Node kind not implemented in printer")
// 	}
// }
