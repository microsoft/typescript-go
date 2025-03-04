package execute_test

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/execute"
	"github.com/microsoft/typescript-go/internal/testutil/baseline"
)

type testTscEdit struct {
	edit            func(execute.System)
	caption         string
	commandLineArgs []string
}

type tscInput struct {
	subScenario     string
	commandLineArgs []string
	sys             *testSys

	// for watch tests
	data       map[string]string
	edits       []*testTscEdit
}

func (test *tscInput) verify(t *testing.T, scenario string) {
	t.Helper()
	t.Run(test.getTestName(scenario), func(t *testing.T) {
		t.Parallel()
		t.Run("baseline for the tsc compiles", func(t *testing.T) {
			t.Parallel()
			// initial test tsc compile
			baselineBuilder := test.startBaseline()

			parsedCommandLine, exit := execute.CommandLineTest(test.sys, nil, test.commandLineArgs)
			baselineBuilder.WriteString("ExitStatus:: " + fmt.Sprint(exit))

			compilerOptionsString, _ := json.MarshalIndent(parsedCommandLine.CompilerOptions(), "", "    ")
			baselineBuilder.WriteString("\n\nCompilerOptions::")
			baselineBuilder.Write(compilerOptionsString)

			test.sys.serializeState(baselineBuilder)
			options, name := test.getBaselineName(scenario, "")
			baseline.Run(t, name, baselineBuilder.String(), options)
		})
	})
}

func (test *tscInput) getTestName(scenario string) string {
	return "tsc " + strings.Join(test.commandLineArgs, " ") + " " + scenario + ":: " + test.subScenario
}

func (test *tscInput) getBaselineName(scenario string, suffix string) (baseline.Options, string) {
	commandName := "tsc"
	// todo build
	// if isBuildCommand(v.data.commandLineArgs) {
	// 	commandName = "tsbuild"
	// }
	watch := ""
	if test.data["watch"] != "" { watch = "Watch" }

	return baseline.Options{Subfolder: filepath.Join(commandName+watch, scenario)},
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
