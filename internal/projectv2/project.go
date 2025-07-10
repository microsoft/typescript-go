package projectv2

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/project"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type Kind int

const (
	KindInferred Kind = iota
	KindConfigured
)

type PendingReload int

const (
	PendingReloadNone PendingReload = iota
	PendingReloadFileNames
	PendingReloadFull
)

var _ ls.Host = (*Project)(nil)

// Project represents a TypeScript project.
// If changing struct fields, also update the Clone method.
type Project struct {
	Name             string
	Kind             Kind
	currentDirectory string
	configFileName   string
	configFilePath   tspath.Path

	dirty         bool
	dirtyFilePath tspath.Path

	host            *compilerHost
	CommandLine     *tsoptions.ParsedCommandLine
	Program         *compiler.Program
	LanguageService *ls.LanguageService

	checkerPool *project.CheckerPool
}

func NewConfiguredProject(
	configFileName string,
	configFilePath tspath.Path,
	builder *projectCollectionBuilder,
) *Project {
	p := NewProject(configFileName, KindConfigured, tspath.GetDirectoryPath(configFileName), builder)
	p.configFileName = configFileName
	p.configFilePath = configFilePath
	return p
}

func NewInferredProject(
	currentDirectory string,
	compilerOptions *core.CompilerOptions,
	rootFileNames []string,
	builder *projectCollectionBuilder,
) *Project {
	p := NewProject("/dev/null/inferredProject", KindInferred, currentDirectory, builder)
	if compilerOptions == nil {
		compilerOptions = &core.CompilerOptions{
			AllowJs:                    core.TSTrue,
			Module:                     core.ModuleKindESNext,
			ModuleResolution:           core.ModuleResolutionKindBundler,
			Target:                     core.ScriptTargetES2022,
			Jsx:                        core.JsxEmitReactJSX,
			AllowImportingTsExtensions: core.TSTrue,
			StrictNullChecks:           core.TSTrue,
			StrictFunctionTypes:        core.TSTrue,
			SourceMap:                  core.TSTrue,
			ESModuleInterop:            core.TSTrue,
			AllowNonTsExtensions:       core.TSTrue,
			ResolveJsonModule:          core.TSTrue,
		}
	}
	p.CommandLine = tsoptions.NewParsedCommandLine(
		compilerOptions,
		rootFileNames,
		tspath.ComparePathsOptions{
			UseCaseSensitiveFileNames: builder.fs.fs.UseCaseSensitiveFileNames(),
			CurrentDirectory:          currentDirectory,
		},
	)
	return p
}

func NewProject(
	name string,
	kind Kind,
	currentDirectory string,
	builder *projectCollectionBuilder,
) *Project {
	project := &Project{
		Name:             name,
		Kind:             kind,
		currentDirectory: currentDirectory,
		dirty:            true,
	}
	host := newCompilerHost(
		currentDirectory,
		project,
		builder,
	)
	project.host = host
	return project
}

// GetLineMap implements ls.Host.
func (p *Project) GetLineMap(fileName string) *ls.LineMap {
	return p.host.overlayFS.getFile(fileName).LineMap()
}

// GetPositionEncoding implements ls.Host.
func (p *Project) GetPositionEncoding() lsproto.PositionEncodingKind {
	return p.host.sessionOptions.PositionEncoding
}

// GetProgram implements ls.Host.
func (p *Project) GetProgram() *compiler.Program {
	return p.Program
}

func (p *Project) containsFile(path tspath.Path) bool {
	if p.isRoot(path) {
		return true
	}
	return p.Program != nil && p.Program.GetSourceFileByPath(path) != nil
}

func (p *Project) isRoot(path tspath.Path) bool {
	if p.CommandLine == nil {
		return false
	}
	_, ok := p.CommandLine.FileNamesByPath()[path]
	return ok
}

func (p *Project) IsSourceFromProjectReference(path tspath.Path) bool {
	return p.Program != nil && p.Program.IsSourceFromProjectReference(path)
}

func (p *Project) Clone() *Project {
	return &Project{
		Name:             p.Name,
		Kind:             p.Kind,
		currentDirectory: p.currentDirectory,
		configFileName:   p.configFileName,
		configFilePath:   p.configFilePath,

		dirty:         p.dirty,
		dirtyFilePath: p.dirtyFilePath,

		host:            p.host,
		CommandLine:     p.CommandLine,
		Program:         p.Program,
		LanguageService: p.LanguageService,

		checkerPool: p.checkerPool,
	}
}

func (p *Project) CreateProgram() (*compiler.Program, *project.CheckerPool) {
	var programCloned bool
	var checkerPool *project.CheckerPool
	var newProgram *compiler.Program
	// oldProgram := p.Program
	if p.dirtyFilePath != "" {
		newProgram, programCloned = p.Program.UpdateProgram(p.dirtyFilePath, p.host)
		if !programCloned {
			// !!! wait until accepting snapshot to release documents!
			// !!! make this less janky
			// UpdateProgram called GetSourceFile (acquiring the document) but was unable to use it directly,
			// so it called NewProgram which acquired it a second time. We need to decrement the ref count
			// for the first acquisition.
			// p.snapshot.parseCache.releaseDocument(newProgram.GetSourceFileByPath(p.dirtyFilePath))
		}
	} else {
		newProgram = compiler.NewProgram(
			compiler.ProgramOptions{
				Host:                        p.host,
				Config:                      p.CommandLine,
				UseSourceOfProjectReference: true,
				TypingsLocation:             p.host.sessionOptions.TypingsLocation,
				JSDocParsingMode:            ast.JSDocParsingModeParseAll,
				CreateCheckerPool: func(program *compiler.Program) compiler.CheckerPool {
					checkerPool = project.NewCheckerPool(4, program, p.log)
					return checkerPool
				},
			},
		)
	}

	// if !programCloned && oldProgram != nil {
	// 	for _, file := range oldProgram.GetSourceFiles() {
	// 		// !!! wait until accepting snapshot to release documents!
	// 		// p.snapshot.parseCache.releaseDocument(file)
	// 	}
	// }

	return newProgram, checkerPool
}

func (p *Project) log(msg string) {
	// !!!
}
