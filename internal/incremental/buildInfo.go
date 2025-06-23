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
	incrementalBuildInfoFileId       int
	incrementalBuildInfoFileIdListId int
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

/**
*  - string if FileInfo.version === FileInfo.signature && !FileInfo.affectsGlobalScope
*  - otherwise encoded FileInfo
*    Signature is
* 	 - undefined if FileInfo.version === FileInfo.signature
* 	 - false if FileInfo has signature as undefined (not calculated)
* 	 - string actual signature
 */
func (f *fileInfo) MarshalJSON() ([]byte, error) {
	if f.version == f.signature && !f.affectsGlobalScope && f.impliedNodeFormat == core.ResolutionModeNone {
		return json.Marshal(f.version)
	}
	info := map[string]any{}
	if f.version != "" {
		info["version"] = f.version
	}
	if f.signature != f.version {
		if f.signature == "" {
			info["signature"] = false
		} else {
			info["signature"] = f.signature
		}
	}
	if f.affectsGlobalScope {
		info["affectsGlobalScope"] = f.affectsGlobalScope
	}
	if f.impliedNodeFormat != core.ResolutionModeNone {
		info["impliedNodeFormat"] = f.impliedNodeFormat
	}
	return json.Marshal(info)
}

func (f *fileInfo) UnmarshalJSON(data []byte) error {
	var version string
	if err := json.Unmarshal(data, &version); err == nil {
		*f = fileInfo{
			version:   version,
			signature: version,
		}
		return nil
	}

	var vFileInfo map[string]any
	if err := json.Unmarshal(data, &vFileInfo); err != nil {
		return fmt.Errorf("invalid fileInfo: %s", data)
	}

	*f = fileInfo{}
	if version, ok := vFileInfo["version"]; ok {
		f.version = version.(string)
	}
	if signature, ok := vFileInfo["signature"]; ok {
		if signature == false {
			f.signature = ""
		} else if f.signature, ok = signature.(string); !ok {
			return fmt.Errorf("invalid signature in fileInfo: expected string or false, got %T", signature)
		}
	} else {
		f.signature = f.version // default to version if no signature provided
	}
	if affectsGlobalScope, ok := vFileInfo["affectsGlobalScope"]; ok {
		if f.affectsGlobalScope, ok = affectsGlobalScope.(bool); !ok {
			return fmt.Errorf("invalid affectsGlobalScope in fileInfo: expected bool, got %T", affectsGlobalScope)
		}
	}
	if impliedNodeFormatV, ok := vFileInfo["impliedNodeFormat"]; ok {
		if impliedNodeFormat, ok := impliedNodeFormatV.(int); ok {
			if impliedNodeFormat != int(core.ResolutionModeCommonJS) && impliedNodeFormat != int(core.ResolutionModeESM) {
				return fmt.Errorf("invalid impliedNodeFormat in fileInfo: %d is out of range", impliedNodeFormat)
			}
			f.impliedNodeFormat = core.ResolutionMode(impliedNodeFormat)
		} else {
			return fmt.Errorf("invalid impliedNodeFormat in fileInfo: expected int, got %T", impliedNodeFormatV)
		}
	}
	return nil
}

type incrementalBuildInfoCompilerOption struct {
	name  string
	value any
}

func (i *incrementalBuildInfoCompilerOption) MarshalJSON() ([]byte, error) {
	nameAndValue := make([]any, 2)
	nameAndValue[0] = i.name
	nameAndValue[1] = i.value
	return json.Marshal(nameAndValue)
}

func (i *incrementalBuildInfoCompilerOption) UnmarshalJSON(data []byte) error {
	var nameAndValue []any
	if err := json.Unmarshal(data, &nameAndValue); err != nil {
		return err
	}
	if len(nameAndValue) != 2 {
		return fmt.Errorf("invalid incrementalBuildInfoCompilerOption: expected array of length 2, got %d", len(nameAndValue))
	}
	if name, ok := nameAndValue[0].(string); ok {
		*i = incrementalBuildInfoCompilerOption{}
		i.name = name
		i.value = nameAndValue[1]
		return nil
	}
	return fmt.Errorf("invalid name in incrementalBuildInfoCompilerOption: expected string, got %T", nameAndValue[0])
}

