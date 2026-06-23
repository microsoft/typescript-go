package ls

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/scanner"
)

// JSDocTagInfo mirrors Strada's `JSDocTagInfo`, but renders the tag's text as a
// plain string instead of `SymbolDisplayPart[]`.
type JSDocTagInfo struct {
	Name string
	Text string
}

// GetSymbolDocumentationComment renders a symbol's documentation comment (the leading
// JSDoc text, without tags) as plain text. It backs the API's Symbol.getDocumentationComment.
func (l *LanguageService) GetSymbolDocumentationComment(c *checker.Checker, symbol *ast.Symbol) string {
	if symbol == nil {
		return ""
	}
	for _, decl := range symbol.Declarations {
		if doc := l.getDocumentationFromDeclaration(c, symbol, decl, decl, lsproto.MarkupKindPlainText, true /*commentOnly*/); doc != "" {
			return doc
		}
	}
	if symbol.Flags&ast.SymbolFlagsAlias != 0 {
		aliased := c.GetAliasedSymbol(symbol)
		if aliased != nil && aliased != c.GetUnknownSymbol() {
			for _, decl := range aliased.Declarations {
				if doc := l.getDocumentationFromDeclaration(c, aliased, decl, decl, lsproto.MarkupKindPlainText, true /*commentOnly*/); doc != "" {
					return doc
				}
			}
		}
	}
	return ""
}

// GetSymbolJSDocTags collects a symbol's JSDoc tags. It backs the API's Symbol.getJsDocTags
// and mirrors Strada's getJsDocTagsFromDeclarations, except each tag's text is rendered as a
// plain string rather than SymbolDisplayPart[]. Tags with no text have an empty Text field.
func (l *LanguageService) GetSymbolJSDocTags(symbol *ast.Symbol) []JSDocTagInfo {
	if symbol == nil {
		return nil
	}
	var infos []JSDocTagInfo
	seen := make(map[*ast.Node]struct{}, len(symbol.Declarations))
	for _, decl := range symbol.Declarations {
		if decl == nil {
			continue
		}
		if _, ok := seen[decl]; ok {
			continue
		}
		seen[decl] = struct{}{}
		tags := declarationJSDocTags(decl)
		// Skip comments containing @typedef/@callback since they're not associated with a
		// particular declaration, unless they also carry @param/@return (treated as local docs).
		hasTypedef := core.Some(tags, func(t *ast.Node) bool {
			return t.Kind == ast.KindJSDocTypedefTag || t.Kind == ast.KindJSDocCallbackTag
		})
		hasParamOrReturn := core.Some(tags, func(t *ast.Node) bool {
			return t.Kind == ast.KindJSDocParameterTag || t.Kind == ast.KindJSDocReturnTag
		})
		if hasTypedef && !hasParamOrReturn {
			continue
		}
		for _, tag := range tags {
			infos = append(infos, JSDocTagInfo{Name: tag.TagName().Text(), Text: getJSDocTagText(tag)})
		}
	}
	return infos
}

// declarationJSDocTags returns the JSDoc tags associated with a declaration, walking the
// JSDoc comment location chain like the checker's getAllJSDocTags.
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

// getJSDocTagText renders the text of a single JSDoc tag as a plain string, mirroring
// Strada's getCommentDisplayParts collapsed from SymbolDisplayPart[] to a string.
func getJSDocTagText(tag *ast.Node) string {
	comment := scanner.GetTextOfJSDocComment(tag.CommentList())
	addComment := func(s string) string {
		if comment == "" {
			return s
		}
		return s + " " + comment
	}
	switch tag.Kind {
	case ast.KindJSDocThrowsTag:
		if te := tag.AsJSDocThrowsTag().TypeExpression; te != nil {
			return addComment(scanner.GetTextOfNode(te))
		}
		return comment
	case ast.KindJSDocImplementsTag:
		return addComment(scanner.GetTextOfNode(tag.AsJSDocImplementsTag().ClassName))
	case ast.KindJSDocAugmentsTag:
		return addComment(scanner.GetTextOfNode(tag.AsJSDocAugmentsTag().ClassName))
	case ast.KindJSDocTemplateTag:
		templateTag := tag.AsJSDocTemplateTag()
		var b strings.Builder
		if templateTag.Constraint != nil {
			b.WriteString(scanner.GetTextOfNode(templateTag.Constraint))
		}
		if templateTag.TypeParameters != nil {
			for i, tp := range templateTag.TypeParameters.Nodes {
				if i == 0 && b.Len() != 0 {
					b.WriteString(" ")
				}
				if i != 0 {
					b.WriteString(", ")
				}
				b.WriteString(scanner.GetTextOfNode(tp))
			}
		}
		if comment != "" {
			if b.Len() != 0 {
				b.WriteString(" ")
			}
			b.WriteString(comment)
		}
		return b.String()
	case ast.KindJSDocTypeTag:
		return addComment(scanner.GetTextOfNode(tag.AsJSDocTypeTag().TypeExpression))
	case ast.KindJSDocSatisfiesTag:
		return addComment(scanner.GetTextOfNode(tag.AsJSDocSatisfiesTag().TypeExpression))
	case ast.KindJSDocSeeTag:
		if ne := tag.AsJSDocSeeTag().NameExpression; ne != nil {
			return addComment(scanner.GetTextOfNode(ne))
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
