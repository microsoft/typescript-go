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

	mapURL := dsm.tryGetSourceMappingURL(content)
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

	var fullSourceMap SourceMap
	if err := json.Unmarshal([]byte(mapContent), &fullSourceMap); err != nil {
		return nil
	}

	if len(fullSourceMap.Sources) == 0 {
		return nil
	}

	mapDir := tspath.GetDirectoryPath(mapFileName)
	var targetSourcePath string
	for _, source := range fullSourceMap.Sources {
		var sourcePath string
		if strings.HasPrefix(source, "/") {
			sourcePath = source
		} else {
			sourcePath = tspath.CombinePaths(mapDir, source)
		}

		if fs.FileExists(sourcePath) {
			targetSourcePath = sourcePath
			break
		}
	}

	if targetSourcePath == "" {
		return nil
	}

	declFileContent, ok := fs.ReadFile(fileName)
	if !ok {
		return nil
	}

	declLineStarts := computeLineStarts(declFileContent)
	declPos := computePositionOfLineAndCharacter(declLineStarts, int(location.Range.Start.Line), int(location.Range.Start.Character))

	sourcePos := dsm.getSourcePosition(DocumentPosition{FileName: fileName, Pos: declPos}, fullSourceMap, mapDir)
	if sourcePos == nil {
		return &lsproto.Location{
			Uri: FileNameToDocumentURI(targetSourcePath),
			Range: lsproto.Range{
				Start: lsproto.Position{Line: 0, Character: 0},
				End:   lsproto.Position{Line: 0, Character: 0},
			},
		}
	}

	sourceFileContent, ok := fs.ReadFile(sourcePos.FileName)
	if !ok {
		return nil
	}

	sourceLineStarts := computeLineStarts(sourceFileContent)
	sourceLineChar := computeLineAndCharacterOfPosition(sourceLineStarts, sourcePos.Pos)

	return &lsproto.Location{
		Uri: FileNameToDocumentURI(sourcePos.FileName),
		Range: lsproto.Range{
			Start: lsproto.Position{Line: uint32(sourceLineChar.Line), Character: uint32(sourceLineChar.Character)},
			End:   lsproto.Position{Line: uint32(sourceLineChar.Line), Character: uint32(sourceLineChar.Character)},
		},
	}
}

func (dsm *DefinitionSourceMapper) tryGetSourceMappingURL(content string) string {
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

func extractSourcesFromMap(mapContent string) []string {
	var sourceMap SourceMap
	if err := json.Unmarshal([]byte(mapContent), &sourceMap); err != nil {
		return nil
	}
	return sourceMap.Sources
}

type LineAndCharacter struct {
	Line      int
	Character int
}

func computeLineStarts(text string) []int {
	lineStarts := []int{0}

	for i, char := range text {
		if char == '\n' {
			lineStarts = append(lineStarts, i+1)
		}
	}

	return lineStarts
}

func computePositionOfLineAndCharacter(lineStarts []int, line int, character int) int {
	if line < 0 || line >= len(lineStarts) {
		return 0
	}
	return lineStarts[line] + character
}

func computeLineAndCharacterOfPosition(lineStarts []int, position int) LineAndCharacter {
	lineNumber := computeLineOfPosition(lineStarts, position)
	return LineAndCharacter{
		Line:      lineNumber,
		Character: position - lineStarts[lineNumber],
	}
}

func computeLineOfPosition(lineStarts []int, position int) int {
	lineNumber := sort.Search(len(lineStarts), func(i int) bool {
		return lineStarts[i] > position
	}) - 1

	if lineNumber < 0 {
		lineNumber = 0
	}

	return lineNumber
}

func decodeMappings(mappings string, declLineStarts []int) []MappedPosition {
	var result []MappedPosition

	decoder := sourcemap.DecodeMappings(mappings)

	for mapping, done := decoder.Next(); !done; mapping, done = decoder.Next() {
		if mapping.IsSourceMapping() {
			generatedPos := computePositionOfLineAndCharacter(declLineStarts, mapping.GeneratedLine, mapping.GeneratedCharacter)

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

func (dsm *DefinitionSourceMapper) getSourcePosition(loc DocumentPosition, sourceMap SourceMap, mapDir string) *DocumentPosition {
	if len(sourceMap.Sources) == 0 || sourceMap.Mappings == "" {
		return nil
	}

	declContent, ok := dsm.program.Host().FS().ReadFile(loc.FileName)
	if !ok {
		return nil
	}
	declLineStarts := computeLineStarts(declContent)

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

	mapping := mappings[targetIndex]
	if mapping.SourceIndex < 0 || mapping.SourceIndex >= len(sourceMap.Sources) {
		return nil
	}

	var sourcePath string
	source := sourceMap.Sources[mapping.SourceIndex]
	if strings.HasPrefix(source, "/") {
		sourcePath = source
	} else {
		sourcePath = tspath.CombinePaths(mapDir, source)
	}

	sourceLine := mapping.SourcePosition / 10000
	sourceColumn := mapping.SourcePosition % 10000

	sourceContent, ok := dsm.program.Host().FS().ReadFile(sourcePath)
	if !ok {
		return &DocumentPosition{
			FileName: sourcePath,
			Pos:      0,
		}
	}

	sourceLineStarts := computeLineStarts(sourceContent)
	sourcePos := computePositionOfLineAndCharacter(sourceLineStarts, sourceLine, sourceColumn)

	return &DocumentPosition{
		FileName: sourcePath,
		Pos:      sourcePos,
	}
}
