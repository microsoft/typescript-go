package api

import (
	"encoding/binary"
	"fmt"
	"slices"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/checker"
)

// EncodedNodeLength is the number of int32 values that represent a single node in the encoded format.
const (
	EncodedKind = iota
	EncodedPos
	EncodedEnd
	EncodedNext
	EncodedParent
	EncodedData
	EncodedNodeLength
)

const (
	NodeDataTypeChildren int32 = iota
	NodeDataTypeStringIndex
	NodeDataTypeExtendedDataIndex
)

const (
	EncodedHeaderLength = 6 * 4
)

type stringTable struct {
	data    *strings.Builder
	offsets []int32
}

func newStringTable(stringLength int, stringCount int) *stringTable {
	builder := &strings.Builder{}
	builder.Grow(stringLength)
	return &stringTable{
		data:    builder,
		offsets: make([]int32, 0, stringCount),
	}
}

func (t *stringTable) add(s string) int32 {
	offset := int32(t.data.Len())
	t.offsets = append(t.offsets, offset)
	t.data.WriteString(s)
	return int32(len(t.offsets) - 1)
}

func (t *stringTable) encode() (result []byte, err error) {
	result = make([]byte, 0, len(t.offsets)*4+t.data.Len())
	if result, err = appendInt32s(result, t.offsets...); err != nil {
		return nil, err
	}
	result = append(result, t.data.String()...)
	return result, nil
}

// EncodeSourceFile encodes a source file into a byte slice.
// The encoded format is a sequence of int32 values, where each node is represented by 5 values:
// - kind: the node kind
// - pos: the start position of the node
// - end: the end position of the node
// - next: the index of the next sibling node (0 if there is no next sibling)
// - parent: the index of the parent node (0 if there is no parent)
// The first encoded node is a zero element that is not part of the source file.
func EncodeSourceFile(sourceFile *ast.SourceFile) ([]byte, error) {
	var err error
	var parentIndex, nodeCount, prevIndex int32

	strs := newStringTable(sourceFile.TextLength, sourceFile.TextCount)
	nodes := make([]byte, 0, (sourceFile.NodeCount+1)*EncodedNodeLength*4)

	visitor := &ast.NodeVisitor{
		Hooks: ast.NodeVisitorHooks{
			VisitNodes: func(nodeList *ast.NodeList, visitor *ast.NodeVisitor) *ast.NodeList {
				if nodeList == nil || len(nodeList.Nodes) == 0 {
					return nodeList
				}

				nodeCount++
				if prevIndex != 0 {
					// this is the next sibling of `prevNode`
					b0, b1, b2, b3 := uint8(nodeCount), uint8(nodeCount>>8), uint8(nodeCount>>16), uint8(nodeCount>>24)
					nodes[prevIndex*EncodedNodeLength*4+EncodedNext*4] = b0
					nodes[prevIndex*EncodedNodeLength*4+EncodedNext*4+1] = b1
					nodes[prevIndex*EncodedNodeLength*4+EncodedNext*4+2] = b2
					nodes[prevIndex*EncodedNodeLength*4+EncodedNext*4+3] = b3
				}

				if nodes, err = appendInt32s(nodes, int32(-1), int32(nodeList.Pos()), int32(nodeList.End()), 0, parentIndex, int32(len(nodeList.Nodes))); err != nil {
					return nil
				}

				saveParentIndex := parentIndex

				currentIndex := nodeCount
				prevIndex = 0
				parentIndex = currentIndex
				visitor.VisitSlice(nodeList.Nodes)
				prevIndex = currentIndex
				parentIndex = saveParentIndex

				return nodeList
			},
			VisitModifiers: func(modifiers *ast.ModifierList, visitor *ast.NodeVisitor) *ast.ModifierList {
				if modifiers != nil && len(modifiers.Nodes) > 0 {
					visitor.Hooks.VisitNodes(&modifiers.NodeList, visitor)
				}
				return modifiers
			},
		},
	}
	visitor.Visit = func(node *ast.Node) *ast.Node {
		nodeCount++
		if prevIndex != 0 {
			// this is the next sibling of `prevNode`
			b0, b1, b2, b3 := uint8(nodeCount), uint8(nodeCount>>8), uint8(nodeCount>>16), uint8(nodeCount>>24)
			nodes[prevIndex*EncodedNodeLength*4+EncodedNext*4] = b0
			nodes[prevIndex*EncodedNodeLength*4+EncodedNext*4+1] = b1
			nodes[prevIndex*EncodedNodeLength*4+EncodedNext*4+2] = b2
			nodes[prevIndex*EncodedNodeLength*4+EncodedNext*4+3] = b3
		}

		if nodes, err = appendInt32s(nodes, int32(node.Kind), int32(node.Pos()), int32(node.End()), 0, parentIndex, getNodeData(node, strs)); err != nil {
			visitor.Visit = nil
			return nil
		}

		saveParentIndex := parentIndex

		currentIndex := nodeCount
		prevIndex = 0
		parentIndex = currentIndex
		visitor.VisitEachChild(node)
		prevIndex = currentIndex
		parentIndex = saveParentIndex
		return node
	}

	// kind, pos, end, next, parent
	if nodes, err = appendInt32s(nodes, 0, 0, 0, 0, 0, 0); err != nil {
		return nil, err
	}

	nodeCount++
	parentIndex++
	if nodes, err = appendInt32s(nodes, int32(sourceFile.Kind), int32(sourceFile.Pos()), int32(sourceFile.End()), 0, 0, 0); err != nil {
		return nil, err
	}

	visitor.VisitEachChild(sourceFile.AsNode())
	if err != nil {
		return nil, err
	}

	headerLength := EncodedHeaderLength
	offsetStringTableOffsets := headerLength
	offsetStringTableData := headerLength + len(strs.offsets)*4
	offsetNodes := offsetStringTableData + strs.data.Len()
	offsetExtendedDataOffsets := 0
	offsetExtendedDataData := 0

	header := []int32{
		0,
		int32(offsetStringTableOffsets),
		int32(offsetStringTableData),
		int32(offsetExtendedDataOffsets),
		int32(offsetExtendedDataData),
		int32(offsetNodes),
	}

	var headerBytes, strsBytes []byte
	if headerBytes, err = appendInt32s(nil, header...); err != nil {
		return nil, err
	}
	if strsBytes, err = strs.encode(); err != nil {
		return nil, err
	}

	return slices.Concat(
		headerBytes,
		strsBytes,
		nodes,
	), nil
}

