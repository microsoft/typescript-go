package lspservertests

import (
	"fmt"
	"slices"
	"testing"

	"github.com/microsoft/typescript-go/internal/bundled"
	"github.com/microsoft/typescript-go/internal/execute/tsctests"
	"github.com/microsoft/typescript-go/internal/lsp/lsproto"
	"github.com/microsoft/typescript-go/internal/testutil/lsptestutil"
	"github.com/microsoft/typescript-go/internal/testutil/stringtestutil"
)

func TestDeclarationMaps(t *testing.T) {
	t.Parallel()

	if !bundled.Embedded {
		t.Skip("bundled files are not embedded")
	}

	testCases := slices.Concat(
		getDeclarationMapTestCasesForProjectReferences(),
		getDeclarationMapTestCasesForMaps(),
		getDeclarationMapTestCasesForRename(),
		[]*lspServerTest{
			{
				subScenario: "findAllReferences definition is in mapped file",
				files: func() map[string]any {
					return map[string]any{
						"/home/src/projects/project/a/a.ts": "export function f() {}",
						"/home/src/projects/project/a/tsconfig.json": stringtestutil.Dedent(`
							{
								"compilerOptions": {
									"outDir": "bin",
									"declarationMap": true,
									"composite": true
								}
							}`),
						"/home/src/projects/project/b/b.ts": stringtestutil.Dedent(`
							import { f } from "../a/bin/a";
							f();`),
						"/home/src/projects/project/b/tsconfig.json": stringtestutil.Dedent(`
							{
								"references": [
									{ "path": "../a" }
								]
							}`),
						"/home/src/projects/project/bin/a.d.ts": stringtestutil.Dedent(`
							export declare function f(): void;
							//# sourceMappingURL=a.d.ts.map`),
						"/home/src/projects/project/bin/a.d.ts.map": stringtestutil.Dedent(`
							{
								"version":3,
								"file":"a.d.ts",
								"sourceRoot":"",
								"sources":["a.ts"],
								"names":[],
								"mappings":"AAAA,wBAAgB,CAAC,SAAK"
							}`),
					}
				},
				test: func(server *testServer) {
					bTs := "/home/src/projects/project/b/b.ts"
					server.openFile(bTs, lsproto.LanguageKindTypeScript)

					// Ref projects are loaded after as part of this command
					server.baselineReferences(bTs, lsptestutil.PositionToLineAndCharacter(bTs, server.content(bTs), "f()", 0))
				},
			},
		},
	)

	for _, test := range testCases {
		test.run(t, "declarationMaps")
	}
}