type incrementalBuildInfoReferenceMapEntry struct {
	fileId       incrementalBuildInfoFileId
	fileIdListId incrementalBuildInfoFileIdListId
}

func (i *incrementalBuildInfoReferenceMapEntry) MarshalJSON() ([]byte, error) {
	return json.Marshal([2]int{int(i.fileId), int(i.fileIdListId)})
}

func (i *incrementalBuildInfoReferenceMapEntry) UnmarshalJSON(data []byte) error {
	var v *[2]int
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*i = incrementalBuildInfoReferenceMapEntry{
		fileId:       incrementalBuildInfoFileId(v[0]),
		fileIdListId: incrementalBuildInfoFileIdListId(v[1]),
	}
	return nil
}

type incrementalBuildInfoDiagnostic struct {
	// false if diagnostic is not for a file,
	// incrementalBuildInfoFileId if it is for a file thats other than its stored for
	file               any
	loc                core.TextRange
	code               int32
	category           diagnostics.Category
	message            string
	messageChain       []*incrementalBuildInfoDiagnostic
	relatedInformation []*incrementalBuildInfoDiagnostic
	reportsUnnecessary bool
	reportsDeprecated  bool
	skippedOnNoEmit    bool
}

func (i *incrementalBuildInfoDiagnostic) toDiagnostic(p *compiler.Program, file *ast.SourceFile) *ast.Diagnostic {
	var fileForDiagnostic *ast.SourceFile
	if i.file != false {
		if i.file == nil {
			fileForDiagnostic = file
		} else {
			fileForDiagnostic = p.GetSourceFileByPath(tspath.Path(i.file.(string)))
		}
	}
	var messageChain []*ast.Diagnostic
	for _, msg := range i.messageChain {
		messageChain = append(messageChain, msg.toDiagnostic(p, fileForDiagnostic))
	}
	var relatedInformation []*ast.Diagnostic
	for _, info := range i.relatedInformation {
		relatedInformation = append(relatedInformation, info.toDiagnostic(p, fileForDiagnostic))
	}
	return ast.NewDiagnosticWith(
		fileForDiagnostic,
		i.loc,
		i.code,
		i.category,
		i.message,
		messageChain,
		relatedInformation,
		i.reportsUnnecessary,
		i.reportsDeprecated,
		i.skippedOnNoEmit,
	)
}

func (i *incrementalBuildInfoDiagnostic) MarshalJSON() ([]byte, error) {
	info := map[string]any{}
	if i.file != "" {
		info["file"] = i.file
		info["pos"] = i.loc.Pos()
		info["end"] = i.loc.End()
	}
	info["code"] = i.code
	info["category"] = i.category
	info["message"] = i.message
	if len(i.messageChain) > 0 {
		info["messageChain"] = i.messageChain
	}
	if len(i.relatedInformation) > 0 {
		info["relatedInformation"] = i.relatedInformation
	}
	if i.reportsUnnecessary {
		info["reportsUnnecessary"] = i.reportsUnnecessary
	}
	if i.reportsDeprecated {
		info["reportsDeprecated"] = i.reportsDeprecated
	}
	if i.skippedOnNoEmit {
		info["skippedOnNoEmit"] = i.skippedOnNoEmit
	}
	return json.Marshal(info)
}

