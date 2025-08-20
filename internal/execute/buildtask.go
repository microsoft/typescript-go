package execute

import (
	"fmt"
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/incremental"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type statusTask struct {
	config       string
	referencedBy string
	status       chan *upToDateStatus
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
		buildResult.statistics.projectsBuilt++
	}
	buildResult.filesToDelete = append(buildResult.filesToDelete, b.filesToDelete...)
	close(b.done)
}

type buildTask struct {
	config       string
	resolved     *tsoptions.ParsedCommandLine
	upStream     []*statusTask
	downStream   []*statusTask
	taskReporter taskReporter
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
