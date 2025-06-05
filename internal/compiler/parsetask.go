package compiler

import (
	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/module"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type parseTask struct {
	normalizedFilePath string
	file               *ast.SourceFile
	isLib              bool
	subTasks           []*parseTask

	metadata                     *ast.SourceFileMetaData
	resolutionsInFile            module.ModeAwareCache[*module.ResolvedModule]
	importHelpersImportSpecifier *ast.Node
	jsxRuntimeImportSpecifier    *jsxRuntimeImportSpecifier
}

func (t *parseTask) start(loader *fileLoader) {
	loader.totalFileCount.Add(1)
	if t.isLib {
		loader.libFileCount.Add(1)
	}

	loader.wg.Queue(func() {
		file := loader.parseSourceFile(t.normalizedFilePath)
		if file == nil {
			return
		}

		t.file = file
		t.metadata = loader.loadSourceFileMetaData(file.FileName())

		// !!! if noResolve, skip all of this
		t.subTasks = make([]*parseTask, 0, len(file.ReferencedFiles)+len(file.Imports())+len(file.ModuleAugmentations))

		for _, ref := range file.ReferencedFiles {
			resolvedPath := loader.resolveTripleslashPathReference(ref.FileName, file.FileName())
			t.addSubTask(resolvedPath, false)
		}

		compilerOptions := loader.opts.Config.CompilerOptions()
		for _, ref := range file.TypeReferenceDirectives {
			resolutionMode := getModeForTypeReferenceDirectiveInFile(ref, file, t.metadata, compilerOptions)
			resolved := loader.resolver.ResolveTypeReferenceDirective(ref.FileName, file.FileName(), resolutionMode, nil)
			if resolved.IsResolved() {
				t.addSubTask(resolved.ResolvedFileName, false)
			}
		}

		if compilerOptions.NoLib != core.TSTrue {
			for _, lib := range file.LibReferenceDirectives {
				name, ok := tsoptions.GetLibFileName(lib.FileName)
				if !ok {
					continue
				}
				t.addSubTask(tspath.CombinePaths(loader.defaultLibraryPath, name), true)
			}
		}

		toParse, resolutionsInFile, importHelpersImportSpecifier, jsxRuntimeImportSpecifier := loader.resolveImportsAndModuleAugmentations(file, t.metadata)
		for _, imp := range toParse {
			t.addSubTask(imp, false)
		}

		t.resolutionsInFile = resolutionsInFile
		t.importHelpersImportSpecifier = importHelpersImportSpecifier
		t.jsxRuntimeImportSpecifier = jsxRuntimeImportSpecifier

		loader.startTasks(t.subTasks)
	})
}

func (t *parseTask) addSubTask(fileName string, isLib bool) {
	normalizedFilePath := tspath.NormalizePath(fileName)
	t.subTasks = append(t.subTasks, &parseTask{normalizedFilePath: normalizedFilePath, isLib: isLib})
}
