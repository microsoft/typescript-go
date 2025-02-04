package printer

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"sync/atomic"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
)

// Stores side-table information used during transformation that can be read by the printer to customize emit
//
// NOTE: EmitContext is not guaranteed to be thread-safe.
type EmitContext struct {
	Factory      *ast.NodeFactory // Required. The NodeFactory to use for creating new nodes
	autoGenerate map[*ast.MemberName]*autoGenerateInfo
	textSource   map[*ast.StringLiteralNode]*ast.Node
	original     map[*ast.Node]*ast.Node
	emitNodes    core.LinkStore[*ast.Node, emitNode]
}

func NewEmitContext() *EmitContext {
	return &EmitContext{Factory: &ast.NodeFactory{}}
}

type AutoGenerateOptions struct {
	Flags  GeneratedIdentifierFlags
	Prefix string
	Suffix string
}

func (c *EmitContext) newGeneratedIdentifier(kind GeneratedIdentifierFlags, text string, options AutoGenerateOptions) (*ast.IdentifierNode, *autoGenerateInfo) {
	name := c.Factory.NewIdentifier(text)
	autoGenerate := &autoGenerateInfo{
		Id:     autoGenerateId(nextAutoGenerateId.Add(1)),
		Flags:  kind | (options.Flags & ^GeneratedIdentifierFlagsKindMask),
		Prefix: options.Prefix,
		Suffix: options.Suffix,
	}
	if c.autoGenerate == nil {
		c.autoGenerate = make(map[*ast.MemberName]*autoGenerateInfo)
	}
	c.autoGenerate[name] = autoGenerate
	return name, autoGenerate
}

func (c *EmitContext) NewTempVariable(options AutoGenerateOptions) *ast.IdentifierNode {
	name, _ := c.newGeneratedIdentifier(GeneratedIdentifierFlagsAuto, "", options)
	return name
}

func (c *EmitContext) NewLoopVariable(options AutoGenerateOptions) *ast.IdentifierNode {
	name, _ := c.newGeneratedIdentifier(GeneratedIdentifierFlagsLoop, "", options)
	return name
}

func (c *EmitContext) NewUniqueName(text string, options AutoGenerateOptions) *ast.IdentifierNode {
	name, _ := c.newGeneratedIdentifier(GeneratedIdentifierFlagsUnique, text, options)
	return name
}

func (c *EmitContext) NewGeneratedNameForNode(node *ast.Node, options AutoGenerateOptions) *ast.IdentifierNode {
	if len(options.Prefix) > 0 || len(options.Suffix) > 0 {
		options.Flags |= GeneratedIdentifierFlagsOptimistic
	}

	// For debugging purposes, set the `text` of the identifier to a reasonable value
	var text string
	switch {
	case node == nil:
		text = formatGeneratedName(false /*privateName*/, options.Prefix, "" /*base*/, options.Suffix)
	case ast.IsMemberName(node):
		text = formatGeneratedName(false, options.Prefix, node.Text(), options.Suffix)
	default:
		text = fmt.Sprintf("generated@%v", ast.GetNodeId(node))
	}

	name, autoGenerate := c.newGeneratedIdentifier(GeneratedIdentifierFlagsNode, text, options)
	autoGenerate.Node = node
	return name
}

func (c *EmitContext) newGeneratedPrivateIdentifier(kind GeneratedIdentifierFlags, text string, options AutoGenerateOptions) (*ast.PrivateIdentifierNode, *autoGenerateInfo) {
	if !strings.HasPrefix(text, "#") {
		panic("First character of private identifier must be #: " + text)
	}

	name := c.Factory.NewPrivateIdentifier(text)
	autoGenerate := &autoGenerateInfo{
		Id:     autoGenerateId(nextAutoGenerateId.Add(1)),
		Flags:  kind | (options.Flags &^ GeneratedIdentifierFlagsKindMask),
		Prefix: options.Prefix,
		Suffix: options.Suffix,
	}
	if c.autoGenerate == nil {
		c.autoGenerate = make(map[*ast.MemberName]*autoGenerateInfo)
	}
	c.autoGenerate[name] = autoGenerate
	return name, autoGenerate
}

func (c *EmitContext) NewUniquePrivateName(text string, options AutoGenerateOptions) *ast.PrivateIdentifierNode {
	name, _ := c.newGeneratedPrivateIdentifier(GeneratedIdentifierFlagsUnique, text, options)
	return name
}

func (c *EmitContext) NewGeneratedPrivateNameForNode(node *ast.Node, options AutoGenerateOptions) *ast.PrivateIdentifierNode {
	if len(options.Prefix) > 0 || len(options.Suffix) > 0 {
		options.Flags |= GeneratedIdentifierFlagsOptimistic
	}

	var text string
	switch {
	case node == nil:
		text = formatGeneratedName(true /*privateName*/, options.Prefix, "" /*base*/, options.Suffix)
	case ast.IsMemberName(node):
		text = formatGeneratedName(true /*privateName*/, options.Prefix, node.Text(), options.Suffix)
	default:
		text = fmt.Sprintf("#generated@%v", ast.GetNodeId(node))
	}

	name, autoGenerate := c.newGeneratedPrivateIdentifier(GeneratedIdentifierFlagsNode, text, options)
	autoGenerate.Node = node
	return name
}

func (c *EmitContext) HasAutoGenerateInfo(node *ast.MemberName) bool {
	if node != nil {
		_, ok := c.autoGenerate[node]
		return ok
	}
	return false
}

var nextAutoGenerateId atomic.Uint32

type autoGenerateId uint32

