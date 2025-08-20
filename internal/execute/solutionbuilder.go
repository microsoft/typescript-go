package execute

import (
	"fmt"
	"io"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/incremental"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type solutionBuilderOptions struct {
	sys     System
	command *tsoptions.ParsedBuildCommandLine
	testing CommandLineTesting
}

type solutionBuilder struct {
	opts                solutionBuilderOptions
	comparePathsOptions tspath.ComparePathsOptions
	host                *solutionBuilderHost
}

func (s *solutionBuilder) buildOrClean(build bool) CommandLineResult {
	s.host = &solutionBuilderHost{
		builder: s,
		host:    compiler.NewCachedFSCompilerHost(s.opts.sys.GetCurrentDirectory(), s.opts.sys.FS(), s.opts.sys.DefaultLibraryPath(), nil, nil),
	}
	orderGenerator := NewBuildOrderGenerator(s.opts.command, s.host, s.opts.command.CompilerOptions.SingleThreaded.IsTrue())
	return orderGenerator.buildOrClean(s, build)
}

func (s *solutionBuilder) relativeFileName(fileName string) string {
	return tspath.ConvertToRelativePath(fileName, s.comparePathsOptions)
}

func (s *solutionBuilder) toPath(fileName string) tspath.Path {
	return tspath.ToPath(fileName, s.comparePathsOptions.CurrentDirectory, s.comparePathsOptions.UseCaseSensitiveFileNames)
}

func (s *solutionBuilder) getWriter(taskReporter *taskReporter) io.Writer {
	if taskReporter == nil {
		return s.opts.sys.Writer()
	}
	return &taskReporter.builder
}

func (s *solutionBuilder) createBuilderStatusReporter(taskReporter *taskReporter) diagnosticReporter {
	return createBuilderStatusReporter(s.opts.sys, s.getWriter(taskReporter), s.opts.command.CompilerOptions, s.opts.testing)
}

func (s *solutionBuilder) createDiagnosticReporter(taskReporter *taskReporter) diagnosticReporter {
	return createDiagnosticReporter(s.opts.sys, s.getWriter(taskReporter), s.opts.command.CompilerOptions)
}

func (s *solutionBuilder) createTaskReporter() *taskReporter {
	var taskReporter taskReporter
	taskReporter.reportStatus = s.createBuilderStatusReporter(&taskReporter)
	taskReporter.diagnosticReporter = s.createDiagnosticReporter(&taskReporter)
	return &taskReporter
}

func (s *solutionBuilder) buildProject(config string, path tspath.Path, task *buildTask) *taskReporter {
	// Wait on upstream tasks to complete
	upStreamStatus := make([]*upToDateStatus, len(task.upStream))
	for i, upstream := range task.upStream {
		if upstream.status != nil {
			upStreamStatus[i] = <-upstream.status
		}
	}

	status := s.getUpToDateStatus(config, path, task, upStreamStatus)
	taskReporter := s.createTaskReporter()

	s.reportUpToDateStatus(config, status, taskReporter)
	if handled := s.handleStatusThatDoesntRequireBuild(config, task, status, taskReporter); handled == nil {
		if s.opts.command.BuildOptions.Verbose.IsTrue() {
			taskReporter.reportStatus(ast.NewCompilerDiagnostic(diagnostics.Building_project_0, s.relativeFileName(config)))
		}

		// Real build
		var compileTimes compileTimes
		configAndTime, _ := s.host.resolvedReferences.Load(path)
		compileTimes.configTime = configAndTime.time
		buildInfoReadStart := s.opts.sys.Now()
		oldProgram := incremental.ReadBuildInfoProgram(task.resolved, s.host, s.host)
		compileTimes.buildInfoReadTime = s.opts.sys.Now().Sub(buildInfoReadStart)
		parseStart := s.opts.sys.Now()
		program := compiler.NewProgram(compiler.ProgramOptions{
			Config: task.resolved,
			Host: &compilerHostForTaskReporter{
				host:  s.host,
				trace: getTraceWithWriterFromSys(&taskReporter.builder, s.opts.testing),
			},
			JSDocParsingMode: ast.JSDocParsingModeParseForTypeErrors,
		})
		compileTimes.parseTime = s.opts.sys.Now().Sub(parseStart)
		changesComputeStart := s.opts.sys.Now()
		taskReporter.program = incremental.NewProgram(program, oldProgram, s.host, s.opts.testing != nil)
		compileTimes.changesComputeTime = s.opts.sys.Now().Sub(changesComputeStart)

		result, statistics := emitAndReportStatistics(
			s.opts.sys,
			taskReporter.program,
			program,
			task.resolved,
			taskReporter.reportDiagnostic,
			quietDiagnosticsReporter,
			&taskReporter.builder,
			compileTimes,
			s.opts.testing,
		)
		taskReporter.exitStatus = result.status
		taskReporter.statistics = statistics
		if (!program.Options().NoEmitOnError.IsTrue() || len(result.diagnostics) == 0) &&
			(len(result.emitResult.EmittedFiles) > 0 || status.kind != upToDateStatusTypeOutOfDateBuildInfoWithErrors) {
			// Update time stamps for rest of the outputs
			s.updateTimeStamps(config, task, taskReporter, result.emitResult.EmittedFiles, diagnostics.Updating_unchanged_output_timestamps_of_project_0)
		}

		if result.status == ExitStatusDiagnosticsPresent_OutputsSkipped || result.status == ExitStatusDiagnosticsPresent_OutputsGenerated {
			status = &upToDateStatus{kind: upToDateStatusTypeBuildErrors}
		} else {
			status = &upToDateStatus{kind: upToDateStatusTypeUpToDate}
		}
	} else {
		status = handled
		if task.resolved != nil {
			for _, diagnostic := range task.resolved.GetConfigFileParsingDiagnostics() {
				taskReporter.reportDiagnostic(diagnostic)
			}
		}
		if len(taskReporter.errors) > 0 {
			taskReporter.exitStatus = ExitStatusDiagnosticsPresent_OutputsSkipped
		}
	}
	for _, downstream := range task.downStream {
		downstream.status <- status
	}
	return taskReporter
}

