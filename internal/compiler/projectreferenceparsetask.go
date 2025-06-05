package compiler

import (
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type projectReferenceParseTask struct {
	configName string
	path       tspath.Path
	resolved   *tsoptions.ParsedCommandLine
	subTasks   []*projectReferenceParseTask
}

func (t *projectReferenceParseTask) FileName() string {
	return t.configName
}

func (t *projectReferenceParseTask) start(loader *fileLoader) {
	t.path = loader.toPath(t.configName)
	t.resolved = loader.opts.Host.GetResolvedProjectReference(t.configName, t.path)
	if t.resolved == nil {
		return
	}
	if t.resolved.SourceToOutput() == nil {
		loader.projectReferenceParseTasks.wg.Queue(func() {
			t.resolved.ParseInputOutputNames()
		})
	}
	subReferences := t.resolved.ProjectReferences()
	if len(subReferences) == 0 {
		return
	}
	t.subTasks = createProjectReferenceParseTasks(subReferences)
}

func getSubTasksOfProjectReferenceParseTask(t *projectReferenceParseTask) []*projectReferenceParseTask {
	return t.subTasks
}

func createProjectReferenceParseTasks(projectReferences []*core.ProjectReference) []*projectReferenceParseTask {
	tasks := make([]*projectReferenceParseTask, 0, len(projectReferences))
	for _, reference := range projectReferences {
		configName := core.ResolveProjectReferencePath(reference)
		tasks = append(tasks, &projectReferenceParseTask{
			configName: configName,
		})
	}
	return tasks
}
