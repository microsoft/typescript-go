package lspservertests

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/execute/tsctests"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil/lsptestutil"
	"github.com/microsoft/typescript-go/internal/testutil/stringtestutil"
)

func TestFindAllReferences(t *testing.T) {
	t.Parallel()

	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	testCases := slices.Concat(
		getFindAllRefsTestCasesForDefaultProjects(),
		getFindAllRefsTestcasesForRootOfReferencedProject(),
		[]*lspServerTest{
			{
				subScenario: "finding local reference doesnt load ancestor sibling projects",
				files: func() map[string]any {
					return getFindAllRefsFileMapForLocalness(false)
				},
				test: func(server *testServer) {
					programFile := "/user/username/projects/solution/compiler/program.ts"
					server.openFile(programFile, lsproto.LanguageKindTypeScript)

					// Find all references for getSourceFile
					// Shouldnt load more projects
					server.baselineReferences(programFile, lsptestutil.PositionToLineAndCharacter(programFile, server.content(programFile), "getSourceFile", 1))

					// Find all references for getSourceFiles
					// Should load more projects
					server.baselineReferences(programFile, lsptestutil.PositionToLineAndCharacter(programFile, server.content(programFile), "getSourceFiles", 0))
				},
			},
			{
				subScenario: "disableSolutionSearching solution and siblings are not loaded",
				files: func() map[string]any {
					return getFindAllRefsFileMapForLocalness(true)
				},
				test: func(server *testServer) {
					programFile := "/user/username/projects/solution/compiler/program.ts"
					server.openFile(programFile, lsproto.LanguageKindTypeScript)

					// Find all references
					// No new solutions/projects loaded
					server.baselineReferences(programFile, lsptestutil.PositionToLineAndCharacter(programFile, server.content(programFile), "getSourceFiles", 0))
				},
			},
			{
				subScenario: "finding references in overlapping projects",
				files: func() map[string]any {
					return map[string]any{
						"/user/username/projects/solution/tsconfig.json": stringtestutil.Dedent(`
							{
								"files": [],
								"include": [],
								"references": [
									{ "path": "./a" },
									{ "path": "./b" },
									{ "path": "./c" },
									{ "path": "./d" },
								],
							}`),
						"/user/username/projects/solution/a/tsconfig.json": stringtestutil.Dedent(`
							{
								"compilerOptions": {
									"composite": true,
								},
								"files": ["./index.ts"]
							}`),
						"/user/username/projects/solution/a/index.ts": stringtestutil.Dedent(`
							export interface I {
								M(): void;
							}`),
						"/user/username/projects/solution/b/tsconfig.json": stringtestutil.Dedent(`
							{
								"compilerOptions": {
									"composite": true
								},
								"files": ["./index.ts"],
								"references": [
									{ "path": "../a" },
								],
							}`),
						"/user/username/projects/solution/b/index.ts": stringtestutil.Dedent(`
							import { I } from "../a";
							export class B implements I {
								M() {}
							}`),
						"/user/username/projects/solution/c/tsconfig.json": stringtestutil.Dedent(`
							{
								"compilerOptions": {
									"composite": true
								},
								"files": ["./index.ts"],
								"references": [
									{ "path": "../b" },
								],
							}`),
						"/user/username/projects/solution/c/index.ts": stringtestutil.Dedent(`
							import { I } from "../a";
							import { B } from "../b";
							export const C: I = new B();`),
						"/user/username/projects/solution/d/tsconfig.json": stringtestutil.Dedent(`
							{
								"compilerOptions": {
									"composite": true
								},
								"files": ["./index.ts"],
								"references": [
									{ "path": "../c" },
								],
							}`),
						"/user/username/projects/solution/d/index.ts": stringtestutil.Dedent(`
							import { I } from "../a";
							import { C } from "../c";
							export const D: I = C;`),
					}
				},
				test: func(server *testServer) {
					bFile := "/user/username/projects/solution/b/index.ts"
					server.openFile(bFile, lsproto.LanguageKindTypeScript)

					// The first search will trigger project loads
					server.baselineReferences(bFile, lsptestutil.PositionToLineAndCharacter(bFile, server.content(bFile), "I", 1))

					// The second search starts with the projects already loaded
					// Formerly, this would search some projects multiple times
					server.baselineReferences(bFile, lsptestutil.PositionToLineAndCharacter(bFile, server.content(bFile), "I", 1))
				},
			},
			{
				subScenario: "files from two projects are open and one project references",
				files: func() map[string]any {
					files := map[string]any{}
					applyPackageConfigAndFile := func(packageName string, references []string, disableReferencedProjectLoad bool) {
						files[fmt.Sprintf("/user/username/projects/myproject/%s/src/file1.ts", packageName)] = fmt.Sprintf(`export const %sConst = 10;`, packageName)
						var extraStr string
						if disableReferencedProjectLoad {
							extraStr = `"disableReferencedProjectLoad": true,`
						}
						var referencesStr strings.Builder
						for _, ref := range references {
							referencesStr.WriteString(fmt.Sprintf(`{ "path": "../%s" },\n`, ref))
						}
						files[fmt.Sprintf("/user/username/projects/myproject/%s/tsconfig.json", packageName)] = stringtestutil.Dedent(fmt.Sprintf(`
							{
								"compilerOptions": {
									"composite": true,
									%s
								},
								%s
							}
						`, extraStr, referencesStr.String()))
					}
					applyPackageConfigAndFile("main", []string{"core", "indirect", "noCoreRef1", "indirectDisabledChildLoad1", "indirectDisabledChildLoad2", "refToCoreRef3", "indirectNoCoreRef"}, false)
					applyPackageConfigAndFile("core", nil, false)
					applyPackageConfigAndFile("noCoreRef1", nil, false)
					applyPackageConfigAndFile("indirect", []string{"coreRef1"}, false)
					applyPackageConfigAndFile("coreRef1", []string{"core"}, false)
					applyPackageConfigAndFile("indirectDisabledChildLoad1", []string{"coreRef2"}, true)
					applyPackageConfigAndFile("coreRef2", []string{"core"}, false)
					applyPackageConfigAndFile("indirectDisabledChildLoad2", []string{"coreRef3"}, true)
					applyPackageConfigAndFile("coreRef3", []string{"core"}, false)
					applyPackageConfigAndFile("refToCoreRef3", []string{"coreRef3"}, false)
					applyPackageConfigAndFile("indirectNoCoreRef", []string{"noCoreRef2"}, false)
					applyPackageConfigAndFile("noCoreRef2", nil, false)
					return files
				},
				test: func(server *testServer) {
					mainFile := "/user/username/projects/myproject/main/src/file1.ts"
					coreFile := "/user/username/projects/myproject/core/src/file1.ts"
					server.openFile(mainFile, lsproto.LanguageKindTypeScript)
					server.openFile(coreFile, lsproto.LanguageKindTypeScript)

					// Find all refs in coreFile
					server.baselineReferences(coreFile, lsptestutil.PositionToLineAndCharacter(coreFile, server.content(coreFile), "coreConst", 0))
				},
			},
			{
				subScenario: "does not try to open a file in a project that was updated and no longer has the file",
				files: func() map[string]any {
					return map[string]any{
						"/home/src/projects/project/packages/babel-loader/tsconfig.json": stringtestutil.Dedent(`
							{
								"compilerOptions": {
									"target": "ES2018",
									"module": "commonjs",
									"strict": true,
									"esModuleInterop": true,
									"composite": true,
									"rootDir": "src",
									"outDir": "dist"
								},
								"include": ["src"],
								"references": [{"path": "../core"}]
							}`),
						"/home/src/projects/project/packages/babel-loader/src/index.ts": stringtestutil.Dedent(`
							import type { Foo } from "../../core/src/index.js";`),
						"/home/src/projects/project/packages/core/tsconfig.json": stringtestutil.Dedent(`
							{
								"compilerOptions": {
									"target": "ES2018",
									"module": "commonjs",
									"strict": true,
									"esModuleInterop": true,
									"composite": true,
									"rootDir": "./src",
									"outDir": "./dist",
								},
								"include": ["./src"]
							}`),
						"/home/src/projects/project/packages/core/src/index.ts": stringtestutil.Dedent(`
							import { Bar } from "./loading-indicator.js";
							export type Foo = {};
							const bar: Bar = {
								prop: 0
							}`),
						"/home/src/projects/project/packages/core/src/loading-indicator.ts": stringtestutil.Dedent(`
							export interface Bar {
								prop: number;
							}
							const bar: Bar = {
								prop: 1
							}`),
					}
				},
				test: func(server *testServer) {
					// Open files in the two configured projects
					indexFile := "/home/src/projects/project/packages/babel-loader/src/index.ts"
					coreFile := "/home/src/projects/project/packages/core/src/index.ts"
					server.openFile(indexFile, lsproto.LanguageKindTypeScript)
					server.openFile(coreFile, lsproto.LanguageKindTypeScript)

					// Now change `babel-loader` project to no longer import `core` project
					server.changeFile(&lsproto.DidChangeTextDocumentParams{
						TextDocument: lsproto.VersionedTextDocumentIdentifier{
							Uri:     lsproto.DocumentUri("file://" + indexFile),
							Version: 2,
						},
						ContentChanges: []lsproto.TextDocumentContentChangePartialOrWholeDocument{
							{
								Partial: &lsproto.TextDocumentContentChangePartial{
									Range: lsproto.Range{
										Start: lsproto.Position{
											Line:      0,
											Character: 0,
										},
										End: lsproto.Position{
											Line:      0,
											Character: 0,
										},
									},
									Text: "// comment",
								},
							},
						},
					})

					// At this point, we haven't updated `babel-loader` project yet,
					// so `babel-loader` is still a containing project of `loading-indicator` file.
					// When calling find all references,
					// we shouldn't crash due to using outdated information on a file's containing projects.
					server.baselineReferences(coreFile, lsptestutil.PositionToLineAndCharacter(coreFile, server.content(coreFile), "prop", 0))
				},
			},
			getFindAllRefsTestCaseForSpecialLocalnessHandling(
				"when using arrow function assignment",
				`export const dog = () => { };`,
				`shared.dog();`,
				"dog",
			),
			getFindAllRefsTestCaseForSpecialLocalnessHandling(
				"when using arrow function as object literal property types",
				`export const foo = { bar: () => { } };`,
				`shared.foo.bar();`,
				"bar",
			),
			getFindAllRefsTestCaseForSpecialLocalnessHandling(
				"when using object literal property",
				`export const foo = {  baz: "BAZ" };`,
				`shared.foo.baz;`,
				"baz",
			),
			getFindAllRefsTestCaseForSpecialLocalnessHandling(
				"when using method of class expression",
				`export const foo = class { fly() {} };`,
				stringtestutil.Dedent(`
					const instance = new shared.foo();
					instance.fly();`),
				"fly",
			),
			getFindAllRefsTestCaseForSpecialLocalnessHandling(
				// when using arrow function as object literal property is loaded through indirect assignment with original declaration local to project is treated as local
				"when using arrow function as object literal property",
				stringtestutil.Dedent(`
					const local = { bar: () => { } };
					export const foo = local;`),
				`shared.foo.bar();`,
				"bar",
			),
			// Pre-loaded = A file from project B is already open when FindAllRefs is invoked
			// dRPL = Project A has disableReferencedProjectLoad
			// dSOPRR = Project A has disableSourceOfProjectReferenceRedirect
			// Map = The declaration map file b/lib/index.d.ts.map exists
			// B refs = files under directory b in which references are found (all scenarios find all references in a/index.ts)

			//                                               Pre-loaded |dRPL|dSOPRR|Map   |     B state      | Notes        | B refs              | Notes
			//                                               -----------+----+------+------+------------------+--------------+---------------------+---------------------------------------------------
			getTestcaseFindAllRefsWithDisableReferencedProjectLoad(true, true, true, true),     // Pre-loaded |              | index.ts, helper.ts | Via map and pre-loaded project
			getTestcaseFindAllRefsWithDisableReferencedProjectLoad(true, true, true, false),    // Pre-loaded |              | lib/index.d.ts      | Even though project is loaded
			getTestcaseFindAllRefsWithDisableReferencedProjectLoad(true, true, false, true),    // Pre-loaded |              | index.ts, helper.ts |
			getTestcaseFindAllRefsWithDisableReferencedProjectLoad(true, true, false, false),   // Pre-loaded |              | index.ts, helper.ts |
			getTestcaseFindAllRefsWithDisableReferencedProjectLoad(true, false, true, true),    // Pre-loaded |              | index.ts, helper.ts | Via map and pre-loaded project
			getTestcaseFindAllRefsWithDisableReferencedProjectLoad(true, false, true, false),   // Pre-loaded |              | lib/index.d.ts      | Even though project is loaded
			getTestcaseFindAllRefsWithDisableReferencedProjectLoad(true, false, false, true),   // Pre-loaded |              | index.ts, helper.ts |
			getTestcaseFindAllRefsWithDisableReferencedProjectLoad(true, false, false, false),  // Pre-loaded |              | index.ts, helper.ts |
			getTestcaseFindAllRefsWithDisableReferencedProjectLoad(false, true, true, true),    // Not loaded |              | lib/index.d.ts      | Even though map is present
			getTestcaseFindAllRefsWithDisableReferencedProjectLoad(false, true, true, false),   // Not loaded |              | lib/index.d.ts      |
			getTestcaseFindAllRefsWithDisableReferencedProjectLoad(false, true, false, true),   // Not loaded |              | index.ts            | But not helper.ts, which is not referenced from a
			getTestcaseFindAllRefsWithDisableReferencedProjectLoad(false, true, false, false),  // Not loaded |              | index.ts            | But not helper.ts, which is not referenced from a
			getTestcaseFindAllRefsWithDisableReferencedProjectLoad(false, false, true, true),   // Loaded     | Via map      | index.ts, helper.ts | Via map and newly loaded project
			getTestcaseFindAllRefsWithDisableReferencedProjectLoad(false, false, true, false),  // Not loaded |              | lib/index.d.ts      |
			getTestcaseFindAllRefsWithDisableReferencedProjectLoad(false, false, false, true),  // Loaded     | Via redirect | index.ts, helper.ts |
			getTestcaseFindAllRefsWithDisableReferencedProjectLoad(false, false, false, false), // Loaded     | Via redirect | index.ts, helper.ts |
		},
	)

	for _, test := range testCases {
		test.run(t, "findAllRefs")
	}
}

