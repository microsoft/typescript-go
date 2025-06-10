package compiler

import (
	"strings"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/module"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type ProjectReferenceFileMapper struct {
	opts   *ProgramOptions
	loader *fileLoader // Only present during populating the mapper and parsing, released after that

	configToProjectReference map[tspath.Path]*tsoptions.ParsedCommandLine // All the resolved references needed
	referencesInConfigFile   map[tspath.Path][]tspath.Path                // Map of config file to its references
	sourceToOutput           map[tspath.Path]*tsoptions.OutputDtsAndProjectReference
	outputDtsToSource        map[tspath.Path]*tsoptions.SourceAndProjectReference

	// Store all the realpath from dts in node_modules to source file from project reference needed during parsing so it can be used later
	realpathDtsToSource collections.SyncMap[tspath.Path, *tsoptions.SourceAndProjectReference]
}

func (mapper *ProjectReferenceFileMapper) init(loader *fileLoader, rootTasks []*projectReferenceParseTask) {
	totalReferences := loader.projectReferenceParseTasks.tasksByFileName.Size() + 1
	mapper.loader = loader
	mapper.configToProjectReference = make(map[tspath.Path]*tsoptions.ParsedCommandLine, totalReferences)
	mapper.referencesInConfigFile = make(map[tspath.Path][]tspath.Path, totalReferences)
	mapper.sourceToOutput = make(map[tspath.Path]*tsoptions.OutputDtsAndProjectReference)
	mapper.outputDtsToSource = make(map[tspath.Path]*tsoptions.SourceAndProjectReference)
	mapper.referencesInConfigFile[mapper.opts.Config.ConfigFile.SourceFile.Path()] = loader.projectReferenceParseTasks.collect(
		loader,
		rootTasks,
		func(task *projectReferenceParseTask, referencesInConfig []tspath.Path) {
			path := loader.toPath(task.configName)
			mapper.configToProjectReference[path] = task.resolved
			if task.resolved == nil || mapper.opts.Config.ConfigFile == task.resolved.ConfigFile {
				return
			}
			mapper.referencesInConfigFile[path] = referencesInConfig
			for key, value := range task.resolved.SourceToOutput() {
				mapper.sourceToOutput[key] = value
			}
			for key, value := range task.resolved.OutputDtsToSource() {
				mapper.outputDtsToSource[key] = value
			}
			if mapper.opts.canUseProjectReferenceSource() {
				declDir := task.resolved.CompilerOptions().DeclarationDir
				if declDir == "" {
					declDir = task.resolved.CompilerOptions().OutDir
				}
				if declDir != "" {
					loader.dtsDirectories.Add(loader.toPath(declDir))
				}
			}
		})
}

func (mapper *ProjectReferenceFileMapper) getParseFileRedirect(file ast.HasFileName) string {
	if mapper.opts.canUseProjectReferenceSource() {
		// Map to source file from project reference
		source := mapper.getSourceAndProjectReference(file.Path())
		if source == nil {
			source = mapper.getSourceToDtsIfSymlink(file)
		}
		if source != nil {
			return source.Source
		}
	} else {
		// Map to dts file from project reference
		output := mapper.getOutputAndProjectReference(file.Path())
		if output != nil && output.OutputDts != "" {
			return output.OutputDts
		}
	}
	return ""
}

func (mapper *ProjectReferenceFileMapper) getResolvedProjectReferences() []*tsoptions.ParsedCommandLine {
	refs, ok := mapper.referencesInConfigFile[mapper.opts.Config.ConfigFile.SourceFile.Path()]
	var result []*tsoptions.ParsedCommandLine
	if ok {
		result = make([]*tsoptions.ParsedCommandLine, 0, len(refs))
		for _, refPath := range refs {
			refConfig, _ := mapper.configToProjectReference[refPath]
			result = append(result, refConfig)
		}
	}
	return result
}

func (mapper *ProjectReferenceFileMapper) getOutputAndProjectReference(path tspath.Path) *tsoptions.OutputDtsAndProjectReference {
	return mapper.sourceToOutput[path]
}

