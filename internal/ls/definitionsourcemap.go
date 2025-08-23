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

// DefinitionSourceMapper handles mapping declaration file definitions to their source files
// This follows TypeScript's exact SourceMapper approach using the existing internal/sourcemap VLQ decoder
type DefinitionSourceMapper struct {
	program *compiler.Program
}

// NewDefinitionSourceMapper creates a new definition source mapper
func NewDefinitionSourceMapper(program *compiler.Program) *DefinitionSourceMapper {
	return &DefinitionSourceMapper{
		program: program,
	}
}

// MapDefinitionLocations maps a list of definition locations, converting declaration file
// locations to their corresponding source file locations when possible.
// This follows TypeScript's session.mapDefinitionInfoLocations pattern
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

// MapSingleLocation maps a single location from declaration file to source file
// This follows TypeScript's exact approach using DocumentPositionMapper and source maps
func (dsm *DefinitionSourceMapper) MapSingleLocation(location lsproto.Location) *lsproto.Location {
	fileName := location.Uri.FileName()

	// Only process declaration files (matching TypeScript's isDeclarationFileName check)
	if !strings.HasSuffix(fileName, ".d.ts") {
		return nil
	}

	// TypeScript's approach: Use DocumentPositionMapper with source maps
	return dsm.tryGetSourcePosition(fileName, location)
}

// tryGetSourcePosition implements TypeScript's exact approach from sourcemaps.ts
// This matches the tryGetSourcePosition function in TypeScript's SourceMapper
// Uses VFS for all file access (matching TypeScript's host.readFile() pattern)
func (dsm *DefinitionSourceMapper) tryGetSourcePosition(fileName string, location lsproto.Location) *lsproto.Location {
	fs := dsm.program.Host().FS()

	// Step 1: Read the declaration file content
	content, ok := fs.ReadFile(fileName)
	if !ok {
		return nil
	}

	// Step 2: Look for sourceMappingURL comment
	mapURL := dsm.tryGetSourceMappingURL(content)
	if mapURL == "" {
		return nil
	}

	// Step 3: Resolve map file path
	var mapFileName string
	if strings.HasPrefix(mapURL, "/") {
		mapFileName = mapURL
	} else {
		dir := tspath.GetDirectoryPath(fileName)
		mapFileName = tspath.CombinePaths(dir, mapURL)
	}

	// Step 4: Read map file (following TypeScript's host.readFile pattern)
	mapContent, ok := fs.ReadFile(mapFileName)
	if !ok {
		return nil
	}

	// Step 5: Parse full source map and create position mapper (following TypeScript's exact approach)
	var fullSourceMap SourceMap
	if err := json.Unmarshal([]byte(mapContent), &fullSourceMap); err != nil {
		return nil
	}

	if len(fullSourceMap.Sources) == 0 {
		return nil
	}

	// Step 6: Find existing source file (following TypeScript's pattern)
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

	// Step 7: Convert LSP position to absolute position (following TypeScript's DocumentPosition pattern)
	declFileContent, ok := fs.ReadFile(fileName)
	if !ok {
		return nil
	}

	declLineStarts := computeLineStarts(declFileContent)
	declPos := computePositionOfLineAndCharacter(declLineStarts, int(location.Range.Start.Line), int(location.Range.Start.Character))

	// Step 8: Map position using source map (following TypeScript's getSourcePosition approach)
	sourcePos := dsm.getSourcePosition(DocumentPosition{FileName: fileName, Pos: declPos}, fullSourceMap, mapDir)
	if sourcePos == nil {
		// If source mapping fails, return file-level mapping
		return &lsproto.Location{
			Uri: FileNameToDocumentURI(targetSourcePath),
			Range: lsproto.Range{
				Start: lsproto.Position{Line: 0, Character: 0},
				End:   lsproto.Position{Line: 0, Character: 0},
			},
		}
	}

	// Step 9: Convert mapped absolute position back to line/character (following TypeScript's computeLineAndCharacterOfPosition)
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

// tryGetSourceMappingURL looks for //# sourceMappingURL= comments in file content
// This matches TypeScript's tryGetSourceMappingURL function
func (dsm *DefinitionSourceMapper) tryGetSourceMappingURL(content string) string {
	lines := strings.Split(content, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "//# sourceMappingURL=") {
			return strings.TrimPrefix(line, "//# sourceMappingURL=")
		}
		// Stop at first non-empty, non-comment line
		if line != "" && !strings.HasPrefix(line, "//") {
			break
		}
	}
	return ""
}

// SourceMap represents the structure of a source map JSON file (matches TypeScript's RawSourceMap)
type SourceMap struct {
	Sources    []string `json:"sources"`
	Mappings   string   `json:"mappings"`
	Names      []string `json:"names"`
	File       string   `json:"file"`
	SourceRoot string   `json:"sourceRoot"`
}

// DocumentPosition matches TypeScript's DocumentPosition interface
type DocumentPosition struct {
	FileName string
	Pos      int // absolute character position in file
}

// MappedPosition represents a position mapping in the source map (matches TypeScript's interface)
type MappedPosition struct {
	GeneratedPosition int // absolute position in generated file
	SourceIndex       int
	SourcePosition    int // absolute position in source file
}

