package tsoptions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/microsoft/typescript-go/internal/diagnosticwriter"
	"github.com/microsoft/typescript-go/internal/parser"
	"github.com/microsoft/typescript-go/internal/repo"
	"github.com/microsoft/typescript-go/internal/testutil/baseline"
	"github.com/microsoft/typescript-go/internal/tspath"
	"github.com/microsoft/typescript-go/internal/vfs"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

type testConfig struct {
	jsonText       string
	configFileName string
	basePath       string
	allFileList    []string
}

func fixRoot(path string) string {
	rootLength := tspath.GetRootLength(path)
	if rootLength == 0 {
		return path
	}
	if len(path) == rootLength {
		return "."
	}
	return path[rootLength:]
}

type vfsParseConfigHost struct {
	fs               vfs.FS
	currentDirectory string
}

var _ ParseConfigHost = (*vfsParseConfigHost)(nil)

func (h *vfsParseConfigHost) FS() vfs.FS {
	return h.fs
}

func (h *vfsParseConfigHost) GetCurrentDirectory() string {
	return h.currentDirectory
}

func newVFSParseConfigHost(files map[string]string, currentDirectory string) *vfsParseConfigHost {
	fs := fstest.MapFS{}
	for name, content := range files {
		fs[fixRoot(name)] = &fstest.MapFile{
			Data: []byte(content),
		}
	}
	return &vfsParseConfigHost{
		fs:               vfstest.FromMapFS(fs, true /*useCaseSensitiveFileNames*/),
		currentDirectory: currentDirectory,
	}
}

var parseConfigFileTextToJsonTests = []struct {
	title string
	input []string
}{
	{
		title: "returns empty config for file with only whitespaces",
		input: []string{
			"",
			" ",
		},
	},
	{
		title: "returns empty config for file with comments only",
		input: []string{
			"// Comment",
			"/* Comment*/",
		},
	},
	{
		title: "returns empty config when config is empty object",
		input: []string{
			`{}`,
		},
	},
	{
		title: "returns config object without comments",
		input: []string{
			`{ // Excluded files
            "exclude": [
                // Exclude d.ts
                "file.d.ts"
            ]
        }`,
			`{
            /* Excluded
                    Files
            */
            "exclude": [
                /* multiline comments can be in the middle of a line */"file.d.ts"
            ]
        }`,
		},
	},
	{
		title: "keeps string content untouched",
		input: []string{
			`{
            "exclude": [
                "xx//file.d.ts"
            ]
        }`,
			`{
            "exclude": [
                "xx/*file.d.ts*/"
            ]
        }`,
		},
	},
	// {
	// 	title: "handles escaped characters in strings correctly",
	// 	input: []string{
	// 		`{
	// 			"exclude": [
	// 				"xx\\"//files"
	// 			]
	// 		}`,
	// 		`{
	// 			"exclude": [
	// 				"xx\\\\" // end of line comment
	// 			]
	// 		}`,
	// 	},
	// 	output: []map[string]any{
	// 		{"exclude": []string{"xx\"//files"}},
	// 		{"exclude": []string{"xx\\"}},
	// 	},
	// },
	{
		title: "returns object when users correctly specify library",
		input: []string{
			`{
            "compilerOptions": {
                "lib": ["es5"]
            }
        }`,
			`{
            "compilerOptions": {
                "lib": ["es5", "es6"]
            }
        }`,
		},
	},
}

func TestParseConfigFileTextToJson(t *testing.T) {
	t.Parallel()
	for _, rec := range parseConfigFileTextToJsonTests {
		t.Run(rec.title, func(t *testing.T) {
			t.Parallel()
			var baselineContent strings.Builder
			for i, jsonText := range rec.input {
				baselineContent.WriteString("Input::\n")
				baselineContent.WriteString(jsonText + "\n")
				parsed, errors := ParseConfigFileTextToJson("/apath/tsconfig.json", "/apath", jsonText)
				if configText, err := jsonToReadableText(parsed); err != nil {
					t.Fatal(err)
				} else {
					baselineContent.WriteString("Config::\n")
					baselineContent.WriteString(configText)
				}
				baselineContent.WriteString("Errors::\n")
				diagnosticwriter.FormatDiagnosticsWithColorAndContext(&baselineContent, errors, &diagnosticwriter.FormattingOptions{
					NewLine: "\n",
					ComparePathsOptions: tspath.ComparePathsOptions{
						CurrentDirectory:          "/",
						UseCaseSensitiveFileNames: true,
					},
				})
				baselineContent.WriteString("\n")
				if i != len(rec.input)-1 {
					baselineContent.WriteString("\n")
				}
			}
			baseline.RunAgainstSubmodule(t, rec.title+" jsonParse.js", baselineContent.String(), baseline.Options{Subfolder: "config/tsconfigParsing"})
		})
	}
}

