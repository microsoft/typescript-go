package incrementaltestutil

import (
	"encoding/json"
	"fmt"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/incremental"
)

type readableBuildInfo struct {
	buildInfo *incremental.BuildInfo
	Version   string

	// Common between incremental and tsc -b buildinfo for non incremental programs
	Errors       bool `json:"errors,omitzero"`
	CheckPending bool `json:"checkPending,omitzero"`
	// Root         []BuildInfoRoot `json:"root,omitzero"`

	// IncrementalProgram info
	FileNames                  []string                                  `json:"fileNames,omitzero"`
	FileInfos                  []*readableBuildInfoFileInfo              `json:"fileInfos,omitzero"`
	FileIdsList                [][]string                                `json:"fileIdsList,omitzero"`
	Options                    *collections.OrderedMap[string, any]      `json:"options,omitzero"`
	ReferencedMap              *collections.OrderedMap[string, []string] `json:"referencedMap,omitzero"`
	SemanticDiagnosticsPerFile []*readableBuildInfoSemanticDiagnostic    `json:"semanticDiagnosticsPerFile,omitzero"`
	EmitDiagnosticsPerFile     []*readableBuildInfoDiagnosticsOfFile     `json:"emitDiagnosticsPerFile,omitzero"`
	ChangeFileSet              []string                                  `json:"changeFileSet,omitzero"` // List of changed files in the program, not the whole set of files
	AffectedFilesPendingEmit   []*incremental.BuildInfoFilePendingEmit   `json:"affectedFilesPendingEmit,omitzero"`
	LatestChangedDtsFile       string                                    `json:"latestChangedDtsFile,omitzero"` // Because this is only output file in the program, we dont need fileId to deduplicate name
	EmitSignatures             []*incremental.BuildInfoEmitSignature     `json:"emitSignatures,omitzero"`
	// resolvedRoot: readonly IncrementalBuildInfoResolvedRoot[] | undefined;
	Size int `json:"size,omitzero"` // Size of the build info file
}

type readableBuildInfoFileInfo struct {
	FileName           string                         `json:"fileName,omitzero"`
	Version            string                         `json:"version,omitzero"`
	Signature          string                         `json:"signature,omitzero"`
	AffectsGlobalScope bool                           `json:"affectsGlobalScope,omitzero"`
	ImpliedNodeFormat  string                         `json:"impliedNodeFormat,omitzero"`
	Original           *incremental.BuildInfoFileInfo `json:"original,omitzero"` // Original file path, if available
}

type readableBuildInfoDiagnostic struct {
	// incrementalBuildInfoFileId if it is for a File thats other than its stored for
	File               string                         `json:"file,omitzero"`
	NoFile             bool                           `json:"noFile,omitzero"`
	Pos                int                            `json:"pos,omitzero"`
	End                int                            `json:"end,omitzero"`
	Code               int32                          `json:"code,omitzero"`
	Category           diagnostics.Category           `json:"category,omitzero"`
	Message            string                         `json:"message,omitzero"`
	MessageChain       []*readableBuildInfoDiagnostic `json:"messageChain,omitzero"`
	RelatedInformation []*readableBuildInfoDiagnostic `json:"relatedInformation,omitzero"`
	ReportsUnnecessary bool                           `json:"reportsUnnecessary,omitzero"`
	ReportsDeprecated  bool                           `json:"reportsDeprecated,omitzero"`
	SkippedOnNoEmit    bool                           `json:"skippedOnNoEmit,omitzero"`
}

type readableBuildInfoDiagnosticsOfFile struct {
	file        string
	diagnostics []*readableBuildInfoDiagnostic
}

func (r *readableBuildInfoDiagnosticsOfFile) MarshalJSON() ([]byte, error) {
	fileIdAndDiagnostics := make([]any, 0, 2)
	fileIdAndDiagnostics = append(fileIdAndDiagnostics, r.file)
	fileIdAndDiagnostics = append(fileIdAndDiagnostics, r.diagnostics)
	return json.Marshal(fileIdAndDiagnostics)
}

