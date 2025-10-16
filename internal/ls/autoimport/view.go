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
	projectTrie, ok := v.registry.projects[v.projectKey]
	if ok {
		results = append(results, projectTrie.Search(prefix)...)
	}
	for directoryPath, nodeModulesTrie := range v.registry.nodeModules {
		if directoryPath.GetDirectoryPath().ContainsPath(v.importingFile.Path()) {
			results = append(results, nodeModulesTrie.Search(prefix)...)
		}
	}
	return results
}