var parseJsonConfigFileTests = []struct {
	title               string
	noSubmoduleBaseline bool
	input               []testConfig
}{
	{
		title: "ignore dotted files and folders",
		input: []testConfig{{
			jsonText:       `{}`,
			configFileName: "tsconfig.json",
			basePath:       "/apath",
			allFileList:    []string{"/apath/test.ts", "/apath/.git/a.ts", "/apath/.b.ts", "/apath/..c.ts"},
		}},
	},
	{
		title: "allow dotted files and folders when explicitly requested",
		input: []testConfig{{
			jsonText: `{
                    "files": ["/apath/.git/a.ts", "/apath/.b.ts", "/apath/..c.ts"]
                }`,
			configFileName: "tsconfig.json",
			basePath:       "/apath",
			allFileList:    []string{"/apath/test.ts", "/apath/.git/a.ts", "/apath/.b.ts", "/apath/..c.ts"},
		}},
	},
	{
		title: "implicitly exclude common package folders",
		input: []testConfig{{
			jsonText:       `{}`,
			configFileName: "tsconfig.json",
			basePath:       "/",
			allFileList:    []string{"/node_modules/a.ts", "/bower_components/b.ts", "/jspm_packages/c.ts", "/d.ts", "/folder/e.ts"},
		}},
	},
	{
		title: "generates errors for empty files list",
		input: []testConfig{{
			jsonText: `{
                "files": []
            }`,
			configFileName: "/apath/tsconfig.json",
			basePath:       "/apath",
			allFileList:    []string{"/apath/a.ts"},
		}},
	},
	{
		title: "generates errors for empty files list when no references are provided",
		input: []testConfig{{
			jsonText: `{
                "files": [],
                "references": []
            }`,
			configFileName: "/apath/tsconfig.json",
			basePath:       "/apath",
			allFileList:    []string{"/apath/a.ts"},
		}},
	},
	{
		title: "generates errors for directory with no .ts files",
		input: []testConfig{{
			jsonText: `{
            }`,
			configFileName: "/apath/tsconfig.json",
			basePath:       "/apath",
			allFileList:    []string{"/apath/a.js"},
		}},
	},
	{
		title: "generates errors for empty include",
		input: []testConfig{{
			jsonText: `{
                "include": []
            }`,
			configFileName: "/apath/tsconfig.json",
			basePath:       "tests/cases/unittests",
			allFileList:    []string{"/apath/a.ts"},
		}},
	},
	{
		title:               "parses tsconfig with compilerOptions, files, include, and exclude",
		noSubmoduleBaseline: true,
		input: []testConfig{{
			jsonText: `{
  "compilerOptions": {
    "outDir": "./dist",
    "strict": true,
    "noImplicitAny": true,
    "target": "ES2017",
    "module": "ESNext",
    "moduleResolution": "bundler",
    "moduleDetection": "auto",
    "jsx": "react",
  },
  "files": ["/apath/src/index.ts", "/apath/src/app.ts"],
  "include": ["/apath/src/**/*"],
  "exclude": ["/apath/node_modules", "/apath/dist"]
}`,
			configFileName: "/apath/tsconfig.json",
			basePath:       "/apath",
			allFileList:    []string{"/apath/src/index.ts", "/apath/src/app.ts", "/apath/node_modules/module.ts", "/apath/dist/output.js"},
		}},
	},
	{
		title: "generates errors when commandline option is in tsconfig",
		input: []testConfig{{
			jsonText: `{
  "compilerOptions": {
    "help": true
  }
}`,
			configFileName: "/apath/tsconfig.json",
			basePath:       "/apath",
			allFileList:    []string{"/apath/a.ts"},
		}},
	},
	{
		title: "does not generate errors for empty files list when one or more references are provided",
		input: []testConfig{{
			jsonText: `{
                "files": [],
                "references": [{ "path": "/apath" }]
            }`,
			configFileName: "/apath/tsconfig.json",
			basePath:       "/apath",
			allFileList:    []string{"/apath/a.ts"},
		}},
	},
	{
		title: "exclude outDir unless overridden",
		input: []testConfig{{
			jsonText: `{
                "compilerOptions": {
                    "outDir": "bin"
                }
            }`,
			configFileName: "tsconfig.json",
			basePath:       "/",
			allFileList:    []string{"/bin/a.ts", "/b.ts"},
		}, {
			jsonText: `{
                "compilerOptions": {
                    "outDir": "bin"
                },
                "exclude": [ "obj" ]
            }`,
			configFileName: "tsconfig.json",
			basePath:       "/",
			allFileList:    []string{"/bin/a.ts", "/b.ts"},
		}},
	},
	{
		title: "exclude declarationDir unless overridden",
		input: []testConfig{{
			jsonText: `{
                "compilerOptions": {
                    "declarationDir": "declarations"
                }
            }`,
			configFileName: "tsconfig.json",
			basePath:       "/",
			allFileList:    []string{"/declarations/a.d.ts", "/a.ts"},
		}, {
			jsonText: `{
                "compilerOptions": {
                    "declarationDir": "declarations"
                },
                "exclude": [ "types" ]
            }`,
			configFileName: "tsconfig.json",
			basePath:       "/",
			allFileList:    []string{"/declarations/a.d.ts", "/a.ts"},
		}},
	},
	{
		title: "generates errors for empty directory",
		input: []testConfig{{
			jsonText: `{
                "compilerOptions": {
                    "allowJs": true
                }
            }`,
			configFileName: "/apath/tsconfig.json",
			basePath:       "/apath",
			allFileList:    []string{},
		}},
	},
	{
		title: "generates errors for includes with outDir",
		input: []testConfig{{
			jsonText: `{
                "compilerOptions": {
                    "outDir": "./"
                },
                "include": ["**/*"]
            }`,
			configFileName: "/apath/tsconfig.json",
			basePath:       "/apath",
			allFileList:    []string{"/apath/a.ts"},
		}},
	},
	{
		title: "generates errors when include is not string",
		input: []testConfig{{
			jsonText: `{
  "include": [
    [
      "./**/*.ts"
    ]
  ]
}`,
			configFileName: "/apath/tsconfig.json",
			basePath:       "/apath",
			allFileList:    []string{"/apath/a.ts"},
		}},
	},
	{
		title: "generates errors when files is not string",
		input: []testConfig{{
			jsonText: `{
  "files": [
    {
      "compilerOptions": {
        "experimentalDecorators": true,
        "allowJs": true
      }
    }
  ]
}`,
			configFileName: "/apath/tsconfig.json",
			basePath:       "/apath",
			allFileList:    []string{"/apath/a.ts"},
		}},
	},
}