func (i *incrementalBuildInfoDiagnostic) UnmarshalJSON(data []byte) error {
	var vIncrementalBuildInfoDiagnostic map[string]any
	if err := json.Unmarshal(data, &vIncrementalBuildInfoDiagnostic); err != nil {
		return fmt.Errorf("invalid incrementalBuildInfoDiagnostic: %s", data)
	}

	*i = incrementalBuildInfoDiagnostic{}
	if file, ok := vIncrementalBuildInfoDiagnostic["file"]; ok {
		if _, ok := file.(float64); !ok {
			if value, ok := file.(bool); !ok || value {
				return fmt.Errorf("invalid file in incrementalBuildInfoDiagnostic: expected false or float64, got %T", file)
			}
		}
		i.file = file
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
		i.loc = core.NewTextRange(int(pos), int(end))
	}
	if codeV, ok := vIncrementalBuildInfoDiagnostic["code"]; ok {
		code, ok := codeV.(float64)
		if !ok {
			return fmt.Errorf("invalid code in incrementalBuildInfoDiagnostic: expected float64, got %T", codeV)
		}
		i.code = int32(code)
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
		i.category = diagnostics.Category(category)
	} else {
		return fmt.Errorf("missing category in incrementalBuildInfoDiagnostic")
	}
	if messageV, ok := vIncrementalBuildInfoDiagnostic["message"]; ok {
		if i.message, ok = messageV.(string); !ok {
			return fmt.Errorf("invalid message in incrementalBuildInfoDiagnostic: expected string, got %T", messageV)
		}
	} else {
		return fmt.Errorf("missing message in incrementalBuildInfoDiagnostic")
	}
	if messageChain, ok := vIncrementalBuildInfoDiagnostic["messageChain"]; ok {
		if messages, ok := messageChain.([]any); ok {
			i.messageChain = make([]*incrementalBuildInfoDiagnostic, len(messages))
			for _, msg := range messages {
				var diagnostic incrementalBuildInfoDiagnostic
				if err := json.Unmarshal([]byte(msg.(string)), &diagnostic); err != nil {
					return fmt.Errorf("invalid messageChain in incrementalBuildInfoDiagnostic: %s", msg)
				}
				i.messageChain = append(i.messageChain, &diagnostic)
			}
		}
	}
	if relatedInformation, ok := vIncrementalBuildInfoDiagnostic["relatedInformation"]; ok {
		if infos, ok := relatedInformation.([]any); ok {
			i.relatedInformation = make([]*incrementalBuildInfoDiagnostic, len(infos))
			for _, info := range infos {
				var diagnostic incrementalBuildInfoDiagnostic
				if err := json.Unmarshal([]byte(info.(string)), &diagnostic); err != nil {
					return fmt.Errorf("invalid relatedInformation in incrementalBuildInfoDiagnostic: %s", info)
				}
				i.relatedInformation = append(i.relatedInformation, &diagnostic)
			}
		}
	}
	if reportsUnnecessary, ok := vIncrementalBuildInfoDiagnostic["reportsUnnecessary"]; ok {
		if i.reportsUnnecessary, ok = reportsUnnecessary.(bool); !ok {
			return fmt.Errorf("invalid reportsUnnecessary in incrementalBuildInfoDiagnostic: expected boolean, got %T", reportsUnnecessary)
		}
	}
	if reportsDeprecated, ok := vIncrementalBuildInfoDiagnostic["reportsDeprecated"]; ok {
		if i.reportsDeprecated, ok = reportsDeprecated.(bool); !ok {
			return fmt.Errorf("invalid reportsDeprecated in incrementalBuildInfoDiagnostic: expected boolean, got %T", reportsDeprecated)
		}
	}
	if skippedOnNoEmit, ok := vIncrementalBuildInfoDiagnostic["skippedOnNoEmit"]; ok {
		if i.skippedOnNoEmit, ok = skippedOnNoEmit.(bool); !ok {
			return fmt.Errorf("invalid skippedOnNoEmit in incrementalBuildInfoDiagnostic: expected boolean, got %T", skippedOnNoEmit)
		}
	}
	return nil
}

type incrementalBuildInfoDiagnosticOfFile struct {
	fileId      incrementalBuildInfoFileId
	diagnostics []*incrementalBuildInfoDiagnostic
}

func (i *incrementalBuildInfoDiagnosticOfFile) MarshalJSON() ([]byte, error) {
	fileIdAndDiagnostics := make([]any, 0, 2)
	fileIdAndDiagnostics = append(fileIdAndDiagnostics, i.fileId)
	fileIdAndDiagnostics = append(fileIdAndDiagnostics, i.diagnostics)
	return json.Marshal(fileIdAndDiagnostics)
}