type autoGenerateInfo struct {
	Flags  GeneratedIdentifierFlags // Specifies whether to auto-generate the text for an identifier.
	Id     autoGenerateId           // Ensures unique generated identifiers get unique names, but clones get the same name.
	Prefix string                   // Optional prefix to apply to the start of the generated name
	Suffix string                   // Optional suffix to apply to the end of the generated name
	Node   *ast.Node                // For a GeneratedIdentifierFlagsNode, the node from which to generate an identifier
}

func (info *autoGenerateInfo) Kind() GeneratedIdentifierFlags {
	return info.Flags & GeneratedIdentifierFlagsKindMask
}

func (c *EmitContext) NewStringLiteralFromNode(textSourceNode *ast.Node) *ast.Node {
	var text string
	if ast.IsMemberName(textSourceNode) || ast.IsJsxNamespacedName(textSourceNode) {
		text = textSourceNode.Text()
	}
	node := c.Factory.NewStringLiteral(text)
	if c.textSource == nil {
		c.textSource = make(map[*ast.StringLiteralNode]*ast.Node)
	}
	c.textSource[node] = textSourceNode
	return node
}

func (c *EmitContext) SetOriginal(node *ast.Node, original *ast.Node) {
	if original == nil {
		panic("Original cannot be nil.")
	}

	if c.original == nil {
		c.original = make(map[*ast.Node]*ast.Node)
	}

	existing, ok := c.original[node]
	if !ok {
		c.original[node] = original
		if emitNode := c.emitNodes.TryGet(node); emitNode != nil {
			c.emitNodes.Get(node).copyFrom(emitNode)
		}
	} else if existing != original {
		panic("Original node already set.")
	}
}

func (c *EmitContext) Original(node *ast.Node) *ast.Node {
	return c.original[node]
}

// Gets the most original node associated with this node by walking Original pointers.
//
// This method is analogous to `getOriginalNode` in the old compiler, but the name has changed to avoid accidental
// conflation with `SetOriginal`/`Original`
func (c *EmitContext) MostOriginal(node *ast.Node) *ast.Node {
	if node != nil {
		original := c.Original(node)
		for original != nil {
			node = original
			original = c.Original(node)
		}
	}
	return node
}

type emitNode struct {
	emitFlags            EmitFlags
	commentRange         *core.TextRange
	sourceMapRange       *core.TextRange
	tokenSourceMapRanges map[ast.Kind]*core.TextRange
}

// NOTE: This method is not guaranteed to be thread-safe
func (e *emitNode) copyFrom(source *emitNode) {
	e.emitFlags = source.emitFlags
	e.commentRange = source.commentRange
	e.sourceMapRange = source.sourceMapRange
	e.tokenSourceMapRanges = maps.Clone(source.tokenSourceMapRanges)
}

func (c *EmitContext) EmitFlags(node *ast.Node) EmitFlags {
	if emitNode := c.emitNodes.TryGet(node); emitNode != nil {
		return emitNode.emitFlags
	}
	return EFNone
}

func (c *EmitContext) SetEmitFlags(node *ast.Node, flags EmitFlags) {
	c.emitNodes.Get(node).emitFlags = flags
}

func (c *EmitContext) AddEmitFlags(node *ast.Node, flags EmitFlags) {
	c.emitNodes.Get(node).emitFlags |= flags
}

// Sets the range to use for a node when emitting comments and source maps.
func (c *EmitContext) SetCommentAndSourceMapRanges(node *ast.Node, loc core.TextRange) {
	emitNode := c.emitNodes.Get(node)
	emitNode.commentRange = &loc
	emitNode.sourceMapRange = &loc
}

// Gets the range to use for a node when emitting comments.
func (c *EmitContext) CommentRange(node *ast.Node) core.TextRange {
	if emitNode := c.emitNodes.TryGet(node); emitNode != nil && emitNode.commentRange != nil {
		return *emitNode.commentRange
	}
	return node.Loc
}

// Sets the range to use for a node when emitting comments.
func (c *EmitContext) SetCommentRange(node *ast.Node, loc core.TextRange) {
	c.emitNodes.Get(node).commentRange = &loc
}

// Gets the range to use for a node when emitting source maps.
func (c *EmitContext) SourceMapRange(node *ast.Node) core.TextRange {
	if emitNode := c.emitNodes.TryGet(node); emitNode != nil && emitNode.sourceMapRange != nil {
		return *emitNode.sourceMapRange
	}
	return node.Loc
}

// Sets the range to use for a node when emitting source maps.
func (c *EmitContext) SetSourceMapRange(node *ast.Node, loc core.TextRange) {
	c.emitNodes.Get(node).sourceMapRange = &loc
}

// Gets the range for a token of a node when emitting source maps.
func (c *EmitContext) TokenSourceMapRange(node *ast.Node, kind ast.Kind) *core.TextRange {
	if emitNode := c.emitNodes.TryGet(node); emitNode != nil {
		if emitNode.tokenSourceMapRanges == nil {
			return nil
		}
		if loc, ok := emitNode.tokenSourceMapRanges[kind]; ok {
			return loc
		}
	}
	return nil
}

// Sets the range for a token of a node when emitting source maps.
func (c *EmitContext) SetTokenSourceMapRange(node *ast.Node, kind ast.Kind, loc *core.TextRange) {
	emitNode := c.emitNodes.Get(node)
	if emitNode.tokenSourceMapRanges == nil {
		emitNode.tokenSourceMapRanges = make(map[ast.Kind]*core.TextRange)
	}
	emitNode.tokenSourceMapRanges[kind] = loc
}

