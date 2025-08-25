package ls

import (
	"sort"
	"strings"

	"github.com/go-json-experiment/json"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/sourcemap"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type DefinitionSourceMapper struct {
	program *compiler.Program
}

func NewDefinitionSourceMapper(program *compiler.Program) *DefinitionSourceMapper {
	return &DefinitionSourceMapper{
		program: program,
	}
}

func (dsm *DefinitionSourceMapper) MapDefinitionLocations(locations []lsproto.Location) []lsproto.Location {
	mappedLocations := make([]lsproto.Location, 0, len(locations))

	for _, location := range locations {
		if mappedLocation := dsm.MapSingleLocation(location); mappedLocation != nil {
			mappedLocations = core.AppendIfUnique(mappedLocations, *mappedLocation)
		} else {
			mappedLocations = core.AppendIfUnique(mappedLocations, location)
		}
	}

	return mappedLocations
}

func (dsm *DefinitionSourceMapper) MapSingleLocation(location lsproto.Location) *lsproto.Location {
	fileName := location.Uri.FileName()

	if !strings.HasSuffix(fileName, ".d.ts") {
		return nil
	}

	if strings.HasPrefix(fileName, "^/bundled/") {
		return nil
	}

	return dsm.tryGetSourcePosition(fileName, location)
}

func (dsm *DefinitionSourceMapper) tryGetSourcePosition(fileName string, location lsproto.Location) *lsproto.Location {
	fs := dsm.program.Host().FS()

	content, ok := fs.ReadFile(fileName)
	if !ok {
		return nil
	}

	mapURL := tryGetSourceMappingURL(content)
	if mapURL == "" {
		return nil
	}

	var mapFileName string
	if strings.HasPrefix(mapURL, "/") {
		mapFileName = mapURL
	} else {
		dir := tspath.GetDirectoryPath(fileName)
		mapFileName = tspath.CombinePaths(dir, mapURL)
	}

	mapContent, ok := fs.ReadFile(mapFileName)
	if !ok {
		return nil
	}

	var sourceMap SourceMap
	if err := json.Unmarshal([]byte(mapContent), &sourceMap); err != nil {
		return nil
	}

	if len(sourceMap.Sources) == 0 {
		return nil
	}

	mapDir := tspath.GetDirectoryPath(mapFileName)
	declLineStarts := core.ComputeLineStarts(content)
	declPos := int(declLineStarts[location.Range.Start.Line]) + int(location.Range.Start.Character)

	sourceRange := dsm.getSourceRangeFromMappings(DocumentPosition{FileName: fileName, Pos: declPos}, sourceMap, mapDir)
	if sourceRange == nil {
		return nil
	}

	sourceFileContent, ok := fs.ReadFile(sourceRange.FileName)
	if !ok {
		return nil
	}

	sourceLineStarts := core.ComputeLineStarts(sourceFileContent)
	sourceStartLine, sourceStartChar := core.PositionToLineAndCharacter(sourceRange.Start, sourceLineStarts)
	sourceEndLine, sourceEndChar := core.PositionToLineAndCharacter(sourceRange.End, sourceLineStarts)

	return &lsproto.Location{
		Uri: FileNameToDocumentURI(sourceRange.FileName),
		Range: lsproto.Range{
			Start: lsproto.Position{Line: uint32(sourceStartLine), Character: uint32(sourceStartChar)},
			End:   lsproto.Position{Line: uint32(sourceEndLine), Character: uint32(sourceEndChar)},
		},
	}
}

type SourceRange struct {
	FileName string
	Start    int
	End      int
}

