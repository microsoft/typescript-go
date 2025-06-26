package incrementaltestutil

import (
	"encoding/json"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/incremental"
)

type readableBuildInfo struct {
	buildInfo *incremental.BuildInfo
	Version   string

	// Common between incremental and tsc -b buildinfo for non incremental programs
	Errors       bool `json:"errors,omitzero"`
	CheckPending bool `json:"checkPending,omitzero"`
	// Root         []BuildInfoRoot `json:"root,omitempty,omitzero"`

	// IncrementalProgram info
	FileNames                  []string                                  `json:"fileNames,omitzero"`
	FileInfos                  []*readableBuildInfoFileInfo              `json:"fileInfos,omitzero"`
	FileIdsList                [][]string                                `json:"fileIdsList,omitzero"`
	Options                    *collections.OrderedMap[string, any]      `json:"options,omitempty"`
	ReferencedMap              []readableBuildInfoReferenceMapEntry      `json:"referencedMap,omitzero"`
	SemanticDiagnosticsPerFile []incremental.BuildInfoSemanticDiagnostic `json:"semanticDiagnosticsPerFile,omitzero"`
	EmitDiagnosticsPerFile     []incremental.BuildInfoDiagnosticOfFile   `json:"emitDiagnosticsPerFile,omitzero"`
	ChangeFileSet              []incremental.BuildInfoFileId             `json:"changeFileSet,omitzero"`
	AffectedFilesPendingEmit   []incremental.BuildInfoFilePendingEmit    `json:"affectedFilesPendingEmit,omitzero"`
	LatestChangedDtsFile       string                                    `json:"latestChangedDtsFile,omitzero"` // Because this is only output file in the program, we dont need fileId to deduplicate name
	EmitSignatures             []incremental.BuildInfoEmitSignature      `json:"emitSignatures,omitzero"`
	// resolvedRoot: readonly IncrementalBuildInfoResolvedRoot[] | undefined;
	Size int `json:"size,omitzero"` // Size of the build info file
}

type readableBuildInfoFileInfo struct {
	Version            string                         `json:"version,omitzero"`
	Signature          string                         `json:"signature,omitzero"`
	AffectsGlobalScope bool                           `json:"affectsGlobalScope,omitzero"`
	ImpliedNodeFormat  string                         `json:"impliedNodeFormat,omitzero"`
	Original           *incremental.BuildInfoFileInfo `json:"original,omitempty"` // Original file path, if available
}

type readableBuildInfoReferenceMapEntry struct {
	file     string
	fileList []string
}

func (r *readableBuildInfoReferenceMapEntry) MarshalJSON() ([]byte, error) {
	return json.Marshal([2]any{r.file, r.fileList})
}

func (r *readableBuildInfoReferenceMapEntry) UnmarshalJSON(data []byte) error {
	var v [2]any
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*r = readableBuildInfoReferenceMapEntry{
		file:     v[0].(string),
		fileList: v[1].([]string),
	}
	return nil
}

// type readableBuildInfoSemanticDiagnostic struct {
// 	file       string                                // File is not in changedSet and still doesnt have cached diagnostics
// 	diagnostic incremental.BuildInfoDiagnosticOfFile // Diagnostics for file
// }

// func (b *readableBuildInfoSemanticDiagnostic) MarshalJSON() ([]byte, error) {
// 	if b.file != 0 {
// 		return json.Marshal(b.file)
// 	}
// 	return json.Marshal(b.diagnostic)
// }

// func (b *readableBuildInfoSemanticDiagnostic) UnmarshalJSON(data []byte) error {
// 	var fileId BuildInfoFileId
// 	if err := json.Unmarshal(data, &fileId); err == nil {
// 		*b = BuildInfoSemanticDiagnostic{
// 			FileId: fileId,
// 		}
// 		return nil
// 	}
// 	var diagnostic BuildInfoDiagnosticOfFile
// 	if err := json.Unmarshal(data, &diagnostic); err == nil {
// 		*b = BuildInfoSemanticDiagnostic{
// 			Diagnostic: diagnostic,
// 		}
// 		return nil
// 	}
// 	return fmt.Errorf("invalid IncrementalBuildInfoDiagnostic: %s", data)
// }

func toReadableBuildInfo(buildInfo *incremental.BuildInfo, buildInfoText string) string {
	readable := readableBuildInfo{
		buildInfo: buildInfo,
		Version:   buildInfo.Version,

		Errors:       buildInfo.Errors,
		CheckPending: buildInfo.CheckPending,
		// Root:         buildInfo.Root,

		FileNames: buildInfo.FileNames,
		Options:   buildInfo.Options,
		// SemanticDiagnosticsPerFile []incremental.BuildInfoSemanticDiagnostic `json:"semanticDiagnosticsPerFile,omitzero"`
		// EmitDiagnosticsPerFile     []incremental.BuildInfoDiagnosticOfFile   `json:"emitDiagnosticsPerFile,omitzero"`
		// ChangeFileSet              []incremental.BuildInfoFileId             `json:"changeFileSet,omitzero"`
		// AffectedFilesPendingEmit   []incremental.BuildInfoFilePendingEmit    `json:"affectedFilesPendingEmit,omitzero"`
		LatestChangedDtsFile: buildInfo.LatestChangedDtsFile,
		// EmitSignatures             []incremental.BuildInfoEmitSignature      `json:"emitSignatures,omitzero"`
		// resolvedRoot: readonly IncrementalBuildInfoResolvedRoot[] | undefined;

		Size: len(buildInfoText),
	}
	readable.setFileInfos()
	readable.setFileIdsList()
	readable.setReferencedMap()
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

func (r *readableBuildInfo) setFileInfos() {
	for _, original := range r.buildInfo.FileInfos {
		fileInfo := original.GetFileInfo()
		// Dont set original for string encoding
		if original.HasSignature() {
			original = nil
		}
		r.FileInfos = append(r.FileInfos, &readableBuildInfoFileInfo{
			Version:            fileInfo.Version(),
			Signature:          fileInfo.Signature(),
			AffectsGlobalScope: fileInfo.AffectsGlobalScope(),
			ImpliedNodeFormat:  fileInfo.ImpliedNodeFormat().String(),
			Original:           original,
		})
	}
}

func (r *readableBuildInfo) setFileIdsList() {
	for _, ids := range r.buildInfo.FileIdsList {
		r.FileIdsList = append(r.FileIdsList, core.Map(ids, r.toFilePath))
	}
}

func (r *readableBuildInfo) setReferencedMap() {
	for _, entry := range r.buildInfo.ReferencedMap {
		r.ReferencedMap = append(r.ReferencedMap, readableBuildInfoReferenceMapEntry{
			file:     r.toFilePath(entry.FileId),
			fileList: r.toFilePathSet(entry.FileIdListId),
		})
	}
}
