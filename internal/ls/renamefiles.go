package ls

import (
	"context"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/modulespecifiers"
	"github.com/microsoft/typescript-go/internal/scanner"
	"github.com/microsoft/typescript-go/internal/stringutil"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type fileRenameIndex struct {
	currentDirectory      string
	useCaseSensitiveNames bool
	renamedFilesByOldPath map[tspath.Path]string
}

type renamedSourceFile struct {
	sourceFile *ast.SourceFile
	fileName   string
	path       tspath.Path
}

func (f *renamedSourceFile) FileName() string {
	return f.fileName
}

func (f *renamedSourceFile) Path() tspath.Path {
	return f.path
}

func (f *renamedSourceFile) Imports() []*ast.StringLiteralLike {
	return f.sourceFile.Imports()
}

func (f *renamedSourceFile) IsJS() bool {
	return f.sourceFile.IsJS()
}

func newFileRenameIndex(files []*lsproto.FileRename, currentDirectory string, useCaseSensitiveNames bool) *fileRenameIndex {
	renamedFilesByOldPath := make(map[tspath.Path]string)
	for _, file := range files {
		if file == nil {
			continue
		}

		oldURI := lsproto.DocumentUri(file.OldUri)
		newURI := lsproto.DocumentUri(file.NewUri)
		if !isFileURI(oldURI) || !isFileURI(newURI) {
			continue
		}

		oldPath := oldURI.Path(useCaseSensitiveNames)
		newPath := newURI.Path(useCaseSensitiveNames)
		if oldPath == newPath {
			continue
		}

		renamedFilesByOldPath[oldPath] = newURI.FileName()
	}

	if len(renamedFilesByOldPath) == 0 {
		return nil
	}

	return &fileRenameIndex{
		currentDirectory:      currentDirectory,
		useCaseSensitiveNames: useCaseSensitiveNames,
		renamedFilesByOldPath: renamedFilesByOldPath,
	}
}

func (r *fileRenameIndex) Apply(fileName string) string {
	if r == nil {
		return fileName
	}

	if renamed, ok := r.renamedFilesByOldPath[tspath.ToPath(fileName, r.currentDirectory, r.useCaseSensitiveNames)]; ok {
		return renamed
	}

	return fileName
}

func isFileURI(uri lsproto.DocumentUri) bool {
	return strings.HasPrefix(string(uri), "file://")
}

func (l *LanguageService) ProvideWillRenameFiles(ctx context.Context, params *lsproto.RenameFilesParams) (lsproto.WorkspaceEditOrNull, error) {
	if params == nil || len(params.Files) == 0 {
		return lsproto.WorkspaceEditOrNull{}, nil
	}

	program := l.GetProgram()
	renames := newFileRenameIndex(params.Files, program.GetCurrentDirectory(), l.UseCaseSensitiveFileNames())
	if renames == nil {
		return lsproto.WorkspaceEditOrNull{}, nil
	}

	changes := make(map[lsproto.DocumentUri][]*lsproto.TextEdit)
	for _, sourceFile := range program.GetSourceFiles() {
		if ctx.Err() != nil {
			return lsproto.WorkspaceEditOrNull{}, ctx.Err()
		}

		if len(sourceFile.Imports()) == 0 || program.IsLibFile(sourceFile) || program.IsSourceFileFromExternalLibrary(sourceFile) {
			continue
		}

		fileEdits := l.getImportEditsForRenamedFiles(sourceFile, program, renames)
		if len(fileEdits) == 0 {
			continue
		}

		changes[lsconv.FileNameToDocumentURI(sourceFile.FileName())] = fileEdits
	}

	if len(changes) == 0 {
		return lsproto.WorkspaceEditOrNull{}, nil
	}

	return lsproto.WorkspaceEditOrNull{
		WorkspaceEdit: &lsproto.WorkspaceEdit{
			Changes: &changes,
		},
	}, nil
}

func (l *LanguageService) getImportEditsForRenamedFiles(sourceFile *ast.SourceFile, program *compiler.Program, renames *fileRenameIndex) []*lsproto.TextEdit {
	currentFileName := sourceFile.FileName()
	renamedFileName := renames.Apply(currentFileName)

	importingSourceFile := modulespecifiers.SourceFileForSpecifierGeneration(sourceFile)
	if renamedFileName != currentFileName {
		importingSourceFile = &renamedSourceFile{
			sourceFile: sourceFile,
			fileName:   renamedFileName,
			path:       tspath.ToPath(renamedFileName, program.GetCurrentDirectory(), l.UseCaseSensitiveFileNames()),
		}
	}

	preferences := l.host.GetPreferences(renamedFileName)
	if preferences == nil {
		preferences = lsutil.NewDefaultUserPreferences()
	}

	var edits []*lsproto.TextEdit
	for _, moduleSpecifier := range sourceFile.Imports() {
		if !tspath.IsExternalModuleNameRelative(moduleSpecifier.Text()) {
			continue
		}

		resolution := program.GetResolvedModuleFromModuleSpecifier(sourceFile, moduleSpecifier)
		if resolution == nil || resolution.ResolvedFileName == "" {
			continue
		}

		renamedTargetFileName := renames.Apply(resolution.ResolvedFileName)
		if renamedFileName == currentFileName && renamedTargetFileName == resolution.ResolvedFileName {
			continue
		}

		newModuleSpecifiers, _ := modulespecifiers.GetModuleSpecifiersForFileWithInfo(
			importingSourceFile,
			renamedTargetFileName,
			program.Options(),
			program,
			preferences.ModuleSpecifierPreferences(),
			modulespecifiers.ModuleSpecifierOptions{},
			false, /*forAutoImports*/
		)
		if len(newModuleSpecifiers) == 0 || newModuleSpecifiers[0] == moduleSpecifier.Text() {
			continue
		}

		edits = append(edits, &lsproto.TextEdit{
			Range:   *l.createLspRangeFromNode(moduleSpecifier, sourceFile),
			NewText: quoteModuleSpecifier(sourceFile, moduleSpecifier, newModuleSpecifiers[0]),
		})
	}

	return edits
}

func quoteModuleSpecifier(sourceFile *ast.SourceFile, moduleSpecifier *ast.StringLiteralLike, text string) string {
	quoted, _ := core.StringifyJson(text, "" /*prefix*/, "" /*indent*/)
	start := scanner.GetTokenPosOfNode(moduleSpecifier, sourceFile, false /*includeJSDoc*/)
	if start >= 0 && start < len(sourceFile.Text()) {
		switch sourceFile.Text()[start] {
		case '\'':
			return "'" + quoteReplacer.Replace(stringutil.StripQuotes(quoted)) + "'"
		case '`':
			return "`" + strings.ReplaceAll(text, "`", "\\`") + "`"
		}
	}

	return quoted
}
