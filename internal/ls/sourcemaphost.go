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


func (h *sourcemapHost) GetSourceFileLike(fileName string) sourcemap.SourceFileLike {
	if content, ok := h.readFileWithFallback(fileName); ok {
		return &sourcemapSourceFile{
			text:       content,
			lineStarts: core.ComputeLineStarts(content),
		}
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


type sourcemapSourceFile struct {
	text       string
	lineStarts []core.TextPos
}

func (f *sourcemapSourceFile) Text() string {
	return f.text
}

func (f *sourcemapSourceFile) LineStarts() []core.TextPos {
	return f.lineStarts
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
	var declStartPos, declEndPos int
	
	if declFile != nil {
		declStartPos = int(languageService.converters.LineAndCharacterToPosition(declFile, location.Range.Start))
		declEndPos = int(languageService.converters.LineAndCharacterToPosition(declFile, location.Range.End))
	} else {
		declContent, ok := host.readFileWithFallback(fileName)
		if !ok {
			return nil
		}
		
		declScript := &textScript{
			fileName: fileName,
			text:     declContent,
		}
		
		declStartPos = int(languageService.converters.LineAndCharacterToPosition(declScript, location.Range.Start))
		declEndPos = int(languageService.converters.LineAndCharacterToPosition(declScript, location.Range.End))
	}
	
	startInput := sourcemap.DocumentPosition{
		FileName: fileName,
		Pos:      declStartPos,
	}

	startResult := sourceMapper.TryGetSourcePosition(startInput)
	if startResult == nil {
		return nil
	}

	originalRangeLength := declEndPos - declStartPos
	sourceEndPos := startResult.Pos + originalRangeLength
	
	sourceFile := program.GetSourceFile(startResult.FileName)
	var sourceStartLSP lsproto.Position
	var sourceEndLSP lsproto.Position
	
	if sourceFile != nil {
		if sourceEndPos > len(sourceFile.Text()) {
			sourceEndPos = len(sourceFile.Text())
		}
		
		sourceStartLSP = languageService.converters.PositionToLineAndCharacter(sourceFile, core.TextPos(startResult.Pos))
		sourceEndLSP = languageService.converters.PositionToLineAndCharacter(sourceFile, core.TextPos(sourceEndPos))

	} else {
		sourceContent, ok := host.readFileWithFallback(startResult.FileName)
		if !ok {
			return nil
		}
		
		if sourceEndPos > len(sourceContent) {
			sourceEndPos = len(sourceContent)
		}
		
		textScript := &textScript{
			fileName: startResult.FileName,
			text:     sourceContent,
		}
		
		sourceStartLSP = languageService.converters.PositionToLineAndCharacter(textScript, core.TextPos(startResult.Pos))
		sourceEndLSP = languageService.converters.PositionToLineAndCharacter(textScript, core.TextPos(sourceEndPos))
	}
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