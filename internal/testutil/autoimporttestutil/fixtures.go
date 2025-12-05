package autoimporttestutil

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/ls/lsconv"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/testutil/projecttestutil"
	"github.com/microsoft/typescript-go/internal/tspath"
)

// FileHandle represents a file created for an autoimport lifecycle test.
type FileHandle struct {
	fileName string
	content  string
}

func (f FileHandle) FileName() string         { return f.fileName }
func (f FileHandle) Content() string          { return f.content }
func (f FileHandle) URI() lsproto.DocumentUri { return lsconv.FileNameToDocumentURI(f.fileName) }

// ProjectFileHandle adds export metadata for TypeScript source files.
type ProjectFileHandle struct {
	FileHandle
	exportIdentifier string
}

func (f ProjectFileHandle) ExportIdentifier() string { return f.exportIdentifier }

// NodeModulesPackageHandle describes a generated package under node_modules.
type NodeModulesPackageHandle struct {
	Name             string
	Directory        string
	ExportIdentifier string
	packageJSON      FileHandle
	declaration      FileHandle
}

func (p NodeModulesPackageHandle) PackageJSONFile() FileHandle { return p.packageJSON }
func (p NodeModulesPackageHandle) DeclarationFile() FileHandle { return p.declaration }

// MonorepoHandle exposes the generated monorepo layout including root and packages.
type MonorepoHandle struct {
	root            string
	rootNodeModules []NodeModulesPackageHandle
	packages        []ProjectHandle
	rootTSConfig    FileHandle
	rootPackageJSON FileHandle
}

func (m MonorepoHandle) Root() string { return m.root }
func (m MonorepoHandle) RootNodeModules() []NodeModulesPackageHandle {
	return slices.Clone(m.rootNodeModules)
}
func (m MonorepoHandle) Packages() []ProjectHandle { return slices.Clone(m.packages) }
func (m MonorepoHandle) Package(index int) ProjectHandle {
	if index < 0 || index >= len(m.packages) {
		panic(fmt.Sprintf("package index %d out of range", index))
	}
	return m.packages[index]
}
func (m MonorepoHandle) RootTSConfig() FileHandle        { return m.rootTSConfig }
func (m MonorepoHandle) RootPackageJSONFile() FileHandle { return m.rootPackageJSON }

// ProjectHandle exposes the generated project layout for a fixture project root.
type ProjectHandle struct {
	root        string
	files       []ProjectFileHandle
	tsconfig    FileHandle
	packageJSON FileHandle
	nodeModules []NodeModulesPackageHandle
}

func (p ProjectHandle) Root() string               { return p.root }
func (p ProjectHandle) Files() []ProjectFileHandle { return slices.Clone(p.files) }
func (p ProjectHandle) File(index int) ProjectFileHandle {
	if index < 0 || index >= len(p.files) {
		panic(fmt.Sprintf("file index %d out of range", index))
	}
	return p.files[index]
}
func (p ProjectHandle) TSConfig() FileHandle        { return p.tsconfig }
func (p ProjectHandle) PackageJSONFile() FileHandle { return p.packageJSON }
func (p ProjectHandle) NodeModules() []NodeModulesPackageHandle {
	return slices.Clone(p.nodeModules)
}

func (p ProjectHandle) NodeModuleByName(name string) *NodeModulesPackageHandle {
	for i := range p.nodeModules {
		if p.nodeModules[i].Name == name {
			return &p.nodeModules[i]
		}
	}
	return nil
}

// Fixture encapsulates a fully-initialized auto import lifecycle test session.
type Fixture struct {
	session  *project.Session
	utils    *projecttestutil.SessionUtils
	projects []ProjectHandle
}

func (f *Fixture) Session() *project.Session            { return f.session }
func (f *Fixture) Utils() *projecttestutil.SessionUtils { return f.utils }
func (f *Fixture) Projects() []ProjectHandle            { return slices.Clone(f.projects) }
func (f *Fixture) Project(index int) ProjectHandle {
	if index < 0 || index >= len(f.projects) {
		panic(fmt.Sprintf("project index %d out of range", index))
	}
	return f.projects[index]
}
func (f *Fixture) SingleProject() ProjectHandle { return f.Project(0) }

// MonorepoFixture encapsulates a fully-initialized monorepo lifecycle test session.
type MonorepoFixture struct {
	session  *project.Session
	utils    *projecttestutil.SessionUtils
	monorepo MonorepoHandle
}

