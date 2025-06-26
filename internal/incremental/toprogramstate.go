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
	for _, fileName := range buildInfo.FileNames {
		path := tspath.ToPath(fileName, to.buildInfoDirectory, config.UseCaseSensitiveFileNames())
		to.filePaths = append(to.filePaths, path)
	}
	for _, fileIdList := range buildInfo.FileIdsList {
		fileSet := collections.NewSetWithSizeHint[tspath.Path](len(fileIdList))
		for _, fileId := range fileIdList {
			fileSet.Add(to.toFilePath(fileId))
		}
		to.filePathSet = append(to.filePathSet, fileSet)
	}
	to.setCompilerOptions()
	to.setFileInfoAndEmitSignatures()
	to.setReferenceMap()
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

func (t *toProgramState) toDiagnosticCompatibleWithProgramState(diagnostics []*BuildInfoDiagnostic) ([]*BuildInfoDiagnostic, bool) {
	var fixedDiagnostics []*BuildInfoDiagnostic
	var changed bool
	for _, d := range diagnostics {
		file := d.file
		if d.file != false && d.file != 0 {
			file = t.toFilePath(BuildInfoFileId(d.file.(float64)))
		}
		messageChain, changedMessageChain := t.toDiagnosticCompatibleWithProgramState(d.messageChain)
		relatedInformation, changedRelatedInformation := t.toDiagnosticCompatibleWithProgramState(d.relatedInformation)
		if file != d.file || changedMessageChain || changedRelatedInformation {
			fixedDiagnostics = append(fixedDiagnostics, &BuildInfoDiagnostic{
				file:               file,
				loc:                d.loc,
				code:               d.code,
				category:           d.category,
				message:            d.message,
				messageChain:       messageChain,
				relatedInformation: relatedInformation,
				reportsUnnecessary: d.reportsUnnecessary,
				reportsDeprecated:  d.reportsDeprecated,
				skippedOnNoEmit:    d.skippedOnNoEmit,
			})
			changed = true
		} else {
			fixedDiagnostics = append(fixedDiagnostics, d)
		}
	}
	return core.IfElse(changed, fixedDiagnostics, diagnostics), changed
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

func (t *toProgramState) setReferenceMap() {
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
	t.state.semanticDiagnosticsPerFile = make(map[tspath.Path]*diagnosticsOrBuildInfoDiagnostics, len(t.state.fileInfos))
	for path := range t.state.fileInfos {
		// Initialize to have no diagnostics if its not changed file
		if !t.state.changedFilesSet.Has(path) {
			t.state.semanticDiagnosticsPerFile[path] = &diagnosticsOrBuildInfoDiagnostics{}
		}
	}
	for _, diagnostic := range t.buildInfo.SemanticDiagnosticsPerFile {
		if diagnostic.FileId != 0 {
			filePath := t.toFilePath(diagnostic.FileId)
			delete(t.state.semanticDiagnosticsPerFile, filePath) // does not have cached diagnostics
		} else {
			filePath := t.toFilePath(diagnostic.Diagnostic.fileId)
			diagnostics, _ := t.toDiagnosticCompatibleWithProgramState(diagnostic.Diagnostic.diagnostics)
			t.state.semanticDiagnosticsPerFile[filePath] = &diagnosticsOrBuildInfoDiagnostics{
				buildInfoDiagnostics: diagnostics,
			}
		}
	}
}

func (t *toProgramState) setEmitDiagnostics() {
	t.state.emitDiagnosticsPerFile = make(map[tspath.Path]*diagnosticsOrBuildInfoDiagnostics, len(t.state.fileInfos))
	for _, diagnostic := range t.buildInfo.EmitDiagnosticsPerFile {
		filePath := t.toFilePath(diagnostic.fileId)
		diagnostics, _ := t.toDiagnosticCompatibleWithProgramState(diagnostic.diagnostics)
		t.state.emitDiagnosticsPerFile[filePath] = &diagnosticsOrBuildInfoDiagnostics{
			buildInfoDiagnostics: diagnostics,
		}
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
