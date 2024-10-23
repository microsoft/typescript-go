package baseline

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/microsoft/typescript-go/internal/repo"
)

type BaselineOptions struct {
	Subfolder string
}

var NoContent string

func init() {
	NoContent = "<no content>"
}

func RunBaseline(fileName string, actual string, opts ...BaselineOptions) error {
	if len(opts) > 1 {
		panic("Too many options")
	}
	var opt BaselineOptions
	if len(opts) == 1 {
		opt = opts[0]
	}
	if actual == "" {
		panic("The generated content was \"\". Return 'baseline.NoContent' if no baselining is required.")
	}

	return writeComparison(actual, fileName, opt)
}

func writeComparison(actual string, relativeFileName string, opts BaselineOptions) error {
	localFileName := localPath(relativeFileName, opts.Subfolder)
	referenceFileName := referencePath(relativeFileName, opts.Subfolder)
	expected := getExpectedContent(relativeFileName, opts)
	if _, err := os.Stat(localFileName); err == nil {
		os.Remove(localFileName)
	}
	if actual != expected {
		os.MkdirAll(filepath.Dir(localFileName), 0755)
		if actual == NoContent {
			os.WriteFile(localFileName+".delete", []byte{}, 0644)
		} else {
			os.WriteFile(localFileName, []byte(actual), 0644)
		}

		if _, err := os.Stat(referenceFileName); err == nil {
			return fmt.Errorf("New baseline created at %s.", localFileName)
		} else {
			return fmt.Errorf("The baseline file %s has changed. (Run `hereby baseline-accept` if the new baseline is correct.)", relativeFileName)
		}
	}
	return nil
}

func getExpectedContent(relativeFileName string, opts BaselineOptions) string {
	refFileName := referencePath(relativeFileName, opts.Subfolder)
	expected := NoContent
	content, err := os.ReadFile(refFileName)
	if err == nil {
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
