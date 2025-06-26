package incremental

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

func programStateToBuildInfo(state *programState, program *compiler.Program, buildInfoFileName string) *BuildInfo {
	to := &toBuildInfo{
		state:              state,
		program:            program,
		buildInfoDirectory: tspath.GetDirectoryPath(buildInfoFileName),
		comparePathsOptions: tspath.ComparePathsOptions{
			CurrentDirectory:          program.GetCurrentDirectory(),
			UseCaseSensitiveFileNames: program.UseCaseSensitiveFileNames(),
		},
		fileNameToFileId:        make(map[string]BuildInfoFileId),
		fileNamesToFileIdListId: make(map[string]BuildInfoFileIdListId),
	}
	to.buildInfo.Version = core.Version()
	if state.options.IsIncremental() {
		to.setFileInfoAndEmitSignatures()
		to.setCompilerOptions()
		to.setReferenceMap()
		to.setChangeFileSet()
		to.setSemanticDiagnostics()
		to.setEmitDiagnostics()
		to.setAffectedFilesPendingEmit()
		if state.latestChangedDtsFile != "" {
			to.buildInfo.LatestChangedDtsFile = to.relativeToBuildInfo(state.latestChangedDtsFile)
		}
	}
	// else {
	//     const buildInfo: NonIncrementalBuildInfo = {
	//         root: arrayFrom(rootFileNames, r => relativeToBuildInfo(r)),
	//     };
	// }
	to.buildInfo.Errors = state.hasErrors.IsTrue()
	to.buildInfo.CheckPending = state.checkPending
	return &to.buildInfo
}

type toBuildInfo struct {
	state                   *programState
	program                 *compiler.Program
	buildInfo               BuildInfo
	buildInfoDirectory      string
	comparePathsOptions     tspath.ComparePathsOptions
	fileNameToFileId        map[string]BuildInfoFileId
	fileNamesToFileIdListId map[string]BuildInfoFileIdListId
}

func (t *toBuildInfo) relativeToBuildInfo(path string) string {
	return tspath.EnsurePathIsNonModuleName(tspath.GetRelativePathFromDirectory(t.buildInfoDirectory, path, t.comparePathsOptions))
}

func (t *toBuildInfo) toFileId(path tspath.Path) BuildInfoFileId {
	fileId := t.fileNameToFileId[string(path)]
	if fileId == 0 {
		t.buildInfo.FileNames = append(t.buildInfo.FileNames, t.relativeToBuildInfo(string(path)))
		fileId = BuildInfoFileId(len(t.buildInfo.FileNames))
		t.fileNameToFileId[string(path)] = fileId
	}
	return fileId
}

func (t *toBuildInfo) toFileIdListId(set *collections.Set[tspath.Path]) BuildInfoFileIdListId {
	fileIds := core.Map(slices.Collect(maps.Keys(set.Keys())), t.toFileId)
	slices.Sort(fileIds)
	key := strings.Join(core.Map(fileIds, func(id BuildInfoFileId) string {
		return fmt.Sprintf("%d", id)
	}), ",")

	fileIdListId := t.fileNamesToFileIdListId[key]
	if fileIdListId == 0 {
		t.buildInfo.FileIdsList = append(t.buildInfo.FileIdsList, fileIds)
		fileIdListId = BuildInfoFileIdListId(len(t.buildInfo.FileIdsList))
		t.fileNamesToFileIdListId[key] = fileIdListId
	}
	return fileIdListId
}

func (t *toBuildInfo) toRelativeToBuildInfoCompilerOptionValue(option *tsoptions.CommandLineOption, v any) any {
	if !option.IsFilePath {
		return v
	}
	if option.Kind == "list" {
		if arr, ok := v.([]string); ok {
			return core.Map(arr, func(item string) string {
				return t.relativeToBuildInfo(item)
			})
		}
	} else {
		return t.relativeToBuildInfo(v.(string))
	}
	return v
}

