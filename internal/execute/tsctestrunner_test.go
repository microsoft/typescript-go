package execute_test

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/execute"
	"github.com/microsoft/typescript-go/internal/testutil/baseline"
)

type testTscEdit struct {
	caption         string
	commandLineArgs []string
	edit            func(*testSys)
}

type tscInput struct {
	subScenario     string
	commandLineArgs []string
	sys             *testSys
	edits           []*testTscEdit
}

func (test *tscInput) run(t *testing.T, scenario string) {
	t.Helper()
	t.Run(test.getTestName(scenario), func(t *testing.T) {
		t.Parallel()
		t.Run("tsc baseline", func(t *testing.T) {
			t.Parallel()
			// initial test tsc compile
			baselineBuilder := test.startBaseline()

			exit, parsedCommandLine, watcher := execute.CommandLineTest(test.sys, test.commandLineArgs)
			baselineBuilder.WriteString("ExitStatus:: " + fmt.Sprint(exit))

			compilerOptionsString, _ := json.MarshalIndent(parsedCommandLine.CompilerOptions(), "", "    ")
			baselineBuilder.WriteString("\n\nCompilerOptions::")
			baselineBuilder.Write(compilerOptionsString)

			if watcher != nil {
				execute.StartForTest(watcher)
			}

			test.sys.serializeState(baselineBuilder)

			for _, do := range test.edits {
				do.edit(test.sys)
				baselineBuilder.WriteString("\n\nEdit:: " + do.caption + "\n")

				if watcher == nil {
					exit, parsedCommandLine, watcher = execute.CommandLineTest(test.sys, test.commandLineArgs)
					baselineBuilder.WriteString("ExitStatus:: " + fmt.Sprint(exit))
				} else {
					execute.RunWatchCycle(watcher)
				}
				test.sys.serializeState(baselineBuilder)
			}

			options, name := test.getBaselineName(scenario, "")
			baseline.Run(t, name, baselineBuilder.String(), options)
		})

		// !!! sheetal TODO :: add incremental correctness
	})
}

func (test *tscInput) getTestNamePrefix() string {
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

func (test *tscInput) getTestName(scenario string) string {
	return test.getTestNamePrefix() + " " + scenario + ":: " + test.subScenario + " " + strings.Join(test.commandLineArgs, " ")
}

func (test *tscInput) getBaselineName(scenario string, suffix string) (baseline.Options, string) {
	return baseline.Options{Subfolder: filepath.Join(test.getTestNamePrefix(), scenario)},
		strings.ReplaceAll(test.subScenario, " ", "-") + suffix + ".js"
}

func (test *tscInput) startBaseline() *strings.Builder {
	s := &strings.Builder{}
	fmt.Fprint(
		s,
		"\ncurrentDirectory::",
		test.sys.GetCurrentDirectory(),
		"\nuseCaseSensitiveFileNames::",
		test.sys.FS().UseCaseSensitiveFileNames(),
		"\nInput::",
	)
	fmt.Fprint(s, strings.Join(test.commandLineArgs, " "), "\n")
	test.sys.baselineFSwithDiff(s)
	return s
}

func (test *tscInput) verifyCommandLineParsing(t *testing.T, scenario string) {
	t.Helper()
	t.Run(test.getTestName(scenario), func(t *testing.T) {
		t.Parallel()
		t.Run("baseline for the tsc compiles", func(t *testing.T) {
			t.Parallel()
			// initial test tsc compile
			baselineBuilder := test.startBaseline()

			exit, parsedCommandLine, _ := execute.CommandLineTest(test.sys, test.commandLineArgs)
			baselineBuilder.WriteString("ExitStatus:: " + fmt.Sprint(exit))
			//nolint:musttag
			parsedCommandLineString, _ := json.MarshalIndent(parsedCommandLine, "", "    ")
			baselineBuilder.WriteString("\n\nParsedCommandLine::")
			baselineBuilder.Write(parsedCommandLineString)

			test.sys.serializeState(baselineBuilder)
			options, name := test.getBaselineName(scenario, "")
			baseline.Run(t, name, baselineBuilder.String(), options)
		})
	})
}
