package sourcemap

import (
	"sort"

	"github.com/go-json-experiment/json"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tspath"
)

// DocumentPosition represents a position in a document
type DocumentPosition struct {
	FileName string
	Pos      int
}

// DocumentPositionMapper maps positions between source and generated files
type DocumentPositionMapper interface {
	GetSourcePosition(input DocumentPosition) DocumentPosition
	GetGeneratedPosition(input DocumentPosition) DocumentPosition
}

// DocumentPositionMapperHost provides file system access for position mapping
type DocumentPositionMapperHost interface {
	GetSourceFileLike(fileName string) SourceFileLike
	GetCanonicalFileName(path string) string
	Log(text string)
}

// SourceFileLike represents a file with text content and line information
type SourceFileLike interface {
	Text() string
	LineStarts() []core.TextPos
}


// MappedPosition represents a position mapping with additional metadata
type MappedPosition struct {
	GeneratedPosition int
	SourceIndex       int
	SourcePosition    int
	SourceFileName    string // Resolved source file name
}

// CreateDocumentPositionMapper creates a document position mapper from a source map
func CreateDocumentPositionMapper(host DocumentPositionMapperHost, mapContent string, mapPath string) DocumentPositionMapper {
	var sourceMap RawSourceMap
	if err := json.Unmarshal([]byte(mapContent), &sourceMap); err != nil {
		host.Log("Failed to parse source map: " + err.Error())
		return &identityMapper{}
	}

	return createDocumentPositionMapperFromRawMap(host, &sourceMap, mapPath)
}

func createDocumentPositionMapperFromRawMap(host DocumentPositionMapperHost, sourceMap *RawSourceMap, mapPath string) DocumentPositionMapper {
	mapDirectory := tspath.GetDirectoryPath(mapPath)
	var sourceRoot string
	if sourceMap.SourceRoot != "" {
		sourceRoot = tspath.GetNormalizedAbsolutePath(sourceMap.SourceRoot, mapDirectory)
	} else {
		sourceRoot = mapDirectory
	}
	
	generatedAbsoluteFilePath := tspath.GetNormalizedAbsolutePath(sourceMap.File, mapDirectory)
	generatedFile := host.GetSourceFileLike(generatedAbsoluteFilePath)
	
	sourceFileAbsolutePaths := make([]string, len(sourceMap.Sources))
	for i, source := range sourceMap.Sources {
		sourceFileAbsolutePaths[i] = tspath.GetNormalizedAbsolutePath(source, sourceRoot)
	}

	return &documentPositionMapper{
		host:                      host,
		sourceMap:                 sourceMap,
		mapPath:                   mapPath,
		generatedFile:             generatedFile,
		sourceFileAbsolutePaths:   sourceFileAbsolutePaths,
		sourceToSourceIndexMap:    createSourceIndexMap(sourceFileAbsolutePaths, host.GetCanonicalFileName),
	}
}

func createSourceIndexMap(sourceFileAbsolutePaths []string, getCanonicalFileName func(string) string) map[string]int {
	result := make(map[string]int, len(sourceFileAbsolutePaths))
	for i, source := range sourceFileAbsolutePaths {
		result[getCanonicalFileName(source)] = i
	}
	return result
}

type documentPositionMapper struct {
	host                    DocumentPositionMapperHost
	sourceMap               *RawSourceMap
	mapPath                 string
	generatedFile           SourceFileLike
	sourceFileAbsolutePaths []string
	sourceToSourceIndexMap  map[string]int
	
	// Cached mappings
	decodedMappings    []MappedPosition
	generatedMappings  []MappedPosition
	sourceMappings     [][]MappedPosition
	mappingsDecoded    bool
}

func (mapper *documentPositionMapper) GetSourcePosition(input DocumentPosition) DocumentPosition {
	if !mapper.ensureMappingsDecoded() {
		return input
	}

	if len(mapper.generatedMappings) == 0 {
		return input
	}

	targetIndex := sort.Search(len(mapper.generatedMappings), func(i int) bool {
		return mapper.generatedMappings[i].GeneratedPosition > input.Pos
	}) - 1

	if targetIndex < 0 {
		targetIndex = 0
	}
	if targetIndex >= len(mapper.generatedMappings) {
		targetIndex = len(mapper.generatedMappings) - 1
	}

	mapping := mapper.generatedMappings[targetIndex]
	if mapping.SourceIndex < 0 || mapping.SourceIndex >= len(mapper.sourceFileAbsolutePaths) {
		return input
	}

	return DocumentPosition{
		FileName: mapper.sourceFileAbsolutePaths[mapping.SourceIndex],
		Pos:      mapping.SourcePosition,
	}
}

