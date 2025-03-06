package execute

import (
	"fmt"
	"time"

	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

func start(w *watcher) ExitStatus {
	if w.configFileName == "" {
		w.host = compiler.NewCompilerHost(w.options.CompilerOptions(), w.sys.GetCurrentDirectory(), w.sys.FS(), w.sys.DefaultLibraryPath())
	}
	watchInterval := 1000 * time.Millisecond
	for {
		watchCycle(w)
		time.Sleep(watchInterval)
	}
}

func watchCycle(w *watcher) {
	if errorsInTsConfig(w) {
		// these are unrecoverable errors--report them and do not build
		return
	}
	// updateProgram()
	w.program = compiler.NewProgramFromParsedCommandLine(w.options, w.host)
	if hasBeenModified(w, w.program) {
		fmt.Fprint(w.sys.Writer(), "build starting at ", w.sys.Now(), w.sys.NewLine())
		timeStart := w.sys.Now()
		w.compileAndEmit()
		fmt.Fprint(w.sys.Writer(), "build finished in", w.sys.Now().Sub(timeStart), w.sys.NewLine())
	} else {
		// print something???
		fmt.Fprint(w.sys.Writer(), "no changes detected at ", w.sys.Now(), w.sys.NewLine())
	}
}

func errorsInTsConfig(w *watcher) bool {
	// only need to check and reparse tsconfig options/update host if we are watching a config file
	if w.configFileName != "" {
		extendedConfigCache := map[tspath.Path]*tsoptions.ExtendedConfigCacheEntry{}
		// !!! need to check that this merges compileroptions correctly. This differs from non-watch, since we allow overriding of previous options
		configParseResult, errors := getParsedCommandLineOfConfigFile(w.configFileName, &core.CompilerOptions{}, w.sys, extendedConfigCache)
		if len(errors) > 0 {
			for _, e := range errors {
				w.reportDiagnostic(e)
			}
			return true
		}
		w.options = configParseResult
		w.host = compiler.NewCompilerHost(w.options.CompilerOptions(), w.sys.GetCurrentDirectory(), w.sys.FS(), w.sys.DefaultLibraryPath())
	}
	return false
}

func hasBeenModified(w *watcher, program *compiler.Program) bool {
	// checks watcher's snapshot against program file modified times
	currState := map[string]time.Time{}
	filesModified := false
	for _, sourceFile := range program.SourceFiles() {
		fileName := sourceFile.FileName()
		s := w.sys.FS().Stat(fileName)
		if s == nil {
			// do nothing; if file is in program.SourceFiles() but is not found when calling Stat, file has been very recently deleted.
			// deleted files are handled outside of this loop
			continue
		}
		currState[fileName] = s.ModTime()
		if !filesModified {
			if currState[fileName] != w.prevModified[fileName] {
				filesModified = true
			}
			delete(w.prevModified, fileName)
		}
	}
	if len(w.prevModified) > 0 {
		filesModified = true
	}
	w.prevModified = currState
	return filesModified
}