func (s *solutionBuilder) handleStatusThatDoesntRequireBuild(config string, task *buildTask, status *upToDateStatus, taskReporter *taskReporter) *upToDateStatus {
	switch status.kind {
	case upToDateStatusTypeUpToDate:
		if s.opts.command.BuildOptions.Dry.IsTrue() {
			taskReporter.reportStatus(ast.NewCompilerDiagnostic(diagnostics.Project_0_is_up_to_date, config))
		}
		return status
	case upToDateStatusTypeUpstreamErrors:
		upstreamStatus := status.data.(*upstreamErrors)
		if s.opts.command.BuildOptions.Verbose.IsTrue() {
			taskReporter.reportStatus(ast.NewCompilerDiagnostic(
				core.IfElse(
					upstreamStatus.refHasUpstreamErrors,
					diagnostics.Skipping_build_of_project_0_because_its_dependency_1_was_not_built,
					diagnostics.Skipping_build_of_project_0_because_its_dependency_1_has_errors,
				),
				s.relativeFileName(config),
				s.relativeFileName(upstreamStatus.ref),
			))
		}
		return status
	case upToDateStatusTypeSolution:
		return status
	case upToDateStatusTypeConfigFileNotFound:
		taskReporter.reportDiagnostic(ast.NewCompilerDiagnostic(diagnostics.File_0_not_found, config))
		return status
	}

	// update timestamps
	if status.IsPseudoBuild() {
		if s.opts.command.BuildOptions.Dry.IsTrue() {
			taskReporter.reportStatus(ast.NewCompilerDiagnostic(diagnostics.A_non_dry_build_would_update_timestamps_for_output_of_project_0, config))
			status = &upToDateStatus{kind: upToDateStatusTypeUpToDate}
			return status
		}

		s.updateTimeStamps(config, task, taskReporter, nil, diagnostics.Updating_output_timestamps_of_project_0)
		status = &upToDateStatus{kind: upToDateStatusTypeUpToDate}
		return status
	}

	if s.opts.command.BuildOptions.Dry.IsTrue() {
		taskReporter.reportStatus(ast.NewCompilerDiagnostic(diagnostics.A_non_dry_build_would_build_project_0, config))
		status = &upToDateStatus{kind: upToDateStatusTypeUpToDate}
		return status
	}
	return nil
}