func (mapper *documentPositionMapper) GetGeneratedPosition(input DocumentPosition) DocumentPosition {
	if !mapper.ensureMappingsDecoded() {
		return input
	}

	sourceIndex, exists := mapper.sourceToSourceIndexMap[mapper.host.GetCanonicalFileName(input.FileName)]
	if !exists {
		return input
	}

	if sourceIndex >= len(mapper.sourceMappings) || len(mapper.sourceMappings[sourceIndex]) == 0 {
		return input
	}

	sourceMappings := mapper.sourceMappings[sourceIndex]
	targetIndex := sort.Search(len(sourceMappings), func(i int) bool {
		return sourceMappings[i].SourcePosition > input.Pos
	}) - 1

	if targetIndex < 0 {
		targetIndex = 0
	}
	if targetIndex >= len(sourceMappings) {
		targetIndex = len(sourceMappings) - 1
	}

	mapping := sourceMappings[targetIndex]
	if mapping.SourceIndex != sourceIndex {
		return input
	}

	return DocumentPosition{
		FileName: tspath.GetNormalizedAbsolutePath(mapper.sourceMap.File, tspath.GetDirectoryPath(mapper.mapPath)),
		Pos:      mapping.GeneratedPosition,
	}
}

func (mapper *documentPositionMapper) ensureMappingsDecoded() bool {
	if mapper.mappingsDecoded {
		return len(mapper.decodedMappings) > 0
	}

	mapper.mappingsDecoded = true
	
	if mapper.sourceMap.Mappings == "" {
		return false
	}

	decoder := DecodeMappings(mapper.sourceMap.Mappings)
	var generatedLineStarts []core.TextPos
	if mapper.generatedFile != nil {
		generatedLineStarts = mapper.generatedFile.LineStarts()
	}

	for mapping, done := decoder.Next(); !done; mapping, done = decoder.Next() {
		if !mapping.IsSourceMapping() {
			continue
		}

		var generatedPos int
		if generatedLineStarts != nil && mapping.GeneratedLine < len(generatedLineStarts) {
			generatedPos = int(generatedLineStarts[mapping.GeneratedLine]) + mapping.GeneratedCharacter
		} else {
			generatedPos = -1
		}

		sourceFile := mapper.host.GetSourceFileLike(mapper.sourceFileAbsolutePaths[mapping.SourceIndex])
		var sourcePos int
		if sourceFile != nil {
			sourceLineStarts := sourceFile.LineStarts()
			if mapping.SourceLine < len(sourceLineStarts) {
				sourcePos = int(sourceLineStarts[mapping.SourceLine]) + mapping.SourceCharacter
			} else {
				sourcePos = -1
			}
		} else {
			sourcePos = -1
		}

		mappedPos := MappedPosition{
			GeneratedPosition: generatedPos,
			SourceIndex:       int(mapping.SourceIndex),
			SourcePosition:    sourcePos,
			SourceFileName:    mapper.sourceFileAbsolutePaths[mapping.SourceIndex],
		}

		mapper.decodedMappings = append(mapper.decodedMappings, mappedPos)
	}

	if err := decoder.Error(); err != nil {
		mapper.host.Log("Error decoding source map mappings: " + err.Error())
		return false
	}

	// Sort and create separate arrays for generated and source mappings
	mapper.generatedMappings = make([]MappedPosition, len(mapper.decodedMappings))
	copy(mapper.generatedMappings, mapper.decodedMappings)
	sort.Slice(mapper.generatedMappings, func(i, j int) bool {
		return mapper.generatedMappings[i].GeneratedPosition < mapper.generatedMappings[j].GeneratedPosition
	})

	// Group by source index
	mapper.sourceMappings = make([][]MappedPosition, len(mapper.sourceFileAbsolutePaths))
	for _, mapping := range mapper.decodedMappings {
		if mapping.SourceIndex >= 0 && mapping.SourceIndex < len(mapper.sourceMappings) {
			mapper.sourceMappings[mapping.SourceIndex] = append(mapper.sourceMappings[mapping.SourceIndex], mapping)
		}
	}

	// Sort each source mapping array
	for i := range mapper.sourceMappings {
		sort.Slice(mapper.sourceMappings[i], func(j, k int) bool {
			return mapper.sourceMappings[i][j].SourcePosition < mapper.sourceMappings[i][k].SourcePosition
		})
	}

	return len(mapper.decodedMappings) > 0
}

// identityMapper is a no-op mapper that returns the input unchanged
type identityMapper struct{}

func (m *identityMapper) GetSourcePosition(input DocumentPosition) DocumentPosition {
	return input
}

func (m *identityMapper) GetGeneratedPosition(input DocumentPosition) DocumentPosition {
	return input
}

// IdentityDocumentPositionMapper returns a mapper that performs no transformation
func IdentityDocumentPositionMapper() DocumentPositionMapper {
	return &identityMapper{}
}