func getDeclarationMapTestCasesForMaps() []*lspServerTest {
	files := func() map[string]any {
		configContent := stringtestutil.Dedent(`
			{
				"compilerOptions": {
					"outDir": "bin",
					"declarationMap": true,
					"composite": true
				}
			}`)
		return map[string]any{
			"/home/src/projects/project/a/a.ts": stringtestutil.Dedent(`
				export function fnA() {}
				export interface IfaceA {}
				export const instanceA: IfaceA = {};
			`),
			"/home/src/projects/project/a/tsconfig.json": configContent,
			"/home/src/projects/project/a/bin/a.d.ts.map": stringtestutil.Dedent(`
				{
					"version": 3,
					"file": "a.d.ts",
					"sourceRoot": "",
					"sources": ["../a.ts"],
					"names": [],
					"mappings": "AAAA,wBAAgB,GAAG,SAAK;AACxB,MAAM,WAAW,MAAM;CAAG;AAC1B,eAAO,MAAM,SAAS,EAAE,MAAW,CAAC"
				}`),
			"/home/src/projects/project/a/bin/a.d.ts": stringtestutil.Dedent(`
				export declare function fnA(): void;
				export interface IfaceA {
				}
				export declare const instanceA: IfaceA;
				//# sourceMappingURL=a.d.ts.map`),
			"/home/src/projects/project/b/tsconfig.json": configContent,
			"/home/src/projects/project/b/bin/b.d.ts.map": stringtestutil.Dedent(`
				{
					"version": 3,
					"file": "b.d.ts",
					"sourceRoot": "",
					"sources": ["../b.ts"],
					"names": [],
					"mappings": "AAAA,wBAAgB,GAAG,SAAK"
				}`),
			"/home/src/projects/project/b/bin/b.d.ts": stringtestutil.Dedent(`
				export declare function fnB(): void;
				//# sourceMappingURL=b.d.ts.map`),
			"/home/src/projects/project/user/user.ts": stringtestutil.Dedent(`
				import * as a from "../a/bin/a";
				import * as b from "../b/bin/b";
				export function fnUser() { a.fnA(); b.fnB(); a.instanceA; }`),
			"/home/src/projects/project/dummy/dummy.ts":      "export const a = 10;",
			"/home/src/projects/project/dummy/tsconfig.json": "{}",
		}
	}
	return []*lspServerTest{
		{
			subScenario: "findAllReferences",
			files:       files,
			test: func(server *testServer) {
				userTs := "/home/src/projects/project/user/user.ts"
				server.openFile(userTs, lsproto.LanguageKindTypeScript)

				// Ref projects are loaded after as part of this command
				server.baselineReferences(userTs, lsptestutil.PositionToLineAndCharacter(userTs, server.content(userTs), "fnA()", 0))

				// Open temp file and verify all projects alive
				server.closeFile(userTs)
				server.openFile("/home/src/projects/project/dummy/dummy.ts", lsproto.LanguageKindTypeScript)
			},
		},
		{
			subScenario: "findAllReferences starting at definition",
			files:       files,
			test: func(server *testServer) {
				userTs := "/home/src/projects/project/user/user.ts"
				server.openFile(userTs, lsproto.LanguageKindTypeScript)
				aTs := "/home/src/projects/project/a/a.ts"
				server.openFile(aTs, lsproto.LanguageKindTypeScript) // If it's not opened, the reference isn't found.

				// Ref projects are loaded after as part of this command
				server.baselineReferences(aTs, lsptestutil.PositionToLineAndCharacter(aTs, server.content(aTs), "fnA", 0))

				// Open temp file and verify all projects alive
				server.closeFile(userTs)
				server.openFile("/home/src/projects/project/dummy/dummy.ts", lsproto.LanguageKindTypeScript)
			},
		},
		{
			subScenario: "findAllReferences target does not exist",
			files:       files,
			test: func(server *testServer) {
				userTs := "/home/src/projects/project/user/user.ts"
				server.openFile(userTs, lsproto.LanguageKindTypeScript)

				// Ref projects are loaded after as part of this command
				server.baselineReferences(userTs, lsptestutil.PositionToLineAndCharacter(userTs, server.content(userTs), "fnB()", 0))

				// Open temp file and verify all projects alive
				server.closeFile(userTs)
				server.openFile("/home/src/projects/project/dummy/dummy.ts", lsproto.LanguageKindTypeScript)
			},
		},
		{
			subScenario: "rename",
			files:       files,
			test: func(server *testServer) {
				userTs := "/home/src/projects/project/user/user.ts"
				server.openFile(userTs, lsproto.LanguageKindTypeScript)

				// Ref projects are loaded after as part of this command
				server.baselineRename(userTs, lsptestutil.PositionToLineAndCharacter(userTs, server.content(userTs), "fnA()", 0))

				// Open temp file and verify all projects alive
				server.closeFile(userTs)
				server.openFile("/home/src/projects/project/dummy/dummy.ts", lsproto.LanguageKindTypeScript)
			},
		},
		{
			subScenario: "rename starting at definition",
			files:       files,
			test: func(server *testServer) {
				userTs := "/home/src/projects/project/user/user.ts"
				server.openFile(userTs, lsproto.LanguageKindTypeScript)
				aTs := "/home/src/projects/project/a/a.ts"
				server.openFile(aTs, lsproto.LanguageKindTypeScript) // If it's not opened, the reference isn't found.

				// Ref projects are loaded after as part of this command
				server.baselineRename(aTs, lsptestutil.PositionToLineAndCharacter(aTs, server.content(aTs), "fnA", 0))

				// Open temp file and verify all projects alive
				server.closeFile(userTs)
				server.openFile("/home/src/projects/project/dummy/dummy.ts", lsproto.LanguageKindTypeScript)
			},
		},
		{
			subScenario: "rename target does not exist",
			files:       files,
			test: func(server *testServer) {
				userTs := "/home/src/projects/project/user/user.ts"
				server.openFile(userTs, lsproto.LanguageKindTypeScript)

				// Ref projects are loaded after as part of this command
				server.baselineRename(userTs, lsptestutil.PositionToLineAndCharacter(userTs, server.content(userTs), "fnB()", 0))

				// Open temp file and verify all projects alive
				server.closeFile(userTs)
				server.openFile("/home/src/projects/project/dummy/dummy.ts", lsproto.LanguageKindTypeScript)
			},
		},
		{
			subScenario: "workspace symbols",
			files: func() map[string]any {
				allFiles := files()
				allFiles["/home/src/projects/project/user/user.ts"] = stringtestutil.Dedent(`
					import * as a from "../a/a";
					import * as b from "../b/b";
					export function fnUser() {
						a.fnA();
						b.fnB();
						a.instanceA;
					}`)
				allFiles["/home/src/projects/project/user/tsconfig.json"] = stringtestutil.Dedent(`
					{
						"references": [{ "path": "../a" }, { "path": "../b" }]
					}`)
				allFiles["/home/src/projects/project/b/b.ts"] = stringtestutil.Dedent(`
					export function fnB() {}`)
				allFiles["/home/src/projects/project/b/c.ts"] = stringtestutil.Dedent(`
					export function fnC() {}`)
				return allFiles
			},
			test: func(server *testServer) {
				userTs := "/home/src/projects/project/user/user.ts"
				server.openFile(userTs, lsproto.LanguageKindTypeScript)
				server.baselineWorkspaceSymbol("fn")
				// Open temp file and verify all projects alive
				server.closeFile(userTs)
				server.openFile("/home/src/projects/project/dummy/dummy.ts", lsproto.LanguageKindTypeScript)
			},
		},
	}
}

