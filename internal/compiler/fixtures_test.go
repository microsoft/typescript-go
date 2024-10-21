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

	srcCompilerCheckerTS = testutil.NewFixture(cwd, "../../_submodules/TypeScript/src/compiler/checker.ts")
)
