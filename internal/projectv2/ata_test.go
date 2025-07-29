package projectv2_test

import (
	"context"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/projectv2"
	"github.com/microsoft/typescript-go/internal/testutil/projectv2testutil"
	"gotest.tools/v3/assert"
)

func TestATA(t *testing.T) {
	t.Parallel()
	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	t.Run("local module should not be picked up", func(t *testing.T) {
		t.Parallel()
		files := map[string]any{
			"/user/username/projects/project/app.js":    `const c = require('./config');`,
			"/user/username/projects/project/config.js": `export let x = 1`,
			"/user/username/projects/project/jsconfig.json": `{
					"compilerOptions": { "moduleResolution": "commonjs" },
					"typeAcquisition": { "enable": true }
			}`,
		}

		testOptions := &projectv2testutil.TestTypingsInstallerOptions{
			TypesRegistry: []string{"config"},
		}

		session, utils := projectv2testutil.SetupWithTypingsInstaller(files, testOptions)
		uri := lsproto.DocumentUri("file:///user/username/projects/project/app.js")
		content := files["/user/username/projects/project/app.js"].(string)

		// Open the file
		awaitNpmInstall := utils.ExpectNpmInstallCalls(1) // types-registry
		session.DidOpenFile(context.Background(), uri, 1, content, lsproto.LanguageKindJavaScript)
		awaitNpmInstall()

		// Get the snapshot and verify the project
		snapshot, release := session.Snapshot()
		defer release()

		projects := snapshot.ProjectCollection.Projects()
		assert.Equal(t, len(projects), 1)

		project := projects[0]
		assert.Equal(t, project.Kind, projectv2.KindConfigured)

		// Verify the local config.js file is included in the program
		program := project.Program
		assert.Assert(t, program != nil)
		configFile := program.GetSourceFile("/user/username/projects/project/config.js")
		assert.Assert(t, configFile != nil, "local config.js should be included")
	})

	t.Run("configured projects", func(t *testing.T) {
		t.Parallel()

		files := map[string]any{
			"/user/username/projects/project/app.js": ``,
			"/user/username/projects/project/tsconfig.json": `{
				"compilerOptions": { "allowJs": true },
				"typeAcquisition": { "enable": true },
			}`,
			"/user/username/projects/project/package.json": `{
				"name": "test",
				"dependencies": {
					"jquery": "^3.1.0"
				}
			}`,
		}

		session, utils := projectv2testutil.SetupWithTypingsInstaller(files, &projectv2testutil.TestTypingsInstallerOptions{
			PackageToFile: map[string]string{
				"jquery": `declare const $: { x: number }`,
			},
		})

		awaitNpmInstall := utils.ExpectNpmInstallCalls(2)
		session.DidOpenFile(context.Background(), lsproto.DocumentUri("file:///user/username/projects/project/app.js"), 1, files["/user/username/projects/project/app.js"].(string), lsproto.LanguageKindJavaScript)
		snapshot, release := session.Snapshot()
		defer release()

		projects := snapshot.ProjectCollection.Projects()
		assert.Equal(t, len(projects), 1)
		npmInstallCalls := awaitNpmInstall()
		assert.Equal(t, npmInstallCalls[0].Cwd, projectv2testutil.TestTypingsLocation)
		assert.DeepEqual(t, npmInstallCalls[0].NpmInstallArgs, []string{"install", "--ignore-scripts", "types-registry@latest"})
		assert.Equal(t, npmInstallCalls[1].Cwd, projectv2testutil.TestTypingsLocation)
		assert.Equal(t, npmInstallCalls[1].NpmInstallArgs[2], "@types/jquery@latest")
	})
}