func getFindAllRefsTestCasesForDefaultProjects() []*lspServerTest {
	filesForSolutionConfigFile := func(solutionRefs []string, disableReferencedProjectLoad bool, ownFiles []string) map[string]any {
		var disableReferencedProjectLoadStr string
		if disableReferencedProjectLoad {
			disableReferencedProjectLoadStr = `"disableReferencedProjectLoad": true`
		}
		var ownFilesStr string
		if len(ownFiles) > 0 {
			ownFilesStr = strings.Join(ownFiles, ",")
		}
		files := map[string]any{
			"/user/username/workspaces/dummy/dummy.ts":      `const x = 1;`,
			"/user/username/workspaces/dummy/tsconfig.json": `{ }`,
			"/user/username/projects/myproject/tsconfig.json": stringtestutil.Dedent(fmt.Sprintf(`
				{
					"compilerOptions": {
						%s
					},
					"files": [%s],
					"references": [
						%s
					]
				}`, disableReferencedProjectLoadStr, ownFilesStr, strings.Join(core.Map(solutionRefs, func(ref string) string {
				return fmt.Sprintf(`{ "path": "%s" }`, ref)
			}), ","))),
			"/user/username/projects/myproject/tsconfig-src.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": true,
						"outDir": "./target",
					},
					"include": ["./src/**/*"]
				}`),
			"/user/username/projects/myproject/src/main.ts": stringtestutil.Dedent(`
				import { foo } from './helpers/functions';
				foo()`),
			"/user/username/projects/myproject/src/helpers/functions.ts": `export function foo() { return 1; }`,
			"/user/username/projects/myproject/indirect3/tsconfig.json":  `{ }`,
			"/user/username/projects/myproject/indirect3/main.ts": stringtestutil.Dedent(`
				import { foo } from '../target/src/main';
				foo()
				export function bar() {}`),
		}
		return files
	}
	filesForIndirectProject := func(projectIndex int, compilerOptions string) map[string]any {
		files := map[string]any{
			fmt.Sprintf("/user/username/projects/myproject/tsconfig-indirect%d.json", projectIndex): fmt.Sprintf(`
			{
				"compilerOptions": {
					"composite": true,
					"outDir": "./target/",
					%s
				},
				"files": [
					"./indirect%d/main.ts"
				],
				"references": [
					{
						"path": "./tsconfig-src.json"
					}
				]
			}`, compilerOptions, projectIndex),
			fmt.Sprintf("/user/username/projects/myproject/indirect%d/main.ts", projectIndex): `export const indirect = 1;`,
		}
		return files
	}
	applyIndirectProjectFiles := func(files map[string]any, projectIndex int, compilerOptions string) {
		maps.Copy(files, filesForIndirectProject(projectIndex, compilerOptions))
	}
	testSolution := func(server *testServer) {
		file := "/user/username/projects/myproject/src/main.ts"
		// Ensure configured project is found for open file
		server.openFile(file, lsproto.LanguageKindTypeScript)

		// !!! TODO Verify errors

		dummyFile := "/user/username/workspaces/dummy/dummy.ts"

		server.openFile(dummyFile, lsproto.LanguageKindTypeScript)

		server.closeFile(dummyFile)
		server.closeFile(file)
		server.openFile(dummyFile, lsproto.LanguageKindTypeScript)

		server.closeFile(dummyFile)
		server.openFile(file, lsproto.LanguageKindTypeScript)

		// Find all ref in default project
		server.baselineReferences(file, lsptestutil.PositionToLineAndCharacter(file, server.content(file), "foo", 1))

		server.closeFile(file)
		file = "/user/username/projects/myproject/indirect3/main.ts"
		server.openFile(file, lsproto.LanguageKindTypeScript)

		// Find all ref in non default project
		server.baselineReferences(file, lsptestutil.PositionToLineAndCharacter(file, server.content(file), "foo", 0))
	}
	return []*lspServerTest{
		{
			subScenario: "project found is solution referencing default project directly",
			files: func() map[string]any {
				return filesForSolutionConfigFile([]string{"./tsconfig-src.json"}, false, nil)
			},
			test: testSolution,
		},
		{
			subScenario: "project found is solution referencing default project indirectly",
			files: func() map[string]any {
				files := filesForSolutionConfigFile([]string{"./tsconfig-indirect1.json", "./tsconfig-indirect2.json"}, false, nil)
				applyIndirectProjectFiles(files, 1, "")
				applyIndirectProjectFiles(files, 2, "")
				return files
			},
			test: testSolution,
		},
		{
			subScenario: "project found is solution with disableReferencedProjectLoad referencing default project directly",
			files: func() map[string]any {
				return filesForSolutionConfigFile([]string{"./tsconfig-src.json"}, true, nil)
			},
			test: testSolution,
		},
		{
			subScenario: "project found is solution referencing default project indirectly through disableReferencedProjectLoad",
			files: func() map[string]any {
				files := filesForSolutionConfigFile([]string{"./tsconfig-indirect1.json", "./tsconfig-indirect2.json"}, false, nil)
				applyIndirectProjectFiles(files, 1, `"disableReferencedProjectLoad": true`)
				return files
			},
			test: testSolution,
		},
		{
			subScenario: "project found is solution referencing default project indirectly through disableReferencedProjectLoad in one but without it in another",
			files: func() map[string]any {
				files := filesForSolutionConfigFile([]string{"./tsconfig-indirect1.json", "./tsconfig-indirect2.json"}, false, nil)
				applyIndirectProjectFiles(files, 1, `"disableReferencedProjectLoad": true`)
				applyIndirectProjectFiles(files, 2, "")
				return files
			},
			test: testSolution,
		},
		{
			subScenario: "project found is project with own files referencing the file from referenced project",
			files: func() map[string]any {
				files := filesForSolutionConfigFile([]string{"./tsconfig-src.json"}, false, []string{`"./own/main.ts"`})
				files["/user/username/projects/myproject/own/main.ts"] = stringtestutil.Dedent(`
					import { foo } from '../src/main';
					foo;
					export function bar() {}
				`)
				return files
			},
			test: testSolution,
		},
	}
}

func getFindAllRefsTestcasesForRootOfReferencedProject() []*lspServerTest {
	files := func(disableSourceOfProjectReferenceRedirect bool) map[string]any {
		return map[string]any{
			"/user/username/projects/project/src/common/input/keyboard.ts": stringtestutil.Dedent(`
				function bar() { return "just a random function so .d.ts location doesnt match"; }
				export function evaluateKeyboardEvent() { }
			`),
			"/user/username/projects/project/src/common/input/keyboard.test.ts": stringtestutil.Dedent(`
				import { evaluateKeyboardEvent } from 'common/input/keyboard';
				function testEvaluateKeyboardEvent() {
					return evaluateKeyboardEvent();
				}`),
			"/user/username/projects/project/src/terminal.ts": stringtestutil.Dedent(`
				import { evaluateKeyboardEvent } from 'common/input/keyboard';
				function foo() {
					return evaluateKeyboardEvent();
				}`),
			"/user/username/projects/project/src/common/tsconfig.json": stringtestutil.Dedent(fmt.Sprintf(`
				{
					"compilerOptions": {
						"composite": true,
						"declarationMap": true,
						"outDir": "../../out",
						"disableSourceOfProjectReferenceRedirect": %v,
						"paths": {
                            "*": ["../*"],
                        },
					},
					"include": ["./**/*"]
				}`, disableSourceOfProjectReferenceRedirect)),
			"/user/username/projects/project/src/tsconfig.json": stringtestutil.Dedent(fmt.Sprintf(`
				{
					"compilerOptions": {
						"composite": true,
						"declarationMap": true,
						"outDir": "../out",
						"disableSourceOfProjectReferenceRedirect": %v,
						"paths": {
                            "common/*": ["./common/*"],
                        },
						"tsBuildInfoFile": "../out/src.tsconfig.tsbuildinfo"
					},
					"include": ["./**/*"],
					"references": [
                        { "path": "./common" },
                    ],
				}`, disableSourceOfProjectReferenceRedirect)),
		}
	}
	testSolution := func(server *testServer) {
		keyboardTs := "/user/username/projects/project/src/common/input/keyboard.ts"
		terminalTs := "/user/username/projects/project/src/terminal.ts"
		server.openFile(keyboardTs, lsproto.LanguageKindTypeScript)
		server.openFile(terminalTs, lsproto.LanguageKindTypeScript)

		// Find all ref in default project
		server.baselineReferences(keyboardTs, lsptestutil.PositionToLineAndCharacter(keyboardTs, server.content(keyboardTs), "evaluateKeyboardEvent", 0))
	}
	return []*lspServerTest{
		{
			subScenario: "root file is file from referenced project",
			files: func() map[string]any {
				return files(false)
			},
			test: testSolution,
		},
		{
			subScenario: "root file is file from referenced project and using declaration maps",
			files: func() map[string]any {
				return tsctests.GetFileMapWithBuild(files(true), []string{"-b", "/user/username/projects/project/src/tsconfig.json"})
			},
			test: testSolution,
		},
	}
}

func getFindAllRefsFileMapForLocalness(disableSolutionSearching bool) map[string]any {
	var extraStr string
	if disableSolutionSearching {
		extraStr = `"disableSolutionSearching": true,`
	}
	return map[string]any{
		"/user/username/projects/solution/tsconfig.json": stringtestutil.Dedent(`
			{
				"files": [],
				"include": [],
				"references": [
					{ "path": "./compiler" },
					{ "path": "./services" },
				],
			}`),
		"/user/username/projects/solution/compiler/tsconfig.json": stringtestutil.Dedent(fmt.Sprintf(`
			{
				"compilerOptions": {
					"composite": true,
					%s
				},
				"files": ["./types.ts", "./program.ts"]
			}`, extraStr)),
		"/user/username/projects/solution/compiler/types.ts": stringtestutil.Dedent(`
			namespace ts {
				export interface Program {
					getSourceFiles(): string[];
				}
			}`),
		"/user/username/projects/solution/compiler/program.ts": stringtestutil.Dedent(`
			namespace ts {
				export const program: Program = {
					getSourceFiles: () => [getSourceFile()]
				};
				function getSourceFile() { return "something"; }
			}`),
		"/user/username/projects/solution/services/tsconfig.json": stringtestutil.Dedent(`
			{
				"compilerOptions": {
					"composite": true
				},
				"files": ["./services.ts"],
				"references": [
					{ "path": "../compiler" },
				],
			}`),
		"/user/username/projects/solution/services/services.ts": stringtestutil.Dedent(`
			/// <reference path="../compiler/types.ts" />
			/// <reference path="../compiler/program.ts" />
			namespace ts {
				const result = program.getSourceFiles();
			}`),
	}
}

func getFindAllRefsTestCaseForSpecialLocalnessHandling(scenario string, definition string, usage string, referenceTerm string) *lspServerTest {
	return &lspServerTest{
		subScenario: "special handling of localness " + scenario,
		files: func() map[string]any {
			return map[string]any{
				"/user/username/projects/solution/tsconfig.json": stringtestutil.Dedent(`
					{
						"files": [],
						"references": [
							{ "path": "./api" },
							{ "path": "./app" },
						],
					}`),
				"/user/username/projects/solution/api/tsconfig.json": stringtestutil.Dedent(`
					{
						"compilerOptions": {
							"composite": true,
							"outDir": "dist",
							"rootDir": "src"
						},
						"include": ["src"],
						"references": [{ "path": "../shared" }],
					}`),
				"/user/username/projects/solution/api/src/server.ts": `import * as shared from "../../shared/dist"` + "\n" + usage,
				"/user/username/projects/solution/app/tsconfig.json": stringtestutil.Dedent(`
					{
						"compilerOptions": {
							"composite": true,
							"outDir": "dist",
							"rootDir": "src"
						},
						"include": ["src"],
						"references": [{ "path": "../shared" }],
					}`),
				"/user/username/projects/solution/app/src/app.ts": `import * as shared from "../../shared/dist"` + "\n" + usage,
				"/user/username/projects/solution/shared/tsconfig.json": stringtestutil.Dedent(`
					{
						"compilerOptions": {
							"composite": true,
							"outDir": "dist",
							"rootDir": "src"
						},
						"include": ["src"],
					}`),
				"/user/username/projects/solution/shared/src/index.ts": definition,
			}
		},
		test: func(server *testServer) {
			apiFile := "/user/username/projects/solution/api/src/server.ts"
			server.openFile(apiFile, lsproto.LanguageKindTypeScript)

			// Find all references
			server.baselineReferences(apiFile, lsptestutil.PositionToLineAndCharacter(apiFile, server.content(apiFile), referenceTerm, 0))
		},
	}
}

func getTestcaseFindAllRefsWithDisableReferencedProjectLoad(
	projectAlreadyLoaded bool,
	disableReferencedProjectLoad bool,
	disableSourceOfProjectReferenceRedirect bool,
	dtsMapPresent bool,
) *lspServerTest {
	subScenario := fmt.Sprintf(`when proj %s loaded`, core.IfElse(projectAlreadyLoaded, "is", "is not")) +
		` and refd proj loading is ` + core.IfElse(disableReferencedProjectLoad, "disabled", "enabled") +
		` and proj ref redirects are ` + core.IfElse(disableSourceOfProjectReferenceRedirect, "disabled", "enabled") +
		` and a decl map is ` + core.IfElse(dtsMapPresent, "present", "missing")

	return &lspServerTest{
		subScenario: "find refs to decl in other proj " + subScenario,
		files: func() map[string]any {
			files := map[string]any{
				"/user/username/projects/myproject/a/tsconfig.json": stringtestutil.Dedent(fmt.Sprintf(`
					{
						"disableReferencedProjectLoad": %t,
						"disableSourceOfProjectReferenceRedirect": %t,
						"composite": true
					}`, disableReferencedProjectLoad, disableSourceOfProjectReferenceRedirect)),
				"/user/username/projects/myproject/a/index.ts": stringtestutil.Dedent(`
					import { B } from "../b/lib";
					const b: B = new B();`),
				"/user/username/projects/myproject/b/tsconfig.json": stringtestutil.Dedent(`
					{
						"declarationMap": true,
						"outDir": "lib",
						"composite": true,
					}`),
				"/user/username/projects/myproject/b/index.ts": stringtestutil.Dedent(`
					export class B {
						M() {}
					}`),
				"/user/username/projects/myproject/b/helper.ts": stringtestutil.Dedent(`
					import { B } from ".";
					const b: B = new B();`),
				"/user/username/projects/myproject/b/lib/index.d.ts": stringtestutil.Dedent(`
					export declare class B {
						M(): void;
					}
					//# sourceMappingURL=index.d.ts.map`),
			}
			if dtsMapPresent {
				files["/user/username/projects/myproject/b/lib/index.d.ts.map"] = stringtestutil.Dedent(`
					{
						"version": 3,
						"file": "index.d.ts",
						"sourceRoot": "",
						"sources": ["../index.ts"],
						"names": [],
						"mappings": "AAAA,qBAAa,CAAC;IACV,CAAC;CACJ"
					}`)
			}
			return files
		},
		test: func(server *testServer) {
			indexA := "/user/username/projects/myproject/a/index.ts"
			server.openFile(indexA, lsproto.LanguageKindTypeScript)
			if projectAlreadyLoaded {
				server.openFile("/user/username/projects/myproject/b/helper.ts", lsproto.LanguageKindTypeScript)
			}
			server.baselineReferences(indexA, lsptestutil.PositionToLineAndCharacter(indexA, server.content(indexA), "B", 1))
		},
	}
}