func (dsm *DefinitionSourceMapper) getSourceRangeFromMappings(loc DocumentPosition, sourceMap SourceMap, mapDir string) *SourceRange {
	if len(sourceMap.Sources) == 0 || sourceMap.Mappings == "" {
		return nil
	}

	declContent, ok := dsm.program.Host().FS().ReadFile(loc.FileName)
	if !ok {
		return nil
	}
	declLineStarts := core.ComputeLineStarts(declContent)

	mappings := decodeMappings(sourceMap.Mappings, declLineStarts)
	if len(mappings) == 0 {
		return nil
	}

	targetIndex := sort.Search(len(mappings), func(i int) bool {
		return mappings[i].GeneratedPosition > loc.Pos
	}) - 1

	if targetIndex < 0 {
		targetIndex = 0
	}
	if targetIndex >= len(mappings) {
		targetIndex = len(mappings) - 1
	}

	currentMapping := mappings[targetIndex]
	if currentMapping.SourceIndex < 0 || currentMapping.SourceIndex >= len(sourceMap.Sources) {
		return nil
	}

	var endMapping *MappedPosition
	if targetIndex+1 < len(mappings) {
		nextMapping := mappings[targetIndex+1]
		if nextMapping.SourceIndex == currentMapping.SourceIndex {
			endMapping = &nextMapping
		}
	}

	var sourcePath string
	source := sourceMap.Sources[currentMapping.SourceIndex]
	if strings.HasPrefix(source, "/") {
		sourcePath = source
	} else {
		sourcePath = tspath.CombinePaths(mapDir, source)
	}

	startLine := currentMapping.SourcePosition / 10000
	startColumn := currentMapping.SourcePosition % 10000

	sourceContent, ok := dsm.program.Host().FS().ReadFile(sourcePath)
	if !ok {
		return &SourceRange{
			FileName: sourcePath,
			Start:    0,
			End:      0,
		}
	}

	sourceLineStarts := core.ComputeLineStarts(sourceContent)
	sourceStartPos := int(sourceLineStarts[startLine]) + startColumn

	var sourceEndPos int
	if endMapping != nil {
		endLine := endMapping.SourcePosition / 10000
		endColumn := endMapping.SourcePosition % 10000
		sourceEndPos = int(sourceLineStarts[endLine]) + endColumn
	} else {
		declNextPos := len(declContent)
		if targetIndex+1 < len(mappings) {
			declNextPos = mappings[targetIndex+1].GeneratedPosition
		}
		rangeLen := declNextPos - currentMapping.GeneratedPosition
		
		sourceEndPos = sourceStartPos + rangeLen
		if sourceEndPos > len(sourceContent) {
			sourceEndPos = len(sourceContent)
		}
		if sourceEndPos <= sourceStartPos {
			sourceEndPos = sourceStartPos + 1
		}
	}

	return &SourceRange{
		FileName: sourcePath,
		Start:    sourceStartPos,
		End:      sourceEndPos,
	}
}

func tryGetSourceMappingURL(content string) string {
	lines := strings.Split(content, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "//# sourceMappingURL=") {
			return strings.TrimPrefix(line, "//# sourceMappingURL=")
		}
		if line != "" && !strings.HasPrefix(line, "//") {
			break
		}
	}
	return ""
}

type SourceMap struct {
	Sources    []string `json:"sources"`
	Mappings   string   `json:"mappings"`
	Names      []string `json:"names"`
	File       string   `json:"file"`
	SourceRoot string   `json:"sourceRoot"`
}

type DocumentPosition struct {
	FileName string
	Pos      int
}

type MappedPosition struct {
	GeneratedPosition int
	SourceIndex       int
	SourcePosition    int
}

func decodeMappings(mappings string, declLineStarts []core.TextPos) []MappedPosition {
	var result []MappedPosition

	decoder := sourcemap.DecodeMappings(mappings)

	for mapping, done := decoder.Next(); !done; mapping, done = decoder.Next() {
		if mapping.IsSourceMapping() {
			generatedPos := int(declLineStarts[mapping.GeneratedLine]) + mapping.GeneratedCharacter

			result = append(result, MappedPosition{
				GeneratedPosition: generatedPos,
				SourceIndex:       int(mapping.SourceIndex),
				SourcePosition:    mapping.SourceLine*10000 + mapping.SourceCharacter,
			})
		}
	}

	if err := decoder.Error(); err != nil {
		return nil
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].GeneratedPosition < result[j].GeneratedPosition
	})

	return result
}
