package execute

import (
	"fmt"
	"time"

	"github.com/microsoft/typescript-go/internal/compiler"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tspath"
)

func start(w *watcher) ExitStatus {
	if w.configFileName == "" {
		w.host = compiler.NewCompilerHost(w.options.CompilerOptions(), w.sys.GetCurrentDirectory(), w.sys.FS(), w.sys.DefaultLibraryPath())
	}
	watchInterval := 500 * time.Millisecond
	for !w.sys.IsTestDone() {
		if w.configFileName != "" {
			// only need to reparse tsconfig options/update host if we are watching a config file
			extendedConfigCache := map[tspath.Path]*tsoptions.ExtendedConfigCacheEntry{}
			configParseResult, errors := getParsedCommandLineOfConfigFile(w.configFileName, w.options.CompilerOptions(), w.sys, extendedConfigCache)
			if len(errors) > 0 {
				// these are unrecoverable errors--report them and do not build
				for _, e := range errors {
					w.reportDiagnostic(e)
				}
				continue
			}
			w.options = configParseResult
			w.host = compiler.NewCompilerHost(w.options.CompilerOptions(), w.sys.GetCurrentDirectory(), w.sys.FS(), w.sys.DefaultLibraryPath())
		}
		w.program = compiler.NewProgramFromParsedCommandLine(w.options, w.host)
		if hasBeenModified(w, w.program) {
			fmt.Fprint(w.sys.Writer(), "build starting at ", w.sys.Now(), w.sys.NewLine())
			w.compileAndEmit()
			fmt.Fprint(w.sys.Writer(), "build finished ", w.sys.Now(), w.sys.NewLine())
		} else {
			// print something???
			fmt.Fprint(w.sys.Writer(), "no changes detected at ", w.sys.Now(), w.sys.NewLine())
		}
		time.Sleep(watchInterval)
	}
	return ExitStatusSuccess
}

func hasBeenModified(w *watcher, program *compiler.Program) bool {
	// checks watcher's snapshot against program file modified times
	currState := map[string]time.Time{}
	filesModified := false
	for _, sourceFile := range program.SourceFiles() {
		fileName := sourceFile.FileName()
		currState[fileName] = w.sys.FS().Stat(fileName).ModTime()
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

func WatchTest(sys System, commandLineArgs []string) {
	commandLine := tsoptions.ParseCommandLine(commandLineArgs, sys)
	configFileName := tspath.NormalizePath(commandLine.CompilerOptions().Project)
	if configFileName != "" || sys.FS().DirectoryExists(configFileName) {
		configFileName = tspath.CombinePaths(configFileName, "tsconfig.json")
	}
	if !sys.FS().FileExists(configFileName) {
		return
	}
	if configFileName != "" {
		extendedConfigCache := map[tspath.Path]*tsoptions.ExtendedConfigCacheEntry{}
		commandLine, _ = getParsedCommandLineOfConfigFile(configFileName, commandLine.CompilerOptions(), sys, extendedConfigCache)
	}

	w := createWatcher(sys, commandLine, createDiagnosticReporter(sys, commandLine.CompilerOptions().Pretty))

	w.host = compiler.NewCompilerHost(w.options.CompilerOptions(), w.sys.GetCurrentDirectory(), w.sys.FS(), w.sys.DefaultLibraryPath())

	totalTime := time.Duration(0)
	for i := range 10 {
		startTime := time.Now()

		extendedConfigCache := map[tspath.Path]*tsoptions.ExtendedConfigCacheEntry{}
		w.options, _ = getParsedCommandLineOfConfigFile(w.configFileName, w.options.CompilerOptions(), w.sys, extendedConfigCache)
		configTime := time.Now()

		w.host = compiler.NewCompilerHost(w.options.CompilerOptions(), w.sys.GetCurrentDirectory(), w.sys.FS(), w.sys.DefaultLibraryPath())
		hostTime := time.Now()

		w.program = compiler.NewProgramFromParsedCommandLine(w.options, w.host)
		timeProgram := time.Now()

		hasBeenModified(w, w.program)
		timeModified := time.Now()

		timeCheck := timeModified.Sub(startTime)
		fmt.Fprint(w.sys.Writer(), i,
			": ", timeCheck,
			"    ", configTime.Sub(startTime),
			"    ", hostTime.Sub(configTime),
			"    ", timeProgram.Sub(hostTime),
			"    ", timeModified.Sub(timeProgram), sys.NewLine())
		totalTime += timeCheck
	}
	fmt.Fprint(w.sys.Writer(), "total time taken ", totalTime, sys.NewLine())
}
