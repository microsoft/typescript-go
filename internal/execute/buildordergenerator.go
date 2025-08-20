package execute

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type solutionBuilderResult struct {
	result        CommandLineResult
	errors        []*ast.Diagnostic
	statistics    statistics
	programStats  []*statistics
	filesToDelete []string
}

func (b *solutionBuilderResult) report(s *solutionBuilder) {
	createReportErrorSummary(s.opts.sys, s.opts.command.CompilerOptions)(b.errors)
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
	if !s.opts.command.CompilerOptions.Diagnostics.IsTrue() && !s.opts.command.CompilerOptions.ExtendedDiagnostics.IsTrue() {
		return
	}
	b.statistics.isAggregate = true
	b.statistics.compileTimes = &compileTimes{}
	for _, stat := range b.programStats {
		// Aggregate statistics
		b.statistics.files += stat.files
		b.statistics.lines += stat.lines
		b.statistics.identifiers += stat.identifiers
		b.statistics.symbols += stat.symbols
		b.statistics.types += stat.types
		b.statistics.instantiations += stat.instantiations
		b.statistics.memoryUsed += stat.memoryUsed
		b.statistics.memoryAllocs += stat.memoryAllocs
		b.statistics.compileTimes.configTime += stat.compileTimes.configTime
		b.statistics.compileTimes.buildInfoReadTime += stat.compileTimes.buildInfoReadTime
		b.statistics.compileTimes.parseTime += stat.compileTimes.parseTime
		b.statistics.compileTimes.bindTime += stat.compileTimes.bindTime
		b.statistics.compileTimes.checkTime += stat.compileTimes.checkTime
		b.statistics.compileTimes.emitTime += stat.compileTimes.emitTime
		b.statistics.compileTimes.changesComputeTime += stat.compileTimes.changesComputeTime
	}
	b.statistics.compileTimes.totalTime = s.opts.sys.SinceStart()
	b.statistics.report(s.opts.sys.Writer(), s.opts.testing)
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
	return core.MapFiltered(task.upStream, func(t *statusTask) (string, bool) {
		return t.config, t.status != nil
	})
}

func (b *buildOrderGenerator) Downstream(configName string) []string {
	path := b.toPath(configName)
	task, ok := b.tasks.Load(path)
	if !ok {
		panic("No build task found for " + configName)
	}
	return core.Map(task.downStream, func(t *statusTask) string {
		return t.referencedBy
	})
}

func NewBuildOrderGenerator(command *tsoptions.ParsedBuildCommandLine, host compiler.CompilerHost, isSingleThreaded bool) *buildOrderGenerator {
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
		b.analyzeConfig(project, nil, false, &completed, &analyzing, circularityStack)
	}
}

func (b *buildOrderGenerator) analyzeConfig(
	configName string,
	downStream *statusTask,
	inCircularContext bool,
	completed *collections.Set[tspath.Path],
	analyzing *collections.Set[tspath.Path],
	circularityStack []string,
) {
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
			return
		}
		analyzing.Add(path)
		circularityStack = append(circularityStack, configName)
		if task.resolved != nil {
			for index, subReference := range task.resolved.ResolvedProjectReferencePaths() {
				statusTask := statusTask{config: subReference, referencedBy: configName}
				task.upStream = append(task.upStream, &statusTask)
				b.analyzeConfig(subReference, &statusTask, inCircularContext || task.resolved.ProjectReferences()[index].Circular, completed, analyzing, circularityStack)
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
		b.order = append(b.order, configName)
	}
	if downStream != nil {
		task.downStream = append(task.downStream, downStream)
		downStream.status = make(chan *upToDateStatus, 1)
	}
}

func (b *buildOrderGenerator) buildOrClean(builder *solutionBuilder, build bool) CommandLineResult {
	if build && builder.opts.command.BuildOptions.Verbose.IsTrue() {
		builder.createBuilderStatusReporter(nil)(ast.NewCompilerDiagnostic(
			diagnostics.Projects_in_this_build_Colon_0,
			strings.Join(core.Map(b.Order(), func(p string) string {
				return "\r\n    * " + builder.relativeFileName(p)
			}), ""),
		))
	}
	var buildResult solutionBuilderResult
	if len(b.errors) == 0 {
		wg := core.NewWorkGroup(builder.opts.command.CompilerOptions.SingleThreaded.IsTrue())
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
		buildResult.statistics.projects = len(b.Order())
	} else {
		buildResult.result.Status = ExitStatusProjectReferenceCycle_OutputsSkipped
		reportDiagnostic := builder.createDiagnosticReporter(nil)
		for _, err := range b.errors {
			reportDiagnostic(err)
		}
		buildResult.errors = b.errors
	}
	buildResult.report(builder)
	return buildResult.result
}