func (t *toBuildInfo) toDiagnosticCompatibleWithBuildInfo(diagnostics []*BuildInfoDiagnostic) ([]*BuildInfoDiagnostic, bool) {
	var fixedDiagnostics []*BuildInfoDiagnostic
	var changed bool
	for _, d := range diagnostics {
		file := d.file
		if d.file != false && d.file != "" {
			file = t.toFileId(d.file.(tspath.Path))
		}
		messageChain, changedMessageChain := t.toDiagnosticCompatibleWithBuildInfo(d.messageChain)
		relatedInformation, changedRelatedInformation := t.toDiagnosticCompatibleWithBuildInfo(d.relatedInformation)
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

func (t *toBuildInfo) toIncrementalBuildInfoDiagnostic(diagnostics []*ast.Diagnostic, file *ast.SourceFile) []*BuildInfoDiagnostic {
	var incrementalDiagnostics []*BuildInfoDiagnostic
	for _, d := range diagnostics {
		var fileValue any
		if d.File() == nil {
			fileValue = false
		} else if d.File() != file {
			fileValue = t.toFileId(d.File().Path())
		}
		incrementalDiagnostics = append(incrementalDiagnostics, &BuildInfoDiagnostic{
			file:               fileValue,
			loc:                d.Loc(),
			code:               d.Code(),
			category:           d.Category(),
			message:            d.Message(),
			messageChain:       t.toIncrementalBuildInfoDiagnostic(d.MessageChain(), file),
			relatedInformation: t.toIncrementalBuildInfoDiagnostic(d.RelatedInformation(), file),
			reportsUnnecessary: d.ReportsUnnecessary(),
			reportsDeprecated:  d.ReportsDeprecated(),
			skippedOnNoEmit:    d.SkippedOnNoEmit(),
		})
	}
	return incrementalDiagnostics
}

func (t *toBuildInfo) setFileInfoAndEmitSignatures() {
	for _, file := range t.program.GetSourceFiles() {
		info := t.state.fileInfos[file.Path()]
		fileId := t.toFileId(file.Path())
		//  tryAddRoot(key, fileId);
		if t.buildInfo.FileNames[fileId-1] != t.relativeToBuildInfo(string(file.Path())) {
			panic(fmt.Sprintf("File name at index %d does not match expected relative path: %s != %s", fileId-1, t.buildInfo.FileNames[fileId-1], t.relativeToBuildInfo(string(file.Path()))))
		}
		var actualSignature string
		if oldSignature, gotOldSignature := t.state.oldSignatures[file.Path()]; gotOldSignature {
			actualSignature = oldSignature
		} else {
			actualSignature = info.signature
		}
		if t.state.options.Composite.IsTrue() {
			if !ast.IsJsonSourceFile(file) && t.program.SourceFileMayBeEmitted(file, false) {
				emitSignature := t.state.emitSignatures[file.Path()]
				if emitSignature == nil {
					t.buildInfo.EmitSignatures = append(t.buildInfo.EmitSignatures, BuildInfoEmitSignature{
						fileId: fileId,
					})
				} else if emitSignature.signature != actualSignature {
					incrementalEmitSignature := BuildInfoEmitSignature{
						fileId: fileId,
					}
					if emitSignature.signature != "" {
						incrementalEmitSignature.signature = emitSignature.signature
					} else if emitSignature.signatureWithDifferentOptions[0] == actualSignature {
						incrementalEmitSignature.differsOnlyInDtsMap = true
					} else {
						incrementalEmitSignature.signature = emitSignature.signatureWithDifferentOptions[0]
						incrementalEmitSignature.differsInOptions = true
					}
					t.buildInfo.EmitSignatures = append(t.buildInfo.EmitSignatures, incrementalEmitSignature)
				}
			}
		}
		if actualSignature == info.signature {
			t.setFileInfo(info)
		} else {
			t.setFileInfo(&fileInfo{
				version:            info.version,
				signature:          actualSignature,
				affectsGlobalScope: info.affectsGlobalScope,
				impliedNodeFormat:  info.impliedNodeFormat,
			})
		}
	}
}

func (t *toBuildInfo) setFileInfo(fileInfo *fileInfo) {
	if fileInfo.version == fileInfo.signature {
		if !fileInfo.affectsGlobalScope && fileInfo.impliedNodeFormat == core.ResolutionModeNone {
			t.buildInfo.FileInfos = append(t.buildInfo.FileInfos, &BuildInfoFileInfo{signature: fileInfo.signature})
			return
		}
	} else if fileInfo.signature == "" {
		t.buildInfo.FileInfos = append(t.buildInfo.FileInfos, &BuildInfoFileInfo{noSignature: &buildInfoFileInfoNoSignature{
			Version:            fileInfo.version,
			NoSignature:        true,
			AffectsGlobalScope: fileInfo.affectsGlobalScope,
			ImpliedNodeFormat:  fileInfo.impliedNodeFormat,
		}})
		return
	}
	t.buildInfo.FileInfos = append(t.buildInfo.FileInfos, &BuildInfoFileInfo{fileInfo: &buildInfoFileInfoWithSignature{
		Version:            fileInfo.version,
		Signature:          core.IfElse(fileInfo.signature == fileInfo.version, "", fileInfo.signature),
		AffectsGlobalScope: fileInfo.affectsGlobalScope,
		ImpliedNodeFormat:  fileInfo.impliedNodeFormat,
	}})
}

func (t *toBuildInfo) setCompilerOptions() {
	tsoptions.ForEachCompilerOptionValue(
		t.state.options, func(option *tsoptions.CommandLineOption) bool {
			return option.AffectsBuildInfo
		},
		func(option *tsoptions.CommandLineOption, value any, i int) bool {
			// Make it relative to buildInfo directory if file path
			t.buildInfo.Options = append(t.buildInfo.Options, BuildInfoCompilerOption{
				name:  option.Name,
				value: t.toRelativeToBuildInfoCompilerOptionValue(option, value),
			})
			return false
		},
	)
}

func (t *toBuildInfo) setReferenceMap() {
	if !t.state.tracksReferences() {
		return
	}
	keys := slices.Collect(maps.Keys(t.state.referencedMap.Keys()))
	slices.Sort(keys)
	for _, filePath := range keys {
		references, _ := t.state.referencedMap.GetValues(filePath)
		t.buildInfo.ReferencedMap = append(t.buildInfo.ReferencedMap, BuildInfoReferenceMapEntry{
			FileId:       t.toFileId(filePath),
			FileIdListId: t.toFileIdListId(references),
		})
	}
}

func (t *toBuildInfo) setChangeFileSet() {
	files := slices.Collect(maps.Keys(t.state.changedFilesSet.Keys()))
	slices.Sort(files)
	for _, filePath := range files {
		t.buildInfo.ChangeFileSet = append(t.buildInfo.ChangeFileSet, t.toFileId(filePath))
	}
}

func (t *toBuildInfo) setSemanticDiagnostics() {
	for _, file := range t.program.GetSourceFiles() {
		value, ok := t.state.semanticDiagnosticsPerFile[file.Path()]
		if !ok {
			if !t.state.changedFilesSet.Has(file.Path()) {
				t.buildInfo.SemanticDiagnosticsPerFile = append(t.buildInfo.SemanticDiagnosticsPerFile, BuildInfoSemanticDiagnostic{
					FileId: t.toFileId(file.Path()),
				})
			}
		} else if len(value.buildInfoDiagnostics) > 0 {
			diagnostics, _ := t.toDiagnosticCompatibleWithBuildInfo(value.buildInfoDiagnostics)
			t.buildInfo.SemanticDiagnosticsPerFile = append(t.buildInfo.SemanticDiagnosticsPerFile, BuildInfoSemanticDiagnostic{
				Diagnostic: BuildInfoDiagnosticOfFile{
					fileId:      t.toFileId(file.Path()),
					diagnostics: diagnostics,
				},
			})
		} else if len(value.diagnostics) > 0 {
			t.buildInfo.SemanticDiagnosticsPerFile = append(t.buildInfo.SemanticDiagnosticsPerFile, BuildInfoSemanticDiagnostic{
				Diagnostic: BuildInfoDiagnosticOfFile{
					fileId:      t.toFileId(file.Path()),
					diagnostics: t.toIncrementalBuildInfoDiagnostic(value.diagnostics, file),
				},
			})
		}
	}
}

func (t *toBuildInfo) setEmitDiagnostics() {
	files := slices.Collect(maps.Keys(t.state.emitDiagnosticsPerFile))
	slices.Sort(files)
	for _, filePath := range files {
		value := t.state.emitDiagnosticsPerFile[filePath]
		if len(value.buildInfoDiagnostics) > 0 {
			diagnostics, _ := t.toDiagnosticCompatibleWithBuildInfo(value.buildInfoDiagnostics)
			t.buildInfo.EmitDiagnosticsPerFile = append(t.buildInfo.EmitDiagnosticsPerFile, BuildInfoDiagnosticOfFile{
				fileId:      t.toFileId(filePath),
				diagnostics: diagnostics,
			})
		} else {
			t.buildInfo.EmitDiagnosticsPerFile = append(t.buildInfo.EmitDiagnosticsPerFile, BuildInfoDiagnosticOfFile{
				fileId:      t.toFileId(filePath),
				diagnostics: t.toIncrementalBuildInfoDiagnostic(value.diagnostics, t.program.GetSourceFileByPath(filePath)),
			})
		}
	}
}

func (t *toBuildInfo) setAffectedFilesPendingEmit() {
	if len(t.state.affectedFilesPendingEmit) == 0 {
		return
	}
	files := slices.Collect(maps.Keys(t.state.affectedFilesPendingEmit))
	slices.Sort(files)
	fullEmitKind := getFileEmitKind(t.state.options)
	for _, filePath := range files {
		file := t.program.GetSourceFileByPath(filePath)
		if file == nil || !t.program.SourceFileMayBeEmitted(file, false) {
			continue
		}
		pendingEmit := t.state.affectedFilesPendingEmit[filePath]
		t.buildInfo.AffectedFilesPendingEmit = append(t.buildInfo.AffectedFilesPendingEmit, BuildInfoFilePendingEmit{
			fileId:   t.toFileId(filePath),
			emitKind: core.IfElse(pendingEmit == fullEmitKind, 0, pendingEmit),
		})
	}
}
