package build

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/execute/tsc"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type solutionBuilderResult struct {
	result        tsc.CommandLineResult
	errors        []*ast.Diagnostic
	statistics    tsc.Statistics
	programStats  []*tsc.Statistics
	filesToDelete []string
}

func (b *solutionBuilderResult) report(s *Orchestrator) {
	tsc.CreateReportErrorSummary(s.opts.Sys, s.opts.Command.CompilerOptions)(b.errors)
	if b.filesToDelete != nil {
		s.createBuilderStatusReporter(nil)(
			ast.NewCompilerDiagnostic(
				diagnostics.A_non_dry_build_would_delete_the_following_files_Colon_0,
				strings.Join(core.Map(b.filesToDelete, func(f string) string {
					return "\r\n * " + f
				}), ""),
			))
	}
	if len(b.programStats) == 0 {
		return
	}
	if !s.opts.Command.CompilerOptions.Diagnostics.IsTrue() && !s.opts.Command.CompilerOptions.ExtendedDiagnostics.IsTrue() {
		return
	}
	b.statistics.Aggregate(b.programStats, s.opts.Sys.SinceStart())
	b.statistics.Report(s.opts.Sys.Writer(), s.opts.Testing)
}

type buildOrderGenerator struct {
	host   compiler.CompilerHost
	tasks  collections.SyncMap[tspath.Path, *buildTask]
	order  []string
	errors []*ast.Diagnostic
}

func (b *buildOrderGenerator) Order() []string {
	return b.order
}

func (b *buildOrderGenerator) Upstream(configName string) []string {
	path := b.toPath(configName)
	task, ok := b.tasks.Load(path)
	if !ok {
		panic("No build task found for " + configName)
	}
	return core.Map(task.upStream, func(t *buildTask) string {
		return t.config
	})
}

func newBuildOrderGenerator(command *tsoptions.ParsedBuildCommandLine, host compiler.CompilerHost, isSingleThreaded bool) *buildOrderGenerator {
	b := &buildOrderGenerator{host: host}

	projects := command.ResolvedProjectPaths()
	// Parse all config files in parallel
	wg := core.NewWorkGroup(isSingleThreaded)
	b.createBuildTasks(projects, wg)
	wg.RunAndWait()

	// Generate the order
	b.generateOrder(projects)

	return b
}

func (b *buildOrderGenerator) toPath(configName string) tspath.Path {
	return tspath.ToPath(configName, b.host.GetCurrentDirectory(), b.host.FS().UseCaseSensitiveFileNames())
}

func (b *buildOrderGenerator) createBuildTasks(projects []string, wg core.WorkGroup) {
	for _, project := range projects {
		b.createBuildTask(project, wg)
	}
}

func (b *buildOrderGenerator) createBuildTask(configName string, wg core.WorkGroup) {
	wg.Queue(func() {
		path := b.toPath(configName)
		task := &buildTask{config: configName}
		if _, loaded := b.tasks.LoadOrStore(path, task); loaded {
			return
		}
		task.resolved = b.host.GetResolvedProjectReference(configName, path)
		if task.resolved != nil {
			b.createBuildTasks(task.resolved.ResolvedProjectReferencePaths(), wg)
		}
	})
}

func (b *buildOrderGenerator) generateOrder(projects []string) {
	completed := collections.Set[tspath.Path]{}
	analyzing := collections.Set[tspath.Path]{}
	circularityStack := []string{}
	for _, project := range projects {
		b.analyzeConfig(project, false, &completed, &analyzing, circularityStack)
	}
}

func (b *buildOrderGenerator) analyzeConfig(
	configName string,
	inCircularContext bool,
	completed *collections.Set[tspath.Path],
	analyzing *collections.Set[tspath.Path],
	circularityStack []string,
) *buildTask {
	path := b.toPath(configName)
	task, ok := b.tasks.Load(path)
	if !ok {
		panic("No build task found for " + configName)
	}
	if !completed.Has(path) {
		if analyzing.Has(path) {
			if !inCircularContext {
				b.errors = append(b.errors, ast.NewCompilerDiagnostic(
					diagnostics.Project_references_may_not_form_a_circular_graph_Cycle_detected_Colon_0,
					strings.Join(circularityStack, "\n"),
				))
			}
			return nil
		}
		analyzing.Add(path)
		circularityStack = append(circularityStack, configName)
		if task.resolved != nil {
			for index, subReference := range task.resolved.ResolvedProjectReferencePaths() {
				upstream := b.analyzeConfig(subReference, inCircularContext || task.resolved.ProjectReferences()[index].Circular, completed, analyzing, circularityStack)
				if upstream != nil {
					task.upStream = append(task.upStream, upstream)
				}
			}
		}
		circularityStack = circularityStack[:len(circularityStack)-1]
		completed.Add(path)
		task.taskReporter.done = make(chan struct{})
		prev := core.LastOrNil(b.order)
		if prev != "" {
			if prevTask, ok := b.tasks.Load(b.toPath(prev)); ok {
				task.taskReporter.prev = &prevTask.taskReporter
			} else {
				panic("No previous task found for " + prev)
			}
		}
		task.done = make(chan struct{})
		b.order = append(b.order, configName)
	}
	return task
}

func (b *buildOrderGenerator) buildOrClean(builder *Orchestrator, build bool) tsc.CommandLineResult {
	if build && builder.opts.Command.BuildOptions.Verbose.IsTrue() {
		builder.createBuilderStatusReporter(nil)(ast.NewCompilerDiagnostic(
			diagnostics.Projects_in_this_build_Colon_0,
			strings.Join(core.Map(b.Order(), func(p string) string {
				return "\r\n    * " + builder.relativeFileName(p)
			}), ""),
		))
	}
	var buildResult solutionBuilderResult
	if len(b.errors) == 0 {
		wg := core.NewWorkGroup(builder.opts.Command.CompilerOptions.SingleThreaded.IsTrue())
		b.tasks.Range(func(path tspath.Path, task *buildTask) bool {
			task.taskReporter.reportStatus = builder.createBuilderStatusReporter(&task.taskReporter)
			task.taskReporter.diagnosticReporter = builder.createDiagnosticReporter(&task.taskReporter)
			wg.Queue(func() {
				if build {
					builder.buildProject(path, task)
				} else {
					builder.cleanProject(path, task)
				}
				task.taskReporter.report(builder, path, &buildResult)
			})
			return true
		})
		wg.RunAndWait()
		buildResult.statistics.Projects = len(b.Order())
	} else {
		buildResult.result.Status = tsc.ExitStatusProjectReferenceCycle_OutputsSkipped
		reportDiagnostic := builder.createDiagnosticReporter(nil)
		for _, err := range b.errors {
			reportDiagnostic(err)
		}
		buildResult.errors = b.errors
	}
	buildResult.report(builder)
	return buildResult.result
}