func (r *readableBuildInfoDiagnosticsOfFile) UnmarshalJSON(data []byte) error {
	var fileIdAndDiagnostics []any
	if err := json.Unmarshal(data, &fileIdAndDiagnostics); err != nil {
		return fmt.Errorf("invalid readableBuildInfoDiagnosticsOfFile: %s", data)
	}
	if len(fileIdAndDiagnostics) != 2 {
		return fmt.Errorf("invalid readableBuildInfoDiagnosticsOfFile: expected 2 elements, got %d", len(fileIdAndDiagnostics))
	}
	file, ok := fileIdAndDiagnostics[0].(string)
	if !ok {
		return fmt.Errorf("invalid fileId in readableBuildInfoDiagnosticsOfFile: expected string, got %T", fileIdAndDiagnostics[0])
	}
	if diagnostics, ok := fileIdAndDiagnostics[1].([]*readableBuildInfoDiagnostic); ok {
		*r = readableBuildInfoDiagnosticsOfFile{
			file:        file,
			diagnostics: diagnostics,
		}
		return nil
	}
	return fmt.Errorf("invalid diagnostics in readableBuildInfoDiagnosticsOfFile: expected []*readableBuildInfoDiagnostic, got %T", fileIdAndDiagnostics[1])
}

type readableBuildInfoSemanticDiagnostic struct {
	file        string                              // File is not in changedSet and still doesnt have cached diagnostics
	diagnostics *readableBuildInfoDiagnosticsOfFile // Diagnostics for file
}

func (r *readableBuildInfoSemanticDiagnostic) MarshalJSON() ([]byte, error) {
	if r.file != "" {
		return json.Marshal(r.file)
	}
	return json.Marshal(r.diagnostics)
}

func (r *readableBuildInfoSemanticDiagnostic) UnmarshalJSON(data []byte) error {
	var file string
	if err := json.Unmarshal(data, &file); err == nil {
		*r = readableBuildInfoSemanticDiagnostic{
			file: file,
		}
		return nil
	}
	var diagnostics readableBuildInfoDiagnosticsOfFile
	if err := json.Unmarshal(data, &diagnostics); err == nil {
		*r = readableBuildInfoSemanticDiagnostic{
			diagnostics: &diagnostics,
		}
		return nil
	}
	return fmt.Errorf("invalid readableBuildInfoSemanticDiagnostic: %s", data)
}

func toReadableBuildInfo(buildInfo *incremental.BuildInfo, buildInfoText string) string {
	readable := readableBuildInfo{
		buildInfo:            buildInfo,
		Version:              buildInfo.Version,
		Errors:               buildInfo.Errors,
		CheckPending:         buildInfo.CheckPending,
		FileNames:            buildInfo.FileNames,
		Options:              buildInfo.Options,
		LatestChangedDtsFile: buildInfo.LatestChangedDtsFile,
		Size:                 len(buildInfoText),
	}
	readable.setFileInfos()
	readable.setFileIdsList()
	readable.setReferencedMap()
	readable.setChangeFileSet()
	readable.setSemanticDiagnostics()
	readable.setEmitDiagnostics()
	readable.setAffectedFilesPendingEmit()
	readable.setEmitSignatures()
	contents, err := json.MarshalIndent(&readable, "", "    ")
	if err != nil {
		panic("readableBuildInfo: failed to marshal readable build info: " + err.Error())
	}
	return string(contents)
}

func (r *readableBuildInfo) toFilePath(fileId incremental.BuildInfoFileId) string {
	return r.buildInfo.FileNames[fileId-1]
}

func (r *readableBuildInfo) toFilePathSet(fileIdListId incremental.BuildInfoFileIdListId) []string {
	return r.FileIdsList[fileIdListId-1]
}

func (r *readableBuildInfo) toReadableBuildInfoDiagnostic(diagnostics []*incremental.BuildInfoDiagnostic) []*readableBuildInfoDiagnostic {
	return core.Map(diagnostics, func(d *incremental.BuildInfoDiagnostic) *readableBuildInfoDiagnostic {
		var file string
		if d.File != 0 {
			file = r.toFilePath(d.File)
		}
		return &readableBuildInfoDiagnostic{
			File:               file,
			NoFile:             d.NoFile,
			Pos:                d.Pos,
			End:                d.End,
			Code:               d.Code,
			Category:           d.Category,
			Message:            d.Message,
			MessageChain:       r.toReadableBuildInfoDiagnostic(d.MessageChain),
			RelatedInformation: r.toReadableBuildInfoDiagnostic(d.RelatedInformation),
			ReportsUnnecessary: d.ReportsUnnecessary,
			ReportsDeprecated:  d.ReportsDeprecated,
			SkippedOnNoEmit:    d.SkippedOnNoEmit,
		}
	})
}