func getDeclarationMapTestCasesForProjectReferences() []*lspServerTest {
	files := func(disableSourceOfProjectReferenceRedirect bool) map[string]any {
		return map[string]any{
			"/user/username/projects/a/a.ts":          "export class A { }",
			"/user/username/projects/a/tsconfig.json": "{}",
			"/user/username/projects/a/a.d.ts": stringtestutil.Dedent(`
				export declare class A {
				}
				//# sourceMappingURL=a.d.ts.map`),
			"/user/username/projects/a/a.d.ts.map": stringtestutil.Dedent(`
				{
					"version": 3,
					"file": "a.d.ts",
					"sourceRoot": "",
					"sources": ["./a.ts"],
					"names": [],
					"mappings": "AAAA,qBAAa,CAAC;CAAI"
				}
			`),
			"/user/username/projects/b/b.ts": stringtestutil.Dedent(`
				import {A} from "../a/a";
				new A();
			`),
			"/user/username/projects/b/tsconfig.json": stringtestutil.Dedent(fmt.Sprintf(`
			{
				"compilerOptions": {
					"disableSourceOfProjectReferenceRedirect": %t
				},
				"references": [
					{ "path": "../a" }
				]
			}`, disableSourceOfProjectReferenceRedirect)),
		}
	}
	test := func(server *testServer) {
		bTs := "/user/username/projects/b/b.ts"
		server.openFile(bTs, lsproto.LanguageKindTypeScript)

		// Ref projects are loaded after as part of this command
		server.baselineReferences(bTs, lsptestutil.PositionToLineAndCharacter(bTs, server.content(bTs), "A();", 0))
	}

	return []*lspServerTest{
		{
			subScenario: "opening original location project",
			files:       func() map[string]any { return files(false) },
			test:        test,
		},
		{
			subScenario: "opening original location project disableSourceOfProjectReferenceRedirect",
			files:       func() map[string]any { return files(true) },
			test:        test,
		},
	}
}

