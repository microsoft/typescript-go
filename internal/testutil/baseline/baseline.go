package baseline

import (
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/repo"
	"github.com/microsoft/typescript-go/internal/tspath"
	delta "github.com/octavore/delta/lib"
	"github.com/pkg/diff"
	edit "github.com/pkg/diff/edit"
	write "github.com/pkg/diff/write"
	"gotest.tools/v3/assert"
)

type Options struct {
	Subfolder           string
	IsSubmodule         bool
	IsSubmoduleAccepted bool
	DiffFixupOld        func(string) string
}

const NoContent = "<no content>"

func Run(t *testing.T, fileName string, actual string, opts Options) {
	origSubfolder := opts.Subfolder

	{
		subfolder := opts.Subfolder
		if opts.IsSubmodule {
			subfolder = filepath.Join("submodule", subfolder)
		}

		localPath := filepath.Join(localRoot, subfolder, fileName)
		referencePath := filepath.Join(referenceRoot, subfolder, fileName)

		writeComparison(t, actual, localPath, referencePath, false)
	}

	if !opts.IsSubmodule {
		// Not a submodule, no diffs.
		return
	}

	submoduleReference := filepath.Join(submoduleReferenceRoot, fileName)
	submoduleExpected := readFileOrNoContent(submoduleReference)

	const (
		submoduleFolder         = "submodule"
		submoduleAcceptedFolder = "submoduleAccepted"
	)

	diffFileName := fileName + ".diff"
	isSubmoduleAccepted := opts.IsSubmoduleAccepted || submoduleAcceptedFileNames().Has(origSubfolder+"/"+diffFileName)

	outRoot := core.IfElse(isSubmoduleAccepted, submoduleAcceptedFolder, submoduleFolder)
	unusedOutRoot := core.IfElse(isSubmoduleAccepted, submoduleFolder, submoduleAcceptedFolder)

	{
		localPath := filepath.Join(localRoot, outRoot, origSubfolder, diffFileName)
		referencePath := filepath.Join(referenceRoot, outRoot, origSubfolder, diffFileName)

		diff := getBaselineDiff(t, actual, submoduleExpected, fileName, opts.DiffFixupOld)
		writeComparison(t, diff, localPath, referencePath, false)
	}

	// Delete the other diff file if it exists
	{
		localPath := filepath.Join(localRoot, unusedOutRoot, origSubfolder, diffFileName)
		referencePath := filepath.Join(referenceRoot, unusedOutRoot, origSubfolder, diffFileName)
		writeComparison(t, NoContent, localPath, referencePath, false)
	}
}

var submoduleAcceptedFileNames = sync.OnceValue(func() *core.Set[string] {
	var set core.Set[string]

	submoduleAccepted := filepath.Join(repo.TestDataPath, "submoduleAccepted.txt")
	if content, err := os.ReadFile(submoduleAccepted); err == nil {
		for line := range strings.SplitSeq(string(content), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || line[0] == '#' {
				continue
			}
			set.Add(line)
		}
	} else {
		panic(fmt.Sprintf("failed to read submodule accepted file: %v", err))
	}

	return &set
})

func readFileOrNoContent(fileName string) string {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return NoContent
	}
	return string(content)
}

type inputPair struct {
	a []string
	b []string
}

func (p *inputPair) WriteATo(w io.Writer, ai int) (int, error) {
	return fmt.Fprint(w, p.a[ai])
}

func (p *inputPair) WriteBTo(w io.Writer, bi int) (int, error) {
	return fmt.Fprint(w, p.b[bi])
}

var newlineRE = regexp.MustCompile("\r?\n")

func inputPairFrom(a string, b string) *inputPair {
	return &inputPair{
		a: newlineRE.Split(a, -1),
		b: newlineRE.Split(b, -1),
	}
}

const maxContext = 3

