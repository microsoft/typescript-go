package incremental

import (
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

func buildInfoToProgramState(buildInfo *BuildInfo, buildInfoFileName string, config *tsoptions.ParsedCommandLine) *programState {
	to := &toProgramState{
		buildInfo:          buildInfo,
		buildInfoDirectory: tspath.GetDirectoryPath(tspath.GetNormalizedAbsolutePath(buildInfoFileName, config.GetCurrentDirectory())),
		filePaths:          make([]tspath.Path, 0, len(buildInfo.FileNames)),
		filePathSet:        make([]*collections.Set[tspath.Path], 0, len(buildInfo.FileIdsList)),
	}
	to.filePaths = core.Map(buildInfo.FileNames, func(fileName string) tspath.Path {
		return tspath.ToPath(fileName, to.buildInfoDirectory, config.UseCaseSensitiveFileNames())
	})
	to.filePathSet = core.Map(buildInfo.FileIdsList, func(fileIdList []BuildInfoFileId) *collections.Set[tspath.Path] {
		fileSet := collections.NewSetWithSizeHint[tspath.Path](len(fileIdList))
		for _, fileId := range fileIdList {
			fileSet.Add(to.toFilePath(fileId))
		}
		return fileSet
	})
	to.setCompilerOptions()
	to.setFileInfoAndEmitSignatures()
	to.setReferencedMap()
	to.setChangeFileSet()
	to.setSemanticDiagnostics()
	to.setEmitDiagnostics()
	to.setAffectedFilesPendingEmit()
	if buildInfo.LatestChangedDtsFile != "" {
		to.state.latestChangedDtsFile = to.toAbsolutePath(buildInfo.LatestChangedDtsFile)
	}
	to.state.hasErrors = core.IfElse(buildInfo.Errors, core.TSTrue, core.TSFalse)
	to.state.checkPending = buildInfo.CheckPending
	return &to.state
}

type toProgramState struct {
	buildInfo          *BuildInfo
	buildInfoDirectory string
	state              programState
	filePaths          []tspath.Path
	filePathSet        []*collections.Set[tspath.Path]
}

func (t *toProgramState) toAbsolutePath(path string) string {
	return tspath.GetNormalizedAbsolutePath(path, t.buildInfoDirectory)
}

func (t *toProgramState) toFilePath(fileId BuildInfoFileId) tspath.Path {
	return t.filePaths[fileId-1]
}

func (t *toProgramState) toFilePathSet(fileIdListId BuildInfoFileIdListId) *collections.Set[tspath.Path] {
	return t.filePathSet[fileIdListId-1]
}

func (t *toProgramState) toBuildInfoDiagnosticsWithFileName(diagnostics []*BuildInfoDiagnostic) []*buildInfoDiagnosticWithFileName {
	return core.Map(diagnostics, func(d *BuildInfoDiagnostic) *buildInfoDiagnosticWithFileName {
		var file tspath.Path
		if d.File != 0 {
			file = t.toFilePath(d.File)
		}
		return &buildInfoDiagnosticWithFileName{
			file:               file,
			noFile:             d.NoFile,
			pos:                d.Pos,
			end:                d.End,
			code:               d.Code,
			category:           d.Category,
			message:            d.Message,
			messageChain:       t.toBuildInfoDiagnosticsWithFileName(d.MessageChain),
			relatedInformation: t.toBuildInfoDiagnosticsWithFileName(d.RelatedInformation),
			reportsUnnecessary: d.ReportsUnnecessary,
			reportsDeprecated:  d.ReportsDeprecated,
			skippedOnNoEmit:    d.SkippedOnNoEmit,
		}
	})
}

func (t *toProgramState) toDiagnosticsOrBuildInfoDiagnosticsWithFileName(dig *BuildInfoDiagnosticsOfFile) *diagnosticsOrBuildInfoDiagnosticsWithFileName {
	return &diagnosticsOrBuildInfoDiagnosticsWithFileName{
		buildInfoDiagnostics: t.toBuildInfoDiagnosticsWithFileName(dig.Diagnostics),
	}
}

func (t *toProgramState) setCompilerOptions() {
	t.state.options = &core.CompilerOptions{}
	for option, value := range t.buildInfo.Options.Entries() {
		result, ok := tsoptions.ConvertOptionToAbsolutePath(option, value, tsoptions.CommandLineCompilerOptionsMap, t.buildInfoDirectory)
		if ok {
			tsoptions.ParseCompilerOptions(option, result, t.state.options)
		} else {
			tsoptions.ParseCompilerOptions(option, value, t.state.options)
		}
	}
}

func (t *toProgramState) setFileInfoAndEmitSignatures() {
	t.state.fileInfos = make(map[tspath.Path]*fileInfo, len(t.buildInfo.FileInfos))
	t.state.createEmitSignaturesMap()
	for index, buildInfoFileInfo := range t.buildInfo.FileInfos {
		path := t.toFilePath(BuildInfoFileId(index + 1))
		info := buildInfoFileInfo.GetFileInfo()
		t.state.fileInfos[path] = info
		// Add default emit signature as file's signature
		if info.signature != "" && len(t.state.emitSignatures) != 0 {
			t.state.emitSignatures[path] = &emitSignature{signature: info.signature}
		}
	}
	// Fix up emit signatures
	for _, value := range t.buildInfo.EmitSignatures {
		if value.noEmitSignature() {
			delete(t.state.emitSignatures, t.toFilePath(value.fileId))
		} else {
			path := t.toFilePath(value.fileId)
			t.state.emitSignatures[path] = value.toEmitSignature(path, t.state.emitSignatures)
		}
	}
}

func (t *toProgramState) setReferencedMap() {
	t.state.createReferenceMap()
	for _, entry := range t.buildInfo.ReferencedMap {
		t.state.referencedMap.Add(t.toFilePath(entry.FileId), t.toFilePathSet(entry.FileIdListId))
	}
}

func (t *toProgramState) setChangeFileSet() {
	t.state.changedFilesSet = collections.NewSetWithSizeHint[tspath.Path](len(t.buildInfo.ChangeFileSet))
	for _, fileId := range t.buildInfo.ChangeFileSet {
		filePath := t.toFilePath(fileId)
		t.state.changedFilesSet.Add(filePath)
	}
}

func (t *toProgramState) setSemanticDiagnostics() {
	t.state.semanticDiagnosticsPerFile = make(map[tspath.Path]*diagnosticsOrBuildInfoDiagnosticsWithFileName, len(t.state.fileInfos))
	for path := range t.state.fileInfos {
		// Initialize to have no diagnostics if its not changed file
		if !t.state.changedFilesSet.Has(path) {
			t.state.semanticDiagnosticsPerFile[path] = &diagnosticsOrBuildInfoDiagnosticsWithFileName{}
		}
	}
	for _, diagnostic := range t.buildInfo.SemanticDiagnosticsPerFile {
		if diagnostic.FileId != 0 {
			filePath := t.toFilePath(diagnostic.FileId)
			delete(t.state.semanticDiagnosticsPerFile, filePath) // does not have cached diagnostics
		} else {
			filePath := t.toFilePath(diagnostic.Diagnostics.FileId)
			t.state.semanticDiagnosticsPerFile[filePath] = t.toDiagnosticsOrBuildInfoDiagnosticsWithFileName(diagnostic.Diagnostics)
		}
	}
}

func (t *toProgramState) setEmitDiagnostics() {
	t.state.emitDiagnosticsPerFile = make(map[tspath.Path]*diagnosticsOrBuildInfoDiagnosticsWithFileName, len(t.state.fileInfos))
	for _, diagnostic := range t.buildInfo.EmitDiagnosticsPerFile {
		filePath := t.toFilePath(diagnostic.FileId)
		t.state.emitDiagnosticsPerFile[filePath] = t.toDiagnosticsOrBuildInfoDiagnosticsWithFileName(diagnostic)
	}
}

func (t *toProgramState) setAffectedFilesPendingEmit() {
	if len(t.buildInfo.AffectedFilesPendingEmit) == 0 {
		return
	}
	ownOptionsEmitKind := getFileEmitKind(t.state.options)
	t.state.affectedFilesPendingEmit = make(map[tspath.Path]fileEmitKind, len(t.buildInfo.AffectedFilesPendingEmit))
	for _, pendingEmit := range t.buildInfo.AffectedFilesPendingEmit {
		t.state.affectedFilesPendingEmit[t.toFilePath(pendingEmit.fileId)] = core.IfElse(pendingEmit.emitKind == 0, ownOptionsEmitKind, pendingEmit.emitKind)
	}
}
