package checker

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/jsnum"
	"github.com/microsoft/typescript-go/internal/nodebuilder"
	"github.com/microsoft/typescript-go/internal/printer"
	"github.com/microsoft/typescript-go/internal/scanner"
)

// TODO: Memoize once per checker to retain threadsafety
func createPrinterWithDefaults(emitContext *printer.EmitContext) *printer.Printer {
	return printer.NewPrinter(printer.PrinterOptions{}, printer.PrintHandlers{}, emitContext)
}

func createPrinterWithRemoveComments(emitContext *printer.EmitContext) *printer.Printer {
	return printer.NewPrinter(printer.PrinterOptions{RemoveComments: true}, printer.PrintHandlers{}, emitContext)
}

func createPrinterWithRemoveCommentsOmitTrailingSemicolon(emitContext *printer.EmitContext) *printer.Printer {
	// TODO: OmitTrailingSemicolon support
	return printer.NewPrinter(printer.PrinterOptions{RemoveComments: true}, printer.PrintHandlers{}, emitContext)
}

func createPrinterWithRemoveCommentsNeverAsciiEscape(emitContext *printer.EmitContext) *printer.Printer {
	// TODO: NeverAsciiEscape support
	return printer.NewPrinter(printer.PrinterOptions{RemoveComments: true}, printer.PrintHandlers{}, emitContext)
}

func getTrailingSemicolonDeferringWriter(writer printer.EmitTextWriter) printer.EmitTextWriter {
	// TODO: wrap arbitrary writer with writer that only commits semicolon writes on following write operations (is OmitTrailingSemicolon printer option redundant?)
	return writer
}

func (c *Checker) TypeToString(type_ *Type) string {
	return c.typeToStringEx(type_, nil, TypeFormatFlagsNone, nil)
}

func toNodeBuilderFlags(flags TypeFormatFlags) nodebuilder.Flags {
	return nodebuilder.Flags(flags & TypeFormatFlagsNodeBuilderFlagsMask)
}

func (c *Checker) typeToStringEx(type_ *Type, enclosingDeclaration *ast.Node, flags TypeFormatFlags, writer printer.EmitTextWriter) string {
	if writer == nil {
		writer = printer.NewTextWriter("")
	}
	noTruncation := (c.compilerOptions.NoErrorTruncation == core.TSTrue) || (flags&TypeFormatFlagsNoTruncation != 0)
	combinedFlags := toNodeBuilderFlags(flags) | nodebuilder.FlagsIgnoreErrors
	if noTruncation {
		combinedFlags = combinedFlags | nodebuilder.FlagsNoTruncation
	}
	typeNode := c.nodeBuilder.typeToTypeNode(type_, enclosingDeclaration, combinedFlags, nodebuilder.InternalFlagsNone, nil)
	if typeNode == nil {
		panic("should always get typenode")
	}
	// The unresolved type gets a synthesized comment on `any` to hint to users that it's not a plain `any`.
	// Otherwise, we always strip comments out.
	var printer *printer.Printer
	if type_ == c.unresolvedType {
		printer = createPrinterWithDefaults(c.diagnosticConstructionContext)
	} else {
		printer = createPrinterWithRemoveComments(c.diagnosticConstructionContext)
	}
	var sourceFile *ast.SourceFile
	if enclosingDeclaration != nil {
		sourceFile = ast.GetSourceFileOfNode(enclosingDeclaration)
	}
	printer.Write(typeNode /*sourceFile*/, sourceFile, writer, nil)
	result := writer.String()

	maxLength := defaultMaximumTruncationLength * 2
	if noTruncation {
		maxLength = noTruncationMaximumTruncationLength * 2
	}
	if maxLength > 0 && result != "" && len(result) >= maxLength {
		return result[0:maxLength-len("...")] + "..."
	}
	return result
}

func (c *Checker) SymbolToString(s *ast.Symbol) string {
	return c.symbolToString(s)
}

func (c *Checker) symbolToString(symbol *ast.Symbol) string {
	return c.symbolToStringEx(symbol, nil, ast.SymbolFlagsAll, SymbolFormatFlagsAllowAnyNodeKind, nil)
}

