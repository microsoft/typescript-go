package project

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls/autoimport"
	"github.com/microsoft/typescript-go/internal/packagejson"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
)

type autoImportRegistryCloneHost struct {
	projectCollection *ProjectCollection
	parseCache        *ParseCache
	fs                *sourceFS
	currentDirectory  string
}

var _ autoimport.RegistryCloneHost = (*autoImportRegistryCloneHost)(nil)

func newAutoImportRegistryCloneHost(
	projectCollection *ProjectCollection,
	parseCache *ParseCache,
	snapshotFSBuilder *snapshotFSBuilder,
	currentDirectory string,
	toPath func(fileName string) tspath.Path,
) *autoImportRegistryCloneHost {
	return &autoImportRegistryCloneHost{
		projectCollection: projectCollection,
		parseCache:        parseCache,
		fs:                &sourceFS{toPath: toPath, source: snapshotFSBuilder},
	}
}

// FS implements autoimport.RegistryCloneHost.
func (a *autoImportRegistryCloneHost) FS() vfs.FS {
	return a.fs
}

// GetCurrentDirectory implements autoimport.RegistryCloneHost.
func (a *autoImportRegistryCloneHost) GetCurrentDirectory() string {
	return a.currentDirectory
}

// GetDefaultProject implements autoimport.RegistryCloneHost.
func (a *autoImportRegistryCloneHost) GetDefaultProject(path tspath.Path) (tspath.Path, *compiler.Program) {
	project := a.projectCollection.GetDefaultProject(path)
	if project == nil {
		return "", nil
	}
	return project.configFilePath, project.GetProgram()
}

// GetPackageJson implements autoimport.RegistryCloneHost.
func (a *autoImportRegistryCloneHost) GetPackageJson(fileName string) *packagejson.InfoCacheEntry {
	// !!! ref-counted shared cache
	fh := a.fs.GetFile(fileName)
	packageDirectory := tspath.GetDirectoryPath(fileName)
	if fh == nil {
		return &packagejson.InfoCacheEntry{
			DirectoryExists:  a.fs.DirectoryExists(packageDirectory),
			PackageDirectory: packageDirectory,
		}
	}
	fields, err := packagejson.Parse([]byte(fh.Content()))
	if err != nil {
		return &packagejson.InfoCacheEntry{
			DirectoryExists:  true,
			PackageDirectory: tspath.GetDirectoryPath(fileName),
			Contents: &packagejson.PackageJson{
				Parseable: false,
			},
		}
	}
	return &packagejson.InfoCacheEntry{
		DirectoryExists:  true,
		PackageDirectory: tspath.GetDirectoryPath(fileName),
		Contents: &packagejson.PackageJson{
			Fields:    fields,
			Parseable: true,
		},
	}
}

// GetProgramForProject implements autoimport.RegistryCloneHost.
func (a *autoImportRegistryCloneHost) GetProgramForProject(projectPath tspath.Path) *compiler.Program {
	project := a.projectCollection.GetProjectByPath(projectPath)
	if project == nil {
		return nil
	}
	return project.GetProgram()
}

// GetSourceFile implements autoimport.RegistryCloneHost.
func (a *autoImportRegistryCloneHost) GetSourceFile(fileName string, path tspath.Path) *ast.SourceFile {
	fh := a.fs.GetFile(fileName)
	if fh == nil {
		return nil
	}
	// !!! andrewbranch/autoimport: this should usually/always be a peek instead of an acquire
	return a.parseCache.Acquire(NewParseCacheKey(ast.SourceFileParseOptions{
		FileName:         fileName,
		Path:             path,
		CompilerOptions:  core.EmptyCompilerOptions.SourceFileAffecting(),
		JSDocParsingMode: ast.JSDocParsingModeParseAll,
		// !!! wrong if we load non-.d.ts files here
		ExternalModuleIndicatorOptions: ast.ExternalModuleIndicatorOptions{},
	}, fh.Hash(), core.GetScriptKindFromFileName(fileName)), fh)
}