func (mapper *ProjectReferenceFileMapper) getSourceAndProjectReference(path tspath.Path) *tsoptions.SourceAndProjectReference {
	return mapper.outputDtsToSource[path]
}

func (mapper *ProjectReferenceFileMapper) isSourceFromProjectReference(path tspath.Path) bool {
	return mapper.opts.canUseProjectReferenceSource() && mapper.getOutputAndProjectReference(path) != nil
}

func (mapper *ProjectReferenceFileMapper) getCompilerOptionsForFile(file ast.HasFileName) *core.CompilerOptions {
	redirect := mapper.getRedirectForResolution(file)
	return module.GetCompilerOptionsWithRedirect(mapper.opts.Config.CompilerOptions(), redirect)
}

func (mapper *ProjectReferenceFileMapper) getRedirectForResolution(file ast.HasFileName) *tsoptions.ParsedCommandLine {
	path := file.Path()
	// Check if outputdts of source file from project reference
	output := mapper.getOutputAndProjectReference(path)
	if output != nil {
		return output.Resolved
	}

	// Source file from project reference
	resultFromDts := mapper.getSourceAndProjectReference(path)
	if resultFromDts != nil {
		return resultFromDts.Resolved
	}

	realpathDtsToSource := mapper.getSourceToDtsIfSymlink(file)
	if realpathDtsToSource != nil {
		return realpathDtsToSource.Resolved
	}
	return nil
}

func (mapper *ProjectReferenceFileMapper) getResolvedReferenceFor(path tspath.Path) (*tsoptions.ParsedCommandLine, bool) {
	config, ok := mapper.configToProjectReference[path]
	return config, ok
}

func (mapper *ProjectReferenceFileMapper) forEachResolvedProjectReference(
	fn func(path tspath.Path, config *tsoptions.ParsedCommandLine) bool,
) {
	if mapper.opts.Config.ConfigFile == nil {
		return
	}
	refs := mapper.referencesInConfigFile[mapper.opts.Config.ConfigFile.SourceFile.Path()]
	mapper.forEachResolvedReferenceWorker(refs, fn)
}

func (mapper *ProjectReferenceFileMapper) forEachResolvedReferenceWorker(
	referenes []tspath.Path,
	fn func(path tspath.Path, config *tsoptions.ParsedCommandLine) bool,
) {
	for _, path := range referenes {
		config, _ := mapper.configToProjectReference[path]
		if !fn(path, config) {
			return
		}
	}
}

func (mapper *ProjectReferenceFileMapper) getSourceToDtsIfSymlink(file ast.HasFileName) *tsoptions.SourceAndProjectReference {
	// If preserveSymlinks is true, module resolution wont jump the symlink
	// but the resolved real path may be the .d.ts from project reference
	// Note:: Currently we try the real path only if the
	// file is from node_modules to avoid having to run real path on all file paths
	path := file.Path()
	realpathDtsToSource, ok := mapper.realpathDtsToSource.Load(path)
	if ok {
		return realpathDtsToSource
	}
	if mapper.loader != nil && mapper.opts.Config.CompilerOptions().PreserveSymlinks == core.TSTrue {
		fileName := file.FileName()
		if !strings.Contains(fileName, "/node_modules/") {
			mapper.realpathDtsToSource.Store(path, nil)
		} else {
			realDeclarationPath := mapper.loader.toPath(mapper.loader.resolver.GetHost().FS().Realpath(fileName))
			if realDeclarationPath == path {
				mapper.realpathDtsToSource.Store(path, nil)
			} else {
				realpathDtsToSource := mapper.getSourceAndProjectReference(realDeclarationPath)
				if realpathDtsToSource != nil {
					mapper.realpathDtsToSource.Store(path, realpathDtsToSource)
					return realpathDtsToSource
				}
				mapper.realpathDtsToSource.Store(path, nil)
			}
		}
	}
	return nil
}
