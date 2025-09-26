package ls

import (
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/debug"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/sourcemap"
	"github.com/microsoft/typescript-go/internal/tspath"
)

func (l *LanguageService) getMappedLocation(location *lsproto.Location) *lsproto.Location {
	uriStart, start := l.tryGetSourceLSPPosition(location.Uri.FileName(), &location.Range.Start)
	if uriStart == nil {
		return location
	}
	uriEnd, end := l.tryGetSourceLSPPosition(location.Uri.FileName(), &location.Range.End)
	debug.Assert(uriEnd == uriStart, "start and end should be in same file")
	debug.Assert(end != nil, "end position should be valid")
	return &lsproto.Location{
		Uri:   *uriStart,
		Range: lsproto.Range{Start: *start, End: *end},
	}
}

func (l *LanguageService) getMappedPosition() {
	// !!! HERE
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

func (l *LanguageService) tryGetSourceLSPPosition(
	genFileName string,
	position *lsproto.Position,
) (*lsproto.DocumentUri, *lsproto.Position) {
	genText, ok := l.ReadFile(genFileName)
	if !ok {
		return nil, nil // That shouldn't happen
	}
	genPos := l.converters.LineAndCharacterToPosition(&script{fileName: genFileName, text: genText}, *position)
	documentPos := l.tryGetSourcePosition(genFileName, genPos)
	if documentPos == nil {
		return nil, nil
	}
	documentURI := FileNameToDocumentURI(documentPos.FileName)
	sourceText, ok := l.ReadFile(documentPos.FileName)
	if !ok {
		return nil, nil
	}
	sourcePos := l.converters.PositionToLineAndCharacter(
		&script{fileName: documentPos.FileName, text: sourceText},
		core.TextPos(documentPos.Pos),
	)
	return &documentURI, &sourcePos
}

func (l *LanguageService) tryGetSourcePosition(
	fileName string,
	genPosition core.TextPos,
) *sourcemap.DocumentPosition {
	if !tspath.IsDeclarationFileName(fileName) {
		return nil
	}

	positionMapper := l.GetDocumentPositionMapper(fileName)
	documentPos := positionMapper.GetSourcePosition(&sourcemap.DocumentPosition{FileName: fileName, Pos: int(genPosition)})
	if documentPos == nil {
		return nil
	}
	if newPos := l.tryGetSourcePosition(documentPos.FileName, core.TextPos(documentPos.Pos)); newPos != nil {
		return newPos
	}
	return documentPos
}