func (f *MonorepoFixture) Session() *project.Session            { return f.session }
func (f *MonorepoFixture) Utils() *projecttestutil.SessionUtils { return f.utils }
func (f *MonorepoFixture) Monorepo() MonorepoHandle             { return f.monorepo }

// MonorepoConfig describes a monorepo package.
type MonorepoConfig struct {
	Name                   string // e.g., "package-a" becomes directory name under packages/
	FileCount              int    // Number of TypeScript source files
	NodeModulePackageCount int    // Number of packages in this package's node_modules
}

// SetupMonorepoLifecycleSession builds a monorepo workspace with root-level node_modules
// and multiple packages, each potentially with their own node_modules.
// The structure is:
//
//	monorepoRoot/
//	├── tsconfig.json (base config)
//	├── package.json
//	├── node_modules/
//	│   └── <rootNodeModuleCount packages>
//	└── packages/
//	    ├── package-a/
//	    │   ├── tsconfig.json
//	    │   ├── package.json
//	    │   ├── node_modules/
//	    │   │   └── <package-specific packages>
//	    │   └── *.ts files
//	    └── package-b/
//	        └── ...
func SetupMonorepoLifecycleSession(t *testing.T, monorepoRoot string, rootNodeModuleCount int, packages []MonorepoConfig) *MonorepoFixture {
	t.Helper()
	builder := newFileMapBuilder(nil)

	// Add root tsconfig.json
	rootTSConfigPath := tspath.CombinePaths(monorepoRoot, "tsconfig.json")
	rootTSConfigContent := "{\n  \"compilerOptions\": {\n    \"module\": \"esnext\",\n    \"target\": \"esnext\",\n    \"strict\": true,\n    \"baseUrl\": \".\"\n  }\n}\n"
	builder.AddTextFile(rootTSConfigPath, rootTSConfigContent)
	rootTSConfig := FileHandle{fileName: rootTSConfigPath, content: rootTSConfigContent}

	// Add root node_modules
	rootNodeModulesDir := tspath.CombinePaths(monorepoRoot, "node_modules")
	rootNodeModules := builder.AddNodeModulesPackages(rootNodeModulesDir, rootNodeModuleCount)

	// Add root package.json with dependencies
	rootPackageJSON := builder.addRootPackageJSON(monorepoRoot, rootNodeModules)

	// Build each package in packages/
	packagesDir := tspath.CombinePaths(monorepoRoot, "packages")
	packageHandles := make([]ProjectHandle, 0, len(packages))
	for _, pkg := range packages {
		pkgDir := tspath.CombinePaths(packagesDir, pkg.Name)
		builder.AddLocalProject(pkgDir, pkg.FileCount)

		// Add package-specific node_modules if requested
		var pkgNodeModules []NodeModulesPackageHandle
		if pkg.NodeModulePackageCount > 0 {
			pkgNodeModulesDir := tspath.CombinePaths(pkgDir, "node_modules")
			pkgNodeModules = builder.AddNodeModulesPackages(pkgNodeModulesDir, pkg.NodeModulePackageCount)
		}

		// Combine root and package-level dependencies for package.json
		allDeps := append(slices.Clone(rootNodeModules), pkgNodeModules...)
		builder.AddPackageJSONWithDependencies(pkgDir, allDeps)
	}

	// Build project handles after all packages are created
	for _, pkg := range packages {
		pkgDir := tspath.CombinePaths(packagesDir, pkg.Name)
		if record, ok := builder.projects[pkgDir]; ok {
			packageHandles = append(packageHandles, record.toHandles())
		}
	}

	session, sessionUtils := projecttestutil.Setup(builder.Files())
	t.Cleanup(session.Close)

	// Build root node_modules handle by looking at the project record for the root
	// (created as side effect of AddNodeModulesPackages)
	var rootNodeModulesHandles []NodeModulesPackageHandle
	if rootRecord, ok := builder.projects[monorepoRoot]; ok {
		rootNodeModulesHandles = rootRecord.nodeModules
	}

	return &MonorepoFixture{
		session: session,
		utils:   sessionUtils,
		monorepo: MonorepoHandle{
			root:            monorepoRoot,
			rootNodeModules: rootNodeModulesHandles,
			packages:        packageHandles,
			rootTSConfig:    rootTSConfig,
			rootPackageJSON: rootPackageJSON,
		},
	}
}

