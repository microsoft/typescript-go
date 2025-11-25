package compiler

import (
	"math"
	"slices"
	"sync"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/module"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type parseTask struct {
	normalizedFilePath          string
	path                        tspath.Path
	file                        *ast.SourceFile
	libFile                     *LibFile
	redirectedParseTask         *parseTask
	subTasks                    []*parseTask
	loaded                      bool
	isForAutomaticTypeDirective bool
	includeReason               *FileIncludeReason

	metadata                     ast.SourceFileMetaData
	resolutionsInFile            module.ModeAwareCache[*module.ResolvedModule]
	resolutionsTrace             []module.DiagAndArgs
	typeResolutionsInFile        module.ModeAwareCache[*module.ResolvedTypeReferenceDirective]
	typeResolutionsTrace         []module.DiagAndArgs
	resolutionDiagnostics        []*ast.Diagnostic
	importHelpersImportSpecifier *ast.Node
	jsxRuntimeImportSpecifier    *jsxRuntimeImportSpecifier
	increaseDepth                bool
	elideOnDepth                 bool

	// Track if this file is from an external library (node_modules)
	// This mirrors the TypeScript currentNodeModulesDepth > 0 check
	fromExternalLibrary bool

	loadedTask        *parseTask
	allIncludeReasons []*FileIncludeReason
}

func (t *parseTask) FileName() string {
	return t.normalizedFilePath
}

func (t *parseTask) Path() tspath.Path {
	return t.path
}

func (t *parseTask) isRoot() bool {
	// Intentionally not checking t.includeReason != nil to ensure we can catch cases for missing include reason
	return !t.isForAutomaticTypeDirective && (t.includeReason.kind == fileIncludeKindRootFile || t.includeReason.kind == fileIncludeKindLibFile)
}

func (t *parseTask) load(loader *fileLoader) {
	t.loaded = true
	if t.isForAutomaticTypeDirective {
		t.loadAutomaticTypeDirectives(loader)
		return
	}
	redirect := loader.projectReferenceFileMapper.getParseFileRedirect(t)
	if redirect != "" {
		t.redirect(loader, redirect)
		return
	}

	loader.totalFileCount.Add(1)
	if t.libFile != nil {
		loader.libFileCount.Add(1)
	}

	t.metadata = loader.loadSourceFileMetaData(t.normalizedFilePath)
	file := loader.parseSourceFile(t)
	if file == nil {
		return
	}

	t.file = file
	t.subTasks = make([]*parseTask, 0, len(file.ReferencedFiles)+len(file.Imports())+len(file.ModuleAugmentations))

	for index, ref := range file.ReferencedFiles {
		resolvedPath := loader.resolveTripleslashPathReference(ref.FileName, file.FileName(), index)
		t.addSubTask(resolvedPath, nil)
	}

	compilerOptions := loader.opts.Config.CompilerOptions()
	loader.resolveTypeReferenceDirectives(t)

	if compilerOptions.NoLib != core.TSTrue {
		for index, lib := range file.LibReferenceDirectives {
			includeReason := &FileIncludeReason{
				kind: fileIncludeKindLibReferenceDirective,
				data: &referencedFileData{
					file:  t.path,
					index: index,
				},
			}
			if name, ok := tsoptions.GetLibFileName(lib.FileName); ok {
				libFile := loader.pathForLibFile(name)
				t.addSubTask(resolvedRef{
					fileName:      libFile.path,
					includeReason: includeReason,
				}, libFile)
			} else {
				loader.includeProcessor.addProcessingDiagnostic(&processingDiagnostic{
					kind: processingDiagnosticKindUnknownReference,
					data: includeReason,
				})
			}
		}
	}

	loader.resolveImportsAndModuleAugmentations(t)
}

