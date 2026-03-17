package fourslash_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/testutil"
)

func TestAutoImport_issue2984_rootWildcardVitest(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /tsconfig.json
{
  "compilerOptions": {
    "module": "node20",
    "moduleResolution": "nodenext",
    "rootDir": "./",
    "outDir": "build"
  }
}
// @Filename: /package.json
{
  "imports": {
    "#/*": {
      "vitest": "./src/*",
      "types": "./src/*",
      "node": "./build/*",
      "default": "./src/*"
    }
  }
}
// @Filename: /src/domain/entities/entity.ts
export const entity = 1;
// @Filename: /feature/very/deep/path/consumer.ts
entit/**/`

	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()

	assertBest := func(prefs *lsutil.UserPreferences, prefName string) {
		f.GoToMarker(t, "")
		completions := f.GetCompletions(t, prefs)
		if completions == nil {
			t.Fatalf("%s: expected completions list", prefName)
		}
		var entitySpecifiers []string
		for _, item := range completions.Items {
			if item.Label != "entity" || item.Data == nil || item.Data.AutoImport == nil {
				continue
			}
			entitySpecifiers = append(entitySpecifiers, item.Data.AutoImport.ModuleSpecifier)
		}
		if len(entitySpecifiers) == 0 {
			t.Fatalf("%s: expected auto-import completion for entity", prefName)
		}
		t.Logf("%s entity specifiers: %v", prefName, entitySpecifiers)
		if entitySpecifiers[0] != "#/domain/entities/entity.js" {
			t.Fatalf("%s: expected top module specifier %q, got %q", prefName, "#/domain/entities/entity.js", entitySpecifiers[0])
		}
	}

	assertBest(&lsutil.UserPreferences{ImportModuleSpecifierPreference: "shortest"}, "shortest")
	assertBest(&lsutil.UserPreferences{ImportModuleSpecifierPreference: "non-relative"}, "non-relative")
}
