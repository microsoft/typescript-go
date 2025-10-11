package ls

import "github.com/microsoft/typescript-go/internal/sourcemap"

type Host interface {
	ReadFile(path string) (contents string, ok bool)
	Converters() *Converters
	GetDocumentPositionMapper(fileName string) *sourcemap.DocumentPositionMapper
}
