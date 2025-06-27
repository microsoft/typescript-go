package incremental

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/microsoft/typescript-go/internal/ast"
	"github.com/microsoft/typescript-go/internal/checker"
	"github.com/microsoft/typescript-go/internal/collections"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type snapshotChange interface {
	commit(snapshot *snapshot)
}

var _ snapshotChange = (*semanticDiagnosticChange)(nil)

type semanticDiagnosticChange struct {
	// new diagnostics
	semanticDiagnosticsPerFile map[tspath.Path][]*ast.Diagnostic
}

func (c *semanticDiagnosticChange) commit(snapshot *snapshot) {
	for file, diagnostics := range c.semanticDiagnosticsPerFile {
		snapshot.semanticDiagnosticsFromOldState.Delete(file)
		snapshot.semanticDiagnosticsPerFile[file] = &diagnosticsOrBuildInfoDiagnosticsWithFileName{diagnostics: diagnostics}
	}
}

type affectedFilesChange struct {
	ctx                                    context.Context
	program                                *Program
	hasAllFilesExcludingDefaultLibraryFile atomic.Bool
	updatedSignatures                      collections.SyncMap[tspath.Path, string]
	allFilesPendingEmit                    []map[tspath.Path]FileEmitKind
	filesToRemoveDiagnostics               collections.SyncSet[tspath.Path]
	cleanedDiagnosticsOfLibFiles           sync.Once
}

var _ snapshotChange = (*affectedFilesChange)(nil)

func (c *affectedFilesChange) commit(snapshot *snapshot) {
	c.updatedSignatures.Range(func(filePath tspath.Path, signature string) bool {
		snapshot.fileInfos[filePath].signature = signature
		return true
	})
	c.filesToRemoveDiagnostics.Range(func(file tspath.Path) bool {
		snapshot.semanticDiagnosticsFromOldState.Delete(file)
		delete(snapshot.semanticDiagnosticsPerFile, file)
		return true
	})
	for _, affectedFilesPendingEmit := range c.allFilesPendingEmit {
		for filePath, emitKind := range affectedFilesPendingEmit {
			snapshot.addFileToAffectedFilesPendingEmit(filePath, emitKind)
		}
	}
	snapshot.changedFilesSet = &collections.Set[tspath.Path]{}
}

func (c *affectedFilesChange) isChangedSignature(path tspath.Path) bool {
	newSignature, _ := c.updatedSignatures.Load(path)
	oldSignature := c.program.snapshot.fileInfos[path].signature
	return newSignature != oldSignature
}

func (c *affectedFilesChange) getDtsMayChange(affectedFilePath tspath.Path, affectedFileEmitKind FileEmitKind) *dtsMayChange {
	affectedFilesPendingEmit := map[tspath.Path]FileEmitKind{affectedFilePath: affectedFileEmitKind}
	c.allFilesPendingEmit = append(c.allFilesPendingEmit, affectedFilesPendingEmit)
	return &dtsMayChange{
		change:                   c,
		affectedFilesPendingEmit: affectedFilesPendingEmit,
	}
}

func (c *affectedFilesChange) removeSemanticDiagnosticsOf(path tspath.Path) {
	if c.program.snapshot.semanticDiagnosticsFromOldState.Has(path) {
		c.filesToRemoveDiagnostics.Add(path)
	}
}

func (c *affectedFilesChange) removeDiagnosticsOfLibraryFiles() {
	c.cleanedDiagnosticsOfLibFiles.Do(func() {
		for _, file := range c.program.GetSourceFiles() {
			if c.program.program.IsSourceFileDefaultLibrary(file.Path()) && !checker.SkipTypeChecking(file, c.program.snapshot.options, c.program.program, true) {
				c.removeSemanticDiagnosticsOf(file.Path())
			}
		}
	})
}

type dtsMayChange struct {
	change                   *affectedFilesChange
	affectedFilesPendingEmit map[tspath.Path]FileEmitKind
}

func (c *dtsMayChange) addFileToAffectedFilesPendingEmit(filePath tspath.Path, emitKind FileEmitKind) {
	c.affectedFilesPendingEmit[filePath] = emitKind
}
