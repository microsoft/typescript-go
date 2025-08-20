package execute

import (
	"fmt"
	"strings"
	"time"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/incremental"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type upToDateStatusType uint16

const (
	// Errors:

	// config file was not found
	upToDateStatusTypeConfigFileNotFound upToDateStatusType = iota
	// found errors during build
	upToDateStatusTypeBuildErrors
	// did not build because upstream project has errors - and we have option to stop build on upstream errors
	upToDateStatusTypeUpstreamErrors

	// Its all good, no work to do
	upToDateStatusTypeUpToDate

	// Pseudo-builds - touch timestamps, no actual build:

	// The project appears out of date because its upstream inputs are newer than its outputs,
	// but all of its outputs are actually newer than the previous identical outputs of its (.d.ts) inputs.
	// This means we can Pseudo-build (just touch timestamps), as if we had actually built this project.
	upToDateStatusTypeUpToDateWithUpstreamTypes
	// The project appears up to date and even though input file changed, its text didnt so just need to update timestamps
	upToDateStatusTypeUpToDateWithInputFileText

	// Needs build:

	// input file is missing
	upToDateStatusTypeInputFileMissing
	// output file is missing
	upToDateStatusTypeOutputMissing
	// input file is newer than output file
	upToDateStatusTypeInputFileNewer
	// build info is out of date as we need to emit some files
	upToDateStatusTypeOutOfDateBuildInfoWithPendingEmit
	// build info indiscates that project has errors and they need to be reported
	upToDateStatusTypeOutOfDateBuildInfoWithErrors
	// build info options indicate there is work to do based on changes in options
	upToDateStatusTypeOutOfDateOptions
	// file was root when built but not any more
	upToDateStatusTypeOutOfDateRoots
	// buildInfo.version mismatch with current ts version
	upToDateStatusTypeTsVersionOutputOfDate
	// build because --force was specified
	upToDateStatusTypeForceBuild

	// solution file
	upToDateStatusTypeSolution
)

type inputOutputName struct {
	input  string
	output string
}

type fileAndTime struct {
	file string
	time time.Time
}

type inputOutputFileAndTime struct {
	input     fileAndTime
	output    fileAndTime
	buildInfo string
}

type upstreamErrors struct {
	ref                  string
	refHasUpstreamErrors bool
}

type upToDateStatus struct {
	kind upToDateStatusType
	data any
}

func (s *upToDateStatus) IsError() bool {
	switch s.kind {
	case upToDateStatusTypeConfigFileNotFound,
		upToDateStatusTypeBuildErrors,
		upToDateStatusTypeUpstreamErrors:
		return true
	default:
		return false
	}
}

func (s *upToDateStatus) NeedsBuild() bool {
	switch s.kind {
	case upToDateStatusTypeInputFileMissing,
		upToDateStatusTypeOutputMissing,
		upToDateStatusTypeInputFileNewer,
		upToDateStatusTypeOutOfDateBuildInfoWithPendingEmit,
		upToDateStatusTypeOutOfDateBuildInfoWithErrors,
		upToDateStatusTypeOutOfDateOptions,
		upToDateStatusTypeOutOfDateRoots,
		upToDateStatusTypeTsVersionOutputOfDate,
		upToDateStatusTypeForceBuild:
		return true
	default:
		return false
	}
}

func (s *upToDateStatus) IsPseudoBuild() bool {
	switch s.kind {
	case upToDateStatusTypeUpToDateWithUpstreamTypes,
		upToDateStatusTypeUpToDateWithInputFileText:
		return true
	default:
		return false
	}
}

func (s *upToDateStatus) InputOutputFileAndTime() *inputOutputFileAndTime {
	data, ok := s.data.(*inputOutputFileAndTime)
	if !ok {
		return nil
	}
	return data
}

type statusTask struct {
	config       string
	referencedBy string
	status       chan *upToDateStatus
}

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

type taskReporter struct {
	builder            strings.Builder
	errors             []*ast.Diagnostic
	reportStatus       diagnosticReporter
	diagnosticReporter diagnosticReporter
	exitStatus         ExitStatus
	statistics         *statistics
	program            *incremental.Program
	filesToDelete      []string
}

func (b *taskReporter) reportDiagnostic(err *ast.Diagnostic) {
	b.errors = append(b.errors, err)
	b.diagnosticReporter(err)
}

func (b *taskReporter) report(s *solutionBuilder, configPath tspath.Path, buildResult *solutionBuilderResult) {
	if len(b.errors) > 0 {
		buildResult.errors = append(core.IfElse(buildResult.errors != nil, buildResult.errors, []*ast.Diagnostic{}), b.errors...)
	}
	fmt.Fprint(s.opts.sys.Writer(), b.builder.String())
	if b.exitStatus > buildResult.result.Status {
		buildResult.result.Status = b.exitStatus
	}
	if b.statistics != nil {
		buildResult.programStats = append(buildResult.programStats, b.statistics)
	}
	if b.program != nil {
		buildResult.result.IncrementalProgram = append(buildResult.result.IncrementalProgram, b.program)
		buildResult.statistics.projectsBuilt++
	}
	buildResult.filesToDelete = append(buildResult.filesToDelete, b.filesToDelete...)
}

type buildTask struct {
	config               string
	resolved             *tsoptions.ParsedCommandLine
	upStream             []*statusTask
	downStream           []*statusTask
	taskReporter         taskReporter
	previousTaskReporter chan *taskReporter
	reporter             chan *taskReporter
}

func (t *buildTask) waitOnUpstream() []*upToDateStatus {
	upStreamStatus := make([]*upToDateStatus, len(t.upStream))
	for i, upstream := range t.upStream {
		if upstream.status != nil {
			upStreamStatus[i] = <-upstream.status
		}
	}
	return upStreamStatus
}

func (t *buildTask) unblockDownstream(status *upToDateStatus) {
	for _, downstream := range t.downStream {
		downstream.status <- status
	}
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
		task.reporter = make(chan *taskReporter, 1)
		prev := core.LastOrNil(b.order)
		if prev != "" {
			if prevTask, ok := b.tasks.Load(b.toPath(prev)); ok {
				task.previousTaskReporter = prevTask.reporter
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
				// Wait for previous build task to complete reporting status, errors etc
				if task.previousTaskReporter != nil {
					<-task.previousTaskReporter
				}
				task.taskReporter.report(builder, path, &buildResult)
				task.reporter <- &task.taskReporter
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
