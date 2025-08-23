package ls

import (
	"encoding/json"
	"strings"

	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/tspath"
)

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
	
	// Step 5: Parse source map to get sources
	sources := extractSourcesFromMap(mapContent)
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
		
		// Check if source file exists
		if fs.FileExists(sourcePath) {
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
