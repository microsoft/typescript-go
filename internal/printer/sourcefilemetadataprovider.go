package printer

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type SourceFileMetaDataProvider interface {
	GetSourceFileMetaData(path tspath.Path) *ast.SourceFileMetaData
}

type sourceFileMetadataProvider struct {
	getSourceFileMetaData func(path tspath.Path) *ast.SourceFileMetaData
}

func (r *sourceFileMetadataProvider) GetSourceFileMetaData(path tspath.Path) *ast.SourceFileMetaData {
	return r.getSourceFileMetaData(path)
}
