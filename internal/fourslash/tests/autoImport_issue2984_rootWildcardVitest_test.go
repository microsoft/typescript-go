package fourslash_test

import (
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/fourslash"
	"github.com/microsoft/typescript-go/internal/ls"
	"github.com/microsoft/typescript-go/internal/ls/lsutil"
	"github.com/microsoft/typescript-go/internal/modulespecifiers"
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
	f.GoToMarker(t, "")

	verifyTopSpecifier := func(pref modulespecifiers.ImportModuleSpecifierPreference, expected string) {
		completions := f.GetCompletions(t, &lsutil.UserPreferences{ImportModuleSpecifierPreference: pref})
		var moduleSpecifiers []string
		for _, item := range completions.Items {
			if item.Label != "entity" || item.SortText == nil || *item.SortText != string(ls.SortTextAutoImportSuggestions) {
				continue
			}
			if item.Data == nil || item.Data.AutoImport == nil {
				continue
			}
			moduleSpecifiers = append(moduleSpecifiers, item.Data.AutoImport.ModuleSpecifier)
		}
		if len(moduleSpecifiers) == 0 {
			t.Fatalf("No auto-import completion specifier found for 'entity' with preference %q", pref)
		}
		if moduleSpecifiers[0] != expected {
			t.Fatalf("Unexpected first auto-import module specifier for preference %q.\nExpected: %s\nActual: %s\nAll: %v", pref, expected, moduleSpecifiers[0], moduleSpecifiers)
		}
	}

	verifyTopSpecifier(modulespecifiers.ImportModuleSpecifierPreferenceShortest, "#/domain/entities/entity.js")
	verifyTopSpecifier(modulespecifiers.ImportModuleSpecifierPreferenceNonRelative, "#/domain/entities/entity.js")
}

func TestAutoImport_issue2984_rootWildcardNotAllowedInNode16(t *testing.T) {
	t.Parallel()
	defer testutil.RecoverAndFail(t, "Panic on fourslash test")
	const content = `// @Filename: /tsconfig.json
{
  "compilerOptions": {
    "module": "node16",
    "moduleResolution": "node16"
  }
}
// @Filename: /package.json
{
  "imports": {
    "#/*": "./src/*"
  }
}
// @Filename: /src/domain/entities/entity.ts
export const entity = 1;
// @Filename: /deep/feature/path/consumer.ts
entit/**/`
	f, done := fourslash.NewFourslash(t, nil /*capabilities*/, content)
	defer done()
	f.GoToMarker(t, "")

	completions := f.GetCompletions(t, &lsutil.UserPreferences{
		ImportModuleSpecifierPreference: modulespecifiers.ImportModuleSpecifierPreferenceShortest,
	})
	var moduleSpecifiers []string
	for _, item := range completions.Items {
		if item.Label != "entity" || item.SortText == nil || *item.SortText != string(ls.SortTextAutoImportSuggestions) {
			continue
		}
		if item.Data == nil || item.Data.AutoImport == nil {
			continue
		}
		moduleSpecifiers = append(moduleSpecifiers, item.Data.AutoImport.ModuleSpecifier)
	}
	if len(moduleSpecifiers) == 0 {
		t.Fatalf("No auto-import completion specifier found for 'entity'")
	}
	if strings.HasPrefix(moduleSpecifiers[0], "#/") {
		t.Fatalf("Expected root wildcard imports key to be ignored outside nodenext/bundler, got %q", moduleSpecifiers[0])
	}
}
