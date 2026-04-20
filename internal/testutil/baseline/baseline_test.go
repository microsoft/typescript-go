package baseline

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSubmoduleAcceptedFilesExist(t *testing.T) {
	t.Parallel()
	for name := range submoduleAcceptedFileNames().Keys() {
		if _, err := os.Stat(filepath.Join(referenceRoot, "submoduleAccepted", name)); err != nil {
			t.Errorf("submoduleAccepted.txt references %q, but the baseline file does not exist", name)
		}
	}
}

func TestSubmoduleTriagedFilesExist(t *testing.T) {
	t.Parallel()
	for name := range submoduleTriagedFileNames().Keys() {
		if _, err := os.Stat(filepath.Join(referenceRoot, "submoduleTriaged", name)); err != nil {
			t.Errorf("submoduleTriaged.txt references %q, but the baseline file does not exist", name)
		}
	}
}

func TestNoFileInBothAcceptedAndTriaged(t *testing.T) {
	t.Parallel()
	accepted := submoduleAcceptedFileNames()
	triaged := submoduleTriagedFileNames()
	for name := range accepted.Keys() {
		if triaged.Has(name) {
			t.Errorf("file %q is in both submoduleAccepted.txt and submoduleTriaged.txt; it should only be in one", name)
		}
	}
}
