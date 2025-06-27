package incremental

import (
	"encoding/json"
	"fmt"

	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type (
	BuildInfoFileId       int
	BuildInfoFileIdListId int
)

// /**
//  * buildInfoRoot is
//  * for incremental program buildinfo
//  * - start and end of FileId for consecutive fileIds to be included as root
//  * - single fileId that is root
//  * for non incremental program buildinfo
//  * - string that is the root file name
//  */
// type buildInfoRoot struct {
// 	incrementalStartEnd *[2]incrementalBuildInfoFileId
// 	incrementalSingle   incrementalBuildInfoFileId
// 	nonIncremental      string
// }

// func (o buildInfoRoot) MarshalJSON() ([]byte, error) {
// 	if o.incrementalStartEnd != nil {
// 		return json.Marshal(o.incrementalStartEnd)
// 	}
// 	if o.incrementalSingle != 0 {
// 		return json.Marshal(o.incrementalSingle)
// 	}
// 	if o.nonIncremental != "" {
// 		return json.Marshal(o.nonIncremental)
// 	}
// 	panic("unknown BuildInfoRoot type")
// }

// func (o *buildInfoRoot) UnmarshalJSON(data []byte) error {
// 	*o = buildInfoRoot{}
// 	var vIncrementalStartEnd [2]incrementalBuildInfoFileId
// 	if err := json.Unmarshal(data, &vIncrementalStartEnd); err == nil {
// 		o.incrementalStartEnd = &vIncrementalStartEnd
// 		return nil
// 	}
// 	var vIncrementalSingle incrementalBuildInfoFileId
// 	if err := json.Unmarshal(data, &vIncrementalSingle); err == nil {
// 		o.incrementalSingle = vIncrementalSingle
// 		return nil
// 	}
// 	var vNonIncremental string
// 	if err := json.Unmarshal(data, &vNonIncremental); err == nil {
// 		o.nonIncremental = vNonIncremental
// 		return nil
// 	}
// 	return fmt.Errorf("invalid BuildInfoRoot: %s", data)
// }

type buildInfoFileInfoNoSignature struct {
	Version            string              `json:"version,omitzero"`
	NoSignature        bool                `json:"noSignature,omitzero"`
	AffectsGlobalScope bool                `json:"affectsGlobalScope,omitzero"`
	ImpliedNodeFormat  core.ResolutionMode `json:"impliedNodeFormat,omitzero"`
}

/**
 *   Signature is
 * 	 - undefined if FileInfo.version === FileInfo.signature
 * 	 - string actual signature
 */
type buildInfoFileInfoWithSignature struct {
	Version            string              `json:"version,omitzero"`
	Signature          string              `json:"signature,omitzero"`
	AffectsGlobalScope bool                `json:"affectsGlobalScope,omitzero"`
	ImpliedNodeFormat  core.ResolutionMode `json:"impliedNodeFormat,omitzero"`
}

type BuildInfoFileInfo struct {
	signature   string
	noSignature *buildInfoFileInfoNoSignature
	fileInfo    *buildInfoFileInfoWithSignature
}

func newBuildInfoFileInfo(fileInfo *fileInfo) *BuildInfoFileInfo {
	if fileInfo.version == fileInfo.signature {
		if !fileInfo.affectsGlobalScope && fileInfo.impliedNodeFormat == core.ResolutionModeCommonJS {
			return &BuildInfoFileInfo{signature: fileInfo.signature}
		}
	} else if fileInfo.signature == "" {
		return &BuildInfoFileInfo{noSignature: &buildInfoFileInfoNoSignature{
			Version:            fileInfo.version,
			NoSignature:        true,
			AffectsGlobalScope: fileInfo.affectsGlobalScope,
			ImpliedNodeFormat:  fileInfo.impliedNodeFormat,
		}}
	}
	return &BuildInfoFileInfo{fileInfo: &buildInfoFileInfoWithSignature{
		Version:            fileInfo.version,
		Signature:          core.IfElse(fileInfo.signature == fileInfo.version, "", fileInfo.signature),
		AffectsGlobalScope: fileInfo.affectsGlobalScope,
		ImpliedNodeFormat:  fileInfo.impliedNodeFormat,
	}}
}

func (b *BuildInfoFileInfo) GetFileInfo() *fileInfo {
	if b.signature != "" {
		return &fileInfo{
			version:           b.signature,
			signature:         b.signature,
			impliedNodeFormat: core.ResolutionModeCommonJS,
		}
	}
	if b.noSignature != nil {
		return &fileInfo{
			version:            b.noSignature.Version,
			affectsGlobalScope: b.noSignature.AffectsGlobalScope,
			impliedNodeFormat:  b.noSignature.ImpliedNodeFormat,
		}
	}
	return &fileInfo{
		version:            b.fileInfo.Version,
		signature:          core.IfElse(b.fileInfo.Signature == "", b.fileInfo.Version, b.fileInfo.Signature),
		affectsGlobalScope: b.fileInfo.AffectsGlobalScope,
		impliedNodeFormat:  b.fileInfo.ImpliedNodeFormat,
	}
}

func (b *BuildInfoFileInfo) HasSignature() bool {
	return b.signature != ""
}

func (b *BuildInfoFileInfo) MarshalJSON() ([]byte, error) {
	if b.signature != "" {
		return json.Marshal(b.signature)
	}
	if b.noSignature != nil {
		return json.Marshal(b.noSignature)
	}
	return json.Marshal(b.fileInfo)
}

func (b *BuildInfoFileInfo) UnmarshalJSON(data []byte) error {
	var vSignature string
	if err := json.Unmarshal(data, &vSignature); err == nil {
		*b = BuildInfoFileInfo{signature: vSignature}
		return nil
	}
	var noSignature buildInfoFileInfoNoSignature
	if err := json.Unmarshal(data, &noSignature); err == nil && noSignature.NoSignature {
		*b = BuildInfoFileInfo{noSignature: &noSignature}
		return nil
	}
	var fileInfo buildInfoFileInfoWithSignature
	if err := json.Unmarshal(data, &fileInfo); err != nil {
		return fmt.Errorf("invalid incrementalBuildInfoFileInfo: %s", data)
	}
	*b = BuildInfoFileInfo{fileInfo: &fileInfo}
	return nil
}

type BuildInfoReferenceMapEntry struct {
	FileId       BuildInfoFileId
	FileIdListId BuildInfoFileIdListId
}

func (b *BuildInfoReferenceMapEntry) MarshalJSON() ([]byte, error) {
	return json.Marshal([2]int{int(b.FileId), int(b.FileIdListId)})
}

func (b *BuildInfoReferenceMapEntry) UnmarshalJSON(data []byte) error {
	var v *[2]int
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*b = BuildInfoReferenceMapEntry{
		FileId:       BuildInfoFileId(v[0]),
		FileIdListId: BuildInfoFileIdListId(v[1]),
	}
	return nil
}

type BuildInfoDiagnostic struct {
	// incrementalBuildInfoFileId if it is for a File thats other than its stored for
	File               BuildInfoFileId        `json:"file,omitzero"`
	NoFile             bool                   `json:"noFile,omitzero"`
	Pos                int                    `json:"pos,omitzero"`
	End                int                    `json:"end,omitzero"`
	Code               int32                  `json:"code,omitzero"`
	Category           diagnostics.Category   `json:"category,omitzero"`
	Message            string                 `json:"message,omitzero"`
	MessageChain       []*BuildInfoDiagnostic `json:"messageChain,omitzero"`
	RelatedInformation []*BuildInfoDiagnostic `json:"relatedInformation,omitzero"`
	ReportsUnnecessary bool                   `json:"reportsUnnecessary,omitzero"`
	ReportsDeprecated  bool                   `json:"reportsDeprecated,omitzero"`
	SkippedOnNoEmit    bool                   `json:"skippedOnNoEmit,omitzero"`
}

type BuildInfoDiagnosticsOfFile struct {
	FileId      BuildInfoFileId
	Diagnostics []*BuildInfoDiagnostic
}

func (b *BuildInfoDiagnosticsOfFile) MarshalJSON() ([]byte, error) {
	fileIdAndDiagnostics := make([]any, 0, 2)
	fileIdAndDiagnostics = append(fileIdAndDiagnostics, b.FileId)
	fileIdAndDiagnostics = append(fileIdAndDiagnostics, b.Diagnostics)
	return json.Marshal(fileIdAndDiagnostics)
}

func (b *BuildInfoDiagnosticsOfFile) UnmarshalJSON(data []byte) error {
	var fileIdAndDiagnostics []any
	if err := json.Unmarshal(data, &fileIdAndDiagnostics); err != nil {
		return fmt.Errorf("invalid BuildInfoDiagnosticsOfFile: %s", data)
	}
	if len(fileIdAndDiagnostics) != 2 {
		return fmt.Errorf("invalid BuildInfoDiagnosticsOfFile: expected 2 elements, got %d", len(fileIdAndDiagnostics))
	}
	var fileId BuildInfoFileId
	if fileIdV, ok := fileIdAndDiagnostics[0].(float64); !ok {
		return fmt.Errorf("invalid fileId in BuildInfoDiagnosticsOfFile: expected float64, got %T", fileIdAndDiagnostics[0])
	} else {
		fileId = BuildInfoFileId(fileIdV)
	}
	if diagnostics, ok := fileIdAndDiagnostics[1].([]*BuildInfoDiagnostic); ok {
		*b = BuildInfoDiagnosticsOfFile{
			FileId:      fileId,
			Diagnostics: diagnostics,
		}
		return nil
	}
	return fmt.Errorf("invalid diagnostics in BuildInfoDiagnosticsOfFile: expected []*BuildInfoDiagnostic, got %T", fileIdAndDiagnostics[1])
}

type BuildInfoSemanticDiagnostic struct {
	FileId      BuildInfoFileId             // File is not in changedSet and still doesnt have cached diagnostics
	Diagnostics *BuildInfoDiagnosticsOfFile // Diagnostics for file
}

func (b *BuildInfoSemanticDiagnostic) MarshalJSON() ([]byte, error) {
	if b.FileId != 0 {
		return json.Marshal(b.FileId)
	}
	return json.Marshal(b.Diagnostics)
}

func (b *BuildInfoSemanticDiagnostic) UnmarshalJSON(data []byte) error {
	var fileId BuildInfoFileId
	if err := json.Unmarshal(data, &fileId); err == nil {
		*b = BuildInfoSemanticDiagnostic{
			FileId: fileId,
		}
		return nil
	}
	var diagnostics BuildInfoDiagnosticsOfFile
	if err := json.Unmarshal(data, &diagnostics); err == nil {
		*b = BuildInfoSemanticDiagnostic{
			Diagnostics: &diagnostics,
		}
		return nil
	}
	return fmt.Errorf("invalid BuildInfoSemanticDiagnostic: %s", data)
}

/**
 * fileId if pending emit is same as what compilerOptions suggest
 * [fileId] if pending emit is only dts file emit
 * [fileId, emitKind] if any other type emit is pending
 */
type BuildInfoFilePendingEmit struct {
	fileId   BuildInfoFileId
	emitKind fileEmitKind
}

func (b *BuildInfoFilePendingEmit) MarshalJSON() ([]byte, error) {
	if b.emitKind == 0 {
		return json.Marshal(b.fileId)
	}
	if b.emitKind == fileEmitKindDts {
		fileListIds := []BuildInfoFileId{b.fileId}
		return json.Marshal(fileListIds)
	}
	fileAndEmitKind := []int{int(b.fileId), int(b.emitKind)}
	return json.Marshal(fileAndEmitKind)
}

func (b *BuildInfoFilePendingEmit) UnmarshalJSON(data []byte) error {
	var fileId BuildInfoFileId
	if err := json.Unmarshal(data, &fileId); err == nil {
		*b = BuildInfoFilePendingEmit{
			fileId: fileId,
		}
		return nil
	}
	var intTuple []int
	if err := json.Unmarshal(data, &intTuple); err == nil {
		if len(intTuple) == 1 {
			*b = BuildInfoFilePendingEmit{
				fileId:   BuildInfoFileId(intTuple[0]),
				emitKind: fileEmitKindDts,
			}
			return nil
		} else if len(intTuple) == 2 {
			*b = BuildInfoFilePendingEmit{
				fileId:   BuildInfoFileId(intTuple[0]),
				emitKind: fileEmitKind(intTuple[1]),
			}
			return nil
		}
		return fmt.Errorf("invalid incrementalBuildInfoFilePendingEmit: expected 1 or 2 integers, got %d", len(intTuple))
	}
	return fmt.Errorf("invalid IncrementalBuildInfoDiagnostic: %s", data)
}

/**
 * [fileId, signature] if different from file's signature
 * fileId if file wasnt emitted
 */
type BuildInfoEmitSignature struct {
	fileId BuildInfoFileId
	// Signature if it is different from file's signature
	signature           string
	differsOnlyInDtsMap bool // true if signature is different only in dtsMap value
	differsInOptions    bool // true if signature is different in options used to emit file
}

func (b *BuildInfoEmitSignature) noEmitSignature() bool {
	return b.signature == "" && !b.differsOnlyInDtsMap && !b.differsInOptions
}

func (b *BuildInfoEmitSignature) toEmitSignature(path tspath.Path, emitSignatures map[tspath.Path]*emitSignature) *emitSignature {
	var signature string
	var signatureWithDifferentOptions []string
	if b.differsOnlyInDtsMap {
		signatureWithDifferentOptions = make([]string, 0, 1)
		signatureWithDifferentOptions = append(signatureWithDifferentOptions, emitSignatures[path].signature)
	} else if b.differsInOptions {
		signatureWithDifferentOptions = make([]string, 0, 1)
		signatureWithDifferentOptions = append(signatureWithDifferentOptions, b.signature)
	} else {
		signature = b.signature
	}
	return &emitSignature{
		signature:                     signature,
		signatureWithDifferentOptions: signatureWithDifferentOptions,
	}
}

func (b *BuildInfoEmitSignature) MarshalJSON() ([]byte, error) {
	if b.noEmitSignature() {
		return json.Marshal(b.fileId)
	}
	fileIdAndSignature := make([]any, 2)
	fileIdAndSignature[0] = b.fileId
	var signature any
	if b.differsOnlyInDtsMap {
		signature = []string{}
	} else if b.differsInOptions {
		signature = []string{b.signature}
	} else {
		signature = b.signature
	}
	fileIdAndSignature[1] = signature
	return json.Marshal(fileIdAndSignature)
}

func (b *BuildInfoEmitSignature) UnmarshalJSON(data []byte) error {
	var fileId BuildInfoFileId
	if err := json.Unmarshal(data, &fileId); err == nil {
		*b = BuildInfoEmitSignature{
			fileId: fileId,
		}
		return nil
	}
	var fileIdAndSignature []any
	if err := json.Unmarshal(data, &fileIdAndSignature); err == nil {
		if len(fileIdAndSignature) == 2 {
			var fileId BuildInfoFileId
			if id, ok := fileIdAndSignature[0].(float64); ok {
				fileId = BuildInfoFileId(id)
			} else {
				return fmt.Errorf("invalid fileId in incrementalBuildInfoEmitSignature: expected float64, got %T", fileIdAndSignature[0])
			}
			var signature string
			var differsOnlyInDtsMap, differsInOptions bool
			if signatureV, ok := fileIdAndSignature[1].(string); !ok {
				if signatureList, ok := fileIdAndSignature[1].([]string); ok {
					if len(signatureList) == 0 {
						differsOnlyInDtsMap = true
					} else if len(signatureList) == 1 {
						signature = signatureList[0]
						differsInOptions = true
					} else {
						return fmt.Errorf("invalid signature in incrementalBuildInfoEmitSignature: expected string or []string with 0 or 1 element, got %d elements", len(signatureList))
					}
				} else {
					return fmt.Errorf("invalid signature in incrementalBuildInfoEmitSignature: expected string or []string, got %T", fileIdAndSignature[1])
				}
			} else {
				signature = signatureV
			}
			*b = BuildInfoEmitSignature{
				fileId:              fileId,
				signature:           signature,
				differsOnlyInDtsMap: differsOnlyInDtsMap,
				differsInOptions:    differsInOptions,
			}
			return nil
		}
		return fmt.Errorf("invalid incrementalBuildInfoEmitSignature: expected 2 elements, got %d", len(fileIdAndSignature))
	}
	return fmt.Errorf("invalid IncrementalBuildInfoDiagnostic: %s", data)
}

type BuildInfo struct {
	Version string

	// Common between incremental and tsc -b buildinfo for non incremental programs
	Errors       bool `json:"errors,omitzero"`
	CheckPending bool `json:"checkPending,omitzero"`
	// Root         []BuildInfoRoot `json:"root,omitzero"`

	// IncrementalProgram info
	FileNames                  []string                             `json:"fileNames,omitzero"`
	FileInfos                  []*BuildInfoFileInfo                 `json:"fileInfos,omitzero"`
	FileIdsList                [][]BuildInfoFileId                  `json:"fileIdsList,omitzero"`
	Options                    *collections.OrderedMap[string, any] `json:"options,omitzero"`
	ReferencedMap              []*BuildInfoReferenceMapEntry        `json:"referencedMap,omitzero"`
	SemanticDiagnosticsPerFile []*BuildInfoSemanticDiagnostic       `json:"semanticDiagnosticsPerFile,omitzero"`
	EmitDiagnosticsPerFile     []*BuildInfoDiagnosticsOfFile        `json:"emitDiagnosticsPerFile,omitzero"`
	ChangeFileSet              []BuildInfoFileId                    `json:"changeFileSet,omitzero"`
	AffectedFilesPendingEmit   []*BuildInfoFilePendingEmit          `json:"affectedFilesPendingEmit,omitzero"`
	LatestChangedDtsFile       string                               `json:"latestChangedDtsFile,omitzero"` // Because this is only output file in the program, we dont need fileId to deduplicate name
	EmitSignatures             []*BuildInfoEmitSignature            `json:"emitSignatures,omitzero"`
	// resolvedRoot: readonly IncrementalBuildInfoResolvedRoot[] | undefined;
}

func (b *BuildInfo) IsValidVersion() bool {
	return b.Version == core.Version()
}

func (b *BuildInfo) IsIncremental() bool {
	return b != nil && len(b.FileNames) != 0
}
