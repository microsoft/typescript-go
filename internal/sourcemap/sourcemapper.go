package sourcemap

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tspath"
)

// SourceMapper provides source mapping functionality for LSP operations
type SourceMapper interface {
	TryGetSourcePosition(info DocumentPosition) *DocumentPosition
	TryGetGeneratedPosition(info DocumentPosition) *DocumentPosition
	ClearCache()
}

// SourceMapperHost provides the necessary dependencies for source mapping
type SourceMapperHost interface {
	DocumentPositionMapperHost
	UseCaseSensitiveFileNames() bool
	GetCurrentDirectory() string
	ReadFile(path string) (string, bool)
	FileExists(path string) bool
}

// FileReader is an interface for reading files, potentially with fallback
type FileReader interface {
	ReadFile(path string) (string, bool)
}

// FallbackFileReader extends FileReader with fallback capability
type FallbackFileReader interface {
	FileReader
	ReadFileWithFallback(path string) (string, bool)
}

// sourceMapper implements the SourceMapper interface
type sourceMapper struct {
	host                     SourceMapperHost
	fileReader               FileReader
	getCanonicalFileName     func(string) string
	currentDirectory         string
	documentPositionMappers  map[string]DocumentPositionMapper
}

// CreateSourceMapper creates a new SourceMapper instance
func CreateSourceMapper(host SourceMapperHost, fileReader FileReader) SourceMapper {
	getCanonicalFileName := createGetCanonicalFileName(host.UseCaseSensitiveFileNames())
	
	return &sourceMapper{
		host:                    host,
		fileReader:              fileReader,
		getCanonicalFileName:    getCanonicalFileName,
		currentDirectory:        host.GetCurrentDirectory(),
		documentPositionMappers: make(map[string]DocumentPositionMapper),
	}
}

func (sm *sourceMapper) TryGetSourcePosition(info DocumentPosition) *DocumentPosition {
	if !isDeclarationFileName(info.FileName) {
		return nil
	}

	mapper := sm.getDocumentPositionMapper(info.FileName, "")
	if mapper == nil {
		return nil
	}

	newLoc := mapper.GetSourcePosition(info)
	if newLoc.FileName == info.FileName && newLoc.Pos == info.Pos {
		return nil // No change
	}

	// Recursively try to map further if needed
	if mapped := sm.TryGetSourcePosition(newLoc); mapped != nil {
		return mapped
	}
	
	return &newLoc
}

func (sm *sourceMapper) TryGetGeneratedPosition(info DocumentPosition) *DocumentPosition {
	if isDeclarationFileName(info.FileName) {
		return nil
	}

	// For generated position mapping, we'd need to know the declaration file path
	// This is more complex and typically handled at a higher level
	return nil
}

func (sm *sourceMapper) ClearCache() {
	sm.documentPositionMappers = make(map[string]DocumentPositionMapper)
}

func (sm *sourceMapper) getDocumentPositionMapper(generatedFileName string, sourceFileName string) DocumentPositionMapper {
	path := sm.toPath(generatedFileName)
	if mapper, exists := sm.documentPositionMappers[path]; exists {
		return mapper
	}

	mapper := sm.createDocumentPositionMapper(generatedFileName, sourceFileName)
	sm.documentPositionMappers[path] = mapper
	return mapper
}

func (sm *sourceMapper) createDocumentPositionMapper(generatedFileName string, sourceFileName string) DocumentPositionMapper {
	// First try to read the generated file to get source mapping URL
	content, ok := sm.fileReader.ReadFile(generatedFileName)
	if !ok {
		return IdentityDocumentPositionMapper()
	}

	mapURL := tryGetSourceMappingURL(content)
	if mapURL == "" {
		return IdentityDocumentPositionMapper()
	}

	var mapFileName string
	if strings.HasPrefix(mapURL, "/") {
		mapFileName = mapURL
	} else {
		dir := tspath.GetDirectoryPath(generatedFileName)
		mapFileName = tspath.CombinePaths(dir, mapURL)
	}

	mapContent, ok := sm.fileReader.ReadFile(mapFileName)
	if !ok {
		return IdentityDocumentPositionMapper()
	}

	return CreateDocumentPositionMapper(sm, mapContent, mapFileName)
}

func (sm *sourceMapper) toPath(fileName string) string {
	return string(tspath.ToPath(fileName, sm.currentDirectory, sm.host.UseCaseSensitiveFileNames()))
}

// DocumentPositionMapperHost implementation
func (sm *sourceMapper) GetSourceFileLike(fileName string) SourceFileLike {
	content, ok := sm.fileReader.ReadFile(fileName)
	if !ok {
		return nil
	}
	
	return &sourceFileLike{
		text:       content,
		lineStarts: core.ComputeLineStarts(content),
	}
}

func (sm *sourceMapper) GetCanonicalFileName(path string) string {
	return sm.getCanonicalFileName(path)
}

func (sm *sourceMapper) Log(text string) {
	// Could be implemented to use host's logging if needed
	// For now, we'll keep it simple
}

// sourceFileLike implements SourceFileLike interface
type sourceFileLike struct {
	text       string
	lineStarts []core.TextPos
}

func (sf *sourceFileLike) Text() string {
	return sf.text
}

func (sf *sourceFileLike) LineStarts() []core.TextPos {
	return sf.lineStarts
}

// Helper functions

func isDeclarationFileName(fileName string) bool {
	return strings.HasSuffix(fileName, ".d.ts")
}

func createGetCanonicalFileName(useCaseSensitiveFileNames bool) func(string) string {
	if useCaseSensitiveFileNames {
		return func(path string) string { return path }
	}
	return strings.ToLower
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