func (c *Checker) symbolToStringEx(symbol *ast.Symbol, enclosingDeclaration *ast.Node, meaning ast.SymbolFlags, flags SymbolFormatFlags, writer printer.EmitTextWriter) string {
	if writer == nil {
		writer = printer.SingleLineStringWriterPool.Get().(printer.EmitTextWriter)
	}

	nodeFlags := nodebuilder.FlagsIgnoreErrors
	internalNodeFlags := nodebuilder.InternalFlagsNone
	if flags&SymbolFormatFlagsUseOnlyExternalAliasing != 0 {
		nodeFlags |= nodebuilder.FlagsUseOnlyExternalAliasing
	}
	if flags&SymbolFormatFlagsWriteTypeParametersOrArguments != 0 {
		nodeFlags |= nodebuilder.FlagsWriteTypeParametersInQualifiedName
	}
	if flags&SymbolFormatFlagsUseAliasDefinedOutsideCurrentScope != 0 {
		nodeFlags |= nodebuilder.FlagsUseAliasDefinedOutsideCurrentScope
	}
	if flags&SymbolFormatFlagsDoNotIncludeSymbolChain != 0 {
		internalNodeFlags |= nodebuilder.InternalFlagsDoNotIncludeSymbolChain
	}
	if flags&SymbolFormatFlagsWriteComputedProps != 0 {
		internalNodeFlags |= nodebuilder.InternalFlagsWriteComputedProps
	}

	var sourceFile *ast.SourceFile
	if enclosingDeclaration != nil {
		sourceFile = ast.GetSourceFileOfNode(enclosingDeclaration)
	}
	if writer == printer.SingleLineStringWriterPool.Get().(printer.EmitTextWriter) {
		// handle uses of the single-line writer during an ongoing write
		existing := writer.String()
		defer writer.Clear()
		if existing != "" {
			defer writer.WriteKeyword(existing)
		}
	}
	var printer_ *printer.Printer
	if enclosingDeclaration != nil && enclosingDeclaration.Kind == ast.KindSourceFile {
		printer_ = createPrinterWithRemoveCommentsNeverAsciiEscape(c.diagnosticConstructionContext)
	} else {
		printer_ = createPrinterWithRemoveComments(c.diagnosticConstructionContext)
	}

	var builder func(symbol *ast.Symbol, meaning ast.SymbolFlags, enclosingDeclaration *ast.Node, flags nodebuilder.Flags, internalFlags nodebuilder.InternalFlags, tracker nodebuilder.SymbolTracker) *ast.Node
	if flags&SymbolFormatFlagsAllowAnyNodeKind != 0 {
		builder = c.nodeBuilder.symbolToNode
	} else {
		builder = c.nodeBuilder.symbolToEntityName
	}
	entity := builder(symbol, meaning, enclosingDeclaration, nodeFlags, internalNodeFlags, nil)         // TODO: GH#18217
	printer_.Write(entity /*sourceFile*/, sourceFile, getTrailingSemicolonDeferringWriter(writer), nil) // TODO: GH#18217
	return writer.String()
}

func (c *Checker) signatureToString(signature *Signature) string {
	return c.signatureToStringEx(signature, nil, TypeFormatFlagsNone, nil, nil)
}

func (c *Checker) signatureToStringEx(signature *Signature, enclosingDeclaration *ast.Node, flags TypeFormatFlags, kind *SignatureKind, writer printer.EmitTextWriter) string {
	var sigOutput ast.Kind
	if flags&TypeFormatFlagsWriteArrowStyleSignature != 0 {
		if kind != nil && *kind == SignatureKindConstruct {
			sigOutput = ast.KindConstructorType
		} else {
			sigOutput = ast.KindFunctionType
		}
	} else {
		if kind != nil && *kind == SignatureKindConstruct {
			sigOutput = ast.KindConstructSignature
		} else {
			sigOutput = ast.KindCallSignature
		}
	}
	if writer == nil {
		writer = printer.SingleLineStringWriterPool.Get().(printer.EmitTextWriter)
	}
	combinedFlags := toNodeBuilderFlags(flags) | nodebuilder.FlagsIgnoreErrors | nodebuilder.FlagsWriteTypeParametersInQualifiedName
	sig := c.nodeBuilder.signatureToSignatureDeclaration(signature, sigOutput, enclosingDeclaration, combinedFlags, nodebuilder.InternalFlagsNone, nil)
	printer_ := createPrinterWithRemoveCommentsOmitTrailingSemicolon(c.diagnosticConstructionContext)
	var sourceFile *ast.SourceFile
	if enclosingDeclaration != nil {
		sourceFile = ast.GetSourceFileOfNode(enclosingDeclaration)
	}
	if writer == printer.SingleLineStringWriterPool.Get().(printer.EmitTextWriter) {
		// handle uses of the single-line writer during an ongoing write
		existing := writer.String()
		defer writer.Clear()
		if existing != "" {
			defer writer.WriteKeyword(existing)
		}
	}
	printer_.Write(sig /*sourceFile*/, sourceFile, getTrailingSemicolonDeferringWriter(writer), nil) // TODO: GH#18217
	return writer.String()
}

func (c *Checker) typePredicateToString(typePredicate *TypePredicate) string {
	return c.typePredicateToStringEx(typePredicate, nil, TypeFormatFlagsUseAliasDefinedOutsideCurrentScope, nil)
}

