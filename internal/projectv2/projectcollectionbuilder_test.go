package projectv2_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil/projectv2testutil"
	"github.com/microsoft/typescript-go/internal/tspath"
	"gotest.tools/v3/assert"
)

func TestProjectCollectionBuilder(t *testing.T) {
	t.Parallel()

	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	t.Run("when project found is solution referencing default project directly", func(t *testing.T) {
		t.Parallel()
		files := filesForSolutionConfigFile([]string{"./tsconfig-src.json"}, "", nil)
		session := projectv2testutil.Setup(files)
		uri := lsproto.DocumentUri("file:///user/username/projects/myproject/src/main.ts")
		content := files["/user/username/projects/myproject/src/main.ts"].(string)

		// Ensure configured project is found for open file
		session.DidOpenFile(context.Background(), uri, 1, content, lsproto.LanguageKindTypeScript)
		snapAfterOpen := session.Snapshot()
		assert.Equal(t, len(snapAfterOpen.ProjectCollection.Projects()), 1)
		assert.Assert(t, snapAfterOpen.ProjectCollection.ConfiguredProject(tspath.Path("/user/username/projects/myproject/tsconfig-src.json")) != nil)

		// Ensure request can use existing snapshot
		_, err := session.GetLanguageService(context.Background(), uri)
		assert.NilError(t, err)
		assert.Equal(t, session.Snapshot(), snapAfterOpen)

		// Close the file and open a different one
		session.DidCloseFile(context.Background(), uri)
		dummyUri := lsproto.DocumentUri("file:///user/username/workspaces/dummy/dummy.ts")
		session.DidOpenFile(context.Background(), dummyUri, 1, "const x = 1;", lsproto.LanguageKindTypeScript)
		assert.Equal(t, len(session.Snapshot().ProjectCollection.Projects()), 1)
		assert.Assert(t, session.Snapshot().ProjectCollection.InferredProject() != nil)

		// Config files should have been released
		assert.Assert(t, session.Snapshot().ConfigFileRegistry.GetConfig("/user/username/projects/myproject/tsconfig.json") == nil)
		assert.Assert(t, session.Snapshot().ConfigFileRegistry.GetConfig("/user/username/projects/myproject/tsconfig-src.json") == nil)
	})
}

func filesForSolutionConfigFile(solutionRefs []string, compilerOptions string, ownFiles []string) map[string]any {
	var compilerOptionsStr string
	if compilerOptions != "" {
		compilerOptionsStr = fmt.Sprintf(`"compilerOptions": {
			%s
		},`, compilerOptions)
	}
	var ownFilesStr string
	if len(ownFiles) > 0 {
		ownFilesStr = strings.Join(ownFiles, ",")
	}
	files := map[string]any{
		"/user/username/projects/myproject/tsconfig.json": fmt.Sprintf(`{
			%s
			"files": [%s],
			"references": [
				%s
			]
		}`, compilerOptionsStr, ownFilesStr, strings.Join(core.Map(solutionRefs, func(ref string) string {
			return fmt.Sprintf(`{ "path": "%s" }`, ref)
		}), ",")),
		"/user/username/projects/myproject/tsconfig-src.json": `{
			"compilerOptions": {
				"composite": true,
				"outDir": "./target",
			},
			"include": ["./src/**/*"]
		}`,
		"/user/username/projects/myproject/src/main.ts": `
			import { foo } from './helpers/functions';
			export { foo };`,
		"/user/username/projects/myproject/src/helpers/functions.ts": `export const foo = 1;`,
	}
	return files
}