func (t *parseTask) redirect(loader *fileLoader, fileName string) {
	t.redirectedParseTask = &parseTask{
		normalizedFilePath:  tspath.NormalizePath(fileName),
		libFile:             t.libFile,
		fromExternalLibrary: t.fromExternalLibrary,
		includeReason:       t.includeReason,
	}
	// increaseDepth and elideOnDepth are not copied to redirects, otherwise their depth would be double counted.
	t.subTasks = []*parseTask{t.redirectedParseTask}
}

func (t *parseTask) loadAutomaticTypeDirectives(loader *fileLoader) {
	toParseTypeRefs, typeResolutionsInFile, typeResolutionsTrace := loader.resolveAutomaticTypeDirectives(t.normalizedFilePath)
	t.typeResolutionsInFile = typeResolutionsInFile
	t.typeResolutionsTrace = typeResolutionsTrace
	for _, typeResolution := range toParseTypeRefs {
		t.addSubTask(typeResolution, nil)
	}
}

type resolvedRef struct {
	fileName              string
	increaseDepth         bool
	elideOnDepth          bool
	isFromExternalLibrary bool
	includeReason         *FileIncludeReason
}

func (t *parseTask) addSubTask(ref resolvedRef, libFile *LibFile) {
	normalizedFilePath := tspath.NormalizePath(ref.fileName)
	subTask := &parseTask{
		normalizedFilePath:  normalizedFilePath,
		libFile:             libFile,
		increaseDepth:       ref.increaseDepth,
		elideOnDepth:        ref.elideOnDepth,
		fromExternalLibrary: ref.isFromExternalLibrary,
		includeReason:       ref.includeReason,
	}
	t.subTasks = append(t.subTasks, subTask)
}

type filesParser struct {
	wg              core.WorkGroup
	tasksByFileName collections.SyncMap[string, *parseTaskData]
	maxDepth        int
}

type parseTaskData struct {
	task                *parseTask
	mu                  sync.Mutex
	isRoot              bool
	lowestDepth         int
	fromExternalLibrary bool
}

func (w *filesParser) parse(loader *fileLoader, tasks []*parseTask) {
	w.start(loader, tasks, 0, false)
	w.wg.RunAndWait()
}

func (w *filesParser) start(loader *fileLoader, tasks []*parseTask, depth int, isFromExternalLibrary bool) {
	for i, task := range tasks {
		task.path = loader.toPath(task.normalizedFilePath)
		taskIsFromExternalLibrary := isFromExternalLibrary || task.fromExternalLibrary
		data, loaded := w.tasksByFileName.LoadOrStore(task.FileName(), &parseTaskData{
			task:        task,
			isRoot:      task.isRoot(),
			lowestDepth: math.MaxInt,
		})
		// task = data.task
		if loaded {
			tasks[i].loadedTask = data.task
			// Add in the loaded task's external-ness.
			taskIsFromExternalLibrary = taskIsFromExternalLibrary || data.fromExternalLibrary
		}

		w.wg.Queue(func() {
			data.mu.Lock()
			defer data.mu.Unlock()

			startSubtasks := false

			currentDepth := core.IfElse(task.increaseDepth, depth+1, depth)
			if currentDepth < data.lowestDepth {
				// If we're seeing this task at a lower depth than before,
				// reprocess its subtasks to ensure they are loaded.
				data.lowestDepth = currentDepth
				startSubtasks = true
			}

			if !data.isRoot && taskIsFromExternalLibrary && !data.fromExternalLibrary {
				// If we're seeing this task now as an external library,
				// reprocess its subtasks to ensure they are also marked as external.
				data.fromExternalLibrary = true
				startSubtasks = true
			}

			if task.elideOnDepth && currentDepth > w.maxDepth {
				return
			}

			if !data.task.loaded {
				data.task.load(loader)
			}

			if startSubtasks {
				w.start(loader, data.task.subTasks, data.lowestDepth, data.fromExternalLibrary)
			}
		})
	}
}

