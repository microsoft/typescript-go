package jsdoc

import (
	"slices"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/scanner"
)

// TagInfo is a JSDoc tag with its text rendered as a plain string.
type TagInfo struct {
	Name string
	Text string
}

// GetSymbolDocumentationComment renders a symbol's documentation comment as plain text.
func GetSymbolDocumentationComment(symbol *ast.Symbol) string {
	if symbol == nil {
		return ""
	}
	var parts []string
	var seen collections.Set[*ast.Node]
	for _, declaration := range symbol.Declarations {
		if declaration == nil || !seen.AddIfAbsent(declaration) {
			continue
		}
		for _, comment := range getJSDocComments(declaration) {
			if shouldSkipComment(declaration, comment) {
				continue
			}
			text := renderComments(comment.Comments())
			if text != "" && !slices.Contains(parts, text) {
				parts = append(parts, text)
			}
		}
	}
	return strings.Join(parts, "\n")
}

// GetSymbolTags collects a symbol's JSDoc tags with their text rendered as plain strings.
func GetSymbolTags(symbol *ast.Symbol) []TagInfo {
	if symbol == nil {
		return nil
	}
	var infos []TagInfo
	var seen collections.Set[*ast.Node]
	for _, declaration := range symbol.Declarations {
		if declaration == nil || !seen.AddIfAbsent(declaration) {
			continue
		}
		tags := declarationJSDocTags(declaration)
		hasTypedef := core.Some(tags, func(tag *ast.Node) bool {
			return tag.Kind == ast.KindJSDocTypedefTag || tag.Kind == ast.KindJSDocCallbackTag
		})
		hasParamOrReturn := core.Some(tags, func(tag *ast.Node) bool {
			return tag.Kind == ast.KindJSDocParameterTag || tag.Kind == ast.KindJSDocReturnTag
		})
		if hasTypedef && !hasParamOrReturn {
			continue
		}
		for _, tag := range tags {
			infos = append(infos, TagInfo{Name: tag.TagName().Text(), Text: getTagText(tag)})
		}
	}
	return infos
}

func getJSDocComments(node *ast.Node) []*ast.Node {
	if node.Flags&ast.NodeFlagsJSDoc != 0 {
		return []*ast.Node{node}
	}
	var comments []*ast.Node
	for current := node; current != nil; current = ast.GetNextJSDocCommentLocation(current) {
		comments = append(comments, current.JSDoc(nil)...)
	}
	return comments
}

func shouldSkipComment(declaration *ast.Node, comment *ast.Node) bool {
	if comment == nil || len(comment.Comments()) == 0 {
		return true
	}
	if declaration.Kind == ast.KindJSDocTypedefTag || declaration.Kind == ast.KindJSDocCallbackTag || comment.Kind != ast.KindJSDoc {
		return false
	}
	tags := comment.AsJSDoc().Tags
	if tags == nil {
		return false
	}
	hasTypedef := core.Some(tags.Nodes, func(tag *ast.Node) bool {
		return tag.Kind == ast.KindJSDocTypedefTag || tag.Kind == ast.KindJSDocCallbackTag
	})
	hasParamOrReturn := core.Some(tags.Nodes, func(tag *ast.Node) bool {
		return tag.Kind == ast.KindJSDocParameterTag || tag.Kind == ast.KindJSDocReturnTag
	})
	return hasTypedef && !hasParamOrReturn
}

func renderComments(comments []*ast.Node) string {
	var builder strings.Builder
	for _, comment := range comments {
		switch comment.Kind {
		case ast.KindJSDocText:
			builder.WriteString(comment.Text())
		case ast.KindJSDocLink, ast.KindJSDocLinkPlain, ast.KindJSDocLinkCode:
			name := comment.Name()
			text := strings.Trim(comment.Text(), " ")
			if name == nil {
				builder.WriteString(text)
			} else if text == "" {
				builder.WriteString(scanner.GetTextOfNode(name))
			} else {
				builder.WriteString(strings.TrimLeft(strings.TrimPrefix(strings.TrimLeft(text, " "), "|"), " "))
			}
		}
	}
	return builder.String()
}

func declarationJSDocTags(node *ast.Node) []*ast.Node {
	if node.Flags&ast.NodeFlagsJSDoc == 0 {
		for current := node; current != nil; current = ast.GetNextJSDocCommentLocation(current) {
			jsdocs := current.JSDoc(nil)
			if len(jsdocs) == 0 {
				continue
			}
			lastJSDoc := jsdocs[len(jsdocs)-1].AsJSDoc()
			if lastJSDoc.Tags != nil {
				return lastJSDoc.Tags.Nodes
			}
		}
	}
	return nil
}

func getTagText(tag *ast.Node) string {
	comment := scanner.GetTextOfJSDocComment(tag.CommentList())
	addComment := func(text string) string {
		if comment == "" {
			return text
		}
		return text + " " + comment
	}
	switch tag.Kind {
	case ast.KindJSDocThrowsTag:
		if typeExpression := tag.AsJSDocThrowsTag().TypeExpression; typeExpression != nil {
			return addComment(scanner.GetTextOfNode(typeExpression))
		}
		return comment
	case ast.KindJSDocImplementsTag:
		return addComment(scanner.GetTextOfNode(tag.AsJSDocImplementsTag().ClassName))
	case ast.KindJSDocAugmentsTag:
		return addComment(scanner.GetTextOfNode(tag.AsJSDocAugmentsTag().ClassName))
	case ast.KindJSDocTemplateTag:
		templateTag := tag.AsJSDocTemplateTag()
		var builder strings.Builder
		if templateTag.Constraint != nil {
			builder.WriteString(scanner.GetTextOfNode(templateTag.Constraint))
		}
		if templateTag.TypeParameters != nil {
			for i, typeParameter := range templateTag.TypeParameters.Nodes {
				if i == 0 && builder.Len() != 0 {
					builder.WriteString(" ")
				}
				if i != 0 {
					builder.WriteString(", ")
				}
				builder.WriteString(scanner.GetTextOfNode(typeParameter))
			}
		}
		if comment != "" {
			if builder.Len() != 0 {
				builder.WriteString(" ")
			}
			builder.WriteString(comment)
		}
		return builder.String()
	case ast.KindJSDocTypeTag:
		return addComment(scanner.GetTextOfNode(tag.AsJSDocTypeTag().TypeExpression))
	case ast.KindJSDocSatisfiesTag:
		return addComment(scanner.GetTextOfNode(tag.AsJSDocSatisfiesTag().TypeExpression))
	case ast.KindJSDocSeeTag:
		if nameExpression := tag.AsJSDocSeeTag().NameExpression; nameExpression != nil {
			return addComment(scanner.GetTextOfNode(nameExpression))
		}
		return comment
	case ast.KindJSDocParameterTag, ast.KindJSDocPropertyTag:
		if name := tag.Name(); name != nil {
			return addComment(scanner.GetTextOfNode(name))
		}
		return comment
	default:
		return comment
	}
}
