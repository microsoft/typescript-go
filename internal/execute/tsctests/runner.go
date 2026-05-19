package tsctests

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/execute"
	"github.com/microsoft/typescript-go/internal/execute/tsc"
	"github.com/microsoft/typescript-go/internal/testutil/baseline"
	"github.com/microsoft/typescript-go/internal/tspath"
)

type tscEdit struct {
	caption         string
	commandLineArgs []string
	edit            func(*TestSys)
	expectedDiff    string
}

var noChange = &tscEdit{
	caption: "no change",
}

var noChangeOnlyEdit = []*tscEdit{
	noChange,
}

type tscInput struct {
	subScenario      string
	commandLineArgs  []string
	files            FileMap
	cwd              string
	edits            []*tscEdit
	env              map[string]string
	ignoreCase       bool
	windowsStyleRoot string
	// extracted is true when this tscInput was constructed from a [TestSpec]
	// (e.g. by a generated test file). When true, run() skips the call to
	// writeTestSourceFile to avoid re-emitting the same generated file.
	extracted bool
}

func (test *tscInput) executeCommand(sys *TestSys, baselineBuilder *strings.Builder, commandLineArgs []string) tsc.CommandLineResult {
	fmt.Fprint(baselineBuilder, "tsgo ", strings.Join(commandLineArgs, " "), "\n")
	result := execute.CommandLine(sys, commandLineArgs, sys)
	switch result.Status {
	case tsc.ExitStatusSuccess:
		baselineBuilder.WriteString("ExitStatus:: Success")
	case tsc.ExitStatusDiagnosticsPresent_OutputsSkipped:
		baselineBuilder.WriteString("ExitStatus:: DiagnosticsPresent_OutputsSkipped")
	case tsc.ExitStatusDiagnosticsPresent_OutputsGenerated:
		baselineBuilder.WriteString("ExitStatus:: DiagnosticsPresent_OutputsGenerated")
	case tsc.ExitStatusInvalidProject_OutputsSkipped:
		baselineBuilder.WriteString("ExitStatus:: InvalidProject_OutputsSkipped")
	case tsc.ExitStatusProjectReferenceCycle_OutputsSkipped:
		baselineBuilder.WriteString("ExitStatus:: ProjectReferenceCycle_OutputsSkipped")
	case tsc.ExitStatusNotImplemented:
		baselineBuilder.WriteString("ExitStatus:: NotImplemented")
	default:
		panic(fmt.Sprintf("UnknownExitStatus %d", result.Status))
	}
	return result
}

func (test *tscInput) run(t *testing.T, scenario string) {
	t.Helper()
	t.Run(test.getBaselineSubFolder()+"/"+test.subScenario, func(t *testing.T) {
		t.Parallel()
		state := test.start(t, scenario)
		for _, do := range test.edits {
			test.runEdit(state, do)
		}
		test.end(state)
	})
}

// tscInputState holds the per-run state created by [tscInput.start] and
// threaded through [tscInput.runEdit] and [tscInput.end]. It is used both by
// the legacy declarative [tscInput.run] entry point and by the imperative
// Start/Edit/End methods on [TestSpec] used by generated test files.
type tscInputState struct {
	t               *testing.T
	scenario        string
	sys             *TestSys
	result          tsc.CommandLineResult
	baselineBuilder *strings.Builder
	unexpectedDiff  *strings.Builder
	// editOps records the file-system operations observed while running each
	// edit. Only populated when the test was not loaded from an extracted
	// source file, since generated tests do not need to re-emit themselves.
	editOps [][]capturedEditOp
	// appliedEdits is the list of edits that have been run so far (in order).
	// The non-incremental comparison build needs to replay every preceding
	// edit before applying the current one, so we keep them around for the
	// lifetime of the test.
	appliedEdits []*tscEdit
	// editIndex is the 0-based index of the next edit to be run.
	editIndex int
	// ended indicates that [tscInput.end] has already run, so subsequent
	// Edit/End calls can be rejected.
	ended bool
}

// start performs the initial build for a test and returns the in-progress
// state. Callers must invoke runEdit zero or more times followed by end.
func (test *tscInput) start(t *testing.T, scenario string) *tscInputState {
	t.Helper()
	state := &tscInputState{
		t:               t,
		scenario:        scenario,
		baselineBuilder: &strings.Builder{},
		unexpectedDiff:  &strings.Builder{},
	}
	state.sys = newTestSys(test, false)
	fmt.Fprint(
		state.baselineBuilder,
		"currentDirectory::",
		state.sys.GetCurrentDirectory(),
		"\nuseCaseSensitiveFileNames::",
		state.sys.FS().UseCaseSensitiveFileNames(),
		"\nInput::\n",
	)
	state.sys.baselineFSwithDiff(state.baselineBuilder)
	state.result = test.executeCommand(state.sys, state.baselineBuilder, test.commandLineArgs)
	state.sys.serializeState(state.baselineBuilder)
	state.unexpectedDiff.WriteString(state.sys.baselinePrograms(state.baselineBuilder, "Initial build"))
	return state
}