func (r *readableBuildInfo) toReadableBuildInfoDiagnosticsOfFile(diagnostics *incremental.BuildInfoDiagnosticsOfFile) *readableBuildInfoDiagnosticsOfFile {
	return &readableBuildInfoDiagnosticsOfFile{
		file:        r.toFilePath(diagnostics.FileId),
		diagnostics: r.toReadableBuildInfoDiagnostic(diagnostics.Diagnostics),
	}
}

func (r *readableBuildInfo) toReadableBuildInfoSemanticDiagnostic(diagnostics *incremental.BuildInfoSemanticDiagnostic) *readableBuildInfoSemanticDiagnostic {
	if diagnostics.FileId != 0 {
		return &readableBuildInfoSemanticDiagnostic{
			file: r.toFilePath(diagnostics.FileId),
		}
	}
	return &readableBuildInfoSemanticDiagnostic{
		diagnostics: r.toReadableBuildInfoDiagnosticsOfFile(diagnostics.Diagnostics),
	}
}

func (r *readableBuildInfo) setFileInfos() {
	r.FileInfos = core.MapIndex(r.buildInfo.FileInfos, func(original *incremental.BuildInfoFileInfo, index int) *readableBuildInfoFileInfo {
		fileInfo := original.GetFileInfo()
		// Dont set original for string encoding
		if original.HasSignature() {
			original = nil
		}
		return &readableBuildInfoFileInfo{
			FileName:           r.toFilePath(incremental.BuildInfoFileId(index + 1)),
			Version:            fileInfo.Version(),
			Signature:          fileInfo.Signature(),
			AffectsGlobalScope: fileInfo.AffectsGlobalScope(),
			ImpliedNodeFormat:  fileInfo.ImpliedNodeFormat().String(),
			Original:           original,
		}
	})
}

func (r *readableBuildInfo) setFileIdsList() {
	r.FileIdsList = core.Map(r.buildInfo.FileIdsList, func(ids []incremental.BuildInfoFileId) []string {
		return core.Map(ids, r.toFilePath)
	})
}

func (r *readableBuildInfo) setReferencedMap() {
	if r.buildInfo.ReferencedMap != nil {
		r.ReferencedMap = &collections.OrderedMap[string, []string]{}
		for _, entry := range r.buildInfo.ReferencedMap {
			r.ReferencedMap.Set(r.toFilePath(entry.FileId), r.toFilePathSet(entry.FileIdListId))
		}
	}
}

func (r *readableBuildInfo) setChangeFileSet() {
	r.ChangeFileSet = core.Map(r.buildInfo.ChangeFileSet, r.toFilePath)
}

func (r *readableBuildInfo) setSemanticDiagnostics() {
	r.SemanticDiagnosticsPerFile = core.Map(r.buildInfo.SemanticDiagnosticsPerFile, r.toReadableBuildInfoSemanticDiagnostic)
}

func (r *readableBuildInfo) setEmitDiagnostics() {
	r.EmitDiagnosticsPerFile = core.Map(r.buildInfo.EmitDiagnosticsPerFile, r.toReadableBuildInfoDiagnosticsOfFile)
}

func (r *readableBuildInfo) setAffectedFilesPendingEmit() {
	// if len(r.buildInfo.AffectedFilesPendingEmit) == 0 {
	// 	return
	// }
	// ownOptionsEmitKind := getFileEmitKind(r.state.options)
	// r.state.affectedFilesPendingEmit = make(map[tspath.Path]fileEmitKind, len(r.buildInfo.AffectedFilesPendingEmit))
	// for _, pendingEmit := range r.buildInfo.AffectedFilesPendingEmit {
	// 	r.state.affectedFilesPendingEmit[r.toFilePath(pendingEmit.fileId)] = core.IfElse(pendingEmit.emitKind == 0, ownOptionsEmitKind, pendingEmit.emitKind)
	// }
}

func (r *readableBuildInfo) setEmitSignatures() {
	// r.EmitSignatures = make([]*incremental.BuildInfoEmitSignature, 0, len(r.buildInfo.EmitSignatures))
	// for _, value := range r.buildInfo.EmitSignatures {
	// 	if value == nil {
	// 		continue
	// 	}
	// 	r.EmitSignatures = append(r.EmitSignatures, value)
	// }
}
