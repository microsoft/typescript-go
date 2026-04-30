package build

import (
	"context"
	"io"
	"maps"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/execute/tsc"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs/vfswatch"
)

type Options struct {
	Sys     tsc.System
	Command *tsoptions.ParsedBuildCommandLine
	Testing tsc.CommandLineTesting
}

type orchestratorResult struct {
	result        tsc.CommandLineResult
	errors        []*ast.Diagnostic
	statistics    tsc.Statistics
	filesToDelete []string
}

func (b *orchestratorResult) report(o *Orchestrator) {
	if o.opts.Command.CompilerOptions.Watch.IsTrue() {
		o.watchStatusReporter(ast.NewCompilerDiagnostic(core.IfElse(len(b.errors) == 1, diagnostics.Found_1_error_Watching_for_file_changes, diagnostics.Found_0_errors_Watching_for_file_changes), len(b.errors)))
	} else {
		o.errorSummaryReporter(b.errors)
	}
	if b.filesToDelete != nil {
		o.createBuilderStatusReporter(nil)(
			ast.NewCompilerDiagnostic(
				diagnostics.A_non_dry_build_would_delete_the_following_files_Colon_0,
				strings.Join(core.Map(b.filesToDelete, func(f string) string {
					return "\r\n * " + f
				}), ""),
			))
	}
	if !o.opts.Command.CompilerOptions.Diagnostics.IsTrue() && !o.opts.Command.CompilerOptions.ExtendedDiagnostics.IsTrue() {
		return
	}
	b.statistics.SetTotalTime(o.opts.Sys.SinceStart())
	b.statistics.Report(o.opts.Sys.Writer(), o.opts.Testing)
}

type Orchestrator struct {
	opts                Options
	comparePathsOptions tspath.ComparePathsOptions
	host                *host

	// order generation result
	tasks  *collections.SyncMap[tspath.Path, *BuildTask]
	order  []string
	errors []*ast.Diagnostic

	errorSummaryReporter tsc.DiagnosticsReporter
	watchStatusReporter  tsc.DiagnosticReporter

	fileWatcher *vfswatch.FileWatcher
}

var _ tsc.Watcher = (*Orchestrator)(nil)

func (o *Orchestrator) relativeFileName(fileName string) string {
	return tspath.ConvertToRelativePath(fileName, o.comparePathsOptions)
}

func (o *Orchestrator) toPath(fileName string) tspath.Path {
	return tspath.ToPath(fileName, o.comparePathsOptions.CurrentDirectory, o.comparePathsOptions.UseCaseSensitiveFileNames)
}

func (o *Orchestrator) Order() []string {
	return o.order
}

func (o *Orchestrator) Upstream(configName string) []string {
	path := o.toPath(configName)
	task := o.getTask(path)
	return core.Map(task.upStream, func(t *upstreamTask) string {
		return t.task.config
	})
}

func (o *Orchestrator) Downstream(configName string) []string {
	path := o.toPath(configName)
	task := o.getTask(path)
	return core.Map(task.downStream, func(t *BuildTask) string {
		return t.config
	})
}

func (o *Orchestrator) getTask(path tspath.Path) *BuildTask {
	task, ok := o.tasks.Load(path)
	if !ok {
		panic("No build task found for " + path)
	}
	return task
}

func (o *Orchestrator) createBuildTasks(oldTasks *collections.SyncMap[tspath.Path, *BuildTask], configs []string, wg core.WorkGroup) {
	for _, config := range configs {
		wg.Queue(func() {
			path := o.toPath(config)
			var task *BuildTask
			var buildInfo *buildInfoEntry
			if oldTasks != nil {
				if existing, ok := oldTasks.Load(path); ok {
					if !existing.dirty {
						// Reuse existing task if config is same
						task = existing
					} else {
						buildInfo = existing.buildInfoEntry
					}
				}
			}
			if task == nil {
				task = &BuildTask{config: config, isInitialCycle: oldTasks == nil}
				task.pending.Store(true)
				task.buildInfoEntry = buildInfo
			}
			if _, loaded := o.tasks.LoadOrStore(path, task); loaded {
				return
			}
			task.resolved = o.host.GetResolvedProjectReference(config, path)
			task.upStream = nil
			if task.resolved != nil {
				o.createBuildTasks(oldTasks, task.resolved.ResolvedProjectReferencePaths(), wg)
			}
		})
	}
}