func (c *Checker) typePredicateToStringEx(typePredicate *TypePredicate, enclosingDeclaration *ast.Node, flags TypeFormatFlags, writer printer.EmitTextWriter) string {
	if writer == nil {
		writer = printer.SingleLineStringWriterPool.Get().(printer.EmitTextWriter)
	}
	combinedFlags := toNodeBuilderFlags(flags) | nodebuilder.FlagsIgnoreErrors | nodebuilder.FlagsWriteTypeParametersInQualifiedName
	predicate := c.nodeBuilder.typePredicateToTypePredicateNode(typePredicate, enclosingDeclaration, combinedFlags, nodebuilder.InternalFlagsNone, nil) // TODO: GH#18217
	printer_ := createPrinterWithRemoveComments(c.diagnosticConstructionContext)
	var sourceFile *ast.SourceFile
	if enclosingDeclaration != nil {
		sourceFile = ast.GetSourceFileOfNode(enclosingDeclaration)
	}
	if writer == printer.SingleLineStringWriterPool.Get().(printer.EmitTextWriter) {
		// handle uses of the single-line writer during an ongoing write
		existing := writer.String()
		defer writer.Clear()
		if existing != "" {
			defer writer.WriteKeyword(existing)
		}
	}
	printer_.Write(predicate /*sourceFile*/, sourceFile, writer, nil)
	return writer.String()
}

func (c *Checker) valueToString(value any) string {
	switch value := value.(type) {
	case string:
		return "\"" + printer.EscapeString(value, '"') + "\""
	case jsnum.Number:
		return value.String()
	case bool:
		return core.IfElse(value, "true", "false")
	case jsnum.PseudoBigInt:
		return value.String() + "n"
	}
	panic("unhandled value type in valueToString")
}

func (c *Checker) formatUnionTypes(types []*Type) []*Type {
	var result []*Type
	var flags TypeFlags
	for i := 0; i < len(types); i++ {
		t := types[i]
		flags |= t.flags
		if t.flags&TypeFlagsNullable == 0 {
			if t.flags&(TypeFlagsBooleanLiteral|TypeFlagsEnumLike) != 0 {
				var baseType *Type
				if t.flags&TypeFlagsBooleanLiteral != 0 {
					baseType = c.booleanType
				} else {
					baseType = c.getBaseTypeOfEnumLikeType(t)
				}
				if baseType.flags&TypeFlagsUnion != 0 {
					count := len(baseType.AsUnionType().types)
					if i+count <= len(types) && c.getRegularTypeOfLiteralType(types[i+count-1]) == c.getRegularTypeOfLiteralType(baseType.AsUnionType().types[count-1]) {
						result = append(result, baseType)
						i += count - 1
						continue
					}
				}
			}
			result = append(result, t)
		}
	}
	if flags&TypeFlagsNull != 0 {
		result = append(result, c.nullType)
	}
	if flags&TypeFlagsUndefined != 0 {
		result = append(result, c.undefinedType)
	}
	return result
}

func (c *Checker) SourceFileWithTypes(sourceFile *ast.SourceFile) string {
	writer := printer.NewTextWriter("\n")
	var pos int
	var visit func(*ast.Node) bool
	var typesPrinted bool
	lineStarts := scanner.GetLineStarts(sourceFile)
	printLinesBefore := func(node *ast.Node) {
		line := scanner.ComputeLineOfPosition(lineStarts, scanner.SkipTrivia(sourceFile.Text(), node.Pos()))
		var nextLineStart int
		if line+1 < len(lineStarts) {
			nextLineStart = int(lineStarts[line+1])
		} else {
			nextLineStart = sourceFile.Loc.End()
		}
		if pos < nextLineStart {
			if typesPrinted {
				writer.WriteLine()
			}
			writer.Write(sourceFile.Text()[pos:nextLineStart])
			pos = nextLineStart
			typesPrinted = false
		}
	}
	visit = func(node *ast.Node) bool {
		text, t, isDeclaration := c.getTextAndTypeOfNode(node)
		if text != "" && !strings.Contains(text, "\n") {
			printLinesBefore(node)
			writer.Write(">")
			writer.Write(text)
			writer.Write(" : ")
			c.typeToStringEx(t, nil, TypeFormatFlagsNone, writer)
			if isDeclaration && t.flags&TypeFlagsEnumLiteral != 0 && t.flags&(TypeFlagsStringLiteral|TypeFlagsNumberLiteral) != 0 {
				writer.Write(" = ")
				writer.Write(c.valueToString(t.AsLiteralType().value))
			}
			writer.WriteLine()
			typesPrinted = true
		}
		return node.ForEachChild(visit)
	}
	visit(sourceFile.AsNode())
	writer.Write(sourceFile.Text()[pos:sourceFile.End()])
	return writer.String()
}

func (c *Checker) getTextAndTypeOfNode(node *ast.Node) (string, *Type, bool) {
	if ast.IsDeclarationNode(node) {
		symbol := node.Symbol()
		if symbol != nil && !isReservedMemberName(symbol.Name) {
			if symbol.Flags&ast.SymbolFlagsValue != 0 {
				return c.symbolToString(symbol), c.getTypeOfSymbol(symbol), true
			}
			if symbol.Flags&ast.SymbolFlagsTypeAlias != 0 {
				return c.symbolToString(symbol), c.getDeclaredTypeOfTypeAlias(symbol), true
			}
		}
	}
	if ast.IsExpressionNode(node) && !isRightSideOfQualifiedNameOrPropertyAccess(node) {
		return scanner.GetTextOfNode(node), c.getTypeOfExpression(node), false
	}
	return "", nil, false
}
