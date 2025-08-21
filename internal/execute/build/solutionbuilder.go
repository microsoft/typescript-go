package build

import (
	"io"

	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/execute/tsc"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type Options struct {
	Sys     tsc.System
	Command *tsoptions.ParsedBuildCommandLine
	Testing tsc.CommandLineTesting
}

type Orchestrator struct {
	opts                Options
	comparePathsOptions tspath.ComparePathsOptions
	host                *solutionBuilderHost
}

func (orchestrator *Orchestrator) GetBuildOrderGenerator() *buildOrderGenerator {
	orchestrator.host = &solutionBuilderHost{
		builder: orchestrator,
		host:    compiler.NewCachedFSCompilerHost(orchestrator.opts.Sys.GetCurrentDirectory(), orchestrator.opts.Sys.FS(), orchestrator.opts.Sys.DefaultLibraryPath(), nil, nil),
	}
	return newBuildOrderGenerator(orchestrator.opts.Command, orchestrator.host, orchestrator.opts.Command.CompilerOptions.SingleThreaded.IsTrue())
}

func (orchestrator *Orchestrator) Start() tsc.CommandLineResult {
	orderGenerator := orchestrator.GetBuildOrderGenerator()
	return orderGenerator.buildOrClean(orchestrator, !orchestrator.opts.Command.BuildOptions.Clean.IsTrue())
}

func (orchestrator *Orchestrator) relativeFileName(fileName string) string {
	return tspath.ConvertToRelativePath(fileName, orchestrator.comparePathsOptions)
}

func (orchestrator *Orchestrator) toPath(fileName string) tspath.Path {
	return tspath.ToPath(fileName, orchestrator.comparePathsOptions.CurrentDirectory, orchestrator.comparePathsOptions.UseCaseSensitiveFileNames)
}

func (orchestrator *Orchestrator) getWriter(task *buildTask) io.Writer {
	if task == nil {
		return orchestrator.opts.Sys.Writer()
	}
	return &task.builder
}

func (orchestrator *Orchestrator) createBuilderStatusReporter(task *buildTask) tsc.DiagnosticReporter {
	return tsc.CreateBuilderStatusReporter(orchestrator.opts.Sys, orchestrator.getWriter(task), orchestrator.opts.Command.CompilerOptions, orchestrator.opts.Testing)
}

func (orchestrator *Orchestrator) createDiagnosticReporter(task *buildTask) tsc.DiagnosticReporter {
	return tsc.CreateDiagnosticReporter(orchestrator.opts.Sys, orchestrator.getWriter(task), orchestrator.opts.Command.CompilerOptions)
}

func NewOrchestrator(opts Options) *Orchestrator {
	return &Orchestrator{
		opts: opts,
		comparePathsOptions: tspath.ComparePathsOptions{
			CurrentDirectory:          opts.Sys.GetCurrentDirectory(),
			UseCaseSensitiveFileNames: opts.Sys.FS().UseCaseSensitiveFileNames(),
		},
	}
}