func appendInt32s(buf []byte, values ...int32) ([]byte, error) {
	for _, value := range values {
		var err error
		if buf, err = binary.Append(buf, binary.LittleEndian, value); err != nil {
			return nil, err
		}
	}
	return buf, nil
}

func readInt32(buf []byte, offset int) int32 {
	return int32(binary.LittleEndian.Uint32(buf[offset : offset+4]))
}

func FormatEncodedSourceFile(encoded []byte) string {
	var result strings.Builder
	var getIndent func(parentIndex int32) string
	offsetNodes := readInt32(encoded, 5*4)
	getIndent = func(parentIndex int32) string {
		if parentIndex == 0 {
			return ""
		}
		return "  " + getIndent(readInt32(encoded, int(offsetNodes)+int(parentIndex)*EncodedNodeLength*4+EncodedParent*4))
	}
	j := 1
	for i := int(offsetNodes) + EncodedNodeLength*4; i < len(encoded); i += EncodedNodeLength * 4 {
		kind := readInt32(encoded, i+EncodedKind*4)
		pos := readInt32(encoded, i+EncodedPos*4)
		end := readInt32(encoded, i+EncodedEnd*4)
		parentIndex := readInt32(encoded, i+EncodedParent*4)
		result.WriteString(getIndent(parentIndex))
		if kind == -1 {
			result.WriteString("NodeList")
		} else {
			result.WriteString(ast.Kind(kind).String())
		}
		fmt.Fprintf(&result, " [%d, %d), i=%d, next=%d", pos, end, j, encoded[i+EncodedNext*4])
		result.WriteString("\n")
		j++
	}
	return result.String()
}

func getNodeData(node *ast.Node, strs *stringTable) int32 {
	t := getNodeDataType(node)
	switch t {
	case NodeDataTypeChildren:
		return t | getNodeDefinedData(node) | int32(getChildrenPropertyMask(node))
	case NodeDataTypeStringIndex:
		return t | getNodeDefinedData(node) | recordNodeStrings(node, strs)
	case NodeDataTypeExtendedDataIndex:
		return t | getNodeDefinedData(node) /* | TODO */
	default:
		panic("unreachable")
	}
}

func getNodeDataType(node *ast.Node) int32 {
	switch node.Kind {
	case ast.KindJsxText,
		ast.KindIdentifier,
		ast.KindPrivateIdentifier,
		ast.KindStringLiteral,
		ast.KindNumericLiteral,
		ast.KindBigIntLiteral,
		ast.KindRegularExpressionLiteral,
		ast.KindNoSubstitutionTemplateLiteral,
		ast.KindJSDocText:
		return NodeDataTypeStringIndex
	case ast.KindTemplateHead,
		ast.KindTemplateMiddle,
		ast.KindTemplateTail:
		return NodeDataTypeExtendedDataIndex
	default:
		return NodeDataTypeChildren
	}
}