func (w *filesParser) getProcessedFiles(loader *fileLoader) processedFiles {
	totalFileCount := int(loader.totalFileCount.Load())
	libFileCount := int(loader.libFileCount.Load())

	var missingFiles []string
	files := make([]*ast.SourceFile, 0, totalFileCount-libFileCount)
	libFiles := make([]*ast.SourceFile, 0, totalFileCount) // totalFileCount here since we append files to it later to construct the final list

	filesByPath := make(map[tspath.Path]*ast.SourceFile, totalFileCount)
	loader.includeProcessor.fileIncludeReasons = make(map[tspath.Path][]*FileIncludeReason, totalFileCount)
	var outputFileToProjectReferenceSource map[tspath.Path]string
	if !loader.opts.canUseProjectReferenceSource() {
		outputFileToProjectReferenceSource = make(map[tspath.Path]string, totalFileCount)
	}
	resolvedModules := make(map[tspath.Path]module.ModeAwareCache[*module.ResolvedModule], totalFileCount+1)
	typeResolutionsInFile := make(map[tspath.Path]module.ModeAwareCache[*module.ResolvedTypeReferenceDirective], totalFileCount)
	sourceFileMetaDatas := make(map[tspath.Path]ast.SourceFileMetaData, totalFileCount)
	var jsxRuntimeImportSpecifiers map[tspath.Path]*jsxRuntimeImportSpecifier
	var importHelpersImportSpecifiers map[tspath.Path]*ast.Node
	var sourceFilesFoundSearchingNodeModules collections.Set[tspath.Path]
	libFilesMap := make(map[tspath.Path]*LibFile, libFileCount)

	var collectFiles func(tasks []*parseTask, seen collections.Set[*parseTaskData])
	collectFiles = func(tasks []*parseTask, seen collections.Set[*parseTaskData]) {
		for _, task := range tasks {
			// Exclude automatic type directive tasks from include reason processing,
			// as these are internal implementation details and should not contribute
			// to the reasons for including files.
			if task.redirectedParseTask == nil && !task.isForAutomaticTypeDirective {
				includeReason := task.includeReason
				if task.loadedTask != nil {
					task = task.loadedTask
				}
				w.addIncludeReason(loader, task, includeReason)
			}
			data, _ := w.tasksByFileName.Load(task.normalizedFilePath)
			// ensure we only walk each task once
			if !task.loaded || !seen.AddIfAbsent(data) {
				continue
			}
			for _, trace := range task.typeResolutionsTrace {
				loader.opts.Host.Trace(trace.Message, trace.Args...)
			}
			for _, trace := range task.resolutionsTrace {
				loader.opts.Host.Trace(trace.Message, trace.Args...)
			}
			if subTasks := task.subTasks; len(subTasks) > 0 {
				collectFiles(subTasks, seen)
			}

			// Exclude automatic type directive tasks from include reason processing,
			// as these are internal implementation details and should not contribute
			// to the reasons for including files.
			if task.redirectedParseTask != nil {
				if !loader.opts.canUseProjectReferenceSource() {
					outputFileToProjectReferenceSource[task.redirectedParseTask.path] = task.FileName()
				}
				continue
			}

			if task.isForAutomaticTypeDirective {
				typeResolutionsInFile[task.path] = task.typeResolutionsInFile
				continue
			}
			file := task.file
			path := task.path
			if file == nil {
				// !!! sheetal file preprocessing diagnostic explaining getSourceFileFromReferenceWorker
				missingFiles = append(missingFiles, task.normalizedFilePath)
				continue
			}

			// !!! sheetal todo porting file case errors
			// if _, ok := filesByPath[path]; ok {
			// 	Check if it differs only in drive letters its ok to ignore that error:
			// 	const checkedAbsolutePath = getNormalizedAbsolutePathWithoutRoot(checkedName, currentDirectory);
			// 	const inputAbsolutePath = getNormalizedAbsolutePathWithoutRoot(fileName, currentDirectory);
			// 	if (checkedAbsolutePath !== inputAbsolutePath) {
			// 	    reportFileNamesDifferOnlyInCasingError(fileName, file, reason);
			// 	}
			// } else if loader.comparePathsOptions.UseCaseSensitiveFileNames {
			// 	pathIgnoreCase := tspath.ToPath(file.FileName(), loader.comparePathsOptions.CurrentDirectory, false)
			// 	// for case-sensitsive file systems check if we've already seen some file with similar filename ignoring case
			// 	if _, ok := filesByNameIgnoreCase[pathIgnoreCase]; ok {
			// 		reportFileNamesDifferOnlyInCasingError(fileName, existingFile, reason);
			// 	} else {
			// 		filesByNameIgnoreCase[pathIgnoreCase] = file
			// 	}
			// }

			if task.libFile != nil {
				libFiles = append(libFiles, file)
				libFilesMap[path] = task.libFile
			} else {
				files = append(files, file)
			}
			filesByPath[path] = file
			resolvedModules[path] = task.resolutionsInFile
			typeResolutionsInFile[path] = task.typeResolutionsInFile
			sourceFileMetaDatas[path] = task.metadata

			if task.jsxRuntimeImportSpecifier != nil {
				if jsxRuntimeImportSpecifiers == nil {
					jsxRuntimeImportSpecifiers = make(map[tspath.Path]*jsxRuntimeImportSpecifier, totalFileCount)
				}
				jsxRuntimeImportSpecifiers[path] = task.jsxRuntimeImportSpecifier
			}
			if task.importHelpersImportSpecifier != nil {
				if importHelpersImportSpecifiers == nil {
					importHelpersImportSpecifiers = make(map[tspath.Path]*ast.Node, totalFileCount)
				}
				importHelpersImportSpecifiers[path] = task.importHelpersImportSpecifier
			}
			if data.fromExternalLibrary {
				sourceFilesFoundSearchingNodeModules.Add(path)
			}
		}
	}

	collectFiles(loader.rootTasks, collections.Set[*parseTaskData]{})
	loader.sortLibs(libFiles)

	allFiles := append(libFiles, files...)

	keys := slices.Collect(loader.pathForLibFileResolutions.Keys())
	slices.Sort(keys)
	for _, key := range keys {
		value, _ := loader.pathForLibFileResolutions.Load(key)
		resolvedModules[key] = module.ModeAwareCache[*module.ResolvedModule]{
			module.ModeAwareCacheKey{Name: value.libraryName, Mode: core.ModuleKindCommonJS}: value.resolution,
		}
		for _, trace := range value.trace {
			loader.opts.Host.Trace(trace.Message, trace.Args...)
		}
	}

	return processedFiles{
		resolver:                             loader.resolver,
		files:                                allFiles,
		filesByPath:                          filesByPath,
		projectReferenceFileMapper:           loader.projectReferenceFileMapper,
		resolvedModules:                      resolvedModules,
		typeResolutionsInFile:                typeResolutionsInFile,
		sourceFileMetaDatas:                  sourceFileMetaDatas,
		jsxRuntimeImportSpecifiers:           jsxRuntimeImportSpecifiers,
		importHelpersImportSpecifiers:        importHelpersImportSpecifiers,
		sourceFilesFoundSearchingNodeModules: sourceFilesFoundSearchingNodeModules,
		libFiles:                             libFilesMap,
		missingFiles:                         missingFiles,
		includeProcessor:                     loader.includeProcessor,
		outputFileToProjectReferenceSource:   outputFileToProjectReferenceSource,
	}
}

func (w *filesParser) addIncludeReason(loader *fileLoader, task *parseTask, reason *FileIncludeReason) {
	if task.redirectedParseTask != nil {
		w.addIncludeReason(loader, task.redirectedParseTask, reason)
	} else if task.loaded {
		if existing, ok := loader.includeProcessor.fileIncludeReasons[task.path]; ok {
			loader.includeProcessor.fileIncludeReasons[task.path] = append(existing, reason)
		} else {
			loader.includeProcessor.fileIncludeReasons[task.path] = []*FileIncludeReason{reason}
		}
	}
}