// SetupLifecycleSession builds a basic single-project workspace configured with the
// requested number of TypeScript files and a single synthetic node_modules package.
func SetupLifecycleSession(t *testing.T, projectRoot string, fileCount int) *Fixture {
	t.Helper()
	builder := newFileMapBuilder(nil)
	builder.AddLocalProject(projectRoot, fileCount)
	nodeModulesDir := tspath.CombinePaths(projectRoot, "node_modules")
	deps := builder.AddNodeModulesPackages(nodeModulesDir, 1)
	builder.AddPackageJSONWithDependencies(projectRoot, deps)
	session, sessionUtils := projecttestutil.Setup(builder.Files())
	t.Cleanup(session.Close)
	return &Fixture{
		session:  session,
		utils:    sessionUtils,
		projects: builder.projectHandles(),
	}
}

type fileMapBuilder struct {
	files         map[string]any
	nextPackageID int
	nextProjectID int
	projects      map[string]*projectRecord
}

type projectRecord struct {
	root        string
	sourceFiles []projectFile
	tsconfig    FileHandle
	packageJSON *FileHandle
	nodeModules []NodeModulesPackageHandle
}

type projectFile struct {
	FileName         string
	ExportIdentifier string
	Content          string
}

func newFileMapBuilder(initial map[string]any) *fileMapBuilder {
	b := &fileMapBuilder{
		files:    make(map[string]any),
		projects: make(map[string]*projectRecord),
	}
	if len(initial) == 0 {
		return b
	}
	for path, content := range initial {
		b.files[normalizeAbsolutePath(path)] = content
	}
	return b
}

func (b *fileMapBuilder) ensureProjectRecord(root string) *projectRecord {
	if record, ok := b.projects[root]; ok {
		return record
	}
	record := &projectRecord{root: root}
	b.projects[root] = record
	return record
}

func (b *fileMapBuilder) projectHandles() []ProjectHandle {
	keys := slices.Collect(maps.Keys(b.projects))
	slices.Sort(keys)
	result := make([]ProjectHandle, 0, len(keys))
	for _, key := range keys {
		result = append(result, b.projects[key].toHandles())
	}
	return result
}

func (r *projectRecord) toHandles() ProjectHandle {
	files := make([]ProjectFileHandle, len(r.sourceFiles))
	for i, file := range r.sourceFiles {
		files[i] = ProjectFileHandle{
			FileHandle:       FileHandle{fileName: file.FileName, content: file.Content},
			exportIdentifier: file.ExportIdentifier,
		}
	}
	packageJSON := FileHandle{}
	if r.packageJSON != nil {
		packageJSON = *r.packageJSON
	}
	return ProjectHandle{
		root:        r.root,
		files:       files,
		tsconfig:    r.tsconfig,
		packageJSON: packageJSON,
		nodeModules: slices.Clone(r.nodeModules),
	}
}

func (b *fileMapBuilder) Files() map[string]any {
	return maps.Clone(b.files)
}

func (b *fileMapBuilder) AddTextFile(path string, contents string) {
	b.ensureFiles()
	b.files[normalizeAbsolutePath(path)] = contents
}

func (b *fileMapBuilder) AddNodeModulesPackages(nodeModulesDir string, count int) []NodeModulesPackageHandle {
	packages := make([]NodeModulesPackageHandle, 0, count)
	for i := 0; i < count; i++ {
		packages = append(packages, b.AddNodeModulesPackage(nodeModulesDir))
	}
	return packages
}

func (b *fileMapBuilder) AddNodeModulesPackage(nodeModulesDir string) NodeModulesPackageHandle {
	b.ensureFiles()
	normalizedDir := normalizeAbsolutePath(nodeModulesDir)
	if tspath.GetBaseFileName(normalizedDir) != "node_modules" {
		panic("nodeModulesDir must point to a node_modules directory: " + nodeModulesDir)
	}
	b.nextPackageID++
	name := fmt.Sprintf("pkg%d", b.nextPackageID)
	exportName := fmt.Sprintf("value%d", b.nextPackageID)
	pkgDir := tspath.CombinePaths(normalizedDir, name)
	packageJSONPath := tspath.CombinePaths(pkgDir, "package.json")
	packageJSONContent := fmt.Sprintf(`{"name":"%s","types":"index.d.ts"}`, name)
	b.files[packageJSONPath] = packageJSONContent
	declarationPath := tspath.CombinePaths(pkgDir, "index.d.ts")
	declarationContent := fmt.Sprintf("export declare const %s: number;\n", exportName)
	b.files[declarationPath] = declarationContent
	packageHandle := NodeModulesPackageHandle{
		Name:             name,
		Directory:        pkgDir,
		ExportIdentifier: exportName,
		packageJSON:      FileHandle{fileName: packageJSONPath, content: packageJSONContent},
		declaration:      FileHandle{fileName: declarationPath, content: declarationContent},
	}
	projectRoot := tspath.GetDirectoryPath(normalizedDir)
	record := b.ensureProjectRecord(projectRoot)
	record.nodeModules = append(record.nodeModules, packageHandle)
	return packageHandle
}

