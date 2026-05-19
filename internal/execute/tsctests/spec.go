package tsctests

import "testing"

// TestSpec is the exported representation of a tscInput. It exists so that
// generated test files in subpackages can construct and drive a single test
// scenario through the [TestSpec.Start], [TestSpec.Edit], and [TestSpec.End]
// methods.
type TestSpec struct {
	// Scenario is the grouping label normally passed to (*tscInput).run, e.g.
	// "commandLine", "noEmit", or "extends". Used as part of the baseline path.
	Scenario string

	// SubScenario is the per-test sub-name used by the baseline runner.
	SubScenario string

	CommandLineArgs  []string
	Files            FileMap
	Cwd              string
	Env              map[string]string
	IgnoreCase       bool
	WindowsStyleRoot string

	// input and state are populated by Start and consumed by Edit/End.
	input *tscInput
	state *tscInputState
}

// TestEdit is the exported representation of a tscEdit.
type TestEdit struct {
	Caption         string
	CommandLineArgs []string
	Edit            func(*TestSys)
	ExpectedDiff    string
}

// Start performs the initial build for the spec on the provided testing.T.
// It must be called exactly once before any call to Edit or End, and is
// responsible for calling t.Parallel so the extracted tests run concurrently.
func (s *TestSpec) Start(t *testing.T) {
	t.Helper()
	if s.input != nil {
		panic("TestSpec.Start called twice")
	}
	s.input = &tscInput{
		subScenario:      s.SubScenario,
		commandLineArgs:  s.CommandLineArgs,
		files:            s.Files,
		cwd:              s.Cwd,
		env:              s.Env,
		ignoreCase:       s.IgnoreCase,
		windowsStyleRoot: s.WindowsStyleRoot,
		extracted:        true,
	}
	t.Parallel()
	s.state = s.input.start(t, s.Scenario)
}

// Edit applies the given edit to the running test and runs the incremental
// and non-incremental builds, accumulating the result into the baseline.
// Passing nil is treated as a "no change" edit.
func (s *TestSpec) Edit(e *TestEdit) {
	if s.state == nil {
		panic("TestSpec.Edit called before Start")
	}
	if e == nil {
		s.input.runEdit(s.state, noChange)
		return
	}
	s.input.runEdit(s.state, &tscEdit{
		caption:         e.Caption,
		commandLineArgs: e.CommandLineArgs,
		edit:            e.Edit,
		expectedDiff:    e.ExpectedDiff,
	})
}

// End finalises the run, writing the accumulated baseline and reporting any
// unexpected diffs against the non-incremental build.
func (s *TestSpec) End() {
	if s.state == nil {
		panic("TestSpec.End called before Start")
	}
	s.input.end(s.state)
}

// WriteFile writes content to the given path on the test file system,
// panicking on error. Exposed so generated tests can describe edits as a
// series of file system operations.
func (s *TestSys) WriteFile(path string, content string) {
	s.writeFileNoError(path, content)
}

// Remove removes the file or directory at path on the test file system,
// panicking on error. Exposed for generated edit functions.
func (s *TestSys) Remove(path string) {
	s.removeNoError(path)
}
