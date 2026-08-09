package ls

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/parser"
)

func TestConvertRefactorToLSPCodeAction_Rename(t *testing.T) {
	t.Parallel()

	text := "const x: { a: number } = { a: 1 };\n"
	sourceFile := parser.ParseSourceFile(ast.SourceFileParseOptions{
		FileName: "/index.ts",
		Path:     "/index.ts",
	}, text, core.ScriptKindTS)

	converters := lsconv.NewConverters(lsproto.PositionEncodingKindUTF16, func(string) *lsconv.LSPLineMap {
		return lsconv.ComputeLSPLineStarts(text)
	})
	l := &LanguageService{converters: converters}

	action := &CodeAction{
		Description: "Extract to type alias",
		Kind:        lsproto.CodeActionKindRefactorExtract,
		Changes: []*lsproto.TextEdit{
			{Range: lsproto.Range{Start: lsproto.Position{Line: 0, Character: 0}, End: lsproto.Position{Line: 0, Character: 0}}, NewText: "type NewType = { a: number };\n"},
			{Range: lsproto.Range{Start: lsproto.Position{Line: 0, Character: 10}, End: lsproto.Position{Line: 0, Character: 23}}, NewText: "NewType"},
		},
		RenameFilename: "/index.ts",
		RenameLocation: 5,
	}

	converted := convertRefactorToLSPCodeAction(l, sourceFile, action, "file:///index.ts")
	codeAction := converted.CodeAction

	switch {
	case codeAction.Command == nil:
		t.Fatal("expected a rename command to be attached to the code action")
	case codeAction.Command.Command != "editor.action.rename":
		t.Errorf("expected command %q, got %q", "editor.action.rename", codeAction.Command.Command)
	case codeAction.Command.Arguments == nil || len(*codeAction.Command.Arguments) != 1:
		t.Fatalf("expected a single argument, got %v", codeAction.Command.Arguments)
	}

	args, ok := (*codeAction.Command.Arguments)[0].([]any)
	if !ok || len(args) != 2 {
		t.Fatalf("expected argument to be [uri, position], got %v", (*codeAction.Command.Arguments)[0])
	}
	if got, ok := args[0].(lsproto.DocumentUri); !ok || got != "file:///index.ts" {
		t.Errorf("expected rename uri %q, got %v", "file:///index.ts", args[0])
	}

	expectedPos := lsproto.Position{Line: 0, Character: 5}
	if args[1] != expectedPos {
		t.Errorf("expected rename position %v, got %v", expectedPos, args[1])
	}
}

func TestConvertRefactorToLSPCodeAction_NoRename(t *testing.T) {
	t.Parallel()

	text := "const x: { a: number } = { a: 1 };\n"
	sourceFile := parser.ParseSourceFile(ast.SourceFileParseOptions{
		FileName: "/index.ts",
		Path:     "/index.ts",
	}, text, core.ScriptKindTS)

	converters := lsconv.NewConverters(lsproto.PositionEncodingKindUTF16, func(string) *lsconv.LSPLineMap {
		return lsconv.ComputeLSPLineStarts(text)
	})
	l := &LanguageService{converters: converters}

	action := &CodeAction{
		Description: "Infer return type",
		Changes: []*lsproto.TextEdit{
			{Range: lsproto.Range{Start: lsproto.Position{Line: 0, Character: 10}, End: lsproto.Position{Line: 0, Character: 23}}, NewText: "NewType"},
		},
	}

	converted := convertRefactorToLSPCodeAction(l, sourceFile, action, "file:///index.ts")
	codeAction := converted.CodeAction
	switch {
	case codeAction.Command != nil:
		t.Errorf("expected no command, got %v", codeAction.Command)
	case codeAction.Edit == nil:
		t.Fatal("expected an edit to be attached to the code action")
	}
}