func (b *fileMapBuilder) AddLocalProject(projectDir string, fileCount int) {
	b.ensureFiles()
	if fileCount <= 0 {
		panic("fileCount must be positive")
	}
	dir := normalizeAbsolutePath(projectDir)
	record := b.ensureProjectRecord(dir)
	b.nextProjectID++
	tsConfigPath := tspath.CombinePaths(dir, "tsconfig.json")
	tsConfigContent := "{\n  \"compilerOptions\": {\n    \"module\": \"esnext\",\n    \"target\": \"esnext\",\n    \"strict\": true\n  }\n}\n"
	b.files[tsConfigPath] = tsConfigContent
	record.tsconfig = FileHandle{fileName: tsConfigPath, content: tsConfigContent}
	for i := 1; i <= fileCount; i++ {
		path := tspath.CombinePaths(dir, fmt.Sprintf("file%d.ts", i))
		exportName := fmt.Sprintf("localExport%d_%d", b.nextProjectID, i)
		content := fmt.Sprintf("export const %s = %d;\n", exportName, i)
		b.files[path] = content
		record.sourceFiles = append(record.sourceFiles, projectFile{FileName: path, ExportIdentifier: exportName, Content: content})
	}
}

func (b *fileMapBuilder) AddPackageJSONWithDependencies(projectDir string, deps []NodeModulesPackageHandle) FileHandle {
	b.ensureFiles()
	dir := normalizeAbsolutePath(projectDir)
	packageJSONPath := tspath.CombinePaths(dir, "package.json")
	b.nextProjectID++
	dependencyLines := make([]string, 0, len(deps))
	for _, dep := range deps {
		dependencyLines = append(dependencyLines, fmt.Sprintf("\"%s\": \"*\"", dep.Name))
	}
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("{\n  \"name\": \"local-project-%d\"", b.nextProjectID))
	if len(dependencyLines) > 0 {
		builder.WriteString(",\n  \"dependencies\": {\n    ")
		builder.WriteString(strings.Join(dependencyLines, ",\n    "))
		builder.WriteString("\n  }\n")
	} else {
		builder.WriteString("\n")
	}
	builder.WriteString("}\n")
	content := builder.String()
	b.files[packageJSONPath] = content
	record := b.ensureProjectRecord(dir)
	packageHandle := FileHandle{fileName: packageJSONPath, content: content}
	record.packageJSON = &packageHandle
	return packageHandle
}

// addRootPackageJSON creates a root package.json for a monorepo without creating a project record.
// This is used to set up the root workspace config without treating it as a project.
func (b *fileMapBuilder) addRootPackageJSON(rootDir string, deps []NodeModulesPackageHandle) FileHandle {
	b.ensureFiles()
	dir := normalizeAbsolutePath(rootDir)
	packageJSONPath := tspath.CombinePaths(dir, "package.json")
	dependencyLines := make([]string, 0, len(deps))
	for _, dep := range deps {
		dependencyLines = append(dependencyLines, fmt.Sprintf("\"%s\": \"*\"", dep.Name))
	}
	var builder strings.Builder
	builder.WriteString("{\n  \"name\": \"monorepo-root\",\n  \"private\": true")
	if len(dependencyLines) > 0 {
		builder.WriteString(",\n  \"dependencies\": {\n    ")
		builder.WriteString(strings.Join(dependencyLines, ",\n    "))
		builder.WriteString("\n  }\n")
	} else {
		builder.WriteString("\n")
	}
	builder.WriteString("}\n")
	content := builder.String()
	b.files[packageJSONPath] = content
	return FileHandle{fileName: packageJSONPath, content: content}
}

func (b *fileMapBuilder) ensureFiles() {
	if b.files == nil {
		b.files = make(map[string]any)
	}
}

func normalizeAbsolutePath(path string) string {
	normalized := tspath.NormalizePath(path)
	if !tspath.PathIsAbsolute(normalized) {
		panic("paths used in lifecycle tests must be absolute: " + path)
	}
	return normalized
}