// extractSourcesFromMap parses source map JSON and extracts source file paths
func extractSourcesFromMap(mapContent string) []string {
	var sourceMap SourceMap
	if err := json.Unmarshal([]byte(mapContent), &sourceMap); err != nil {
		return nil
	}
	return sourceMap.Sources
}

// LineAndCharacter matches TypeScript's LineAndCharacter interface
type LineAndCharacter struct {
	Line      int
	Character int
}

// computeLineStarts computes line start positions (matches TypeScript's computeLineStarts)
func computeLineStarts(text string) []int {
	lineStarts := []int{0} // First line starts at position 0

	for i, char := range text {
		if char == '\n' {
			lineStarts = append(lineStarts, i+1) // Line starts after newline
		}
	}

	return lineStarts
}

// computePositionOfLineAndCharacter converts line/character to absolute position (matches TypeScript's function)
func computePositionOfLineAndCharacter(lineStarts []int, line int, character int) int {
	if line < 0 || line >= len(lineStarts) {
		return 0 // Clamp to valid range
	}
	return lineStarts[line] + character
}

// computeLineAndCharacterOfPosition converts absolute position to line/character (matches TypeScript's function)
func computeLineAndCharacterOfPosition(lineStarts []int, position int) LineAndCharacter {
	lineNumber := computeLineOfPosition(lineStarts, position)
	return LineAndCharacter{
		Line:      lineNumber,
		Character: position - lineStarts[lineNumber],
	}
}

// computeLineOfPosition finds line number for a position using binary search (matches TypeScript's function)
func computeLineOfPosition(lineStarts []int, position int) int {
	lineNumber := sort.Search(len(lineStarts), func(i int) bool {
		return lineStarts[i] > position
	}) - 1

	if lineNumber < 0 {
		lineNumber = 0
	}

	return lineNumber
}

// decodeMappings decodes source map VLQ mappings using the existing sourcemap decoder (matches TypeScript's decodeMappings logic)
func decodeMappings(mappings string, declLineStarts []int) []MappedPosition {
	var result []MappedPosition

	// Use the existing sourcemap decoder from the codebase
	decoder := sourcemap.DecodeMappings(mappings)

	for mapping, done := decoder.Next(); !done; mapping, done = decoder.Next() {
		if mapping.IsSourceMapping() {
			// Convert line/column to absolute positions (following TypeScript's approach)
			generatedPos := computePositionOfLineAndCharacter(declLineStarts, mapping.GeneratedLine, mapping.GeneratedCharacter)

			result = append(result, MappedPosition{
				GeneratedPosition: generatedPos,
				SourceIndex:       int(mapping.SourceIndex),                           // Convert from sourcemap.SourceIndex to int
				SourcePosition:    mapping.SourceLine*10000 + mapping.SourceCharacter, // Store line/column encoded
			})
		}
	}

	// Check for decoder errors
	if err := decoder.Error(); err != nil {
		return nil
	}

	// Sort by generated position (matching TypeScript's approach)
	sort.Slice(result, func(i, j int) bool {
		return result[i].GeneratedPosition < result[j].GeneratedPosition
	})

	return result
}

// getSourcePosition maps position using source map (matches TypeScript's DocumentPositionMapper.getSourcePosition)
func (dsm *DefinitionSourceMapper) getSourcePosition(loc DocumentPosition, sourceMap SourceMap, mapDir string) *DocumentPosition {
	if len(sourceMap.Sources) == 0 || sourceMap.Mappings == "" {
		return nil
	}

	// Get line starts for the declaration file
	declContent, ok := dsm.program.Host().FS().ReadFile(loc.FileName)
	if !ok {
		return nil
	}
	declLineStarts := computeLineStarts(declContent)

	// Decode all mappings using existing sourcemap decoder (following TypeScript's approach)
	mappings := decodeMappings(sourceMap.Mappings, declLineStarts)
	if len(mappings) == 0 {
		return nil
	}

	// Binary search for the closest mapping (matching TypeScript's binarySearchKey logic)
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

	// Resolve source file path
	var sourcePath string
	source := sourceMap.Sources[mapping.SourceIndex]
	if strings.HasPrefix(source, "/") {
		sourcePath = source
	} else {
		sourcePath = tspath.CombinePaths(mapDir, source)
	}

	// Decode line/column from encoded position (we stored it as line*10000 + column)
	sourceLine := mapping.SourcePosition / 10000
	sourceColumn := mapping.SourcePosition % 10000

	// Read source file to get line starts for accurate position calculation
	sourceContent, ok := dsm.program.Host().FS().ReadFile(sourcePath)
	if !ok {
		// Fallback to line 0 if we can't read the source file
		return &DocumentPosition{
			FileName: sourcePath,
			Pos:      0,
		}
	}

	// Convert line/column to absolute position (following TypeScript's exact approach)
	sourceLineStarts := computeLineStarts(sourceContent)
	sourcePos := computePositionOfLineAndCharacter(sourceLineStarts, sourceLine, sourceColumn)

	return &DocumentPosition{
		FileName: sourcePath,
		Pos:      sourcePos,
	}
}