// getChildrenPropertyMask returns a mask of which children properties are present in the node.
// It is defined for node kinds that have more than one property that is a pointer to a child node.
// Example: QualifiedName has two children properties: Left and Right, which are visited in that order.
// result&1 is non-zero if Left is present, and result&2 is non-zero if Right is present. If the client
// knows that QualifiedName has properties ["Left", "Right"] and sees an encoded node with only one
// child, it can use the mask to determine which property is present.
func getChildrenPropertyMask(node *ast.Node) uint8 {
	switch node.Kind {
	case ast.KindQualifiedName:
		n := node.AsQualifiedName()
		return (boolToByte(n.Left != nil) << 0) | (boolToByte(n.Right != nil) << 1)
	case ast.KindTypeParameter:
		n := node.AsTypeParameter()
		return (boolToByte(n.Modifiers() != nil) << 0) | (boolToByte(n.Name() != nil) << 1) | (boolToByte(n.Constraint != nil) << 2) | (boolToByte(n.DefaultType != nil) << 3)
	case ast.KindIfStatement:
		n := node.AsIfStatement()
		return (boolToByte(n.Expression != nil) << 0) | (boolToByte(n.ThenStatement != nil) << 1) | (boolToByte(n.ElseStatement != nil) << 2)
	case ast.KindDoStatement:
		n := node.AsDoStatement()
		return (boolToByte(n.Statement != nil) << 0) | (boolToByte(n.Expression != nil) << 1)
	case ast.KindWhileStatement:
		n := node.AsWhileStatement()
		return (boolToByte(n.Expression != nil) << 0) | (boolToByte(n.Statement != nil) << 1)
	case ast.KindForStatement:
		n := node.AsForStatement()
		return (boolToByte(n.Initializer != nil) << 0) | (boolToByte(n.Condition != nil) << 1) | (boolToByte(n.Incrementor != nil) << 2) | (boolToByte(n.Statement != nil) << 3)
	case ast.KindForInStatement, ast.KindForOfStatement:
		n := node.AsForInOrOfStatement()
		return (boolToByte(n.AwaitModifier != nil) << 0) | (boolToByte(n.Initializer != nil) << 1) | (boolToByte(n.Expression != nil) << 2) | (boolToByte(n.Statement != nil) << 3)
	case ast.KindWithStatement:
		n := node.AsWithStatement()
		return (boolToByte(n.Expression != nil) << 0) | (boolToByte(n.Statement != nil) << 1)
	case ast.KindSwitchStatement:
		n := node.AsSwitchStatement()
		return (boolToByte(n.Expression != nil) << 0) | (boolToByte(n.CaseBlock != nil) << 1)
	case ast.KindCaseClause, ast.KindDefaultClause:
		n := node.AsCaseOrDefaultClause()
		return (boolToByte(n.Expression != nil) << 0) | (boolToByte(n.Statements != nil) << 1)
	case ast.KindTryStatement:
		n := node.AsTryStatement()
		return (boolToByte(n.TryBlock != nil) << 0) | (boolToByte(n.CatchClause != nil) << 1) | (boolToByte(n.FinallyBlock != nil) << 2)
	case ast.KindCatchClause:
		n := node.AsCatchClause()
		return (boolToByte(n.VariableDeclaration != nil) << 0) | (boolToByte(n.Block != nil) << 1)
	case ast.KindLabeledStatement:
		n := node.AsLabeledStatement()
		return (boolToByte(n.Label != nil) << 0) | (boolToByte(n.Statement != nil) << 1)
	case ast.KindVariableStatement:
		n := node.AsVariableStatement()
		return (boolToByte(n.Modifiers() != nil) << 0) | (boolToByte(n.DeclarationList != nil) << 1)
	case ast.KindVariableDeclaration:
		n := node.AsVariableDeclaration()
		return (boolToByte(n.Name() != nil) << 0) | (boolToByte(n.ExclamationToken != nil) << 1) | (boolToByte(n.Type != nil) << 2) | (boolToByte(n.Initializer != nil) << 3)
	case ast.KindParameter:
		n := node.AsParameterDeclaration()
		return (boolToByte(n.Modifiers() != nil) << 0) | (boolToByte(n.DotDotDotToken != nil) << 1) | (boolToByte(n.Name() != nil) << 2) | (boolToByte(n.QuestionToken != nil) << 3) | (boolToByte(n.Type != nil) << 4) | (boolToByte(n.Initializer != nil) << 5)
	case ast.KindBindingElement:
		n := node.AsBindingElement()
		return (boolToByte(n.DotDotDotToken != nil) << 0) | (boolToByte(n.PropertyName != nil) << 1) | (boolToByte(n.Name() != nil) << 2) | (boolToByte(n.Initializer != nil) << 3)
	case ast.KindFunctionDeclaration:
		n := node.AsFunctionDeclaration()
		return (boolToByte(n.Modifiers() != nil) << 0) | (boolToByte(n.AsteriskToken != nil) << 1) | (boolToByte(n.Name() != nil) << 2) | (boolToByte(n.TypeParameters != nil) << 3) | (boolToByte(n.Parameters != nil) << 4) | (boolToByte(n.Type != nil) << 5) | (boolToByte(n.Body != nil) << 6)
	case ast.KindInterfaceDeclaration:
		n := node.AsInterfaceDeclaration()
		return (boolToByte(n.Modifiers() != nil) << 0) | (boolToByte(n.Name() != nil) << 1) | (boolToByte(n.TypeParameters != nil) << 2) | (boolToByte(n.HeritageClauses != nil) << 3) | (boolToByte(n.Members != nil) << 4)
	case ast.KindTypeAliasDeclaration:
		n := node.AsTypeAliasDeclaration()
		return (boolToByte(n.Modifiers() != nil) << 0) | (boolToByte(n.Name() != nil) << 1) | (boolToByte(n.TypeParameters != nil) << 2) | (boolToByte(n.Type != nil) << 3)
	case ast.KindEnumMember:
		n := node.AsEnumMember()
		return (boolToByte(n.Name() != nil) << 0) | (boolToByte(n.Initializer != nil) << 1)
	case ast.KindEnumDeclaration:
		n := node.AsEnumDeclaration()
		return (boolToByte(n.Modifiers() != nil) << 0) | (boolToByte(n.Name() != nil) << 1) | (boolToByte(n.Members != nil) << 2)
	case ast.KindModuleDeclaration:
		n := node.AsModuleDeclaration()
		return (boolToByte(n.Modifiers() != nil) << 0) | (boolToByte(n.Name() != nil) << 1) | (boolToByte(n.Body != nil) << 2)
	case ast.KindImportEqualsDeclaration:
		n := node.AsImportEqualsDeclaration()
		return (boolToByte(n.Modifiers() != nil) << 0) | (boolToByte(n.Name() != nil) << 1) | (boolToByte(n.ModuleReference != nil) << 2)
	case ast.KindImportDeclaration:
		n := node.AsImportDeclaration()
		return (boolToByte(n.Modifiers() != nil) << 0) | (boolToByte(n.ImportClause != nil) << 1) | (boolToByte(n.ModuleSpecifier != nil) << 2) | (boolToByte(n.Attributes != nil) << 3)
	case ast.KindImportSpecifier:
		n := node.AsImportSpecifier()
		return (boolToByte(n.PropertyName != nil) << 0) | (boolToByte(n.Name() != nil) << 1)
	case ast.KindImportClause:
		n := node.AsImportClause()
		return (boolToByte(n.Name() != nil) << 0) | (boolToByte(n.NamedBindings != nil) << 1)
	case ast.KindExportAssignment:
		n := node.AsExportAssignment()
		return (boolToByte(n.Modifiers() != nil) << 0) | (boolToByte(n.Expression != nil) << 1)
	case ast.KindNamespaceExportDeclaration:
		n := node.AsNamespaceExportDeclaration()
		return (boolToByte(n.Modifiers() != nil) << 0) | (boolToByte(n.Name() != nil) << 1)
	case ast.KindExportDeclaration:
		n := node.AsExportDeclaration()
		return (boolToByte(n.Modifiers() != nil) << 0) | (boolToByte(n.ExportClause != nil) << 1) | (boolToByte(n.ModuleSpecifier != nil) << 2) | (boolToByte(n.Attributes != nil) << 3)
	case ast.KindExportSpecifier:
		n := node.AsExportSpecifier()
		return (boolToByte(n.PropertyName != nil) << 0) | (boolToByte(n.Name() != nil) << 1)
	case ast.KindCallSignature:
		n := node.AsCallSignatureDeclaration()
		return (boolToByte(n.TypeParameters != nil) << 0) | (boolToByte(n.Parameters != nil) << 1) | (boolToByte(n.Type != nil) << 2)
	case ast.KindConstructSignature:
		n := node.AsConstructSignatureDeclaration()
		return (boolToByte(n.TypeParameters != nil) << 0) | (boolToByte(n.Parameters != nil) << 1) | (boolToByte(n.Type != nil) << 2)
	case ast.KindConstructor:
		n := node.AsConstructorDeclaration()
		return (boolToByte(n.Modifiers() != nil) << 0) | (boolToByte(n.TypeParameters != nil) << 1) | (boolToByte(n.Parameters != nil) << 2) | (boolToByte(n.Type != nil) << 3) | (boolToByte(n.Body != nil) << 4)
	case ast.KindGetAccessor:
		n := node.AsGetAccessorDeclaration()
		return (boolToByte(n.Modifiers() != nil) << 0) | (boolToByte(n.Name() != nil) << 1) | (boolToByte(n.TypeParameters != nil) << 2) | (boolToByte(n.Parameters != nil) << 3) | (boolToByte(n.Type != nil) << 4) | (boolToByte(n.Body != nil) << 5)
	case ast.KindSetAccessor:
		n := node.AsSetAccessorDeclaration()
		return (boolToByte(n.Modifiers() != nil) << 0) | (boolToByte(n.Name() != nil) << 1) | (boolToByte(n.TypeParameters != nil) << 2) | (boolToByte(n.Parameters != nil) << 3) | (boolToByte(n.Type != nil) << 4) | (boolToByte(n.Body != nil) << 5)
	case ast.KindIndexSignature:
		n := node.AsIndexSignatureDeclaration()
		return (boolToByte(n.Modifiers() != nil) << 0) | (boolToByte(n.Parameters != nil) << 1) | (boolToByte(n.Type != nil) << 2)
	case ast.KindMethodSignature:
		n := node.AsMethodSignatureDeclaration()
		return (boolToByte(n.Modifiers() != nil) << 0) | (boolToByte(n.Name() != nil) << 1) | (boolToByte(n.PostfixToken != nil) << 2) | (boolToByte(n.TypeParameters != nil) << 3) | (boolToByte(n.Parameters != nil) << 4) | (boolToByte(n.Type != nil) << 5)
	case ast.KindMethodDeclaration:
		n := node.AsMethodDeclaration()
		return (boolToByte(n.Modifiers() != nil) << 0) | (boolToByte(n.AsteriskToken != nil) << 1) | (boolToByte(n.Name() != nil) << 2) | (boolToByte(n.PostfixToken != nil) << 3) | (boolToByte(n.TypeParameters != nil) << 4) | (boolToByte(n.Parameters != nil) << 5) | (boolToByte(n.Type != nil) << 6) | (boolToByte(n.Body != nil) << 7)
	case ast.KindPropertySignature:
		n := node.AsPropertySignatureDeclaration()
		return (boolToByte(n.Modifiers() != nil) << 0) | (boolToByte(n.Name() != nil) << 1) | (boolToByte(n.PostfixToken != nil) << 2) | (boolToByte(n.Type != nil) << 3) | (boolToByte(n.Initializer != nil) << 4)
	case ast.KindPropertyDeclaration:
		n := node.AsPropertyDeclaration()
		return (boolToByte(n.Modifiers() != nil) << 0) | (boolToByte(n.Name() != nil) << 1) | (boolToByte(n.PostfixToken != nil) << 2) | (boolToByte(n.Type != nil) << 3) | (boolToByte(n.Initializer != nil) << 4)
	case ast.KindBinaryExpression:
		n := node.AsBinaryExpression()
		return (boolToByte(n.Left != nil) << 0) | (boolToByte(n.OperatorToken != nil) << 1) | (boolToByte(n.Right != nil) << 2)
	case ast.KindYieldExpression:
		n := node.AsYieldExpression()
		return (boolToByte(n.AsteriskToken != nil) << 0) | (boolToByte(n.Expression != nil) << 1)
	case ast.KindArrowFunction:
		n := node.AsArrowFunction()
		return (boolToByte(n.Modifiers() != nil) << 0) | (boolToByte(n.TypeParameters != nil) << 1) | (boolToByte(n.Parameters != nil) << 2) | (boolToByte(n.Type != nil) << 3) | (boolToByte(n.EqualsGreaterThanToken != nil) << 4) | (boolToByte(n.Body != nil) << 5)
	case ast.KindFunctionExpression:
		n := node.AsFunctionExpression()
		return (boolToByte(n.Modifiers() != nil) << 0) | (boolToByte(n.AsteriskToken != nil) << 1) | (boolToByte(n.Name() != nil) << 2) | (boolToByte(n.TypeParameters != nil) << 3) | (boolToByte(n.Parameters != nil) << 4) | (boolToByte(n.Type != nil) << 5) | (boolToByte(n.Body != nil) << 6)
	case ast.KindAsExpression:
		n := node.AsAsExpression()
		return (boolToByte(n.Expression != nil) << 0) | (boolToByte(n.Type != nil) << 1)
	case ast.KindSatisfiesExpression:
		n := node.AsSatisfiesExpression()
		return (boolToByte(n.Expression != nil) << 0) | (boolToByte(n.Type != nil) << 1)
	case ast.KindConditionalExpression:
		n := node.AsConditionalExpression()
		return (boolToByte(n.Condition != nil) << 0) | (boolToByte(n.QuestionToken != nil) << 1) | (boolToByte(n.WhenTrue != nil) << 2) | (boolToByte(n.ColonToken != nil) << 3) | (boolToByte(n.WhenFalse != nil) << 4)
	case ast.KindPropertyAccessExpression:
		n := node.AsPropertyAccessExpression()
		return (boolToByte(n.Expression != nil) << 0) | (boolToByte(n.QuestionDotToken != nil) << 1) | (boolToByte(n.Name() != nil) << 2)
	case ast.KindElementAccessExpression:
		n := node.AsElementAccessExpression()
		return (boolToByte(n.Expression != nil) << 0) | (boolToByte(n.QuestionDotToken != nil) << 1) | (boolToByte(n.ArgumentExpression != nil) << 2)
	case ast.KindCallExpression:
		n := node.AsCallExpression()
		return (boolToByte(n.Expression != nil) << 0) | (boolToByte(n.QuestionDotToken != nil) << 1) | (boolToByte(n.TypeArguments != nil) << 2) | (boolToByte(n.Arguments != nil) << 3)
	case ast.KindNewExpression:
		n := node.AsNewExpression()
		return (boolToByte(n.Expression != nil) << 0) | (boolToByte(n.TypeArguments != nil) << 1) | (boolToByte(n.Arguments != nil) << 2)
	case ast.KindTemplateExpression:
		n := node.AsTemplateExpression()
		return (boolToByte(n.Head != nil) << 0) | (boolToByte(n.TemplateSpans != nil) << 1)
	case ast.KindTemplateSpan:
		n := node.AsTemplateSpan()
		return (boolToByte(n.Expression != nil) << 0) | (boolToByte(n.Literal != nil) << 1)
	case ast.KindTaggedTemplateExpression:
		n := node.AsTaggedTemplateExpression()
		return (boolToByte(n.Tag != nil) << 0) | (boolToByte(n.QuestionDotToken != nil) << 1) | (boolToByte(n.TypeArguments != nil) << 2) | (boolToByte(n.Template != nil) << 3)
	case ast.KindPropertyAssignment:
		n := node.AsPropertyAssignment()
		return (boolToByte(n.Modifiers() != nil) << 0) | (boolToByte(n.Name() != nil) << 1) | (boolToByte(n.PostfixToken != nil) << 2) | (boolToByte(n.Initializer != nil) << 3)
	case ast.KindShorthandPropertyAssignment:
		n := node.AsShorthandPropertyAssignment()
		return (boolToByte(n.Modifiers() != nil) << 0) | (boolToByte(n.Name() != nil) << 1) | (boolToByte(n.PostfixToken != nil) << 2) | (boolToByte(n.EqualsToken != nil) << 3) | (boolToByte(n.ObjectAssignmentInitializer != nil) << 4)
	case ast.KindTypeAssertionExpression:
		n := node.AsTypeAssertion()
		return (boolToByte(n.Type != nil) << 0) | (boolToByte(n.Expression != nil) << 1)
	case ast.KindConditionalType:
		n := node.AsConditionalTypeNode()
		return (boolToByte(n.CheckType != nil) << 0) | (boolToByte(n.ExtendsType != nil) << 1) | (boolToByte(n.TrueType != nil) << 2) | (boolToByte(n.FalseType != nil) << 3)
	case ast.KindIndexedAccessType:
		n := node.AsIndexedAccessTypeNode()
		return (boolToByte(n.ObjectType != nil) << 0) | (boolToByte(n.IndexType != nil) << 1)
	case ast.KindTypeReference:
		n := node.AsTypeReferenceNode()
		return (boolToByte(n.TypeName != nil) << 0) | (boolToByte(n.TypeArguments != nil) << 1)
	case ast.KindExpressionWithTypeArguments:
		n := node.AsExpressionWithTypeArguments()
		return (boolToByte(n.Expression != nil) << 0) | (boolToByte(n.TypeArguments != nil) << 1)
	case ast.KindTypePredicate:
		n := node.AsTypePredicateNode()
		return (boolToByte(n.AssertsModifier != nil) << 0) | (boolToByte(n.ParameterName != nil) << 1) | (boolToByte(n.Type != nil) << 2)
	case ast.KindImportType:
		n := node.AsImportTypeNode()
		return (boolToByte(n.Argument != nil) << 0) | (boolToByte(n.Attributes != nil) << 1) | (boolToByte(n.Qualifier != nil) << 2) | (boolToByte(n.TypeArguments != nil) << 3)
	case ast.KindImportAttribute:
		n := node.AsImportAttribute()
		return (boolToByte(n.Name() != nil) << 0) | (boolToByte(n.Value != nil) << 1)
	case ast.KindTypeQuery:
		n := node.AsTypeQueryNode()
		return (boolToByte(n.ExprName != nil) << 0) | (boolToByte(n.TypeArguments != nil) << 1)
	case ast.KindMappedType:
		n := node.AsMappedTypeNode()
		return (boolToByte(n.ReadonlyToken != nil) << 0) | (boolToByte(n.TypeParameter != nil) << 1) | (boolToByte(n.NameType != nil) << 2) | (boolToByte(n.QuestionToken != nil) << 3) | (boolToByte(n.Type != nil) << 4) | (boolToByte(n.Members != nil) << 5)
	case ast.KindNamedTupleMember:
		n := node.AsNamedTupleMember()
		return (boolToByte(n.DotDotDotToken != nil) << 0) | (boolToByte(n.Name() != nil) << 1) | (boolToByte(n.QuestionToken != nil) << 2) | (boolToByte(n.Type != nil) << 3)
	case ast.KindFunctionType:
		n := node.AsFunctionTypeNode()
		return (boolToByte(n.TypeParameters != nil) << 0) | (boolToByte(n.Parameters != nil) << 1) | (boolToByte(n.Type != nil) << 2)
	case ast.KindConstructorType:
		n := node.AsConstructorTypeNode()
		return (boolToByte(n.Modifiers() != nil) << 0) | (boolToByte(n.TypeParameters != nil) << 1) | (boolToByte(n.Parameters != nil) << 2) | (boolToByte(n.Type != nil) << 3)
	case ast.KindTemplateLiteralType:
		n := node.AsTemplateLiteralTypeNode()
		return (boolToByte(n.Head != nil) << 0) | (boolToByte(n.TemplateSpans != nil) << 1)
	case ast.KindTemplateLiteralTypeSpan:
		n := node.AsTemplateLiteralTypeSpan()
		return (boolToByte(n.Type != nil) << 0) | (boolToByte(n.Literal != nil) << 1)
	case ast.KindJsxElement:
		n := node.AsJsxElement()
		return (boolToByte(n.OpeningElement != nil) << 0) | (boolToByte(n.Children != nil) << 1) | (boolToByte(n.ClosingElement != nil) << 2)
	case ast.KindJsxNamespacedName:
		n := node.AsJsxNamespacedName()
		return (boolToByte(n.Name() != nil) << 0) | (boolToByte(n.Namespace != nil) << 1)
	case ast.KindJsxOpeningElement:
		n := node.AsJsxOpeningElement()
		return (boolToByte(n.TagName != nil) << 0) | (boolToByte(n.TypeArguments != nil) << 1) | (boolToByte(n.Attributes != nil) << 2)
	case ast.KindJsxSelfClosingElement:
		n := node.AsJsxSelfClosingElement()
		return (boolToByte(n.TagName != nil) << 0) | (boolToByte(n.TypeArguments != nil) << 1) | (boolToByte(n.Attributes != nil) << 2)
	case ast.KindJsxFragment:
		n := node.AsJsxFragment()
		return (boolToByte(n.OpeningFragment != nil) << 0) | (boolToByte(n.Children != nil) << 1) | (boolToByte(n.ClosingFragment != nil) << 2)
	case ast.KindJsxAttribute:
		n := node.AsJsxAttribute()
		return (boolToByte(n.Name() != nil) << 0) | (boolToByte(n.Initializer != nil) << 1)
	case ast.KindJsxExpression:
		n := node.AsJsxExpression()
		return (boolToByte(n.DotDotDotToken != nil) << 0) | (boolToByte(n.Expression != nil) << 1)
	case ast.KindJSDoc:
		n := node.AsJSDoc()
		return (boolToByte(n.Comment != nil) << 0) | (boolToByte(n.Tags != nil) << 1)
	case ast.KindJSDocTypeTag:
		n := node.AsJSDocTypeTag()
		return (boolToByte(n.TagName != nil) << 0) | (boolToByte(n.TypeExpression != nil) << 1) | (boolToByte(n.Comment != nil) << 2)
	case ast.KindJSDocTag:
		n := node.AsJSDocUnknownTag()
		return (boolToByte(n.TagName != nil) << 0) | (boolToByte(n.Comment != nil) << 1)
	case ast.KindJSDocTemplateTag:
		n := node.AsJSDocTemplateTag()
		return (boolToByte(n.TagName != nil) << 0) | (boolToByte(n.Constraint != nil) << 1) | (boolToByte(n.TypeParameters() != nil) << 2) | (boolToByte(n.Comment != nil) << 3)
	case ast.KindJSDocReturnTag:
		n := node.AsJSDocReturnTag()
		return (boolToByte(n.TagName != nil) << 0) | (boolToByte(n.TypeExpression != nil) << 1) | (boolToByte(n.Comment != nil) << 2)
	case ast.KindJSDocPublicTag:
		n := node.AsJSDocPublicTag()
		return (boolToByte(n.TagName != nil) << 0) | (boolToByte(n.Comment != nil) << 1)
	case ast.KindJSDocPrivateTag:
		n := node.AsJSDocPrivateTag()
		return (boolToByte(n.TagName != nil) << 0) | (boolToByte(n.Comment != nil) << 1)
	case ast.KindJSDocProtectedTag:
		n := node.AsJSDocProtectedTag()
		return (boolToByte(n.TagName != nil) << 0) | (boolToByte(n.Comment != nil) << 1)
	case ast.KindJSDocReadonlyTag:
		n := node.AsJSDocReadonlyTag()
		return (boolToByte(n.TagName != nil) << 0) | (boolToByte(n.Comment != nil) << 1)
	case ast.KindJSDocOverrideTag:
		n := node.AsJSDocOverrideTag()
		return (boolToByte(n.TagName != nil) << 0) | (boolToByte(n.Comment != nil) << 1)
	case ast.KindJSDocDeprecatedTag:
		n := node.AsJSDocDeprecatedTag()
		return (boolToByte(n.TagName != nil) << 0) | (boolToByte(n.Comment != nil) << 1)
	case ast.KindJSDocSeeTag:
		n := node.AsJSDocSeeTag()
		return (boolToByte(n.TagName != nil) << 0) | (boolToByte(n.NameExpression != nil) << 1) | (boolToByte(n.Comment != nil) << 2)
	case ast.KindJSDocImplementsTag:
		n := node.AsJSDocImplementsTag()
		return (boolToByte(n.TagName != nil) << 0) | (boolToByte(n.ClassName != nil) << 1) | (boolToByte(n.Comment != nil) << 2)
	case ast.KindJSDocAugmentsTag:
		n := node.AsJSDocAugmentsTag()
		return (boolToByte(n.TagName != nil) << 0) | (boolToByte(n.ClassName != nil) << 1) | (boolToByte(n.Comment != nil) << 2)
	case ast.KindJSDocSatisfiesTag:
		n := node.AsJSDocSatisfiesTag()
		return (boolToByte(n.TagName != nil) << 0) | (boolToByte(n.TypeExpression != nil) << 1) | (boolToByte(n.Comment != nil) << 2)
	case ast.KindJSDocThisTag:
		n := node.AsJSDocThisTag()
		return (boolToByte(n.TagName != nil) << 0) | (boolToByte(n.TypeExpression != nil) << 1) | (boolToByte(n.Comment != nil) << 2)
	case ast.KindJSDocImportTag:
		n := node.AsJSDocImportTag()
		return (boolToByte(n.TagName != nil) << 0) | (boolToByte(n.ImportClause != nil) << 1) | (boolToByte(n.ModuleSpecifier != nil) << 2) | (boolToByte(n.Attributes != nil) << 3) | (boolToByte(n.Comment != nil) << 4)
	case ast.KindJSDocCallbackTag:
		n := node.AsJSDocCallbackTag()
		return (boolToByte(n.TagName != nil) << 0) | (boolToByte(n.TypeExpression != nil) << 1) | (boolToByte(n.FullName != nil) << 2) | (boolToByte(n.Comment != nil) << 3)
	case ast.KindJSDocOverloadTag:
		n := node.AsJSDocOverloadTag()
		return (boolToByte(n.TagName != nil) << 0) | (boolToByte(n.TypeExpression != nil) << 1) | (boolToByte(n.Comment != nil) << 2)
	case ast.KindJSDocTypedefTag:
		n := node.AsJSDocTypedefTag()
		return (boolToByte(n.TagName != nil) << 0) | (boolToByte(n.TypeExpression != nil) << 1) | (boolToByte(n.FullName != nil) << 2) | (boolToByte(n.Comment != nil) << 3)
	case ast.KindJSDocSignature:
		n := node.AsJSDocSignature()
		return (boolToByte(n.TypeParameters() != nil) << 0) | (boolToByte(n.Parameters != nil) << 1) | (boolToByte(n.Type != nil) << 2)
	case ast.KindClassStaticBlockDeclaration:
		n := node.AsClassStaticBlockDeclaration()
		return (boolToByte(n.Modifiers() != nil) << 0) | (boolToByte(n.Body != nil) << 1)
	case ast.KindClassDeclaration:
		n := node.AsClassDeclaration()
		return (boolToByte(n.Modifiers() != nil) << 0) | (boolToByte(n.Name() != nil) << 1) | (boolToByte(n.TypeParameters != nil) << 2) | (boolToByte(n.HeritageClauses != nil) << 3) | (boolToByte(n.Members != nil) << 4)
	case ast.KindJSDocPropertyTag:
		n := node.AsJSDocPropertyTag()
		if n.IsNameFirst {
			return (boolToByte(n.Name() != nil) << 0) | (boolToByte(n.TypeExpression != nil) << 1)
		}
		return (boolToByte(n.TypeExpression != nil) << 0) | (boolToByte(n.Name() != nil) << 1)
	case ast.KindJSDocParameterTag:
		n := node.AsJSDocParameterTag()
		if n.IsNameFirst {
			return (boolToByte(n.TagName != nil) << 0) | (boolToByte(n.Name() != nil) << 1) | (boolToByte(n.TypeExpression != nil) << 2) | (boolToByte(n.Comment != nil) << 3)
		}
		return (boolToByte(n.TagName != nil) << 0) | (boolToByte(n.TypeExpression != nil) << 1) | (boolToByte(n.Name() != nil) << 2) | (boolToByte(n.Comment != nil) << 3)
	default:
		return 0
	}
}