func TestParseJsonConfigFileContent(t *testing.T) {
	t.Parallel()
	for _, rec := range parseJsonConfigFileTests {
		t.Run(rec.title+" with json api", func(t *testing.T) {
			t.Parallel()
			baselineParseConfigWith(t, rec.title+" with json api.js", rec.noSubmoduleBaseline, rec.input, func(config testConfig, host ParseConfigHost, basePath string) ParsedCommandLine {
				parsed, _ := ParseConfigFileTextToJson(config.configFileName, config.basePath, config.jsonText)
				return ParseJsonConfigFileContent(
					parsed,
					host,
					basePath,
					nil,
					tspath.GetNormalizedAbsolutePath(config.configFileName, basePath),
					/*resolutionStack*/ nil,
					/*extraFileExtensions*/ nil,
					/*extendedConfigCache*/ nil,
				)
			})
		})
	}
}

func TestParseJsonSourceFileConfigFileContent(t *testing.T) {
	t.Parallel()
	for _, rec := range parseJsonConfigFileTests {
		t.Run(rec.title+" with jsonSourceFile api", func(t *testing.T) {
			t.Parallel()
			baselineParseConfigWith(t, rec.title+" with jsonSourceFile api.js", rec.noSubmoduleBaseline, rec.input, func(config testConfig, host ParseConfigHost, basePath string) ParsedCommandLine {
				parsed := parser.ParseJSONText(config.configFileName, config.jsonText)
				tsConfigSourceFile := &tsConfigSourceFile{
					sourceFile: parsed,
				}
				return ParseJsonSourceFileConfigFileContent(
					tsConfigSourceFile,
					host,
					host.GetCurrentDirectory(),
					nil,
					tspath.GetNormalizedAbsolutePath(config.configFileName, basePath),
					/*resolutionStack*/ nil,
					/*extraFileExtensions*/ nil,
					/*extendedConfigCache*/ nil,
				)
			})
		})
	}
}