func TestConvertRefactorToLSPCodeAction_Disabled(t *testing.T) {
	t.Parallel()

	text := "const x: { a: number } = { a: 1 };\n"
	sourceFile := parser.ParseSourceFile(ast.SourceFileParseOptions{
		FileName: "/index.ts",
		Path:     "/index.ts",
	}, text, core.ScriptKindTS)

	converters := lsconv.NewConverters(lsproto.PositionEncodingKindUTF16, func(string) *lsconv.LSPLineMap {
		return lsconv.ComputeLSPLineStarts(text)
	})
	l := &LanguageService{converters: converters}

	action := &CodeAction{
		Description:    "Extract to type alias",
		DisabledReason: "Selection is not a valid type node",
		RenameFilename: "/index.ts",
		RenameLocation: 5,
	}

	converted := convertRefactorToLSPCodeAction(l, sourceFile, action, "file:///index.ts")
	codeAction := converted.CodeAction

	switch {
	case codeAction.Disabled == nil || codeAction.Disabled.Reason != "Selection is not a valid type node":
		t.Errorf("expected disabled reason, got %v", codeAction.Disabled)
	case codeAction.Edit != nil:
		t.Errorf("expected no edit for a disabled action, got %v", codeAction.Edit)
	case codeAction.Command != nil:
		t.Errorf("expected no rename command for a disabled action, got %v", codeAction.Command)
	}
}

var errRefactorFactoryFailed = errors.New("refactor factory failed")

func newRefactorPipelineTestLS(t *testing.T, text string) (*LanguageService, *ast.SourceFile) {
	t.Helper()

	sourceFile := parser.ParseSourceFile(ast.SourceFileParseOptions{
		FileName: "/index.ts",
		Path:     "/index.ts",
	}, text, core.ScriptKindTS)
	converters := lsconv.NewConverters(lsproto.PositionEncodingKindUTF16, func(string) *lsconv.LSPLineMap {
		return lsconv.ComputeLSPLineStarts(text)
	})

	return &LanguageService{converters: converters, activeConfig: lsutil.NewDefaultUserPreferences()}, sourceFile
}

func setRefactorProviders(t *testing.T, providers []*RefactorProvider) {
	t.Helper()

	old := refactorProviders

	refactorProviders = providers

	t.Cleanup(func() { refactorProviders = old })
}

func refactorPipelineParams() *RefactorContext {
	return &RefactorContext{
		SourceFile: nil,
		Params: &lsproto.CodeActionParams{
			TextDocument: lsproto.TextDocumentIdentifier{Uri: "file:///index.ts"},
			Context:      &lsproto.CodeActionContext{},
		},
	}
}

func TestProvideRefactorActionsForProviders_OnlyFilter(t *testing.T) {
	text := "const x: { a: number } = { a: 1 };\n"
	l, sourceFile := newRefactorPipelineTestLS(t, text)
	refactorContext := refactorPipelineParams()

	refactorContext.SourceFile = sourceFile

	extractFactoryCalled := false

	setRefactorProviders(t, []*RefactorProvider{{
		RefactorActions: []RefactorAction{
			{
				ID:    "extract-type",
				Title: "Extract type",
				Kinds: []lsproto.CodeActionKind{lsproto.CodeActionKindRefactorExtract},
				Factory: func(ctx context.Context, refactorContext *RefactorContext, refactorID string) ([]*CodeAction, error) {
					extractFactoryCalled = true
					if refactorID != "extract-type" {
						t.Errorf("expected refactorID %q, got %q", "extract-type", refactorID)
					}
					return []*CodeAction{
						{Description: "Extract to type alias", Changes: []*lsproto.TextEdit{{
							Range:   lsproto.Range{Start: lsproto.Position{Line: 0, Character: 10}, End: lsproto.Position{Line: 0, Character: 23}},
							NewText: "NewType",
						}}},
					}, nil
				},
			},
			{
				ID:    "infer-return-type",
				Title: "Infer return type",
				Kinds: []lsproto.CodeActionKind{lsproto.CodeActionKindRefactorRewrite},
				Factory: func(ctx context.Context, refactorContext *RefactorContext, refactorID string) ([]*CodeAction, error) {
					t.Errorf("factory for %q should not have been called", refactorID)
					return nil, nil
				},
			},
		},
	}})

	only := &[]lsproto.CodeActionKind{lsproto.CodeActionKindRefactorExtract}

	var actions []lsproto.CommandOrCodeAction

	if err := l.provideRefactorActionsForProviders(context.Background(), refactorContext, only, &actions); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !extractFactoryCalled {
		t.Fatal("expected the matching refactor factory to be called")
	}
	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}
	if got := actions[0].CodeAction.Title; got != "Extract to type alias" {
		t.Errorf("expected title %q, got %q", "Extract to type alias", got)
	}
}

