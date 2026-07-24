package tsoptions_test

import (
	"slices"
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/microsoft/typescript-go/internal/tsoptions"
	"github.com/microsoft/typescript-go/internal/tsoptions/tsoptionstest"
	"github.com/microsoft/typescript-go/internal/vfs/vfstest"
	"gotest.tools/v3/assert"
)

func TestParsedCommandLine(t *testing.T) {
	t.Parallel()
	t.Run("PossiblyMatchesFileName", func(t *testing.T) {
		t.Parallel()

		noFiles := map[string]string{}
		noFilesFS := vfstest.FromMap(noFiles, true)

		files := map[string]string{
			"/dev/a.ts":         "",
			"/dev/a.d.ts":       "",
			"/dev/a.js":         "",
			"/dev/b.ts":         "",
			"/dev/b.js":         "",
			"/dev/c.d.ts":       "",
			"/dev/z/a.ts":       "",
			"/dev/z/abz.ts":     "",
			"/dev/z/aba.ts":     "",
			"/dev/z/b.ts":       "",
			"/dev/z/bbz.ts":     "",
			"/dev/z/bba.ts":     "",
			"/dev/x/a.ts":       "",
			"/dev/x/aa.ts":      "",
			"/dev/x/b.ts":       "",
			"/dev/x/y/a.ts":     "",
			"/dev/x/y/b.ts":     "",
			"/dev/js/a.js":      "",
			"/dev/js/b.js":      "",
			"/dev/js/d.min.js":  "",
			"/dev/js/ab.min.js": "",
			"/ext/ext.ts":       "",
			"/ext/b/a..b.ts":    "",
		}

		assertMatches := func(t *testing.T, parsedCommandLine *tsoptions.ParsedCommandLine, files map[string]string, matches []string) {
			t.Helper()
			for fileName := range files {
				actual := parsedCommandLine.PossiblyMatchesFileName(fileName)
				expected := slices.Contains(matches, fileName)
				assert.Equal(t, actual, expected, "fileName: %s", fileName)
			}
			for _, fileName := range matches {
				if _, ok := files[fileName]; !ok {
					actual := parsedCommandLine.PossiblyMatchesFileName(fileName)
					assert.Equal(t, actual, true, "fileName: %s", fileName)
				}
			}
		}

		t.Run("with literal file list", func(t *testing.T) {
			t.Parallel()
			t.Run("without exclude", func(t *testing.T) {
				t.Parallel()
				parsedCommandLine := tsoptionstest.GetParsedCommandLine(
					t,
					`{
						"files": [
							"a.ts",
							"b.ts"
						]
					}`,
					files,
					"/dev",
					/*useCaseSensitiveFileNames*/ true,
				)

				assertMatches(t, parsedCommandLine, files, []string{
					"/dev/a.ts",
					"/dev/b.ts",
				})
			})

			t.Run("are not removed due to excludes", func(t *testing.T) {
				t.Parallel()
				parsedCommandLine := tsoptionstest.GetParsedCommandLine(
					t,
					`{
						"files": [
							"a.ts",
							"b.ts"
						],
						"exclude": [
							"b.ts"
						]
					}`,
					files,
					"/dev",
					/*useCaseSensitiveFileNames*/ true,
				)

				assertMatches(t, parsedCommandLine, files, []string{
					"/dev/a.ts",
					"/dev/b.ts",
				})

				emptyParsedCommandLine := parsedCommandLine.ReloadFileNamesOfParsedCommandLine(noFilesFS)
				assertMatches(t, emptyParsedCommandLine, noFiles, []string{
					"/dev/a.ts",
					"/dev/b.ts",
				})
			})

			t.Run("duplicates", func(t *testing.T) {
				t.Parallel()
				parsedCommandLine := tsoptionstest.GetParsedCommandLine(
					t,
					`{
						"files": [
							"a.ts",
							"a.ts",
							"b.ts",
						]
					}`,
					files,
					"/dev",
					/*useCaseSensitiveFileNames*/ true,
				)

				assert.DeepEqual(t, parsedCommandLine.LiteralFileNames(), []string{
					"/dev/a.ts",
					"/dev/b.ts",
				})
			})
		})

		t.Run("with literal include list", func(t *testing.T) {
			t.Parallel()
			t.Run("without exclude", func(t *testing.T) {
				t.Parallel()
				parsedCommandLine := tsoptionstest.GetParsedCommandLine(
					t,
					`{
						"include": [
							"a.ts",
							"b.ts"
						]
					}`,
					files,
					"/dev",
					/*useCaseSensitiveFileNames*/ true,
				)

				assertMatches(t, parsedCommandLine, files, []string{
					"/dev/a.ts",
					"/dev/b.ts",
				})

				emptyParsedCommandLine := parsedCommandLine.ReloadFileNamesOfParsedCommandLine(noFilesFS)
				assertMatches(t, emptyParsedCommandLine, noFiles, []string{
					"/dev/a.ts",
					"/dev/b.ts",
				})
			})
		})

		t.Run("PossiblyMatchesFileName with content mapper extensions", func(t *testing.T) {
			t.Parallel()

			host := tsoptionstest.NewVFSParseConfigHost(map[string]string{}, "/dev", true)
			configFileName := "/dev/tsconfig.json"
			jsonText := `{
			"include": ["src"],
			"contentMappers": [ { "package": "mapper", "extensions": [".box"] } ]
		}`
			tsconfigSourceFile := tsoptions.NewTsconfigSourceFileFromFilePath(configFileName, "/dev/tsconfig.json", jsonText)
			parsedCommandLine := tsoptions.ParseJsonSourceFileConfigFileContent(
				tsconfigSourceFile,
				host,
				"/dev",
				&core.CompilerOptions{LoadExternalPlugins: core.TSTrue},
				nil,
				configFileName,
				nil,
				nil,
			)

			// A created content-mapped file under an included directory must be recognized as a
			// possible root file, or the config's root files are never reloaded for it.
			assert.Assert(t, parsedCommandLine.PossiblyMatchesFileName("/dev/src/new.box"))
			assert.Assert(t, parsedCommandLine.PossiblyMatchesFileName("/dev/src/new.ts"))
			assert.Assert(t, !parsedCommandLine.PossiblyMatchesFileName("/dev/src/new.vue"))
			assert.Assert(t, !parsedCommandLine.PossiblyMatchesFileName("/dev/other/new.box"))
		})
	})
}
