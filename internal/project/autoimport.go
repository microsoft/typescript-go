package project

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/ls/autoimport"
	"github.com/microsoft/typescript-go/internal/packagejson"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type autoImportRegistryCloneHost struct {
	projectCollection *ProjectCollection
	snapshotFSBuilder *snapshotFSBuilder
	parseCache        *ParseCache
}

// FS implements autoimport.RegistryCloneHost.
func (a *autoImportRegistryCloneHost) FS() vfs.FS {
	return a.snapshotFSBuilder
}

// GetCurrentDirectory implements autoimport.RegistryCloneHost.
func (a *autoImportRegistryCloneHost) GetCurrentDirectory() string {
	panic("unimplemented")
}

// GetDefaultProject implements autoimport.RegistryCloneHost.
func (a *autoImportRegistryCloneHost) GetDefaultProject(fileName string) (tspath.Path, *compiler.Program) {
	panic("unimplemented")
}

// GetPackageJson implements autoimport.RegistryCloneHost.
func (a *autoImportRegistryCloneHost) GetPackageJson(fileName string) *packagejson.InfoCacheEntry {
	panic("unimplemented")
}

// GetProgramForProject implements autoimport.RegistryCloneHost.
func (a *autoImportRegistryCloneHost) GetProgramForProject(projectPath tspath.Path) *compiler.Program {
	panic("unimplemented")
}

// GetSourceFile implements autoimport.RegistryCloneHost.
func (a *autoImportRegistryCloneHost) GetSourceFile(fileName string, path tspath.Path) *ast.SourceFile {
	panic("unimplemented")
}

var _ autoimport.RegistryCloneHost = (*autoImportRegistryCloneHost)(nil)