func (s *solutionBuilder) getUpToDateStatus(config string, configPath tspath.Path, task *buildTask, upStreamStatus []*upToDateStatus) *upToDateStatus {
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

		if s.opts.command.BuildOptions.StopBuildOnErrors.IsTrue() && upstreamStatus.IsError() {
			// Upstream project has errors, so we cannot build this project
			return &upToDateStatus{kind: upToDateStatusTypeUpstreamErrors, data: &upstreamErrors{task.resolved.ProjectReferences()[index].Path, upstreamStatus.kind == upToDateStatusTypeUpstreamErrors}}
		}
	}

	if s.opts.command.BuildOptions.Force.IsTrue() {
		return &upToDateStatus{kind: upToDateStatusTypeForceBuild}
	}

	// Check the build info
	buildInfoPath := task.resolved.GetBuildInfoFileName()
	buildInfo := s.host.readOrStoreBuildInfo(configPath, buildInfoPath)
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
		if buildInfo.IsEmitPending(task.resolved, tspath.GetDirectoryPath(tspath.GetNormalizedAbsolutePath(buildInfoPath, s.comparePathsOptions.CurrentDirectory))) {
			return &upToDateStatus{kind: upToDateStatusTypeOutOfDateOptions, data: buildInfoPath}
		}
	}
	var inputTextUnchanged bool
	oldestOutputFileAndTime := fileAndTime{buildInfoPath, s.host.GetMTime(buildInfoPath)}
	var newestInputFileAndTime fileAndTime
	var seenRoots collections.Set[tspath.Path]
	var buildInfoRootInfoReader *incremental.BuildInfoRootInfoReader
	for _, inputFile := range task.resolved.FileNames() {
		inputTime := s.host.GetMTime(inputFile)
		if inputTime.IsZero() {
			return &upToDateStatus{kind: upToDateStatusTypeInputFileMissing, data: inputFile}
		}
		inputPath := s.toPath(inputFile)
		if inputTime.After(oldestOutputFileAndTime.time) {
			var version string
			var currentVersion string
			if buildInfo.IsIncremental() {
				if buildInfoRootInfoReader == nil {
					buildInfoRootInfoReader = buildInfo.GetBuildInfoRootInfoReader(tspath.GetDirectoryPath(tspath.GetNormalizedAbsolutePath(buildInfoPath, s.comparePathsOptions.CurrentDirectory)), s.comparePathsOptions)
				}
				buildInfoFileInfo, resolvedInputPath := buildInfoRootInfoReader.GetBuildInfoFileInfo(inputPath)
				if fileInfo := buildInfoFileInfo.GetFileInfo(); fileInfo != nil && fileInfo.Version() != "" {
					version = fileInfo.Version()
					if text, ok := s.host.FS().ReadFile(string(resolvedInputPath)); ok {
						currentVersion = incremental.ComputeHash(text, s.opts.testing != nil)
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
		buildInfoRootInfoReader = buildInfo.GetBuildInfoRootInfoReader(tspath.GetDirectoryPath(tspath.GetNormalizedAbsolutePath(buildInfoPath, s.comparePathsOptions.CurrentDirectory)), s.comparePathsOptions)
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
			outputTime := s.host.GetMTime(outputFile)
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
		refInputOutputFileAndTime := upstreamStatus.InputOutputFileAndTime()
		if refInputOutputFileAndTime != nil && !refInputOutputFileAndTime.input.time.IsZero() && refInputOutputFileAndTime.input.time.Before(oldestOutputFileAndTime.time) {
			continue
		}

		// Check if tsbuildinfo path is shared, then we need to rebuild
		if s.host.hasConflictingBuildInfo(configPath) {
			return &upToDateStatus{kind: upToDateStatusTypeInputFileNewer, data: &inputOutputName{task.resolved.ProjectReferences()[index].Path, oldestOutputFileAndTime.file}}
		}

		// If the upstream project has only change .d.ts files, and we've built
		// *after* those files, then we're "pseudo up to date" and eligible for a fast rebuild
		newestDtsChangeTime := s.host.getLatestChangedDtsMTime(task.resolved.ResolvedProjectReferencePaths()[index])
		if !newestDtsChangeTime.IsZero() && newestDtsChangeTime.Before(oldestOutputFileAndTime.time) {
			refDtsUnchanged = true
			continue
		}

		// We have an output older than an upstream output - we are out of date
		return &upToDateStatus{kind: upToDateStatusTypeInputFileNewer, data: &inputOutputName{task.resolved.ProjectReferences()[index].Path, oldestOutputFileAndTime.file}}
	}

	configStatus := s.checkInputFileTime(config, &oldestOutputFileAndTime)
	if configStatus != nil {
		return configStatus
	}

	for _, extendedConfig := range task.resolved.ExtendedSourceFiles() {
		extendedConfigStatus := s.checkInputFileTime(extendedConfig, &oldestOutputFileAndTime)
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

func (s *solutionBuilder) checkInputFileTime(inputFile string, oldestOutputFileAndTime *fileAndTime) *upToDateStatus {
	inputTime := s.host.GetMTime(inputFile)
	if inputTime.After(oldestOutputFileAndTime.time) {
		// Output file is older than input file
		return &upToDateStatus{kind: upToDateStatusTypeInputFileNewer, data: &inputOutputName{inputFile, oldestOutputFileAndTime.file}}
	}
	return nil
}

func (s *solutionBuilder) reportUpToDateStatus(config string, status *upToDateStatus, taskReporter *taskReporter) {
	if !s.opts.command.BuildOptions.Verbose.IsTrue() {
		return
	}
	switch status.kind {
	case upToDateStatusTypeConfigFileNotFound:
		taskReporter.reportStatus(ast.NewCompilerDiagnostic(
			diagnostics.Project_0_is_out_of_date_because_config_file_does_not_exist,
			s.relativeFileName(config),
		))
	case upToDateStatusTypeUpstreamErrors:
		upstreamStatus := status.data.(*upstreamErrors)
		taskReporter.reportStatus(ast.NewCompilerDiagnostic(
			core.IfElse(
				upstreamStatus.refHasUpstreamErrors,
				diagnostics.Project_0_can_t_be_built_because_its_dependency_1_was_not_built,
				diagnostics.Project_0_can_t_be_built_because_its_dependency_1_has_errors,
			),
			s.relativeFileName(config),
			s.relativeFileName(upstreamStatus.ref),
		))
	case upToDateStatusTypeUpToDate:
		inputOutputFileAndTime := status.data.(*inputOutputFileAndTime)
		taskReporter.reportStatus(ast.NewCompilerDiagnostic(
			diagnostics.Project_0_is_up_to_date_because_newest_input_1_is_older_than_output_2,
			s.relativeFileName(config),
			s.relativeFileName(inputOutputFileAndTime.input.file),
			s.relativeFileName(inputOutputFileAndTime.output.file),
		))
	case upToDateStatusTypeUpToDateWithUpstreamTypes:
		taskReporter.reportStatus(ast.NewCompilerDiagnostic(
			diagnostics.Project_0_is_up_to_date_with_d_ts_files_from_its_dependencies,
			s.relativeFileName(config),
		))
	case upToDateStatusTypeUpToDateWithInputFileText:
		taskReporter.reportStatus(ast.NewCompilerDiagnostic(
			diagnostics.Project_0_is_up_to_date_but_needs_to_update_timestamps_of_output_files_that_are_older_than_input_files,
			s.relativeFileName(config),
		))
	case upToDateStatusTypeInputFileMissing:
		taskReporter.reportStatus(ast.NewCompilerDiagnostic(
			diagnostics.Project_0_is_out_of_date_because_input_1_does_not_exist,
			s.relativeFileName(config),
			s.relativeFileName(status.data.(string)),
		))
	case upToDateStatusTypeOutputMissing:
		taskReporter.reportStatus(ast.NewCompilerDiagnostic(
			diagnostics.Project_0_is_out_of_date_because_output_file_1_does_not_exist,
			s.relativeFileName(config),
			s.relativeFileName(status.data.(string)),
		))
	case upToDateStatusTypeInputFileNewer:
		inputOutput := status.data.(*inputOutputName)
		taskReporter.reportStatus(ast.NewCompilerDiagnostic(
			diagnostics.Project_0_is_out_of_date_because_output_1_is_older_than_input_2,
			s.relativeFileName(config),
			s.relativeFileName(inputOutput.output),
			s.relativeFileName(inputOutput.input),
		))
	case upToDateStatusTypeOutOfDateBuildInfoWithPendingEmit:
		taskReporter.reportStatus(ast.NewCompilerDiagnostic(
			diagnostics.Project_0_is_out_of_date_because_buildinfo_file_1_indicates_that_some_of_the_changes_were_not_emitted,
			s.relativeFileName(config),
			s.relativeFileName(status.data.(string)),
		))
	case upToDateStatusTypeOutOfDateBuildInfoWithErrors:
		taskReporter.reportStatus(ast.NewCompilerDiagnostic(
			diagnostics.Project_0_is_out_of_date_because_buildinfo_file_1_indicates_that_program_needs_to_report_errors,
			s.relativeFileName(config),
			s.relativeFileName(status.data.(string)),
		))
	case upToDateStatusTypeOutOfDateOptions:
		taskReporter.reportStatus(ast.NewCompilerDiagnostic(
			diagnostics.Project_0_is_out_of_date_because_buildinfo_file_1_indicates_there_is_change_in_compilerOptions,
			s.relativeFileName(config),
			s.relativeFileName(status.data.(string)),
		))
	case upToDateStatusTypeOutOfDateRoots:
		inputOutput := status.data.(*inputOutputName)
		taskReporter.reportStatus(ast.NewCompilerDiagnostic(
			diagnostics.Project_0_is_out_of_date_because_buildinfo_file_1_indicates_that_file_2_was_root_file_of_compilation_but_not_any_more,
			s.relativeFileName(config),
			s.relativeFileName(inputOutput.output),
			s.relativeFileName(inputOutput.input),
		))
	case upToDateStatusTypeTsVersionOutputOfDate:
		taskReporter.reportStatus(ast.NewCompilerDiagnostic(
			diagnostics.Project_0_is_out_of_date_because_output_for_it_was_generated_with_version_1_that_differs_with_current_version_2,
			s.relativeFileName(config),
			s.relativeFileName(status.data.(string)),
			core.Version(),
		))
	case upToDateStatusTypeForceBuild:
		taskReporter.reportStatus(ast.NewCompilerDiagnostic(
			diagnostics.Project_0_is_being_forcibly_rebuilt,
			s.relativeFileName(config),
		))
	case upToDateStatusTypeSolution:
		// Does not need to report status
	default:
		panic(fmt.Sprintf("Unknown up to date status kind: %v", status.kind))
	}
}

func (s *solutionBuilder) updateTimeStamps(config string, task *buildTask, taskReporter *taskReporter, emittedFiles []string, verboseMessage *diagnostics.Message) {
	if task.resolved.CompilerOptions().NoEmit.IsTrue() {
		return
	}
	emitted := collections.NewSetFromItems(emittedFiles...)
	var verboseMessageReported bool
	updateTimeStamp := func(file string) {
		if emitted.Has(file) {
			return
		}
		if !verboseMessageReported && s.opts.command.BuildOptions.Verbose.IsTrue() {
			taskReporter.reportStatus(ast.NewCompilerDiagnostic(verboseMessage, s.relativeFileName(config)))
			verboseMessageReported = true
		}
		err := s.host.SetMTime(file, s.opts.sys.Now())
		if err != nil {
			taskReporter.reportDiagnostic(ast.NewCompilerDiagnostic(diagnostics.Failed_to_update_timestamp_of_file_0, file))
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

func (s *solutionBuilder) cleanProject(config string, path tspath.Path, task *buildTask) *taskReporter {
	taskReporter := s.createTaskReporter()
	if task.resolved == nil {
		taskReporter.reportDiagnostic(ast.NewCompilerDiagnostic(diagnostics.File_0_not_found, config))
		taskReporter.exitStatus = ExitStatusDiagnosticsPresent_OutputsSkipped
		return taskReporter
	}

	inputs := collections.NewSetFromItems(core.Map(task.resolved.FileNames(), s.toPath)...)
	for outputFile := range task.resolved.GetOutputFileNames() {
		s.cleanProjectOutput(outputFile, inputs, taskReporter)
	}
	s.cleanProjectOutput(task.resolved.GetBuildInfoFileName(), inputs, taskReporter)
	return taskReporter
}

func (s *solutionBuilder) cleanProjectOutput(outputFile string, inputs *collections.Set[tspath.Path], taskReporter *taskReporter) {
	outputPath := s.toPath(outputFile)
	// If output name is same as input file name, do not delete and ignore the error
	if inputs.Has(outputPath) {
		return
	}
	if s.host.FS().FileExists(outputFile) {
		if !s.opts.command.BuildOptions.Dry.IsTrue() {
			err := s.host.FS().Remove(outputFile)
			if err != nil {
				taskReporter.reportDiagnostic(ast.NewCompilerDiagnostic(diagnostics.Failed_to_delete_file_0, outputFile))
			}
		} else {
			taskReporter.filesToDelete = append(taskReporter.filesToDelete, outputFile)
		}
	}
}

func newSolutionBuilder(opts solutionBuilderOptions) *solutionBuilder {
	solutionBuilder := &solutionBuilder{
		opts: opts,
		comparePathsOptions: tspath.ComparePathsOptions{
			CurrentDirectory:          opts.sys.GetCurrentDirectory(),
			UseCaseSensitiveFileNames: opts.sys.FS().UseCaseSensitiveFileNames(),
		},
	}
	return solutionBuilder
}
