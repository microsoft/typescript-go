package ls

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/sourcemap"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type sourcemapHost struct {
	program *compiler.Program
}

func (h *sourcemapHost) GetSource(fileName string) sourcemap.Source {
	if content, ok := h.readFileWithFallback(fileName); ok {
		return sourcemap.NewSimpleSourceFile(fileName, content)
	}
	return nil
}

func (h *sourcemapHost) GetCanonicalFileName(path string) string {
	return tspath.GetCanonicalFileName(path, h.program.UseCaseSensitiveFileNames())
}

func (h *sourcemapHost) Log(text string) {
}

func (h *sourcemapHost) UseCaseSensitiveFileNames() bool {
	return h.program.UseCaseSensitiveFileNames()
}

func (h *sourcemapHost) GetCurrentDirectory() string {
	return h.program.Host().GetCurrentDirectory()
}

func (h *sourcemapHost) ReadFile(path string) (string, bool) {
	return h.readFileWithFallback(path)
}

func (h *sourcemapHost) FileExists(path string) bool {
	return h.program.Host().FS().FileExists(path)
}

// Tries to read a file from the program's FS first, and falls back to underlying FS for files not tracked by the program.
func (h *sourcemapHost) readFileWithFallback(fileName string) (string, bool) {
	if content, ok := h.program.Host().FS().ReadFile(fileName); ok {
		return content, true
	}

	if fallbackFS, ok := h.program.Host().FS().(sourcemap.FallbackFileReader); ok {
		return fallbackFS.ReadFileWithFallback(fileName)
	}

	return "", false
}

type sourcemapFileReader struct {
	host *sourcemapHost
}

func (r *sourcemapFileReader) ReadFile(path string) (string, bool) {
	return r.host.readFileWithFallback(path)
}

// Creates a SourceMapper for the given program.
func CreateSourceMapperForProgram(program *compiler.Program) sourcemap.SourceMapper {
	host := &sourcemapHost{program: program}
	fileReader := &sourcemapFileReader{host: host}
	return sourcemap.CreateSourceMapper(host, fileReader)
}

// Maps a single definition location using source maps.
func MapSingleDefinitionLocation(program *compiler.Program, location lsproto.Location, languageService *LanguageService) *lsproto.Location {
	fileName := location.Uri.FileName()

	if !strings.HasSuffix(fileName, ".d.ts") {
		return nil
	}

	host := &sourcemapHost{program: program}
	sourceMapper := CreateSourceMapperForProgram(program)

	return tryMapLocation(sourceMapper, host, location, languageService)
}

func tryMapLocation(sourceMapper sourcemap.SourceMapper, host *sourcemapHost, location lsproto.Location, languageService *LanguageService) *lsproto.Location {
	fileName := location.Uri.FileName()
	program := host.program

	declFile := program.GetSourceFile(fileName)
	var declStartPos, declEndPos core.TextPos

	// Get the script interface for the declaration file
	var declScript Script
	if declFile != nil {
		declScript = declFile
	} else {
		declContent, ok := host.readFileWithFallback(fileName)
		if !ok {
			return nil
		}

		declScript = &textScript{
			fileName: fileName,
			text:     declContent,
		}
	}

	// Convert both positions using the same script interface
	declStartPos = languageService.converters.LineAndCharacterToPosition(declScript, location.Range.Start)
	declEndPos = languageService.converters.LineAndCharacterToPosition(declScript, location.Range.End)

	startInput := sourcemap.DocumentPosition{
		FileName: fileName,
		Pos:      declStartPos,
	}

	startResult := sourceMapper.TryGetSourcePosition(startInput)
	if startResult == nil {
		return nil
	}

	// Map the end position individually through the source map
	endInput := sourcemap.DocumentPosition{
		FileName: fileName,
		Pos:      declEndPos,
	}

	endResult := sourceMapper.TryGetSourcePosition(endInput)
	var sourceEndPos core.TextPos
	if endResult != nil && endResult.FileName == startResult.FileName {
		// Both positions mapped to the same source file
		sourceEndPos = endResult.Pos
	} else {
		// Fallback: use original range length (this shouldn't happen often)
		originalRangeLength := declEndPos - declStartPos
		sourceEndPos = startResult.Pos + originalRangeLength
	}

	// Get the script interface for the source file
	sourceFile := program.GetSourceFile(startResult.FileName)
	var sourceScript Script
	if sourceFile != nil {
		sourceScript = sourceFile
	} else {
		sourceContent, ok := host.readFileWithFallback(startResult.FileName)
		if !ok {
			return nil
		}

		sourceScript = &textScript{
			fileName: startResult.FileName,
			text:     sourceContent,
		}
	}

	// Convert both positions using the same script interface
	sourceStartLSP := languageService.converters.PositionToLineAndCharacter(sourceScript, startResult.Pos)
	sourceEndLSP := languageService.converters.PositionToLineAndCharacter(sourceScript, sourceEndPos)
	return &lsproto.Location{
		Uri: FileNameToDocumentURI(startResult.FileName),
		Range: lsproto.Range{
			Start: sourceStartLSP,
			End:   sourceEndLSP,
		},
	}
}

// textScript is a simple wrapper that implements the Script interface for raw text content
type textScript struct {
	fileName string
	text     string
}

func (t *textScript) FileName() string {
	return t.fileName
}

func (t *textScript) Text() string {
	return t.text
}