// runEdit applies a single edit to the in-progress state, runs the
// incremental and non-incremental builds, and appends the result to the
// baseline being accumulated.
func (test *tscInput) runEdit(state *tscInputState, do *tscEdit) {
	if state.ended {
		panic("tscInput.runEdit called after end")
	}
	index := state.editIndex
	state.editIndex++
	state.appliedEdits = append(state.appliedEdits, do)
	extractEdits := !test.extracted

	sys := state.sys
	sys.clearOutput()
	var nonIncrementalSys *TestSys
	commandLineArgs := core.IfElse(do.commandLineArgs == nil, test.commandLineArgs, do.commandLineArgs)
	state.baselineBuilder.WriteString(fmt.Sprintf("\n\nEdit [%d]:: %s\n", index, do.caption))
	var beforeSnap map[string]capturedFsEntry
	if extractEdits {
		beforeSnap = captureFsSnapshot(sys)
	}
	if do.edit != nil {
		do.edit(sys)
	}
	if extractEdits {
		for len(state.editOps) <= index {
			state.editOps = append(state.editOps, nil)
		}
		state.editOps[index] = diffFsSnapshots(beforeSnap, captureFsSnapshot(sys))
	}
	sys.baselineFSwithDiff(state.baselineBuilder)

	if state.result.Watcher == nil {
		test.executeCommand(sys, state.baselineBuilder, commandLineArgs)
	} else {
		state.result.Watcher.DoCycle()
	}
	sys.serializeState(state.baselineBuilder)
	state.unexpectedDiff.WriteString(sys.baselinePrograms(state.baselineBuilder, fmt.Sprintf("Edit [%d]:: %s\n", index, do.caption)))

	// Compute build with all the edits applied so far in order.
	nonIncrementalSys = newTestSys(test, true)
	for _, e := range state.appliedEdits {
		if e.edit != nil {
			e.edit(nonIncrementalSys)
		}
	}
	execute.CommandLine(nonIncrementalSys, commandLineArgs, nonIncrementalSys)

	diff := getDiffForIncremental(sys, nonIncrementalSys)
	if diff != "" {
		state.baselineBuilder.WriteString(fmt.Sprintf("\n\nDiff:: %s\n", core.IfElse(do.expectedDiff == "", "!!! Unexpected diff, please review and either fix or write explanation as expectedDiff !!!", do.expectedDiff)))
		state.baselineBuilder.WriteString(diff)
		if do.expectedDiff == "" {
			state.unexpectedDiff.WriteString(fmt.Sprintf("Edit [%d]:: %s\n!!! Unexpected diff, please review and either fix or write explanation as expectedDiff !!!\n%s\n", index, do.caption, diff))
		}
	} else if do.expectedDiff != "" {
		state.baselineBuilder.WriteString(fmt.Sprintf("\n\nDiff:: %s !!! Diff not found but explanation present, please review and remove the explanation !!!\n", do.expectedDiff))
		state.unexpectedDiff.WriteString(fmt.Sprintf("Edit [%d]:: %s\n!!! Diff not found but explanation present, please review and remove the explanation !!!\n", index, do.caption))
	}
}

// end finalises the run by writing the extracted test source file (when
// applicable), running the baseline comparison, and reporting unexpected
// diffs against the non-incremental build.
func (test *tscInput) end(state *tscInputState) {
	if state.ended {
		panic("tscInput.end called twice")
	}
	state.ended = true
	if !test.extracted {
		test.writeTestSourceFile(state.scenario, state.editOps)
	}
	baseline.Run(state.t, strings.ReplaceAll(test.subScenario, " ", "-")+".js", state.baselineBuilder.String(), baseline.Options{Subfolder: filepath.Join(test.getBaselineSubFolder(), state.scenario)})
	if state.unexpectedDiff.String() != "" {
		state.t.Errorf("Test %s has unexpected diff %s with incremental build, please review the baseline file", test.subScenario, state.unexpectedDiff.String())
	}
}

func getDiffForIncremental(incrementalSys *TestSys, nonIncrementalSys *TestSys) string {
	var diffBuilder strings.Builder

	nonIncrementalOutputs := nonIncrementalSys.fs.writtenFiles.ToSlice()
	slices.Sort(nonIncrementalOutputs)
	for _, nonIncrementalOutput := range nonIncrementalOutputs {
		if tspath.FileExtensionIs(nonIncrementalOutput, tspath.ExtensionTsBuildInfo) ||
			strings.HasSuffix(nonIncrementalOutput, ".readable.baseline.txt") {
			// Just check existence
			if !incrementalSys.fsFromFileMap().FileExists(nonIncrementalOutput) {
				diffBuilder.WriteString(baseline.DiffText("nonIncremental "+nonIncrementalOutput, "incremental "+nonIncrementalOutput, "Exists", ""))
				diffBuilder.WriteString("\n")
			}
		} else {
			nonIncrementalText, ok := nonIncrementalSys.fsFromFileMap().ReadFile(nonIncrementalOutput)
			if !ok {
				panic("Written file not found " + nonIncrementalOutput)
			}
			incrementalText, ok := incrementalSys.fsFromFileMap().ReadFile(nonIncrementalOutput)
			if !ok || incrementalText != nonIncrementalText {
				diffBuilder.WriteString(baseline.DiffText("nonIncremental "+nonIncrementalOutput, "incremental "+nonIncrementalOutput, nonIncrementalText, incrementalText))
				diffBuilder.WriteString("\n")
			}
		}
	}

	incrementalOutput := incrementalSys.getOutput(true)
	nonIncrementalOutput := nonIncrementalSys.getOutput(true)
	if incrementalOutput != nonIncrementalOutput {
		diffBuilder.WriteString(baseline.DiffText("nonIncremental.output.txt", "incremental.output.txt", nonIncrementalOutput, incrementalOutput))
	}
	return diffBuilder.String()
}

func (test *tscInput) getBaselineSubFolder() string {
	commandName := "tsc"
	if slices.ContainsFunc(test.commandLineArgs, func(arg string) bool {
		switch arg {
		case "-b", "--b", "-build", "--build":
			return true
		}
		return false
	}) {
		commandName = "tsbuild"
	}
	w := ""
	if slices.ContainsFunc(test.commandLineArgs, func(arg string) bool {
		switch arg {
		case "-w", "--w", "-watch", "--watch":
			return true
		}
		return false
	}) {
		w = "Watch"
	}
	return commandName + w
}
