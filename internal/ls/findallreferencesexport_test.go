package ls

import (
	"github.com/microsoft/typescript-go/internal/astnav"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil/lstestutil"
)

func (l *LanguageService) GetExpectedReferenceFromMarker(marker *lstestutil.Marker) *lsproto.Location {
	// Temporary testing function--this function only works for markers that are on symbols/names.
	// We won't need this once marker ranges are implemented, or once reference tests are baselined
	_, sourceFile := l.getProgramAndFile(marker.Filename)
	node := astnav.GetTouchingPropertyName(sourceFile, marker.Position)
	return &lsproto.Location{
		Uri:   FileNameToDocumentURI(marker.Filename),
		Range: *l.createLspRangeFromNode(node, sourceFile),
	}
}
