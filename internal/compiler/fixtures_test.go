package compiler

import (
	"os"

	"github.com/microsoft/typescript-go/internal/testutil"
)

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

var (
	// Test binaries always start in the package directory; grab this early.
	cwd = must(os.Getwd())

	emptyFile = testutil.NewFileFixtureFromString("empty.ts", "empty.ts", "")
	checkerTs = testutil.NewFileFixtureFromFile("checker.ts", []string{cwd, "../../_submodules/TypeScript/src/compiler/checker.ts"})
	// largeTsxFile = testutil.NewFixture("jsxComplexSignatureHasApplicabilityError.tsx", []string{cwd, "../../_submodules/TypeScript/tests/cases/compiler/jsxComplexSignatureHasApplicabilityError.tsx"})

	benchFixtures = []testutil.FileFixture{
		emptyFile,
		checkerTs,
		// This crashes in bind.
		// largeTsxFile,
	}
)