func baselineParseConfigWith(t *testing.T, baselineFileName string, noSubmoduleBaseline bool, input []testConfig, getParsed func(config testConfig, host ParseConfigHost, basePath string) ParsedCommandLine) {
	var baselineContent strings.Builder
	for i, config := range input {
		basePath := config.basePath
		if basePath == "" {
			basePath = tspath.GetNormalizedAbsolutePath(tspath.GetDirectoryPath(config.configFileName), "")
		}
		configFileName := tspath.CombinePaths(basePath, config.configFileName)
		allFileLists := make(map[string]string, len(config.allFileList)+1)
		for _, file := range config.allFileList {
			allFileLists[file] = ""
		}
		allFileLists[configFileName] = config.jsonText
		host := newVFSParseConfigHost(allFileLists, config.basePath)
		parsedConfigFileContent := getParsed(config, host, basePath)

		baselineContent.WriteString("Fs::\n")
		printFS(&baselineContent, host.FS(), "/")
		baselineContent.WriteString("\n")
		baselineContent.WriteString("configFileName:: " + config.configFileName + "\n")
		baselineContent.WriteString("FileNames::\n")
		baselineContent.WriteString(strings.Join(parsedConfigFileContent.Options.FileNames, ",") + "\n")
		baselineContent.WriteString("Errors::\n")
		diagnosticwriter.FormatDiagnosticsWithColorAndContext(&baselineContent, parsedConfigFileContent.Errors, &diagnosticwriter.FormattingOptions{
			NewLine: "\r\n",
			ComparePathsOptions: tspath.ComparePathsOptions{
				CurrentDirectory:          basePath,
				UseCaseSensitiveFileNames: true,
			},
		})
		baselineContent.WriteString("\n")
		if i != len(input)-1 {
			baselineContent.WriteString("\n")
		}
	}
	if noSubmoduleBaseline {
		baseline.Run(t, baselineFileName, baselineContent.String(), baseline.Options{Subfolder: "config/tsconfigParsing"})
	} else {
		baseline.RunAgainstSubmodule(t, baselineFileName, baselineContent.String(), baseline.Options{Subfolder: "config/tsconfigParsing"})
	}
}

func jsonToReadableText(input any) (string, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(input); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func printFS(output io.Writer, files vfs.FS, root string) error {
	return files.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type().IsRegular() {
			if content, ok := files.ReadFile(path); !ok {
				return fmt.Errorf("failed to read file %s", path)
			} else {
				output.Write([]byte(fmt.Sprintf("//// [%s]\r\n%s\r\n\r\n", path, content)))
			}
		}
		return nil
	})
}

func TestParseSrcCompiler(t *testing.T) {
	t.Parallel()

	repo.SkipIfNoTypeScriptSubmodule(t)

	compilerDir := tspath.NormalizeSlashes(filepath.Join(repo.TypeScriptSubmodulePath, "src", "compiler"))
	tsconfigPath := tspath.CombinePaths(compilerDir, "tsconfig.json")

	fs := vfs.FromOS()
	host := &vfsParseConfigHost{
		fs:               fs,
		currentDirectory: compilerDir,
	}

	jsonText, ok := fs.ReadFile(tsconfigPath)
	assert.Assert(t, ok)
	parsed := parser.ParseJSONText(tsconfigPath, jsonText)

	if len(parsed.Diagnostics()) > 0 {
		for _, error := range parsed.Diagnostics() {
			t.Log(error.Message())
		}
		t.FailNow()
	}

	tsConfigSourceFile := &tsConfigSourceFile{
		sourceFile: parsed,
	}

	parseConfigFileContent := ParseJsonSourceFileConfigFileContent(
		tsConfigSourceFile,
		host,
		host.GetCurrentDirectory(),
		nil,
		tsconfigPath,
		/*resolutionStack*/ nil,
		/*extraFileExtensions*/ nil,
		/*extendedConfigCache*/ nil,
	)

	if len(parseConfigFileContent.Errors) > 0 {
		for _, error := range parseConfigFileContent.Errors {
			t.Log(error.Message())
		}
		t.FailNow()
	}

	// opts := parseConfigFileContent.CompilerOptions()
	// assert.DeepEqual(t, opts, &core.CompilerOptions{}) // TODO: fill out

	fileNames := parseConfigFileContent.Options.FileNames
	fmt.Println(fileNames)
	// assert.DeepEqual(t, fileNames, []string{}) // TODO: fill out (make paths relative to cwd)
}
