package baseline

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/microsoft/typescript-go/internal/repo"
)

type Options struct {
	Subfolder string
	Gold      bool
	Expected  bool
}

const NoContent = "<no content>"

func Run(t testing.TB, fileName string, actual string, opts Options) {
	opts.Expected = false
	writeComparison(t, "", actual, fileName, opts)
}

func RunFromText(t testing.TB, fileName string, expected string, actual string, opts Options) {
	opts.Expected = true
	writeComparison(t, expected, actual, fileName, opts)
}

func writeComparison(t testing.TB, expected string, actual string, relativeFileName string, opts Options) {
	if actual == "" || opts.Expected && expected == "" {
		panic("The generated content was \"\". Return 'baseline.NoContent' if no baselining is required.")
	}

	localFileName := localPath(relativeFileName, opts.Subfolder)
	var referenceFileName string
	if opts.Gold {
		referenceFileName = goldPath(relativeFileName, opts.Subfolder)
	} else {
		referenceFileName = referencePath(relativeFileName, opts.Subfolder)
	}
	if !opts.Expected {
		expected = getExpectedContent(referenceFileName)
	}
	if _, err := os.Stat(localFileName); err == nil {
		if err := os.Remove(localFileName); err != nil {
			t.Fatal(fmt.Errorf("failed to remove the local baseline file %s: %w", localFileName, err))
		}
	}
	if actual != expected {
		if err := os.MkdirAll(filepath.Dir(localFileName), 0755); err != nil {
			t.Fatal(fmt.Errorf("failed to create directories for the local baseline file %s: %w", localFileName, err))
		}
		if actual == NoContent {
			if err := os.WriteFile(localFileName+".delete", []byte{}, 0644); err != nil {
				t.Fatal(fmt.Errorf("failed to write the local baseline file %s: %w", localFileName+".delete", err))
			}
		} else if err := os.WriteFile(localFileName, []byte(actual), 0644); err != nil {
			t.Fatal(fmt.Errorf("failed to write the local baseline file %s: %w", localFileName, err))
		}

		if _, err := os.Stat(referenceFileName); err != nil {
			t.Errorf("New baseline created at %s.", localFileName)
		} else {
			t.Errorf("The baseline file %s has changed. (Run `hereby baseline-accept` if the new baseline is correct.)", relativeFileName)
		}
	}
}

func getExpectedContent(referenceFileName string) string {
	expected := NoContent
	if content, err := os.ReadFile(referenceFileName); err == nil {
		expected = string(content)
	}
	return expected
}

func localPath(fileName string, subfolder string) string {
	return filepath.Join(repo.TestDataPath, "baselines", "local", subfolder, fileName)
}

func referencePath(fileName string, subfolder string) string {
	return filepath.Join(repo.TestDataPath, "baselines", "reference", subfolder, fileName)
}

func goldPath(fileName string, subfolder string) string {
	return filepath.Join(repo.TestDataPath, "baselines", "gold", subfolder, fileName)
}
