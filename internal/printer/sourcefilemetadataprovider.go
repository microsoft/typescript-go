package printer

import (
	"github.com/microsoft/typescript-go/internal/ast"
)

type SourceFileMetaDataProvider interface {
	GetSourceFileMetaData(sourceFile *ast.SourceFile) *ast.SourceFileMetaData
}
