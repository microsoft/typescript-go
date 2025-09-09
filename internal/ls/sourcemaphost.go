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
func MapSingleDefinitionLocation(program *compiler.Program, location lsproto.Location) *lsproto.Location {
	fileName := location.Uri.FileName()

	if !strings.HasSuffix(fileName, ".d.ts") {
		return nil
	}

	if strings.HasPrefix(fileName, "^/bundled/") {
		return nil
	}

	host := &sourcemapHost{program: program}
	sourceMapper := CreateSourceMapperForProgram(program)
	
	return tryMapLocation(sourceMapper, host, location)
}

func tryMapLocation(sourceMapper sourcemap.SourceMapper, host *sourcemapHost, location lsproto.Location) *lsproto.Location {
	fileName := location.Uri.FileName()
	
	declContent, ok := host.readFileWithFallback(fileName)
	if !ok {
		return nil
	}
	
	declLineStarts := core.ComputeLineStarts(declContent)
	
	if int(location.Range.Start.Line) >= len(declLineStarts) || int(location.Range.End.Line) >= len(declLineStarts) {
		return nil
	}
	
	declStartPos := int(declLineStarts[location.Range.Start.Line]) + int(location.Range.Start.Character)
	declEndPos := int(declLineStarts[location.Range.End.Line]) + int(location.Range.End.Character)
	
	startInput := sourcemap.DocumentPosition{
		FileName: fileName,
		Pos:      declStartPos,
	}

	startResult := sourceMapper.TryGetSourcePosition(startInput)
	if startResult == nil {
		return nil
	}

	sourceContent, ok := host.readFileWithFallback(startResult.FileName)
	if !ok {
		return nil
	}
	
	sourceLineStarts := core.ComputeLineStarts(sourceContent)
	
	sourceStartLine, sourceStartChar := core.PositionToLineAndCharacter(startResult.Pos, sourceLineStarts)
	
	originalRangeLength := declEndPos - declStartPos
	
	if declStartPos >= 0 && declEndPos <= len(declContent) {
		originalText := strings.TrimSpace(declContent[declStartPos:declEndPos])
		
		if isSimpleIdentifier(originalText) {
			sourceEndPos := startResult.Pos + len(originalText)
			
			if sourceEndPos > len(sourceContent) {
				sourceEndPos = len(sourceContent)
			}
			
			sourceEndLine, sourceEndChar := core.PositionToLineAndCharacter(sourceEndPos, sourceLineStarts)
			
			return &lsproto.Location{
				Uri: FileNameToDocumentURI(startResult.FileName),
				Range: lsproto.Range{
					Start: lsproto.Position{Line: uint32(sourceStartLine), Character: uint32(sourceStartChar)},
					End:   lsproto.Position{Line: uint32(sourceEndLine), Character: uint32(sourceEndChar)},
				},
			}
		}
	}
	
	sourceEndPos := startResult.Pos + originalRangeLength
	
	if sourceEndPos > len(sourceContent) {
		sourceEndPos = len(sourceContent)
	}
	
	sourceEndLine, sourceEndChar := core.PositionToLineAndCharacter(sourceEndPos, sourceLineStarts)

	return &lsproto.Location{
		Uri: FileNameToDocumentURI(startResult.FileName),
		Range: lsproto.Range{
			Start: lsproto.Position{Line: uint32(sourceStartLine), Character: uint32(sourceStartChar)},
			End:   lsproto.Position{Line: uint32(sourceEndLine), Character: uint32(sourceEndChar)},
		},
	}
}


func isIdentifierChar(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || 
		   (ch >= '0' && ch <= '9') || ch == '_' || ch == '$'
}

func isSimpleIdentifier(text string) bool {
	if len(text) == 0 {
		return false
	}
	
	first := text[0]
	if !((first >= 'a' && first <= 'z') || (first >= 'A' && first <= 'Z') || first == '_' || first == '$') {
		return false
	}
	
	for i := 1; i < len(text); i++ {
		if !isIdentifierChar(text[i]) {
			return false
		}
	}
	
	return true
}