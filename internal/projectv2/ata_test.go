package projectv2_test

import (
	"context"
	"slices"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
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
		session.DidOpenFile(context.Background(), uri, 1, content, lsproto.LanguageKindJavaScript)
		session.WaitForBackgroundTasks()
		ls, err := session.GetLanguageService(context.Background(), uri)
		assert.NilError(t, err)
		// Verify the local config.js file is included in the program
		program := ls.GetProgram()
		assert.Assert(t, program != nil)
		configFile := program.GetSourceFile("/user/username/projects/project/config.js")
		assert.Assert(t, configFile != nil, "local config.js should be included")

		// Verify that only types-registry was installed (no @types/config since it's a local module)
		npmCalls := utils.NpmExecutor().NpmInstallCalls()
		assert.Equal(t, len(npmCalls), 1)
		assert.Equal(t, npmCalls[0].Args[2], "types-registry@latest")
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

		session.DidOpenFile(context.Background(), lsproto.DocumentUri("file:///user/username/projects/project/app.js"), 1, files["/user/username/projects/project/app.js"].(string), lsproto.LanguageKindJavaScript)
		session.WaitForBackgroundTasks()
		npmCalls := utils.NpmExecutor().NpmInstallCalls()
		assert.Equal(t, len(npmCalls), 2)
		assert.Equal(t, npmCalls[0].Cwd, projectv2testutil.TestTypingsLocation)
		assert.Equal(t, npmCalls[0].Args[2], "types-registry@latest")
		assert.Equal(t, npmCalls[1].Cwd, projectv2testutil.TestTypingsLocation)
		assert.Assert(t, slices.Contains(npmCalls[1].Args, "@types/jquery@latest"))
	})

	t.Run("inferred projects", func(t *testing.T) {
		t.Parallel()

		files := map[string]any{
			"/user/username/projects/project/app.js": ``,
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

		session.DidOpenFile(context.Background(), lsproto.DocumentUri("file:///user/username/projects/project/app.js"), 1, files["/user/username/projects/project/app.js"].(string), lsproto.LanguageKindJavaScript)
		session.WaitForBackgroundTasks()
		// Check that npm install was called twice
		calls := utils.NpmExecutor().NpmInstallCalls()
		assert.Equal(t, 2, len(calls), "Expected exactly 2 npm install calls")
		assert.Equal(t, calls[0].Cwd, projectv2testutil.TestTypingsLocation)
		assert.DeepEqual(t, calls[0].Args, []string{"install", "--ignore-scripts", "types-registry@latest"})
		assert.Equal(t, calls[1].Cwd, projectv2testutil.TestTypingsLocation)
		assert.Equal(t, calls[1].Args[2], "@types/jquery@latest")

		// Verify the types file was installed
		ls, err := session.GetLanguageService(context.Background(), lsproto.DocumentUri("file:///user/username/projects/project/app.js"))
		assert.NilError(t, err)
		program := ls.GetProgram()
		jqueryTypesFile := program.GetSourceFile(projectv2testutil.TestTypingsLocation + "/node_modules/@types/jquery/index.d.ts")
		assert.Assert(t, jqueryTypesFile != nil, "jquery types should be installed")
	})

	t.Run("type acquisition with disableFilenameBasedTypeAcquisition:true", func(t *testing.T) {
		t.Parallel()

		files := map[string]any{
			"/user/username/projects/project/jquery.js": ``,
			"/user/username/projects/project/tsconfig.json": `{
				"compilerOptions": { "allowJs": true },
				"typeAcquisition": { "enable": true, "disableFilenameBasedTypeAcquisition": true }
			}`,
		}

		session, utils := projectv2testutil.SetupWithTypingsInstaller(files, &projectv2testutil.TestTypingsInstallerOptions{
			TypesRegistry: []string{"jquery"},
		})

		// Should only get types-registry install, no jquery install since filename-based acquisition is disabled
		session.DidOpenFile(context.Background(), lsproto.DocumentUri("file:///user/username/projects/project/jquery.js"), 1, files["/user/username/projects/project/jquery.js"].(string), lsproto.LanguageKindJavaScript)
		session.WaitForBackgroundTasks()

		// Check that npm install was called once (only types-registry)
		calls := utils.NpmExecutor().NpmInstallCalls()
		assert.Equal(t, 1, len(calls), "Expected exactly 1 npm install call")
		assert.Equal(t, calls[0].Cwd, projectv2testutil.TestTypingsLocation)
		assert.DeepEqual(t, calls[0].Args, []string{"install", "--ignore-scripts", "types-registry@latest"})
	})

	t.Run("discover from node_modules", func(t *testing.T) {
		t.Parallel()

		files := map[string]any{
			"/user/username/projects/project/app.js": "",
			"/user/username/projects/project/package.json": `{
			    "dependencies": {
					"jquery": "1.0.0"
				}
			}`,
			"/user/username/projects/project/jsconfig.json":                           `{}`,
			"/user/username/projects/project/node_modules/commander/index.js":         "",
			"/user/username/projects/project/node_modules/commander/package.json":     `{ "name": "commander" }`,
			"/user/username/projects/project/node_modules/jquery/index.js":            "",
			"/user/username/projects/project/node_modules/jquery/package.json":        `{ "name": "jquery" }`,
			"/user/username/projects/project/node_modules/jquery/nested/package.json": `{ "name": "nested" }`,
		}

		session, utils := projectv2testutil.SetupWithTypingsInstaller(files, &projectv2testutil.TestTypingsInstallerOptions{
			TypesRegistry: []string{"nested", "commander"},
			PackageToFile: map[string]string{
				"jquery": "declare const jquery: { x: number }",
			},
		})

		session.DidOpenFile(context.Background(), lsproto.DocumentUri("file:///user/username/projects/project/app.js"), 1, files["/user/username/projects/project/app.js"].(string), lsproto.LanguageKindJavaScript)
		session.WaitForBackgroundTasks()

		// Check that npm install was called twice
		calls := utils.NpmExecutor().NpmInstallCalls()
		assert.Equal(t, 2, len(calls), "Expected exactly 2 npm install calls")
		assert.Equal(t, calls[0].Cwd, projectv2testutil.TestTypingsLocation)
		assert.DeepEqual(t, calls[0].Args, []string{"install", "--ignore-scripts", "types-registry@latest"})
		assert.Equal(t, calls[1].Cwd, projectv2testutil.TestTypingsLocation)
		assert.Equal(t, calls[1].Args[2], "@types/jquery@latest")
	})

	t.Run("discover from bower_components", func(t *testing.T) {
		t.Parallel()

		files := map[string]any{
			"/user/username/projects/project/app.js":                             ``,
			"/user/username/projects/project/jsconfig.json":                      `{}`,
			"/user/username/projects/project/bower_components/jquery/index.js":   "",
			"/user/username/projects/project/bower_components/jquery/bower.json": `{ "name": "jquery" }`,
		}

		session, utils := projectv2testutil.SetupWithTypingsInstaller(files, &projectv2testutil.TestTypingsInstallerOptions{
			PackageToFile: map[string]string{
				"jquery": "declare const jquery: { x: number }",
			},
		})

		session.DidOpenFile(context.Background(), lsproto.DocumentUri("file:///user/username/projects/project/app.js"), 1, files["/user/username/projects/project/app.js"].(string), lsproto.LanguageKindJavaScript)
		session.WaitForBackgroundTasks()

		// Check that npm install was called twice
		calls := utils.NpmExecutor().NpmInstallCalls()
		assert.Equal(t, 2, len(calls), "Expected exactly 2 npm install calls")
		assert.Equal(t, calls[0].Cwd, projectv2testutil.TestTypingsLocation)
		assert.DeepEqual(t, calls[0].Args, []string{"install", "--ignore-scripts", "types-registry@latest"})
		assert.Equal(t, calls[1].Cwd, projectv2testutil.TestTypingsLocation)
		assert.Equal(t, calls[1].Args[2], "@types/jquery@latest")

		// Verify the types file was installed
		ls, err := session.GetLanguageService(context.Background(), lsproto.DocumentUri("file:///user/username/projects/project/app.js"))
		assert.NilError(t, err)
		jqueryTypesFile := ls.GetProgram().GetSourceFile(projectv2testutil.TestTypingsLocation + "/node_modules/@types/jquery/index.d.ts")
		assert.Assert(t, jqueryTypesFile != nil, "jquery types should be installed")
	})

	t.Run("discover from bower.json", func(t *testing.T) {
		t.Parallel()

		files := map[string]any{
			"/user/username/projects/project/app.js":        ``,
			"/user/username/projects/project/jsconfig.json": `{}`,
			"/user/username/projects/project/bower.json": `{
				"dependencies": {
                    "jquery": "^3.1.0"
                }
			}`,
		}

		session, utils := projectv2testutil.SetupWithTypingsInstaller(files, &projectv2testutil.TestTypingsInstallerOptions{
			PackageToFile: map[string]string{
				"jquery": "declare const jquery: { x: number }",
			},
		})

		session.DidOpenFile(context.Background(), lsproto.DocumentUri("file:///user/username/projects/project/app.js"), 1, files["/user/username/projects/project/app.js"].(string), lsproto.LanguageKindJavaScript)
		session.WaitForBackgroundTasks()

		// Check that npm install was called twice
		calls := utils.NpmExecutor().NpmInstallCalls()
		assert.Equal(t, 2, len(calls), "Expected exactly 2 npm install calls")
		assert.Equal(t, calls[0].Cwd, projectv2testutil.TestTypingsLocation)
		assert.DeepEqual(t, calls[0].Args, []string{"install", "--ignore-scripts", "types-registry@latest"})
		assert.Equal(t, calls[1].Cwd, projectv2testutil.TestTypingsLocation)
		assert.Equal(t, calls[1].Args[2], "@types/jquery@latest")

		// Verify the types file was installed
		ls, err := session.GetLanguageService(context.Background(), lsproto.DocumentUri("file:///user/username/projects/project/app.js"))
		assert.NilError(t, err)
		jqueryTypesFile := ls.GetProgram().GetSourceFile(projectv2testutil.TestTypingsLocation + "/node_modules/@types/jquery/index.d.ts")
		assert.Assert(t, jqueryTypesFile != nil, "jquery types should be installed")
	})

	t.Run("should install typings for unresolved imports", func(t *testing.T) {
		t.Parallel()

		files := map[string]any{
			"/user/username/projects/project/app.js": `
				import * as fs from "fs";
                import * as commander from "commander";
                import * as component from "@ember/component";
			`,
		}

		session, utils := projectv2testutil.SetupWithTypingsInstaller(files, &projectv2testutil.TestTypingsInstallerOptions{
			PackageToFile: map[string]string{
				"node":             "export let node: number",
				"commander":        "export let commander: number",
				"ember__component": "export let ember__component: number",
			},
		})

		session.DidOpenFile(context.Background(), lsproto.DocumentUri("file:///user/username/projects/project/app.js"), 1, files["/user/username/projects/project/app.js"].(string), lsproto.LanguageKindJavaScript)
		session.WaitForBackgroundTasks()

		// Check that npm install was called twice
		calls := utils.NpmExecutor().NpmInstallCalls()
		assert.Equal(t, 2, len(calls), "Expected exactly 2 npm install calls")
		assert.Equal(t, calls[0].Cwd, projectv2testutil.TestTypingsLocation)
		assert.DeepEqual(t, calls[0].Args, []string{"install", "--ignore-scripts", "types-registry@latest"})

		// The second call should install all three packages at once
		assert.Equal(t, calls[1].Cwd, projectv2testutil.TestTypingsLocation)
		assert.Equal(t, calls[1].Args[0], "install")
		assert.Equal(t, calls[1].Args[1], "--ignore-scripts")
		// Check that all three packages are in the install command
		installArgs := calls[1].Args
		assert.Assert(t, slices.Contains(installArgs, "@types/ember__component@latest"))
		assert.Assert(t, slices.Contains(installArgs, "@types/commander@latest"))
		assert.Assert(t, slices.Contains(installArgs, "@types/node@latest"))

		// Verify the types files were installed
		ls, err := session.GetLanguageService(context.Background(), lsproto.DocumentUri("file:///user/username/projects/project/app.js"))
		assert.NilError(t, err)
		program := ls.GetProgram()
		nodeTypesFile := program.GetSourceFile(projectv2testutil.TestTypingsLocation + "/node_modules/@types/node/index.d.ts")
		assert.Assert(t, nodeTypesFile != nil, "node types should be installed")
		commanderTypesFile := program.GetSourceFile(projectv2testutil.TestTypingsLocation + "/node_modules/@types/commander/index.d.ts")
		assert.Assert(t, commanderTypesFile != nil, "commander types should be installed")
		emberComponentTypesFile := program.GetSourceFile(projectv2testutil.TestTypingsLocation + "/node_modules/@types/ember__component/index.d.ts")
		assert.Assert(t, emberComponentTypesFile != nil, "ember__component types should be installed")
	})
}