func intoEditScript(sln *delta.DiffSolution) edit.Script {
	var ranges []edit.Range
	aIdx := 0
	bIdx := 0
	lastAIdx := 0
	lastBIdx := 0
	lastMode := ""
	// accumulate all the histogram clusters into edits
	for i, line := range sln.Lines {
		if i == 0 {
			lastMode = line[2]
		}
		if lastMode != line[2] {
			// span kind changed on this line - record span as range
			ranges = addRange(ranges, lastMode, aIdx, bIdx, lastAIdx, lastBIdx, false)
			lastMode = line[2]
			lastAIdx = aIdx
			lastBIdx = bIdx
		}
		switch line[2] {
		case string(delta.LineFromA):
			aIdx++
		case string(delta.LineFromB):
			bIdx++
		case string(delta.LineFromBoth), string(delta.LineFromBothEdit):
			aIdx++
			bIdx++
		}
	}

	ranges = addRange(ranges, lastMode, aIdx, bIdx, lastAIdx, lastBIdx, true)
	return edit.NewScript(ranges...)
}

func addRange(ranges []edit.Range, lastMode string, aIdx int, bIdx int, lastAIdx int, lastBIdx int, eof bool) []edit.Range {
	if lastMode == string(delta.LineFromBothEdit) {
		// record as deleta A, add B
		ranges = append(ranges, edit.Range{
			LowA: lastAIdx, HighA: aIdx,
			LowB: lastBIdx, HighB: lastBIdx,
		})
		return append(ranges, edit.Range{
			LowA: aIdx, HighA: aIdx,
			LowB: lastBIdx, HighB: bIdx,
		})
	} else if lastMode != string(delta.LineFromBoth) {
		return append(ranges, edit.Range{
			LowA: lastAIdx, HighA: aIdx,
			LowB: lastBIdx, HighB: bIdx,
		})
	} else {
		// Ranges in both files are skipped over outside a given number of context lines
		// End of file - one span
		if eof {
			// skip EOF itself
			aIdx--
			bIdx--
			return append(ranges, edit.Range{
				LowA: lastAIdx, HighA: min(aIdx, lastAIdx+maxContext),
				LowB: lastBIdx, HighB: min(bIdx, lastBIdx+maxContext),
			})
		}
		// Start of file - one span
		if lastAIdx == 0 {
			return append(ranges, edit.Range{
				LowA: max(lastAIdx, aIdx-maxContext), HighA: aIdx,
				LowB: max(lastBIdx, bIdx-maxContext), HighB: bIdx,
			})
		}
		// Between two other spans - check span length - less than 2x context limit? Yield one span without eliding any context
		if aIdx-lastAIdx < (maxContext*2 + 1) {
			return append(ranges, edit.Range{
				LowA: lastAIdx, HighA: aIdx,
				LowB: lastBIdx, HighB: bIdx,
			})
		}
		// Two spans - one for the start, one for the end, middle elided
		ranges = append(ranges, edit.Range{
			LowA: lastAIdx, HighA: lastAIdx + maxContext,
			LowB: lastBIdx, HighB: lastBIdx + maxContext,
		})
		return append(ranges, edit.Range{
			LowA: aIdx - maxContext, HighA: aIdx,
			LowB: bIdx - maxContext, HighB: bIdx,
		})
	}
}

func min(a int, b int) int {
	if a <= b {
		return a
	}
	return b
}

func max(a int, b int) int {
	if a >= b {
		return a
	}
	return b
}

func diffText(oldName string, newName string, expected string, actual string, w io.Writer) error {
	// diff.Text uses the Meyers diff algorithm in quadratic space, which performs poorly on very large diffs
	// So instead, we leverage an implementation of the historgram diff algorithm which handles long spans of additions much more quickly
	// But only on tests with output length differences above a certain length (indicating many additions or deletions) so most
	// baselines aren't impacted by the change in diff algorithm
	if math.Abs(float64(strings.Count(expected, "\n")-strings.Count(actual, "\n"))) < 10000 {
		return diff.Text(oldName, newName, expected, actual, w)
	}
	sln := delta.HistogramDiff(expected, actual)
	sln.PostProcess()
	return write.Unified(intoEditScript(sln), w, inputPairFrom(expected, actual), write.Names(oldName, newName))
}

