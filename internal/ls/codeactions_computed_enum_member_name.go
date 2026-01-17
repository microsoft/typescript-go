package ls

import (
	"context"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/locale"
	"github.com/microsoft/typescript-go/internal/ls/change"
)

var computedEnumMemberNameFixErrorCodes = []int32{
	diagnostics.Using_a_string_literal_as_an_enum_member_name_via_a_computed_property_is_deprecated_Use_a_simple_string_literal_instead.Code(),
}

const computedEnumMemberNameFixID = "convertComputedEnumMemberName"

var computedEnumMemberNameFixProvider = &CodeFixProvider{
	ErrorCodes:     computedEnumMemberNameFixErrorCodes,
	GetCodeActions: getComputedEnumMemberNameCodeActions,
	FixIds:         []string{computedEnumMemberNameFixID},
}

type computedEnumMemberNameInfo struct {
	computedName *ast.Node
	literalText  string
}

func getComputedEnumMemberNameCodeActions(ctx context.Context, fixContext *CodeFixContext) ([]CodeAction, error) {
	info := getComputedEnumMemberNameInfo(fixContext.SourceFile, fixContext.Span.Pos())
	if info == nil {
		return nil, nil
	}

	tracker := change.NewTracker(ctx, fixContext.Program.Options(), fixContext.LS.FormatOptions(), fixContext.LS.converters)
	tracker.ReplaceNode(
		fixContext.SourceFile,
		info.computedName,
		tracker.NodeFactory.NewStringLiteral(info.literalText, ast.TokenFlagsNone),
		nil,
	)

	locale := locale.FromContext(ctx)
	description := diagnostics.Remove_unnecessary_computed_property_name_syntax.Localize(locale)
	return []CodeAction{{
		Description: description,
		Changes:     tracker.GetChanges()[fixContext.SourceFile.FileName()],
	}}, nil
}

func getComputedEnumMemberNameInfo(sourceFile *ast.SourceFile, pos int) *computedEnumMemberNameInfo {
	token := astnav.GetTokenAtPosition(sourceFile, pos)
	node := token
	for node != nil && !ast.IsComputedPropertyName(node) {
		node = node.Parent
	}

	if node == nil || !ast.IsComputedPropertyName(node) {
		return nil
	}
	if node.Parent == nil || !ast.IsEnumMember(node.Parent) {
		return nil
	}

	expression := node.AsComputedPropertyName().Expression
	if !ast.IsStringLiteralLike(expression) {
		return nil
	}

	return &computedEnumMemberNameInfo{
		computedName: node,
		literalText:  expression.Text(),
	}
}
