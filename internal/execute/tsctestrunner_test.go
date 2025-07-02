package execute_test

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/execute"
	"github.com/microsoft/typescript-go/internal/incremental"
	"github.com/microsoft/typescript-go/internal/testutil/baseline"
)

type testTscEdit struct {
	caption         string
	commandLineArgs []string
	edit            func(*testSys)
}

var noChange = &testTscEdit{
	caption: "no change",
}

var noChangeOnlyEdit = []*testTscEdit{
	noChange,
}

type tscInput struct {
	subScenario     string
	commandLineArgs []string
	sys             *testSys
	edits           []*testTscEdit
}

func (test *tscInput) executeCommand(baselineBuilder *strings.Builder, commandLineArgs []string) (*incremental.Program, *execute.Watcher) {
	fmt.Fprint(baselineBuilder, "tsgo ", strings.Join(commandLineArgs, " "), "\n")
	exit, incrementalProgram, watcher := execute.CommandLine(test.sys, commandLineArgs, true)
	switch exit {
	case execute.ExitStatusSuccess:
		baselineBuilder.WriteString("ExitStatus:: Success")
	case execute.ExitStatusDiagnosticsPresent_OutputsSkipped:
		baselineBuilder.WriteString("ExitStatus:: DiagnosticsPresent_OutputsSkipped")
	case execute.ExitStatusDiagnosticsPresent_OutputsGenerated:
		baselineBuilder.WriteString("ExitStatus:: DiagnosticsPresent_OutputsGenerated")
	case execute.ExitStatusInvalidProject_OutputsSkipped:
		baselineBuilder.WriteString("ExitStatus:: InvalidProject_OutputsSkipped")
	case execute.ExitStatusProjectReferenceCycle_OutputsSkipped:
		baselineBuilder.WriteString("ExitStatus:: ProjectReferenceCycle_OutputsSkipped")
	case execute.ExitStatusNotImplemented:
		baselineBuilder.WriteString("ExitStatus:: NotImplemented")
	default:
		panic(fmt.Sprintf("UnknownExitStatus %d", exit))
	}
	return incrementalProgram, watcher
}

func (test *tscInput) run(t *testing.T, scenario string) {
	t.Helper()
	// !!! sheetal TODO :: add incremental correctness
	t.Run(test.subScenario+" tsc baseline", func(t *testing.T) {
		t.Parallel()
		// initial test tsc compile
		baselineBuilder := &strings.Builder{}
		fmt.Fprint(
			baselineBuilder,
			"currentDirectory::",
			test.sys.GetCurrentDirectory(),
			"\nuseCaseSensitiveFileNames::",
			test.sys.FS().UseCaseSensitiveFileNames(),
			"\nInput::\n",
		)
		test.sys.baselineFSwithDiff(baselineBuilder)
		incrementalProgram, watcher := test.executeCommand(baselineBuilder, test.commandLineArgs)
		test.sys.serializeState(baselineBuilder)
		test.sys.baselineProgram(baselineBuilder, incrementalProgram, watcher)

		for index, do := range test.edits {
			baselineBuilder.WriteString(fmt.Sprintf("\n\nEdit [%d]:: %s\n", index, do.caption))
			if do.edit != nil {
				do.edit(test.sys)
			}
			test.sys.baselineFSwithDiff(baselineBuilder)

			var incrementalProgram *incremental.Program
			if watcher == nil {
				commandLineArgs := core.IfElse(do.commandLineArgs == nil, test.commandLineArgs, do.commandLineArgs)
				incrementalProgram, watcher = test.executeCommand(baselineBuilder, commandLineArgs)
			} else {
				watcher.DoCycle()
			}
			test.sys.serializeState(baselineBuilder)
			test.sys.baselineProgram(baselineBuilder, incrementalProgram, watcher)
		}

		baseline.Run(t, strings.ReplaceAll(test.subScenario, " ", "-")+".js", baselineBuilder.String(), baseline.Options{Subfolder: filepath.Join(test.getBaselineSubFolder(), scenario)})
	})
}

func (test *tscInput) getBaselineSubFolder() string {
	commandName := "tsc"
	if slices.ContainsFunc(test.commandLineArgs, func(arg string) bool {
		return arg == "--build" || arg == "-b"
	}) {
		commandName = "tsbuild"
	}
	w := ""
	if slices.ContainsFunc(test.commandLineArgs, func(arg string) bool {
		return arg == "--watch" || arg == "-w"
	}) {
		w = "Watch"
	}
	return commandName + w
}