func getNodeDefinedData(node *ast.Node) int32 {
	switch node.Kind {
	case ast.KindJSDocTypeLiteral:
		n := node.AsJSDocTypeLiteral()
		return int32(boolToByte(n.IsArrayType)) << 24
	case ast.KindImportSpecifier:
		n := node.AsImportSpecifier()
		return int32(boolToByte(n.IsTypeOnly)) << 24
	case ast.KindImportClause:
		n := node.AsImportClause()
		return int32(boolToByte(n.IsTypeOnly)) << 24
	case ast.KindExportSpecifier:
		n := node.AsExportSpecifier()
		return int32(boolToByte(n.IsTypeOnly)) << 24
	case ast.KindImportType:
		n := node.AsImportTypeNode()
		return int32(boolToByte(n.IsTypeOf)) << 24
	case ast.KindBlock:
		n := node.AsBlock()
		return int32(boolToByte(n.Multiline)) << 24
	case ast.KindImportEqualsDeclaration:
		n := node.AsImportEqualsDeclaration()
		return int32(boolToByte(n.IsTypeOnly)) << 24
	case ast.KindExportAssignment:
		n := node.AsExportAssignment()
		return int32(boolToByte(n.IsExportEquals)) << 24
	case ast.KindExportDeclaration:
		n := node.AsExportDeclaration()
		return int32(boolToByte(n.IsTypeOnly)) << 24
	case ast.KindArrayLiteralExpression:
		n := node.AsArrayLiteralExpression()
		return int32(boolToByte(n.MultiLine)) << 24
	case ast.KindObjectLiteralExpression:
		n := node.AsObjectLiteralExpression()
		return int32(boolToByte(n.MultiLine)) << 24
	case ast.KindJSDocPropertyTag:
		n := node.AsJSDocPropertyTag()
		return int32(boolToByte(n.IsBracketed))<<24 | int32(boolToByte(n.IsNameFirst))<<25
	case ast.KindJSDocParameterTag:
		n := node.AsJSDocParameterTag()
		return int32(boolToByte(n.IsBracketed))<<24 | int32(boolToByte(n.IsNameFirst))<<25
	case ast.KindJsxText:
		n := node.AsJsxText()
		return int32(boolToByte(n.ContainsOnlyTriviaWhiteSpaces)) << 24
	case ast.KindVariableDeclarationList:
		n := node.AsVariableDeclarationList()
		return int32(n.Flags & (ast.NodeFlagsLet | ast.NodeFlagsConst) << 24)
	case ast.KindImportAttributes:
		n := node.AsImportAttributes()
		return int32(boolToByte(n.MultiLine))<<24 | int32(boolToByte(n.Token == ast.KindAssertKeyword))<<25
	}
	return 0
}

