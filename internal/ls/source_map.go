package ls

import (
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/outputpaths"
	"github.com/microsoft/typescript-go/internal/sourcemap"
	"github.com/microsoft/typescript-go/internal/spanmap"
	"github.com/microsoft/typescript-go/internal/tspath"
)

func (l *LanguageService) getMappedLocation(fileName string, fileRange core.TextRange) (lsproto.Location, spanmap.Fidelity) {
	startPos := l.tryGetSourcePosition(fileName, core.TextPos(fileRange.Pos()))
	if startPos == nil {
		lspRange, fidelity := l.createLspRangeFromRange(fileRange, l.getScript(fileName))
		return lsproto.Location{
			Uri:   lsconv.FileNameToDocumentURI(fileName),
			Range: lspRange,
		}, fidelity
	}
	endPos := l.tryGetSourcePosition(fileName, core.TextPos(fileRange.End()))
	if endPos == nil || endPos.FileName != startPos.FileName || endPos.Pos < startPos.Pos {
		// When end doesn't map, maps to a different source file (e.g. in a .d.ts with a
		// multi-source source map from --outFile compilation), or maps to a position before
		// start (non-monotonic source map mappings), approximate the end position.
		endPos = &sourcemap.DocumentPosition{
			FileName: startPos.FileName,
			Pos:      startPos.Pos + fileRange.Len(),
		}
	}
	newRange := core.NewTextRange(startPos.Pos, endPos.Pos)
	lspRange, fidelity := l.createLspRangeFromRange(newRange, l.getScript(startPos.FileName))
	return lsproto.Location{
		Uri:   lsconv.FileNameToDocumentURI(startPos.FileName),
		Range: lspRange,
	}, fidelity
}

type script struct {
	fileName string
	text     string
}

func (s *script) FileName() string {
	return s.fileName
}

func (s *script) Text() string {
	return s.text
}

// SpanMap and OriginalText satisfy lsconv.Script for a plain (non-content-mapped) file: it carries no
// span map and its original text is its own text.
func (s *script) SpanMap() *spanmap.SpanMap { return nil }

func (s *script) OriginalText() string { return s.text }

func (l *LanguageService) getScript(fileName string) lsconv.Script {
	if program := l.GetProgram(); program != nil {
		if file := program.GetSourceFile(fileName); file != nil {
			// Use a program lookup first so we get mappable scripts for content-mapped files
			return file
		}
	}
	// Fall back to getting the plain text from the file system. This happens when fileName
	// is the result of a .d.ts.map mapping back to a source file whose declaration file is
	// part of the program instead of the source.
	text, ok := l.host.ReadFile(fileName)
	if !ok {
		return nil
	}
	return &script{fileName: fileName, text: text}
}

func (l *LanguageService) tryGetSourcePosition(
	fileName string,
	position core.TextPos,
) *sourcemap.DocumentPosition {
	newPos := l.tryGetSourcePositionWorker(fileName, position)
	if newPos != nil {
		if _, ok := l.ReadFile(newPos.FileName); !ok { // File doesn't exist
			return nil
		}
	}
	return newPos
}

func (l *LanguageService) tryGetSourcePositionWorker(
	fileName string,
	position core.TextPos,
) *sourcemap.DocumentPosition {
	if !tspath.IsDeclarationFileName(fileName) {
		return nil
	}

	positionMapper := l.GetDocumentPositionMapper(fileName)
	documentPos := positionMapper.GetSourcePosition(&sourcemap.DocumentPosition{FileName: fileName, Pos: int(position)})
	if documentPos == nil {
		return nil
	}
	if newPos := l.tryGetSourcePositionWorker(documentPos.FileName, core.TextPos(documentPos.Pos)); newPos != nil {
		return newPos
	}
	return documentPos
}

func (l *LanguageService) tryGetGeneratedPosition(
	fileName string,
	position core.TextPos,
) *sourcemap.DocumentPosition {
	newPos := l.tryGetGeneratedPositionWorker(fileName, position)
	if newPos != nil {
		if _, ok := l.ReadFile(newPos.FileName); !ok { // File doesn't exist
			return nil
		}
	}
	return newPos
}

func (l *LanguageService) tryGetGeneratedPositionWorker(
	fileName string,
	position core.TextPos,
) *sourcemap.DocumentPosition {
	if tspath.IsDeclarationFileName(fileName) {
		return nil
	}

	program := l.GetProgram()
	if program == nil || program.GetSourceFile(fileName) == nil {
		return nil
	}

	path := l.toPath(fileName)
	// If this is source file of project reference source (instead of redirect) there is no generated position
	if program.IsSourceFromProjectReference(path) {
		return nil
	}

	declarationFileName := outputpaths.GetOutputDeclarationFileNameWorker(fileName, program.Options(), program)
	positionMapper := l.GetDocumentPositionMapper(declarationFileName)
	documentPos := positionMapper.GetGeneratedPosition(&sourcemap.DocumentPosition{FileName: fileName, Pos: int(position)})
	if documentPos == nil {
		return nil
	}
	if newPos := l.tryGetGeneratedPositionWorker(documentPos.FileName, core.TextPos(documentPos.Pos)); newPos != nil {
		return newPos
	}
	return documentPos
}
