package execute

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

type taskReporter struct {
	builder            strings.Builder
	errors             []*ast.Diagnostic
	reportStatus       tsc.DiagnosticReporter
	diagnosticReporter tsc.DiagnosticReporter
	exitStatus         tsc.ExitStatus
	statistics         *tsc.Statistics
	program            *incremental.Program
	pseudoBuild        bool
	filesToDelete      []string
	prev               *taskReporter
	done               chan struct{}
}

func (b *taskReporter) reportDiagnostic(err *ast.Diagnostic) {
	b.errors = append(b.errors, err)
	b.diagnosticReporter(err)
}

func (b *taskReporter) report(s *solutionBuilder, configPath tspath.Path, buildResult *solutionBuilderResult) {
	if b.prev != nil {
		<-b.prev.done
	}
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
		buildResult.statistics.ProjectsBuilt++
	}
	if b.pseudoBuild {
		buildResult.statistics.TimestampUpdates++
	}
	buildResult.filesToDelete = append(buildResult.filesToDelete, b.filesToDelete...)
	close(b.done)
}

type buildTask struct {
	config       string
	resolved     *tsoptions.ParsedCommandLine
	upStream     []*buildTask
	status       *upToDateStatus
	done         chan struct{}
	taskReporter taskReporter
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
