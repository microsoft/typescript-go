package autoimport

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type View struct {
	registry      *Registry
	importingFile *ast.SourceFile
	projectKey    tspath.Path
}

func NewView(registry *Registry, importingFile *ast.SourceFile, projectKey tspath.Path) *View {
	return &View{
		registry:      registry,
		importingFile: importingFile,
		projectKey:    projectKey,
	}
}

func (v *View) Search(prefix string) []*RawExport {
	// !!! deal with duplicates due to symlinks
	var results []*RawExport
	bucket, ok := v.registry.projects[v.projectKey]
	if ok {
		results = append(results, bucket.Index.Search(prefix)...)
	}
	for directoryPath, nodeModulesBucket := range v.registry.nodeModules {
		// !!! better to iterate by ancestor directory?
		if directoryPath.GetDirectoryPath().ContainsPath(v.importingFile.Path()) {
			results = append(results, nodeModulesBucket.Index.Search(prefix)...)
		}
	}
	return results
}