func TestProvideRefactorActionsForProviders_DisabledGating(t *testing.T) {
	text := "const x: { a: number } = { a: 1 };\n"
	l, sourceFile := newRefactorPipelineTestLS(t, text)
	refactorContext := refactorPipelineParams()

	refactorContext.SourceFile = sourceFile

	setRefactorProviders(t, []*RefactorProvider{{
		RefactorActions: []RefactorAction{{
			ID:    "extract-type",
			Title: "Extract type",
			Kinds: []lsproto.CodeActionKind{lsproto.CodeActionKindRefactorExtract},
			Factory: func(ctx context.Context, refactorContext *RefactorContext, refactorID string) ([]*CodeAction, error) {
				return []*CodeAction{
					{Description: "Extract to type alias", Changes: []*lsproto.TextEdit{{
						Range:   lsproto.Range{Start: lsproto.Position{Line: 0, Character: 10}, End: lsproto.Position{Line: 0, Character: 23}},
						NewText: "NewType",
					}}},
					{Description: "Extract to type alias", DisabledReason: "Selection is not a valid type node"},
				}, nil
			},
		}},
	}})

	only := &[]lsproto.CodeActionKind{lsproto.CodeActionKindRefactorExtract}
	disabledCtx := lsproto.WithClientCapabilities(context.Background(), &lsproto.ResolvedClientCapabilities{
		TextDocument: lsproto.ResolvedTextDocumentClientCapabilities{
			CodeAction: lsproto.ResolvedCodeActionClientCapabilities{DisabledSupport: true},
		},
	})

	var actions []lsproto.CommandOrCodeAction
	if err := l.provideRefactorActionsForProviders(disabledCtx, refactorContext, only, &actions); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(actions) != 2 {
		t.Fatalf("expected 2 actions when the client supports disabled reasons, got %d", len(actions))
	}
	if actions[1].CodeAction.Disabled == nil {
		t.Fatal("expected the disabled action to carry a Disabled reason")
	}

	actions = nil
	if err := l.provideRefactorActionsForProviders(context.Background(), refactorContext, only, &actions); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(actions) != 1 {
		t.Fatalf("expected 1 action when the client does not support disabled reasons, got %d", len(actions))
	}

	l.activeConfig.ProvideRefactorNotApplicableReason = core.TSFalse
	actions = nil
	if err := l.provideRefactorActionsForProviders(disabledCtx, refactorContext, only, &actions); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(actions) != 1 {
		t.Fatalf("expected 1 action when the preference is disabled, got %d", len(actions))
	}
}

func TestProvideRefactorActionsForProviders_FactoryError(t *testing.T) {
	text := "const x: { a: number } = { a: 1 };\n"
	l, sourceFile := newRefactorPipelineTestLS(t, text)
	refactorContext := refactorPipelineParams()

	refactorContext.SourceFile = sourceFile

	setRefactorProviders(t, []*RefactorProvider{{
		RefactorActions: []RefactorAction{{
			ID:    "extract-type",
			Title: "Extract type",
			Kinds: []lsproto.CodeActionKind{lsproto.CodeActionKindRefactorExtract},
			Factory: func(ctx context.Context, refactorContext *RefactorContext, refactorID string) ([]*CodeAction, error) {
				return nil, errRefactorFactoryFailed
			},
		}},
	}})

	var actions []lsproto.CommandOrCodeAction

	err := l.provideRefactorActionsForProviders(context.Background(), refactorContext, nil, &actions)
	if !errors.Is(err, errRefactorFactoryFailed) {
		t.Fatalf("expected factory error, got %v", err)
	}
	if !strings.Contains(err.Error(), "extract-type") {
		t.Errorf("expected error to reference the action ID, got %v", err)
	}
}