func getDeclarationMapTestCasesForRename() []*lspServerTest {
	dependencyTs := "/user/username/projects/myproject/dependency/FnS.ts"
	filesWithRef := func() map[string]any {
		return map[string]any{
			dependencyTs: stringtestutil.Dedent(`
				export function fn1() { }
				export function fn2() { }
				export function fn3() { }
				export function fn4() { }
				export function fn5() { }
			`) + "\n",
			"/user/username/projects/myproject/dependency/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": true,
						"declarationMap": true,
						"declarationDir": "../decls"
					}
				}`),
			"/user/username/projects/myproject/main/main.ts": stringtestutil.Dedent(`
				import {
					fn1,
					fn2,
					fn3,
					fn4,
					fn5
				} from "../decls/FnS";

				fn1();
				fn2();
				fn3();
				fn4();
				fn5();
			`),
			"/user/username/projects/myproject/main/tsconfig.json": stringtestutil.Dedent(`
				{
					"compilerOptions": {
						"composite": true,
						"declarationMap": true,
					},
					"references": [
						{ "path": "../dependency" }
					]
				}`),
			"/user/username/projects/myproject/tsconfig.json": stringtestutil.Dedent(`
				{
					"references": [
						{ "path": "main" }
					]
				}`),
			"/user/username/projects/random/random.ts":     "export const a = 10;",
			"/user/username/projects/random/tsconfig.json": "{}",
		}
	}
	filesWithBuiltRef := func() map[string]any {
		allFiles := filesWithRef()
		tsctests.GetFileMapWithBuild(allFiles, []string{"-b", "/user/username/projects/myproject/tsconfig.json"})
		return allFiles
	}

	fileWithNoRefs := func() map[string]any {
		allFiles := filesWithBuiltRef()
		allFiles["/user/username/projects/myproject/main/tsconfig.json"] = stringtestutil.Dedent(`
			{
				"compilerOptions": {
					"composite": true,
					"declarationMap": true,
				},
			}`)
		return allFiles
	}

	filesWithDisableProjectRefSource := func() map[string]any {
		allFiles := filesWithBuiltRef()
		allFiles["/user/username/projects/myproject/main/tsconfig.json"] = stringtestutil.Dedent(`
			{
				"compilerOptions": {
					"composite": true,
					"declarationMap": true,
					"disableSourceOfProjectReferenceRedirect": true
				},
				"references": [
					{ "path": "../dependency" }
				 ]
			}`)
		return allFiles
	}

	baselineRename := func(server *testServer) {
		server.baselineRename(dependencyTs, lsptestutil.PositionToLineAndCharacter(dependencyTs, server.content(dependencyTs), "fn3", 0))
	}

	testRename := func(server *testServer) {
		dummyTs := "/user/username/projects/random/random.ts"
		server.openFile(dependencyTs, lsproto.LanguageKindTypeScript)
		server.openFile(dummyTs, lsproto.LanguageKindTypeScript)

		// Ref projects are loaded after as part of this command
		baselineRename(server)

		// Collecting at this point retains dependency.d.ts and map
		server.closeFile(dummyTs)
		server.openFile(dummyTs, lsproto.LanguageKindTypeScript)

		// Closing open file, removes dependencies too
		server.closeFile(dependencyTs)
		server.closeFile(dummyTs)
		server.openFile(dummyTs, lsproto.LanguageKindTypeScript)
	}

	openFileAndBaselineRename := func(server *testServer) {
		server.openFile(dependencyTs, lsproto.LanguageKindTypeScript)
		baselineRename(server)
	}

	testPrefix := func(server *testServer) {
		openFileAndBaselineRename(server)
		server.changeFile(&lsproto.DidChangeTextDocumentParams{
			TextDocument: lsproto.VersionedTextDocumentIdentifier{
				Uri:     lsproto.DocumentUri("file://" + dependencyTs),
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
						Text: "function fooBar() { }\n",
					},
				},
			},
		})
		baselineRename(server)
	}

	testSuffix := func(server *testServer) {
		openFileAndBaselineRename(server)
		server.changeFile(&lsproto.DidChangeTextDocumentParams{
			TextDocument: lsproto.VersionedTextDocumentIdentifier{
				Uri:     lsproto.DocumentUri("file://" + dependencyTs),
				Version: 2,
			},
			ContentChanges: []lsproto.TextDocumentContentChangePartialOrWholeDocument{
				{
					Partial: &lsproto.TextDocumentContentChangePartial{
						Range: lsproto.Range{
							Start: lsproto.Position{
								Line:      5,
								Character: 0,
							},
							End: lsproto.Position{
								Line:      5,
								Character: 0,
							},
						},
						Text: "const x = 10;",
					},
				},
			},
		})
		baselineRename(server)
	}

	getAllConfigKindTests := func(subScenario string, testFn func(server *testServer)) []*lspServerTest {
		return []*lspServerTest{
			{
				subScenario: subScenario + " with project references",
				files:       filesWithBuiltRef,
				test:        testFn,
			},
			{
				subScenario: subScenario + " with disableSourceOfProjectReferenceRedirect",
				files:       filesWithDisableProjectRefSource,
				test:        testFn,
			},
			{
				subScenario: subScenario + " with source maps",
				files:       fileWithNoRefs,
				test:        testFn,
			},
		}
	}

	return slices.Concat(
		getAllConfigKindTests("rename", testRename),
		getAllConfigKindTests("rename on edit", testPrefix),
		getAllConfigKindTests("rename on edit at end", testSuffix),
		[]*lspServerTest{
			{
				subScenario: "rename before project is built",
				files:       filesWithRef,
				test:        testRename,
			},
		},
	)
}