func (i *incrementalBuildInfoDiagnosticOfFile) UnmarshalJSON(data []byte) error {
	var fileIdAndDiagnostics []any
	if err := json.Unmarshal(data, &fileIdAndDiagnostics); err != nil {
		return fmt.Errorf("invalid IncrementalBuildInfoDiagnostic: %s", data)
	}
	if len(fileIdAndDiagnostics) != 2 {
		return fmt.Errorf("invalid IncrementalBuildInfoDiagnostic: expected 2 elements, got %d", len(fileIdAndDiagnostics))
	}
	var fileId incrementalBuildInfoFileId
	if fileIdV, ok := fileIdAndDiagnostics[0].(float64); !ok {
		return fmt.Errorf("invalid fileId in IncrementalBuildInfoDiagnostic: expected float64, got %T", fileIdAndDiagnostics[0])
	} else {
		fileId = incrementalBuildInfoFileId(fileIdV)
	}
	if diagnostics, ok := fileIdAndDiagnostics[1].([]*incrementalBuildInfoDiagnostic); ok {
		*i = incrementalBuildInfoDiagnosticOfFile{
			fileId:      fileId,
			diagnostics: diagnostics,
		}
		return nil
	}
	return fmt.Errorf("invalid diagnostics in IncrementalBuildInfoDiagnostic: expected []*incrementalBuildInfoDiagnostic, got %T", fileIdAndDiagnostics[1])
}

type incrementalBuildInfoSemanticDiagnostic struct {
	fileId     incrementalBuildInfoFileId           // File is not in changedSet and still doesnt have cached diagnostics
	diagnostic incrementalBuildInfoDiagnosticOfFile // Diagnostics for file
}

func (i *incrementalBuildInfoSemanticDiagnostic) MarshalJSON() ([]byte, error) {
	if i.fileId != 0 {
		return json.Marshal(i.fileId)
	}
	return json.Marshal(i.diagnostic)
}

