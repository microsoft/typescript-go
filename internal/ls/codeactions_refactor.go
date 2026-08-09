package ls

import (
	"context"
	"fmt"
	"slices"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
)

// RefactorActionFactory is a function that produces refactoring code actions.
type RefactorActionFactory func(ctx context.Context, refactorContext *RefactorContext, refactorID string) ([]*CodeAction, error)

// RefactorAction describes a single refactoring action offered by a RefactorProvider.
type RefactorAction struct {
	ID      string
	Title   string
	Kinds   []lsproto.CodeActionKind
	Factory RefactorActionFactory
}

// RefactorProvider provides refactoring code actions.
type RefactorProvider struct {
	RefactorActions []RefactorAction
}

// RefactorContext contains the context needed to generate refactoring actions.
type RefactorContext struct {
	SourceFile *ast.SourceFile
	Range      core.TextRange
	Program    *compiler.Program
	LS         *LanguageService
	Params     *lsproto.CodeActionParams
}

// CodeActionTriggerKind returns the trigger kind of the code action request, or nil when not provided.
func (c *RefactorContext) CodeActionTriggerKind() *lsproto.CodeActionTriggerKind {
	if c.Params == nil || c.Params.Context == nil {
		return nil
	}

	return c.Params.Context.TriggerKind
}

var refactorProviders = []*RefactorProvider{}

// provideRefactorActions adds contextual (non-diagnostic-driven) refactoring code actions for the given range.
func (l *LanguageService) provideRefactorActions(ctx context.Context, params *lsproto.CodeActionParams, program *compiler.Program, file *ast.SourceFile, actions *[]lsproto.CommandOrCodeAction) error {
	if params.Context == nil || !wantsRefactors(params.Context.Only) {
		return nil
	}

	refactorContext := &RefactorContext{
		SourceFile: file,
		Range: core.NewTextRange(
			int(l.converters.LineAndCharacterToPosition(file, params.Range.Start)),
			int(l.converters.LineAndCharacterToPosition(file, params.Range.End)),
		),
		Program: program,
		LS:      l,
		Params:  params,
	}

	return l.provideRefactorActionsForProviders(ctx, refactorContext, params.Context.Only, actions)
}

func (l *LanguageService) provideRefactorActionsForProviders(ctx context.Context, refactorContext *RefactorContext, only *[]lsproto.CodeActionKind, actions *[]lsproto.CommandOrCodeAction) error {
	for _, provider := range refactorProviders {
		for _, action := range provider.RefactorActions {
			lspActions, err := l.convertRefactorAction(ctx, refactorContext, action, only)
			if err != nil {
				return err
			}

			*actions = append(*actions, lspActions...)
		}
	}

	return nil
}

func (l *LanguageService) convertRefactorAction(ctx context.Context, refactorContext *RefactorContext, action RefactorAction, only *[]lsproto.CodeActionKind) ([]lsproto.CommandOrCodeAction, error) {
	if !refactorActionMatchesOnly(action, only) {
		return nil, nil
	}

	providerActions, err := action.Factory(ctx, refactorContext, action.ID)
	if err != nil {
		return nil, fmt.Errorf("refactoring action %q: %w", action.ID, err)
	}

	showDisabled := showNotApplicableReasons(ctx, l)
	file := refactorContext.SourceFile
	uri := refactorContext.Params.TextDocument.Uri

	var lspActions []lsproto.CommandOrCodeAction

	for _, a := range providerActions {
		if a.DisabledReason != "" && !showDisabled {
			continue
		}

		if a.RenameFilename != "" && lsconv.FileNameToDocumentURI(a.RenameFilename) != uri {
			return nil, fmt.Errorf("refactoring action %q: rename target %q does not match requested document %q", action.ID, a.RenameFilename, uri)
		}

		lspActions = append(lspActions, convertRefactorToLSPCodeAction(l, file, a, uri))
	}

	return lspActions, nil
}

func showNotApplicableReasons(ctx context.Context, l *LanguageService) bool {
	return l.UserPreferences().ProvideRefactorNotApplicableReason.IsTrue() &&
		lsproto.GetClientCapabilities(ctx).TextDocument.CodeAction.DisabledSupport
}

func wantsRefactors(only *[]lsproto.CodeActionKind) bool {
	if only == nil || len(*only) == 0 {
		return true
	}

	for _, kind := range *only {
		if kind == lsproto.CodeActionKindEmpty || codeActionKindContains(lsproto.CodeActionKindRefactor, kind) {
			return true
		}
	}

	return false
}

func refactorActionMatchesOnly(action RefactorAction, only *[]lsproto.CodeActionKind) bool {
	if only == nil || len(*only) == 0 {
		return true
	}

	for _, requestedKind := range *only {
		for _, actionKind := range action.Kinds {
			if codeActionKindContains(requestedKind, actionKind) {
				return true
			}
		}
	}

	return false
}

func convertRefactorToLSPCodeAction(l *LanguageService, file *ast.SourceFile, action *CodeAction, uri lsproto.DocumentUri) lsproto.CommandOrCodeAction {
	kind := action.Kind
	if kind == "" {
		kind = lsproto.CodeActionKindRefactorRewrite
	}

	lspAction := &lsproto.CodeAction{
		Title: action.Description,
		Kind:  &kind,
	}

	if action.DisabledReason != "" {
		lspAction.Disabled = &lsproto.CodeActionDisabled{Reason: action.DisabledReason}
		return lsproto.CommandOrCodeAction{CodeAction: lspAction}
	}

	setRefactorEdit(l, file, action, uri, lspAction)

	return lsproto.CommandOrCodeAction{
		CodeAction: lspAction,
	}
}

func setRefactorEdit(l *LanguageService, file *ast.SourceFile, action *CodeAction, uri lsproto.DocumentUri, lspAction *lsproto.CodeAction) {
	changes := map[lsproto.DocumentUri][]*lsproto.TextEdit{
		uri: action.Changes,
	}

	lspAction.Edit = &lsproto.WorkspaceEdit{Changes: &changes}

	if action.RenameFilename != "" {
		lspAction.Command = refactorToRenameCommand(l, file, action)
	}
}

func refactorToRenameCommand(l *LanguageService, file *ast.SourceFile, action *CodeAction) *lsproto.Command {
	// The rename location is an offset into the file after the action's edits have been applied.
	postEditText := applyTextEdits(file, action.Changes, l)
	renamePosition := l.converters.PositionToLineAndCharacterForText(postEditText, core.TextPos(action.RenameLocation))

	return &lsproto.Command{
		Title:   "",
		Command: "editor.action.rename",
		Arguments: &[]any{
			[]any{lsconv.FileNameToDocumentURI(action.RenameFilename), renamePosition},
		},
	}
}

func applyTextEdits(file *ast.SourceFile, edits []*lsproto.TextEdit, l *LanguageService) string {
	changes := make([]core.TextChange, 0, len(edits))

	for _, edit := range edits {
		changes = append(changes, core.TextChange{
			TextRange: core.NewTextRange(
				int(l.converters.LineAndCharacterToPosition(file, edit.Range.Start)),
				int(l.converters.LineAndCharacterToPosition(file, edit.Range.End)),
			),
			NewText: edit.NewText,
		})
	}

	slices.SortStableFunc(changes, func(a, b core.TextChange) int {
		return a.Pos() - b.Pos()
	})

	return core.ApplyBulkEdits(file.Text(), changes)
}
