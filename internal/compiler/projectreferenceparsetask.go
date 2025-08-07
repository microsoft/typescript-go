package compiler

import (
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type projectReferenceParseTask struct {
	loaded     bool
	configPath tspath.Path
	configName string
	resolved   *tsoptions.ParsedCommandLine
	subTasks   []*projectReferenceParseTask
}

func (t *projectReferenceParseTask) Path() tspath.Path {
	return t.configPath
}

func (t *projectReferenceParseTask) load(loader *fileLoader) {
	t.loaded = true

	t.resolved = loader.opts.Host.GetResolvedProjectReference(t.configName, t.configPath)
	if t.resolved == nil {
		return
	}
	if t.resolved.SourceToOutput() == nil {
		loader.projectReferenceParseTasks.wg.Queue(func() {
			t.resolved.ParseInputOutputNames()
		})
	}
	subReferences := t.resolved.ResolvedProjectReferencePaths()
	if len(subReferences) == 0 {
		return
	}
	t.subTasks = createProjectReferenceParseTasks(subReferences, loader)
}

func (t *projectReferenceParseTask) getSubTasks() []*projectReferenceParseTask {
	return t.subTasks
}

func (t *projectReferenceParseTask) shouldIncreaseDepth() bool {
	return false
}

func (t *projectReferenceParseTask) shouldElideOnDepth() bool {
	return false
}

func (t *projectReferenceParseTask) isLoaded() bool {
	return t.loaded
}

func (t *projectReferenceParseTask) isRoot() bool {
	return true
}

func (t *projectReferenceParseTask) isFromExternalLibrary() bool {
	return false
}

func (t *projectReferenceParseTask) markFromExternalLibrary() {}

func createProjectReferenceParseTasks(projectReferences []string, loader *fileLoader) []*projectReferenceParseTask {
	return core.Map(projectReferences, func(configName string) *projectReferenceParseTask {
		return &projectReferenceParseTask{
			configName: configName,
			configPath: loader.toPath(configName),
		}
	})
}
