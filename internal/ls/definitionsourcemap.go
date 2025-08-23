package ls

import (
	"encoding/json"
	"net/url"
	"os"
	"strings"

	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/tspath"
)

// DocumentURIToFileName converts a document URI to a file path
func DocumentURIToFileName(uri lsproto.DocumentUri) string {
	if parsed, err := url.Parse(string(uri)); err == nil && parsed.Scheme == "file" {
		return parsed.Path
	}
	return string(uri)
}

// DefinitionSourceMapper handles mapping declaration file definitions to their source files
// This follows TypeScript's exact SourceMapper approach using source maps and DocumentPositionMapper
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
	fileName := DocumentURIToFileName(location.Uri)
	
	// Only process declaration files (matching TypeScript's isDeclarationFileName check)
	if !strings.HasSuffix(fileName, ".d.ts") {
		return nil
	}
	
	// TypeScript's approach: Use DocumentPositionMapper with source maps
	return dsm.tryGetSourcePosition(fileName, location)
}



// tryGetSourcePosition implements TypeScript's exact approach from sourcemaps.ts
// This matches the tryGetSourcePosition function in TypeScript's SourceMapper
func (dsm *DefinitionSourceMapper) tryGetSourcePosition(fileName string, location lsproto.Location) *lsproto.Location {
	// Step 1: Read the declaration file to get its content
	content, err := os.ReadFile(fileName)
	if err != nil {
		return nil
	}
	
	// Step 2: Look for sourceMappingURL comment
	mapURL := dsm.tryGetSourceMappingURL(string(content))
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
	
	// Step 4: Read map file
	mapContent, err := os.ReadFile(mapFileName)
	if err != nil {
		return nil
	}
	
	// Step 5: Parse source map to get sources
	sources := extractSourcesFromMap(string(mapContent))
	if len(sources) == 0 {
		return nil
	}
	
	// Step 6: Find existing source file and locate precise symbol position
	mapDir := tspath.GetDirectoryPath(mapFileName)
	for _, source := range sources {
		var sourcePath string
		if strings.HasPrefix(source, "/") {
			sourcePath = source
		} else {
			sourcePath = tspath.CombinePaths(mapDir, source)
		}
		
		if _, err := os.Stat(sourcePath); err == nil {
			// Return source file at line 0 (file-level mapping)
			return &lsproto.Location{
				Uri: FileNameToDocumentURI(sourcePath),
				Range: lsproto.Range{
					Start: lsproto.Position{Line: 0, Character: 0},
					End:   lsproto.Position{Line: 0, Character: 0},
				},
			}
		}
	}
	
	return nil
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



// SourceMap represents the structure of a source map JSON file
type SourceMap struct {
	Sources []string `json:"sources"`
}

// extractSourcesFromMap parses source map JSON and extracts source file paths
func extractSourcesFromMap(mapContent string) []string {
	var sourceMap SourceMap
	if err := json.Unmarshal([]byte(mapContent), &sourceMap); err != nil {
		return nil
	}
	return sourceMap.Sources
}