func getBaselineDiff(t *testing.T, actual string, expected string, fileName string, fixupOld func(string) string) string {
	if fixupOld != nil {
		expected = fixupOld(expected)
	}
	if actual == expected {
		return NoContent
	}
	var b strings.Builder
	if err := diffText("old."+fileName, "new."+fileName, expected, actual, &b); err != nil {
		return fmt.Sprintf("failed to diff the actual and expected content: %v\n", err)
	}

	// Remove line numbers from unified diff headers; this avoids adding/deleting
	// lines in our baselines from causing knock-on header changes later in the diff.
	s := b.String()

	aCurLine := 1
	bCurLine := 1
	s = fixUnifiedDiff.ReplaceAllStringFunc(s, func(match string) string {
		var aLine, aLineCount, bLine, bLineCount int
		if _, err := fmt.Sscanf(match, "@@ -%d,%d +%d,%d @@", &aLine, &aLineCount, &bLine, &bLineCount); err != nil {
			panic(fmt.Sprintf("failed to parse unified diff header: %v", err))
		}
		aDiff := aLine - aCurLine
		bDiff := bLine - bCurLine
		aCurLine = aLine
		bCurLine = bLine

		// Keep surrounded by @@, to make GitHub's grammar happy.
		// https://github.com/textmate/diff.tmbundle/blob/0593bb775eab1824af97ef2172fd38822abd97d7/Syntaxes/Diff.plist#L68
		return fmt.Sprintf("@@= skipped -%d, +%d lines =@@", aDiff, bDiff)
	})

	return s
}

var fixUnifiedDiff = regexp.MustCompile(`@@ -\d+,\d+ \+\d+,\d+ @@`)

func RunAgainstSubmodule(t *testing.T, fileName string, actual string, opts Options) {
	local := filepath.Join(localRoot, opts.Subfolder, fileName)
	reference := filepath.Join(submoduleReferenceRoot, opts.Subfolder, fileName)
	writeComparison(t, actual, local, reference, true)
}

func writeComparison(t *testing.T, actualContent string, local, reference string, comparingAgainstSubmodule bool) {
	if actualContent == "" {
		panic("the generated content was \"\". Return 'baseline.NoContent' if no baselining is required.")
	}

	if err := os.MkdirAll(filepath.Dir(local), 0o755); err != nil {
		t.Error(fmt.Errorf("failed to create directories for the local baseline file %s: %w", local, err))
		return
	}

	if _, err := os.Stat(local); err == nil {
		if err := os.Remove(local); err != nil {
			t.Error(fmt.Errorf("failed to remove the local baseline file %s: %w", local, err))
			return
		}
	}

	expected := NoContent
	foundExpected := false
	if content, err := os.ReadFile(reference); err == nil {
		expected = string(content)
		foundExpected = true
	}

	if expected != actualContent || actualContent == NoContent && foundExpected {
		if actualContent == NoContent {
			if err := os.WriteFile(local+".delete", []byte{}, 0o644); err != nil {
				t.Error(fmt.Errorf("failed to write the local baseline file %s: %w", local+".delete", err))
				return
			}
		} else {
			if err := os.WriteFile(local, []byte(actualContent), 0o644); err != nil {
				t.Error(fmt.Errorf("failed to write the local baseline file %s: %w", local, err))
				return
			}
		}

		relReference, err := filepath.Rel(repo.RootPath, reference)
		assert.NilError(t, err)
		relReference = tspath.NormalizeSlashes(relReference)

		relLocal, err := filepath.Rel(repo.RootPath, local)
		assert.NilError(t, err)
		relLocal = tspath.NormalizeSlashes(relLocal)

		if _, err := os.Stat(reference); err != nil {
			if comparingAgainstSubmodule {
				t.Errorf("the baseline file %s does not exist in the TypeScript submodule", relReference)
			} else {
				t.Errorf("new baseline created at %s.", relLocal)
			}
		} else if comparingAgainstSubmodule {
			t.Errorf("the baseline file %s does not match the reference in the TypeScript submodule", relReference)
		} else {
			t.Errorf("the baseline file %s has changed. (Run `hereby baseline-accept` if the new baseline is correct.)", relReference)
		}
	}
}

var (
	localRoot              = filepath.Join(repo.TestDataPath, "baselines", "local")
	referenceRoot          = filepath.Join(repo.TestDataPath, "baselines", "reference")
	submoduleReferenceRoot = filepath.Join(repo.TypeScriptSubmodulePath, "tests", "baselines", "reference")
)
