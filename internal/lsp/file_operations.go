package lsp

import (
	"context"
	"strings"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project"
)

func newRenameFileOperationRegistrationOptions() *lsproto.FileOperationRegistrationOptions {
	fileScheme := "file"
	fileMatch := lsproto.FileOperationPatternKindFile
	return &lsproto.FileOperationRegistrationOptions{
		Filters: []*lsproto.FileOperationFilter{
			{
				Scheme: &fileScheme,
				Pattern: &lsproto.FileOperationPattern{
					Glob:    "**/*",
					Matches: &fileMatch,
				},
			},
		},
	}
}

func combineWorkspaceEditResults(results []lsproto.WorkspaceEditOrNull) lsproto.WorkspaceEditOrNull {
	combined := make(map[lsproto.DocumentUri][]*lsproto.TextEdit)
	seenChanges := make(map[lsproto.DocumentUri]*collections.Set[lsproto.Range])

	for _, result := range results {
		if result.WorkspaceEdit == nil || result.WorkspaceEdit.Changes == nil {
			continue
		}

		for uri, edits := range *result.WorkspaceEdit.Changes {
			seenSet := seenChanges[uri]
			if seenSet == nil {
				seenSet = &collections.Set[lsproto.Range]{}
				seenChanges[uri] = seenSet
			}

			for _, edit := range edits {
				if edit == nil || !seenSet.AddIfAbsent(edit.Range) {
					continue
				}
				combined[uri] = append(combined[uri], edit)
			}
		}
	}

	if len(combined) == 0 {
		return lsproto.WorkspaceEditOrNull{}
	}

	return lsproto.WorkspaceEditOrNull{
		WorkspaceEdit: &lsproto.WorkspaceEdit{
			Changes: &combined,
		},
	}
}

func renameFilesToWatchedFileEvents(files []*lsproto.FileRename) []*lsproto.FileEvent {
	type uriState struct {
		old bool
		new bool
	}

	states := make(map[lsproto.DocumentUri]uriState)
	for _, file := range files {
		if file == nil {
			continue
		}

		oldURI := lsproto.DocumentUri(file.OldUri)
		newURI := lsproto.DocumentUri(file.NewUri)
		if !strings.HasPrefix(string(oldURI), "file://") || !strings.HasPrefix(string(newURI), "file://") {
			continue
		}

		oldState := states[oldURI]
		oldState.old = true
		states[oldURI] = oldState

		newState := states[newURI]
		newState.new = true
		states[newURI] = newState
	}

	events := make([]*lsproto.FileEvent, 0, len(states))
	for uri, state := range states {
		event := &lsproto.FileEvent{Uri: uri}
		switch {
		case state.old && state.new:
			event.Type = lsproto.FileChangeTypeChanged
		case state.new:
			event.Type = lsproto.FileChangeTypeCreated
		case state.old:
			event.Type = lsproto.FileChangeTypeDeleted
		default:
			continue
		}
		events = append(events, event)
	}

	return events
}

func activeFileForRenameProject(p *project.Project, files []*lsproto.FileRename) string {
	for _, file := range files {
		if file == nil {
			continue
		}

		oldURI := lsproto.DocumentUri(file.OldUri)
		if !strings.HasPrefix(string(oldURI), "file://") {
			continue
		}

		if p.HasFile(oldURI.FileName()) {
			return oldURI.FileName()
		}
	}

	sourceFiles := p.GetProgram().GetSourceFiles()
	if len(sourceFiles) == 0 {
		return ""
	}

	return sourceFiles[0].FileName()
}

func (s *Server) handleWillRenameFiles(ctx context.Context, params *lsproto.RenameFilesParams, req *lsproto.RequestMessage) (lsproto.WillRenameFilesResponse, error) {
	defer s.recover(req)

	if params == nil || len(params.Files) == 0 {
		return lsproto.WorkspaceEditOrNull{}, nil
	}

	documents := make([]lsproto.DocumentUri, 0, len(params.Files))
	for _, file := range params.Files {
		if file == nil {
			continue
		}

		oldURI := lsproto.DocumentUri(file.OldUri)
		if !strings.HasPrefix(string(oldURI), "file://") {
			continue
		}

		documents = append(documents, oldURI)
	}

	snapshot := s.session.GetSnapshotLoadingDocumentsAndProjectTree(ctx, documents, nil)
	var results []lsproto.WorkspaceEditOrNull
	for _, project := range snapshot.ProjectCollection.Projects() {
		if ctx.Err() != nil {
			return lsproto.WorkspaceEditOrNull{}, ctx.Err()
		}

		activeFile := activeFileForRenameProject(project, params.Files)
		if activeFile == "" {
			continue
		}

		languageService := ls.NewLanguageService(project.Id(), project.GetProgram(), snapshot, activeFile)
		result, err := languageService.ProvideWillRenameFiles(ctx, params)
		if err != nil {
			return lsproto.WorkspaceEditOrNull{}, err
		}
		results = append(results, result)
	}

	return combineWorkspaceEditResults(results), nil
}

func (s *Server) handleDidRenameFiles(ctx context.Context, params *lsproto.RenameFilesParams) error {
	if params == nil || len(params.Files) == 0 {
		return nil
	}

	fileEvents := renameFilesToWatchedFileEvents(params.Files)
	if len(fileEvents) == 0 {
		return nil
	}

	s.session.DidChangeWatchedFiles(ctx, fileEvents)
	return nil
}
