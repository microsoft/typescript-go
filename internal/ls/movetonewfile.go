package ls

// movetonewfile.go implements the "Move to a new file" refactoring.
//
// !!! TODO: This is a stub implementation. The full implementation requires
// complex interactions with the type checker, module resolution, and import
// management that have not yet been fully ported.

import (
	"context"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/change"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
)

// GetMoveToNewFileEdits returns the text edits needed to move the selected
// statements to a new file. The selection is given by span in sourceFile.
//
// Returns a map from filename to list of text edits. The map will contain
// edits for the original file (to remove the moved statements and update
// imports) as well as edits for the new file (to create it with the moved
// statements).
//
// !!! TODO: Implement full logic from TypeScript's moveToNewFile refactoring.
func (l *LanguageService) GetMoveToNewFileEdits(
	ctx context.Context,
	sourceFile *ast.SourceFile,
	program *compiler.Program,
	span core.TextRange,
) map[string][]*lsproto.TextEdit {
	// !!! TODO: Implement full logic:
	// 1. Call getStatementsToMove to determine what to move
	// 2. Call getUsageInfo to analyze symbol usage
	// 3. Create a new filename based on the moved symbols
	// 4. Use a change tracker to:
	//    a. Remove moved statements from original file
	//    b. Update imports in original file
	//    c. Add exports for symbols used by moved code
	//    d. Create new file with moved statements and necessary imports
	// 5. Return the changes

	toMove := getStatementsToMove(sourceFile, span)
	if toMove == nil {
		return nil
	}

	changeTracker := change.NewTracker(ctx, program.Options(), l.FormatOptions(), l.converters)
	doMoveToNewFile(ctx, sourceFile, program, toMove, changeTracker)
	return changeTracker.GetChanges()
}

// doMoveToNewFile performs the actual work of moving statements to a new file.
//
// !!! TODO: Implement full logic:
// 1. Get type checker
// 2. Compute usage info (movedSymbols, imports needed, etc.)
// 3. Determine new filename
// 4. Delete moved statements from old file
// 5. Add/update imports in old file
// 6. Add exports for things needed by the new file
// 7. Create new file content with moved statements and their imports
func doMoveToNewFile(
	_ context.Context,
	_ *ast.SourceFile,
	_ *compiler.Program,
	_ *ToMove,
	_ *change.Tracker,
) {
	// Not yet implemented.
}
