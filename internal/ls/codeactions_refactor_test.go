package ls

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
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