func (i *incrementalBuildInfoSemanticDiagnostic) UnmarshalJSON(data []byte) error {
	var fileId incrementalBuildInfoFileId
	if err := json.Unmarshal(data, &fileId); err == nil {
		*i = incrementalBuildInfoSemanticDiagnostic{
			fileId: fileId,
		}
		return nil
	}
	var diagnostic incrementalBuildInfoDiagnosticOfFile
	if err := json.Unmarshal(data, &diagnostic); err == nil {
		*i = incrementalBuildInfoSemanticDiagnostic{
			diagnostic: diagnostic,
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
type incrementalBuildInfoFilePendingEmit struct {
	fileId   incrementalBuildInfoFileId
	emitKind fileEmitKind
}

func (i *incrementalBuildInfoFilePendingEmit) MarshalJSON() ([]byte, error) {
	if i.emitKind == 0 {
		return json.Marshal(i.fileId)
	}
	if i.emitKind == fileEmitKindDts {
		fileListIds := []incrementalBuildInfoFileId{i.fileId}
		return json.Marshal(fileListIds)
	}
	fileAndEmitKind := []int{int(i.fileId), int(i.emitKind)}
	return json.Marshal(fileAndEmitKind)
}

func (i *incrementalBuildInfoFilePendingEmit) UnmarshalJSON(data []byte) error {
	var fileId incrementalBuildInfoFileId
	if err := json.Unmarshal(data, &fileId); err == nil {
		*i = incrementalBuildInfoFilePendingEmit{
			fileId: fileId,
		}
		return nil
	}
	var intTuple []int
	if err := json.Unmarshal(data, &intTuple); err == nil {
		if len(intTuple) == 1 {
			*i = incrementalBuildInfoFilePendingEmit{
				fileId:   incrementalBuildInfoFileId(intTuple[0]),
				emitKind: fileEmitKindDts,
			}
			return nil
		} else if len(intTuple) == 2 {
			*i = incrementalBuildInfoFilePendingEmit{
				fileId:   incrementalBuildInfoFileId(intTuple[0]),
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
type incrementalBuildInfoEmitSignature struct {
	fileId incrementalBuildInfoFileId
	// Signature if it is different from file's signature
	signature           string
	differsOnlyInDtsMap bool // true if signature is different only in dtsMap value
	differsInOptions    bool // true if signature is different in options used to emit file
}

func (i *incrementalBuildInfoEmitSignature) noEmitSignature() bool {
	return i.signature == "" && !i.differsOnlyInDtsMap && !i.differsInOptions
}

func (i *incrementalBuildInfoEmitSignature) toEmitSignature(path tspath.Path, emitSignatures map[tspath.Path]*emitSignature) *emitSignature {
	var signature string
	var signatureWithDifferentOptions []string
	if i.differsOnlyInDtsMap {
		signatureWithDifferentOptions = make([]string, 0, 1)
		signatureWithDifferentOptions = append(signatureWithDifferentOptions, emitSignatures[path].signature)
	} else if i.differsInOptions {
		signatureWithDifferentOptions = make([]string, 0, 1)
		signatureWithDifferentOptions = append(signatureWithDifferentOptions, i.signature)
	} else {
		signature = i.signature
	}
	return &emitSignature{
		signature:                     signature,
		signatureWithDifferentOptions: signatureWithDifferentOptions,
	}
}

func (i *incrementalBuildInfoEmitSignature) MarshalJSON() ([]byte, error) {
	if i.noEmitSignature() {
		return json.Marshal(i.fileId)
	}
	fileIdAndSignature := make([]any, 2)
	fileIdAndSignature[0] = i.fileId
	var signature any
	if i.differsOnlyInDtsMap {
		signature = []string{}
	} else if i.differsInOptions {
		signature = []string{i.signature}
	} else {
		signature = i.signature
	}
	fileIdAndSignature[1] = signature
	return json.Marshal(fileIdAndSignature)
}

func (i *incrementalBuildInfoEmitSignature) UnmarshalJSON(data []byte) error {
	var fileId incrementalBuildInfoFileId
	if err := json.Unmarshal(data, &fileId); err == nil {
		*i = incrementalBuildInfoEmitSignature{
			fileId: fileId,
		}
		return nil
	}
	var fileIdAndSignature []any
	if err := json.Unmarshal(data, &fileIdAndSignature); err == nil {
		if len(fileIdAndSignature) == 2 {
			var fileId incrementalBuildInfoFileId
			if id, ok := fileIdAndSignature[0].(float64); ok {
				fileId = incrementalBuildInfoFileId(id)
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
			*i = incrementalBuildInfoEmitSignature{
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
	FileNames                  []string                                 `json:"fileNames,omitzero"`
	FileInfos                  []*fileInfo                              `json:"fileInfos,omitzero"`
	FileIdsList                [][]incrementalBuildInfoFileId           `json:"fileIdsList,omitzero"`
	Options                    []incrementalBuildInfoCompilerOption     `json:"options,omitzero"`
	ReferencedMap              []incrementalBuildInfoReferenceMapEntry  `json:"referencedMap,omitzero"`
	SemanticDiagnosticsPerFile []incrementalBuildInfoSemanticDiagnostic `json:"semanticDiagnosticsPerFile,omitzero"`
	EmitDiagnosticsPerFile     []incrementalBuildInfoDiagnosticOfFile   `json:"emitDiagnosticsPerFile,omitzero"`
	ChangeFileSet              []incrementalBuildInfoFileId             `json:"changeFileSet,omitzero"`
	AffectedFilesPendingEmit   []incrementalBuildInfoFilePendingEmit    `json:"affectedFilesPendingEmit,omitzero"`
	LatestChangedDtsFile       string                                   `json:"latestChangedDtsFile,omitzero"` // Because this is only output file in the program, we dont need fileId to deduplicate name
	EmitSignatures             []incrementalBuildInfoEmitSignature      `json:"emitSignatures,omitzero"`
	// resolvedRoot: readonly IncrementalBuildInfoResolvedRoot[] | undefined;
}

func (b *BuildInfo) IsValidVersion() bool {
	return b.Version == core.Version()
}

func (b *BuildInfo) IsIncremental() bool {
	return b != nil && len(b.FileNames) != 0
}
