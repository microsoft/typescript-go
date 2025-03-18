package transformers

import "github.com/microsoft/typescript-go/internal/ast"

type MetaDataProvider interface {
	GetSourceFileMetaData(file *ast.SourceFile) *ast.SourceFileMetaData
}

type metaDataProvider struct {
	getSourceFileMetaData func(file *ast.SourceFile) *ast.SourceFileMetaData
}

func NewMetaDataProvider(getSourceFileMetaData func(file *ast.SourceFile) *ast.SourceFileMetaData) MetaDataProvider {
	return &metaDataProvider{
		getSourceFileMetaData: getSourceFileMetaData,
	}
}

func (r *metaDataProvider) GetSourceFileMetaData(file *ast.SourceFile) *ast.SourceFileMetaData {
	return r.getSourceFileMetaData(file)
}
