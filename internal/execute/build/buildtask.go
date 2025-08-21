package build

import (
	"fmt"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
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

func (t *buildTask) report(s *Orchestrator, configPath tspath.Path, buildResult *solutionBuilderResult) {
	if t.prevReporter != nil {
		<-t.prevReporter.reportDone
	}
	if len(t.errors) > 0 {
		buildResult.errors = append(core.IfElse(buildResult.errors != nil, buildResult.errors, []*ast.Diagnostic{}), t.errors...)
	}
	fmt.Fprint(s.opts.Sys.Writer(), t.builder.String())
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
