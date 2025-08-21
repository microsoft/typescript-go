package build

import (
	"fmt"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/diagnostics"
	"github.com/microsoft/typescript-go/internal/execute/tsc"
	"github.com/microsoft/typescript-go/internal/incremental"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type buildTask struct {
	config   string
	resolved *tsoptions.ParsedCommandLine
	upStream []*buildTask
	status   *upToDateStatus
	done     chan struct{}

	// task reporting
	builder            strings.Builder
	errors             []*ast.Diagnostic
	reportStatus       tsc.DiagnosticReporter
	diagnosticReporter tsc.DiagnosticReporter
	exitStatus         tsc.ExitStatus
	statistics         *tsc.Statistics
	program            *incremental.Program
	pseudoBuild        bool
	filesToDelete      []string
	prevReporter       *buildTask
	reportDone         chan struct{}
}

func (t *buildTask) waitOnUpstream() []*upToDateStatus {
	upStreamStatus := make([]*upToDateStatus, len(t.upStream))
	for i, upstream := range t.upStream {
		<-upstream.done
		upStreamStatus[i] = upstream.status
	}
	return upStreamStatus
}

func (t *buildTask) unblockDownstream(status *upToDateStatus) {
	t.status = status
	close(t.done)
}

func (t *buildTask) reportDiagnostic(err *ast.Diagnostic) {
	t.errors = append(t.errors, err)
	t.diagnosticReporter(err)
}

func (t *buildTask) report(orchestrator *Orchestrator, configPath tspath.Path, buildResult *solutionBuilderResult) {
	if t.prevReporter != nil {
		<-t.prevReporter.reportDone
	}
	if len(t.errors) > 0 {
		buildResult.errors = append(core.IfElse(buildResult.errors != nil, buildResult.errors, []*ast.Diagnostic{}), t.errors...)
	}
	fmt.Fprint(orchestrator.opts.Sys.Writer(), t.builder.String())
	if t.exitStatus > buildResult.result.Status {
		buildResult.result.Status = t.exitStatus
	}
	if t.statistics != nil {
		buildResult.programStats = append(buildResult.programStats, t.statistics)
	}
	if t.program != nil {
		buildResult.result.IncrementalProgram = append(buildResult.result.IncrementalProgram, t.program)
		buildResult.statistics.ProjectsBuilt++
	}
	if t.pseudoBuild {
		buildResult.statistics.TimestampUpdates++
	}
	buildResult.filesToDelete = append(buildResult.filesToDelete, t.filesToDelete...)
	close(t.reportDone)
}

func (t *buildTask) reportUpToDateStatus(orchestrator *Orchestrator, status *upToDateStatus) {
	if !orchestrator.opts.Command.BuildOptions.Verbose.IsTrue() {
		return
	}
	switch status.kind {
	case upToDateStatusTypeConfigFileNotFound:
		t.reportStatus(ast.NewCompilerDiagnostic(
			diagnostics.Project_0_is_out_of_date_because_config_file_does_not_exist,
			orchestrator.relativeFileName(t.config),
		))
	case upToDateStatusTypeUpstreamErrors:
		upstreamStatus := status.data.(*upstreamErrors)
		t.reportStatus(ast.NewCompilerDiagnostic(
			core.IfElse(
				upstreamStatus.refHasUpstreamErrors,
				diagnostics.Project_0_can_t_be_built_because_its_dependency_1_was_not_built,
				diagnostics.Project_0_can_t_be_built_because_its_dependency_1_has_errors,
			),
			orchestrator.relativeFileName(t.config),
			orchestrator.relativeFileName(upstreamStatus.ref),
		))
	case upToDateStatusTypeUpToDate:
		inputOutputFileAndTime := status.data.(*inputOutputFileAndTime)
		t.reportStatus(ast.NewCompilerDiagnostic(
			diagnostics.Project_0_is_up_to_date_because_newest_input_1_is_older_than_output_2,
			orchestrator.relativeFileName(t.config),
			orchestrator.relativeFileName(inputOutputFileAndTime.input.file),
			orchestrator.relativeFileName(inputOutputFileAndTime.output.file),
		))
	case upToDateStatusTypeUpToDateWithUpstreamTypes:
		t.reportStatus(ast.NewCompilerDiagnostic(
			diagnostics.Project_0_is_up_to_date_with_d_ts_files_from_its_dependencies,
			orchestrator.relativeFileName(t.config),
		))
	case upToDateStatusTypeUpToDateWithInputFileText:
		t.reportStatus(ast.NewCompilerDiagnostic(
			diagnostics.Project_0_is_up_to_date_but_needs_to_update_timestamps_of_output_files_that_are_older_than_input_files,
			orchestrator.relativeFileName(t.config),
		))
	case upToDateStatusTypeInputFileMissing:
		t.reportStatus(ast.NewCompilerDiagnostic(
			diagnostics.Project_0_is_out_of_date_because_input_1_does_not_exist,
			orchestrator.relativeFileName(t.config),
			orchestrator.relativeFileName(status.data.(string)),
		))
	case upToDateStatusTypeOutputMissing:
		t.reportStatus(ast.NewCompilerDiagnostic(
			diagnostics.Project_0_is_out_of_date_because_output_file_1_does_not_exist,
			orchestrator.relativeFileName(t.config),
			orchestrator.relativeFileName(status.data.(string)),
		))
	case upToDateStatusTypeInputFileNewer:
		inputOutput := status.data.(*inputOutputName)
		t.reportStatus(ast.NewCompilerDiagnostic(
			diagnostics.Project_0_is_out_of_date_because_output_1_is_older_than_input_2,
			orchestrator.relativeFileName(t.config),
			orchestrator.relativeFileName(inputOutput.output),
			orchestrator.relativeFileName(inputOutput.input),
		))
	case upToDateStatusTypeOutOfDateBuildInfoWithPendingEmit:
		t.reportStatus(ast.NewCompilerDiagnostic(
			diagnostics.Project_0_is_out_of_date_because_buildinfo_file_1_indicates_that_some_of_the_changes_were_not_emitted,
			orchestrator.relativeFileName(t.config),
			orchestrator.relativeFileName(status.data.(string)),
		))
	case upToDateStatusTypeOutOfDateBuildInfoWithErrors:
		t.reportStatus(ast.NewCompilerDiagnostic(
			diagnostics.Project_0_is_out_of_date_because_buildinfo_file_1_indicates_that_program_needs_to_report_errors,
			orchestrator.relativeFileName(t.config),
			orchestrator.relativeFileName(status.data.(string)),
		))
	case upToDateStatusTypeOutOfDateOptions:
		t.reportStatus(ast.NewCompilerDiagnostic(
			diagnostics.Project_0_is_out_of_date_because_buildinfo_file_1_indicates_there_is_change_in_compilerOptions,
			orchestrator.relativeFileName(t.config),
			orchestrator.relativeFileName(status.data.(string)),
		))
	case upToDateStatusTypeOutOfDateRoots:
		inputOutput := status.data.(*inputOutputName)
		t.reportStatus(ast.NewCompilerDiagnostic(
			diagnostics.Project_0_is_out_of_date_because_buildinfo_file_1_indicates_that_file_2_was_root_file_of_compilation_but_not_any_more,
			orchestrator.relativeFileName(t.config),
			orchestrator.relativeFileName(inputOutput.output),
			orchestrator.relativeFileName(inputOutput.input),
		))
	case upToDateStatusTypeTsVersionOutputOfDate:
		t.reportStatus(ast.NewCompilerDiagnostic(
			diagnostics.Project_0_is_out_of_date_because_output_for_it_was_generated_with_version_1_that_differs_with_current_version_2,
			orchestrator.relativeFileName(t.config),
			orchestrator.relativeFileName(status.data.(string)),
			core.Version(),
		))
	case upToDateStatusTypeForceBuild:
		t.reportStatus(ast.NewCompilerDiagnostic(
			diagnostics.Project_0_is_being_forcibly_rebuilt,
			orchestrator.relativeFileName(t.config),
		))
	case upToDateStatusTypeSolution:
		// Does not need to report status
	default:
		panic(fmt.Sprintf("Unknown up to date status kind: %v", status.kind))
	}
}
