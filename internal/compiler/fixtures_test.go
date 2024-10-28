package compiler

import (
	"path/filepath"

	"github.com/microsoft/typescript-go/internal/repo"
	"github.com/microsoft/typescript-go/internal/testutil"
)

var benchFixtures = []testutil.FileFixture{
	testutil.NewFileFixtureFromString("empty.ts", "empty.ts", ""),
	testutil.NewFileFixtureFromFile("checker.ts", filepath.Join(repo.TypeScriptSubmodulePath, "src/compiler/checker.ts")),
	testutil.NewFileFixtureFromFile("dom.generated.d.ts", filepath.Join(repo.TypeScriptSubmodulePath, "src/lib/dom.generated.d.ts")),
	testutil.NewFileFixtureFromFile("Herebyfile.mjs", filepath.Join(repo.TypeScriptSubmodulePath, "Herebyfile.mjs")),
	testutil.NewFileFixtureFromFile("jsxComplexSignatureHasApplicabilityError.tsx", filepath.Join(repo.TypeScriptSubmodulePath, "tests/cases/compiler/jsxComplexSignatureHasApplicabilityError.tsx")),
}