func (o *Orchestrator) setupBuildTask(
	configName string,
	downStream *BuildTask,
	inCircularContext bool,
	completed *collections.Set[tspath.Path],
	analyzing *collections.Set[tspath.Path],
	circularityStack []string,
) *BuildTask {
	path := o.toPath(configName)
	task := o.getTask(path)
	if !completed.Has(path) {
		if analyzing.Has(path) {
			if !inCircularContext {
				o.errors = append(o.errors, ast.NewCompilerDiagnostic(
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
				upstream := o.setupBuildTask(subReference, task, inCircularContext || task.resolved.ProjectReferences()[index].Circular, completed, analyzing, circularityStack)
				if upstream != nil {
					task.upStream = append(task.upStream, &upstreamTask{task: upstream, refIndex: index})
				}
			}
		}
		circularityStack = circularityStack[:len(circularityStack)-1]
		completed.Add(path)
		task.reportDone = make(chan struct{})
		prev := core.LastOrNil(o.order)
		if prev != "" {
			task.prevReporter = o.getTask(o.toPath(prev))
		}
		task.done = make(chan struct{})
		o.order = append(o.order, configName)
	}
	if o.opts.Command.CompilerOptions.Watch.IsTrue() && downStream != nil {
		task.downStream = append(task.downStream, downStream)
	}
	return task
}

func (o *Orchestrator) GenerateGraphReusingOldTasks() {
	tasks := o.tasks
	o.tasks = &collections.SyncMap[tspath.Path, *BuildTask]{}
	o.order = nil
	o.errors = nil
	o.GenerateGraph(tasks)
}

func (o *Orchestrator) GenerateGraph(oldTasks *collections.SyncMap[tspath.Path, *BuildTask]) {
	projects := o.opts.Command.ResolvedProjectPaths()
	// Parse all config files in parallel
	wg := core.NewWorkGroup(o.opts.Command.CompilerOptions.SingleThreaded.IsTrue())
	o.createBuildTasks(oldTasks, projects, wg)
	wg.RunAndWait()

	// Generate the graph
	completed := collections.Set[tspath.Path]{}
	analyzing := collections.Set[tspath.Path]{}
	circularityStack := []string{}
	for _, project := range projects {
		o.setupBuildTask(project, nil, false, &completed, &analyzing, circularityStack)
	}
}

func (o *Orchestrator) Start() tsc.CommandLineResult {
	if o.opts.Command.CompilerOptions.Watch.IsTrue() {
		o.watchStatusReporter(ast.NewCompilerDiagnostic(diagnostics.Starting_compilation_in_watch_mode))
	}
	o.GenerateGraph(nil)
	result := o.buildOrClean()
	if o.opts.Command.CompilerOptions.Watch.IsTrue() {
		o.Watch()
		result.Watcher = o
	}
	return result
}

func (o *Orchestrator) Watch() {
	o.fileWatcher = vfswatch.NewFileWatcher(
		o.opts.Sys.FS(),
		o.opts.Command.WatchOptions.WatchInterval(),
		o.opts.Testing != nil,
		nil, // build mode uses ScanForChanges, not callback
	)
	o.updateWatch()
	o.resetCaches()

	// Start watching for file changes
	if o.opts.Testing == nil {
		watchInterval := o.opts.Command.WatchOptions.WatchInterval()
		for {
			time.Sleep(watchInterval)
			o.DoCycle()
		}
	}
}

func (o *Orchestrator) updateWatch() {
	if c, ok := o.host.host.FS().(interface{ ClearCache() }); ok {
		c.ClearCache()
	}
	oldCache := o.host.mTimes
	o.host.mTimes = &collections.SyncMap[tspath.Path, time.Time]{}

	o.rangeTask(func(path tspath.Path, task *BuildTask) {
		task.updateWatchState(o, oldCache)
	})

	mergedDirs := o.computeWatchDirectories()
	o.fileWatcher.UpdateWatchedDirectories(mergedDirs)
}

func (o *Orchestrator) computeWatchDirectories() map[string]bool {
	dirs := make(map[string]bool)
	var mu sync.Mutex

	o.rangeTask(func(_ tspath.Path, task *BuildTask) {
		localDirs := make(map[string]bool)
		if task.resolved != nil {
			maps.Copy(localDirs, task.resolved.WildcardDirectories())
		}

		projectRoot := o.comparePathsOptions.CurrentDirectory
		if task.resolved != nil && task.resolved.ConfigFile != nil {
			projectRoot = tspath.GetDirectoryPath(task.resolved.ConfigFile.SourceFile.FileName())
		}

		internalDirs := make(map[string]struct{})
		externalDirs := make(map[string]struct{})
		if task.seenFiles != nil {
			task.seenFiles.Range(func(path string) bool {
				if !tspath.HasExtension(path) {
					return true
				}
				dir := tspath.GetDirectoryPath(path)
				if dir == "" {
					return true
				}
				if tspath.ContainsPath(projectRoot, path, o.comparePathsOptions) {
					internalDirs[dir] = struct{}{}
				} else {
					externalDirs[dir] = struct{}{}
				}
				return true
			})
		}

		// Coarsen internal directories to common parents
		if len(internalDirs) > 0 {
			fileDirs := make([]string, 0, len(internalDirs))
			for dir := range internalDirs {
				fileDirs = append(fileDirs, dir)
			}
			commonParents, _ := tspath.GetCommonParents(
				fileDirs,
				tspath.MinWatchLocationDepth,
				tspath.GetPathComponentsForWatching,
				o.comparePathsOptions,
			)
			for _, parent := range commonParents {
				if !isCoveredByWildcardDir(parent, localDirs, o.comparePathsOptions) {
					localDirs[parent] = true
				}
			}
		}

		// Add external directories non-recursively
		for dir := range externalDirs {
			if !isCoveredByWildcardDir(dir, localDirs, o.comparePathsOptions) {
				if _, already := localDirs[dir]; !already {
					localDirs[dir] = false
				}
			}
		}

		mu.Lock()
		for dir, recursive := range localDirs {
			if recursive {
				dirs[dir] = true
			} else if _, already := dirs[dir]; !already {
				dirs[dir] = false
			}
		}
		mu.Unlock()
	})

	return dirs
}

func isCoveredByWildcardDir(dir string, wildcardDirs map[string]bool, opts tspath.ComparePathsOptions) bool {
	for wdir, recursive := range wildcardDirs {
		if tspath.ComparePaths(wdir, dir, opts) == 0 {
			return true
		}
		if recursive && tspath.ContainsPath(wdir, dir, opts) {
			return true
		}
	}
	return false
}

func (o *Orchestrator) resetCaches() {
	// Clean out all the caches
	if c, ok := o.host.host.FS().(interface{ ClearCache() }); ok {
		c.ClearCache()
	}
	o.host.extendedConfigCache = tsc.ExtendedConfigCache{}
	o.host.sourceFiles.reset()
	o.host.configTimes = collections.SyncMap[tspath.Path, time.Duration]{}
}

func (o *Orchestrator) DoCycle() {
	event := o.fileWatcher.ScanForChanges()
	if !event.HasChanges() {
		return
	}

	o.fileWatcher.WaitForSettled(context.Background())
	event = o.fileWatcher.ScanForChanges()

	var needsConfigUpdate atomic.Bool
	var needsUpdate atomic.Bool
	mTimes := o.host.mTimes.Clone()
	o.rangeTask(func(path tspath.Path, task *BuildTask) {
		if updateKind := task.hasUpdate(o, path, event); updateKind != updateKindNone {
			needsUpdate.Store(true)
			if updateKind == updateKindConfig {
				needsConfigUpdate.Store(true)
			}
		}
	})

	if !needsUpdate.Load() {
		o.host.mTimes = mTimes
		o.resetCaches()
		o.fileWatcher.UpdateWatchedDirectories(o.computeWatchDirectories())
		return
	}

	o.watchStatusReporter(ast.NewCompilerDiagnostic(diagnostics.File_change_detected_Starting_incremental_compilation))

	if c, ok := o.host.host.FS().(interface{ ClearCache() }); ok {
		c.ClearCache()
	}
	o.host.mTimes = &collections.SyncMap[tspath.Path, time.Time]{}

	if needsConfigUpdate.Load() {
		o.GenerateGraphReusingOldTasks()
	}

	o.buildOrClean()
	o.updateWatch()
	o.resetCaches()
}

func (o *Orchestrator) buildOrClean() tsc.CommandLineResult {
	if !o.opts.Command.BuildOptions.Clean.IsTrue() && o.opts.Command.BuildOptions.Verbose.IsTrue() {
		o.createBuilderStatusReporter(nil)(ast.NewCompilerDiagnostic(
			diagnostics.Projects_in_this_build_Colon_0,
			strings.Join(core.Map(o.Order(), func(p string) string {
				return "\r\n    * " + o.relativeFileName(p)
			}), ""),
		))
	}
	var buildResult orchestratorResult
	if len(o.errors) == 0 {
		buildResult.statistics.Projects = len(o.Order())
		o.rangeTask(func(path tspath.Path, task *BuildTask) {
			o.buildOrCleanProject(task, path, &buildResult)
		})
	} else {
		// Circularity errors prevent any project from being built
		buildResult.result.Status = tsc.ExitStatusProjectReferenceCycle_OutputsSkipped
		reportDiagnostic := o.createDiagnosticReporter(nil)
		for _, err := range o.errors {
			reportDiagnostic(err)
		}
		buildResult.errors = o.errors
	}
	buildResult.report(o)
	return buildResult.result
}

func (o *Orchestrator) rangeTask(f func(path tspath.Path, task *BuildTask)) {
	numRoutines := 4
	if o.opts.Command.CompilerOptions.SingleThreaded.IsTrue() {
		numRoutines = 1
	} else if builders := o.opts.Command.BuildOptions.Builders; builders != nil {
		numRoutines = *builders
	}

	var currentTaskIndex atomic.Int64
	getNextTask := func() (tspath.Path, *BuildTask, bool) {
		index := int(currentTaskIndex.Add(1) - 1)
		if index >= len(o.order) {
			return "", nil, false
		}
		config := o.order[index]
		path := o.toPath(config)
		task := o.getTask(path)
		return path, task, true
	}
	runTask := func() {
		for path, task, ok := getNextTask(); ok; path, task, ok = getNextTask() {
			f(path, task)
		}
	}

	if numRoutines == 1 {
		runTask()
	} else {
		wg := core.NewWorkGroup(false)
		for range numRoutines {
			wg.Queue(runTask)
		}
		wg.RunAndWait()
	}
}

func (o *Orchestrator) buildOrCleanProject(task *BuildTask, path tspath.Path, buildResult *orchestratorResult) {
	task.result = &taskResult{}
	task.result.reportStatus = o.createBuilderStatusReporter(task)
	task.result.diagnosticReporter = o.createDiagnosticReporter(task)
	if !o.opts.Command.BuildOptions.Clean.IsTrue() {
		task.buildProject(o, path)
	} else {
		task.cleanProject(o, path)
	}
	task.report(o, path, buildResult)
}

func (o *Orchestrator) getWriter(task *BuildTask) io.Writer {
	if task == nil {
		return o.opts.Sys.Writer()
	}
	return &task.result.builder
}

func (o *Orchestrator) createBuilderStatusReporter(task *BuildTask) tsc.DiagnosticReporter {
	return tsc.CreateBuilderStatusReporter(o.opts.Sys, o.getWriter(task), o.opts.Command.Locale(), o.opts.Command.CompilerOptions, o.opts.Testing)
}

func (o *Orchestrator) createDiagnosticReporter(task *BuildTask) tsc.DiagnosticReporter {
	return tsc.CreateDiagnosticReporter(o.opts.Sys, o.getWriter(task), o.opts.Command.Locale(), o.opts.Command.CompilerOptions)
}

func NewOrchestrator(opts Options) *Orchestrator {
	orchestrator := &Orchestrator{
		opts: opts,
		comparePathsOptions: tspath.ComparePathsOptions{
			CurrentDirectory:          opts.Sys.GetCurrentDirectory(),
			UseCaseSensitiveFileNames: opts.Sys.FS().UseCaseSensitiveFileNames(),
		},
		tasks: &collections.SyncMap[tspath.Path, *BuildTask]{},
	}
	orchestrator.host = &host{
		orchestrator: orchestrator,
		host: compiler.NewCachedFSCompilerHost(
			orchestrator.opts.Sys.GetCurrentDirectory(),
			orchestrator.opts.Sys.FS(),
			orchestrator.opts.Sys.DefaultLibraryPath(),
			nil,
			nil,
		),
		mTimes: &collections.SyncMap[tspath.Path, time.Time]{},
	}
	if opts.Command.CompilerOptions.Watch.IsTrue() {
		orchestrator.watchStatusReporter = tsc.CreateWatchStatusReporter(opts.Sys, opts.Command.Locale(), opts.Command.CompilerOptions, opts.Testing)
	} else {
		orchestrator.errorSummaryReporter = tsc.CreateReportErrorSummary(opts.Sys, opts.Command.Locale(), opts.Command.CompilerOptions)
	}
	return orchestrator
}
