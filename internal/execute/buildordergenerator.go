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

type upToDateStatusType uint16

const (
	// building current project
	upToDateStatusTypeUnknown upToDateStatusType = iota
	// config file was not found
	upToDateStatusTypeConfigFileNotFound
	// upToDateStatusTypeUnbuildable
	// upToDateStatusTypeUpToDate
	// // The project appears out of date because its upstream inputs are newer than its outputs,
	// // but all of its outputs are actually newer than the previous identical outputs of its (.d.ts) inputs.
	// // This means we can Pseudo-build (just touch timestamps), as if we had actually built this project.
	// upToDateStatusTypeUpToDateWithUpstreamTypes
	// upToDateStatusTypeOutputMissing
	// upToDateStatusTypeErrorReadingFile
	// upToDateStatusTypeOutOfDateWithSelf
	// upToDateStatusTypeOutOfDateWithUpstream
	// upToDateStatusTypeOutOfDateBuildInfoWithPendingEmit
	// upToDateStatusTypeOutOfDateBuildInfoWithErrors
	// upToDateStatusTypeOutOfDateOptions
	// upToDateStatusTypeOutOfDateRoots
	// upToDateStatusTypeUpstreamOutOfDate
	// upToDateStatusTypeUpstreamBlocked
	// upToDateStatusTypeTsVersionOutputOfDate
	// upToDateStatusTypeUpToDateWithInputFileText
	// // solution file
	upToDateStatusTypeSolution
	// upToDateStatusTypeForceBuild
)

type upToDateStatus struct {
	kind upToDateStatusType
}

type statusTask struct {
	config       string
	referencedBy string
	status       chan *upToDateStatus
}

type buildTask struct {
	config     string
	resolved   *tsoptions.ParsedCommandLine
	upStream   []*statusTask
	downStream []*statusTask
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
		b.order = append(b.order, configName)
	}
	if downStream != nil {
		task.downStream = append(task.downStream, downStream)
		downStream.status = make(chan *upToDateStatus, 1)
	}
}