func recordNodeStrings(node *ast.Node, strs *stringTable) int32 {
	switch node.Kind {
	case ast.KindJsxText:
		return strs.add(node.AsJsxText().Text)
	case ast.KindIdentifier:
		return strs.add(node.AsIdentifier().Text)
	case ast.KindPrivateIdentifier:
		return strs.add(node.AsPrivateIdentifier().Text)
	case ast.KindStringLiteral:
		return strs.add(node.AsStringLiteral().Text)
	case ast.KindNumericLiteral:
		return strs.add(node.AsNumericLiteral().Text)
	case ast.KindBigIntLiteral:
		return strs.add(node.AsBigIntLiteral().Text)
	case ast.KindRegularExpressionLiteral:
		return strs.add(node.AsRegularExpressionLiteral().Text)
	case ast.KindNoSubstitutionTemplateLiteral:
		return strs.add(node.AsNoSubstitutionTemplateLiteral().Text)
	case ast.KindJSDocText:
		return strs.add(node.AsJSDocText().Text)
	default:
		panic(fmt.Sprintf("Unexpected node kind %v", node.Kind))
	}
}

func boolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

func EncodeSymbolResponse(symbol *ast.Symbol, handle Handle[ast.Symbol]) ([]byte, error) {
	var err error
	result := make([]byte, 0, 12+len(symbol.Name))
	if result, err = appendInt32s(result, int32(handle), int32(symbol.Flags), int32(symbol.CheckFlags)); err != nil {
		return nil, err
	}
	result = append(result, symbol.Name...)
	return result, nil
}

func EncodeTypeResponse(t *checker.Type, handle Handle[checker.Type]) ([]byte, error) {
	var err error
	result := make([]byte, 0, 8)
	if result, err = appendInt32s(result, int32(handle), int32(t.Flags())); err != nil {
		return nil, err
	}
	return result, nil
}
