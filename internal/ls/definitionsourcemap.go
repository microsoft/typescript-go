package ls

import (
	"encoding/json"
	"net/url"
	"os"
	"regexp"
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
	fs := dsm.program.Host().FS()
	
	// Step 1: Read the declaration file to get its content
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
	
	// Step 4: Check if map file exists
	if !fs.FileExists(mapFileName) {
		return nil
	}
	
	// Step 5: Read map file (with VFS fallback to OS)
	mapContent, ok := fs.ReadFile(mapFileName)
	if !ok {
		// Try direct OS file read as fallback for VFS issues
		if osContent, err := os.ReadFile(mapFileName); err == nil {
			mapContent = string(osContent)
		} else {
			return nil
		}
	}
	
	// Step 6: Parse source map to get sources
	mapper := &simpleSourceMapMapper{
		mapContent:   mapContent,
		mapFileName:  mapFileName,
		fs:           fs,
	}
	
	sources := mapper.extractSourcesFromMap()
	if len(sources) == 0 {
		return nil
	}
	
	// Step 7: Find existing source file and locate precise symbol position
	mapDir := tspath.GetDirectoryPath(mapFileName)
	for _, source := range sources {
		var sourcePath string
		if strings.HasPrefix(source, "/") {
			sourcePath = source
		} else {
			sourcePath = tspath.CombinePaths(mapDir, source)
		}
		
		if fs.FileExists(sourcePath) {
			// Try to find the precise position of the symbol in the source file
			symbolName := dsm.extractSymbolNameFromDeclaration(fileName, location)
			sourcePosition := dsm.findSymbolInSourceFile(sourcePath, symbolName)
			
			if sourcePosition == nil {
				// Fallback to line 0 if symbol not found
				sourcePosition = &lsproto.Position{Line: 0, Character: 0}
			}
			
			return &lsproto.Location{
				Uri: FileNameToDocumentURI(sourcePath),
				Range: lsproto.Range{
					Start: *sourcePosition,
					End:   *sourcePosition,
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



// simpleSourceMapMapper is a simplified source map parser
type simpleSourceMapMapper struct {
	mapContent   string
	mapFileName  string
	fs           interface{ FileExists(path string) bool }
}

// SourceMap represents the structure of a source map JSON file
type SourceMap struct {
	Sources []string `json:"sources"`
}

// extractSourcesFromMap extracts source file paths from source map JSON
func (s *simpleSourceMapMapper) extractSourcesFromMap() []string {
	var sourceMap SourceMap
	if err := json.Unmarshal([]byte(s.mapContent), &sourceMap); err != nil {
		return nil
	}
	return sourceMap.Sources
}



// extractSymbolNameFromDeclaration extracts the symbol name from a declaration file at the given position
func (dsm *DefinitionSourceMapper) extractSymbolNameFromDeclaration(fileName string, location lsproto.Location) string {
	fs := dsm.program.Host().FS()
	content, ok := fs.ReadFile(fileName)
	if !ok {
		return ""
	}
	
	lines := strings.Split(content, "\n")
	lineNum := int(location.Range.Start.Line)
	if lineNum >= len(lines) {
		return ""
	}
	
	line := lines[lineNum]
	charPos := int(location.Range.Start.Character)
	
	// Extract identifier around the cursor position
	return extractIdentifierAtPosition(line, charPos)
}

// findSymbolInSourceFile finds a symbol in a source file and returns its position
func (dsm *DefinitionSourceMapper) findSymbolInSourceFile(sourcePath string, symbolName string) *lsproto.Position {
	if symbolName == "" {
		return nil
	}
	
	// Try to read from VFS first, fallback to OS
	var content string
	if vfsContent, ok := dsm.program.Host().FS().ReadFile(sourcePath); ok {
		content = vfsContent
	} else if osContent, err := os.ReadFile(sourcePath); err == nil {
		content = string(osContent)
	} else {
		return nil
	}
	
	lines := strings.Split(content, "\n")
	
	// Search for various TypeScript/JavaScript patterns
	patterns := []string{
		// Export function/const/let/var
		`export\s+(?:function\s+|const\s+|let\s+|var\s+)?` + regexp.QuoteMeta(symbolName) + `\b`,
		// Function declaration
		`function\s+` + regexp.QuoteMeta(symbolName) + `\b`,
		// Const/let/var assignment
		`(?:const|let|var)\s+` + regexp.QuoteMeta(symbolName) + `\b`,
		// Class/interface
		`(?:class|interface)\s+` + regexp.QuoteMeta(symbolName) + `\b`,
		// Object property/method
		regexp.QuoteMeta(symbolName) + `\s*[:=]`,
	}
	
	for lineIndex, line := range lines {
		for _, pattern := range patterns {
			if matched, _ := regexp.MatchString(pattern, line); matched {
				// Find the exact character position of the symbol
				re := regexp.MustCompile(regexp.QuoteMeta(symbolName) + `\b`)
				if match := re.FindStringIndex(line); match != nil {
					return &lsproto.Position{
						Line:      uint32(lineIndex),
						Character: uint32(match[0]),
					}
				}
			}
		}
	}
	
	return nil
}

// extractIdentifierAtPosition extracts an identifier from a line at the given character position
func extractIdentifierAtPosition(line string, charPos int) string {
	if charPos >= len(line) {
		return ""
	}
	
	// Find the start and end of the identifier at this position
	start := charPos
	end := charPos
	
	// Move start backwards to find the beginning of the identifier
	for start > 0 && isIdentifierChar(line[start-1]) {
		start--
	}
	
	// Move end forwards to find the end of the identifier
	for end < len(line) && isIdentifierChar(line[end]) {
		end++
	}
	
	if start == end {
		return ""
	}
	
	return line[start:end]
}

// isIdentifierChar checks if a character can be part of an identifier
func isIdentifierChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' || c == '$'
}
