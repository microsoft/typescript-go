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

const inferredProjectName = "/dev/null/inferredProject"

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
	return NewProject(configFileName, KindConfigured, tspath.GetDirectoryPath(configFileName), builder)
}

func NewInferredProject(
	currentDirectory string,
	compilerOptions *core.CompilerOptions,
	rootFileNames []string,
	builder *projectCollectionBuilder,
) *Project {
	p := NewProject(inferredProjectName, KindInferred, currentDirectory, builder)
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
	configFileName string,
	kind Kind,
	currentDirectory string,
	builder *projectCollectionBuilder,
) *Project {
	project := &Project{
		configFileName:   configFileName,
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
	project.configFilePath = tspath.ToPath(configFileName, currentDirectory, builder.fs.fs.UseCaseSensitiveFileNames())
	return project
}

func (p *Project) Name() string {
	return p.configFileName
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
	return p.Program != nil && p.Program.GetSourceFileByPath(path) != nil
}

func (p *Project) IsSourceFromProjectReference(path tspath.Path) bool {
	return p.Program != nil && p.Program.IsSourceFromProjectReference(path)
}

func (p *Project) Clone() *Project {
	return &Project{
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
	if p.dirtyFilePath != "" && p.Program != nil && p.Program.CommandLine() == p.CommandLine {
		newProgram, programCloned = p.Program.UpdateProgram(p.dirtyFilePath, p.host)
		if programCloned {
			for _, file := range newProgram.GetSourceFiles() {
				if file.Path() != p.dirtyFilePath {
					// UpdateProgram only called host.GetSourceFile for the dirty file.
					// Increment ref count for all other files.
					p.host.builder.parseCache.Ref(file)
				}
			}
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

	return newProgram, checkerPool
}

func (p *Project) log(msg string) {
	// !!!
}
