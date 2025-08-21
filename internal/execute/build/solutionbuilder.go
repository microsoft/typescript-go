package build

import (
	"io"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/execute/tsc"
	"github.com/microsoft/typescript-go/internal/incremental"
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

func (o *Orchestrator) GetBuildOrderGenerator() *buildOrderGenerator {
	o.host = &solutionBuilderHost{
		builder: o,
		host:    compiler.NewCachedFSCompilerHost(o.opts.Sys.GetCurrentDirectory(), o.opts.Sys.FS(), o.opts.Sys.DefaultLibraryPath(), nil, nil),
	}
	return newBuildOrderGenerator(o.opts.Command, o.host, o.opts.Command.CompilerOptions.SingleThreaded.IsTrue())
}

func (o *Orchestrator) Start() tsc.CommandLineResult {
	o.host = &solutionBuilderHost{
		builder: o,
		host:    compiler.NewCachedFSCompilerHost(o.opts.Sys.GetCurrentDirectory(), o.opts.Sys.FS(), o.opts.Sys.DefaultLibraryPath(), nil, nil),
	}
	orderGenerator := newBuildOrderGenerator(o.opts.Command, o.host, o.opts.Command.CompilerOptions.SingleThreaded.IsTrue())
	return orderGenerator.buildOrClean(o, !o.opts.Command.BuildOptions.Clean.IsTrue())
}

func (o *Orchestrator) relativeFileName(fileName string) string {
	return tspath.ConvertToRelativePath(fileName, o.comparePathsOptions)
}

func (o *Orchestrator) toPath(fileName string) tspath.Path {
	return tspath.ToPath(fileName, o.comparePathsOptions.CurrentDirectory, o.comparePathsOptions.UseCaseSensitiveFileNames)
}

func (o *Orchestrator) getWriter(task *buildTask) io.Writer {
	if task == nil {
		return o.opts.Sys.Writer()
	}
	return &task.builder
}

func (o *Orchestrator) createBuilderStatusReporter(task *buildTask) tsc.DiagnosticReporter {
	return tsc.CreateBuilderStatusReporter(o.opts.Sys, o.getWriter(task), o.opts.Command.CompilerOptions, o.opts.Testing)
}

func (o *Orchestrator) createDiagnosticReporter(task *buildTask) tsc.DiagnosticReporter {
	return tsc.CreateDiagnosticReporter(o.opts.Sys, o.getWriter(task), o.opts.Command.CompilerOptions)
}

func (o *Orchestrator) buildProject(path tspath.Path, task *buildTask) {
	// Wait on upstream tasks to complete
	upStreamStatus := task.waitOnUpstream()
	status := o.getUpToDateStatus(path, task, upStreamStatus)
	task.reportUpToDateStatus(o, status)
	if handled := o.handleStatusThatDoesntRequireBuild(task, status); handled == nil {
		if o.opts.Command.BuildOptions.Verbose.IsTrue() {
			task.reportStatus(ast.NewCompilerDiagnostic(diagnostics.Building_project_0, o.relativeFileName(task.config)))
		}

		// Real build
		var compileTimes tsc.CompileTimes
		configAndTime, _ := o.host.resolvedReferences.Load(path)
		compileTimes.ConfigTime = configAndTime.time
		buildInfoReadStart := o.opts.Sys.Now()
		oldProgram := incremental.ReadBuildInfoProgram(task.resolved, o.host, o.host)
		compileTimes.BuildInfoReadTime = o.opts.Sys.Now().Sub(buildInfoReadStart)
		parseStart := o.opts.Sys.Now()
		program := compiler.NewProgram(compiler.ProgramOptions{
			Config: task.resolved,
			Host: &compilerHostForTaskReporter{
				host:  o.host,
				trace: tsc.GetTraceWithWriterFromSys(&task.builder, o.opts.Testing),
			},
			JSDocParsingMode: ast.JSDocParsingModeParseForTypeErrors,
		})
		compileTimes.ParseTime = o.opts.Sys.Now().Sub(parseStart)
		changesComputeStart := o.opts.Sys.Now()
		task.program = incremental.NewProgram(program, oldProgram, o.host, o.opts.Testing != nil)
		compileTimes.ChangesComputeTime = o.opts.Sys.Now().Sub(changesComputeStart)

		result, statistics := tsc.EmitAndReportStatistics(
			o.opts.Sys,
			task.program,
			program,
			task.resolved,
			task.reportDiagnostic,
			tsc.QuietDiagnosticsReporter,
			&task.builder,
			compileTimes,
			o.opts.Testing,
		)
		task.exitStatus = result.Status
		task.statistics = statistics
		if (!program.Options().NoEmitOnError.IsTrue() || len(result.Diagnostics) == 0) &&
			(len(result.EmitResult.EmittedFiles) > 0 || status.kind != upToDateStatusTypeOutOfDateBuildInfoWithErrors) {
			// Update time stamps for rest of the outputs
			o.updateTimeStamps(task, result.EmitResult.EmittedFiles, diagnostics.Updating_unchanged_output_timestamps_of_project_0)
		}

		if result.Status == tsc.ExitStatusDiagnosticsPresent_OutputsSkipped || result.Status == tsc.ExitStatusDiagnosticsPresent_OutputsGenerated {
			status = &upToDateStatus{kind: upToDateStatusTypeBuildErrors}
		} else {
			status = &upToDateStatus{kind: upToDateStatusTypeUpToDate}
		}
	} else {
		status = handled
		if task.resolved != nil {
			for _, diagnostic := range task.resolved.GetConfigFileParsingDiagnostics() {
				task.reportDiagnostic(diagnostic)
			}
		}
		if len(task.errors) > 0 {
			task.exitStatus = tsc.ExitStatusDiagnosticsPresent_OutputsSkipped
		}
	}
	task.unblockDownstream(status)
}

func (o *Orchestrator) handleStatusThatDoesntRequireBuild(task *buildTask, status *upToDateStatus) *upToDateStatus {
	switch status.kind {
	case upToDateStatusTypeUpToDate:
		if o.opts.Command.BuildOptions.Dry.IsTrue() {
			task.reportStatus(ast.NewCompilerDiagnostic(diagnostics.Project_0_is_up_to_date, task.config))
		}
		return status
	case upToDateStatusTypeUpstreamErrors:
		upstreamStatus := status.data.(*upstreamErrors)
		if o.opts.Command.BuildOptions.Verbose.IsTrue() {
			task.reportStatus(ast.NewCompilerDiagnostic(
				core.IfElse(
					upstreamStatus.refHasUpstreamErrors,
					diagnostics.Skipping_build_of_project_0_because_its_dependency_1_was_not_built,
					diagnostics.Skipping_build_of_project_0_because_its_dependency_1_has_errors,
				),
				o.relativeFileName(task.config),
				o.relativeFileName(upstreamStatus.ref),
			))
		}
		return status
	case upToDateStatusTypeSolution:
		return status
	case upToDateStatusTypeConfigFileNotFound:
		task.reportDiagnostic(ast.NewCompilerDiagnostic(diagnostics.File_0_not_found, task.config))
		return status
	}

	// update timestamps
	if status.isPseudoBuild() {
		if o.opts.Command.BuildOptions.Dry.IsTrue() {
			task.reportStatus(ast.NewCompilerDiagnostic(diagnostics.A_non_dry_build_would_update_timestamps_for_output_of_project_0, task.config))
			status = &upToDateStatus{kind: upToDateStatusTypeUpToDate}
			return status
		}

		o.updateTimeStamps(task, nil, diagnostics.Updating_output_timestamps_of_project_0)
		status = &upToDateStatus{kind: upToDateStatusTypeUpToDate}
		task.pseudoBuild = true
		return status
	}

	if o.opts.Command.BuildOptions.Dry.IsTrue() {
		task.reportStatus(ast.NewCompilerDiagnostic(diagnostics.A_non_dry_build_would_build_project_0, task.config))
		status = &upToDateStatus{kind: upToDateStatusTypeUpToDate}
		return status
	}
	return nil
}

func (o *Orchestrator) getUpToDateStatus(configPath tspath.Path, task *buildTask, upStreamStatus []*upToDateStatus) *upToDateStatus {
	// Config file not found
	if task.resolved == nil {
		return &upToDateStatus{kind: upToDateStatusTypeConfigFileNotFound}
	}

	// Solution - nothing to build
	if len(task.resolved.FileNames()) == 0 && task.resolved.ProjectReferences() != nil {
		return &upToDateStatus{kind: upToDateStatusTypeSolution}
	}

	for index, upstreamStatus := range upStreamStatus {
		if upstreamStatus == nil {
			// Not dependent on this upstream project (expected cycle was detected and hence skipped)
			continue
		}

		if o.opts.Command.BuildOptions.StopBuildOnErrors.IsTrue() && upstreamStatus.isError() {
			// Upstream project has errors, so we cannot build this project
			return &upToDateStatus{kind: upToDateStatusTypeUpstreamErrors, data: &upstreamErrors{task.resolved.ProjectReferences()[index].Path, upstreamStatus.kind == upToDateStatusTypeUpstreamErrors}}
		}
	}

	if o.opts.Command.BuildOptions.Force.IsTrue() {
		return &upToDateStatus{kind: upToDateStatusTypeForceBuild}
	}

	// Check the build info
	buildInfoPath := task.resolved.GetBuildInfoFileName()
	buildInfo := o.host.readOrStoreBuildInfo(configPath, buildInfoPath)
	if buildInfo == nil {
		return &upToDateStatus{kind: upToDateStatusTypeOutputMissing, data: buildInfoPath}
	}

	// build info version
	if !buildInfo.IsValidVersion() {
		return &upToDateStatus{kind: upToDateStatusTypeTsVersionOutputOfDate, data: buildInfo.Version}
	}

	// Report errors if build info indicates errors
	if buildInfo.Errors || // Errors that need to be reported irrespective of "--noCheck"
		(!task.resolved.CompilerOptions().NoCheck.IsTrue() && (buildInfo.SemanticErrors || buildInfo.CheckPending)) { // Errors without --noCheck
		return &upToDateStatus{kind: upToDateStatusTypeOutOfDateBuildInfoWithErrors, data: buildInfoPath}
	}

	if task.resolved.CompilerOptions().IsIncremental() {
		if !buildInfo.IsIncremental() {
			// Program options out of date
			return &upToDateStatus{kind: upToDateStatusTypeOutOfDateOptions, data: buildInfoPath}
		}

		// Errors need to be reported if build info has errors
		if (task.resolved.CompilerOptions().GetEmitDeclarations() && buildInfo.EmitDiagnosticsPerFile != nil) || // Always reported errors
			(!task.resolved.CompilerOptions().NoCheck.IsTrue() && // Semantic errors if not --noCheck
				(buildInfo.ChangeFileSet != nil || buildInfo.SemanticDiagnosticsPerFile != nil)) {
			return &upToDateStatus{kind: upToDateStatusTypeOutOfDateBuildInfoWithErrors, data: buildInfoPath}
		}

		// Pending emit files
		if !task.resolved.CompilerOptions().NoEmit.IsTrue() &&
			(buildInfo.ChangeFileSet != nil || buildInfo.AffectedFilesPendingEmit != nil) {
			return &upToDateStatus{kind: upToDateStatusTypeOutOfDateBuildInfoWithPendingEmit, data: buildInfoPath}
		}

		// Some of the emit files like source map or dts etc are not yet done
		if buildInfo.IsEmitPending(task.resolved, tspath.GetDirectoryPath(tspath.GetNormalizedAbsolutePath(buildInfoPath, o.comparePathsOptions.CurrentDirectory))) {
			return &upToDateStatus{kind: upToDateStatusTypeOutOfDateOptions, data: buildInfoPath}
		}
	}
	var inputTextUnchanged bool
	oldestOutputFileAndTime := fileAndTime{buildInfoPath, o.host.GetMTime(buildInfoPath)}
	var newestInputFileAndTime fileAndTime
	var seenRoots collections.Set[tspath.Path]
	var buildInfoRootInfoReader *incremental.BuildInfoRootInfoReader
	for _, inputFile := range task.resolved.FileNames() {
		inputTime := o.host.GetMTime(inputFile)
		if inputTime.IsZero() {
			return &upToDateStatus{kind: upToDateStatusTypeInputFileMissing, data: inputFile}
		}
		inputPath := o.toPath(inputFile)
		if inputTime.After(oldestOutputFileAndTime.time) {
			var version string
			var currentVersion string
			if buildInfo.IsIncremental() {
				if buildInfoRootInfoReader == nil {
					buildInfoRootInfoReader = buildInfo.GetBuildInfoRootInfoReader(tspath.GetDirectoryPath(tspath.GetNormalizedAbsolutePath(buildInfoPath, o.comparePathsOptions.CurrentDirectory)), o.comparePathsOptions)
				}
				buildInfoFileInfo, resolvedInputPath := buildInfoRootInfoReader.GetBuildInfoFileInfo(inputPath)
				if fileInfo := buildInfoFileInfo.GetFileInfo(); fileInfo != nil && fileInfo.Version() != "" {
					version = fileInfo.Version()
					if text, ok := o.host.FS().ReadFile(string(resolvedInputPath)); ok {
						currentVersion = incremental.ComputeHash(text, o.opts.Testing != nil)
						if version == currentVersion {
							inputTextUnchanged = true
						}
					}
				}
			}

			if version == "" || version != currentVersion {
				return &upToDateStatus{kind: upToDateStatusTypeInputFileNewer, data: &inputOutputName{inputFile, buildInfoPath}}
			}
		}
		if inputTime.After(newestInputFileAndTime.time) {
			newestInputFileAndTime = fileAndTime{inputFile, inputTime}
		}
		seenRoots.Add(inputPath)
	}

	if buildInfoRootInfoReader == nil {
		buildInfoRootInfoReader = buildInfo.GetBuildInfoRootInfoReader(tspath.GetDirectoryPath(tspath.GetNormalizedAbsolutePath(buildInfoPath, o.comparePathsOptions.CurrentDirectory)), o.comparePathsOptions)
	}
	for root := range buildInfoRootInfoReader.Roots() {
		if !seenRoots.Has(root) {
			// File was root file when project was built but its not any more
			return &upToDateStatus{kind: upToDateStatusTypeOutOfDateRoots, data: &inputOutputName{string(root), buildInfoPath}}
		}
	}

	if !task.resolved.CompilerOptions().IsIncremental() {
		// Check output file stamps
		for outputFile := range task.resolved.GetOutputFileNames() {
			outputTime := o.host.GetMTime(outputFile)
			if outputTime.IsZero() {
				// Output file missing
				return &upToDateStatus{kind: upToDateStatusTypeOutputMissing, data: outputFile}
			}

			if outputTime.Before(newestInputFileAndTime.time) {
				// Output file is older than input file
				return &upToDateStatus{kind: upToDateStatusTypeInputFileNewer, data: &inputOutputName{newestInputFileAndTime.file, outputFile}}
			}

			if outputTime.Before(oldestOutputFileAndTime.time) {
				oldestOutputFileAndTime = fileAndTime{outputFile, outputTime}
			}
		}
	}

	var refDtsUnchanged bool
	for index, upstreamStatus := range upStreamStatus {
		if upstreamStatus == nil || upstreamStatus.kind == upToDateStatusTypeSolution {
			// Not dependent on the status or this upstream project
			// (eg: expected cycle was detected and hence skipped, or is solution)
			continue
		}

		// If the upstream project's newest file is older than our oldest output,
		// we can't be out of date because of it
		// inputTime will not be present if we just built this project or updated timestamps
		// - in that case we do want to either build or update timestamps
		refInputOutputFileAndTime := upstreamStatus.inputOutputFileAndTime()
		if refInputOutputFileAndTime != nil && !refInputOutputFileAndTime.input.time.IsZero() && refInputOutputFileAndTime.input.time.Before(oldestOutputFileAndTime.time) {
			continue
		}

		// Check if tsbuildinfo path is shared, then we need to rebuild
		if o.host.hasConflictingBuildInfo(configPath) {
			return &upToDateStatus{kind: upToDateStatusTypeInputFileNewer, data: &inputOutputName{task.resolved.ProjectReferences()[index].Path, oldestOutputFileAndTime.file}}
		}

		// If the upstream project has only change .d.ts files, and we've built
		// *after* those files, then we're "pseudo up to date" and eligible for a fast rebuild
		newestDtsChangeTime := o.host.getLatestChangedDtsMTime(task.resolved.ResolvedProjectReferencePaths()[index])
		if !newestDtsChangeTime.IsZero() && newestDtsChangeTime.Before(oldestOutputFileAndTime.time) {
			refDtsUnchanged = true
			continue
		}

		// We have an output older than an upstream output - we are out of date
		return &upToDateStatus{kind: upToDateStatusTypeInputFileNewer, data: &inputOutputName{task.resolved.ProjectReferences()[index].Path, oldestOutputFileAndTime.file}}
	}

	configStatus := o.checkInputFileTime(task.config, &oldestOutputFileAndTime)
	if configStatus != nil {
		return configStatus
	}

	for _, extendedConfig := range task.resolved.ExtendedSourceFiles() {
		extendedConfigStatus := o.checkInputFileTime(extendedConfig, &oldestOutputFileAndTime)
		if extendedConfigStatus != nil {
			return extendedConfigStatus
		}
	}

	// !!! sheetal TODO : watch??
	// // Check package file time
	// const packageJsonLookups = state.lastCachedPackageJsonLookups.get(resolvedPath);
	// const dependentPackageFileStatus = packageJsonLookups && forEachKey(
	//     packageJsonLookups,
	//     path => checkConfigFileUpToDateStatus(state, path, oldestOutputFileTime, oldestOutputFileName),
	// );
	// if (dependentPackageFileStatus) return dependentPackageFileStatus;

	return &upToDateStatus{
		kind: core.IfElse(
			refDtsUnchanged,
			upToDateStatusTypeUpToDateWithUpstreamTypes,
			core.IfElse(inputTextUnchanged, upToDateStatusTypeUpToDateWithInputFileText, upToDateStatusTypeUpToDate),
		),
		data: &inputOutputFileAndTime{newestInputFileAndTime, oldestOutputFileAndTime, buildInfoPath},
	}
}

func (o *Orchestrator) checkInputFileTime(inputFile string, oldestOutputFileAndTime *fileAndTime) *upToDateStatus {
	inputTime := o.host.GetMTime(inputFile)
	if inputTime.After(oldestOutputFileAndTime.time) {
		// Output file is older than input file
		return &upToDateStatus{kind: upToDateStatusTypeInputFileNewer, data: &inputOutputName{inputFile, oldestOutputFileAndTime.file}}
	}
	return nil
}

func (o *Orchestrator) updateTimeStamps(task *buildTask, emittedFiles []string, verboseMessage *diagnostics.Message) {
	if task.resolved.CompilerOptions().NoEmit.IsTrue() {
		return
	}
	emitted := collections.NewSetFromItems(emittedFiles...)
	var verboseMessageReported bool
	updateTimeStamp := func(file string) {
		if emitted.Has(file) {
			return
		}
		if !verboseMessageReported && o.opts.Command.BuildOptions.Verbose.IsTrue() {
			task.reportStatus(ast.NewCompilerDiagnostic(verboseMessage, o.relativeFileName(task.config)))
			verboseMessageReported = true
		}
		err := o.host.SetMTime(file, o.opts.Sys.Now())
		if err != nil {
			task.reportDiagnostic(ast.NewCompilerDiagnostic(diagnostics.Failed_to_update_timestamp_of_file_0, file))
		}
	}

	if task.resolved.CompilerOptions().IsIncremental() {
		updateTimeStamp(task.resolved.GetBuildInfoFileName())
	} else {
		for outputFile := range task.resolved.GetOutputFileNames() {
			updateTimeStamp(outputFile)
		}
	}
}

func (o *Orchestrator) cleanProject(path tspath.Path, task *buildTask) {
	if task.resolved == nil {
		task.reportDiagnostic(ast.NewCompilerDiagnostic(diagnostics.File_0_not_found, task.config))
		task.exitStatus = tsc.ExitStatusDiagnosticsPresent_OutputsSkipped
		return
	}

	inputs := collections.NewSetFromItems(core.Map(task.resolved.FileNames(), o.toPath)...)
	for outputFile := range task.resolved.GetOutputFileNames() {
		o.cleanProjectOutput(task, outputFile, inputs)
	}
	o.cleanProjectOutput(task, task.resolved.GetBuildInfoFileName(), inputs)
}

func (o *Orchestrator) cleanProjectOutput(task *buildTask, outputFile string, inputs *collections.Set[tspath.Path]) {
	outputPath := o.toPath(outputFile)
	// If output name is same as input file name, do not delete and ignore the error
	if inputs.Has(outputPath) {
		return
	}
	if o.host.FS().FileExists(outputFile) {
		if !o.opts.Command.BuildOptions.Dry.IsTrue() {
			err := o.host.FS().Remove(outputFile)
			if err != nil {
				task.reportDiagnostic(ast.NewCompilerDiagnostic(diagnostics.Failed_to_delete_file_0, outputFile))
			}
		} else {
			task.filesToDelete = append(task.filesToDelete, outputFile)
		}
	}
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
