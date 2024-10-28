package compiler

import (
	"github.com/microsoft/typescript-go/internal/repo"
	"github.com/microsoft/typescript-go/internal/testutil"
)

var benchFixtures = []testutil.FileFixture{
	testutil.NewFileFixtureFromString("empty.ts", "empty.ts", ""),
	testutil.NewFileFixtureFromFile("checker.ts", []string{repo.TypeScriptSubmodulePath, "src/compiler/checker.ts"}),
	testutil.NewFileFixtureFromFile("dom.generated.d.ts", []string{repo.TypeScriptSubmodulePath, "src/lib/dom.generated.d.ts"}),
	testutil.NewFileFixtureFromFile("Herebyfile.mjs", []string{repo.TypeScriptSubmodulePath, "Herebyfile.mjs"}),
	// This crashes in bind.
	// testutil.NewFixture("jsxComplexSignatureHasApplicabilityError.tsx", []string{cwd, "../../_submodules/TypeScript/tests/cases/compiler/jsxComplexSignatureHasApplicabilityError.tsx"}),
}
