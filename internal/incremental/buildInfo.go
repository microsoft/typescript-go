package incremental

import (
	"encoding/json"
	"fmt"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler"
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

func (b *BuildInfoFileInfo) HasSignature() bool {
	return b.signature != ""
}

func (b *BuildInfoFileInfo) GetFileInfo() *fileInfo {
	if b.signature != "" {
		return &fileInfo{
			version:   b.signature,
			signature: b.signature,
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

type BuildInfoCompilerOption struct {
	name  string
	value any
}

func (b *BuildInfoCompilerOption) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{b.name, b.value})
}

func (b *BuildInfoCompilerOption) UnmarshalJSON(data []byte) error {
	var nameAndValue []any
	if err := json.Unmarshal(data, &nameAndValue); err != nil {
		return err
	}
	if len(nameAndValue) != 2 {
		return fmt.Errorf("invalid incrementalBuildInfoCompilerOption: expected array of length 2, got %d", len(nameAndValue))
	}
	if name, ok := nameAndValue[0].(string); ok {
		*b = BuildInfoCompilerOption{}
		b.name = name
		b.value = nameAndValue[1]
		return nil
	}
	return fmt.Errorf("invalid name in incrementalBuildInfoCompilerOption: expected string, got %T", nameAndValue[0])
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
	// false if diagnostic is not for a file,
	// incrementalBuildInfoFileId if it is for a file thats other than its stored for
	file               any
	loc                core.TextRange
	code               int32
	category           diagnostics.Category
	message            string
	messageChain       []*BuildInfoDiagnostic
	relatedInformation []*BuildInfoDiagnostic
	reportsUnnecessary bool
	reportsDeprecated  bool
	skippedOnNoEmit    bool
}

func (b *BuildInfoDiagnostic) toDiagnostic(p *compiler.Program, file *ast.SourceFile) *ast.Diagnostic {
	var fileForDiagnostic *ast.SourceFile
	if b.file != false {
		if b.file == nil {
			fileForDiagnostic = file
		} else {
			fileForDiagnostic = p.GetSourceFileByPath(tspath.Path(b.file.(string)))
		}
	}
	var messageChain []*ast.Diagnostic
	for _, msg := range b.messageChain {
		messageChain = append(messageChain, msg.toDiagnostic(p, fileForDiagnostic))
	}
	var relatedInformation []*ast.Diagnostic
	for _, info := range b.relatedInformation {
		relatedInformation = append(relatedInformation, info.toDiagnostic(p, fileForDiagnostic))
	}
	return ast.NewDiagnosticWith(
		fileForDiagnostic,
		b.loc,
		b.code,
		b.category,
		b.message,
		messageChain,
		relatedInformation,
		b.reportsUnnecessary,
		b.reportsDeprecated,
		b.skippedOnNoEmit,
	)
}

func (b *BuildInfoDiagnostic) MarshalJSON() ([]byte, error) {
	info := map[string]any{}
	if b.file != "" {
		info["file"] = b.file
		info["pos"] = b.loc.Pos()
		info["end"] = b.loc.End()
	}
	info["code"] = b.code
	info["category"] = b.category
	info["message"] = b.message
	if len(b.messageChain) > 0 {
		info["messageChain"] = b.messageChain
	}
	if len(b.relatedInformation) > 0 {
		info["relatedInformation"] = b.relatedInformation
	}
	if b.reportsUnnecessary {
		info["reportsUnnecessary"] = b.reportsUnnecessary
	}
	if b.reportsDeprecated {
		info["reportsDeprecated"] = b.reportsDeprecated
	}
	if b.skippedOnNoEmit {
		info["skippedOnNoEmit"] = b.skippedOnNoEmit
	}
	return json.Marshal(info)
}

func (b *BuildInfoDiagnostic) UnmarshalJSON(data []byte) error {
	var vIncrementalBuildInfoDiagnostic map[string]any
	if err := json.Unmarshal(data, &vIncrementalBuildInfoDiagnostic); err != nil {
		return fmt.Errorf("invalid incrementalBuildInfoDiagnostic: %s", data)
	}

	*b = BuildInfoDiagnostic{}
	if file, ok := vIncrementalBuildInfoDiagnostic["file"]; ok {
		if _, ok := file.(float64); !ok {
			if value, ok := file.(bool); !ok || value {
				return fmt.Errorf("invalid file in incrementalBuildInfoDiagnostic: expected false or float64, got %T", file)
			}
		}
		b.file = file
		var pos float64
		posV, ok := vIncrementalBuildInfoDiagnostic["pos"]
		if ok {
			pos, ok = posV.(float64)
			if !ok {
				return fmt.Errorf("invalid pos in incrementalBuildInfoDiagnostic: expected float64, got %T", posV)
			}
		} else {
			return fmt.Errorf("missing pos in incrementalBuildInfoDiagnostic")
		}
		var end float64
		endv, ok := vIncrementalBuildInfoDiagnostic["end"]
		if ok {
			end, ok = endv.(float64)
			if !ok {
				return fmt.Errorf("invalid end in incrementalBuildInfoDiagnostic: expected float64, got %T", endv)
			}
		} else {
			return fmt.Errorf("missing end in incrementalBuildInfoDiagnostic")
		}
		b.loc = core.NewTextRange(int(pos), int(end))
	}
	if codeV, ok := vIncrementalBuildInfoDiagnostic["code"]; ok {
		code, ok := codeV.(float64)
		if !ok {
			return fmt.Errorf("invalid code in incrementalBuildInfoDiagnostic: expected float64, got %T", codeV)
		}
		b.code = int32(code)
	} else {
		return fmt.Errorf("missing code in incrementalBuildInfoDiagnostic")
	}
	if categoryV, ok := vIncrementalBuildInfoDiagnostic["category"]; ok {
		category, ok := categoryV.(float64)
		if !ok {
			return fmt.Errorf("invalid category in incrementalBuildInfoDiagnostic: expected float64, got %T", categoryV)
		}
		if category < 0 || category > float64(diagnostics.CategoryMessage) {
			return fmt.Errorf("invalid category in incrementalBuildInfoDiagnostic: %f is out of range", category)
		}
		b.category = diagnostics.Category(category)
	} else {
		return fmt.Errorf("missing category in incrementalBuildInfoDiagnostic")
	}
	if messageV, ok := vIncrementalBuildInfoDiagnostic["message"]; ok {
		if b.message, ok = messageV.(string); !ok {
			return fmt.Errorf("invalid message in incrementalBuildInfoDiagnostic: expected string, got %T", messageV)
		}
	} else {
		return fmt.Errorf("missing message in incrementalBuildInfoDiagnostic")
	}
	if messageChain, ok := vIncrementalBuildInfoDiagnostic["messageChain"]; ok {
		if messages, ok := messageChain.([]any); ok {
			b.messageChain = make([]*BuildInfoDiagnostic, len(messages))
			for _, msg := range messages {
				var diagnostic BuildInfoDiagnostic
				if err := json.Unmarshal([]byte(msg.(string)), &diagnostic); err != nil {
					return fmt.Errorf("invalid messageChain in incrementalBuildInfoDiagnostic: %s", msg)
				}
				b.messageChain = append(b.messageChain, &diagnostic)
			}
		}
	}
	if relatedInformation, ok := vIncrementalBuildInfoDiagnostic["relatedInformation"]; ok {
		if infos, ok := relatedInformation.([]any); ok {
			b.relatedInformation = make([]*BuildInfoDiagnostic, len(infos))
			for _, info := range infos {
				var diagnostic BuildInfoDiagnostic
				if err := json.Unmarshal([]byte(info.(string)), &diagnostic); err != nil {
					return fmt.Errorf("invalid relatedInformation in incrementalBuildInfoDiagnostic: %s", info)
				}
				b.relatedInformation = append(b.relatedInformation, &diagnostic)
			}
		}
	}
	if reportsUnnecessary, ok := vIncrementalBuildInfoDiagnostic["reportsUnnecessary"]; ok {
		if b.reportsUnnecessary, ok = reportsUnnecessary.(bool); !ok {
			return fmt.Errorf("invalid reportsUnnecessary in incrementalBuildInfoDiagnostic: expected boolean, got %T", reportsUnnecessary)
		}
	}
	if reportsDeprecated, ok := vIncrementalBuildInfoDiagnostic["reportsDeprecated"]; ok {
		if b.reportsDeprecated, ok = reportsDeprecated.(bool); !ok {
			return fmt.Errorf("invalid reportsDeprecated in incrementalBuildInfoDiagnostic: expected boolean, got %T", reportsDeprecated)
		}
	}
	if skippedOnNoEmit, ok := vIncrementalBuildInfoDiagnostic["skippedOnNoEmit"]; ok {
		if b.skippedOnNoEmit, ok = skippedOnNoEmit.(bool); !ok {
			return fmt.Errorf("invalid skippedOnNoEmit in incrementalBuildInfoDiagnostic: expected boolean, got %T", skippedOnNoEmit)
		}
	}
	return nil
}

type BuildInfoDiagnosticOfFile struct {
	fileId      BuildInfoFileId
	diagnostics []*BuildInfoDiagnostic
}

func (b *BuildInfoDiagnosticOfFile) MarshalJSON() ([]byte, error) {
	fileIdAndDiagnostics := make([]any, 0, 2)
	fileIdAndDiagnostics = append(fileIdAndDiagnostics, b.fileId)
	fileIdAndDiagnostics = append(fileIdAndDiagnostics, b.diagnostics)
	return json.Marshal(fileIdAndDiagnostics)
}

func (b *BuildInfoDiagnosticOfFile) UnmarshalJSON(data []byte) error {
	var fileIdAndDiagnostics []any
	if err := json.Unmarshal(data, &fileIdAndDiagnostics); err != nil {
		return fmt.Errorf("invalid IncrementalBuildInfoDiagnostic: %s", data)
	}
	if len(fileIdAndDiagnostics) != 2 {
		return fmt.Errorf("invalid IncrementalBuildInfoDiagnostic: expected 2 elements, got %d", len(fileIdAndDiagnostics))
	}
	var fileId BuildInfoFileId
	if fileIdV, ok := fileIdAndDiagnostics[0].(float64); !ok {
		return fmt.Errorf("invalid fileId in IncrementalBuildInfoDiagnostic: expected float64, got %T", fileIdAndDiagnostics[0])
	} else {
		fileId = BuildInfoFileId(fileIdV)
	}
	if diagnostics, ok := fileIdAndDiagnostics[1].([]*BuildInfoDiagnostic); ok {
		*b = BuildInfoDiagnosticOfFile{
			fileId:      fileId,
			diagnostics: diagnostics,
		}
		return nil
	}
	return fmt.Errorf("invalid diagnostics in IncrementalBuildInfoDiagnostic: expected []*incrementalBuildInfoDiagnostic, got %T", fileIdAndDiagnostics[1])
}

type BuildInfoSemanticDiagnostic struct {
	FileId     BuildInfoFileId           // File is not in changedSet and still doesnt have cached diagnostics
	Diagnostic BuildInfoDiagnosticOfFile // Diagnostics for file
}

func (b *BuildInfoSemanticDiagnostic) MarshalJSON() ([]byte, error) {
	if b.FileId != 0 {
		return json.Marshal(b.FileId)
	}
	return json.Marshal(b.Diagnostic)
}

func (b *BuildInfoSemanticDiagnostic) UnmarshalJSON(data []byte) error {
	var fileId BuildInfoFileId
	if err := json.Unmarshal(data, &fileId); err == nil {
		*b = BuildInfoSemanticDiagnostic{
			FileId: fileId,
		}
		return nil
	}
	var diagnostic BuildInfoDiagnosticOfFile
	if err := json.Unmarshal(data, &diagnostic); err == nil {
		*b = BuildInfoSemanticDiagnostic{
			Diagnostic: diagnostic,
		}
		return nil
	}
	return fmt.Errorf("invalid IncrementalBuildInfoDiagnostic: %s", data)
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
	// Root         []BuildInfoRoot `json:"root,omitempty,omitzero"`

	// IncrementalProgram info
	FileNames                  []string                      `json:"fileNames,omitzero"`
	FileInfos                  []*BuildInfoFileInfo          `json:"fileInfos,omitzero"`
	FileIdsList                [][]BuildInfoFileId           `json:"fileIdsList,omitzero"`
	Options                    []BuildInfoCompilerOption     `json:"options,omitzero"`
	ReferencedMap              []BuildInfoReferenceMapEntry  `json:"referencedMap,omitzero"`
	SemanticDiagnosticsPerFile []BuildInfoSemanticDiagnostic `json:"semanticDiagnosticsPerFile,omitzero"`
	EmitDiagnosticsPerFile     []BuildInfoDiagnosticOfFile   `json:"emitDiagnosticsPerFile,omitzero"`
	ChangeFileSet              []BuildInfoFileId             `json:"changeFileSet,omitzero"`
	AffectedFilesPendingEmit   []BuildInfoFilePendingEmit    `json:"affectedFilesPendingEmit,omitzero"`
	LatestChangedDtsFile       string                        `json:"latestChangedDtsFile,omitzero"` // Because this is only output file in the program, we dont need fileId to deduplicate name
	EmitSignatures             []BuildInfoEmitSignature      `json:"emitSignatures,omitzero"`
	// resolvedRoot: readonly IncrementalBuildInfoResolvedRoot[] | undefined;
}

func (b *BuildInfo) IsValidVersion() bool {
	return b.Version == core.Version()
}

func (b *BuildInfo) IsIncremental() bool {
	return b != nil && len(b.FileNames) != 0
}
